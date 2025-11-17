import { ref, computed } from 'vue'
import type { AnalyticsSettings } from './types'

/**
 * 埋点统计设置子模块
 * 管理网站埋点统计配置（百度统计、Google Analytics等）
 */
export function useAnalyticsSettingsModule() {
  const settings = ref<AnalyticsSettings>({})

  const baiduEnabled = computed(() => settings.value.baidu_analytics_enabled === true || settings.value.baidu_analytics_enabled === 'true')

  const baiduSiteId = computed(() => settings.value.baidu_analytics_site_id || '')

  const googleEnabled = computed(() => settings.value.google_analytics_enabled === true || settings.value.google_analytics_enabled === 'true')

  const googleMeasurementId = computed(() => settings.value.google_analytics_measurement_id || '')

  const updateSettings = (newSettings: AnalyticsSettings) => {
    settings.value = { ...settings.value, ...newSettings }
  }

  const reset = () => {
    settings.value = {}
  }

  return {
    analyticsSettings: settings,

    baiduEnabled,
    baiduSiteId,
    googleEnabled,
    googleMeasurementId,

    updateAnalyticsSettings: updateSettings,
    resetAnalyticsSettings: reset,
  }
}
