/**
 * 统一日志管理器
 * 提供结构化的日志记录和管理功能
 */
import { DEBUG_CONFIG, isDevelopment, isProduction } from '@/constants/env'
import { StorageUtil } from '@/utils/storage'

/* 日志级别 */
export enum LogLevel {
  DEBUG = 0,
  INFO = 1,
  WARN = 2,
  ERROR = 3,
  FATAL = 4,
}

/* 日志类型 */
export enum LogType {
  SYSTEM = 'system',
  API = 'api',
  USER = 'user',
  PERFORMANCE = 'performance',
  SECURITY = 'security',
  BUSINESS = 'business',
}

/* 日志条目接口 */
export interface LogEntry {
  level: LogLevel
  type: LogType
  message: string
  data?: any
  timestamp: number
  source?: string
  userId?: string
  sessionId?: string
  requestId?: string
  stack?: string
}

/* 日志配置接口 */
interface LoggerConfig {
  level: LogLevel
  enableConsole: boolean
  enableStorage: boolean
  enableRemote: boolean
  remoteUrl?: string
  maxStorageEntries: number
  bufferSize: number
  flushInterval: number
}

/* 默认配置 */
const defaultConfig: LoggerConfig = {
  level: isDevelopment ? LogLevel.DEBUG : LogLevel.INFO,
  enableConsole: isDevelopment,
  enableStorage: true,
  enableRemote: isProduction,
  maxStorageEntries: 1000,
  bufferSize: 50,
  flushInterval: 30000, // 30秒
}

class Logger {
  private config: LoggerConfig
  private buffer: LogEntry[] = []
  private storageKey = 'app_logs'
  private sessionId: string
  private flushTimer?: NodeJS.Timeout

  constructor(config: Partial<LoggerConfig> = {}) {
    this.config = { ...defaultConfig, ...config }
    this.sessionId = this.generateSessionId()
    this.initFlushTimer()
  }

  private generateSessionId(): string {
    return `${Date.now()}-${Math.random().toString(36).substr(2, 9)}`
  }

  private initFlushTimer() {
    if (this.config.enableRemote && this.config.flushInterval > 0) {
      this.flushTimer = setInterval(() => {
        this.flush()
      }, this.config.flushInterval)
    }
  }

  private createLogEntry(level: LogLevel, type: LogType, message: string, data?: any, source?: string): LogEntry {
    return {
      level,
      type,
      message,
      data,
      timestamp: Date.now(),
      source,
      sessionId: this.sessionId,
      userId: this.getCurrentUserId(),
      requestId: this.getCurrentRequestId(),
      stack: level >= LogLevel.ERROR ? new Error().stack : undefined,
    }
  }

  private getCurrentUserId(): string | undefined {
    try {
      return undefined
    } catch {
      return undefined
    }
  }

  private getCurrentRequestId(): string | undefined {
    return undefined
  }

  private shouldLog(level: LogLevel): boolean {
    return level >= this.config.level
  }

  private formatConsoleOutput(entry: LogEntry): any[] {
    const timestamp = new Date(entry.timestamp).toISOString()
    const levelName = LogLevel[entry.level]
    const prefix = `[${timestamp}] [${levelName}] [${entry.type.toUpperCase()}]`

    const outputs = [`%c${prefix} ${entry.message}`, this.getConsoleStyle(entry.level)]

    if (entry.data) {
      outputs.push('\nData:', entry.data)
    }

    if (entry.source) {
      outputs.push('\nSource:', entry.source)
    }

    if (entry.stack && entry.level >= LogLevel.ERROR) {
      outputs.push('\nStack:', entry.stack)
    }

    return outputs
  }

  private getConsoleStyle(level: LogLevel): string {
    const styles = {
      [LogLevel.DEBUG]: 'color: #888; font-size: 12px;',
      [LogLevel.INFO]: 'color: #2196F3; font-weight: bold;',
      [LogLevel.WARN]: 'color: #FF9800; font-weight: bold;',
      [LogLevel.ERROR]: 'color: #F44336; font-weight: bold;',
      [LogLevel.FATAL]: 'color: #fff; background: #F44336; font-weight: bold; padding: 2px 4px;',
    }
    return styles[level] || ''
  }

  private outputToConsole(entry: LogEntry) {
    if (!this.config.enableConsole) {
      return
    }

    const outputs = this.formatConsoleOutput(entry)

    switch (entry.level) {
      case LogLevel.DEBUG:
        break
      case LogLevel.INFO:
        break
      case LogLevel.WARN:
        console.warn(...outputs)
        break
      case LogLevel.ERROR:
      case LogLevel.FATAL:
        console.error(...outputs)
        break
    }
  }

  private storeLocally(entry: LogEntry) {
    if (!this.config.enableStorage) {
      return
    }

    try {
      const existingLogs = this.getStoredLogs()
      existingLogs.push(entry)

      if (existingLogs.length > this.config.maxStorageEntries) {
        existingLogs.splice(0, existingLogs.length - this.config.maxStorageEntries)
      }

      StorageUtil.set(this.storageKey, existingLogs)
    } catch (error) {
      console.error('Failed to store log entry:', error)
    }
  }

  private getStoredLogs(): LogEntry[] {
    try {
      return StorageUtil.get<LogEntry[]>(this.storageKey) || []
    } catch {
      return []
    }
  }

  private addToBuffer(entry: LogEntry) {
    if (!this.config.enableRemote) {
      return
    }

    this.buffer.push(entry)

    if (this.buffer.length >= this.config.bufferSize) {
      this.flush()
    }
  }

  private async flush() {
    if (!this.config.enableRemote || !this.config.remoteUrl || this.buffer.length === 0) {
      return
    }

    const logsToSend = [...this.buffer]
    this.buffer = []

    try {
      await fetch(this.config.remoteUrl, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          logs: logsToSend,
          sessionId: this.sessionId,
          userAgent: navigator.userAgent,
          url: window.location.href,
        }),
      })
    } catch (error) {
      this.buffer.unshift(...logsToSend)
      console.error('Failed to send logs to remote:', error)
    }
  }

  private log(level: LogLevel, type: LogType, message: string, data?: any, source?: string) {
    if (!this.shouldLog(level)) {
      return
    }

    const entry = this.createLogEntry(level, type, message, data, source)

    this.outputToConsole(entry)

    this.storeLocally(entry)

    this.addToBuffer(entry)
  }

  debug(message: string, data?: any, source?: string) {
    this.log(LogLevel.DEBUG, LogType.SYSTEM, message, data, source)
  }

  info(message: string, data?: any, source?: string) {
    this.log(LogLevel.INFO, LogType.SYSTEM, message, data, source)
  }

  warn(message: string, data?: any, source?: string) {
    this.log(LogLevel.WARN, LogType.SYSTEM, message, data, source)
  }

  error(message: string, data?: any, source?: string) {
    this.log(LogLevel.ERROR, LogType.SYSTEM, message, data, source)
  }

  fatal(message: string, data?: any, source?: string) {
    this.log(LogLevel.FATAL, LogType.SYSTEM, message, data, source)
  }

  api(message: string, data?: any, source?: string) {
    this.log(LogLevel.INFO, LogType.API, message, data, source)
  }

  user(message: string, data?: any, source?: string) {
    this.log(LogLevel.INFO, LogType.USER, message, data, source)
  }

  performance(message: string, data?: any, source?: string) {
    this.log(LogLevel.INFO, LogType.PERFORMANCE, message, data, source)
  }

  security(message: string, data?: any, source?: string) {
    this.log(LogLevel.WARN, LogType.SECURITY, message, data, source)
  }

  business(message: string, data?: any, source?: string) {
    this.log(LogLevel.INFO, LogType.BUSINESS, message, data, source)
  }

  batch(entries: Array<Omit<LogEntry, 'timestamp' | 'sessionId' | 'userId' | 'requestId'>>) {
    entries.forEach((entry) => {
      this.log(entry.level, entry.type, entry.message, entry.data, entry.source)
    })
  }

  private timers = new Map<string, number>()

  time(label: string) {
    this.timers.set(label, performance.now())
  }

  timeEnd(label: string) {
    const startTime = this.timers.get(label)
    if (startTime !== undefined) {
      const duration = performance.now() - startTime
      this.timers.delete(label)
      this.performance(`Timer "${label}" took ${duration.toFixed(2)}ms`, { duration })
      return duration
    }
    return 0
  }

  getStats() {
    const logs = this.getStoredLogs()
    const stats = {
      total: logs.length,
      byLevel: {} as Record<string, number>,
      byType: {} as Record<string, number>,
      sessionId: this.sessionId,
      bufferSize: this.buffer.length,
    }

    logs.forEach((log) => {
      const levelName = LogLevel[log.level]
      const typeName = log.type

      stats.byLevel[levelName] = (stats.byLevel[levelName] || 0) + 1
      stats.byType[typeName] = (stats.byType[typeName] || 0) + 1
    })

    return stats
  }

  clear() {
    StorageUtil.remove(this.storageKey)
    this.buffer = []
  }

  async flushNow() {
    await this.flush()
  }

  destroy() {
    if (this.flushTimer) {
      clearInterval(this.flushTimer)
    }
    this.flush()
  }

  updateConfig(config: Partial<LoggerConfig>) {
    this.config = { ...this.config, ...config }

    if (this.flushTimer) {
      clearInterval(this.flushTimer)
    }
    this.initFlushTimer()
  }
}

export const logger = new Logger({
  enableConsole: DEBUG_CONFIG.enableNetworkLogging,
  enableRemote: isProduction,
  remoteUrl: import.meta.env.VITE_LOG_ENDPOINT,
})

export const log = logger.info.bind(logger)
export const logError = logger.error.bind(logger)
export const logApi = logger.api.bind(logger)
export const logUser = logger.user.bind(logger)
export const logPerformance = logger.performance.bind(logger)

export function logExecutionTime(target: any, propertyName: string, descriptor: PropertyDescriptor) {
  const method = descriptor.value

  descriptor.value = function (...args: any[]) {
    const label = `${target.constructor.name}.${propertyName}`
    logger.time(label)

    try {
      const result = method.apply(this, args)

      if (result && typeof result.then === 'function') {
        return result.finally(() => {
          logger.timeEnd(label)
        })
      }
      logger.timeEnd(label)
      return result
    } catch (error) {
      logger.timeEnd(label)
      logger.error(`Error in ${label}`, error)
      throw error
    }
  }

  return descriptor
}

export default logger
