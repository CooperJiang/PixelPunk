<script setup lang="ts">
  import { computed, onMounted, onUnmounted, ref } from 'vue'
  import { useTexts } from '@/composables/useTexts'
  import type { ImageLoadingProps } from './types'

  defineOptions({
    name: 'CyberFileLoading',
  })

  const { $t } = useTexts()

  const props = withDefaults(defineProps<ImageLoadingProps>(), {
    progress: 0,
    loadingText: undefined,
    showDataStream: true,
  })

  const progressOffset = computed(() => {
    const circumference = 2 * Math.PI * 54
    return circumference - (props.progress / 100) * circumference
  })

  const dataStreams = ref([
    ['L', 'O', 'A', 'D', 'I', 'N', 'G', '.', '.', '.'],
    ['0', '1', '0', '1', '1', '0', '0', '1', '0', '1'],
    ['>', '>', '>', 'I', 'M', 'A', 'G', 'E', '<', '<'],
    ['#', '#', '#', 'D', 'A', 'T', 'A', '#', '#', '#'],
  ])

  const simulatedProgress = ref(0)
  let progressInterval: number | null = null

  onMounted(() => {
    if (props.progress === 0) {
      progressInterval = setInterval(() => {
        if (simulatedProgress.value < 90) {
          simulatedProgress.value += Math.random() * 10
        }
      }, 200)
    }
  })

  onUnmounted(() => {
    if (progressInterval) {
      clearInterval(progressInterval)
    }
  })

  const progress = computed(() => (props.progress > 0 ? props.progress : simulatedProgress.value))

  const displayLoadingText = computed(() => props.loadingText || $t('fileLoading.loading'))
</script>

<template>
  <div class="cyber-image-loading">
    <div class="loading-container">
      <div class="loading-spinner">
        <div class="spinner-outer">
          <div class="spinner-inner" />
          <div class="spinner-dots">
            <div v-for="(dot, index) in 8" :key="index" class="dot" :style="{ '--delay': index * 0.1 + 's' }" />
          </div>
        </div>

        <svg class="progress-ring" viewBox="0 0 120 120">
          <circle class="progress-ring-background" cx="60" cy="60" r="54" />
          <circle class="progress-ring-progress" cx="60" cy="60" r="54" :style="{ strokeDashoffset: progressOffset }" />
        </svg>
      </div>

      <div class="loading-info">
        <div class="loading-text">{{ displayLoadingText }}</div>
        <div class="loading-progress">
          <div class="progress-bar">
            <div class="progress-fill" :style="{ width: progress + '%' }" />
            <div class="progress-glow" />
          </div>
          <div class="progress-percentage">{{ Math.round(progress) }}%</div>
        </div>

        <div class="data-stream">
          <div
            v-for="(stream, index) in dataStreams"
            :key="index"
            class="stream-line"
            :style="{ '--stream-delay': index * 0.3 + 's' }"
          >
            <span v-for="(char, charIndex) in stream" :key="charIndex">{{ char }}</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
  .cyber-image-loading {
    @apply absolute inset-0 flex items-center justify-center;
    background: rgba(var(--color-background-900-rgb), 0.95);
    backdrop-filter: blur(10px);
    z-index: 10;
  }

  .loading-container {
    @apply flex flex-col items-center gap-8;
  }

  .loading-spinner {
    @apply relative;
    width: 120px;
    height: 120px;
  }

  .spinner-outer {
    @apply absolute h-full w-full rounded-full;
    border: 2px solid rgba(var(--color-brand-500-rgb), 0.1);
    animation: rotate 2s linear infinite;
  }

  .spinner-inner {
    @apply absolute left-1/2 top-1/2 h-3/5 w-3/5 rounded-full;
    transform: translate(-50%, -50%);
    border: 2px solid rgba(var(--color-brand-500-rgb), 0.3);
    border-top: 2px solid rgba(var(--color-brand-500-rgb), 0.8);
    animation: rotate 1.5s linear infinite reverse;
  }

  .spinner-dots {
    @apply absolute h-full w-full;
  }

  .dot {
    @apply absolute left-1/2 h-1.5 w-1.5 rounded-full;
    background: rgba(var(--color-brand-500-rgb), 0.8);
    top: 2px;
    transform: translateX(-50%);
    transform-origin: 50% 58px;
    animation:
      rotate 2s linear infinite,
      pulse 1s ease-in-out infinite;
    animation-delay: var(--delay);
    box-shadow:
      0 0 6px rgba(var(--color-brand-500-rgb), 0.6),
      0 0 12px rgba(var(--color-brand-500-rgb), 0.3);
  }

  .progress-ring {
    @apply absolute h-full w-full;
    transform: rotate(-90deg);
  }

  .progress-ring-background {
    fill: none;
    stroke: rgba(var(--color-brand-500-rgb), 0.1);
    stroke-width: 2;
  }

  .progress-ring-progress {
    fill: none;
    stroke: rgba(var(--color-brand-500-rgb), 0.9);
    stroke-width: 3;
    stroke-linecap: round;
    stroke-dasharray: 339.292;
    transition: stroke-dashoffset 0.3s ease;
    filter: drop-shadow(0 0 6px rgba(var(--color-brand-500-rgb), 0.6));
  }

  .loading-info {
    @apply flex flex-col items-center gap-4;
    min-width: 300px;
  }

  .loading-text {
    @apply text-lg font-medium uppercase tracking-wider;
    color: rgba(var(--color-brand-500-rgb), 0.9);
    animation: textGlow 2s ease-in-out infinite alternate;
  }

  .loading-progress {
    @apply flex w-full items-center gap-4;
  }

  .progress-bar {
    @apply relative h-1.5 flex-1 overflow-hidden rounded-sm;
    background: rgba(var(--color-brand-500-rgb), 0.1);
  }

  .progress-fill {
    @apply relative h-full rounded-sm;
    background: linear-gradient(
      90deg,
      rgba(var(--color-brand-500-rgb), 0.6) 0%,
      rgba(var(--color-brand-500-rgb), 0.9) 50%,
      rgba(var(--color-brand-500-rgb), 1) 100%
    );
    transition: width 0.3s ease;
  }

  .progress-fill::after {
    @apply absolute right-0 top-0 h-full w-5;
    content: '';
    background: linear-gradient(90deg, transparent, rgba(var(--color-white-rgb), 0.3));
    animation: shimmer 1.5s infinite;
  }

  .progress-glow {
    @apply absolute left-0 right-0 rounded-sm;
    top: -2px;
    height: 10px;
    background: rgba(var(--color-brand-500-rgb), 0.3);
    filter: blur(4px);
    opacity: 0.6;
  }

  .progress-percentage {
    @apply min-w-11 text-right text-sm font-semibold;
    color: rgba(var(--color-brand-500-rgb), 0.9);
    font-family: 'Courier New', monospace;
  }

  .data-stream {
    @apply mt-4 flex flex-col gap-2 opacity-70;
  }

  .stream-line {
    @apply flex gap-1 text-xs;
    font-family: 'Courier New', monospace;
    color: rgba(var(--color-brand-500-rgb), 0.6);
    animation: streamFlow 3s linear infinite;
    animation-delay: var(--stream-delay);
  }

  .stream-line span {
    animation: charGlow 0.5s ease-in-out infinite alternate;
    animation-delay: calc(var(--stream-delay) + var(--char-delay, 0s));
  }

  @media (max-width: 768px) {
    .loading-container {
      @apply gap-6;
    }
    .loading-spinner {
      @apply h-20 w-20;
    }
    .loading-info {
      min-width: 250px;
    }
    .loading-text {
      @apply text-base;
    }
    .data-stream {
      @apply hidden;
    }
  }
</style>
