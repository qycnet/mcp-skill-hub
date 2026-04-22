import { useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { Link } from 'react-router-dom'
import { paymentApi } from '../api/client'
import { Check, X, Star, Zap, Building2 } from 'lucide-react'

export default function PricingPage() {
  const [billingInterval, setBillingInterval] = useState('monthly')
  const [couponCode, setCouponCode] = useState('')

  // 获取价格计划
  const { data: plans } = useQuery({
    queryKey: ['pricing-plans'],
    queryFn: paymentApi.getPlans,
    select: (res) => res.data.plans,
  })

  const handleSubscribe = async (planId) => {
    // TODO: 创建结账会话并重定向到支付页面
    console.log('订阅计划:', planId)
  }

  const handleValidateCoupon = async () => {
    // TODO: 验证优惠券
    console.log('验证优惠券:', couponCode)
  }

  return (
    <div className="space-y-12">
      {/* 页面头部 */}
      <div className="text-center">
        <h1 className="text-4xl font-bold text-gray-900 dark:text-white mb-4">
          选择适合你的计划
        </h1>
        <p className="text-xl text-gray-600 dark:text-gray-400 mb-8">
          从免费开始，根据需要随时升级
        </p>

        {/* 计费周期切换 */}
        <div className="flex items-center justify-center gap-4">
          <span className={`text-sm ${billingInterval === 'monthly' ? 'font-semibold' : 'text-gray-500'}`}>
            月付
          </span>
          <button
            onClick={() => setBillingInterval(billingInterval === 'monthly' ? 'yearly' : 'monthly')}
            className="relative w-14 h-7 bg-gray-200 dark:bg-gray-700 rounded-full transition-colors"
          >
            <div
              className={`absolute top-1 w-5 h-5 bg-white rounded-full shadow transition-transform ${
                billingInterval === 'yearly' ? 'translate-x-8' : 'translate-x-1'
              }`}
            />
          </button>
          <span className={`text-sm ${billingInterval === 'yearly' ? 'font-semibold' : 'text-gray-500'}`}>
            年付 <span className="text-green-500 text-xs">(省 20%)</span>
          </span>
        </div>
      </div>

      {/* 价格卡片 */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
        {plans?.map((plan) => (
          <PricingCard
            key={plan.id}
            plan={plan}
            interval={billingInterval}
            onSubscribe={handleSubscribe}
          />
        ))}
      </div>

      {/* 功能对比表 */}
      <div className="card">
        <h2 className="text-2xl font-bold text-center mb-8">功能对比</h2>
        <div className="overflow-x-auto">
          <table className="w-full">
            <thead>
              <tr className="border-b">
                <th className="text-left py-4 px-4">功能</th>
                <th className="text-center py-4 px-4">免费版</th>
                <th className="text-center py-4 px-4 bg-primary-50 dark:bg-primary-900/20">专业版</th>
                <th className="text-center py-4 px-4">企业版</th>
              </tr>
            </thead>
            <tbody>
              <FeatureRow
                feature="技能发布数量"
                free="5 个"
                pro="无限"
                enterprise="无限"
              />
              <FeatureRow
                feature="月下载量"
                free="100 次"
                pro="10,000 次"
                enterprise="无限"
              />
              <FeatureRow
                feature="存储空间"
                free="1 GB"
                pro="10 GB"
                enterprise="100 GB"
              />
              <FeatureRow
                feature="数据分析"
                free={false}
                pro={true}
                enterprise={true}
              />
              <FeatureRow
                feature="优先支持"
                free={false}
                pro={true}
                enterprise={true}
              />
              <FeatureRow
                feature="自定义域名"
                free={false}
                pro={false}
                enterprise={true}
              />
              <FeatureRow
                feature="SSO 集成"
                free={false}
                pro={false}
                enterprise={true}
              />
              <FeatureRow
                feature="SLA 保障"
                free={false}
                pro={false}
                enterprise="99.9%"
              />
            </tbody>
          </table>
        </div>
      </div>

      {/* 常见问题 */}
      <div>
        <h2 className="text-2xl font-bold text-center mb-8">常见问题</h2>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          <FAQItem
            question="我可以随时取消订阅吗？"
            answer="是的，你可以随时取消订阅。取消后，你的订阅将在当前计费周期结束时失效。"
          />
          <FAQItem
            question="支持哪些支付方式？"
            answer="我们支持信用卡、支付宝、微信支付等多种支付方式。"
          />
          <FAQItem
            question="有试用期吗？"
            answer="专业版和企业版都提供 14 天免费试用，无需绑定信用卡。"
          />
          <FAQItem
            question="学生有优惠吗？"
            answer="是的，学生可以享受 50% 的优惠。请联系我们的支持团队验证学生身份。"
          />
        </div>
      </div>

      {/* 优惠券 */}
      <div className="card max-w-md mx-auto">
        <h3 className="text-lg font-semibold mb-4">有优惠券？</h3>
        <div className="flex gap-2">
          <input
            type="text"
            value={couponCode}
            onChange={(e) => setCouponCode(e.target.value)}
            placeholder="输入优惠码"
            className="input-field flex-1"
          />
          <button
            onClick={handleValidateCoupon}
            className="btn-secondary"
          >
            验证
          </button>
        </div>
      </div>
    </div>
  )
}

// 价格卡片组件
function PricingCard({ plan, interval, onSubscribe }) {
  const isPro = plan.name === 'pro'
  const isEnterprise = plan.name === 'enterprise'

  const icon = isEnterprise ? Building2 : isPro ? Zap : Star

  return (
    <div
      className={`card relative ${
        isPro ? 'border-2 border-primary-500 shadow-xl' : ''
      }`}
    >
      {isPro && (
        <div className="absolute -top-4 left-1/2 -translate-x-1/2">
          <span className="bg-primary-500 text-white px-4 py-1 rounded-full text-sm font-semibold">
            最受欢迎
          </span>
        </div>
      )}

      <div className="text-center mb-6">
        <div className="flex justify-center mb-4">
          <div className={`w-12 h-12 rounded-full flex items-center justify-center ${
            isPro ? 'bg-primary-100 dark:bg-primary-900/30 text-primary-600' :
            isEnterprise ? 'bg-purple-100 dark:bg-purple-900/30 text-purple-600' :
            'bg-gray-100 dark:bg-gray-700 text-gray-600'
          }`}>
            {icon({ className: 'w-6 h-6' })}
          </div>
        </div>
        <h3 className="text-xl font-bold text-gray-900 dark:text-white mb-2">
          {plan.display_name}
        </h3>
        <p className="text-gray-600 dark:text-gray-400 text-sm mb-4">
          {plan.description}
        </p>
        <div className="flex items-baseline justify-center gap-1">
          <span className="text-4xl font-bold text-gray-900 dark:text-white">
            ${interval === 'yearly' ? (plan.price * 0.8).toFixed(0) : plan.price}
          </span>
          <span className="text-gray-500">/{interval === 'yearly' ? '年' : '月'}</span>
        </div>
      </div>

      <ul className="space-y-3 mb-8">
        <FeatureItem included={true}>
          {plan.max_skills < 0 ? '无限' : plan.max_skills} 个技能发布
        </FeatureItem>
        <FeatureItem included={true}>
          {plan.max_downloads < 0 ? '无限' : plan.max_downloads} 次下载/月
        </FeatureItem>
        <FeatureItem included={true}>
          {formatStorage(plan.max_storage)} 存储空间
        </FeatureItem>
        <FeatureItem included={plan.has_analytics}>
          数据分析
        </FeatureItem>
        <FeatureItem included={plan.has_priority_support}>
          优先支持
        </FeatureItem>
        {isEnterprise && (
          <>
            <FeatureItem included={true}>
              自定义域名
            </FeatureItem>
            <FeatureItem included={true}>
              SSO 集成
            </FeatureItem>
            <FeatureItem included={true}>
              99.9% SLA 保障
            </FeatureItem>
          </>
        )}
      </ul>

      <button
        onClick={() => onSubscribe(plan.id)}
        className={`btn-primary w-full ${
          isPro ? 'bg-primary-600 hover:bg-primary-700' : ''
        }`}
      >
        {plan.price === 0 ? '免费开始' : '开始试用'}
      </button>
    </div>
  )
}

// 功能项组件
function FeatureItem({ included, children }) {
  return (
    <li className="flex items-center gap-3">
      {included ? (
        <Check className="w-5 h-5 text-green-500 flex-shrink-0" />
      ) : (
        <X className="w-5 h-5 text-gray-300 flex-shrink-0" />
      )}
      <span className={included ? '' : 'text-gray-400 line-through'}>
        {children}
      </span>
    </li>
  )
}

// 功能行组件
function FeatureRow({ feature, free, pro, enterprise }) {
  return (
    <tr className="border-b hover:bg-gray-50 dark:hover:bg-gray-800/50">
      <td className="py-4 px-4 font-medium text-gray-900 dark:text-white">{feature}</td>
      <td className="text-center py-4 px-4">{renderValue(free)}</td>
      <td className="text-center py-4 px-4 bg-primary-50 dark:bg-primary-900/20">
        {renderValue(pro)}
      </td>
      <td className="text-center py-4 px-4">{renderValue(enterprise)}</td>
    </tr>
  )
}

// 渲染值
function renderValue(value) {
  if (typeof value === 'boolean') {
    return value ? (
      <Check className="w-5 h-5 text-green-500 mx-auto" />
    ) : (
      <X className="w-5 h-5 text-gray-300 mx-auto" />
    )
  }
  return value
}

// 格式化存储
function formatStorage(bytes) {
  if (bytes >= 1073741824) {
    return `${(bytes / 1073741824).toFixed(0)} GB`
  }
  if (bytes >= 1048576) {
    return `${(bytes / 1048576).toFixed(0)} MB`
  }
  return `${(bytes / 1024).toFixed(0)} KB`
}
