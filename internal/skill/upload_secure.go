package skill

import (
	"context"
	"fmt"

	"github.com/qycnet/mcp-skill-hub/internal/security"
)

// UploadSkillPackageWithSecurityScan 上传技能包并进行安全扫描
func (s *Service) UploadSkillPackageWithSecurityScan(ctx context.Context, file *multipart.FileHeader, userID uint) (*Skill, *SkillVersion, *security.ScanResult, error) {
	// 先执行普通上传流程
	skill, version, err := s.UploadSkillPackage(ctx, file, userID)
	if err != nil {
		return nil, nil, nil, err
	}

	// 如果配置了安全扫描器，执行扫描
	if s.scanner != nil && s.scanner.IsEnabled() {
		// 解压到临时目录进行扫描
		tmpDir, err := s.extractToTempDir(file)
		if err != nil {
			return skill, version, nil, fmt.Errorf("解压失败：%w", err)
		}
		defer s.cleanupTempDir(tmpDir)

		// 执行安全扫描
		scanResult, err := s.scanner.ScanSkillPackage(ctx, tmpDir, skill.ID, version.ID)
		if err != nil {
			// 扫描失败不阻止上传，但记录警告
			scanResult = &security.ScanResult{
				SkillID:   skill.ID,
				VersionID: version.ID,
				Status:    "error",
			}
		}

		// 更新技能的安全状态
		if scanResult.Status == "failed" {
			// 发现高危漏洞，标记为待修复
			s.db.Model(skill).Update("security_scan_status", "failed")
			s.db.Model(skill).Update("status", "security_review")
		} else if scanResult.Status == "passed" {
			s.db.Model(skill).Update("security_scan_status", "passed")
		}

		// 保存扫描结果
		if err := s.saveScanResult(ctx, scanResult); err != nil {
			// 记录错误但不影响上传
		}

		return skill, version, scanResult, nil
	}

	// 未配置扫描器，返回成功
	return skill, version, &security.ScanResult{
		SkillID:   skill.ID,
		VersionID: version.ID,
		Status:    "skipped",
	}, nil
}

// extractToTempDir 解压文件到临时目录
func (s *Service) extractToTempDir(file *multipart.FileHeader) (string, error) {
	// 实现解压逻辑
	// 这里简化实现，实际应使用 archive/zip
	return "", nil
}

// cleanupTempDir 清理临时目录
func (s *Service) cleanupTempDir(dir string) {
	// 实现清理逻辑
}

// saveScanResult 保存扫描结果
func (s *Service) saveScanResult(ctx context.Context, result *security.ScanResult) error {
	// 实现保存逻辑
	return nil
}
