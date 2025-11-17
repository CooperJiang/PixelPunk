<script setup lang="ts">
  import type { DropdownOption } from '../types'

  interface DropdownOptionProps {
    option: DropdownOption
    isSelected?: boolean
    isHighlighted?: boolean
    isColorMode?: boolean
    multiple?: boolean
  }

  interface DropdownOptionEmits {
    (e: 'select', option: DropdownOption): void
    (e: 'mouseover'): void
  }

  defineProps<DropdownOptionProps>()
  defineEmits<DropdownOptionEmits>()
</script>

<template>
  <div
    class="cyber-dropdown-option"
    :class="{
      'is-selected': isSelected,
      'is-highlighted': isHighlighted,
      'is-color-option': isColorMode && option.color,
      'is-create-option': option.isCreateOption,
    }"
    @click.stop="$emit('select', option)"
    @mouseover="$emit('mouseover')"
  >
    <template v-if="multiple">
      <div class="checkbox-wrapper" @click.stop="$emit('select', option)">
        <span class="checkbox">
          <i v-if="isSelected" class="fas fa-check" />
        </span>
      </div>
    </template>
    <span v-if="$slots['option-icon'] || (isColorMode && option.color) || option.isCreateOption" class="option-icon">
      <slot name="option-icon" :option="option" :is-selected="isSelected" :is-permanent="option.value === 'permanent'">
        <span v-if="isColorMode && option.color" class="color-block" :style="{ backgroundColor: option.color }" />
        <i v-else-if="option.isCreateOption" class="fas fa-plus-circle create-icon" />
      </slot>
    </span>
    <span class="option-label">{{ option.label }}</span>
  </div>
</template>

<style scoped>
  .cyber-dropdown-option {
    @apply flex cursor-pointer items-center px-2.5 py-1 text-sm transition-colors duration-200;
  }

  .cyber-dropdown-option:hover,
  .cyber-dropdown-option.is-highlighted {
    @apply bg-[rgba(5,217,232,0.1)];
  }

  .cyber-dropdown-option.is-selected {
    @apply bg-[rgba(5,217,232,0.15)];
  }

  .cyber-dropdown-option.is-selected:hover,
  .cyber-dropdown-option.is-selected.is-highlighted {
    @apply bg-[rgba(5,217,232,0.2)];
  }

  .checkbox-wrapper {
    @apply mr-2.5 flex cursor-pointer items-center justify-center;
  }

  .checkbox {
    @apply flex h-4 w-4 items-center justify-center rounded border border-[var(--color-brand-500)] bg-[rgba(5,217,232,0.05)] transition-all duration-200;
  }

  .is-selected .checkbox {
    @apply border-[var(--color-error-500)] bg-[rgba(5,217,232,0.2)];
  }

  .checkbox i {
    @apply text-xs text-[var(--color-error-500)];
  }

  .option-label {
    @apply flex-1 overflow-hidden text-ellipsis whitespace-nowrap;
  }

  .is-color-option {
    @apply flex items-center;
  }

  .is-create-option {
    @apply mt-1 border-t border-dashed border-[rgba(5,217,232,0.3)] pt-2.5;
  }

  .create-icon {
    @apply mr-2 text-sm text-[var(--color-error-500)];
  }

  .option-icon {
    @apply mr-2 inline-flex min-w-4 shrink-0 items-center justify-center;
  }

  .color-block {
    @apply inline-block h-4 w-4 overflow-hidden rounded border border-default;
    position: relative;
  }

  .color-block::after {
    @apply absolute inset-0;
    content: '';
    background: linear-gradient(135deg, transparent 0%, rgba(0, 0, 0, 0.2) 100%);
  }
</style>
