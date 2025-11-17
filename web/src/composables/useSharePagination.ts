import { computed, nextTick, ref, type ComputedRef } from 'vue'
import { INITIAL_LOAD_COUNT, PAGE_SIZE } from '@/constants'

/**
 * 分享分页功能
 * 处理文件列表的分页加载和显示
 */
export function useSharePagination<T>(sourceList: ComputedRef<T[]>, preloadImageFn?: (url: string) => void) {
  const currentPage = ref(1)

  const showAllImages = ref(false)

  const isLoading = ref(false)
  const pendingLoadMore = ref(false) // 标记是否有待处理的加载请求

  const paginatedList = computed(() => {
    if (!showAllImages.value) {
      return sourceList.value.slice(0, currentPage.value * PAGE_SIZE)
    }

    return sourceList.value
  })

  const canLoadMore = computed(
    () => !showAllImages.value && currentPage.value * PAGE_SIZE < sourceList.value.length && !isLoading.value
  )

  const loadMoreImages = async () => {
    if (isLoading.value) {
      pendingLoadMore.value = true
      return
    }

    if (showAllImages.value || currentPage.value * PAGE_SIZE >= sourceList.value.length) {
      showAllImages.value = true
      pendingLoadMore.value = false
      return
    }

    isLoading.value = true
    pendingLoadMore.value = false

    try {
      currentPage.value++

      if (preloadImageFn) {
        await nextTick()
        const start = (currentPage.value - 1) * PAGE_SIZE
        const end = Math.min(currentPage.value * PAGE_SIZE, sourceList.value.length)

        const preloadPromises = []
        for (let i = start; i < end; i++) {
          const image = sourceList.value[i] as any
          if (image) {
            const imageUrl = image.full_thumb_url || image.thumb_url || image.full_url || image.url
            if (imageUrl) {
              preloadPromises.push(preloadImageFn(imageUrl))
            }
          }
        }

        if (preloadPromises.length > 0) {
          try {
            await Promise.allSettled(preloadPromises)
          } catch (error) {
            console.warn('预加载部分文件失败:', error)
          }
        }
      }

      if (currentPage.value * PAGE_SIZE >= sourceList.value.length) {
        showAllImages.value = true
      }
    } finally {
      isLoading.value = false

      if (pendingLoadMore.value && canLoadMore.value) {
        nextTick(() => {
          loadMoreImages()
        })
      }
    }
  }

  const showLoadMoreButton = computed(() => sourceList.value.length > PAGE_SIZE && !showAllImages.value)

  const resetPagination = () => {
    currentPage.value = 1
    showAllImages.value = false
    isLoading.value = false
    pendingLoadMore.value = false
  }

  const preloadInitialImages = () => {
    if (!preloadImageFn) {
      return
    }

    nextTick(() => {
      const count = Math.min(INITIAL_LOAD_COUNT, sourceList.value.length)

      for (let i = 0; i < count; i++) {
        const image = sourceList.value[i] as any
        if (image) {
          const imageUrl = image.full_thumb_url || image.thumb_url || image.full_url || image.url
          if (imageUrl) {
            preloadImageFn(imageUrl)
          }
        }
      }
    })
  }

  return {
    currentPage,
    showAllImages,
    paginatedList,
    loadMoreImages,
    showLoadMoreButton,
    resetPagination,
    preloadInitialImages,
    isLoading,
    canLoadMore,
    pendingLoadMore,
  }
}
