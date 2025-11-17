/**
 * Switch开关组件属性
 */
export interface SwitchProps {
  modelValue: boolean

  label?: string

  disabled?: boolean

  size?: 'small' | 'medium' | 'large'
}

/**
 * Switch开关组件事件
 */
export interface SwitchEmits {
  (e: 'update:modelValue', value: boolean): void

  (e: 'change', value: boolean): void
}
