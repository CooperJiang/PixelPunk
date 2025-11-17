import { createPinia } from 'pinia'
import type { App } from 'vue'

const pinia = createPinia()

/* 导出Pinia插件安装函数 */
export function setupStore(app: App) {
  app.use(pinia)
}

export { useAuthStore } from './auth'
export { useLayoutStore } from './layout'
export { useSettingsStore } from './settings'
export { useThemeStore } from './theme'
export { useWebSocketStore } from './websocket'
export { useVisualThemeStore } from './visualTheme'
export { useTextThemeStore } from './textTheme'
export { useBreadcrumbStore } from './breadcrumb'

export default pinia
