<script setup lang="ts">
  import { computed } from 'vue'

  defineOptions({
    name: 'ProgressRing',
  })

  interface ProgressRingProps {
    progress: number
    active?: boolean
  }

  const props = withDefaults(defineProps<ProgressRingProps>(), {
    progress: 0,
    active: false,
  })

  const progressStyle = computed(() => {
    const circumference = 2 * Math.PI * 36
    const offset = circumference - (props.progress / 100) * circumference
    return {
      strokeDasharray: `${circumference}`,
      strokeDashoffset: `${offset}`,
      transition: 'stroke-dashoffset 0.3s ease',
    }
  })
</script>

<template>
  <div class="progress-ring-wrapper">
    <svg viewBox="0 0 86 86" class="progress-svg">
      <circle cx="43" cy="43" r="36" stroke="rgba(100, 116, 139, 0.4)" stroke-width="3" fill="none" class="ring-bg" />

      <circle
        cx="43"
        cy="43"
        r="36"
        stroke="url(#uploadGradient)"
        stroke-width="3"
        fill="none"
        stroke-linecap="round"
        class="ring-progress"
        :style="progressStyle"
      />

      <defs>
        <linearGradient id="uploadGradient" x1="0%" y1="0%" x2="100%" y2="100%">
          <stop offset="0%" style="stop-color: var(--color-brand-400); stop-opacity: 1" />
          <stop offset="50%" style="stop-color: var(--color-brand-500); stop-opacity: 1" />
          <stop offset="100%" style="stop-color: var(--color-brand-600); stop-opacity: 1" />
        </linearGradient>
        <filter id="glow">
          <feGaussianBlur stdDeviation="2.5" result="coloredBlur" />
          <feMerge>
            <feMergeNode in="coloredBlur" />
            <feMergeNode in="SourceGraphic" />
          </feMerge>
        </filter>
      </defs>
    </svg>

    <div class="pulse-ring" :class="{ active }" />
  </div>
</template>

<style scoped>
  .progress-ring-wrapper {
    @apply relative h-full w-full;
  }

  .progress-svg {
    @apply absolute left-1/2 top-1/2;
    width: 86px;
    height: 86px;
    transform: translate(-50%, -50%) rotate(-90deg);
    filter: drop-shadow(0 0 8px var(--color-brand-500));
  }

  .ring-bg {
    opacity: 0.5;
  }

  .ring-progress {
    filter: url(#glow);
    animation: pulse-glow 2s ease-in-out infinite alternate;
    stroke-width: 4;
  }

  .pulse-ring {
    @apply pointer-events-none absolute left-1/2 top-1/2 h-[100%] w-[100%] rounded-full border-2;
    border-color: transparent;
    transform: translate(-50%, -50%);
  }

  .pulse-ring.active {
    animation: pulse-expand 2s ease-out infinite;
    border-color: rgba(var(--color-brand-500-rgb), 0.8);
    box-shadow: 0 0 20px rgba(var(--color-brand-500-rgb), 0.6);
  }

  @keyframes pulse-glow {
    0% {
      filter: drop-shadow(0 0 5px var(--color-brand-500));
    }
    100% {
      filter: drop-shadow(0 0 15px var(--color-brand-400));
    }
  }

  @keyframes pulse-expand {
    0% {
      transform: translate(-50%, -50%) scale(1);
      opacity: 0.6;
    }
    50% {
      opacity: 0.2;
    }
    100% {
      transform: translate(-50%, -50%) scale(1.3);
      opacity: 0;
    }
  }
</style>
