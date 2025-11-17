/**
 * File 组件类型定义
 */
import type { FileType } from '@/constants/file'

export type FileFitMode = 'cover' | 'contain' | 'fill'
export type FileBackgroundPattern = 'grid' | 'dots' | 'diagonal' | 'cyber' | 'circuit' | 'glitch' | 'none'

export interface FileInfo {
  width: number
  height: number
  aspectRatio: number
}

export interface FileProps {
  fileType?: FileType
  src: string
  alt?: string
  width?: string
  height?: string
  fitMode?: FileFitMode
  rounded?: boolean
  borderRadius?: string
  loadingText?: string
  errorText?: string
  showLoadingText?: boolean
  className?: string
  placeholder?: string
  retryCount?: number
  retryInterval?: number
  maintainAspectRatio?: boolean
  showDimensions?: boolean
  backgroundColor?: string
  backgroundPattern?: FileBackgroundPattern
  isNsfw?: boolean
}

export interface FileEmits {
  load: []
  error: []
  loading: [isLoading: boolean]
  dimensions: [width: number, height: number]
  'view-nsfw': []
}
