export interface SmartTagContainerProps {
  tags: string[]
  variant?: 'primary' | 'success' | 'warning' | 'error' | 'info' | 'secondary' | 'outline'
  size?: 'small' | 'medium' | 'large'
  truncate?: boolean
  minTags?: number
  moreVariant?: 'primary' | 'success' | 'warning' | 'error' | 'info' | 'secondary' | 'outline'
}

export interface SmartTagContainerEmits {
  (e: 'tag-click', tag: string, index: number): void
}
