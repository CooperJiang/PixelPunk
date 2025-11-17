/**
 * 响应式断点配置
 */

export const RESPONSIVE_BREAKPOINTS = {
  mobile: 768,
  tablet: 1024,
  desktop: 1280,
  large: 1536,
  xlarge: 1920,
} as const

export type ResponsiveBreakpointKey = keyof typeof RESPONSIVE_BREAKPOINTS
