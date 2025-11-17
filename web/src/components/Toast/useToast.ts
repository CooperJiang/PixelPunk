import { h, reactive, ref, render } from 'vue'
import Toast from './index.vue'
import { TIMING, Z_INDEX } from '@/constants'
import type { ToastOptions } from './types'

const toastList = reactive<
  Array<{
    id: number
    el: HTMLElement
    height: number
    visible: boolean
    hide: () => void
    top: number
  }>
>([])

let toastIdCounter = 0

const calculateToastPositions = () => {
  let totalHeight = TIMING.TOAST.INITIAL_TOP
  const visibleToasts = toastList.filter((toast) => toast.visible)

  visibleToasts.forEach((toast) => {
    toast.top = totalHeight

    if (toast.el) {
      toast.el.style.transition = 'top 0.3s ease'
      toast.el.style.top = `${totalHeight}px`
    }

    totalHeight += (toast.height || 60) + TIMING.TOAST.MARGIN
  })
}

export function useToast() {
  const addToast = (options: ToastOptions) => {
    const { message, type = 'info', duration = TIMING.TOAST.DEFAULT_DURATION } = options
    const id = toastIdCounter++

    const container = document.createElement('div')
    container.className = 'cyber-toast-wrapper'
    container.style.position = 'fixed'
    container.style.left = '50%'
    container.style.transform = 'translateX(-50%)'
    container.style.zIndex = String(Z_INDEX.TOAST)
    container.style.top = `${TIMING.TOAST.INITIAL_TOP - 20}px`
    container.style.opacity = '0'
    container.style.transition = 'top 0.3s ease, opacity 0.3s ease'

    const visible = ref(true)

    const hide = () => {
      const index = toastList.findIndex((item) => item.id === id)
      if (index !== -1) {
        toastList[index].visible = false
        visible.value = false

        if (container && document.body.contains(container)) {
          container.style.opacity = '0'
          container.style.transform = 'translateX(-50%) translateY(-20px)'
        }

        calculateToastPositions()

        setTimeout(() => {
          if (document.body.contains(container)) {
            try {
              render(null, container)
              document.body.removeChild(container)
              const removeIndex = toastList.findIndex((item) => item.id === id)
              if (removeIndex !== -1) {
                toastList.splice(removeIndex, 1)
              }
            } catch {}
          }
        }, TIMING.TOAST.ANIMATION_DELAY)
      }
    }

    const toastInstance = {
      id,
      el: container,
      height: 0,
      visible: true,
      hide,
      top: 0,
    }

    document.body.appendChild(container)

    const vnode = h(Toast, {
      type,
      message,
      duration,
      visible: visible.value,
      onUpdate: (val: boolean) => {
        visible.value = val
      },
      onClose: () => {
        hide()
      },
    })

    render(vnode, container)
    toastList.push(toastInstance)

    setTimeout(() => {
      if (container.firstElementChild) {
        toastInstance.height = container.firstElementChild.getBoundingClientRect().height
        calculateToastPositions()
        container.style.opacity = '1'
        container.style.top = `${toastInstance.top}px`
      }
    }, TIMING.TOAST.SHOW_DELAY)

    if (duration > 0) {
      setTimeout(hide, duration)
    }

    return { close: hide }
  }

  const clearAll = () => {
    ;[...toastList].forEach((toast) => {
      toast.hide()
    })
  }

  return {
    success: (message: string, duration?: number) => addToast({ message, type: 'success', duration }),
    error: (message: string, duration?: number) => addToast({ message, type: 'error', duration }),
    info: (message: string, duration?: number) => addToast({ message, type: 'info', duration }),
    warning: (message: string, duration?: number) => addToast({ message, type: 'warning', duration }),
    show: addToast,
    clearAll,
  }
}
