import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { useThemeStore } from '@/store/theme'
import { useTexts } from '@/composables/useTexts'

export interface NavItem {
  path: string
  name: string
  icon: string
  external?: boolean
  target?: string
  requiresAuth?: boolean
}

/**
 * 导航配置组合式函数
 */
export function useNavigation() {
  const route = useRoute()
  const themeStore = useThemeStore()
  const { $t } = useTexts()

  const baseNavItems = computed<NavItem[]>(() => [
    {
      path: '/',
      name: $t('navigation.home'),
      icon: 'fas fa-home',
    },
    {
      path: '/explore',
      name: $t('navigation.explore'),
      icon: 'fas fa-images',
    },
    {
      path: '/random',
      name: $t('navigation.random'),
      icon: 'fas fa-random',
    },
    {
      path: '/upload',
      name: $t('navigation.upload'),
      icon: 'fas fa-upload',
    },
    {
      path: '/hive',
      name: $t('navigation.hive'),
      icon: 'fas fa-th',
      external: true,
      target: '_blank',
    },
    {
      path: '/docs',
      name: $t('navigation.api'),
      icon: 'fas fa-code',
      external: true,
      target: '_blank',
    },
  ])

  const _authNavItems = computed<NavItem[]>(() => [
    {
      path: '/dashboard',
      name: $t('navigation.dashboard'),
      icon: 'fas fa-tachometer-alt',
      requiresAuth: true,
    },
    {
      path: '/my-files',
      name: $t('navigation.myFiles'),
      icon: 'fas fa-image',
      requiresAuth: true,
    },
    {
      path: '/folders',
      name: $t('navigation.folders'),
      icon: 'fas fa-folder',
      requiresAuth: true,
    },
    {
      path: '/shares',
      name: $t('navigation.shares'),
      icon: 'fas fa-share-alt',
      requiresAuth: true,
    },
    {
      path: '/settings',
      name: $t('navigation.settings'),
      icon: 'fas fa-user-cog',
      requiresAuth: true,
    },
  ])

  const navigationItems = computed(() => baseNavItems.value)

  const isRouteActive = (path: string): boolean => {
    if (path === '/') {
      return route.path === '/'
    }

    return route.path.startsWith(path) || route.path === path
  }

  const getHomePath = () => '/'

  const getPath = (basePath: string) => basePath

  return {
    navigationItems,
    isRouteActive,
    getHomePath,
    getPath,
    siteMode: computed(() => themeStore.getSiteMode()),
  }
}
