package skill

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/qycnet/mcp-skill-hub/internal/storage"
	"gorm.io/gorm"
)

var (
	ErrSkillNotFound    = errors.New("技能不存在")
	ErrSkillAlreadyExists = errors.New("技能已存在")
	ErrUnauthorized       = errors.New("未授权操作")
	ErrInvalidVersion     = errors.New("无效的版本号")
)

// Service 技能服务
type Service struct {
	db      *gorm.DB
	storage storage.ObjectStorage
}

// NewService 创建技能服务
func NewService(db *gorm.DB, storage storage.ObjectStorage) *Service {
	return &Service{
		db:      db,
		storage: storage,
	}
}

// CreateSkill 创建技能
func (s *Service) CreateSkill(ctx context.Context, skill *Skill) error {
	// 检查名称是否已存在
	var existing Skill
	if err := s.db.Where("name = ?", skill.Name).First(&existing).Error; err == nil {
		return ErrSkillAlreadyExists
	}

	// 生成 UUID
	skill.UUID = uuid.New().String()
	skill.Status = "pending" // 待审核

	// 创建技能
	return s.db.WithContext(ctx).Create(skill).Error
}

// GetSkillByID 根据 ID 获取技能
func (s *Service) GetSkillByID(ctx context.Context, id uint) (*Skill, error) {
	var skill Skill
	err := s.db.WithContext(ctx).
		Preload("Author").
		Preload("Tags").
		Preload("Versions").
		Preload("SecurityIssues").
		First(&skill, id).Error
	
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrSkillNotFound
		}
		return nil, err
	}

	return &skill, nil
}

// GetSkillByUUID 根据 UUID 获取技能
func (s *Service) GetSkillByUUID(ctx context.Context, uuid string) (*Skill, error) {
	var skill Skill
	err := s.db.WithContext(ctx).
		Preload("Author").
		Preload("Tags").
		Preload("Versions").
		Where("uuid = ?", uuid).
		First(&skill).Error
	
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrSkillNotFound
		}
		return nil, err
	}

	return &skill, nil
}

// ListSkills 获取技能列表
func (s *Service) ListSkills(ctx context.Context, opts ListOptions) (*SkillList, error) {
	var skills []Skill
	var total int64

	query := s.db.WithContext(ctx).
		Preload("Author").
		Preload("Tags").
		Where("status = ?", "approved").
		Where("is_private = ?", false)

	// 分类过滤
	if opts.Category != "" {
		query = query.Where("category = ?", opts.Category)
	}

	// 标签过滤
	if len(opts.Tags) > 0 {
		query = query.Joins("JOIN skill_tags ON skill_tags.skill_id = skills.id").
			Joins("JOIN tags ON tags.id = skill_tags.tag_id").
			Where("tags.name IN ?", opts.Tags)
	}

	// 搜索
	if opts.Query != "" {
		searchPattern := "%" + opts.Query + "%"
		query = query.Where("name LIKE ? OR display_name LIKE ? OR description LIKE ?",
			searchPattern, searchPattern, searchPattern)
	}

	// 获取总数
	if err := query.Model(&Skill{}).Count(&total).Error; err != nil {
		return nil, err
	}

	// 排序
	switch opts.SortBy {
	case "downloads":
		query = query.Order("downloads DESC")
	case "rating":
		query = query.Order("rating DESC")
	case "created_at":
		query = query.Order("created_at DESC")
	case "updated_at":
		query = query.Order("updated_at DESC")
	case "quality_score":
		query = query.Order("quality_score DESC")
	default:
		query = query.Order("downloads DESC")
	}

	// 分页
	offset := (opts.Page - 1) * opts.PageSize
	query = query.Offset(offset).Limit(opts.PageSize)

	if err := query.Find(&skills).Error; err != nil {
		return nil, err
	}

	return &SkillList{
		Skills: skills,
		Total:  total,
		Page:   opts.Page,
		PageSize: opts.PageSize,
	}, nil
}

// UpdateSkill 更新技能
func (s *Service) UpdateSkill(ctx context.Context, id uint, updates map[string]interface{}, userID uint) error {
	// 检查权限
	var skill Skill
	if err := s.db.First(&skill, id).Error; err != nil {
		return err
	}

	if skill.AuthorID != userID {
		return ErrUnauthorized
	}

	// 不允许更新某些字段
	delete(updates, "uuid")
	delete(updates, "author_id")
	delete(updates, "status")
	delete(updates, "is_verified")

	return s.db.WithContext(ctx).Model(&skill).Updates(updates).Error
}

// DeleteSkill 删除技能
func (s *Service) DeleteSkill(ctx context.Context, id uint, userID uint) error {
	var skill Skill
	if err := s.db.First(&skill, id).Error; err != nil {
		return err
	}

	if skill.AuthorID != userID {
		return ErrUnauthorized
	}

	// 软删除
	return s.db.WithContext(ctx).Delete(&skill).Error
}

// PublishVersion 发布新版本
func (s *Service) PublishVersion(ctx context.Context, skillID uint, version *SkillVersion) error {
	// 检查技能是否存在
	var skill Skill
	if err := s.db.First(&skill, skillID).Error; err != nil {
		return err
	}

	// 检查版本是否已存在
	var existing SkillVersion
	if err := s.db.Where("skill_id = ? AND version = ?", skillID, version.Version).First(&existing).Error; err == nil {
		return ErrSkillAlreadyExists
	}

	// 设置版本关联
	version.SkillID = skillID

	// 如果是第一个版本或版本号更大，设为 latest
	var latestVersion SkillVersion
	if err := s.db.Where("skill_id = ? AND is_latest = ?", skillID, true).First(&latestVersion).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			version.IsLatest = true
		}
	} else {
		if compareVersions(version.Version, latestVersion.Version) > 0 {
			// 更新旧版本为 non-latest
			s.db.Model(&latestVersion).Update("is_latest", false)
			version.IsLatest = true
			
			// 更新技能的最新版本号
			s.db.Model(&skill).Update("latest_version", version.Version)
		}
	}

	// 创建版本
	return s.db.WithContext(ctx).Create(version).Error
}

// RateSkill 评分
func (s *Service) RateSkill(ctx context.Context, skillID uint, userID uint, rating int, comment string) error {
	// 检查评分范围
	if rating < 1 || rating > 5 {
		return errors.New("评分必须在 1-5 之间")
	}

	// 创建或更新评分
	var skillRating SkillRating
	result := s.db.Where("skill_id = ? AND user_id = ?", skillID, userID).FirstOrCreate(&skillRating, SkillRating{
		SkillID: skillID,
		UserID:  userID,
		Rating:  rating,
		Comment: comment,
	})

	if result.RowsAffected == 0 {
		// 更新已有评分
		skillRating.Rating = rating
		skillRating.Comment = comment
		return s.db.Save(&skillRating).Error
	}

	// 重新计算平均评分
	return s.recalculateRating(ctx, skillID)
}

// RecalculateRating 重新计算平均评分
func (s *Service) recalculateRating(ctx context.Context, skillID uint) error {
	var avgRating float64
	var count int64

	err := s.db.Model(&SkillRating{}).
		Where("skill_id = ?", skillID).
		Select("AVG(rating), COUNT(*)").
		Scan(&avgRating, &count).Error
	
	if err != nil {
		return err
	}

	return s.db.Model(&Skill{}).Where("id = ?", skillID).Updates(map[string]interface{}{
		"rating":       avgRating,
		"rating_count": count,
	}).Error
}

// SearchSkills 搜索技能
func (s *Service) SearchSkills(ctx context.Context, query string, opts ListOptions) (*SkillList, error) {
	opts.Query = query
	return s.ListSkills(ctx, opts)
}

// ListCategories 获取分类列表
func (s *Service) ListCategories(ctx context.Context) ([]CategoryStat, error) {
	var categories []CategoryStat

	err := s.db.Model(&Skill{}).
		Select("category, COUNT(*) as count").
		Where("status = ?", "approved").
		Where("is_private = ?", false).
		Group("category").
		Order("count DESC").
		Scan(&categories).Error

	return categories, err
}

// ListOptions 列表选项
type ListOptions struct {
	Page     int
	PageSize int
	SortBy   string
	Category string
	Tags     []string
	Query    string
}

// SkillList 技能列表
type SkillList struct {
	Skills   []Skill `json:"skills"`
	Total    int64   `json:"total"`
	Page     int     `json:"page"`
	PageSize int     `json:"page_size"`
}

// CategoryStat 分类统计
type CategoryStat struct {
	Category string `json:"category"`
	Count    int64  `json:"count"`
}

// compareVersions 比较语义化版本号
// 返回值：1 (v1>v2), 0 (v1==v2), -1 (v1<v2)
func compareVersions(v1, v2 string) int {
	// 简化实现，生产环境应使用 github.com/Masterminds/semver
	if v1 == v2 {
		return 0
	}
	if v1 > v2 {
		return 1
	}
	return -1
}

// IncrementDownload 增加下载计数
func (s *Service) IncrementDownload(ctx context.Context, skillID uint, versionID uint, userID uint, ip string, ua string) error {
	// 记录下载
	download := SkillDownload{
		SkillID:   skillID,
		VersionID: versionID,
		UserID:    userID,
		IPAddress: ip,
		UserAgent: ua,
	}
	if err := s.db.Create(&download).Error; err != nil {
		return err
	}

	// 增加计数
	return s.db.Model(&Skill{}).Where("id = ?", skillID).UpdateColumn("downloads", gorm.Expr("downloads + ?", 1)).Error
}

// CalculateQualityScore 计算质量评分
func (s *Service) CalculateQualityScore(skill *Skill) float64 {
	score := 0.0

	// 代码质量 (25 分) - 简化实现
	score += 20.0

	// 安全性 (25 分)
	if skill.SecurityScanStatus == "passed" {
		score += 25.0
	} else if skill.SecurityScanStatus == "pending" {
		score += 10.0
	}

	// 社区活跃 (20 分)
	if skill.Downloads > 1000 {
		score += 10.0
	} else if skill.Downloads > 100 {
		score += 5.0
	}
	if skill.RatingCount > 50 {
		score += 10.0
	} else if skill.RatingCount > 10 {
		score += 5.0
	}

	// 兼容性 (15 分)
	if skill.IsVerified {
		score += 15.0
	}

	// 维护性 (15 分)
	if skill.Status == "approved" {
		score += 15.0
	}

	return score
}
