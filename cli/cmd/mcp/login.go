package main

import (
	"fmt"
	"os"

	"github.com/mcp-skill-hub/cli/internal/api"
	"github.com/mcp-skill-hub/cli/internal/config"
	"github.com/spf13/cobra"
)

var (
	loginUsername string
	loginPassword string
)

var loginCmdDetailed = &cobra.Command{
	Use:   "login",
	Short: "登录到 MCP Skill Hub",
	Long:  `登录到 MCP Skill Hub 服务器，保存认证令牌用于后续操作。`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// 加载配置
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("加载配置失败：%w", err)
		}

		// 如果未提供用户名密码，提示输入
		if loginUsername == "" || loginPassword == "" {
			fmt.Print("用户名：")
			fmt.Scanln(&loginUsername)
			fmt.Print("密码：")
			fmt.Scanln(&loginPassword)
		}

		// 创建客户端
		client := api.NewClient(cfg.Server.URL, "")

		// 登录
		fmt.Println("🔐 正在登录...")
		resp, err := client.Login(loginUsername, loginPassword)
		if err != nil {
			return fmt.Errorf("登录失败：%w", err)
		}

		// 保存令牌
		cfg.Server.Token = resp.AccessToken
		if err := config.Save(cfg); err != nil {
			return fmt.Errorf("保存配置失败：%w", err)
		}

		fmt.Println("✅ 登录成功！")
		fmt.Printf("📝 令牌已保存到 ~/.mcp/config.yaml\n")
		fmt.Printf("⏰ 令牌有效期：%d 小时\n", resp.ExpiresIn/3600)

		return nil
	},
}

func init() {
	loginCmdDetailed.Flags().StringVarP(&loginUsername, "username", "u", "", "用户名")
	loginCmdDetailed.Flags().StringVarP(&loginPassword, "password", "p", "", "密码")
}
