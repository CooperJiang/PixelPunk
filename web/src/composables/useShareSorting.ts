import { computed, ref, type Ref } from 'vue'
import { SORT_DIRECTIONS } from '@/constants'

/**
 * 分享排序功能
 * 处理文件列表的排序逻辑
 */
export function useShareSorting<T>(sourceList: Ref<T[]>) {
  const currentSort = ref('date')

  const sortDirection = ref(SORT_DIRECTIONS.DESC)

  const handleSortChange = (sortType: string) => {
    if (currentSort.value === sortType) {
      sortDirection.value = sortDirection.value === SORT_DIRECTIONS.ASC ? SORT_DIRECTIONS.DESC : SORT_DIRECTIONS.ASC
    } else {
      currentSort.value = sortType
      sortDirection.value = sortType === 'date' ? SORT_DIRECTIONS.DESC : SORT_DIRECTIONS.ASC
    }
  }

  const sortedList = computed(() => {
    if (!sourceList.value || sourceList.value.length === 0) {
      return []
    }

    return [...sourceList.value].sort((a: any, b: any) => {
      switch (currentSort.value) {
        case 'date':
          return sortDirection.value === SORT_DIRECTIONS.ASC
            ? new Date(a.created_at).getTime() - new Date(b.created_at).getTime()
            : new Date(b.created_at).getTime() - new Date(a.created_at).getTime()
        case 'name':
          const nameA = a.display_name || a.original_name || ''
          const nameB = b.display_name || b.original_name || ''
          return sortDirection.value === SORT_DIRECTIONS.ASC ? nameA.localeCompare(nameB) : nameB.localeCompare(nameA)
        case 'size':
          const sizeA = a.size || 0
          const sizeB = b.size || 0
          return sortDirection.value === SORT_DIRECTIONS.ASC ? sizeA - sizeB : sizeB - sizeA
        default:
          return 0
      }
    })
  })

  return {
    currentSort,
    sortDirection,
    handleSortChange,
    sortedList,
  }
}
