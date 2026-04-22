package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config CLI 配置
type Config struct {
	Server ServerConfig `yaml:"server"`
	Proxy  ProxyConfig  `yaml:"proxy"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	URL   string `yaml:"url"`
	Token string `yaml:"token"`
}

// ProxyConfig 代理配置
type ProxyConfig struct {
	HTTP  string `yaml:"http"`
	HTTPS string `yaml:"https"`
}

// Load 加载配置文件
func Load() (*Config, error) {
	// 获取配置路径
	configPath := os.Getenv("MCP_CONFIG")
	if configPath == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("获取用户目录失败：%w", err)
		}
		configPath = filepath.Join(home, ".mcp", "config.yaml")
	}

	// 检查文件是否存在
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// 返回默认配置
		return &Config{
			Server: ServerConfig{
				URL: "http://localhost:8080",
			},
		}, nil
	}

	// 读取文件
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败：%w", err)
	}

	// 解析 YAML
	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("解析配置文件失败：%w", err)
	}

	// 设置默认值
	if config.Server.URL == "" {
		config.Server.URL = "http://localhost:8080"
	}

	return &config, nil
}

// Save 保存配置文件
func Save(cfg *Config) error {
	// 创建目录
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("获取用户目录失败：%w", err)
	}

	configDir := filepath.Join(home, ".mcp")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("创建配置目录失败：%w", err)
	}

	configPath := filepath.Join(configDir, "config.yaml")

	// 序列化 YAML
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("序列化配置失败：%w", err)
	}

	// 写入文件
	if err := os.WriteFile(configPath, data, 0600); err != nil {
		return fmt.Errorf("写入配置文件失败：%w", err)
	}

	return nil
}

// EnsureConfigDir 确保配置目录存在
func EnsureConfigDir() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configDir := filepath.Join(home, ".mcp")
	return os.MkdirAll(configDir, 0755)
}
