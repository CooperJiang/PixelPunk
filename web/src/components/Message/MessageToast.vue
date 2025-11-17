<script setup lang="ts">
  import { computed, onMounted, onUnmounted, ref } from 'vue'
  import { useRouter } from 'vue-router'
  import { useTexts } from '@/composables/useTexts'

  interface Props {
    type: 'success' | 'error' | 'warning' | 'info'
    title: string
    message?: string
    duration?: number
    actionText?: string
    actionUrl?: string
    actionStyle?: 'primary' | 'secondary' | 'success' | 'warning' | 'danger'
    id?: string
  }

  const props = withDefaults(defineProps<Props>(), {
    type: 'info',
    duration: 5000,
    actionStyle: 'primary',
  })

  const emit = defineEmits<{
    close: []
  }>()

  const visible = ref(false)
  let timer: NodeJS.Timeout | null = null

  const router = useRouter()
  const { $t } = useTexts()

  const isActionable = computed(() => Boolean(props.actionText && props.actionUrl))

  const iconClass = computed(() => {
    const icons = {
      success: 'fas fa-check-circle',
      error: 'fas fa-times-circle',
      warning: 'fas fa-exclamation-triangle',
      info: 'fas fa-info-circle',
    }
    return icons[props.type]
  })

  const show = () => {
    visible.value = true

    if (props.duration > 0) {
      timer = setTimeout(() => {
        handleClose()
      }, props.duration)
    }
  }

  const handleClose = () => {
    if (timer) {
      clearTimeout(timer)
      timer = null
    }
    visible.value = false
  }

  const handleAfterLeave = () => {
    emit('close')
  }

  const _handleAction = () => {
    if (props.actionUrl) {
      if (props.actionUrl.startsWith('http')) {
        window.open(props.actionUrl, '_blank')
      } else {
        router.push(props.actionUrl)
      }
      handleClose()
    }
  }

  onMounted(() => {
    requestAnimationFrame(() => {
      show()
    })
  })

  onUnmounted(() => {
    if (timer) {
      clearTimeout(timer)
    }
  })

  defineExpose({
    close: handleClose,
  })
</script>

<template>
  <Transition
    enter-active-class="transition ease-out duration-300"
    enter-from-class="transform opacity-0 translate-y-2 sm:translate-x-2 sm:translate-y-0"
    enter-to-class="transform opacity-100 translate-y-0 sm:translate-x-0"
    leave-active-class="transition ease-in duration-200"
    leave-from-class="transform opacity-100 translate-y-0 sm:translate-x-0"
    leave-to-class="transform opacity-0 translate-y-2 sm:translate-x-2 sm:translate-y-0"
    @after-leave="handleAfterLeave"
  >
    <div
      v-if="visible"
      class="pointer-events-auto relative w-full max-w-sm sm:w-auto"
      role="alert"
      aria-live="assertive"
      aria-atomic="true"
    >
      <div class="message-toast" :class="[`toast-${type}`, { 'toast-actionable': isActionable }]">
        <div class="toast-content">
          <div class="toast-icon">
            <i :class="iconClass" class="text-lg" />
          </div>

          <div class="toast-text">
            <p class="toast-title">{{ title }}</p>
            <p v-if="message" class="toast-message">{{ message }}</p>
          </div>

          <button class="toast-close" :aria-label="$t('components.messageToast.close')" @click="handleClose">
            <i class="fas fa-times text-sm" />
          </button>
        </div>
      </div>
    </div>
  </Transition>
</template>

<style scoped>
  .message-toast {
    background: rgba(var(--color-background-900-rgb), 0.95);
    backdrop-filter: blur(20px);
    border-radius: var(--radius-sm);
    box-shadow:
      0 20px 25px -5px rgba(var(--color-background-900-rgb), 0.3),
      0 10px 10px -5px rgba(var(--color-background-900-rgb), 0.2),
      0 0 0 1px rgba(var(--color-white-rgb), 0.05);
    overflow: hidden;
    border: 1px solid var(--color-border-default);
  }

  .toast-content {
    display: flex;
    align-items: flex-start;
    gap: 0.75rem;
    padding: 1rem;
  }

  .toast-icon {
    flex-shrink: 0;
    margin-top: 0.125rem;
  }

  .toast-success .toast-icon {
    color: var(--color-success-500);
  }

  .toast-error .toast-icon {
    color: var(--color-error-500);
  }

  .toast-warning .toast-icon {
    color: var(--color-warning-500);
  }

  .toast-info .toast-icon {
    color: var(--color-brand-400);
  }

  .toast-text {
    flex: 1;
    min-width: 0;
  }

  .toast-title {
    font-size: 0.875rem;
    font-weight: 600;
    color: var(--color-content-default);
    margin: 0;
    line-height: 1.25;
  }

  .toast-message {
    font-size: 0.8rem;
    color: var(--color-content-muted);
    margin: 0.25rem 0 0 0;
    line-height: 1.4;
    word-wrap: break-word;
  }

  .toast-close {
    flex-shrink: 0;
    margin-top: -0.125rem;
    padding: 0.25rem;
    color: var(--color-content-muted);
    background: transparent;
    border: none;
    border-radius: var(--radius-sm);
    cursor: pointer;
    transition: all 0.2s ease;
  }

  .toast-close:hover {
    color: var(--color-content-default);
    background: rgba(var(--color-background-900-rgb), 0.3);
  }

  .toast-actions {
    border-top: 1px solid var(--color-border-subtle);
    padding: 0.75rem 1rem;
  }

  .toast-action-btn {
    width: 100%;
    padding: 0.5rem 1rem;
    font-size: 0.8rem;
    font-weight: 500;
    border: 1px solid;
    border-radius: var(--radius-sm);
    cursor: pointer;
    transition: all 0.2s ease;
  }

  .action-primary {
    background: var(--color-brand-500);
    border-color: var(--color-brand-500);
    color: var(--color-text-on-brand);
  }

  .action-primary:hover {
    background: var(--color-brand-400);
    border-color: var(--color-brand-400);
  }

  .action-secondary {
    background: rgba(var(--color-background-900-rgb), 0.25);
    border-color: rgba(var(--color-brand-500-rgb), 0.3);
    color: var(--color-content-default);
  }

  .action-secondary:hover {
    background: rgba(var(--color-background-900-rgb), 0.35);
    border-color: rgba(var(--color-brand-500-rgb), 0.4);
  }

  .action-success {
    background: var(--color-success-500);
    border-color: var(--color-success-500);
    color: var(--color-text-on-success);
  }

  .action-success:hover {
    background: var(--color-success-400);
    border-color: var(--color-success-400);
  }

  .action-warning {
    background: var(--color-warning-500);
    border-color: var(--color-warning-500);
    color: var(--color-text-on-warning);
  }

  .action-warning:hover {
    background: var(--color-warning-400);
    border-color: var(--color-warning-400);
  }

  .action-danger {
    background: var(--color-error-500);
    border-color: var(--color-error-500);
    color: var(--color-text-on-error);
  }

  .action-danger:hover {
    background: var(--color-error-400);
    border-color: var(--color-error-400);
  }

  .toast-success {
    border-left: 4px solid var(--color-success-500);
  }

  .toast-error {
    border-left: 4px solid var(--color-error-500);
  }

  .toast-warning {
    border-left: 4px solid var(--color-warning-500);
  }

  .toast-info {
    border-left: 4px solid var(--color-brand-400);
  }

  @media (max-width: 640px) {
    .message-toast {
      margin: 0 1rem;
    }

    .toast-content {
      padding: 0.875rem;
    }

    .toast-actions {
      padding: 0.625rem 0.875rem;
    }

    .toast-title {
      font-size: 0.8rem;
    }

    .toast-message {
      font-size: 0.75rem;
    }
  }
</style>
