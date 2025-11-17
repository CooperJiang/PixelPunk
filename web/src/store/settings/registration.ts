import { ref, computed } from 'vue'
import type { RegistrationSettings } from './types'

/**
 * 注册设置子模块
 * 管理用户注册相关的配置
 */
export function useRegistrationSettingsModule() {
  const settings = ref<RegistrationSettings>({})

  const isRegistrationEnabled = computed(() => settings.value.enable_registration ?? true)

  const isEmailVerificationEnabled = computed(() => settings.value.email_verification ?? true)

  const updateSettings = (newSettings: RegistrationSettings) => {
    settings.value = { ...settings.value, ...newSettings }
  }

  const reset = () => {
    settings.value = {}
  }

  return {
    registrationSettings: settings,

    isRegistrationEnabled,
    isEmailVerificationEnabled,

    updateRegistrationSettings: updateSettings,
    resetRegistrationSettings: reset,
  }
}
