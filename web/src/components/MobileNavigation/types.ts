import type { FolderInfo } from '@/api/types'

export interface MobileNavigationProps {
  items?: FolderInfo[]
  title?: string
  hasPassword?: boolean
  showNavigation?: boolean
}

export interface MobileNavigationEmits {
  (e: 'click', folder: FolderInfo | null): void
}
