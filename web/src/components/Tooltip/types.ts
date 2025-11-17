export type TooltipPlacement = 'top' | 'bottom' | 'left' | 'right' | 'top-start' | 'top-end' | 'bottom-start' | 'bottom-end'

export type TooltipTrigger = 'hover' | 'click' | 'manual'
export type TooltipTheme = 'dark' | 'light'

export interface TooltipProps {
  content?: string
  placement?: TooltipPlacement
  trigger?: TooltipTrigger
  disabled?: boolean
  offset?: number | [number, number]
  showDelay?: number
  hideDelay?: number
  theme?: TooltipTheme
  maxWidth?: string | number
}
