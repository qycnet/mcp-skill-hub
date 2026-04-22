package skill

import (
	"time"

	"gorm.io/gorm"
)

// Skill 代表一个 MCP 技能
type Skill struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	UUID        string         `gorm:"uniqueIndex;size:36" json:"uuid"`
	Name        string         `gorm:"uniqueIndex;size:100;not null" json:"name"`
	DisplayName string         `gorm:"size:200" json:"display_name"`
	Description string         `gorm:"type:text" json:"description"`
	Category    string         `gorm:"size:50;index" json:"category"`
	Tags        []Tag          `gorm:"many2many:skill_tags;" json:"tags"`
	AuthorID    uint           `gorm:"index" json:"author_id"`
	Author      User           `gorm:"foreignKey:AuthorID" json:"author,omitempty"`
	Repository  string         `gorm:"size:500" json:"repository"`
	Homepage    string         `gorm:"size:500" json:"homepage"`
	License     string         `gorm:"size:50" json:"license"`
	
	// 版本信息
	LatestVersion string       `gorm:"size:20" json:"latest_version"`
	Versions      []SkillVersion `gorm:"foreignKey:SkillID" json:"versions,omitempty"`
	
	// 统计信息
	Downloads     int64        `gorm:"default:0" json:"downloads"`
	Rating        float64      `gorm:"default:0" json:"rating"`
	RatingCount   int64        `gorm:"default:0" json:"rating_count"`
	
	// 质量评分
	QualityScore  float64      `gorm:"default:0" json:"quality_score"`
	
	// 状态
	Status        string       `gorm:"size:20;default:pending" json:"status"` // pending, approved, rejected, archived
	IsPrivate     bool         `gorm:"default:false" json:"is_private"`
	IsFeatured    bool         `gorm:"default:false" json:"is_featured"`
	IsVerified    bool         `gorm:"default:false" json:"is_verified"`
	
	// 安全
	SecurityScanStatus string `gorm:"size:20;default:pending" json:"security_scan_status"` // pending, passed, failed
	SecurityIssues     []SecurityIssue `gorm:"foreignKey:SkillID" json:"security_issues,omitempty"`
	
	// 时间戳
	CreatedAt     time.Time    `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time    `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

// Tag 技能标签
type Tag struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `gorm:"uniqueIndex;size:50;not null" json:"name"`
}

// SkillVersion 技能版本
type SkillVersion struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	SkillID     uint           `gorm:"index;not null" json:"skill_id"`
	Version     string         `gorm:"size:20;not null" json:"version"` // 语义化版本：1.2.3
	Description string         `gorm:"type:text" json:"description"`
	Changelog   string         `gorm:"type:text" json:"changelog"`
	
	// 构建信息
	BuildHash   string       `gorm:"size:64" json:"build_hash"`
	DownloadURL string       `gorm:"size:500" json:"download_url"`
	Size        int64        `gorm:"default:0" json:"size"` // 字节
	
	// 依赖
	Dependencies  []Dependency `gorm:"foreignKey:VersionID" json:"dependencies,omitempty"`
	MCPVersion    string       `gorm:"size:20" json:"mcp_version"` // 兼容的 MCP 协议版本
	
	// 平台支持
	Platforms     []Platform   `gorm:"many2many:version_platforms;" json:"platforms,omitempty"`
	
	// 状态
	IsLatest      bool         `gorm:"default:false" json:"is_latest"`
	IsDeprecated  bool         `gorm:"default:false" json:"is_deprecated"`
	
	// 时间戳
	CreatedAt     time.Time    `gorm:"autoCreateTime" json:"created_at"`
}

// Dependency 依赖
type Dependency struct {
	ID         uint   `gorm:"primaryKey" json:"id"`
	VersionID  uint   `gorm:"index;not null" json:"version_id"`
	Name       string `gorm:"size:100;not null" json:"name"`
	Version    string `gorm:"size:20;not null" json:"version"`
	Optional   bool   `gorm:"default:false" json:"optional"`
}

// Platform 支持的平台
type Platform struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `gorm:"uniqueIndex;size:50;not null" json:"name"` // linux, darwin, windows, etc.
	Arch string `gorm:"size:20" json:"arch"` // amd64, arm64, etc.
}

// SecurityIssue 安全问题
type SecurityIssue struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	SkillID     uint      `gorm:"index;not null" json:"skill_id"`
	Severity    string    `gorm:"size:20" json:"severity"` // low, medium, high, critical
	CVE         string    `gorm:"size:50" json:"cve"`
	Description string    `gorm:"type:text" json:"description"`
	FixedIn     string    `gorm:"size:20" json:"fixed_in"`
	ReportedAt  time.Time `gorm:"autoCreateTime" json:"reported_at"`
}

// User 作者信息（简化版，完整见 auth 包）
type User struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Username string `gorm:"uniqueIndex;size:50;not null" json:"username"`
	Email    string `gorm:"uniqueIndex;size:255" json:"email"`
	Avatar   string `gorm:"size:500" json:"avatar"`
	Bio      string `gorm:"type:text" json:"bio"`
}

// SkillRating 技能评分
type SkillRating struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	SkillID   uint      `gorm:"index;not null" json:"skill_id"`
	UserID    uint      `gorm:"index;not null" json:"user_id"`
	Rating    int       `gorm:"not null" json:"rating"` // 1-5
	Comment   string    `gorm:"type:text" json:"comment"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// SkillDownload 下载记录（用于统计）
type SkillDownload struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	SkillID   uint      `gorm:"index;not null" json:"skill_id"`
	VersionID uint      `gorm:"index" json:"version_id"`
	UserID    uint      `gorm:"index" json:"user_id"`
	IPAddress string    `gorm:"size:45" json:"ip_address"`
	UserAgent string    `gorm:"type:text" json:"user_agent"`
	DownloadedAt time.Time `gorm:"autoCreateTime" json:"downloaded_at"`
}
