/* 消息系统相关的类型定义 - Legacy */
import type { TranslationFunction } from '@/composables/useTexts'

export interface Message {
  id: number
  user_id: number
  title: string
  content: string
  type: string
  status: number // 1:未读 2:已读 3:已删除
  priority: number // 1:高 2:中 3:低
  read_at?: string
  expires_at?: string
  created_at: string

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

export interface MessageTemplate {
  id: number
  type: string
  title: string
  content: string
  description: string
  is_enabled: boolean
  send_email: boolean
  show_toast: boolean
  toast_type: 'success' | 'error' | 'warning' | 'info'
  default_action_type: string
  default_action_text: string
  default_action_style: 'primary' | 'secondary' | 'success' | 'warning' | 'danger' | 'info'
  action_url_template: string
  created_at: string
  updated_at: string
}

export interface MessageListParams {
  page?: number
  page_size?: number
  status?: number // 1:未读 2:已读
  type?: string
}

export interface MessageListResponse {
  messages: Message[]
  total: number
  page: number
  page_size: number
}

export interface UnreadCountResponse {
  count: number
}

export interface TemplateListResponse {
  templates: MessageTemplate[]
}

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

export interface CreateTemplateRequest {
  type: string
  title: string
  content: string
  description?: string
  is_enabled?: boolean
  send_email?: boolean
  show_toast?: boolean
  toast_type?: 'success' | 'error' | 'warning' | 'info'
  default_action_type?: string
  default_action_text?: string
  default_action_style?: 'primary' | 'secondary' | 'success' | 'warning' | 'danger' | 'info'
  action_url_template?: string
}

export interface UpdateTemplateRequest {
  title?: string
  content?: string
  description?: string
  is_enabled?: boolean
  send_email?: boolean
  show_toast?: boolean
  toast_type?: 'success' | 'error' | 'warning' | 'info'
  default_action_type?: string
  default_action_text?: string
  default_action_style?: 'primary' | 'secondary' | 'success' | 'warning' | 'danger' | 'info'
  action_url_template?: string
}

/* 消息状态常量 */
export const MessageStatus = {
  UNREAD: 1,
  READ: 2,
  DELETED: 3,
} as const

/* 消息优先级常量 */
export const MessagePriority = {
  HIGH: 1,
  NORMAL: 2,
  LOW: 3,
} as const

/* 消息类型常量 */
export const MessageType = {
  SYSTEM_MAINTENANCE: 'system.maintenance',
  SYSTEM_ANNOUNCEMENT: 'system.announcement',
  ACCOUNT_REGISTER: 'account.register',
  ACCOUNT_STORAGE_GRANTED: 'account.storage_granted',
  ACCOUNT_BANDWIDTH_GRANTED: 'account.bandwidth_granted',
  CONTENT_REVIEW_APPROVED: 'content.review_approved',
  CONTENT_REVIEW_REJECTED: 'content.review_rejected',
  CONTENT_REVIEW_PENDING: 'content.review_pending',
  STORAGE_QUOTA_WARNING: 'storage.quota_warning',
  STORAGE_QUOTA_INCREASED: 'storage.quota_increased',
  STORAGE_QUOTA_DECREASED: 'storage.quota_decreased',
  FILE_DELETED_BY_ADMIN: 'file.deleted_by_admin',
  FILE_BATCH_DELETED_BY_ADMIN: 'file.batch_deleted_by_admin',
  FILE_HARD_DELETED_BY_ADMIN: 'file.hard_deleted_by_admin',
  SECURITY_LOGIN_ALERT: 'security.login_alert',
} as const

/* Toast类型常量 */
export const ToastType = {
  SUCCESS: 'success',
  ERROR: 'error',
  WARNING: 'warning',
  INFO: 'info',
} as const

/* 操作样式常量 */
export const ActionStyle = {
  PRIMARY: 'primary',
  SECONDARY: 'secondary',
  SUCCESS: 'success',
  WARNING: 'warning',
  DANGER: 'danger',
  INFO: 'info',
} as const

/* 消息类型显示配置 */
export const getMessageTypeConfig = ($t: TranslationFunction) => ({
  [MessageType.SYSTEM_MAINTENANCE]: {
    label: $t('api.message.types.system.maintenance'),
    icon: '🔧',
    color: 'orange',
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
  [MessageType.SECURITY_LOGIN_ALERT]: {
    label: $t('api.message.types.security.login_alert'),
    icon: '🔒',
    color: 'red',
  },
})

/* 优先级显示配置 */
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

/* 状态显示配置 */
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
