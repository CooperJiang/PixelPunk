<script setup lang="ts">
  import { computed } from 'vue'
  import { useLayoutStore } from '@/store/layout'
  import { useTexts } from '@/composables/useTexts'

  defineOptions({
    name: 'LayoutSwitcher',
  })

  const { $t } = useTexts()

  const layoutStore = useLayoutStore()

  const isTopLayout = computed(() => layoutStore.isTopLayout)
  const isLeftLayout = computed(() => layoutStore.isLeftLayout)

  const switchToTop = () => {
    layoutStore.setLayoutMode('top')
  }

  const switchToLeft = () => {
    layoutStore.setLayoutMode('left')
  }
</script>

<template>
  <div class="layout-switcher">
    <div class="switcher-title">{{ $t('components.layoutSwitcher.title') }}</div>

    <div class="switcher-buttons">
      <button class="layout-button" :class="{ active: isTopLayout }" @click="switchToTop">
        <i class="fas fa-align-justify" />
        <span>{{ $t('components.layoutSwitcher.topNav') }}</span>
      </button>

      <button class="layout-button" :class="{ active: isLeftLayout }" @click="switchToLeft">
        <i class="fas fa-grip-vertical" />
        <span>{{ $t('components.layoutSwitcher.sidebar') }}</span>
      </button>
    </div>
  </div>
</template>

<style scoped>
  .layout-switcher {
    @apply rounded-lg border border-subtle bg-background-800 p-4;
  }

  .switcher-title {
    @apply mb-3 text-sm font-medium text-content;
  }

  .switcher-buttons {
    @apply flex gap-2;
  }

  .layout-button {
    @apply flex items-center gap-2 rounded-lg border border-subtle bg-background-700 px-4 py-2 text-content-muted;
    transition: all 0.3s ease;
    flex: 1;
    justify-content: center;
  }

  .layout-button:hover {
    @apply border-content-muted bg-background-600 text-content;
  }

  .layout-button.active {
    background: rgba(var(--color-brand-500), 0.1);
    border-color: rgba(var(--color-brand-500), 0.3);
    color: var(--color-brand-500);
  }

  .layout-button.active:hover {
    background: rgba(var(--color-brand-500), 0.15);
  }

  .layout-button i {
    @apply text-sm;
  }

  .layout-button span {
    @apply text-sm font-medium;
  }
</style>
