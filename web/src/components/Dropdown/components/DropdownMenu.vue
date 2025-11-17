<script setup lang="ts">
  import DropdownOption from './DropdownOption.vue'
  import type { DropdownOption as DropdownOptionType } from '../types'
  import { useTexts } from '@/composables/useTexts'

  const { $t } = useTexts()

  interface DropdownMenuProps {
    isOpen?: boolean
    searchable?: boolean
    searchQuery?: string
    maxHeight?: string
    filteredOptions?: DropdownOptionType[]
    highlightedIndex?: number
    multiple?: boolean
    isColorMode?: boolean
    dropdownPosition?: {
      top: number
      left: number
      width: number
      zIndex: number
    }
    selectedValues?: (string | number | boolean)[]
  }

  interface DropdownMenuEmits {
    (e: 'update:searchQuery', value: string): void
    (e: 'selectOption', option: DropdownOptionType): void
    (e: 'update:highlightedIndex', index: number): void
    (e: 'clickStop'): void
  }

  const props = withDefaults(defineProps<DropdownMenuProps>(), {
    isOpen: false,
    searchable: false,
    searchQuery: '',
    maxHeight: '250px',
    filteredOptions: () => [],
    highlightedIndex: -1,
    multiple: false,
    isColorMode: false,
    dropdownPosition: () => ({ top: 0, left: 0, width: 0, zIndex: 1000 }),
    selectedValues: () => [],
  })

  defineEmits<DropdownMenuEmits>()

  const isSelected = (value: string | number | boolean) => {
    return props.selectedValues?.some((selectedValue) => String(selectedValue) === String(value))
  }
</script>

<template>
  <Teleport to="body">
    <transition name="cyber-dropdown">
      <div
        v-if="isOpen"
        class="cyber-dropdown-menu custom-scrollbar"
        :style="{
          maxHeight: maxHeight,
          position: 'fixed',
          top: dropdownPosition?.top + 'px',
          left: dropdownPosition?.left + 'px',
          width: dropdownPosition?.width + 'px',
          zIndex: dropdownPosition?.zIndex,
        }"
        @click.stop="$emit('clickStop')"
      >
        <div v-if="searchable" class="cyber-dropdown-search">
          <input
            :value="searchQuery"
            type="text"
            :placeholder="$t('components.dropdown.searchPlaceholder')"
            @input="$emit('update:searchQuery', ($event.target as HTMLInputElement).value)"
            @click.stop
          />
        </div>
        <div class="cyber-dropdown-options">
          <template v-if="filteredOptions.length">
            <template v-if="$slots.option">
              <slot
                v-for="(option, index) in filteredOptions"
                :key="option.value"
                name="option"
                :option="option"
                :is-selected="isSelected(option.value)"
                :is-color-mode="isColorMode"
                :on-select="() => $emit('selectOption', option)"
                :is-highlighted="highlightedIndex === index"
                :is-create-option="option.isCreateOption"
              />
            </template>
            <template v-else>
              <DropdownOption
                v-for="(option, index) in filteredOptions"
                :key="`default-${option.value}`"
                :option="option"
                :is-selected="isSelected(option.value)"
                :is-highlighted="highlightedIndex === index"
                :is-color-mode="isColorMode"
                :multiple="multiple"
                @select="$emit('selectOption', option)"
                @mouseover="$emit('update:highlightedIndex', index)"
              >
                <template #option-icon="{ option: iconOption, isSelected: iconIsSelected, isPermanent }">
                  <slot name="option-icon" :option="iconOption" :is-selected="iconIsSelected" :is-permanent="isPermanent" />
                </template>
              </DropdownOption>
            </template>
          </template>
          <template v-else>
            <div class="cyber-dropdown-empty">{{ $t('components.dropdown.empty') }}</div>
          </template>
        </div>
      </div>
    </transition>
  </Teleport>
</template>

<style scoped>
  .cyber-dropdown-menu {
    @apply mt-0.5 flex flex-col overflow-y-auto rounded border border-[rgba(5,217,232,0.5)] backdrop-blur-sm;
    background-color: rgba(15, 20, 25, 0.98);
    box-shadow:
      0 4px 12px rgba(0, 0, 0, 0.3),
      0 0 0 1px rgba(var(--color-brand-500-rgb), 0.2);
  }

  .cyber-dropdown-search {
    @apply sticky top-0 z-[1] flex-shrink-0 border-b border-[var(--color-background-800)] bg-[rgba(15,20,25,0.95)] p-1;
  }

  .cyber-dropdown-search input {
    @apply min-h-7 w-full rounded border border-[var(--color-brand-500)] bg-[var(--color-background-800)] px-2 py-1 text-sm text-[var(--color-cyber-light)] outline-none;
    font-family: 'Roboto Mono', monospace;
  }

  .cyber-dropdown-search input:focus {
    @apply border-[var(--color-error-500)];
  }

  .cyber-dropdown-options {
    @apply min-h-0 flex-grow overflow-y-visible py-1;
  }

  .cyber-dropdown-empty {
    @apply p-4 text-center text-content-muted;
  }

  .custom-scrollbar::-webkit-scrollbar {
    @apply w-1.5;
  }

  .custom-scrollbar::-webkit-scrollbar-track {
    @apply rounded bg-background-800;
  }

  .custom-scrollbar::-webkit-scrollbar-thumb {
    @apply rounded bg-[rgba(5,217,232,0.4)];
  }

  .custom-scrollbar::-webkit-scrollbar-thumb:hover {
    @apply bg-[rgba(5,217,232,0.6)];
  }

  .cyber-dropdown-enter-active,
  .cyber-dropdown-leave-active {
    @apply transition-all duration-200;
  }

  .cyber-dropdown-enter-from,
  .cyber-dropdown-leave-to {
    @apply -translate-y-2.5 opacity-0;
  }
</style>
