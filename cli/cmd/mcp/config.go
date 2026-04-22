package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mcp-skill-hub/cli/internal/config"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "管理配置",
	Long:  `查看、编辑和管理 MCP CLI 配置文件。`,
}

var configViewCmd = &cobra.Command{
	Use:   "view",
	Short: "查看当前配置",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("加载配置失败：%w", err)
		}

		data, err := yaml.Marshal(cfg)
		if err != nil {
			return err
		}

		fmt.Println(string(data))
		return nil
	},
}

var configEditCmd = &cobra.Command{
	Use:   "edit",
	Short: "编辑配置文件",
	RunE: func(cmd *cobra.Command, args []string) error {
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}

		configPath := filepath.Join(home, ".mcp", "config.yaml")
		
		// 确保目录存在
		if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
			return err
		}

		// 如果文件不存在，创建默认配置
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			defaultCfg := &config.Config{
				Server: config.ServerConfig{
					URL: "http://localhost:8080",
				},
			}
			if err := config.Save(defaultCfg); err != nil {
				return err
			}
			fmt.Printf("已创建默认配置文件：%s\n", configPath)
		}

		// 打开编辑器
		editor := os.Getenv("EDITOR")
		if editor == "" {
			editor = "vi"
		}

		cmdExec := exec.Command(editor, configPath)
		cmdExec.Stdin = os.Stdin
		cmdExec.Stdout = os.Stdout
		cmdExec.Stderr = os.Stderr
		return cmdExec.Run()
	},
}

var configSetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "设置配置项",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("加载配置失败：%w", err)
		}

		key, value := args[0], args[1]

		switch key {
		case "server.url":
			cfg.Server.URL = value
		case "server.token":
			cfg.Server.Token = value
		default:
			return fmt.Errorf("未知的配置项：%s", key)
		}

		if err := config.Save(cfg); err != nil {
			return fmt.Errorf("保存配置失败：%w", err)
		}

		fmt.Printf("✅ 已设置 %s = %s\n", key, value)
		return nil
	},
}

var configGetCmd = &cobra.Command{
	Use:   "get <key>",
	Short: "获取配置项",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("加载配置失败：%w", err)
		}

		key := args[0]
		switch key {
		case "server.url":
			fmt.Println(cfg.Server.URL)
		case "server.token":
			fmt.Println(cfg.Server.Token)
		default:
			return fmt.Errorf("未知的配置项：%s", key)
		}

		return nil
	},
}

func init() {
	configCmd.AddCommand(configViewCmd)
	configCmd.AddCommand(configEditCmd)
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configGetCmd)
	rootCmd.AddCommand(configCmd)
}
