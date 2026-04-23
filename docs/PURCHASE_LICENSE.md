# 💰 购买企业版许可证

**官网**: https://enterprise.codevault.cn

---

## 🎯 许可证类型

| 类型 | 价格 | 有效期 | 功能 |
|------|------|--------|------|
| **试用版** | 免费 | 14 天 | 全部企业功能 |
| **专业版** | $25/月 | 按月 | 支付 + 订阅 + 分析 |
| **企业版** | $99/月 | 按年 | 全部功能 + 专属支持 |

---

## 🛒 购买流程

### 方式一：在线购买（推荐）⭐

```bash
# 1. 访问企业版官网
访问：https://enterprise.codevault.cn

# 2. 选择价格方案
- 试用版：免费 14 天
- 专业版：$25/月
- 企业版：$99/月

# 3. 完成支付
支持：Stripe、支付宝、微信支付

# 4. 获取许可证密钥
邮件自动发送许可证密钥
```

### 方式二：邮件购买

```bash
# 1. 发送邮件至
Email: tian@qycnet.cn

# 2. 邮件内容
主题：购买企业版许可证
内容：
- 公司名称：
- 联系人：
- 邮箱：
- 许可证类型：专业版/企业版

# 3. 等待回复
我们会发送支付信息和许可证密钥
```

### 方式三：微信购买

```bash
# 1. 添加微信
微信号：[待提供]

# 2. 发送消息
"购买企业版许可证"

# 3. 扫码支付
微信/支付宝扫码支付

# 4. 获取许可证
自动发送许可证密钥
```

---

## 🔑 激活许可证

### 方式一：升级脚本激活

```bash
# 运行升级脚本
cd mcp-skill-hub
export ENTERPRISE_LICENSE_KEY=your-license-key
./scripts/upgrade-to-enterprise.sh

# 许可证自动激活
```

### 方式二：手动激活

```bash
# 1. 设置许可证密钥
export ENTERPRISE_LICENSE_KEY=your-license-key

# 2. 编辑企业版配置
cd mcp-skill-hub-enterprise
nano .env

# 3. 添加许可证
ENTERPRISE_LICENSE_KEY=your-license-key

# 4. 重启企业版
go run cmd/server-enterprise/main.go
```

### 方式三：管理后台激活

```bash
# 1. 访问管理后台
访问：http://localhost:8081/admin/license

# 2. 输入许可证密钥
输入：your-license-key

# 3. 点击激活
激活成功！
```

---

## ✅ 验证许可证

### 命令行验证

```bash
# 检查许可证状态
curl -X POST http://localhost:8081/api/v1/license/verify \
  -H "Content-Type: application/json" \
  -d '{"license_key": "your-license-key"}'

# 返回结果
{
  "valid": true,
  "type": "enterprise",
  "features": ["payment", "subscription", "sso", ...],
  "valid_until": "2027-04-23"
}
```

### 代码验证

```go
import "github.com/qycnet/mcp-skill-hub-enterprise/internal/license"

// 验证许可证
lic, err := license.Verify("your-license-key")
if err != nil {
    log.Fatal("许可证无效:", err)
}

// 检查功能
if lic.HasFeature("payment") {
    // 启用支付功能
}
```

---

## ❓ 常见问题

### Q1: 许可证可以退款吗？

**A**: 可以！购买后 7 天内无条件退款。

### Q2: 许可证可以转移吗？

**A**: 可以！许可证可以转移至其他账户，联系 tian@qycnet.cn 办理。

### Q3: 许可证过期了怎么办？

**A**: 
- 过期后 30 天内：可以续费恢复
- 过期 30 天后：自动降级为开源版
- 数据不会丢失，企业功能不可用

### Q4: 可以升级许可证吗？

**A**: 可以！专业版可以补差价升级为企业版。

### Q5: 支持哪些支付方式？

**A**: 
- **国际**: Stripe（信用卡/借记卡）
- **国内**: 支付宝、微信支付

---

## 📞 获取帮助

购买或使用过程中遇到问题？

- 🌐 官网：https://enterprise.codevault.cn
- 📧 Email: tian@qycnet.cn
- 💬 社区：https://github.com/qycnet/mcp-skill-hub/discussions

---

**祝你使用愉快！** 🎉
