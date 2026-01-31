package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// RedisRateLimiter Redis 限流器
type RedisRateLimiter struct {
	client *redis.Client
	rate   int           // 每秒请求数
	burst  int           // 突发请求数
	window time.Duration // 时间窗口
}

// NewRedisRateLimiter 创建 Redis 限流器
func NewRedisRateLimiter(client *redis.Client, rate, burst int) *RedisRateLimiter {
	return &RedisRateLimiter{
		client: client,
		rate:   rate,
		burst:  burst,
		window: time.Second, // 1 秒窗口
	}
}

// Allow 检查是否允许请求（滑动窗口算法）
func (rl *RedisRateLimiter) Allow(ctx context.Context, key string) (bool, error) {
	// Redis Lua 脚本：原子性操作
	script := `
		-- 获取当前时间戳（毫秒）
		local now = tonumber(ARGV[1])
		local window = tonumber(ARGV[2])
		local limit = tonumber(ARGV[3])

		-- 删除时间窗口之外的记录
		redis.call('ZREMRANGEBYSCORE', KEYS[1], '-inf', now - window)

		-- 统计窗口内的请求数
		local count = redis.call('ZCARD', KEYS[1])

		-- 检查是否超过限制
		if count < limit then
			-- 添加当前请求到有序集合
			redis.call('ZADD', KEYS[1], now, now)
			-- 设置过期时间为窗口时间 + 1 秒
			redis.call('EXPIRE', KEYS[1], window/1000000000 + 1)
			return 1
		else
			return 0
		end
	`

	now := time.Now().UnixNano()
	result, err := rl.client.Eval(ctx, script, []string{key}, now, rl.window.Milliseconds(), rl.burst).Result()
	if err != nil {
		return false, err
	}

	// 返回 1 表示允许，0 表示拒绝
	return result.(int64) == 1, nil
}

// 全局 Redis 限流器实例
var globalRedisLimiter *RedisRateLimiter

// InitRedisRateLimiter 初始化 Redis 限流器
func InitRedisRateLimiter(client *redis.Client, rate, burst int) {
	globalRedisLimiter = NewRedisRateLimiter(client, rate, burst)
}

// IsRedisEnabled 检查 Redis 限流器是否已启用
func IsRedisEnabled() bool {
	return globalRedisLimiter != nil
}

// RedisRateLimitMiddleware Redis 限流中间件
func RedisRateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 如果 Redis 限流器未初始化，跳过限流
		if globalRedisLimiter == nil {
			c.Next()
			return
		}

		// 获取客户端 IP
		ip := c.ClientIP()
		ctx := context.Background()

		// 生成 Redis key
		key := fmt.Sprintf("rate_limit:ip:%s", ip)

		// 检查是否允许请求
		allowed, err := globalRedisLimiter.Allow(ctx, key)
		if err != nil {
			// Redis 错误时记录日志，但允许请求通过（降级处理）
			// 这样即使 Redis 故障也不会影响服务
			c.Next()
			return
		}

		if !allowed {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "请求过于频繁，请稍后再试",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
