<script setup lang="ts">
  import { onMounted, reactive, watch } from 'vue'
  import { defaultSettings, type Setting } from '@/api/admin/settings'
  import { useTexts } from '@/composables/useTexts'
  import SiteModeSelector from './SiteModeSelector.vue'

  const props = defineProps<{
    settings: Setting[]
  }>()

  const emit = defineEmits<{
    (e: 'update', settings: Setting[]): void
  }>()

  const { $t } = useTexts()

  /* 从defaultSettings获取theme设置的默认值 */
  const themeDefaults =
    defaultSettings.theme?.reduce(
      (acc, setting) => {
        acc[setting.key] = setting.value
        return acc
      },
      {} as Record<string, unknown>
    ) || {}

  /* 本地设置对象 */
  const localSettings = reactive({
    site_mode: themeDefaults.site_mode || 'website',
  })

  /* 将扁平化的设置转换为Setting数组 */
  const getSettingsArray = (): Setting[] => {
    const relevantKeys = ['site_mode']

    if (!defaultSettings.theme) {
      return [
        {
          key: 'site_mode',
          value: localSettings.site_mode,
          type: 'string',
          group: 'theme',
          description: $t('admin.construction.siteMode.label'),
        },
      ]
    }

    return defaultSettings.theme
      .filter((defaultSetting) => relevantKeys.includes(defaultSetting.key))
      .map((defaultSetting) => {
        const settingKey = defaultSetting.key as keyof typeof localSettings
        return {
          key: defaultSetting.key,
          value: localSettings[settingKey],
          type: defaultSetting.type,
          group: 'theme',
          description: defaultSetting.description,
        }
      })
  }

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
  <div class="theme-settings">
    <div class="rounded-xl border border-subtle bg-background-600 p-6">
      <div class="space-y-6">
        <div class="form-group">
          <label class="form-label">{{ $t('admin.construction.siteMode.formLabel') }}</label>
          <div class="form-help mb-4">{{ $t('admin.construction.siteMode.helpText') }}</div>

          <SiteModeSelector v-model="localSettings.site_mode" />
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
  .theme-settings {
    margin-bottom: 2rem;
  }
</style>
