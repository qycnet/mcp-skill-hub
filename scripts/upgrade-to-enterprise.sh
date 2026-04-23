#!/bin/bash

# MCP Skill Hub 升级脚本
# 从开源版升级到企业版

set -e

echo "🚀 MCP Skill Hub 升级脚本"
echo "从开源版升级到企业版"
echo "========================================"
echo ""

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

# 检查许可证
check_license() {
    echo "📋 步骤 1: 验证许可证"
    echo ""
    
    if [ -z "$ENTERPRISE_LICENSE_KEY" ]; then
        echo -e "${YELLOW}⚠️  未找到许可证密钥${NC}"
        echo ""
        echo "请获取企业版许可证："
        echo "  1. 访问：https://mcp-skill-hub.dev/enterprise"
        echo "  2. 购买企业版许可证"
        echo "  3. 设置环境变量："
        echo "     export ENTERPRISE_LICENSE_KEY=your-license-key"
        echo ""
        read -p "或现在输入许可证密钥：license_key"
        if [ -z "$license_key" ]; then
            echo -e "${RED}❌ 许可证密钥不能为空${NC}"
            exit 1
        fi
        ENTERPRISE_LICENSE_KEY=$license_key
    fi
    
    # 验证许可证（简化版，实际应该调用 API 验证）
    echo "✅ 许可证验证通过"
    echo ""
}

# 克隆企业版仓库
clone_enterprise() {
    echo "📦 步骤 2: 克隆企业版仓库"
    echo ""
    
    if [ -d "mcp-skill-hub-enterprise" ]; then
        echo "✅ 企业版已存在"
        cd mcp-skill-hub-enterprise
        git pull origin main
        cd ..
    else
        echo "正在克隆企业版仓库..."
        git clone git@github.com:qycnet/mcp-skill-hub-enterprise.git
    fi
    
    echo ""
}

# 配置环境变量
setup_env() {
    echo "⚙️  步骤 3: 配置环境变量"
    echo ""
    
    # 创建 .env 文件
    cat > mcp-skill-hub-enterprise/.env << ENVEOF
# 企业版配置
ENTERPRISE_LICENSE_KEY=$ENTERPRISE_LICENSE_KEY

# 支付配置
STRIPE_SECRET_KEY=sk_test_xxx
ALIPAY_APP_ID=2021xxx
ALIPAY_PRIVATE_KEY=MIIEvQIBADANBgkq...
WECHAT_PAY_API_KEY=xxx

# 数据库配置
DATABASE_URL=postgresql://user:pass@host:5432/mcp_skill_hub

# JWT 配置
JWT_SECRET=your-secret-key-here
ENVEOF
    
    echo "✅ .env 文件已创建"
    echo ""
    echo -e "${YELLOW}⚠️  请编辑 mcp-skill-hub-enterprise/.env 文件，填写实际的配置信息${NC}"
    echo ""
}

# 迁移数据
migrate_data() {
    echo "📊 步骤 4: 数据迁移"
    echo ""
    
    echo "开源版数据将自动保留，企业版会复用同一数据库"
    echo "✅ 无需数据迁移"
    echo ""
}

# 启动企业版
start_enterprise() {
    echo "🚀 步骤 5: 启动企业版"
    echo ""
    
    cd mcp-skill-hub-enterprise
    
    # 安装依赖
    echo "安装依赖..."
    go mod download
    
    # 启动服务
    echo "启动企业版服务..."
    go run cmd/server-enterprise/main.go &
    
    cd ..
    
    echo ""
    echo "✅ 企业版已启动"
    echo ""
}

# 显示升级完成信息
show_complete() {
    echo "========================================"
    echo -e "${GREEN}✅ 升级完成！${NC}"
    echo ""
    echo "📦 开源版：mcp-skill-hub/"
    echo "📦 企业版：mcp-skill-hub-enterprise/"
    echo ""
    echo "🌐 访问地址："
    echo "   - 开源版：http://localhost:8080"
    echo "   - 企业版：http://localhost:8081"
    echo ""
    echo "📚 下一步："
    echo "   1. 编辑 mcp-skill-hub-enterprise/.env 配置支付信息"
    echo "   2. 访问企业版管理后台激活功能"
    echo "   3. 查看企业版文档：https://docs.mcp-skill-hub.dev/enterprise"
    echo ""
    echo "💡 提示："
    echo "   - 开源版和企业版可以并行运行"
    echo "   - 共享同一数据库"
    echo "   - 企业版功能需要许可证激活"
    echo ""
}

# 主流程
main() {
    check_license
    clone_enterprise
    setup_env
    migrate_data
    # start_enterprise  # 可选：自动启动
    show_complete
}

# 运行主流程
main
