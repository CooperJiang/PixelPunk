import { computed, getCurrentInstance, type App } from 'vue'
import { useTextThemeStore } from '@/store/textTheme'

/**
 * 文案主题全局插件
 * 提供全局的 $t 方法
 */
export function installTextTheme(app: App) {
  const textCache = new Map<string, any>()

  app.config.globalProperties.$t = (key: string) => {
    const store = useTextThemeStore()

    const currentInstance = getCurrentInstance()
    if (currentInstance && currentInstance.isUnmounted === false) {
      if (!textCache.has(key)) {
        textCache.set(
          key,
          computed(() => store.getText(key))
        )
      }
      return textCache.get(key).value // 直接返回字符串值
    } else {
      return computed(() => store.getText(key))
    }
  }

  app.config.globalProperties.$toggleTheme = (): void => {
    const store = useTextThemeStore()
    store.toggleTheme()
  }

  app.config.globalProperties.$setTheme = (theme: 'normal' | 'cyber'): void => {
    const store = useTextThemeStore()
    store.setTheme(theme)
  }
}

declare module '@vue/runtime-core' {
  interface ComponentCustomProperties {
    $t: (key: string) => string | import('vue').ComputedRef<string>
    $toggleTheme: () => void
    $setTheme: (theme: 'normal' | 'cyber') => void
  }
}
