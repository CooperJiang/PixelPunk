/** Upload status type */
export type FloatUploadStatus = 'preparing' | 'uploading' | 'success' | 'error' | 'paused'

/** Position coordinates interface */
export interface Position {
  x: number
  y: number
}

/** Upload statistics interface */
export interface FloatUploadStatistics {
  uploading: number
  pending: number
  completed: number
  failed: number
  total: number
}

/** GlobalUploadFloat component props */
export interface GlobalUploadFloatProps {
  initialPosition?: Position
  hideOnRoutes?: string[]
  draggable?: boolean
  showFileInfo?: boolean
}

/** GlobalUploadFloat component events */
export interface GlobalUploadFloatEmits {
  (e: 'open-drawer'): void
}

/** Exposed methods interface */
export interface GlobalUploadFloatExposed {
  hide: () => void
  show: () => void
}

/** Float display status constants */
export const FloatDisplayStatus = {
  PREPARING: 'preparing',
  UPLOADING: 'uploading',
  SUCCESS: 'success',
  ERROR: 'error',
  PAUSED: 'paused',
} as const

export type FloatDisplayStatusType = keyof typeof FloatDisplayStatus

/** Storage keys constants */
export const StorageKeys = {
  POSITION: 'upload-float-position',
} as const
