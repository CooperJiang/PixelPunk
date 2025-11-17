<script setup lang="ts">
  import { useTexts } from '@/composables/useTexts'
  import type { ControlsHintProps } from '../types'

  defineOptions({
    name: 'ControlsHint',
  })

  defineProps<ControlsHintProps>()

  const { $t } = useTexts()
</script>

<template>
  <div class="enhanced-controls-hint" :class="{ 'auto-hide': !showControlsHint }">
    <div class="hint-item">
      <kbd>ESC</kbd> {{ $t(`components.enhancedFilePreview.${isBrowserFullscreen ? 'exitFullscreen' : 'exitPreview'}`) }}
    </div>
    <div v-if="!isBrowserFullscreen" class="hint-item"><kbd>F</kbd> {{ $t('components.enhancedFilePreview.enterFullscreen') }}</div>
    <div class="hint-item">
      <kbd>{{ $t('components.enhancedFilePreview.spaceKey') }}</kbd>
      {{ shouldUseFillMode ? $t('components.enhancedFilePreview.switchToFit') : $t('components.enhancedFilePreview.switchToFill')
      }}{{ $t('components.enhancedFilePreview.mode') }}
    </div>
    <div class="hint-item"><i class="fas fa-mouse" /> {{ $t('components.enhancedFilePreview.wheelZoom') }}</div>
    <div class="hint-item"><i class="fas fa-hand-rock" /> {{ $t('components.enhancedFilePreview.dragMove') }}</div>
  </div>
</template>

<style scoped>
  .enhanced-controls-hint {
    @apply pointer-events-none absolute bottom-6 left-6 z-10 space-y-2 rounded-lg border p-4 text-sm text-content backdrop-blur-md transition-opacity duration-300;
    border-color: rgba(var(--color-brand-500-rgb), 0.2);
    background: rgba(var(--color-background-900-rgb), 0.75);
  }

  .enhanced-controls-hint.auto-hide {
    @apply opacity-0;
  }

  .hint-item {
    @apply flex items-center gap-2 text-content;
  }

  .hint-item kbd {
    @apply rounded border px-2 py-1 font-mono text-xs text-content;
    border-color: rgba(var(--color-brand-500-rgb), 0.3);
    background: rgba(var(--color-brand-500-rgb), 0.1);
  }

  .hint-item i {
    @apply w-4 text-center;
  }
</style>
