import type { FileType } from '@/constants/file'

/** Enhanced file preview component props */
export interface EnhancedFilePreviewProps {
  visible: boolean
  fileUrl: string
  fileName?: string
  fileType?: FileType
  fileWidth?: number
  fileHeight?: number
}

export interface EnhancedFilePreviewEmits {
  'update:visible': [value: boolean]
  close: []
}

/** File container component props */
export interface FileContainerProps {
  fileUrl: string
  fileName?: string
  fileType?: FileType
  shouldUseFillMode: boolean
  fileStyle: Record<string, string | number>
}

export interface FileContainerEmits {
  wheel: [event: WheelEvent]
  mousedown: [event: MouseEvent]
  mousemove: [event: MouseEvent]
  mouseup: [event: MouseEvent]
  mouseleave: [event: MouseEvent]
  touchstart: [event: TouchEvent]
  touchmove: [event: TouchEvent]
  touchend: [event: TouchEvent]
  fileClick: [event: MouseEvent]
  fileLoad: [event: Event]
  fileError: [event: Event]
}

/** Control buttons component props */
export interface ControlButtonsProps {
  isBrowserFullscreen: boolean
}

export interface ControlButtonsEmits {
  close: []
  exitFullscreen: []
}

/** Controls hint component props */
export interface ControlsHintProps {
  showControlsHint: boolean
  isBrowserFullscreen: boolean
  shouldUseFillMode: boolean
}

/** Indicators component props */
export interface IndicatorsProps {
  showZoomIndicator: boolean
  showModeIndicator: boolean
  scale: number
  shouldUseFillMode: boolean
}

/* ==================== 向后兼容别名 ==================== */
/* 为了保持向后兼容，提供旧属性名的映射 */
export interface EnhancedImagePreviewProps
  extends Omit<EnhancedFilePreviewProps, 'fileUrl' | 'fileName' | 'fileWidth' | 'fileHeight'> {
  imageUrl: string
  imageName?: string
  imageWidth?: number
  imageHeight?: number
}

export interface ImageContainerProps extends Omit<FileContainerProps, 'fileUrl' | 'fileName' | 'fileStyle'> {
  imageUrl: string
  imageName?: string
  imageStyle: Record<string, string | number>
}
