<script setup lang="ts">
  import { computed, onMounted, reactive, watch } from 'vue'
  import { type Setting, defaultSettings } from '@/api/admin/settings'
  import { useSettingsStore } from '@/store/settings'
  import { useTexts } from '@/composables/useTexts'

  const { $t } = useTexts()

  const props = defineProps<{
    settings: Setting[]
    aiSettings?: Setting[] // AI设置，用于检查配置状态
  }>()

  const emit = defineEmits<{
    (e: 'update', settings: Setting[]): void
  }>()

  const settingsStore = useSettingsStore()

  const getFieldName = (key: string): string => {
    return $t(`admin.ai.contentDetection.fieldNames.${key}`) || key
  }

  const uploadDefaults = defaultSettings.upload.reduce(
    (acc, setting) => {
      acc[setting.key] = setting.value
      return acc
    },
    {} as Record<string, unknown>
  )

  const aiDefaults = defaultSettings.ai.reduce(
    (acc, setting) => {
      acc[setting.key] = setting.value
      return acc
    },
    {} as Record<string, unknown>
  )

  const localSettings = reactive({
    content_detection_enabled: uploadDefaults.content_detection_enabled ?? true,
    sensitive_content_handling: uploadDefaults.sensitive_content_handling || 'mark_only',
    ai_analysis_enabled: uploadDefaults.ai_analysis_enabled ?? true,
    nsfw_threshold: aiDefaults.nsfw_threshold ?? 0.6,
  })

  const sensitiveContentOptions = computed(() => [
    {
      label: $t('admin.ai.contentDetection.sensitiveOptions.markOnly'),
      value: 'mark_only',
      description: $t('admin.ai.contentDetection.sensitiveOptions.markOnlyDesc'),
    },
    {
      label: $t('admin.ai.contentDetection.sensitiveOptions.autoDelete'),
      value: 'auto_delete',
      description: $t('admin.ai.contentDetection.sensitiveOptions.autoDeleteDesc'),
    },
    {
      label: $t('admin.ai.contentDetection.sensitiveOptions.pendingReview'),
      value: 'pending_review',
      description: $t('admin.ai.contentDetection.sensitiveOptions.pendingReviewDesc'),
    },
  ])

  const getThresholdLabel = (threshold: number) => {
    if (threshold <= 0.3) {
      return { text: $t('admin.ai.contentDetection.thresholds.lenient'), color: 'text-green-400' }
    } else if (threshold <= 0.5) {
      return { text: $t('admin.ai.contentDetection.thresholds.moderate'), color: 'text-yellow-400' }
    } else if (threshold <= 0.7) {
      return { text: $t('admin.ai.contentDetection.thresholds.strict'), color: 'text-orange-400' }
    }
    return { text: $t('admin.ai.contentDetection.thresholds.veryStrict'), color: 'text-red-400' }
  }

  const aiConfigStatus = computed(() => {
    const { isAIEnabled } = settingsStore

    if (!isAIEnabled) {
      return {
        isConfigured: false,
        missingFields: ['ai_enabled'],
        hasApiKey: false,
        hasModel: false,
      }
    }

    if (!props.aiSettings || props.aiSettings.length === 0) {
      return {
        isConfigured: true, // AI已启用就认为配置完成
        missingFields: [],
        hasApiKey: true, // 假设已配置（因为AI已启用）
        hasModel: true, // 假设已配置（因为AI已启用）
      }
    }

    const settingsMap = props.aiSettings.reduce(
      (acc, setting) => {
        acc[setting.key] = setting.value
        return acc
      },
      {} as Record<string, unknown>
    )

    const missingFields: string[] = []

    const hasApiKey = Boolean(settingsMap.ai_api_key && settingsMap.ai_api_key.toString().trim())
    if (!hasApiKey) {
      missingFields.push('ai_api_key')
    }

    const hasModel = Boolean(settingsMap.ai_model && settingsMap.ai_model.toString().trim())
    if (!hasModel) {
      missingFields.push('ai_model')
    }

    const isConfigured = missingFields.length === 0

    return {
      isConfigured,
      missingFields,
      hasApiKey,
      hasModel,
    }
  })

  const getSettingsArray = () =>
    Object.entries(localSettings).map(([key, value]) => ({
      key,
      value,
      type: typeof value === 'number' ? 'number' : typeof value === 'boolean' ? 'boolean' : 'string',
      group: key === 'nsfw_threshold' ? 'ai' : 'upload', // NSFW阈值属于AI组，其他属于upload组
      description: getSettingDescription(key),
      is_system: true,
    }))

  const getSettingDescription = (key: string): string => {
    const descriptions: Record<string, string> = {
      content_detection_enabled: $t('admin.ai.contentDetection.desc.enabled'),
      sensitive_content_handling: $t('admin.ai.contentDetection.desc.handling'),
      ai_analysis_enabled: $t('admin.ai.contentDetection.desc.aiAnalysis'),
      nsfw_threshold: $t('admin.ai.contentDetection.desc.nsfwThreshold'),
    }
    return descriptions[key] || ''
  }

  const initializeSettings = () => {
    props.settings.forEach((setting) => {
      if (setting.key in localSettings) {
        ;(localSettings as Record<string, unknown>)[setting.key] = setting.value
      }
    })

    if (props.aiSettings) {
      props.aiSettings.forEach((setting) => {
        if (setting.key in localSettings) {
          ;(localSettings as Record<string, unknown>)[setting.key] = setting.value
        }
      })
    }
  }

  watch(
    localSettings,
    () => {
      emit('update', getSettingsArray())
    },
    { deep: true }
  )

  watch(
    () => props.settings,
    (newSettings) => {
      if (newSettings && newSettings.length > 0) {
        initializeSettings()
      }
    },
    { immediate: true }
  )

  watch(
    () => props.aiSettings,
    (newAISettings) => {
      if (newAISettings && newAISettings.length > 0) {
        initializeSettings()
      }
    },
    { immediate: true }
  )

  onMounted(() => {
    initializeSettings()
  })
</script>

<template>
  <div class="content-detection-container">
    <div class="settings-list">
      <div class="settings-group">
        <div v-if="!aiConfigStatus.isConfigured" class="warning-panel">
          <div class="flex items-start">
            <i class="warning-icon fas fa-exclamation-triangle" />
            <div>
              <h4 class="warning-title">{{ $t('admin.ai.contentDetection.aiNotConfigured') }}</h4>
              <p class="warning-text">{{ $t('admin.ai.contentDetection.aiNotConfiguredDesc') }}</p>
              <div class="missing-fields">
                <span v-for="field in aiConfigStatus.missingFields" :key="field" class="missing-field-badge">
                  {{ getFieldName(field) }}
                </span>
              </div>
              <router-link to="/admin/ai?tab=ai" class="config-link">
                <i class="fas fa-arrow-right mr-2" />
                {{ $t('admin.ai.contentDetection.goToAIConfig') }}
              </router-link>
            </div>
          </div>
        </div>

        <div class="flex flex-col md:flex-row md:items-center">
          <label class="flex w-full items-center text-sm text-content md:w-48 md:py-2">
            <i class="fas fa-toggle-on mr-2 text-brand-700" />{{ $t('admin.ai.contentDetection.enableDetection') }}
          </label>
          <div class="flex-1">
            <div class="flex items-center">
              <CyberSwitch v-model="localSettings.content_detection_enabled" :disabled="!aiConfigStatus.isConfigured" />
              <span class="ml-3 text-sm text-content-muted">{{
                localSettings.content_detection_enabled
                  ? $t('admin.ai.contentDetection.detectionEnabled')
                  : $t('admin.ai.contentDetection.detectionDisabled')
              }}</span>
            </div>
            <p class="setting-description">{{ $t('admin.ai.contentDetection.detectionDesc') }}</p>
          </div>
        </div>

        <div class="flex flex-col md:flex-row md:items-center">
          <label class="flex w-full items-center text-sm text-content md:w-48 md:py-2">
            <i class="fas fa-exclamation-triangle mr-2 text-brand-700" />{{ $t('admin.ai.contentDetection.sensitiveHandling') }}
          </label>
          <div class="flex-1">
            <CyberDropdown
              v-model="localSettings.sensitive_content_handling"
              :options="sensitiveContentOptions"
              :placeholder="$t('admin.ai.contentDetection.sensitiveHandlingPlaceholder')"
              width="500px"
              :disabled="!aiConfigStatus.isConfigured || !localSettings.content_detection_enabled"
            />
            <p class="setting-description">{{ $t('admin.ai.contentDetection.sensitiveHandlingDesc') }}</p>
          </div>
        </div>

        <div class="flex flex-col md:flex-row md:items-center">
          <label class="flex w-full items-center text-sm text-content md:w-48 md:py-2">
            <i class="fas fa-brain mr-2 text-brand-700" />{{ $t('admin.ai.contentDetection.enableAIAnalysis') }}
          </label>
          <div class="flex-1">
            <div class="flex items-center">
              <CyberSwitch v-model="localSettings.ai_analysis_enabled" :disabled="!aiConfigStatus.isConfigured" />
              <span class="ml-3 text-sm text-content-muted">{{
                localSettings.ai_analysis_enabled
                  ? $t('admin.ai.contentDetection.aiAnalysisEnabled')
                  : $t('admin.ai.contentDetection.aiAnalysisDisabled')
              }}</span>
            </div>
            <p class="setting-description">{{ $t('admin.ai.contentDetection.aiAnalysisDesc') }}</p>
          </div>
        </div>

        <div class="flex flex-col md:flex-row md:items-center">
          <label class="flex w-full items-center text-sm text-content md:w-48 md:py-2">
            <i class="fas fa-eye-slash mr-2 text-brand-700" />{{ $t('admin.ai.contentDetection.nsfwThreshold') }}
          </label>
          <div class="flex-1">
            <div class="slider-container flex items-center">
              <CyberSlider
                v-model="localSettings.nsfw_threshold"
                :min="0"
                :max="0.95"
                :step="0.01"
                :disabled="!aiConfigStatus.isConfigured || !localSettings.content_detection_enabled"
                :description="$t('admin.ai.contentDetection.nsfwThresholdDesc')"
                width="365px"
              />
              <span class="ml-3 flex h-full items-center text-sm text-content-muted">
                {{ localSettings.nsfw_threshold.toFixed(2) }}
              </span>
              <div class="ml-3 flex items-center text-xs text-content-muted">
                <span class="mr-2">{{ $t('admin.ai.contentDetection.strictnessLevel') }}:</span>
                <span :class="getThresholdLabel(localSettings.nsfw_threshold).color">
                  {{ getThresholdLabel(localSettings.nsfw_threshold).text }}
                </span>
              </div>
            </div>
          </div>
        </div>

        <ul class="feature-list">
          <li>• {{ $t('admin.ai.contentDetection.features.autoDetect') }}</li>
          <li>• {{ $t('admin.ai.contentDetection.features.multiStrategy') }}</li>
          <li>• {{ $t('admin.ai.contentDetection.features.adjustableThreshold') }}</li>
          <li>• {{ $t('admin.ai.contentDetection.features.safeguard') }}</li>
        </ul>
      </div>
    </div>
  </div>
</template>

<style scoped lang="scss">
  .content-detection-container {
    display: flex;
    flex-direction: column;
    gap: var(--space-lg);
  }

  .settings-list {
    display: flex;
    flex-direction: column;
    gap: var(--space-lg);
  }

  .settings-group {
    display: flex;
    flex-direction: column;
    gap: var(--space-lg);
  }

  .warning-panel {
    border-radius: var(--radius-lg);
    border: 1px solid rgba(var(--color-warning-rgb), 0.3);
    background: rgba(var(--color-warning-rgb), 0.1);
    padding: var(--space-md);
  }

  .warning-icon {
    margin-right: var(--space-sm);
    margin-top: var(--space-xs);
    color: var(--color-warning-400);
  }

  .warning-title {
    margin-bottom: var(--space-xs);
    font-weight: var(--font-medium);
    color: var(--color-warning-300);
  }

  .warning-text {
    margin-bottom: var(--space-sm);
    font-size: var(--text-sm);
    color: var(--color-content-muted);
  }

  .missing-fields {
    margin-bottom: var(--space-sm);
    display: flex;
    flex-wrap: wrap;
    gap: var(--space-xs);
  }

  .missing-field-badge {
    border-radius: var(--radius-sm);
    border: 1px solid rgba(var(--color-warning-rgb), 0.4);
    background: rgba(var(--color-warning-rgb), 0.2);
    padding: var(--space-xs) var(--space-sm);
    font-size: var(--text-xs);
    color: var(--color-warning-300);
  }

  .config-link {
    display: inline-flex;
    align-items: center;
    font-size: var(--text-sm);
    color: var(--color-brand-500);
    transition: color var(--transition-normal) var(--ease-out);

    &:hover {
      color: var(--color-brand-400);
    }
  }

  .setting-description {
    margin-top: var(--space-xs);
    font-size: var(--text-xs);
    color: var(--color-content-subtle);
  }

  .feature-list {
    color: var(--color-content-muted);
    margin-top: var(--space-xl);
    display: flex;
    flex-direction: column;
    gap: var(--space-xs);
    font-size: var(--text-sm);
  }

  .hover:shadow-glow:hover {
    box-shadow: var(--shadow-glow-md);
  }

  .slider-container {
    margin-bottom: var(--space-sm);
  }
</style>
