<script setup lang="ts">
  import { ref } from 'vue'
  import type { IconButtonProps } from './types'

  defineOptions({
    name: 'CyberIconButton',
  })

  const props = withDefaults(defineProps<IconButtonProps>(), {
    type: 'default',
    size: 'medium',
    loading: false,
    disabled: false,
    tooltipPlacement: 'top',
  })

  const emit = defineEmits<{
    click: [event: MouseEvent]
  }>()

  const buttonRef = ref<HTMLButtonElement>()

  const handleClick = (event: MouseEvent) => {
    if (props.disabled || props.loading) {
      return
    }
    emit('click', event)
  }

  defineExpose({
    $el: buttonRef,
  })
</script>

<template>
  <CyberTooltip v-if="tooltip" :content="tooltip" :placement="tooltipPlacement">
    <button
      ref="buttonRef"
      class="icon-button"
      :class="[
        `icon-button--${type}`,
        `icon-button--${size}`,
        {
          'icon-button--loading': loading,
          'icon-button--disabled': disabled,
        },
      ]"
      :disabled="disabled || loading"
      @click="handleClick"
    >
      <i v-if="loading" class="fas fa-spinner fa-spin" />
      <slot v-else />
    </button>
  </CyberTooltip>
  <button
    v-else
    ref="buttonRef"
    class="icon-button"
    :class="[
      `icon-button--${type}`,
      `icon-button--${size}`,
      {
        'icon-button--loading': loading,
        'icon-button--disabled': disabled,
      },
    ]"
    :disabled="disabled || loading"
    @click="handleClick"
  >
    <i v-if="loading" class="fas fa-spinner fa-spin" />
    <slot v-else />
  </button>
</template>

<style scoped lang="scss">
  .icon-button {
    @apply relative inline-flex cursor-pointer select-none items-center justify-center border-0 align-middle font-medium outline-none;
    transition: all var(--transition-normal) var(--ease-in-out);
    backdrop-filter: var(--backdrop-blur-md) saturate(120%);
    border-radius: var(--radius-sm);
    position: relative;
    overflow: hidden;

    &::before {
      content: '';
      position: absolute;
      inset: 0;
      border-radius: inherit;
      background: linear-gradient(135deg, rgba(255, 255, 255, 0.08), rgba(255, 255, 255, 0.03));
      opacity: 0;
      transition: opacity var(--transition-normal) ease;
    }

    &:focus-visible {
      @apply outline-none;
      box-shadow:
        0 0 0 2px var(--color-focus-ring),
        0 0 20px rgba(var(--color-brand-500-rgb), 0.1);
    }

    &:active {
      transform: translateY(1px) scale(0.98);
    }

    &--tiny {
      width: 24px;
      height: 24px;
      font-size: var(--text-xs);
      border-radius: var(--radius-sm);

      i {
        font-size: 10px;
      }
    }

    &--small {
      width: var(--size-button-sm);
      height: var(--size-button-sm);
      font-size: var(--text-sm);
      border-radius: var(--radius-sm);

      i {
        font-size: var(--size-icon-sm);
      }
    }

    &--medium {
      width: var(--size-button-md);
      height: var(--size-button-md);
      font-size: var(--text-base);
      border-radius: var(--radius-sm);

      i {
        font-size: var(--size-icon-sm);
      }
    }

    &--large {
      width: var(--size-button-lg);
      height: var(--size-button-lg);
      font-size: var(--text-lg);
      border-radius: var(--radius-lg);

      i {
        font-size: var(--size-icon-md);
      }
    }

    &--default {
      background: linear-gradient(135deg, rgba(var(--color-background-700-rgb), 0.6), rgba(var(--color-background-800-rgb), 0.8));
      border: 1px solid var(--color-border-subtle);
      color: var(--color-content-muted);
      box-shadow:
        var(--shadow-sm),
        inset 0 1px 0 rgba(255, 255, 255, 0.1);

      &:hover:not(.icon-button--disabled):not(.icon-button--loading) {
        background: linear-gradient(
          135deg,
          rgba(var(--color-background-600-rgb), 0.7),
          rgba(var(--color-background-700-rgb), 0.9)
        );
        border-color: var(--color-border-default);
        color: var(--color-content-heading);
        transform: translateY(-2px);
        box-shadow:
          var(--shadow-md),
          inset 0 1px 0 rgba(255, 255, 255, 0.15);

        &::before {
          opacity: 1;
        }
      }
    }

    &--primary {
      background: var(--color-hover-bg);
      border: 1px solid rgba(var(--color-brand-500-rgb), 0.3);
      color: var(--color-brand-500);
      box-shadow:
        0 4px 12px rgba(var(--color-brand-500-rgb), 0.15),
        inset 0 1px 0 rgba(255, 255, 255, 0.1);

      &:hover:not(.icon-button--disabled):not(.icon-button--loading) {
        background: var(--color-active-bg);
        border-color: var(--color-hover-border);
        transform: translateY(-2px);
        box-shadow:
          0 8px 25px rgba(var(--color-brand-500-rgb), 0.25),
          0 0 20px rgba(var(--color-brand-500-rgb), 0.15),
          inset 0 1px 0 rgba(255, 255, 255, 0.15);

        &::before {
          opacity: 1;
        }
      }
    }

    &--secondary {
      background: linear-gradient(135deg, rgba(var(--color-background-600-rgb), 0.4), rgba(var(--color-background-700-rgb), 0.6));
      border: 1px solid var(--color-border-subtle);
      color: var(--color-content-muted);
      box-shadow:
        var(--shadow-xs),
        inset 0 1px 0 rgba(255, 255, 255, 0.08);

      &:hover:not(.icon-button--disabled):not(.icon-button--loading) {
        background: linear-gradient(
          135deg,
          rgba(var(--color-background-500-rgb), 0.5),
          rgba(var(--color-background-600-rgb), 0.7)
        );
        border-color: var(--color-border-default);
        color: var(--color-content-heading);
        transform: translateY(-1px);
        box-shadow:
          var(--shadow-sm),
          inset 0 1px 0 rgba(255, 255, 255, 0.12);
      }
    }

    &--cyber {
      background: linear-gradient(
        135deg,
        rgba(var(--color-brand-500-rgb), 0.08),
        rgba(var(--color-background-800-rgb), 0.9),
        rgba(var(--color-brand-500-rgb), 0.05)
      );
      border: 1px solid rgba(var(--color-brand-500-rgb), 0.2);
      color: var(--color-content-muted);
      position: relative;
      box-shadow:
        var(--shadow-sm),
        0 0 15px rgba(var(--color-brand-500-rgb), 0.05),
        inset 0 1px 0 rgba(255, 255, 255, 0.1);

      &::after {
        content: '';
        position: absolute;
        inset: -1px;
        border-radius: inherit;
        background: linear-gradient(45deg, transparent, rgba(var(--color-brand-500-rgb), 0.1), transparent);
        opacity: 0;
        transition: opacity var(--transition-normal) ease;
        z-index: -1;
      }

      &:hover:not(.icon-button--disabled):not(.icon-button--loading) {
        background: linear-gradient(
          135deg,
          rgba(var(--color-brand-500-rgb), 0.15),
          rgba(var(--color-background-700-rgb), 0.95),
          rgba(var(--color-brand-500-rgb), 0.08)
        );
        border-color: rgba(var(--color-brand-500-rgb), 0.4);
        color: var(--color-brand-500);
        transform: translateY(-2px);
        box-shadow:
          var(--shadow-lg),
          0 0 25px rgba(var(--color-brand-500-rgb), 0.15),
          inset 0 1px 0 rgba(255, 255, 255, 0.15);

        &::before {
          opacity: 1;
        }

        &::after {
          opacity: 1;
        }
      }

      &:focus-visible {
        box-shadow:
          0 0 0 2px var(--color-focus-ring),
          0 0 25px rgba(var(--color-brand-500-rgb), 0.2);
      }
    }

    &--success {
      background: linear-gradient(135deg, rgba(var(--color-success-rgb), 0.15), rgba(var(--color-success-rgb), 0.08));
      border: 1px solid rgba(var(--color-success-rgb), 0.3);
      color: var(--color-success-500);
      box-shadow:
        0 4px 12px rgba(var(--color-success-rgb), 0.15),
        inset 0 1px 0 rgba(255, 255, 255, 0.1);

      &:hover:not(.icon-button--disabled):not(.icon-button--loading) {
        background: linear-gradient(135deg, rgba(var(--color-success-rgb), 0.25), rgba(var(--color-success-rgb), 0.15));
        border-color: rgba(var(--color-success-rgb), 0.5);
        transform: translateY(-2px);
        box-shadow:
          0 8px 25px rgba(var(--color-success-rgb), 0.25),
          0 0 20px rgba(var(--color-success-rgb), 0.15),
          inset 0 1px 0 rgba(255, 255, 255, 0.15);
      }
    }

    &--warning {
      background: linear-gradient(135deg, rgba(var(--color-warning-rgb), 0.15), rgba(var(--color-warning-rgb), 0.08));
      border: 1px solid rgba(var(--color-warning-rgb), 0.3);
      color: var(--color-warning-500);
      box-shadow:
        0 4px 12px rgba(var(--color-warning-rgb), 0.15),
        inset 0 1px 0 rgba(255, 255, 255, 0.1);

      &:hover:not(.icon-button--disabled):not(.icon-button--loading) {
        background: linear-gradient(135deg, rgba(var(--color-warning-rgb), 0.25), rgba(var(--color-warning-rgb), 0.15));
        border-color: rgba(var(--color-warning-rgb), 0.5);
        transform: translateY(-2px);
        box-shadow:
          0 8px 25px rgba(var(--color-warning-rgb), 0.25),
          0 0 20px rgba(var(--color-warning-rgb), 0.15),
          inset 0 1px 0 rgba(255, 255, 255, 0.15);
      }
    }

    &--danger {
      background: linear-gradient(135deg, rgba(var(--color-error-rgb), 0.15), rgba(var(--color-error-rgb), 0.08));
      border: 1px solid rgba(var(--color-error-rgb), 0.3);
      color: var(--color-error-500);
      box-shadow:
        0 4px 12px rgba(var(--color-error-rgb), 0.15),
        inset 0 1px 0 rgba(255, 255, 255, 0.1);

      &:hover:not(.icon-button--disabled):not(.icon-button--loading) {
        background: linear-gradient(135deg, rgba(var(--color-error-rgb), 0.25), rgba(var(--color-error-rgb), 0.15));
        border-color: rgba(var(--color-error-rgb), 0.5);
        transform: translateY(-2px);
        box-shadow:
          0 8px 25px rgba(var(--color-error-rgb), 0.25),
          0 0 20px rgba(var(--color-error-rgb), 0.15),
          inset 0 1px 0 rgba(255, 255, 255, 0.15);
      }
    }

    &--loading,
    &--disabled {
      @apply cursor-not-allowed;
      opacity: var(--opacity-disabled);
      transform: none !important;
      filter: grayscale(0.3);

      &:hover {
        transform: none !important;
        box-shadow: none !important;

        &::before {
          opacity: 0 !important;
        }
      }
    }

    &--loading {
      @apply cursor-wait;

      i {
        animation: spin 1s linear infinite;
      }
    }
  }

  @keyframes spin {
    from {
      transform: rotate(0deg);
    }
    to {
      transform: rotate(360deg);
    }
  }
</style>
