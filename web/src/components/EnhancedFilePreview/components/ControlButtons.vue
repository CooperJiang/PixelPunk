<script setup lang="ts">
  import { useTexts } from '@/composables/useTexts'
  import type { ControlButtonsProps, ControlButtonsEmits } from '../types'

  defineOptions({
    name: 'ControlButtons',
  })

  defineProps<ControlButtonsProps>()
  defineEmits<ControlButtonsEmits>()

  const { $t } = useTexts()
</script>

<template>
  <div class="control-buttons">
    <button class="enhanced-close-button" :title="$t('components.enhancedFilePreview.exitPreviewEsc')" @click="$emit('close')">
      <i class="fas fa-times" />
    </button>

    <button
      v-if="isBrowserFullscreen"
      class="enhanced-fullscreen-exit-button"
      :title="$t('components.enhancedFilePreview.exitFullscreenEsc')"
      @click="$emit('exitFullscreen')"
    >
      <i class="fas fa-compress" />
    </button>
  </div>
</template>

<style scoped>
  .control-buttons {
    @apply pointer-events-none absolute inset-0;
  }

  .enhanced-close-button,
  .enhanced-fullscreen-exit-button {
    @apply pointer-events-auto absolute z-20 flex h-10 w-10 items-center justify-center border border-[rgba(var(--color-brand-500-rgb),0.2)] bg-[rgba(var(--color-background-900-rgb),0.7)] text-content backdrop-blur-sm transition-all duration-200 hover:scale-110;
    border-radius: var(--radius-lg);
  }

  .enhanced-close-button:hover,
  .enhanced-fullscreen-exit-button:hover {
    background: rgba(var(--color-brand-500-rgb), 0.2);
    border-color: rgba(var(--color-brand-500-rgb), 0.4);
    @apply text-brand-400;
    box-shadow: 0 0 20px rgba(var(--color-brand-500-rgb), 0.3);
  }

  .enhanced-close-button {
    @apply right-6 top-6;
  }

  .enhanced-fullscreen-exit-button {
    @apply right-20 top-6;
  }

  .enhanced-close-button i,
  .enhanced-fullscreen-exit-button i {
    @apply text-lg;
  }

  @media (max-width: 768px) {
    .enhanced-close-button,
    .enhanced-fullscreen-exit-button {
      @apply h-9 w-9;
    }

    .enhanced-close-button {
      @apply right-4 top-4;
    }

    .enhanced-fullscreen-exit-button {
      @apply right-16 top-4;
    }

    .enhanced-close-button i,
    .enhanced-fullscreen-exit-button i {
      @apply text-base;
    }
  }
</style>
