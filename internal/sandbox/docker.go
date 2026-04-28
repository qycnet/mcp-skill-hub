package sandbox

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Sandbox 沙箱执行环境
type Sandbox struct {
	dockerPath   string
	networkName  string
	timeout      time.Duration
	maxMemory    string // 例如 "512m"
	maxCPU       string // 例如 "0.5"
	readOnlyRoot bool
}

// SandboxConfig 沙箱配置
type SandboxConfig struct {
	DockerPath   string
	NetworkName  string
	Timeout      time.Duration
	MaxMemory    string
	MaxCPU       string
	ReadOnlyRoot bool
}

// ExecutionResult 执行结果
type ExecutionResult struct {
	ID         string        `json:"id"`
	Status     string        `json:"status"` // success, failed, timeout
	ExitCode   int           `json:"exit_code"`
	Stdout     string        `json:"stdout"`
	Stderr     string        `json:"stderr"`
	Duration   time.Duration `json:"duration"`
	ContainerID string       `json:"container_id"`
}

// NewSandbox 创建沙箱环境
func NewSandbox(config SandboxConfig) *Sandbox {
	if config.DockerPath == "" {
		config.DockerPath = "docker"
	}
	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}
	if config.MaxMemory == "" {
		config.MaxMemory = "512m"
	}
	if config.MaxCPU == "" {
		config.MaxCPU = "1.0"
	}

	return &Sandbox{
		dockerPath:   config.DockerPath,
		networkName:  config.NetworkName,
		timeout:      config.Timeout,
		maxMemory:    config.MaxMemory,
		maxCPU:       config.MaxCPU,
		readOnlyRoot: config.ReadOnlyRoot,
	}
}

// ExecuteInSandbox 在沙箱中执行技能
func (s *Sandbox) ExecuteInSandbox(ctx context.Context, skillPath, entrypoint string, args []string, input io.Reader) (*ExecutionResult, error) {
	executionID := uuid.New().String()
	containerName := fmt.Sprintf("skill-%s", executionID[:8])

	// 构建安全的 Docker 运行命令
	dockerArgs := []string{
		"run",
		"--rm", // 执行后自动删除
		"--name", containerName,
		
		// 资源限制
		"--memory", s.maxMemory,
		"--memory-swap", s.maxMemory,
		"--cpus", s.maxCPU,
		"--pids-limit", "100",
		
		// 安全选项
		"--security-opt", "no-new-privileges:true",
		"--security-opt", "seccomp=unconfined", // 可根据需要配置 seccomp profile
		
		// 网络隔离
		"--network", s.getNetworkOption(),
		
		// 文件系统限制
		"--read-only", // 只读根文件系统
		"--tmpfs", "/tmp:size=100m,mode=1777",
		
		// 用户权限
		"--user", "1000:1000", // 非 root 用户
		
		// 环境变量
		"-e", "SKILL_EXECUTION_ID=" + executionID,
		"-e", "SKILL_TIMEOUT=" + s.timeout.String(),
		
		// 挂载技能包（只读）
		"-v", fmt.Sprintf("%s:/skill:ro", skillPath),
		
		// 工作目录
		"-w", "/skill",
	}

	// 如果需要可写目录
	if !s.readOnlyRoot {
		dockerArgs = append(dockerArgs, 
			"--tmpfs", "/var:size=50m",
			"--tmpfs", "/run:size=10m",
		)
	}

	// 添加入口点和参数
	dockerArgs = append(dockerArgs, entrypoint)
	dockerArgs = append(dockerArgs, args...)

	// 创建带超时的上下文
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	// 执行命令
	cmd := exec.CommandContext(ctx, s.dockerPath, dockerArgs...)
	
	if input != nil {
		cmd.Stdin = input
	}

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("创建 stdout 管道失败：%w", err)
	}

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("创建 stderr 管道失败：%w", err)
	}

	start := time.Now()
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("启动容器失败：%w", err)
	}

	// 读取输出
	stdoutBytes, _ := io.ReadAll(stdoutPipe)
	stderrBytes, _ := io.ReadAll(stderrPipe)

	err = cmd.Wait()
	duration := time.Since(start)

	result := &ExecutionResult{
		ID:         executionID,
		Duration:   duration,
		Stdout:     string(stdoutBytes),
		Stderr:     string(stderrBytes),
		ContainerID: containerName,
	}

	if ctx.Err() == context.DeadlineExceeded {
		result.Status = "timeout"
		result.ExitCode = -1
		result.Stderr = "执行超时"
		// 强制停止容器
		exec.Command(s.dockerPath, "stop", "-t", "2", containerName).Run()
		return result, nil
	}

	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			result.Status = "failed"
			result.ExitCode = exitErr.ExitCode()
		} else {
			return nil, fmt.Errorf("执行失败：%w", err)
		}
	} else {
		result.Status = "success"
		result.ExitCode = 0
	}

	return result, nil
}

// ExecuteSkillWithImage 使用指定镜像执行技能
func (s *Sandbox) ExecuteSkillWithImage(ctx context.Context, image, skillPath, entrypoint string, args []string) (*ExecutionResult, error) {
	executionID := uuid.New().String()
	containerName := fmt.Sprintf("skill-%s", executionID[:8])

	dockerArgs := []string{
		"run",
		"--rm",
		"--name", containerName,
		"--memory", s.maxMemory,
		"--cpus", s.maxCPU,
		"--network", s.getNetworkOption(),
		"--security-opt", "no-new-privileges:true",
		"--user", "1000:1000",
		"-v", fmt.Sprintf("%s:/skill:ro", skillPath),
		"-w", "/skill",
		image,
		entrypoint,
	}
	dockerArgs = append(dockerArgs, args...)

	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	start := time.Now()
	output, err := exec.CommandContext(ctx, s.dockerPath, dockerArgs...).CombinedOutput()
	duration := time.Since(start)

	result := &ExecutionResult{
		ID:         executionID,
		Duration:   duration,
		Stdout:     string(output),
		ContainerID: containerName,
	}

	if ctx.Err() == context.DeadlineExceeded {
		result.Status = "timeout"
		result.ExitCode = -1
		exec.Command(s.dockerPath, "stop", "-t", "2", containerName).Run()
		return result, nil
	}

	if err != nil {
		result.Status = "failed"
		if exitErr, ok := err.(*exec.ExitError); ok {
			result.ExitCode = exitErr.ExitCode()
		}
	} else {
		result.Status = "success"
		result.ExitCode = 0
	}

	return result, nil
}

// getNetworkOption 获取网络配置
func (s *Sandbox) getNetworkOption() string {
	if s.networkName != "" {
		return s.networkName
	}
	return "none" // 默认无网络访问
}

// BuildSkillImage 为技能构建隔离镜像
func (s *Sandbox) BuildSkillImage(ctx context.Context, skillPath, imageName string) error {
	// 创建临时 Dockerfile
	dockerfile := `FROM alpine:3.19
RUN apk add --no-cache ca-certificates tzdata
RUN addgroup -g 1000 skill && adduser -u 1000 -G skill -D skill
WORKDIR /skill
COPY . /skill/
RUN chown -R skill:skill /skill
USER skill
`

	tmpDir, err := os.MkdirTemp("", "skill-build-*")
	if err != nil {
		return fmt.Errorf("创建临时目录失败：%w", err)
	}
	defer os.RemoveAll(tmpDir)

	// 写入 Dockerfile
	if err := os.WriteFile(filepath.Join(tmpDir, "Dockerfile"), []byte(dockerfile), 0644); err != nil {
		return fmt.Errorf("写入 Dockerfile 失败：%w", err)
	}

	// 复制技能文件
	if err := copyDir(skillPath, filepath.Join(tmpDir, "skill")); err != nil {
		return fmt.Errorf("复制技能文件失败：%w", err)
	}

	// 构建镜像
	cmd := exec.CommandContext(ctx, s.dockerPath, "build", "-t", imageName, tmpDir)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("构建镜像失败：%w, output: %s", err, string(output))
	}

	return nil
}

// Cleanup 清理资源
func (s *Sandbox) Cleanup(ctx context.Context) error {
	// 清理所有技能容器
	cmd := exec.CommandContext(ctx, s.dockerPath, 
		"container", "prune", "-f", "--filter", "label=skill=true")
	return cmd.Run()
}

// copyDir 复制目录
func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		dstPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(dstPath, info.Mode())
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(dstPath, data, info.Mode())
	})
}

// ValidateSkillPackage 验证技能包安全性
func (s *Sandbox) ValidateSkillPackage(skillPath string) ([]string, error) {
	var warnings []string

	// 检查可疑文件
	suspiciousFiles := []string{
		".ssh", ".gnupg", ".aws",
		"id_rsa", "id_ed25519",
		".env", ".env.local",
	}

	err := filepath.Walk(skillPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		name := info.Name()
		for _, suspicious := range suspiciousFiles {
			if strings.Contains(name, suspicious) {
				warnings = append(warnings, fmt.Sprintf("发现可疑文件：%s", path))
			}
		}

		// 检查文件权限
		if info.Mode().Perm()&0077 != 0 {
			warnings = append(warnings, fmt.Sprintf("文件权限过于宽松：%s (%o)", path, info.Mode().Perm()))
		}

		return nil
	})

	return warnings, err
}
