<script setup lang="ts">
  import { computed, onMounted, watch } from 'vue'
  import { TIMING } from '@/constants'
  import type { ToastProps } from './types'

  defineOptions({
    name: 'CyberToast',
  })

  const props = withDefaults(defineProps<ToastProps>(), {
    type: 'info',
    duration: TIMING.TOAST.DEFAULT_DURATION,
    visible: true,
  })

  const emit = defineEmits<{
    update: [value: boolean]
    close: []
  }>()

  const closeToast = () => {
    emit('close')
    emit('update', false)
  }

  onMounted(() => {
    emit('update', true)
  })

  watch(
    () => props.visible,
    (val) => {
      emit('update', val)
    }
  )

  const toastClass = computed(() => ({
    [`toast-${props.type}`]: true,
  }))

  const iconClass = computed(() => {
    const iconMap: Record<string, string> = {
      success: 'fas fa-check-circle',
      error: 'fas fa-times-circle',
      warning: 'fas fa-exclamation-triangle',
      info: 'fas fa-info-circle',
    }
    return iconMap[props.type] || 'fas fa-info-circle'
  })
</script>

<template>
  <div v-show="visible" class="toast" :class="toastClass">
    <div v-if="type !== 'default'" class="toast-icon">
      <i :class="iconClass" />
    </div>
    <div class="toast-content">
      {{ message }}
    </div>
    <div class="toast-close" @click="closeToast">
      <i class="fas fa-times" />
    </div>
  </div>
</template>

<style scoped>
  .toast {
    @apply relative flex w-full min-w-[300px] max-w-[500px] items-center overflow-hidden rounded border-l-4 px-4 py-3 text-sm text-content transition-all duration-300;
    background-color: rgba(var(--color-background-800-rgb), 0.9);
    box-shadow:
      0 0 15px rgba(var(--color-shadow-rgb), 0.3),
      0 0 8px rgba(var(--color-shadow-rgb), 0.5);
  }

  @media (max-width: 768px) {
    .toast {
      @apply max-w-[90%];
    }
  }

  .toast::after {
    @apply absolute bottom-0 right-0 top-0 w-1.5;
    content: '';
    background: inherit;
    filter: brightness(1.5);
  }

  .toast-icon {
    @apply mr-3 flex-shrink-0 text-lg;
  }

  .toast-content {
    @apply mr-3 flex-1 break-words leading-relaxed;
    text-shadow: 0 0 2px rgba(var(--color-content-rgb), 0.5);
  }

  .toast-close {
    @apply z-10 flex-shrink-0 cursor-pointer p-1 opacity-80 transition-all duration-200;
  }

  .toast-close:hover {
    @apply scale-110 opacity-100;
  }

  .toast-success {
    border-color: var(--color-success-500);
    box-shadow:
      0 0 15px rgba(var(--color-success-rgb), 0.3),
      0 0 8px rgba(var(--color-shadow-rgb), 0.5);
  }

  .toast-success .toast-icon {
    color: var(--color-success-500);
    text-shadow: 0 0 5px var(--color-success-500);
  }

  .toast-error {
    border-color: var(--color-error-500);
    box-shadow:
      0 0 15px rgba(var(--color-error-rgb), 0.3),
      0 0 8px rgba(var(--color-shadow-rgb), 0.5);
  }

  .toast-error .toast-icon {
    color: var(--color-error-500);
    text-shadow: 0 0 5px var(--color-error-500);
  }

  .toast-warning {
    border-color: var(--color-warning-500);
    box-shadow:
      0 0 15px rgba(var(--color-warning-rgb), 0.3),
      0 0 8px rgba(var(--color-shadow-rgb), 0.5);
  }

  .toast-warning .toast-icon {
    color: var(--color-warning-500);
    text-shadow: 0 0 5px var(--color-warning-500);
  }

  .toast-info {
    border-color: var(--color-brand-500);
    box-shadow:
      0 0 15px rgba(var(--color-brand-500-rgb), 0.3),
      0 0 8px rgba(var(--color-shadow-rgb), 0.5);
  }

  .toast-info .toast-icon {
    color: var(--color-brand-500);
    text-shadow: 0 0 5px var(--color-brand-500);
  }
</style>
