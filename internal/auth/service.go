package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrUserNotFound      = errors.New("用户不存在")
	ErrUserAlreadyExists = errors.New("用户已存在")
	ErrInvalidPassword   = errors.New("密码错误")
	ErrInvalidToken      = errors.New("无效令牌")
	ErrTokenExpired      = errors.New("令牌已过期")
)

// Service 认证服务
type Service struct {
	db       *gorm.DB
	jwtSecret []byte
}

// NewService 创建认证服务
func NewService(db *gorm.DB, jwtSecret string) *Service {
	return &Service{
		db:       db,
		jwtSecret: []byte(jwtSecret),
	}
}

// Register 用户注册
func (s *Service) Register(username, email, password string) (*User, error) {
	// 检查用户名是否已存在
	var existing User
	if err := s.db.Where("username = ? OR email = ?", username, email).First(&existing).Error; err == nil {
		return nil, ErrUserAlreadyExists
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// 创建用户
	user := &User{
		Username: username,
		Email:    email,
		Password: string(hashedPassword),
		Role:     "user",
		IsActive: true,
	}

	if err := s.db.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

// Login 用户登录
func (s *Service) Login(username, password string) (string, string, error) {
	// 查找用户
	var user User
	if err := s.db.Where("username = ? OR email = ?", username, username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", "", ErrUserNotFound
		}
		return "", "", err
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", "", ErrInvalidPassword
	}

	// 生成 JWT
	accessToken, err := s.generateJWT(user.ID, user.Username, user.Role, 24*time.Hour)
	if err != nil {
		return "", "", err
	}

	// 生成刷新令牌
	refreshToken, err := s.generateRefreshToken()
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// RefreshToken 刷新令牌
func (s *Service) RefreshToken(refreshToken string) (string, error) {
	// 验证刷新令牌（简化实现，生产环境应存储并验证）
	claims, err := s.validateToken(refreshToken)
	if err != nil {
		return "", err
	}

	// 生成新访问令牌
	return s.generateJWT(claims.UserID, claims.Username, claims.Role, 24*time.Hour)
}

// ValidateToken 验证令牌
func (s *Service) ValidateToken(tokenString string) (*Claims, error) {
	return s.validateToken(tokenString)
}

// CreateAPIKey 创建 API 密钥
func (s *Service) CreateAPIKey(userID uint, name, description string, expiresAt *time.Time) (*APIKey, error) {
	// 生成密钥
	key, err := s.generateAPIKey()
	if err != nil {
		return nil, err
	}

	apiKey := &APIKey{
		UserID:      userID,
		Key:         key,
		Name:        name,
		Description: description,
		ExpiresAt:   expiresAt,
	}

	if err := s.db.Create(apiKey).Error; err != nil {
		return nil, err
	}

	return apiKey, nil
}

// RevokeAPIKey 撤销 API 密钥
func (s *Service) RevokeAPIKey(userID uint, keyID uint) error {
	var apiKey APIKey
	if err := s.db.First(&apiKey, keyID).Error; err != nil {
		return err
	}

	// 验证所有权
	if apiKey.UserID != userID {
		return errors.New("无权撤销此 API 密钥")
	}

	return s.db.Delete(&apiKey).Error
}

// ListAPIKeys 列出用户的 API 密钥
func (s *Service) ListAPIKeys(userID uint) ([]APIKey, error) {
	var keys []APIKey
	err := s.db.Where("user_id = ?", userID).Find(&keys).Error
	return keys, err
}

// GetUserByID 根据 ID 获取用户
func (s *Service) GetUserByID(id uint) (*User, error) {
	var user User
	err := s.db.First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

// UpdateUser 更新用户信息
func (s *Service) UpdateUser(id uint, updates map[string]interface{}) error {
	// 不允许更新敏感字段
	delete(updates, "password")
	delete(updates, "role")
	delete(updates, "is_active")

	return s.db.Model(&User{}).Where("id = ?", id).Updates(updates).Error
}

// ChangePassword 修改密码
func (s *Service) ChangePassword(id uint, oldPassword, newPassword string) error {
	var user User
	if err := s.db.First(&user, id).Error; err != nil {
		return err
	}

	// 验证旧密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)); err != nil {
		return ErrInvalidPassword
	}

	// 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return s.db.Model(&user).Update("password", string(hashedPassword)).Error
}

// Claims JWT 声明
type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// generateJWT 生成 JWT
func (s *Service) generateJWT(userID uint, username, role string, duration time.Duration) (string, error) {
	claims := Claims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "mcp-skill-hub",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

// validateToken 验证令牌
func (s *Service) validateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return s.jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	// 检查过期
	if claims.ExpiresAt.Before(time.Now()) {
		return nil, ErrTokenExpired
	}

	return claims, nil
}

// generateRefreshToken 生成刷新令牌
func (s *Service) generateRefreshToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// generateAPIKey 生成 API 密钥
func (s *Service) generateAPIKey() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return "mcp_" + hex.EncodeToString(bytes), nil
}

// UpdateLastUsed 更新 API 密钥最后使用时间
func (s *Service) UpdateLastUsed(key string) error {
	now := time.Now()
	return s.db.Model(&APIKey{}).Where("key = ?", key).Update("last_used_at", &now).Error
}

// ValidateAPIKey 验证 API 密钥
func (s *Service) ValidateAPIKey(key string) (*APIKey, error) {
	var apiKey APIKey
	err := s.db.Where("key = ?", key).First(&apiKey).Error
	if err != nil {
		return nil, err
	}

	// 检查是否过期
	if apiKey.ExpiresAt != nil && apiKey.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("API 密钥已过期")
	}

	// 更新最后使用时间
	go s.UpdateLastUsed(key)

	return &apiKey, nil
}
