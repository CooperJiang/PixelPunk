<script setup lang="ts">
  import type { CheckboxEmits, CheckboxProps } from './types'

  defineOptions({
    name: 'CyberCheckbox',
  })

  withDefaults(defineProps<CheckboxProps>(), {
    modelValue: false,
    disabled: false,
    indeterminate: false,
    size: 'medium',
  })

  const emit = defineEmits<CheckboxEmits>()

  const handleChange = (e: Event) => {
    const target = e.target as HTMLInputElement
    const value = target.checked
    emit('update:modelValue', value)
    emit('change', value)
  }
</script>

<template>
  <label
    class="cyber-checkbox"
    :class="{
      'cyber-checkbox--disabled': disabled,
      'cyber-checkbox--checked': modelValue,
      'cyber-checkbox--indeterminate': indeterminate,
      [`cyber-checkbox--${size}`]: size,
    }"
  >
    <span class="cyber-checkbox__input">
      <input type="checkbox" :checked="modelValue" :disabled="disabled" :name="name" @change="handleChange" />
      <span class="cyber-checkbox__inner">
        <i v-if="modelValue && !indeterminate" class="cyber-checkbox__check fas fa-check" />
        <i v-if="indeterminate" class="cyber-checkbox__indeterminate fas fa-minus" />
      </span>
    </span>
    <span v-if="$slots.default" class="cyber-checkbox__label">
      <slot />
    </span>
  </label>
</template>

<style scoped>
  .cyber-checkbox {
    @apply flex cursor-pointer select-none items-center text-sm leading-[18px] text-content outline-none;
  }

  .cyber-checkbox--disabled {
    @apply cursor-not-allowed text-content-disabled;
  }

  .cyber-checkbox__input {
    @apply relative mr-2 inline-flex h-5 w-5 flex-shrink-0 items-center justify-center;
  }

  .cyber-checkbox__input input {
    @apply absolute left-0 top-0 z-[2] m-0 h-full w-full cursor-pointer p-0 opacity-0;
    -webkit-appearance: none;
    appearance: none;
    outline: none;
    border: none;
    background: none;
    box-shadow: none;
  }

  .cyber-checkbox--disabled .cyber-checkbox__input input {
    @apply cursor-not-allowed;
  }

  .cyber-checkbox__inner {
    @apply relative z-[1] box-border flex h-5 w-5 items-center justify-center border-2;
    border-radius: var(--radius-sm);
    border-color: rgba(var(--color-brand-500-rgb), 0.4);
    background: linear-gradient(
      135deg,
      rgba(var(--color-background-900-rgb), 0.9) 0%,
      rgba(var(--color-background-900-rgb), 0.7) 100%
    );
    transition: all 0.3s cubic-bezier(0.25, 0.8, 0.25, 1);
    box-shadow:
      inset 0 1px 0 rgba(var(--color-white-rgb), 0.1),
      0 1px 3px rgba(var(--color-background-900-rgb), 0.3);
  }

  .cyber-checkbox--checked .cyber-checkbox__inner {
    background: linear-gradient(135deg, rgba(var(--color-brand-500-rgb), 0.3) 0%, rgba(var(--color-brand-500-rgb), 0.2) 100%);
    border-color: var(--color-brand-500);
    box-shadow:
      inset 0 1px 0 rgba(var(--color-white-rgb), 0.2),
      0 0 8px rgba(var(--color-brand-500-rgb), 0.3),
      0 1px 2px rgba(var(--color-background-900-rgb), 0.2);
  }

  .cyber-checkbox--disabled .cyber-checkbox__inner {
    background: linear-gradient(
      135deg,
      rgba(var(--color-background-900-rgb), 0.5) 0%,
      rgba(var(--color-background-900-rgb), 0.4) 100%
    );
    border-color: rgba(var(--color-brand-500-rgb), 0.2);
    box-shadow:
      inset 0 1px 0 rgba(var(--color-white-rgb), 0.05),
      0 1px 3px rgba(var(--color-background-900-rgb), 0.2);
  }

  .cyber-checkbox__check {
    @apply text-[10px] leading-none;
    color: var(--color-brand-500);
    text-shadow:
      0 0 4px rgba(var(--color-brand-500-rgb), 0.8),
      0 0 2px rgba(var(--color-white-rgb), 0.3);
  }

  .cyber-checkbox__indeterminate {
    @apply text-[8px] leading-none;
    color: var(--color-brand-500);
    text-shadow:
      0 0 4px rgba(var(--color-brand-500-rgb), 0.8),
      0 0 2px rgba(var(--color-white-rgb), 0.3);
  }

  .cyber-checkbox--disabled .cyber-checkbox__check,
  .cyber-checkbox--disabled .cyber-checkbox__indeterminate {
    color: rgba(var(--color-brand-500-rgb), 0.4);
    text-shadow: none;
  }

  .cyber-checkbox:hover:not(.cyber-checkbox--disabled) .cyber-checkbox__inner {
    border-color: var(--color-brand-500);
    box-shadow:
      inset 0 1px 0 rgba(var(--color-white-rgb), 0.15),
      0 0 6px rgba(var(--color-brand-500-rgb), 0.2),
      0 1px 2px rgba(var(--color-background-900-rgb), 0.2);
  }

  .cyber-checkbox__label {
    @apply inline-flex items-center text-sm leading-[18px];
  }

  .cyber-checkbox--small {
    @apply text-[11px] leading-[14px];
  }

  .cyber-checkbox--small .cyber-checkbox__input {
    @apply mr-1.5 h-[14px] w-[14px];
  }

  .cyber-checkbox--small .cyber-checkbox__inner {
    @apply h-[14px] w-[14px];
    border-radius: var(--radius-sm);
    border-width: 1.5px;
  }

  .cyber-checkbox--small .cyber-checkbox__check {
    @apply text-[7px];
  }

  .cyber-checkbox--small .cyber-checkbox__indeterminate {
    @apply text-[6px];
  }

  .cyber-checkbox--small .cyber-checkbox__label {
    @apply text-[11px] leading-[14px];
  }

  .cyber-checkbox--large {
    @apply text-base leading-[22px];
  }

  .cyber-checkbox--large .cyber-checkbox__input {
    @apply mr-2.5 h-6 w-6;
  }

  .cyber-checkbox--large .cyber-checkbox__inner {
    @apply h-6 w-6;
    border-radius: var(--radius-sm);
    border-width: 2.5px;
  }

  .cyber-checkbox--large .cyber-checkbox__check {
    @apply text-[12px];
  }

  .cyber-checkbox--large .cyber-checkbox__indeterminate {
    @apply text-[10px];
  }

  .cyber-checkbox--large .cyber-checkbox__label {
    @apply text-base leading-[22px];
  }
</style>
