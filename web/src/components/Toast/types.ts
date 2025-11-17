/**
 * Toast通知类型
 */
export type ToastType = 'success' | 'error' | 'info' | 'warning'

/**
 * Toast组件属性
 */
export interface ToastProps {
  type: ToastType

  message: string

  duration?: number

  visible: boolean
}

/**
 * Toast配置选项
 */
export interface ToastOptions {
  message: string

  type?: ToastType

  duration?: number
}

/**
 * Toast服务接口
 */
export interface ToastService {
  show: (options: ToastOptions) => void

  success: (message: string, duration?: number) => void

  error: (message: string, duration?: number) => void

  info: (message: string, duration?: number) => void

  warning: (message: string, duration?: number) => void
}
