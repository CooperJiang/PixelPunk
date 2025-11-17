import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/store/auth'
import { useSettingsStore } from '@/store/settings'
import { useLayoutStore } from '@/store/layout'
import { useTextThemeStore } from '@/store/textTheme'
import { useTexts } from '@/composables/useTexts'
import routes from './routes'
import { StorageUtil } from '@/utils/storage/storage'
import setupApi from '@/api/setup'

const IP_RESTRICTED_KEY = 'ip_restricted'
const USER_DISABLED_KEY = 'user_disabled'

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior(to, from, savedPosition) {
    if (savedPosition) {
      return savedPosition
    }
    return { top: 0, left: 0, behavior: 'instant' }
  },
})

router.beforeEach(async (to, from, next) => {
  const authStore = useAuthStore()
  const layoutStore = useLayoutStore()
  const requiresAuth = to.matched.some((record) => record.meta.requiresAuth)
  const requiresAdmin = to.matched.some((record) => record.meta.requiresAdmin)
  const skipAuth = to.matched.some((record) => record.meta.skipAuth)

  const isIpRestricted = StorageUtil.get(IP_RESTRICTED_KEY) === 'true'
  const isUserDisabled = StorageUtil.get(USER_DISABLED_KEY) === 'true'
  const isAccessRestricted = isIpRestricted || isUserDisabled

  if (isAccessRestricted && to.name !== 'refuse') {
    return next({ name: 'refuse' })
  }

  if (to.name === 'refuse' && !isAccessRestricted) {
    return next({ path: '/' })
  }

  // Only check installation status when accessing setup page
  // For other pages, HTTP interceptor will handle SYSTEM_NOT_INSTALLED (7001) error
  if (to.name === 'setup') {
    try {
      const response = await setupApi.getInstallStatus()
      if (response.data.installed) {
        return next({ path: '/' })
      }
    } catch (_error) {
      // If API call fails, allow access to setup page
    }
  }

  if (!authStore.initialized) {
    authStore.initAuth()
  }

  const isAdmin = [1, 2].includes(authStore.userInfo?.role)

  if (to.path === '/' && layoutStore.isLeftLayout) {
    return next({ path: '/dashboard' })
  }

  if (to.name === 'auth' && authStore.isLoggedIn) {
    const targetPath = layoutStore.isLeftLayout ? '/dashboard' : '/'
    return next({ path: targetPath })
  }

  if (skipAuth) {
    next()
  } else if (requiresAuth && !authStore.isLoggedIn) {
    const redirectQuery: any = { redirect: to.fullPath }
    next({ path: '/auth', query: redirectQuery })
  } else if (requiresAdmin && !isAdmin) {
    next({ path: '/' })
  } else {
    next()
  }
})

router.afterEach(async (to) => {
  const titleKey = to.meta.title as string
  if (titleKey) {
    const settingsStore = useSettingsStore()
    const textThemeStore = useTextThemeStore()
    const siteName = settingsStore.siteName

    if (textThemeStore.isInitializing) {
      const checkInterval = setInterval(() => {
        if (!textThemeStore.isInitializing && textThemeStore.localeData) {
          clearInterval(checkInterval)
          const { $t } = useTexts()
          const title = $t(titleKey)
          document.title = `${title} - ${siteName}`
        }
      }, 50)

      setTimeout(() => clearInterval(checkInterval), 3000)
    } else if (textThemeStore.localeData) {
      const { $t } = useTexts()
      const title = $t(titleKey)
      document.title = `${title} - ${siteName}`
    } else {
      document.title = `${titleKey} - ${siteName}`
    }
  }

  if (typeof window !== 'undefined' && window._hmt) {
    window._hmt.push(['_trackPageview', to.fullPath])
  }

  if (typeof window !== 'undefined') {
    window.scrollTo(0, 0)
    document.documentElement.scrollTop = 0
    document.body.scrollTop = 0
    setTimeout(() => {
      const contentBody = document.querySelector('.content-body')
      if (contentBody) {
        contentBody.scrollTop = 0
      }
      const adminContentWrapper = document.querySelector('.admin-content-wrapper')
      if (adminContentWrapper) {
        adminContentWrapper.scrollTop = 0
      }
    }, 0)
  }
})

export default router
