export interface ShareFolderProps {
  folder: {
    id: string
    name: string
    folder_count?: number
    file_count?: number
    created_at?: string
  }
}

export interface ShareFolderEmits {
  (e: 'click', folder: ShareFolderProps['folder']): void
}
