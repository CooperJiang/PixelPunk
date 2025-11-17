import { watch } from 'vue'
import { galleryVectorSearch } from '@/api/search'
import type { UseGallerySearchOptions } from '../types'
import { useTexts } from '@/composables/useTexts'

export function useGallerySearch(opts: UseGallerySearchOptions) {
  const { $t } = useTexts()
  const {
    images,
    currentPage,
    hasMore,
    hasReachedLimit,
    withPageLoading,
    loadImages,
    searchKeyword,
    searchMode,
    isVectorSearchMode,
    lastVectorSearchQuery,
    toast,
  } = opts

  const handleNormalSearch = () => {
    if (isVectorSearchMode.value) {
      isVectorSearchMode.value = false
      lastVectorSearchQuery.value = ''
    }
    loadImages(false)
  }

  const handleVectorSearch = async () => {
    const query = searchKeyword.value.trim()

    // 允许空搜索，回退到正常模式
    if (!query) {
      isVectorSearchMode.value = false
      lastVectorSearchQuery.value = ''
      loadImages(false)
      return
    }

    images.value = []
    currentPage.value = 1
    hasMore.value = false
    hasReachedLimit.value = false

    await withPageLoading(
      'isLoading',
      async () => {
        const result = await galleryVectorSearch({ query, page: 1 })
        if (result && result.data && result.data.items) {
          const { data } = result
          const newImages = data.items.map((item: any) => {
            if (item.ai_info?.description) item.description = item.ai_info.description
            if (item.similarity !== undefined) item.similarity = item.similarity
            return item
          })
          images.value = newImages
          isVectorSearchMode.value = true
          lastVectorSearchQuery.value = query
          toast.success($t('explore.toast.foundResults').replace('{count}', data.pagination.total.toString()))
        } else {
          images.value = []
          toast.info($t('explore.toast.noResultsFound'))
        }
      },
      {
        onError: (_error) => {
          toast.error($t('explore.toast.searchFailed'))
        },
      }
    )
  }

  const handleAdvancedSearch = ({ keyword, mode }: { keyword: string; mode: 'normal' | 'vector' }) => {
    searchKeyword.value = keyword
    searchMode.value = mode
    if (mode === 'vector') handleVectorSearch()
    else handleNormalSearch()
  }

  watch(searchMode, (newMode) => {
    if (newMode === 'normal' && isVectorSearchMode.value) {
      isVectorSearchMode.value = false
      lastVectorSearchQuery.value = ''
      if (searchKeyword.value) handleNormalSearch()
      else loadImages(false)
    }
  })

  return {
    handleNormalSearch,
    handleVectorSearch,
    handleAdvancedSearch,
  }
}
