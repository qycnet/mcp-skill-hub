import { useState } from 'react'
import { useQuery, useMutation } from '@tanstack/react-query'
import { Link, useNavigate } from 'react-router-dom'
import { paymentApi } from '../api/client'
import { CreditCard, CheckCircle, AlertCircle, Download, Calendar, RefreshCw } from 'lucide-react'

export default function SubscriptionPage() {
  const navigate = useNavigate()
  const [showCancelModal, setShowCancelModal] = useState(false)

  // 获取活跃订阅
  const { data: subscription, refetch } = useQuery({
    queryKey: ['subscription'],
    queryFn: paymentApi.getActiveSubscription,
    select: (res) => res.data.subscription,
  })

  // 获取支付历史
  const { data: payments } = useQuery({
    queryKey: ['payments'],
    queryFn: paymentApi.getPaymentHistory,
    select: (res) => res.data.payments,
  })

  // 获取使用量
  const { data: usage } = useQuery({
    queryKey: ['usage'],
    queryFn: paymentApi.getUsage,
    select: (res) => res.data.usage,
  })

  // 取消订阅突变
  const cancelMutation = useMutation({
    mutationFn: paymentApi.cancelSubscription,
    onSuccess: () => {
      refetch()
      setShowCancelModal(false)
      alert('订阅已取消，将在当前周期结束后失效')
    },
  })

  // 恢复订阅突变
  const resumeMutation = useMutation({
    mutationFn: paymentApi.resumeSubscription,
    onSuccess: () => {
      refetch()
      alert('订阅已恢复')
    },
  })

  const handleCancel = () => {
    if (subscription?.id) {
      cancelMutation.mutate(subscription.id)
    }
  }

  const handleResume = () => {
    if (subscription?.id) {
      resumeMutation.mutate(subscription.id)
    }
  }

  const handleUpgrade = (planId) => {
    navigate(`/pricing?plan=${planId}`)
  }

  const handleDownloadInvoice = (invoiceUrl) => {
    window.open(invoiceUrl, '_blank')
  }

  return (
    <div className="max-w-4xl mx-auto space-y-8">
      {/* 页面标题 */}
      <div>
        <h1 className="text-3xl font-bold text-gray-900 dark:text-white mb-2">
          订阅管理
        </h1>
        <p className="text-gray-600 dark:text-gray-400">
          管理你的订阅计划、支付历史和使用情况
        </p>
      </div>

      {/* 当前订阅状态 */}
      <div className="card">
        <h2 className="text-xl font-semibold mb-6">当前订阅</h2>

        {subscription ? (
          <div className="space-y-6">
            {/* 订阅信息 */}
            <div className="flex items-start justify-between">
              <div>
                <div className="flex items-center gap-3 mb-2">
                  <h3 className="text-2xl font-bold text-gray-900 dark:text-white">
                    {subscription.plan.display_name}
                  </h3>
                  <span className={`px-3 py-1 rounded-full text-sm font-semibold ${
                    subscription.status === 'active'
                      ? 'bg-green-100 dark:bg-green-900/30 text-green-600 dark:text-green-400'
                      : 'bg-yellow-100 dark:bg-yellow-900/30 text-yellow-600 dark:text-yellow-400'
                  }`}>
                    {subscription.status === 'active' ? '活跃' : '已取消'}
                  </span>
                </div>
                <p className="text-gray-600 dark:text-gray-400">
                  {subscription.plan.description}
                </p>
              </div>
              <div className="text-right">
                <div className="text-2xl font-bold text-gray-900 dark:text-white">
                  ${subscription.plan.price}
                  <span className="text-sm font-normal text-gray-500">/月</span>
                </div>
                <div className="text-sm text-gray-500">
                  下次账单：{new Date(subscription.current_period_end).toLocaleDateString('zh-CN')}
                </div>
              </div>
            </div>

            {/* 订阅详情 */}
            <div className="grid grid-cols-2 gap-4 bg-gray-50 dark:bg-gray-800/50 rounded-lg p-4">
              <div>
                <div className="text-sm text-gray-500 mb-1">计费周期</div>
                <div className="font-medium">
                  {new Date(subscription.current_period_start).toLocaleDateString('zh-CN')} - 
                  {new Date(subscription.current_period_end).toLocaleDateString('zh-CN')}
                </div>
              </div>
              <div>
                <div className="text-sm text-gray-500 mb-1">支付方式</div>
                <div className="font-medium flex items-center gap-2">
                  <CreditCard className="w-4 h-4" />
                  {subscription.provider === 'stripe' ? '信用卡' : subscription.provider}
                </div>
              </div>
            </div>

            {/* 操作按钮 */}
            <div className="flex gap-4">
              {subscription.status === 'active' ? (
                <>
                  <button
                    onClick={() => setShowCancelModal(true)}
                    className="btn-secondary text-red-600 hover:text-red-700"
                  >
                    取消订阅
                  </button>
                  <button
                    onClick={() => handleUpgrade(subscription.plan.id)}
                    className="btn-primary"
                  >
                    更改计划
                  </button>
                </>
              ) : (
                <button
                  onClick={handleResume}
                  className="btn-primary flex items-center"
                >
                  <RefreshCw className="w-4 h-4 mr-2" />
                  恢复订阅
                </button>
              )}
            </div>
          </div>
        ) : (
          <div className="text-center py-12">
            <div className="w-16 h-16 bg-gray-100 dark:bg-gray-700 rounded-full flex items-center justify-center mx-auto mb-4">
              <CreditCard className="w-8 h-8 text-gray-400" />
            </div>
            <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-2">
              暂无活跃订阅
            </h3>
            <p className="text-gray-600 dark:text-gray-400 mb-6">
              订阅专业版或企业版，解锁更多功能
            </p>
            <Link to="/pricing" className="btn-primary">
              查看价格计划
            </Link>
          </div>
        )}
      </div>

      {/* 使用情况 */}
      {subscription && (
        <div className="card">
          <h2 className="text-xl font-semibold mb-6">使用情况</h2>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
            <UsageCard
              icon={CheckCircle}
              title="技能发布"
              current={usage?.skills || 0}
              limit={usage?.limits?.skills || '∞'}
              unit="个"
            />
            <UsageCard
              icon={Download}
              title="下载量"
              current={usage?.downloads || 0}
              limit={usage?.limits?.downloads || '∞'}
              unit="次"
              period="/月"
            />
            <UsageCard
              icon={Calendar}
              title="存储空间"
              current={formatStorage(usage?.storage_used || 0)}
              limit={formatStorage(usage?.limits?.storage || 0)}
              period=""
            />
          </div>
        </div>
      )}

      {/* 支付历史 */}
      <div className="card">
        <h2 className="text-xl font-semibold mb-6">支付历史</h2>
        {payments && payments.length > 0 ? (
          <div className="space-y-4">
            {payments.map((payment) => (
              <div
                key={payment.id}
                className="flex items-center justify-between p-4 bg-gray-50 dark:bg-gray-800/50 rounded-lg"
              >
                <div className="flex items-center gap-4">
                  <div className={`w-10 h-10 rounded-full flex items-center justify-center ${
                    payment.status === 'completed'
                      ? 'bg-green-100 dark:bg-green-900/30 text-green-600'
                      : 'bg-red-100 dark:bg-red-900/30 text-red-600'
                  }`}>
                    {payment.status === 'completed' ? (
                      <CheckCircle className="w-5 h-5" />
                    ) : (
                      <AlertCircle className="w-5 h-5" />
                    )}
                  </div>
                  <div>
                    <div className="font-semibold">
                      ${payment.amount} - {payment.type === 'subscription' ? '订阅费用' : '一次性支付'}
                    </div>
                    <div className="text-sm text-gray-500">
                      {new Date(payment.created_at).toLocaleDateString('zh-CN')}
                    </div>
                  </div>
                </div>
                <div className="flex items-center gap-4">
                  <span className={`px-3 py-1 rounded-full text-sm ${
                    payment.status === 'completed'
                      ? 'bg-green-100 dark:bg-green-900/30 text-green-600'
                      : 'bg-red-100 dark:bg-red-900/30 text-red-600'
                  }`}>
                    {payment.status === 'completed' ? '已完成' : payment.status}
                  </span>
                  {payment.invoice_url && (
                    <button
                      onClick={() => handleDownloadInvoice(payment.invoice_url)}
                      className="btn-secondary text-sm"
                    >
                      下载发票
                    </button>
                  )}
                </div>
              </div>
            ))}
          </div>
        ) : (
          <div className="text-center py-8 text-gray-500">
            暂无支付记录
          </div>
        )}
      </div>

      {/* 取消确认对话框 */}
      {showCancelModal && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
          <div className="bg-white dark:bg-gray-800 rounded-lg p-6 max-w-md mx-4">
            <h3 className="text-xl font-bold mb-4">确认取消订阅？</h3>
            <p className="text-gray-600 dark:text-gray-400 mb-6">
              取消后，你的订阅将在当前计费周期结束（{new Date(subscription?.current_period_end).toLocaleDateString('zh-CN')}）时失效。
              你仍然可以继续使用到那时。
            </p>
            <div className="flex gap-4">
              <button
                onClick={() => setShowCancelModal(false)}
                className="btn-secondary flex-1"
              >
                保留订阅
              </button>
              <button
                onClick={handleCancel}
                disabled={cancelMutation.isPending}
                className="btn-primary bg-red-600 hover:bg-red-700 flex-1"
              >
                {cancelMutation.isPending ? '处理中...' : '确认取消'}
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}

// 使用卡片组件
function UsageCard({ icon: Icon, title, current, limit, unit, period }) {
  const percentage = limit === '∞' ? 100 : Math.min((current / limit) * 100, 100)

  return (
    <div>
      <div className="flex items-center gap-3 mb-4">
        <div className="w-10 h-10 bg-primary-100 dark:bg-primary-900/30 rounded-lg flex items-center justify-center">
          <Icon className="w-5 h-5 text-primary-600 dark:text-primary-400" />
        </div>
        <div>
          <div className="text-sm text-gray-500">{title}</div>
          <div className="text-lg font-semibold">
            {current}{unit}{period && <span className="text-sm font-normal text-gray-500">{period}</span>}
            {limit !== '∞' && <span className="text-sm text-gray-500"> / {limit}{unit}</span>}
          </div>
        </div>
      </div>
      <div className="w-full bg-gray-200 dark:bg-gray-700 rounded-full h-2">
        <div
          className={`h-2 rounded-full ${
            percentage > 90 ? 'bg-red-500' : percentage > 70 ? 'bg-yellow-500' : 'bg-green-500'
          }`}
          style={{ width: `${percentage}%` }}
        />
      </div>
    </div>
  )
}

// 格式化存储
function formatStorage(bytes) {
  if (bytes >= 1073741824) {
    return `${(bytes / 1073741824).toFixed(1)} GB`
  }
  if (bytes >= 1048576) {
    return `${(bytes / 1048576).toFixed(0)} MB`
  }
  return `${(bytes / 1024).toFixed(0)} KB`
}
