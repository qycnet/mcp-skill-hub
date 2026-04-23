# 多阶段构建 - 生产镜像
FROM golang:1.21-alpine AS builder

# 安装必要工具
RUN apk add --no-cache git ca-certificates

# 设置工作目录
WORKDIR /build

# 复制 go.mod 和 go.sum
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 编译（静态链接，支持 musl）
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-w -s -X main.version=0.1.0 -X main.buildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
    -o /build/mcp-skill-hub \
    ./cmd/server/main.go

# 生产镜像
FROM alpine:3.19

# 安装运行时依赖
RUN apk add --no-cache ca-certificates tzdata curl

# 创建非 root 用户
RUN addgroup -g 1000 appgroup && \
    adduser -u 1000 -G appgroup -D appuser

# 设置工作目录
WORKDIR /app

# 从 builder 复制二进制文件
COPY --from=builder /build/mcp-skill-hub /app/mcp-skill-hub
COPY --from=builder /build/config.yaml.example /app/config.yaml.example

# 设置权限
RUN chown -R appuser:appgroup /app && \
    chmod +x /app/mcp-skill-hub

# 切换到非 root 用户
USER appuser

# 暴露端口
EXPOSE 8080

# 健康检查
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8080/health || exit 1

# 启动命令
ENTRYPOINT ["/app/mcp-skill-hub"]
