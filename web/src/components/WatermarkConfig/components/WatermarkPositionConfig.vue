<script setup lang="ts">
  import { computed } from 'vue'
  import type { WatermarkConfig } from '../types'

  defineOptions({
    name: 'WatermarkPositionConfig',
  })

  interface Props {
    config: WatermarkConfig
  }

  interface Emits {
    (e: 'update:config', config: WatermarkConfig): void
  }

  const props = defineProps<Props>()
  const emit = defineEmits<Emits>()

  const positionOptions = [
    { value: 'top-left', labelKey: 'watermark.position.topLeft', icon: 'fas fa-arrow-up text-xs' },
    { value: 'top-center', labelKey: 'watermark.position.topCenter', icon: 'fas fa-arrow-up text-xs' },
    { value: 'top-right', labelKey: 'watermark.position.topRight', icon: 'fas fa-arrow-up text-xs' },
    { value: 'middle-left', labelKey: 'watermark.position.middleLeft', icon: 'fas fa-arrow-left text-xs' },
    { value: 'middle-center', labelKey: 'watermark.position.center', icon: 'fas fa-circle text-xs' },
    { value: 'middle-right', labelKey: 'watermark.position.middleRight', icon: 'fas fa-arrow-right text-xs' },
    { value: 'bottom-left', labelKey: 'watermark.position.bottomLeft', icon: 'fas fa-arrow-down text-xs' },
    { value: 'bottom-center', labelKey: 'watermark.position.bottomCenter', icon: 'fas fa-arrow-down text-xs' },
    { value: 'bottom-right', labelKey: 'watermark.position.bottomRight', icon: 'fas fa-arrow-down text-xs' },
  ]

  const updatePosition = (position: string) => {
    emit('update:config', { ...props.config, position })
  }

  const updateOffsetX = (value: number) => {
    emit('update:config', { ...props.config, offsetX: value, offsetUnit: 'px' })
  }

  const updateOffsetY = (value: number) => {
    emit('update:config', { ...props.config, offsetY: value, offsetUnit: 'px' })
  }

  const offsetX = computed({
    get: () => props.config.offsetX ?? 20,
    set: (value) => updateOffsetX(value),
  })

  const offsetY = computed({
    get: () => props.config.offsetY ?? 20,
    set: (value) => updateOffsetY(value),
  })

  const showOffsetWarning = computed(() => {
    const x = offsetX.value
    const y = offsetY.value
    return x > 1000 || y > 1000 || x < -200 || y < -200
  })
</script>

<template>
  <div class="section-group">
    <div class="section-header">
      <h4 class="section-title">{{ $t('components.watermark.position.title') }}</h4>
    </div>
    <div class="section-content">
      <div class="config-item">
        <label class="config-label">{{ $t('components.watermark.position.label') }}</label>
        <div class="position-grid">
          <button
            v-for="pos in positionOptions"
            :key="pos.value"
            class="position-btn"
            :class="[config.position === pos.value ? 'active' : '']"
            :title="$t(pos.labelKey)"
            @click="updatePosition(pos.value)"
          >
            <i :class="pos.icon" />
          </button>
        </div>
      </div>

      <div class="config-item">
        <label class="config-label">
          <i class="fas fa-ruler mr-1.5 text-xs" />
          <span>{{ $t('components.watermark.offset.title') }}</span>
        </label>

        <div class="offset-controls">
          <div class="offset-input-group">
            <label class="input-label">{{ $t('components.watermark.offset.horizontal') }}</label>
            <cyberInput v-model.number="offsetX" type="number" :min="-500" :max="5000" class="offset-input" />
            <span class="input-unit">px</span>
          </div>

          <div class="offset-input-group">
            <label class="input-label">{{ $t('components.watermark.offset.vertical') }}</label>
            <cyberInput v-model.number="offsetY" type="number" :min="-500" :max="5000" class="offset-input" />
            <span class="input-unit">px</span>
          </div>
        </div>

        <p class="hint-text">
          <i class="fas fa-info-circle mr-1 text-xs" />
          {{ $t('components.watermark.offset.hint') }}
        </p>
        <p v-if="showOffsetWarning" class="warning-text">
          <i class="fas fa-exclamation-triangle mr-1 text-xs" />
          {{ $t('components.watermark.offset.warning') }}
        </p>
      </div>
    </div>
  </div>
</template>

<style scoped>
  .section-group {
    @apply rounded-lg border;
    border-color: var(--color-border-default);
    background: var(--color-background-700);
    backdrop-filter: blur(4px);
  }

  .section-header {
    @apply border-b px-3 py-2;
    border-color: var(--color-border-subtle);
    background: var(--color-background-800);
  }

  .section-title {
    @apply flex items-center text-sm font-semibold;
    color: var(--color-content-heading);
  }

  .section-content {
    @apply space-y-3 p-3;
  }

  .config-item {
    @apply space-y-1.5;
  }

  .config-label {
    @apply block text-sm font-medium;
    color: var(--color-content-default);
  }

  .position-grid {
    @apply grid grid-cols-3 gap-2;
  }

  .position-btn {
    @apply flex h-12 w-12 cursor-pointer items-center justify-center rounded-lg border transition-all duration-200;
    border-color: var(--color-border-default);
    background: var(--color-background-800);
    color: var(--color-brand-500);
  }

  .position-btn:hover {
    border-color: var(--color-hover-border);
    background: var(--color-hover-bg);
    color: var(--color-content-heading);
  }

  .position-btn.active {
    @apply scale-105 shadow-lg;
    border-color: var(--color-brand-500);
    background: var(--color-brand-500);
    color: var(--color-text-on-brand);
  }

  .offset-controls {
    @apply flex flex-col gap-2;
  }

  .offset-input-group {
    @apply flex items-center gap-2;
  }

  .input-label {
    @apply w-12 text-sm font-medium;
    color: var(--color-content-default);
  }

  .offset-input {
    @apply flex-1;
  }

  .offset-input :deep(input) {
    @apply h-8 text-sm;
  }

  .input-unit {
    @apply w-10 text-right text-xs font-medium;
    color: var(--color-content-subtle);
  }

  .unit-toggle {
    @apply mt-1 flex gap-1 rounded-lg border p-1;
    border-color: var(--color-border-default);
    background: var(--color-background-800);
  }

  .unit-btn {
    @apply flex-1 rounded px-2 py-1 text-xs font-medium transition-all duration-200;
    color: var(--color-content-default);
  }

  .unit-btn:hover {
    background: var(--color-hover-bg);
  }

  .unit-btn.active {
    background: var(--color-brand-500);
    color: var(--color-text-on-brand);
  }

  .hint-text {
    @apply mt-2 flex items-center text-xs;
    color: var(--color-brand-500);
  }

  .warning-text {
    @apply mt-2 flex items-center text-xs;
    color: var(--color-warning-500);
  }
</style>
