/* 组件导出 */
export { default as MessageList } from './MessageList.vue'
export { default as MessageNotification } from './MessageNotification.vue'
export { default as MessageDetail } from './MessageDetail.vue'
export { default as MessageToast } from './MessageToast.vue'
export { default as TemplateManagement } from './TemplateManagement.vue'
export { default as TemplateModal } from './TemplateModal.vue'

/* 功能模块导出 */
export { useMessageToast, messageToast } from './useMessageToast'
export { useMessageModal, messageModal } from './useMessageModal'
export { useMessageEvents, messageEvents } from './messageEvents'
export { useMessageSystem, messageSystem } from './messageSystem'

/* 类型导出 */
export type { ToastOptions, ToastInstance } from './useMessageToast'
export type { MessageModalInstance } from './useMessageModal'
export type {
  MessageListProps,
  MessageNotificationProps,
  MessageDetailProps,
  MessageDetailEmits,
  MessageToastProps,
  MessageToastEmits,
  TemplateManagementProps,
  TemplateModalProps,
  TemplateModalEmits,
  MessageFilterOption,
  MessageTableColumn,
  ToastMessageOptions,
} from './types'

/* 统一的消息系统入口 */
export { messageSystem as default } from './messageSystem'
