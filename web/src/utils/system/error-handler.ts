/**
 * 全局错误处理器
 * 统一处理各种类型的错误和异常
 */
import type { TranslationFunction } from '@/composables/useTexts'
import { useToast } from '@/components/Toast/useToast'
import { HTTP_STATUS } from '@/constants'
import { logger } from './logger'

/* 错误类型枚举 */
export enum ErrorType {
  API_ERROR = 'API_ERROR',
  NETWORK_ERROR = 'NETWORK_ERROR',
  RUNTIME_ERROR = 'RUNTIME_ERROR',
  VALIDATION_ERROR = 'VALIDATION_ERROR',
  AUTH_ERROR = 'AUTH_ERROR',
  PERMISSION_ERROR = 'PERMISSION_ERROR',
  UNKNOWN_ERROR = 'UNKNOWN_ERROR',
}

/* 错误级别枚举 */
export enum ErrorLevel {
  LOW = 'low',
  MEDIUM = 'medium',
  HIGH = 'high',
  CRITICAL = 'critical',
}

/* 错误信息接口 */
export interface ErrorInfo {
  type: ErrorType
  level: ErrorLevel
  message: string
  detail?: any
  stack?: string
  timestamp: number
  url?: string
  userAgent?: string
  userId?: string
  requestId?: string
}

/* 错误处理配置 */
interface ErrorHandlerConfig {
  enableConsoleLog: boolean
  enableToast: boolean
  enableReport: boolean
  reportUrl?: string
  maxErrors: number
  reportThrottle: number
}

/* 默认配置 */
const defaultConfig: ErrorHandlerConfig = {
  enableConsoleLog: true,
  enableToast: true,
  enableReport: false,
  maxErrors: 100,
  reportThrottle: 5000,
}

class GlobalErrorHandler {
  private config: ErrorHandlerConfig
  private errorQueue: ErrorInfo[] = []
  private lastReportTime = 0

  constructor(config: Partial<ErrorHandlerConfig> = {}) {
    this.config = { ...defaultConfig, ...config }
    this.init()
  }

  private init() {
    if (typeof window !== 'undefined') {
      window.addEventListener('error', this.handleGlobalError.bind(this))

      window.addEventListener('unhandledrejection', this.handlePromiseError.bind(this))

      window.addEventListener('error', this.handleResourceError.bind(this), true)
    }
  }

  private handleGlobalError(event: ErrorEvent) {
    const errorInfo: ErrorInfo = {
      type: ErrorType.RUNTIME_ERROR,
      level: ErrorLevel.HIGH,
      message: event.message,
      detail: {
        filename: event.filename,
        lineno: event.lineno,
        colno: event.colno,
      },
      stack: event.error?.stack,
      timestamp: Date.now(),
      url: window.location.href,
      userAgent: navigator.userAgent,
    }

    this.handleError(errorInfo)
  }

  private handlePromiseError(event: PromiseRejectionEvent) {
    const errorInfo: ErrorInfo = {
      type: ErrorType.RUNTIME_ERROR,
      level: ErrorLevel.MEDIUM,
      message: event.reason?.message || 'Unhandled Promise Rejection',
      detail: event.reason,
      stack: event.reason?.stack,
      timestamp: Date.now(),
      url: window.location.href,
      userAgent: navigator.userAgent,
    }

    this.handleError(errorInfo)
    event.preventDefault() // 阻止默认的控制台输出
  }

  private handleResourceError(event: Event) {
    const target = event.target as HTMLElement
    if (target && target !== window) {
      const errorInfo: ErrorInfo = {
        type: ErrorType.NETWORK_ERROR,
        level: ErrorLevel.LOW,
        message: `Resource loading failed: ${target.tagName}`,
        detail: {
          tagName: target.tagName,
          src: (target as any).src || (target as any).href,
          currentSrc: (target as any).currentSrc,
        },
        timestamp: Date.now(),
        url: window.location.href,
      }

      this.handleError(errorInfo)
    }
  }

  public handleError(errorInfo: ErrorInfo) {
    if (this.config.enableConsoleLog) {
      this.logError(errorInfo)
    }

    if (this.config.enableToast && this.shouldShowToast(errorInfo)) {
      this.showErrorToast(errorInfo)
    }

    if (this.config.enableReport) {
      this.reportError(errorInfo)
    }

    this.addToQueue(errorInfo)
  }

  public handleApiError(_error: error, requestConfig?: any, $t?: TranslationFunction): ErrorInfo {
    const status = error.response?.status
    const message = error.response?.data?.message || error.message

    let errorType = ErrorType.API_ERROR
    let errorLevel = ErrorLevel.MEDIUM

    switch (status) {
      case HTTP_STATUS.UNAUTHORIZED:
        errorType = ErrorType.AUTH_ERROR
        errorLevel = ErrorLevel.HIGH
        break
      case HTTP_STATUS.FORBIDDEN:
        errorType = ErrorType.PERMISSION_ERROR
        errorLevel = ErrorLevel.HIGH
        break
      case HTTP_STATUS.NOT_FOUND:
        errorLevel = ErrorLevel.LOW
        break
      case HTTP_STATUS.INTERNAL_SERVER_ERROR:
        errorLevel = ErrorLevel.CRITICAL
        break
      case HTTP_STATUS.TOO_MANY_REQUESTS:
        errorLevel = ErrorLevel.MEDIUM
        break
    }

    const errorInfo: ErrorInfo = {
      type: errorType,
      level: errorLevel,
      message: this.formatApiErrorMessage(status, message, $t),
      detail: {
        status,
        url: requestConfig?.url,
        method: requestConfig?.method,
        data: requestConfig?.data,
        response: error.response?.data,
      },
      timestamp: Date.now(),
      requestId: error.response?.headers?.['x-request-id'],
    }

    this.handleError(errorInfo)
    return errorInfo
  }

  public handleValidationError(errors: Record<string, string[]>, context?: string): ErrorInfo {
    const errorInfo: ErrorInfo = {
      type: ErrorType.VALIDATION_ERROR,
      level: ErrorLevel.LOW,
      message: `Validation failed${context ? ` in ${context}` : ''}`,
      detail: errors,
      timestamp: Date.now(),
    }

    this.handleError(errorInfo)
    return errorInfo
  }

  /* 格式化 API 错误消息（带翻译） */
  private formatApiErrorMessage(status: number, message: string, $t?: TranslationFunction): string {
    if ($t) {
      const statusMessages: Record<number, string> = {
        [HTTP_STATUS.BAD_REQUEST]: $t('utils.errorHandler.httpStatus.400'),
        [HTTP_STATUS.UNAUTHORIZED]: $t('utils.errorHandler.httpStatus.401'),
        [HTTP_STATUS.FORBIDDEN]: $t('utils.errorHandler.httpStatus.403'),
        [HTTP_STATUS.NOT_FOUND]: $t('utils.errorHandler.httpStatus.404'),
        [HTTP_STATUS.TOO_MANY_REQUESTS]: $t('utils.errorHandler.httpStatus.429'),
        [HTTP_STATUS.INTERNAL_SERVER_ERROR]: $t('utils.errorHandler.httpStatus.500'),
        [HTTP_STATUS.BAD_GATEWAY]: $t('utils.errorHandler.httpStatus.502'),
        [HTTP_STATUS.SERVICE_UNAVAILABLE]: $t('utils.errorHandler.httpStatus.503'),
      }
      return statusMessages[status] || message || $t('utils.errorHandler.fallback.unknown')
    }

    // Backward compatibility - default English fallback
    const statusMessages: Record<number, string> = {
      [HTTP_STATUS.BAD_REQUEST]: 'Bad request',
      [HTTP_STATUS.UNAUTHORIZED]: 'Unauthorized or session expired',
      [HTTP_STATUS.FORBIDDEN]: 'Forbidden',
      [HTTP_STATUS.NOT_FOUND]: 'Resource not found',
      [HTTP_STATUS.TOO_MANY_REQUESTS]: 'Too many requests, please try again later',
      [HTTP_STATUS.INTERNAL_SERVER_ERROR]: 'Internal server error',
      [HTTP_STATUS.BAD_GATEWAY]: 'Bad gateway',
      [HTTP_STATUS.SERVICE_UNAVAILABLE]: 'Service unavailable',
    }

    return statusMessages[status] || message || 'Unknown error'
  }

  private logError(errorInfo: ErrorInfo) {
    const logMessage = `[${errorInfo.type}] ${errorInfo.message}`

    const logData = {
      ...errorInfo,
      formattedMessage: `🚨 ${logMessage}`,
    }

    switch (errorInfo.level) {
      case ErrorLevel.LOW:
        logger.info(logMessage, logData)
        break
      case ErrorLevel.MEDIUM:
        logger.warn(logMessage, logData)
        break
      case ErrorLevel.HIGH:
      case ErrorLevel.CRITICAL:
        logger.error(logMessage, logData)
        break
      default:
        logger.error(logMessage, logData)
    }
  }

  private getConsoleLogLevel(level: ErrorLevel): 'log' | 'warn' | 'error' {
    switch (level) {
      case ErrorLevel.LOW:
        return 'log'
      case ErrorLevel.MEDIUM:
        return 'warn'
      case ErrorLevel.HIGH:
      case ErrorLevel.CRITICAL:
        return 'error'
      default:
        return 'error'
    }
  }

  private shouldShowToast(errorInfo: ErrorInfo): boolean {
    if (errorInfo.level === ErrorLevel.LOW) {
      return false
    }

    if (errorInfo.type === ErrorType.NETWORK_ERROR) {
      return false
    }

    return true
  }

  private showErrorToast(errorInfo: ErrorInfo) {
    try {
      const { showToast } = useToast()
      const toastType = errorInfo.level === ErrorLevel.CRITICAL ? 'error' : 'warning'

      showToast({
        type: toastType,
        message: errorInfo.message,
        duration: errorInfo.level === ErrorLevel.CRITICAL ? 0 : 4000, // 严重错误不自动关闭
      })
    } catch {
      console.error('Toast unavailable, error message:', errorInfo.message)
    }
  }

  private reportError(errorInfo: ErrorInfo) {
    const now = Date.now()

    if (now - this.lastReportTime < this.config.reportThrottle) {
      return
    }

    this.lastReportTime = now

    if (errorInfo.level === ErrorLevel.LOW) {
      return
    }

    if (this.config.reportUrl) {
      fetch(this.config.reportUrl, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(errorInfo),
      }).catch(() => {})
    }
  }

  private addToQueue(errorInfo: ErrorInfo) {
    this.errorQueue.push(errorInfo)

    if (this.errorQueue.length > this.config.maxErrors) {
      this.errorQueue.shift()
    }
  }

  public getErrorStats() {
    const stats = {
      total: this.errorQueue.length,
      byType: {} as Record<ErrorType, number>,
      byLevel: {} as Record<ErrorLevel, number>,
      recent: this.errorQueue.slice(-10),
    }

    this.errorQueue.forEach((error) => {
      stats.byType[error.type] = (stats.byType[error.type] || 0) + 1
      stats.byLevel[error.level] = (stats.byLevel[error.level] || 0) + 1
    })

    return stats
  }

  public clearErrors() {
    this.errorQueue = []
  }

  public updateConfig(config: Partial<ErrorHandlerConfig>) {
    this.config = { ...this.config, ...config }
  }
}

export const globalErrorHandler = new GlobalErrorHandler({
  enableConsoleLog: import.meta.env.DEV,
  enableToast: true,
  enableReport: import.meta.env.PROD,
  reportUrl: import.meta.env.VITE_ERROR_REPORT_URL,
})

export const errorHandlerPlugin = {
  install(app: any) {
    app.config.errorHandler = (err: Error, instance: Event, info: string) => {
      const errorInfo: ErrorInfo = {
        type: ErrorType.RUNTIME_ERROR,
        level: ErrorLevel.HIGH,
        message: err.message,
        detail: {
          componentName: instance?.$options?.name || 'Unknown',
          info,
          props: instance?.$props,
        },
        stack: err.stack,
        timestamp: Date.now(),
        url: typeof window !== 'undefined' ? window.location.href : undefined,
      }

      globalErrorHandler.handleError(errorInfo)
    }

    app.config.warnHandler = (msg: string, instance: Event, trace: string) => {
      if (import.meta.env.DEV) {
        console.warn(`[Vue Warn]: ${msg}`, { instance, trace })
      }
    }
  },
}

export const reportError = globalErrorHandler.handleError.bind(globalErrorHandler)
export const reportApiError = globalErrorHandler.handleApiError.bind(globalErrorHandler)
export const reportValidationError = globalErrorHandler.handleValidationError.bind(globalErrorHandler)
