# MCP Skill Hub

🚀 **企业级 MCP 技能市场与分发平台**

[![Go Version](https://img.shields.io/badge/go-1.21+-00ADD8?style=flat-square&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)](LICENSE)
[![Docker](https://img.shields.io/badge/docker-ready-2496ED?style=flat-square&logo=docker)](https://hub.docker.com)

---

## 📖 简介

MCP Skill Hub 是一个标准化的 **Model Context Protocol (MCP) 技能发布与发现平台**，帮助企业和个人开发者：

- 📦 **发布技能**: 将自定义能力封装为标准 MCP 技能
- 🔍 **发现技能**: 浏览、搜索、评分高质量技能
- 🔐 **安全分发**: 企业级权限控制和审计日志
- 🚀 **一键部署**: Docker 自托管，数据完全自控

> 类似 "npm for Agent Skills" —— 让 AI Agent 能力像 Node.js 包一样易于共享和复用

---

## ✨ 核心特性

### 对于技能消费者
- 🔍 智能搜索：按类别、评分、下载量筛选
- 📊 质量评分：基于代码质量、社区活跃度、安全扫描
- 🔐 安全验证：自动扫描恶意代码和权限滥用
- 📥 一键安装：`mcp install <skill-name>`

### 对于技能开发者
- 📦 标准化发布：遵循 MCP 协议规范
- 📈 数据分析：下载量、用户反馈、使用统计
- 🔖 版本管理：语义化版本控制
- 💰 商业化支持：付费技能、订阅制（后续版本）

### 对于企业
- 👥 权限管理：RBAC 角色权限系统
- 📋 审计日志：完整的操作追踪
- 🔒 私有仓库：内网隔离部署
- 🛡️ 安全合规：SOC2 就绪架构

---

## 🚀 快速开始

### Docker 部署（推荐）

```bash
# 1. 克隆仓库
git clone https://github.com/qycnet/mcp-skill-hub.git
cd mcp-skill-hub

# 2. 配置环境变量
cp .env.example .env
# 编辑 .env 文件，填入真实的配置信息

# 3. 安装 Trivy（安全扫描必需）
# macOS
brew install trivy

# Ubuntu/Debian
sudo apt-get install trivy

# 或访问 https://aquasecurity.github.io/trivy/latest/getting-started/installation/

# 4. 一键启动
make docker-up

# 5. 访问服务
# Web 界面：http://localhost:8080
# API: http://localhost:8080/api/v1
# MinIO Console: http://localhost:9001
```

> ⚠️ **重要安全配置**：
> - 必须配置 `JWT_SECRET`（至少 32 字符）
> - 必须配置 `DATABASE_PASSWORD`（强密码）
> - 必须配置 `STORAGE_SECRET_KEY`（MinIO 密钥）
> - 生产环境必须配置 `CORS_ALLOWED_ORIGINS`（具体域名）
> - 否则服务无法启动

### 本地开发

```bash
# 安装依赖
go mod download

# 启动服务
make dev

# 运行测试
make test

# 前端开发
cd web && npm install && npm run dev
```

### CLI 工具

```bash
# 安装
cd cli && go install

# 登录
mcp login

# 搜索技能
mcp search "code analysis"

# 安装技能
mcp install claude-context

# 发布技能
mcp publish ./my-skill
```

---

## 📦 使用示例

### Web 界面

访问 http://localhost:8080 浏览技能市场：
- 首页：热门技能、高评分技能、搜索
- 技能库：筛选、排序、分页
- 技能详情：信息、版本、评分、评论
- 个人中心：资料、我的技能、API Keys
- 发布技能：上传 ZIP 或手动填写

### REST API

```bash
# 获取技能列表
curl http://localhost:8080/api/v1/skills

# 搜索技能
curl "http://localhost:8080/api/v1/search?q=code"

# 获取技能详情
curl http://localhost:8080/api/v1/skills/1

# 用户注册
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"test","email":"test@example.com","password":"secure123"}'

# 用户登录
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"test","password":"secure123"}'
```

---

## 🏗️ 架构设计

```
┌─────────────────────────────────────────────────────────────┐
│                      MCP Skill Hub                          │
├─────────────────────────────────────────────────────────────┤
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐  │
│  │   Web UI    │  │   CLI Tool  │  │   REST API /gRPC    │  │
│  │  (React)    │  │    (Go)     │  │      (Go)           │  │
│  └──────┬──────┘  └──────┬──────┘  └──────────┬──────────┘  │
│         │                │                      │            │
│         └────────────────┼──────────────────────┘            │
│                          │                                   │
│              ┌───────────▼───────────┐                       │
│              │    Core Services      │                       │
│              ├───────────────────────┤                       │
│              │  - Skill Registry     │                       │
│              │  - Search Engine      │                       │
│              │  - Auth & RBAC        │                       │
│              │  - Security Scanner   │                       │
│              │  - Analytics          │                       │
│              └───────────┬───────────┘                       │
│                          │                                   │
│         ┌────────────────┼────────────────┐                 │
│         │                │                │                 │
│  ┌──────▼──────┐ ┌──────▼──────┐ ┌──────▼──────┐           │
│  │ PostgreSQL  │ │   MinIO     │ │   Redis     │           │
│  │  (Metadata) │ │  (Storage)  │ │   (Cache)   │           │
│  └─────────────┘ └─────────────┘ └─────────────┘           │
└─────────────────────────────────────────────────────────────┘
```

---

## 📁 项目结构

```
mcp-skill-hub/ (开源版)
├── cmd/
│   └── server/
│       └── main.go                  # 主服务入口
├── internal/
│   ├── api/                         # HTTP handlers
│   │   ├── handlers.go              # 技能 CRUD 接口
│   │   └── auth_handlers.go         # 认证接口
│   ├── auth/                        # 认证服务 (JWT)
│   │   ├── model.go                 # User/APIKey 模型
│   │   ├── service.go               # JWT/密码加密
│   │   ├── middleware.go            # JWT 中间件
│   │   └── service_test.go          # 单元测试
│   ├── skill/                       # 技能管理
│   │   ├── model.go                 # Skill 数据模型
│   │   ├── service.go               # CRUD 逻辑
│   │   ├── service_test.go          # 单元测试
│   │   └── upload.go                # ZIP 上传/解析
│   ├── cache/                       # Redis 缓存
│   │   └── service.go               # 缓存服务
│   ├── i18n/                        # 国际化
│   │   ├── locales.go               # 翻译器
│   │   └── translations/            # 翻译文件
│   │       ├── locale_zh.json       # 中文
│   │       └── locale_en.json       # 英文
│   ├── middleware/                  # 中间件
│   │   ├── ratelimit.go             # 速率限制
│   │   └── audit.go                 # 审计日志
│   ├── storage/                     # 对象存储
│   │   └── minio.go                 # MinIO 集成
│   └── email/                       # 邮件服务
│       └── service.go               # SMTP 邮件
├── cli/                             # 命令行工具 (16 个文件)
│   ├── cmd/mcp/                     # 11 个命令
│   │   ├── main.go                  # CLI 入口
│   │   ├── login.go                 # 登录
│   │   ├── search.go                # 搜索
│   │   ├── install.go               # 安装
│   │   ├── list.go                  # 列出
│   │   ├── info.go                  # 详情
│   │   ├── publish.go               # 发布
│   │   ├── update.go                # 更新
│   │   ├── uninstall.go             # 卸载
│   │   ├── version.go               # 版本
│   │   ├── config.go                # 配置
│   │   └── whoami.go                # 当前用户
│   ├── internal/                    # CLI 内部库
│   │   ├── api/client.go            # API 客户端
│   │   └── config/config.go         # 配置加载
│   └── README.md                    # CLI 文档
├── web/                             # React 前端 (15 个文件)
│   ├── src/
│   │   ├── pages/                   # 8 个页面
│   │   │   ├── HomePage.jsx         # 首页
│   │   │   ├── SkillListPage.jsx    # 技能列表
│   │   │   ├── SkillDetailPage.jsx  # 技能详情
│   │   │   ├── LoginPage.jsx        # 登录/注册
│   │   │   ├── ProfilePage.jsx      # 个人中心
│   │   │   ├── PublishPage.jsx      # 发布技能
│   │   │   ├── PricingPage.jsx      # 价格
│   │   │   └── SubscriptionPage.jsx # 订阅管理
│   │   ├── components/              # 可复用组件
│   │   │   ├── Layout.jsx           # 布局
│   │   │   └── LanguageSwitcher.jsx # 语言切换
│   │   ├── stores/                  # 状态管理
│   │   │   └── authStore.js         # 认证状态
│   │   └── api/                     # API 客户端
│   │       └── client.js            # Axios 封装
│   └── (配置文件)
├── docs/                            # 文档
│   ├── QUICKSTART.md                # 5 分钟快速开始
│   └── SNYK_SETUP.md                # Snyk 配置
├── .github/workflows/               # CI/CD
│   ├── ci.yml                       # 持续集成
│   └── snyk.yml                     # 安全扫描
├── docker-compose.yml               # Docker 编排
├── Dockerfile                       # 生产镜像
├── Makefile                         # 构建命令
├── go.mod                           # Go 依赖
├── README.md                        # 项目介绍
├── CONTRIBUTING.md                  # 贡献指南
└── LICENSE                          # MIT 许可
```

---

## 🔐 安全特性

### 核心安全功能
- ✅ **安全扫描**: 集成 Trivy 自动扫描漏洞（HIGH/CRITICAL 自动阻断）
- ✅ **沙箱隔离**: Docker 容器隔离执行，资源限制（内存/CPU/PID）
- ✅ **密码加密**: 使用 bcrypt 进行密码哈希
- ✅ **JWT 认证**: 基于 JWT 的无状态认证机制（最少 32 字符密钥）
- ✅ **权限隔离**: 基于 RBAC 的细粒度权限控制
- ✅ **速率限制**: 多级限制（IP/用户/API 端点），防止滥用和 DDoS
- ✅ **CORS 加固**: 生产环境强制配置具体域名
- ✅ **安全配置**: 敏感信息强制要求环境变量配置，无默认弱密码
- ✅ **事务处理**: 数据库事务保证数据一致性
- ✅ **审计日志**: 完整操作追踪（企业版）

### 沙箱安全措施
- 🐳 容器隔离执行
- 🔒 只读文件系统挂载
- 👤 非 root 用户运行（UID 1000）
- 🌐 网络隔离（默认无网络访问）
- ⏱️ 超时控制（默认 30 秒）
- 📊 资源限制（内存 512m，CPU 1.0，PIDs 100）

### 计划中的功能
- 🚧 **技能签名**: GPG 签名验证
- 🚧 **WebAssembly 运行时**: 替代 Docker 的轻量级隔离

---

## 📊 质量评分系统

技能质量评分（满分 100）：

| 维度 | 权重 | 评估指标 |
|------|------|----------|
| 代码质量 | 25% | 测试覆盖率、Lint 评分、文档完整度 |
| 安全性 | 25% | 漏洞扫描、权限最小化、签名验证 |
| 社区活跃 | 20% | 下载量、评分、Issue 响应速度 |
| 兼容性 | 15% | MCP 协议版本、跨平台支持 |
| 维护性 | 15% | 版本更新频率、向后兼容 |

---

## 🧪 测试

```bash
# 运行所有测试
make test

# 带覆盖率
go test -cover ./...

# 生成覆盖率报告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

---

## 🛡️ 安全最佳实践

### 生产环境部署清单

- [ ] **配置强密钥**
  ```bash
  # 生成强 JWT 密钥（至少 32 字符）
  openssl rand -base64 32
  
  # 生成数据库密码
  openssl rand -base64 24
  ```

- [ ] **配置 CORS**
  ```bash
  # .env
  CORS_ALLOWED_ORIGINS=https://yourdomain.com,https://app.yourdomain.com
  ```

- [ ] **启用安全扫描**
  ```bash
  # 安装 Trivy
  brew install trivy  # macOS
  sudo apt-get install trivy  # Ubuntu
  
  # 配置环境变量
  SECURITY_SCAN_ENABLED=true
  TRIVY_PATH=trivy
  ```

- [ ] **配置速率限制**
  ```bash
  # 默认限制已在代码中配置
  # 登录：1 req/s
  # 搜索：5 req/s
  # 管理：20 req/s
  ```

- [ ] **启用审计日志**（企业版）
  ```bash
  AUDIT_ENABLED=true
  AUDIT_RETENTION_DAYS=90
  ```

### 安全扫描工作流

```
技能上传 → ZIP 解析 → Trivy 扫描 → 漏洞评估
                                    ↓
                            HIGH/CRITICAL → 阻止上传
                            MEDIUM/LOW → 允许上传 + 警告
```

### 沙箱执行流程

```
技能执行请求 → 创建容器 → 资源限制 → 执行 → 清理
                    ↓
            - 内存：512m
            - CPU：1.0
            - 网络：无
            - 用户：UID 1000
            - 文件系统：只读
```

---

## 📚 文档

- [快速开始](docs/QUICKSTART.md) - 5 分钟上手
- [贡献指南](CONTRIBUTING.md) - 如何贡献代码

---

## 🤝 贡献

我们欢迎各种形式的贡献！

```bash
# Fork 并克隆
git clone https://github.com/YOUR_USERNAME/mcp-skill-hub.git
cd mcp-skill-hub

# 创建分支
git checkout -b feature/your-feature

# 开发并提交
git commit -m "feat: add new feature"

# 提交 PR
```

详见 [CONTRIBUTING.md](CONTRIBUTING.md)

---

## 📄 许可证

MIT License - 详见 [LICENSE](LICENSE)

---

## 📬 联系方式

- 🐛 问题反馈：[GitHub Issues](https://github.com/qycnet/mcp-skill-hub/issues)
- 💬 讨论区：[GitHub Discussions](https://github.com/qycnet/mcp-skill-hub/discussions)
- 📧 邮件：tian@qycnet.cn

---

## 🗺️ 路线图

### v0.2.0 (当前 - 安全加固版)
- [x] 安全扫描集成（Trivy）
- [x] Docker 沙箱隔离执行
- [x] 数据库事务支持
- [x] 速率限制中间件
- [x] CORS 安全加固
- [x] 移除默认弱密码
- [x] 审计日志完整实现（企业版）
- [x] 数据分析完整实现（企业版）

### v0.1.0
- [x] 基础 API 框架
- [x] 用户认证（JWT）
- [x] 技能 CRUD
- [x] Docker 部署
- [x] CORS 支持

### v0.3.0 (计划中)
- [ ] WebAssembly 运行时支持
- [ ] 集成测试完善
- [ ] 测试覆盖率 80%
- [ ] GPG 签名验证

### v0.4.0 (计划中)
- [ ] 付费技能支持
- [ ] 订阅制
- [ ] 多租户支持

### v1.0.0 (目标)
- [ ] 生产就绪
- [ ] 多语言支持
- [ ] 企业 SSO
- [ ] Kubernetes Helm Chart

---

<div align="center">

**Made with ❤️ by the MCP Skill Hub Team**

[📚 完整文档](docs/) | [🚀 快速开始](docs/QUICKSTART.md) | [📝 更新日志](CHANGELOG.md) | [💬 社区讨论](https://github.com/qycnet/mcp-skill-hub/discussions)

</div>
