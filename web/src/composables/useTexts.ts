import { computed } from 'vue'
import { useTextThemeStore } from '@/store/textTheme'
import type { TextKey, TextTheme, SupportedLocale } from '@/locales/zh-CN'

/**
 * Translation function type
 */
export type TranslationFunction = (key: string, params?: Record<string, string | number>) => string

/**
 * 文案获取组合函数
 * 提供简单易用的API来获取主题化文案
 */
export function useTexts() {
  const store = useTextThemeStore()

  const currentTheme = computed(() => store.currentTheme)

  const currentLocale = computed(() => store.currentLocale)

  const isCyberTheme = computed(() => store.isCyberTheme)

  const texts = computed(() => store.themeTexts)

  const common = computed(() => store.commonTexts)

  const $t = (key: string, params?: Record<string, string | number>): string => {
    let text = store.getText(key)

    if (params) {
      Object.keys(params).forEach((paramKey) => {
        const placeholder = `{${paramKey}}`
        const value = String(params[paramKey])
        // 使用 replaceAll 避免正则表达式转义问题
        text = text.replaceAll(placeholder, value)
      })
    }

    return text
  }

  const $tr = (key: string) => {
    return computed(() => store.getText(key))
  }

  const getThemeText = (key: TextKey): string => {
    return store.getThemeText(key)
  }

  const getCommonText = (key: keyof typeof common.value): string => {
    return store.getCommonText(key)
  }

  const setTheme = (theme: TextTheme): void => {
    store.setTheme(theme)
  }

  const toggleTheme = (): void => {
    store.toggleTheme()
  }

  const setLocale = (locale: SupportedLocale): void => {
    store.setLocale(locale)
  }

  return {
    currentTheme,
    currentLocale,
    isCyberTheme,
    texts,
    common,

    $t, // 主要的文案获取方法
    $tr, // 响应式文案获取方法
    getThemeText, // 专门获取主题文案
    getCommonText, // 专门获取通用文案

    setTheme,
    toggleTheme,

    setLocale,
  }
}

export function useQuickTexts() {
  const { $t } = useTexts()

  return {
    upload: computed(() => $t('upload')),
    download: computed(() => $t('download')),
    delete: computed(() => $t('delete')),
    save: computed(() => $t('save')),
    cancel: computed(() => $t('cancel')),
    confirm: computed(() => $t('confirm')),

    loading: computed(() => $t('loading')),
    success: computed(() => $t('success')),
    error: computed(() => $t('error')),

    zoomIn: computed(() => $t('zoomIn')),
    zoomOut: computed(() => $t('zoomOut')),
    previous: computed(() => $t('previous')),
    next: computed(() => $t('next')),
    fullscreen: computed(() => $t('fullscreen')),

    text: $t,
  }
}
