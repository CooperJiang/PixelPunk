<script setup lang="ts">
  import { computed, onMounted, ref, watch } from 'vue'
  import { useToast } from '@/components/Toast/useToast'
  import { useRoute, useRouter } from 'vue-router'
  import { useTexts } from '@/composables/useTexts'
  import WebsiteSettings from './components/WebsiteSettings.vue'
  import RegistrationSettings from './components/RegistrationSettings.vue'
  import MailSettings from './components/MailSettings.vue'
  import SecuritySettings from './components/SecuritySettings.vue'
  import UploadSettings from './components/UploadSettings.vue'
  import GuestSettings from './components/GuestSettings.vue'
  import OAuthSettings from './components/OAuthSettings.vue'
  import AnalyticsSettings from './components/AnalyticsSettings.vue'
  import { batchUpsertSettings, defaultSettings, getSettings, type Setting } from '@/api/admin/settings'

  /* 工具 */
  const toast = useToast()
  const route = useRoute()
  const router = useRouter()
  const { $t } = useTexts()

  const isLoading = ref(false)
  const isSaving = ref(false)

  interface TabItem {
    key: string
    name: string
    icon?: string
  }

  /* 选项卡数据 */
  const tabs = computed<TabItem[]>(() => [
    { key: 'website', name: $t('admin.settings.tabs.website'), icon: 'fas fa-globe' },
    { key: 'registration', name: $t('admin.settings.tabs.registration'), icon: 'fas fa-user-plus' },
    { key: 'mail', name: $t('admin.settings.tabs.mail'), icon: 'fas fa-envelope' },
    { key: 'upload', name: $t('admin.settings.tabs.upload'), icon: 'fas fa-upload' },
    { key: 'security', name: $t('admin.settings.tabs.security'), icon: 'fas fa-shield-alt' },
    { key: 'guest', name: $t('admin.settings.tabs.guest'), icon: 'fas fa-user-secret' },
    { key: 'oauth', name: $t('admin.settings.tabs.oauth'), icon: 'fas fa-plug' },
    { key: 'analytics', name: $t('admin.settings.tabs.analytics'), icon: 'fas fa-chart-line' },
  ])

  /* 从查询参数获取初始tab，如果没有则默认为'website' */
  const activeTab = ref(route.query.tab?.toString() || 'website')

  const updateActiveTab = (tabId: string) => {
    activeTab.value = tabId
    router.push({
      name: route.name,
      query: { ...route.query, tab: tabId },
    })
  }

  const websiteSettings = ref<Setting[]>([])
  const registrationSettings = ref<Setting[]>([])
  const mailSettings = ref<Setting[]>([])
  const securitySettings = ref<Setting[]>([])
  const uploadSettings = ref<Setting[]>([])
  const guestSettings = ref<Setting[]>([])
  const oauthSettings = ref<Setting[]>([])
  const analyticsSettings = ref<Setting[]>([])

  const modifiedWebsiteSettings = ref<Setting[]>([])
  const modifiedRegistrationSettings = ref<Setting[]>([])
  const modifiedMailSettings = ref<Setting[]>([])
  const modifiedSecuritySettings = ref<Setting[]>([])
  const modifiedUploadSettings = ref<Setting[]>([])
  const modifiedGuestSettings = ref<Setting[]>([])
  const modifiedOAuthSettings = ref<Setting[]>([])
  const modifiedAnalyticsSettings = ref<Setting[]>([])

  const handleWebsiteSettingsUpdate = (updatedSettings: Setting[]) => {
    modifiedWebsiteSettings.value = updatedSettings
  }

  const handleRegistrationSettingsUpdate = (updatedSettings: Setting[]) => {
    modifiedRegistrationSettings.value = updatedSettings
  }

  const handleMailSettingsUpdate = (updatedSettings: Setting[]) => {
    modifiedMailSettings.value = updatedSettings
  }

  const handleSecuritySettingsUpdate = (updatedSettings: Setting[]) => {
    modifiedSecuritySettings.value = updatedSettings
  }

  const handleUploadSettingsUpdate = (updatedSettings: Setting[]) => {
    modifiedUploadSettings.value = updatedSettings
  }

  const handleGuestSettingsUpdate = (updatedSettings: Setting[]) => {
    modifiedGuestSettings.value = updatedSettings
  }

  const handleOAuthSettingsUpdate = (updatedSettings: Setting[]) => {
    modifiedOAuthSettings.value = updatedSettings
  }

  const handleAnalyticsSettingsUpdate = (updatedSettings: Setting[]) => {
    modifiedAnalyticsSettings.value = updatedSettings
  }

  const resetSettings = async () => {
    isLoading.value = true
    await loadSettingsForTab(activeTab.value as SettingGroup)

    if (activeTab.value === 'website') {
      modifiedWebsiteSettings.value = []
    } else if (activeTab.value === 'registration') {
      modifiedRegistrationSettings.value = []
    } else if (activeTab.value === 'mail') {
      modifiedMailSettings.value = []
    } else if (activeTab.value === 'security') {
      modifiedSecuritySettings.value = []
    } else if (activeTab.value === 'upload') {
      modifiedUploadSettings.value = []
    } else if (activeTab.value === 'guest') {
      modifiedGuestSettings.value = []
    } else if (activeTab.value === 'oauth') {
      modifiedOAuthSettings.value = []
    } else if (activeTab.value === 'analytics') {
      modifiedAnalyticsSettings.value = []
    }

    toast.success($t('admin.settings.toast.reset'))
    isLoading.value = false
  }

  const saveSettings = async () => {
    let settingsToUpdate: Setting[] = []
    if (activeTab.value === 'website' && modifiedWebsiteSettings.value.length > 0) {
      settingsToUpdate = modifiedWebsiteSettings.value
    } else if (activeTab.value === 'registration' && modifiedRegistrationSettings.value.length > 0) {
      settingsToUpdate = modifiedRegistrationSettings.value
    } else if (activeTab.value === 'mail' && modifiedMailSettings.value.length > 0) {
      settingsToUpdate = modifiedMailSettings.value
    } else if (activeTab.value === 'security' && modifiedSecuritySettings.value.length > 0) {
      settingsToUpdate = modifiedSecuritySettings.value
    } else if (activeTab.value === 'upload' && modifiedUploadSettings.value.length > 0) {
      settingsToUpdate = modifiedUploadSettings.value
    } else if (activeTab.value === 'guest' && modifiedGuestSettings.value.length > 0) {
      settingsToUpdate = modifiedGuestSettings.value
    } else if (activeTab.value === 'oauth' && modifiedOAuthSettings.value.length > 0) {
      settingsToUpdate = modifiedOAuthSettings.value
    } else if (activeTab.value === 'analytics' && modifiedAnalyticsSettings.value.length > 0) {
      settingsToUpdate = modifiedAnalyticsSettings.value
    }

    if (settingsToUpdate.length === 0) {
      toast.info($t('admin.settings.toast.noChanges'))
      return
    }

    try {
      isSaving.value = true
      const result = await batchUpsertSettings(settingsToUpdate)

      if (result.success) {
        if (activeTab.value === 'website') {
          modifiedWebsiteSettings.value = []
        } else if (activeTab.value === 'registration') {
          modifiedRegistrationSettings.value = []
        } else if (activeTab.value === 'mail') {
          modifiedMailSettings.value = []
        } else if (activeTab.value === 'security') {
          modifiedSecuritySettings.value = []
        } else if (activeTab.value === 'upload') {
          modifiedUploadSettings.value = []
        } else if (activeTab.value === 'guest') {
          modifiedGuestSettings.value = []
        } else if (activeTab.value === 'oauth') {
          modifiedOAuthSettings.value = []
        } else if (activeTab.value === 'analytics') {
          modifiedAnalyticsSettings.value = []
        }
        toast.success($t('admin.settings.toast.saved'))
        await loadSettingsForTab(activeTab.value as SettingGroup)
      }
    } catch (error) {
      toast.error(error.message)
    } finally {
      isSaving.value = false
    }
  }

  const loadSettingsForTab = async (tabId: SettingGroup) => {
    isLoading.value = true

    try {
      const result = await getSettings({ group: tabId })

      const settings = result.success ? result.data?.settings || [] : []

      if (settings.length === 0) {
        if (tabId === 'website') {
          websiteSettings.value = defaultSettings.website || []
        } else if (tabId === 'registration') {
          registrationSettings.value = defaultSettings.registration || []
        } else if (tabId === 'mail') {
          mailSettings.value = defaultSettings.mail || []
        } else if (tabId === 'security') {
          securitySettings.value = defaultSettings.security || []
        } else if (tabId === 'upload') {
          uploadSettings.value = defaultSettings.upload || []
        } else if (tabId === 'guest') {
          guestSettings.value = defaultSettings.guest || []
        } else if (tabId === 'oauth') {
          oauthSettings.value = []
        } else if (tabId === 'analytics') {
          analyticsSettings.value = []
        }
      } else {
        if (tabId === 'website') {
          websiteSettings.value = settings
        } else if (tabId === 'registration') {
          registrationSettings.value = settings
        } else if (tabId === 'mail') {
          mailSettings.value = settings
        } else if (tabId === 'security') {
          securitySettings.value = settings
        } else if (tabId === 'upload') {
          uploadSettings.value = settings
        } else if (tabId === 'guest') {
          guestSettings.value = settings
        } else if (tabId === 'oauth') {
          oauthSettings.value = settings
        } else if (tabId === 'analytics') {
          analyticsSettings.value = settings
        }
      }
    } catch (error) {
      if (tabId === 'website') {
        websiteSettings.value = defaultSettings.website || []
      } else if (tabId === 'registration') {
        registrationSettings.value = defaultSettings.registration || []
      } else if (tabId === 'mail') {
        mailSettings.value = defaultSettings.mail || []
      } else if (tabId === 'security') {
        securitySettings.value = defaultSettings.security || []
      } else if (tabId === 'upload') {
        uploadSettings.value = defaultSettings.upload || []
      } else if (tabId === 'guest') {
        guestSettings.value = defaultSettings.guest || []
      } else if (tabId === 'oauth') {
        oauthSettings.value = []
      } else if (tabId === 'analytics') {
        analyticsSettings.value = []
      }
      toast.error(error.message)
    } finally {
      isLoading.value = false
    }
  }

  watch(activeTab, async (newTab) => {
    if (
      (newTab === 'website' && websiteSettings.value.length === 0) ||
      (newTab === 'registration' && registrationSettings.value.length === 0) ||
      (newTab === 'mail' && mailSettings.value.length === 0) ||
      (newTab === 'security' && securitySettings.value.length === 0) ||
      (newTab === 'upload' && uploadSettings.value.length === 0) ||
      (newTab === 'guest' && guestSettings.value.length === 0) ||
      (newTab === 'oauth' && oauthSettings.value.length === 0) ||
      (newTab === 'analytics' && analyticsSettings.value.length === 0)
    ) {
      await loadSettingsForTab(newTab as SettingGroup)
    }
  })

  watch(
    () => route.query.tab,
    (newTab) => {
      if (newTab && newTab !== activeTab.value) {
        activeTab.value = newTab.toString()
      }
    },
    { immediate: false }
  )

  onMounted(async () => {
    await loadSettingsForTab(activeTab.value as SettingGroup)
  })
</script>

<template>
  <div class="admin-settings-page admin-page-container">
    <CyberAdminWrapper
      :title="$t('admin.settings.title')"
      :subtitle="$t('admin.settings.subtitle')"
      icon="fas fa-cog"
      :sidebar-layout="true"
    >
      <template #actions>
        <CyberButton type="text" class="min-w-24" @click="resetSettings">
          {{ $t('admin.settings.buttons.reset') }}
        </CyberButton>
        <CyberButton type="primary" :loading="isSaving" class="min-w-24" @click="saveSettings">
          <i class="fas fa-save mr-2" /> {{ $t('admin.settings.buttons.save') }}
        </CyberButton>
      </template>

      <template #sidebar>
        <CyberSidebarNav :tabs="tabs" :active-tab="activeTab" @tab-change="updateActiveTab" />
      </template>

      <template #content>
        <div v-show="activeTab === 'website'" class="settings-section animate-fade-in">
          <WebsiteSettings :settings="websiteSettings" @update="handleWebsiteSettingsUpdate" />
        </div>

        <div v-show="activeTab === 'registration'" class="settings-section animate-fade-in">
          <RegistrationSettings :settings="registrationSettings" @update="handleRegistrationSettingsUpdate" />
        </div>

        <div v-show="activeTab === 'mail'" class="settings-section animate-fade-in">
          <MailSettings :settings="mailSettings" @update="handleMailSettingsUpdate" />
        </div>

        <div v-show="activeTab === 'security'" class="settings-section animate-fade-in">
          <SecuritySettings :settings="securitySettings" @update="handleSecuritySettingsUpdate" />
        </div>

        <div v-show="activeTab === 'upload'" class="settings-section animate-fade-in">
          <UploadSettings :settings="uploadSettings" @update="handleUploadSettingsUpdate" />
        </div>

        <div v-show="activeTab === 'guest'" class="settings-section animate-fade-in">
          <GuestSettings :settings="guestSettings" @update="handleGuestSettingsUpdate" />
        </div>

        <div v-show="activeTab === 'oauth'" class="settings-section animate-fade-in">
          <OAuthSettings :settings="oauthSettings" @update="handleOAuthSettingsUpdate" />
        </div>

        <div v-show="activeTab === 'analytics'" class="settings-section animate-fade-in">
          <AnalyticsSettings :settings="analyticsSettings" @update="handleAnalyticsSettingsUpdate" />
        </div>

        <CyberLoading v-if="isLoading" />
      </template>
    </CyberAdminWrapper>
  </div>
</template>

<style scoped>
  .admin-settings-page {
    color: var(--color-content);
  }

  .settings-section {
    display: flex;
    flex-direction: column;
  }

  :deep(.form-container) {
    background: linear-gradient(
      135deg,
      rgba(var(--color-background-800-rgb), 0.6) 0%,
      rgba(var(--color-background-900-rgb), 0.4) 100%
    );
    border-radius: var(--radius-sm);
    padding: 24px;
    margin-bottom: 24px;
    box-shadow:
      0 4px 20px rgba(0, 0, 0, 0.3),
      0 0 40px rgba(var(--color-brand-500-rgb), 0.02);
    backdrop-filter: blur(10px);
    transition: all 0.3s ease;
  }

  :deep(.form-container:hover) {
    box-shadow:
      0 8px 32px rgba(0, 0, 0, 0.4),
      0 0 60px rgba(var(--color-brand-500-rgb), 0.05);
  }

  :deep(.form-section) {
    margin-bottom: 32px;
    padding: 20px;
    background: rgba(var(--color-background-800-rgb), 0.3);
    border-radius: var(--radius-sm);
  }

  :deep(.form-section-title) {
    color: var(--color-brand-500);
    font-weight: 600;
    font-size: 1.1rem;
    margin-bottom: 16px;
    text-shadow: 0 0 10px rgba(var(--color-brand-500-rgb), 0.3);
    display: flex;
    align-items: center;
    gap: 8px;
  }

  :deep(.form-section-title i) {
    color: var(--color-brand-500);
    text-shadow: 0 0 10px rgba(var(--color-brand-500-rgb), 0.5);
  }

  :deep(.form-group) {
    margin-bottom: 24px;
  }

  :deep(.form-label) {
    color: var(--color-content-default);
    font-weight: 500;
    margin-bottom: 8px;
    display: block;
    font-size: 0.9rem;
  }

  :deep(.form-input),
  :deep(.form-select),
  :deep(.form-textarea) {
    width: 100%;
    padding: 12px 16px;
    background: rgba(var(--color-background-900-rgb), 0.8);
    border-radius: var(--radius-sm);
    color: var(--color-content);
    font-size: 0.9rem;
    transition: all 0.3s ease;
    backdrop-filter: blur(5px);
  }

  :deep(.form-input:focus),
  :deep(.form-select:focus),
  :deep(.form-textarea:focus) {
    outline: none;
    box-shadow:
      0 0 0 3px rgba(var(--color-brand-500-rgb), 0.1),
      0 0 20px rgba(var(--color-brand-500-rgb), 0.2);
    background: rgba(var(--color-background-900-rgb), 0.9);
  }

  :deep(.form-input::placeholder),
  :deep(.form-textarea::placeholder) {
    color: var(--color-content-subtle);
  }

  :deep(.form-textarea) {
    min-height: 120px;
    resize: vertical;
  }

  :deep(.form-switch) {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 12px 0;
  }

  :deep(.switch-container) {
    position: relative;
    width: 50px;
    height: 24px;
    background: rgba(var(--color-background-900-rgb), 0.8);
    border-radius: var(--radius-sm);
    cursor: pointer;
    transition: all 0.3s ease;
  }

  :deep(.switch-container.active) {
    background: linear-gradient(45deg, var(--color-brand-500), var(--color-error-500));
    box-shadow: 0 0 20px rgba(var(--color-brand-500-rgb), 0.3);
  }

  :deep(.switch-slider) {
    position: absolute;
    top: 1px;
    left: 1px;
    width: 20px;
    height: 20px;
    background: var(--color-content);
    border-radius: var(--radius-full);
    transition: all 0.3s ease;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
  }

  :deep(.switch-container.active .switch-slider) {
    transform: translateX(26px);
    box-shadow: 0 0 10px rgba(var(--color-brand-500-rgb), 0.5);
  }

  :deep(.switch-label) {
    color: rgba(var(--color-white-rgb), 0.8);
    font-size: 0.9rem;
    cursor: pointer;
  }

  :deep(.form-help) {
    color: var(--color-content-muted);
    font-size: 0.8rem;
    margin-top: 6px;
    line-height: 1.4;
  }

  :deep(.form-error) {
    color: var(--color-error-400);
    font-size: 0.8rem;
    margin-top: 6px;
    display: flex;
    align-items: center;
    gap: 4px;
  }

  :deep(.form-error i) {
    font-size: 0.75rem;
  }

  :deep(.form-success) {
    color: var(--color-success);
    font-size: 0.8rem;
    margin-top: 6px;
    display: flex;
    align-items: center;
    gap: 4px;
  }

  :deep(.form-button-group) {
    display: flex;
    gap: 12px;
    margin-top: 24px;
    padding-top: 20px;
  }

  :deep(.form-button) {
    padding: 10px 20px;
    border-radius: var(--radius-sm);
    font-weight: 500;
    cursor: pointer;
    transition: all 0.3s ease;
    display: flex;
    align-items: center;
    gap: 8px;
    font-size: 0.9rem;
    min-width: 120px;
    justify-content: center;
  }

  :deep(.form-button-primary) {
    background: linear-gradient(135deg, var(--color-brand-500), var(--color-error-500));
    color: var(--color-content);
  }

  :deep(.form-button-primary:hover) {
    background: linear-gradient(135deg, var(--color-error-500), var(--color-error-400));
    box-shadow: 0 0 20px rgba(var(--color-brand-500-rgb), 0.3);
    transform: translateY(-2px);
  }

  :deep(.form-button-secondary) {
    background: rgba(var(--color-background-900-rgb), 0.8);
    color: rgba(var(--color-white-rgb), 0.8);
  }

  :deep(.form-button-secondary:hover) {
    background: rgba(var(--color-background-900-rgb), 0.9);
    color: var(--color-content);
    box-shadow: 0 0 15px rgba(var(--color-brand-500-rgb), 0.2);
  }

  :deep(.form-button:disabled) {
    opacity: 0.5;
    cursor: not-allowed;
    transform: none !important;
    box-shadow: none !important;
  }

  :deep(.form-row) {
    display: flex;
    gap: 16px;
  }

  :deep(.form-col) {
    flex: 1;
  }

  :deep(.form-col-2) {
    flex: 2;
  }

  :deep(.form-col-3) {
    flex: 3;
  }

  :deep(.form-radio-group) {
    display: flex;
    flex-direction: column;
    gap: 12px;
  }

  :deep(.form-radio-item) {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 8px 12px;
    background: rgba(var(--color-background-900-rgb), 0.5);
    border-radius: var(--radius-sm);
    cursor: pointer;
    transition: all 0.3s ease;
  }

  :deep(.form-radio-item:hover) {
    background: rgba(var(--color-background-900-rgb), 0.7);
  }

  :deep(.form-radio-item.selected) {
    background: rgba(var(--color-brand-500-rgb), 0.1);
    box-shadow: 0 0 10px rgba(var(--color-brand-500-rgb), 0.1);
  }

  :deep(.form-radio) {
    width: 16px;
    height: 16px;
    border-radius: var(--radius-full);
    position: relative;
    flex-shrink: 0;
    background: rgba(var(--color-brand-500-rgb), 0.1);
  }

  :deep(.form-radio.selected) {
    background: rgba(var(--color-brand-500-rgb), 0.3);
    box-shadow: 0 0 10px rgba(var(--color-brand-500-rgb), 0.3);
  }

  :deep(.form-radio.selected::after) {
    content: '';
    position: absolute;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    width: 6px;
    height: 6px;
    background: var(--color-brand-500);
    border-radius: var(--radius-full);
    box-shadow: 0 0 8px rgba(var(--color-brand-500-rgb), 0.6);
  }

  :deep(.form-checkbox-group) {
    display: flex;
    flex-direction: column;
    gap: 12px;
  }

  :deep(.form-checkbox-item) {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 8px 12px;
    background: rgba(var(--color-background-900-rgb), 0.5);
    border-radius: var(--radius-sm);
    cursor: pointer;
    transition: all 0.3s ease;
  }

  :deep(.form-checkbox-item:hover) {
    background: rgba(var(--color-background-900-rgb), 0.7);
  }

  :deep(.form-checkbox-item.selected) {
    background: rgba(var(--color-brand-500-rgb), 0.1);
    box-shadow: 0 0 10px rgba(var(--color-brand-500-rgb), 0.1);
  }

  :deep(.form-checkbox) {
    width: 16px;
    height: 16px;
    border-radius: var(--radius-sm);
    position: relative;
    flex-shrink: 0;
    display: flex;
    align-items: center;
    justify-content: center;
    background: rgba(var(--color-brand-500-rgb), 0.1);
  }

  :deep(.form-checkbox.selected) {
    background: linear-gradient(135deg, var(--color-brand-500), var(--color-error-500));
    box-shadow: 0 0 10px rgba(var(--color-brand-500-rgb), 0.3);
  }

  :deep(.form-checkbox.selected::after) {
    content: '✓';
    color: var(--color-content);
    font-size: 10px;
    font-weight: bold;
    text-shadow: 0 0 4px rgba(var(--color-white-rgb), 0.8);
  }

  :deep(.form-file-input) {
    display: none;
  }

  :deep(.form-file-label) {
    display: inline-flex;
    align-items: center;
    gap: 8px;
    padding: 10px 16px;
    background: rgba(var(--color-background-900-rgb), 0.8);
    border-radius: var(--radius-sm);
    color: rgba(var(--color-white-rgb), 0.8);
    cursor: pointer;
    transition: all 0.3s ease;
    font-size: 0.9rem;
  }

  :deep(.form-file-label:hover) {
    background: rgba(var(--color-background-900-rgb), 0.9);
    color: var(--color-content);
    box-shadow: 0 0 15px rgba(var(--color-brand-500-rgb), 0.2);
  }

  :deep(.form-file-preview) {
    margin-top: 12px;
    padding: 12px;
    background: rgba(var(--color-background-800-rgb), 0.5);
    border-radius: var(--radius-sm);
  }

  :deep(.form-file-name) {
    color: var(--color-brand-500);
    font-size: 0.9rem;
    display: flex;
    align-items: center;
    gap: 8px;
  }

  :deep(.form-divider) {
    height: 1px;
    background: linear-gradient(90deg, transparent, rgba(var(--color-brand-500-rgb), 0.3), transparent);
    margin: 32px 0;
  }

  :deep(.form-number-input) {
    display: flex;
    align-items: center;
    gap: 8px;
  }

  :deep(.form-number-input input) {
    flex: 1;
    min-width: 0;
  }

  :deep(.form-number-buttons) {
    display: flex;
    flex-direction: column;
    gap: 2px;
  }

  :deep(.form-number-button) {
    width: 24px;
    height: 18px;
    background: rgba(var(--color-background-900-rgb), 0.8);
    color: var(--color-content-muted);
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 10px;
    transition: all 0.3s ease;
  }

  :deep(.form-number-button:first-child) {
    border-radius: var(--radius-sm) var(--radius-sm) 0 0;
  }

  :deep(.form-number-button:last-child) {
    border-radius: 0 0 var(--radius-sm) var(--radius-sm);
  }

  :deep(.form-number-button:hover) {
    background: rgba(var(--color-background-900-rgb), 0.9);
    color: var(--color-content);
  }

  :deep(.form-unit) {
    color: var(--color-content-muted);
    font-size: 0.9rem;
    margin-left: 8px;
  }

  :deep(.form-validation-message) {
    display: flex;
    align-items: center;
    gap: 6px;
    padding: 8px 12px;
    margin-top: 12px;
    border-radius: var(--radius-sm);
    font-size: 0.85rem;
    backdrop-filter: blur(5px);
  }

  :deep(.form-validation-message.error) {
    background: rgba(var(--color-error-rgb), 0.1);
    color: var(--color-error-400);
  }

  :deep(.form-validation-message.success) {
    background: rgba(16, 185, 129, 0.1);
    color: var(--color-success);
  }

  :deep(.form-validation-message.warning) {
    background: rgba(245, 158, 11, 0.1);
    color: var(--color-warning);
  }

  :deep(.form-validation-message.info) {
    background: rgba(var(--color-brand-500-rgb), 0.1);
    color: var(--color-brand-500);
  }

  @media (max-width: 768px) {
    :deep(.form-row) {
      flex-direction: column;
    }

    :deep(.form-button-group) {
      flex-direction: column;
    }

    :deep(.form-button) {
      width: 100%;
    }

    :deep(.form-container) {
      padding: 16px;
    }

    :deep(.form-section) {
      padding: 16px;
    }
  }
</style>
