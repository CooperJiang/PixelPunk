<script setup lang="ts">
  import { onMounted, reactive, watch } from 'vue'
  import { defaultSettings, type Setting } from '@/api/admin/settings'
  import { useTexts } from '@/composables/useTexts'

  const { $t } = useTexts()

  const props = defineProps<{
    settings: Setting[]
  }>()

  const emit = defineEmits<{
    (e: 'update', settings: Setting[]): void
  }>()

  /* 从defaultSettings获取注册设置的默认值 */
  const registrationDefaults = defaultSettings.registration.reduce(
    (acc, setting) => {
      acc[setting.key] = setting.value
      return acc
    },
    {} as Record<string, unknown>
  )

  /* 本地设置对象（扁平化）- 使用导入的默认值 */
  const localSettings = reactive({
    user_initial_storage: registrationDefaults.user_initial_storage || 100,
    enable_registration: registrationDefaults.enable_registration ?? true,
    email_verification: registrationDefaults.email_verification ?? true,
    user_initial_bandwidth: registrationDefaults.user_initial_bandwidth || 1024,
  })

  /* 将扁平化的设置转换为Setting数组 - 使用defaultSettings中的类型和描述信息 */
  const getSettingsArray = (): Setting[] =>
    defaultSettings.registration.map((defaultSetting) => {
      const settingKey = defaultSetting.key as keyof typeof localSettings
      return {
        key: defaultSetting.key,
        value: localSettings[settingKey],
        type: defaultSetting.type,
        group: 'registration',
        description: defaultSetting.description,
      }
    })

  const applySettings = (settings: Setting[]) => {
    settings.forEach((setting) => {
      const key = setting.key as keyof typeof localSettings

      if (key in localSettings) {
        const settingsRecord = localSettings as Record<string, unknown>
        settingsRecord[key] = setting.value
      }
    })
  }

  watch(
    () => props.settings,
    (newSettings) => {
      if (newSettings && newSettings.length > 0) {
        applySettings(newSettings)
      }
    },
    { deep: true, immediate: true }
  )

  watch(
    localSettings,
    () => {
      emit('update', getSettingsArray())
    },
    { deep: true }
  )

  onMounted(() => {
    if (props.settings && props.settings.length > 0) {
      applySettings(props.settings)
    }
  })
</script>

<template>
  <div class="space-y-5">
    <div class="space-y-5">
      <div class="space-y-5">
        <div class="flex flex-col md:flex-row md:items-center">
          <label class="flex w-full items-center text-sm text-content md:w-48 md:py-2">
            <i class="fas fa-door-open mr-2 text-brand-600" />{{ $t('admin.settings.registration.enableRegistration.label') }}
          </label>
          <div class="flex flex-1 items-center">
            <CyberSwitch v-model="localSettings.enable_registration" />
            <span class="ml-3 text-sm text-content-muted">{{
              localSettings.enable_registration
                ? $t('admin.settings.registration.enableRegistration.enabled')
                : $t('admin.settings.registration.enableRegistration.disabled')
            }}</span>
          </div>
        </div>

        <div class="flex flex-col md:flex-row md:items-center">
          <label class="flex w-full items-center text-sm text-content md:w-48 md:py-2">
            <i class="fas fa-check-circle mr-2 text-brand-600" />{{ $t('admin.settings.registration.emailVerification.label') }}
          </label>
          <div class="flex flex-1 items-center">
            <CyberSwitch v-model="localSettings.email_verification" />
            <span class="ml-3 text-sm text-content-muted">{{
              localSettings.email_verification
                ? $t('admin.settings.registration.emailVerification.enabled')
                : $t('admin.settings.registration.emailVerification.disabled')
            }}</span>
          </div>
        </div>

        <div class="flex flex-col md:flex-row md:items-center">
          <label class="flex w-full items-center text-sm text-content md:w-48 md:py-2">
            <i class="fas fa-hdd mr-2 text-brand-600" />{{ $t('admin.settings.registration.initialStorage.label') }}
          </label>
          <div class="flex-1">
            <div class="flex items-center">
              <CyberInput
                v-model="localSettings.user_initial_storage"
                type="number"
                :placeholder="$t('admin.settings.registration.initialStorage.placeholder')"
                width="500px"
              >
                <template #unit>MB</template>
              </CyberInput>
            </div>
            <p class="text-content-subtle mt-1 text-xs">{{ $t('admin.settings.registration.initialStorage.hint') }}</p>
          </div>
        </div>

        <div class="flex flex-col md:flex-row md:items-center">
          <label class="flex w-full items-center text-sm text-content md:w-48 md:py-2">
            <i class="fas fa-tachometer-alt mr-2 text-brand-600" />{{ $t('admin.settings.registration.initialBandwidth.label') }}
          </label>
          <div class="flex-1">
            <div class="flex items-center">
              <CyberInput
                v-model="localSettings.user_initial_bandwidth"
                type="number"
                :placeholder="$t('admin.settings.registration.initialBandwidth.placeholder')"
                width="500px"
              >
                <template #unit>MB</template>
              </CyberInput>
            </div>
            <p class="text-content-subtle mt-1 text-xs">{{ $t('admin.settings.registration.initialBandwidth.hint') }}</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped lang="scss">
  .hover:shadow-glow:hover {
    box-shadow: 0 0 15px rgba(var(--color-brand-500-rgb), 0.15);
  }
</style>
