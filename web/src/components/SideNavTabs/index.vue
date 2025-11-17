<script setup lang="ts">
  import { computed } from 'vue'
  import type { TabItem, SideNavTabsProps, SideNavTabsEmits } from './types'

  defineOptions({
    name: 'CyberSideNavTabs',
  })

  const props = withDefaults(defineProps<SideNavTabsProps>(), {
    activeTab: '',
    syncWithUrl: false,
    urlParamName: 'tab',
  })

  const emit = defineEmits<SideNavTabsEmits>()

  const visibleTabs = computed(() => props.tabs.filter((tab) => tab.visible !== false))

  const handleTabClick = (tab: TabItem) => {
    if (tab.disabled) {
      return
    }

    const tabKey = tab.key || tab.id || ''
    emit('update:activeTab', tabKey)
    emit('tab-change', tabKey, tab)
    emit('change', tabKey)
  }
</script>

<template>
  <div class="cyber-side-nav-layout">
    <div class="side-nav-header">
      <div class="header-content">
        <div class="header-left">
          <div class="title-section">
            <i v-if="icon" :class="icon" class="title-icon" />
            <h1 class="page-title">{{ title }}</h1>
            <div v-if="subtitle" class="page-subtitle">{{ subtitle }}</div>
          </div>
        </div>
        <div v-if="$slots.actions" class="header-actions">
          <slot name="actions" />
        </div>
      </div>
    </div>

    <div class="side-nav-main-content">
      <div class="side-nav-sidebar">
        <nav class="nav-menu">
          <div
            v-for="tab in visibleTabs"
            :key="tab.key || tab.id"
            class="nav-item"
            :class="{ active: activeTab === (tab.key || tab.id) }"
            @click="handleTabClick(tab)"
          >
            <i v-if="tab.icon" :class="tab.icon" class="nav-icon" />
            <span class="nav-label">{{ tab.name || tab.label }}</span>
            <span v-if="tab.badge" class="nav-badge">{{ tab.badge }}</span>
            <div class="nav-indicator" />
          </div>
        </nav>
      </div>

      <div class="side-nav-content-area">
        <div class="content-wrapper">
          <slot />
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
  .cyber-side-nav-layout {
    @apply flex h-screen flex-col overflow-hidden bg-background-900;
  }

  .side-nav-header {
    @apply relative z-[100] flex-shrink-0;
    background: linear-gradient(
      135deg,
      rgba(var(--color-brand-500-rgb), 0.08) 0%,
      rgba(var(--color-brand-600-rgb), 0.05) 50%,
      rgba(var(--color-warning-rgb), 0.03) 100%
    );
    backdrop-filter: blur(16px);
    border-bottom: 1px solid rgba(var(--color-brand-500-rgb), 0.2);
    box-shadow:
      0 2px 8px rgba(var(--color-shadow-rgb), 0.15),
      0 0 20px rgba(var(--color-brand-500-rgb), 0.1),
      inset 0 1px 0 rgba(var(--color-brand-500-rgb), 0.1);
  }
  .header-content {
    @apply relative z-[1] flex items-center justify-between gap-8 px-8 py-5;
    min-height: 72px;
  }

  .header-left {
    @apply flex flex-1 items-center;
  }

  .title-section {
    @apply flex items-center gap-3;
  }

  .title-icon {
    @apply text-[1.375rem] text-brand-500;
    filter: drop-shadow(0 0 8px rgba(var(--color-brand-500-rgb), 0.4));
  }

  .page-title {
    @apply m-0 text-2xl font-bold leading-tight text-content;
    text-shadow:
      0 0 15px rgba(var(--color-brand-500-rgb), 0.3),
      0 0 30px rgba(var(--color-brand-500-rgb), 0.15);
    background: linear-gradient(
      135deg,
      var(--color-content) 0%,
      rgba(var(--color-brand-500-rgb), 0.9) 50%,
      var(--color-content) 100%
    );
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    background-clip: text;
  }

  .page-subtitle {
    @apply ml-4 whitespace-nowrap text-xs font-normal text-content-muted opacity-90;
  }

  .header-actions {
    @apply flex flex-shrink-0 items-center gap-3;
  }

  .side-nav-main-content {
    @apply flex min-h-0 flex-1 overflow-hidden;
  }

  .side-nav-sidebar {
    @apply relative w-[220px] flex-shrink-0;
    background: linear-gradient(
      180deg,
      rgba(var(--color-brand-500-rgb), 0.03) 0%,
      rgba(var(--color-brand-600-rgb), 0.02) 50%,
      rgba(var(--color-warning-rgb), 0.01) 100%
    );
    border-right: 1px solid rgba(var(--color-border-rgb), 0.3);
  }
  .nav-menu {
    @apply py-4;
  }

  .nav-item {
    @apply relative my-0.5 flex cursor-pointer items-center gap-3 px-5 py-3 text-content-muted transition-all duration-300;
    transition-timing-function: cubic-bezier(0.4, 0, 0.2, 1);
  }
  .nav-item:hover:not(.disabled) {
    @apply translate-x-0.5 text-content;
    background: linear-gradient(90deg, rgba(var(--color-brand-500-rgb), 0.08) 0%, rgba(var(--color-brand-500-rgb), 0.02) 100%);
  }
  .nav-item.active {
    @apply translate-x-1 font-semibold text-brand-500;
    background: linear-gradient(
      90deg,
      rgba(var(--color-brand-500-rgb), 0.12) 0%,
      rgba(var(--color-brand-600-rgb), 0.06) 50%,
      rgba(var(--color-warning-rgb), 0.03) 100%
    );
    text-shadow: 0 0 8px rgba(var(--color-brand-500-rgb), 0.3);
  }
  .nav-icon {
    @apply flex-shrink-0 text-base transition-all duration-300;
  }

  .nav-item.active .nav-icon {
    @apply text-brand-500;
    filter: drop-shadow(0 0 6px rgba(var(--color-brand-500-rgb), 0.4));
  }

  .nav-label {
    @apply flex-1 whitespace-nowrap text-sm font-medium transition-all duration-300;
  }

  .nav-badge {
    @apply flex h-[18px] min-w-[18px] items-center justify-center rounded-full px-1.5 py-0.5 text-[0.7rem] font-semibold text-content-heading;
    background: linear-gradient(135deg, var(--color-warning-500), var(--color-brand-600));
    box-shadow:
      0 1px 3px rgba(var(--color-shadow-rgb), 0.2),
      0 0 8px rgba(var(--color-warning-rgb), 0.3);
  }

  .nav-indicator {
    @apply ml-auto h-1.5 w-1.5 rounded-full bg-transparent transition-all duration-300;
  }

  .nav-item.active .nav-indicator {
    background: linear-gradient(135deg, var(--color-brand-500), var(--color-brand-600));
    box-shadow:
      0 0 8px rgba(var(--color-brand-500-rgb), 0.6),
      0 0 16px rgba(var(--color-brand-500-rgb), 0.3);
  }

  .side-nav-content-area {
    @apply relative flex-1 overflow-y-auto overflow-x-hidden bg-background-900;
    scrollbar-width: thin;
    scrollbar-color: var(--color-border) transparent;
    scrollbar-gutter: stable;
  }

  .content-wrapper {
    @apply p-6;
    min-height: calc(100vh - 140px);
  }

  .side-nav-content-area::-webkit-scrollbar {
    @apply w-2 bg-transparent;
  }

  .side-nav-content-area::-webkit-scrollbar-track {
    @apply rounded;
    background: rgba(var(--color-border-rgb), 0.1);
  }

  .side-nav-content-area::-webkit-scrollbar-thumb {
    @apply rounded transition-all duration-300;
    background: var(--color-background-400);
  }

  .side-nav-content-area::-webkit-scrollbar-thumb:hover {
    @apply bg-brand-500;
  }

  @media (max-width: 1024px) {
    .side-nav-sidebar {
      @apply w-[200px];
    }

    .header-content {
      @apply px-5 py-3;
    }

    .content-wrapper {
      @apply p-6;
    }

    .nav-item {
      @apply px-4 py-2.5;
    }
  }

  @media (max-width: 768px) {
    .side-nav-main-content {
      @apply flex-col;
    }

    .side-nav-sidebar {
      @apply w-full border-r-0;
      border-bottom: 1px solid rgba(var(--color-border-rgb), 0.3);
    }

    .nav-menu {
      @apply flex gap-2 overflow-x-auto p-4;
      scrollbar-width: none;
      -ms-overflow-style: none;
    }

    .nav-menu::-webkit-scrollbar {
      @apply hidden;
    }

    .nav-item {
      @apply m-0 min-w-[80px] flex-shrink-0 flex-col gap-2 rounded-lg px-4 py-3 text-center;
    }
  }

  @media (max-width: 480px) {
    .page-title {
      @apply text-xl;
    }

    .title-icon {
      @apply text-lg;
    }
  }
</style>

<style>
  .cyber-side-nav-layout .header-actions .cyber-button {
    @apply min-h-8 px-3.5 py-2 text-xs;
  }

  @media (max-width: 480px) {
    .cyber-side-nav-layout .header-actions .cyber-button {
      @apply px-3 py-1.5 text-xs;
    }
  }
</style>
