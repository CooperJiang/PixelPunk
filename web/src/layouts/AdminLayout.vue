<script setup lang="ts">
  import { computed, onMounted } from 'vue'
  import { useRoute } from 'vue-router'
  import { useAuthStore } from '@/store/auth'
  import { useSettingsStore } from '@/store/settings'
  import { useResponsiveSidebar } from '@/composables/useResponsiveSidebar'
  import { useGlobalWebSocket } from '@/composables/useGlobalWebSocket'
  import { useTexts } from '@/composables/useTexts'
  import { getAdminMenuItems } from '@/constants'
  import AdminSidebar from '@/layouts/components/AdminSidebar.vue'

  const { $t } = useTexts()

  const ADMIN_MENU_ITEMS = computed(() => getAdminMenuItems($t))

  const route = useRoute()
  const authStore = useAuthStore()
  const settingsStore = useSettingsStore()
  const _userInfo = computed(() => authStore.userInfo)

  /* 使用响应式侧边栏组合式函数 */
  const {
    isCollapsed,
    isMobileOverlayOpen,
    isTransitioning,
    isMobile,
    isTablet,
    isDesktop,
    _currentSidebarWidth,
    toggleSidebar,
    closeMobileSidebar,
    handleRouteChange,
  } = useResponsiveSidebar()

  const { initWebSocketIfNeeded } = useGlobalWebSocket()

  /* 在组件挂载时加载设置并初始化WebSocket */
  onMounted(async () => {
    if (!settingsStore.isLoaded) {
      await settingsStore.loadGlobalSettings()
    }

    initWebSocketIfNeeded()
  })

  const pageTitle = computed(() => {
    const { path } = route
    if (path.includes('/admin/dashboard')) {
      return $t('admin.dashboard.title')
    }
    if (path.includes('/admin/users')) {
      return $t('admin.users.title')
    }
    if (path.includes('/admin/channels')) {
      return $t('admin.channels.title')
    }
    if (path.includes('/admin/categories')) {
      return $t('admin.categories.title')
    }
    if (path.includes('/admin/announcements')) {
      return $t('admin.announcements.title')
    }
    if (path.includes('/admin/settings')) {
      return $t('admin.settings.title')
    }
    return $t('admin.common.adminCenter')
  })

  const menuItems = computed(() => {
    let items = [...ADMIN_MENU_ITEMS.value]

    if (!settingsStore.isAIEnabled) {
      items = items.filter((item) => item.path !== '/admin/tagging')
    }

    if (!settingsStore.isVectorEnabled) {
      items = items.filter((item) => item.path !== '/admin/vectors')
    }

    return items
  })
</script>

<template>
  <div
    class="admin-layout"
    :class="{
      'is-transitioning': isTransitioning,
      mobile: isMobile,
      tablet: isTablet,
      desktop: isDesktop,
    }"
  >
    <div v-if="isMobile && isMobileOverlayOpen" class="mobile-overlay" @click="closeMobileSidebar" />

    <div class="admin-container">
      <AdminSidebar
        :menu-items="menuItems"
        :is-collapsed="isCollapsed"
        :is-mobile="isMobile"
        :is-mobile-open="isMobileOverlayOpen"
        :site-name="settingsStore.siteName"
        @toggle-sidebar="toggleSidebar"
        @close-mobile="closeMobileSidebar"
        @route-change="handleRouteChange"
      />

      <div
        class="admin-content"
        :class="{
          'sidebar-collapsed': isCollapsed,
          'mobile-mode': isMobile,
          'tablet-mode': isTablet,
        }"
      >
        <div v-if="isMobile" class="mobile-header">
          <button class="mobile-menu-btn" @click="toggleSidebar">
            <i class="fas fa-bars" />
          </button>
          <h1 class="mobile-title">{{ pageTitle }}</h1>
        </div>

        <div :key="route.path" class="content-body">
          <router-view />
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
  .admin-layout {
    --sidebar-width: 240px;
    --sidebar-collapsed-width: 60px;

    height: 100vh;
    background-color: var(--color-background-800);
    position: relative;
    overflow: hidden;
    font-family:
      'Inter',
      -apple-system,
      BlinkMacSystemFont,
      'Segoe UI',
      Roboto,
      sans-serif;
    color: var(--color-content-heading);
    transition: all var(--transition-base);
  }

  .admin-layout.mobile {
  }

  .admin-layout.tablet {
  }

  .admin-layout.desktop {
  }

  .mobile-overlay {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: var(--color-overlay-light);
    backdrop-filter: blur(4px);
    z-index: 1040;
    animation: fadeIn var(--transition-fast) forwards;
  }

  @keyframes fadeIn {
    from {
      opacity: 0;
    }
    to {
      opacity: 1;
    }
  }

  .admin-layout::-webkit-scrollbar {
    display: none;
  }

  .admin-container {
    display: flex;
    height: 100vh;
    overflow-x: hidden;
  }

  .admin-content {
    padding: 1rem 1rem 1rem 0.75rem;
    flex-grow: 1;
    transition:
      margin-left var(--transition-base),
      width var(--transition-base),
      padding-left var(--transition-base);
    box-sizing: border-box;
    will-change: margin-left, width;
    position: relative;
    overflow: hidden;
    z-index: 10;
  }

  .admin-content:not(.mobile-mode):not(.tablet-mode) {
    margin-left: var(--sidebar-width);
    width: calc(100% - var(--sidebar-width));
    height: 100vh;
  }

  .admin-content.sidebar-collapsed:not(.mobile-mode):not(.tablet-mode) {
    margin-left: var(--sidebar-collapsed-width);
    width: calc(100% - var(--sidebar-collapsed-width));
    height: 100vh;
  }

  .admin-content.mobile-mode {
    margin-left: 0 !important;
    width: 100% !important;
    height: 100vh;
    padding-top: 0;
    padding-left: 1rem;
  }

  .admin-content.tablet-mode {
    margin-left: 0 !important;
    width: 100% !important;
    height: 100vh;
    padding-left: 1rem;
  }

  .mobile-header {
    display: flex;
    align-items: center;
    gap: var(--spacing-md);
    padding: var(--spacing-md);
    background: var(--color-background-700);
    border-bottom: 1px solid var(--color-border-subtle);
    margin: calc(-1 * var(--spacing-md)) calc(-1 * var(--spacing-md)) var(--spacing-md);
    backdrop-filter: blur(8px);
    position: sticky;
    top: 0;
    z-index: 1030;
  }

  .mobile-menu-btn {
    background: none;
    border: none;
    color: var(--color-content-heading);
    cursor: pointer;
    padding: var(--spacing-sm);
    border-radius: var(--radius-md);
    transition: all var(--transition-fast);
    display: flex;
    align-items: center;
    justify-content: center;
    width: 40px;
    height: 40px;
  }

  .mobile-menu-btn:hover {
    background: var(--color-hover-bg);
    color: var(--color-brand-500);
  }

  .mobile-title {
    font-size: 1.125rem;
    font-weight: 600;
    color: var(--color-content-heading);
    margin: 0;
  }

  .content-body {
    height: calc(100vh - 2rem);
    overflow: hidden;
    position: relative;
    transition: all var(--transition-base);
    margin: 0;
    contain: layout;
    background: var(--color-background-600);
    border-radius: var(--radius-lg);
    border: 1px solid var(--color-border-default);
    box-shadow:
      0 2px 8px rgba(0, 0, 0, 0.06),
      0 4px 16px rgba(0, 0, 0, 0.04),
      inset 0 1px 0 rgba(255, 255, 255, 0.1);
    padding: 0;
  }

  .admin-content.mobile-mode .content-body {
    height: calc(100vh - 80px);
    margin: 0.25rem;
    padding: 0;
    border-radius: var(--radius-lg);
  }

  .admin-content.tablet-mode .content-body {
    height: calc(100vh - 2rem);
    margin: 0.5rem;
    padding: 0;
  }

  .content-body .admin-content-area {
    background: var(--color-background-700);
    border-radius: var(--radius-lg);
    box-shadow: var(--shadow-cyber-md);
    backdrop-filter: blur(6px);
  }

  .content-body .admin-card {
    background: var(--color-background-700);
    border: 1px solid var(--color-border-default);
    border-radius: var(--radius-lg);
    padding: 1.25rem;
    box-shadow: var(--shadow-cyber-sm);
    transition: all var(--transition-fast);
  }

  .content-body .admin-card:hover {
    border-color: var(--color-hover-border);
    box-shadow: var(--shadow-cyber-md);
    transform: translateY(-2px);
  }

  .content-body .admin-table {
    background: var(--color-background-700);
    border: 1px solid var(--color-border-subtle);
    border-radius: var(--radius-md);
    overflow: hidden;
    box-shadow: var(--shadow-cyber-sm);
  }

  .content-body .table-container {
    overflow-x: auto;
    overflow-y: hidden;
    max-width: 100%;
    margin: 0 auto;
  }

  .content-body .cyber-table-container {
    overflow-x: auto;
    overflow-y: hidden;
    max-width: 100%;
    margin: 0 auto;
  }

  .content-body .cyber-table {
    min-width: 100%;
    table-layout: fixed;
  }

  .content-body .admin-page-container {
    max-width: 100%;
    overflow: hidden;
    height: 100%;
    display: flex;
    flex-direction: column;
  }

  .content-body .admin-table th {
    background: var(--color-background-800);
    color: var(--color-content-default);
    border-bottom: 1px solid var(--color-border-subtle);
  }

  .content-body .admin-table td {
    border-bottom: 1px solid var(--color-divider);
  }

  .content-body .admin-table tr:hover {
    background: var(--color-hover-bg);
  }

  .content-body .admin-tabs {
    border-bottom: 1px solid var(--color-border-subtle);
    background: var(--color-background-800);
    border-radius: var(--radius-md) var(--radius-md) 0 0;
  }

  .content-body .admin-tab {
    color: var(--color-content-default);
    border-bottom: 2px solid transparent;
    transition: all var(--transition-fast);
  }

  .content-body .admin-tab:hover {
    color: var(--color-brand-500);
    border-bottom-color: var(--color-brand-500);
  }

  .content-body .admin-tab.active {
    color: var(--color-brand-500);
    border-bottom-color: var(--color-brand-500);
    background: var(--color-hover-bg);
  }

  .content-body .admin-empty {
    background: var(--color-background-800);
    border: 1px dashed var(--color-border-subtle);
    border-radius: var(--radius-lg);
    padding: 3rem;
    text-align: center;
    color: var(--color-content-muted);
  }

  .content-body .admin-empty-icon {
    font-size: 3rem;
    color: var(--color-brand-500);
    margin-bottom: 1rem;
    opacity: 0.6;
  }

  .content-body .cyber-btn.primary {
    background: linear-gradient(135deg, var(--color-brand-500), var(--color-brand-600));
    border-color: var(--color-brand-500);
    box-shadow: var(--shadow-glow-primary-sm);
  }

  .content-body .cyber-btn.primary:hover {
    background: linear-gradient(135deg, var(--color-brand-600), var(--color-brand-500));
    box-shadow: var(--shadow-glow-primary-md);
    transform: translateY(-1px);
  }

  .content-body .cyber-btn.secondary {
    background: var(--color-background-700);
    border-color: var(--color-border-default);
    color: var(--color-content-default);
  }

  .content-body .cyber-btn.secondary:hover {
    background: var(--color-hover-bg);
    border-color: var(--color-hover-border);
    color: var(--color-brand-500);
  }

  .content-body::-webkit-scrollbar {
    width: 6px;
  }

  .content-body::-webkit-scrollbar-track {
    background: var(--color-background-800);
    border-radius: var(--radius-sm);
  }

  .content-body::-webkit-scrollbar-thumb {
    background-color: var(--color-border-default);
    border-radius: var(--radius-sm);
  }

  .content-body::-webkit-scrollbar-thumb:hover {
    background-color: var(--color-brand-500);
  }

  .admin-layout {
    scrollbar-width: none;
    -ms-overflow-style: none;
  }

  .admin-layout::-webkit-scrollbar {
    display: none;
  }

  .admin-route-container {
    animation: slideUp var(--transition-base) var(--ease-out);
  }
</style>
