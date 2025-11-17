<script setup lang="ts">
  import { reactive, watch, onMounted } from 'vue'
  import type { Setting } from '@/api/admin/settings'
  import { useToast } from '@/components/Toast/useToast'
  import { useTexts } from '@/composables/useTexts'

  const { $t } = useTexts()

  const props = defineProps<{
    settings: Setting[]
  }>()

  const emit = defineEmits<{
    (e: 'update', settings: Setting[]): void
  }>()

  const toast = useToast()

  const localSettings = reactive({
    google_oauth_enabled: false,
    google_oauth_client_id: '',
    google_oauth_client_secret: '',
    google_oauth_redirect_uri: '',
    google_oauth_scope: 'openid email profile',
    google_oauth_proxy_enabled: false,
  })

  const getSettingsArray = (): Setting[] => {
    return [
      {
        key: 'google_oauth_enabled',
        value: localSettings.google_oauth_enabled,
        type: 'boolean',
        group: 'oauth',
        description: $t('admin.settings.oauth.google.desc.enabled'),
      },
      {
        key: 'google_oauth_client_id',
        value: localSettings.google_oauth_client_id,
        type: 'string',
        group: 'oauth',
        description: $t('admin.settings.oauth.google.desc.clientId'),
      },
      {
        key: 'google_oauth_client_secret',
        value: localSettings.google_oauth_client_secret,
        type: 'string',
        group: 'oauth',
        description: $t('admin.settings.oauth.google.desc.clientSecret'),
      },
      {
        key: 'google_oauth_redirect_uri',
        value: localSettings.google_oauth_redirect_uri,
        type: 'string',
        group: 'oauth',
        description: $t('admin.settings.oauth.google.desc.redirectUri'),
      },
      {
        key: 'google_oauth_scope',
        value: 'openid email profile', // 固定值，始终使用 openid email profile
        type: 'string',
        group: 'oauth',
        description: $t('admin.settings.oauth.google.desc.scope'),
      },
      {
        key: 'google_oauth_proxy_enabled',
        value: localSettings.google_oauth_proxy_enabled,
        type: 'boolean',
        group: 'oauth',
        description: $t('admin.settings.oauth.google.desc.proxyEnabled'),
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

    if (!localSettings.google_oauth_redirect_uri) {
      const baseUrl = window.location.origin
      localSettings.google_oauth_redirect_uri = `${baseUrl}/auth`
    }

    localSettings.google_oauth_scope = 'openid email profile'
  }

  const copyToClipboard = async (text: string) => {
    try {
      await navigator.clipboard.writeText(text)
      toast.success($t('admin.settings.oauth.toast.copied'))
    } catch (_err) {
      toast.error($t('admin.settings.oauth.toast.copyFailed'))
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
        <i class="fab fa-google mr-2 text-brand-500" />{{ $t('admin.settings.oauth.google.enabled.labelFull') }}
      </label>
      <div class="flex-1">
        <CyberSwitch v-model="localSettings.google_oauth_enabled" />
      </div>
    </div>

    <div class="flex flex-col md:flex-row md:items-start">
      <label class="flex w-full items-center text-sm text-content md:w-48 md:py-2">
        <i class="fas fa-key mr-2 text-brand-500" />{{ $t('admin.settings.oauth.google.clientId.label') }}
      </label>
      <div class="flex-1">
        <CyberInput
          v-model="localSettings.google_oauth_client_id"
          :placeholder="$t('admin.settings.oauth.google.clientId.placeholderFull')"
          width="500px"
        />
        <div class="text-content-subtle mt-1 text-xs">{{ $t('admin.settings.oauth.google.clientId.descriptionFull') }}</div>
      </div>
    </div>

    <div class="flex flex-col md:flex-row md:items-start">
      <label class="flex w-full items-center text-sm text-content md:w-48 md:py-2">
        <i class="fas fa-lock mr-2 text-brand-500" />{{ $t('admin.settings.oauth.google.clientSecret.label') }}
      </label>
      <div class="flex-1">
        <CyberInput
          v-model="localSettings.google_oauth_client_secret"
          type="password"
          :placeholder="$t('admin.settings.oauth.google.clientSecret.placeholderFull')"
          width="500px"
        />
        <div class="text-content-subtle mt-1 text-xs">{{ $t('admin.settings.oauth.google.clientSecret.descriptionFull') }}</div>
      </div>
    </div>

    <div class="flex flex-col md:flex-row md:items-start">
      <label class="flex w-full items-center text-sm text-content md:w-48 md:py-2">
        <i class="fas fa-link mr-2 text-brand-500" />{{ $t('admin.settings.oauth.google.callbackUrl.label') }}
      </label>
      <div class="flex-1">
        <div class="flex items-center gap-2">
          <CyberInput
            v-model="localSettings.google_oauth_redirect_uri"
            :placeholder="$t('admin.settings.oauth.google.callbackUrl.placeholderSimple')"
            width="500px"
          />
          <button
            class="rounded bg-background-700 px-3 py-2 text-sm text-content transition-colors hover:bg-background-600"
            @click="copyToClipboard(localSettings.google_oauth_redirect_uri)"
          >
            <i class="fas fa-copy mr-1" />{{ $t('actions.copy') }}
          </button>
        </div>
        <div class="text-content-subtle mt-1 text-xs">
          {{ $t('admin.settings.oauth.common.callback.format') }}
        </div>
      </div>
    </div>

    <div class="flex flex-col md:flex-row md:items-start">
      <label class="flex w-full items-center text-sm text-content md:w-48 md:py-2">
        <i class="fas fa-shield-alt mr-2 text-brand-500" />{{ $t('admin.settings.oauth.common.scopes.label') }}
      </label>
      <div class="flex-1">
        <div class="flex items-center gap-2">
          <code class="rounded bg-background-800 px-3 py-2 text-sm text-brand-400" style="font-family: 'Courier New', monospace">
            openid email profile
          </code>
        </div>
        <div class="text-content-subtle mt-1 text-xs">{{ $t('admin.settings.oauth.google.scopes.fixed') }}</div>
      </div>
    </div>

    <div class="border-brand-500/20 my-4 border-t"></div>

    <div class="flex flex-col md:flex-row md:items-center">
      <label class="flex w-full items-center text-sm text-content md:w-48 md:py-2">
        <i class="fas fa-network-wired mr-2 text-brand-500" />{{ $t('admin.settings.oauth.common.proxy.label') }}
      </label>
      <div class="flex-1">
        <CyberSwitch v-model="localSettings.google_oauth_proxy_enabled" />
        <div class="text-content-subtle mt-1 text-xs">
          {{ $t('admin.settings.oauth.common.proxy.googleDescription') }}
        </div>
      </div>
    </div>

    <div class="border-brand-500/20 bg-background-800/50 mt-4 rounded-lg border p-3">
      <div class="mb-2 flex items-center gap-2 text-xs font-medium text-brand-400">
        <i class="fas fa-info-circle text-xs" />
        <span>{{ $t('admin.settings.oauth.common.setup.title') }}</span>
      </div>
      <ol class="text-content-subtle space-y-1.5 text-xs">
        <li class="flex gap-2">
          <span class="text-brand-500">1.</span>
          <span
            >{{ $t('admin.settings.oauth.google.setup.step1') }}
            <a
              href="https://console.cloud.google.com/apis/credentials"
              target="_blank"
              class="text-brand-400 hover:text-brand-300"
              >{{ $t('admin.settings.oauth.google.setup.step1Link') }}</a
            ></span
          >
        </li>
        <li class="flex gap-2">
          <span class="text-brand-500">2.</span>
          <span>{{ $t('admin.settings.oauth.google.setup.step2') }}</span>
        </li>
        <li class="flex gap-2">
          <span class="text-brand-500">3.</span>
          <div class="flex-1">
            <div>
              {{ $t('admin.settings.oauth.google.setup.step3') }}
              <span class="text-brand-400">{{ $t('admin.settings.oauth.google.setup.step3Field') }}</span>
              {{ $t('admin.settings.oauth.google.setup.step3Suffix') }}
            </div>
            <div class="mt-1 rounded bg-background-900 px-2 py-1 font-mono text-brand-300">
              {{ $t('admin.settings.oauth.common.callback.exampleUrl') }}
            </div>
            <div class="mt-1 text-[11px] text-content-muted">{{ $t('admin.settings.oauth.common.callback.examplePrefix') }}</div>
          </div>
        </li>
        <li class="flex gap-2">
          <span class="text-brand-500">4.</span>
          <span
            >{{ $t('admin.settings.oauth.google.setup.step4') }}
            <span class="text-brand-400">{{ $t('admin.settings.oauth.google.setup.step4Field1') }}</span>
            {{ $t('admin.settings.oauth.google.setup.step4And') }}
            <span class="text-brand-400">{{ $t('admin.settings.oauth.google.setup.step4Field2') }}</span>
            {{ $t('admin.settings.oauth.google.setup.step4Suffix') }}</span
          >
        </li>
        <li class="flex gap-2">
          <span class="text-brand-500">5.</span>
          <span>{{ $t('admin.settings.oauth.google.setup.step5') }}</span>
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
