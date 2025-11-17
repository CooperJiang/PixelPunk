import type { FileListParams } from '@/api/types'
import type { UseGalleryFiltersOptions } from '../types'

export function useGalleryFilters(options: UseGalleryFiltersOptions) {
  const { filters, loadFiles, defaultPageSize, searchKeyword } = options

  const getFilterCount = (): number => {
    let count = 0
    if (filters.sort && filters.sort !== 'newest') count++
    if (filters.dominant_color) count++
    if (filters.tags) count++
    if (filters.is_recommended !== undefined) count++
    if (filters.resolution) count++
    if (filters.min_width || filters.max_width || filters.min_height || filters.max_height) count++
    return count
  }

  const quickSort = (sortValue: string) => {
    filters.sort = sortValue
    loadFiles(false)
  }

  const applyFilter = (newFilters: FileListParams) => {
    Object.assign(filters, newFilters)
    loadFiles(false)
  }

  const resetFilter = () => {
    Object.assign(filters, {
      page: 1,
      size: defaultPageSize,
      sort: 'newest',
      keyword: '',
      tags: '',
      dominant_color: '',
      is_recommended: undefined,
      resolution: '',
      min_width: undefined,
      max_width: undefined,
      min_height: undefined,
      max_height: undefined,
    })

    searchKeyword.value = ''
    loadFiles(false)
  }

  return {
    getFilterCount,
    quickSort,
    applyFilter,
    resetFilter,
  }
}
