package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type RateLimiter struct {
	limiters map[string]*rate.Limiter
	mu       sync.RWMutex
	rate     rate.Limit
	burst    int
}

func NewRateLimiter(rps int, burst int) *RateLimiter {
	return &RateLimiter{
		limiters: make(map[string]*rate.Limiter),
		rate:     rate.Limit(rps),
		burst:    burst,
	}
}

func (rl *RateLimiter) GetLimiter(ip string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	limiter, exists := rl.limiters[ip]
	if !exists {
		limiter = rate.NewLimiter(rl.rate, rl.burst)
		rl.limiters[ip] = limiter
	}

	return limiter
}

func (rl *RateLimiter) CleanupOldLimiters() {
	ticker := time.NewTicker(time.Minute)
	go func() {
		for range ticker.C {
			rl.mu.Lock()
			rl.mu.Unlock()
		}
	}()
}

func RateLimitMiddleware() gin.HandlerFunc {
	limiter := NewRateLimiter(10, 20)
	limiter.CleanupOldLimiters()

	return func(c *gin.Context) {
		ip := c.ClientIP()
		l := limiter.GetLimiter(ip)

		if !l.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "rate limit exceeded",
				"message": "too many requests, please try again later",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

func RateLimitWithConfig(rps, burst int) gin.HandlerFunc {
	limiter := NewRateLimiter(rps, burst)
	limiter.CleanupOldLimiters()

	return func(c *gin.Context) {
		ip := c.ClientIP()
		l := limiter.GetLimiter(ip)

		if !l.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "rate limit exceeded",
				"message": "too many requests, please try again later",
				"retry_after": "1s",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

func RateLimitByKey(keyFunc func(*gin.Context) string, rps, burst int) gin.HandlerFunc {
	limiter := NewRateLimiter(rps, burst)
	limiter.CleanupOldLimiters()

	return func(c *gin.Context) {
		key := keyFunc(c)
		if key == "" {
			key = c.ClientIP() // Fallback para IP
		}

		l := limiter.GetLimiter(key)

		if !l.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "rate limit exceeded",
				"message": "too many requests, please try again later",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
