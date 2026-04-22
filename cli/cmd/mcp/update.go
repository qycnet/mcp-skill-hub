package main

import (
	"archive/zip"
	"encoding/json"
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
	updateAll    bool
	updateDryRun bool
)

var updateCmdDetailed = &cobra.Command{
	Use:   "update [技能名]",
	Short: "更新技能",
	Long:  `更新已安装的 MCP 技能到最新版本。`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// 加载配置
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("加载配置失败：%w", err)
		}

		client := api.NewClient(cfg.Server.URL, cfg.Server.Token)

		// 获取已安装的技能列表
		home, _ := os.UserHomeDir()
		skillsDir := filepath.Join(home, ".mcp", "skills")

		var skillsToUpdate []string

		if updateAll {
			// 更新所有技能
			entries, err := os.ReadDir(skillsDir)
			if err != nil {
				return fmt.Errorf("读取技能目录失败：%w", err)
			}

			for _, entry := range entries {
				if entry.IsDir() {
					skillsToUpdate = append(skillsToUpdate, entry.Name())
				}
			}

			if len(skillsToUpdate) == 0 {
				fmt.Println("暂无已安装的技能")
				return nil
			}
		} else {
			if len(args) == 0 {
				return fmt.Errorf("请指定要更新的技能名，或使用 --all 更新所有技能")
			}
			skillsToUpdate = []string{args[0]}
		}

		updated := 0
		failed := 0

		for _, skillName := range skillsToUpdate {
			fmt.Printf("\n🔍 检查技能：%s\n", skillName)

			// 获取本地版本
			localVersion, err := getLocalVersion(skillName)
			if err != nil {
				fmt.Printf("⚠️  跳过（无法获取本地版本）: %v\n", err)
				failed++
				continue
			}

			// 获取远程最新版本
			remoteSkill, err := client.GetSkill(skillName)
			if err != nil {
				fmt.Printf("⚠️  跳过（无法获取远程信息）: %v\n", err)
				failed++
				continue
			}

			remoteVersion := remoteSkill.LatestVersion

			// 比较版本
			if localVersion == remoteVersion {
				fmt.Printf("✅ 已是最新版本 (%s)\n", localVersion)
				continue
			}

			fmt.Printf("📦 发现新版本：%s -> %s\n", localVersion, remoteVersion)

			if updateDryRun {
				fmt.Println("   (dry-run, 跳过更新)")
				continue
			}

			// 下载并更新
			fmt.Println("⬇️  正在下载...")
			resp, err := client.InstallSkill(fmt.Sprintf("%d", remoteSkill.ID))
			if err != nil {
				fmt.Printf("❌ 下载失败：%v\n", err)
				failed++
				continue
			}

			// 下载文件
			zipPath := filepath.Join(skillsDir, fmt.Sprintf("%s-%s.zip", skillName, remoteVersion))
			if err := downloadFile(resp.DownloadURL, zipPath); err != nil {
				fmt.Printf("❌ 下载文件失败：%v\n", err)
				failed++
				continue
			}

			// 备份旧版本
			backupPath := skillName + ".bak"
			oldPath := filepath.Join(skillsDir, skillName)
			if err := os.Rename(oldPath, backupPath); err != nil {
				fmt.Printf("⚠️  备份失败：%v\n", err)
			}

			// 解压新版本
			fmt.Println("📦 正在解压...")
			if err := unzip(zipPath, oldPath); err != nil {
				// 恢复旧版本
				if _, err := os.Stat(backupPath); err == nil {
					os.RemoveAll(oldPath)
					os.Rename(backupPath, oldPath)
				}
				fmt.Printf("❌ 解压失败：%v\n", err)
				os.Remove(zipPath)
				failed++
				continue
			}

			// 清理
			os.Remove(zipPath)
			os.RemoveAll(backupPath)

			fmt.Printf("✅ 更新成功：%s -> %s\n", localVersion, remoteVersion)
			updated++
		}

		fmt.Printf("\n📊 更新完成：%d 个成功，%d 个失败\n", updated, failed)
		return nil
	},
}

func init() {
	updateCmdDetailed.Flags().BoolVarP(&updateAll, "all", "a", false, "更新所有技能")
	updateCmdDetailed.Flags().BoolVar(&updateDryRun, "dry-run", false, "仅检查，不实际更新")
}

// getLocalVersion 获取本地版本
func getLocalVersion(skillName string) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	manifestPath := filepath.Join(home, ".mcp", "skills", skillName, "mcp-manifest.json")
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		return "", err
	}

	var manifest struct {
		Version string `json:"version"`
	}

	if err := json.Unmarshal(data, &manifest); err != nil {
		return "", err
	}

	return manifest.Version, nil
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
		if f.FileInfo().IsDir() {
			continue
		}

		path := filepath.Join(dest, f.Name)
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return err
		}

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

// contains 字符串包含检查
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}
