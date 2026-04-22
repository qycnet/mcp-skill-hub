# 贡献指南

感谢你对 MCP Skill Hub 的兴趣！我们欢迎各种形式的贡献。

## 🚀 快速开始

### 1. Fork 仓库

```bash
# Fork 并克隆
git clone https://github.com/YOUR_USERNAME/mcp-skill-hub.git
cd mcp-skill-hub

# 添加上游远程仓库
git remote add upstream https://github.com/qycnet/mcp-skill-hub.git
```

### 2. 设置开发环境

```bash
# 安装 Go 1.21+
go version

# 安装依赖
go mod download

# 启动开发环境
make dev
```

### 3. 创建分支

```bash
# 从 main 分支创建功能分支
git checkout -b feature/your-feature-name
```

## 📝 开发规范

### 代码风格

- 遵循 [Go 代码审查评论](https://github.com/golang/go/wiki/CodeReviewComments)
- 使用 `go fmt` 格式化代码
- 使用 `golangci-lint` 进行代码检查

```bash
make fmt
make lint
```

### 提交规范

遵循 [Conventional Commits](https://www.conventionalcommits.org/)：

```
feat: 添加新功能
fix: 修复 bug
docs: 文档更新
style: 代码格式调整
refactor: 代码重构
test: 测试相关
chore: 构建/工具相关
```

示例：
```bash
git commit -m "feat: 添加技能搜索功能"
git commit -m "fix: 修复版本比较逻辑"
```

### 测试要求

- 新功能必须包含单元测试
- 保持测试覆盖率 > 80%
- 运行所有测试通过

```bash
make test
```

## 🔄 提交流程

1. **创建 Issue** - 描述问题或功能建议
2. **Fork 仓库** - 创建你自己的副本
3. **创建分支** - `feature/xxx` 或 `fix/xxx`
4. **开发** - 编写代码和测试
5. **提交** - 遵循提交规范
6. **Push** - 推送到你的 Fork
7. **Pull Request** - 提交 PR 到主仓库
8. **Code Review** - 等待审查和反馈
9. **合并** - 通过后合并到 main

## 📚 资源

- [Go 官方文档](https://golang.org/doc/)
- [Gin Web 框架](https://gin-gonic.com/docs/)
- [GORM 文档](https://gorm.io/docs/)
- [MinIO 文档](https://min.io/docs/minio/linux/index.html)

## 💬 联系方式

- GitHub Issues: 报告 bug 或提出功能建议
- GitHub Discussions: 一般讨论和问题咨询

## 📄 许可证

MIT License - 贡献即表示你同意此许可
