package auth

import (
	"time"

	"gorm.io/gorm"
)

// User 用户
type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Username  string         `gorm:"uniqueIndex;size:50;not null" json:"username"`
	Email     string         `gorm:"uniqueIndex;size:255;not null" json:"email"`
	Password  string         `gorm:"size:255;not null" json:"-"`
	Avatar    string         `gorm:"size:500" json:"avatar"`
	Bio       string         `gorm:"type:text" json:"bio"`
	Role      string         `gorm:"size:20;default:user" json:"role"` // user, admin
	IsActive  bool           `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// APIKey API 密钥
type APIKey struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      uint      `gorm:"index;not null" json:"user_id"`
	Key         string    `gorm:"uniqueIndex;size:64;not null" json:"key"`
	Name        string    `gorm:"size:100" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	LastUsedAt  *time.Time `json:"last_used_at"`
	ExpiresAt   *time.Time `json:"expires_at"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
}

// Service 认证服务
type Service struct {
	db *gorm.DB
}

// NewService 创建认证服务
func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

// TODO: 实现完整的认证逻辑
// - Register: 用户注册
// - Login: 用户登录（JWT）
// - RefreshToken: 刷新 Token
// - Logout: 登出
// - CreateAPIKey: 创建 API 密钥
// - RevokeAPIKey: 撤销 API 密钥
// - AuthMiddleware: JWT 认证中间件
// - AdminMiddleware: 管理员权限中间件
