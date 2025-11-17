export interface FolderInfo {
  id: string
  name: string
  fullPath?: string
  path?: string
  children?: FolderInfo[]
  has_children?: boolean
}

export interface FolderTreeProps {
  modelValue?: string
  borderColor?: string
  presetPath?: string
  maxPathLength?: number
  width?: string
  height?: string
  dropdownZIndex?: number
}

export interface FolderTreeEmits {
  (e: 'update:modelValue', value: string): void
  (e: 'folderSelected', folder: { id: string; name: string; path: string }): void
}

export interface FolderTreeItemProps {
  folder: FolderInfo
  level: number
  selectedId: string
  borderColor?: string
  isExpanded: boolean
  expandedFolders: Set<string>
  cachedChildren: FolderInfo[]
  getFolderChildren: (folderId: string) => FolderInfo[]
}

export interface FolderTreeItemEmits {
  (e: 'select', folder: FolderInfo): void
  (e: 'toggle', folderId: string, isExpanded: boolean): void
  (e: 'children-loaded', folderId: string, children: FolderInfo[]): void
  (e: 'subfolder-request-start'): void
  (e: 'subfolder-request-end'): void
}
