<script setup lang="ts">
  import { computed, onMounted, ref, watch } from 'vue'
  import { useToast } from '@/components/Toast/useToast'
  import { useRoute, useRouter } from 'vue-router'
  import { useSettingsStore } from '@/store/settings'
  import { useTexts } from '@/composables/useTexts'
  import AISettings from './components/AISettings.vue'
  import VectorSettings from './components/VectorSettings.vue'
  import ContentDetectionSettings from './components/ContentDetectionSettings.vue'
  import {
    batchUpsertSettings,
    checkAIConfiguration,
    defaultSettings,
    getSettings,
    testAIConfig,
    regenerateVectors as apiRegenerateVectors,
    testVectorConfig as apiTestVectorConfig,
    type AITestResult,
    type Setting,
    type SettingGroup,
  } from '@/api/admin/settings'

  defineOptions({
    name: 'AIManagementPage',
  })

  const { $t } = useTexts()
  const toast = useToast()
  const route = useRoute()
  const router = useRouter()
  const settingsStore = useSettingsStore()

  const isLoading = ref(false)
  const isSaving = ref(false)

  const testing = ref(false)
  const testResult = ref<AITestResult | null>(null)
  const showTestDialog = ref(false)

  const vectorTesting = ref(false)
  const vectorTestResult = ref<{
    success: boolean
    message?: string
    error?: string
    details?: string
    saveTestFailed?: boolean
    config?: {
      model?: string
      dimension?: number
      provider?: string
      baseUrl?: string
    }
    response_time?: number
  } | null>(null)
  const showVectorTestDialog = ref(false)
  const regeneratingVectors = ref(false)
  const showRegenerateDialog = ref(false)

  const showAIFeaturePromptDialog = ref(false)
  const aiFeaturePromptData = ref<{
    needContentDetection: boolean
    needAIAnalysis: boolean
  } | null>(null)

  interface TabItem {
    key: string
    name: string
    icon?: string
  }

  const tabs = computed<TabItem[]>(() => [
    { key: 'ai', name: $t('admin.ai.tabs.ai'), icon: 'fas fa-robot' },
    { key: 'vector', name: $t('admin.ai.tabs.vector'), icon: 'fas fa-search' },
    { key: 'content-detection', name: $t('admin.ai.tabs.contentDetection'), icon: 'fas fa-shield-alt' },
  ])

  const activeTab = ref(route.query.tab?.toString() || 'ai')

  const updateActiveTab = (tabId: string) => {
    activeTab.value = tabId
    router.push({
      name: route.name,
      query: { ...route.query, tab: tabId },
    })
  }

  const aiSettings = ref<Setting[]>([])
  const vectorSettings = ref<Setting[]>([])
  const contentDetectionSettings = ref<Setting[]>([])

  const modifiedAISettings = ref<Setting[]>([])
  const modifiedVectorSettings = ref<Setting[]>([])
  const modifiedContentDetectionSettings = ref<Setting[]>([])

  const handleAISettingsUpdate = (updatedSettings: Setting[]) => {
    modifiedAISettings.value = updatedSettings
  }

  const handleVectorSettingsUpdate = (updatedSettings: Setting[]) => {
    modifiedVectorSettings.value = updatedSettings
  }

  const handleContentDetectionSettingsUpdate = (updatedSettings: Setting[]) => {
    modifiedContentDetectionSettings.value = updatedSettings
  }

  const resetSettings = async () => {
    isLoading.value = true
    await loadSettingsForTab(activeTab.value as SettingGroup)

    if (activeTab.value === 'ai') {
      modifiedAISettings.value = []
    } else if (activeTab.value === 'vector') {
      modifiedVectorSettings.value = []
    } else if (activeTab.value === 'content-detection') {
      modifiedContentDetectionSettings.value = []
    }

    toast.success($t('admin.ai.messages.settingsReset'))
    isLoading.value = false
  }

  const saveSettings = async () => {
    let settingsToUpdate: Setting[] = []
    if (activeTab.value === 'ai' && modifiedAISettings.value.length > 0) {
      settingsToUpdate = modifiedAISettings.value
    } else if (activeTab.value === 'vector' && modifiedVectorSettings.value.length > 0) {
      settingsToUpdate = modifiedVectorSettings.value
    } else if (activeTab.value === 'content-detection' && modifiedContentDetectionSettings.value.length > 0) {
      settingsToUpdate = modifiedContentDetectionSettings.value
    }

    if (settingsToUpdate.length === 0) {
      toast.info($t('admin.ai.messages.noChanges'))
      return
    }

    try {
      isSaving.value = true
      const result = await batchUpsertSettings(settingsToUpdate)

      if (result.success) {
        if (activeTab.value === 'ai') {
          modifiedAISettings.value = []
          await settingsStore.loadGlobalSettings()
        } else if (activeTab.value === 'vector') {
          modifiedVectorSettings.value = []
          await settingsStore.loadGlobalSettings()
        } else if (activeTab.value === 'content-detection') {
          modifiedContentDetectionSettings.value = []
        }
        toast.success($t('admin.ai.messages.settingsSaved'))
        await loadSettingsForTab(activeTab.value as SettingGroup)

        await testConfigurationAfterSave()
      }
    } catch (error: unknown) {
      toast.error(error.message)
    } finally {
      isSaving.value = false
    }
  }

  const checkAIFeatureStatus = async () => {
    try {
      const [uploadResult, aiResult] = await Promise.all([getSettings({ group: 'upload' }), getSettings({ group: 'ai' })])

      const uploadSettings = uploadResult.success ? uploadResult.data?.settings || [] : []
      const aiSettings = aiResult.success ? aiResult.data?.settings || [] : []
      const allSettings = [...uploadSettings, ...aiSettings]

      const settingsMap = allSettings.reduce(
        (acc, setting) => {
          acc[setting.key] = setting.value
          return acc
        },
        {} as Record<string, unknown>
      )

      const contentDetectionEnabled = settingsMap.content_detection_enabled
      const aiAnalysisEnabled = settingsMap.ai_analysis_enabled

      if (!contentDetectionEnabled || !aiAnalysisEnabled) {
        aiFeaturePromptData.value = {
          needContentDetection: !contentDetectionEnabled,
          needAIAnalysis: !aiAnalysisEnabled,
        }
        showAIFeaturePromptDialog.value = true
      }
    } catch {}
  }

  const goToContentDetectionSettings = () => {
    showAIFeaturePromptDialog.value = false
    router.push('/admin/ai?tab=content-detection')
  }

  const testConfigurationAfterSave = async () => {
    try {
      if (activeTab.value === 'ai') {
        const currentAISettings = aiSettings.value // 保存后使用最新的数据库数据
        const aiSettingsMap = currentAISettings.reduce(
          (acc, setting) => {
            acc[setting.key] = setting.value
            return acc
          },
          {} as Record<string, unknown>
        )

        if (!aiSettingsMap.ai_enabled) {
          return // AI功能未启用，跳过测试
        }

        const config = checkAIConfiguration(currentAISettings)
        if (!config.hasApiKey || !config.hasModel) {
          return // 配置不完整，跳过测试
        }

        const testConfig = {
          ai_api_key: aiSettingsMap.ai_api_key,
          ai_proxy: aiSettingsMap.ai_proxy || 'https://api.openai.com',
          ai_model: aiSettingsMap.ai_model,
          ai_temperature: aiSettingsMap.ai_temperature,
          ai_max_tokens: aiSettingsMap.ai_max_tokens,
        }

        const response = await testAIConfig(testConfig)
        const result = response.data?.data || response.data

        if (!response.success || !result?.success) {
          testResult.value = {
            success: false,
            message: result?.message || response.message || $t('admin.ai.dialogs.saveSuccessButTestFailed'),
            error: result?.error || 'config_test_failed',
            details: result?.details,
            saveTestFailed: true, // 标记为保存后测试失败
          }
          showTestDialog.value = true
        } else {
          await checkAIFeatureStatus()
        }
      } else if (activeTab.value === 'vector') {
        const currentVectorSettings = vectorSettings.value // 保存后使用最新的数据库数据
        const vectorSettingsMap = currentVectorSettings.reduce(
          (acc, setting) => {
            acc[setting.key] = setting.value
            return acc
          },
          {} as Record<string, unknown>
        )

        if (!vectorSettingsMap.vector_enabled) {
          return // 向量功能未启用，跳过测试
        }

        if (!vectorSettingsMap.vector_api_key || !vectorSettingsMap.vector_model) {
          return // 配置不完整，跳过测试
        }

        const testParams = {
          provider: vectorSettingsMap.vector_provider || 'openai',
          api_key: vectorSettingsMap.vector_api_key,
          model: vectorSettingsMap.vector_model,
          base_url: vectorSettingsMap.vector_base_url || 'https://api.openai.com/v1',
          timeout: vectorSettingsMap.vector_timeout || 30,
          test_text: $t('admin.ai.testText'),
        }

        const result = await apiTestVectorConfig(testParams)

        if (!result.success || !result.data?.success) {
          vectorTestResult.value = {
            success: false,
            message: result.data?.message || result.message || $t('admin.ai.dialogs.saveSuccessButTestFailed'),
            error: result.data?.error || 'config_test_failed',
            details: result.data?.details,
            saveTestFailed: true, // 标记为保存后测试失败
          }
          showVectorTestDialog.value = true
        }
      }
    } catch (error) {
      if (activeTab.value === 'ai') {
        testResult.value = {
          success: false,
          message: $t('admin.ai.dialogs.saveSuccessButTestFailed'),
          error: 'network_error',
          details: (error as any).message || $t('admin.ai.messages.networkError'),
          saveTestFailed: true, // 标记为保存后测试失败
        }
        showTestDialog.value = true
      } else if (activeTab.value === 'vector') {
        vectorTestResult.value = {
          success: false,
          message: $t('admin.ai.dialogs.saveSuccessButTestFailed'),
          error: 'network_error',
          details: (error as any).message || $t('admin.ai.messages.networkError'),
          saveTestFailed: true, // 标记为保存后测试失败
        }
        showVectorTestDialog.value = true
      }
    }
  }

  const testCurrentConfiguration = async () => {
    if (activeTab.value !== 'ai') {
      toast.error($t('admin.ai.messages.onlyInAITab'))
      return
    }

    const currentAISettings = modifiedAISettings.value.length > 0 ? modifiedAISettings.value : aiSettings.value
    const aiSettingsMap = currentAISettings.reduce(
      (acc, setting) => {
        acc[setting.key] = setting.value
        return acc
      },
      {} as Record<string, unknown>
    )

    const config = checkAIConfiguration(currentAISettings)
    if (!config.hasApiKey || !config.hasModel) {
      toast.error($t('admin.ai.messages.needApiKeyAndModel'))
      return
    }

    testing.value = true
    testResult.value = null

    try {
      const testConfig = {
        ai_api_key: aiSettingsMap.ai_api_key,
        ai_proxy: aiSettingsMap.ai_proxy || 'https://api.openai.com',
        ai_model: aiSettingsMap.ai_model,
        ai_temperature: aiSettingsMap.ai_temperature,
        ai_max_tokens: aiSettingsMap.ai_max_tokens,
      }

      const [response] = await Promise.all([testAIConfig(testConfig), new Promise((resolve) => setTimeout(resolve, 1000))])

      if (response.success) {
        const result = response.data?.data || response.data
        testResult.value = result

        if (result?.success) {
          toast.success($t('admin.ai.messages.aiTestSuccess'))
        } else {
          toast.error(
            $t('admin.ai.messages.aiTestFailed').replace('{message}', result?.message || $t('admin.ai.messages.unknownError'))
          )
        }
        showTestDialog.value = true
      } else {
        testResult.value = {
          success: false,
          message: response.message || $t('admin.ai.messages.testRequestFailed'),
          error: 'api_error',
        }
        toast.error(
          $t('admin.ai.messages.aiTestFailed').replace('{message}', response.message || $t('admin.ai.messages.unknownError'))
        )
        showTestDialog.value = true
      }
    } catch (error: unknown) {
      const errorMessage = (error as any).message || $t('admin.ai.messages.networkError')
      testResult.value = {
        success: false,
        message: $t('admin.ai.messages.testRequestFailed'),
        error: 'network_error',
        details: errorMessage,
      }
      toast.error(error.message)
      showTestDialog.value = true
    }

    testing.value = false
  }

  const testVectorConfiguration = async () => {
    if (activeTab.value !== 'vector') {
      toast.error($t('admin.ai.messages.onlyInVectorTab'))
      return
    }

    const currentVectorSettings = modifiedVectorSettings.value.length > 0 ? modifiedVectorSettings.value : vectorSettings.value
    const vectorSettingsMap = currentVectorSettings.reduce(
      (acc, setting) => {
        acc[setting.key] = setting.value
        return acc
      },
      {} as Record<string, unknown>
    )

    if (!vectorSettingsMap.vector_api_key || !vectorSettingsMap.vector_model) {
      toast.error($t('admin.ai.messages.needVectorApiKeyAndModel'))
      return
    }

    vectorTesting.value = true
    vectorTestResult.value = null

    try {
      const testParams = {
        provider: vectorSettingsMap.vector_provider || 'openai',
        api_key: vectorSettingsMap.vector_api_key,
        model: vectorSettingsMap.vector_model,
        base_url: vectorSettingsMap.vector_base_url || 'https://api.openai.com/v1',
        timeout: vectorSettingsMap.vector_timeout || 30,
        test_text: $t('admin.ai.testText'),
      }

      const result = await apiTestVectorConfig(testParams)

      if (result.success) {
        vectorTestResult.value = result.data
        if (result.data?.success) {
          toast.success($t('admin.ai.messages.vectorTestSuccess'))
        } else {
          toast.error($t('admin.ai.messages.vectorTestFailed').replace('{message}', ''))
        }
        showVectorTestDialog.value = true
      } else {
        vectorTestResult.value = {
          success: false,
          message: result.message || $t('admin.ai.messages.testRequestFailed'),
          error: 'api_error',
        }
        toast.error(
          $t('admin.ai.messages.vectorTestFailed').replace('{message}', result.message || $t('admin.ai.messages.unknownError'))
        )
        showVectorTestDialog.value = true
      }
    } catch (error: unknown) {
      vectorTestResult.value = {
        success: false,
        message: $t('admin.ai.messages.vectorTestFailed').replace('{message}', ''),
        error: 'network_error',
        details: (error as any).message,
      }
      toast.error(
        $t('admin.ai.messages.vectorTestFailed').replace('{message}', error.message || $t('admin.ai.messages.unknownError'))
      )
      showVectorTestDialog.value = true
    }

    vectorTesting.value = false
  }

  const regenerateVectors = async () => {
    regeneratingVectors.value = true

    try {
      const result = await apiRegenerateVectors()

      if (result.success) {
        const updateCount = result.data?.updated_count || 0
        toast.success($t('admin.ai.messages.vectorRegenerateSuccess').replace('{count}', String(updateCount)))
      } else {
        toast.error(
          $t('admin.ai.messages.vectorRegenerateFailed').replace(
            '{error}',
            result.message || $t('admin.ai.messages.unknownError')
          )
        )
      }
    } catch (error: unknown) {
      toast.error(
        $t('admin.ai.messages.vectorRegenerateFailed').replace('{error}', error.message || $t('admin.ai.messages.unknownError'))
      )
    }

    regeneratingVectors.value = false
    showRegenerateDialog.value = false
  }

  const getTestConfig = () => {
    const currentAISettings = modifiedAISettings.value.length > 0 ? modifiedAISettings.value : aiSettings.value
    if (!currentAISettings.length) return {}

    const aiSettingsMap = currentAISettings.reduce(
      (acc, setting) => {
        acc[setting.key] = setting.value
        return acc
      },
      {} as Record<string, unknown>
    )

    return {
      ai_api_key: aiSettingsMap.ai_api_key || '',
      ai_proxy: aiSettingsMap.ai_proxy || 'https://api.openai.com',
      ai_model: aiSettingsMap.ai_model || '',
      ai_temperature: aiSettingsMap.ai_temperature || 0.1,
      ai_max_tokens: aiSettingsMap.ai_max_tokens || 8096,
    }
  }

  const formatConfigKey = (key: string) => {
    return $t(`admin.ai.configFields.${key}`) || key
  }

  const loadSettingsForTab = async (tabId: SettingGroup) => {
    isLoading.value = true

    try {
      if (tabId === 'content-detection') {
        const [uploadResult, aiResult] = await Promise.all([getSettings({ group: 'upload' }), getSettings({ group: 'ai' })])

        const uploadSettings = uploadResult.success ? uploadResult.data?.settings || [] : []
        const aiSettings = aiResult.success ? aiResult.data?.settings || [] : []

        const allSettings = [...uploadSettings, ...aiSettings]
        const contentDetectionKeys = [
          'content_detection_enabled',
          'sensitive_content_handling',
          'ai_analysis_enabled',
          'nsfw_threshold',
        ]

        contentDetectionSettings.value = allSettings.filter((s) => contentDetectionKeys.includes(s.key))

        const existingKeys = contentDetectionSettings.value.map((s) => s.key)
        const missingKeys = contentDetectionKeys.filter((key) => !existingKeys.includes(key))

        if (missingKeys.length > 0) {
          const uploadDefaults = defaultSettings.upload || []
          const aiDefaults = defaultSettings.ai || []
          const allDefaults = [...uploadDefaults, ...aiDefaults]

          const missingSettings = allDefaults.filter((s) => missingKeys.includes(s.key))
          contentDetectionSettings.value = [...contentDetectionSettings.value, ...missingSettings]
        }
      } else {
        const group = tabId as SettingGroup
        const result = await getSettings({ group })
        const settings = result.success ? result.data?.settings || [] : []

        if (settings.length === 0) {
          const defaults = defaultSettings[group] || []
          if (tabId === 'ai') {
            aiSettings.value = defaults
          } else if (tabId === 'vector') {
            vectorSettings.value = defaults
          }
        } else {
          if (tabId === 'ai') {
            aiSettings.value = settings
          } else if (tabId === 'vector') {
            vectorSettings.value = settings
          }
        }
      }
    } catch (error: unknown) {
      if (tabId === 'content-detection') {
        const uploadDefaults = defaultSettings.upload || []
        const aiDefaults = defaultSettings.ai || []
        const allDefaults = [...uploadDefaults, ...aiDefaults]
        const contentDetectionKeys = [
          'content_detection_enabled',
          'sensitive_content_handling',
          'ai_analysis_enabled',
          'nsfw_threshold',
        ]

        contentDetectionSettings.value = allDefaults.filter((s) => contentDetectionKeys.includes(s.key))
      } else {
        const group = tabId as SettingGroup
        const defaults = defaultSettings[group] || []

        if (tabId === 'ai') {
          aiSettings.value = defaults
        } else if (tabId === 'vector') {
          vectorSettings.value = defaults
        }
      }
      toast.error(error.message)
    } finally {
      isLoading.value = false
    }
  }

  watch(
    () => route.query.tab,
    (newTab) => {
      if (newTab && typeof newTab === 'string' && newTab !== activeTab.value) {
        activeTab.value = newTab
      }
    }
  )

  watch(activeTab, async (newTab) => {
    const shouldLoad =
      (newTab === 'ai' && aiSettings.value.length === 0) ||
      (newTab === 'vector' && vectorSettings.value.length === 0) ||
      (newTab === 'content-detection' && contentDetectionSettings.value.length === 0)

    if (shouldLoad) {
      await loadSettingsForTab(newTab as SettingGroup)
    }

    if (newTab === 'content-detection' && aiSettings.value.length === 0) {
      await loadSettingsForTab('ai')
    }
  })

  onMounted(async () => {
    await loadSettingsForTab(activeTab.value as SettingGroup)
  })
</script>

<template>
  <div class="admin-ai-page admin-page-container">
    <CyberAdminWrapper
      :title="$t('admin.ai.title')"
      :subtitle="$t('admin.ai.subtitle')"
      icon="fas fa-robot"
      :sidebar-layout="true"
    >
      <template #actions>
        <CyberButton
          v-if="activeTab === 'ai'"
          type="outlined"
          :loading="testing"
          loading-mode="inline"
          icon="flask"
          :disabled="aiSettings.length === 0"
          class="min-w-24"
          @click="testCurrentConfiguration"
        >
          {{ testing ? $t('admin.ai.actions.testing') : $t('admin.ai.actions.test') }}
        </CyberButton>

        <template v-if="activeTab === 'vector'">
          <CyberButton
            type="outlined"
            :loading="vectorTesting"
            loading-mode="inline"
            icon="flask"
            :disabled="vectorSettings.length === 0"
            class="min-w-24"
            @click="testVectorConfiguration"
          >
            {{ vectorTesting ? $t('admin.ai.actions.testing') : $t('admin.ai.actions.test') }}
          </CyberButton>
          <CyberButton
            type="warning"
            :loading="regeneratingVectors"
            loading-mode="inline"
            icon="sync"
            class="min-w-32"
            @click="showRegenerateDialog = true"
          >
            {{ regeneratingVectors ? $t('admin.ai.actions.regenerating') : $t('admin.ai.actions.regenerateVectors') }}
          </CyberButton>
        </template>

        <CyberButton type="text" class="min-w-24" @click="resetSettings"> {{ $t('admin.ai.actions.reset') }} </CyberButton>
        <CyberButton type="primary" :loading="isSaving" loading-mode="inline" icon="save" class="min-w-24" @click="saveSettings">
          {{ isSaving ? $t('admin.ai.actions.saving') : $t('admin.ai.actions.save') }}
        </CyberButton>
      </template>

      <template #sidebar>
        <CyberSidebarNav :tabs="tabs" :active-tab="activeTab" @tab-change="updateActiveTab" />
      </template>

      <template #content>
        <div v-show="activeTab === 'ai'" class="content-panel">
          <AISettings :settings="aiSettings" @update="handleAISettingsUpdate" />
        </div>

        <div v-show="activeTab === 'vector'" class="content-panel">
          <VectorSettings :settings="vectorSettings" @update="handleVectorSettingsUpdate" />
        </div>

        <div v-show="activeTab === 'content-detection'" class="content-panel">
          <ContentDetectionSettings
            :settings="contentDetectionSettings"
            :ai-settings="aiSettings"
            @update="handleContentDetectionSettingsUpdate"
          />
        </div>
      </template>
    </CyberAdminWrapper>

    <CyberLoading v-if="isLoading" />

    <CyberDialog
      v-model="showTestDialog"
      width="600px"
      :title="testResult?.saveTestFailed ? $t('admin.ai.dialogs.aiTestFailedAfterSave') : $t('admin.ai.dialogs.aiTestResult')"
      :show-default-footer="false"
      :show-footer="false"
    >
      <div v-if="testResult" class="space-y-4">
        <div v-if="testResult.saveTestFailed" class="status-warning rounded-lg border p-4">
          <div class="flex items-start space-x-3">
            <i class="fas fa-exclamation-triangle mt-0.5 text-warning" />
            <div>
              <h4 class="text-warning-heading mb-2 font-medium">{{ $t('admin.ai.dialogs.saveSuccessButTestFailed') }}</h4>
              <p class="text-muted text-sm">
                {{ $t('admin.ai.dialogs.saveSuccessButTestFailedDesc') }}
              </p>
            </div>
          </div>
        </div>

        <div class="flex items-center space-x-3">
          <div v-if="testResult.success" class="flex items-center space-x-2">
            <i class="fas fa-check-circle text-success" />
            <span class="text-success-heading font-medium">{{ $t('admin.ai.dialogs.testSuccess') }}</span>
          </div>
          <div v-else class="flex items-center space-x-2">
            <i class="fas fa-times-circle text-error" />
            <span class="text-error-heading font-medium">{{ $t('admin.ai.dialogs.testFailed') }}</span>
          </div>
        </div>

        <div v-if="testResult.success" class="status-success rounded-lg border p-4">
          <h4 class="text-success-heading mb-3 font-medium">{{ $t('admin.ai.dialogs.testDetails') }}</h4>
          <div class="grid grid-cols-1 gap-3 text-sm">
            <div class="flex items-center">
              <i class="fas fa-robot icon-muted mr-2 text-success" />
              <span class="text-muted min-w-20">{{ $t('admin.ai.testResults.model') }}:</span>
              <span
                class="ml-2 truncate font-medium text-success"
                :title="testResult.model || $t('admin.ai.testResults.unknown')"
                >{{ testResult.model || $t('admin.ai.testResults.unknown') }}</span
              >
            </div>
            <div class="flex items-center">
              <i class="fas fa-coins icon-muted mr-2 text-success" />
              <span class="text-muted min-w-20">{{ $t('admin.ai.testResults.tokensUsed') }}:</span>
              <span class="ml-2 font-medium text-success">{{ testResult.tokens_used || '0' }}</span>
            </div>
            <div class="flex items-center">
              <i class="fas fa-clock icon-muted mr-2 text-success" />
              <span class="text-muted min-w-20">{{ $t('admin.ai.testResults.status') }}:</span>
              <span class="ml-2 font-medium text-success">{{ testResult.status || $t('admin.ai.testResults.connected') }}</span>
            </div>
            <div v-if="testResult.test_response" class="flex items-center">
              <i class="fas fa-comment-dots icon-muted mr-2 text-success" />
              <span class="text-muted min-w-20">{{ $t('admin.ai.testResults.response') }}:</span>
              <span class="ml-2 font-medium text-success">{{ testResult.test_response }}</span>
            </div>
          </div>
        </div>

        <div v-else class="status-error rounded-lg border p-4">
          <h4 class="text-error-heading mb-3 font-medium">{{ $t('admin.ai.dialogs.errorDetails') }}</h4>
          <div class="space-y-3 text-sm">
            <div class="flex items-start">
              <i class="fas fa-times-circle icon-muted mr-2 mt-0.5 text-error" />
              <div>
                <span class="text-muted">{{ $t('admin.ai.testResults.errorReason') }}:</span>
                <p class="ml-0 mt-1 font-medium text-error">{{ testResult.message || $t('admin.ai.messages.unknownError') }}</p>
              </div>
            </div>
            <div v-if="testResult.details" class="flex items-start">
              <i class="fas fa-info-circle icon-muted mr-2 mt-0.5 text-error" />
              <div>
                <span class="text-muted">{{ $t('admin.ai.testResults.errorDetails') }}:</span>
                <p class="text-error-heading ml-0 mt-1">{{ testResult.details }}</p>
              </div>
            </div>
            <div v-if="testResult.error" class="flex items-start">
              <i class="fas fa-bug icon-muted mr-2 mt-0.5 text-error" />
              <div>
                <span class="text-muted">{{ $t('admin.ai.testResults.errorType') }}:</span>
                <p class="text-error-heading ml-0 mt-1 font-mono text-xs">{{ testResult.error }}</p>
              </div>
            </div>
          </div>
        </div>

        <div class="status-info rounded-lg border p-4">
          <h4 class="text-info-heading mb-3 font-medium">{{ $t('admin.ai.dialogs.testParams') }}</h4>
          <div class="grid grid-cols-1 gap-3 text-sm">
            <div v-for="(value, key) in getTestConfig()" :key="key" class="flex items-center">
              <span class="text-muted min-w-20">{{ formatConfigKey(key) }}:</span>
              <span class="text-info-heading ml-2 font-mono text-xs">
                {{ key === 'ai_api_key' ? '***' + String(value).slice(-6) : value }}
              </span>
            </div>
          </div>
        </div>
      </div>
    </CyberDialog>

    <CyberDialog
      v-model="showVectorTestDialog"
      width="600px"
      :title="
        vectorTestResult?.saveTestFailed
          ? $t('admin.ai.dialogs.vectorTestFailedAfterSave')
          : $t('admin.ai.dialogs.vectorTestResult')
      "
      :show-default-footer="false"
      :show-footer="false"
    >
      <div v-if="vectorTestResult" class="space-y-4">
        <div v-if="vectorTestResult.saveTestFailed" class="status-warning rounded-lg border p-4">
          <div class="flex items-start space-x-3">
            <i class="fas fa-exclamation-triangle mt-0.5 text-warning" />
            <div>
              <h4 class="text-warning-heading mb-2 font-medium">{{ $t('admin.ai.dialogs.saveSuccessButTestFailed') }}</h4>
              <p class="text-muted text-sm">
                {{ $t('admin.ai.dialogs.vectorSaveSuccessButTestFailedDesc') }}
              </p>
            </div>
          </div>
        </div>

        <div class="flex items-center space-x-3">
          <div v-if="vectorTestResult.success" class="flex items-center space-x-2">
            <i class="fas fa-check-circle text-success" />
            <span class="text-success-heading font-medium">{{ $t('admin.ai.dialogs.testSuccess') }}</span>
          </div>
          <div v-else class="flex items-center space-x-2">
            <i class="fas fa-times-circle text-error" />
            <span class="text-error-heading font-medium">{{ $t('admin.ai.dialogs.testFailed') }}</span>
          </div>
        </div>

        <div v-if="vectorTestResult.success" class="status-success rounded-lg border p-4">
          <h4 class="text-success-heading mb-3 font-medium">{{ $t('admin.ai.dialogs.testDetails') }}</h4>
          <div class="grid grid-cols-1 gap-3 text-sm">
            <div class="flex items-center">
              <i class="fas fa-cube icon-muted mr-2 w-4 text-success" />
              <span class="text-muted min-w-20">{{ $t('admin.ai.testResults.vectorModel') }}:</span>
              <span class="ml-2 font-medium text-success">{{
                vectorTestResult.config?.model || $t('admin.ai.testResults.unknown')
              }}</span>
            </div>
            <div class="flex items-center">
              <i class="fas fa-ruler icon-muted mr-2 w-4 text-success" />
              <span class="text-muted min-w-20">{{ $t('admin.ai.testResults.vectorDimension') }}:</span>
              <span class="ml-2 font-medium text-success">{{
                vectorTestResult.config?.dimension || $t('admin.ai.testResults.unknown')
              }}</span>
            </div>
            <div class="flex items-center">
              <i class="fas fa-server icon-muted mr-2 w-4 text-success" />
              <span class="text-muted min-w-20">{{ $t('admin.ai.testResults.vectorProvider') }}:</span>
              <span class="ml-2 font-medium text-success">{{
                vectorTestResult.config?.provider || $t('admin.ai.testResults.unknown')
              }}</span>
            </div>
            <div class="flex items-center">
              <i class="fas fa-network-wired icon-muted mr-2 w-4 text-success" />
              <span class="text-muted min-w-20">{{ $t('admin.ai.testResults.apiUrl') }}:</span>
              <span
                class="ml-2 truncate font-medium text-success"
                :title="vectorTestResult.config?.baseUrl || $t('admin.ai.testResults.unknown')"
                >{{ vectorTestResult.config?.baseUrl || $t('admin.ai.testResults.unknown') }}</span
              >
            </div>
            <div v-if="vectorTestResult.response_time" class="flex items-center">
              <i class="fas fa-clock icon-muted mr-2 w-4 text-success" />
              <span class="text-muted min-w-20">{{ $t('admin.ai.testResults.responseTime') }}:</span>
              <span class="ml-2 font-medium text-success">{{ vectorTestResult.response_time }}ms</span>
            </div>
          </div>
        </div>

        <div v-else class="status-error rounded-lg border p-4">
          <h4 class="text-error-heading mb-3 font-medium">{{ $t('admin.ai.dialogs.errorDetails') }}</h4>
          <div class="space-y-3 text-sm">
            <div class="flex items-start">
              <i class="fas fa-times-circle icon-muted mr-2 mt-0.5 text-error" />
              <div>
                <span class="text-muted">{{ $t('admin.ai.testResults.errorReason') }}:</span>
                <p class="ml-0 mt-1 font-medium text-error">
                  {{ vectorTestResult.message || $t('admin.ai.messages.unknownError') }}
                </p>
              </div>
            </div>
            <div v-if="vectorTestResult.details" class="flex items-start">
              <i class="fas fa-info-circle icon-muted mr-2 mt-0.5 text-error" />
              <div>
                <span class="text-muted">{{ $t('admin.ai.testResults.errorDetails') }}:</span>
                <p class="text-error-heading ml-0 mt-1">{{ vectorTestResult.details }}</p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </CyberDialog>

    <CyberDialog
      v-model="showRegenerateDialog"
      width="500px"
      :title="$t('admin.ai.dialogs.regenerateVectors')"
      :show-default-footer="false"
      :show-footer="false"
    >
      <div class="space-y-4">
        <div class="status-warning rounded-lg border p-4">
          <div class="flex items-start space-x-3">
            <i class="fas fa-exclamation-triangle mt-0.5 text-warning" />
            <div>
              <h4 class="text-warning-heading font-medium">{{ $t('admin.ai.dialogs.regenerateWarning') }}</h4>
              <p class="text-muted mt-2 text-sm">
                <strong class="text-error">{{ $t('admin.ai.dialogs.regenerateDestructive') }}</strong>
              </p>
              <ul class="text-muted mt-2 space-y-1 text-sm">
                <li>• {{ $t('admin.ai.dialogs.regeneratePoint1') }}</li>
                <li>• {{ $t('admin.ai.dialogs.regeneratePoint2') }}</li>
                <li>• {{ $t('admin.ai.dialogs.regeneratePoint3') }}</li>
                <li>• {{ $t('admin.ai.dialogs.regeneratePoint4') }}</li>
              </ul>
            </div>
          </div>
        </div>

        <div class="flex justify-end space-x-3">
          <CyberButton type="outlined" @click="showRegenerateDialog = false"> {{ $t('admin.ai.actions.cancel') }} </CyberButton>
          <CyberButton
            type="danger"
            :loading="regeneratingVectors"
            loading-mode="inline"
            icon="exclamation-triangle"
            @click="regenerateVectors"
          >
            {{ regeneratingVectors ? $t('admin.ai.actions.regenerating') : $t('admin.ai.actions.confirm') }}
          </CyberButton>
        </div>
      </div>
    </CyberDialog>

    <CyberDialog
      v-model="showAIFeaturePromptDialog"
      width="600px"
      :title="$t('admin.ai.dialogs.aiFeatureSuccess')"
      :show-default-footer="false"
      :show-footer="false"
    >
      <div class="space-y-4">
        <div class="status-success rounded-lg border p-4">
          <div class="flex items-start space-x-3">
            <i class="fas fa-check-circle mt-0.5 text-success" />
            <div>
              <h4 class="text-success-heading mb-2 font-medium">{{ $t('admin.ai.dialogs.aiConfigComplete') }}</h4>
              <p class="text-muted text-sm">{{ $t('admin.ai.dialogs.aiConfigCompleteDesc') }}</p>
            </div>
          </div>
        </div>

        <div class="space-y-3">
          <div
            v-if="aiFeaturePromptData?.needContentDetection"
            class="bg-background-600-bg flex items-center space-x-3 rounded-lg border border-default p-3"
          >
            <i class="fas fa-shield-alt text-warning" />
            <div class="flex-1">
              <h5 class="font-medium text-content-heading">{{ $t('admin.ai.contentDetection.contentDetectionFeature') }}</h5>
              <p class="text-muted text-sm">{{ $t('admin.ai.contentDetection.contentDetectionFeatureDesc') }}</p>
            </div>
            <div class="text-sm text-error">{{ $t('admin.ai.dialogs.featureNotEnabled') }}</div>
          </div>

          <div
            v-if="aiFeaturePromptData?.needAIAnalysis"
            class="bg-background-600-bg flex items-center space-x-3 rounded-lg border border-default p-3"
          >
            <i class="fas fa-brain text-info" />
            <div class="flex-1">
              <h5 class="font-medium text-content-heading">{{ $t('admin.ai.contentDetection.aiAnalysisFeature') }}</h5>
              <p class="text-muted text-sm">{{ $t('admin.ai.contentDetection.aiAnalysisFeatureDesc') }}</p>
            </div>
            <div class="text-sm text-error">{{ $t('admin.ai.dialogs.featureNotEnabled') }}</div>
          </div>
        </div>

        <div class="status-info rounded-lg border p-4">
          <div class="flex items-start space-x-3">
            <i class="fas fa-info-circle mt-0.5 text-info" />
            <div>
              <h4 class="text-info-heading mb-2 font-medium">{{ $t('admin.ai.dialogs.featureExplanation') }}</h4>
              <p class="text-muted text-sm">
                {{ $t('admin.ai.dialogs.featureExplanationDesc') }}
              </p>
            </div>
          </div>
        </div>

        <div class="flex justify-end space-x-3">
          <CyberButton type="outlined" @click="showAIFeaturePromptDialog = false">
            {{ $t('admin.ai.actions.later') }}
          </CyberButton>
          <CyberButton type="primary" icon="cog" @click="goToContentDetectionSettings">
            {{ $t('admin.ai.actions.goToSettings') }}
          </CyberButton>
        </div>
      </div>
    </CyberDialog>
  </div>
</template>

<style scoped lang="scss">
  .admin-ai-page {
    color: var(--color-content);
  }

  .content-panel {
    opacity: 1;
    transform: translateY(0);
    transition: all var(--transition-normal) var(--ease-out);
    display: flex;
    flex-direction: column;
    flex: 1;
    min-height: 0;
    overflow-y: auto;
    overflow-x: hidden;
  }

  .content-panel::-webkit-scrollbar {
    width: var(--space-sm);
    background: transparent;
  }

  .content-panel::-webkit-scrollbar-track {
    background: rgba(var(--color-brand-500-rgb), 0.1);
    border-radius: var(--radius-sm);
  }

  .content-panel::-webkit-scrollbar-thumb {
    background: rgba(var(--color-brand-500-rgb), 0.3);
    border-radius: var(--radius-sm);
    transition: all var(--transition-normal);
  }

  .content-panel::-webkit-scrollbar-thumb:hover {
    background: rgba(var(--color-brand-500-rgb), 0.5);
  }

  .status-warning {
    border-color: rgba(var(--color-warning-rgb), 0.3);
    background: rgba(var(--color-warning-rgb), 0.1);
  }

  .status-success {
    border-color: rgba(var(--color-success-rgb), 0.3);
    background: rgba(var(--color-success-rgb), 0.1);
  }

  .status-error {
    border-color: rgba(var(--color-error-rgb), 0.3);
    background: rgba(var(--color-error-rgb), 0.1);
  }

  .status-info {
    border-color: rgba(var(--color-brand-500-rgb), 0.3);
    background: rgba(var(--color-brand-500-rgb), 0.1);
  }

  .text-warning {
    color: var(--color-warning-400);
  }

  .text-warning-heading {
    color: var(--color-warning-300);
  }

  .text-success {
    color: var(--color-success-400);
  }

  .text-success-heading {
    color: var(--color-success-300);
  }

  .text-error {
    color: var(--color-error-400);
  }

  .text-error-heading {
    color: var(--color-error-300);
  }

  .text-info {
    color: var(--color-brand-500);
  }

  .text-info-heading {
    color: var(--color-brand-600);
  }

  .text-muted {
    color: var(--color-content-muted);
  }

  .text-subtle {
    color: var(--color-content-subtle);
  }

  .icon-muted {
    opacity: var(--opacity-60);
  }
</style>
