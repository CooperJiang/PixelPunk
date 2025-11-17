/**
 * 多语言支持入口 - 动态加载优化版
 * 参考 themeLoader 的设计，按需加载语言包
 */
import type { ZhCNLocale } from './zh-CN'
import type { JaJPLocale } from './ja-JP'

export const DEFAULT_LOCALE = 'zh-CN'
export const DEFAULT_THEME = 'normal'

export type SupportedLocale = 'zh-CN' | 'en-US' | 'ja-JP'
export type LocaleData = ZhCNLocale | JaJPLocale

// 语言包缓存
const localeCache = new Map<SupportedLocale, LocaleData>()
// 正在加载的语言包Promise，防止重复加载
const loadingLocales = new Map<SupportedLocale, Promise<LocaleData>>()

/**
 * 动态加载语言包
 */
async function loadLocale(locale: SupportedLocale): Promise<LocaleData> {
  switch (locale) {
    case 'zh-CN':
      return (await import('./zh-CN')).zhCN
    case 'en-US':
      return (await import('./en-US')).enUS
    case 'ja-JP':
      return (await import('./ja-JP')).jaJP
    default:
      return (await import('./zh-CN')).zhCN
  }
}

/**
 * 获取语言包（带缓存）
 */
export async function getLocale(locale: SupportedLocale): Promise<LocaleData> {
  // 如果已缓存，直接返回
  if (localeCache.has(locale)) {
    return localeCache.get(locale)!
  }

  // 如果正在加载，等待加载完成
  if (loadingLocales.has(locale)) {
    return loadingLocales.get(locale)!
  }

  // 开始加载
  const loadPromise = loadLocale(locale).then((data) => {
    localeCache.set(locale, data)
    loadingLocales.delete(locale)
    return data
  }).catch((error) => {
    loadingLocales.delete(locale)
    throw error
  })

  loadingLocales.set(locale, loadPromise)
  return loadPromise
}

/**
 * 获取支持的语言列表
 */
export function getSupportedLocales(): SupportedLocale[] {
  return ['zh-CN', 'en-US', 'ja-JP']
}

/**
 * 检查语言是否支持
 */
export function isLocaleSupported(locale: string): locale is SupportedLocale {
  return ['zh-CN', 'en-US', 'ja-JP'].includes(locale)
}

/**
 * 预加载语言包（在空闲时加载，提升用户体验）
 */
export function preloadLocales(locales: SupportedLocale[]): void {
  if (typeof requestIdleCallback !== 'undefined') {
    requestIdleCallback(
      () => {
        locales.forEach((locale) => {
          getLocale(locale).catch(() => {})
        })
      },
      { timeout: 2000 }
    )
  } else {
    setTimeout(() => {
      locales.forEach((locale) => {
        getLocale(locale).catch(() => {})
      })
    }, 1000)
  }
}

/**
 * 检查语言包是否已加载
 */
export function isLocaleLoaded(locale: SupportedLocale): boolean {
  return localeCache.has(locale)
}

/**
 * 清除缓存（主要用于开发/测试）
 */
export function clearLocaleCache(): void {
  localeCache.clear()
  loadingLocales.clear()
}
