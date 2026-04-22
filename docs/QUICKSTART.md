# 快速开始指南

## 🚀 5 分钟上手

### 方式一：Docker 部署（推荐）

```bash
# 1. 克隆仓库
git clone https://github.com/mcp-skill-hub/mcp-skill-hub.git
cd mcp-skill-hub

# 2. 启动所有服务
make docker-up

# 3. 验证安装
curl http://localhost:8080/health

# 4. 访问管理界面
open http://localhost:8080
```

### 方式二：本地开发

```bash
# 1. 环境要求
# - Go 1.21+
# - PostgreSQL 15+
# - MinIO (可选，可用本地存储替代)
# - Redis 7+

# 2. 克隆并安装依赖
git clone https://github.com/mcp-skill-hub/mcp-skill-hub.git
cd mcp-skill-hub
go mod download

# 3. 配置数据库
# 创建数据库
createdb mcp_skill_hub

# 4. 复制配置文件
cp config.yaml.example config.yaml

# 5. 启动服务
make dev
```

## 📋 API 快速参考

### 公开 API

```bash
# 获取技能列表
curl http://localhost:8080/api/v1/skills

# 搜索技能
curl "http://localhost:8080/api/v1/search?q=code"

# 获取技能详情
curl http://localhost:8080/api/v1/skills/1

# 获取分类
curl http://localhost:8080/api/v1/categories
```

### 认证 API

```bash
# 注册
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"test","email":"test@example.com","password":"secure123"}'

# 登录
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"test","password":"secure123"}'
```

### 发布技能（需要认证）

```bash
curl -X POST http://localhost:8080/api/v1/skills \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "name": "my-awesome-skill",
    "display_name": "My Awesome Skill",
    "description": "这是一个很棒的 MCP 技能",
    "category": "developer-tools",
    "tags": ["ai", "productivity"],
    "license": "MIT"
  }'
```

## 🔧 常见问题

### 端口被占用

修改 `docker-compose.yml` 中的端口映射：
```yaml
ports:
  - "8081:8080"  # 改为其他端口
```

### 数据库连接失败

检查 PostgreSQL 是否运行：
```bash
docker-compose ps
docker-compose logs postgres
```

### 查看日志

```bash
# 查看所有服务日志
docker-compose logs -f

# 查看特定服务日志
docker-compose logs -f app
```

## 📚 下一步

- [架构设计](docs/architecture.md) - 了解系统架构
- [API 文档](docs/api.md) - 完整 API 参考
- [开发指南](docs/development.md) - 贡献代码
- [部署指南](docs/deployment.md) - 生产环境部署

## 💬 获取帮助

- 📖 [完整文档](docs/)
- 🐛 [问题反馈](https://github.com/mcp-skill-hub/mcp-skill-hub/issues)
- 💬 [社区讨论](https://github.com/mcp-skill-hub/mcp-skill-hub/discussions)
