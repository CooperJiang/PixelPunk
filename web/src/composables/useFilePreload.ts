import { ref } from 'vue'

/**
 * 文件预加载Composable
 * 提供文件预加载功能，优化文件加载体验
 */
export function useImagePreload() {
  const preloadedUrls = ref(new Set<string>())

  const preloadImage = (url: string): Promise<void> => {
    if (preloadedUrls.value.has(url)) {
      return Promise.resolve()
    }

    return new Promise((resolve, reject) => {
      const img = new Image()

      img.onload = () => {
        preloadedUrls.value.add(url)
        resolve()
      }

      img.onerror = () => {
        reject(new Error(`Failed to preload image: ${url}`))
      }

      img.src = url
    })
  }

  const preloadImages = (urls: string[]): Promise<void[]> => {
    const urlsToLoad = urls.filter((url) => !preloadedUrls.value.has(url))

    if (urlsToLoad.length === 0) {
      return Promise.resolve([])
    }

    return Promise.all(urlsToLoad.map((url) => preloadImage(url)))
  }

  const isPreloaded = (url: string): boolean => preloadedUrls.value.has(url)

  const clearPreloadCache = () => {
    preloadedUrls.value.clear()
  }

  return {
    preloadImage,
    preloadImages,
    isPreloaded,
    clearPreloadCache,
    preloadedUrls,
  }
}
