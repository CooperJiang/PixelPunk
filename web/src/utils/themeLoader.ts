import type { VisualTheme } from '@/store/visualTheme'

const loadedThemes = new Set<VisualTheme>()
const loadingThemes = new Map<VisualTheme, Promise<void>>()
const DEFAULT_THEMES: VisualTheme[] = ['dark', 'light', 'cyberpunk-cyan']
const themeLinkElements = new Map<VisualTheme, HTMLLinkElement>()

const themeModules = import.meta.glob('@/styles/design-system/theme-*.css', {
  eager: false,
  query: '?url',
  import: 'default',
})

export async function loadTheme(theme: VisualTheme): Promise<void> {
  if (DEFAULT_THEMES.includes(theme)) {
    return
  }

  if (loadedThemes.has(theme)) {
    return
  }

  if (loadingThemes.has(theme)) {
    return loadingThemes.get(theme)
  }

  const loadPromise = new Promise<void>((resolve, reject) => {
    ;(async () => {
      try {
        const possiblePaths = [`/src/styles/design-system/theme-${theme}.css`, `@/styles/design-system/theme-${theme}.css`]

        let themeLoader = null

        for (const path of possiblePaths) {
          if (themeModules[path]) {
            themeLoader = themeModules[path]
            break
          }
        }

        if (!themeLoader) {
          throw new Error(`Theme module not found: ${theme}`)
        }

        const themeUrl = (await themeLoader()) as string
        const link = document.createElement('link')
        link.rel = 'stylesheet'
        link.href = themeUrl
        link.setAttribute('data-theme', theme)

        link.onload = () => {
          loadedThemes.add(theme)
          themeLinkElements.set(theme, link)
          loadingThemes.delete(theme)
          resolve()
        }

        link.onerror = () => {
          document.head.removeChild(link)
          loadingThemes.delete(theme)
          reject(new Error(`Failed to load theme: ${theme}`))
        }

        document.head.appendChild(link)
      } catch (error) {
        loadingThemes.delete(theme)
        reject(error)
      }
    })()
  })

  loadingThemes.set(theme, loadPromise)
  return loadPromise
}

export function preloadThemes(themes: VisualTheme[]): void {
  if (typeof requestIdleCallback !== 'undefined') {
    requestIdleCallback(
      () => {
        themes.forEach((theme) => {
          loadTheme(theme).catch(() => {})
        })
      },
      { timeout: 2000 }
    )
  } else {
    setTimeout(() => {
      themes.forEach((theme) => {
        loadTheme(theme).catch(() => {})
      })
    }, 1000)
  }
}

export function isThemeLoaded(theme: VisualTheme): boolean {
  return DEFAULT_THEMES.includes(theme) || loadedThemes.has(theme)
}

export function getLoadedThemes(): VisualTheme[] {
  return [...DEFAULT_THEMES, ...Array.from(loadedThemes)]
}

export function clearThemeCache(): void {
  loadedThemes.clear()
  loadingThemes.clear()
}
