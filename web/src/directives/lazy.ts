import type { Directive, DirectiveBinding } from 'vue'

/* 存储所有需要懒加载的元素及其处理函数 */
const lazyElements = new Map<HTMLElement, () => void>()

/* 使用 IntersectionObserver 监听元素是否进入视口 */
let observer: IntersectionObserver | null = null

/**
 * 创建并初始化 IntersectionObserver
 */
const createObserver = () => {
  if ('IntersectionObserver' in window) {
    observer = new IntersectionObserver(
      (entries) => {
        entries.forEach((entry) => {
          if (entry.isIntersecting) {
            const element = entry.target as HTMLElement
            const loadHandler = lazyElements.get(element)

            if (loadHandler) {
              loadHandler()

              observer?.unobserve(element)
              lazyElements.delete(element)
            }
          }
        })
      },
      {
        rootMargin: '100px',
        threshold: 0.1,
      }
    )
  }
}

export const lazy: Directive = {
  mounted(el: HTMLElement, binding: DirectiveBinding) {
    if (!observer) {
      createObserver()
    }

    if (!observer) {
      if (el.tagName === 'IMG') {
        const img = el as HTMLImageElement
        const dataSrc = img.getAttribute('data-src')
        if (dataSrc) {
          img.src = dataSrc
          img.removeAttribute('data-src')
        }
      }
      el.style.opacity = '1'
      return
    }

    const shouldLazy = binding.value !== false

    if (!shouldLazy) {
      if (el.tagName === 'IMG') {
        const img = el as HTMLImageElement
        const dataSrc = img.getAttribute('data-src')
        if (dataSrc) {
          img.src = dataSrc
          img.removeAttribute('data-src')
        }
      }
      el.style.opacity = '1'
      return
    }

    if (el.tagName === 'IMG') {
      const img = el as HTMLImageElement
      if (img.src) {
        img.setAttribute('data-src', img.src)
        img.removeAttribute('src')
      }
    }

    el.style.opacity = '0'
    el.style.transition = 'opacity 0.3s ease-in-out'

    const loadHandler = () => {
      if (el.tagName === 'IMG') {
        const img = el as HTMLImageElement
        const dataSrc = img.getAttribute('data-src')
        if (dataSrc) {
          img.src = dataSrc
          img.removeAttribute('data-src')

          img.onload = () => {
            el.style.opacity = '1'
          }
          img.onerror = () => {
            el.style.opacity = '1' // 即使加载失败也显示
          }
        } else {
          el.style.opacity = '1'
        }
      } else {
        el.style.opacity = '1'
      }
    }

    lazyElements.set(el, loadHandler)

    observer.observe(el)
  },

  updated(el: HTMLElement, binding: DirectiveBinding) {
    const shouldLazy = binding.value !== false

    if (!shouldLazy && el.style.opacity === '0') {
      if (el.tagName === 'IMG') {
        const img = el as HTMLImageElement
        const dataSrc = img.getAttribute('data-src')
        if (dataSrc) {
          img.src = dataSrc
          img.removeAttribute('data-src')
        }
      }
      el.style.opacity = '1'

      observer?.unobserve(el)
      lazyElements.delete(el)
    }
  },

  beforeUnmount(el: HTMLElement) {
    observer?.unobserve(el)
    lazyElements.delete(el)
  },
}

export default {
  install(app: any) {
    app.directive('lazy', lazy)
  },
}
