export interface BadgeProps {
  text?: string
  color?: string
  icon?: string
  position?: 'top-left' | 'top-right' | 'bottom-left' | 'bottom-right'
  variant?: 'success' | 'error' | 'warning' | 'info' | 'danger' | 'primary'
  inline?: boolean
}

export const BadgePosition = {
  TOP_LEFT: 'top-left',
  TOP_RIGHT: 'top-right',
  BOTTOM_LEFT: 'bottom-left',
  BOTTOM_RIGHT: 'bottom-right',
} as const

export type BadgePositionType = keyof typeof BadgePosition

export const BadgeVariant = {
  SUCCESS: 'success',
  ERROR: 'error',
  WARNING: 'warning',
  INFO: 'info',
  DANGER: 'danger',
  PRIMARY: 'primary',
} as const

export type BadgeVariantType = keyof typeof BadgeVariant
