<script setup lang="ts">
  import { computed } from 'vue'
  import { useRoute } from 'vue-router'
  import { useLayoutStore } from '@/store/layout'
  import { useBreadcrumbStore } from '@/store/breadcrumb'
  import { useFolderBreadcrumb } from '@/hooks/useFolderBreadcrumb'
  import { useTexts } from '@/composables/useTexts'
  import UserTerminal from '@/components/Navbar/UserTerminal.vue'
  import Breadcrumb from '@/components/Breadcrumb/index.vue'
  import FolderBreadcrumb from '@/pages/folders/components/FolderBreadcrumb.vue'

  defineOptions({
    name: 'SidebarTopbar',
  })

  const emit = defineEmits<{
    (e: 'login'): void
  }>()

  const { $t } = useTexts()
  const route = useRoute()
  const layoutStore = useLayoutStore()
  const breadcrumbStore = useBreadcrumbStore()
  const { handleBreadcrumbClick } = useFolderBreadcrumb()

  const topbarStyle = computed(() => ({
    left: `${layoutStore.currentSidebarWidth}px`,
    right: 0,
  }))

  /* 是否为 folders 页面 */
  const isFoldersPage = computed(() => {
    return route.name === 'folders' || route.name === 'foldersWithPath'
  })

  /* folders 页面的面包屑数据（添加首页和我的文件夹前缀） */
  const foldersBreadcrumbItems = computed(() => {
    if (!isFoldersPage.value) return []

    const items = []

    items.push({ id: 'home', name: $t('components.sidebar.topbar.home'), isSpecial: true })

    items.push({ id: null, name: $t('components.sidebar.topbar.myFolders') })

    if (breadcrumbStore.items.length > 0) {
      items.push(...breadcrumbStore.items)
    }

    return items
  })

  const handleLogin = () => {
    emit('login')
  }
</script>

<template>
  <div class="sidebar-topbar" :style="topbarStyle">
    <div class="topbar-left">
      <FolderBreadcrumb
        v-if="isFoldersPage"
        :items="foldersBreadcrumbItems"
        :show-icon="false"
        :auto-add-root="false"
        @click="handleBreadcrumbClick"
      />

      <Breadcrumb v-else />
    </div>

    <div class="topbar-right">
      <UserTerminal :mobile-menu-open="false" @login="handleLogin" @toggle-mobile-menu="() => {}" />
    </div>
  </div>
</template>

<style scoped>
  .sidebar-topbar {
    @apply fixed top-0 flex h-16 items-center justify-between bg-background-800 px-6;
    z-index: 100;
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);

    background: var(--color-background-800);

    border-bottom: 1px solid var(--color-border-subtle);
    box-sizing: border-box;
  }

  .topbar-left {
    @apply flex items-center gap-3;
  }

  .sidebar-toggle-btn {
    @apply flex h-8 w-8 items-center justify-center;
    transition: all 0.2s ease;
    font-size: 12px;
    border-radius: var(--radius-sm);
  }

  .sidebar-toggle-btn:hover {
    transform: scale(1.05);
  }

  .topbar-right {
    @apply flex items-center gap-3;
  }

  .message-notification-wrapper {
    @apply flex items-center;
  }

  .layout-toggle-btn {
    @apply flex h-8 w-8 items-center justify-center;
    transition: all 0.2s ease;
    font-size: 12px;
    border-radius: var(--radius-sm);
  }

  .layout-toggle-btn:hover {
    transform: scale(1.05);
  }

  .user-section {
    @apply relative;
  }

  .user-dropdown {
    @apply relative;
  }

  .user-trigger {
    @apply flex items-center gap-3 px-3 py-2;
    transition: all 0.2s ease;
    border-radius: var(--radius-sm);
  }

  .user-avatar {
    @apply relative h-8 w-8 overflow-hidden rounded-full;
  }

  .avatar-image {
    @apply h-full w-full object-cover;
  }

  .avatar-fallback {
    @apply flex h-full w-full items-center justify-center;
    font-size: 12px;
  }

  .online-indicator {
    @apply absolute -bottom-0.5 -right-0.5 h-3 w-3 rounded-full;
    box-shadow: 0 0 6px rgba(var(--color-success-rgb), 0.6);
  }

  .user-info {
    @apply flex flex-col;
  }

  .username {
    @apply text-sm font-medium;
    line-height: 1.2;
  }

  .user-role {
    @apply text-xs;
    line-height: 1.2;
  }

  .dropdown-arrow {
    @apply text-xs text-content-muted transition-transform duration-200;
  }

  .dropdown-open .dropdown-arrow {
    transform: rotate(180deg);
  }

  .user-menu {
    @apply absolute right-0 top-full mt-2 min-w-48 py-2;
    transform: translateY(-10px) scale(0.95);
    opacity: 0;
    visibility: hidden;
    transition: all 0.2s ease;
    z-index: 100;
    border-radius: var(--radius-sm);
  }

  .dropdown-open .user-menu {
    transform: translateY(0) scale(1);
    opacity: 1;
    visibility: visible;
  }

  .menu-item {
    @apply flex cursor-pointer items-center gap-3 px-4 py-3;
    transition: background-color 0.2s ease;
  }

  .menu-item i {
    @apply w-4 text-sm text-content-muted;
  }

  .login-btn {
    @apply flex items-center gap-2 px-4 py-2;
    transition: all 0.2s ease;
    font-weight: 500;
    border-radius: var(--radius-sm);
  }
</style>
