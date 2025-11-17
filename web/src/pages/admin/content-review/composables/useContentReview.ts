import { getCurrentLocale } from '@/utils/locale'
import { computed, ref } from 'vue'
import {
  getReviewLogs,
  getReviewQueue,
  getReviewStats,
  type ReviewImage,
  type ReviewLog,
  type ReviewLogQuery,
  type ReviewQueueQuery,
} from '@/api/admin/content-review'
import { getNsfwScoreConfig, getReviewMessages, REVIEW_PAGINATION, TEXT_TRUNCATE } from '@/constants'
import { useTexts } from '@/composables/useTexts'
import { useToast } from '@/components/Toast/useToast'

export function useContentReview() {
  const { $t } = useTexts()

  const NSFW_SCORE_CONFIG = getNsfwScoreConfig($t)
  const REVIEW_MESSAGES = getReviewMessages($t)

  const isLoading = ref(false)
  const currentPage = ref(1)
  const totalPages = ref(1)
  const pageSize = ref(REVIEW_PAGINATION.PAGE_SIZE)

  const formatDateTime = (dateStr: string) => {
    const date = new Date(dateStr)
    return date.toLocaleDateString(getCurrentLocale(), {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    })
  }

  const formatDate = (date: Date | string, format: string = 'MM-DD HH:mm') => {
    const d = new Date(date)
    const month = (d.getMonth() + 1).toString().padStart(2, '0')
    const day = d.getDate().toString().padStart(2, '0')
    const hours = d.getHours().toString().padStart(2, '0')
    const minutes = d.getMinutes().toString().padStart(2, '0')

    return format.replace('MM', month).replace('DD', day).replace('HH', hours).replace('mm', minutes)
  }

  const truncateText = (text: string, maxLength: number = TEXT_TRUNCATE.REASON) => {
    if (text.length <= maxLength) {
      return text
    }
    return `${text.substring(0, maxLength)}...`
  }

  const getNsfwScoreClass = (score: number) => {
    if (score >= NSFW_SCORE_CONFIG.HIGH.threshold) {
      return NSFW_SCORE_CONFIG.HIGH.class
    }
    if (score >= NSFW_SCORE_CONFIG.MEDIUM.threshold) {
      return NSFW_SCORE_CONFIG.MEDIUM.class
    }
    if (score >= NSFW_SCORE_CONFIG.LOW.threshold) {
      return NSFW_SCORE_CONFIG.LOW.class
    }
    return NSFW_SCORE_CONFIG.SAFE.class
  }

  const handleImageError = (event: Event) => {
    const img = event.target as HTMLImageElement
    img.src = '/placeholder-image.png'
  }

  return {
    isLoading,
    currentPage,
    totalPages,
    pageSize,

    formatDateTime,
    formatDate,
    truncateText,
    getNsfwScoreClass,
    handleImageError,

    getReviewQueue,
    getReviewLogs,
    getReviewStats,

    REVIEW_MESSAGES,
    TEXT_TRUNCATE,
  }
}

export function useReviewSelection() {
  const selectMode = ref(false)
  const selectedItems = ref<(string | number)[]>([])
  const totalItems = ref(0)

  const startSelectMode = () => {
    selectMode.value = true
    selectedItems.value = []
  }

  const exitSelectMode = () => {
    selectMode.value = false
    selectedItems.value = []
    totalItems.value = 0
  }

  const toggleItemSelection = (itemId: string | number) => {
    const index = selectedItems.value.indexOf(itemId)
    if (index > -1) {
      selectedItems.value.splice(index, 1)
    } else {
      selectedItems.value.push(itemId)
    }
  }

  const toggleSelectAll = <T extends { id: string | number }>(items: T[]) => {
    totalItems.value = items.length
    if (isAllSelected.value) {
      selectedItems.value = []
    } else {
      selectedItems.value = items.map((item) => item.id)
    }
  }

  const selectInvert = <T extends { id: string | number }>(items: T[]) => {
    totalItems.value = items.length
    const allIds = items.map((item) => item.id)
    selectedItems.value = allIds.filter((id) => !selectedItems.value.includes(id))
  }

  const isAllSelected = computed(() => totalItems.value > 0 && selectedItems.value.length === totalItems.value)

  const hasSelection = computed(() => selectedItems.value.length > 0)

  return {
    selectMode,
    selectedItems,
    startSelectMode,
    exitSelectMode,
    toggleItemSelection,
    toggleSelectAll,
    selectInvert,
    isAllSelected,
    hasSelection,
  }
}

export function useReviewQueue() {
  const { isLoading, currentPage, totalPages, pageSize, REVIEW_MESSAGES } = useContentReview()
  const toast = useToast()

  const images = ref<ReviewImage[]>([])
  const totalImages = ref(0)
  const sortBy = ref('newest')

  const loadData = async (page = 1) => {
    try {
      isLoading.value = true

      const params: ReviewQueueQuery = {
        page,
        size: pageSize.value,
        sort: sortBy.value,
      }

      const response = await getReviewQueue(params)

      if (response.success && response.data) {
        images.value = response.data.data || []
        totalImages.value = response.data.pagination?.total || 0
        const total = response.data.pagination?.total || 0
        const pageSizeFromResponse = response.data.pagination?.page_size || REVIEW_PAGINATION.PAGE_SIZE
        totalPages.value = total > 0 ? Math.ceil(total / pageSizeFromResponse) : 1
        currentPage.value = response.data.pagination?.page || 1
      } else {
        resetData()
      }
    } catch {
      toast.error(REVIEW_MESSAGES.LOAD_ERROR)
      resetData()
    } finally {
      isLoading.value = false
    }
  }

  const resetData = () => {
    images.value = []
    totalImages.value = 0
    totalPages.value = 1
    currentPage.value = 1
  }

  const refreshData = () => {
    loadData(currentPage.value)
  }

  const handlePageChange = (page: number) => {
    loadData(page)
  }

  return {
    images,
    totalImages,
    sortBy,
    isLoading,
    currentPage,
    totalPages,
    pageSize,
    loadData,
    refreshData,
    handlePageChange,
  }
}

export function useReviewLogs() {
  const { isLoading, currentPage, totalPages, pageSize, REVIEW_MESSAGES } = useContentReview()
  const { $t } = useTexts()
  const toast = useToast()

  const logs = ref<ReviewLog[]>([])
  const totalLogs = ref(0)

  const searchKeyword = ref('')
  const actionFilter = ref('')
  const auditorFilter = ref(0)
  const dateRange = ref<[Date | null, Date | null]>([null, null])

  const formatDateForAPI = (date: Date): string => {
    const year = date.getFullYear()
    const month = (date.getMonth() + 1).toString().padStart(2, '0')
    const day = date.getDate().toString().padStart(2, '0')
    return `${year}-${month}-${day}`
  }

  const buildQueryParams = (): ReviewLogQuery => ({
    page: currentPage.value,
    size: pageSize.value,
    keyword: searchKeyword.value || undefined,
    action: (actionFilter.value as string | undefined) || undefined,
    auditor_id: auditorFilter.value || undefined,
    date_from: dateRange.value[0] ? formatDateForAPI(dateRange.value[0]) : undefined,
    date_to: dateRange.value[1] ? formatDateForAPI(dateRange.value[1]) : undefined,
  })

  const loadData = async (page = 1) => {
    try {
      isLoading.value = true
      currentPage.value = page

      const response = await getReviewLogs(buildQueryParams())

      if (response.success && response.data) {
        logs.value = response.data.data || []
        totalLogs.value = response.data.pagination?.total || 0
        const total = response.data.pagination?.total || 0
        const pageSizeFromResponse = response.data.pagination?.page_size || REVIEW_PAGINATION.PAGE_SIZE
        totalPages.value = total > 0 ? Math.ceil(total / pageSizeFromResponse) : 1
        currentPage.value = response.data.pagination?.page || 1
      } else {
        resetData()
      }
    } catch {
      toast.error(REVIEW_MESSAGES.LOAD_ERROR)
      resetData()
    } finally {
      isLoading.value = false
    }
  }

  const resetData = () => {
    logs.value = []
    totalLogs.value = 0
    totalPages.value = 1
    currentPage.value = 1
  }

  const refreshData = () => {
    loadData(currentPage.value)
  }

  const handlePageChange = (page: number) => {
    loadData(page)
  }

  const applyFilters = () => {
    loadData(1)
  }

  const resetFilters = () => {
    searchKeyword.value = ''
    actionFilter.value = ''
    auditorFilter.value = 0
    dateRange.value = [null, null]
    loadData(1)
  }

  const auditorOptions = computed(() => {
    const auditors = new Map()
    logs.value.forEach((log) => {
      if (log.auditor) {
        auditors.set(log.auditor.id, log.auditor)
      }
    })

    const options = [{ label: $t('admin.contentReview.allAuditors'), value: 0 }]
    auditors.forEach((auditor) => {
      options.push({
        label: auditor.username,
        value: auditor.id,
      })
    })

    return options
  })

  return {
    logs,
    totalLogs,
    searchKeyword,
    actionFilter,
    auditorFilter,
    dateRange,
    auditorOptions,
    isLoading,
    currentPage,
    totalPages,
    pageSize,
    loadData,
    refreshData,
    handlePageChange,
    applyFilters,
    resetFilters,
  }
}
