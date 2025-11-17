import type { VisualTheme } from '@/composables/useTheme'

/**
 * 主题切换组件 Props
 */
export interface ThemeToggleProps {
  size?: 'small' | 'medium' | 'large'
  showDropdown?: boolean
  showText?: boolean
  placement?: 'top' | 'bottom'
  align?: 'left' | 'right'
}

/**
 */
export interface ThemeOption {
  value: VisualTheme
  label: string
  icon: string
  description: string
}
