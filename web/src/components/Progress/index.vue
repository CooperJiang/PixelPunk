<script setup lang="ts">
  import { computed } from 'vue'

  defineOptions({
    name: 'CyberProgress',
  })

  interface ProgressProps {
    percent?: number
    showPercent?: boolean
    height?: string
    type?: 'primary' | 'success' | 'warning' | 'error'
    animated?: boolean
  }

  const props = withDefaults(defineProps<ProgressProps>(), {
    percent: 0,
    showPercent: true,
    height: '6px',
    type: 'primary',
    animated: true,
  })

  /* 限制进度在 0-100 之间 */
  const normalizedPercent = computed(() => {
    return Math.max(0, Math.min(100, props.percent))
  })

  /* 进度条颜色配置 */
  const colorMap = {
    primary: {
      bg: 'rgba(var(--color-brand-500-rgb), 0.1)',
      gradient: 'linear-gradient(90deg, var(--color-brand-500), var(--color-brand-400))',
    },
    success: {
      bg: 'rgba(var(--color-success-rgb), 0.1)',
      gradient: 'linear-gradient(90deg, var(--color-success-500), var(--color-success-400))',
    },
    warning: {
      bg: 'rgba(var(--color-warning-rgb), 0.1)',
      gradient: 'linear-gradient(90deg, var(--color-warning-500), var(--color-warning-400))',
    },
    error: {
      bg: 'rgba(var(--color-error-rgb), 0.1)',
      gradient: 'linear-gradient(90deg, var(--color-error-500), var(--color-error-400))',
    },
  }

  const colorConfig = computed(() => colorMap[props.type])
</script>

<template>
  <div class="cyber-progress">
    <div class="progress-track" :style="{ height: height, background: colorConfig.bg }">
      <div
        class="progress-bar"
        :class="{ animated: animated }"
        :style="{
          width: `${normalizedPercent}%`,
          background: colorConfig.gradient,
        }"
      />
    </div>
  </div>
</template>

<style scoped lang="scss">
  .cyber-progress {
    width: 100%;
  }

  .progress-track {
    position: relative;
    width: 100%;
    border-radius: var(--radius-full);
    overflow: hidden;
    background: rgba(var(--color-background-700-rgb), 0.8);
  }

  .progress-bar {
    height: 100%;
    border-radius: var(--radius-full);
    transition: width 0.5s cubic-bezier(0.4, 0, 0.2, 1);
    position: relative;

    &.animated::after {
      content: '';
      position: absolute;
      top: 0;
      left: 0;
      right: 0;
      bottom: 0;
      background: linear-gradient(90deg, transparent 0%, rgba(255, 255, 255, 0.2) 50%, transparent 100%);
      animation: shimmer 2s infinite;
    }
  }

  @keyframes shimmer {
    0% {
      transform: translateX(-100%);
    }
    100% {
      transform: translateX(100%);
    }
  }
</style>
