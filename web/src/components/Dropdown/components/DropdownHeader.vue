<script setup lang="ts">
  import type { DropdownOption } from '../types'
  import { useTexts } from '@/composables/useTexts'

  const { $t } = useTexts()

  interface DropdownHeaderProps {
    isOpen?: boolean
    disabled?: boolean
    multiple?: boolean
    placeholder?: string
    borderColor?: string
    textColor?: string
    bgColor?: string
    isColorMode?: boolean
    clearable?: boolean
    selectedOption?: DropdownOption | null
    selectedOptions?: DropdownOption[]
    showClearButton?: boolean
    hoveredTag?: string | number | boolean | null
  }

  interface DropdownHeaderEmits {
    (e: 'toggle'): void
    (e: 'clearAll'): void
    (e: 'removeSelected', value: string | number | boolean): void
    (e: 'mouseenterControls'): void
    (e: 'mouseleaveControls'): void
    (e: 'mouseenterTag', value: string | number | boolean): void
    (e: 'mouseleaveTag'): void
  }

  defineProps<DropdownHeaderProps>()
  defineEmits<DropdownHeaderEmits>()
</script>

<template>
  <div
    class="cyber-dropdown-header"
    :style="{
      '--border-color': borderColor,
      '--text-color': textColor,
      '--bg-color': bgColor,
    }"
    @click="$emit('toggle')"
  >
    <div class="cyber-dropdown-selected">
      <template v-if="multiple && selectedOptions && selectedOptions.length">
        <div class="selected-tags">
          <div
            v-for="option in selectedOptions"
            :key="option.value"
            class="selected-tag"
            @mouseenter="$emit('mouseenterTag', option.value)"
            @mouseleave="$emit('mouseleaveTag')"
          >
            <span v-if="isColorMode && option.color" class="color-block" :style="{ backgroundColor: option.color }" />
            {{ option.label }}
            <span
              class="remove-tag"
              :class="{ visible: hoveredTag === option.value }"
              @click.stop="$emit('removeSelected', option.value)"
              >×</span
            >
          </div>
        </div>
      </template>
      <template v-else-if="!multiple && selectedOption">
        <div class="selected-option-container">
          <span v-if="$slots['selected-icon'] || (isColorMode && selectedOption.color)" class="selected-icon">
            <slot name="selected-icon" :option="selectedOption" :is-permanent="selectedOption.value === 'permanent'">
              <span
                v-if="isColorMode && selectedOption.color"
                class="color-block"
                :style="{ backgroundColor: selectedOption.color }"
              />
            </slot>
          </span>
          <span class="selected-label">{{ selectedOption.label }}</span>
        </div>
      </template>
      <template v-else>
        <span class="placeholder">{{ placeholder }}</span>
      </template>
    </div>
    <div class="cyber-dropdown-controls" @mouseenter="$emit('mouseenterControls')" @mouseleave="$emit('mouseleaveControls')">
      <span
        v-if="(multiple ? selectedOptions && selectedOptions.length > 0 : selectedOption) && showClearButton && clearable"
        class="clear-all"
        :title="$t('components.dropdown.clearAll')"
        @click.stop="$emit('clearAll')"
      >
        <i class="fas fa-times-circle"></i>
      </span>
      <i v-else class="fas fa-chevron-down cyber-dropdown-arrow" :class="{ 'is-open': isOpen }"></i>
    </div>
  </div>
</template>

<style scoped>
  .cyber-dropdown-header {
    @apply box-border flex h-8 min-h-8 cursor-pointer items-center justify-between rounded px-3 text-sm leading-none transition-all duration-200;
    background-color: var(--bg-color, var(--color-background-800));
    border: 1px solid var(--border-color, var(--color-brand-500));
    color: var(--text-color, var(--color-cyber-light));
  }

  .cyber-dropdown-header:hover {
    @apply border-[var(--color-error-500)] shadow-[0_0_10px_rgba(5,217,232,0.2)];
  }

  .cyber-dropdown-selected {
    @apply mr-2.5 flex h-[calc(32px-2px)] max-h-[calc(32px-2px)] flex-1 items-center overflow-hidden text-ellipsis whitespace-nowrap text-sm leading-none;
  }

  .cyber-dropdown-controls {
    @apply flex h-6 shrink-0 items-center justify-center;
  }

  .cyber-dropdown-arrow {
    @apply inline-flex items-center justify-center text-xs transition-transform duration-200;
    transition-timing-function: cubic-bezier(0.175, 0.885, 0.32, 1.275);
  }

  .cyber-dropdown-arrow.is-open {
    @apply rotate-180;
  }

  .placeholder {
    @apply text-content-muted;
  }

  .selected-tags {
    @apply flex h-full w-full items-center overflow-x-hidden whitespace-nowrap;
  }

  .selected-tag {
    @apply mr-1.5 inline-flex min-w-5 items-center overflow-hidden text-ellipsis whitespace-nowrap rounded border border-[var(--color-brand-500)] bg-[rgba(5,217,232,0.2)] px-1 text-xs text-[var(--color-cyber-light)] transition-all duration-200;
    height: 18px;
    vertical-align: middle;
  }

  .selected-tag:hover {
    @apply border-[var(--color-error-500)] bg-[rgba(5,217,232,0.3)];
  }

  .remove-tag {
    @apply ml-1.5 inline-flex w-2 cursor-pointer items-center justify-center font-bold text-[var(--color-cyber-light)] opacity-0 transition-opacity duration-200;
  }

  .remove-tag.visible {
    @apply opacity-70;
  }

  .remove-tag.visible:hover {
    @apply text-[var(--color-error-500)] opacity-100;
  }

  .clear-all {
    @apply flex h-6 w-6 cursor-pointer items-center justify-center text-content-muted transition-all duration-200;
  }

  .clear-all:hover {
    @apply text-[var(--color-error-500)];
  }

  .clear-all i {
    @apply text-sm;
  }

  .color-block {
    @apply mr-1.5 inline-block h-4 w-4 overflow-hidden rounded border border-default;
    position: relative;
  }

  .color-block::after {
    @apply absolute inset-0;
    content: '';
    background: linear-gradient(135deg, transparent 0%, rgba(0, 0, 0, 0.2) 100%);
  }

  .selected-option-container {
    @apply flex w-full items-center;
  }

  .selected-icon {
    @apply mr-2 inline-flex min-w-4 shrink-0 items-center justify-center;
  }

  .selected-label {
    @apply flex-1 overflow-hidden text-ellipsis whitespace-nowrap;
  }
</style>
