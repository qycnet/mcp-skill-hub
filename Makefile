.PHONY: help dev build test run docker-build docker-up docker-down clean security-scan coverage

# 初始化环境
init:
	@echo "🔐 初始化 MCP Skill Hub 环境..."
	@chmod +x scripts/init-env.sh
	@./scripts/init-env.sh
	@echo ""
	@echo "✅ 初始化完成！"
	@echo "🚀 运行 make docker-up 启动服务"

# 默认目标
help:
	@echo "MCP Skill Hub - Makefile 命令"
	@echo ""
	@echo "开发命令:"
	@echo "  make dev          - 启动开发环境（热重载）"
	@echo "  make run          - 运行服务器"
	@echo "  make test         - 运行测试"
	@echo "  make lint         - 代码检查"
	@echo ""
	@echo "测试命令:"
	@echo "  make test         - 运行单元测试"
	@echo "  make test-integration - 运行集成测试"
	@echo "  make coverage     - 运行测试并生成覆盖率报告"
	@echo "  make coverage-html - 生成 HTML 覆盖率报告"
	@echo "  make coverage-check - 检查覆盖率是否达标"
	@echo ""
	@echo "安全命令:"
	@echo "  make security-scan     - 运行完整安全扫描"
	@echo "  make install-security-tools - 安装安全扫描工具"
	@echo "  make govulncheck       - 仅运行 govulncheck"
	@echo "  make osv-scan          - 仅运行 OSV-Scanner"
	@echo "  make security-audit    - 运行全面安全审计"
	@echo ""
	@echo "构建命令:"
	@echo "  make build        - 编译二进制文件"
	@echo "  make docker-build - 构建 Docker 镜像"
	@echo ""
	@echo "Docker 命令:"
	@echo "  make docker-up    - 启动 Docker Compose"
	@echo "  make docker-down  - 停止 Docker Compose"
	@echo "  make docker-logs  - 查看日志"
	@echo ""
	@echo "清理命令:"
	@echo "  make clean        - 清理构建产物"

# 开发环境
dev:
	@echo "🚀 启动开发环境..."
	go run ./cmd/server/main.go

# 运行
run:
	go run ./cmd/server/main.go

# 测试
test:
	@echo "🧪 运行测试..."
	go test -v -race -cover ./...

# 代码检查
lint:
	@echo "🔍 代码检查..."
	golangci-lint run

# 构建
build:
	@echo "🔨 编译二进制文件..."
	CGO_ENABLED=0 go build -ldflags="-w -s -X main.version=0.1.0" -o bin/mcp-skill-hub ./cmd/server/main.go
	@echo "✅ 构建完成：bin/mcp-skill-hub"

# Docker 构建
docker-build:
	@echo "🐳 构建 Docker 镜像..."
	docker build -t mcp-skill-hub:latest .

# Docker 启动
docker-up:
	@echo "🚀 启动 Docker Compose..."
	docker-compose up -d
	@echo "✅ 服务已启动"
	@echo "📡 访问地址：http://localhost:8080"
	@echo "🗄️  PostgreSQL: localhost:5432"
	@echo "💾 MinIO Console: http://localhost:9001"

# Docker 停止
docker-down:
	@echo "🛑 停止 Docker Compose..."
	docker-compose down

# Docker 日志
docker-logs:
	docker-compose logs -f app

# 清理
clean:
	@echo "🧹 清理构建产物..."
	rm -rf bin/
	rm -rf dist/
	go clean -cache

# 数据库迁移
migrate:
	@echo "📦 运行数据库迁移..."
	go run ./cmd/migrate/main.go

# 生成 API 文档
docs:
	@echo "📚 生成 API 文档..."
	swag init -g cmd/server/main.go -o docs/api

# 格式化代码
fmt:
	@echo "✨ 格式化代码..."
	go fmt ./...
	goimports -w .

# 检查依赖
deps:
	@echo "📦 检查依赖..."
	go mod tidy
	go mod verify

# 安全扫描
security-scan:
	@echo "🔒 安全扫描..."
	govulncheck ./...

# 性能分析
profile:
	@echo "📊 性能分析..."
	go test -cpuprofile=cpu.prof -memprofile=mem.prof ./...
