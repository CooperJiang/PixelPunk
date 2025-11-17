<script setup lang="ts">
  import { ref, watch, onMounted } from 'vue'
  import BaiduAnalyticsSettings from './BaiduAnalyticsSettings.vue'
  import GoogleAnalyticsSettings from './GoogleAnalyticsSettings.vue'
  import type { Setting } from '@/api/admin/settings'
  import { useTexts } from '@/composables/useTexts'

  const { $t } = useTexts()

  const props = defineProps<{
    settings: Setting[]
  }>()

  const emit = defineEmits<{
    (e: 'update', settings: Setting[]): void
  }>()

  const activeTab = ref('baidu')
  const baiduSettings = ref<Setting[]>([])
  const googleSettings = ref<Setting[]>([])

  const filterSettingsByProvider = (settings: Setting[], provider: string) => {
    return settings.filter((s) => s.key.startsWith(`${provider}_analytics_`))
  }

  const applySettings = (settings: Setting[]) => {
    baiduSettings.value = filterSettingsByProvider(settings, 'baidu')
    googleSettings.value = filterSettingsByProvider(settings, 'google')
  }

  const handleBaiduUpdate = (settings: Setting[]) => {
    baiduSettings.value = settings
    emitAllSettings()
  }

  const handleGoogleUpdate = (settings: Setting[]) => {
    googleSettings.value = settings
    emitAllSettings()
  }

  const emitAllSettings = () => {
    const allSettings = [...baiduSettings.value, ...googleSettings.value]
    emit('update', allSettings)
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

  onMounted(() => {
    if (props.settings && props.settings.length > 0) {
      applySettings(props.settings)
    }
  })
</script>

<template>
  <div class="analytics-settings">
    <CyberTabs v-model="activeTab">
      <CyberTabPane name="baidu" :label="$t('admin.settings.analytics.baidu.title')">
        <BaiduAnalyticsSettings :settings="baiduSettings" @update="handleBaiduUpdate" />
      </CyberTabPane>
      <CyberTabPane name="google" :label="$t('admin.settings.analytics.google.title')">
        <GoogleAnalyticsSettings :settings="googleSettings" @update="handleGoogleUpdate" />
      </CyberTabPane>
    </CyberTabs>
  </div>
</template>

<style scoped>
  .analytics-settings {
    width: 100%;
  }

  .analytics-settings :deep(.cyber-tabs) {
    gap: 16px;
  }

  .analytics-settings :deep(.tab-item) {
    padding: 8px 20px;
  }

  .analytics-settings :deep(.tab-label) {
    font-size: 13px;
    font-weight: 500;
  }

  .analytics-settings :deep(.tabs-header) {
    padding-bottom: 6px;
  }

  .analytics-settings :deep(.space-y-5) {
    row-gap: 16px;
  }
</style>
