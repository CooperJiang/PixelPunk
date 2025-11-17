import { ref, computed } from 'vue'
import type { VersionSettings } from './types'

/**
 * 版本设置子模块
 * 管理系统版本信息
 */
export function useVersionSettingsModule() {
  const settings = ref<VersionSettings>({})

  const currentVersion = computed(() => settings.value.current_version || '1.0.0')

  const buildTime = computed(() => settings.value.build_time || '')

  const isUpdateAvailable = computed(() => settings.value.update_available ?? false)

  const lastUpdateCheck = computed(() => settings.value.last_update_check || '')

  const updateSettings = (newSettings: VersionSettings) => {
    settings.value = { ...settings.value, ...newSettings }
  }

  const reset = () => {
    settings.value = {}
  }

  return {
    versionSettings: settings,

    currentVersion,
    buildTime,
    isUpdateAvailable,
    lastUpdateCheck,

    updateVersionSettings: updateSettings,
    resetVersionSettings: reset,
  }
}
