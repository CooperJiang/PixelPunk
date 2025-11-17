<template>
  <div
    class="cyber-accordion-v2"
    :class="[
      `cyber-accordion-v2--${size}`,
      { 'cyber-accordion-v2--bordered': bordered },
      size === 'small' ? 'text-sm' : size === 'large' ? 'text-lg' : 'text-base',
    ]"
    :style="{ gap }"
  >
    <div
      v-for="(item, index) in items"
      :key="item.id"
      class="accordion-item"
      :class="{
        'accordion-item--active': isItemActive(item.id),
        'accordion-item--disabled': item.disabled,
      }"
    >
      <button
        class="accordion-header"
        :disabled="item.disabled"
        @click="handleToggle(item)"
        @keydown.enter="handleToggle(item)"
        @keydown.space.prevent="handleToggle(item)"
      >
        <div class="header-content">
          <div class="header-icon-wrapper">
            <i v-if="item.icon" :class="item.icon" class="header-icon" />
            <div class="icon-glow"></div>
          </div>
          <div class="header-text">
            <h3 class="header-title">
              {{ item.title }}
            </h3>
            <div class="header-line"></div>
          </div>
        </div>
        <div class="toggle-area">
          <div class="toggle-icon" :class="{ 'toggle-icon--rotated': isItemActive(item.id) }">
            <div class="chevron-line chevron-line--1"></div>
            <div class="chevron-line chevron-line--2"></div>
          </div>
          <div class="toggle-glow"></div>
        </div>
      </button>

      <div
        class="accordion-content"
        :class="{ 'accordion-content--expanded': isItemActive(item.id) }"
        :style="getContentStyle(item.id)"
      >
        <div class="content-wrapper">
          <div class="content-inner">
            <slot :item="item" :index="index" :expanded="isItemActive(item.id)">
              <div v-if="item.content" class="default-content" v-html="item.content" />
            </slot>
          </div>
          <div class="content-glow"></div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
  import type { AccordionProps, AccordionEmits, AccordionItem } from './types'

  defineOptions({ name: 'CyberAccordion' })

  const props = withDefaults(defineProps<AccordionProps>(), {
    modelValue: null,
    multiple: false,
    expandedItems: () => [],
    bordered: true,
    gap: '8px',
    animationDuration: 250,
    card: true,
    size: 'medium',
  })

  const emit = defineEmits<AccordionEmits>()

  const isItemActive = (itemId: string | number): boolean => {
    if (props.multiple) {
      return props.expandedItems.includes(itemId)
    }
    return props.modelValue === itemId
  }

  const handleToggle = (item: AccordionItem) => {
    if (item.disabled) return
    const isExpanded = props.multiple ? props.expandedItems.includes(item.id) : props.modelValue === item.id
    if (props.multiple) {
      const newItems = [...props.expandedItems]
      if (isExpanded) {
        const index = newItems.indexOf(item.id)
        if (index > -1) newItems.splice(index, 1)
      } else {
        newItems.push(item.id)
      }
      emit('update:expandedItems', newItems)
    } else {
      emit('update:modelValue', isExpanded ? null : item.id)
    }
    emit('item-click', item, !isExpanded)
    if (isExpanded) {
      emit('item-collapse', item)
    } else {
      emit('item-expand', item)
    }
  }

  const getContentStyle = (itemId: string | number) => {
    const isExpanded = isItemActive(itemId)

    if (isExpanded) {
      return {
        '--expand-duration': `${props.animationDuration}ms`,
      }
    } else {
      return {
        '--expand-duration': '0ms',
      }
    }
  }
</script>

<style scoped>
  .cyber-accordion-v2 {
    @apply flex flex-col;
    font-family:
      'SF Pro Display',
      -apple-system,
      BlinkMacSystemFont,
      sans-serif;
  }

  .cyber-accordion-v2--small {
    @apply text-sm;
  }

  .cyber-accordion-v2--medium {
    @apply text-base;
  }

  .cyber-accordion-v2--large {
    @apply text-lg;
  }

  .accordion-item {
    @apply relative transform-gpu overflow-hidden rounded-xl;

    background: linear-gradient(
      135deg,
      rgba(var(--color-background-900-rgb), 0.95) 0%,
      rgba(var(--color-background-800-rgb), 0.9) 50%,
      rgba(var(--color-background-900-rgb), 0.95) 100%
    );
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.12);
    transition: all 0.3s ease;
  }

  .accordion-item:hover {
    @apply -translate-y-0.5;
    border-color: rgba(var(--color-brand-500-rgb), 0.25);
    box-shadow: 0 4px 20px rgba(var(--color-brand-500-rgb), 0.1);
  }

  .accordion-item--active {
    border-color: rgba(var(--color-brand-500-rgb), 0.4);
    background: linear-gradient(
      135deg,
      rgba(var(--color-background-900-rgb), 0.98) 0%,
      rgba(var(--color-background-800-rgb), 0.95) 50%,
      rgba(var(--color-background-900-rgb), 0.98) 100%
    );
    box-shadow:
      0 12px 35px rgba(var(--color-brand-500-rgb), 0.15),
      inset 0 1px 0 rgba(var(--color-brand-500-rgb), 0.1);
  }

  .accordion-item--disabled {
    @apply cursor-not-allowed opacity-50;
    transform: none !important;
  }

  .accordion-header {
    @apply relative flex w-full cursor-pointer items-center justify-between border-0 bg-transparent px-6 py-5 outline-none;
    font: inherit;
    color: inherit;
    transition: background-color 0.2s ease;
  }

  .accordion-header:focus-visible {
    background: rgba(var(--color-brand-500-rgb), 0.08);
  }

  .accordion-header:disabled {
    @apply cursor-not-allowed;
  }

  .header-content {
    @apply flex min-w-0 flex-1 items-center gap-4;
  }

  .header-icon-wrapper {
    @apply relative flex h-8 w-8 flex-shrink-0 items-center justify-center;
  }

  .header-icon {
    @apply relative z-[2];
    color: var(--color-brand-500);
    font-size: 1.25em;
    transition: all 0.3s ease;
  }

  .icon-glow {
    @apply absolute -inset-1 rounded-full opacity-0;
    background: radial-gradient(circle, rgba(var(--color-brand-500-rgb), 0.15) 0%, transparent 70%);
    transition: opacity 0.3s ease;
  }

  .accordion-item--active .header-icon {
    color: var(--color-brand-500);
    transform: scale(1.1);
    filter: drop-shadow(0 0 8px rgba(var(--color-brand-500-rgb), 0.6));
  }

  .accordion-item--active .icon-glow,
  .accordion-header:hover .icon-glow {
    @apply opacity-100;
  }

  .header-text {
    @apply relative min-w-0 flex-1;
  }

  .header-title {
    @apply m-0 break-words text-base font-semibold;
    color: var(--color-text-content);
    line-height: 1.4;
    transition: all 0.3s ease;
  }

  .accordion-item--active .header-title {
    color: var(--color-brand-500);
    text-shadow: 0 0 10px rgba(var(--color-brand-500-rgb), 0.3);
  }

  .header-line {
    @apply mt-2 h-px origin-left;
    background: linear-gradient(
      90deg,
      rgba(var(--color-brand-500-rgb), 0.3) 0%,
      rgba(var(--color-brand-500-rgb), 0.1) 50%,
      transparent 100%
    );
    transform: scaleX(0);
    transition: transform 0.3s ease;
  }

  .accordion-item--active .header-line,
  .accordion-header:hover .header-line {
    transform: scaleX(1);
  }

  .toggle-area {
    @apply relative flex h-8 w-8 flex-shrink-0 items-center justify-center;
  }

  .toggle-icon {
    @apply relative h-4 w-4;
    transition: transform 0.3s ease;
  }

  .toggle-icon--rotated {
    @apply rotate-180;
  }

  .chevron-line {
    @apply absolute w-2 rounded-sm;
    height: 1.5px;
    background: var(--color-text-secondary);
    transform-origin: center;
    transition: all 0.3s ease;
  }

  .chevron-line--1 {
    @apply left-0.5 top-1.5;
    transform: rotate(45deg);
  }

  .chevron-line--2 {
    @apply right-0.5 top-1.5;
    transform: rotate(-45deg);
  }

  .accordion-item--active .chevron-line {
    background: var(--color-brand-500);
    box-shadow: 0 0 4px rgba(var(--color-brand-500-rgb), 0.5);
  }

  .toggle-glow {
    @apply absolute -inset-1.5 rounded-full opacity-0;
    background: radial-gradient(circle, rgba(var(--color-brand-500-rgb), 0.1) 0%, transparent 70%);
    transition: opacity 0.3s ease;
  }

  .accordion-item--active .toggle-glow,
  .accordion-header:hover .toggle-glow {
    @apply opacity-100;
  }

  .accordion-content {
    @apply max-h-0 overflow-hidden;
  }

  .accordion-content--expanded {
    max-height: 500px;
  }

  .content-wrapper {
    @apply relative;
    border-top: 1px solid rgba(var(--color-brand-500-rgb), 0.08);
  }

  .content-inner {
    @apply relative z-[1] px-6 py-5;
    background: rgba(var(--color-background-900-rgb), 0.15);
  }

  .content-glow {
    @apply absolute left-0 right-0 top-0 h-px;
    background: linear-gradient(90deg, transparent 0%, rgba(var(--color-brand-500-rgb), 0.3) 50%, transparent 100%);
  }

  .default-content {
    @apply m-0 opacity-95;
    color: var(--color-content-muted);
    font-size: 0.9em;
    line-height: 1.7;
  }

  @media (max-width: 768px) {
    .accordion-header {
      @apply px-5 py-4;
    }

    .content-inner {
      @apply px-5 py-4;
    }

    .header-content {
      @apply gap-3;
    }

    .header-icon-wrapper {
      @apply h-7 w-7;
    }

    .header-icon {
      font-size: 1.1em;
    }
  }

  .cyber-accordion-v2:not(.cyber-accordion-v2--bordered) .accordion-item {
    @apply rounded-none border-l-0 border-r-0 border-t-0;
    border-bottom: 1px solid rgba(var(--color-brand-500-rgb), 0.1);
  }

  .cyber-accordion-v2:not(.cyber-accordion-v2--bordered) .accordion-item:last-child {
    @apply border-b-0;
  }

  .cyber-accordion-v2--small .accordion-header {
    @apply px-4 py-3.5;
  }

  .cyber-accordion-v2--small .content-inner {
    @apply px-4 py-3.5;
  }

  .cyber-accordion-v2--large .accordion-header {
    @apply px-7 py-6;
  }

  .cyber-accordion-v2--large .content-inner {
    @apply px-7 py-6;
  }
</style>
