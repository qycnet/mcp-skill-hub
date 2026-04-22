package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/mcp-skill-hub/cli/internal/api"
	"github.com/mcp-skill-hub/cli/internal/config"
	"github.com/spf13/cobra"
)

var infoDetailed bool

var infoCmdDetailed = &cobra.Command{
	Use:   "info [技能名]",
	Short: "查看技能信息",
	Long:  `查看已安装技能的详细信息，包括版本、描述、依赖等。`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		skillName := args[0]

		// 首先检查本地安装
		localInfo, err := getLocalSkillInfo(skillName)
		if err == nil {
			if infoDetailed {
				printDetailedInfo(localInfo)
			} else {
				printBasicInfo(localInfo)
			}
			return nil
		}

		// 本地未安装，查询远程
		fmt.Printf("🔍 本地未安装，查询远程技能：%s\n", skillName)

		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("加载配置失败：%w", err)
		}

		client := api.NewClient(cfg.Server.URL, cfg.Server.Token)

		// 先搜索
		searchResult, err := client.SearchSkills(skillName, 1, 10)
		if err != nil {
			return fmt.Errorf("搜索失败：%w", err)
		}

		if len(searchResult.Skills) == 0 {
			return fmt.Errorf("未找到技能：%s", skillName)
		}

		// 获取第一个匹配的详情
		skill := searchResult.Skills[0]
		info := &SkillInfo{
			Name:        skill.Name,
			Version:     "", // 需要进一步获取
			Description: skill.Description,
			Category:    skill.Category,
			Author:      "",
			License:     "",
			Homepage:    "",
			Downloads:   skill.Downloads,
			Rating:      skill.Rating,
			Status:      "available",
		}

		printBasicInfo(info)
		fmt.Printf("\n💡 提示：使用 'mcp install %s' 安装此技能\n", skillName)

		return nil
	},
}

func init() {
	infoCmdDetailed.Flags().BoolVarP(&infoDetailed, "detailed", "d", false, "显示详细信息")
}

// getLocalSkillInfo 获取本地技能信息
func getLocalSkillInfo(skillName string) (*SkillInfo, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	skillPath := filepath.Join(home, ".mcp", "skills", skillName)

	// 检查目录是否存在
	if _, err := os.Stat(skillPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("技能未安装")
	}

	// 读取 manifest
	manifestPath := filepath.Join(skillPath, "mcp-manifest.json")
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		// 没有 manifest，返回基本信息
		return &SkillInfo{
			Name:   skillName,
			Status: "installed (no manifest)",
		}, nil
	}

	var manifest struct {
		Name        string   `json:"name"`
		Version     string   `json:"version"`
		DisplayName string   `json:"display_name"`
		Description string   `json:"description"`
		Category    string   `json:"category"`
		Author      string   `json:"author"`
		License     string   `json:"license"`
		Homepage    string   `json:"homepage"`
		Tags        []string `json:"tags"`
	}

	if err := json.Unmarshal(data, &manifest); err != nil {
		return nil, fmt.Errorf("解析 manifest 失败：%w", err)
	}

	return &SkillInfo{
		Name:        manifest.Name,
		Version:     manifest.Version,
		Description: manifest.Description,
		Category:    manifest.Category,
		Author:      manifest.Author,
		License:     manifest.License,
		Homepage:    manifest.Homepage,
		Tags:        manifest.Tags,
		Status:      "installed",
	}, nil
}

// printBasicInfo 打印基本信息
func printBasicInfo(info *SkillInfo) {
	fmt.Printf("\n📦 %s\n", info.Name)
	if info.Version != "" {
		fmt.Printf("   版本：%s\n", info.Version)
	}
	if info.Description != "" {
		fmt.Printf("   描述：%s\n", info.Description)
	}
	if info.Category != "" {
		fmt.Printf("   分类：%s\n", info.Category)
	}
	if info.Author != "" {
		fmt.Printf("   作者：%s\n", info.Author)
	}
	if info.License != "" {
		fmt.Printf("   许可：%s\n", info.License)
	}
	if info.Downloads > 0 {
		fmt.Printf("   下载：%d\n", info.Downloads)
	}
	if info.Rating > 0 {
		fmt.Printf("   评分：%.1f\n", info.Rating)
	}
	fmt.Printf("   状态：%s\n", info.Status)
}

// printDetailedInfo 打印详细信息
func printDetailedInfo(info *SkillInfo) {
	fmt.Printf("\n📦 技能详情：%s\n", info.Name)
	fmt.Println("=" + string(make([]byte, 50)))

	if info.Version != "" {
		fmt.Printf("版本号：%s\n", info.Version)
	}
	if info.Description != "" {
		fmt.Printf("描述：%s\n", info.Description)
	}
	if info.Category != "" {
		fmt.Printf("分类：%s\n", info.Category)
	}
	if info.Author != "" {
		fmt.Printf("作者：%s\n", info.Author)
	}
	if info.License != "" {
		fmt.Printf("许可证：%s\n", info.License)
	}
	if info.Homepage != "" {
		fmt.Printf("主页：%s\n", info.Homepage)
	}
	if len(info.Tags) > 0 {
		fmt.Printf("标签：%s\n", joinStrings(info.Tags, ", "))
	}
	if info.Downloads > 0 {
		fmt.Printf("下载量：%d\n", info.Downloads)
	}
	if info.Rating > 0 {
		fmt.Printf("评分：%.1f\n", info.Rating)
	}
	fmt.Printf("状态：%s\n", info.Status)
	fmt.Printf("安装路径：%s\n", info.InstallPath)
}

// joinStrings 连接字符串
func joinStrings(strs []string, sep string) string {
	if len(strs) == 0 {
		return ""
	}
	result := strs[0]
	for _, s := range strs[1:] {
		result += sep + s
	}
	return result
}
