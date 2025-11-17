<script setup lang="ts">
  import type { TagProps } from './types'

  defineOptions({
    name: 'CyberTag',
  })

  const props = withDefaults(defineProps<TagProps>(), {
    variant: 'primary',
    size: 'small',
    closable: false,
    disabled: false,
    round: false,
  })

  defineEmits<{
    (e: 'close'): void
    (e: 'click'): void
  }>()

  const getVariantClasses = () => {
    switch (props.variant) {
      case 'success':
        return 'tag-success'
      case 'warning':
        return 'tag-warning'
      case 'error':
        return 'tag-error'
      case 'info':
        return 'tag-info'
      case 'secondary':
        return 'tag-secondary'
      case 'outline':
        return 'bg-transparent text-brand-500 border-brand-500'
      case 'primary':
      default:
        return 'tag-primary'
    }
  }

  const getSizeClasses = () => {
    switch (props.size) {
      case 'large':
        return 'px-3 py-1.5 text-sm'
      case 'medium':
        return 'px-2.5 py-1 text-sm'
      case 'small':
      default:
        return 'px-1.5 py-0.5'
    }
  }
</script>

<template>
  <span
    class="cyber-tag inline-flex items-center gap-1 border transition-all duration-200"
    :class="[
      getVariantClasses(),
      getSizeClasses(),
      {
        'rounded-full': round,
        rounded: !round,
        'cursor-not-allowed opacity-50': disabled,
        'cursor-pointer hover:opacity-80': !disabled,
        'max-w-20 overflow-hidden text-ellipsis whitespace-nowrap': props.truncate,
      },
    ]"
    :title="truncate ? String($slots.default?.()[0]?.children || '') : undefined"
    @click="!disabled && $emit('click')"
  >
    <slot />
    <i v-if="closable && !disabled" class="fas fa-times cursor-pointer text-xs hover:opacity-70" @click.stop="$emit('close')" />
  </span>
</template>

<style scoped>
  .cyber-tag {
    font-weight: 500;
    letter-spacing: 0.025em;
    flex-shrink: 0;
    font-size: 10px;
    border-radius: var(--radius-sm);
  }

  .cyber-tag:hover:not(.cursor-not-allowed) {
    transform: scale(1.05);
  }

  .tag-primary {
    background: var(--color-badge-primary-bg);
    color: var(--color-badge-primary-text);
    border: 1px solid var(--color-badge-primary-border);
  }

  .tag-success {
    background: var(--color-badge-success-bg);
    color: var(--color-badge-success-text);
    border: 1px solid var(--color-badge-success-border);
  }

  .tag-warning {
    background: var(--color-badge-warning-bg);
    color: var(--color-badge-warning-text);
    border: 1px solid var(--color-badge-warning-border);
  }

  .tag-error {
    background: var(--color-badge-accent-bg);
    color: var(--color-badge-accent-text);
    border: 1px solid var(--color-badge-accent-border);
  }

  .tag-info {
    background: var(--color-badge-primary-bg);
    color: var(--color-badge-primary-text);
    border: 1px solid var(--color-badge-primary-border);
  }

  .tag-secondary {
    background: var(--color-badge-neutral-bg);
    color: var(--color-badge-neutral-text);
    border: 1px solid var(--color-badge-neutral-border);
  }

  .tag-primary:hover {
    background: rgba(var(--color-brand-500-rgb), 0.25);
  }

  .tag-success:hover {
    background: rgba(var(--color-success-rgb), 0.25);
  }

  .tag-warning:hover {
    background: rgba(var(--color-warning-rgb), 0.25);
  }

  .tag-error:hover {
    background: rgba(var(--color-error-rgb), 0.25);
  }

  .tag-info:hover {
    background: rgba(var(--color-brand-500-rgb), 0.25);
  }

  .tag-secondary:hover {
    background: rgba(var(--color-content-rgb), 0.15);
  }
</style>
