<script setup lang="ts">
  import type { WatermarkConfig } from '../types'

  defineOptions({
    name: 'WatermarkEffectsConfig',
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
      <h4 class="section-title">{{ $t('components.watermark.effects.title') }}</h4>
    </div>
    <div class="section-content">
      <div class="config-item">
        <cyberCheckbox :model-value="config.shadow" @update:model-value="updateConfig('shadow', $event)">
          {{ $t('components.watermark.effects.shadow') }}
        </cyberCheckbox>
      </div>

      <template v-if="config.shadow">
        <div class="config-item">
          <label class="config-label">{{ $t('components.watermark.effects.shadowColor') }}</label>
          <cyberColorPicker :model-value="config.shadowColor" @update:model-value="updateConfig('shadowColor', $event)" />
        </div>

        <div class="config-item">
          <label class="config-label">{{ $t('components.watermark.effects.shadowBlur') }}</label>
          <div class="slider-control">
            <cyberSlider
              :model-value="config.shadowBlur"
              :min="0"
              :max="10"
              @update:model-value="updateConfig('shadowBlur', $event)"
            />
            <span class="slider-value">{{ config.shadowBlur }}</span>
          </div>
        </div>
      </template>
    </div>
  </div>
</template>

<style scoped>
  @import '../styles/common.css';

  .slider-control :deep(.cyber-slider) {
    @apply flex-1;
  }

  .color-text-input {
    @apply flex-1;
  }
</style>
