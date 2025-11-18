/**
 * 上传Store类型定义
 */

export type UploadStatus =
  | 'pending'
  | 'uploading'
  | 'paused'
  | 'completed'
  | 'failed'
  | 'preparing'
  | 'analyzing'
  | 'retrying'
  | 'instant'

export interface UploadItem {
  id: string
  file: File
  status: UploadStatus
  progress: number
  error?: string
  result?: any
  speed: number
  remainingTime: number
  statusMessage?: string
  type: 'regular' | 'chunked'
  totalChunks?: number
  uploadedChunks?: number
  preview?: string
  dimensions?: {
    width: number
    height: number
  }
  uploadSessionId?: string
}

/**
 * 上传配置接口
 */
export interface GlobalUploadOptions {
  folderId?: string | null
  accessLevel?: 'public' | 'private' | 'protected'
  optimize?: boolean
  autoRemove?: boolean
  storageDuration?: string
  watermarkEnabled?: boolean
  watermarkConfig?: any
}
