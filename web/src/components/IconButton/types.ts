export interface IconButtonProps {
  type?: 'default' | 'primary' | 'secondary' | 'success' | 'warning' | 'danger' | 'cyber'
  size?: 'tiny' | 'small' | 'medium' | 'large'
  loading?: boolean
  disabled?: boolean
  tooltip?: string
  tooltipPlacement?: 'top' | 'bottom' | 'left' | 'right'
}

export interface IconButtonEvents {
  click: (event: MouseEvent) => void
}
