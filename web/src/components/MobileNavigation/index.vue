<script setup lang="ts">
  import { onMounted, ref, watch } from 'vue'
  import type { FolderInfo } from '@/api/types'
  import type { MobileNavigationProps, MobileNavigationEmits } from './types'
  import { useTexts } from '@/composables/useTexts'

  const props = withDefaults(defineProps<MobileNavigationProps>(), {
    items: () => [],
    title: '',
    hasPassword: false,
    showNavigation: true,
  })

  const emit = defineEmits<MobileNavigationEmits>()
  const { $t } = useTexts()
  const pathsContainer = ref<HTMLElement | null>(null)

  const handleClick = (folder: FolderInfo | null) => {
    emit('click', folder)
  }

  const scrollToEnd = () => {
    if (pathsContainer.value) {
      const container = pathsContainer.value
      setTimeout(() => {
        container.scrollLeft = container.scrollWidth
      }, 0)
    }
  }

  onMounted(() => {
    scrollToEnd()
  })

  watch(
    () => props.items?.length,
    () => {
      scrollToEnd()
    }
  )
</script>

<template>
  <div class="mobile-navigation">
    <div v-if="title" class="nav-header">
      <div class="nav-title">{{ title }}</div>
      <div v-if="hasPassword" class="nav-badge">
        <i class="fas fa-lock" />
      </div>
    </div>

    <div v-if="showNavigation" class="nav-paths">
      <div class="nav-home-item" :class="{ active: !items || items.length === 0 }" @click="handleClick(null)">
        <i class="fas fa-home" />
        <span class="path-name">{{ $t('navigation.home') }}</span>
      </div>

      <div ref="pathsContainer" class="nav-path-items-container">
        <div
          v-for="(item, index) in items"
          :key="index"
          class="nav-path-item"
          :class="{ active: index === items.length - 1 }"
          @click="handleClick(item)"
        >
          <div class="path-arrow">
            <i class="fas fa-chevron-right" />
          </div>
          <div class="path-content">
            <span class="path-name">{{ item.name }}</span>
          </div>
        </div>
      </div>
    </div>

    <div v-if="$slots.actions" class="nav-actions">
      <slot name="actions" />
    </div>
  </div>
</template>

<style scoped>
  .mobile-navigation {
    @apply flex w-full flex-col text-content;
    background: linear-gradient(to bottom, rgba(10, 14, 20, 0.95), rgba(10, 14, 20, 0.9));
    border-bottom: 1px solid rgba(var(--color-brand-500-rgb), 0.2);
  }

  .nav-header {
    @apply flex items-center justify-between px-3 py-2;
    border-bottom: 1px solid rgba(255, 255, 255, 0.05);
  }

  .nav-title {
    @apply max-w-[80%] overflow-hidden text-ellipsis whitespace-nowrap text-sm font-semibold;
    background: linear-gradient(90deg, var(--color-brand-500), var(--color-error-500));
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
  }

  .nav-badge {
    @apply text-xs;
    color: var(--color-brand-500);
  }

  .nav-paths {
    @apply flex w-full items-center overflow-x-auto whitespace-nowrap px-2 py-1.5;
    scrollbar-width: none;
    -ms-overflow-style: none;
  }

  .nav-paths::-webkit-scrollbar {
    @apply hidden;
  }

  .nav-home-item {
    @apply mr-1.5 flex flex-shrink-0 items-center rounded px-2 py-0.5 text-content-muted transition-all duration-200;
    background: rgba(var(--color-brand-500-rgb), 0.05);
  }

  .nav-home-item.active {
    background: rgba(var(--color-brand-500-rgb), 0.15);
    color: var(--color-brand-500);
  }

  .nav-home-item i {
    @apply mr-1 text-xs;
  }

  .nav-path-items-container {
    @apply flex overflow-x-auto;
    scrollbar-width: none;
    -ms-overflow-style: none;
  }

  .nav-path-items-container::-webkit-scrollbar {
    @apply hidden;
  }

  .nav-path-item {
    @apply flex flex-shrink-0 items-center;
  }

  .path-arrow {
    @apply mx-0.5 text-xs text-content-disabled;
  }

  .path-content {
    @apply max-w-[140px] overflow-hidden text-ellipsis whitespace-nowrap rounded px-2 py-0.5 transition-all duration-200;
    background: rgba(255, 255, 255, 0.05);
  }

  .nav-path-item.active .path-content {
    background: rgba(var(--color-brand-500-rgb), 0.15);
    color: var(--color-brand-500);
  }

  .nav-path-item:not(.active) .path-content:hover {
    background: rgba(255, 255, 255, 0.1);
  }

  .path-name {
    @apply text-xs;
  }

  .nav-actions {
    @apply flex gap-2 px-3 py-1.5;
    border-top: 1px solid rgba(255, 255, 255, 0.05);
  }
</style>
