<script setup lang="ts">
  import type { FloatUploadStatus } from '../types'

  defineOptions({
    name: 'CenterContent',
  })

  interface CenterContentProps {
    status: FloatUploadStatus
    progress: number
    queueSize: number
  }

  defineProps<CenterContentProps>()
</script>

<template>
  <div class="ring-center">
    <div v-if="status === 'uploading'" class="center-content">
      <div class="percent-line">
        <span class="percent-display">{{ Math.round(progress) }}</span>
        <span class="percent-symbol">%</span>
      </div>
      <div v-if="queueSize > 0" class="queue-display">{{ queueSize }}</div>
    </div>

    <div v-else-if="status === 'success'" class="center-content success">
      <div class="status-icon">✓</div>
      <div v-if="queueSize > 0" class="queue-display">{{ queueSize }}</div>
    </div>

    <div v-else-if="status === 'error'" class="center-content error">
      <div class="status-icon">✗</div>
      <div v-if="queueSize > 0" class="queue-display">{{ queueSize }}</div>
    </div>

    <div v-else-if="status === 'paused'" class="center-content paused">
      <div class="status-icon">❚❚</div>
      <div v-if="queueSize > 0" class="queue-display">{{ queueSize }}</div>
    </div>

    <div v-else class="center-content preparing">
      <div class="loading-dots">
        <span />
        <span />
        <span />
      </div>
      <div v-if="queueSize > 0" class="queue-display">{{ queueSize }}</div>
    </div>
  </div>
</template>

<style scoped>
  .ring-center {
    @apply absolute left-1/2 top-1/2 z-[2] -translate-x-1/2 -translate-y-1/2 transform text-center;
  }

  .center-content {
    @apply flex flex-col items-center gap-0.5;
  }

  .percent-line {
    @apply flex items-baseline gap-0.5 leading-none;
  }

  .percent-display {
    @apply font-mono text-xl font-black text-brand-300;
    text-shadow:
      0 0 12px var(--color-brand-500),
      0 0 24px var(--color-brand-500),
      0 0 36px var(--color-brand-400);
    filter: drop-shadow(0 2px 4px rgba(0, 0, 0, 0.5));
  }

  .percent-symbol {
    @apply ml-px text-sm font-bold text-error-400;
    text-shadow:
      0 0 8px var(--color-error-500),
      0 0 16px var(--color-error-500);
  }

  .queue-display {
    @apply mt-0.5 rounded px-1 py-px font-mono text-[10px] font-semibold;
    background: rgba(239, 68, 68, 0.25);
    color: rgb(254, 215, 215);
    text-shadow:
      0 0 8px var(--color-error-400),
      0 0 12px var(--color-error-500);
    box-shadow: 0 0 10px rgba(239, 68, 68, 0.4);
    animation: queue-pulse 1.5s ease-in-out infinite;
  }

  .status-icon {
    @apply text-3xl font-black;
    filter: drop-shadow(0 0 8px currentColor);
  }

  .success .status-icon {
    @apply text-success-400;
    text-shadow:
      0 0 20px var(--color-success-500),
      0 0 40px var(--color-success-500),
      0 0 60px var(--color-success-400);
    animation:
      scaleIn 0.3s ease-out,
      pulse-success 2s ease-in-out infinite;
  }

  .error .status-icon {
    @apply text-error-400;
    text-shadow:
      0 0 20px var(--color-error-500),
      0 0 40px var(--color-error-500),
      0 0 60px var(--color-error-400);
    animation:
      error-shake 0.6s ease-out,
      pulse-error 2s ease-in-out infinite;
  }

  .paused .status-icon {
    @apply text-warning-400;
    text-shadow:
      0 0 20px var(--color-warning-500),
      0 0 40px var(--color-warning-500);
  }

  .loading-dots {
    @apply mb-1 flex justify-center gap-1;
  }

  .loading-dots span {
    @apply h-1.5 w-1.5 rounded-full;
    background: var(--color-brand-400);
    animation: loading-bounce 1.4s ease-in-out infinite both;
    box-shadow:
      0 0 8px var(--color-brand-500),
      0 0 16px var(--color-brand-500);
    filter: brightness(1.2);
  }

  .loading-dots span:nth-child(1) {
    animation-delay: -0.32s;
  }

  .loading-dots span:nth-child(2) {
    animation-delay: -0.16s;
  }

  .loading-dots span:nth-child(3) {
    animation-delay: 0;
  }

  @keyframes loading-bounce {
    0%,
    80%,
    100% {
      transform: scale(0.8);
      opacity: 0.5;
    }
    40% {
      transform: scale(1.2);
      opacity: 1;
    }
  }

  @keyframes scaleIn {
    from {
      transform: scale(0);
      opacity: 0;
    }
    to {
      transform: scale(1);
      opacity: 1;
    }
  }

  @keyframes error-shake {
    0%,
    100% {
      transform: translateX(0);
    }
    25% {
      transform: translateX(-5px);
    }
    75% {
      transform: translateX(5px);
    }
  }

  @keyframes pulse-success {
    0%,
    100% {
      text-shadow:
        0 0 20px var(--color-success-500),
        0 0 40px var(--color-success-500),
        0 0 60px var(--color-success-400);
    }
    50% {
      text-shadow:
        0 0 10px var(--color-success-500),
        0 0 20px var(--color-success-500),
        0 0 30px var(--color-success-400);
    }
  }

  @keyframes pulse-error {
    0%,
    100% {
      text-shadow:
        0 0 20px var(--color-error-500),
        0 0 40px var(--color-error-500),
        0 0 60px var(--color-error-400);
    }
    50% {
      text-shadow:
        0 0 10px var(--color-error-500),
        0 0 20px var(--color-error-500),
        0 0 30px var(--color-error-400);
    }
  }

  @keyframes queue-pulse {
    0%,
    100% {
      box-shadow: 0 0 8px rgba(239, 68, 68, 0.3);
      opacity: 0.9;
    }
    50% {
      box-shadow: 0 0 12px rgba(239, 68, 68, 0.5);
      opacity: 1;
    }
  }
</style>
