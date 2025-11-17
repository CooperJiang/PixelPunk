import { ref, computed } from 'vue'
import type { WebsiteInfoSettings } from './types'
import { useTexts } from '@/composables/useTexts'

/**
 * 网站信息设置子模块
 * 管理网站基本信息、SEO、品牌等配置
 */
export function useWebsiteInfoSettingsModule() {
  const { $t } = useTexts()
  const settings = ref<WebsiteInfoSettings>({})

  const siteName = computed(() => settings.value.site_name || 'PixelPunk')

  const siteDescription = computed(() => settings.value.site_description || $t('store.settings.defaults.siteDescription'))

  const siteKeywords = computed(() => settings.value.site_keywords || '')

  const icpNumber = computed(() => settings.value.icp_number || '')

  const siteLogoUrl = computed(() => settings.value.site_logo_url || '/logo.png')

  const faviconUrl = computed(() => settings.value.favicon_url || '/favicon.ico')

  const copyrightText = computed(() => settings.value.copyright_text || '')

  const contactEmail = computed(() => settings.value.contact_email || '')

  const footerCustomText = computed(() => settings.value.footer_custom_text || '')

  const siteHeroTitle = computed(() => settings.value.site_hero_title || '')

  const siteFeaturesText = computed(() => settings.value.site_features_text || '')

  const siteSlogan = computed(() => settings.value.site_slogan || '')

  const showImageCount = computed(() => settings.value.show_file_count ?? false)

  const showStorageUsage = computed(() => settings.value.show_storage_usage ?? false)

  const updateSettings = (newSettings: WebsiteInfoSettings) => {
    settings.value = { ...settings.value, ...newSettings }
  }

  const reset = () => {
    settings.value = {}
  }

  return {
    websiteInfoSettings: settings,

    siteName,
    siteDescription,
    siteKeywords,
    icpNumber,
    siteLogoUrl,
    faviconUrl,
    copyrightText,
    contactEmail,
    footerCustomText,
    siteHeroTitle,
    siteFeaturesText,
    siteSlogan,
    showImageCount,
    showStorageUsage,

    updateWebsiteInfoSettings: updateSettings,
    resetWebsiteInfoSettings: reset,
  }
}
