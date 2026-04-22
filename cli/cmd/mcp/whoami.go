package main

import (
	"fmt"

	"github.com/mcp-skill-hub/cli/internal/api"
	"github.com/mcp-skill-hub/cli/internal/config"
	"github.com/spf13/cobra"
)

var whoamiCmd = &cobra.Command{
	Use:   "whoami",
	Short: "显示当前用户信息",
	Long:  `显示当前登录的用户信息。`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("加载配置失败：%w", err)
		}

		if cfg.Server.Token == "" {
			fmt.Println("未登录")
			fmt.Println("请使用 'mcp login' 登录")
			return nil
		}

		client := api.NewClient(cfg.Server.URL, cfg.Server.Token)

		profile, err := client.GetProfile()
		if err != nil {
			fmt.Printf("⚠️  无法获取用户信息：%v\n", err)
			fmt.Println("Token 可能已过期，请重新登录")
			return nil
		}

		fmt.Printf("👤 用户：%s\n", profile.Username)
		fmt.Printf("📧 邮箱：%s\n", profile.Email)
		fmt.Printf("🎫 订阅：%s\n", profile.SubscriptionStatus)
		fmt.Printf("🔑 Token: %s...\n", cfg.Server.Token[:20])

		return nil
	},
}

func init() {
	rootCmd.AddCommand(whoamiCmd)
}
