/**
 * 对话框组件属性
 */
export interface DialogProps {
  modelValue: boolean

  title?: string

  width?: string

  height?: string

  maxWidth?: string

  maxHeight?: string

  showDefaultFooter?: boolean

  showCloseButton?: boolean

  noPadding?: boolean

  noScroll?: boolean

  showFooter?: boolean

  cancelText?: string

  confirmText?: string

  loading?: boolean

  closeOnEsc?: boolean

  closeOnClickOverlay?: boolean

  hideBorder?: boolean

  fullscreen?: boolean
}

/**
 * 对话框组件事件
 */
export interface DialogEmits {
  (e: 'update:modelValue', value: boolean): void

  (e: 'confirm'): void

  (e: 'cancel'): void

  (e: 'close'): void
}
