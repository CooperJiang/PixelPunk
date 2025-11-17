export enum UploadStatus {
  PENDING = 'pending',
  PREPARING = 'preparing',
  UPLOADING = 'uploading',
  PAUSED = 'paused',
  COMPLETED = 'completed',
  FAILED = 'failed',
  CANCELLED = 'cancelled',
  RETRYING = 'retrying',
  ANALYZING = 'analyzing',
}

export interface UploadItem {
  id: string
  file: File
  sessionId?: string
  status: UploadStatus
  statusMessage?: string
  progress: number
  totalChunks: number
  uploadedChunks: number
  uploadedChunksList: number[]
  failedChunks: Set<number>
  speed: number
  remainingTime: number
  error?: string
  result?: any
  chunkSize: number
  folderId?: string
  accessLevel: 'public' | 'private' | 'protected'
  optimize: boolean
  watermarkEnabled?: boolean
  watermarkConfig?: any
  fileMD5?: string
  startTime?: number
  endTime?: number
  retryCount: number
  maxRetries: number
  totalBytesUploaded?: number
  lastUploadTime?: number
  averageSpeed?: number
  dimensions?: {
    width: number
    height: number
  }
}

export interface ChunkedUploadOptions {
  chunkSize?: number
  concurrency?: number
  retryCount?: number
  folderId?: string | any
  accessLevel?: 'public' | 'private' | 'protected' | any
  optimize?: boolean | any
  watermarkEnabled?: boolean | any
  watermarkConfig?: any
}
