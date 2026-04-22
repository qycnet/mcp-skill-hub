package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/mcp-skill-hub/cli/internal/api"
	"github.com/mcp-skill-hub/cli/internal/config"
	"github.com/spf13/cobra"
)

var (
	searchPage     int
	searchPageSize int
	searchCategory string
	searchSort     string
)

var searchCmdDetailed = &cobra.Command{
	Use:   "search [关键词]",
	Short: "搜索技能",
	Long:  `搜索 MCP Skill Hub 中的技能，支持分类过滤和排序。`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// 加载配置
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("加载配置失败：%w", err)
		}

		// 创建客户端
		client := api.NewClient(cfg.Server.URL, cfg.Server.Token)

		// 搜索
		query := args[0]
		fmt.Printf("🔍 搜索技能：%s\n", query)
		if searchCategory != "" {
			fmt.Printf("📁 分类：%s\n", searchCategory)
		}

		result, err := client.SearchSkills(query, searchPage, searchPageSize)
		if err != nil {
			return fmt.Errorf("搜索失败：%w", err)
		}

		if result.Total == 0 {
			fmt.Println("❌ 未找到相关技能")
			return nil
		}

		// 显示结果
		fmt.Printf("\n✅ 找到 %d 个技能 (第 %d 页，共 %d 页)\n\n",
			result.Total, result.Page, (result.Total+int64(result.PageSize)-1)/int64(result.PageSize))

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "NAME\tDISPLAY NAME\tCATEGORY\t⭐ RATING\t⬇️ DOWNLOADS\tSCORE")
		fmt.Fprintln(w, "----\t------------\t--------\t--------\t---------\t-----")

		for _, skill := range result.Skills {
			fmt.Fprintf(w, "%s\t%s\t%s\t%.1f\t%d\t%.0f\n",
				skill.Name,
				skill.DisplayName,
				skill.Category,
				skill.Rating,
				skill.Downloads,
				skill.QualityScore)
		}
		w.Flush()

		fmt.Printf("\n💡 提示：使用 'mcp info <技能名>' 查看详情\n")

		return nil
	},
}

func init() {
	searchCmdDetailed.Flags().IntVarP(&searchPage, "page", "p", 1, "页码")
	searchCmdDetailed.Flags().IntVar(&searchPageSize, "page-size", 20, "每页数量")
	searchCmdDetailed.Flags().StringVarP(&searchCategory, "category", "c", "", "分类过滤")
	searchCmdDetailed.Flags().StringVar(&searchSort, "sort", "downloads", "排序方式 (downloads, rating, quality_score)")
}
