<script setup lang="ts">
  import { computed } from 'vue'
  import { useTexts } from '@/composables/useTexts'
  import type { DuplicateBadgeProps } from './types'

  defineOptions({
    name: 'DuplicateBadge',
  })

  const { $t } = useTexts()

  const props = withDefaults(defineProps<DuplicateBadgeProps>(), {
    text: '',
    color: '', // 为空时使用主题变量
    icon: 'fas fa-clone',
    position: 'top-left',
  })

  const displayText = computed(() => props.text || $t('duplicateBadge.duplicate'))

  const positionClass = computed(() => {
    switch (props.position) {
      case 'top-left':
        return 'pos-top-left'
      case 'top-right':
        return 'pos-top-right'
      case 'bottom-left':
        return 'pos-bottom-left'
      case 'bottom-right':
        return 'pos-bottom-right'
      default:
        return 'pos-top-left'
    }
  })

  const badgeStyle = computed(() => ({
    backgroundColor: props.color || 'rgba(var(--color-brand-500-rgb), 0.9)',
  }))
</script>

<template>
  <div class="duplicate-badge" :class="positionClass" :style="badgeStyle">
    <i v-if="icon" :class="icon" class="mr-1" />
    <span class="badge-text">{{ displayText }}</span>
  </div>
</template>

<style scoped>
  .duplicate-badge {
    @apply absolute z-20 inline-flex items-center gap-1.5 px-2 py-0.5 text-xs leading-none tracking-wide;
    height: 20px;
    border-radius: var(--radius-sm);
    color: var(--color-text-on-brand);
    font-weight: 500;
    box-shadow:
      0 1px 2px rgba(0, 0, 0, 0.2),
      0 0 0 1px rgba(var(--color-content-rgb), 0.1) inset;
    backdrop-filter: blur(4px);
    border: 1px solid rgba(var(--color-content-rgb), 0.15);
  }

  .duplicate-badge i {
    font-size: 10px;
    opacity: 0.95;
  }

  .badge-text {
    font-weight: 500;
    letter-spacing: 0.3px;
    font-size: 11px;
  }

  .pos-top-left {
    @apply left-2 top-2;
  }

  .pos-top-right {
    @apply right-2 top-2;
  }

  .pos-bottom-left {
    @apply bottom-2 left-2;
  }

  .pos-bottom-right {
    @apply bottom-2 right-2;
  }
</style>
