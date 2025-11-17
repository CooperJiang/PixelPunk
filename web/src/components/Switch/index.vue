<script setup lang="ts">
  import type { SwitchEmits, SwitchProps } from './types'

  defineOptions({
    name: 'CyberSwitch',
  })

  const props = withDefaults(defineProps<SwitchProps>(), {
    modelValue: false,
    label: '',
    disabled: false,
    size: 'medium',
  })

  const emit = defineEmits<SwitchEmits>()

  const toggle = () => {
    if (props.disabled) {
      return
    }

    const newValue = !props.modelValue
    emit('update:modelValue', newValue)
    emit('change', newValue)
  }
</script>

<template>
  <label
    class="cyber-switch"
    :class="{
      'cyber-switch--disabled': disabled,
      'cyber-switch--checked': modelValue,
      [`cyber-switch--${size}`]: size,
    }"
  >
    <input type="checkbox" :checked="modelValue" :disabled="disabled" @change="toggle" />
    <div class="cyber-switch__track">
      <div class="cyber-switch__knob">
        <div class="cyber-switch__glow"></div>
        <div class="cyber-switch__flare"></div>
      </div>
    </div>
    <span v-if="label" class="cyber-switch__label">{{ label }}</span>
  </label>
</template>

<style scoped>
  .cyber-switch {
    @apply inline-flex cursor-pointer select-none items-center;
    perspective: 1200px;
    position: relative;
  }

  .cyber-switch input {
    @apply hidden;
  }

  .cyber-switch--disabled {
    @apply pointer-events-none cursor-not-allowed opacity-50;
  }

  .cyber-switch__track {
    @apply relative cursor-pointer;
    width: 48px;
    height: 24px;
    background: rgba(var(--color-background-800-rgb), 0.8);
    border-radius: 12px;
    border: 1.5px solid rgba(var(--color-brand-500-rgb), 0.5);
    box-shadow:
      inset 0 2px 4px rgba(0, 0, 0, 0.3),
      0 1px 2px rgba(0, 0, 0, 0.2);
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  }

  .cyber-switch__knob {
    @apply absolute;
    width: 18px;
    height: 18px;
    top: 2px;
    left: 2px;
    border-radius: var(--radius-full);
    background: linear-gradient(145deg, var(--color-background-600), var(--color-background-700));
    border: 1.5px solid rgba(var(--color-brand-500-rgb), 0.6);
    box-shadow:
      0 2px 4px rgba(0, 0, 0, 0.3),
      inset 0 1px 1px rgba(255, 255, 255, 0.1);
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  }

  .cyber-switch__glow {
    @apply pointer-events-none absolute;
    width: 100%;
    height: 100%;
    background: radial-gradient(circle, rgba(var(--color-brand-500-rgb), 0.4), transparent);
    border-radius: var(--radius-full);
    top: 0;
    left: 0;
    opacity: 0;
    transition: all 0.3s ease;
  }

  .cyber-switch__flare {
    @apply pointer-events-none absolute;
    width: 4px;
    height: 4px;
    background: radial-gradient(circle, rgba(255, 255, 255, 0.8), transparent);
    border-radius: var(--radius-full);
    top: 2px;
    left: 2px;
    opacity: 0;
    transition: all 0.3s ease;
  }

  .cyber-switch--checked .cyber-switch__track {
    background: linear-gradient(135deg, rgba(var(--color-brand-500-rgb), 0.3), rgba(var(--color-brand-600-rgb), 0.5));
    border-color: var(--color-brand-500);
    box-shadow:
      inset 0 2px 4px rgba(var(--color-brand-500-rgb), 0.3),
      0 0 8px rgba(var(--color-brand-500-rgb), 0.3),
      0 1px 2px rgba(0, 0, 0, 0.2);
  }

  .cyber-switch--checked .cyber-switch__knob {
    left: 26px;
    background: linear-gradient(145deg, var(--color-brand-400), var(--color-brand-500));
    border-color: var(--color-brand-400);
    box-shadow:
      0 2px 6px rgba(var(--color-brand-500-rgb), 0.5),
      0 0 12px rgba(var(--color-brand-500-rgb), 0.4),
      inset 0 1px 1px rgba(255, 255, 255, 0.3);
  }

  .cyber-switch--checked .cyber-switch__glow {
    opacity: 0.6;
    animation: pulse-glow 2s infinite ease-in-out;
  }

  .cyber-switch--checked .cyber-switch__flare {
    opacity: 0.8;
  }

  .cyber-switch:not(.cyber-switch--disabled):hover .cyber-switch__track {
    border-color: rgba(var(--color-brand-500-rgb), 0.8);
  }

  .cyber-switch:not(.cyber-switch--disabled):hover .cyber-switch__knob {
    transform: scale(1.05);
    box-shadow:
      0 3px 8px rgba(0, 0, 0, 0.4),
      inset 0 1px 2px rgba(255, 255, 255, 0.15);
  }

  .cyber-switch--checked:not(.cyber-switch--disabled):hover .cyber-switch__knob {
    box-shadow:
      0 3px 10px rgba(var(--color-brand-500-rgb), 0.6),
      0 0 16px rgba(var(--color-brand-500-rgb), 0.5),
      inset 0 1px 2px rgba(255, 255, 255, 0.4);
  }

  .cyber-switch--small .cyber-switch__track {
    width: 48px;
    height: 24px;
    border-radius: var(--radius-sm);
  }

  .cyber-switch--small .cyber-switch__knob {
    width: 18px;
    height: 18px;
    top: 2px;
    left: 2px;
  }

  .cyber-switch--small.cyber-switch--checked .cyber-switch__knob {
    left: 26px;
  }

  .cyber-switch--large .cyber-switch__track {
    width: 60px;
    height: 30px;
    border-radius: var(--radius-sm);
  }

  .cyber-switch--large .cyber-switch__knob {
    width: 24px;
    height: 24px;
    top: 2px;
    left: 2px;
  }

  .cyber-switch--large.cyber-switch--checked .cyber-switch__knob {
    left: 32px;
  }

  .cyber-switch--large .cyber-switch__flare {
    width: 5px;
    height: 5px;
    top: 3px;
    left: 3px;
  }

  .cyber-switch__label {
    @apply ml-3 inline-flex items-center text-sm leading-[18px];
    color: var(--color-content-default);
    font-weight: 500;
    transition: color 0.3s ease;
  }

  .cyber-switch:hover .cyber-switch__label {
    color: var(--color-brand-500);
  }

  @keyframes pulse-glow {
    0%,
    100% {
      opacity: 0.4;
      transform: scale(1);
    }
    50% {
      opacity: 0.7;
      transform: scale(1.05);
    }
  }

  @media (max-width: 768px) {
    .cyber-switch__track {
      width: 50px;
      height: 26px;
    }

    .cyber-switch__knob {
      width: 20px;
      height: 20px;
    }

    .cyber-switch--checked .cyber-switch__knob {
      left: 28px;
    }
  }

  @media (prefers-reduced-motion: reduce) {
    .cyber-switch__track,
    .cyber-switch__knob,
    .cyber-switch__glow,
    .cyber-switch__flare {
      animation: none !important;
      transition: none !important;
    }
  }
</style>
