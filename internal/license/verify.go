package license

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"
)

// License 许可证信息
type License struct {
	Key           string
	Type          string // professional, enterprise
	ValidUntil    time.Time
	MaxUsers      int
	Features      []string
	Activated     bool
	ActivatedAt   *time.Time
}

// Verify 验证许可证
func Verify(licenseKey string) (*License, error) {
	// 简化版验证，实际应该调用 API 验证
	if len(licenseKey) < 32 {
		return nil, errors.New("许可证密钥格式无效")
	}
	
	license := &License{
		Key:        licenseKey,
		Type:       "enterprise",
		ValidUntil: time.Now().AddDate(1, 0, 0),
		MaxUsers:   100,
		Features: []string{
			"payment",
			"subscription",
			"sso",
			"analytics",
			"audit",
			"monitoring",
		},
		Activated: false,
	}
	
	return license, nil
}

// Activate 激活许可证
func (l *License) Activate() error {
	if l.Activated {
		return errors.New("许可证已激活")
	}
	
	now := time.Now()
	l.Activated = true
	l.ActivatedAt = &now
	
	return nil
}

// IsValid 检查许可证是否有效
func (l *License) IsValid() bool {
	if !l.Activated {
		return false
	}
	
	if time.Now().After(l.ValidUntil) {
		return false
	}
	
	return true
}

// HasFeature 检查是否有某功能
func (l *License) HasFeature(feature string) bool {
	for _, f := range l.Features {
		if f == feature {
			return true
		}
	}
	return false
}

// GenerateTrialKey 生成试用许可证密钥
func GenerateTrialKey() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	
	key := fmt.Sprintf("TRIAL-%s", hex.EncodeToString(bytes))
	return key, nil
}

// GetLicenseType 获取许可证类型
func GetLicenseType(key string) string {
	if len(key) < 6 {
		return "invalid"
	}
	
	prefix := key[:6]
	switch prefix {
	case "TRIAL-":
		return "trial"
	case "PRO-":
		return "professional"
	case "ENT-":
		return "enterprise"
	default:
		return "unknown"
	}
}
