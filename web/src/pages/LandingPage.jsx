import { Link } from 'react-router-dom'
import { useTranslation } from 'react-i18next'
import { 
  Zap, Shield, Globe, Users, TrendingUp, 
  CheckCircle, ArrowRight, Star, Code 
} from 'lucide-react'

export default function LandingPage() {
  const { t } = useTranslation()

  return (
    <div className="min-h-screen">
      {/* Hero 区域 */}
      <section className="bg-gradient-to-r from-primary-600 to-primary-700 text-white py-20">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center">
            <h1 className="text-5xl font-bold mb-6">
              MCP Skill Hub
            </h1>
            <p className="text-2xl mb-8 text-primary-100">
              发现和分享 MCP 技能，让 AI Agent 更强大
            </p>
            <div className="flex justify-center gap-4">
              <Link to="/pricing" className="bg-white text-primary-600 px-8 py-4 rounded-lg font-semibold hover:bg-primary-50 transition-colors">
                开始免费使用
              </Link>
              <Link to="/skills" className="border-2 border-white text-white px-8 py-4 rounded-lg font-semibold hover:bg-white/10 transition-colors">
                浏览技能库
              </Link>
            </div>
          </div>
        </div>
      </section>

      {/* 特性展示 */}
      <section className="py-20 bg-gray-50 dark:bg-gray-900">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <h2 className="text-3xl font-bold text-center mb-12">
            为什么选择 MCP Skill Hub？
          </h2>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
            <FeatureCard
              icon={Zap}
              title="快速集成"
              description="一键安装技能，立即提升 AI Agent 能力"
            />
            <FeatureCard
              icon={Shield}
              title="安全可靠"
              description="所有技能经过安全扫描，企业级权限控制"
            />
            <FeatureCard
              icon={Globe}
              title="全球市场"
              description="发现和分享技能，连接全球开发者"
            />
            <FeatureCard
              icon={Users}
              title="社区驱动"
              description="评分、评论、反馈，共同成长"
            />
            <FeatureCard
              icon={TrendingUp}
              title="数据分析"
              description="实时追踪技能使用情况和趋势"
            />
            <FeatureCard
              icon={Code}
              title="开放生态"
              description="基于 MCP 协议，兼容主流 AI 框架"
            />
          </div>
        </div>
      </section>

      {/* 价格预览 */}
      <section className="py-20">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <h2 className="text-3xl font-bold text-center mb-4">
            简单透明的定价
          </h2>
          <p className="text-center text-gray-600 dark:text-gray-400 mb-12">
            从免费开始，根据需要随时升级
          </p>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
            <PricingCard
              name="免费版"
              price="$0"
              period="/月"
              description="适合个人开发者"
              features={[
                "5 个技能发布",
                "100 次下载/月",
                "1GB 存储空间",
                "社区支持"
              ]}
              cta="免费开始"
              popular={false}
            />
            <PricingCard
              name="专业版"
              price="$25"
              period="/月"
              description="适合专业开发者"
              features={[
                "无限技能发布",
                "10,000 次下载/月",
                "10GB 存储空间",
                "数据分析",
                "优先支持",
                "14 天免费试用"
              ]}
              cta="开始试用"
              popular={true}
            />
            <PricingCard
              name="企业版"
              price="$99"
              period="/月"
              description="适合企业团队"
              features={[
                "无限技能发布",
                "无限下载",
                "100GB 存储空间",
                "数据分析",
                "专属支持",
                "自定义域名",
                "SSO 集成",
                "99.9% SLA"
              ]}
              cta="联系销售"
              popular={false}
            />
          </div>
          <div className="text-center mt-12">
            <Link to="/pricing" className="text-primary-600 hover:text-primary-700 font-semibold">
              查看所有功能对比 →
            </Link>
          </div>
        </div>
      </section>

      {/* 用户评价 */}
      <section className="py-20 bg-gray-50 dark:bg-gray-900">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <h2 className="text-3xl font-bold text-center mb-12">
            用户怎么说
          </h2>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
            <TestimonialCard
              quote="MCP Skill Hub 让我们的 AI 产品开发效率提升了 300%！"
              author="张三"
              role="CTO @ 某科技公司"
              avatar="Z"
            />
            <TestimonialCard
              quote="技能市场生态太棒了，找到了很多现成的解决方案"
              author="李四"
              role="独立开发者"
              avatar="L"
            />
            <TestimonialCard
              quote="企业版功能强大，SSO 集成和权限管理非常专业"
              author="王五"
              role="技术总监 @ 某企业"
              avatar="W"
            />
          </div>
        </div>
      </section>

      {/* CTA 区域 */}
      <section className="py-20 bg-primary-600 text-white">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 text-center">
          <h2 className="text-4xl font-bold mb-6">
            准备好开始了吗？
          </h2>
          <p className="text-xl mb-8 text-primary-100">
            立即注册，享受 14 天免费试用
          </p>
          <Link 
            to="/register" 
            className="bg-white text-primary-600 px-12 py-5 rounded-lg font-semibold text-lg hover:bg-primary-50 transition-colors inline-flex items-center"
          >
            免费注册
            <ArrowRight className="ml-2 w-5 h-5" />
          </Link>
          <p className="mt-4 text-primary-100 text-sm">
            无需信用卡 · 随时取消
          </p>
        </div>
      </section>

      {/* 页脚 */}
      <footer className="bg-gray-900 text-gray-400 py-12">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="grid grid-cols-1 md:grid-cols-4 gap-8">
            <div>
              <h3 className="text-white font-semibold mb-4">产品</h3>
              <ul className="space-y-2">
                <li><Link to="/skills" className="hover:text-white">技能库</Link></li>
                <li><Link to="/pricing" className="hover:text-white">价格</Link></li>
                <li><Link to="/docs" className="hover:text-white">文档</Link></li>
              </ul>
            </div>
            <div>
              <h3 className="text-white font-semibold mb-4">公司</h3>
              <ul className="space-y-2">
                <li><Link to="/about" className="hover:text-white">关于我们</Link></li>
                <li><Link to="/blog" className="hover:text-white">博客</Link></li>
                <li><Link to="/careers" className="hover:text-white">招聘</Link></li>
              </ul>
            </div>
            <div>
              <h3 className="text-white font-semibold mb-4">支持</h3>
              <ul className="space-y-2">
                <li><Link to="/help" className="hover:text-white">帮助中心</Link></li>
                <li><Link to="/contact" className="hover:text-white">联系我们</Link></li>
                <li><Link to="/status" className="hover:text-white">状态</Link></li>
              </ul>
            </div>
            <div>
              <h3 className="text-white font-semibold mb-4">法律</h3>
              <ul className="space-y-2">
                <li><Link to="/privacy" className="hover:text-white">隐私政策</Link></li>
                <li><Link to="/terms" className="hover:text-white">服务条款</Link></li>
                <li><Link to="/security" className="hover:text-white">安全</Link></li>
              </ul>
            </div>
          </div>
          <div className="border-t border-gray-800 mt-12 pt-8 text-center">
            <p>&copy; 2026 MCP Skill Hub. All rights reserved.</p>
          </div>
        </div>
      </footer>
    </div>
  )
}

// 特性卡片组件
function FeatureCard({ icon: Icon, title, description }) {
  return (
    <div className="text-center p-6">
      <div className="w-16 h-16 bg-primary-100 dark:bg-primary-900/30 rounded-full flex items-center justify-center mx-auto mb-4">
        <Icon className="w-8 h-8 text-primary-600 dark:text-primary-400" />
      </div>
      <h3 className="text-xl font-semibold mb-2">{title}</h3>
      <p className="text-gray-600 dark:text-gray-400">{description}</p>
    </div>
  )
}

// 价格卡片组件
function PricingCard({ name, price, period, description, features, cta, popular }) {
  return (
    <div className={`relative p-8 rounded-2xl border-2 ${
      popular 
        ? 'border-primary-500 bg-primary-50 dark:bg-primary-900/20' 
        : 'border-gray-200 dark:border-gray-700'
    }`}>
      {popular && (
        <div className="absolute -top-4 left-1/2 -translate-x-1/2">
          <span className="bg-primary-500 text-white px-4 py-1 rounded-full text-sm font-semibold">
            最受欢迎
          </span>
        </div>
      )}
      <div className="text-center mb-6">
        <h3 className="text-2xl font-bold mb-2">{name}</h3>
        <p className="text-gray-600 dark:text-gray-400 mb-4">{description}</p>
        <div className="flex items-baseline justify-center">
          <span className="text-5xl font-bold">{price}</span>
          <span className="text-gray-500 ml-2">{period}</span>
        </div>
      </div>
      <ul className="space-y-4 mb-8">
        {features.map((feature, i) => (
          <li key={i} className="flex items-center">
            <CheckCircle className="w-5 h-5 text-green-500 mr-3 flex-shrink-0" />
            <span>{feature}</span>
          </li>
        ))}
      </ul>
      <button className={`w-full py-3 rounded-lg font-semibold ${
        popular
          ? 'bg-primary-600 hover:bg-primary-700 text-white'
          : 'bg-gray-100 dark:bg-gray-800 hover:bg-gray-200 dark:hover:bg-gray-700'
      }`}>
        {cta}
      </button>
    </div>
  )
}

// 用户评价卡片
function TestimonialCard({ quote, author, role, avatar }) {
  return (
    <div className="bg-white dark:bg-gray-800 p-8 rounded-2xl shadow-lg">
      <div className="flex items-center mb-4">
        <div className="w-12 h-12 bg-gray-200 dark:bg-gray-700 rounded-full flex items-center justify-center mr-4">
          <span className="text-lg font-semibold">{avatar}</span>
        </div>
        <div>
          <div className="font-semibold">{author}</div>
          <div className="text-sm text-gray-500">{role}</div>
        </div>
      </div>
      <div className="flex text-yellow-500 mb-4">
        {[1, 2, 3, 4, 5].map(i => (
          <Star key={i} className="w-5 h-5 fill-current" />
        ))}
      </div>
      <p className="text-gray-700 dark:text-gray-300 italic">"{quote}"</p>
    </div>
  )
}
