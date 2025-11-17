<script setup lang="ts">
  import { computed } from 'vue'
  import type { BackgroundProps } from './types'

  defineOptions({
    name: 'CyberBackground',
  })

  const props = withDefaults(defineProps<BackgroundProps>(), {
    pattern: 'cyber',
    backgroundColor: 'rgba(var(--color-background-900-rgb), 0.7)',
    className: '',
  })

  const backgroundStyle = computed(() => ({
    backgroundColor: props.backgroundColor,
  }))

  const backgroundClass = computed(() => ['cyber-background', `cyber-bg-${props.pattern}`, props.className])
</script>

<template>
  <div :class="backgroundClass" :style="backgroundStyle">
    <slot />
  </div>
</template>

<style scoped>
  .cyber-background {
    @apply relative h-full w-full overflow-hidden;
  }

  .cyber-bg-grid {
    background-image:
      linear-gradient(rgba(var(--color-brand-500-rgb), 0.1) 1px, transparent 1px),
      linear-gradient(90deg, rgba(var(--color-brand-500-rgb), 0.1) 1px, transparent 1px);
    background-size: 20px 20px;
  }

  .cyber-bg-dots {
    background-image: radial-gradient(rgba(var(--color-brand-500-rgb), 0.15) 1px, transparent 1px);
    background-size: 10px 10px;
  }

  .cyber-bg-diagonal {
    background-image: repeating-linear-gradient(
      45deg,
      rgba(var(--color-brand-500-rgb), 0.05),
      rgba(var(--color-brand-500-rgb), 0.05) 5px,
      rgba(var(--color-error-rgb), 0.05) 5px,
      rgba(var(--color-error-rgb), 0.05) 10px
    );
  }

  .cyber-bg-cyber {
    background:
      linear-gradient(
        to bottom,
        rgba(var(--color-error-rgb), 0.05) 0%,
        transparent 50%,
        rgba(var(--color-brand-500-rgb), 0.05) 100%
      ),
      repeating-linear-gradient(
        to right,
        transparent 0,
        transparent 10px,
        rgba(var(--color-error-rgb), 0.1) 10px,
        rgba(var(--color-error-rgb), 0.1) 12px
      ),
      repeating-linear-gradient(
        to bottom,
        transparent 0,
        transparent 10px,
        rgba(var(--color-brand-500-rgb), 0.1) 10px,
        rgba(var(--color-brand-500-rgb), 0.1) 12px
      );
    background-size:
      100% 100%,
      20px 100%,
      100% 20px;
    background-position: center;
    @apply relative;
  }

  .cyber-bg-cyber::after {
    content: '';
    @apply absolute bottom-0 left-0 top-0 w-px;
    background: linear-gradient(to bottom, transparent, rgba(var(--color-error-rgb), 0.5), transparent);
    box-shadow: 0 0 10px rgba(var(--color-error-rgb), 0.5);
    animation: cyber-scan-vert 4s linear infinite;
  }

  .cyber-bg-circuit {
    background-image:
      radial-gradient(circle at 25% 25%, rgba(var(--color-brand-500-rgb), 0.2) 1px, transparent 1px),
      radial-gradient(circle at 75% 75%, rgba(var(--color-error-rgb), 0.2) 1px, transparent 1px),
      linear-gradient(to right, rgba(var(--color-brand-500-rgb), 0.1) 1px, transparent 1px),
      linear-gradient(to bottom, rgba(var(--color-brand-500-rgb), 0.1) 1px, transparent 1px);
    background-size:
      20px 20px,
      20px 20px,
      10px 10px,
      10px 10px;
    @apply relative overflow-hidden;
  }

  .cyber-bg-glitch {
    @apply relative overflow-hidden;
  }

  .cyber-bg-glitch::after {
    content: '';
    @apply pointer-events-none absolute inset-0 opacity-50;
    background-image: linear-gradient(
      90deg,
      rgba(var(--color-error-rgb), 0) 0%,
      rgba(var(--color-error-rgb), 0.2) 50%,
      rgba(var(--color-error-rgb), 0) 100%
    );
    animation: glitch-scan 2s ease-in-out infinite alternate;
  }

  .cyber-bg-none {
    background: none;
  }
</style>
