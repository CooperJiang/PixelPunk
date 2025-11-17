import type { Directive } from 'vue'
import { useLayoutStore } from '@/store/layout'

interface LayoutAdaptiveOptions {
  top?: Record<string, string>
  left?: Record<string, string>
  common?: Record<string, string>
  pageType?: string
}

/**
 * 布局自适应指令
 * 根据当前布局模式自动应用不同的样式
 */
export const vLayoutAdaptive: Directive<HTMLElement, LayoutAdaptiveOptions> = {
  mounted(el, binding) {
    updateElementStyles(el, binding.value)

    const store = useLayoutStore()
    const unwatch = store.$subscribe(() => {
      updateElementStyles(el, binding.value)
    })

    ;(el as any)._layoutAdaptiveUnwatch = unwatch
  },

  updated(el, binding) {
    updateElementStyles(el, binding.value)
  },

  unmounted(el) {
    const unwatch = (el as any)._layoutAdaptiveUnwatch
    if (unwatch) {
      unwatch()
      delete (el as any)._layoutAdaptiveUnwatch
    }
  },
}

function updateElementStyles(el: HTMLElement, options: LayoutAdaptiveOptions = {}) {
  const store = useLayoutStore()
  const { top = {}, left = {}, common = {} } = options

  Object.entries(common).forEach(([property, value]) => {
    el.style.setProperty(property, value)
  })

  const targetStyles = store.isTopLayout ? top : left
  Object.entries(targetStyles).forEach(([property, value]) => {
    el.style.setProperty(property, value)
  })

  el.classList.remove('layout-top', 'layout-left')
  el.classList.add(store.isTopLayout ? 'layout-top' : 'layout-left')

  if (options.pageType) {
    Array.from(el.classList).forEach((cls) => {
      if (cls.startsWith('page-type-')) {
        el.classList.remove(cls)
      }
    })
    el.classList.add(`page-type-${options.pageType}`)
  }
}

export const vResponsiveSpacing: Directive<HTMLElement, string> = {
  mounted(el, binding) {
    updateSpacing(el, binding.value)

    const store = useLayoutStore()
    const unwatch = store.$subscribe(() => {
      updateSpacing(el, binding.value)
    })

    ;(el as any)._spacingUnwatch = unwatch
  },

  updated(el, binding) {
    updateSpacing(el, binding.value)
  },

  unmounted(el) {
    const unwatch = (el as any)._spacingUnwatch
    if (unwatch) {
      unwatch()
      delete (el as any)._spacingUnwatch
    }
  },
}

function updateSpacing(el: HTMLElement, spacingType: string = 'element') {
  const store = useLayoutStore()

  const spacingMap = {
    element: store.isTopLayout ? '1rem' : '0.75rem',
    card: store.isTopLayout ? '1.5rem' : '1rem',
    section: store.isTopLayout ? '2rem' : '1.5rem',
  }

  const spacing = spacingMap[spacingType as keyof typeof spacingMap] || spacingMap.element
  el.style.setProperty('--responsive-spacing', spacing)
}

export const vResponsiveGrid: Directive<HTMLElement, { minWidth?: number; gap?: string }> = {
  mounted(el, binding) {
    updateGrid(el, binding.value)

    const store = useLayoutStore()
    const unwatch = store.$subscribe(() => {
      updateGrid(el, binding.value)
    })

    ;(el as any)._gridUnwatch = unwatch
  },

  updated(el, binding) {
    updateGrid(el, binding.value)
  },

  unmounted(el) {
    const unwatch = (el as any)._gridUnwatch
    if (unwatch) {
      unwatch()
      delete (el as any)._gridUnwatch
    }
  },
}

function updateGrid(el: HTMLElement, options: { minWidth?: number; gap?: string } = {}) {
  const store = useLayoutStore()
  const { minWidth = 300, gap } = options

  const adjustedMinWidth = store.isTopLayout ? minWidth : Math.max(minWidth - 50, 200)

  el.style.display = 'grid'
  el.style.gridTemplateColumns = `repeat(auto-fill, minmax(${adjustedMinWidth}px, 1fr))`

  if (gap) {
    el.style.gap = gap
  } else {
    el.style.gap = store.isTopLayout ? '1.5rem' : '1rem'
  }
}
