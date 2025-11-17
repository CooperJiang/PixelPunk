import { ref } from 'vue'
import type { ImageListParams } from '@/api/admin/files'

export function useFileFilter(filters: ImageListParams, onFilterChange: () => void) {
  const showFilter = ref(false)

  const handleFilter = (newFilters: Partial<ImageListParams>) => {
    Object.assign(filters, newFilters)
    filters.page = 1
    onFilterChange()
  }

  const resetFilter = () => {
    Object.assign(filters, {
      page: 1,
      page_size: filters.page_size, // 保持页面大小
      keyword: '',
      uploader_id: null,
      channel: '',
      is_recommended: null,
      start_date: '',
      end_date: '',
      sort: 'created_at',
      order: 'desc',
    })
    showFilter.value = false
    onFilterChange()
  }

  const toggleFilter = () => {
    showFilter.value = !showFilter.value
  }

  return {
    showFilter,
    handleFilter,
    resetFilter,
    toggleFilter,
  }
}
