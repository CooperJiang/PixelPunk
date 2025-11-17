import type { WatermarkBackendConfig, WatermarkConfig } from '../types'

/**
 * 转换前端配置为后端所需的最小配置
 * 移除前端预览相关的冗余参数（文字样式、变换效果、阴影等）
 * 这些参数已经烘焙到 fileBase64 中，后端不需要
 */
export function toBackendConfig(config: WatermarkConfig, base64: string): WatermarkBackendConfig {
  return {
    enabled: config.enabled,
    type: config.type,
    fileBase64: base64,
    position: config.position,
    offsetX: config.offsetX ?? 20,
    offsetY: config.offsetY ?? 20,
    offsetUnit: 'px',
    opacity: config.opacity,
    scale: config.scale ?? 1.0,
  }
}
