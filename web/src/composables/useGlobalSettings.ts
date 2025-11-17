import { computed, onMounted, ref } from 'vue'
import { getGlobalSettings, type GlobalSettingsResponse } from '@/api/admin/settings'
import { useToast } from '@/components/Toast/useToast'
import { useSettingsStore } from '@/store/settings'

export function useGlobalSettings() {
  const _toast = useToast()
  const settingsStore = useSettingsStore()
  const isLoading = ref(false)
  const settings = ref<GlobalSettingsResponse | null>(null)

  const fetchSettings = async () => {
    if (settingsStore.isLoaded && settingsStore.rawSettings) {
      settings.value = settingsStore.rawSettings
      return
    }

    if (settingsStore.loading) {
      let waited = 0
      const maxWait = 5000
      const checkInterval = setInterval(() => {
        waited += 100
        if (!settingsStore.loading && settingsStore.isLoaded && settingsStore.rawSettings) {
          clearInterval(checkInterval)
          settings.value = settingsStore.rawSettings
        } else if (waited >= maxWait) {
          clearInterval(checkInterval)
          doFetch()
        }
      }, 100)
      return
    }

    await doFetch()
  }

  const doFetch = async () => {
    try {
      isLoading.value = true
      const response = await getGlobalSettings()
      if (response.code === 200 && response.data) {
        settings.value = response.data
      }
    } catch (error) {
      console.error('获取全局设置失败:', error)
    } finally {
      isLoading.value = false
    }
  }

  onMounted(() => {
    fetchSettings()
  })

  const settingsRef = computed(() => {
    if (settingsStore.isLoaded && settingsStore.rawSettings) {
      return settingsStore.rawSettings
    }
    return settings.value
  })

  return {
    isLoading,
    settings: settingsRef,
    fetchSettings,
  }
}
