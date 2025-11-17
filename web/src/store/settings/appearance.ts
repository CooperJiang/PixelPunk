import { ref, computed } from 'vue'
import type { AppearanceSettings } from './types'

/**
 * 外观设置子模块
 * 管理外观相关的配置（社交链接、布局、语言等）
 */
export function useAppearanceSettingsModule() {
  const settings = ref<AppearanceSettings>({})

  const showOfficialSite = computed(() => settings.value.show_official_site ?? false)

  const officialSiteUrl = computed(() => settings.value.official_site_url || '')

  const showGithubLink = computed(() => settings.value.show_github_link ?? false)

  const githubUrl = computed(() => settings.value.github_url || '')

  const showWechatGroup = computed(() => settings.value.show_wechat_group ?? false)

  const wechatQrImageUrl = computed(() => settings.value.wechat_qr_image_url || '')

  const wechatContactAccount = computed(() => settings.value.wechat_contact_account || '')

  const showQqGroup = computed(() => settings.value.show_qq_group ?? false)

  const qqQrImageUrl = computed(() => settings.value.qq_qr_image_url || '')

  const qqGroupNumber = computed(() => settings.value.qq_group_number || '')

  const enableMultiLayout = computed(() => settings.value.enable_multi_layout ?? false)

  const defaultLayout = computed(() => settings.value.default_layout || 'default')

  const enableMultiLanguage = computed(() => settings.value.enable_multi_language ?? false)

  const defaultLanguage = computed(() => settings.value.default_language || 'zh-CN')

  const defaultTheme = computed(() => settings.value.default_theme || '')

  const updateSettings = (newSettings: AppearanceSettings) => {
    settings.value = { ...settings.value, ...newSettings }
  }

  const reset = () => {
    settings.value = {}
  }

  return {
    appearanceSettings: settings,

    showOfficialSite,
    officialSiteUrl,
    showGithubLink,
    githubUrl,
    showWechatGroup,
    wechatQrImageUrl,
    wechatContactAccount,
    showQqGroup,
    qqQrImageUrl,
    qqGroupNumber,
    enableMultiLayout,
    defaultLayout,
    enableMultiLanguage,
    defaultLanguage,
    defaultTheme,

    updateAppearanceSettings: updateSettings,
    resetAppearanceSettings: reset,
  }
}
