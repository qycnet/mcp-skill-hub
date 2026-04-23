package auth

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupTestDB 创建测试数据库
func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("数据库初始化失败：%v", err)
	}

	// 自动迁移
	err = db.AutoMigrate(&User{}, &APIKey{})
	if err != nil {
		t.Fatalf("数据库迁移失败：%v", err)
	}

	return db
}

// TestRegister 测试用户注册
func TestRegister(t *testing.T) {
	db := setupTestDB(t)
	service := NewService(db, "test-secret-key")

	// 测试正常注册
	user, err := service.Register("testuser", "test@example.com", "password123")
	assert.NoError(t, err)
	assert.NotEmpty(t, user.ID)
	assert.Equal(t, "testuser", user.Username)
	assert.Equal(t, "user", user.Role)
	assert.True(t, user.IsActive)

	// 测试重复用户名
	_, err = service.Register("testuser", "another@example.com", "password456")
	assert.Error(t, err)
	assert.Equal(t, ErrUserAlreadyExists, err)

	// 测试重复邮箱
	_, err = service.Register("anotheruser", "test@example.com", "password456")
	assert.Error(t, err)
	assert.Equal(t, ErrUserAlreadyExists, err)
}

// TestLogin 测试用户登录
func TestLogin(t *testing.T) {
	db := setupTestDB(t)
	service := NewService(db, "test-secret-key")

	// 先注册
	_, err := service.Register("loginuser", "login@example.com", "password123")
	assert.NoError(t, err)

	// 测试正确登录
	accessToken, refreshToken, err := service.Login("loginuser", "password123")
	assert.NoError(t, err)
	assert.NotEmpty(t, accessToken)
	assert.NotEmpty(t, refreshToken)

	// 测试错误密码
	_, _, err = service.Login("loginuser", "wrongpassword")
	assert.Error(t, err)
	assert.Equal(t, ErrInvalidPassword, err)

	// 测试用户不存在
	_, _, err = service.Login("nonexistent", "password123")
	assert.Error(t, err)
	assert.Equal(t, ErrUserNotFound, err)
}

// TestValidateToken 测试令牌验证
func TestValidateToken(t *testing.T) {
	db := setupTestDB(t)
	service := NewService(db, "test-secret-key")

	// 注册并登录
	_, err := service.Register("tokenuser", "token@example.com", "password123")
	assert.NoError(t, err)

	accessToken, _, err := service.Login("tokenuser", "password123")
	assert.NoError(t, err)

	// 测试有效令牌
	claims, err := service.ValidateToken(accessToken)
	assert.NoError(t, err)
	assert.Equal(t, "tokenuser", claims.Username)
	assert.Equal(t, "user", claims.Role)

	// 测试无效令牌
	_, err = service.ValidateToken("invalid-token")
	assert.Error(t, err)

	// 测试过期令牌（使用短时效令牌）
	service2 := NewService(db, "test-secret-key")
	expiredToken, err := service2.generateJWT(1, "test", "user", 1*time.Millisecond)
	assert.NoError(t, err)
	time.Sleep(10 * time.Millisecond) // 等待过期

	_, err = service2.ValidateToken(expiredToken)
	assert.Error(t, err)
	assert.Equal(t, ErrTokenExpired, err)
}

// TestCreateAPIKey 测试创建 API 密钥
func TestCreateAPIKey(t *testing.T) {
	db := setupTestDB(t)
	service := NewService(db, "test-secret-key")

	// 先创建用户
	user, _ := service.Register("apiuser", "api@example.com", "password123")

	// 测试创建 API 密钥
	apiKey, err := service.CreateAPIKey(user.ID, "test-key", "测试密钥", nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, apiKey.ID)
	assert.NotEmpty(t, apiKey.Key)
	assert.Equal(t, "test-key", apiKey.Name)
	assert.Equal(t, "测试密钥", apiKey.Description)
	assert.True(t, apiKey.Key != "")
}

// TestRevokeAPIKey 测试撤销 API 密钥
func TestRevokeAPIKey(t *testing.T) {
	db := setupTestDB(t)
	service := NewService(db, "test-secret-key")

	// 创建用户和 API 密钥
	user, _ := service.Register("revokeuser", "revoke@example.com", "password123")
	apiKey, _ := service.CreateAPIKey(user.ID, "revoke-key", "要撤销的密钥", nil)

	// 测试正常撤销
	err := service.RevokeAPIKey(user.ID, apiKey.ID)
	assert.NoError(t, err)

	// 验证已删除
	var count int64
	db.Model(&APIKey{}).Where("id = ?", apiKey.ID).Count(&count)
	assert.Equal(t, int64(0), count)

	// 测试无权撤销
	apiKey2, _ := service.CreateAPIKey(user.ID, "revoke-key-2", "另一个密钥", nil)
	err = service.RevokeAPIKey(999, apiKey2.ID) // 不同用户 ID
	assert.Error(t, err)
}

// TestListAPIKeys 测试列出 API 密钥
func TestListAPIKeys(t *testing.T) {
	db := setupTestDB(t)
	service := NewService(db, "test-secret-key")

	// 创建用户和多个 API 密钥
	user, _ := service.Register("listuser", "list@example.com", "password123")
	service.CreateAPIKey(user.ID, "key-1", "第一个密钥", nil)
	service.CreateAPIKey(user.ID, "key-2", "第二个密钥", nil)
	service.CreateAPIKey(user.ID, "key-3", "第三个密钥", nil)

	// 测试列出
	keys, err := service.ListAPIKeys(user.ID)
	assert.NoError(t, err)
	assert.Len(t, keys, 3)

	// 测试空列表
	keys2, err := service.ListAPIKeys(999)
	assert.NoError(t, err)
	assert.Len(t, keys2, 0)
}

// TestValidateAPIKey 测试验证 API 密钥
func TestValidateAPIKey(t *testing.T) {
	db := setupTestDB(t)
	service := NewService(db, "test-secret-key")

	// 创建用户和 API 密钥
	user, _ := service.Register("validateuser", "validate@example.com", "password123")
	apiKey, _ := service.CreateAPIKey(user.ID, "validate-key", "验证密钥", nil)

	// 测试有效密钥
	validatedKey, err := service.ValidateAPIKey(apiKey.Key)
	assert.NoError(t, err)
	assert.Equal(t, apiKey.Key, validatedKey.Key)
	assert.Equal(t, user.ID, validatedKey.UserID)

	// 测试无效密钥
	_, err = service.ValidateAPIKey("invalid-key")
	assert.Error(t, err)
}

// TestChangePassword 测试修改密码
func TestChangePassword(t *testing.T) {
	db := setupTestDB(t)
	service := NewService(db, "test-secret-key")

	// 创建用户
	user, _ := service.Register("changeuser", "change@example.com", "oldpassword")

	// 测试正确修改
	err := service.ChangePassword(user.ID, "oldpassword", "newpassword")
	assert.NoError(t, err)

	// 验证新密码可以登录
	_, _, err = service.Login("changeuser", "newpassword")
	assert.NoError(t, err)

	// 验证旧密码不能登录
	_, _, err = service.Login("changeuser", "oldpassword")
	assert.Error(t, err)

	// 测试错误旧密码
	err = service.ChangePassword(user.ID, "wrongpassword", "anotherpassword")
	assert.Error(t, err)
	assert.Equal(t, ErrInvalidPassword, err)
}

// TestUpdateUser 测试更新用户信息
func TestUpdateUser(t *testing.T) {
	db := setupTestDB(t)
	service := NewService(db, "test-secret-key")

	// 创建用户
	user, _ := service.Register("updateuser", "update@example.com", "password123")

	// 测试正常更新
	updates := map[string]interface{}{
		"avatar": "https://example.com/avatar.jpg",
		"bio":    "这是一个测试用户",
	}
	err := service.UpdateUser(user.ID, updates)
	assert.NoError(t, err)

	// 验证更新
	updated, _ := service.GetUserByID(user.ID)
	assert.Equal(t, "https://example.com/avatar.jpg", updated.Avatar)
	assert.Equal(t, "这是一个测试用户", updated.Bio)

	// 测试不允许更新敏感字段
	updates = map[string]interface{}{
		"role":      "admin", // 不允许
		"is_active": false,   // 不允许
	}
	err = service.UpdateUser(user.ID, updates)
	assert.NoError(t, err) // 不报错，但会忽略这些字段

	// 验证敏感字段未变
	updated, _ = service.GetUserByID(user.ID)
	assert.Equal(t, "user", updated.Role)
	assert.True(t, updated.IsActive)
}

// TestGetUserByID 测试获取用户
func TestGetUserByID(t *testing.T) {
	db := setupTestDB(t)
	service := NewService(db, "test-secret-key")

	// 创建用户
	user, _ := service.Register("getuser", "get@example.com", "password123")

	// 测试正常获取
	retrieved, err := service.GetUserByID(user.ID)
	assert.NoError(t, err)
	assert.Equal(t, "getuser", retrieved.Username)

	// 测试用户不存在
	_, err = service.GetUserByID(999)
	assert.Error(t, err)
	assert.Equal(t, ErrUserNotFound, err)
}

// TestGenerateAPIKey 测试 API 密钥生成格式
func TestGenerateAPIKey(t *testing.T) {
	db := setupTestDB(t)
	service := NewService(db, "test-secret-key")

	key, err := service.generateAPIKey()
	assert.NoError(t, err)
	assert.NotEmpty(t, key)
	assert.Len(t, key, 68) // "mcp_" + 64 字符 hex
	assert.HasPrefix(t, key, "mcp_")
}
