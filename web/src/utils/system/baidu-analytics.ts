/**
 * 初始化百度统计
 * @param siteId 百度统计站点ID
 */
export const initBaiduAnalytics = (siteId: string): void => {
  try {
    if (typeof window === 'undefined') {
      return
    }

    if (window._hmt && document.querySelector(`script[src*="hm.baidu.com/hm.js?${siteId}"]`)) {
      return
    }

    window._hmt = window._hmt || []
    const hm = document.createElement('script')
    hm.src = `https://hm.baidu.com/hm.js?${siteId}`

    hm.onerror = () => {
      console.error('[Baidu Analytics] Failed to load script')
    }

    const s = document.getElementsByTagName('script')[0]
    s.parentNode?.insertBefore(hm, s)
  } catch (error) {
    console.error('[Baidu Analytics] Failed to initialize:', error)
  }
}

/**
 * 检查百度统计是否已初始化
 */
export const checkBaiduAnalytics = (): { initialized: boolean; siteId?: string; scriptLoaded: boolean } => {
  const initialized = !!(window._hmt && window._hmt.length >= 0)
  const script = document.querySelector('script[src*="hm.baidu.com/hm.js"]') as HTMLScriptElement
  const scriptLoaded = !!script
  const siteId = script?.src.match(/hm\.js\?(.+)$/)?.[1]

  return {
    initialized,
    siteId,
    scriptLoaded
  }
}

export class BaiduAnalytics {
  static trackPageview(pageUrl: string): void {
    if (window._hmt) {
      window._hmt.push(['_trackPageview', pageUrl])
    }
  }

  static trackEvent(category: string, action: string, label?: string, value?: number): void {
    if (window._hmt) {
      window._hmt.push(['_trackEvent', category, action, label || '', value || ''])
    }
  }

  static setCustomVar(name: string, value: string): void {
    if (window._hmt) {
      window._hmt.push(['_setCustomVar', name, value])
    }
  }
}

declare global {
  interface Window {
    _hmt: any[]
  }
}
