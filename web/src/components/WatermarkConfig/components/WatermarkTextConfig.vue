<script setup lang="ts">
  import type { WatermarkConfig } from '../types'

  defineOptions({
    name: 'WatermarkTextConfig',
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
      <h4 class="section-title">{{ $t('components.watermark.textConfig.title') }}</h4>
    </div>
    <div class="section-content">
      <div class="config-item">
        <label class="config-label">{{ $t('components.watermark.textConfig.text') }}</label>
        <cyberInput
          :model-value="config.text"
          :placeholder="$t('components.watermark.textConfig.textPlaceholder')"
          @update:model-value="updateConfig('text', $event)"
        />
      </div>

      <div class="config-item">
        <label class="config-label">{{ $t('components.watermark.textConfig.fontSize') }}</label>
        <div class="slider-control">
          <cyberSlider
            :model-value="config.fontSize"
            :min="12"
            :max="200"
            @update:model-value="updateConfig('fontSize', $event)"
          />
          <span class="slider-value">{{ config.fontSize }}px</span>
        </div>
      </div>

      <div class="config-item">
        <label class="config-label">{{ $t('components.watermark.textConfig.fontColor') }}</label>
        <cyberColorPicker :model-value="config.fontColor" @update:model-value="updateConfig('fontColor', $event)" />
      </div>

      <div class="config-item">
        <label class="config-label">{{ $t('components.watermark.textConfig.fontWeight') }}</label>
        <cyberRadioGroup
          :model-value="config.fontWeight ?? 'normal'"
          :options="[
            { label: $t('components.watermark.textConfig.normal'), value: 'normal' },
            { label: $t('components.watermark.textConfig.bold'), value: 'bold' },
          ]"
          @update:model-value="updateConfig('fontWeight', $event)"
        />
      </div>
    </div>
  </div>
</template>

<style scoped>
  @import '../styles/common.css';

  .config-item :deep(.cyber-input) {
    @apply w-full;
  }

  .slider-control :deep(.cyber-slider) {
    @apply flex-1;
  }

  .color-text-input {
    @apply flex-1;
  }

  .config-item :deep(.cyber-radio-group) {
    @apply flex space-x-4;
  }

  @media (max-width: 640px) {
    .config-item :deep(.cyber-radio-group) {
      @apply flex-col space-x-0 space-y-2;
    }
  }
</style>
