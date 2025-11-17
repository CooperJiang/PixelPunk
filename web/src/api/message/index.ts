import { getCurrentLocale } from '@/utils/locale'
import { del, get, post, put } from '@/utils/network/http'
import type { TranslationFunction } from '@/composables/useTexts'

/* ============ 类型定义 ============ */
export interface Message {
  id: number
  user_id: number
  type: string
  data: Record<string, any> // 结构化消息数据，前端根据 type 动态生成标题和内容
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
  toast_type: string
  default_action_type: string
  default_action_text: string
  default_action_style: string
  action_url_template: string
  created_at: string
  updated_at: string
}

export interface MessageListResponse {
  items: Message[]
  pagination: {
    total: number
    size: number
    current_page: number
    last_page: number
  }
}

export interface UnreadCountResponse {
  count: number
}

export interface TemplateListResponse {
  templates: MessageTemplate[]
}

/* ============ 用户消息接口 ============ */

/**
 * 获取用户消息列表
 */
export function getUserMessages(params?: {
  page?: number
  pageSize?: number
  status?: number // 1:未读 2:已读
  type?: string
}): Promise<MessageListResponse> {
  return get<MessageListResponse>('/messages', params)
}

export function getUnreadCount(): Promise<UnreadCountResponse> {
  return get<UnreadCountResponse>('/messages/unread-count')
}

export function markMessageRead(messageId: string): Promise<void> {
  return post<void>(`/messages/${messageId}/read`)
}

export function markAllMessagesRead(): Promise<void> {
  return post<void>('/messages/read-all')
}

export function deleteMessage(messageId: string): Promise<void> {
  return del<void>(`/messages/${messageId}`)
}

export function sendMessage(data: {
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
}): Promise<void> {
  return post<void>('/admin/messages/send', data)
}

export function getMessageStats(): Promise<{
  total_messages: string | number
  unread_messages: string | number
  daily_messages: string | number
  template_count: string | number
}> {
  return get('/admin/messages/stats')
}

export function getMessageTemplates(): Promise<TemplateListResponse> {
  return get<TemplateListResponse>('/admin/messages/templates')
}

export function getMessageTemplate(templateId: number): Promise<MessageTemplate> {
  return get<MessageTemplate>(`/admin/messages/templates/${templateId}`)
}

export function createMessageTemplate(data: {
  type: string
  title: string
  content: string
  description?: string
  is_enabled?: boolean
  send_email?: boolean
  show_toast?: boolean
  toast_type?: string
  default_action_type?: string
  default_action_text?: string
  default_action_style?: string
  action_url_template?: string
}): Promise<MessageTemplate> {
  return post<MessageTemplate>('/admin/messages/templates', data)
}

export function updateMessageTemplate(
  templateId: number,
  data: {
    title?: string
    content?: string
    description?: string
    is_enabled?: boolean
    send_email?: boolean
    show_toast?: boolean
    toast_type?: string
    default_action_type?: string
    default_action_text?: string
    default_action_style?: string
    action_url_template?: string
  }
): Promise<void> {
  return put<void>(`/admin/messages/templates/${templateId}`, data)
}

export function deleteMessageTemplate(templateId: number): Promise<void> {
  return del<void>(`/admin/messages/templates/${templateId}`)
}

export function toggleMessageTemplate(templateId: number, enabled: boolean): Promise<void> {
  return post<void>(`/admin/messages/templates/${templateId}/toggle`, { enabled })
}

/* 获取消息类型标签（带翻译） */
export function getMessageTypeLabel(type: string, $t: TranslationFunction): string {
  const typeLabels: Record<string, string> = {
    'system.maintenance': $t('api.message.types.system.maintenance'),
    'system.announcement': $t('api.message.types.system.announcement'),
    'account.register': $t('api.message.types.account.register'),
    'account.storage_granted': $t('api.message.types.account.storage_granted'),
    'account.bandwidth_granted': $t('api.message.types.account.bandwidth_granted'),
    'content.review_approved': $t('api.message.types.content.review_approved'),
    'content.review_rejected': $t('api.message.types.content.review_rejected'),
    'content.review_pending': $t('api.message.types.content.review_pending'),
    'storage.quota_warning': $t('api.message.types.storage.quota_warning'),
    'security.login_alert': $t('api.message.types.security.login_alert'),
  }
  return typeLabels[type] || type
}

/* 获取消息优先级信息（带翻译） */
export function getMessagePriorityInfo(
  priority: number,
  $t: TranslationFunction
): {
  label: string
  color: string
  class: string
} {
  switch (priority) {
    case 1:
      return { label: $t('api.message.priority.high'), color: '#f56565', class: 'text-red-500' }
    case 2:
      return { label: $t('api.message.priority.normal'), color: '#ed8936', class: 'text-orange-500' }
    case 3:
      return { label: $t('api.message.priority.low'), color: '#38b2ac', class: 'text-teal-500' }
    default:
      return { label: $t('api.message.priority.normal'), color: '#ed8936', class: 'text-orange-500' }
  }
}

/* 获取消息状态信息（带翻译） */
export function getMessageStatusInfo(
  status: number,
  $t: TranslationFunction
): {
  label: string
  color: string
  class: string
} {
  switch (status) {
    case 1:
      return { label: $t('api.message.status.unread'), color: '#3182ce', class: 'text-blue-500' }
    case 2:
      return { label: $t('api.message.status.read'), color: '#38a169', class: 'text-green-500' }
    case 3:
      return { label: $t('api.message.status.deleted'), color: '#a0aec0', class: 'text-gray-400' }
    default:
      return { label: $t('api.message.status.unknown'), color: '#a0aec0', class: 'text-gray-400' }
  }
}

/* 格式化相对时间（带翻译） */
export function formatRelativeTime(dateString: string, $t: TranslationFunction): string {
  const date = new Date(dateString)
  const now = new Date()
  const diffMs = now.getTime() - date.getTime()
  const diffMinutes = Math.floor(diffMs / (1000 * 60))
  const diffHours = Math.floor(diffMs / (1000 * 60 * 60))
  const diffDays = Math.floor(diffMs / (1000 * 60 * 60 * 24))

  if (diffMinutes < 1) {
    return $t('api.message.time.justNow')
  } else if (diffMinutes < 60) {
    return $t('api.message.time.minutesAgo', { n: diffMinutes })
  } else if (diffHours < 24) {
    return $t('api.message.time.hoursAgo', { n: diffHours })
  } else if (diffDays < 7) {
    return $t('api.message.time.daysAgo', { n: diffDays })
  }
  return date.toLocaleDateString(getCurrentLocale(), {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
  })
}
