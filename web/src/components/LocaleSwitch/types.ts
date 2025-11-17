import type { SupportedLocale } from '@/composables/useLocale'

/**
 * 语言切换组件 Props
 */
export interface LocaleSwitchProps {
  size?: 'small' | 'medium' | 'large'
  showDropdown?: boolean
  showText?: boolean
  placement?: 'top' | 'bottom'
  align?: 'left' | 'right'
}

/**
 * 语言选项
 */
export interface LocaleOption {
  value: SupportedLocale
  label: string
  emoji: string
  description: string
  nativeName: string
}
