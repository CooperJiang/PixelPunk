<script setup lang="ts">
  import { ref, watch, onMounted } from 'vue'
  import GithubOAuthSettings from './GithubOAuthSettings.vue'
  import GoogleOAuthSettings from './GoogleOAuthSettings.vue'
  import LinuxdoOAuthSettings from './LinuxdoOAuthSettings.vue'
  import ProxySettings from './ProxySettings.vue'
  import type { Setting } from '@/api/admin/settings'
  import { useTexts } from '@/composables/useTexts'

  const { $t } = useTexts()

  const props = defineProps<{
    settings?: Setting[]
  }>()

  const emit = defineEmits<{
    (e: 'update', settings: Setting[]): void
  }>()

  const activeTab = ref('github')
  const githubSettings = ref<Setting[]>([])
  const googleSettings = ref<Setting[]>([])
  const proxySettings = ref<Setting[]>([])
  const linuxdoSettings = ref<Setting[]>([])

  const filterSettingsByProvider = (settings: Setting[], provider: string) => {
    if (provider === 'proxy') {
      return settings.filter((s) => s.key.startsWith('oauth_proxy_'))
    }
    return settings.filter((s) => s.key.startsWith(`${provider}_oauth_`))
  }

  const applySettings = (settings: Setting[]) => {
    githubSettings.value = filterSettingsByProvider(settings, 'github')
    googleSettings.value = filterSettingsByProvider(settings, 'google')
    proxySettings.value = filterSettingsByProvider(settings, 'proxy')
    linuxdoSettings.value = filterSettingsByProvider(settings, 'linuxdo')
  }

  const handleGithubUpdate = (settings: Setting[]) => {
    githubSettings.value = settings
    emitAllSettings()
  }

  const handleGoogleUpdate = (settings: Setting[]) => {
    googleSettings.value = settings
    emitAllSettings()
  }

  const handleProxyUpdate = (settings: Setting[]) => {
    proxySettings.value = settings
    emitAllSettings()
  }

  const handleLinuxdoUpdate = (settings: Setting[]) => {
    linuxdoSettings.value = settings
    emitAllSettings()
  }

  const emitAllSettings = () => {
    const allSettings = [...githubSettings.value, ...googleSettings.value, ...proxySettings.value, ...linuxdoSettings.value]
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
  <div class="oauth-settings">
    <CyberTabs v-model="activeTab">
      <CyberTabPane name="github" label="GitHub">
        <GithubOAuthSettings :settings="githubSettings" @update="handleGithubUpdate" />
      </CyberTabPane>

      <CyberTabPane name="google" label="Google">
        <GoogleOAuthSettings :settings="googleSettings" @update="handleGoogleUpdate" />
      </CyberTabPane>

      <CyberTabPane name="linuxdo" label="Linux.DO">
        <LinuxdoOAuthSettings :settings="linuxdoSettings" @update="handleLinuxdoUpdate" />
      </CyberTabPane>

      <CyberTabPane name="proxy" :label="$t('admin.settings.oauth.proxyConfig')">
        <ProxySettings :settings="proxySettings" @update="handleProxyUpdate" />
      </CyberTabPane>
    </CyberTabs>
  </div>
</template>

<style scoped>
  .oauth-settings {
    width: 100%;
  }

  .oauth-settings :deep(.cyber-tabs) {
    gap: 16px;
  }

  .oauth-settings :deep(.tab-item) {
    padding: 8px 20px;
  }

  .oauth-settings :deep(.tab-label) {
    font-size: 13px;
    font-weight: 500;
  }

  .oauth-settings :deep(.tabs-header) {
    padding-bottom: 6px;
  }

  .oauth-settings :deep(.space-y-5) {
    row-gap: 16px;
  }
</style>
