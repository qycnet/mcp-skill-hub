package main

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

var (
	version   = "1.0.0"
	buildTime = "unknown"
	gitCommit = "unknown"
	goVersion = runtime.Version()
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "显示版本信息",
	Long:  `显示 MCP CLI 的版本、构建时间和 Git 提交信息。`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("MCP CLI Version %s\n", version)
		fmt.Printf("  Go Version: %s\n", goVersion)
		fmt.Printf("  OS/Arch:    %s/%s\n", runtime.GOOS, runtime.GOARCH)
		fmt.Printf("  Git Commit: %s\n", gitCommit)
		fmt.Printf("  Built:      %s\n", buildTime)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
