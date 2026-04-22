package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/mcp-skill-hub/cli/internal/config"
	"github.com/spf13/cobra"
)

var (
	publishManifest string
	publishZip      string
)

var publishCmdDetailed = &cobra.Command{
	Use:   "publish [技能路径]",
	Short: "发布技能",
	Long:  `发布 MCP 技能到 Skill Hub，支持从目录或 ZIP 文件发布。`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		skillPath := args[0]

		// 加载配置
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("加载配置失败：%w", err)
		}

		// 检查 token
		if cfg.Server.Token == "" {
			return fmt.Errorf("未登录，请先使用 'mcp login' 登录")
		}

		// 确定发布方式
		var zipPath string
		if publishZip != "" {
			// 使用指定的 ZIP 文件
			zipPath = publishZip
		} else {
			// 从目录创建 ZIP
			zipPath, err = createZipFromDir(skillPath)
			if err != nil {
				return fmt.Errorf("创建 ZIP 失败：%w", err)
			}
			defer os.Remove(zipPath) // 清理临时文件
		}

		// 验证 manifest
		manifest, err := readManifest(skillPath)
		if err != nil {
			return fmt.Errorf("读取 manifest 失败：%w", err)
		}

		fmt.Printf("📦 准备发布技能：%s v%s\n", manifest.Name, manifest.Version)
		fmt.Printf("   显示名称：%s\n", manifest.DisplayName)
		fmt.Printf("   分类：%s\n", manifest.Category)

		// 上传
		fmt.Println("\n⬆️  正在上传...")
		uploadedSkill, err := uploadSkill(cfg.Server.URL, cfg.Server.Token, zipPath, manifest)
		if err != nil {
			return fmt.Errorf("上传失败：%w", err)
		}

		fmt.Printf("\n✅ 发布成功！\n")
		fmt.Printf("   技能 ID: %d\n", uploadedSkill.ID)
		fmt.Printf("   UUID: %s\n", uploadedSkill.UUID)
		fmt.Printf("   状态：%s (待审核)\n", uploadedSkill.Status)
		fmt.Printf("\n💡 提示：技能审核通过后将上线到技能市场\n")

		return nil
	},
}

func init() {
	publishCmdDetailed.Flags().StringVarP(&publishManifest, "manifest", "m", "", "manifest 文件路径（默认从技能路径读取）")
	publishCmdDetailed.Flags().StringVarP(&publishZip, "zip", "z", "", "ZIP 文件路径（如指定则直接上传，不从目录创建）")
}

// createZipFromDir 从目录创建 ZIP
func createZipFromDir(srcDir string) (string, error) {
	// 创建临时文件
	tmpFile, err := os.CreateTemp("", "mcp-skill-*.zip")
	if err != nil {
		return "", err
	}
	defer tmpFile.Close()

	// 创建 ZIP 写入器
	zipWriter := zip.NewWriter(tmpFile)
	defer zipWriter.Close()

	// 遍历目录
	err = filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 跳过目录
		if info.IsDir() {
			return nil
		}

		// 跳过某些文件
		if shouldSkip(path) {
			return nil
		}

		// 计算相对路径
		relPath, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}

		// 添加到 ZIP
		return addToZip(zipWriter, path, relPath)
	})

	if err != nil {
		return "", err
	}

	return tmpFile.Name(), nil
}

// shouldSkip 判断是否跳过文件
func shouldSkip(path string) bool {
	skipPatterns := []string{
		".git",
		"node_modules",
		"__pycache__",
		".DS_Store",
		"*.log",
		"*.tmp",
	}

	base := filepath.Base(path)
	for _, pattern := range skipPatterns {
		if matchPattern(base, pattern) {
			return true
		}
	}
	return false
}

// matchPattern 简单模式匹配
func matchPattern(s, pattern string) bool {
	if pattern[0] == '*' && pattern[len(pattern)-1] == '*' {
		return contains(s, pattern[1:len(pattern)-1])
	}
	if pattern[0] == '.' {
		return s == pattern
	}
	return s == pattern
}

// addToZip 添加文件到 ZIP
func addToZip(zipWriter *zip.Writer, srcPath, relPath string) error {
	// 创建 ZIP 条目
	header, err := zip.FileInfoHeader(os.FileInfo(nil))
	if err != nil {
		return err
	}
	header.Name = relPath
	header.Method = zip.Deflate

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}

	// 读取源文件
	file, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// 复制内容
	_, err = io.Copy(writer, file)
	return err
}

// readManifest 读取 manifest
func readManifest(skillPath string) (*Manifest, error) {
	manifestPath := publishManifest
	if manifestPath == "" {
		manifestPath = filepath.Join(skillPath, "mcp-manifest.json")
	}

	data, err := os.ReadFile(manifestPath)
	if err != nil {
		return nil, err
	}

	var manifest Manifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return nil, err
	}

	return &manifest, nil
}

// uploadSkill 上传技能
func uploadSkill(baseURL, token, zipPath string, manifest *Manifest) (*UploadedSkill, error) {
	// 创建 multipart 表单
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// 添加文件
	file, err := os.Open(zipPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	part, err := writer.CreateFormFile("package", filepath.Base(zipPath))
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return nil, err
	}

	// 关闭 writer
	if err := writer.Close(); err != nil {
		return nil, err
	}

	// 创建请求
	url := fmt.Sprintf("%s/api/v1/skills/upload", baseURL)
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("上传失败 (状态码 %d): %s", resp.StatusCode, string(respBody))
	}

	var result UploadedSkill
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// Manifest 技能清单
type Manifest struct {
	Name        string   `json:"name"`
	DisplayName string   `json:"display_name"`
	Version     string   `json:"version"`
	Description string   `json:"description"`
	Category    string   `json:"category"`
	Tags        []string `json:"tags"`
	Author      string   `json:"author"`
	License     string   `json:"license"`
	Homepage    string   `json:"homepage"`
	Repository  string   `json:"repository"`
	MCPVersion  string   `json:"mcp_version"`
}

// UploadedSkill 上传结果
type UploadedSkill struct {
	ID       uint   `json:"id"`
	UUID     string `json:"uuid"`
	Name     string `json:"name"`
	Status   string `json:"status"`
	Message  string `json:"message"`
}

// contains 字符串包含
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && findSubstring(s, substr))
}

// findSubstring 查找子串
func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
