/**
 * Message组件类型定义

/* MessageList组件Props */
export interface MessageListProps {
  pageSize?: number
}

/* MessageNotification组件Props */
export interface MessageNotificationProps {
  offsetX?: number

  offsetY?: number
}

/* MessageDetail组件Props */
export interface MessageDetailProps {
  modelValue: boolean

  messageInfo: any | null
}

/* MessageDetail组件Emits */
export interface MessageDetailEmits {
  'update:modelValue': [value: boolean]

  delete: [id?: number]

  markRead: [id?: number]
}

/* MessageToast组件Props */
export interface MessageToastProps {
  id: string

  title?: string

  message: string

  type?: 'success' | 'error' | 'warning' | 'info'

  duration?: number

  closable?: boolean

  actionUrl?: string

  actionText?: string

  actionStyle?: string
}

/* MessageToast组件Emits */
export interface MessageToastEmits {
  close: [id: string]

  action: [id: string]
}

/* TemplateManagement组件Props */
export type TemplateManagementProps = Record<string, never>

/* TemplateModal组件Props */
export interface TemplateModalProps {
  modelValue: boolean

  templateInfo?: any | null
}

/* TemplateModal组件Emits */
export interface TemplateModalEmits {
  'update:modelValue': [value: boolean]

  save: [template: any]
}

/* 消息过滤选项类型 */
export interface MessageFilterOption {
  label: string
  value: string | number
}

/* 表格列配置类型 */
export interface MessageTableColumn {
  key: string
  dataIndex?: string
  title: string
  width?: number
  align?: 'left' | 'center' | 'right'
  ellipsis?: boolean
  render?: (value: any, record: any) => any
}

/* Toast选项类型 */
export interface ToastMessageOptions {
  title?: string
  content: string
  type?: 'success' | 'error' | 'warning' | 'info'
  action_url?: string
  action_text?: string
  action_style?: string
}
