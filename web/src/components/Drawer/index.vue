<script setup lang="ts">
  import { computed, onUnmounted, watch } from 'vue'
  import type { DrawerEmits, DrawerProps } from './types'
  import { useTexts } from '@/composables/useTexts'

  defineOptions({
    name: 'CyberDrawer',
  })

  const props = withDefaults(defineProps<DrawerProps>(), {
    width: '600px',
    position: 'right',
    maskClosable: true,
    escClosable: true,
    showHeader: true,
    showKeyboardTip: true,
  })

  const emit = defineEmits<DrawerEmits>()
  const { $t } = useTexts()

  const visible = computed({
    get: () => props.modelValue,
    set: (value: boolean) => {
      emit('update:modelValue', value)
    },
  })

  const sizeStyle = computed(() => {
    const width = typeof props.width === 'number' ? `${props.width}px` : props.width
    return {
      width,
      maxWidth: 'min(90vw, 800px)',
      minWidth: 'min(400px, 90vw)',
    }
  })

  const positionClass = computed(() => {
    return props.position === 'left' ? 'drawer-left' : 'drawer-right'
  })

  const handleClose = () => {
    visible.value = false
    emit('close')
  }

  const handleMaskClick = () => {
    if (props.maskClosable) {
      handleClose()
    }
  }

  const handleKeydown = (event: KeyboardEvent) => {
    if (visible.value && props.escClosable && event.key === 'Escape') {
      event.preventDefault()
      handleClose()
    }
  }

  watch(
    visible,
    (newVisible) => {
      if (newVisible && props.escClosable) {
        document.addEventListener('keydown', handleKeydown)
      } else {
        document.removeEventListener('keydown', handleKeydown)
      }
    },
    { immediate: true }
  )

  onUnmounted(() => {
    document.removeEventListener('keydown', handleKeydown)
  })
</script>

<template>
  <Teleport to="body">
    <Transition name="drawer">
      <div v-if="visible" class="drawer-wrapper fixed inset-0 flex" @click="handleMaskClick">
        <div class="drawer-mask absolute inset-0 bg-overlay-heavy backdrop-blur-sm" />

        <div
          :style="sizeStyle"
          :class="['drawer-panel', positionClass]"
          class="relative flex h-full flex-col overflow-hidden"
          @click.stop
        >
          <template v-if="props.showHeader">
            <div v-if="$slots.header" class="drawer-header">
              <slot name="header" :close="handleClose" />
            </div>
            <div v-else class="drawer-header flex items-center justify-between px-5 py-4">
              <h2 class="drawer-title flex items-center gap-2 text-lg font-semibold text-content">
                <slot name="title-icon" />
                {{ title }}
              </h2>
              <button
                class="close-btn"
                :title="props.escClosable ? $t('components.drawer.closeWithEsc') : $t('components.drawer.close')"
                @click="handleClose"
              >
                <i class="fas fa-times" />
              </button>
            </div>
          </template>

          <div class="drawer-content flex-1 overflow-y-auto">
            <slot />
          </div>

          <div v-if="$slots.footer" class="border-t border-default">
            <slot name="footer" />
          </div>

          <!-- Keyboard Tip 区域 - 可通过 showKeyboardTip 配置隐藏 -->
          <div v-else-if="props.escClosable && props.showKeyboardTip" class="drawer-keyboard-tip">
            <div class="flex items-center justify-center text-xs">
              <i class="fas fa-keyboard mr-2" />
              <span>{{ $t('components.drawer.keyboardHint', { key: 'ESC' }) }}</span>
            </div>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
  .drawer-wrapper {
    z-index: var(--z-index-drawer);
  }

  .drawer-panel {
    background: rgba(var(--color-background-900-rgb), 0.98);
    backdrop-filter: blur(20px);
    position: relative;
  }

  .drawer-right {
    margin-left: auto;
    border-left: 1px solid rgba(var(--color-brand-500-rgb), 0.3);
    box-shadow:
      -4px 0 30px rgba(var(--color-brand-500-rgb), 0.2),
      -2px 0 10px rgba(0, 0, 0, 0.3),
      inset 1px 0 0 rgba(var(--color-brand-500-rgb), 0.1);
  }

  .drawer-left {
    margin-right: auto;
    border-right: 1px solid rgba(var(--color-brand-500-rgb), 0.3);
    box-shadow:
      4px 0 30px rgba(var(--color-brand-500-rgb), 0.2),
      2px 0 10px rgba(0, 0, 0, 0.3),
      inset -1px 0 0 rgba(var(--color-brand-500-rgb), 0.1);
  }

  @keyframes drawerBorderGlow {
    0%,
    100% {
      opacity: 0.5;
    }
    50% {
      opacity: 0.8;
    }
  }

  .drawer-header {
    border-bottom: 1px solid rgba(var(--color-brand-500-rgb), 0.12);
    background: linear-gradient(
      180deg,
      rgba(var(--color-brand-500-rgb), 0.05) 0%,
      rgba(var(--color-background-800-rgb), 0.98) 100%
    );
    position: relative;
    backdrop-filter: blur(10px);
  }

  .drawer-header::after {
    content: '';
    position: absolute;
    bottom: 0;
    left: 0;
    right: 0;
    height: 1px;
    background: linear-gradient(90deg, transparent, rgba(var(--color-brand-500-rgb), 0.5), transparent);
    opacity: 0.3;
  }

  .drawer-title {
    color: var(--color-content-heading);
    text-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
    letter-spacing: 0.025em;
  }

  .close-btn {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 32px;
    height: 32px;
    border-radius: var(--radius-sm);
    border: 1px solid var(--color-border-default);
    background: var(--color-hover-bg-neutral);
    color: rgba(var(--color-content-rgb), 0.6);
    transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
    cursor: pointer;
    position: relative;
    font-size: 14px;
    flex-shrink: 0;
  }

  .close-btn:hover {
    border-color: rgba(var(--color-error-rgb), 0.4);
    background: rgba(var(--color-error-rgb), 0.2);
    color: var(--color-error-500);
    transform: scale(1.05);
  }

  .close-btn:active {
    transform: scale(0.95);
  }

  .drawer-content {
    background: rgba(var(--color-background-900-rgb), 0.5);
  }

  .drawer-keyboard-tip {
    border-top: 1px solid rgba(var(--color-brand-500-rgb), 0.12);
    background: linear-gradient(
      180deg,
      rgba(var(--color-background-800-rgb), 0.98) 0%,
      rgba(var(--color-brand-500-rgb), 0.03) 100%
    );
    padding: 0.5rem 1.25rem;
    position: relative;
    color: var(--color-brand-400);
  }

  .drawer-keyboard-tip i {
    animation: keyboardPulse 2s ease-in-out infinite;
  }

  @keyframes keyboardPulse {
    0%,
    100% {
      transform: scale(1);
      opacity: 0.8;
    }
    50% {
      transform: scale(1.1);
      opacity: 1;
    }
  }

  .drawer-enter-active,
  .drawer-leave-active {
    @apply transition-all duration-300;
    transition-timing-function: cubic-bezier(0.25, 0.8, 0.25, 1);
  }

  .drawer-enter-from,
  .drawer-leave-to {
    @apply opacity-0;
  }

  .drawer-enter-to,
  .drawer-leave-from {
    @apply opacity-100;
  }

  .drawer-enter-active .drawer-mask,
  .drawer-leave-active .drawer-mask {
    @apply transition-opacity duration-300 ease-in-out;
  }

  .drawer-enter-from .drawer-mask,
  .drawer-leave-to .drawer-mask {
    @apply opacity-0;
  }

  .drawer-enter-to .drawer-mask,
  .drawer-leave-from .drawer-mask {
    @apply opacity-100;
  }

  .drawer-enter-active .drawer-panel,
  .drawer-leave-active .drawer-panel {
    @apply transition-transform duration-300;
    transition-timing-function: cubic-bezier(0.25, 0.8, 0.25, 1);
  }

  .drawer-enter-from .drawer-right {
    @apply translate-x-full;
  }

  .drawer-leave-to .drawer-right {
    @apply translate-x-full;
  }

  .drawer-enter-from .drawer-left {
    @apply -translate-x-full;
  }

  .drawer-leave-to .drawer-left {
    @apply -translate-x-full;
  }

  .drawer-enter-to .drawer-panel,
  .drawer-leave-from .drawer-panel {
    @apply translate-x-0;
  }

  .drawer-enter-active .drawer-content {
    animation: drawerContentIn 0.4s ease 0.15s both;
  }

  .drawer-content::-webkit-scrollbar {
    @apply w-2;
  }

  .drawer-content::-webkit-scrollbar-track {
    @apply rounded;
    background: rgba(var(--color-brand-500-rgb), 0.08);
  }

  .drawer-content::-webkit-scrollbar-thumb {
    @apply rounded;
    background: rgba(var(--color-brand-500-rgb), 0.3);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.1);
    transition: background 0.2s ease;
  }

  .drawer-content::-webkit-scrollbar-thumb:hover {
    background: rgba(var(--color-brand-500-rgb), 0.5);
    box-shadow: 0 0 8px rgba(var(--color-brand-500-rgb), 0.3);
  }

  .drawer-keyboard-tip .kbd {
    @apply mx-0.5 inline-block rounded-sm px-1 py-0.5 text-[10px] font-medium leading-none;
    @apply border shadow-sm;
    color: var(--color-content-default);
    background: rgba(var(--color-brand-500-rgb), 0.15);
    border-color: rgba(var(--color-brand-500-rgb), 0.3);
    box-shadow: 0 1px 3px rgba(var(--color-brand-500-rgb), 0.2);
  }

  @media (max-width: 768px) {
    .drawer-enter-active,
    .drawer-leave-active {
      @apply duration-[250ms];
    }

    .drawer-enter-active .drawer-panel,
    .drawer-leave-active .drawer-panel {
      @apply duration-[250ms];
    }

    .drawer-keyboard-tip {
      @apply hidden;
    }
  }
</style>
