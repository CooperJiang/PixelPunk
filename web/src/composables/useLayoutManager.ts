import { computed, onMounted, onUnmounted } from 'vue'
import { useRoute } from 'vue-router'
import { useLayoutStore } from '@/store/layout'

/**
 * 布局管理组合式函数
 * 统一管理布局相关的逻辑和状态
 */
export function useLayoutManager() {
  const route = useRoute()
  const layoutStore = useLayoutStore()

  const isTopLayout = computed(() => layoutStore.isTopLayout)
  const isLeftLayout = computed(() => layoutStore.isLeftLayout)
  const layoutClasses = computed(() => layoutStore.layoutClasses)

  const isFullWidthPage = computed(() => ['/', '/home', '/random', '/explore'].includes(route.path))
  const isHomePage = computed(() => route.path === '/')
  const isDocsPage = computed(() => route.path.startsWith('/docs'))

  const mainContentStyle = computed(() => {
    const baseStyle = {
      paddingTop: getPaddingTop(),
    }

    if (isLeftLayout.value) {
      return {
        ...baseStyle,
        ...layoutStore.containerStyle,
      }
    }

    return baseStyle
  })

  function getPaddingTop(): string {
    if (isLeftLayout.value) {
      if (route.path === '/') return '64px'
      if (route.path === '/random') return '64px'
      return '60px'
    }

    if (route.path === '/') return '20px'
    if (route.path === '/random') return '0px'
    return '4.5rem'
  }

  function handleResize() {
    layoutStore.handleResize(window.innerWidth)
  }

  onMounted(() => {
    window.addEventListener('resize', handleResize)
    handleResize() // 初始化
  })

  onUnmounted(() => {
    window.removeEventListener('resize', handleResize)
  })

  return {
    isTopLayout,
    isLeftLayout,
    layoutClasses,

    isFullWidthPage,
    isHomePage,
    isDocsPage,

    mainContentStyle,

    toggleLayout: layoutStore.toggleLayoutMode,
    toggleSidebar: layoutStore.toggleSidebar,
  }
}

export function useParticleBackground() {
  function getParticleStyle(_index: number) {
    const left = Math.random() * 100
    const animationDelay = Math.random() * 10
    const animationDuration = 8 + Math.random() * 12
    const size = 1 + Math.random() * 3

    return {
      left: `${left}%`,
      animationDelay: `${animationDelay}s`,
      animationDuration: `${animationDuration}s`,
      width: `${size}px`,
      height: `${size}px`,
    }
  }

  return {
    getParticleStyle,
  }
}
