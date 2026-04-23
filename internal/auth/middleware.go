package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware JWT 认证中间件
func AuthMiddleware(s *Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取 Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "缺少认证令牌"})
			c.Abort()
			return
		}

		// 解析 Bearer Token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的认证格式"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// 验证令牌
		claims, err := s.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		// 将用户信息存入上下文
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("user_role", claims.Role)

		c.Next()
	}
}

// AdminMiddleware 管理员权限中间件
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("user_role")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "未授权"})
			c.Abort()
			return
		}

		if role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "需要管理员权限"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// APIKeyMiddleware API 密钥认证中间件
func APIKeyMiddleware(s *Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 支持 X-API-Key header 或 ?api_key=query 参数
		apiKey := c.GetHeader("X-API-Key")
		if apiKey == "" {
			apiKey = c.Query("api_key")
		}

		if apiKey == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "缺少 API 密钥"})
			c.Abort()
			return
		}

		// 验证 API 密钥
		key, err := s.ValidateAPIKey(apiKey)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的 API 密钥"})
			c.Abort()
			return
		}

		// 将用户信息存入上下文
		c.Set("user_id", key.UserID)
		c.Set("api_key", key.Key)

		c.Next()
	}
}

// OptionalAuth 可选认证中间件（有令牌则验证，无令牌则跳过）
func OptionalAuth(s *Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 && parts[0] == "Bearer" {
			claims, err := s.ValidateToken(parts[1])
			if err == nil {
				c.Set("user_id", claims.UserID)
				c.Set("username", claims.Username)
				c.Set("user_role", claims.Role)
			}
		}

		c.Next()
	}
}
