package middleware

import (
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// RateLimiter 速率限制器
type RateLimiter struct {
	visitors map[string]*visitor
	mu       sync.RWMutex
	rate     rate.Limit
	burst    int
}

type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// NewRateLimiter 创建速率限制器
func NewRateLimiter(rate rate.Limit, burst int) *RateLimiter {
	return &RateLimiter{
		visitors: make(map[string]*visitor),
		rate:     rate,
		burst:    burst,
	}
}

// Allow 检查是否允许请求
func (rl *RateLimiter) Allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	v, exists := rl.visitors[ip]
	if !exists {
		limiter := rate.NewLimiter(rl.rate, rl.burst)
		v = &visitor{limiter: &limiter, lastSeen: time.Now()}
		rl.visitors[ip] = v
	} else {
		v.lastSeen = time.Now()
	}

	// 清理过期访问者（可选）
	go rl.cleanup()

	return v.limiter.Allow()
}

// cleanup 清理超过 3 分钟未访问的记录
func (rl *RateLimiter) cleanup() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	for ip, v := range rl.visitors {
		if time.Since(v.lastSeen) > 3*time.Minute {
			delete(rl.visitors, ip)
		}
	}
}

// RateLimitMiddleware 速率限制中间件
func RateLimitMiddleware(rl *RateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		if !rl.Allow(ip) {
			c.JSON(429, gin.H{
				"error": "请求过于频繁，请稍后再试",
				"retry_after": 60, // 秒
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// AllowWait 等待直到允许请求
func (rl *RateLimiter) AllowWait(ip string) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	v, exists := rl.visitors[ip]
	if !exists {
		limiter := rate.NewLimiter(rl.rate, rl.burst)
		v = &visitor{limiter: &limiter, lastSeen: time.Now()}
		rl.visitors[ip] = v
	}

	// 阻塞直到允许
	v.limiter.Wait(c.Context())
}

// GetLimiter 获取指定 IP 的限制器
func (rl *RateLimiter) GetLimiter(ip string) *rate.Limiter {
	rl.mu.RLock()
	defer rl.mu.RUnlock()

	if v, exists := rl.visitors[ip]; exists {
		return v.limiter
	}
	return nil
}
