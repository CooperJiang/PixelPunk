import type { FileInfo, FileListParams } from '@/api/types'
import type { Ref } from 'vue'

export type Density = 'compact' | 'normal' | 'comfortable'
export type SearchMode = 'normal' | 'vector'

export interface SortOption {
  value: string
  label: string
  icon: string
}

export interface DensityOption {
  value: Density
  label: string
  icon: string
}

export interface UseGalleryFiltersOptions {
  filters: FileListParams
  loadFiles: (append?: boolean) => Promise<void>
  defaultPageSize: number
  searchKeyword: Ref<string>
}

export interface UseGalleryListOptions {
  images: Ref<FileInfo[]>
  filters: FileListParams
  searchKeyword: Ref<string>
  isVectorSearchMode: Ref<boolean>
  pageSize?: number
  maxImages?: number
}

export interface UseGallerySearchOptions {
  images: Ref<FileInfo[]>
  currentPage: Ref<number>
  hasMore: Ref<boolean>
  hasReachedLimit: Ref<boolean>
  withPageLoading: <T>(key: string, fn: () => Promise<T>) => Promise<T>
  loadImages: (append?: boolean) => Promise<void>
  searchKeyword: Ref<string>
  searchMode: Ref<SearchMode>
  isVectorSearchMode: Ref<boolean>
  lastVectorSearchQuery: Ref<string>
  toast: {
    info: (message: string) => void
    warning: (message: string) => void
    error: (message: string) => void
  }
}

export interface ColorOption {
  color: string
  count: number
  label?: string
}

export interface PaginatedColorsState {
  allColors: ColorOption[]
  visibleColors: ColorOption[]
  showAllColors: boolean
}

export interface TagOption {
  tag: string
  count: number
}

export interface PaginatedTagsState {
  allTags: TagOption[]
  visibleTags: TagOption[]
  showAllTags: boolean
}

export interface ShareButtonRef {
  selectedFiles?: FileInfo[]
  selectedFolders?: unknown[]
  isFileSelected?: (file: FileInfo) => boolean
  toggleFileSelect?: (file: FileInfo) => void
  cancelSelectMode?: () => void
}
