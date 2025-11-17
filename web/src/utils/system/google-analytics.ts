/**
 * 初始化 Google Analytics 4
 * @param measurementId Google Analytics 测量ID (格式: G-XXXXXXXXXX)
 */
export const initGoogleAnalytics = (measurementId: string): void => {
  try {
    if (typeof window === 'undefined') {
      return
    }

    if (window.gtag && document.querySelector(`script[src*="googletagmanager.com/gtag/js?id=${measurementId}"]`)) {
      return
    }

    window.dataLayer = window.dataLayer || []
    window.gtag = function () {
      // eslint-disable-next-line prefer-rest-params
      window.dataLayer.push(arguments)
    }

    const script = document.createElement('script')
    script.async = true
    script.src = `https://www.googletagmanager.com/gtag/js?id=${measurementId}`

    script.onload = () => {
      window.gtag('js', new Date())
      window.gtag('config', measurementId)
    }

    script.onerror = () => {
      console.error('[Google Analytics] Failed to load script')
    }

    const firstScript = document.getElementsByTagName('script')[0]
    firstScript.parentNode?.insertBefore(script, firstScript)
  } catch (error) {
    console.error('[Google Analytics] Failed to initialize:', error)
  }
}

/**
 * 检查 Google Analytics 是否已初始化
 */
export const checkGoogleAnalytics = (): { initialized: boolean; measurementId?: string; scriptLoaded: boolean } => {
  const initialized = !!(window.gtag && window.dataLayer)
  const script = document.querySelector('script[src*="googletagmanager.com/gtag/js"]') as HTMLScriptElement
  const scriptLoaded = !!script
  const measurementId = script?.src.match(/id=(G-[A-Z0-9]+)/)?.[1]

  return {
    initialized,
    measurementId,
    scriptLoaded
  }
}

export class GoogleAnalytics {
  /**
   * 跟踪页面浏览
   */
  static trackPageview(pageUrl: string, pageTitle?: string): void {
    if (window.gtag) {
      window.gtag('event', 'page_view', {
        page_path: pageUrl,
        page_title: pageTitle || document.title
      })
    }
  }

  /**
   * 跟踪事件
   */
  static trackEvent(eventName: string, eventParams?: Record<string, any>): void {
    if (window.gtag) {
      window.gtag('event', eventName, eventParams)
    }
  }

  /**
   * 设置用户属性
   */
  static setUserProperties(properties: Record<string, any>): void {
    if (window.gtag) {
      window.gtag('set', 'user_properties', properties)
    }
  }

  /**
   * 设置用户ID
   */
  static setUserId(userId: string): void {
    if (window.gtag) {
      window.gtag('set', { user_id: userId })
    }
  }
}

declare global {
  interface Window {
    dataLayer: any[]
    gtag: (...args: any[]) => void
  }
}
