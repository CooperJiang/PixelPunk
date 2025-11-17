<script setup lang="ts">
  import { reactive, watch, onMounted } from 'vue'
  import type { Setting } from '@/api/admin/settings'
  import { useToast } from '@/components/Toast/useToast'
  import { useTexts } from '@/composables/useTexts'
  import { checkGoogleAnalytics } from '@/utils/system/google-analytics'

  const { $t } = useTexts()

  const props = defineProps<{
    settings: Setting[]
  }>()

  const emit = defineEmits<{
    (e: 'update', settings: Setting[]): void
  }>()

  const toast = useToast()

  const localSettings = reactive({
    google_analytics_enabled: false,
    google_analytics_measurement_id: '',
  })

  const getSettingsArray = (): Setting[] => {
    return [
      {
        key: 'google_analytics_enabled',
        value: localSettings.google_analytics_enabled,
        type: 'boolean',
        group: 'analytics',
        description: $t('admin.settings.analytics.google.desc.enabled'),
      },
      {
        key: 'google_analytics_measurement_id',
        value: localSettings.google_analytics_measurement_id,
        type: 'string',
        group: 'analytics',
        description: $t('admin.settings.analytics.google.desc.measurementId'),
      },
    ]
  }

  const applySettings = (settings: Setting[]) => {
    settings.forEach((setting) => {
      const key = setting.key as keyof typeof localSettings
      if (key in localSettings) {
        const settingsRecord = localSettings as Record<string, any>
        if (setting.type === 'boolean') {
          settingsRecord[key] = setting.value === true || setting.value === 'true'
        } else {
          settingsRecord[key] = setting.value
        }
      }
    })
  }

  const checkStatus = () => {
    const status = checkGoogleAnalytics()

    if (status.scriptLoaded && status.initialized) {
      toast.success($t('admin.settings.analytics.google.status.success', { measurementId: status.measurementId || 'unknown' }))
    } else if (status.initialized && !status.scriptLoaded) {
      toast.warning($t('admin.settings.analytics.google.status.loading'))
    } else {
      toast.error($t('admin.settings.analytics.google.status.notInitialized'))
    }
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
  <div class="space-y-4 py-3">
    <div class="flex flex-col md:flex-row md:items-center">
      <label class="flex w-full items-center text-sm text-content md:w-48 md:py-2">
        <i class="fas fa-chart-line mr-2 text-brand-500" />{{ $t('admin.settings.analytics.google.enabled.label') }}
      </label>
      <div class="flex-1">
        <CyberSwitch v-model="localSettings.google_analytics_enabled" />
      </div>
    </div>

    <div class="flex flex-col md:flex-row md:items-start">
      <label class="flex w-full items-center text-sm text-content md:w-48 md:py-2">
        <i class="fas fa-fingerprint mr-2 text-brand-500" />{{ $t('admin.settings.analytics.google.measurementId.label') }}
      </label>
      <div class="flex-1">
        <div class="flex gap-2">
          <CyberInput
            v-model="localSettings.google_analytics_measurement_id"
            :placeholder="$t('admin.settings.analytics.google.measurementId.placeholder')"
            width="500px"
          />
          <CyberButton size="small" variant="secondary" @click="checkStatus">
            <i class="fas fa-check-circle mr-1.5" />
            {{ $t('admin.settings.analytics.google.checkStatus') }}
          </CyberButton>
        </div>
        <div class="text-content-subtle mt-1 text-xs">{{ $t('admin.settings.analytics.google.measurementId.description') }}</div>
      </div>
    </div>

    <div class="border-brand-500/20 bg-background-800/50 mt-4 rounded-lg border p-3">
      <div class="mb-2 flex items-center gap-2 text-xs font-medium text-brand-400">
        <i class="fas fa-info-circle text-xs" />
        <span>{{ $t('admin.settings.analytics.google.setup.title') }}</span>
      </div>
      <ol class="text-content-subtle space-y-1.5 text-xs">
        <li class="flex gap-2">
          <span class="text-brand-500">1.</span>
          <span
            >{{ $t('admin.settings.analytics.google.setup.step1') }}
            <a href="https://analytics.google.com" target="_blank" class="text-brand-400 hover:text-brand-300">{{
              $t('admin.settings.analytics.google.setup.step1Link')
            }}</a></span
          >
        </li>
        <li class="flex gap-2">
          <span class="text-brand-500">2.</span>
          <span>{{ $t('admin.settings.analytics.google.setup.step2') }}</span>
        </li>
        <li class="flex gap-2">
          <span class="text-brand-500">3.</span>
          <div class="flex-1">
            <div>{{ $t('admin.settings.analytics.google.setup.step3') }}</div>
            <div class="mt-1 rounded bg-background-900 px-2 py-1 font-mono text-brand-300">
              <span class="text-error-400">G-XXXXXXXXXX</span>
            </div>
            <div class="mt-1 text-[11px] text-content-muted">
              {{ $t('admin.settings.analytics.google.setup.step3Hint') }}
            </div>
          </div>
        </li>
        <li class="flex gap-2">
          <span class="text-brand-500">4.</span>
          <span>{{ $t('admin.settings.analytics.google.setup.step4') }}</span>
        </li>
        <li class="flex gap-2">
          <span class="text-brand-500">5.</span>
          <span>{{ $t('admin.settings.analytics.google.setup.step5') }}</span>
        </li>
      </ol>
    </div>
  </div>
</template>

<style scoped>
  a {
    text-decoration: underline;
    text-underline-offset: 2px;
    transition: color 0.2s ease;
  }
</style>
