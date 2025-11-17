import { computed } from 'vue'
import { useVisualThemeStore, type VisualTheme, type AppliedTheme } from '@/store/visualTheme'
import { useTextThemeStore } from '@/store/textTheme'
import { useTexts } from '@/composables/useTexts'
import type { TextTheme } from '@/locales/zh-CN'

/**
 * 🎨 统一主题管理 Composable
 * 整合 visualTheme + textTheme，提供一个入口
 * 简化主题系统的使用，避免导入多个 store
 */
export function useTheme() {
  const visualThemeStore = useVisualThemeStore()
  const textThemeStore = useTextThemeStore()

  const selectedTheme = computed(() => visualThemeStore.selectedTheme)

  const appliedTheme = computed(() => visualThemeStore.appliedTheme)

  const isInitializing = computed(() => visualThemeStore.isInitializing)

  const isLight = computed(() => visualThemeStore.isLight)

  const isDark = computed(() => visualThemeStore.isDark)

  const isCyberpunk = computed(() => visualThemeStore.isCyberpunk)

  function setVisualTheme(theme: VisualTheme) {
    visualThemeStore.setTheme(theme)
  }

  function toggleVisualTheme() {
    visualThemeStore.toggleTheme()
  }

  function resetVisualTheme() {
    visualThemeStore.reset()
  }

  const currentTextTheme = computed(() => textThemeStore.currentTheme)

  const isCyberText = computed(() => textThemeStore.isCyberTheme)

  function setTextTheme(theme: TextTheme) {
    textThemeStore.setTheme(theme)
  }

  function toggleTextTheme() {
    textThemeStore.toggleTheme()
  }

  function getText(key: string): string {
    return textThemeStore.getText(key)
  }

  async function initialize(systemDefaultTheme?: string) {
    await Promise.all([visualThemeStore.initialize(systemDefaultTheme), textThemeStore.initialize()])
  }

  function reset() {
    visualThemeStore.reset()
    textThemeStore.reset()
  }

  function setTheme(theme: VisualTheme) {
    setVisualTheme(theme)
  }

  function toggleTheme() {
    toggleVisualTheme()
  }

  function resetTheme() {
    resetVisualTheme()
  }

  interface ThemeConfig {
    label: string // 主题显示名称
    value: VisualTheme // 主题值（必须与 CSS 文件名一致）
    description: string // 主题描述
    icon: string // FontAwesome 图标名（不含 fa- 前缀）
    isDark: boolean // 是否为暗色主题
  }

  const { $t } = useTexts()

  const THEME_CONFIGS = computed<ThemeConfig[]>(() => [
    { label: $t('themes.light.label'), value: 'light', description: $t('themes.light.description'), icon: 'sun', isDark: false },
    { label: $t('themes.dark.label'), value: 'dark', description: $t('themes.dark.description'), icon: 'moon', isDark: true },

    {
      label: $t('themes.cyberpunk-cyan.label'),
      value: 'cyberpunk-cyan',
      description: $t('themes.cyberpunk-cyan.description'),
      icon: 'bolt',
      isDark: true,
    },
    {
      label: $t('themes.cyberpunk-neon.label'),
      value: 'cyberpunk-neon',
      description: $t('themes.cyberpunk-neon.description'),
      icon: 'fire',
      isDark: true,
    },

    {
      label: $t('themes.midnight-ocean.label'),
      value: 'midnight-ocean',
      description: $t('themes.midnight-ocean.description'),
      icon: 'water',
      isDark: true,
    },
    {
      label: $t('themes.forest-night.label'),
      value: 'forest-night',
      description: $t('themes.forest-night.description'),
      icon: 'tree',
      isDark: true,
    },
    {
      label: $t('themes.aurora-night.label'),
      value: 'aurora-night',
      description: $t('themes.aurora-night.description'),
      icon: 'star',
      isDark: true,
    },

    {
      label: $t('themes.sky-azure.label'),
      value: 'sky-azure',
      description: $t('themes.sky-azure.description'),
      icon: 'cloud-sun',
      isDark: false,
    },
    {
      label: $t('themes.lavender-dream.label'),
      value: 'lavender-dream',
      description: $t('themes.lavender-dream.description'),
      icon: 'wand-magic-sparkles',
      isDark: false,
    },
    {
      label: $t('themes.lemon-yellow.label'),
      value: 'lemon-yellow',
      description: $t('themes.lemon-yellow.description'),
      icon: 'lemon',
      isDark: false,
    },
    {
      label: $t('themes.cotton-candy.label'),
      value: 'cotton-candy',
      description: $t('themes.cotton-candy.description'),
      icon: 'candy-cane',
      isDark: false,
    },
    {
      label: $t('themes.coral-sunset.label'),
      value: 'coral-sunset',
      description: $t('themes.coral-sunset.description'),
      icon: 'palette',
      isDark: false,
    },
  ])

  function getThemeLabel(theme: VisualTheme): string {
    return THEME_CONFIGS.value.find((t) => t.value === theme)?.label || theme
  }

  function getThemeDescription(theme: VisualTheme): string {
    return THEME_CONFIGS.value.find((t) => t.value === theme)?.description || ''
  }

  function getThemeIcon(theme: VisualTheme): string {
    return THEME_CONFIGS.value.find((t) => t.value === theme)?.icon || 'palette'
  }

  const allThemes = computed(() => THEME_CONFIGS.value.map((t) => t.value))

  const themeOptions = THEME_CONFIGS

  const currentThemeInfo = computed(() => {
    const theme = selectedTheme.value
    return THEME_CONFIGS.value.find((t) => t.value === theme) || THEME_CONFIGS.value[1]
  })

  return {
    selectedTheme,
    appliedTheme,
    isInitializing,

    isLight,
    isDark,
    isCyberpunk,

    setVisualTheme,
    toggleVisualTheme,
    resetVisualTheme,

    currentTextTheme,
    isCyberText,

    setTextTheme,
    toggleTextTheme,
    getText,

    initialize,
    reset,

    getThemeLabel,
    getThemeDescription,
    getThemeIcon,
    themeOptions,
    currentThemeInfo,
    allThemes,

    setTheme,
    toggleTheme,
    resetTheme,
  }
}

export type { VisualTheme, AppliedTheme, TextTheme }
