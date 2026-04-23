package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AuditLog 审计日志
type AuditLog struct {
	ID         string    `gorm:"primaryKey;size:36" json:"id"`
	UserID     uint      `gorm:"index" json:"user_id"`
	Username   string    `gorm:"size:50" json:"username"`
	Action     string    `gorm:"size:50;index" json:"action"`
	Resource   string    `gorm:"size:100" json:"resource"`
	ResourceID string    `gorm:"size:50" json:"resource_id"`
	Method     string    `gorm:"size:10" json:"method"`
	Path       string    `gorm:"size:500" json:"path"`
	IPAddress  string    `gorm:"size:45" json:"ip_address"`
	UserAgent  string    `gorm:"type:text" json:"user_agent"`
	StatusCode int       `gorm:"index" json:"status_code"`
	RequestBody  string  `gorm:"type:text" json:"request_body"`
	ResponseBody string  `gorm:"type:text" json:"response_body"`
	Duration   int64     `json:"duration"` // 毫秒
	CreatedAt  time.Time `gorm:"autoCreateTime;index" json:"created_at"`
}

// AuditMiddleware 审计日志中间件
func AuditMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 记录开始时间
		start := time.Now()

		// 读取请求体
		var requestBody string
		if c.Request.Body != nil {
			bodyBytes, _ := io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			requestBody = string(bodyBytes)
		}

		// 创建响应写入器拦截
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		// 继续处理请求
		c.Next()

		// 计算耗时
		duration := time.Since(start).Milliseconds()

		// 获取用户信息
		userID, _ := c.Get("user_id")
		username, _ := c.Get("username")

		// 创建审计日志
		log := &AuditLog{
			ID:           uuid.New().String(),
			UserID:       getUint(userID),
			Username:     getString(username),
			Action:       getAction(c.Request.Method, c.Request.URL.Path),
			Resource:     getResource(c.Request.URL.Path),
			ResourceID:   getResourceID(c.Request.URL.Path),
			Method:       c.Request.Method,
			Path:         c.Request.URL.Path,
			IPAddress:    c.ClientIP(),
			UserAgent:    c.Request.UserAgent(),
			StatusCode:   c.Writer.Status(),
			RequestBody:  truncate(requestBody, 10000),
			ResponseBody: truncate(blw.body.String(), 10000),
			Duration:     duration,
			CreatedAt:    start,
		}

		// 异步写入数据库
		go func() {
			db.WithContext(context.Background()).Create(log)
		}()
	}
}

// bodyLogWriter 响应体拦截器
type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// getAction 从方法和路径推断动作
func getAction(method, path string) string {
	switch method {
	case "GET":
		return "read"
	case "POST":
		return "create"
	case "PUT", "PATCH":
		return "update"
	case "DELETE":
		return "delete"
	default:
		return "unknown"
	}
}

// getResource 从路径提取资源
func getResource(path string) string {
	// /api/v1/skills/123 -> skills
	// /api/v1/auth/login -> auth
	if len(path) < 10 {
		return "unknown"
	}

	// 跳过 /api/v1/
	parts := splitPath(path)
	if len(parts) >= 3 {
		return parts[2]
	}
	return "unknown"
}

// getResourceID 从路径提取资源 ID
func getResourceID(path string) string {
	parts := splitPath(path)
	for i, part := range parts {
		if _, err := parseUint(part); err == nil && i > 2 {
			return part
		}
	}
	return ""
}

// splitPath 分割路径
func splitPath(path string) []string {
	var parts []string
	for _, part := range bytes.Split([]byte(path), []byte("/")) {
		if len(part) > 0 {
			parts = append(parts, string(part))
		}
	}
	return parts
}

// getUint 安全转换 uint
func getUint(v interface{}) uint {
	if u, ok := v.(uint); ok {
		return u
	}
	return 0
}

// getString 安全转换 string
func getString(v interface{}) string {
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

// truncate 截断字符串
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "... (truncated)"
}

// parseUint 尝试解析无符号整数
func parseUint(s string) (uint64, error) {
	var v uint64
	_, err := json.Number(s).Int64()
	if err != nil {
		return 0, err
	}
	return v, nil
}

// GetAuditLogs 获取审计日志（管理员）
func GetAuditLogs(db *gorm.DB, opts AuditLogOptions) ([]AuditLog, int64, error) {
	var logs []AuditLog
	var total int64

	query := db.Model(&AuditLog{})

	if opts.UserID > 0 {
		query = query.Where("user_id = ?", opts.UserID)
	}

	if opts.Action != "" {
		query = query.Where("action = ?", opts.Action)
	}

	if opts.Resource != "" {
		query = query.Where("resource = ?", opts.Resource)
	}

	if !opts.StartAt.IsZero() {
		query = query.Where("created_at >= ?", opts.StartAt)
	}

	if !opts.EndAt.IsZero() {
		query = query.Where("created_at <= ?", opts.EndAt)
	}

	// 获取总数
	query.Count(&total)

	// 排序和分页
	query = query.Order("created_at DESC").
		Offset(opts.Offset).
		Limit(opts.Limit)

	err := query.Find(&logs).Error
	return logs, total, err
}

// AuditLogOptions 审计日志查询选项
type AuditLogOptions struct {
	UserID   uint
	Action   string
	Resource string
	StartAt  time.Time
	EndAt    time.Time
	Offset   int
	Limit    int
}

// CreateAuditLog 直接创建审计日志
func CreateAuditLog(db *gorm.DB, log *AuditLog) error {
	return db.Create(log).Error
}
