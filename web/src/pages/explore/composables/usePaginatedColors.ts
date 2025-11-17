import { ref, computed } from 'vue'
import { getGuestColorList } from '@/api/common'

interface PaginatedColorsState {
  colors: string[]
  loading: boolean
  loadingMore: boolean
  hasMore: boolean
  currentPage: number
  pageSize: number
  searchQuery: string
}

export function usePaginatedColors(initialPageSize = 20) {
  const state = ref<PaginatedColorsState>({
    colors: [],
    loading: false,
    loadingMore: false,
    hasMore: true,
    currentPage: 1,
    pageSize: initialPageSize,
    searchQuery: '',
  })

  const colorOptions = computed(() =>
    state.value.colors.map((color) => ({
      label: color,
      value: color,
      color,
    }))
  )

  const loadColors = async (page = 1, append = false) => {
    if (page === 1) {
      state.value.loading = true
    } else {
      state.value.loadingMore = true
    }

    try {
      const result = await getGuestColorList(page, state.value.pageSize)
      if (result.success && result.data) {
        const newColors = result.data.items || []

        if (append) {
          state.value.colors = [...state.value.colors, ...newColors]
        } else {
          state.value.colors = newColors
          state.value.currentPage = 1
        }

        const pagination = result.data.pagination
        if (pagination) {
          state.value.hasMore = pagination.current_page < pagination.last_page
          state.value.currentPage = pagination.current_page
        } else {
          state.value.hasMore = newColors.length >= state.value.pageSize
        }
      }
    } catch (error) {
      console.error('Load colors failed:', error)
      state.value.hasMore = false
    } finally {
      state.value.loading = false
      state.value.loadingMore = false
    }
  }

  const loadMore = async () => {
    if (state.value.loadingMore || !state.value.hasMore) return
    await loadColors(state.value.currentPage + 1, true)
  }

  const searchColors = async (query: string) => {
    state.value.searchQuery = query
    state.value.currentPage = 1
    state.value.hasMore = true
    await loadColors(1, false)
  }

  const initialize = async () => {
    await loadColors(1, false)
  }

  return {
    colors: computed(() => state.value.colors),
    colorOptions,
    loading: computed(() => state.value.loading),
    loadingMore: computed(() => state.value.loadingMore),
    hasMore: computed(() => state.value.hasMore),

    loadMore,
    searchColors,
    initialize,
  }
}
