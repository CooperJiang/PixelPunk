/**
 * Admin 动画管理 Composable
 * 提供统一的动画控制接口

 */
import { nextTick, onMounted, ref } from 'vue'

export interface AdminAnimationOptions {
  enablePageAnimation?: boolean
  enableStagger?: boolean
  staggerDelay?: number
  duration?: 'fast' | 'normal' | 'slow' | 'extra-slow'
  easing?: 'out' | 'in' | 'in-out' | 'back' | 'elastic'
}

export function useAdminAnimations(options: AdminAnimationOptions = {}) {
  const {
    enablePageAnimation = true,
    enableStagger = false,
    staggerDelay = 50,
    duration = 'normal',
    easing: _easing = 'out',
  } = options

  const isAnimating = ref(false)
  const animationContainer = ref<HTMLElement>()

  const triggerPageAnimation = async () => {
    if (!enablePageAnimation) {
      return
    }

    isAnimating.value = true
    await nextTick()

    if (animationContainer.value) {
      animationContainer.value.classList.add('admin-page-enter')

      setTimeout(() => {
        isAnimating.value = false
        animationContainer.value?.classList.remove('admin-page-enter')
      }, getDurationMs(duration))
    }
  }

  const triggerStaggerAnimation = async (selector: string = '.admin-stagger-item') => {
    if (!enableStagger) {
      return
    }

    await nextTick()

    const items = document.querySelectorAll(selector)
    items.forEach((item, index) => {
      const element = item as HTMLElement
      element.style.animationDelay = `${index * staggerDelay}ms`
      element.classList.add('admin-stagger-item')
    })
  }

  const addCardHoverAnimation = (element: HTMLElement) => {
    element.classList.add('admin-card')
  }

  const triggerButtonPress = (element: HTMLElement) => {
    element.classList.add('admin-btn')
    element.addEventListener(
      'click',
      () => {
        element.classList.add('admin-btn-press')
        setTimeout(() => {
          element.classList.remove('admin-btn-press')
        }, 150)
      },
      { once: true }
    )
  }

  const triggerSuccessAnimation = (element: HTMLElement) => {
    element.classList.add('admin-btn-success')
    setTimeout(() => {
      element.classList.remove('admin-btn-success')
    }, 600)
  }

  const triggerErrorShake = (element: HTMLElement) => {
    element.classList.add('admin-error-shake')
    setTimeout(() => {
      element.classList.remove('admin-error-shake')
    }, 500)
  }

  const addPulseAnimation = (element: HTMLElement) => {
    element.classList.add('admin-card-pulse')
  }

  const removePulseAnimation = (element: HTMLElement) => {
    element.classList.remove('admin-card-pulse')
  }

  const addSkeletonAnimation = (element: HTMLElement, dark = false) => {
    element.classList.add(dark ? 'admin-skeleton-dark' : 'admin-skeleton')
  }

  const removeSkeletonAnimation = (element: HTMLElement) => {
    element.classList.remove('admin-skeleton', 'admin-skeleton-dark')
  }

  function getDurationMs(duration: string): number {
    const durations = {
      fast: 150,
      normal: 250,
      slow: 400,
      'extra-slow': 600,
    }
    return durations[duration as keyof typeof durations] || 250
  }

  const applyCustomAnimation = (
    element: HTMLElement,
    animationName: string,
    duration: number = 250,
    easing: string = 'ease-out'
  ) => {
    element.style.animation = `${animationName} ${duration}ms ${easing}`

    setTimeout(() => {
      element.style.animation = ''
    }, duration)
  }

  const applyAnimationToElements = (selector: string, animationClass: string, delay: number = 0) => {
    const elements = document.querySelectorAll(selector)

    elements.forEach((element, index) => {
      setTimeout(() => {
        element.classList.add(animationClass)
      }, index * delay)
    })
  }

  const cleanupAnimations = (element: HTMLElement) => {
    const animationClasses = [
      'admin-page-enter',
      'admin-page-fade-enter',
      'admin-stagger-item',
      'admin-card',
      'admin-card-pulse',
      'admin-btn',
      'admin-btn-success',
      'admin-error-shake',
      'admin-error-border',
      'admin-skeleton',
      'admin-skeleton-dark',
    ]

    element.classList.remove(...animationClasses)
    element.style.animation = ''
    element.style.animationDelay = ''
  }

  const fadeTransition = {
    enterActiveClass: 'admin-fade-enter-active',
    leaveActiveClass: 'admin-fade-leave-active',
    enterFromClass: 'admin-fade-enter-from',
    leaveToClass: 'admin-fade-leave-to',
  }

  const slideDownTransition = {
    enterActiveClass: 'admin-slide-down-enter-active',
    leaveActiveClass: 'admin-slide-down-leave-active',
    enterFromClass: 'admin-slide-down-enter-from',
    leaveToClass: 'admin-slide-down-leave-to',
  }

  const scaleTransition = {
    enterActiveClass: 'admin-scale-enter-active',
    leaveActiveClass: 'admin-scale-leave-active',
    enterFromClass: 'admin-scale-enter-from',
    leaveToClass: 'admin-scale-leave-to',
  }

  const modalTransition = {
    backdrop: {
      enterActiveClass: 'admin-modal-backdrop-enter-active',
      leaveActiveClass: 'admin-modal-backdrop-leave-active',
      enterFromClass: 'admin-modal-backdrop-enter-from',
      leaveToClass: 'admin-modal-backdrop-leave-to',
    },
    content: {
      enterActiveClass: 'admin-modal-content-enter-active',
      leaveActiveClass: 'admin-modal-content-leave-active',
      enterFromClass: 'admin-modal-content-enter-from',
      leaveToClass: 'admin-modal-content-leave-to',
    },
  }

  onMounted(() => {
    if (enablePageAnimation) {
      triggerPageAnimation()
    }
    if (enableStagger) {
      triggerStaggerAnimation()
    }
  })

  return {
    isAnimating,
    animationContainer,

    triggerPageAnimation,
    triggerStaggerAnimation,
    addCardHoverAnimation,
    triggerButtonPress,
    triggerSuccessAnimation,
    triggerErrorShake,
    addPulseAnimation,
    removePulseAnimation,
    addSkeletonAnimation,
    removeSkeletonAnimation,
    applyCustomAnimation,
    applyAnimationToElements,
    cleanupAnimations,

    fadeTransition,
    slideDownTransition,
    scaleTransition,
    modalTransition,

    getDurationMs,
  }
}

export default useAdminAnimations
