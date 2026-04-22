import axios from 'axios'
import { useAuthStore } from '../stores/authStore'

// 创建 Axios 实例
const apiClient = axios.create({
  baseURL: '/api/v1',
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json',
  },
})

// 请求拦截器 - 添加认证头
apiClient.interceptors.request.use(
  (config) => {
    const token = useAuthStore.getState().token
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// 响应拦截器 - 处理错误
apiClient.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      // 令牌过期，清除认证状态
      useAuthStore.getState().logout()
      window.location.href = '/login'
    }
    return Promise.reject(error)
  }
)

// API 方法

export const skillsApi = {
  // 获取技能列表
  list: (params) => apiClient.get('/skills', { params }),

  // 搜索技能
  search: (query, params) => apiClient.get('/search', { params: { q: query, ...params } }),

  // 获取技能详情
  getById: (id) => apiClient.get(`/skills/${id}`),

  // 发布技能
  publish: (data) => apiClient.post('/skills', data),

  // 更新技能
  update: (id, data) => apiClient.put(`/skills/${id}`, data),

  // 删除技能
  delete: (id) => apiClient.delete(`/skills/${id}`),

  // 评分
  rate: (id, rating, comment) => apiClient.post(`/skills/${id}/rate`, { rating, comment }),

  // 下载技能
  download: (id) => apiClient.get(`/skills/${id}/download`),

  // 获取分类
  categories: () => apiClient.get('/categories'),
}

export const authApi = {
  // 注册
  register: (data) => apiClient.post('/auth/register', data),

  // 登录
  login: (data) => apiClient.post('/auth/login', data),

  // 刷新令牌
  refreshToken: (data) => apiClient.post('/auth/refresh', data),

  // 登出
  logout: () => apiClient.post('/auth/logout'),

  // 获取用户资料
  getProfile: () => apiClient.get('/user/profile'),

  // 更新用户资料
  updateProfile: (data) => apiClient.put('/user/profile', data),

  // 创建 API Key
  createApiKey: (data) => apiClient.post('/user/api-keys', data),

  // 撤销 API Key
  revokeApiKey: (id) => apiClient.delete(`/user/api-keys/${id}`),

  // 列出 API Keys
  listApiKeys: () => apiClient.get('/user/api-keys'),
}

export const adminApi = {
  // 管理技能列表
  listSkills: (params) => apiClient.get('/admin/skills', { params }),

  // 批准技能
  approveSkill: (id) => apiClient.put(`/admin/skills/${id}/approve`),

  // 拒绝技能
  rejectSkill: (id, reason) => apiClient.put(`/admin/skills/${id}/reject`, { reason }),

  // 用户列表
  listUsers: (params) => apiClient.get('/admin/users', { params }),

  // 分析数据
  analytics: () => apiClient.get('/admin/analytics'),
}

export default apiClient
