# MCP CLI

🚀 MCP Skill Hub 命令行工具 - 管理和使用 MCP 技能

## 安装

### 方式一：Go 安装

```bash
go install github.com/mcp-skill-hub/cli@latest
```

### 方式二：下载二进制文件

```bash
# Linux
curl -LO https://github.com/qycnet/mcp-skill-hub/releases/latest/download/mcp-linux-amd64
chmod +x mcp-linux-amd64
sudo mv mcp-linux-amd64 /usr/local/bin/mcp

# macOS
curl -LO https://github.com/qycnet/mcp-skill-hub/releases/latest/download/mcp-darwin-amd64
chmod +x mcp-darwin-amd64
sudo mv mcp-darwin-amd64 /usr/local/bin/mcp
```

### 方式三：从源码构建

```bash
git clone https://github.com/qycnet/mcp-skill-hub.git
cd mcp-skill-hub/cli
go build -o mcp ./cmd/mcp
sudo mv mcp /usr/local/bin/
```

## 快速开始

```bash
# 1. 登录
mcp login

# 2. 搜索技能
mcp search "code analysis"

# 3. 安装技能
mcp install claude-context

# 4. 查看已安装
mcp list
```

## 命令参考

### 认证
```bash
mcp login
mcp login -u username -p password
```

### 技能管理
```bash
# 搜索
mcp search "关键词"
mcp search "ai" -c developer-tools --sort rating

# 安装
mcp install claude-context
mcp install claude-context@1.2.0  # 指定版本

# 查看详情
mcp info claude-context
mcp info claude-context -d  # 详细信息

# 列出已安装
mcp list
mcp list -f json  # JSON 格式

# 更新
mcp update claude-context
mcp update --all  # 更新所有

# 卸载
mcp uninstall claude-context
mcp uninstall claude-context -f  # 强制卸载
```

### 发布技能
```bash
# 从目录发布
mcp publish ./my-skill

# 从 ZIP 发布
mcp publish ./my-skill -z my-skill.zip
```

## 配置文件

配置文件位于 `~/.mcp/config.yaml`

```yaml
server:
  url: http://localhost:8080
  token: your-api-token

proxy:
  http: http://proxy.example.com:8080
  https: http://proxy.example.com:8080
```

## 环境变量

| 变量 | 描述 | 默认值 |
|------|------|--------|
| `MCP_CONFIG` | 配置文件路径 | `~/.mcp/config.yaml` |
| `MCP_SERVER` | 服务器 URL | `http://localhost:8080` |
| `MCP_TOKEN` | API Token | - |
| `HTTP_PROXY` | HTTP 代理 | - |
| `MCP_DEBUG` | 调试模式 | `false` |

## 示例

### 搜索并安装 AI 工具
```bash
mcp search "ai" -c developer-tools
mcp install ai-code-assistant
mcp info ai-code-assistant
```

### 发布自己的技能
```bash
# 准备技能目录
mkdir my-skill && cd my-skill

# 创建 manifest
cat > mcp-manifest.json << EOF
{
  "name": "my-awesome-skill",
  "display_name": "My Awesome Skill",
  "version": "1.0.0",
  "description": "这是一个很棒的技能",
  "category": "developer-tools",
  "tags": ["ai", "productivity"],
  "author": "Your Name",
  "license": "MIT"
}
EOF

# 发布
mcp publish .
```

## 故障排查

**Q: 登录失败**
```bash
curl -I https://mcp-skill-hub.dev
cat ~/.mcp/config.yaml
mcp logout && mcp login
```

**Q: 安装失败**
```bash
df -h ~/.mcp
rm -rf ~/.mcp/cache/*
mcp install <skill> -v
```

## 许可证

MIT License
