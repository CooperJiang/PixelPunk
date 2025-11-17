/**
 * 通知相关类型定义
 * 包含通知消息、提示等类型

 */
import type { ID, Timestamp, VoidFunction } from './base'

/* 通知类型 */
export type NotificationType = 'info' | 'success' | 'warning' | 'error'

/* 通知信息 */
export interface Notification {
  id: ID
  type: NotificationType
  title: string
  message: string
  duration?: number
  closable?: boolean
  actions?: Array<{
    label: string
    action: VoidFunction
  }>
  created_at: Timestamp
}
