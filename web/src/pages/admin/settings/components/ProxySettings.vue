<script setup lang="ts">
  import { reactive, watch, onMounted, ref } from 'vue'
  import type { Setting } from '@/api/admin/settings'
  import { useToast } from '@/components/Toast/useToast'
  import { testProxy } from '@/api/admin/settings/oauth'
  import { useTexts } from '@/composables/useTexts'
  import Validator, { getValidationRules } from '@/utils/validation/validator'

  const { $t } = useTexts()
  const R = getValidationRules($t)

  const props = defineProps<{
    settings: Setting[]
  }>()

  const emit = defineEmits<{
    (e: 'update', settings: Setting[]): void
  }>()

  const toast = useToast()

  const localSettings = reactive({
    oauth_proxy_dynamic: false,
    oauth_proxy_api_url: '',
    oauth_proxy_type: 'http',
    oauth_proxy_host: '',
    oauth_proxy_port: '',
    oauth_proxy_username: '',
    oauth_proxy_password: '',
  })

  const proxyTypeOptions = [
    { value: 'http', label: 'HTTP' },
    { value: 'https', label: 'HTTPS' },
    { value: 'socks5', label: 'SOCKS5' },
  ]

  const isTesting = ref(false)
  const testResult = ref<{ success: boolean; message: string } | null>(null)

  const getSettingsArray = (): Setting[] => {
    return [
      {
        key: 'oauth_proxy_dynamic',
        value: localSettings.oauth_proxy_dynamic,
        type: 'boolean',
        group: 'oauth',
        description: $t('admin.settings.proxy.desc.dynamic'),
      },
      {
        key: 'oauth_proxy_api_url',
        value: localSettings.oauth_proxy_api_url,
        type: 'string',
        group: 'oauth',
        description: $t('admin.settings.proxy.desc.apiUrl'),
      },
      {
        key: 'oauth_proxy_type',
        value: localSettings.oauth_proxy_type,
        type: 'string',
        group: 'oauth',
        description: $t('admin.settings.proxy.desc.type'),
      },
      {
        key: 'oauth_proxy_host',
        value: localSettings.oauth_proxy_host,
        type: 'string',
        group: 'oauth',
        description: $t('admin.settings.proxy.desc.host'),
      },
      {
        key: 'oauth_proxy_port',
        value: localSettings.oauth_proxy_port,
        type: 'string',
        group: 'oauth',
        description: $t('admin.settings.proxy.desc.port'),
      },
      {
        key: 'oauth_proxy_username',
        value: localSettings.oauth_proxy_username,
        type: 'string',
        group: 'oauth',
        description: $t('admin.settings.proxy.desc.username'),
      },
      {
        key: 'oauth_proxy_password',
        value: localSettings.oauth_proxy_password,
        type: 'string',
        group: 'oauth',
        description: $t('admin.settings.proxy.desc.password'),
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

  const handleTestProxy = async () => {
    if (localSettings.oauth_proxy_dynamic) {
      const apiUrlCheck = Validator.validate(localSettings.oauth_proxy_api_url || '', [R.required], $t)
      if (!apiUrlCheck.valid) return toast.error($t('admin.settings.proxy.toast.apiUrlRequired'))
      const userCheck = Validator.validate(localSettings.oauth_proxy_username || '', [R.required], $t)
      const passCheck = Validator.validate(localSettings.oauth_proxy_password || '', [R.required], $t)
      if (!userCheck.valid || !passCheck.valid) return toast.error($t('admin.settings.proxy.toast.credentialsRequired'))
    } else {
      const hostCheck = Validator.validate(localSettings.oauth_proxy_host || '', [R.required], $t)
      const portCheck = Validator.validate(localSettings.oauth_proxy_port || '', [R.required], $t)
      if (!hostCheck.valid || !portCheck.valid) return toast.error($t('admin.settings.proxy.toast.hostPortRequired'))
    }

    isTesting.value = true
    testResult.value = null

    try {
      const result = await testProxy({
        proxy_dynamic: localSettings.oauth_proxy_dynamic,
        proxy_api_url: localSettings.oauth_proxy_api_url || undefined,
        proxy_type: localSettings.oauth_proxy_type,
        proxy_host: localSettings.oauth_proxy_host,
        proxy_port: localSettings.oauth_proxy_port,
        proxy_username: localSettings.oauth_proxy_username || undefined,
        proxy_password: localSettings.oauth_proxy_password || undefined,
      })

      if (result.success && result.data) {
        testResult.value = {
          success: result.data.success,
          message: result.data.message,
        }

        if (result.data.success) {
          if (result.data.latency) {
            toast.success(
              $t('admin.settings.proxy.toast.testSuccessWithLatency').replace('{latency}', result.data.latency.toString())
            )
          } else {
            toast.success($t('admin.settings.proxy.toast.testSuccess'))
          }
        } else {
          toast.error($t('admin.settings.proxy.toast.testFailed').replace('{message}', result.data.message))
        }
      } else {
        testResult.value = {
          success: false,
          message: result.message || $t('admin.settings.proxy.toast.testRequestFailed'),
        }
        toast.error(result.message || $t('admin.settings.proxy.toast.testRequestFailed'))
      }
    } catch (error: any) {
      const errorMsg = error?.response?.data?.message || error?.message || $t('admin.settings.proxy.toast.networkRequestFailed')
      testResult.value = {
        success: false,
        message: errorMsg,
      }
      toast.error($t('admin.settings.proxy.toast.testError').replace('{message}', errorMsg))
    } finally {
      isTesting.value = false
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
    <div class="border-brand-500/20 bg-background-800/50 rounded-lg border p-3">
      <div class="mb-2 flex items-center gap-2 text-xs font-medium text-brand-400">
        <i class="fas fa-info-circle text-xs" />
        <span>{{ $t('admin.settings.proxy.info.title') }}</span>
      </div>
      <p class="text-content-subtle text-xs">
        {{ $t('admin.settings.proxy.info.desc1') }}
        <br />
        {{ $t('admin.settings.proxy.info.desc2') }}
      </p>
    </div>

    <div class="flex flex-col md:flex-row md:items-center">
      <label class="flex w-full items-center text-sm text-content md:w-48 md:py-2">
        <i class="fas fa-sync-alt mr-2 text-brand-500" />{{ $t('admin.settings.proxy.dynamic.label') }}
      </label>
      <div class="flex-1">
        <CyberSwitch v-model="localSettings.oauth_proxy_dynamic" />
        <div class="text-content-subtle mt-1 text-xs">{{ $t('admin.settings.proxy.dynamic.hint') }}</div>
      </div>
    </div>

    <div v-if="localSettings.oauth_proxy_dynamic" class="flex flex-col md:flex-row md:items-start">
      <label class="flex w-full items-center text-sm text-content md:w-48 md:py-2">
        <i class="fas fa-link mr-2 text-brand-500" />{{ $t('admin.settings.proxy.apiUrl.label') }}
      </label>
      <div class="flex-1">
        <CyberInput
          v-model="localSettings.oauth_proxy_api_url"
          :placeholder="$t('admin.settings.proxy.apiUrl.placeholder')"
          width="500px"
        />
        <div class="text-content-subtle mt-1 text-xs">{{ $t('admin.settings.proxy.apiUrl.hint') }}</div>
      </div>
    </div>

    <div v-if="!localSettings.oauth_proxy_dynamic" class="flex flex-col md:flex-row md:items-start">
      <label class="flex w-full items-center text-sm text-content md:w-48 md:py-2">
        <i class="fas fa-layer-group mr-2 text-brand-500" />{{ $t('admin.settings.proxy.type.label') }}
      </label>
      <div class="flex-1">
        <CyberDropdown
          v-model="localSettings.oauth_proxy_type"
          :options="proxyTypeOptions"
          :placeholder="$t('admin.settings.proxy.type.placeholder')"
          width="500px"
        />
        <div class="text-content-subtle mt-1 text-xs">{{ $t('admin.settings.proxy.type.hint') }}</div>
      </div>
    </div>

    <div v-if="!localSettings.oauth_proxy_dynamic" class="flex flex-col md:flex-row md:items-start">
      <label class="flex w-full items-center text-sm text-content md:w-48 md:py-2">
        <i class="fas fa-server mr-2 text-brand-500" />{{ $t('admin.settings.proxy.host.label') }}
      </label>
      <div class="flex-1">
        <CyberInput
          v-model="localSettings.oauth_proxy_host"
          :placeholder="$t('admin.settings.proxy.host.placeholder')"
          width="500px"
        />
        <div class="text-content-subtle mt-1 text-xs">{{ $t('admin.settings.proxy.host.hint') }}</div>
      </div>
    </div>

    <div v-if="!localSettings.oauth_proxy_dynamic" class="flex flex-col md:flex-row md:items-start">
      <label class="flex w-full items-center text-sm text-content md:w-48 md:py-2">
        <i class="fas fa-plug mr-2 text-brand-500" />{{ $t('admin.settings.proxy.port.label') }}
      </label>
      <div class="flex-1">
        <CyberInput
          v-model="localSettings.oauth_proxy_port"
          :placeholder="$t('admin.settings.proxy.port.placeholder')"
          width="500px"
        />
        <div class="text-content-subtle mt-1 text-xs">{{ $t('admin.settings.proxy.port.hint') }}</div>
      </div>
    </div>

    <div class="flex flex-col md:flex-row md:items-start">
      <label class="flex w-full items-center text-sm text-content md:w-48 md:py-2">
        <i class="fas fa-user mr-2 text-brand-500" />{{ $t('admin.settings.proxy.username.label') }}
      </label>
      <div class="flex-1">
        <CyberInput
          v-model="localSettings.oauth_proxy_username"
          :placeholder="$t('admin.settings.proxy.username.placeholder')"
          width="500px"
        />
        <div class="text-content-subtle mt-1 text-xs">{{ $t('admin.settings.proxy.username.hint') }}</div>
      </div>
    </div>

    <div class="flex flex-col md:flex-row md:items-start">
      <label class="flex w-full items-center text-sm text-content md:w-48 md:py-2">
        <i class="fas fa-key mr-2 text-brand-500" />{{ $t('admin.settings.proxy.password.label') }}
      </label>
      <div class="flex-1">
        <CyberInput
          v-model="localSettings.oauth_proxy_password"
          type="password"
          :placeholder="$t('admin.settings.proxy.password.placeholder')"
          width="500px"
        />
        <div class="text-content-subtle mt-1 text-xs">{{ $t('admin.settings.proxy.password.hint') }}</div>
      </div>
    </div>

    <div class="flex flex-col md:flex-row md:items-start">
      <label class="flex w-full items-center text-sm text-content md:w-48 md:py-2">
        <i class="fas fa-vial mr-2 text-brand-500" />{{ $t('admin.settings.proxy.test.label') }}
      </label>
      <div class="flex-1">
        <CyberButton type="primary" icon="flask" :loading="isTesting" :disabled="isTesting" @click="handleTestProxy">
          {{ isTesting ? $t('admin.settings.proxy.test.testing') : $t('admin.settings.proxy.test.button') }}
        </CyberButton>

        <div class="text-content-subtle mt-1 text-xs">{{ $t('admin.settings.proxy.test.hint') }}</div>

        <div v-if="testResult && !isTesting" class="mt-2 flex items-center gap-2 text-xs">
          <i v-if="testResult.success" class="fas fa-check-circle text-success" />
          <i v-else class="fas fa-times-circle text-error-500" />
          <span :class="testResult.success ? 'text-success' : 'text-error-500'">
            {{ testResult.message }}
          </span>
        </div>
      </div>
    </div>

    <div class="border-brand-500/20 bg-background-800/50 mt-4 rounded-lg border p-3">
      <div class="mb-2 flex items-center gap-2 text-xs font-medium text-brand-400">
        <i class="fas fa-book mr-1 text-xs" />
        <span>{{ $t('admin.settings.proxy.guide.title') }}</span>
      </div>
      <ol class="text-content-subtle space-y-1.5 text-xs">
        <li class="flex gap-2">
          <span class="text-brand-500">1.</span>
          <span v-html="$t('admin.settings.proxy.guide.step1')"></span>
        </li>
        <li class="flex gap-2">
          <span class="text-brand-500">2.</span>
          <span v-html="$t('admin.settings.proxy.guide.step2')"></span>
        </li>
        <li class="flex gap-2">
          <span class="text-brand-500">3.</span>
          <span>{{ $t('admin.settings.proxy.guide.step3') }}</span>
        </li>
        <li class="flex gap-2">
          <span class="text-brand-500">4.</span>
          <span>{{ $t('admin.settings.proxy.guide.step4') }}</span>
        </li>
        <li class="flex gap-2">
          <span class="text-brand-500">5.</span>
          <span>{{ $t('admin.settings.proxy.guide.step5') }}</span>
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
