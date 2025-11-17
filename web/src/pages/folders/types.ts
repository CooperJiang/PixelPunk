import type { FolderInfo, ImageInfo } from '@/api/types'
import type { Ref } from 'vue'

export interface ShareButtonRef {
  shareButtonRef?: ShareButtonInstance
  selectedFolders?: FolderInfo[]
  selectedImages?: ImageInfo[]
  isFolderSelected?: (folder: FolderInfo) => boolean
  isImageSelected?: (image: ImageInfo) => boolean
  toggleSelectAllImages?: (images: ImageInfo[]) => void
  toggleSelectAllFolders?: (folders: FolderInfo[]) => void
  invertFolderSelection?: (folders: FolderInfo[]) => void
  invertImageSelection?: (images: ImageInfo[]) => void
  toggleImageSelect?: (image: ImageInfo) => void
  toggleFolderSelect?: (folder: FolderInfo) => void
  cancelSelectMode?: () => void
}

export interface ShareButtonInstance {
  selectedFolders: FolderInfo[]
  selectedImages: ImageInfo[]
  isFolderSelected: (folder: FolderInfo) => boolean
  isImageSelected: (image: ImageInfo) => boolean
  toggleSelectAllImages: (images: ImageInfo[]) => void
  toggleSelectAllFolders: (folders: FolderInfo[]) => void
  invertFolderSelection: (folders: FolderInfo[]) => void
  invertImageSelection: (images: ImageInfo[]) => void
  toggleImageSelect: (image: ImageInfo) => void
  toggleFolderSelect: (folder: FolderInfo) => void
  cancelSelectMode: () => void
}

export type ViewMode = 'grid' | 'list'

export type DeleteType = 'folder' | 'image'

export type DialogMode = 'create' | 'edit'

export interface BreadcrumbItem {
  id: string | null
  name: string
  path?: string
}

export interface DragData {
  event: DragEvent
  sortedFolders?: FolderInfo[]
  sortedFiles?: ImageInfo[]
}

export interface FolderMoveData {
  fileId: string
  targetFolderId: string
}

export interface KeyboardEventsOptions {
  selectMode: Ref<boolean>
  shareButtonRef: Ref<ShareButtonRef | undefined>
}
