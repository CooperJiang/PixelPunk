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

  /* 从defaultSettings获取网站设置的默认值 */
  const websiteDefaults = defaultSettings.website.reduce(
    (acc, setting) => {
      acc[setting.key] = setting.value
      return acc
    },
    {} as Record<string, unknown>
  )

  /* 本地设置对象（扁平化）- 使用导入的默认值 */
  const localSettings = reactive({
    admin_email: websiteDefaults.admin_email || '',
    site_base_url: websiteDefaults.site_base_url || '',
  })

  /* 将扁平化的设置转换为Setting数组 - 使用defaultSettings中的类型和描述信息 */
  const getSettingsArray = (): Setting[] =>
    defaultSettings.website.map((defaultSetting) => {
      const settingKey = defaultSetting.key as keyof typeof localSettings
      return {
        key: defaultSetting.key,
        value: localSettings[settingKey],
        type: defaultSetting.type,
        group: 'website',
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
        <div class="flex flex-col md:flex-row md:items-start">
          <label class="flex w-full items-center text-sm text-content md:w-48 md:py-2">
            <i class="fas fa-globe mr-2 text-brand-600" />{{ $t('admin.settings.website.siteUrl.label') }}
          </label>
          <div class="flex-1">
            <CyberInput
              v-model="localSettings.site_base_url"
              :placeholder="$t('admin.settings.website.siteUrl.placeholder')"
              width="500px"
            />
            <div class="text-content-subtle mt-1 text-xs">{{ $t('admin.settings.website.siteUrl.description') }}</div>
          </div>
        </div>

        <div class="flex flex-col md:flex-row md:items-center">
          <label class="flex w-full items-center text-sm text-content md:w-48 md:py-2">
            <i class="fas fa-envelope mr-2 text-brand-600" />{{ $t('admin.settings.website.contactEmail.label') }}
          </label>
          <div class="flex-1">
            <CyberInput
              v-model="localSettings.admin_email"
              :placeholder="$t('admin.settings.website.contactEmail.placeholder')"
              width="500px"
            />
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
