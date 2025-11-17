/**
 * DOM操作相关工具函数
 */

/**
 * 切换全屏模式
 * @param element 要全屏显示的元素，默认为document.documentElement
 * @returns Promise<boolean> 是否成功进入/退出全屏
 */
export const toggleFullScreen = async (element: HTMLElement = document.documentElement): Promise<boolean> => {
  try {
    if (!isFullScreen()) {
      if (element.requestFullscreen) {
        await element.requestFullscreen()
      } else if ((element as any).mozRequestFullScreen) {
        await (element as any).mozRequestFullScreen()
      } else if ((element as any).webkitRequestFullscreen) {
        await (element as any).webkitRequestFullscreen()
      } else if ((element as any).msRequestFullscreen) {
        await (element as any).msRequestFullscreen()
      } else {
        return false
      }
      return true
    }
    if (document.exitFullscreen) {
      await document.exitFullscreen()
    } else if ((document as any).mozCancelFullScreen) {
      await (document as any).mozCancelFullScreen()
    } else if ((document as any).webkitExitFullscreen) {
      await (document as any).webkitExitFullscreen()
    } else if ((document as any).msExitFullscreen) {
      await (document as any).msExitFullscreen()
    } else {
      return false
    }
    return true
  } catch (error) {
    console.error('切换全屏模式失败:', error)
    return false
  }
}

export const isFullScreen = (): boolean =>
  Boolean(
    document.fullscreenElement ||
      (document as any).mozFullScreenElement ||
      (document as any).webkitFullscreenElement ||
      (document as any).msFullscreenElement
  )

export const saveScrollPosition = (key: string, positions: Record<string, number>): Record<string, number> => {
  const newPositions = { ...positions }
  newPositions[key] = window.scrollY
  return newPositions
}

export const restoreScrollPosition = (key: string, positions: Record<string, number>, smooth: boolean = false): void => {
  const savedPosition = positions[key]

  if (savedPosition !== undefined) {
    requestAnimationFrame(() => {
      window.scrollTo({
        top: savedPosition,
        behavior: smooth ? 'smooth' : 'auto',
      })
    })
  } else {
    window.scrollTo(0, 0)
  }
}

export const setCssVariable = (
  name: string,
  value: string,
  element: HTMLElement = document.documentElement,
  debug: boolean = false
): void => {
  element.style.setProperty(`--${name}`, value)
  if (debug) {
  }
}

export const getCssVariable = (name: string, element: HTMLElement = document.documentElement): string =>
  getComputedStyle(element).getPropertyValue(`--${name}`).trim()

export const addGlobalClass = (className: string, element: HTMLElement = document.body): void => {
  element.classList.add(className)
}

export const removeGlobalClass = (className: string, element: HTMLElement = document.body): void => {
  element.classList.remove(className)
}

export const toggleGlobalClass = (className: string, force?: boolean, element: HTMLElement = document.body): boolean => {
  if (force !== undefined) {
    if (force) {
      addGlobalClass(className, element)
      return true
    }
    removeGlobalClass(className, element)
    return false
  }
  element.classList.toggle(className)
  return element.classList.contains(className)
}
