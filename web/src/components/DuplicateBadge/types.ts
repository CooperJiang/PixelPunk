/**
 * DuplicateBadge组件类型定义
 */
export interface DuplicateBadgeProps {
  text?: string

  color?: string

  icon?: string

  position?: 'top-left' | 'top-right' | 'bottom-left' | 'bottom-right'
}

/**
 * 重复徽章位置枚举
 */
export const DuplicateBadgePosition = {
  TOP_LEFT: 'top-left',
  TOP_RIGHT: 'top-right',
  BOTTOM_LEFT: 'bottom-left',
  BOTTOM_RIGHT: 'bottom-right',
} as const

export type DuplicateBadgePositionType = keyof typeof DuplicateBadgePosition

/**
 * 预定义颜色常量
 */
export const DuplicateBadgeColors = {
  ORANGE: '#FF9500',
  RED: '#FF3B30',
  YELLOW: '#FFCC02',
  BLUE: '#007AFF',
  GREEN: '#34C759',
} as const
