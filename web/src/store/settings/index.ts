import { defineStore } from 'pinia'
import { ref } from 'vue'
import { getGlobalSettings, type GlobalSettingsResponse } from '@/api/admin/settings'
import FaviconManager from '@/utils/favicon'
import { SEOManager } from '@/utils/seo'
import { useTexts } from '@/composables/useTexts'

/* 导入所有子模块 */
import { useWebsiteSettingsModule } from './website'
import { useWebsiteInfoSettingsModule } from './websiteInfo'
import { useUploadSettingsModule } from './upload'
import { useRegistrationSettingsModule } from './registration'
import { useVersionSettingsModule } from './version'
import { useAppearanceSettingsModule } from './appearance'
import { useAISettingsModule } from './ai'
import { useAnalyticsSettingsModule } from './analytics'

export * from './types'

/**
 * 🎛️ 统一设置管理 Store
 * 采用模块化架构，将 450+ 行的大 Store 拆分为多个子模块
 * 每个子模块负责特定领域的设置，便于维护和扩展
 */
export const useSettingsStore = defineStore('settings', () => {
  const { $t } = useTexts()

  const website = useWebsiteSettingsModule()
  const websiteInfo = useWebsiteInfoSettingsModule()
  const upload = useUploadSettingsModule()
  const registration = useRegistrationSettingsModule()
  const version = useVersionSettingsModule()
  const appearance = useAppearanceSettingsModule()
  const ai = useAISettingsModule()
  const analytics = useAnalyticsSettingsModule()

  const isLoaded = ref(false)
  const loading = ref(false)
  const rawSettings = ref<GlobalSettingsResponse | null>(null) // 保存原始响应数据

  async function loadGlobalSettings() {
    if (loading.value) return

    loading.value = true
    try {
      const response = await getGlobalSettings()
      if (response.code === 200 && response.data) {
        const data = response.data as GlobalSettingsResponse

        rawSettings.value = data

        if (data.website) {
          website.updateWebsiteSettings(data.website)
        }
        if (data.website_info) {
          websiteInfo.updateWebsiteInfoSettings(data.website_info)
        }
        if (data.upload) {
          upload.updateUploadSettings(data.upload)
        }
        if (data.registration) {
          registration.updateRegistrationSettings(data.registration)
        }
        if (data.version) {
          version.updateVersionSettings(data.version)
        }
        if (data.appearance) {
          appearance.updateAppearanceSettings(data.appearance)
        }
        if (data.ai) {
          ai.updateAISettings(data.ai)
        }
        if (data.vector) {
          ai.updateVectorSettings(data.vector)
        }
        if (data.analytics) {
          analytics.updateAnalyticsSettings(data.analytics)
        }

        if (data.website_info?.favicon_url) {
          FaviconManager.update(data.website_info.favicon_url)
        }

        SEOManager.setSEO({
          siteName: data.website_info?.site_name || 'PixelPunk',
          description: data.website_info?.site_description || $t('store.settings.defaults.siteDescription'),
          keywords: data.website_info?.site_keywords || '',
        })

        isLoaded.value = true
      }
    } catch (error) {
      console.error('加载全局设置失败:', error)
    } finally {
      loading.value = false
    }
  }

  function reset() {
    website.resetWebsiteSettings()
    websiteInfo.resetWebsiteInfoSettings()
    upload.resetUploadSettings()
    registration.resetRegistrationSettings()
    version.resetVersionSettings()
    appearance.resetAppearanceSettings()
    ai.resetAISettings()
    analytics.resetAnalyticsSettings()

    isLoaded.value = false
    loading.value = false
    rawSettings.value = null
  }

  return {
    isLoaded,
    loading,
    rawSettings, // 导出原始响应数据

    ...website,
    ...websiteInfo,
    ...upload,
    ...registration,
    ...version,
    ...appearance,
    ...ai,
    ...analytics,

    loadGlobalSettings,
    reset,
  }
})
