import { create } from 'zustand'
import { persist } from 'zustand/middleware'

export const useAuthStore = create(
  persist(
    (set, get) => ({
      // 状态
      user: null,
      token: null,
      isAuthenticated: false,

      // 登录
      login: (userData, token) => {
        set({
          user: userData,
          token,
          isAuthenticated: true,
        })
      },

      // 登出
      logout: () => {
        set({
          user: null,
          token: null,
          isAuthenticated: false,
        })
      },

      // 更新用户信息
      updateUser: (userData) => {
        set({ user: { ...get().user, ...userData } })
      },

      // 检查是否过期
      isTokenExpired: () => {
        // TODO: 实现令牌过期检查
        return false
      },
    }),
    {
      name: 'mcp-auth', // localStorage 键名
      partialize: (state) => ({
        user: state.user,
        token: state.token,
        isAuthenticated: state.isAuthenticated,
      }),
    }
  )
)
