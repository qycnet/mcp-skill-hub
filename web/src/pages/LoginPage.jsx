import { useState } from 'react'
import { useNavigate, Link } from 'react-router-dom'
import { useMutation } from '@tanstack/react-query'
import { authApi } from '../api/client'
import { useAuthStore } from '../stores/authStore'
import { Mail, Lock, Eye, EyeOff } from 'lucide-react'

export default function LoginPage() {
  const [isLogin, setIsLogin] = useState(true)
  const [showPassword, setShowPassword] = useState(false)
  const [formData, setFormData] = useState({
    username: '',
    email: '',
    password: '',
    confirmPassword: '',
  })

  const navigate = useNavigate()
  const login = useAuthStore((state) => state.login)

  // 登录突变
  const loginMutation = useMutation({
    mutationFn: authApi.login,
    onSuccess: (res) => {
      const { access_token, token_type, expires_in } = res.data
      // 获取用户信息
      authApi.getProfile().then((profileRes) => {
        login(profileRes.data, access_token)
        navigate('/')
      })
    },
  })

  // 注册突变
  const registerMutation = useMutation({
    mutationFn: authApi.register,
    onSuccess: () => {
      // 注册成功后自动登录
      loginMutation.mutate({
        username: formData.username,
        password: formData.password,
      })
    },
  })

  const handleSubmit = (e) => {
    e.preventDefault()

    if (isLogin) {
      loginMutation.mutate({
        username: formData.username,
        password: formData.password,
      })
    } else {
      if (formData.password !== formData.confirmPassword) {
        alert('两次输入的密码不一致')
        return
      }
      registerMutation.mutate({
        username: formData.username,
        email: formData.email,
        password: formData.password,
      })
    }
  }

  const handleChange = (e) => {
    setFormData({
      ...formData,
      [e.target.name]: e.target.value,
    })
  }

  return (
    <div className="max-w-md mx-auto">
      <div className="card">
        <div className="text-center mb-8">
          <h1 className="text-2xl font-bold text-gray-900 dark:text-white">
            {isLogin ? '欢迎回来' : '创建账户'}
          </h1>
          <p className="text-gray-600 dark:text-gray-400 mt-2">
            {isLogin ? '登录到你的账户' : '注册一个新账户'}
          </p>
        </div>

        <form onSubmit={handleSubmit} className="space-y-6">
          {!isLogin && (
            <>
              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                  用户名
                </label>
                <div className="relative">
                  <input
                    type="text"
                    name="username"
                    value={formData.username}
                    onChange={handleChange}
                    required={!isLogin}
                    className="input-field pl-10"
                    placeholder="请输入用户名"
                  />
                  <Mail className="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-400" />
                </div>
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                  邮箱
                </label>
                <div className="relative">
                  <input
                    type="email"
                    name="email"
                    value={formData.email}
                    onChange={handleChange}
                    required={!isLogin}
                    className="input-field pl-10"
                    placeholder="请输入邮箱"
                  />
                  <Mail className="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-400" />
                </div>
              </div>
            </>
          )}

          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
              密码
            </label>
            <div className="relative">
              <input
                type={showPassword ? 'text' : 'password'}
                name="password"
                value={formData.password}
                onChange={handleChange}
                required
                className="input-field pl-10 pr-10"
                placeholder="请输入密码"
              />
              <Lock className="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-400" />
              <button
                type="button"
                onClick={() => setShowPassword(!showPassword)}
                className="absolute right-3 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-600"
              >
                {showPassword ? <EyeOff className="w-5 h-5" /> : <Eye className="w-5 h-5" />}
              </button>
            </div>
          </div>

          {!isLogin && (
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                确认密码
              </label>
              <div className="relative">
                <input
                  type={showPassword ? 'text' : 'password'}
                  name="confirmPassword"
                  value={formData.confirmPassword}
                  onChange={handleChange}
                  required={!isLogin}
                  className="input-field pl-10"
                  placeholder="请再次输入密码"
                />
                <Lock className="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-400" />
              </div>
            </div>
          )}

          <button
            type="submit"
            disabled={loginMutation.isPending || registerMutation.isPending}
            className="btn-primary w-full py-3"
          >
            {(loginMutation.isPending || registerMutation.isPending) ? '处理中...' : (isLogin ? '登录' : '注册')}
          </button>
        </form>

        {/* 切换登录/注册 */}
        <div className="mt-6 text-center">
          <p className="text-gray-600 dark:text-gray-400">
            {isLogin ? '还没有账户？' : '已有账户？'}
            <button
              onClick={() => setIsLogin(!isLogin)}
              className="text-primary-600 hover:text-primary-700 font-medium ml-1"
            >
              {isLogin ? '立即注册' : '返回登录'}
            </button>
          </p>
        </div>
      </div>
    </div>
  )
}
