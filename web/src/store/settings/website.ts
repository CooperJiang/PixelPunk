import { ref, computed } from 'vue'
import type { WebsiteSettings } from './types'

/**
 * 网站基础设置子模块
 * 管理网站基础配置（管理员邮箱、基础URL等）
 */
export function useWebsiteSettingsModule() {
  const settings = ref<WebsiteSettings>({})

  const adminEmail = computed(() => settings.value.admin_email || '')

  const siteBaseUrl = computed(() => settings.value.site_base_url || '')

  const updateSettings = (newSettings: WebsiteSettings) => {
    settings.value = { ...settings.value, ...newSettings }
  }

  const reset = () => {
    settings.value = {}
  }

  return {
    websiteSettings: settings,

    adminEmail,
    siteBaseUrl,

    updateWebsiteSettings: updateSettings,
    resetWebsiteSettings: reset,
  }
}
