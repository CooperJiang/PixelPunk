/**
 * 打标历史记录管理
 */
import { ref, reactive } from 'vue'
import { useToast } from '@/components/Toast/useToast'
import { useTexts } from '@/composables/useTexts'
import { getUserTaggingTasks, type TaggingTaskItem } from '@/api/user/automation'

/* 模块级别的状态，确保单例 */
const sharedState = {
  isLoading: ref(false),
  data: ref<TaggingTaskItem[]>([]),
  total: ref(0),
  query: reactive({
    page: 1,
    limit: 10,
    status: '', // 状态筛选
  }),
  hasLoaded: ref(false), // 标记是否已经加载过（使用ref使其响应式）
}

export function useTaggingHistory() {
  const { $t } = useTexts()
  const toast = useToast()

  const loadData = async () => {
    if (sharedState.isLoading.value) {
      return
    }

    try {
      sharedState.isLoading.value = true

      const result = await getUserTaggingTasks({
        status: sharedState.query.status || undefined,
        page: sharedState.query.page,
        limit: sharedState.query.limit,
      })

      if (result.success) {
        sharedState.data.value = result.data.data || []
        sharedState.total.value = result.data.total || 0
        sharedState.hasLoaded.value = true
      }
    } catch (error: any) {
      console.error('[useTaggingHistory] Failed to load tagging history', error)
      toast.error(error.message || $t('automation.toast.loadHistoryError'))
    } finally {
      sharedState.isLoading.value = false
    }
  }

  const handlePageChange = (page: number) => {
    sharedState.query.page = page
    loadData()
  }

  const handleSizeChange = (size: number) => {
    sharedState.query.limit = size
    sharedState.query.page = 1
    loadData()
  }

  const handleStatusFilter = (status: string) => {
    sharedState.query.status = status
    sharedState.query.page = 1
    loadData()
  }

  const refresh = () => {
    loadData()
  }

  return {
    isLoading: sharedState.isLoading,
    data: sharedState.data,
    total: sharedState.total,
    query: sharedState.query,
    hasLoaded: sharedState.hasLoaded,
    loadData,
    handlePageChange,
    handleSizeChange,
    handleStatusFilter,
    refresh,
  }
}
