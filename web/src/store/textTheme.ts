import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { getLocale, DEFAULT_LOCALE, DEFAULT_THEME, isLocaleSupported, preloadLocales } from '@/locales'
import type { SupportedLocale, TextTheme, TextKey, LocaleData } from '@/locales/zh-CN'
import { detectBrowserLanguage } from '@/utils/language'

/**
 * 文案主题状态管理 - 支持动态加载语言包
 */
export const useTextThemeStore = defineStore('textTheme', () => {
  const currentLocale = ref<SupportedLocale>(DEFAULT_LOCALE)

  const currentTheme = ref<TextTheme>(DEFAULT_THEME)

  const isInitializing = ref(true)

  // 改为ref以支持异步加载
  const localeData = ref<LocaleData | null>(null)

  const themeTexts = computed(() => localeData.value?.themes[currentTheme.value] || {})

  const commonTexts = computed(() => localeData.value?.common || {})

  const isCyberTheme = computed(() => currentTheme.value === 'cyber')

  async function initialize() {
    try {
      const savedTheme = localStorage.getItem('text-theme')
      const savedLocale = localStorage.getItem('locale')

      if (savedTheme && (savedTheme === 'normal' || savedTheme === 'cyber')) {
        currentTheme.value = savedTheme as TextTheme
      }

      if (savedLocale && isLocaleSupported(savedLocale)) {
        currentLocale.value = savedLocale as SupportedLocale
      } else {
        try {
          const { getGlobalSettings } = await import('@/api/admin/settings')
          const response = await getGlobalSettings()

          if (response.code === 200 && response.data?.appearance) {
            const defaultLanguage = response.data.appearance.default_language
            const enableMultiLanguage = response.data.appearance.enable_multi_language

            if (enableMultiLanguage && defaultLanguage) {
              let targetLocale: SupportedLocale

              if (defaultLanguage === 'auto') {
                targetLocale = detectBrowserLanguage()
              } else if (isLocaleSupported(defaultLanguage as string)) {
                targetLocale = defaultLanguage as SupportedLocale
              } else {
                targetLocale = DEFAULT_LOCALE
              }

              currentLocale.value = targetLocale
              localStorage.setItem('locale', targetLocale)
            }
          }
        } catch (error) {
          // 静默处理错误，使用系统默认语言
        }
      }

      // 异步加载语言包
      localeData.value = await getLocale(currentLocale.value)

      // 预加载其他语言包（在空闲时）
      const otherLocales = (['zh-CN', 'en-US', 'ja-JP'] as SupportedLocale[]).filter(
        (locale) => locale !== currentLocale.value
      )
      preloadLocales(otherLocales)
    } catch (error) {
      // 静默处理错误，使用默认设置
      currentTheme.value = DEFAULT_THEME
      currentLocale.value = DEFAULT_LOCALE
      // 加载默认语言包
      try {
        localeData.value = await getLocale(DEFAULT_LOCALE)
      } catch (fallbackError) {
        // 最后的fallback，设置一个空对象避免报错
        localeData.value = { common: {}, themes: { normal: {}, cyber: {} } } as any
      }
    } finally {
      isInitializing.value = false
    }
  }

  function setTheme(theme: TextTheme) {
    currentTheme.value = theme
    localStorage.setItem('text-theme', theme)
  }

  async function setLocale(locale: SupportedLocale) {
    // 显示加载状态
    isInitializing.value = true
    try {
      currentLocale.value = locale
      localStorage.setItem('locale', locale)
      // 异步加载新的语言包
      localeData.value = await getLocale(locale)
    } catch (error) {
      // 加载失败时回退到默认语言
      currentLocale.value = DEFAULT_LOCALE
      localStorage.setItem('locale', DEFAULT_LOCALE)
      try {
        localeData.value = await getLocale(DEFAULT_LOCALE)
      } catch (fallbackError) {
        // 最后的fallback
        localeData.value = { common: {}, themes: { normal: {}, cyber: {} } } as any
      }
    } finally {
      isInitializing.value = false
    }
  }

  function toggleTheme() {
    const newTheme = currentTheme.value === 'normal' ? 'cyber' : 'normal'
    setTheme(newTheme)
  }

  function getThemeText(key: TextKey): string {
    return themeTexts.value[key] || key
  }

  function getCommonText(key: keyof typeof commonTexts.value): string {
    return commonTexts.value[key] || key
  }

  function getNestedValue(obj: any, path: string): string | undefined {
    return path.split('.').reduce((current, key) => {
      return current && typeof current === 'object' ? current[key] : undefined
    }, obj)
  }

  function getText(key: string): string {
    if (!localeData.value) {
      return key
    }

    if (key.includes('.')) {
      const themeValue = getNestedValue(themeTexts.value, key)
      if (themeValue) {
        return themeValue
      }

      const commonValue = getNestedValue(commonTexts.value, key)
      if (commonValue) {
        return commonValue
      }
    } else {
      if (key in themeTexts.value) {
        return themeTexts.value[key as TextKey]
      }

      if (key in commonTexts.value) {
        return commonTexts.value[key as keyof typeof commonTexts.value]
      }
    }

    return key
  }

  async function reset() {
    setTheme(DEFAULT_THEME)
    await setLocale(DEFAULT_LOCALE)
  }

  return {
    currentLocale,
    currentTheme,
    isInitializing,
    localeData,
    themeTexts,
    commonTexts,
    isCyberTheme,

    initialize,
    setTheme,
    setLocale,
    toggleTheme,
    getThemeText,
    getCommonText,
    getText,
    reset,
  }
})

export async function initializeTextTheme() {
  const store = useTextThemeStore()
  await store.initialize()
}
