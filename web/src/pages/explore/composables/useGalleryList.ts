import { computed, nextTick, ref, onMounted, onUnmounted } from 'vue'
import type { FileInfo } from '@/api/types'
import { getGuestFileList } from '@/api/common'
import { processParams } from '@/utils/formatting/format'
import { useMultipleLoading } from '@/hooks/useLoading'
import type { UseGalleryListOptions } from '../types'

export function useGalleryList(options: UseGalleryListOptions) {
  const { images, filters, searchKeyword, isVectorSearchMode, pageSize: initialPageSize = 24, maxImages = 5000 } = options

  const currentPage = ref(1)
  const pageSize = ref(initialPageSize)
  const hasMore = ref(true)
  const hasReachedLimit = ref(false)

  const { loading: loadingStates, withLoading } = useMultipleLoading({
    isLoading: false,
    isLoadingMore: false,
  })

  const isLoading = computed(() => loadingStates.value?.isLoading || false)
  const isLoadingMore = computed(() => loadingStates.value?.isLoadingMore || false)

  const checkAndLoadMoreIfNeeded = async () => {
    if (
      !isVectorSearchMode ||
      !images ||
      (isVectorSearchMode && isVectorSearchMode.value) ||
      !images.value ||
      images.value.length === 0 ||
      hasReachedLimit.value ||
      !hasMore.value
    ) {
      return
    }

    await nextTick()

    const { height, scroll } = getScrollMetricsFrom()
    const contentHeight = scroll
    const viewportHeight = height
    const threshold = viewportHeight * 1.5

    if (contentHeight < threshold && hasMore.value && !isLoadingMore.value) {
      currentPage.value++
      await loadImages(true)
      await checkAndLoadMoreIfNeeded()
    }
  }

  const loadImages = async (append = false) => {
    if (!append) {
      images.value = []
      currentPage.value = 1
      hasMore.value = true
      hasReachedLimit.value = false
    }

    const loadingKey = append ? 'isLoadingMore' : 'isLoading'

    await withLoading(
      loadingKey as any,
      async () => {
        const params = processParams(
          {
            ...filters,
            page: currentPage.value,
            size: pageSize.value,
            keyword: (searchKeyword && searchKeyword.value) || '',
          },
          {
            removeUndefined: true,
            removeEmpty: true,
          }
        )

        const result = await getGuestFileList(params)

        if (result.success) {
          const { data } = result
          const newFiles = (data.items || []).map((file: FileInfo) => {
            if (file.ai_info?.description) {
              file.description = file.ai_info.description
            }
            return file
          }) as FileInfo[]

          if (append) {
            images.value = [...images.value, ...newFiles]
          } else {
            images.value = newFiles
          }
          const currentImageCount = images.value.length
          const total = data.pagination?.total
          if (typeof total === 'number' && total >= 0) {
            hasMore.value = currentImageCount < total && newFiles.length > 0
          } else {
            hasMore.value = newFiles.length >= pageSize.value
          }

          if (images.value.length >= maxImages) {
            hasReachedLimit.value = true
            hasMore.value = false
            images.value = images.value.slice(0, maxImages)
          }

          if (!append && newFiles.length > 0) {
            setTimeout(() => {
              checkAndLoadMoreIfNeeded()
            }, 300)
          }
        } else {
          hasMore.value = false
        }
      },
      {
        onError: (error) => {
          console.error('Load files failed:', error)
          hasMore.value = false
        },
      }
    )
  }

  const loadMore = async () => {
    if (isLoadingMore.value || isLoading.value || !hasMore.value || hasReachedLimit.value) return
    currentPage.value++
    await loadImages(true)
  }

  const hasUserScrolled = ref(false)
  const markUserScrolled = () => {
    hasUserScrolled.value = true
  }
  const onKeyDown = (e: KeyboardEvent) => {
    const keys = ['ArrowDown', 'ArrowUp', 'PageDown', 'PageUp', 'Space']
    if (keys.includes(e.code)) {
      hasUserScrolled.value = true
    }
  }
  onMounted(() => {
    window.addEventListener('wheel', markUserScrolled, { passive: true })
    window.addEventListener('touchmove', markUserScrolled, { passive: true })
    window.addEventListener('keydown', onKeyDown)
  })
  onUnmounted(() => {
    window.removeEventListener('wheel', markUserScrolled as any)
    window.removeEventListener('touchmove', markUserScrolled as any)
    window.removeEventListener('keydown', onKeyDown)
  })

  const scrollTarget = ref<HTMLElement | Window>(window)
  const pickScrollTarget = () => {
    const candidates: (HTMLElement | null)[] = [
      document.querySelector('.main-content'),
      document.querySelector('.gallery-images-fullscreen'),
      document.querySelector('.gallery-container'),
      document.querySelector('.gallery-page'),
    ]
    for (const el of candidates) {
      if (!el) continue
      const style = window.getComputedStyle(el)
      const overflowY = style.overflowY
      const isScrollable = overflowY === 'auto' || overflowY === 'scroll'
      if (isScrollable) return el
    }
    return window as unknown as Window
  }

  const getScrollMetricsFrom = (el?: any) => {
    const t = el || scrollTarget.value
    if (t === window || t === document || t === document.documentElement || t === document.body) {
      const docEl = document.documentElement
      const body = document.body
      const top = window.scrollY || docEl.scrollTop || body.scrollTop || 0
      const height = window.innerHeight
      const scroll = docEl.scrollHeight
      return { top, height, scroll }
    }
    const node = t as HTMLElement
    return { top: node.scrollTop, height: node.clientHeight, scroll: node.scrollHeight }
  }
  const BOTTOM_THRESHOLD = 40
  let scrollRafPending = false
  const checkShouldLoad = (source?: any) => {
    if (!hasUserScrolled.value) return false
    if (isLoadingMore.value || isLoading.value) return false
    if (!hasMore.value || hasReachedLimit.value) return false
    if (isVectorSearchMode && isVectorSearchMode.value) return false

    const { top, height, scroll } = getScrollMetricsFrom(source)
    const distanceToBottom = scroll - (top + height)
    return distanceToBottom <= BOTTOM_THRESHOLD
  }

  const onScroll = (e?: Event) => {
    hasUserScrolled.value = true
    if (scrollRafPending) return
    scrollRafPending = true
    requestAnimationFrame(async () => {
      scrollRafPending = false
      const src = (e && (e.target || (e.currentTarget as any))) || undefined
      if (checkShouldLoad(src)) await loadMore()
    })
  }

  onMounted(() => {
    scrollTarget.value = pickScrollTarget()
    const tgt: any = scrollTarget.value
    tgt.addEventListener('scroll', onScroll, { passive: true })
  })
  onUnmounted(() => {
    const tgt: any = scrollTarget.value
    if (tgt) {
      tgt.removeEventListener('scroll', onScroll)
    }
  })

  return {
    images,
    currentPage,
    pageSize,
    hasMore,
    hasReachedLimit,
    isLoading,
    isLoadingMore,

    loadImages,
    loadMore,
    withPageLoading: withLoading,
  }
}
