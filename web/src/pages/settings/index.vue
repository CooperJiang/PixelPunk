<script setup lang="ts">
  import { computed, onMounted, ref, watch } from 'vue'
  import { useRoute } from 'vue-router'
  import ProfileSection from './components/ProfileSection.vue'
  import SecuritySection from './components/SecuritySection.vue'
  import ApiKeySection from './components/ApiKeySection.vue'
  import AccessControlSection from './components/AccessControlSection.vue'
  import PreferencesSection from './components/PreferencesSection.vue'
  import { useTexts } from '@/composables/useTexts'

  const route = useRoute()
  const { $t } = useTexts()

  const navItems = computed(() => [
    { id: 'profile', label: $t('settings.nav.profile'), icon: 'user', iconClass: 'text-content' },
    { id: 'security', label: $t('settings.nav.security'), icon: 'lock', iconClass: 'text-warning-400' },
    { id: 'api', label: $t('settings.nav.apiKeys'), icon: 'key', iconClass: 'text-error-400' },
    {
      id: 'access-control',
      label: $t('settings.nav.accessControl'),
      icon: 'shield-alt',
      iconClass: 'text-success-400',
    },
    {
      id: 'preferences',
      label: $t('settings.nav.preferences'),
      icon: 'cog',
      iconClass: 'text-purple-400',
    },
  ])

  const profileSectionRef = ref(null)
  const apiKeySectionRef = ref(null)
  const accessControlSectionRef = ref(null)
  const preferencesSectionRef = ref(null)

  const activeSection = ref('profile')

  const updateUrlHash = (hash: string) => {
    window.history.replaceState(null, '', `#${hash}`)
  }

  const loadComponentData = (tabId: string) => {
    switch (tabId) {
      case 'profile':
        if (profileSectionRef.value) {
          profileSectionRef.value.fetchUserProfile?.()
        }
        break
      case 'api':
        if (apiKeySectionRef.value) {
          const { apiKeyManagerRef } = apiKeySectionRef.value
          if (apiKeyManagerRef) {
            apiKeyManagerRef.fetchApiKeys()
          }
        }
        break
      case 'access-control':
        if (accessControlSectionRef.value) {
          accessControlSectionRef.value.fetchAccessControlConfig?.()
        }
        break
      case 'preferences':
        break
    }
  }

  const switchTab = (tabId: string) => {
    activeSection.value = tabId
    updateUrlHash(tabId)
    loadComponentData(tabId)
  }

  watch(
    () => route.hash,
    (newHash) => {
      const hash = newHash.replace('#', '')
      if (navItems.value.some((item) => item.id === hash)) {
        activeSection.value = hash
        loadComponentData(hash)
      }
    },
    { immediate: true }
  )

  onMounted(() => {
    const initialHash = window.location.hash.substring(1)
    if (initialHash && navItems.value.some((item) => item.id === initialHash)) {
      activeSection.value = initialHash
    } else {
      activeSection.value = 'profile'
      updateUrlHash('profile')
    }
  })
</script>

<template>
  <div class="settings-page">
    <header class="settings-header">
      <div class="header-main">
        <div class="header-left">
          <div class="title-section">
            <h1 class="page-title">
              <i class="fas fa-cog title-icon" />
              <span class="title-text">{{ $t('settings.page.title') }}</span>
            </h1>
            <div class="subtitle">{{ $t('settings.page.subtitle') }}</div>
          </div>
        </div>
      </div>

      <div class="header-topbar">
        <div class="nav-tabs">
          <button
            v-for="item in navItems"
            :key="item.id"
            class="nav-tab"
            :class="{ active: activeSection === item.id }"
            @click="switchTab(item.id)"
          >
            <i :class="['fas', `fa-${item.icon}`, item.iconClass]" />
            <span>{{ item.label }}</span>
          </button>
        </div>
      </div>
    </header>

    <main class="settings-content">
      <div class="content-container">
        <ProfileSection v-show="activeSection === 'profile'" ref="profileSectionRef" />
        <SecuritySection v-show="activeSection === 'security'" />
        <ApiKeySection v-show="activeSection === 'api'" ref="apiKeySectionRef" />
        <AccessControlSection v-show="activeSection === 'access-control'" ref="accessControlSectionRef" />
        <PreferencesSection v-show="activeSection === 'preferences'" ref="preferencesSectionRef" />
      </div>
    </main>
  </div>
</template>

<style scoped>
  .settings-page {
    @apply relative flex h-full flex-col overflow-hidden;
    animation: layoutFadeIn 0.4s ease-out;
  }

  @keyframes layoutFadeIn {
    from {
      opacity: 0;
      transform: translateY(8px);
    }
    to {
      opacity: 1;
      transform: translateY(0);
    }
  }

  .settings-header {
    @apply relative z-[100] flex-shrink-0;
    border-bottom: 1px solid var(--color-border-subtle);
  }

  .header-main {
    @apply flex min-h-12 items-center justify-between gap-6 py-2.5;
  }

  .header-left {
    @apply flex min-w-0 flex-1 items-center gap-5;
  }

  .title-section {
    @apply min-w-0 shrink-0;
  }

  .page-title {
    @apply m-0 flex items-center gap-3 text-base font-bold leading-tight;
  }

  .title-icon {
    @apply shrink-0 text-sm text-brand-400;
  }

  .title-text {
    @apply overflow-hidden text-ellipsis whitespace-nowrap text-content-heading;
    letter-spacing: 0.025em;
  }

  .subtitle {
    @apply mt-0.5 text-xs font-normal text-content-muted opacity-85;
  }

  .header-topbar {
    @apply py-2;
    border-top: 1px solid var(--color-border-subtle);
    background: transparent;
  }

  .nav-tabs {
    @apply flex gap-2 overflow-x-auto pb-2;
  }

  .nav-tab {
    @apply flex items-center gap-2 border px-4 py-2 text-sm font-medium transition-all duration-200;
    border-radius: var(--radius-sm);
    border-color: var(--color-border-subtle);
    color: var(--color-content-muted);
    background: transparent;
    cursor: pointer;
    white-space: nowrap;
  }

  .nav-tab:hover {
    border-color: var(--color-border-default);
    color: var(--color-content-default);
  }

  .nav-tab.active {
    border-color: var(--color-brand-400);
    color: var(--color-brand-400);
    background-color: rgba(var(--color-brand-500-rgb), 0.12);
    box-shadow: 0 6px 20px rgba(var(--color-brand-500-rgb), 0.12);
  }

  .nav-tab i {
    @apply text-xs;
  }

  .settings-content {
    @apply relative flex-1 overflow-y-auto overflow-x-hidden;
    scrollbar-gutter: stable;
    scrollbar-width: thin;
  }

  .content-container {
    @apply space-y-4;
  }

  @media (max-width: 768px) {
    .header-main {
      @apply flex-col items-stretch gap-4 py-3;
    }

    .header-left {
      @apply flex-col items-start gap-3;
    }

    .header-topbar {
      @apply py-3;
    }

    .content-container {
      @apply py-5;
    }
  }

  @media (max-width: 480px) {
    .content-container {
      @apply py-4;
    }
  }
</style>
