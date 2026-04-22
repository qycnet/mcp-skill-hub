package main

import (
	"fmt"
	"os"
	"path/filepath"
	"text/tabwriter"

	"github.com/mcp-skill-hub/cli/internal/config"
	"github.com/spf13/cobra"
)

var (
	listFormat string
	listAll    bool
)

var listCmdDetailed = &cobra.Command{
	Use:   "list",
	Short: "列出已安装的技能",
	Long:  `列出本地已安装的所有 MCP 技能，支持按状态筛选和格式化输出。`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// 获取技能安装目录
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("获取用户目录失败：%w", err)
		}

		skillsDir := filepath.Join(home, ".mcp", "skills")

		// 检查目录是否存在
		if _, err := os.Stat(skillsDir); os.IsNotExist(err) {
			fmt.Println("暂无已安装的技能")
			fmt.Println("使用 'mcp install <技能名>' 安装技能")
			return nil
		}

		// 读取技能目录
		entries, err := os.ReadDir(skillsDir)
		if err != nil {
			return fmt.Errorf("读取技能目录失败：%w", err)
		}

		if len(entries) == 0 {
			fmt.Println("暂无已安装的技能")
			return nil
		}

		// 加载配置
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("加载配置失败：%w", err)
		}

		// 获取远程技能信息（可选）
		var skills []SkillInfo
		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}

			skillPath := filepath.Join(skillsDir, entry.Name())
			info := SkillInfo{
				Name:      entry.Name(),
				InstallPath: skillPath,
				Status:    "installed",
			}

			// 尝试读取本地 manifest
			manifestPath := filepath.Join(skillPath, "mcp-manifest.json")
			if _, err := os.Stat(manifestPath); err == nil {
				info.HasManifest = true
			}

			skills = append(skills, info)
		}

		// 输出结果
		switch listFormat {
		case "table":
			printTable(skills)
		case "json":
			printJSON(skills)
		default:
			printTable(skills)
		}

		return nil
	},
}

func init() {
	listCmdDetailed.Flags().StringVarP(&listFormat, "format", "f", "table", "输出格式 (table, json)")
	listCmdDetailed.Flags().BoolVarP(&listAll, "all", "a", false, "显示所有状态（包括已禁用）")
}

// SkillInfo 技能信息
type SkillInfo struct {
	Name        string `json:"name"`
	Version     string `json:"version,omitempty"`
	InstallPath string `json:"install_path"`
	Status      string `json:"status"`
	HasManifest bool   `json:"has_manifest"`
}

// printTable 表格输出
func printTable(skills []SkillInfo) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "NAME\tVERSION\tSTATUS\tPATH")
	fmt.Fprintln(w, "----\t-------\t------\t----")

	for _, skill := range skills {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
			skill.Name,
			skill.Version,
			skill.Status,
			skill.InstallPath)
	}
	w.Flush()

	fmt.Printf("\n共 %d 个技能\n", len(skills))
}

// printJSON JSON 输出
func printJSON(skills []SkillInfo) {
	fmt.Println("[")
	for i, skill := range skills {
		comma := ","
		if i == len(skills)-1 {
			comma = ""
		}
		fmt.Printf("  {\"name\": \"%s\", \"version\": \"%s\", \"status\": \"%s\"}%s\n",
			skill.Name, skill.Version, skill.Status, comma)
	}
	fmt.Println("]")
}
