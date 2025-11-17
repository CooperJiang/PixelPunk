export type ButtonType = 'primary' | 'secondary' | 'outlined' | 'text' | 'danger' | 'success' | 'warning' | 'info'

export type ButtonSize = 'small' | 'medium' | 'large'

export interface ButtonProps {
  type?: ButtonType
  size?: ButtonSize
  icon?: string
  rightIcon?: string
  loading?: boolean
  loadingMode?: 'replace' | 'inline'
  disabled?: boolean
  block?: boolean
  customClass?: string
}

export interface ButtonEmits {
  (e: 'click', event: MouseEvent): void
  (e: 'mouseover', event: MouseEvent): void
  (e: 'mouseleave', event: MouseEvent): void
}
