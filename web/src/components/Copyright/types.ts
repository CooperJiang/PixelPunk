/**
 * Copyright组件类型定义
 */
export interface CopyrightProps {
  theme?: 'dark' | 'transparent'
}

/**
 * 版权主题枚举
 */
export const CopyrightTheme = {
  DARK: 'dark',
  TRANSPARENT: 'transparent',
} as const

export type CopyrightThemeType = keyof typeof CopyrightTheme

/**
 * 状态指示器配置
 */
export interface StatusIndicator {
  type: 'online' | 'secure'

  text: string

  icon?: string

  showDot?: boolean
}
