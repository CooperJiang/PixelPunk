import type { FileInfo, ImageInfo } from '@/api/types/index'

export interface FileGridProps {
  files: FileInfo[]
  pagination?: {
    total: number
    page: number
    limit: number
  }
  selectMode: boolean
  batchMode: boolean
  viewMode: 'grid' | 'list'
  selectedBatchFiles: string[]
  isAllFilesSelected: boolean
  isFileSelected?: (file: FileInfo) => boolean
  isBatchFileSelected?: (file: FileInfo) => boolean
  preview?: boolean
}

export interface FileGridEmits {
  'preview-file': [file: FileInfo]
  'copy-link': [file: FileInfo]
  'download-file': [file: FileInfo]
  'delete-file': [file: FileInfo]
  'toggle-file-visibility': [file: FileInfo]
  'toggle-file-select': [file: FileInfo]
  'toggle-batch-file-select': [file: FileInfo]
  'toggle-select-all-files': []
  'invert-file-selection': []
  'start-batch-mode': []
  'cancel-batch-mode': []
  'batch-delete': []
  'change-view-mode': [mode: 'grid' | 'list']
  'file-drag-start': []
  'file-drag-end': [{ event: Event; sortedFiles: FileInfo[] }]
  'file-moved': [fileId: string, targetFolderId?: string]
  'file-deleted': [fileId: string]
}

export interface ImageGridProps {
  images: ImageInfo[]
  pagination?: {
    total: number
    page: number
    limit: number
  }
  selectMode: boolean
  batchMode: boolean
  viewMode: 'grid' | 'list'
  selectedBatchImages: string[]
  isAllImagesSelected: boolean
  isImageSelected?: (image: ImageInfo) => boolean
  isBatchImageSelected?: (image: ImageInfo) => boolean
  preview?: boolean
}

export interface ImageGridEmits {
  'preview-image': [image: ImageInfo]
  'view-details': [image: ImageInfo]
  'copy-link': [image: ImageInfo]
  'download-image': [image: ImageInfo]
  'delete-image': [image: ImageInfo]
  'toggle-image-visibility': [image: ImageInfo]
  'toggle-image-select': [image: ImageInfo]
  'toggle-batch-image-select': [image: ImageInfo]
  'toggle-select-all-images': []
  'invert-image-selection': []
  'start-batch-mode': []
  'cancel-batch-mode': []
  'batch-delete': []
  'change-view-mode': [mode: 'grid' | 'list']
  'image-drag-start': []
  'image-drag-end': [{ event: Event; sortedFiles: ImageInfo[] }]
  'image-moved': [fileId: string, targetFolderId?: string]
  'image-deleted': [fileId: string]
}
