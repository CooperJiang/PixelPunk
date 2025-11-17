import { computed } from 'vue'
import { useSettingsStore } from '@/store/settings'

export function useAppearanceSettings() {
  const settingsStore = useSettingsStore()

  const loadSettings = async () => {
    if (!settingsStore.isLoaded) {
      await settingsStore.loadGlobalSettings()
    }
  }

  const socialLinks = computed(() => ({
    official: {
      show: settingsStore.showOfficialSite,
      url: settingsStore.officialSiteUrl,
    },
    github: {
      show: settingsStore.showGithubLink,
      url: settingsStore.githubUrl,
    },
    wechat: {
      show: settingsStore.showWechatGroup,
      qrImage: settingsStore.wechatQrImageUrl,
      account: settingsStore.wechatContactAccount,
      hasConfig: !!(settingsStore.wechatQrImageUrl || settingsStore.wechatContactAccount),
    },
    qq: {
      show: settingsStore.showQqGroup,
      qrImage: settingsStore.qqQrImageUrl,
      groupNumber: settingsStore.qqGroupNumber,
      hasConfig: !!(settingsStore.qqQrImageUrl || settingsStore.qqGroupNumber),
    },
  }))

  const layoutSettings = computed(() => ({
    multiLayoutEnabled: settingsStore.enableMultiLayout,
    defaultLayout: settingsStore.defaultLayout,
  }))

  const languageSettings = computed(() => ({
    multiLanguageEnabled: settingsStore.enableMultiLanguage,
    defaultLanguage: settingsStore.defaultLanguage,
  }))

  return {
    socialLinks,
    layoutSettings,
    languageSettings,
    isLoading: computed(() => settingsStore.loading),
    isLoaded: computed(() => settingsStore.isLoaded),
    loadSettings,
  }
}
