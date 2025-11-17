export interface DownloadOptions {
  fileId: string
  fileName?: string
  quality?: 'original' | 'compressed'
  format?: 'original' | 'webp' | 'jpeg' | 'png'
  position?: {
    x?: number | string
    y?: number | string
    side?: 'left' | 'right' | 'center'
  }
  autoClose?: boolean
  closeDelay?: number
  queueSize?: number
  shareKey?: string
  accessToken?: string
}

export interface DownloadProgressInfo {
  loaded: number
  total: number
  percent: number
}

/** Download status types */
export type DownloadStatus = 'preparing' | 'downloading' | 'success' | 'error'
