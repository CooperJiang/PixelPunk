<script setup lang="ts">
  import { computed } from 'vue'
  import type { BadgeProps } from './types'

  defineOptions({
    name: 'Badge',
  })

  const props = withDefaults(defineProps<BadgeProps>(), {
    text: '',
    color: '',
    icon: '',
    position: 'top-right',
    inline: false,
  })

  const positionClass = computed(() => {
    if (props.inline) {
      return ''
    }

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
        return 'pos-top-right'
    }
  })

  const badgeStyle = computed(() => {
    if (props.variant) {
      return {}
    }
    return {
      backgroundColor: props.color,
    }
  })
</script>

<template>
  <div
    class="cyber-badge"
    :class="[positionClass, props.variant ? `variant-${props.variant}` : 'variant-default', inline ? 'inline' : '']"
    :style="badgeStyle"
  >
    <i v-if="icon" :class="icon" class="mr-1" />
    <span class="badge-text">{{ text }}</span>
  </div>
</template>

<style scoped>
  .cyber-badge {
    @apply absolute z-20 inline-flex h-[18px] items-center gap-1 rounded px-1.5 py-0.5 text-[11px] leading-none tracking-wide shadow-lg;
  }

  .cyber-badge.inline {
    @apply static h-auto min-h-5 rounded-md border p-2 text-xs transition-all duration-200;
  }

  .cyber-badge.inline:hover {
    @apply -translate-y-0.5;
  }

  .cyber-badge.variant-default {
    @apply border-subtle bg-background-700 text-content;
  }

  .cyber-badge.variant-success {
    @apply border-badge-success-border bg-badge-success-bg text-badge-success-text;
  }

  .cyber-badge.variant-error,
  .cyber-badge.variant-danger {
    @apply border-badge-accent-border bg-badge-accent-bg text-badge-accent-text;
  }

  .cyber-badge.variant-warning {
    @apply border-badge-warning-border bg-badge-warning-bg text-badge-warning-text;
  }

  .cyber-badge.variant-info {
    @apply border-badge-primary-border bg-badge-primary-bg text-badge-primary-text;
  }

  .cyber-badge.variant-primary {
    @apply border-badge-primary-border bg-badge-primary-bg text-badge-primary-text;
  }

  .cyber-badge.variant-success:hover {
    @apply border-success-400 bg-success-200;
  }

  .cyber-badge.variant-error:hover,
  .cyber-badge.variant-danger:hover {
    @apply border-error-400 bg-error-200;
  }

  .cyber-badge.variant-warning:hover {
    @apply border-warning-400 bg-warning-200;
  }

  .cyber-badge.variant-info:hover,
  .cyber-badge.variant-primary:hover {
    @apply border-brand-400 bg-brand-200;
  }

  .badge-text {
    @apply font-semibold;
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
