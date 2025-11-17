<script setup lang="ts">
  import { computed, onMounted, ref } from 'vue'
  import { useToast } from '@/components/Toast/useToast'
  import { useRoute, useRouter } from 'vue-router'
  import { useTexts } from '@/composables/useTexts'
  import WebsiteSettings from './components/WebsiteSettings.vue'
  import AppearanceSettings from './components/AppearanceSettings.vue'
  import PrivacySettings from './components/PrivacySettings.vue'
  import TermsSettings from './components/TermsSettings.vue'
  /* import ThemeSettings from './components/ThemeSettings.vue' // 已停用多模式功能 */
  import { batchUpsertSettings, defaultSettings, getSettings, type Setting } from '@/api/admin/settings'

  defineOptions({
    name: 'AdminConstructionPage',
  })

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
    { key: 'website', name: $t('admin.construction.tabs.website'), icon: 'fas fa-globe' },
    { key: 'appearance', name: $t('admin.construction.tabs.appearance'), icon: 'fas fa-paint-brush' },
    { key: 'privacy', name: $t('admin.construction.tabs.privacy'), icon: 'fas fa-shield-alt' },
    { key: 'terms', name: $t('admin.construction.tabs.terms'), icon: 'fas fa-file-contract' },
  ])

  /* 从查询参数获取初始tab，如果没有或者不存在则默认为'website' */
  const validTabs = ['website', 'appearance', 'privacy', 'terms']
  const requestedTab = route.query.tab?.toString() || 'website'
  const activeTab = ref(validTabs.includes(requestedTab) ? requestedTab : 'website')

  const updateActiveTab = (tabId: string) => {
    if (validTabs.includes(tabId)) {
      activeTab.value = tabId
      router.push({
        name: route.name,
        query: { ...route.query, tab: tabId },
      })
    }
  }

  const websiteSettings = ref<Setting[]>([])
  const appearanceSettings = ref<Setting[]>([])
  const privacySettings = ref<Setting[]>([])
  const termsSettings = ref<Setting[]>([])

  const modifiedWebsiteSettings = ref<Setting[]>([])
  const modifiedAppearanceSettings = ref<Setting[]>([])
  const modifiedPrivacySettings = ref<Setting[]>([])
  const modifiedTermsSettings = ref<Setting[]>([])

  const handleWebsiteSettingsUpdate = (updatedSettings: Setting[]) => {
    modifiedWebsiteSettings.value = updatedSettings
  }

  const handleAppearanceSettingsUpdate = (updatedSettings: Setting[]) => {
    modifiedAppearanceSettings.value = updatedSettings
  }

  const handlePrivacySettingsUpdate = (updatedSettings: Setting[]) => {
    modifiedPrivacySettings.value = updatedSettings
  }

  const handleTermsSettingsUpdate = (updatedSettings: Setting[]) => {
    modifiedTermsSettings.value = updatedSettings
  }

  const resetSettings = async () => {
    isLoading.value = true
    await loadConstructionSettings()

    if (activeTab.value === 'website') {
      modifiedWebsiteSettings.value = []
    } else if (activeTab.value === 'appearance') {
      modifiedAppearanceSettings.value = []
    } else if (activeTab.value === 'privacy') {
      modifiedPrivacySettings.value = []
    } else if (activeTab.value === 'terms') {
      modifiedTermsSettings.value = []
    }

    toast.success($t('admin.construction.toast.reset'))
    isLoading.value = false
  }

  const saveSettings = async () => {
    let settingsToUpdate: Setting[] = []

    if (activeTab.value === 'website') {
      settingsToUpdate = modifiedWebsiteSettings.value
    } else if (activeTab.value === 'appearance') {
      settingsToUpdate = modifiedAppearanceSettings.value
    } else if (activeTab.value === 'privacy') {
      settingsToUpdate = modifiedPrivacySettings.value
    } else if (activeTab.value === 'terms') {
      settingsToUpdate = modifiedTermsSettings.value
    }

    if (settingsToUpdate.length === 0) {
      toast.info($t('admin.construction.toast.noChanges'))
      return
    }

    try {
      isSaving.value = true
      const result = await batchUpsertSettings(settingsToUpdate)

      if (result.success) {
        if (activeTab.value === 'website') {
          modifiedWebsiteSettings.value = []
        } else if (activeTab.value === 'appearance') {
          modifiedAppearanceSettings.value = []
        } else if (activeTab.value === 'privacy') {
          modifiedPrivacySettings.value = []
        } else if (activeTab.value === 'terms') {
          modifiedTermsSettings.value = []
        }

        toast.success($t('admin.construction.toast.saved'))
      }
    } catch (error) {
      toast.error(error.message)
    } finally {
      isSaving.value = false
    }
  }

  const loadConstructionSettings = async () => {
    isLoading.value = true

    try {
      // const constructionResult = await getSettings({ group: 'construction' })
      const websiteInfoResult = await getSettings({ group: 'website_info' })
      const appearanceResult = await getSettings({ group: 'appearance' })
      const legalResult = await getSettings({ group: 'legal' })

      // const constructionSettingsData = constructionResult.success ? constructionResult.data?.settings || [] : []

      const websiteInfoSettingsData = websiteInfoResult.success ? websiteInfoResult.data?.settings || [] : []

      const appearanceSettingsData = appearanceResult.success ? appearanceResult.data?.settings || [] : []

      const legalSettingsData = legalResult.success ? legalResult.data?.settings || [] : []

      if (websiteInfoSettingsData.length === 0) {
        websiteSettings.value = defaultSettings.website_info || []
      } else {
        websiteSettings.value = websiteInfoSettingsData
      }

      if (appearanceSettingsData.length === 0) {
        appearanceSettings.value = defaultSettings.appearance || []
      } else {
        appearanceSettings.value = appearanceSettingsData
      }

      if (legalSettingsData.length === 0) {
        privacySettings.value = [
          {
            key: 'privacy_policy_content',
            value: '',
            type: 'text',
            group: 'legal',
            description: $t('descriptions.privacyPolicyContent'),
          },
        ]
      } else {
        privacySettings.value = legalSettingsData.filter((s) => s.key === 'privacy_policy_content')
        if (privacySettings.value.length === 0) {
          privacySettings.value = [
            {
              key: 'privacy_policy_content',
              value: '',
              type: 'text',
              group: 'legal',
              description: $t('descriptions.privacyPolicyContent'),
            },
          ]
        }
      }

      if (legalSettingsData.length === 0) {
        termsSettings.value = [
          {
            key: 'terms_of_service_content',
            value: '',
            type: 'text',
            group: 'legal',
            description: $t('descriptions.termsOfServiceContent'),
          },
        ]
      } else {
        termsSettings.value = legalSettingsData.filter((s) => s.key === 'terms_of_service_content')
        if (termsSettings.value.length === 0) {
          termsSettings.value = [
            {
              key: 'terms_of_service_content',
              value: '',
              type: 'text',
              group: 'legal',
              description: $t('descriptions.termsOfServiceContent'),
            },
          ]
        }
      }
    } catch (error) {
      websiteSettings.value = defaultSettings.website_info || []
      appearanceSettings.value = defaultSettings.appearance || []
      privacySettings.value = [
        {
          key: 'privacy_policy_content',
          value: '',
          type: 'text',
          group: 'legal',
          description: $t('descriptions.privacyPolicyContent'),
        },
      ]
      termsSettings.value = [
        {
          key: 'terms_of_service_content',
          value: '',
          type: 'text',
          group: 'legal',
          description: $t('descriptions.termsOfServiceContent'),
        },
      ]
      toast.error(error.message)
    } finally {
      isLoading.value = false
    }
  }

  onMounted(async () => {
    await loadConstructionSettings()
  })
</script>

<template>
  <div class="admin-construction-page admin-page-container">
    <CyberAdminWrapper
      :title="$t('admin.construction.title')"
      :subtitle="$t('admin.construction.subtitle')"
      icon="fas fa-tools"
      :sidebar-layout="true"
    >
      <template #actions>
        <CyberButton type="text" class="min-w-24" @click="resetSettings">{{
          $t('admin.construction.buttons.reset')
        }}</CyberButton>
        <CyberButton type="primary" :loading="isSaving" loading-mode="inline" icon="save" class="min-w-24" @click="saveSettings">
          {{ isSaving ? $t('admin.construction.buttons.saving') : $t('admin.construction.buttons.save') }}
        </CyberButton>
      </template>

      <template #sidebar>
        <CyberSidebarNav :tabs="tabs" :active-tab="activeTab" @tab-change="updateActiveTab" />
      </template>

      <template #content>
        <div v-show="activeTab === 'website'" class="content-panel">
          <WebsiteSettings :settings="websiteSettings" @update="handleWebsiteSettingsUpdate" />
        </div>

        <div v-show="activeTab === 'appearance'" class="content-panel">
          <AppearanceSettings :settings="appearanceSettings" @update="handleAppearanceSettingsUpdate" />
        </div>

        <div v-show="activeTab === 'privacy'" class="content-panel">
          <PrivacySettings :settings="privacySettings" @update="handlePrivacySettingsUpdate" />
        </div>

        <div v-show="activeTab === 'terms'" class="content-panel">
          <TermsSettings :settings="termsSettings" @update="handleTermsSettingsUpdate" />
        </div>

        <CyberLoading v-if="isLoading" />
      </template>
    </CyberAdminWrapper>
  </div>
</template>

<style scoped>
  .admin-construction-page {
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
  }

  .homepage-settings-placeholder {
    background: linear-gradient(
      135deg,
      rgba(var(--color-background-800-rgb), 0.8) 0%,
      rgba(var(--color-background-900-rgb), 0.6) 100%
    );
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.2);
    border-radius: var(--radius-sm);
    backdrop-filter: var(--backdrop-blur-md);
    box-shadow:
      var(--shadow-lg),
      0 0 40px rgba(var(--color-brand-500-rgb), 0.05);
  }
</style>
