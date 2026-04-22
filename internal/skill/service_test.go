package skill

import (
	"context"
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
	err = db.AutoMigrate(&Skill{}, &SkillVersion{}, &Tag{}, &SkillRating{})
	if err != nil {
		t.Fatalf("数据库迁移失败：%v", err)
	}

	return db
}

// TestCreateSkill 测试创建技能
func TestCreateSkill(t *testing.T) {
	db := setupTestDB(t)
	service := NewService(db, nil)

	ctx := context.Background()

	// 测试正常创建
	skill := &Skill{
		Name:        "test-skill",
		DisplayName: "Test Skill",
		Description: "这是一个测试技能",
		Category:    "developer-tools",
		AuthorID:    1,
	}

	err := service.CreateSkill(ctx, skill)
	assert.NoError(t, err)
	assert.NotEmpty(t, skill.ID)
	assert.NotEmpty(t, skill.UUID)
	assert.Equal(t, "pending", skill.Status)

	// 测试重复创建
	duplicateSkill := &Skill{
		Name:        "test-skill",
		DisplayName: "Duplicate Skill",
		Description: "重复的技能名称",
		Category:    "developer-tools",
		AuthorID:    1,
	}

	err = service.CreateSkill(ctx, duplicateSkill)
	assert.Error(t, err)
	assert.Equal(t, ErrSkillAlreadyExists, err)
}

// TestGetSkillByID 测试获取技能
func TestGetSkillByID(t *testing.T) {
	db := setupTestDB(t)
	service := NewService(db, nil)

	ctx := context.Background()

	// 创建测试技能
	skill := &Skill{
		Name:        "get-test",
		DisplayName: "Get Test Skill",
		Description: "用于测试获取",
		Category:    "test",
		AuthorID:    1,
		Status:      "approved",
	}
	db.Create(skill)

	// 测试正常获取
	retrieved, err := service.GetSkillByID(ctx, skill.ID)
	assert.NoError(t, err)
	assert.Equal(t, skill.ID, retrieved.ID)
	assert.Equal(t, "get-test", retrieved.Name)

	// 测试不存在的技能
	_, err = service.GetSkillByID(ctx, 999)
	assert.Error(t, err)
	assert.Equal(t, ErrSkillNotFound, err)
}

// TestListSkills 测试技能列表
func TestListSkills(t *testing.T) {
	db := setupTestDB(t)
	service := NewService(db, nil)

	ctx := context.Background()

	// 创建多个测试技能
	skills := []Skill{
		{Name: "skill-1", DisplayName: "Skill 1", Category: "dev-tools", AuthorID: 1, Status: "approved", Downloads: 100},
		{Name: "skill-2", DisplayName: "Skill 2", Category: "dev-tools", AuthorID: 1, Status: "approved", Downloads: 200},
		{Name: "skill-3", DisplayName: "Skill 3", Category: "ai", AuthorID: 1, Status: "approved", Downloads: 150},
		{Name: "skill-4", DisplayName: "Skill 4", Category: "ai", AuthorID: 1, Status: "pending", Downloads: 50}, // 未审核
	}

	for i := range skills {
		db.Create(&skills[i])
	}

	// 测试获取全部（只显示 approved）
	list, err := service.ListSkills(ctx, ListOptions{Page: 1, PageSize: 10})
	assert.NoError(t, err)
	assert.Equal(t, int64(3), list.Total) // 排除 pending 的

	// 测试分类过滤
	list, err = service.ListSkills(ctx, ListOptions{Page: 1, PageSize: 10, Category: "dev-tools"})
	assert.NoError(t, err)
	assert.Equal(t, int64(2), list.Total)

	// 测试搜索
	list, err = service.ListSkills(ctx, ListOptions{Page: 1, PageSize: 10, Query: "Skill 1"})
	assert.NoError(t, err)
	assert.Equal(t, int64(1), list.Total)

	// 测试排序（按下载量）
	list, err = service.ListSkills(ctx, ListOptions{Page: 1, PageSize: 10, SortBy: "downloads"})
	assert.NoError(t, err)
	assert.Len(t, list.Skills, 3)
	assert.Equal(t, "skill-2", list.Skills[0].Name) // 下载量最高
}

// TestUpdateSkill 测试更新技能
func TestUpdateSkill(t *testing.T) {
	db := setupTestDB(t)
	service := NewService(db, nil)

	ctx := context.Background()

	// 创建测试技能
	skill := &Skill{
		Name:        "update-test",
		DisplayName: "Original Name",
		Description: "Original Description",
		Category:    "test",
		AuthorID:    1,
	}
	db.Create(skill)

	// 测试正常更新
	updates := map[string]interface{}{
		"display_name":  "Updated Name",
		"description":   "Updated Description",
		"license":       "MIT",
	}

	err := service.UpdateSkill(ctx, skill.ID, updates, 1)
	assert.NoError(t, err)

	// 验证更新
	updated, _ := service.GetSkillByID(ctx, skill.ID)
	assert.Equal(t, "Updated Name", updated.DisplayName)
	assert.Equal(t, "Updated Description", updated.Description)
	assert.Equal(t, "MIT", updated.License)

	// 测试无权更新
	err = service.UpdateSkill(ctx, skill.ID, map[string]interface{}{"display_name": "Hacked"}, 2)
	assert.Error(t, err)
	assert.Equal(t, ErrUnauthorized, err)
}

// TestDeleteSkill 测试删除技能
func TestDeleteSkill(t *testing.T) {
	db := setupTestDB(t)
	service := NewService(db, nil)

	ctx := context.Background()

	// 创建测试技能
	skill := &Skill{
		Name:        "delete-test",
		DisplayName: "To Delete",
		Category:    "test",
		AuthorID:    1,
	}
	db.Create(skill)

	// 测试正常删除
	err := service.DeleteSkill(ctx, skill.ID, 1)
	assert.NoError(t, err)

	// 验证软删除（通过 Unscoped 才能查到）
	var count int64
	db.Model(&Skill{}).Count(&count)
	assert.Equal(t, int64(0), count)

	// 测试无权删除
	skill2 := &Skill{
		Name:        "delete-test-2",
		DisplayName: "To Delete 2",
		Category:    "test",
		AuthorID:    2,
	}
	db.Create(skill2)

	err = service.DeleteSkill(ctx, skill2.ID, 1)
	assert.Error(t, err)
	assert.Equal(t, ErrUnauthorized, err)
}

// TestRateSkill 测试评分
func TestRateSkill(t *testing.T) {
	db := setupTestDB(t)
	service := NewService(db, nil)

	ctx := context.Background()

	// 创建测试技能
	skill := &Skill{
		Name:        "rate-test",
		DisplayName: "Rate Test",
		Category:    "test",
		AuthorID:    1,
		Status:      "approved",
	}
	db.Create(skill)

	// 测试正常评分
	err := service.RateSkill(ctx, skill.ID, 1, 5, "很棒的技能！")
	assert.NoError(t, err)

	// 验证评分
	updated, _ := service.GetSkillByID(ctx, skill.ID)
	assert.Equal(t, float64(5), updated.Rating)
	assert.Equal(t, int64(1), updated.RatingCount)

	// 测试无效评分
	err = service.RateSkill(ctx, skill.ID, 2, 6, "无效评分")
	assert.Error(t, err)

	err = service.RateSkill(ctx, skill.ID, 2, 0, "无效评分")
	assert.Error(t, err)

	// 测试更新评分
	err = service.RateSkill(ctx, skill.ID, 1, 4, "更新评分")
	assert.NoError(t, err)

	// 验证平均分
	updated, _ = service.GetSkillByID(ctx, skill.ID)
	assert.Equal(t, float64(4), updated.Rating) // 只有一个用户的评分
}

// TestPublishVersion 测试发布版本
func TestPublishVersion(t *testing.T) {
	db := setupTestDB(t)
	service := NewService(db, nil)

	ctx := context.Background()

	// 创建测试技能
	skill := &Skill{
		Name:        "version-test",
		DisplayName: "Version Test",
		Category:    "test",
		AuthorID:    1,
		Status:      "approved",
	}
	db.Create(skill)

	// 发布 v1.0.0
	v1 := &SkillVersion{
		Version:     "1.0.0",
		Description: "初始版本",
		MCPVersion:  "1.0",
	}
	err := service.PublishVersion(ctx, skill.ID, v1)
	assert.NoError(t, err)
	assert.True(t, v1.IsLatest)

	// 发布 v1.1.0
	v2 := &SkillVersion{
		Version:     "1.1.0",
		Description: "新增功能",
		MCPVersion:  "1.0",
	}
	err = service.PublishVersion(ctx, skill.ID, v2)
	assert.NoError(t, err)
	assert.True(t, v2.IsLatest)

	// 验证 v1 不再是 latest
	var ver1 SkillVersion
	db.Where("skill_id = ? AND version = ?", skill.ID, "1.0.0").First(&ver1)
	assert.False(t, ver1.IsLatest)

	// 测试重复版本
	v3 := &SkillVersion{
		Version:     "1.0.0",
		Description: "重复版本",
	}
	err = service.PublishVersion(ctx, skill.ID, v3)
	assert.Error(t, err)
	assert.Equal(t, ErrSkillAlreadyExists, err)
}

// TestSearchSkills 测试搜索
func TestSearchSkills(t *testing.T) {
	db := setupTestDB(t)
	service := NewService(db, nil)

	ctx := context.Background()

	// 创建测试数据
	skills := []Skill{
		{Name: "ai-code-assistant", DisplayName: "AI Code Assistant", Description: "AI 编程助手", Category: "ai", AuthorID: 1, Status: "approved"},
		{Name: "python-linter", DisplayName: "Python Linter", Description: "Python 代码检查", Category: "dev-tools", AuthorID: 1, Status: "approved"},
		{Name: "ai-translator", DisplayName: "AI Translator", Description: "AI 翻译工具", Category: "ai", AuthorID: 1, Status: "approved"},
	}

	for i := range skills {
		db.Create(&skills[i])
	}

	// 测试关键词搜索
	list, err := service.SearchSkills(ctx, "AI", ListOptions{Page: 1, PageSize: 10})
	assert.NoError(t, err)
	assert.Equal(t, int64(2), list.Total)

	// 测试分类 + 搜索
	list, err = service.SearchSkills(ctx, "AI", ListOptions{Page: 1, PageSize: 10, Category: "ai"})
	assert.NoError(t, err)
	assert.Equal(t, int64(2), list.Total)
}

// TestListCategories 测试分类列表
func TestListCategories(t *testing.T) {
	db := setupTestDB(t)
	service := NewService(db, nil)

	ctx := context.Background()

	// 创建测试数据
	skills := []Skill{
		{Name: "s1", Category: "ai", AuthorID: 1, Status: "approved"},
		{Name: "s2", Category: "ai", AuthorID: 1, Status: "approved"},
		{Name: "s3", Category: "dev-tools", AuthorID: 1, Status: "approved"},
		{Name: "s4", Category: "dev-tools", AuthorID: 1, Status: "approved"},
		{Name: "s5", Category: "dev-tools", AuthorID: 1, Status: "approved"},
		{Name: "s6", Category: "security", AuthorID: 1, Status: "pending"}, // 不计数
	}

	for i := range skills {
		db.Create(&skills[i])
	}

	// 测试分类统计
	categories, err := service.ListCategories(ctx)
	assert.NoError(t, err)
	assert.Len(t, categories, 2) // 只统计 approved 的

	// 验证排序（按数量降序）
	assert.Equal(t, "dev-tools", categories[0].Category)
	assert.Equal(t, int64(3), categories[0].Count)
	assert.Equal(t, "ai", categories[1].Category)
	assert.Equal(t, int64(2), categories[1].Count)
}

// TestCalculateQualityScore 测试质量评分计算
func TestCalculateQualityScore(t *testing.T) {
	db := setupTestDB(t)
	service := NewService(db, nil)

	// 测试高质量技能
	highQuality := &Skill{
		SecurityScanStatus: "passed",
		Downloads:          5000,
		RatingCount:        100,
		IsVerified:         true,
		Status:             "approved",
	}
	score := service.CalculateQualityScore(highQuality)
	assert.Greater(t, score, 80.0)

	// 测试低质量技能
	lowQuality := &Skill{
		SecurityScanStatus: "pending",
		Downloads:          10,
		RatingCount:        0,
		IsVerified:         false,
		Status:             "pending",
	}
	score = service.CalculateQualityScore(lowQuality)
	assert.Less(t, score, 30.0)
}

// TestCompareVersions 测试版本比较
func TestCompareVersions(t *testing.T) {
	tests := []struct {
		v1       string
		v2       string
		expected int
	}{
		{"1.0.0", "1.0.0", 0},
		{"1.0.1", "1.0.0", 1},
		{"1.1.0", "1.0.0", 1},
		{"2.0.0", "1.9.9", 1},
		{"1.0.0", "1.0.1", -1},
		{"1.0.0", "1.1.0", -1},
	}

	for _, tt := range tests {
		result := compareVersions(tt.v1, tt.v2)
		assert.Equal(t, tt.expected, result, "%s vs %s", tt.v1, tt.v2)
	}
}

// TestIncrementDownload 测试下载计数
func TestIncrementDownload(t *testing.T) {
	db := setupTestDB(t)
	service := NewService(db, nil)

	ctx := context.Background()

	// 创建测试技能
	skill := &Skill{
		Name:        "download-test",
		DisplayName: "Download Test",
		Category:    "test",
		AuthorID:    1,
		Status:      "approved",
	}
	db.Create(skill)

	// 模拟多次下载
	for i := 0; i < 5; i++ {
		err := service.IncrementDownload(ctx, skill.ID, 0, uint(i+1), "127.0.0.1", "test-agent")
		assert.NoError(t, err)
	}

	// 验证下载计数
	updated, _ := service.GetSkillByID(ctx, skill.ID)
	assert.Equal(t, int64(5), updated.Downloads)

	// 验证下载记录
	var count int64
	db.Model(&SkillDownload{}).Where("skill_id = ?", skill.ID).Count(&count)
	assert.Equal(t, int64(5), count)
}
