export type BackgroundPattern = 'grid' | 'dots' | 'diagonal' | 'cyber' | 'circuit' | 'glitch' | 'none'

export interface BackgroundProps {
  pattern?: BackgroundPattern
  backgroundColor?: string
  className?: string
}
