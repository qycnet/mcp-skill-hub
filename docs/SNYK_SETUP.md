# Snyk 配置指南

**版本**: v0.4.0  
**最后更新**: 2026-04-22

---

## 📋 目录

1. [Snyk 简介](#snyk 简介)
2. [安装 Snyk CLI](#安装-snyk-cli)
3. [认证配置](#认证配置)
4. [本地使用](#本地使用)
5. [CI/CD 集成](#cicd-集成)
6. [故障排查](#故障排查)

---

## Snyk 简介

[Snyk](https://snyk.io/) 是一个开发者安全平台，提供：

- 🔍 **漏洞扫描** - 检测依赖中的已知漏洞
- 📊 **许可证合规** - 检查开源许可证风险
- 🐳 **容器扫描** - 扫描 Docker 镜像漏洞
- 🏗️ **IaC 扫描** - 检查基础设施配置安全
- 📡 **持续监控** - 监控新披露的漏洞

### 为什么使用 Snyk

| 功能 | OSV-Scanner | Snyk |
|------|-------------|------|
| 漏洞数据库 | 开源 | 商业 + 开源 |
| 更新频率 | 每日 | 实时 |
| 修复建议 | 基础 | 详细 + 自动化 |
| CI/CD 集成 | 简单 | 完整 |
| 监控告警 | ❌ | ✅ |
| 容器扫描 | ❌ | ✅ |
| IaC 扫描 | ❌ | ✅ |

---

## 安装 Snyk CLI

### 方式一：npm（推荐）

```bash
npm install -g snyk
```

### 方式二：Homebrew（macOS）

```bash
brew install snyk/tap/snyk
```

### 方式三：二进制下载

```bash
# Linux
curl -Lo snyk https://github.com/snyk/cli/releases/latest/download/snyk-linux
chmod +x snyk
sudo mv snyk /usr/local/bin/

# macOS
curl -Lo snyk https://github.com/snyk/cli/releases/latest/download/snyk-macos
chmod +x snyk
sudo mv snyk /usr/local/bin/

# Windows
curl -Lo snyk.exe https://github.com/snyk/cli/releases/latest/download/snyk-win.exe
move snyk.exe C:\Windows\System32\snyk.exe
```

### 验证安装

```bash
snyk --version
```

---

## 认证配置

### 快速认证

运行认证脚本：

```bash
./scripts/snyk-auth.sh
```

### 手动认证

#### 方式一：浏览器认证

```bash
snyk auth
```

这会打开浏览器，登录 Snyk 后自动完成认证。

#### 方式二：API Token

1. 登录 [Snyk Console](https://app.snyk.io/)
2. 进入 **Account Settings** → **General**
3. 复制 **API Token**
4. 运行认证：

```bash
echo "YOUR_API_TOKEN" | snyk auth
```

### 配置组织

```bash
# 查看可用组织
snyk orgs

# 设置默认组织
snyk org <org-id>
```

### 环境变量（可选）

```bash
# ~/.bashrc 或 ~/.zshrc
export SNYK_TOKEN="your-api-token"
export SNYK_ORG="your-org-id"
```

---

## 本地使用

### 运行安全扫描

```bash
# 使用脚本（推荐）
./scripts/snyk-test.sh

# 或直接使用 Makefile
make security-scan
```

### 常用命令

```bash
# 测试漏洞
snyk test

# 仅高严重性
snyk test --severity-threshold=high

# JSON 输出
snyk test --json > results.json

# 监控项目
snyk monitor

# 查看组织
snyk org

# 重新认证
snyk auth
```

### 扫描 Docker 镜像

```bash
# 构建镜像
docker build -t mcp-skill-hub:latest .

# 扫描镜像
snyk container test mcp-skill-hub:latest --file=Dockerfile
```

### 扫描基础设施配置

```bash
# 扫描 Docker Compose
snyk iac test docker-compose.yml

# 扫描 Kubernetes 配置
snyk iac test k8s/
```

---

## CI/CD 集成

### GitHub Actions

Snyk 已集成到 GitHub Actions：

- **推送扫描**: 每次 push 自动运行
- **定时扫描**: 每天 UTC 2:00
- **PR 检查**: Pull Request 自动扫描
- **监控**: main 分支自动监控

### 配置仓库密钥

1. 进入 GitHub 仓库 **Settings** → **Secrets and variables** → **Actions**
2. 添加以下密钥：

| 密钥名 | 值 | 说明 |
|--------|-----|------|
| `SNYK_TOKEN` | 你的 API Token | Snyk 认证 |
| `SNYK_ORG_ID` | 组织 ID | 可选，多组织时使用 |

### 工作流文件

`.github/workflows/snyk.yml` 已配置：

- `snyk-test`: 依赖漏洞扫描
- `snyk-container`: Docker 镜像扫描
- `snyk-iac`: 基础设施配置扫描
- `security-summary`: 扫描摘要

### 忽略误报

创建 `.snyk` 文件：

```
# 格式：vuln-id | reason | expires
GO-2023-XXXX | 误报，实际不受影响 | 2024-12-31
```

---

## 故障排查

### Q: 认证失败

```bash
# 检查网络
curl -I https://snyk.io

# 清除旧认证
rm -rf ~/.config/configstore/snyk.json

# 重新认证
snyk auth
```

### Q: 扫描超时

```bash
# 增加超时
snyk test --timeout=600

# 仅扫描直接依赖
snyk test --strict-out-of-sync=false
```

### Q: 发现漏洞

1. **审查漏洞**
   ```bash
   snyk test --json | jq '.vulnerabilities[] | {id, severity, packageName}'
   ```

2. **查看修复建议**
   ```bash
   snyk wizard
   ```

3. **更新依赖**
   ```bash
   go get -u github.com/problematic/package
   ```

4. **重新扫描验证**
   ```bash
   snyk test
   ```

### Q: CI/CD 失败

检查以下几点：

1. **密钥配置**
   ```bash
   # 在 GitHub Actions 日志中查看
   echo $SNYK_TOKEN | head -c 10
   ```

2. **权限问题**
   - 确保 Token 有效
   - 确保组织权限正确

3. **网络问题**
   - 检查 GitHub Actions runner 网络
   - 考虑使用自托管 runner

---

## 最佳实践

### ✅ 推荐

- 定期运行安全扫描（至少每周）
- 在 CI/CD 中集成 Snyk
- 设置邮件告警
- 及时修复高严重性漏洞
- 使用 `snyk monitor` 持续监控

### ❌ 避免

- 忽略所有警告
- 在生产环境使用未扫描的镜像
- 将 Token 提交到代码库
- 跳过 PR 安全检查

---

## 定价和计划

| 计划 | 价格 | 功能 |
|------|------|------|
| **Free** | $0 | 200 次测试/月，基础功能 |
| **Pro** | $25/月 | 无限测试，PR 检查，邮件告警 |
| **Enterprise** | 联系销售 | SSO，高级支持，定制功能 |

**开源项目**: 免费 Pro 计划

申请：https://snyk.io/plans/

---

## 相关资源

- [Snyk 文档](https://docs.snyk.io/)
- [漏洞数据库](https://security.snyk.io/)
- [GitHub 集成](https://docs.snyk.io/integrations/git-repository-scm-integrations/github-integration)
- [Slack 告警](https://docs.snyk.io/integrations/chat-ops-integrations/slack-integration)

---

*持续更新中...*
