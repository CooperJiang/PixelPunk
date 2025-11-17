<template>
  <button
    :class="buttonClasses"
    :disabled="disabled || loading"
    @click="handleClick"
    @mouseover="onMouseover"
    @mouseleave="onMouseleave"
  >
    <div v-if="loading && loadingMode === 'replace'" class="spinner"></div>

    <template v-else>
      <div v-if="loading && loadingMode === 'inline'" class="spinner inline-spinner"></div>

      <span v-if="!loading && (icon || $slots.icon)" class="icon-wrap left-icon">
        <slot name="icon">
          <i v-if="icon" :class="['fas', `fa-${icon}`]"></i>
        </slot>
      </span>

      <span v-if="loadingMode === 'replace' || !loading" class="button-content">
        <slot></slot>
      </span>

      <span v-if="loadingMode === 'inline' && loading" class="button-content">
        <slot></slot>
      </span>
    </template>

    <span v-if="rightIcon || $slots.rightIcon" class="icon-wrap right-icon">
      <slot name="rightIcon">
        <i v-if="rightIcon" :class="['fas', `fa-${rightIcon}`]"></i>
      </slot>
    </span>
  </button>
</template>

<script setup lang="ts">
  import { computed } from 'vue'
  import type { ButtonProps, ButtonEmits } from './types'

  defineOptions({
    name: 'CyberButton',
  })

  const props = withDefaults(defineProps<ButtonProps>(), {
    type: 'primary',
    size: 'medium',
    icon: '',
    rightIcon: '',
    loading: false,
    loadingMode: 'inline',
    disabled: false,
    block: false,
    customClass: '',
  })

  const emit = defineEmits<ButtonEmits>()

  const buttonClasses = computed(() => [
    'cyber-btn',
    `cyber-btn--${props.type}`,
    `cyber-btn--${props.size}`,
    {
      'cyber-btn--loading': props.loading,
      'cyber-btn--disabled': props.disabled,
      'cyber-btn--block': props.block,
      'cyber-btn--has-icon': props.icon || props.rightIcon,
      'cyber-btn--replace-mode': props.loading && props.loadingMode === 'replace',
    },
    props.customClass,
  ])

  const handleClick = (event: MouseEvent) => {
    if (props.loading || props.disabled) return
    emit('click', event)
  }

  const onMouseover = (event: MouseEvent) => {
    emit('mouseover', event)
  }

  const onMouseleave = (event: MouseEvent) => {
    emit('mouseleave', event)
  }
</script>

<style scoped>
  .cyber-btn {
    @apply relative inline-flex cursor-pointer select-none items-center justify-center font-normal outline-none;
    border-radius: var(--radius-sm);
    border: 2px solid transparent;
    transition: all 0.15s ease;
    box-sizing: border-box;
    font-weight: 500;
    letter-spacing: 0.025em;
    white-space: nowrap;
  }

  .cyber-btn--small {
    @apply gap-1 px-2.5 py-1;
    min-height: 24px;
    font-size: 12px;
  }

  .cyber-btn--sm {
    @apply gap-1 px-2.5 py-1;
    min-height: 24px;
    font-size: 12px;
  }

  .cyber-btn--medium {
    @apply gap-1.5 px-3;
    min-height: 32px;
    font-size: 13px;
  }

  .cyber-btn--large {
    @apply gap-2 px-4 py-2;
    min-height: 36px;
    font-size: 14px;
  }

  .cyber-btn--primary {
    background: var(--color-brand-500);
    color: var(--color-text-on-brand);
    border-color: var(--color-brand-500);
    position: relative;
    overflow: hidden;
    box-shadow: 0.15em 0.15em 0 rgba(var(--color-brand-500-rgb), 0.6);
  }

  .cyber-btn--primary:hover:not(.cyber-btn--disabled):not(.cyber-btn--loading) {
    background: rgba(var(--color-brand-500-rgb), 0.9);
    border-color: rgba(var(--color-brand-500-rgb), 0.9);
    transform: translate(-0.05em, -0.05em);
    box-shadow:
      0.2em 0.2em 0 rgba(var(--color-brand-500-rgb), 0.7),
      0 0 20px rgba(var(--color-brand-500-rgb), 0.3);
  }

  .cyber-btn--primary:active:not(.cyber-btn--disabled):not(.cyber-btn--loading) {
    transform: translate(0.05em, 0.05em);
    box-shadow:
      0.1em 0.1em 0 rgba(var(--color-brand-500-rgb), 0.5),
      0 0 10px rgba(var(--color-brand-500-rgb), 0.2);
  }

  .cyber-btn--secondary {
    background: var(--color-background-700);
    color: var(--color-content-default);
    border-color: var(--color-border-default);
    position: relative;
    overflow: hidden;
    box-shadow: 0.15em 0.15em 0 rgba(var(--color-brand-500-rgb), 0.3);
  }

  .cyber-btn--secondary:hover:not(.cyber-btn--disabled):not(.cyber-btn--loading) {
    background: var(--color-hover-bg);
    border-color: var(--color-border-strong);
    color: var(--color-content-heading);
    transform: translate(-0.05em, -0.05em);
    box-shadow:
      0.2em 0.2em 0 rgba(var(--color-brand-500-rgb), 0.4),
      0 0 15px rgba(var(--color-brand-500-rgb), 0.15);
  }

  .cyber-btn--secondary:active:not(.cyber-btn--disabled):not(.cyber-btn--loading) {
    background: var(--color-background-800);
    transform: translate(0.05em, 0.05em);
    box-shadow:
      0.1em 0.1em 0 rgba(var(--color-brand-500-rgb), 0.25),
      0 0 8px rgba(var(--color-brand-500-rgb), 0.1);
  }

  .cyber-btn--outlined {
    background: transparent;
    color: var(--color-content-heading);
    border-color: rgba(var(--color-brand-500-rgb), 0.4);
    position: relative;
    overflow: hidden;
    box-shadow: 0.15em 0.15em 0 rgba(var(--color-brand-500-rgb), 0.3);
  }

  .cyber-btn--outlined:hover:not(.cyber-btn--disabled):not(.cyber-btn--loading) {
    border-color: var(--color-brand-500);
    color: var(--color-content-heading);
    transform: translate(-0.05em, -0.05em);
    box-shadow:
      0.2em 0.2em 0 rgba(var(--color-brand-500-rgb), 0.4),
      0 0 15px rgba(var(--color-brand-500-rgb), 0.2);
  }

  .cyber-btn--outlined:active:not(.cyber-btn--disabled):not(.cyber-btn--loading) {
    transform: translate(0.05em, 0.05em);
    box-shadow:
      0.1em 0.1em 0 rgba(var(--color-brand-500-rgb), 0.2),
      0 0 8px rgba(var(--color-brand-500-rgb), 0.15);
  }

  .cyber-btn--text {
    background: transparent;
    color: var(--color-brand-500);
    border-color: transparent;
    position: relative;
    overflow: hidden;
  }

  .cyber-btn--text:hover:not(.cyber-btn--disabled):not(.cyber-btn--loading) {
    color: var(--color-brand-500);
    text-shadow: 0 0 4px rgba(var(--color-brand-500-rgb), 0.3);
  }

  .cyber-btn--danger {
    background: var(--color-error-500);
    color: var(--color-text-on-error);
    border-color: var(--color-error-500);
    position: relative;
    overflow: hidden;
    box-shadow: 0.15em 0.15em 0 rgba(var(--color-error-rgb), 0.6);
  }

  .cyber-btn--danger:hover:not(.cyber-btn--disabled):not(.cyber-btn--loading) {
    background: var(--color-error-600);
    border-color: var(--color-error-600);
    transform: translate(-0.05em, -0.05em);
    box-shadow:
      0.2em 0.2em 0 rgba(var(--color-error-rgb), 0.7),
      0 0 20px rgba(var(--color-error-rgb), 0.4);
  }

  .cyber-btn--danger:active:not(.cyber-btn--disabled):not(.cyber-btn--loading) {
    transform: translate(0.05em, 0.05em);
    box-shadow:
      0.1em 0.1em 0 rgba(var(--color-error-rgb), 0.5),
      0 0 10px rgba(var(--color-error-rgb), 0.3);
  }

  .cyber-btn--success {
    background: var(--color-success-500);
    color: var(--color-text-on-success);
    border-color: var(--color-success-500);
    position: relative;
    overflow: hidden;
    box-shadow: 0.15em 0.15em 0 rgba(var(--color-success-rgb), 0.6);
  }

  .cyber-btn--success:hover:not(.cyber-btn--disabled):not(.cyber-btn--loading) {
    background: var(--color-success-600);
    border-color: var(--color-success-600);
    transform: translate(-0.05em, -0.05em);
    box-shadow:
      0.2em 0.2em 0 rgba(var(--color-success-rgb), 0.7),
      0 0 20px rgba(var(--color-success-rgb), 0.3);
  }

  .cyber-btn--success:active:not(.cyber-btn--disabled):not(.cyber-btn--loading) {
    transform: translate(0.05em, 0.05em);
    box-shadow:
      0.1em 0.1em 0 rgba(var(--color-success-rgb), 0.5),
      0 0 10px rgba(var(--color-success-rgb), 0.2);
  }

  .cyber-btn--warning {
    background: var(--color-warning-500);
    color: var(--color-text-on-warning);
    border-color: var(--color-warning-500);
    position: relative;
    overflow: hidden;
    box-shadow: 0.15em 0.15em 0 rgba(var(--color-warning-rgb), 0.6);
  }

  .cyber-btn--warning:hover:not(.cyber-btn--disabled):not(.cyber-btn--loading) {
    background: var(--color-warning-600);
    border-color: var(--color-warning-600);
    transform: translate(-0.05em, -0.05em);
    box-shadow:
      0.2em 0.2em 0 rgba(var(--color-warning-rgb), 0.7),
      0 0 20px rgba(var(--color-warning-rgb), 0.3);
  }

  .cyber-btn--warning:active:not(.cyber-btn--disabled):not(.cyber-btn--loading) {
    transform: translate(0.05em, 0.05em);
    box-shadow:
      0.1em 0.1em 0 rgba(var(--color-warning-rgb), 0.5),
      0 0 10px rgba(var(--color-warning-rgb), 0.2);
  }

  .cyber-btn--info {
    background: rgba(var(--color-brand-500-rgb), 0.15);
    color: var(--color-brand-500);
    border-color: rgba(var(--color-brand-500-rgb), 0.3);
    position: relative;
    overflow: hidden;
    box-shadow: none;
  }

  .cyber-btn--info:hover:not(.cyber-btn--disabled):not(.cyber-btn--loading) {
    background: rgba(var(--color-brand-500-rgb), 0.25);
    border-color: var(--color-brand-500);
    color: var(--color-brand-500);
  }

  .cyber-btn--info:active:not(.cyber-btn--disabled):not(.cyber-btn--loading) {
    background: rgba(var(--color-brand-500-rgb), 0.3);
  }

  .cyber-btn--block {
    @apply w-full;
  }

  .cyber-btn--disabled {
    @apply pointer-events-none cursor-not-allowed opacity-50;
    filter: grayscale(0.3);
  }

  .cyber-btn--loading {
    @apply pointer-events-none;
    opacity: 0.8;
  }

  .cyber-btn--replace-mode {
    color: transparent;
  }

  .cyber-btn--replace-mode .button-content,
  .cyber-btn--replace-mode .icon-wrap {
    opacity: 0;
    visibility: hidden;
  }

  .spinner {
    @apply absolute;
    width: 14px;
    height: 14px;
    border: 2px solid transparent;
    border-top: 2px solid currentColor;
    border-radius: var(--radius-full);
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    filter: drop-shadow(0 0 2px currentColor);
    animation: spin 1s linear infinite;
  }

  .inline-spinner {
    @apply relative inline-block;
    width: 12px;
    height: 12px;
    border-width: 1.5px;
    top: auto;
    left: auto;
    transform: none;
    margin-right: 6px;
    flex-shrink: 0;
  }

  .button-content {
    @apply inline-flex items-center;
  }

  .icon-wrap {
    @apply inline-flex items-center;
    filter: drop-shadow(0 0 2px rgba(var(--color-brand-500-rgb), 0.3));
  }

  .left-icon {
    margin-right: 0.25rem;
  }

  .right-icon {
    margin-left: 0.25rem;
  }

  @media (max-width: 768px) {
    .cyber-btn--block {
      @apply flex;
    }

    .cyber-btn--small,
    .cyber-btn--sm {
      @apply gap-0.5 px-2 py-0.5;
      min-height: 22px;
      font-size: 11px;
    }

    .cyber-btn--medium {
      @apply gap-1 px-2.5 py-1;
      min-height: 28px;
      font-size: 12px;
    }

    .cyber-btn--large {
      @apply gap-1 px-3 py-1.5;
      min-height: 32px;
      font-size: 13px;
    }
  }

  @media (prefers-reduced-motion: reduce) {
    .cyber-btn::before,
    .spinner {
      animation: none !important;
      transition: none !important;
    }

    .cyber-btn:hover {
      transform: none !important;
    }
  }
</style>
