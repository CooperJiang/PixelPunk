import type { FolderInfo, ImageInfo } from '@/api/types/index'

export interface FolderContentViewProps {
  folders: FolderInfo[]
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
  isAllFoldersSelected: boolean
  isAllImagesSelected: boolean
  isFolderSelected?: (folder: FolderInfo) => boolean
  isImageSelected?: (image: ImageInfo) => boolean
  isBatchImageSelected?: (image: ImageInfo) => boolean
  preview?: boolean // 新增：纯预览模式，隐藏所有编辑功能
}

export interface FolderContentViewEmits {
  'folder-click': [folder: FolderInfo]
  'folder-double-click': [folder: FolderInfo]
  'start-batch-mode': []
  'cancel-batch-mode': []
  'batch-delete': []
  'toggle-select-all-images': []
  'toggle-select-all-folders': []
  'invert-folder-selection': []
  'invert-image-selection': []
  'toggle-batch-image-select': [image: ImageInfo]
  'toggle-image-select': [image: ImageInfo]
  'toggle-folder-select': [folder: FolderInfo]
  'change-view-mode': [mode: 'grid' | 'list']
  'preview-image': [image: ImageInfo]
  'view-details': [image: ImageInfo]
  'copy-link': [image: ImageInfo]
  'download-image': [image: ImageInfo]
  'delete-image': [image: ImageInfo]
  'delete-folder': [folder: FolderInfo]
  'edit-folder': [folder: FolderInfo]
  'toggle-image-visibility': [image: ImageInfo]
  'toggle-folder-visibility': [folder: FolderInfo, event: Event]
  'create-folder': []
  'folder-drag-start': []
  'folder-drag-end': [data: { event: Event; sortedFolders: FolderInfo[] }]
  'image-drag-start': []
  'image-drag-end': [data: { event: Event; sortedFiles: ImageInfo[] }]
  'image-moved': [fileId: string, targetFolderId?: string]
  'image-deleted': [fileId: string]
}
