import { onMounted, onUnmounted, ref, watch, type Ref } from 'vue'

interface CountAnimationOptions {
  duration?: number // 动画持续时间（毫秒）
  startValue?: number // 起始值
  decimalPlaces?: number // 小数位数
  separator?: string // 千位分隔符
  prefix?: string // 前缀
  suffix?: string // 后缀
  easing?: (t: number) => number // 缓动函数
  autoStart?: boolean // 是否自动开始
  startOnVisible?: boolean // 是否在可见时开始
}

/**
 * 数字动画 composable
 * @param endValue 目标值
 * @param options 配置选项
 */
export function useCountAnimation(endValue: Ref<number> | number, options: CountAnimationOptions = {}) {
  const {
    duration = 2000,
    startValue = 0,
    decimalPlaces = 0,
    separator = ',',
    prefix = '',
    suffix = '',
    easing = easeOutCubic,
    autoStart = true,
    startOnVisible = false,
  } = options

  const displayValue = ref(`${prefix}${startValue}${suffix}`)
  const currentValue = ref(startValue)
  const isAnimating = ref(false)
  let animationId: number | null = null
  let startTime: number | null = null
  let observer: IntersectionObserver | null = null
  const hasStarted = ref(false)

  function easeOutCubic(t: number): number {
    return 1 - Math.pow(1 - t, 3)
  }

  function formatNumber(num: number): string {
    const fixed = num.toFixed(decimalPlaces)
    if (separator) {
      const parts = fixed.split('.')
      parts[0] = parts[0].replace(/\B(?=(\d{3})+(?!\d))/g, separator)
      return parts.join('.')
    }
    return fixed
  }

  function animate(timestamp: number) {
    if (!startTime) {
      startTime = timestamp
    }
    const progress = Math.min((timestamp - startTime) / duration, 1)

    const end = typeof endValue === 'object' ? endValue.value : endValue
    const easedProgress = easing(progress)
    currentValue.value = startValue + (end - startValue) * easedProgress

    displayValue.value = `${prefix}${formatNumber(currentValue.value)}${suffix}`

    if (progress < 1) {
      animationId = requestAnimationFrame(animate)
    } else {
      isAnimating.value = false
      startTime = null
    }
  }

  function start() {
    if (isAnimating.value || hasStarted.value) {
      return
    }

    hasStarted.value = true
    isAnimating.value = true
    startTime = null
    currentValue.value = startValue
    animationId = requestAnimationFrame(animate)
  }

  function stop() {
    if (animationId) {
      cancelAnimationFrame(animationId)
      animationId = null
    }
    isAnimating.value = false
    startTime = null
  }

  function reset() {
    stop()
    hasStarted.value = false
    currentValue.value = startValue
    displayValue.value = `${prefix}${formatNumber(startValue)}${suffix}`
  }

  function finish() {
    stop()
    const end = typeof endValue === 'object' ? endValue.value : endValue
    currentValue.value = end
    displayValue.value = `${prefix}${formatNumber(end)}${suffix}`
  }

  if (typeof endValue === 'object') {
    watch(endValue, () => {
      if (hasStarted.value) {
        reset()
        start()
      }
    })
  }

  function createObserver(element: Element) {
    if (!startOnVisible) {
      return
    }

    observer = new IntersectionObserver(
      (entries) => {
        entries.forEach((entry) => {
          if (entry.isIntersecting && !hasStarted.value) {
            start()
          }
        })
      },
      {
        threshold: 0.1, // 当元素10%可见时触发
      }
    )

    observer.observe(element)
  }

  function observeElement(element: Element) {
    if (observer) {
      observer.disconnect()
    }
    createObserver(element)
  }

  onMounted(() => {
    if (autoStart && !startOnVisible) {
      start()
    }
  })

  onUnmounted(() => {
    stop()
    if (observer) {
      observer.disconnect()
      observer = null
    }
  })

  return {
    displayValue,
    currentValue,
    isAnimating,
    start,
    stop,
    reset,
    finish,
    observeElement,
  }
}

export const easings = {
  linear: (t: number) => t,
  easeInQuad: (t: number) => t * t,
  easeOutQuad: (t: number) => t * (2 - t),
  easeInOutQuad: (t: number) => (t < 0.5 ? 2 * t * t : -1 + (4 - 2 * t) * t),
  easeInCubic: (t: number) => t * t * t,
  easeOutCubic: (t: number) => 1 - Math.pow(1 - t, 3),
  easeInOutCubic: (t: number) => (t < 0.5 ? 4 * t * t * t : (t - 1) * (2 * t - 2) * (2 * t - 2) + 1),
  easeInQuart: (t: number) => t * t * t * t,
  easeOutQuart: (t: number) => 1 - Math.pow(1 - t, 4),
  easeInOutQuart: (t: number) => (t < 0.5 ? 8 * t * t * t * t : 1 - 8 * Math.pow(t - 1, 4)),
  easeOutBounce: (t: number) => {
    if (t < 1 / 2.75) {
      return 7.5625 * t * t
    } else if (t < 2 / 2.75) {
      t -= 1.5 / 2.75
      return 7.5625 * t * t + 0.75
    } else if (t < 2.5 / 2.75) {
      t -= 2.25 / 2.75
      return 7.5625 * t * t + 0.9375
    }
    t -= 2.625 / 2.75
    return 7.5625 * t * t + 0.984375
  },
}
