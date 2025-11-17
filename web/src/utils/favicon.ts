/**
 * Favicon缓存管理器
 * 解决单页面应用中动态设置favicon的问题
 */
class FaviconManager {
  private static STORAGE_KEY = 'site_favicon_url'
  private static defaultFavicon = '/favicon.ico'

  static setImmediate(): boolean {
    const cachedUrl = localStorage.getItem(this.STORAGE_KEY)
    if (cachedUrl && cachedUrl !== 'null' && cachedUrl.trim()) {
      this.updateFavicon(cachedUrl)
      return true
    }
    return false
  }

  static updateAndCache(url: string): void {
    if (url && url.trim()) {
      this.updateFavicon(url)
      localStorage.setItem(this.STORAGE_KEY, url)
    } else {
      this.updateFavicon(this.defaultFavicon)
      localStorage.removeItem(this.STORAGE_KEY)
    }
  }

  private static updateFavicon(url: string): void {
    let link = document.querySelector("link[rel*='icon']") as HTMLLinkElement

    if (!link) {
      link = document.createElement('link')
      link.rel = 'icon'
      link.type = 'image/x-icon'
      document.head.appendChild(link)
    }

    link.href = url
  }

  static clearCache(): void {
    localStorage.removeItem(this.STORAGE_KEY)
    this.updateFavicon(this.defaultFavicon)
  }

  static getCachedUrl(): string | null {
    return localStorage.getItem(this.STORAGE_KEY)
  }

  static hasCache(): boolean {
    const cached = this.getCachedUrl()
    return cached !== null && cached !== 'null' && cached.trim() !== ''
  }
}

export default FaviconManager
