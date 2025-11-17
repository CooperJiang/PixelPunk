/**
 * 自动任务数据管理
 */
import { ref, onMounted, onUnmounted } from 'vue'
import { useToast } from '@/components/Toast/useToast'
import { useTexts } from '@/composables/useTexts'
import { getUserAutomationOverview, type AutomationOverview } from '@/api/user/automation'

export function useAutomationData() {
  const { $t } = useTexts()
  const toast = useToast()

  const isLoading = ref(false)
  const overview = ref<AutomationOverview | null>(null)
  const lastUpdated = ref<Date | null>(null)

  const pollingInterval = ref(5000) // 默认5秒轮询一次
  let pollingTimer: number | null = null

  const loadData = async (showLoading = true) => {
    try {
      if (showLoading) {
        isLoading.value = true
      }

      const result = await getUserAutomationOverview()
      if (result.success) {
        overview.value = result.data
        lastUpdated.value = new Date()
      }
    } catch (error: any) {
      if (showLoading) {
        toast.error(error.message || $t('automation.toast.loadDataError'))
      }
    } finally {
      if (showLoading) {
        isLoading.value = false
      }
    }
  }

  const startPolling = () => {
    if (pollingTimer) {
      return
    }

    pollingTimer = window.setInterval(() => {
      loadData(false) // 轮询时不显示loading状态
    }, pollingInterval.value)
  }

  const stopPolling = () => {
    if (pollingTimer) {
      clearInterval(pollingTimer)
      pollingTimer = null
    }
  }

  const refresh = async () => {
    await loadData(true)
  }

  onMounted(() => {
    loadData()
    startPolling()
  })

  onUnmounted(() => {
    stopPolling()
  })

  return {
    isLoading,
    overview,
    lastUpdated,
    loadData,
    refresh,
    startPolling,
    stopPolling,
  }
}
