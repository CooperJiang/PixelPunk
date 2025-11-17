<script setup lang="ts">
  import type { LoadingProps } from './types'
  import { useTexts } from '@/composables/useTexts'
  import { computed } from 'vue'

  defineOptions({
    name: 'CyberLoading',
  })

  const props = withDefaults(defineProps<LoadingProps>(), {
    visible: false,
    fullScreen: false,
    text: '',
    background: 'var(--color-overlay-heavy)',
    textColor: 'var(--color-brand-500)',
  })

  const { $t } = useTexts()

  const displayText = computed(() => props.text || $t('status.loading'))
</script>

<template>
  <div
    v-if="visible"
    class="cyber-loading"
    :class="[{ 'full-screen': fullScreen }]"
    :style="{
      background: fullScreen ? background : 'transparent',
    }"
  >
    <div class="loading-container">
      <div class="loading-animation">
        <div class="hexagon">
          <div class="hex-inner" />
        </div>
        <div class="glitch-circles">
          <div class="circle circle-1" />
          <div class="circle circle-2" />
          <div class="circle circle-3" />
        </div>
        <div class="scan-line" />
      </div>
      <div class="fallback-spinner" aria-hidden="true">
        <i class="fas fa-circle-notch fa-spin" />
      </div>
      <div v-if="displayText" class="loading-text" :style="{ color: textColor }">
        <span class="text-glitch" :data-text="displayText">{{ displayText }}</span>
      </div>
    </div>
  </div>
</template>

<style scoped>
  .cyber-loading {
    display: flex;
    align-items: center;
    justify-content: center;
    backdrop-filter: blur(3px);
  }

  .cyber-loading.full-screen {
    position: fixed !important;
    top: 0 !important;
    left: 0 !important;
    right: 0 !important;
    bottom: 0 !important;
    width: 100vw !important;
    height: 100vh !important;
    z-index: 10000 !important;
    margin: 0 !important;
    padding: 0 !important;
  }

  .cyber-loading:not(.full-screen) {
    position: relative;
    width: 100%;
    height: 100%;
    min-height: 200px;
  }

  .loading-container {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: 1.5rem;
    transform: none;
    position: relative;
  }

  .fallback-spinner {
    @apply flex items-center justify-center;
    color: var(--color-brand-500);
    font-size: 20px;
    line-height: 1;
  }

  .loading-animation {
    position: relative;
    display: flex;
    align-items: center;
    justify-content: center;
    width: 120px;
    height: 120px;
    transform: none;
  }

  .hexagon {
    @apply relative;
    width: 70px;
    height: 80px;
  }

  .full-screen .hexagon {
    transform: rotate(30deg);
  }

  .cyber-loading:not(.full-screen) .hexagon {
    transform: rotate(0deg);
  }

  .hexagon::before,
  .hexagon::after,
  .hex-inner::before,
  .hex-inner::after {
    @apply absolute inset-0 border-2 border-transparent;
    content: '';
    background: transparent;
    border-color: var(--color-brand-500);
    animation: hexPulse 1.5s infinite alternate;
  }

  .full-screen .hexagon::before,
  .full-screen .hexagon::after,
  .full-screen .hex-inner::before,
  .full-screen .hex-inner::after {
    animation:
      hexPulse 1.5s infinite alternate,
      hexGlitchRotated 3s infinite;
  }

  .cyber-loading:not(.full-screen) .hexagon::before,
  .cyber-loading:not(.full-screen) .hexagon::after,
  .cyber-loading:not(.full-screen) .hex-inner::before,
  .cyber-loading:not(.full-screen) .hex-inner::after {
    animation:
      hexPulse 1.5s infinite alternate,
      hexGlitchFlat 3s infinite;
  }

  .hexagon::after {
    clip-path: polygon(50% 0%, 100% 25%, 100% 75%, 50% 100%, 0% 75%, 0% 25%);
    border-color: var(--color-error-500);
    opacity: 0.7;
    animation-delay: 0.3s;
  }

  .hex-inner {
    @apply absolute inset-0;
  }

  .hex-inner::after {
    clip-path: polygon(50% 20%, 80% 35%, 80% 65%, 50% 80%, 20% 65%, 20% 35%);
    border-color: var(--color-error-500);
    opacity: 0.5;
    transform: scale(0.6);
    animation-delay: 0.9s;
  }

  .glitch-circles {
    @apply absolute inset-0 flex items-center justify-center;
  }

  .circle {
    @apply absolute rounded-full border border-transparent;
    animation: circlePulse 2s infinite;
  }

  .circle-1 {
    @apply h-12 w-12;
    border-top-color: var(--color-brand-500);
    border-right-color: var(--color-brand-500);
    animation-delay: 0s;
  }

  .circle-2 {
    width: 70px;
    height: 70px;
    border-right-color: var(--color-error-500);
    border-left-color: var(--color-error-500);
    animation-delay: 0.3s;
  }

  .circle-3 {
    width: 90px;
    height: 90px;
    border-bottom-color: var(--color-brand-500);
    border-left-color: var(--color-brand-500);
    animation-delay: 0.6s;
  }

  .scan-line {
    @apply absolute w-full;
    height: 2px;
    background: linear-gradient(to right, transparent, var(--color-brand-500), var(--color-error-500), transparent);
    animation: scanLine 1.5s infinite;
  }

  .loading-text {
    font-family: 'Orbitron', sans-serif;
    font-size: 1.2rem;
    font-weight: 700;
    letter-spacing: 1px;
    text-shadow: 0 0 10px currentColor;
  }

  .text-glitch {
    @apply relative;
    animation: textFlicker 3s infinite;
  }

  /* 移除文字叠加效果 */

  .cyber-loading-fade-enter-active,
  .cyber-loading-fade-leave-active {
    @apply transition-opacity duration-300;
  }

  .cyber-loading-fade-enter-from,
  .cyber-loading-fade-leave-to {
    @apply opacity-0;
  }

  /* Keyframes animations */
  @keyframes hexPulse {
    0% {
      opacity: 0.6;
      transform: scale(1);
    }
    100% {
      opacity: 1;
      transform: scale(1.05);
    }
  }

  @keyframes hexGlitchRotated {
    0%,
    100% {
      transform: rotate(30deg) translate(0, 0);
    }
    20% {
      transform: rotate(30deg) translate(-2px, 2px);
    }
    40% {
      transform: rotate(30deg) translate(2px, -2px);
    }
    60% {
      transform: rotate(30deg) translate(-2px, -2px);
    }
    80% {
      transform: rotate(30deg) translate(2px, 2px);
    }
  }

  @keyframes hexGlitchFlat {
    0%,
    100% {
      transform: translate(0, 0);
    }
    20% {
      transform: translate(-2px, 2px);
    }
    40% {
      transform: translate(2px, -2px);
    }
    60% {
      transform: translate(-2px, -2px);
    }
    80% {
      transform: translate(2px, 2px);
    }
  }

  @keyframes circlePulse {
    0% {
      transform: rotate(0deg) scale(1);
      opacity: 0.6;
    }
    50% {
      opacity: 1;
    }
    100% {
      transform: rotate(360deg) scale(1.1);
      opacity: 0.6;
    }
  }

  @keyframes scanLine {
    0% {
      top: 0%;
      opacity: 0;
    }
    10% {
      opacity: 1;
    }
    90% {
      opacity: 1;
    }
    100% {
      top: 100%;
      opacity: 0;
    }
  }

  @keyframes textFlicker {
    0%,
    100% {
      opacity: 1;
    }
    41.99% {
      opacity: 1;
    }
    42% {
      opacity: 0;
    }
    43% {
      opacity: 0;
    }
    43.01% {
      opacity: 1;
    }
    47.99% {
      opacity: 1;
    }
    48% {
      opacity: 0;
    }
    49% {
      opacity: 0;
    }
    49.01% {
      opacity: 1;
    }
  }

  @keyframes textGlitch {
    0% {
      clip: rect(42px, 9999px, 44px, 0);
    }
    5% {
      clip: rect(12px, 9999px, 59px, 0);
    }
    10% {
      clip: rect(48px, 9999px, 29px, 0);
    }
    15% {
      clip: rect(42px, 9999px, 73px, 0);
    }
    20% {
      clip: rect(63px, 9999px, 27px, 0);
    }
    25% {
      clip: rect(34px, 9999px, 55px, 0);
    }
    30% {
      clip: rect(86px, 9999px, 73px, 0);
    }
    35% {
      clip: rect(20px, 9999px, 20px, 0);
    }
    40% {
      clip: rect(26px, 9999px, 60px, 0);
    }
    45% {
      clip: rect(25px, 9999px, 66px, 0);
    }
    50% {
      clip: rect(57px, 9999px, 98px, 0);
    }
    55% {
      clip: rect(5px, 9999px, 46px, 0);
    }
    60% {
      clip: rect(82px, 9999px, 31px, 0);
    }
    65% {
      clip: rect(54px, 9999px, 27px, 0);
    }
    70% {
      clip: rect(28px, 9999px, 99px, 0);
    }
    75% {
      clip: rect(45px, 9999px, 69px, 0);
    }
    80% {
      clip: rect(23px, 9999px, 85px, 0);
    }
    85% {
      clip: rect(54px, 9999px, 84px, 0);
    }
    90% {
      clip: rect(45px, 9999px, 47px, 0);
    }
    95% {
      clip: rect(37px, 9999px, 20px, 0);
    }
    100% {
      clip: rect(4px, 9999px, 91px, 0);
    }
  }
</style>
