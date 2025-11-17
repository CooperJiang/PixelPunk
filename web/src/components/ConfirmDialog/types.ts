/** 确认对话框属性接口 */
export interface ConfirmDialogProps {
  modelValue: boolean
  title?: string
  message: string | string[]
  type?: 'warning' | 'danger' | 'info' | 'success'
  width?: string
  confirmText?: string
  cancelText?: string
  loading?: boolean
  closable?: boolean
  requireInput?: boolean
  inputLabel?: string
  inputPlaceholder?: string
  expectedInput?: string
}

/** 确认对话框事件接口 */
export interface ConfirmDialogEmits {
  (e: 'update:modelValue', value: boolean): void
  (e: 'confirm', inputValue?: string): void
  (e: 'cancel'): void
}
