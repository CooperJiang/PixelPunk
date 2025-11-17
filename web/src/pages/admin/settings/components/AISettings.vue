<script setup lang="ts">
  import { computed, onMounted, reactive, ref, watch } from 'vue'
  import { checkAIConfiguration, defaultSettings, testAIConfig, type AITestResult, type Setting } from '@/api/admin/settings'
  import { useToast } from '@/components/Toast/useToast'
  import { useTexts } from '@/composables/useTexts'

  const { $t } = useTexts()

  /* Props 和 Emits */
  const props = defineProps<{
    settings: Setting[]
  }>()

  const emit = defineEmits<{
    (e: 'update', settings: Setting[]): void
  }>()

  /* 工具 */
  const toast = useToast()

  /* 测试相关状态 */
  const testing = ref(false)
  const testResult = ref<AITestResult | null>(null)

  /* 从defaultSettings获取AI设置的默认值 */
  const aiDefaults = defaultSettings.ai.reduce(
    (acc, setting) => {
      acc[setting.key] = setting.value
      return acc
    },
    {} as Record<string, unknown>
  )

  /* 本地设置对象（扁平化）- 使用导入的默认值 */
  const localSettings = reactive({
    ai_enabled: aiDefaults.ai_enabled ?? false,
    ai_proxy: aiDefaults.ai_proxy || '',
    ai_model: aiDefaults.ai_model || 'gpt-4o-mini',
    ai_api_key: aiDefaults.ai_api_key || '',
    ai_temperature: aiDefaults.ai_temperature ?? 0.1,
    ai_max_tokens: aiDefaults.ai_max_tokens || 8096,
    ai_concurrency: aiDefaults.ai_concurrency || 5,
  })

  const sanitizeInput = (value: string): string => {
    if (!value) return ''
    return value
      .trim() // 去除首尾空格
      .replace(/[\s\n\r\t]/g, '')
      .replace(/["'`]/g, '')
  }

  watch(
    () => localSettings.ai_proxy,
    (newValue) => {
      const sanitized = sanitizeInput(newValue as string)
      if (sanitized !== newValue) {
        localSettings.ai_proxy = sanitized
      }
    }
  )

  watch(
    () => localSettings.ai_api_key,
    (newValue) => {
      const sanitized = sanitizeInput(newValue as string)
      if (sanitized !== newValue) {
        localSettings.ai_api_key = sanitized
      }
    }
  )

  const modelOptions = computed(() =>
    [
      'gpt-5-mini',
      'gpt-5-nano',
      'gpt-4.1-nano',
      'gpt-4o-mini',
      'gpt-4.1-mini',
      'o3-mini',
      'o4-mini',
      'o4-mini-high',
      'gpt-4o',
      'gpt-4-turbo-with-vision',
      'gpt-4.1-mini-2025-04-14',
      'gpt-4.1-nano-2025-04-14',
    ].map((model) => ({
      label: model,
      value: model,
    }))
  )

  const canTest = computed(() => {
    const config = checkAIConfiguration(getSettingsArray())
    return config.hasApiKey && config.hasModel
  })

  const testAIConfiguration = async () => {
    if (!canTest.value) {
      toast.error($t('admin.settings.ai.toast.configRequired'))
      return
    }

    testing.value = true
    testResult.value = null

    try {
      const testConfig = {
        ai_api_key: localSettings.ai_api_key,
        ai_proxy: localSettings.ai_proxy || 'https://api.openai.com',
        ai_model: localSettings.ai_model,
        ai_temperature: localSettings.ai_temperature,
        ai_max_tokens: localSettings.ai_max_tokens,
      }

      const response = await testAIConfig(testConfig)

      if (response.success) {
        const result = response.data?.data || response.data
        testResult.value = result

        if (result?.success) {
          toast.success($t('admin.settings.ai.toast.testSuccess'))
        } else {
          toast.error(
            $t('admin.settings.ai.toast.testFailed').replace(
              '{message}',
              result?.message || $t('admin.settings.ai.test.errorUnknown')
            )
          )
        }
      } else {
        testResult.value = {
          success: false,
          message: response.message || $t('admin.settings.ai.test.apiCallFailed'),
          error: 'api_error',
        }
        toast.error(
          $t('admin.settings.ai.toast.testFailed').replace(
            '{message}',
            response.message || $t('admin.settings.ai.test.errorUnknown')
          )
        )
      }
    } catch (error: unknown) {
      const errorMessage = (error as any).message || $t('admin.settings.ai.test.networkError')
      testResult.value = {
        success: false,
        message: $t('admin.settings.ai.test.requestFailed'),
        error: 'network_error',
        details: errorMessage,
      }
      toast.error(error.message)
    }

    testing.value = false
  }

  const getSettingsArray = () =>
    Object.entries(localSettings).map(([key, value]) => ({
      key,
      value,
      type: typeof value === 'number' ? 'number' : typeof value === 'boolean' ? 'boolean' : 'string',
      group: 'ai',
      description: '',
      is_system: true,
    }))

  const initializeSettings = () => {
    props.settings.forEach((setting) => {
      if (setting.key in localSettings) {
        const settingsRecord = localSettings as Record<string, unknown>
        settingsRecord[setting.key] = setting.value
      }
    })
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

  onMounted(() => {
    initializeSettings()
  })
</script>

<template>
  <div class="space-y-5">
    <div class="space-y-5">
      <div class="flex flex-col md:flex-row md:items-center">
        <label class="flex w-full items-center text-sm text-content md:w-48 md:py-2">
          <i class="fas fa-power-off mr-2 text-brand-600" />{{ $t('admin.settings.ai.enabled.label') }}
        </label>
        <div class="flex-1">
          <div class="flex items-center">
            <CyberSwitch v-model="localSettings.ai_enabled" />
            <span class="ml-3 text-sm text-content-muted">{{
              localSettings.ai_enabled ? $t('admin.settings.ai.enabled.activeText') : $t('admin.settings.ai.enabled.inactiveText')
            }}</span>
          </div>
          <p class="text-content-subtle mt-1 text-xs">{{ $t('admin.settings.ai.enabled.description') }}</p>
        </div>
      </div>

      <div class="flex flex-col md:flex-row md:items-center">
        <label class="flex w-full items-center text-sm text-content md:w-48 md:py-2">
          <i class="fas fa-network-wired mr-2 text-brand-600" />{{ $t('admin.settings.ai.proxy.label') }}
        </label>
        <div class="flex-1">
          <CyberInput v-model="localSettings.ai_proxy" :placeholder="$t('admin.settings.ai.proxy.placeholder')" width="500px" />
          <p class="text-content-subtle mt-1 text-xs">{{ $t('admin.settings.ai.proxy.description') }}</p>
        </div>
      </div>

      <div class="flex flex-col md:flex-row md:items-center">
        <label class="flex w-full items-center text-sm text-content md:w-48 md:py-2">
          <i class="fas fa-key mr-2 text-brand-600" />{{ $t('admin.settings.ai.apiKey.label') }}
        </label>
        <div class="flex-1">
          <CyberInput
            v-model="localSettings.ai_api_key"
            type="password"
            :placeholder="$t('admin.settings.ai.apiKey.placeholder')"
            width="500px"
          />
          <p class="text-content-subtle mt-1 text-xs">{{ $t('admin.settings.ai.apiKey.description') }}</p>
        </div>
      </div>

      <div class="flex flex-col md:flex-row md:items-center">
        <label class="flex w-full items-center text-sm text-content md:w-48 md:py-2">
          <i class="fas fa-robot mr-2 text-brand-600" />{{ $t('admin.settings.ai.model.label') }}
        </label>
        <div class="flex-1">
          <CyberDropdown
            v-model="localSettings.ai_model"
            :options="modelOptions"
            :placeholder="$t('admin.settings.ai.model.placeholder')"
            width="500px"
          />
          <p class="text-content-subtle mt-1 text-xs">{{ $t('admin.settings.ai.model.description') }}</p>
        </div>
      </div>

      <div class="flex flex-col md:flex-row md:items-center">
        <label class="flex w-full items-center text-sm text-content md:w-48 md:py-2">
          <i class="fas fa-thermometer-half mr-2 text-brand-600" />{{ $t('admin.settings.ai.temperature.label') }}
        </label>
        <div class="flex-1">
          <div class="slider-container flex items-center">
            <CyberSlider
              v-model="localSettings.ai_temperature"
              :min="0"
              :max="1"
              :step="0.01"
              :description="$t('admin.settings.ai.temperature.description')"
              width="460px"
            />
            <span class="ml-3 flex h-full items-center text-sm text-content-muted">{{
              localSettings.ai_temperature.toFixed(2)
            }}</span>
          </div>
        </div>
      </div>

      <div class="flex flex-col md:flex-row md:items-center">
        <label class="flex w-full items-center text-sm text-content md:w-48 md:py-2">
          <i class="fas fa-text-width mr-2 text-brand-600" />{{ $t('admin.settings.ai.maxTokens.label') }}
        </label>
        <div class="flex-1">
          <CyberInput
            v-model.number="localSettings.ai_max_tokens"
            type="number"
            :placeholder="$t('admin.settings.ai.maxTokens.placeholder')"
            width="500px"
          />
          <p class="text-content-subtle mt-1 text-xs">{{ $t('admin.settings.ai.maxTokens.description') }}</p>
        </div>
      </div>

      <div class="flex flex-col md:flex-row md:items-center">
        <label class="flex w-full items-center text-sm text-content md:w-48 md:py-2">
          <i class="fas fa-tasks mr-2 text-brand-600" />{{ $t('admin.settings.ai.concurrency.label') }}
        </label>
        <div class="flex-1">
          <CyberInput
            v-model.number="localSettings.ai_concurrency"
            type="number"
            :placeholder="$t('admin.settings.ai.concurrency.placeholder')"
            width="500px"
            :min="1"
            :max="20"
          >
            <template #unit>{{ $t('admin.settings.ai.concurrency.unit') }}</template>
          </CyberInput>
          <p class="text-content-subtle mt-1 text-xs">{{ $t('admin.settings.ai.concurrency.description') }}</p>
        </div>
      </div>

      <!-- 测试AI配置 -->
      <div class="mt-6 border-t border-subtle pt-5">
        <div class="flex flex-col md:flex-row md:items-center md:justify-between">
          <div class="mb-4 md:mb-0">
            <h4 class="mb-1 font-medium text-content-heading">{{ $t('admin.settings.ai.test.title') }}</h4>
            <p class="text-sm text-content-muted">{{ $t('admin.settings.ai.test.description') }}</p>
          </div>
          <div class="flex items-center space-x-3">
            <div v-if="testResult" class="flex items-center space-x-2">
              <div
                v-if="testResult.success"
                class="flex items-center rounded-lg border border-green-500/40 bg-green-900/30 px-3 py-1"
              >
                <i class="fas fa-check-circle mr-2 text-green-400" />
                <span class="text-sm font-medium text-green-300">{{ $t('admin.settings.ai.test.successStatus') }}</span>
              </div>
              <div v-else class="flex items-center rounded-lg border border-red-500/40 bg-red-900/30 px-3 py-1">
                <i class="fas fa-times-circle mr-2 text-red-400" />
                <span class="text-sm font-medium text-red-300">{{ $t('admin.settings.ai.test.failedStatus') }}</span>
              </div>
            </div>
            <CyberButton
              type="secondary"
              :loading="testing"
              :show-spinner="false"
              :disabled="!canTest"
              @click="testAIConfiguration"
            >
              <i class="fas fa-flask mr-2" />
              {{ testing ? $t('admin.settings.ai.test.testing') : $t('admin.settings.ai.test.button') }}
            </CyberButton>
          </div>
        </div>

        <!-- 测试结果详情 -->
        <div v-if="testResult && testResult.success" class="mt-4 rounded-lg border border-green-500/30 bg-green-900/20 p-4">
          <h5 class="mb-3 flex items-center font-medium text-green-300">
            <i class="fas fa-info-circle mr-2" />
            {{ $t('admin.settings.ai.test.detailsTitle') }}
          </h5>
          <div class="grid grid-cols-1 gap-4 text-sm md:grid-cols-2 lg:grid-cols-3">
            <div class="flex items-center">
              <i class="fas fa-robot mr-2 text-green-400/60" />
              <span class="text-content-muted">{{ $t('admin.settings.ai.test.modelLabel') }}</span>
              <span class="ml-2 font-medium text-green-400">{{ testResult.model }}</span>
            </div>
            <div class="flex items-center">
              <i class="fas fa-coins mr-2 text-green-400/60" />
              <span class="text-content-muted">{{ $t('admin.settings.ai.test.tokensLabel') }}</span>
              <span class="ml-2 font-medium text-green-400">{{ testResult.tokens_used }}</span>
            </div>
            <div class="flex items-center">
              <i class="fas fa-clock mr-2 text-green-400/60" />
              <span class="text-content-muted">{{ $t('admin.settings.ai.test.statusLabel') }}</span>
              <span class="ml-2 font-medium text-green-400">{{
                testResult.status || $t('admin.settings.ai.test.statusConnected')
              }}</span>
            </div>
            <div class="flex items-start md:col-span-2 lg:col-span-3">
              <i class="fas fa-comment-dots mr-2 mt-0.5 text-green-400/60" />
              <span class="text-content-muted">{{ $t('admin.settings.ai.test.responseLabel') }}</span>
              <span class="ml-2 font-medium text-green-400">{{
                testResult.test_response || $t('admin.settings.ai.test.responseNormal')
              }}</span>
            </div>
          </div>
        </div>

        <div v-else-if="testResult && !testResult.success" class="mt-4 rounded-lg border border-red-500/30 bg-red-900/20 p-4">
          <h5 class="mb-3 flex items-center font-medium text-red-300">
            <i class="fas fa-exclamation-triangle mr-2" />
            {{ $t('admin.settings.ai.test.errorTitle') }}
          </h5>
          <div class="space-y-3 text-sm">
            <div class="flex items-start">
              <i class="fas fa-times-circle mr-2 mt-0.5 text-red-400/60" />
              <div>
                <span class="text-content-muted">{{ $t('admin.settings.ai.test.errorReasonLabel') }}</span>
                <span class="ml-2 font-medium text-red-400">{{
                  testResult.message || $t('admin.settings.ai.test.errorUnknown')
                }}</span>
              </div>
            </div>
            <div v-if="testResult.details" class="flex items-start">
              <i class="fas fa-info-circle mr-2 mt-0.5 text-red-400/60" />
              <div>
                <span class="text-content-muted">{{ $t('admin.settings.ai.test.errorDetailsLabel') }}</span>
                <span class="ml-2 text-red-300">{{ testResult.details }}</span>
              </div>
            </div>
            <div v-if="testResult.error" class="flex items-start">
              <i class="fas fa-bug mr-2 mt-0.5 text-red-400/60" />
              <div>
                <span class="text-content-muted">{{ $t('admin.settings.ai.test.errorTypeLabel') }}</span>
                <span class="ml-2 font-mono text-xs text-red-300">{{ testResult.error }}</span>
              </div>
            </div>
          </div>
        </div>

        <!-- 功能说明 -->
        <div
          class="mt-8 overflow-hidden rounded-xl border border-emerald-500/20 bg-gradient-to-br from-emerald-900/10 to-green-900/10 p-6 backdrop-blur-sm"
        >
          <h4 class="mb-4 flex items-center text-lg font-semibold text-content-heading">
            <div class="mr-3 flex h-8 w-8 items-center justify-center rounded-lg bg-emerald-500/20">
              <i class="fas fa-robot text-emerald-400" />
            </div>
            {{ $t('admin.settings.ai.features.title') }}
          </h4>

          <div class="mb-4 rounded-lg border border-emerald-500/10 bg-emerald-900/10 p-4">
            <div class="mb-2 flex items-center">
              <i class="fas fa-tags mr-2 text-emerald-400" />
              <span class="font-medium text-emerald-300">{{ $t('admin.settings.ai.features.overview.title') }}</span>
            </div>
            <p class="text-sm leading-relaxed text-content">
              {{ $t('admin.settings.ai.features.overview.description') }}
            </p>
          </div>

          <div class="grid gap-4 md:grid-cols-2">
            <div class="rounded-lg border border-blue-500/20 bg-blue-900/10 p-4">
              <div class="mb-3 flex items-center">
                <i class="fas fa-eye mr-2 text-blue-400" />
                <span class="font-medium text-blue-300">{{ $t('admin.settings.ai.features.imageRecognition.title') }}</span>
              </div>
              <div class="space-y-2 text-sm text-content-muted">
                <p>
                  •
                  <strong class="text-content-heading">{{
                    $t('admin.settings.ai.features.imageRecognition.objectDetection')
                  }}</strong>
                  - {{ $t('admin.settings.ai.features.imageRecognition.objectDetectionDesc') }}
                </p>
                <p>
                  •
                  <strong class="text-content-heading">{{
                    $t('admin.settings.ai.features.imageRecognition.sceneAnalysis')
                  }}</strong>
                  - {{ $t('admin.settings.ai.features.imageRecognition.sceneAnalysisDesc') }}
                </p>
                <p>
                  •
                  <strong class="text-content-heading">{{
                    $t('admin.settings.ai.features.imageRecognition.personDetection')
                  }}</strong>
                  - {{ $t('admin.settings.ai.features.imageRecognition.personDetectionDesc') }}
                </p>
                <p>
                  •
                  <strong class="text-content-heading">{{ $t('admin.settings.ai.features.imageRecognition.colorStyle') }}</strong>
                  - {{ $t('admin.settings.ai.features.imageRecognition.colorStyleDesc') }}
                </p>
              </div>
            </div>

            <div class="rounded-lg border border-purple-500/20 bg-purple-900/10 p-4">
              <div class="mb-3 flex items-center">
                <i class="fas fa-brain mr-2 text-purple-400" />
                <span class="font-medium text-purple-300">{{ $t('admin.settings.ai.features.smartClassification.title') }}</span>
              </div>
              <div class="space-y-2 text-sm text-content-muted">
                <p>
                  •
                  <strong class="text-content-heading">{{ $t('admin.settings.ai.features.smartClassification.autoTag') }}</strong>
                  - {{ $t('admin.settings.ai.features.smartClassification.autoTagDesc') }}
                </p>
                <p>
                  •
                  <strong class="text-content-heading">{{
                    $t('admin.settings.ai.features.smartClassification.contentDesc')
                  }}</strong>
                  - {{ $t('admin.settings.ai.features.smartClassification.contentDescDesc') }}
                </p>
                <p>
                  •
                  <strong class="text-content-heading">{{
                    $t('admin.settings.ai.features.smartClassification.categoryArchive')
                  }}</strong>
                  - {{ $t('admin.settings.ai.features.smartClassification.categoryArchiveDesc') }}
                </p>
                <p>
                  •
                  <strong class="text-content-heading">{{
                    $t('admin.settings.ai.features.smartClassification.similarity')
                  }}</strong>
                  - {{ $t('admin.settings.ai.features.smartClassification.similarityDesc') }}
                </p>
              </div>
            </div>
          </div>

          <div class="mt-4 rounded-lg border border-amber-500/20 bg-amber-900/10 p-4">
            <div class="mb-3 flex items-center">
              <i class="fas fa-cogs mr-2 text-amber-400" />
              <span class="font-medium text-amber-300">{{ $t('admin.settings.ai.features.modelConfig.title') }}</span>
            </div>
            <div class="grid gap-3 text-sm text-content-muted md:grid-cols-3">
              <div class="flex items-center">
                <div class="mr-2 h-2 w-2 rounded-full bg-green-400"></div>
                <span
                  ><strong class="text-green-300">{{ $t('admin.settings.ai.features.modelConfig.temperature') }}</strong> -
                  {{ $t('admin.settings.ai.features.modelConfig.temperatureDesc') }}</span
                >
              </div>
              <div class="flex items-center">
                <div class="mr-2 h-2 w-2 rounded-full bg-blue-400"></div>
                <span
                  ><strong class="text-blue-300">{{ $t('admin.settings.ai.features.modelConfig.concurrency') }}</strong> -
                  {{ $t('admin.settings.ai.features.modelConfig.concurrencyDesc') }}</span
                >
              </div>
              <div class="flex items-center">
                <div class="mr-2 h-2 w-2 rounded-full bg-purple-400"></div>
                <span
                  ><strong class="text-purple-300">{{ $t('admin.settings.ai.features.modelConfig.tokens') }}</strong> -
                  {{ $t('admin.settings.ai.features.modelConfig.tokensDesc') }}</span
                >
              </div>
            </div>
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

  .slider-container {
    margin-bottom: 6px;
  }
</style>
