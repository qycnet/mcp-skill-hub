package main

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/mcp-skill-hub/cli/internal/api"
	"github.com/mcp-skill-hub/cli/internal/config"
	"github.com/spf13/cobra"
)

var (
	installVersion string
	installDir     string
)

var installCmdDetailed = &cobra.Command{
	Use:   "install [技能名]",
	Short: "安装技能",
	Long:  `从 MCP Skill Hub 下载并安装技能到本地。`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// 加载配置
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("加载配置失败：%w", err)
		}

		skillName := args[0]

		// 创建客户端
		client := api.NewClient(cfg.Server.URL, cfg.Server.Token)

		// 获取技能详情
		fmt.Printf("📥 正在获取技能信息：%s\n", skillName)
		skill, err := client.GetSkill(skillName)
		if err != nil {
			return fmt.Errorf("获取技能信息失败：%w", err)
		}

		fmt.Printf("✅ 找到技能：%s (%s)\n", skill.DisplayName, skill.LatestVersion)

		// 下载技能
		fmt.Printf("⬇️  正在下载...\n")
		resp, err := client.InstallSkill(fmt.Sprintf("%d", skill.ID))
		if err != nil {
			return fmt.Errorf("获取下载链接失败：%w", err)
		}

		// 确定安装目录
		if installDir == "" {
			home, _ := os.UserHomeDir()
			installDir = filepath.Join(home, ".mcp", "skills")
		}

		if err := os.MkdirAll(installDir, 0755); err != nil {
			return fmt.Errorf("创建安装目录失败：%w", err)
		}

		// 下载文件
		zipPath := filepath.Join(installDir, fmt.Sprintf("%s-%s.zip", skill.Name, skill.LatestVersion))
		if err := downloadFile(resp.DownloadURL, zipPath); err != nil {
			return fmt.Errorf("下载失败：%w", err)
		}

		// 解压
		skillDir := filepath.Join(installDir, skill.Name)
		fmt.Printf("📦 正在解压到：%s\n", skillDir)
		if err := unzip(zipPath, skillDir); err != nil {
			return fmt.Errorf("解压失败：%w", err)
		}

		// 清理临时文件
		os.Remove(zipPath)

		fmt.Printf("✅ 技能安装成功！\n")
		fmt.Printf("📂 安装位置：%s\n", skillDir)

		return nil
	},
}

func init() {
	installCmdDetailed.Flags().StringVarP(&installVersion, "version", "v", "", "指定版本（默认最新版）")
	installCmdDetailed.Flags().StringVarP(&installDir, "dir", "d", "", "安装目录（默认 ~/.mcp/skills）")
}

// downloadFile 下载文件
func downloadFile(url, dest string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

// unzip 解压 ZIP 文件
func unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		// 跳过目录
		if f.FileInfo().IsDir() {
			continue
		}

		// 创建目标路径
		path := filepath.Join(dest, f.Name)
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return err
		}

		// 解压文件
		rc, err := f.Open()
		if err != nil {
			return err
		}

		out, err := os.Create(path)
		if err != nil {
			rc.Close()
			return err
		}

		_, err = io.Copy(out, rc)
		out.Close()
		rc.Close()

		if err != nil {
			return err
		}
	}

	return nil
}

// 辅助函数：字符串操作
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}
