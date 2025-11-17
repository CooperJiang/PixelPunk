export interface IconOption {
  value: string
  label: string
}

export interface IconPickerProps {
  modelValue: string
  disabled?: boolean
  options?: IconOption[]
}

export interface IconPickerEmits {
  (e: 'update:modelValue', value: string): void
}
