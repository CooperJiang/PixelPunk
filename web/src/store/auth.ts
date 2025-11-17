import { defineStore } from 'pinia'
import { userApi } from '@/api'
import type { UserInfo, UserLoginRequest, UserRegisterRequest } from '@/api/types/index'
import { StorageUtil } from '@/utils/storage/storage'
import { TOKEN_EXPIRES, TOKEN_KEY, USER_INFO_KEY } from '@/constants'

interface AuthState {
  user: UserInfo | null
  token: string | null
  isAuthenticated: boolean
  avatarUrl: string | null
  initialized: boolean
}

export const useAuthStore = defineStore('auth', {
  state: (): AuthState => {
    const storedUser = typeof window !== 'undefined' ? StorageUtil.get<UserInfo>(USER_INFO_KEY) : null
    const storedToken = typeof window !== 'undefined' ? StorageUtil.get<string>(TOKEN_KEY) : null

    return {
      user: storedUser,
      token: storedToken,
      isAuthenticated: Boolean(storedToken && storedUser),
      avatarUrl: storedUser?.avatarFullPath || storedUser?.avatar || null,
      initialized: false,
    }
  },

  getters: {
    userInfo: (state) => state.user,
    isLoggedIn: (state) => state.isAuthenticated && Boolean(state.token) && Boolean(state.user),
    userAvatar: (state) => state.avatarUrl,
    isAdmin: (state) => [1, 2].includes(state.user?.role),
  },

  actions: {
    initAuth() {
      if (this.initialized || typeof window === 'undefined') {
        return
      }

      this.checkAuth()
      this.refreshCookie()
      this.initialized = true
    },

    async login(loginData: UserLoginRequest) {
      const result = await userApi.login(loginData)
      if (result.success) {
        const { userInfo, token } = result.data
        this.setUserInfo(userInfo)
        this.setToken(token)
        return result.data
      }
      // 优先使用后端返回的错误消息，fallback使用默认消息
      // 翻译: store.auth.errors.loginFailed
      const errorMessage = result.message || 'Login failed'
      const error = new Error(errorMessage)
      ;(error as any).response = { data: result }
      throw error
    },

    async register(registerData: UserRegisterRequest) {
      const result = await userApi.register(registerData)
      if (result.success) {
        return result.data
      }
      // 优先使用后端返回的错误消息，fallback使用默认消息
      // 翻译: store.auth.errors.registerFailed
      throw new Error(result.message || 'Registration failed')
    },

    setUserInfo(user: UserInfo) {
      this.user = user
      this.isAuthenticated = true

      this.avatarUrl = user.avatarFullPath || user.avatar

      StorageUtil.set<UserInfo>(USER_INFO_KEY, user, TOKEN_EXPIRES)
    },

    setToken(token: string) {
      this.token = token
      StorageUtil.set<string>(TOKEN_KEY, token, TOKEN_EXPIRES)

      const tokenExpireTime = this.getTokenExpireTime(token)
      if (tokenExpireTime > 0) {
        document.cookie = `token=${token}; path=/; max-age=${tokenExpireTime}; SameSite=Lax`
      } else {
        const maxAgeSeconds = TOKEN_EXPIRES * 3600
        document.cookie = `token=${token}; path=/; max-age=${maxAgeSeconds}; SameSite=Lax`
      }
    },

    logout() {
      this.user = null
      this.token = null
      this.isAuthenticated = false
      this.avatarUrl = null
      this.initialized = false
      StorageUtil.remove(TOKEN_KEY)
      StorageUtil.remove(USER_INFO_KEY)
      StorageUtil.remove('rememberedLogin')

      document.cookie = 'token=; path=/; expires=Thu, 01 Jan 1970 00:00:01 GMT;'

      const currentPath = window.location.pathname
      if (currentPath !== '/auth' && currentPath !== '/login') {
        window.location.replace('/auth')
      }
    },

    checkAuth() {
      if (typeof window === 'undefined') {
        return
      }

      const token = StorageUtil.get<string>(TOKEN_KEY)
      const userInfo = StorageUtil.get<UserInfo>(USER_INFO_KEY)

      if (token && userInfo) {
        this.token = token
        this.user = userInfo
        this.isAuthenticated = true
        this.avatarUrl = userInfo.avatarFullPath || userInfo.avatar
      } else {
        this.isAuthenticated = false
        this.user = null
        this.token = null
        this.avatarUrl = null
        if (!token || !userInfo) {
          StorageUtil.remove(TOKEN_KEY)
          StorageUtil.remove(USER_INFO_KEY)
        }
      }
    },

    updateUserAvatar(avatarUrl: string, fullUrl?: string) {
      if (this.user) {
        this.user.avatar = avatarUrl

        if (fullUrl) {
          this.user.avatarFullPath = fullUrl
          this.avatarUrl = fullUrl
        } else {
          this.avatarUrl = avatarUrl
        }

        StorageUtil.set(USER_INFO_KEY, this.user, TOKEN_EXPIRES)
      }
    },

    updateUserInfo(newUserInfo: Partial<UserInfo>) {
      if (!this.user) {
        return
      }
      this.user = { ...this.user, ...newUserInfo }

      if (newUserInfo.avatarFullPath) {
        this.avatarUrl = newUserInfo.avatarFullPath
      } else if (newUserInfo.avatar) {
        this.avatarUrl = newUserInfo.avatar
      }

      StorageUtil.set(USER_INFO_KEY, this.user, TOKEN_EXPIRES)
    },

    refreshCookie() {
      if (typeof window === 'undefined') {
        return
      }

      if (this.token && this.isAuthenticated) {
        const tokenExpireTime = this.getTokenExpireTime(this.token)
        if (tokenExpireTime > 0) {
          document.cookie = `token=${this.token}; path=/; max-age=${tokenExpireTime}; SameSite=Lax`
        } else {
          const maxAgeSeconds = TOKEN_EXPIRES * 3600
          document.cookie = `token=${this.token}; path=/; max-age=${maxAgeSeconds}; SameSite=Lax`
        }
      }
    },

    getTokenExpireTime(token: string): number {
      try {
        const parts = token.split('.')
        if (parts.length !== 3) {
          return 0
        }

        const payload = parts[1]
        const paddedPayload = payload + '=='.substring(0, (4 - (payload.length % 4)) % 4)
        const decodedPayload = atob(paddedPayload.replace(/-/g, '+').replace(/_/g, '/'))
        const claims = JSON.parse(decodedPayload)

        if (claims.exp && typeof claims.exp === 'number') {
          const now = Math.floor(Date.now() / 1000)
          const maxAge = claims.exp - now
          return maxAge > 0 ? maxAge : 0
        }

        return 0
      } catch {
        return 0
      }
    },
  },
})
