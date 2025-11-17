/**
 * 事件相关类型定义
 * 包含自定义事件、事件监听器等类型

 */
import type { Timestamp } from './base'

/* 自定义事件 */
export interface CustomEvent<T = any> {
  type: string
  data?: T
  timestamp: Timestamp
  source?: string
}

/* 事件监听器 */
export type EventListener<T = any> = (event: CustomEvent<T>) => void
