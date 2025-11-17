/**
 * GlobalUploadDrawer组件类型定义
 */

export type UploadStatus = 'preparing' | 'pending' | 'uploading' | 'paused' | 'completed' | 'failed'

export type UploadType = 'normal' | 'chunked'
export interface UploadItem {
  id: string

  file: File

  status: UploadStatus

  progress: number

  type: UploadType

  preview?: string

  statusMessage?: string

  speed?: number

  error?: string

  result?: {
    full_url?: string
    [key: string]: unknown
  }
}

/**
 * 上传统计信息接口
 */
export interface UploadStatistics {
  uploading: number

  pending: number

  completed: number

  failed: number
}

/**
 * GlobalUploadDrawer组件事件定义
 */
export interface GlobalUploadDrawerEmits {
  close: []
}

/**
 * 批量操作类型
 */
export const BatchActionType = {
  PAUSE_ALL: 'pauseAll',
  RESUME_ALL: 'resumeAll',
  CLEAR_COMPLETED: 'clearCompleted',
  CLEAR_ALL: 'clearAll',
} as const

export type BatchActionTypeKeys = keyof typeof BatchActionType

/**
 * 文件操作类型
 */
export const FileActionType = {
  PAUSE: 'pause',
  RESUME: 'resume',
  RETRY: 'retry',
  COPY_URL: 'copyUrl',
  REMOVE: 'remove',
} as const

export type FileActionTypeKeys = keyof typeof FileActionType
