import type { FileType } from '@/constants/file'

export interface FileInfo {
  id: string
  url: string
  full_url: string
  display_name?: string
  description?: string
  width?: number
  height?: number
  size_formatted?: string
  resolution?: string
  format?: string
  file_type?: FileType
  tags?: FileTag[]
  ai_info?: {
    description?: string
  }
}

export interface FileTag {
  id: string
  name: string
}

export interface FileViewerProps {
  modelValue: boolean
  file: FileInfo | null
  files?: FileInfo[]
  initialIndex?: number
  showSideNav?: boolean
  showKeyboardTips?: boolean
  controlsHideTimeout?: number
  searchScope?: 'gallery' | 'user' | 'admin' | 'default'
}

export interface FileTransformState {
  scale: number
  rotation: number
  translateX: number
  translateY: number
}

export interface FileViewerState {
  currentIndex: number
  isLoading: boolean
  hasError: boolean
  isDragging: boolean
  isControlsHidden: boolean
  isFirstView: boolean
  isAllContentVisible: boolean
  isImmersiveMode: boolean
  isFullscreen: boolean
  isFillMode: boolean
  isHorizontalFile: boolean
  showModeIndicator: boolean
  showZoomIndicator: boolean
  loadingProgress: number
  loadingText: string
}

export interface ViewerEmits {
  (event: 'update:modelValue', value: boolean): void
  (event: 'close'): void
  (event: 'change', file: FileInfo, index: number): void
  (event: 'error', error: { file: FileInfo; event: Event }): void
  (event: 'load', file: FileInfo): void
  (event: 'prev'): void
  (event: 'next'): void
  (event: 'tag-click', tag: FileTag): void
}
