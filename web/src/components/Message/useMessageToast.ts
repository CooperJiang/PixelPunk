import { createApp, type App } from 'vue'
import MessageToast from './MessageToast.vue'
import router from '@/router'

interface ToastOptions {
  type: 'success' | 'error' | 'warning' | 'info'
  title: string
  message?: string
  duration?: number
  actionText?: string
  actionUrl?: string
  actionStyle?: 'primary' | 'secondary' | 'success' | 'warning' | 'danger'
  id?: string // 添加唯一标识
}

interface ToastInstance {
  id: string
  app: App
  close: () => void
  mountPoint: HTMLElement
}

class ToastManager {
  private toasts = new Map<string, ToastInstance>()
  private container: HTMLElement | null = null
  private toastIds = new Set<string>() // 跟踪所有toast ID

  private generateId(): string {
    return `toast-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`
  }

  private ensureContainer(): void {
    if (!this.container) {
      this.container = document.createElement('div')
      this.container.id = 'message-toast-container'
      this.container.className = 'fixed top-0 right-0 z-[9999] p-4 space-y-3 pointer-events-none'
      this.container.style.cssText = `
        position: fixed;
        top: 1rem;
        right: 1rem;
        z-index: 9999;
        display: flex;
        flex-direction: column;
        pointer-events: none;
        max-height: calc(100vh - 2rem);
        overflow-y: auto;
        overflow-x: visible;
        width: auto;
        min-width: 320px;
        max-width: 420px;
      `
      document.body.appendChild(this.container)
    }
  }

  private createToast(options: ToastOptions): ToastInstance {
    this.ensureContainer()

    const id = options.id || this.generateId()

    if (this.toastIds.has(id)) {
      const existingToast = this.toasts.get(id)
      if (existingToast) {
        return existingToast
      }
    }

    this.toastIds.add(id)

    const mountPoint = document.createElement('div')
    mountPoint.dataset.toastId = id
    mountPoint.style.cssText = `
      pointer-events: auto;
      margin-bottom: 0.75rem;
      transform: translateX(100%);
      opacity: 0;
      transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    `

    if (!this.container) return

    if (this.container.firstChild) {
      this.container.insertBefore(mountPoint, this.container.firstChild)
    } else {
      this.container.appendChild(mountPoint)
    }

    const duration = options.duration || 4000

    const app = createApp(MessageToast, {
      ...options,
      id,
      duration,
      onClose: () => {
        this.removeToast(id)
      },
    })

    app.use(router)

    app.mount(mountPoint)

    requestAnimationFrame(() => {
      mountPoint.style.transform = 'translateX(0)'
      mountPoint.style.opacity = '1'
    })

    const instance: ToastInstance = {
      id,
      app,
      close: () => {
        this.removeToast(id)
      },
      mountPoint,
    }

    this.toasts.set(id, instance)

    return instance
  }

  private removeToast(id: string): void {
    const toast = this.toasts.get(id)
    if (!toast) {
      return
    }

    const { mountPoint } = toast
    if (mountPoint) {
      setTimeout(() => {
        try {
          toast.app.unmount()
          if (mountPoint.parentElement) {
            mountPoint.remove()
          }
        } catch {}

        if (this.container && this.container.children.length === 0) {
          this.container.remove()
          this.container = null
        }
      }, 0)
    }

    this.toasts.delete(id)
    this.toastIds.delete(id) // 从ID集合中删除
  }

  success(title: string, message?: string, options?: Partial<ToastOptions>): ToastInstance {
    return this.createToast({
      type: 'success',
      title,
      message,
      duration: 4000,
      ...options,
    })
  }

  error(title: string, message?: string, options?: Partial<ToastOptions>): ToastInstance {
    return this.createToast({
      type: 'error',
      title,
      message,
      duration: 6000,
      ...options,
    })
  }

  warning(title: string, message?: string, options?: Partial<ToastOptions>): ToastInstance {
    return this.createToast({
      type: 'warning',
      title,
      message,
      duration: 5000,
      ...options,
    })
  }

  info(title: string, message?: string, options?: Partial<ToastOptions>): ToastInstance {
    return this.createToast({
      type: 'info',
      title,
      message,
      duration: 4000,
      ...options,
    })
  }

  showMessage(messageData: {
    title: string
    content?: string
    type?: string
    action_url?: string
    action_text?: string
    action_style?: string
    id?: string // 添加可选的ID参数
  }): ToastInstance {
    let toastType: ToastOptions['type'] = 'info'

    if (messageData.type) {
      if (messageData.type.includes('success') || messageData.type.includes('approved') || messageData.type.includes('granted')) {
        toastType = 'success'
      } else if (
        messageData.type.includes('error') ||
        messageData.type.includes('rejected') ||
        messageData.type.includes('alert')
      ) {
        toastType = 'error'
      } else if (messageData.type.includes('warning') || messageData.type.includes('quota')) {
        toastType = 'warning'
      } else {
        toastType = 'info'
      }
    }

    const messageId = messageData.id || `msg-${messageData.title}-${Date.now()}-${Math.random().toString(36).substr(2, 5)}`

    let duration = 4000 // 默认4秒
    if (toastType === 'error') {
      duration = 6000
    } else if (toastType === 'warning') {
      duration = 5000
    }

    return this.createToast({
      type: toastType,
      title: messageData.title,
      message: messageData.content,
      actionText: messageData.action_text,
      actionUrl: messageData.action_url,
      actionStyle: (messageData.action_style as ToastOptions['actionStyle']) || 'primary',
      id: messageId,
      duration, // 明确传递duration
    })
  }

  clear(): void {
    this.toasts.forEach((toast) => {
      toast.close()
    })
  }
}

export const messageToast = new ToastManager()

export type { ToastOptions, ToastInstance }

export function useMessageToast() {
  return {
    success: messageToast.success.bind(messageToast),
    error: messageToast.error.bind(messageToast),
    warning: messageToast.warning.bind(messageToast),
    info: messageToast.info.bind(messageToast),
    showMessage: messageToast.showMessage.bind(messageToast),
    clear: messageToast.clear.bind(messageToast),
  }
}

export default messageToast
