import { ref, computed } from 'vue'
import { getGuestTagList } from '@/api/common'

interface TagOption {
  id: number
  name: string
  count: number
}

interface PaginatedTagsState {
  tags: TagOption[]
  loading: boolean
  loadingMore: boolean
  hasMore: boolean
  currentPage: number
  pageSize: number
  searchQuery: string
}

export function usePaginatedTags(initialPageSize = 20) {
  const state = ref<PaginatedTagsState>({
    tags: [],
    loading: false,
    loadingMore: false,
    hasMore: true,
    currentPage: 1,
    pageSize: initialPageSize,
    searchQuery: '',
  })

  const tagOptions = computed(() =>
    state.value.tags.map((tag) => ({
      label: `${tag.name} (${tag.count})`,
      value: tag.name,
    }))
  )

  const loadTags = async (page = 1, append = false) => {
    if (page === 1) {
      state.value.loading = true
    } else {
      state.value.loadingMore = true
    }

    try {
      const result = await getGuestTagList(page, state.value.pageSize)
      if (result.success && result.data) {
        const newTags = result.data.items || []

        if (append) {
          state.value.tags = [...state.value.tags, ...newTags]
        } else {
          state.value.tags = newTags
          state.value.currentPage = 1
        }

        const pagination = result.data.pagination
        if (pagination) {
          state.value.hasMore = pagination.current_page < pagination.last_page
          state.value.currentPage = pagination.current_page
        } else {
          state.value.hasMore = newTags.length >= state.value.pageSize
        }
      }
    } catch (error) {
      console.error('Load tags failed:', error)
      state.value.hasMore = false
    } finally {
      state.value.loading = false
      state.value.loadingMore = false
    }
  }

  const loadMore = async () => {
    if (state.value.loadingMore || !state.value.hasMore) return
    await loadTags(state.value.currentPage + 1, true)
  }

  const searchTags = async (query: string) => {
    state.value.searchQuery = query
    state.value.currentPage = 1
    state.value.hasMore = true
    await loadTags(1, false)
  }

  const initialize = async () => {
    await loadTags(1, false)
  }

  return {
    tags: computed(() => state.value.tags),
    tagOptions,
    loading: computed(() => state.value.loading),
    loadingMore: computed(() => state.value.loadingMore),
    hasMore: computed(() => state.value.hasMore),

    loadMore,
    searchTags,
    initialize,
  }
}
