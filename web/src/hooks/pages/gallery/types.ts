/* Gallery页面通用类型定义 */

export interface ImageDetail {
  id: string
  url: string
  full_url: string
  full_thumb_url: string
  original_name: string
  display_name: string
  size: number
  width: number
  height: number
  format: string
  access_level: string
  created_at: string
  updated_at: string
  ai_info?: {
    description: string
    tags: string[]
    dominant_color: string
    resolution: string
    is_nsfw: boolean
    nsfw_score: number
    nsfw_evaluation: string
  }
  description?: string // 兼容字段
}

export interface Pagination {
  current_page: number
  last_page: number
  size: number
  total: number
}

export interface FilterParams {
  page: number
  size: number
  access_level: string
  sort: string
  keyword: string
  tags: string
  dominant_color: string
  resolution: string
  min_width?: number
  min_height?: number
}

export interface GalleryState {
  images: ImageDetail[]
  pagination: Pagination
  isLoading: boolean
  currentPage: number
  showFilter: boolean
  filterParams: FilterParams
}

export interface ImageModalState {
  showFilePreview: boolean
  showDetailModal: boolean
  selectedImageDetails: ImageDetail | null
  currentPreviewIndex: number
}

export interface ImageSelectionState {
  selectMode: boolean
  selectedImages: Set<string>
}

export interface GalleryOptions {
  defaultPageSize?: number
  enableInfiniteScroll?: boolean
  apiEndpoint?: string
  enableSelection?: boolean
}
