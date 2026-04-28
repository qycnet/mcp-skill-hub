package skill

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrInvalidSkillPackage = errors.New("无效的技能包格式")
	ErrMissingManifest     = errors.New("缺少 mcp-manifest.json")
	ErrInvalidManifest     = errors.New("manifest.json 格式无效")
	ErrPackageTooLarge     = errors.New("技能包过大（最大 100MB）")
)

const (
	MaxPackageSize = 100 << 20 // 100MB
)

// SkillManifest 技能包清单
type SkillManifest struct {
	Name        string   `json:"name"`
	DisplayName string   `json:"display_name"`
	Version     string   `json:"version"`
	Description string   `json:"description"`
	Category    string   `json:"category"`
	Tags        []string `json:"tags"`
	Author      string   `json:"author"`
	License     string   `json:"license"`
	Homepage    string   `json:"homepage"`
	Repository  string   `json:"repository"`
	MCPVersion  string   `json:"mcp_version"`
	Entrypoint  string   `json:"entrypoint"`
	Platforms   []string `json:"platforms"`
}

// UploadSkillPackage 上传技能包（带事务处理）
func (s *Service) UploadSkillPackage(ctx context.Context, file *multipart.FileHeader, userID uint) (*Skill, *SkillVersion, error) {
	// 检查文件大小
	if file.Size > MaxPackageSize {
		return nil, nil, ErrPackageTooLarge
	}

	// 打开上传的文件
	src, err := file.Open()
	if err != nil {
		return nil, nil, err
	}
	defer src.Close()

	// 读取文件内容
	data, err := io.ReadAll(src)
	if err != nil {
		return nil, nil, err
	}

	// 解析 ZIP 包
	zipReader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return nil, nil, ErrInvalidSkillPackage
	}

	// 查找并验证 manifest
	manifest, err := s.extractManifest(zipReader)
	if err != nil {
		return nil, nil, err
	}

	// 验证技能名称
	if !isValidSkillName(manifest.Name) {
		return nil, nil, errors.New("无效的技能名称（只能包含小写字母、数字和连字符）")
	}

	var skill *Skill
	var version *SkillVersion
	var uploadedObjectName string

	// 使用事务确保数据一致性
	err = s.db.Transaction(func(tx *gorm.DB) error {
		// 检查技能是否已存在且用户有权限
		existingSkill, err := s.GetSkillByUUID(ctx, manifest.Name)
		if err != nil {
			if err == ErrSkillNotFound {
				// 创建新技能
				skill = &Skill{
					UUID:        uuid.New().String(),
					Name:        manifest.Name,
					DisplayName: manifest.DisplayName,
					Description: manifest.Description,
					Category:    manifest.Category,
					AuthorID:    userID,
					License:     manifest.License,
					Homepage:    manifest.Homepage,
					Repository:  manifest.Repository,
					Status:      "pending", // 待审核
				}

				if err := tx.Create(skill).Error; err != nil {
					return fmt.Errorf("创建技能失败：%w", err)
				}
			} else {
				return err
			}
		} else {
			skill = existingSkill
			// 验证权限
			if skill.AuthorID != userID {
				return ErrUnauthorized
			}
		}

		// 上传技能包到对象存储
		uploadedObjectName = fmt.Sprintf("skills/%s/%s/%s.zip", skill.UUID, manifest.Version, uuid.New().String())
		err = s.storage.Upload(ctx, uploadedObjectName, bytes.NewReader(data), file.Size, "application/zip")
		if err != nil {
			return fmt.Errorf("上传技能包失败：%w", err)
		}

		// 创建版本记录
		version = &SkillVersion{
			SkillID:     skill.ID,
			Version:     manifest.Version,
			Description: manifest.Description,
			DownloadURL: uploadedObjectName,
			Size:        file.Size,
			MCPVersion:  manifest.MCPVersion,
			IsLatest:    true,
		}

		// 更新旧版本的 is_latest 标记
		if err := tx.Model(&SkillVersion{}).
			Where("skill_id = ? AND is_latest = ?", skill.ID, true).
			Update("is_latest", false).Error; err != nil {
			// 回滚时删除已上传的文件
			go s.storage.Delete(context.Background(), uploadedObjectName)
			return fmt.Errorf("更新版本标记失败：%w", err)
		}

		// 创建新版本
		if err := tx.Create(version).Error; err != nil {
			// 回滚时删除已上传的文件
			go s.storage.Delete(context.Background(), uploadedObjectName)
			return fmt.Errorf("创建版本失败：%w", err)
		}

		// 更新技能信息
		if err := tx.Model(skill).Updates(map[string]interface{}{
			"latest_version": manifest.Version,
			"updated_at":     time.Now(),
		}).Error; err != nil {
			go s.storage.Delete(context.Background(), uploadedObjectName)
			return fmt.Errorf("更新技能信息失败：%w", err)
		}

		return nil
	})

	if err != nil {
		return nil, nil, err
	}

	return skill, version, nil
}

// extractManifest 从 ZIP 中提取并验证 manifest
func (s *Service) extractManifest(zipReader *zip.Reader) (*SkillManifest, error) {
	var manifestFile *zip.File

	// 查找 mcp-manifest.json
	for _, file := range zipReader.File {
		if file.Name == "mcp-manifest.json" || file.Name == "manifest.json" {
			manifestFile = file
			break
		}
	}

	if manifestFile == nil {
		return nil, ErrMissingManifest
	}

	// 读取 manifest 内容
	rc, err := manifestFile.Open()
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	var manifest SkillManifest
	if err := json.NewDecoder(rc).Decode(&manifest); err != nil {
		return nil, ErrInvalidManifest
	}

	// 验证必填字段
	if manifest.Name == "" {
		return nil, errors.New("缺少必填字段：name")
	}
	if manifest.Version == "" {
		return nil, errors.New("缺少必填字段：version")
	}
	if manifest.Description == "" {
		return nil, errors.New("缺少必填字段：description")
	}

	// 验证版本号格式
	if !isValidVersion(manifest.Version) {
		return nil, errors.New("版本号格式无效（应为语义化版本，如 1.0.0）")
	}

	return &manifest, nil
}

// ValidateSkillPackage 验证技能包（不上传，仅验证）
func (s *Service) ValidateSkillPackage(ctx context.Context, file *multipart.FileHeader) (*SkillManifest, error) {
	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	data, err := io.ReadAll(src)
	if err != nil {
		return nil, err
	}

	zipReader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return nil, ErrInvalidSkillPackage
	}

	return s.extractManifest(zipReader)
}

// isValidSkillName 验证技能名称格式
func isValidSkillName(name string) bool {
	matched, _ := regexp.MatchString("^[a-z][a-z0-9-]*$", name)
	return matched
}

// isValidVersion 验证语义化版本号
func isValidVersion(version string) bool {
	matched, _ := regexp.MatchString(`^v?\d+\.\d+\.\d+(-[a-zA-Z0-9]+)?(\+[a-zA-Z0-9]+)?$`, version)
	return matched
}

// DownloadSkillPackage 下载技能包
func (s *Service) DownloadSkillPackage(ctx context.Context, skillID uint, version string, userID uint, ip string, ua string) (io.ReadCloser, int64, error) {
	// 获取技能
	skill, err := s.GetSkillByID(ctx, skillID)
	if err != nil {
		return nil, 0, err
	}

	// 获取版本
	var skillVersion SkillVersion
	query := s.db.Where("skill_id = ?", skillID)
	if version != "" {
		query = query.Where("version = ?", version)
	} else {
		query = query.Where("is_latest = ?", true)
	}

	if err := query.First(&skillVersion).Error; err != nil {
		return nil, 0, errors.New("版本不存在")
	}

	// 记录下载
	if userID > 0 {
		go s.IncrementDownload(ctx, skillID, skillVersion.ID, userID, ip, ua)
	}

	// 从对象存储获取文件
	reader, err := s.storage.Download(ctx, skillVersion.DownloadURL)
	if err != nil {
		return nil, 0, fmt.Errorf("下载技能包失败：%w", err)
	}

	return reader, skillVersion.Size, nil
}

// ListSkillVersions 列出技能的所有版本
func (s *Service) ListSkillVersions(ctx context.Context, skillID uint) ([]SkillVersion, error) {
	var versions []SkillVersion
	err := s.db.Where("skill_id = ?", skillID).
		Order("created_at DESC").
		Find(&versions).Error
	return versions, err
}

// GetVersionByNumber 获取指定版本号
func (s *Service) GetVersionByNumber(ctx context.Context, skillID uint, version string) (*SkillVersion, error) {
	var v SkillVersion
	err := s.db.Where("skill_id = ? AND version = ?", skillID, version).First(&v).Error
	if err != nil {
		return nil, errors.New("版本不存在")
	}
	return &v, nil
}

// DeleteVersion 删除版本（仅管理员或作者）
func (s *Service) DeleteVersion(ctx context.Context, skillID uint, version string, userID uint) error {
	var skill Skill
	if err := s.db.First(&skill, skillID).Error; err != nil {
		return err
	}

	// 验证权限
	if skill.AuthorID != userID {
		return ErrUnauthorized
	}

	var v SkillVersion
	if err := s.db.Where("skill_id = ? AND version = ?", skillID, version).First(&v).Error; err != nil {
		return err
	}

	// 从存储中删除
	if err := s.storage.Delete(ctx, v.DownloadURL); err != nil {
		return err
	}

	// 删除数据库记录
	return s.db.Delete(&v).Error
}

// SearchVersions 搜索特定条件的版本
func (s *Service) SearchVersions(ctx context.Context, opts VersionSearchOptions) ([]SkillVersion, error) {
	var versions []SkillVersion

	query := s.db.Model(&SkillVersion{}).
		Joins("JOIN skills ON skills.id = skill_versions.skill_id").
		Where("skills.status = ?", "approved").
		Where("skills.is_private = ?", false)

	if opts.Category != "" {
		query = query.Where("skills.category = ?", opts.Category)
	}

	if opts.MCPVersion != "" {
		query = query.Where("mcp_version = ?", opts.MCPVersion)
	}

	if opts.Platform != "" {
		// 简化实现，实际应关联 platform 表
		query = query.Where("platforms LIKE ?", "%"+opts.Platform+"%")
	}

	err := query.Order("skill_versions.created_at DESC").
		Limit(opts.Limit).
		Find(&versions).Error

	return versions, err
}

// VersionSearchOptions 版本搜索选项
type VersionSearchOptions struct {
	Category   string
	MCPVersion string
	Platform   string
	Limit      int
}

// GetSkillStats 获取技能统计信息
func (s *Service) GetSkillStats(ctx context.Context, skillID uint) (*SkillStats, error) {
	var stats SkillStats

	// 获取版本数量
	var versionCount int64
	s.db.Model(&SkillVersion{}).Where("skill_id = ?", skillID).Count(&versionCount)
	stats.VersionCount = versionCount

	// 获取下载趋势（最近 7 天）
	// 简化实现
	stats.TotalDownloads = 0
	stats.TrendingScore = 0.0

	// 获取评分分布
	var ratings []struct {
		Rating int   `json:"rating"`
		Count  int64 `json:"count"`
	}
	s.db.Model(&SkillRating{}).
		Where("skill_id = ?", skillID).
		Select("rating, COUNT(*) as count").
		Group("rating").
		Scan(&ratings)
	stats.RatingDistribution = ratings

	return &stats, nil
}

// SkillStats 技能统计
type SkillStats struct {
	VersionCount       int64   `json:"version_count"`
	TotalDownloads     int64   `json:"total_downloads"`
	TrendingScore      float64 `json:"trending_score"`
	RatingDistribution []struct {
		Rating int   `json:"rating"`
		Count  int64 `json:"count"`
	} `json:"rating_distribution"`
}

// extractFileName 从路径提取文件名
func extractFileName(path string) string {
	return filepath.Base(strings.TrimSuffix(path, "/"))
}
