export type InputType = 'text' | 'password' | 'number' | 'email' | 'tel' | 'url' | 'textarea'

export interface InputProps {
  modelValue: string | number
  type?: InputType
  placeholder?: string
  disabled?: boolean
  readonly?: boolean
  clearable?: boolean
  autofocus?: boolean
  maxlength?: number
  autocomplete?: string
  name?: string
  rows?: number
  error?: boolean
  errorMessage?: string
  prefixIcon?: string
  suffixIcon?: string
  inputId?: string
  width?: string
  height?: string
  no3d?: boolean // 禁用 3D 效果
}

export interface InputEmits {
  (e: 'update:modelValue', value: string | number): void
  (e: 'input', event: Event): void
  (e: 'change', event: Event): void
  (e: 'focus', event: FocusEvent): void
  (e: 'blur', event: FocusEvent): void
  (e: 'keydown', event: KeyboardEvent): void
  (e: 'keyup', event: KeyboardEvent): void
  (e: 'enter', event: KeyboardEvent): void
  (e: 'clear'): void
}
