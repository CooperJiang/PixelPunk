export interface SliderProps {
  modelValue: number
  min?: number
  max?: number
  step?: number | string
  disabled?: boolean
  description?: string
  showValue?: boolean
  width?: string
}

export interface SliderEmits {
  (e: 'update:modelValue', value: number): void
}
