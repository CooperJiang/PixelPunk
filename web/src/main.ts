import { createApp, nextTick } from 'vue'
import App from './App.vue'
import router from './router'
import './styles/main.css'
/* 导入布局适配变量系统 */
import './styles/layout-vars.css'
/* 导入赛博朋克动画样式库 */
import './styles/design-system/animations.css'
/* 导入路由切换动画样式 */
import './styles/route-animations.css'
/* 导入vue-easy-lightbox样式 */
import 'vue-easy-lightbox/dist/external-css/vue-easy-lightbox.css'
/* 导入Markdown编辑器主题覆盖样式 */
import './styles/markdown-editor-theme.css'
import { createPinia } from 'pinia'
/* 导入所有指令 */
import directivesPlugin from './directives'
import CyberComponentsPlugin from './components'
/* 导入认证store */
import { useAuthStore } from './store/auth'
/* 导入设置store */
import { useSettingsStore } from './store/settings'
import { useLayoutStore } from './store/layout'
/* 导入上传存储管理器 */
import { UploadStorageManager } from './utils/storage/uploadStorage'
/* 导入消息系统 */
import { messageSystem } from './components/Message'
/* 导入文案主题系统 */
import { installTextTheme } from './plugins/textTheme'
/* 导入统一主题管理 */
import { useTheme } from './composables/useTheme'
/* 导入favicon管理器 */
import FaviconManager from './utils/favicon'

/* 异步初始化应用 */
async function initializeApp() {
  FaviconManager.setImmediate()

  const app = createApp(App)

  const pinia = createPinia()
  app.use(pinia)

  const settingsStore = useSettingsStore()
  const layoutStore = useLayoutStore()

  try {
    await settingsStore.loadGlobalSettings()
    await layoutStore.initializeLayout()
  } catch (error) {
    console.error('Settings initialization failed:', error)
  }

  app.use(router)

  app.use(CyberComponentsPlugin)

  app.use(directivesPlugin)

  app.use(installTextTheme)

  const theme = useTheme()
  const systemDefaultTheme = settingsStore.defaultTheme
  await theme.initialize(systemDefaultTheme)

  messageSystem.init()

  app.mount('#app')

  return app
}

initializeApp().catch((error) => {
  console.error('App initialization failed:', error)
})

nextTick(() => {
  const initializeAuth = () => {
    const authStore = useAuthStore()
    authStore.initAuth()
  }

  const stopCleanup = UploadStorageManager.startPeriodicCleanup()

  window.addEventListener('beforeunload', stopCleanup)

  if (document.readyState === 'complete' || document.readyState === 'interactive') {
    setTimeout(initializeAuth, 0)
  } else {
    document.addEventListener('DOMContentLoaded', initializeAuth)
  }

  window.addEventListener('load', () => {
    document.dispatchEvent(new Event('render-event'))
  })
})
