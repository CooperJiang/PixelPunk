<script setup lang="ts">
  import { computed } from 'vue'
  import { useRoute, useRouter } from 'vue-router'
  import { useAuthStore } from '@/store/auth'
  import { useTexts } from '@/composables/useTexts'

  defineOptions({
    name: 'SidebarNav',
  })

  interface Props {
    collapsed?: boolean
  }

  withDefaults(defineProps<Props>(), {
    collapsed: false,
  })

  const { $t } = useTexts()
  const route = useRoute()
  const router = useRouter()
  const authStore = useAuthStore()

  /* 分类导航菜单配置 */
  const navSections = computed(() => [
    {
      title: null, // 主要功能区，不显示分类标题
      items: [
        { name: $t('components.sidebar.nav.dashboard'), path: '/dashboard', icon: 'fas fa-tachometer-alt' },
        { name: $t('components.sidebar.nav.folders'), path: '/folders', icon: 'fas fa-folder' },
        { name: $t('components.sidebar.nav.resource'), path: '/resource', icon: 'fas fa-images' },
        { name: $t('components.sidebar.nav.explore'), path: '/explore', icon: 'fas fa-compass' },
        { name: $t('components.sidebar.nav.upload'), path: '/upload', icon: 'fas fa-cloud-upload-alt' },
        { name: $t('components.sidebar.nav.messages'), path: '/messages', icon: 'fas fa-envelope' },
        { name: $t('components.sidebar.nav.shares'), path: '/shares', icon: 'fas fa-share-alt' },
      ],
    },
    {
      title: $t('components.sidebar.nav.management'),
      icon: 'fas fa-cogs',
      items: [
        { name: $t('components.sidebar.nav.tagManage'), path: '/tag-manage', icon: 'fas fa-tags' },
        { name: $t('components.sidebar.nav.categoryManage'), path: '/category-manage', icon: 'fas fa-sitemap' },
        { name: $t('components.sidebar.nav.automation'), path: '/automation', icon: 'fas fa-robot' },
        { name: $t('components.sidebar.nav.openApi'), path: '/open-api', icon: 'fas fa-random' },
      ],
    },
    {
      title: $t('components.sidebar.nav.system'),
      icon: 'fas fa-cog',
      items: [{ name: $t('components.sidebar.nav.settings'), path: '/settings', icon: 'fas fa-user-cog' }],
    },
  ])

  /* 过滤权限相关的导航项 */
  const filteredNavSections = computed(() => {
    return navSections.value
      .map((section) => {
        if (section.requiresAdmin && !authStore.isAdmin) {
          return null
        }

        const filteredItems = section.items.filter((item) => {
          return !item.requiresAdmin || authStore.isAdmin
        })

        return {
          ...section,
          items: filteredItems,
        }
      })
      .filter((section) => section !== null && section.items.length > 0)
  })

  const isActive = (path: string) => {
    if (path === '/dashboard') {
      return route.path === '/dashboard' || route.path === '/'
    }
    return route.path.startsWith(path)
  }

  const navigateTo = (path: string) => {
    router.push(path)
  }
</script>

<template>
  <div class="sidebar-nav-container" :class="{ 'no-padding': collapsed }">
    <div
      v-for="(section, index) in filteredNavSections"
      :key="section.title || `section-${index}`"
      class="nav-section"
      :class="{ 'nav-section-collapsed': collapsed }"
    >
      <div v-if="!collapsed && section.title" class="nav-section-header">
        <div class="nav-section-title">
          <i :class="section.icon" class="section-icon" />
          <span class="section-text">{{ section.title }}</span>
        </div>
      </div>

      <div class="nav-items" :class="{ 'nav-items-collapsed': collapsed }">
        <CyberTooltip
          v-for="item in section.items"
          :key="item.path"
          :content="item.name"
          placement="right"
          :disabled="!collapsed"
        >
          <div
            class="nav-item"
            :class="[{ 'nav-item-active': isActive(item.path) }, { 'nav-item-collapsed': collapsed }]"
            @click="navigateTo(item.path)"
          >
            <div class="nav-item-icon">
              <i :class="item.icon" />
            </div>
            <span class="nav-item-text" :class="{ 'nav-text-collapsed': collapsed }">
              {{ item.name }}
            </span>
          </div>
        </CyberTooltip>
      </div>
    </div>
  </div>
</template>

<style scoped>
  .sidebar-nav-container {
    @apply px-3 py-2;
  }

  .sidebar-nav-container.no-padding {
    @apply px-1 py-2;
  }

  .nav-section {
    @apply mb-6 space-y-1;
  }

  .nav-section:last-child {
    @apply mb-0;
  }

  .nav-section-collapsed {
    @apply mb-4;
  }

  .nav-section-header {
    @apply mb-2;
  }

  .nav-section-title {
    @apply flex items-center px-2 py-1;
    gap: 6px;
    font-size: 11px;
    font-weight: 600;
    text-transform: uppercase;
    letter-spacing: 0.5px;
    @apply text-content-muted;
  }

  .section-icon {
    font-size: 10px;
    @apply text-content-muted;
  }

  .section-text {
    user-select: none;
  }

  .nav-items {
    @apply space-y-1;
  }

  .nav-items-collapsed {
    @apply flex flex-col items-center;
    gap: 6px;
  }

  .nav-items-collapsed > * {
    margin-top: 0 !important;
  }

  .nav-item {
    @apply relative flex cursor-pointer items-center text-content-muted;
    padding: 10px 12px;
    border: 1px solid transparent;
    border-radius: var(--radius-sm);
    transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
    position: relative;
    overflow: hidden;
  }

  .nav-item::after {
    content: '';
    position: absolute;
    left: 0;
    top: 50%;
    transform: translateY(-50%);
    width: 3px;
    height: 0;
    background: var(--color-brand-500);
    transition: height 0.3s ease;
    border-radius: 0 var(--radius-sm) var(--radius-sm) 0;
  }

  .nav-item:hover {
    @apply text-content;
    background: rgba(var(--color-brand-500-rgb), 0.08);
    border-color: rgba(var(--color-brand-500-rgb), 0.2);
    transform: translateX(2px);
    box-shadow: 0 2px 8px rgba(var(--color-brand-500-rgb), 0.1);
  }

  .nav-item:hover::after {
    height: 60%;
  }

  .nav-item-collapsed {
    @apply items-center justify-center;

    width: 36px;
    height: 36px;
    padding: 0;
    margin: 0 auto 6px;
    border-radius: var(--radius-sm);

    display: flex;
    align-items: center;
    justify-content: center;

    background: rgba(var(--color-background-700-rgb), 0.5);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.1);
    position: relative;
    overflow: hidden;
  }

  .nav-item-collapsed::after {
    display: none;
  }

  .nav-item-icon {
    @apply flex h-5 w-5 flex-shrink-0 items-center justify-center;
    font-size: 14px;
  }

  .nav-item-collapsed .nav-item-icon {
    width: 100%;
    height: 100%;
    border-radius: var(--radius-sm);
    font-size: 14px;
    display: flex;
    align-items: center;
    justify-content: center;

    @apply text-content-muted;
  }

  .nav-item-collapsed .nav-item-icon i {
    line-height: 1;
    display: block;
  }

  .nav-item-text {
    @apply ml-3 text-sm font-medium;
    transition: all 0.3s ease;
    white-space: nowrap;
    overflow: hidden;
    width: auto;
    opacity: 1;
  }

  .nav-text-collapsed {
    @apply ml-0 w-0;
    opacity: 0;
    pointer-events: none;
  }

  .nav-item-active {
    background: linear-gradient(135deg, rgba(var(--color-brand-500-rgb), 0.2) 0%, rgba(var(--color-brand-500-rgb), 0.1) 100%);
    border-color: rgba(var(--color-brand-500-rgb), 0.4);
    @apply text-content;
    box-shadow:
      0 0 15px rgba(var(--color-brand-500-rgb), 0.2),
      inset 0 0 15px rgba(var(--color-brand-500-rgb), 0.1);
    position: relative;
  }

  .nav-item-active::after {
    height: 70%;
    box-shadow: 0 0 10px rgba(var(--color-brand-500-rgb), 0.6);
  }

  .nav-item-collapsed:hover {
    @apply text-content;
    background: rgba(var(--color-brand-500-rgb), 0.12);
    border-color: rgba(var(--color-brand-500-rgb), 0.3);
    transform: scale(1.05);
    box-shadow:
      0 4px 12px rgba(var(--color-brand-500-rgb), 0.2),
      0 0 20px rgba(var(--color-brand-500-rgb), 0.1);
  }

  .nav-item-collapsed.nav-item-active {
    background: linear-gradient(135deg, rgba(var(--color-brand-500-rgb), 0.25) 0%, rgba(var(--color-brand-500-rgb), 0.15) 100%);
    border-color: rgba(var(--color-brand-500-rgb), 0.5);
    box-shadow:
      0 4px 16px rgba(var(--color-brand-500-rgb), 0.3),
      0 0 25px rgba(var(--color-brand-500-rgb), 0.2),
      inset 0 0 15px rgba(var(--color-brand-500-rgb), 0.1);
    position: relative;
  }

  .nav-item-collapsed.nav-item-active .nav-item-icon {
    @apply text-content;
    animation: iconGlow 2s ease-in-out infinite alternate;
  }

  @keyframes iconGlow {
    0% {
      filter: drop-shadow(0 0 3px rgba(var(--color-brand-500-rgb), 0.6));
    }
    100% {
      filter: drop-shadow(0 0 6px rgba(var(--color-brand-500-rgb), 0.8));
    }
  }

  .nav-item-active .nav-item-icon {
    @apply text-content;
  }
</style>
