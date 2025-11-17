<script setup lang="ts">
  import { computed } from 'vue'
  import type { SliderEmits, SliderProps } from './types'

  defineOptions({
    name: 'CyberSlider',
  })

  const props = withDefaults(defineProps<SliderProps>(), {
    min: 0,
    max: 100,
    step: 1,
    disabled: false,
    description: '',
    showValue: false,
    width: '100%', // 默认宽度
  })

  const emit = defineEmits<SliderEmits>()

  const localValue = computed({
    get: () => props.modelValue,
    set: (value) => emit('update:modelValue', parseFloat(value as unknown as string)),
  })
</script>

<template>
  <div class="cyber-slider-wrapper" :style="{ width: props.width }">
    <div class="flex items-center">
      <input v-model="localValue" type="range" :min="min" :max="max" :step="step" :disabled="disabled" class="cyber-slider" />
      <span v-if="showValue" class="cyber-slider-value">{{ localValue }}</span>
    </div>
    <p v-if="description" class="cyber-slider-description">{{ description }}</p>
  </div>
</template>

<style scoped>
  .cyber-slider-wrapper {
    @apply w-full;
  }

  .cyber-slider-value {
    @apply ml-3 text-sm;
    color: var(--color-content-muted);
  }

  .cyber-slider-description {
    @apply mt-2 text-xs;
    color: var(--color-content-muted);
  }

  .cyber-slider {
    @apply w-full rounded-md outline-none transition-all duration-300;
    -webkit-appearance: none;
    height: 6px;
    background: rgba(var(--color-brand-500-rgb), 0.1);
  }

  .cyber-slider::-webkit-slider-thumb {
    @apply h-[18px] w-[18px] cursor-pointer rounded-full transition-all duration-300;
    -webkit-appearance: none;
    appearance: none;
    background: rgba(var(--color-brand-500-rgb), 0.8);
    box-shadow: 0 0 10px rgba(var(--color-brand-500-rgb), 0.7);
  }

  .cyber-slider::-moz-range-thumb {
    @apply h-[18px] w-[18px] cursor-pointer rounded-full border-0 transition-all duration-300;
    background: rgba(var(--color-brand-500-rgb), 0.8);
    box-shadow: 0 0 10px rgba(var(--color-brand-500-rgb), 0.7);
  }

  .cyber-slider:hover::-webkit-slider-thumb {
    background: var(--color-brand-500);
    box-shadow: 0 0 15px rgba(var(--color-brand-500-rgb), 0.9);
  }

  .cyber-slider:hover::-moz-range-thumb {
    background: var(--color-brand-500);
    box-shadow: 0 0 15px rgba(var(--color-brand-500-rgb), 0.9);
  }

  .cyber-slider:disabled {
    @apply cursor-not-allowed opacity-50;
  }

  .cyber-slider:disabled::-webkit-slider-thumb {
    @apply cursor-not-allowed;
    background: rgba(var(--color-brand-500-rgb), 0.4);
    box-shadow: 0 0 5px rgba(var(--color-brand-500-rgb), 0.3);
  }

  .cyber-slider:disabled::-moz-range-thumb {
    @apply cursor-not-allowed;
    background: rgba(var(--color-brand-500-rgb), 0.4);
    box-shadow: 0 0 5px rgba(var(--color-brand-500-rgb), 0.3);
  }
</style>
