import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { loadTheme, preloadThemes } from '@/utils/themeLoader'

/**
 * 视觉主题类型
 * - light: 亮色主题
 * - dark: 暗色主题 (默认 - Midnight Black 午夜黑色)
 * - cyberpunk-cyan: 赛博朋克 - 霓虹青主题
 * - cyberpunk-neon: 赛博朋克 - 霓虹都市主题
 * - midnight-ocean: Midnight Ocean 午夜海洋主题
 * - forest-night: Forest Night 森林夜色主题
 * 🌸 可爱系列亮色主题
 * - sakura-pink: Sakura Pink 樱花粉主题
 * - lavender-dream: Lavender Dream 薰衣草梦境主题
 * - lemon-yellow: Lemon Yellow 柠檬黄主题
 * - cotton-candy: Cotton Candy 棉花糖主题
 * - coral-sunset: Coral Sunset 珊瑚日落主题
 * 🌟 新增主题
 * - aurora-night: Aurora Night 极光之夜主题
 * - sky-azure: Sky Azure 天空蔚蓝主题
 */
export type VisualTheme =
  | 'light'
  | 'dark'
  | 'cyberpunk-cyan'
  | 'cyberpunk-neon'
  | 'midnight-ocean'
  | 'forest-night'
  | 'lavender-dream'
  | 'lemon-yellow'
  | 'cotton-candy'
  | 'coral-sunset'
  | 'aurora-night'
  | 'sky-azure'

/**
 * 实际应用的主题
 */
export type AppliedTheme = VisualTheme

/**
 * 视觉主题状态管理
 * 负责控制整个应用的颜色主题 (亮色/暗色/赛博朋克)
 * 与 textTheme (文案风格) 完全独立
 */
export const useVisualThemeStore = defineStore('visualTheme', () => {
  const selectedTheme = ref<VisualTheme>('dark')
  const appliedTheme = ref<AppliedTheme>('dark')
  const isInitializing = ref(true)

  const isLight = computed(() => {
    const lightThemes: VisualTheme[] = ['light', 'lavender-dream', 'lemon-yellow', 'cotton-candy', 'coral-sunset', 'sky-azure']
    return lightThemes.includes(appliedTheme.value)
  })

  const isDark = computed(() => !isLight.value)
  const isCyberpunk = computed(() => appliedTheme.value === 'cyberpunk-cyan' || appliedTheme.value === 'cyberpunk-neon')
  const isCyberpunkCyan = computed(() => appliedTheme.value === 'cyberpunk-cyan')
  const isCyberpunkNeon = computed(() => appliedTheme.value === 'cyberpunk-neon')
  const isSynthwave = computed(() => appliedTheme.value === 'synthwave')

  function applyThemeToDOM(theme: AppliedTheme) {
    if (typeof document === 'undefined') return

    const html = document.documentElement
    html.removeAttribute('data-theme')
    html.classList.remove('theme-light', 'theme-dark', 'theme-cyberpunk')
    html.setAttribute('data-theme', theme)
    html.classList.add(`theme-${theme}`)
    updateMetaThemeColor(theme)
  }

  function updateMetaThemeColor(_theme: AppliedTheme) {
    if (typeof document === 'undefined') return

    const html = document.documentElement
    const computedStyle = getComputedStyle(html)
    const bgColor = computedStyle.getPropertyValue('--color-background-900').trim()
    const themeColor = bgColor || '#0a0e17'

    let metaThemeColor = document.querySelector('meta[name="theme-color"]')
    if (!metaThemeColor) {
      metaThemeColor = document.createElement('meta')
      metaThemeColor.setAttribute('name', 'theme-color')
      document.head.appendChild(metaThemeColor)
    }
    metaThemeColor.setAttribute('content', themeColor)
  }

  async function initialize(systemDefaultTheme?: string) {
    try {
      const savedTheme = localStorage.getItem('visual-theme')

      const validThemes = [
        'light',
        'dark',
        'cyberpunk-cyan',
        'cyberpunk-neon',
        'midnight-ocean',
        'forest-night',
        'lavender-dream',
        'lemon-yellow',
        'cotton-candy',
        'coral-sunset',
        'aurora-night',
        'sky-azure',
      ]

      if (savedTheme && validThemes.includes(savedTheme)) {
        selectedTheme.value = savedTheme as VisualTheme
      } else if (systemDefaultTheme && validThemes.includes(systemDefaultTheme)) {
        selectedTheme.value = systemDefaultTheme as VisualTheme
      } else {
        selectedTheme.value = 'dark'
      }

      await loadTheme(selectedTheme.value)
      appliedTheme.value = selectedTheme.value
      applyThemeToDOM(selectedTheme.value)
      preloadThemes(['cyberpunk-cyan', 'light'])
    } catch (error) {
      selectedTheme.value = 'dark'
      appliedTheme.value = 'dark'
      applyThemeToDOM('dark')
    } finally {
      isInitializing.value = false
    }
  }

  async function setTheme(theme: VisualTheme) {
    try {
      await loadTheme(theme)
      selectedTheme.value = theme
      appliedTheme.value = theme
      applyThemeToDOM(theme)

      try {
        localStorage.setItem('visual-theme', theme)
      } catch (error) {
        // 存储失败时静默处理
      }
    } catch (error) {
      throw error
    }
  }

  function toggleTheme() {
    const newTheme = appliedTheme.value === 'dark' ? 'light' : 'dark'
    setTheme(newTheme)
  }

  function reset() {
    setTheme('dark')
  }

  function cleanup() {}

  return {
    selectedTheme,
    appliedTheme,
    isInitializing,
    isLight,
    isDark,
    isCyberpunk,
    isCyberpunkCyan,
    isCyberpunkNeon,
    isSynthwave,
    initialize,
    setTheme,
    toggleTheme,
    reset,
    cleanup,
  }
})
