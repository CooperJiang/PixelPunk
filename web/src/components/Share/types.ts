export interface ShareItem {
  id: string | number
  name?: string
  [key: string]: unknown
}

export interface ShareButtonProps {
  folders?: ShareItem[]
  images?: ShareItem[]
  requireLogin?: boolean
}

export interface ShareButtonEmits {
  'update:selectMode': [value: boolean]
  'update:shareDialogVisible': [value: boolean]
}

export interface ShareState {
  selectMode: boolean
  selectedFolders: ShareItem[]
  selectedImages: ShareItem[]
}
