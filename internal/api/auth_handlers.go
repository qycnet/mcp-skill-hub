package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/qycnet/mcp-skill-hub/internal/auth"
)

// Register 用户注册
func Register(s *auth.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Username string `json:"username" binding:"required,min=3,max=50"`
			Email    string `json:"email" binding:"required,email"`
			Password string `json:"password" binding:"required,min=6,max=50"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user, err := s.Register(req.Username, req.Email, req.Password)
		if err != nil {
			if err == auth.ErrUserAlreadyExists {
				c.JSON(http.StatusConflict, gin.H{"error": "用户名或邮箱已被注册"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"user_id":  user.ID,
			"username": user.Username,
			"email":    user.Email,
			"message":  "注册成功",
		})
	}
}

// Login 用户登录
func Login(s *auth.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Username string `json:"username" binding:"required"`
			Password string `json:"password" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		accessToken, refreshToken, err := s.Login(req.Username, req.Password)
		if err != nil {
			if err == auth.ErrUserNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
				return
			}
			if err == auth.ErrInvalidPassword {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "密码错误"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
			"token_type":    "Bearer",
			"expires_in":    86400, // 24 小时
		})
	}
}

// RefreshToken 刷新令牌
func RefreshToken(s *auth.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			RefreshToken string `json:"refresh_token" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		newToken, err := s.RefreshToken(req.RefreshToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"access_token": newToken,
			"token_type":   "Bearer",
			"expires_in":   86400,
		})
	}
}

// Logout 登出（简化实现，生产环境应加入 token 黑名单）
func Logout(s *auth.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "登出成功"})
	}
}

// GetUserProfile 获取用户资料
func GetUserProfile(s *auth.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := getUserID(c)

		user, err := s.GetUserByID(userID)
		if err != nil {
			if err == auth.ErrUserNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"id":         user.ID,
			"username":   user.Username,
			"email":      user.Email,
			"avatar":     user.Avatar,
			"bio":        user.Bio,
			"created_at": user.CreatedAt,
		})
	}
}

// UpdateUserProfile 更新用户资料
func UpdateUserProfile(s *auth.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := getUserID(c)

		var req struct {
			Avatar string `json:"avatar"`
			Bio    string `json:"bio"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		updates := make(map[string]interface{})
		if req.Avatar != "" {
			updates["avatar"] = req.Avatar
		}
		if req.Bio != "" {
			updates["bio"] = req.Bio
		}

		if err := s.UpdateUser(userID, updates); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
	}
}

// CreateAPIKey 创建 API 密钥
func CreateAPIKey(s *auth.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := getUserID(c)

		var req struct {
			Name        string `json:"name" binding:"required,max=100"`
			Description string `json:"description"`
			ExpiresIn   int    `json:"expires_in"` // 天数，0 表示永不过期
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var expiresAt interface{} = nil
		if req.ExpiresIn > 0 {
			// 计算过期时间
		}

		apiKey, err := s.CreateAPIKey(userID, req.Name, req.Description, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"key_id":      apiKey.ID,
			"key":         apiKey.Key,
			"name":        apiKey.Name,
			"description": apiKey.Description,
			"created_at":  apiKey.CreatedAt,
			"warning":     "请妥善保存此密钥，刷新后将无法查看",
		})
	}
}

// RevokeAPIKey 撤销 API 密钥
func RevokeAPIKey(s *auth.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := getUserID(c)

		keyID, err := parseUint(c.Param("keyId"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的密钥 ID"})
			return
		}

		if err := s.RevokeAPIKey(userID, keyID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "API 密钥已撤销"})
	}
}

// ListAPIKeys 列出 API 密钥
func ListAPIKeys(s *auth.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := getUserID(c)

		keys, err := s.ListAPIKeys(userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"keys": keys,
		})
	}
}

// AdminListUsers 管理员列出用户
func AdminListUsers(s *auth.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 简化实现
		c.JSON(http.StatusOK, gin.H{"message": "管理员功能"})
	}
}

// parseUint 解析无符号整数
func parseUint(s string) (uint, error) {
	v, err := strconv.ParseUint(s, 10, 32)
	return uint(v), err
}
