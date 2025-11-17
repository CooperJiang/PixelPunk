import type { HexagonStyle, LoaderHexStyle } from '../types'
import {
  HEX_SIZE_MIN,
  HEX_SIZE_RANGE,
  HEX_POSITION_RANGE,
  HEX_ANIMATION_DELAY_MAX,
  LOADER_HEX_ANGLE_STEP,
  LOADER_HEX_RADIUS,
  LOADER_HEX_ANIMATION_DELAY_STEP,
} from '../constants'

/* 生成背景六边形样式 */
export function getHexagonStyle(_i: number): HexagonStyle {
  const size = Math.random() * HEX_SIZE_RANGE + HEX_SIZE_MIN
  const left = Math.random() * HEX_POSITION_RANGE
  const top = Math.random() * HEX_POSITION_RANGE
  const animationDelay = Math.random() * HEX_ANIMATION_DELAY_MAX

  return {
    width: `${size}px`,
    height: `${size}px`,
    left: `${left}%`,
    top: `${top}%`,
    animationDelay: `${animationDelay}s`,
  }
}

export function getLoaderHexStyle(i: number): LoaderHexStyle {
  const angle = (i - 1) * LOADER_HEX_ANGLE_STEP
  const x = Math.cos((angle * Math.PI) / 180) * LOADER_HEX_RADIUS
  const y = Math.sin((angle * Math.PI) / 180) * LOADER_HEX_RADIUS

  return {
    transform: `translate(${x}px, ${y}px)`,
    animationDelay: `${(i - 1) * LOADER_HEX_ANIMATION_DELAY_STEP}s`,
  }
}
