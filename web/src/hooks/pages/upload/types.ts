import type { ComputedRef, Ref } from 'vue'

/**
 * 上传页面通用类型定义
 * 供所有上传相关的布局和组件使用

 */
export interface UploadFile {
  file: File
  name: string
  size: number
  progress: number
  status: 'pending' | 'uploading' | 'success' | 'error' | 'analyzing' | 'retrying' | 'paused'
  full_url?: string
  full_thumb_url?: string
  url?: string
  id?: string
  preview?: string
  error?: string
  uniqueId: string
  previewLoaded: boolean
  itemLoaded: boolean
  statusMessage?: string
  speed?: number
  remainingTime?: number
}

export interface UploadSettings {
  folderId: string | null
  accessLevel: 'public' | 'private' | 'protected'
  optimize: boolean
  autoRemove: boolean
}

export interface PasteShortcut {
  key: string
  isMac: boolean
}

export interface UploadStats {
  pendingCount: number
  uploadingCount: number
  successCount: number
  errorCount: number
  totalFileSize: number
}

/**
 * 通用上传Hook返回的接口
 */
export interface UseUploadReturn {
  fileInput: Ref<HTMLInputElement | null>
  isDragging: Ref<boolean>
  folderId: Ref<string | null>
  accessLevel: Ref<'public' | 'private' | 'protected'>
  optimize: Ref<boolean>
  autoRemove: Ref<boolean>

  smartUpload: any
  uploadQueue: ComputedRef<UploadFile[]>
  globalProgress: ComputedRef<number>
  globalSpeed: ComputedRef<number>

  hasPendingFiles: ComputedRef<boolean>
  hasUploadingFiles: ComputedRef<boolean>
  hasSuccessFiles: ComputedRef<boolean>
  pendingCount: ComputedRef<number>
  uploadingCount: ComputedRef<number>
  successCount: ComputedRef<number>
  errorCount: ComputedRef<number>
  totalFileSize: ComputedRef<number>
  acceptFormats: ComputedRef<string>
  pasteShortcut: ComputedRef<PasteShortcut>

  triggerFileInput: () => void
  onDragOver: () => void
  onDragLeave: () => void
  onDrop: (event: DragEvent) => Promise<void>
  handleFileChange: (event: Event) => Promise<void>

  startUpload: () => Promise<void>
  cancelUpload: () => void
  retryUpload: (index: number) => void
  resumeUploadWrapper: (index: number) => void
  removeFile: (index: number) => void
  showNoFilesMessage: () => void
  copyAllUrls: () => Promise<void>
  onSettingsChange: (settings: any) => void
  formatFileSize: (size: number) => string

  setupEventListeners: () => Promise<void>
  cleanupEventListeners: () => void

  settingsStore: Event

  clearQueue: () => void
  clearAllSessions: () => void
}
