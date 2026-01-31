package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimiter 速率限制器
type RateLimiter struct {
	visitors map[string]*Visitor
	mu       sync.RWMutex
	rate     int // 每秒请求数
	burst    int // 突发请求数
}

// Visitor 访问者
type Visitor struct {
	limiter  chan struct{}
	lastSeen time.Time
}

// NewRateLimiter 创建速率限制器
func NewRateLimiter(rate, burst int) *RateLimiter {
	rl := &RateLimiter{
		visitors: make(map[string]*Visitor),
		rate:     rate,
		burst:    burst,
	}

	// 清理过期访问者
	go rl.cleanupVisitors()

	return rl
}

// Allow 检查是否允许请求
func (rl *RateLimiter) Allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	visitor, exists := rl.visitors[ip]
	if !exists {
		visitor = &Visitor{
			limiter:  make(chan struct{}, rl.burst),
			lastSeen: time.Now(),
		}
		rl.visitors[ip] = visitor
	}

	visitor.lastSeen = time.Now()

	// 尝试向通道发送，如果通道已满则超过限流
	select {
	case visitor.limiter <- struct{}{}:
		// 启动一个goroutine定期释放令牌
		go func() {
			time.Sleep(time.Second / time.Duration(rl.rate))
			<-visitor.limiter
		}()
		return true
	default:
		return false
	}
}

// cleanupVisitors 清理3分钟未活动的访问者
func (rl *RateLimiter) cleanupVisitors() {
	for {
		time.Sleep(time.Minute)
		rl.mu.Lock()
		for ip, visitor := range rl.visitors {
			if time.Since(visitor.lastSeen) > 3*time.Minute {
				delete(rl.visitors, ip)
			}
		}
		rl.mu.Unlock()
	}
}

// 全局限流器实例
var globalLimiter = NewRateLimiter(10, 20) // 每秒10个请求，最多20个突发（默认值）

// InitRateLimiter 初始化或重新初始化限流器
func InitRateLimiter(rate, burst int) {
	globalLimiter = NewRateLimiter(rate, burst)
}

// RateLimitMiddleware API限流中间件
func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		if !globalLimiter.Allow(ip) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "请求过于频繁，请稍后再试",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
