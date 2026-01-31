package middleware

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimiterConfig holds rate limiter configuration
type RateLimiterConfig struct {
	RequestsPerMinute int
	RequestsPerHour   int
	RequestsPerDay    int
}

// SensitiveOperationRateLimiter tracks request counts for sensitive operations
type SensitiveOperationRateLimiter struct {
	mu    sync.Mutex
	store map[string]*operationCount
}

// operationCount holds count information for an operation
type operationCount struct {
	minuteCount int
	hourCount   int
	dayCount    int
	lastMinute  time.Time
	lastHour    time.Time
	lastDay     time.Time
}

// NewSensitiveOperationRateLimiter creates a new rate limiter
func NewSensitiveOperationRateLimiter() *SensitiveOperationRateLimiter {
	limiter := &SensitiveOperationRateLimiter{
		store: make(map[string]*operationCount),
	}

	// Start cleanup goroutine
	go limiter.cleanup()

	return limiter
}

// GetConfig returns rate limit config for operation type
func (l *SensitiveOperationRateLimiter) GetConfig(operation string) *RateLimiterConfig {
	switch operation {
	case "login":
		return &RateLimiterConfig{
			RequestsPerMinute: 5,
			RequestsPerHour:   20,
			RequestsPerDay:    50,
		}
	case "password_change":
		return &RateLimiterConfig{
			RequestsPerMinute: 2,
			RequestsPerHour:   5,
			RequestsPerDay:    10,
		}
	case "admin_create":
		return &RateLimiterConfig{
			RequestsPerMinute: 1,
			RequestsPerHour:   10,
			RequestsPerDay:    50,
		}
	default:
		return &RateLimiterConfig{
			RequestsPerMinute: 60,
			RequestsPerHour:   1000,
			RequestsPerDay:    10000,
		}
	}
}

// CheckRateLimit checks if the request should be rate limited
func (l *SensitiveOperationRateLimiter) CheckRateLimit(key string, config *RateLimiterConfig) (bool, string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	count, exists := l.store[key]

	if !exists {
		count = &operationCount{
			lastMinute: now,
			lastHour:   now,
			lastDay:    now,
		}
		l.store[key] = count
	}

	// Reset counters if time has passed
	if now.Sub(count.lastMinute) >= time.Minute {
		count.minuteCount = 0
		count.lastMinute = now
	}
	if now.Sub(count.lastHour) >= time.Hour {
		count.hourCount = 0
		count.lastHour = now
	}
	if now.Sub(count.lastDay) >= 24*time.Hour {
		count.dayCount = 0
		count.lastDay = now
	}

	// Check limits
	if config.RequestsPerMinute > 0 && count.minuteCount >= config.RequestsPerMinute {
		return true, fmt.Sprintf("请求过于频繁，每分钟最多%d次请求", config.RequestsPerMinute)
	}
	if config.RequestsPerHour > 0 && count.hourCount >= config.RequestsPerHour {
		return true, fmt.Sprintf("请求过于频繁，每小时最多%d次请求", config.RequestsPerHour)
	}
	if config.RequestsPerDay > 0 && count.dayCount >= config.RequestsPerDay {
		return true, fmt.Sprintf("请求过于频繁，每天最多%d次请求", config.RequestsPerDay)
	}

	// Increment counters
	count.minuteCount++
	count.hourCount++
	count.dayCount++

	return false, ""
}

// cleanup removes old entries from the store
func (l *SensitiveOperationRateLimiter) cleanup() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		l.mu.Lock()
		now := time.Now()

		for key, count := range l.store {
			// Remove entries not used in 24 hours
			if now.Sub(count.lastDay) > 24*time.Hour {
				delete(l.store, key)
			}
		}

		l.mu.Unlock()
	}
}

// Global rate limiter instance
var sensitiveRateLimiter = NewSensitiveOperationRateLimiter()

// LoginRateLimitMiddleware applies rate limiting to login attempts
func LoginRateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		key := getClientKey(c)
		config := sensitiveRateLimiter.GetConfig("login")

		limited, message := sensitiveRateLimiter.CheckRateLimit("login:"+key, config)
		if limited {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": message,
				"code":  "RATE_LIMIT_EXCEEDED",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// PasswordChangeRateLimitMiddleware applies rate limiting to password changes
func PasswordChangeRateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		key := getClientKey(c)
		config := sensitiveRateLimiter.GetConfig("password_change")

		limited, message := sensitiveRateLimiter.CheckRateLimit("password_change:"+key, config)
		if limited {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": message,
				"code":  "RATE_LIMIT_EXCEEDED",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// AdminCreateRateLimitMiddleware applies rate limiting to admin creation
func AdminCreateRateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		key := getClientKey(c)
		config := sensitiveRateLimiter.GetConfig("admin_create")

		limited, message := sensitiveRateLimiter.CheckRateLimit("admin_create:"+key, config)
		if limited {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": message,
				"code":  "RATE_LIMIT_EXCEEDED",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// getClientKey gets a unique key for the client
func getClientKey(c *gin.Context) string {
	// Try to get user ID first
	if userID, exists := c.Get("user_id"); exists {
		return fmt.Sprintf("user:%v", userID)
	}

	// Fall back to IP address
	return fmt.Sprintf("ip:%s", c.ClientIP())
}
