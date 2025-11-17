<script setup lang="ts">
  import { computed } from 'vue'
  import { useTexts } from '@/composables/useTexts'

  interface Props {
    error: string
  }

  interface Emits {
    (e: 'retry'): void
  }

  const props = defineProps<Props>()
  const emit = defineEmits<Emits>()
  const { $t } = useTexts()

  const isNoContentError = computed(
    () =>
      props.error.includes($t('random.errorState.keywords.noFiles1')) ||
      props.error.includes($t('random.errorState.keywords.noFiles2')) ||
      props.error.includes($t('random.errorState.keywords.explore')) ||
      props.error.includes($t('random.errorState.keywords.noRecommend1')) ||
      props.error.includes($t('random.errorState.keywords.noRecommend2'))
  )

  const iconClass = computed(() => (isNoContentError.value ? 'fas fa-images' : 'fas fa-exclamation-triangle'))

  const displayTitle = computed(() =>
    isNoContentError.value ? $t('random.errorState.noContentTitle') : $t('random.errorState.fetchFailedTitle')
  )

  const displayMessage = computed(() => {
    if (isNoContentError.value) {
      return $t('random.errorState.noContentMessage')
    }
    return $t('random.errorState.fetchFailedMessage')
  })

  const buttonText = computed(() =>
    isNoContentError.value ? $t('random.errorState.refreshButton') : $t('random.errorState.retryButton')
  )

  const buttonIconClass = computed(() => (isNoContentError.value ? 'fas fa-sync-alt' : 'fas fa-redo'))

  const handleRetry = () => {
    emit('retry')
  }
</script>

<template>
  <div class="fullscreen-center">
    <div class="empty-state-container">
      <div class="empty-content">
        <div class="empty-icon">
          <i :class="iconClass" />
        </div>

        <h3 class="empty-title">{{ displayTitle }}</h3>

        <p class="empty-description">{{ displayMessage }}</p>

        <div class="empty-actions">
          <button class="retry-btn" @click="handleRetry">
            <i :class="buttonIconClass + ' mr-2'" />
            {{ buttonText }}
          </button>
        </div>

        <div class="empty-decorations">
          <div class="decoration-circle circle-1" :class="{ 'error-theme': !isNoContentError }" />
          <div class="decoration-circle circle-2" :class="{ 'error-theme': !isNoContentError }" />
          <div class="decoration-circle circle-3" :class="{ 'error-theme': !isNoContentError }" />
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
  .fullscreen-center {
    height: calc(100vh - 4.5rem);
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 0 1rem;
  }

  .empty-state-container {
    @apply flex min-h-[400px] items-center justify-center px-4 py-12;
    background: transparent;
  }

  .empty-content {
    @apply relative mx-auto max-w-md text-center;
    animation: fadeInUp 0.6s ease-out;
  }

  .empty-icon {
    @apply relative mb-6;
  }

  .empty-icon i {
    @apply text-6xl;
    filter: drop-shadow(0 0 20px rgba(var(--color-brand-500-rgb), 0.3));
    animation: floatIcon 3s ease-in-out infinite;
    color: var(--color-brand-500);
  }

  .empty-icon i.fa-exclamation-triangle {
    color: var(--color-error-500);
    filter: drop-shadow(0 0 20px rgba(var(--color-error-rgb), 0.3));
  }

  .empty-title {
    @apply mb-4 text-xl font-semibold text-content-heading;
    background: linear-gradient(135deg, var(--color-brand-500) 0%, var(--color-content-heading) 50%, var(--color-error-500) 100%);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    background-clip: text;
  }

  .empty-description {
    @apply mb-8 text-sm leading-relaxed text-content-muted;
    line-height: 1.6;
    white-space: pre-line;
  }

  .empty-actions {
    @apply flex justify-center;
  }

  .retry-btn {
    @apply flex items-center rounded-lg px-6 py-3 text-sm font-medium transition-all duration-300;
    background: linear-gradient(135deg, rgba(var(--color-brand-500-rgb), 0.1), rgba(var(--color-brand-500-rgb), 0.05));
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.3);
    color: rgba(var(--color-content-rgb), 0.9);
    backdrop-filter: blur(10px);
  }

  .retry-btn:hover {
    @apply -translate-y-1 transform;
    background: linear-gradient(135deg, rgba(var(--color-brand-500-rgb), 0.2), rgba(var(--color-brand-500-rgb), 0.1));
    border-color: rgba(var(--color-brand-500-rgb), 0.5);
    box-shadow: 0 8px 25px rgba(var(--color-brand-500-rgb), 0.3);
  }

  .empty-decorations {
    @apply pointer-events-none absolute inset-0 overflow-hidden;
  }

  .decoration-circle {
    @apply absolute rounded-full opacity-20;
    background: linear-gradient(135deg, rgba(var(--color-brand-500-rgb), 0.1), rgba(var(--color-error-rgb), 0.1));
    filter: blur(1px);
  }

  .decoration-circle.error-theme {
    background: linear-gradient(135deg, rgba(var(--color-error-rgb), 0.1), rgba(var(--color-error-rgb), 0.05));
  }

  .circle-1 {
    @apply -left-8 -top-8 h-32 w-32;
    animation: floatCircle1 4s ease-in-out infinite;
  }

  .circle-2 {
    @apply -bottom-4 -right-4 h-24 w-24;
    animation: floatCircle2 3.5s ease-in-out infinite;
  }

  .circle-3 {
    @apply -left-4 top-1/2 h-16 w-16;
    animation: floatCircle3 4.5s ease-in-out infinite;
  }

  @keyframes fadeInUp {
    from {
      opacity: 0;
      transform: translateY(30px);
    }
    to {
      opacity: 1;
      transform: translateY(0);
    }
  }

  @keyframes floatIcon {
    0%,
    100% {
      transform: translateY(0px);
    }
    50% {
      transform: translateY(-10px);
    }
  }

  @keyframes floatCircle1 {
    0%,
    100% {
      transform: translate(0, 0) rotate(0deg);
    }
    33% {
      transform: translate(10px, -10px) rotate(120deg);
    }
    66% {
      transform: translate(-5px, 10px) rotate(240deg);
    }
  }

  @keyframes floatCircle2 {
    0%,
    100% {
      transform: translate(0, 0) rotate(0deg);
    }
    50% {
      transform: translate(-8px, -8px) rotate(180deg);
    }
  }

  @keyframes floatCircle3 {
    0%,
    100% {
      transform: translate(0, 0) scale(1);
    }
    50% {
      transform: translate(5px, -5px) scale(1.1);
    }
  }

  @media (max-width: 768px) {
    .fullscreen-center {
      height: calc(100vh - 3.5rem);
    }

    .empty-state-container {
      @apply min-h-[300px] py-8;
    }

    .empty-icon i {
      @apply text-5xl;
    }

    .empty-title {
      @apply text-lg;
    }

    .empty-description {
      @apply px-2 text-xs;
    }

    .decoration-circle {
      @apply opacity-10;
    }

    .circle-1 {
      @apply h-24 w-24;
    }

    .circle-2 {
      @apply h-16 w-16;
    }

    .circle-3 {
      @apply h-12 w-12;
    }
  }
</style>
