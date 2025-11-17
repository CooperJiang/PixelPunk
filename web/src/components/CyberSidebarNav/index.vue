<script setup lang="ts">
  import type { CyberSidebarNavProps, CyberSidebarNavEmits, TabItem } from './types'

  defineOptions({
    name: 'CyberSidebarNav',
  })

  defineProps<CyberSidebarNavProps>()
  const emit = defineEmits<CyberSidebarNavEmits>()

  const handleTabClick = (tab: TabItem) => {
    emit('tab-change', tab.key)
  }
</script>

<template>
  <nav class="nav-menu">
    <div
      v-for="tab in tabs"
      :key="tab.key"
      class="nav-item"
      :class="{ active: activeTab === tab.key }"
      @click="handleTabClick(tab)"
    >
      <i v-if="tab.icon" :class="tab.icon" class="nav-icon" />
      <span class="nav-label">{{ tab.name }}</span>
      <span v-if="tab.badge && tab.badge !== '0'" class="nav-badge">{{ tab.badge }}</span>
      <div class="nav-indicator" />
    </div>
  </nav>
</template>

<style scoped>
  .nav-menu {
    @apply py-4;
  }

  .nav-item {
    @apply relative my-0.5 flex cursor-pointer items-center gap-3 px-5 py-3;
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    border-radius: 0;
    color: var(--color-content-default-secondary);
  }

  .nav-item:hover:not(.active) {
    @apply text-content-heading;
    background: linear-gradient(90deg, rgba(var(--color-brand-500-rgb), 0.08) 0%, rgba(var(--color-brand-500-rgb), 0.02) 100%);
    transform: translateX(2px);
  }

  .nav-item.active {
    @apply font-semibold;
    color: var(--color-brand-500);
    background: linear-gradient(
      90deg,
      rgba(var(--color-brand-500-rgb), 0.12) 0%,
      rgba(var(--color-brand-500-rgb), 0.06) 50%,
      rgba(var(--color-brand-500-rgb), 0.03) 100%
    );
    text-shadow: 0 0 8px rgba(var(--color-brand-500-rgb), 0.3);
    transform: translateX(4px);
  }

  .nav-icon {
    @apply flex-shrink-0 text-base;
    transition: all 0.3s ease;
  }

  .nav-item.active .nav-icon {
    color: var(--color-brand-500);
    filter: drop-shadow(0 0 6px rgba(var(--color-brand-500-rgb), 0.4));
  }

  .nav-label {
    @apply flex-1 whitespace-nowrap text-sm font-medium;
    transition: all 0.3s ease;
  }

  .nav-badge {
    @apply flex h-[18px] min-w-[18px] items-center justify-center rounded-full px-1.5 text-xs font-semibold text-content-heading;
    background: linear-gradient(135deg, var(--color-error-500), var(--color-purple-500));
    box-shadow:
      0 1px 3px rgba(0, 0, 0, 0.2),
      0 0 8px rgba(var(--color-error-rgb), 0.3);
  }

  .nav-indicator {
    @apply ml-auto h-1.5 w-1.5 rounded-full bg-transparent;
    transition: all 0.3s ease;
  }

  .nav-item.active .nav-indicator {
    background: linear-gradient(135deg, var(--color-brand-500), var(--color-purple-500));
    box-shadow:
      0 0 8px rgba(var(--color-brand-500-rgb), 0.6),
      0 0 16px rgba(var(--color-brand-500-rgb), 0.3);
  }

  @media (max-width: 1024px) {
    .nav-item {
      @apply px-4 py-2.5;
    }
  }

  @media (max-width: 768px) {
    .nav-menu {
      @apply flex gap-2 overflow-x-auto p-4;
      scrollbar-width: none;
      -ms-overflow-style: none;
    }

    .nav-menu::-webkit-scrollbar {
      display: none;
    }

    .nav-item {
      @apply m-0 min-w-[80px] flex-shrink-0 flex-col gap-2 rounded-lg px-4 py-3 text-center;
      transform: none;
    }
  }
</style>
