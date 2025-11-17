/**
 * 动画工具函数
 */

/**
 * 数字动画滚动效果
 * @param element HTML元素
 * @param endValue 目标数值
 * @param duration 动画时长（毫秒），默认2000ms
 * @param formatter 数字格式化函数
 */
export function animateNumber(
  element: HTMLElement | null,
  endValue: number,
  duration: number = 2000,
  formatter?: (num: number) => string
): void {
  if (!element || endValue === undefined || endValue === null || isNaN(endValue)) {
    return
  }

  const startValue = 0
  const startTime = performance.now()
  const safeEndValue = Number(endValue) || 0

  const updateNumber = (currentTime: number) => {
    const elapsed = currentTime - startTime
    const progress = Math.min(elapsed / duration, 1)

    const easeOutCubic = 1 - Math.pow(1 - progress, 3)
    const currentValue = Math.floor(startValue + (safeEndValue - startValue) * easeOutCubic)

    if (element && element.textContent !== null) {
      element.textContent = formatter ? formatter(currentValue) : currentValue.toString()
    }

    if (progress < 1) {
      requestAnimationFrame(updateNumber)
    } else {
      if (element && element.textContent !== null) {
        element.textContent = formatter ? formatter(safeEndValue) : safeEndValue.toString()
      }
    }
  }

  requestAnimationFrame(updateNumber)
}

export function fadeIn(element: HTMLElement, duration: number = 300): Promise<void> {
  return new Promise((resolve) => {
    element.style.opacity = '0'
    element.style.transition = `opacity ${duration}ms ease`

    void element.offsetHeight

    element.style.opacity = '1'

    setTimeout(() => {
      element.style.transition = ''
      resolve()
    }, duration)
  })
}

export function slideIn(
  element: HTMLElement,
  direction: 'left' | 'right' | 'up' | 'down' = 'down',
  duration: number = 300
): Promise<void> {
  return new Promise((resolve) => {
    const transforms = {
      left: 'translateX(-100%)',
      right: 'translateX(100%)',
      up: 'translateY(-100%)',
      down: 'translateY(100%)',
    }

    element.style.transform = transforms[direction]
    element.style.transition = `transform ${duration}ms ease`

    void element.offsetHeight

    element.style.transform = 'translate(0, 0)'

    setTimeout(() => {
      element.style.transition = ''
      element.style.transform = ''
      resolve()
    }, duration)
  })
}

export function scaleIn(
  element: HTMLElement,
  fromScale: number = 0.8,
  toScale: number = 1,
  duration: number = 300
): Promise<void> {
  return new Promise((resolve) => {
    element.style.transform = `scale(${fromScale})`
    element.style.transition = `transform ${duration}ms ease`

    void element.offsetHeight

    element.style.transform = `scale(${toScale})`

    setTimeout(() => {
      element.style.transition = ''
      element.style.transform = ''
      resolve()
    }, duration)
  })
}

export async function staggerAnimation(
  elements: HTMLElement[],
  animationFn: (element: HTMLElement) => Promise<void>,
  delay: number = 100
): Promise<void> {
  for (let i = 0; i < elements.length; i++) {
    if (i > 0) {
      await new Promise((resolve) => setTimeout(resolve, delay))
    }
    animationFn(elements[i])
  }
}

export function cyberGlow(element: HTMLElement, glowColor: string = '#05d9e8', duration: number = 1000): void {
  const originalBoxShadow = element.style.boxShadow

  element.style.transition = `box-shadow ${duration}ms ease-in-out`
  element.style.boxShadow = `0 0 20px ${glowColor}, 0 0 40px ${glowColor}40`

  setTimeout(() => {
    element.style.boxShadow = originalBoxShadow
    setTimeout(() => {
      element.style.transition = ''
    }, duration)
  }, duration)
}

export default {
  animateNumber,
  fadeIn,
  slideIn,
  scaleIn,
  staggerAnimation,
  cyberGlow,
}
