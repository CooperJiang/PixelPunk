<script setup lang="ts">
  import { useTexts } from '@/composables/useTexts'

  interface Props {
    type: 'initial' | 'image'
  }

  defineProps<Props>()
  const { $t } = useTexts()
</script>

<template>
  <div class="fullscreen-center">
    <div v-if="type === 'initial'" class="cyber-loading-custom">
      <div class="loading-hexagon">
        <div class="hexagon-inner">
          <div class="hexagon-core" />
        </div>
      </div>
      <div class="loading-text">
        <span class="text-content">{{ $t('random.loadingState.initial.searching') }}</span>
        <div class="loading-dots"><span /><span /><span /></div>
      </div>
      <p class="loading-subtitle">{{ $t('random.loadingState.initial.selecting') }}</p>
    </div>

    <div v-else-if="type === 'image'" class="image-loading-overlay">
      <div class="pulse-loader">
        <div class="pulse-ring pulse-ring-1" />
        <div class="pulse-ring pulse-ring-2" />
        <div class="pulse-ring pulse-ring-3" />
        <div class="pulse-center">
          <i class="fas fa-image" />
        </div>
      </div>
      <p class="loading-tip">{{ $t('random.loadingState.image.loading') }}</p>
      <p class="loading-subtitle">{{ $t('random.loadingState.image.loadingHD') }}</p>
    </div>
  </div>
</template>

<style scoped>
  .fullscreen-center {
    height: calc(100vh - 4.5rem);
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .cyber-loading-custom {
    text-align: center;
  }

  .loading-hexagon {
    position: relative;
    width: 80px;
    height: 80px;
    margin: 0 auto 2rem;
  }

  .hexagon-inner {
    width: 100%;
    height: 100%;
    position: relative;
    transform: rotate(30deg);
    animation: hexRotate 2s linear infinite;
  }

  .hexagon-core {
    width: 100%;
    height: 100%;
    background: linear-gradient(45deg, var(--color-brand-500), var(--color-error-500));
    clip-path: polygon(50% 0%, 100% 25%, 100% 75%, 50% 100%, 0% 75%, 0% 25%);
    position: relative;
  }

  @keyframes hexRotate {
    0% {
      transform: rotate(30deg);
    }
    100% {
      transform: rotate(390deg);
    }
  }

  .loading-text {
    color: var(--color-brand-500);
    font-size: 1.1rem;
    margin-bottom: 0.5rem;
  }

  .loading-subtitle {
    color: var(--color-content-muted);
    font-size: 0.85rem;
    margin-top: 0.5rem;
    opacity: 0.8;
  }

  .loading-dots {
    display: inline-flex;
    gap: 4px;
    margin-left: 8px;
  }

  .loading-dots span {
    width: 4px;
    height: 4px;
    background: var(--color-brand-500);
    border-radius: var(--radius-full);
    animation: dotPulse 1.4s infinite ease-in-out;
  }

  .loading-dots span:nth-child(1) {
    animation-delay: -0.32s;
  }
  .loading-dots span:nth-child(2) {
    animation-delay: -0.16s;
  }
  .loading-dots span:nth-child(3) {
    animation-delay: 0s;
  }

  @keyframes dotPulse {
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

  .image-loading-overlay {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    height: 100%;
  }

  .pulse-loader {
    position: relative;
    width: 120px;
    height: 120px;
    margin-bottom: 2rem;
  }

  .pulse-ring {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    border: 2px solid transparent;
    border-radius: var(--radius-full);
    animation: pulseRing 2s linear infinite;
  }

  .pulse-ring-1 {
    border-top-color: var(--color-brand-500);
    animation-delay: 0s;
  }

  .pulse-ring-2 {
    border-right-color: var(--color-error-500);
    animation-delay: 0.4s;
  }

  .pulse-ring-3 {
    border-bottom-color: var(--color-brand-500);
    animation-delay: 0.8s;
  }

  .pulse-center {
    position: absolute;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    width: 40px;
    height: 40px;
    background: linear-gradient(45deg, var(--color-brand-500), var(--color-error-500));
    border-radius: var(--radius-full);
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 18px;
    color: var(--color-text-on-brand);
    animation: centerPulse 1s ease-in-out infinite alternate;
  }

  @keyframes pulseRing {
    0% {
      transform: scale(0.8);
      opacity: 1;
    }
    100% {
      transform: scale(1.2);
      opacity: 0;
    }
  }

  @keyframes centerPulse {
    0% {
      transform: translate(-50%, -50%) scale(1);
    }
    100% {
      transform: translate(-50%, -50%) scale(1.1);
    }
  }

  .loading-tip {
    color: var(--color-brand-500);
    font-size: 1.1rem;
    text-shadow: 0 0 10px rgba(var(--color-brand-500-rgb), 0.5);
    margin-bottom: 0.5rem;
  }

  .image-loading-overlay .loading-subtitle {
    color: var(--color-content-muted);
    font-size: 0.85rem;
    margin-top: 0.5rem;
    opacity: 0.8;
  }

  @media (max-width: 768px) {
    .fullscreen-center {
      height: calc(100vh - 3.5rem);
    }
  }
</style>
