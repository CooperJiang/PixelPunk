import { nextTick, ref } from 'vue'
import { useTexts } from '@/composables/useTexts'
import { getGuestFileList } from '@/api/common'
import type { FileInfo, FileListParams } from '@/api/types'
import type { ViewportDimensions } from '../types'
import {
  PAGE_SIZE,
  PROGRESS_INITIAL,
  PROGRESS_DATA_LOADED,
  PROGRESS_PRELOAD_START,
  PROGRESS_PRELOAD_RANGE,
  PROGRESS_COMPLETE,
  ERROR_DISPLAY_DURATION,
  ITEM_SIZE,
  ROW_HEIGHT_RATIO,
  ITEMS_PER_ROW_BUFFER,
  VISIBLE_ROWS_BUFFER,
  MAX_VISIBLE_COUNT,
} from '../constants'

export function useHiveData() {
  const { $t } = useTexts()
  const images = ref<FileInfo[]>([])
  const isLoading = ref(false)
  const isInitialLoading = ref(true)
  const isLoadingFadingOut = ref(false)
  const error = ref('')
  const currentPage = ref(1)
  const hasMorePages = ref(true)
  const loadingProgress = ref(0)
  const loadingStartTime = ref(Date.now())

  let canvasRefCallback: (() => void) | null = null

  const setCanvasUpdateCallback = (callback: () => void) => {
    canvasRefCallback = callback
  }

  const getViewportDimensions = (): ViewportDimensions => ({
    width: window.innerWidth,
    height: window.innerHeight,
  })

  const calculateVisibleImageCount = (): number => {
    const { width, height } = getViewportDimensions()
    const itemsPerRow = Math.floor(width / ITEM_SIZE) + ITEMS_PER_ROW_BUFFER
    const rowsVisible = Math.floor(height / (ITEM_SIZE * ROW_HEIGHT_RATIO)) + VISIBLE_ROWS_BUFFER
    return Math.min(itemsPerRow * rowsVisible, MAX_VISIBLE_COUNT)
  }

  const loadFiles = async (page = 1, append = false) => {
    if (isLoading.value || !hasMorePages.value) {
      return
    }

    isLoading.value = true
    error.value = ''

    if (!append && isInitialLoading.value) {
      loadingProgress.value = PROGRESS_INITIAL
    }

    try {
      const params: FileListParams = {
        page,
        size: PAGE_SIZE,
        sort: 'newest',
      }

      const response = await getGuestFileList(params)
      const newImages = response?.data?.items || []

      if (!append && isInitialLoading.value) {
        loadingProgress.value = PROGRESS_DATA_LOADED
      }

      if (newImages.length === 0) {
        hasMorePages.value = false
        return
      }

      if (append) {
        images.value = [...images.value, ...newImages]
      } else {
        images.value = newImages
        currentPage.value = page
        if (isInitialLoading.value) {
          loadingProgress.value = PROGRESS_PRELOAD_START
        }
      }

      if (response.data.pagination) {
        hasMorePages.value = response.data.pagination.current_page < response.data.pagination.last_page
      }

      await nextTick()
      canvasRefCallback?.()
    } catch (err) {
      const message = err instanceof Error ? err.message : $t('hive.errors.unknown')
      error.value = $t('hive.errors.loadFailedWithReason', { message })
      if (isInitialLoading.value) {
        isInitialLoading.value = false
      }
      setTimeout(() => {
        error.value = ''
      }, ERROR_DISPLAY_DURATION)
    } finally {
      isLoading.value = false
    }
  }

  const loadMoreImages = () => {
    if (hasMorePages.value && !isLoading.value) {
      currentPage.value++
      loadFiles(currentPage.value, true)
    }
  }

  const preloadVisibleFiles = async (fileList: FileInfo[]) => {
    const visibleCount = calculateVisibleImageCount()
    const filesToLoad = fileList.slice(0, visibleCount)

    if (filesToLoad.length === 0) {
      loadingProgress.value = PROGRESS_COMPLETE
      return
    }

    let loadedCount = 0
    const totalCount = filesToLoad.length

    const loadFile = (file: FileInfo): Promise<void> =>
      new Promise<void>((resolve) => {
        const img = new window.Image()
        const updateProgress = () => {
          loadedCount++
          const progress = PROGRESS_PRELOAD_START + Math.round((loadedCount / totalCount) * PROGRESS_PRELOAD_RANGE)
          loadingProgress.value = progress
          resolve()
        }

        img.onload = updateProgress
        img.onerror = updateProgress

        const fileUrl = file.full_thumb_url || file.thumb_url
        if (fileUrl) {
          img.src = fileUrl
        } else {
          updateProgress()
        }
      })

    await Promise.all(filesToLoad.map(loadFile))
    loadingProgress.value = PROGRESS_COMPLETE
  }

  return {
    images,
    isLoading,
    isInitialLoading,
    isLoadingFadingOut,
    error,
    currentPage,
    hasMorePages,
    loadingProgress,
    loadingStartTime,
    loadFiles,
    loadMoreImages,
    preloadVisibleFiles,
    setCanvasUpdateCallback,
  }
}
