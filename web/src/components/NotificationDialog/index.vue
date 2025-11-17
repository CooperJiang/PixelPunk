<script setup lang="ts">
  import { ref } from 'vue'
  import type { NotificationDialogProps } from './types'

  defineOptions({
    name: 'CyberNotificationDialog',
  })

  withDefaults(defineProps<NotificationDialogProps>(), {
    width: 360,
    showCloseButton: true,
  })

  const visible = ref(false)

  const show = () => {
    visible.value = true
  }

  const hide = () => {
    visible.value = false
  }

  const handleClose = () => {
    hide()
  }

  defineExpose({
    show,
    hide,
  })
</script>

<template>
  <Teleport to="body">
    <Transition
      enter-active-class="transition-all duration-300 ease-out"
      enter-from-class="opacity-0"
      enter-to-class="opacity-100"
      leave-active-class="transition-all duration-200 ease-in"
      leave-from-class="opacity-100"
      leave-to-class="opacity-0"
    >
      <div v-if="visible" class="notification-dialog-overlay" @click="handleClose">
        <Transition
          enter-active-class="transition-all duration-300 ease-out"
          enter-from-class="opacity-0 scale-95 translate-y-4"
          enter-to-class="opacity-100 scale-100 translate-y-0"
          leave-active-class="transition-all duration-200 ease-in"
          leave-from-class="opacity-100 scale-100 translate-y-0"
          leave-to-class="opacity-0 scale-95 translate-y-4"
        >
          <div
            v-if="visible"
            class="notification-dialog"
            :style="{ width: typeof width === 'number' ? `${width}px` : width }"
            @click.stop
          >
            <div class="notification-dialog-header">
              <h3 class="notification-dialog-title">{{ title }}</h3>
              <button v-if="showCloseButton" class="notification-dialog-close" @click="handleClose">
                <i class="fas fa-times" />
              </button>
            </div>

            <div class="notification-dialog-content">
              <slot />
            </div>
          </div>
        </Transition>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
  .notification-dialog-overlay {
    @apply fixed inset-0 z-[9999] flex items-center justify-center p-4;
    background: rgba(var(--color-background-900-rgb), 0.6);
    backdrop-filter: blur(4px);
  }

  .notification-dialog {
    @apply overflow-hidden rounded-xl bg-background-800 shadow-2xl;
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.3);
    box-shadow:
      0 20px 40px rgba(var(--color-background-900-rgb), 0.8),
      0 0 30px rgba(var(--color-brand-500-rgb), 0.2),
      inset 0 1px 0 rgba(var(--color-white-rgb), 0.05);
    backdrop-filter: blur(20px) saturate(180%);
    max-width: 90vw;
    max-height: 90vh;
  }

  .notification-dialog-header {
    @apply flex items-center justify-between px-5 py-4;
  }

  .notification-dialog-title {
    @apply m-0 text-lg font-semibold leading-tight text-content-heading;
  }

  .notification-dialog-close {
    @apply flex h-8 w-8 cursor-pointer items-center justify-center rounded-lg border-0 transition-all duration-200;
    @apply bg-background-700 text-content-muted;
    font-size: 0.875rem;
  }

  .notification-dialog-close:hover {
    color: var(--color-error-400);
    background: rgba(var(--color-error-rgb), 0.12);
    transform: scale(1.05);
  }

  .notification-dialog-content {
    @apply px-5 py-4 leading-relaxed text-content;
  }
</style>
