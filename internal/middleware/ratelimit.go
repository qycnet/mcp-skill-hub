package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// RateLimiter 速率限制器
type RateLimiter struct {
	limiters map[string]*rate.Limiter
	mu       sync.RWMutex
	rate     rate.Limit // 每秒允许的请求数
	burst    int        // 突发容量
}

// NewRateLimiter 创建速率限制器
func NewRateLimiter(requestsPerSecond int, burst int) *RateLimiter {
	return &RateLimiter{
		limiters: make(map[string]*rate.Limiter),
		rate:     rate.Limit(requestsPerSecond),
		burst:    burst,
	}
}

// getLimiter 获取或创建限制器
func (rl *RateLimiter) getLimiter(key string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	limiter, exists := rl.limiters[key]
	if !exists {
		limiter = rate.NewLimiter(rl.rate, rl.burst)
		rl.limiters[key] = limiter
	}

	return limiter
}

// RateLimitMiddleware IP 级别速率限制中间件
func RateLimitMiddleware(requestsPerSecond int, burst int) gin.HandlerFunc {
	limiter := NewRateLimiter(requestsPerSecond, burst)

	return func(c *gin.Context) {
		ip := c.ClientIP()
		
		if !limiter.getLimiter(ip).Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "请求过于频繁，请稍后再试",
				"retry_after": "1s",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// UserRateLimitMiddleware 用户级别速率限制中间件
func UserRateLimitMiddleware(requestsPerSecond int, burst int) gin.HandlerFunc {
	limiter := NewRateLimiter(requestsPerSecond, burst)

	return func(c *gin.Context) {
		// 优先使用用户 ID，其次使用 IP
		var key string
		if userID, exists := c.Get("user_id"); exists {
			key = string(rune(userID.(uint)))
		} else {
			key = c.ClientIP()
		}

		if !limiter.getLimiter(key).Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "请求过于频繁，请稍后再试",
				"retry_after": "1s",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// APIRateLimit API 级别速率限制（不同端点不同限制）
func APIRateLimit(limits map[string]RateLimitConfig) gin.HandlerFunc {
	limiters := make(map[string]*RateLimiter)
	for endpoint, config := range limits {
		limiters[endpoint] = NewRateLimiter(config.RequestsPerSecond, config.Burst)
	}

	return func(c *gin.Context) {
		endpoint := c.FullPath()
		if endpoint == "" {
			endpoint = c.Request.URL.Path
		}

		config, exists := limits[endpoint]
		if !exists {
			c.Next()
			return
		}

		ip := c.ClientIP()
		if !limiters[endpoint].getLimiter(ip).Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "请求过于频繁，请稍后再试",
				"endpoint": endpoint,
				"limit": config.RequestsPerSecond,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RateLimitConfig 速率限制配置
type RateLimitConfig struct {
	RequestsPerSecond int
	Burst             int
}

// DefaultRateLimits 默认速率限制配置
var DefaultRateLimits = map[string]RateLimitConfig{
	"/api/v1/auth/login":    {RequestsPerSecond: 1, Burst: 5},    // 登录：1 req/s
	"/api/v1/auth/register": {RequestsPerSecond: 1, Burst: 3},    // 注册：1 req/s
	"/api/v1/skills":        {RequestsPerSecond: 10, Burst: 20},  // 技能列表：10 req/s
	"/api/v1/search":        {RequestsPerSecond: 5, Burst: 10},   // 搜索：5 req/s
	"/api/v1/admin":         {RequestsPerSecond: 20, Burst: 40},  // 管理接口：20 req/s
}

// SlidingWindowRateLimiter 滑动窗口速率限制器
type SlidingWindowRateLimiter struct {
	windows map[string]*SlidingWindow
	mu      sync.RWMutex
	limit   int           // 窗口内最大请求数
	window  time.Duration // 窗口大小
}

// SlidingWindow 滑动窗口
type SlidingWindow struct {
	timestamps []time.Time
	mu         sync.Mutex
}

// NewSlidingWindowRateLimiter 创建滑动窗口限制器
func NewSlidingWindowRateLimiter(limit int, window time.Duration) *SlidingWindowRateLimiter {
	return &SlidingWindowRateLimiter{
		windows: make(map[string]*SlidingWindow),
		limit:   limit,
		window:  window,
	}
}

// Allow 检查是否允许请求
func (sw *SlidingWindowRateLimiter) Allow(key string) bool {
	sw.mu.Lock()
	window, exists := sw.windows[key]
	if !exists {
		window = &SlidingWindow{timestamps: []time.Time{}}
		sw.windows[key] = window
	}
	sw.mu.Unlock()

	window.mu.Lock()
	defer window.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-sw.window)

	// 移除过期的请求记录
	validIdx := 0
	for i, ts := range window.timestamps {
		if ts.After(cutoff) {
			validIdx = i
			break
		}
	}
	window.timestamps = window.timestamps[validIdx:]

	// 检查是否超过限制
	if len(window.timestamps) >= sw.limit {
		return false
	}

	window.timestamps = append(window.timestamps, now)
	return true
}

// SlidingWindowMiddleware 滑动窗口速率限制中间件
func SlidingWindowMiddleware(limit int, window time.Duration) gin.HandlerFunc {
	limiter := NewSlidingWindowRateLimiter(limit, window)

	return func(c *gin.Context) {
		ip := c.ClientIP()

		if !limiter.Allow(ip) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "请求过于频繁，请稍后再试",
				"limit": limit,
				"window": window.String(),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// CleanupLimiters 定期清理不活跃的限制器
func (rl *RateLimiter) CleanupLimiters(interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			rl.mu.Lock()
			// 简单实现：清空所有限制器
			// 生产环境应根据最后访问时间清理
			rl.limiters = make(map[string]*rate.Limiter)
			rl.mu.Unlock()
		}
	}()
}
