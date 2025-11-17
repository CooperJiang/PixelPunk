export type WatermarkType = 'text' | 'image'

export type WatermarkPosition =
  | 'top-left'
  | 'top-center'
  | 'top-right'
  | 'middle-left'
  | 'middle-center'
  | 'middle-right'
  | 'bottom-left'
  | 'bottom-center'
  | 'bottom-right'
  | 'custom'

/* 后端配置接口（前后端统一使用 'image' 类型） */
export interface WatermarkBackendConfig {
  enabled: boolean
  type: WatermarkType // 统一使用：'text' | 'image'
  fileBase64: string
  position: WatermarkPosition
  offsetX: number
  offsetY: number
  offsetUnit: 'px'
  opacity: number
  scale: number // 水印在目标图片上的缩放比例（前端生成的水印图片需要在目标图片上二次缩放）
}

export interface WatermarkConfig {
  enabled: boolean
  type: WatermarkType

  text: string
  fontSize: number
  fontFamily: string
  fontColor: string
  fontWeight: 'normal' | 'bold'

  fileBase64?: string
  imageFile?: File

  position: WatermarkPosition
  offsetX: number // 距离参考边缘的距离（px）
  offsetY: number // 距离参考边缘的距离（px）
  offsetUnit: 'px' | 'percent' // 单位

  opacity: number // 0-1
  scale: number // 0.1-2
  rotation: number // -180 to 180 degrees

  shadow: boolean
  shadowColor: string
  shadowBlur: number
  shadowOffsetX: number
  shadowOffsetY: number

  relativeSize?: number
}

export const DEFAULT_WATERMARK_CONFIG: WatermarkConfig = {
  enabled: false,
  type: 'text',

  text: 'PixelPunk',
  fontSize: 48,
  fontFamily: 'Arial',
  fontColor: '#ffffff',
  fontWeight: 'normal',

  fileBase64: '',

  position: 'bottom-right',
  offsetX: 20, // 距离右边缘20px
  offsetY: 20, // 距离底边缘20px
  offsetUnit: 'px',

  opacity: 0.7,
  scale: 1.0,
  rotation: 0,

  shadow: true,
  shadowColor: '#000000',
  shadowBlur: 4,
  shadowOffsetX: 1,
  shadowOffsetY: 1,

  relativeSize: 0.15,
}

export interface WatermarkPreviewProps {
  config: WatermarkConfig
  previewImage?: string
  onConfigChange: (config: WatermarkConfig) => void
}

export interface WatermarkPanelProps {
  config: WatermarkConfig
  onConfigChange: (config: WatermarkConfig) => void
}
