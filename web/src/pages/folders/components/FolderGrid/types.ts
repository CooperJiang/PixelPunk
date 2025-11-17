import type { FolderInfo } from '@/api/types/index'

export interface FolderGridProps {
  folders: FolderInfo[]
  selectMode: boolean
  isAllFoldersSelected: boolean
  isFolderSelected?: (folder: FolderInfo) => boolean
  preview?: boolean
}

export interface FolderGridEmits {
  'folder-click': [folder: FolderInfo]
  'folder-double-click': [folder: FolderInfo]
  'toggle-folder-select': [folder: FolderInfo]
  'toggle-select-all-folders': []
  'invert-folder-selection': []
  'delete-folder': [folder: FolderInfo]
  'edit-folder': [folder: FolderInfo]
  'toggle-folder-visibility': [folder: FolderInfo, event: Event]
  'folder-drag-start': []
  'folder-drag-end': [{ event: Event; sortedFolders: FolderInfo[] }]
  'folder-moved': [folderId: string, targetParentId?: string]
}
