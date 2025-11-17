export interface TagProps {
  variant?: 'primary' | 'secondary' | 'success' | 'warning' | 'error' | 'info' | 'outline'

  size?: 'small' | 'medium' | 'large'

  closable?: boolean

  disabled?: boolean

  round?: boolean

  truncate?: boolean
}

export interface TagEmits {
  (e: 'close'): void
  (e: 'click'): void
}
