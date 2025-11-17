<script setup lang="ts">
  import { useTexts } from '@/composables/useTexts'
  import type { IndicatorsProps } from '../types'

  defineOptions({
    name: 'Indicators',
  })

  defineProps<IndicatorsProps>()

  const { $t } = useTexts()
</script>

<template>
  <div class="indicators">
    <div v-if="showZoomIndicator" class="enhanced-zoom-indicator">{{ Math.round(scale * 100) }}%</div>

    <div v-if="showModeIndicator" class="enhanced-mode-indicator">
      {{ shouldUseFillMode ? $t('components.enhancedFilePreview.fillMode') : $t('components.enhancedFilePreview.fitMode') }}
    </div>
  </div>
</template>

<style scoped>
  .indicators {
    @apply pointer-events-none absolute inset-0;
  }

  .enhanced-zoom-indicator,
  .enhanced-mode-indicator {
    @apply absolute left-1/2 top-1/2 -translate-x-1/2 -translate-y-1/2 rounded-lg border px-4 py-2 text-lg font-bold text-content backdrop-blur-md;
    border-color: rgba(var(--color-brand-500-rgb), 0.3);
    background: rgba(var(--color-background-900-rgb), 0.8);
    animation: indicator-fade-in 0.3s ease-out;
  }

  .enhanced-zoom-indicator {
    @apply translate-y-[-70px] transform;
  }

  .enhanced-mode-indicator {
    @apply translate-y-[70px] transform;
  }

  @keyframes indicator-fade-in {
    from {
      @apply scale-75 opacity-0;
    }
    to {
      @apply scale-100 opacity-100;
    }
  }
</style>
