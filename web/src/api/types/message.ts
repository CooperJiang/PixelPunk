/**
 * 消息系统相关类型定义
 */
import type { PaginationParams, TimeStamps } from './common'
import type { TranslationFunction } from '@/composables/useTexts'

/* ==================== 消息相关类型 ==================== */
export interface Message extends TimeStamps {
  id: number
  user_id: number
  title: string
  content: string
  type: string
  status: number // 1:未读 2:已读 3:已删除
  priority: number // 1:高 2:中 3:低
  read_at?: string
  expires_at?: string

  related_type?: string
  related_id?: string
  related_data?: Record<string, unknown>

  is_actionable: boolean
  action_type?: string
  action_url?: string
  action_text?: string
  action_style?: string
  requires_permission?: string
  meta_data?: Record<string, unknown>
}

/* ==================== 消息模板类型 ==================== */
export type ToastType = 'success' | 'error' | 'warning' | 'info'
export type ActionStyle = 'primary' | 'secondary' | 'success' | 'warning' | 'danger' | 'info'

export interface MessageTemplate extends TimeStamps {
  id: number
  type: string
  title: string
  content: string
  description: string
  is_enabled: boolean
  send_email: boolean
  show_toast: boolean
  toast_type: ToastType
  default_action_type: string
  default_action_text: string
  default_action_style: ActionStyle
  action_url_template: string
}

/* ==================== 消息查询参数类型 ==================== */
export interface MessageListParams extends PaginationParams {
  page_size?: number // 兼容现有API
  status?: number // 1:未读 2:已读
  type?: string
}

/* ==================== 消息列表响应类型 ==================== */
export interface MessageListResponse {
  messages: Message[]
  total: number
  page: number
  page_size: number
}

/* ==================== 消息操作响应类型 ==================== */
export interface UnreadCountResponse {
  count: number
}

export interface TemplateListResponse {
  templates: MessageTemplate[]
}

/* ==================== 消息发送请求类型 ==================== */
export interface SendMessageRequest {
  user_id: number
  type: string
  title: string
  content: string
  priority?: number
  expires_at?: string
  related_type?: string
  related_id?: string
  related_data?: Record<string, unknown>
  is_actionable?: boolean
  action_type?: string
  action_url?: string
  action_text?: string
  action_style?: string
  requires_permission?: string
  meta_data?: Record<string, unknown>
}

/* ==================== 消息模板操作类型 ==================== */
export interface CreateTemplateRequest {
  type: string
  title: string
  content: string
  description?: string
  is_enabled?: boolean
  send_email?: boolean
  show_toast?: boolean
  toast_type?: ToastType
  default_action_type?: string
  default_action_text?: string
  default_action_style?: ActionStyle
  action_url_template?: string
}

export interface UpdateTemplateRequest {
  title?: string
  content?: string
  description?: string
  is_enabled?: boolean
  send_email?: boolean
  show_toast?: boolean
  toast_type?: ToastType
  default_action_type?: string
  default_action_text?: string
  default_action_style?: ActionStyle
  action_url_template?: string
}

/* ==================== 消息常量定义 ==================== */
export const MessageStatus = {
  UNREAD: 1,
  READ: 2,
  DELETED: 3,
} as const

export const MessagePriority = {
  HIGH: 1,
  NORMAL: 2,
  LOW: 3,
} as const

export const MessageType = {
  SYSTEM_MAINTENANCE: 'system.maintenance',
  SYSTEM_UPDATE: 'system.update',
  SYSTEM_ANNOUNCEMENT: 'system.announcement',
  ACCOUNT_REGISTER: 'account.register',
  ACCOUNT_STORAGE_GRANTED: 'account.storage_granted',
  ACCOUNT_BANDWIDTH_GRANTED: 'account.bandwidth_granted',
  CONTENT_REVIEW_APPROVED: 'content.review_approved',
  CONTENT_REVIEW_REJECTED: 'content.review_rejected',
  CONTENT_REVIEW_PENDING: 'content.review_pending',
  STORAGE_QUOTA_WARNING: 'storage.quota_warning',
  STORAGE_QUOTA_EXCEEDED: 'storage.quota_exceeded',
  STORAGE_QUOTA_INCREASED: 'storage.quota_increased',
  STORAGE_QUOTA_DECREASED: 'storage.quota_decreased',
  FILE_DELETED_BY_ADMIN: 'file.deleted_by_admin',
  FILE_BATCH_DELETED_BY_ADMIN: 'file.batch_deleted_by_admin',
  FILE_HARD_DELETED_BY_ADMIN: 'file.hard_deleted_by_admin',
  FILE_EXPIRY_WARNING: 'file.expiry_warning',
  FILE_THUMBNAIL_FAILED: 'file.thumbnail_failed',
  SECURITY_LOGIN_ALERT: 'security.login_alert',
  SECURITY_PASSWORD_CHANGED: 'security.password_changed',
  APIKEY_CREATED: 'apikey.created',
  APIKEY_DELETED: 'apikey.deleted',
  APIKEY_REGENERATED: 'apikey.regenerated',
  APIKEY_DISABLED: 'apikey.disabled',
  APIKEY_ENABLED: 'apikey.enabled',
  RANDOM_API_CREATED: 'random_api.created',
  RANDOM_API_DELETED: 'random_api.deleted',
  RANDOM_API_DISABLED: 'random_api.disabled',
  RANDOM_API_ENABLED: 'random_api.enabled',
  SHARE_EXPIRY_WARNING: 'share.expiry_warning',
} as const

export const ToastTypes = {
  SUCCESS: 'success',
  ERROR: 'error',
  WARNING: 'warning',
  INFO: 'info',
} as const

export const ActionStyles = {
  PRIMARY: 'primary',
  SECONDARY: 'secondary',
  SUCCESS: 'success',
  WARNING: 'warning',
  DANGER: 'danger',
  INFO: 'info',
} as const

/* ==================== 消息类型显示配置 ==================== */
export const getMessageTypeConfig = ($t: TranslationFunction) => ({
  [MessageType.SYSTEM_MAINTENANCE]: {
    label: $t('api.message.types.system.maintenance'),
    icon: '🔧',
    color: 'orange',
  },
  [MessageType.SYSTEM_UPDATE]: {
    label: $t('api.message.types.system.update'),
    icon: '🔄',
    color: 'blue',
  },
  [MessageType.SYSTEM_ANNOUNCEMENT]: {
    label: $t('api.message.types.system.announcement'),
    icon: '📢',
    color: 'blue',
  },
  [MessageType.ACCOUNT_REGISTER]: {
    label: $t('api.message.types.account.register'),
    icon: '🎉',
    color: 'green',
  },
  [MessageType.ACCOUNT_STORAGE_GRANTED]: {
    label: $t('api.message.types.account.storage_granted'),
    icon: '💾',
    color: 'purple',
  },
  [MessageType.ACCOUNT_BANDWIDTH_GRANTED]: {
    label: $t('api.message.types.account.bandwidth_granted'),
    icon: '🚀',
    color: 'cyan',
  },
  [MessageType.CONTENT_REVIEW_APPROVED]: {
    label: $t('api.message.types.content.review_approved'),
    icon: '✅',
    color: 'green',
  },
  [MessageType.CONTENT_REVIEW_REJECTED]: {
    label: $t('api.message.types.content.review_rejected'),
    icon: '❌',
    color: 'red',
  },
  [MessageType.CONTENT_REVIEW_PENDING]: {
    label: $t('api.message.types.content.review_pending'),
    icon: '⏳',
    color: 'yellow',
  },
  [MessageType.STORAGE_QUOTA_WARNING]: {
    label: $t('api.message.types.storage.quota_warning'),
    icon: '⚠️',
    color: 'yellow',
  },
  [MessageType.STORAGE_QUOTA_EXCEEDED]: {
    label: $t('api.message.types.storage.quota_exceeded'),
    icon: '🚫',
    color: 'red',
  },
  [MessageType.STORAGE_QUOTA_INCREASED]: {
    label: $t('api.message.types.storage.quota_increased'),
    icon: '📈',
    color: 'green',
  },
  [MessageType.STORAGE_QUOTA_DECREASED]: {
    label: $t('api.message.types.storage.quota_decreased'),
    icon: '📉',
    color: 'orange',
  },
  [MessageType.FILE_DELETED_BY_ADMIN]: {
    label: $t('api.message.types.file.deleted_by_admin'),
    icon: '🗑️',
    color: 'red',
  },
  [MessageType.FILE_BATCH_DELETED_BY_ADMIN]: {
    label: $t('api.message.types.file.batch_deleted_by_admin'),
    icon: '🗂️',
    color: 'red',
  },
  [MessageType.FILE_HARD_DELETED_BY_ADMIN]: {
    label: $t('api.message.types.file.hard_deleted_by_admin'),
    icon: '⚠️',
    color: 'red',
  },
  [MessageType.FILE_EXPIRY_WARNING]: {
    label: $t('api.message.types.file.expiry_warning'),
    icon: '⏰',
    color: 'yellow',
  },
  [MessageType.FILE_THUMBNAIL_FAILED]: {
    label: $t('api.message.types.file.thumbnail_failed'),
    icon: '🖼️',
    color: 'orange',
  },
  [MessageType.SECURITY_LOGIN_ALERT]: {
    label: $t('api.message.types.security.login_alert'),
    icon: '🔒',
    color: 'red',
  },
  [MessageType.SECURITY_PASSWORD_CHANGED]: {
    label: $t('api.message.types.security.password_changed'),
    icon: '🔑',
    color: 'orange',
  },
  [MessageType.APIKEY_CREATED]: {
    label: $t('api.message.types.apikey.created'),
    icon: '🔑',
    color: 'green',
  },
  [MessageType.APIKEY_DELETED]: {
    label: $t('api.message.types.apikey.deleted'),
    icon: '🗑️',
    color: 'red',
  },
  [MessageType.APIKEY_REGENERATED]: {
    label: $t('api.message.types.apikey.regenerated'),
    icon: '🔄',
    color: 'blue',
  },
  [MessageType.APIKEY_DISABLED]: {
    label: $t('api.message.types.apikey.disabled'),
    icon: '⛔',
    color: 'orange',
  },
  [MessageType.APIKEY_ENABLED]: {
    label: $t('api.message.types.apikey.enabled'),
    icon: '✅',
    color: 'green',
  },
  [MessageType.RANDOM_API_CREATED]: {
    label: $t('api.message.types.random_api.created'),
    icon: '🎲',
    color: 'green',
  },
  [MessageType.RANDOM_API_DELETED]: {
    label: $t('api.message.types.random_api.deleted'),
    icon: '🗑️',
    color: 'red',
  },
  [MessageType.RANDOM_API_DISABLED]: {
    label: $t('api.message.types.random_api.disabled'),
    icon: '⛔',
    color: 'orange',
  },
  [MessageType.RANDOM_API_ENABLED]: {
    label: $t('api.message.types.random_api.enabled'),
    icon: '✅',
    color: 'green',
  },
  [MessageType.SHARE_EXPIRY_WARNING]: {
    label: $t('api.message.types.share.expiry_warning'),
    icon: '📤',
    color: 'yellow',
  },
})

/* ==================== 优先级显示配置 ==================== */
export const getPriorityConfig = ($t: TranslationFunction) => ({
  [MessagePriority.HIGH]: {
    label: $t('api.message.priority.high'),
    color: 'red',
    class: 'text-red-500 bg-red-50',
  },
  [MessagePriority.NORMAL]: {
    label: $t('api.message.priority.normal'),
    color: 'orange',
    class: 'text-orange-500 bg-orange-50',
  },
  [MessagePriority.LOW]: {
    label: $t('api.message.priority.low'),
    color: 'gray',
    class: 'text-gray-500 bg-gray-50',
  },
})

/* ==================== 状态显示配置 ==================== */
export const getStatusConfig = ($t: TranslationFunction) => ({
  [MessageStatus.UNREAD]: {
    label: $t('api.message.status.unread'),
    color: 'blue',
    class: 'text-blue-500 bg-blue-50',
  },
  [MessageStatus.READ]: {
    label: $t('api.message.status.read'),
    color: 'gray',
    class: 'text-gray-500 bg-gray-50',
  },
  [MessageStatus.DELETED]: {
    label: $t('api.message.status.deleted'),
    color: 'gray',
    class: 'text-gray-400 bg-gray-50',
  },
})
