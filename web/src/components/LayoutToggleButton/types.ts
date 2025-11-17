export type LayoutToggleSize = 'sm' | 'md' | 'lg'
export type LayoutToggleVariant = 'icon' | 'text' | 'both'

export interface LayoutToggleButtonProps {
  size?: LayoutToggleSize
  variant?: LayoutToggleVariant
  tooltip?: boolean
}
