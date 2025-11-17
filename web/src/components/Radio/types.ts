export interface RadioProps {
  modelValue: string | number | boolean
  value: string | number | boolean
  disabled?: boolean
}

export interface RadioEmits {
  'update:modelValue': [value: string | number | boolean]
}

export interface RadioGroupOption {
  label: string
  value: string | number | boolean
}

export interface RadioGroupProps {
  modelValue: string | number | boolean
  options: RadioGroupOption[]
  disabled?: boolean
  layout?: 'horizontal' | 'vertical'
}

export interface RadioGroupEmits {
  'update:modelValue': [value: string | number | boolean]
}
