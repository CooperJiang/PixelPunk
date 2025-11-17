<script setup lang="ts">
  import type { WatermarkConfig } from '../types'

  defineOptions({
    name: 'WatermarkStyleConfig',
  })

  interface Props {
    config: WatermarkConfig
  }

  interface Emits {
    (e: 'update:config', config: WatermarkConfig): void
  }

  const props = defineProps<Props>()
  const emit = defineEmits<Emits>()

  const updateConfig = (key: keyof WatermarkConfig, value: string | number | boolean) => {
    emit('update:config', { ...props.config, [key]: value })
  }
</script>

<template>
  <div class="section-group">
    <div class="section-header">
      <h4 class="section-title">{{ $t('components.watermark.style.title') }}</h4>
    </div>
    <div class="section-content">
      <div class="config-item">
        <label class="config-label">{{ $t('components.watermark.style.opacity') }}</label>
        <div class="slider-control">
          <cyberSlider
            :model-value="config.opacity"
            :min="0.1"
            :max="1"
            :step="0.1"
            @update:model-value="updateConfig('opacity', $event)"
          />
          <span class="slider-value">{{ Math.round(config.opacity * 100) }}%</span>
        </div>
      </div>

      <div class="config-item">
        <label class="config-label">
          <i class="fas fa-expand-arrows-alt mr-1.5 text-xs" />
          <span>{{ $t('components.watermark.style.relativeSize') }}</span>
        </label>
        <div class="slider-control">
          <cyberSlider
            :model-value="config.relativeSize || 0.15"
            :min="0.05"
            :max="0.5"
            :step="0.01"
            @update:model-value="updateConfig('relativeSize', $event)"
          />
          <span class="slider-value">{{ Math.round((config.relativeSize || 0.15) * 100) }}%</span>
        </div>
        <p class="hint-text">
          <i class="fas fa-info-circle mr-1 text-xs" />
          {{ $t('components.watermark.style.relativeSizeHint') }}
        </p>
      </div>

      <div class="config-item">
        <label class="config-label">{{ $t('components.watermark.style.scale') }}</label>
        <div class="slider-control">
          <cyberSlider
            :model-value="config.scale ?? 1.0"
            :min="0.01"
            :max="2"
            :step="0.01"
            @update:model-value="updateConfig('scale', $event)"
          />
          <span class="slider-value">{{ Math.round((config.scale ?? 1.0) * 100) }}%</span>
        </div>
        <p class="hint-text">
          <i class="fas fa-info-circle mr-1 text-xs" />
          {{ $t('components.watermark.style.scaleHint') }}
        </p>
      </div>

      <div class="config-item">
        <label class="config-label">{{ $t('components.watermark.style.rotation') }}</label>
        <div class="slider-control">
          <cyberSlider
            :model-value="config.rotation ?? 0"
            :min="-180"
            :max="180"
            :step="5"
            @update:model-value="updateConfig('rotation', $event)"
          />
          <span class="slider-value">{{ config.rotation ?? 0 }}°</span>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
  @import '../styles/common.css';

  .slider-control :deep(.cyber-slider) {
    @apply flex-1;
  }
</style>
