export interface PaginationProps {
  currentPage: number
  total: number
  pageSize?: number
  pageSizeOptions?: number[]
  showPageSizeSelector?: boolean
  showQuickJumper?: boolean
  showTotal?: boolean
  maxVisiblePages?: number
  autoScrollToTop?: boolean
  scrollTarget?: string
}

export interface PaginationEmits {
  'update:currentPage': [page: number]
  'update:pageSize': [size: number]
  'page-change': [page: number]
  'page-size-change': [size: number]
}
