import { ref } from 'vue'
import { useToast } from '@/components/Toast/useToast'
import {
  getVectorList,
  getQdrantRealStats,
  getAvailableModels,
  getVectorLogs,
  type VectorItem,
  type QdrantRealStats,
  type VectorLogItem,
  type PaginationInfo,
  type VectorListParams,
} from '@/api/admin/vectors'
import { useTexts } from '@/composables/useTexts'

export function useVectorData() {
  const toast = useToast()
  const { $t } = useTexts()

  const vectors = ref<VectorItem[]>([])
  const qdrantRealStats = ref<QdrantRealStats | null>(null)
  const logs = ref<VectorLogItem[]>([])
  const availableModels = ref<string[]>([])

  const pagination = ref<PaginationInfo>({
    page: 1,
    page_size: 10,
    total: 0,
    total_page: 0,
  })

  const logPagination = ref<PaginationInfo>({
    page: 1,
    page_size: 20,
    total: 0,
    total_page: 0,
  })

  const isLoading = ref(false)
  const selectedFileId = ref<string>('')

  const loadVectors = async (filters: VectorListParams) => {
    try {
      isLoading.value = true
      const result = await getVectorList(filters)

      if (result && (result.success === true || result.code === 200)) {
        if (result.data && typeof result.data === 'object') {
          const dataArray = Array.isArray(result.data.data) ? result.data.data : []

          if (result.data.pagination) {
            vectors.value = dataArray
            pagination.value = {
              page: result.data.pagination.page || 1,
              page_size: result.data.pagination.page_size || 10,
              total: result.data.pagination.total || 0,
              total_page: result.data.pagination.total_page || 0,
            }
          }
        }
      }
    } catch {
      toast.error($t('admin.vectors.toast.loadFailed'))
    } finally {
      isLoading.value = false
    }
  }

  const loadQdrantStats = async () => {
    try {
      const result = await getQdrantRealStats()

      if (result && (result.success === true || result.code === 200)) {
        if (result.data && typeof result.data === 'object') {
          qdrantRealStats.value = {
            qdrant_vector_count: result.data.qdrant_vector_count || 0,
            qdrant_indexed_count: result.data.qdrant_indexed_count || 0,
            mysql_total_count: result.data.mysql_total_count || 0,
            mysql_completed_count: result.data.mysql_completed_count || 0,
            sync_ratio: result.data.sync_ratio || 0,
            is_healthy: result.data.is_healthy || false,
            collection_status: result.data.collection_status || 'unknown',
            last_checked: result.data.last_checked || new Date().toISOString(),
          }
        }
      }
    } catch {}
  }

  const loadAvailableModels = async () => {
    try {
      const result = await getAvailableModels()

      if ((result.success || result.code === 200) && result.data) {
        availableModels.value = result.data.models || []
      }
    } catch {
      availableModels.value = ['text-embedding-3-small', 'text-embedding-3-large']
    }
  }

  const loadVectorLogs = async (fileId?: string, page = 1, append = false) => {
    try {
      const params: Record<string, unknown> = {
        page,
        page_size: logPagination.value.page_size,
      }

      if (fileId) {
        params.file_id = fileId
      }

      const result = await getVectorLogs(params)

      if ((result.success || result.code === 200) && result.data) {
        const newLogs = Array.isArray(result.data.data) ? result.data.data : []

        if (append) {
          logs.value = [...logs.value, ...newLogs]
        } else {
          logs.value = newLogs
        }

        logPagination.value = result.data.pagination || {
          page,
          page_size: 20,
          total: 0,
          total_page: 0,
        }
      }
    } catch {}
  }

  const cleanup = () => {}

  return {
    vectors,
    qdrantRealStats,
    logs,
    availableModels,
    pagination,
    logPagination,
    isLoading,
    selectedFileId,

    loadVectors,
    loadQdrantStats,
    loadAvailableModels,
    loadVectorLogs,
    cleanup,
  }
}
