export interface CheckboxProps {
  modelValue?: boolean
  disabled?: boolean
  label?: string | number | boolean
  name?: string
  indeterminate?: boolean
  size?: 'small' | 'medium' | 'large'
}

export interface CheckboxEmits {
  (e: 'update:modelValue', value: boolean): void
  (e: 'change', value: boolean): void
}
