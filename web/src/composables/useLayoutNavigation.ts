import { useRouter, useRoute } from 'vue-router'
import { useLayoutStore } from '@/store/layout'

/**
 * 布局导航组合式函数
 * 处理用户端布局切换逻辑（Admin后台不参与布局切换）
 */
export function useLayoutNavigation() {
  const router = useRouter()
  const route = useRoute()
  const layoutStore = useLayoutStore()

  const toggleLayoutWithNavigation = async () => {
    if (!layoutStore.canToggleLayout) {
      return {
        newMode: layoutStore.mode,
        targetPath: route.path,
        switched: false,
        reason: 'Multi-layout mode is disabled',
      }
    }

    const currentMode = layoutStore.mode
    const newMode = currentMode === 'left' ? 'top' : 'left'

    layoutStore.setLayoutMode(newMode)

    const targetPath = newMode === 'left' ? '/dashboard' : '/'

    try {
      await router.push(targetPath)
    } catch (error) {
      console.warn('Failed to navigate:', error)
    }

    return {
      newMode,
      targetPath,
      switched: true,
    }
  }

  return {
    currentLayout: layoutStore.mode,
    canToggleLayout: layoutStore.canToggleLayout,

    toggleLayoutWithNavigation,
  }
}
