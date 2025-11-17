export type ButtonType = 'primary' | 'success' | 'warning' | 'danger' | 'default' | 'text'
export type PlacementType = 'top' | 'top-start' | 'top-end' | 'bottom' | 'bottom-start' | 'bottom-end' | 'left' | 'right'
export type TriggerType = 'hover' | 'click'

/**
 * 气泡确认框组件属性
 */
export interface PopconfirmProps {
  title?: string

  description?: string

  placement?: PlacementType

  trigger?: TriggerType

  visible?: boolean

  showArrow?: boolean

  icon?: string

  type?: ButtonType

  cancelText?: string

  confirmText?: string

  confirmLoading?: boolean

  hideAfterConfirm?: boolean

  zIndex?: number

  width?: number | string

  offset?: number
}

/**
 * 气泡确认框组件事件
 */
export interface PopconfirmEmits {
  (e: 'update:visible', value: boolean): void

  (e: 'cancel'): void

  (e: 'confirm'): void

  (e: 'show'): void

  (e: 'hide'): void
}
