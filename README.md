# MCP Skill Hub

🚀 **企业级 MCP 技能市场与分发平台**

[![Go Version](https://img.shields.io/badge/go-1.21+-00ADD8?style=flat-square&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)](LICENSE)
[![CI/CD](https://github.com/mcp-skill-hub/mcp-skill-hub/actions/workflows/ci.yml/badge.svg)](https://github.com/mcp-skill-hub/mcp-skill-hub/actions)
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
git clone https://github.com/mcp-skill-hub/mcp-skill-hub.git
cd mcp-skill-hub

# 2. 一键启动
make docker-up

# 3. 访问服务
# Web 界面：http://localhost:8080
# API: http://localhost:8080/api/v1
# MinIO Console: http://localhost:9001
```

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
mcp-skill-hub/
├── cmd/                    # 应用入口
│   ├── server/             # 主服务
│   └── cli/mcp/            # CLI 工具
├── internal/               # 内部包
│   ├── api/                # HTTP handlers
│   ├── auth/               # 认证服务
│   ├── skill/              # 技能管理
│   ├── storage/            # 对象存储
│   ├── middleware/         # 中间件
│   └── email/              # 邮件服务
├── web/                    # React 前端
│   ├── src/
│   │   ├── components/     # 组件
│   │   ├── pages/          # 页面
│   │   ├── stores/         # 状态管理
│   │   └── api/            # API 客户端
│   └── package.json
├── cli/                    # CLI 工具
├── deployments/            # 部署配置
├── docs/                   # 文档
├── .github/workflows/      # CI/CD
├── docker-compose.yml      # Docker 编排
├── Dockerfile              # 生产镜像
├── Makefile                # 构建命令
└── README.md               # 本文件
```

---

## 🔐 安全特性

- ✅ **技能签名**: 所有发布技能使用 GPG 签名验证
- ✅ **依赖扫描**: 自动检测已知漏洞（集成 OSV、Snyk）
- ✅ **权限隔离**: 基于 RBAC 的细粒度权限控制
- ✅ **审计日志**: 所有操作记录到不可变日志
- ✅ **速率限制**: 防止滥用和 DDoS
- ✅ **私有网络**: 支持 VPC 内网部署

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

## 📚 文档

- [快速开始](docs/QUICKSTART.md) - 5 分钟上手
- [开发文档](DEVELOPMENT.md) - API、数据库、测试指南
- [贡献指南](CONTRIBUTING.md) - 如何贡献代码
- [项目状态](PROJECT_STATUS.md) - 进度跟踪
- [发布说明](RELEASE_NOTES.md) - 版本历史

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

- 🐛 问题反馈：[GitHub Issues](https://github.com/mcp-skill-hub/mcp-skill-hub/issues)
- 💬 讨论区：[GitHub Discussions](https://github.com/mcp-skill-hub/mcp-skill-hub/discussions)
- 📧 邮件：team@mcp-skill-hub.dev

---

## 🗺️ 路线图

### v0.4.0 (计划中)
- [ ] 安全扫描集成
- [ ] 集成测试
- [ ] 测试覆盖率 80%

### v0.5.0 (计划中)
- [ ] 付费技能支持
- [ ] 订阅制

### v1.0.0 (目标)
- [ ] 生产就绪
- [ ] 多语言支持
- [ ] 企业 SSO

---

<div align="center">

**Made with ❤️ by the MCP Skill Hub Team**

[📚 完整文档](docs/) | [🚀 快速开始](docs/QUICKSTART.md) | [💬 社区讨论](https://github.com/mcp-skill-hub/mcp-skill-hub/discussions)

</div>
