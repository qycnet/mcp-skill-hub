package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	uninstallForce bool
)

var uninstallCmdDetailed = &cobra.Command{
	Use:   "uninstall [技能名]",
	Short: "卸载技能",
	Long:  `从本地卸载已安装的 MCP 技能。`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		skillName := args[0]

		// 获取用户目录
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("获取用户目录失败：%w", err)
		}

		skillsDir := filepath.Join(home, ".mcp", "skills")
		skillPath := filepath.Join(skillsDir, skillName)

		// 检查是否已安装
		if _, err := os.Stat(skillPath); os.IsNotExist(err) {
			return fmt.Errorf("技能未安装：%s", skillName)
		}

		// 确认卸载
		if !uninstallForce {
			fmt.Printf("⚠️  确定要卸载技能 '%s' 吗？\n", skillName)
			fmt.Print("输入 'y' 确认：")
			var confirm string
			fmt.Scanln(&confirm)
			if confirm != "y" && confirm != "Y" {
				fmt.Println("取消卸载")
				return nil
			}
		}

		// 删除目录
		fmt.Printf("🗑️  正在卸载：%s\n", skillName)
		if err := os.RemoveAll(skillPath); err != nil {
			return fmt.Errorf("卸载失败：%w", err)
		}

		fmt.Printf("✅ 技能 '%s' 已卸载\n", skillName)
		return nil
	},
}

func init() {
	uninstallCmdDetailed.Flags().BoolVarP(&uninstallForce, "force", "f", false, "强制卸载，不确认")
}
