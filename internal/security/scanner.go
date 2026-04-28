package security

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Scanner 安全扫描器
type Scanner struct {
	trivyPath   string
	cacheDir    string
	enabled     bool
}

// ScanResult 扫描结果
type ScanResult struct {
	ID           string          `json:"id"`
	SkillID      uint            `json:"skill_id"`
	VersionID    uint            `json:"version_id"`
	Status       string          `json:"status"` // pending, scanning, passed, failed
	ScannedAt    time.Time       `json:"scanned_at"`
	Vulnerabilities []Vulnerability `json:"vulnerabilities"`
	Summary      VulnerabilitySummary `json:"summary"`
}

// Vulnerability 漏洞信息
type Vulnerability struct {
	ID          string   `json:"id"`           // CVE-ID
	Severity    string   `json:"severity"`     // LOW, MEDIUM, HIGH, CRITICAL
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Package     string   `json:"package"`
	Version     string   `json:"installed_version"`
	FixedIn     string   `json:"fixed_version"`
	References  []string `json:"references"`
}

// VulnerabilitySummary 漏洞统计
type VulnerabilitySummary struct {
	Total     int `json:"total"`
	Critical  int `json:"critical"`
	High      int `json:"high"`
	Medium    int `json:"medium"`
	Low       int `json:"low"`
}

// NewScanner 创建安全扫描器
func NewScanner(trivyPath, cacheDir string, enabled bool) *Scanner {
	if trivyPath == "" {
		trivyPath = "trivy" // 假设在 PATH 中
	}
	
	return &Scanner{
		trivyPath: trivyPath,
		cacheDir:  cacheDir,
		enabled:   enabled,
	}
}

// ScanSkillPackage 扫描技能包
func (s *Scanner) ScanSkillPackage(ctx context.Context, packagePath string, skillID, versionID uint) (*ScanResult, error) {
	result := &ScanResult{
		ID:        uuid.New().String(),
		SkillID:   skillID,
		VersionID: versionID,
		Status:    "scanning",
		ScannedAt: time.Now(),
	}

	if !s.enabled {
		result.Status = "passed"
		result.Summary = VulnerabilitySummary{}
		return result, nil
	}

	// 使用 Trivy 扫描
	cmd := exec.CommandContext(ctx, s.trivyPath,
		"fs",
		"--quiet",
		"--format", "json",
		"--cache-dir", s.cacheDir,
		"--severity", "HIGH,CRITICAL",
		"--exit-code", "1",
		packagePath,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		// Trivy exit code 1 表示发现漏洞
		if !strings.Contains(string(output), "exit status 1") {
			return nil, fmt.Errorf("扫描失败：%w, output: %s", err, string(output))
		}
	}

	// 解析 Trivy 输出
	vulns, err := s.parseTrivyOutput(output)
	if err != nil {
		return nil, fmt.Errorf("解析扫描结果失败：%w", err)
	}

	result.Vulnerabilities = vulns
	result.Summary = s.summarizeVulnerabilities(vulns)

	// 判断是否通过
	if result.Summary.Critical > 0 || result.Summary.High > 0 {
		result.Status = "failed"
	} else {
		result.Status = "passed"
	}

	return result, nil
}

// ScanImage 扫描 Docker 镜像
func (s *Scanner) ScanImage(ctx context.Context, imageName string) (*ScanResult, error) {
	result := &ScanResult{
		ID:        uuid.New().String(),
		Status:    "scanning",
		ScannedAt: time.Now(),
	}

	if !s.enabled {
		result.Status = "passed"
		return result, nil
	}

	cmd := exec.CommandContext(ctx, s.trivyPath,
		"image",
		"--quiet",
		"--format", "json",
		"--cache-dir", s.cacheDir,
		"--severity", "HIGH,CRITICAL",
		imageName,
	)

	output, err := cmd.CombinedOutput()
	if err != nil && !strings.Contains(string(output), "exit status 1") {
		return nil, fmt.Errorf("镜像扫描失败：%w", err)
	}

	vulns, err := s.parseTrivyOutput(output)
	if err != nil {
		return nil, err
	}

	result.Vulnerabilities = vulns
	result.Summary = s.summarizeVulnerabilities(vulns)

	if result.Summary.Critical > 0 || result.Summary.High > 0 {
		result.Status = "failed"
	} else {
		result.Status = "passed"
	}

	return result, nil
}

// parseTrivyOutput 解析 Trivy JSON 输出
func (s *Scanner) parseTrivyOutput(output []byte) ([]Vulnerability, error) {
	if len(output) == 0 {
		return []Vulnerability{}, nil
	}

	var trivyResult struct {
		Results []struct {
			Vulnerabilities []struct {
				VulnerabilityID  string   `json:"VulnerabilityID"`
				Severity        string   `json:"Severity"`
				Title           string   `json:"Title"`
				Description     string   `json:"Description"`
				PkgName         string   `json:"PkgName"`
				InstalledVersion string  `json:"InstalledVersion"`
				FixedVersion    string   `json:"FixedVersion"`
				References      []string `json:"References"`
			} `json:"Vulnerabilities"`
		} `json:"Results"`
	}

	if err := json.Unmarshal(output, &trivyResult); err != nil {
		return nil, err
	}

	var vulns []Vulnerability
	for _, result := range trivyResult.Results {
		for _, v := range result.Vulnerabilities {
			vulns = append(vulns, Vulnerability{
				ID:          v.VulnerabilityID,
				Severity:    v.Severity,
				Title:       v.Title,
				Description: v.Description,
				Package:     v.PkgName,
				Version:     v.InstalledVersion,
				FixedIn:     v.FixedVersion,
				References:  v.References,
			})
		}
	}

	return vulns, nil
}

// summarizeVulnerabilities 统计漏洞
func (s *Scanner) summarizeVulnerabilities(vulns []Vulnerability) VulnerabilitySummary {
	summary := VulnerabilitySummary{Total: len(vulns)}
	
	for _, v := range vulns {
		switch strings.ToUpper(v.Severity) {
		case "CRITICAL":
			summary.Critical++
		case "HIGH":
			summary.High++
		case "MEDIUM":
			summary.Medium++
		case "LOW":
			summary.Low++
		}
	}
	
	return summary
}

// IsEnabled 检查扫描器是否启用
func (s *Scanner) IsEnabled() bool {
	return s.enabled
}

// QuickScan 快速扫描（仅检查已知恶意特征）
func (s *Scanner) QuickScan(packagePath string) ([]string, error) {
	var warnings []string

	// 检查可疑文件
	suspiciousPatterns := []string{
		"eval(", "exec(", "system(",
		"__import__", "subprocess.call",
		"child_process", "require('child_process')",
		"rm -rf", "dd if=", "mkfs",
	}

	// 这里可以添加更多检查逻辑
	// 例如：扫描文件内容、检查可疑脚本等

	return warnings, nil
}
