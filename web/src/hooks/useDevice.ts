import { computed } from 'vue'
import { useBreakpoints, useMediaQuery, useWindowSize } from '@vueuse/core'

/* 响应式断点 */
export const breakpoints = {
  xs: 0,
  sm: 640,
  md: 768,
  lg: 1024,
  xl: 1280,
  '2xl': 1536,
}

export type DeviceType = 'mobile' | 'tablet' | 'desktop'
export type BreakpointKey = keyof typeof breakpoints

/**
 * 响应式获取当前设备类型和屏幕尺寸
 * 使用VueUse库实现，更加简洁和可靠
 */
export function useDevice() {
  const { width, height } = useWindowSize()

  const breakpointObj = useBreakpoints(breakpoints)

  const isPortrait = useMediaQuery('(orientation: portrait)')
  const isTouchScreen = useMediaQuery('(hover: none) and (pointer: coarse)')

  const breakpointState = {
    xs: breakpointObj.between('xs', 'sm'),
    sm: breakpointObj.between('sm', 'md'),
    md: breakpointObj.between('md', 'lg'),
    lg: breakpointObj.between('lg', 'xl'),
    xl: breakpointObj.between('xl', '2xl'),
    '2xl': breakpointObj.greater('2xl'),
  }

  const isMobile = computed(() => width.value < breakpoints.sm)
  const isTablet = computed(() => width.value >= breakpoints.sm && width.value < breakpoints.lg)
  const isDesktop = computed(() => width.value >= breakpoints.lg)

  const deviceType = computed<DeviceType>(() => {
    if (isMobile.value) {
      return 'mobile'
    }
    if (isTablet.value) {
      return 'tablet'
    }
    return 'desktop'
  })

  const isMin = (bp: BreakpointKey) => breakpointObj.greater(bp)

  const isMax = (bp: BreakpointKey) => breakpointObj.smaller(bp)

  const isBreakpoint = (bp: BreakpointKey) => breakpointState[bp]

  return {
    width,
    height,

    deviceType,
    isMobile,
    isTablet,
    isDesktop,

    isPortrait,
    isTouchScreen,

    breakpointState,
    isBreakpoint,
    isMin,
    isMax,
  }
}
