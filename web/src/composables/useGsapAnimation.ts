import { onMounted, onUnmounted, ref, type Ref } from 'vue'
import gsap from 'gsap'

/**
 * GSAP 动画 composable
 * 提供常用的动画效果和工具函数
 */
export function useGsapAnimation() {
  const ctx = ref<gsap.Context>()

  const fadeIn = (element: Ref<HTMLElement | undefined> | HTMLElement | string, options = {}) => {
    const target = typeof element === 'string' ? element : (element as any).value || element
    if (!target) {
      return
    }

    return gsap.from(target, {
      opacity: 0,
      y: 30,
      duration: 1,
      ease: 'power3.out',
      ...options,
    })
  }

  const fadeOut = (element: Ref<HTMLElement | undefined> | HTMLElement | string, options = {}) => {
    const target = typeof element === 'string' ? element : (element as any).value || element
    if (!target) {
      return
    }

    return gsap.to(target, {
      opacity: 0,
      y: -30,
      duration: 0.8,
      ease: 'power3.in',
      ...options,
    })
  }

  const scaleIn = (element: Ref<HTMLElement | undefined> | HTMLElement | string, options = {}) => {
    const target = typeof element === 'string' ? element : (element as any).value || element
    if (!target) {
      return
    }

    return gsap.from(target, {
      scale: 0,
      opacity: 0,
      duration: 0.5,
      ease: 'back.out(1.7)',
      ...options,
    })
  }

  const slideIn = (element: Ref<HTMLElement | undefined> | HTMLElement | string, direction = 'left', options = {}) => {
    const target = typeof element === 'string' ? element : (element as any).value || element
    if (!target) {
      return
    }

    const directionMap = {
      left: { x: -100 },
      right: { x: 100 },
      top: { y: -100 },
      bottom: { y: 100 },
    }

    return gsap.from(target, {
      ...directionMap[direction],
      opacity: 0,
      duration: 0.8,
      ease: 'power3.out',
      ...options,
    })
  }

  const glitchEffect = (element: Ref<HTMLElement | undefined> | HTMLElement | string, options = {}) => {
    const target = typeof element === 'string' ? element : (element as any).value || element
    if (!target) {
      return
    }

    const tl = gsap.timeline({
      repeat: 2,
      repeatDelay: 0.1,
      ...options,
    })

    tl.to(target, {
      skewX: 10,
      x: 5,
      duration: 0.1,
      ease: 'power2.inOut',
    })
      .to(target, {
        skewX: -10,
        x: -5,
        duration: 0.1,
        ease: 'power2.inOut',
      })
      .to(target, {
        skewX: 0,
        x: 0,
        duration: 0.1,
        ease: 'power2.inOut',
      })

    return tl
  }

  const neonFlicker = (element: Ref<HTMLElement | undefined> | HTMLElement | string, options = {}) => {
    const target = typeof element === 'string' ? element : (element as any).value || element
    if (!target) {
      return
    }

    return gsap.to(target, {
      opacity: 0.8,
      duration: 0.05,
      repeat: 10,
      yoyo: true,
      ease: 'none',
      onComplete: () => {
        gsap.set(target, { opacity: 1 })
      },
      ...options,
    })
  }

  const countUp = (element: Ref<HTMLElement | undefined> | HTMLElement, endValue: number, options = {}) => {
    const target = (element as any).value || element
    if (!target) {
      return
    }

    const obj = { value: 0 }

    return gsap.to(obj, {
      value: endValue,
      duration: 2,
      ease: 'power3.out',
      onUpdate: () => {
        target.textContent = Math.round(obj.value).toLocaleString()
      },
      ...options,
    })
  }

  const typewriter = (element: Ref<HTMLElement | undefined> | HTMLElement, text: string, options = {}) => {
    const target = (element as any).value || element
    if (!target) {
      return
    }

    target.textContent = ''

    return gsap.to(target, {
      duration: text.length * 0.05,
      text: {
        value: text,
        delimiter: '',
      },
      ease: 'none',
      ...options,
    })
  }

  const staggerAnimation = (selector: string, animationProps: any, staggerOptions = 0.1) =>
    gsap.from(selector, {
      ...animationProps,
      stagger: staggerOptions,
    })

  const createTimeline = (options = {}) => gsap.timeline(options)

  const mouseFollow = (element: Ref<HTMLElement | undefined> | HTMLElement | string, speed = 0.5) => {
    const target = typeof element === 'string' ? element : (element as any).value || element
    if (!target) {
      return
    }

    const handleMouseMove = (e: MouseEvent) => {
      gsap.to(target, {
        x: e.clientX,
        y: e.clientY,
        duration: speed,
        ease: 'power2.out',
      })
    }

    window.addEventListener('mousemove', handleMouseMove)

    return () => {
      window.removeEventListener('mousemove', handleMouseMove)
    }
  }

  const parallax = (element: Ref<HTMLElement | undefined> | HTMLElement | string, speed = 0.5) => {
    const target = typeof element === 'string' ? element : (element as any).value || element
    if (!target) {
      return
    }

    const handleScroll = () => {
      const scrolled = window.pageYOffset
      const yPos = -(scrolled * speed)

      gsap.to(target, {
        y: yPos,
        duration: 0.5,
        ease: 'power2.out',
      })
    }

    window.addEventListener('scroll', handleScroll)

    return () => {
      window.removeEventListener('scroll', handleScroll)
    }
  }

  const killAll = () => {
    ctx.value?.revert()
  }

  onMounted(() => {
    ctx.value = gsap.context(() => {})
  })

  onUnmounted(() => {
    killAll()
  })

  return {
    fadeIn,
    fadeOut,
    scaleIn,
    slideIn,
    glitchEffect,
    neonFlicker,
    countUp,
    typewriter,
    staggerAnimation,
    createTimeline,
    mouseFollow,
    parallax,
    killAll,

    gsap,
  }
}
