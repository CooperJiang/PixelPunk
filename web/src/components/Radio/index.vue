<script setup lang="ts">
  import type { RadioProps, RadioEmits } from './types'

  defineOptions({
    name: 'CyberRadio',
  })

  const props = withDefaults(defineProps<RadioProps>(), {
    disabled: false,
  })

  const emit = defineEmits<RadioEmits>()

  const handleChange = () => {
    if (!props.disabled) {
      emit('update:modelValue', props.value)
    }
  }
</script>

<template>
  <label class="cyber-radio-wrapper" :class="{ 'is-disabled': props.disabled }">
    <input
      type="radio"
      class="cyber-radio-input"
      :checked="props.modelValue === props.value"
      :disabled="props.disabled"
      @change="handleChange"
    />
    <span class="cyber-radio">
      <span class="cyber-radio-inner" />
    </span>
    <span class="cyber-radio-label">
      <slot />
    </span>
  </label>
</template>

<style scoped>
  .cyber-radio-wrapper {
    @apply relative mr-4 inline-flex cursor-pointer items-center text-sm text-content;
  }

  .cyber-radio-input {
    @apply absolute m-0 h-0 w-0 opacity-0;
  }

  .cyber-radio {
    @apply relative mr-2 inline-flex h-[18px] w-[18px] items-center justify-center rounded-full border transition-all duration-200;
    border-color: rgba(var(--color-brand-500-rgb), 0.5);
    background-color: rgba(var(--color-brand-500-rgb), 0.08);
  }

  .cyber-radio-inner {
    @apply relative block h-0 w-0 scale-0 rounded-full transition-all duration-200;
    background-color: var(--color-brand-500);
    transition-timing-function: cubic-bezier(0.78, 0.14, 0.15, 0.86);
  }

  .cyber-radio-input:checked + .cyber-radio {
    border-color: var(--color-brand-500);
    background-color: rgba(var(--color-brand-500-rgb), 0.15);
    box-shadow: 0 0 8px rgba(var(--color-brand-500-rgb), 0.25);
  }

  .cyber-radio-input:checked + .cyber-radio .cyber-radio-inner {
    @apply h-[10px] w-[10px] scale-100;
  }

  .cyber-radio-input:focus + .cyber-radio {
    box-shadow: 0 0 0 2px rgba(var(--color-brand-500-rgb), 0.2);
  }

  .cyber-radio-wrapper:hover .cyber-radio {
    border-color: var(--color-brand-500);
  }

  .cyber-radio-label {
    @apply select-none;
  }

  .cyber-radio-wrapper.is-disabled {
    @apply cursor-not-allowed opacity-60;
  }

  .cyber-radio-wrapper.is-disabled:hover .cyber-radio {
    border-color: rgba(var(--color-brand-500-rgb), 0.3);
  }

  .cyber-radio-wrapper.is-disabled .cyber-radio-input:checked + .cyber-radio {
    background-color: rgba(var(--color-brand-500-rgb), 0.1);
    box-shadow: none;
  }
</style>
