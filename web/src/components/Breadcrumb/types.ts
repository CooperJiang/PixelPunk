import type { FolderInfo } from '@/api/types'

export interface BreadcrumbProps {
  items: FolderInfo[]
  maxItems?: number
  startVisible?: number
  endVisible?: number
}

export interface BreadcrumbEmits {
  (e: 'click', folder: FolderInfo | null): void
}
