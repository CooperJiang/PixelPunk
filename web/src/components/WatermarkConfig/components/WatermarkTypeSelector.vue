<script setup lang="ts">
  import type { WatermarkConfig } from '../types'

  defineOptions({
    name: 'WatermarkTypeSelector',
  })

  interface Props {
    config: WatermarkConfig
  }

  interface Emits {
    (e: 'update:config', config: WatermarkConfig): void
  }

  const props = defineProps<Props>()
  const emit = defineEmits<Emits>()

  const updateType = (type: 'text' | 'image') => {
    emit('update:config', { ...props.config, type })
  }
</script>

<template>
  <div class="section-group">
    <div class="section-header">
      <h4 class="section-title">{{ $t('components.watermark.type.label') }}</h4>
    </div>
    <div class="section-content">
      <div class="config-item">
        <div class="type-selector">
          <label class="type-option" :class="{ active: config.type === 'text' }">
            <input :checked="config.type === 'text'" type="radio" value="text" class="type-radio" @change="updateType('text')" />
            <div class="type-card">
              <i class="fas fa-font text-sm" />
              <span class="type-label">{{ $t('components.watermark.type.text') }}</span>
            </div>
          </label>
          <label class="type-option" :class="{ active: config.type === 'image' }">
            <input
              :checked="config.type === 'image'"
              type="radio"
              value="image"
              class="type-radio"
              @change="updateType('image')"
            />
            <div class="type-card">
              <i class="fas fa-image text-sm" />
              <span class="type-label">{{ $t('components.watermark.type.image') }}</span>
            </div>
          </label>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
  @import '../styles/common.css';

  .type-selector {
    @apply grid grid-cols-2 gap-3;
  }

  .type-option {
    @apply cursor-pointer;
  }

  .type-option.active .type-card {
    border-color: var(--color-brand-500);
    background: var(--color-brand-500);
    color: var(--color-text-on-brand);
  }

  .type-card {
    @apply flex flex-col items-center space-y-1 rounded-md border p-2 text-center transition-all duration-200;
    border-color: var(--color-border-default);
    background: var(--color-background-800);
    color: var(--color-content-default);
  }

  .type-card:hover {
    border-color: var(--color-hover-border);
    background: var(--color-hover-bg);
  }

  .type-radio {
    @apply sr-only;
  }

  .type-label {
    @apply text-xs font-medium;
  }

  @media (max-width: 640px) {
    .type-selector {
      @apply grid-cols-1;
    }
  }
</style>
