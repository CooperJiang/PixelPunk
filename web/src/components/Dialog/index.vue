<script setup lang="ts">
  import { onMounted, onUnmounted, watch } from 'vue'
  import { useTexts } from '@/composables/useTexts'
  import type { DialogEmits, DialogProps } from './types'

  defineOptions({
    name: 'CyberDialog',
    inheritAttrs: false,
  })

  const { $t } = useTexts()

  const props = withDefaults(defineProps<DialogProps>(), {
    modelValue: false,
    title: '',
    width: '420px',
    height: 'auto',
    maxWidth: '90vw',
    maxHeight: '90vh',
    showDefaultFooter: true,
    showCloseButton: true,
    noPadding: false,
    noScroll: false,
    cancelText: '',
    confirmText: '',
    loading: false,
    showFooter: false,
    closeOnEsc: true,
    closeOnClickOverlay: true,
    hideBorder: false,
  })

  const emit = defineEmits<DialogEmits>()

  const dialogTitle = () => props.title || $t('dialog.title')
  const dialogCancelText = () => props.cancelText || $t('actions.cancel')
  const dialogConfirmText = () => props.confirmText || $t('actions.confirm')

  const handleOverlayClick = () => {
    if (props.closeOnClickOverlay) {
      closeDialog()
    }
  }
  const lockBodyScroll = () => {
    const y = window.scrollY || window.pageYOffset
    document.body.style.position = 'fixed'
    document.body.style.top = `-${y}px`
    document.body.style.width = '100%'
  }

  const unlockBodyScroll = () => {
    const scrollY = document.body.style.top
    document.body.style.position = ''
    document.body.style.top = ''
    document.body.style.width = ''
    window.scrollTo(0, parseInt(scrollY || '0') * -1)
  }

  const handleKeyDown = (event: KeyboardEvent) => {
    if (event.key === 'Escape' && props.closeOnEsc) {
      closeDialog()
    }
    if (event.key === 'Enter' && props.showDefaultFooter) {
      confirmDialog()
    }
  }

  watch(
    () => props.modelValue,
    (val) => {
      if (val) {
        lockBodyScroll()
      } else {
        unlockBodyScroll()
      }
    },
    { immediate: true }
  )

  onMounted(() => {
    if (props.closeOnEsc) {
      document.addEventListener('keydown', handleKeyDown)
    }
    if (props.modelValue) {
      lockBodyScroll()
    }
  })

  onUnmounted(() => {
    if (props.closeOnEsc) {
      document.removeEventListener('keydown', handleKeyDown)
    }
    unlockBodyScroll()
  })

  const closeDialog = () => {
    emit('update:modelValue', false)
    emit('close')
    emit('cancel')
  }

  const confirmDialog = () => {
    emit('confirm')
  }
</script>

<template>
  <Teleport to="body">
    <Transition name="cyber-dialog">
      <div v-if="modelValue" class="dialog-overlay" @click.self="handleOverlayClick">
        <div
          class="dialog-container cyber-panel"
          :style="{
            width: width,
            height: height,
            maxWidth: maxWidth,
            maxHeight: maxHeight,
            border: hideBorder ? 'none' : '1px solid var(--color-border-default)',
          }"
          v-bind="$attrs"
        >
          <div class="dialog-header">
            <slot name="header">
              <h3 class="dialog-title">{{ dialogTitle() }}</h3>
            </slot>
            <button v-if="showCloseButton" class="dialog-close-btn" @click="closeDialog">
              <i class="fas fa-times" />
            </button>
          </div>

          <div class="dialog-content" :class="{ 'no-padding': noPadding, 'no-scroll': noScroll }">
            <slot />
          </div>

          <div v-if="$slots.footer || showFooter">
            <slot name="footer" />
          </div>
          <div v-else-if="showDefaultFooter" class="dialog-footer">
            <CyberButton type="outlined" @click="closeDialog">
              {{ dialogCancelText() }}
            </CyberButton>
            <CyberButton type="primary" :loading="loading" @click="confirmDialog">
              {{ dialogConfirmText() }}
            </CyberButton>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
  .dialog-overlay {
    @apply fixed inset-0 flex items-center justify-center;
    background: rgba(0, 0, 0, 0.7);
    backdrop-filter: blur(12px);
    z-index: var(--z-index-drawer-modal);
  }

  .dialog-container {
    @apply relative flex min-w-[320px] flex-col overflow-hidden;
    background: linear-gradient(
      135deg,
      rgba(var(--color-background-900-rgb), 0.95) 0%,
      rgba(var(--color-background-800-rgb), 0.9) 100%
    );
    color: var(--color-content);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.3);
    border-radius: var(--radius-sm);
    box-shadow:
      0 8px 32px rgba(0, 0, 0, 0.4),
      0 0 0 1px rgba(var(--color-brand-500-rgb), 0.1) inset,
      0 0 80px rgba(var(--color-brand-500-rgb), 0.15);
    backdrop-filter: blur(20px);
  }

  .dialog-container::before {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    height: 2px;
    background: linear-gradient(90deg, transparent 0%, rgba(var(--color-brand-500-rgb), 0.6) 50%, transparent 100%);
    box-shadow: 0 0 20px rgba(var(--color-brand-500-rgb), 0.5);
  }

  .dialog-header {
    @apply flex items-center justify-between border-b px-5 py-4;
    border-color: rgba(var(--color-brand-500-rgb), 0.2);
    background: linear-gradient(135deg, rgba(var(--color-brand-500-rgb), 0.08) 0%, rgba(var(--color-brand-500-rgb), 0.03) 100%);
    position: relative;
  }

  .dialog-header::after {
    content: '';
    position: absolute;
    bottom: 0;
    left: 20px;
    right: 20px;
    height: 1px;
    background: linear-gradient(90deg, transparent 0%, rgba(var(--color-brand-500-rgb), 0.3) 50%, transparent 100%);
  }

  .dialog-title {
    @apply m-0 text-lg font-semibold;
    background: linear-gradient(135deg, var(--color-brand-500) 0%, rgba(var(--color-brand-500-rgb), 0.8) 100%);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    background-clip: text;
    letter-spacing: 0.5px;
    text-shadow: 0 0 20px rgba(var(--color-brand-500-rgb), 0.3);
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    max-width: calc(100% - 40px); /* 为关闭按钮留出空间 */
  }

  .dialog-close-btn {
    @apply flex cursor-pointer items-center justify-center border-none bg-transparent p-1 text-lg outline-none transition-all duration-300;
    color: var(--color-content-muted);
    border-radius: var(--radius-sm);
    width: 32px;
    height: 32px;
    border: 1px solid transparent;
  }

  .dialog-close-btn:hover {
    color: var(--color-error-500);
    background: linear-gradient(135deg, rgba(var(--color-error-rgb), 0.15) 0%, rgba(var(--color-error-rgb), 0.08) 100%);
    border-color: rgba(var(--color-error-rgb), 0.3);
    box-shadow: 0 0 12px rgba(var(--color-error-rgb), 0.3);
    transform: scale(1.05);
  }

  .dialog-content {
    @apply flex-grow overflow-y-auto px-5 py-5;
    background: linear-gradient(
      180deg,
      rgba(var(--color-background-900-rgb), 0.3) 0%,
      rgba(var(--color-background-900-rgb), 0.5) 100%
    );
    position: relative;
  }

  .dialog-content::-webkit-scrollbar {
    width: 6px;
  }

  .dialog-content::-webkit-scrollbar-track {
    background: rgba(var(--color-brand-500-rgb), 0.05);
    border-radius: var(--radius-sm);
  }

  .dialog-content::-webkit-scrollbar-thumb {
    background: linear-gradient(180deg, rgba(var(--color-brand-500-rgb), 0.6) 0%, rgba(var(--color-brand-500-rgb), 0.4) 100%);
    border-radius: var(--radius-sm);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.2);
  }

  .dialog-content::-webkit-scrollbar-thumb:hover {
    background: var(--color-brand-500);
    box-shadow: 0 0 8px rgba(var(--color-brand-500-rgb), 0.5);
  }

  .dialog-content.no-padding {
    @apply p-0;
  }

  .dialog-content.no-scroll {
    @apply overflow-hidden;
  }

  .dialog-footer {
    @apply flex justify-end gap-3 border-t px-5 py-4;
    border-color: rgba(var(--color-brand-500-rgb), 0.2);
    background: linear-gradient(135deg, rgba(var(--color-brand-500-rgb), 0.05) 0%, rgba(var(--color-brand-500-rgb), 0.02) 100%);
    position: relative;
  }

  .dialog-footer::before {
    content: '';
    position: absolute;
    top: 0;
    left: 20px;
    right: 20px;
    height: 1px;
    background: linear-gradient(90deg, transparent 0%, rgba(var(--color-brand-500-rgb), 0.3) 50%, transparent 100%);
  }

  .cyber-dialog-enter-active,
  .cyber-dialog-leave-active {
    @apply transition-all duration-300;
  }

  .cyber-dialog-enter-from,
  .cyber-dialog-leave-to {
    @apply opacity-0;
  }

  .cyber-dialog-enter-active .dialog-container {
    animation: dialogIn 0.4s cubic-bezier(0.34, 1.56, 0.64, 1);
  }

  .cyber-dialog-leave-active .dialog-container {
    animation: dialogOut 0.3s cubic-bezier(0.4, 0, 0.6, 1) forwards;
  }

  @keyframes dialogIn {
    0% {
      opacity: 0;
      transform: scale(0.85) translateY(-30px);
      filter: blur(4px);
    }
    60% {
      opacity: 1;
      transform: scale(1.02) translateY(5px);
      filter: blur(0);
    }
    100% {
      opacity: 1;
      transform: scale(1) translateY(0);
      filter: blur(0);
    }
  }

  @keyframes dialogOut {
    0% {
      opacity: 1;
      transform: scale(1) translateY(0);
      filter: blur(0);
    }
    100% {
      opacity: 0;
      transform: scale(0.9) translateY(15px);
      filter: blur(2px);
    }
  }

  @media (max-width: 768px) {
    .dialog-overlay {
      padding: 16px;
      backdrop-filter: blur(8px);
    }

    .dialog-container {
      @apply min-w-[280px] max-w-[calc(100vw-32px)];
      margin: 0 auto;
      border-radius: var(--radius-sm);
    }

    .dialog-header {
      @apply px-4 py-3;
    }

    .dialog-title {
      @apply text-base;
    }

    .dialog-content {
      @apply px-4 py-4;
    }

    .dialog-footer {
      @apply gap-3 px-4 py-3;
      flex-direction: column-reverse;
    }

    .dialog-footer :deep(.cyber-btn) {
      @apply w-full;
    }
  }

  @media (max-width: 480px) {
    .dialog-overlay {
      padding: 12px;
    }

    .dialog-container {
      @apply min-w-[260px] max-w-[calc(100vw-24px)];
      border-radius: var(--radius-sm);
    }

    .dialog-header {
      @apply px-3 py-2;
    }

    .dialog-content {
      @apply px-3 py-3;
    }

    .dialog-footer {
      @apply gap-2 px-3 py-2;
    }

    .dialog-close-btn {
      width: 28px;
      height: 28px;
      @apply text-sm;
    }
  }

  @media (prefers-reduced-motion: reduce) {
    .cyber-dialog-enter-active,
    .cyber-dialog-leave-active {
      @apply transition-none;
    }

    .cyber-dialog-enter-active .dialog-container,
    .cyber-dialog-leave-active .dialog-container {
      animation: none;
    }
  }
</style>
