import { nextTick, onMounted, onUnmounted, ref } from 'vue'

export interface UseScrollNavigationOptions {
  threshold?: number
  debounceTime?: number
  transitionDuration?: string
}

export function useScrollNavigation(options: UseScrollNavigationOptions = {}) {
  const {
    threshold = 80, // 增加滚动阈值到80px
    debounceTime = 100,
    transitionDuration = '0.5s ease-in-out',
  } = options

  const isScrolled = ref(false)
  const showNavigation = ref(false)

  let debounceTimer: NodeJS.Timeout | null = null

  const handleScroll = () => {
    if (debounceTimer) {
      clearTimeout(debounceTimer)
    }

    debounceTimer = setTimeout(() => {
      const scrollTop = window.pageYOffset || document.documentElement.scrollTop || document.body.scrollTop
      const scrolled = scrollTop > threshold

      if (scrolled !== isScrolled.value) {
        isScrolled.value = scrolled
        showNavigation.value = scrolled
        toggleNavAndFooter(scrolled)
      }
    }, debounceTime)
  }

  const toggleNavAndFooter = (show: boolean) => {
    nextTick(() => {
      const navbar = document.querySelector('.cyber-navbar') as HTMLElement
      const footer = document.querySelector('.footer-container') as HTMLElement

      if (navbar) {
        applyStyles(navbar, show, 'translateY(-100%)', 'translateY(0)')
      }

      if (footer) {
        applyStyles(footer, show, 'translateY(100%)', 'translateY(0)')
      }
    })
  }

  const applyStyles = (element: HTMLElement, show: boolean, hideTransform: string, showTransform: string) => {
    const styles = show
      ? {
          opacity: '1',
          transform: showTransform,
          pointerEvents: 'auto',
        }
      : {
          opacity: '0',
          transform: hideTransform,
          pointerEvents: 'none',
        }

    Object.entries(styles).forEach(([property, value]) => {
      element.style.setProperty(property, value, 'important')
    })

    element.style.setProperty('transition', `all ${transitionDuration}`, 'important')
  }

  const immediateHideNavigation = () => {
    const navbar = document.querySelector('.cyber-navbar') as HTMLElement
    const footer = document.querySelector('.footer-container') as HTMLElement

    if (navbar) {
      navbar.style.setProperty('opacity', '0', 'important')
      navbar.style.setProperty('transform', 'translateY(-100%)', 'important')
      navbar.style.setProperty('pointer-events', 'none', 'important')
      navbar.style.setProperty('transition', 'none', 'important')
    }

    if (footer) {
      footer.style.setProperty('opacity', '0', 'important')
      footer.style.setProperty('transform', 'translateY(100%)', 'important')
      footer.style.setProperty('pointer-events', 'none', 'important')
      footer.style.setProperty('transition', 'none', 'important')
    }
  }

  const restoreNavigation = () => {
    const navbar = document.querySelector('.cyber-navbar') as HTMLElement
    const footer = document.querySelector('.footer-container') as HTMLElement

    const propertiesToRemove = ['opacity', 'transform', 'pointer-events', 'transition']

    if (navbar) {
      propertiesToRemove.forEach((prop) => navbar.style.removeProperty(prop))
    }

    if (footer) {
      propertiesToRemove.forEach((prop) => footer.style.removeProperty(prop))
    }
  }

  onMounted(() => {
    immediateHideNavigation()

    nextTick(() => {
      immediateHideNavigation()
    })

    setTimeout(() => {
      const navbar = document.querySelector('.cyber-navbar') as HTMLElement
      const footer = document.querySelector('.footer-container') as HTMLElement

      if (navbar) {
        navbar.style.setProperty('transition', `all ${transitionDuration}`, 'important')
      }
      if (footer) {
        footer.style.setProperty('transition', `all ${transitionDuration}`, 'important')
      }

      window.addEventListener('scroll', handleScroll, { passive: true })
    }, 100)
  })

  onUnmounted(() => {
    if (debounceTimer) {
      clearTimeout(debounceTimer)
    }

    window.removeEventListener('scroll', handleScroll)

    restoreNavigation()
  })

  return {
    isScrolled,
    showNavigation,
    toggleNavAndFooter,
    immediateHideNavigation,
    restoreNavigation,
  }
}
