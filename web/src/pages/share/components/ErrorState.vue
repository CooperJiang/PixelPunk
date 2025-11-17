<script setup lang="ts">
  import { computed } from 'vue'
  import { useTexts } from '@/composables/useTexts'

  const { $t } = useTexts()

  const props = defineProps({
    type: {
      type: String,
      default: 'error',
      validator: (value: string) => ['error', 'expired', 'max-views', 'not-found', 'unauthorized'].includes(value),
    },
    title: {
      type: String,
      default: '',
    },
    message: {
      type: String,
      default: '',
    },
    actionText: {
      type: String,
      default: '',
    },
    homeIcon: {
      type: Boolean,
      default: true,
    },
  })

  const _emit = defineEmits(['home-click'])

  const displayTitle = computed(() => props.title || $t('share.error.defaultTitle'))
  const displayMessage = computed(() => props.message || $t('share.error.defaultMessage'))
  const displayActionText = computed(() => props.actionText || $t('share.error.backHome'))

  const iconClass = computed(() => {
    switch (props.type) {
      case 'expired':
        return 'fas fa-hourglass-end'
      case 'max-views':
        return 'fas fa-eye-slash'
      case 'not-found':
        return 'fas fa-unlink'
      case 'unauthorized':
        return 'fas fa-lock'
      default:
        return 'fas fa-exclamation-triangle'
    }
  })

  const errorClass = computed(() => `error-type-${props.type}`)
</script>

<template>
  <div class="error-container">
    <div class="error-card" :class="errorClass">
      <div class="error-icon">
        <i :class="iconClass" />
        <div class="error-glow" />
      </div>

      <div class="error-content">
        <h2 class="error-title">{{ displayTitle }}</h2>
        <p class="error-message">{{ displayMessage }}</p>

        <div class="error-actions">
          <slot name="actions">
            <div class="cyber-home-link-wrapper">
              <div class="cyber-home-link" @click="$emit('home-click')">
                <div class="cyber-link-content">
                  <i v-if="homeIcon" class="fas fa-home" />
                  <span>{{ displayActionText }}</span>
                </div>
                <div class="cyber-link-glow" />
              </div>
              <div class="cyber-link-lines">
                <div class="link-line line-left" />
                <div class="link-line line-right" />
              </div>
            </div>
          </slot>
        </div>
      </div>

      <div class="decoration-line line-top" />
      <div class="decoration-line line-bottom" />
      <div class="decoration-dots">
        <div class="dot dot-1" />
        <div class="dot dot-2" />
        <div class="dot dot-3" />
        <div class="dot dot-4" />
      </div>
    </div>
  </div>
</template>

<style scoped lang="scss">
  .error-container {
    display: flex;
    justify-content: center;
    align-items: center;
    min-height: 100vh;
    padding: var(--space-xl);
  }

  .error-card {
    position: relative;
    width: 100%;
    max-width: 600px;
    background: rgba(var(--color-background-800-rgb), 0.7);
    backdrop-filter: var(--backdrop-blur-md);
    border-radius: var(--radius-sm);
    overflow: hidden;
    display: flex;
    flex-direction: column;
    align-items: center;
    padding: var(--space-3xl) var(--space-xl);
    box-shadow:
      var(--shadow-xl),
      0 0 20px rgba(var(--color-brand-500-rgb), 0.1);
    border: 1.5px solid rgba(var(--color-brand-500-rgb), 0.15);
    animation: card-appear 0.5s ease-out forwards;
  }

  .error-icon {
    position: relative;
    width: 100px;
    height: 100px;
    border-radius: var(--radius-full);
    display: flex;
    align-items: center;
    justify-content: center;
    margin-bottom: 2rem;
    font-size: 3rem;
    z-index: 1;
  }

  .error-glow {
    position: absolute;
    width: 100%;
    height: 100%;
    border-radius: var(--radius-full);
    animation: pulse 2s infinite;
    z-index: -1;
  }

  .error-content {
    text-align: center;
    z-index: 1;
  }

  .error-title {
    margin: 0 0 1rem;
    font-size: 2rem;
    font-weight: 700;
    letter-spacing: 0.5px;
    background: linear-gradient(90deg, var(--color-content-heading), rgba(var(--color-content-rgb), 0.8));
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    animation: title-appear 0.6s ease-out forwards;
  }

  .error-message {
    color: var(--color-content-default);
    font-size: 1.1rem;
    line-height: 1.6;
    margin-bottom: 2rem;
    max-width: 500px;
    animation: message-appear 0.7s ease-out forwards;
  }

  .error-actions {
    animation: actions-appear 0.8s ease-out forwards;
  }

  .cyber-home-link-wrapper {
    position: relative;
    display: inline-block;
    margin-top: 0.5rem;
  }

  .cyber-home-link {
    position: relative;
    display: inline-flex;
    align-items: center;
    justify-content: center;
    padding: 0.7rem 2rem;
    cursor: pointer;
    transition: all 0.3s ease;
    color: currentColor;
    border-radius: var(--radius-sm);
    text-decoration: none;
    background: rgba(var(--color-background-800-rgb), 0.5);
    -webkit-clip-path: polygon(10px 0, 100% 0, calc(100% - 10px) 100%, 0 100%);
    clip-path: polygon(10px 0, 100% 0, calc(100% - 10px) 100%, 0 100%);
  }

  .cyber-link-content {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 0.8rem;
    position: relative;
    z-index: 2;
    font-weight: 600;
    letter-spacing: 1px;
    font-size: 1rem;
    text-transform: uppercase;
  }

  .cyber-link-content i {
    font-size: 1.2rem;
    text-shadow: 0 0 5px currentColor;
  }

  .cyber-link-glow {
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(var(--color-background-800-rgb), 0.5);
    box-shadow:
      0 0 15px currentColor,
      inset 0 0 5px var(--color-content-default);
    opacity: 0;
    transition: opacity 0.3s ease;
    z-index: 1;
  }

  .cyber-home-link-wrapper:hover .cyber-link-glow {
    opacity: 0.5;
  }

  .cyber-home-link-wrapper:hover .cyber-home-link {
    box-shadow: 0 0 15px var(--color-content-default);
    text-shadow: 0 0 5px currentColor;
  }

  .cyber-link-lines {
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    pointer-events: none;
  }

  .link-line {
    position: absolute;
    background-color: currentColor;
    transition: all 0.4s ease;
    border-radius: var(--radius-sm);
  }

  .line-left,
  .line-right {
    width: 2px;
    height: 0;
    top: 50%;
    transform: translateY(-50%);
  }

  .line-left {
    left: 4px;
  }

  .line-right {
    right: 4px;
  }

  .cyber-home-link-wrapper:hover .line-left,
  .cyber-home-link-wrapper:hover .line-right {
    height: 60%;
    box-shadow: 0 0 8px currentColor;
  }

  .cyber-home-link-wrapper::before,
  .cyber-home-link-wrapper::after {
    content: '';
    position: absolute;
    height: 2px;
    width: 0;
    background-color: currentColor;
    transition: all 0.4s ease;
    box-shadow: 0 0 5px currentColor;
    border-radius: var(--radius-sm);
    z-index: 1;
  }

  .cyber-home-link-wrapper::after {
    bottom: 5px;
    left: 8px;
  }

  .cyber-home-link-wrapper:hover::before,
  .cyber-home-link-wrapper:hover::after {
    width: calc(100% - 20px);
  }

  .cyber-home-link:active {
    transform: scale(0.95);
  }

  .cyber-home-link-wrapper::after {
    animation: glitch 2s infinite;
  }

  @keyframes glitch {
    0% {
      opacity: 1;
    }
    7% {
      opacity: 0.75;
    }
    10% {
      opacity: 1;
    }
    27% {
      opacity: 1;
    }
    30% {
      opacity: 0.75;
    }
    35% {
      opacity: 1;
    }
    52% {
      opacity: 1;
    }
    55% {
      opacity: 0.75;
    }
    60% {
      opacity: 1;
    }
    85% {
      opacity: 1;
    }
    88% {
      opacity: 0.75;
    }
    90% {
      opacity: 1;
    }
    100% {
      opacity: 1;
    }
  }

  .error-type-error .cyber-home-link-wrapper {
    color: var(--color-error-500);
  }

  .error-type-expired .cyber-home-link-wrapper {
    color: var(--color-warning-500);
  }

  .error-type-max-views .cyber-home-link-wrapper {
    color: var(--color-purple-500);
  }

  .error-type-not-found .cyber-home-link-wrapper {
    color: var(--color-cyan-500);
  }

  .error-type-unauthorized .cyber-home-link-wrapper {
    color: var(--color-brand-500);
  }

  .cyber-home-link-wrapper .cyber-home-link {
    color: inherit;
  }

  .decoration-line {
    position: absolute;
    height: 2px;
    width: 150px;
    background: linear-gradient(90deg, transparent, currentColor, transparent);
    opacity: 0.3;
  }

  .line-top {
    top: 20px;
    left: 50%;
    transform: translateX(-50%);
  }

  .line-bottom {
    bottom: 20px;
    left: 50%;
    transform: translateX(-50%);
  }

  .decoration-dots {
    position: absolute;
    width: 100%;
    height: 100%;
    z-index: 0;
  }

  .dot {
    position: absolute;
    width: 6px;
    height: 6px;
    border-radius: var(--radius-full);
    background-color: currentColor;
    opacity: 0.2;
  }

  .dot-1 {
    top: 15%;
    left: 10%;
    animation: float 3s infinite ease-in-out;
  }
  .dot-2 {
    top: 80%;
    left: 15%;
    animation: float 4s infinite ease-in-out;
  }
  .dot-3 {
    top: 25%;
    right: 15%;
    animation: float 3.5s infinite ease-in-out;
  }
  .dot-4 {
    top: 70%;
    right: 10%;
    animation: float 4.5s infinite ease-in-out;
  }

  .error-type-error {
    color: var(--color-error-500);
  }

  .error-type-error .error-glow {
    box-shadow:
      0 0 30px rgba(var(--color-error-rgb), 0.6),
      0 0 50px rgba(var(--color-error-rgb), 0.3),
      0 0 70px rgba(var(--color-error-rgb), 0.1);
  }

  .error-type-expired {
    color: var(--color-warning-500);
  }

  .error-type-expired .error-glow {
    box-shadow:
      0 0 30px rgba(var(--color-warning-rgb), 0.6),
      0 0 50px rgba(var(--color-warning-rgb), 0.3),
      0 0 70px rgba(var(--color-warning-rgb), 0.1);
  }

  .error-type-max-views {
    color: var(--color-purple-500);
  }

  .error-type-max-views .error-glow {
    box-shadow:
      0 0 30px rgba(var(--color-brand-500-rgb), 0.6),
      0 0 50px rgba(var(--color-brand-500-rgb), 0.3),
      0 0 70px rgba(var(--color-brand-500-rgb), 0.1);
  }

  .error-type-not-found {
    color: var(--color-cyan-500);
  }

  .error-type-not-found .error-glow {
    box-shadow:
      0 0 30px rgba(var(--color-info-rgb), 0.6),
      0 0 50px rgba(var(--color-info-rgb), 0.3),
      0 0 70px rgba(var(--color-info-rgb), 0.1);
  }

  .error-type-unauthorized {
    color: var(--color-brand-500);
  }

  .error-type-unauthorized .error-glow {
    box-shadow:
      0 0 30px rgba(var(--color-brand-500-rgb), 0.6),
      0 0 50px rgba(var(--color-brand-500-rgb), 0.3),
      0 0 70px rgba(var(--color-brand-500-rgb), 0.1);
  }

  @keyframes pulse {
    0% {
      transform: scale(0.95);
      opacity: 0.7;
    }
    50% {
      transform: scale(1.05);
      opacity: 0.9;
    }
    100% {
      transform: scale(0.95);
      opacity: 0.7;
    }
  }

  @keyframes float {
    0% {
      transform: translateY(0);
    }
    50% {
      transform: translateY(10px);
    }
    100% {
      transform: translateY(0);
    }
  }

  @keyframes card-appear {
    0% {
      transform: translateY(30px);
      opacity: 0;
    }
    100% {
      transform: translateY(0);
      opacity: 1;
    }
  }

  @keyframes title-appear {
    0% {
      transform: translateY(20px);
      opacity: 0;
    }
    100% {
      transform: translateY(0);
      opacity: 1;
    }
  }

  @keyframes message-appear {
    0% {
      transform: translateY(15px);
      opacity: 0;
    }
    100% {
      transform: translateY(0);
      opacity: 1;
    }
  }

  @keyframes actions-appear {
    0% {
      transform: translateY(10px);
      opacity: 0;
    }
    100% {
      transform: translateY(0);
      opacity: 1;
    }
  }

  @media (max-width: 768px) {
    .error-container {
      padding: 1rem;
    }

    .error-card {
      padding: 2rem 1.5rem;
    }

    .error-icon {
      width: 80px;
      height: 80px;
      font-size: 2.5rem;
      margin-bottom: 1.5rem;
    }

    .error-title {
      font-size: 1.75rem;
    }

    .error-message {
      font-size: 1rem;
    }

    .cyber-link-content {
      font-size: 0.9rem;
    }
  }
</style>
