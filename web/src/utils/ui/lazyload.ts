/**
 * 文件懒加载工具
 * 提供高性能的文件延迟加载功能

 */
interface LazyLoadOptions {
  rootMargin?: string
  threshold?: number | number[]
  placeholder?: string
  errorImage?: string
  fadeIn?: boolean
  retryTimes?: number
  retryDelay?: number
}

interface LazyImageElement extends HTMLImageElement {
  dataset: {
    src?: string
    srcset?: string
    loaded?: string
    retryCount?: string
  }
}

class LazyImageLoader {
  private observer: IntersectionObserver | null = null
  private options: Required<LazyLoadOptions>
  private loadingQueue = new Set<LazyImageElement>()
  private retryQueue = new Map<LazyImageElement, number>()

  constructor(options: LazyLoadOptions = {}) {
    this.options = {
      rootMargin: '50px',
      threshold: 0.1,
      placeholder:
        'data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iMTAwIiBoZWlnaHQ9IjEwMCIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj48cmVjdCB3aWR0aD0iMTAwIiBoZWlnaHQ9IjEwMCIgZmlsbD0iI2Y1ZjVmNSIvPjx0ZXh0IHg9IjUwIiB5PSI1MCIgZm9udC1mYW1pbHk9IkFyaWFsIiBmb250LXNpemU9IjE0IiBmaWxsPSIjOTk5IiB0ZXh0LWFuY2hvcj0ibWlkZGxlIiBkeT0iMC4zZW0iPkxvYWRpbmcuLi48L3RleHQ+PC9zdmc+',
      errorImage:
        'data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iMTAwIiBoZWlnaHQ9IjEwMCIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj48cmVjdCB3aWR0aD0iMTAwIiBoZWlnaHQ9IjEwMCIgZmlsbD0iI2ZmZWJlZSIvPjx0ZXh0IHg9IjUwIiB5PSI1MCIgZm9udC1mYW1pbHk9IkFyaWFsIiBmb250LXNpemU9IjEyIiBmaWxsPSIjZGM0NDQ0IiB0ZXh0LWFuY2hvcj0ibWlkZGxlIiBkeT0iMC4zZW0iPkVycm9yPC90ZXh0Pjwvc3ZnPg==',
      fadeIn: true,
      retryTimes: 3,
      retryDelay: 1000,
      ...options,
    }

    this.initObserver()
  }

  private initObserver() {
    if (!window.IntersectionObserver) {
      this.fallbackLoad()
      return
    }

    this.observer = new IntersectionObserver(
      (entries) => {
        entries.forEach((entry) => {
          if (entry.isIntersecting) {
            const img = entry.target as LazyImageElement
            this.loadImage(img)
            this.observer?.unobserve(img)
          }
        })
      },
      {
        rootMargin: this.options.rootMargin,
        threshold: this.options.threshold,
      }
    )
  }

  private async loadImage(img: LazyImageElement) {
    if (img.dataset.loaded === 'true') {
      return
    }

    const originalSrc = img.dataset.src
    const originalSrcset = img.dataset.srcset

    if (!originalSrc) {
      return
    }

    this.loadingQueue.add(img)

    if (!img.src || img.src === this.options.errorImage) {
      img.src = this.options.placeholder
    }

    try {
      await this.preloadImage(originalSrc)

      img.src = originalSrc
      if (originalSrcset) {
        img.srcset = originalSrcset
      }

      if (this.options.fadeIn) {
        img.style.opacity = '0'
        img.style.transition = 'opacity 0.3s ease-in-out'

        setTimeout(() => {
          img.style.opacity = '1'
        }, 10)
      }

      img.dataset.loaded = 'true'
      this.loadingQueue.delete(img)
      this.retryQueue.delete(img)

      img.dispatchEvent(new CustomEvent('lazyloaded', { detail: { src: originalSrc } }))
    } catch {
      this.handleLoadError(img, originalSrc)
    }
  }

  private preloadImage(src: string): Promise<void> {
    return new Promise((resolve, reject) => {
      const tempImg = new Image()

      tempImg.onload = () => resolve()
      tempImg.onerror = () => reject(new Error(`Failed to load image: ${src}`))

      tempImg.src = src
    })
  }

  private handleLoadError(img: LazyImageElement, originalSrc: string) {
    const retryCount = this.retryQueue.get(img) || 0

    if (retryCount < this.options.retryTimes) {
      this.retryQueue.set(img, retryCount + 1)

      setTimeout(
        () => {
          this.loadImage(img)
        },
        this.options.retryDelay * (retryCount + 1)
      )
    } else {
      img.src = this.options.errorImage
      img.dataset.loaded = 'error'
      this.loadingQueue.delete(img)
      this.retryQueue.delete(img)

      img.dispatchEvent(new CustomEvent('lazyerror', { detail: { src: originalSrc } }))
    }
  }

  private fallbackLoad() {
    document.addEventListener('DOMContentLoaded', () => {
      const images = document.querySelectorAll('img[data-src]') as NodeListOf<LazyImageElement>
      images.forEach((img) => this.loadImage(img))
    })
  }

  observe(img: LazyImageElement) {
    if (!this.observer) {
      this.loadImage(img)
      return
    }

    this.observer.observe(img)
  }

  unobserve(img: LazyImageElement) {
    if (this.observer) {
      this.observer.unobserve(img)
    }
    this.loadingQueue.delete(img)
    this.retryQueue.delete(img)
  }

  loadAll() {
    this.loadingQueue.forEach((img) => this.loadImage(img))
  }

  getStats() {
    return {
      loading: this.loadingQueue.size,
      retrying: this.retryQueue.size,
    }
  }

  destroy() {
    if (this.observer) {
      this.observer.disconnect()
      this.observer = null
    }
    this.loadingQueue.clear()
    this.retryQueue.clear()
  }
}

let defaultLoader: LazyImageLoader | null = null

export function createLazyLoader(options?: LazyLoadOptions) {
  return new LazyImageLoader(options)
}

export function getDefaultLazyLoader(options?: LazyLoadOptions) {
  if (!defaultLoader) {
    defaultLoader = new LazyImageLoader(options)
  }
  return defaultLoader
}

export const lazyloadDirective = {
  mounted(el: LazyImageElement, binding: any) {
    const loader = getDefaultLazyLoader(binding.value)

    if (binding.arg) {
      el.dataset.src = binding.arg
    }

    loader.observe(el)
  },

  unmounted(el: LazyImageElement) {
    const loader = getDefaultLazyLoader()
    loader.unobserve(el)
  },
}

export function useLazyload(options?: LazyLoadOptions) {
  const loader = ref<LazyImageLoader | null>(null)

  onMounted(() => {
    loader.value = createLazyLoader(options)
  })

  onUnmounted(() => {
    loader.value?.destroy()
  })

  const observe = (img: LazyImageElement) => {
    loader.value?.observe(img)
  }

  const unobserve = (img: LazyImageElement) => {
    loader.value?.unobserve(img)
  }

  return {
    observe,
    unobserve,
    loader: readonly(loader),
  }
}
