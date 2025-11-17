<script setup lang="ts">
  import { useTexts } from '@/composables/useTexts'

  const { $t } = useTexts()

  const _props = defineProps<{
    horizontalGap: number
    verticalGap: number
    strength3D: number
    clipGap: number
    imageCount: number
    isLoading: boolean
  }>()

  const emit = defineEmits<{
    updateHorizontalGap: [value: number]
    updateVerticalGap: [value: number]
    updateStrength3D: [value: number]
    updateClipGap: [value: number]
  }>()

  const updateHorizontalGap = (e: Event) => {
    const value = parseFloat((e.target as HTMLInputElement).value)
    emit('updateHorizontalGap', value)
  }

  const updateVerticalGap = (e: Event) => {
    const value = parseFloat((e.target as HTMLInputElement).value)
    emit('updateVerticalGap', value)
  }

  const updateStrength3D = (e: Event) => {
    const value = parseFloat((e.target as HTMLInputElement).value)
    emit('updateStrength3D', value)
  }

  const updateClipGap = (e: Event) => {
    const value = parseFloat((e.target as HTMLInputElement).value)
    emit('updateClipGap', value)
  }
</script>

<template>
  <div class="controls">
    <div class="control-group">
      <label class="control-label">{{ $t('hive.controls.horizontalGap') }}</label>
      <input
        type="range"
        class="control-slider"
        :min="0"
        :max="0.3"
        :step="0.01"
        :value="horizontalGap"
        @input="updateHorizontalGap"
      />
      <span class="control-value">{{ horizontalGap.toFixed(2) }}</span>
    </div>

    <div class="control-group">
      <label class="control-label">{{ $t('hive.controls.verticalGap') }}</label>
      <input
        type="range"
        class="control-slider"
        :min="0"
        :max="0.3"
        :step="0.01"
        :value="verticalGap"
        @input="updateVerticalGap"
      />
      <span class="control-value">{{ verticalGap.toFixed(2) }}</span>
    </div>

    <div class="control-group">
      <label class="control-label">{{ $t('hive.controls.strength3D') }}</label>
      <input
        type="range"
        class="control-slider"
        :min="0.5"
        :max="2.0"
        :step="0.1"
        :value="strength3D"
        @input="updateStrength3D"
      />
      <span class="control-value">{{ strength3D.toFixed(1) }}</span>
    </div>

    <div class="control-group">
      <label class="control-label">{{ $t('hive.controls.clipGap') }}</label>
      <input type="range" class="control-slider" :min="0" :max="10" :step="1" :value="clipGap" @input="updateClipGap" />
      <span class="control-value">{{ clipGap.toFixed(0) }}</span>
    </div>

    <div class="stats-section">
      <div>{{ $t('hive.controls.mode') }}</div>
      <div>
        {{ $t('hive.controls.loaded') }}: <span class="highlight">{{ imageCount }}</span> {{ $t('hive.controls.resources') }}
      </div>
      <div v-if="isLoading" class="loading-text">
        <div class="spinner" />
        {{ $t('hive.controls.loading') }}
      </div>
    </div>
  </div>
</template>

<style scoped>
  .controls {
    position: fixed;
    bottom: 20px;
    right: 20px;
    background: rgba(0, 0, 0, 0.8);
    color: var(--color-content-heading);
    padding: 20px;
    border-radius: var(--radius-sm);
    backdrop-filter: blur(10px);
    font-size: 0.9rem;
    z-index: 100;
    min-width: 280px;
    border: 1px solid rgba(var(--color-info-rgb), 0.3);
  }

  .control-group {
    margin-bottom: 15px;
  }

  .control-group:last-child {
    margin-bottom: 0;
  }

  .control-label {
    display: block;
    margin-bottom: 5px;
    font-weight: bold;
    color: var(--color-white);
  }

  .control-slider {
    width: 100%;
    height: 6px;
    border-radius: var(--radius-sm);
    background: var(--color-hover-bg-neutral);
    outline: none;
    -webkit-appearance: none;
  }

  .control-slider::-webkit-slider-thumb {
    -webkit-appearance: none;
    appearance: none;
    width: 18px;
    height: 18px;
    border-radius: var(--radius-full);
    background: var(--color-info-400);
    cursor: pointer;
    border: 2px solid white;
  }

  .control-slider::-moz-range-thumb {
    width: 18px;
    height: 18px;
    border-radius: var(--radius-full);
    background: var(--color-info-400);
    cursor: pointer;
    border: 2px solid white;
  }

  .control-value {
    display: inline-block;
    margin-left: 10px;
    color: var(--color-info-400);
    font-size: 0.8rem;
    font-weight: bold;
  }

  .stats-section {
    margin-top: 15px;
    padding-top: 15px;
    border-top: 1px solid var(--color-border-default);
    font-size: 0.85rem;
    line-height: 1.4;
  }

  .highlight {
    color: var(--color-error-500);
    font-weight: bold;
  }

  .loading-text {
    display: flex;
    align-items: center;
    color: var(--color-info-400);
    margin-top: 5px;
  }

  .spinner {
    border: 2px solid var(--color-border-default);
    border-radius: var(--radius-full);
    border-top: 2px solid var(--color-info-400);
    width: 12px;
    height: 12px;
    animation: spin 1s linear infinite;
    margin-right: 8px;
  }

  @keyframes spin {
    0% {
      transform: rotate(0deg);
    }
    100% {
      transform: rotate(360deg);
    }
  }
</style>
