package i18n

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"strings"

	"golang.org/x/text/language"
)

//go:embed translations/*.json
var localeFiles embed.FS

// Translator 翻译器
type Translator struct {
	locales map[string]map[string]string
}

// NewTranslator 创建翻译器
func NewTranslator() (*Translator, error) {
	t := &Translator{
		locales: make(map[string]map[string]string),
	}

	// 加载所有语言文件
	err := fs.WalkDir(localeFiles, "translations", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}

		// 读取文件
		data, err := localeFiles.ReadFile(path)
		if err != nil {
			return err
		}

		// 解析 JSON
		var translations map[string]string
		if err := json.Unmarshal(data, &translations); err != nil {
			return err
		}

		// 提取语言代码
		lang := strings.TrimSuffix(strings.TrimPrefix(d.Name(), "locale_"), ".json")
		t.locales[lang] = translations

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("加载语言文件失败：%w", err)
	}

	return t, nil
}

// T 翻译文本
func (t *Translator) T(lang, key string, args ...interface{}) string {
	translations, ok := t.locales[lang]
	if !ok {
		// 回退到英语
		translations = t.locales["en"]
	}

	message, ok := translations[key]
	if !ok {
		// 回退到键名
		return key
	}

	// 替换参数
	for i, arg := range args {
		placeholder := fmt.Sprintf("{%d}", i)
		message = strings.ReplaceAll(message, placeholder, fmt.Sprintf("%v", arg))
	}

	return message
}

// MustT 翻译文本（不处理错误）
func (t *Translator) MustT(lang, key string, args ...interface{}) string {
	return t.T(lang, key, args...)
}

// GetSupportedLanguages 获取支持的语言
func (t *Translator) GetSupportedLanguages() []language.Tag {
	tags := []language.Tag{
		language.Chinese,
		language.English,
		language.Japanese,
		language.Spanish,
		language.French,
	}
	return tags
}

// GetLanguageName 获取语言名称
func (t *Translator) GetLanguageName(lang string) string {
	names := map[string]string{
		"zh": "中文",
		"en": "English",
		"ja": "日本語",
		"es": "Español",
		"fr": "Français",
	}
	return names[lang]
}

// DetectLanguage 从 Accept-Language header 检测语言
func (t *Translator) DetectLanguage(acceptLang string) string {
	// 简单实现，生产环境应使用更复杂的检测逻辑
	if strings.Contains(acceptLang, "zh") {
		return "zh"
	}
	if strings.Contains(acceptLang, "ja") {
		return "ja"
	}
	if strings.Contains(acceptLang, "es") {
		return "es"
	}
	if strings.Contains(acceptLang, "fr") {
		return "fr"
	}
	return "en"
}

// GlobalTranslator 全局翻译器
var GlobalTranslator *Translator

// InitGlobalTranslator 初始化全局翻译器
func InitGlobalTranslator() error {
	var err error
	GlobalTranslator, err = NewTranslator()
	return err
}

// T 全局翻译函数
func T(lang, key string, args ...interface{}) string {
	if GlobalTranslator == nil {
		return key
	}
	return GlobalTranslator.T(lang, key, args...)
}
