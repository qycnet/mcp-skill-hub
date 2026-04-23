# 🚀 从开源版升级至企业版指南

**更新时间**: 2026-04-23 10:30 GMT+8

---

## 📋 升级方式

### 方式一：一键升级脚本（推荐）⭐

```bash
# 进入开源版目录
cd mcp-skill-hub

# 运行升级脚本
./scripts/upgrade-to-enterprise.sh
```

脚本会自动：
1. ✅ 验证许可证
2. ✅ 克隆企业版仓库
3. ✅ 配置环境变量
4. ✅ 准备数据库

---

### 方式二：手动升级

#### 步骤 1: 获取许可证

访问：https://mcp-skill-hub.dev/enterprise

| 版本 | 价格 | 功能 |
|------|------|------|
| **专业版** | $25/月 | 支付 + 订阅 + 分析 |
| **企业版** | $99/月 | 全部功能 + SSO + 审计 |

获取许可证密钥后：

```bash
export ENTERPRISE_LICENSE_KEY=your-license-key
```

#### 步骤 2: 克隆企业版

```bash
# 克隆企业版仓库
git clone git@github.com:qycnet/mcp-skill-hub-enterprise.git

# 进入企业版目录
cd mcp-skill-hub-enterprise
```

#### 步骤 3: 配置环境

创建 `.env` 文件：

```bash
# 企业版配置
ENTERPRISE_LICENSE_KEY=your-license-key

# 支付配置
STRIPE_SECRET_KEY=sk_test_xxx
ALIPAY_APP_ID=2021xxx
ALIPAY_PRIVATE_KEY=MIIEvQIBADANBgkq...
WECHAT_PAY_API_KEY=xxx

# 数据库配置（与开源版相同）
DATABASE_URL=postgresql://user:pass@host:5432/mcp_skill_hub

# JWT 配置
JWT_SECRET=your-secret-key-here
```

#### 步骤 4: 启动企业版

```bash
# 安装依赖
go mod download

# 启动服务
go run cmd/server-enterprise/main.go
```

---

## 🎯 升级后的目录结构

```
workspace/
├── mcp-skill-hub/                 # 开源版
│   ├── cmd/server/
│   ├── internal/
│   ├── cli/
│   ├── web/
│   └── ...
│
└── mcp-skill-hub-enterprise/      # 企业版
    ├── cmd/server-enterprise/
    ├── internal/
    │   ├── payment/
    │   ├── auth-pro/
    │   ├── billing/
    │   └── ...
    └── ...
```

---

## 📊 功能对比

| 功能 | 开源版 | 企业版 |
|------|--------|--------|
| **技能管理** | ✅ | ✅ |
| **用户认证** | ✅ (JWT) | ✅ (JWT + SSO) |
| **CLI 工具** | ✅ | ✅ |
| **Web 前端** | ✅ | ✅ |
| **支付集成** | ❌ | ✅ (Stripe/支付宝/微信) |
| **订阅管理** | ❌ | ✅ |
| **企业 SSO** | ❌ | ✅ (SAML/OIDC) |
| **数据分析** | ❌ | ✅ |
| **审计日志** | ❌ | ✅ (完整合规) |
| **监控指标** | ❌ | ✅ (Prometheus) |
| **邮件通知** | ✅ (基础) | ✅ (企业) |
| **技术支持** | 社区 | 专属支持 |

---

## 💡 常见问题

### Q1: 升级后开源版还能用吗？

**A**: 可以！开源版和企业版可以并行运行，互不影响。

```bash
# 开源版（端口 8080）
cd mcp-skill-hub
go run cmd/server/main.go

# 企业版（端口 8081）
cd mcp-skill-hub-enterprise
go run cmd/server-enterprise/main.go
```

### Q2: 数据会丢失吗？

**A**: 不会！企业版和开源版共享同一数据库，数据完全保留。

```yaml
# 两个版本使用相同的数据库配置
DATABASE_URL=postgresql://user:pass@host:5432/mcp_skill_hub
```

### Q3: 可以随时降级吗？

**A**: 可以！企业版功能通过许可证控制，停止续费后自动降级为开源版功能。

### Q4: 支持哪些支付方式？

**A**: 企业版支持：
- **Stripe** - 国际信用卡/借记卡
- **支付宝** - 国内网页支付
- **微信支付** - V2/V3 版本

### Q5: 企业 SSO 支持哪些服务商？

**A**: 支持：
- **SAML 2.0**: ADFS, Okta, OneLogin
- **OIDC**: Google, Azure AD, Keycloak

---

## 🔐 许可证验证

### 本地验证

```bash
# 检查许可证
curl -X POST http://localhost:8081/api/v1/license/verify \
  -H "Content-Type: application/json" \
  -d '{"license_key": "your-license-key"}'
```

### API 验证

```go
import "github.com/qycnet/mcp-skill-hub-enterprise/internal/license"

// 验证许可证
valid, err := license.Verify(ENTERPRISE_LICENSE_KEY)
if err != nil {
    log.Fatal("许可证无效:", err)
}
```

---

## 📞 获取帮助

升级过程中遇到问题？

- 📧 Email: tian@qycnet.cn
- 📚 文档：https://docs.mcp-skill-hub.dev/enterprise
- 💬 社区：https://github.com/qycnet/mcp-skill-hub/discussions

---

**祝你升级顺利！** 🎉
