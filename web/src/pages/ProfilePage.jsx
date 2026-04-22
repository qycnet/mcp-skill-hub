import { useState } from 'react'
import { useQuery, useMutation } from '@tanstack/react-query'
import { authApi, skillsApi } from '../api/client'
import { useAuthStore } from '../stores/authStore'
import { 
  User, Mail, Key, Package, Star, Settings,
  Plus, Trash2, Copy, Check, ExternalLink
} from 'lucide-react'

export default function ProfilePage() {
  const [activeTab, setActiveTab] = useState('profile')
  const { user, token, updateUser, logout } = useAuthStore()
  const [editMode, setEditMode] = useState(false)
  const [formData, setFormData] = useState({
    avatar: user?.avatar || '',
    bio: user?.bio || '',
  })

  // 获取用户资料
  const { data: profile, refetch } = useQuery({
    queryKey: ['profile'],
    queryFn: authApi.getProfile,
    select: (res) => res.data,
  })

  // 获取用户发布的技能
  const { data: mySkills } = useQuery({
    queryKey: ['my-skills'],
    queryFn: () => skillsApi.list({ author_id: user?.id }),
    select: (res) => res.data.skills,
  })

  // 获取 API Keys
  const { data: apiKeys, refetch: refetchKeys } = useQuery({
    queryKey: ['api-keys'],
    queryFn: authApi.listApiKeys,
    select: (res) => res.data.keys,
  })

  // 更新资料突变
  const updateProfileMutation = useMutation({
    mutationFn: authApi.updateProfile,
    onSuccess: () => {
      refetch()
      setEditMode(false)
      alert('资料更新成功')
    },
  })

  // 创建 API Key 突变
  const createKeyMutation = useMutation({
    mutationFn: authApi.createApiKey,
    onSuccess: () => {
      refetchKeys()
      alert('API Key 创建成功')
    },
  })

  // 撤销 API Key 突变
  const revokeKeyMutation = useMutation({
    mutationFn: authApi.revokeApiKey,
    onSuccess: () => {
      refetchKeys()
      alert('API Key 已撤销')
    },
  })

  const handleUpdateProfile = (e) => {
    e.preventDefault()
    updateProfileMutation.mutate(formData)
  }

  const handleCreateApiKey = () => {
    const name = prompt('请输入 API Key 名称：')
    if (name) {
      createKeyMutation.mutate({ name, description: '自动创建' })
    }
  }

  const handleCopyToken = () => {
    navigator.clipboard.writeText(token)
    alert('Token 已复制到剪贴板')
  }

  const tabs = [
    { id: 'profile', label: '个人资料', icon: User },
    { id: 'skills', label: '我的技能', icon: Package },
    { id: 'api-keys', label: 'API Keys', icon: Key },
    { id: 'settings', label: '设置', icon: Settings },
  ]

  return (
    <div className="max-w-4xl mx-auto space-y-6">
      {/* 页面标题 */}
      <div className="flex items-center justify-between">
        <h1 className="text-3xl font-bold text-gray-900 dark:text-white">个人中心</h1>
        <button
          onClick={logout}
          className="text-red-600 hover:text-red-700 text-sm"
        >
          退出登录
        </button>
      </div>

      {/* 标签页 */}
      <div className="border-b border-gray-200 dark:border-gray-700">
        <nav className="flex space-x-8">
          {tabs.map((tab) => (
            <button
              key={tab.id}
              onClick={() => setActiveTab(tab.id)}
              className={`flex items-center py-4 px-1 border-b-2 font-medium text-sm transition-colors ${
                activeTab === tab.id
                  ? 'border-primary-500 text-primary-600 dark:text-primary-400'
                  : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
              }`}
            >
              <tab.icon className="w-4 h-4 mr-2" />
              {tab.label}
            </button>
          ))}
        </nav>
      </div>

      {/* 个人资料 */}
      {activeTab === 'profile' && (
        <div className="card">
          <div className="flex justify-between items-center mb-6">
            <h2 className="text-xl font-semibold flex items-center">
              <User className="w-5 h-5 mr-2" />
              基本信息
            </h2>
            <button
              onClick={() => setEditMode(!editMode)}
              className="text-primary-600 hover:text-primary-700 text-sm"
            >
              {editMode ? '取消' : '编辑'}
            </button>
          </div>

          <div className="space-y-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                用户名
              </label>
              <input
                type="text"
                value={profile?.username || ''}
                disabled
                className="input-field bg-gray-100 dark:bg-gray-700 cursor-not-allowed"
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                邮箱
              </label>
              <div className="flex items-center">
                <Mail className="w-5 h-5 text-gray-400 mr-2" />
                <input
                  type="email"
                  value={profile?.email || ''}
                  disabled
                  className="input-field bg-gray-100 dark:bg-gray-700 cursor-not-allowed flex-1"
                />
              </div>
            </div>

            {editMode && (
              <>
                <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                    头像 URL
                  </label>
                  <input
                    type="url"
                    value={formData.avatar}
                    onChange={(e) => setFormData({ ...formData, avatar: e.target.value })}
                    className="input-field"
                    placeholder="https://..."
                  />
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                    个人简介
                  </label>
                  <textarea
                    value={formData.bio}
                    onChange={(e) => setFormData({ ...formData, bio: e.target.value })}
                    rows={4}
                    className="input-field"
                    placeholder="介绍一下自己..."
                  />
                </div>

                <button
                  onClick={handleUpdateProfile}
                  disabled={updateProfileMutation.isPending}
                  className="btn-primary"
                >
                  {updateProfileMutation.isPending ? '保存中...' : '保存更改'}
                </button>
              </>
            )}

            {!editMode && profile?.bio && (
              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                  个人简介
                </label>
                <p className="text-gray-700 dark:text-gray-300">{profile.bio}</p>
              </div>
            )}
          </div>
        </div>
      )}

      {/* 我的技能 */}
      {activeTab === 'skills' && (
        <div className="space-y-4">
          <div className="flex justify-between items-center">
            <h2 className="text-xl font-semibold flex items-center">
              <Package className="w-5 h-5 mr-2" />
              我的技能
            </h2>
            <a href="/publish" className="btn-primary flex items-center">
              <Plus className="w-4 h-4 mr-2" />
              发布技能
            </a>
          </div>

          {mySkills?.length === 0 ? (
            <div className="card text-center py-12">
              <Package className="w-16 h-16 text-gray-300 mx-auto mb-4" />
              <p className="text-gray-600 dark:text-gray-400 mb-4">还没有发布技能</p>
              <a href="/publish" className="btn-primary">发布第一个技能</a>
            </div>
          ) : (
            <div className="grid gap-4">
              {mySkills?.map((skill) => (
                <div key={skill.id} className="card">
                  <div className="flex justify-between items-start">
                    <div>
                      <h3 className="font-semibold text-lg">{skill.display_name}</h3>
                      <p className="text-gray-600 dark:text-gray-400 text-sm">{skill.description}</p>
                      <div className="flex items-center gap-4 mt-2 text-sm text-gray-500">
                        <span className="flex items-center">
                          <Star className="w-4 h-4 mr-1 text-yellow-500" />
                          {skill.rating.toFixed(1)}
                        </span>
                        <span className="flex items-center">
                          <Download className="w-4 h-4 mr-1" />
                          {skill.downloads}
                        </span>
                      </div>
                    </div>
                    <div className="flex gap-2">
                      <a
                        href={`/skills/${skill.id}`}
                        className="btn-secondary text-sm"
                      >
                        查看
                      </a>
                      <button className="btn-secondary text-sm">
                        编辑
                      </button>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>
      )}

      {/* API Keys */}
      {activeTab === 'api-keys' && (
        <div className="space-y-4">
          <div className="flex justify-between items-center">
            <h2 className="text-xl font-semibold flex items-center">
              <Key className="w-5 h-5 mr-2" />
              API Keys
            </h2>
            <button onClick={handleCreateApiKey} className="btn-primary flex items-center">
              <Plus className="w-4 h-4 mr-2" />
              创建新 Key
            </button>
          </div>

          {/* 当前 Token */}
          <div className="card bg-yellow-50 dark:bg-yellow-900/20 border-yellow-200 dark:border-yellow-800">
            <h3 className="font-semibold mb-2">当前会话 Token</h3>
            <div className="flex items-center gap-2">
              <code className="flex-1 bg-white dark:bg-gray-800 px-3 py-2 rounded text-sm font-mono overflow-hidden text-ellipsis">
                {token?.substring(0, 50)}...
              </code>
              <button onClick={handleCopyToken} className="btn-secondary">
                <Copy className="w-4 h-4" />
              </button>
            </div>
          </div>

          {/* API Key 列表 */}
          {apiKeys?.length === 0 ? (
            <div className="card text-center py-12">
              <Key className="w-16 h-16 text-gray-300 mx-auto mb-4" />
              <p className="text-gray-600 dark:text-gray-400">还没有创建 API Key</p>
            </div>
          ) : (
            <div className="space-y-3">
              {apiKeys?.map((key) => (
                <div key={key.id} className="card">
                  <div className="flex justify-between items-center">
                    <div>
                      <div className="font-semibold">{key.name}</div>
                      {key.description && (
                        <p className="text-sm text-gray-600 dark:text-gray-400">{key.description}</p>
                      )}
                      <div className="text-xs text-gray-500 mt-1">
                        创建：{new Date(key.created_at).toLocaleDateString('zh-CN')}
                        {key.last_used_at && ` · 最后使用：${new Date(key.last_used_at).toLocaleDateString('zh-CN')}`}
                      </div>
                    </div>
                    <button
                      onClick={() => {
                        if (confirm('确定要撤销这个 API Key 吗？')) {
                          revokeKeyMutation.mutate(key.id)
                        }
                      }}
                      className="text-red-600 hover:text-red-700 p-2"
                    >
                      <Trash2 className="w-5 h-5" />
                    </button>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>
      )}

      {/* 设置 */}
      {activeTab === 'settings' && (
        <div className="card">
          <h2 className="text-xl font-semibold mb-6 flex items-center">
            <Settings className="w-5 h-5 mr-2" />
            账户设置
          </h2>

          <div className="space-y-6">
            <div>
              <h3 className="font-semibold mb-2">修改密码</h3>
              <p className="text-gray-600 dark:text-gray-400 text-sm mb-4">
                如需修改密码，请使用密码重置功能
              </p>
              <button className="btn-secondary">
                发送重置邮件
              </button>
            </div>

            <div className="border-t border-gray-200 dark:border-gray-700 pt-6">
              <h3 className="font-semibold mb-2">账户状态</h3>
              <div className="flex items-center gap-2">
                <span className="w-3 h-3 bg-green-500 rounded-full"></span>
                <span className="text-gray-700 dark:text-gray-300">账户正常</span>
              </div>
              <p className="text-sm text-gray-500 mt-2">
                注册时间：{profile?.created_at ? new Date(profile.created_at).toLocaleDateString('zh-CN') : '未知'}
              </p>
            </div>

            <div className="border-t border-gray-200 dark:border-gray-700 pt-6">
              <h3 className="font-semibold mb-2 text-red-600">危险区域</h3>
              <p className="text-gray-600 dark:text-gray-400 text-sm mb-4">
                注销账户将删除所有数据，此操作不可恢复
              </p>
              <button className="bg-red-600 hover:bg-red-700 text-white font-medium py-2 px-4 rounded-lg">
                注销账户
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}

// 下载图标组件
function Download({ className }) {
  return (
    <svg className={className} fill="none" viewBox="0 0 24 24" stroke="currentColor">
      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
    </svg>
  )
}
