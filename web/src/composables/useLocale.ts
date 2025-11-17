import { computed } from 'vue'
import { useTextThemeStore } from '@/store/textTheme'
import { useTexts } from '@/composables/useTexts'
import type { SupportedLocale } from '@/locales'

/**
 * 🌍 语言管理 Composable
 * 提供语言切换功能
 */
export function useLocale() {
  const textThemeStore = useTextThemeStore()
  const { $t } = useTexts()

  const currentLocale = computed(() => textThemeStore.currentLocale)

  const isInitializing = computed(() => textThemeStore.isInitializing)

  function setLocale(locale: SupportedLocale) {
    textThemeStore.setLocale(locale)
  }

  interface LocaleConfig {
    label: string
    value: SupportedLocale
    description: string
    emoji: string
    nativeName: string
  }

  const LOCALE_CONFIGS = computed<LocaleConfig[]>(() => [
    {
      label: $t('locale.zhCN.label'),
      value: 'zh-CN',
      description: $t('locale.zhCN.description'),
      emoji: '🇨🇳',
      nativeName: '简体中文',
    },
    {
      label: $t('locale.enUS.label'),
      value: 'en-US',
      description: $t('locale.enUS.description'),
      emoji: '🇺🇸',
      nativeName: 'English',
    },
    {
      label: $t('locale.jaJP.label'),
      value: 'ja-JP',
      description: $t('locale.jaJP.description'),
      emoji: '🇯🇵',
      nativeName: '日本語',
    },
  ])

  function getLocaleLabel(locale: SupportedLocale): string {
    return LOCALE_CONFIGS.value.find((l) => l.value === locale)?.label || locale
  }

  function getLocaleDescription(locale: SupportedLocale): string {
    return LOCALE_CONFIGS.value.find((l) => l.value === locale)?.description || ''
  }

  function getLocaleIcon(locale: SupportedLocale): string {
    return LOCALE_CONFIGS.value.find((l) => l.value === locale)?.emoji || '🌐'
  }

  function getLocaleNativeName(locale: SupportedLocale): string {
    return LOCALE_CONFIGS.value.find((l) => l.value === locale)?.nativeName || locale
  }

  function getLocaleEmoji(locale: SupportedLocale): string {
    return LOCALE_CONFIGS.value.find((l) => l.value === locale)?.emoji || '🌐'
  }

  const allLocales = computed(() => LOCALE_CONFIGS.value.map((l) => l.value))

  const localeOptions = LOCALE_CONFIGS

  const currentLocaleInfo = computed(() => {
    const locale = currentLocale.value
    return LOCALE_CONFIGS.value.find((l) => l.value === locale) || LOCALE_CONFIGS.value[0]
  })

  return {
    currentLocale,
    isInitializing,
    setLocale,
    getLocaleLabel,
    getLocaleDescription,
    getLocaleIcon,
    getLocaleNativeName,
    getLocaleEmoji,
    localeOptions,
    currentLocaleInfo,
    allLocales,
  }
}

export type { SupportedLocale }
