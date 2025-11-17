import type { SupportedLocale } from '@/locales'

/**
 * 检测浏览器语言并映射到支持的语言
 */
export function detectBrowserLanguage(): SupportedLocale {
  const browserLang = navigator.language || (navigator as any).userLanguage

  const langMap: Record<string, SupportedLocale> = {
    'zh': 'zh-CN',
    'zh-CN': 'zh-CN',
    'zh-Hans': 'zh-CN',
    'zh-SG': 'zh-CN',
    'zh-TW': 'zh-CN',
    'zh-HK': 'zh-CN',
    'zh-Hant': 'zh-CN',
    'en': 'en-US',
    'en-US': 'en-US',
    'en-GB': 'en-US',
    'en-AU': 'en-US',
    'en-CA': 'en-US',
    'en-NZ': 'en-US',
    'ja': 'ja-JP',
    'ja-JP': 'ja-JP',
  }

  if (langMap[browserLang]) {
    return langMap[browserLang]
  }

  const langPrefix = browserLang.split('-')[0]
  if (langMap[langPrefix]) {
    return langMap[langPrefix]
  }

  return 'zh-CN'
}
