<script setup lang="ts">
  import { useTexts } from '@/composables/useTexts'

  const { $t } = useTexts()

  defineOptions({
    name: 'DocsIntro',
  })

  const props = defineProps<Props>()

  interface Props {
    scrollToSection: (target: string) => void
  }

  const onScrollToSection = (target: string) => {
    props.scrollToSection(target)
  }
</script>

<template>
  <div class="docs-intro">
    <div class="intro-background">
      <div class="intro-grid" />
      <div class="intro-particles" />
    </div>

    <div class="title-section">
      <div class="title-wrapper">
        <h1 class="page-title">
          <i class="fas fa-code title-icon" />
          <span class="title-text">{{ $t('docs.intro.pageTitle') }}</span>
          <div class="title-glow" />
        </h1>
        <p class="page-subtitle">{{ $t('docs.intro.subtitle') }}</p>
      </div>

      <div class="status-indicators">
        <div class="status-item">
          <div class="status-dot online" />
          <span>{{ $t('docs.intro.status.online') }}</span>
        </div>
        <div class="status-item">
          <div class="status-dot" />
          <span>{{ $t('docs.intro.status.version') }}</span>
        </div>
        <div class="status-item">
          <div class="status-dot perf" />
          <span>{{ $t('docs.intro.status.detailed') }}</span>
        </div>
      </div>
    </div>

    <div class="quick-actions">
      <CyberButton type="primary" @click="onScrollToSection('#api-examples')">
        <i class="fas fa-code" />
        <span>{{ $t('docs.intro.actions.examples') }}</span>
      </CyberButton>
      <CyberButton type="success" @click="onScrollToSection('#api-upload')">
        <i class="fas fa-upload" />
        <span>{{ $t('docs.intro.actions.start') }}</span>
      </CyberButton>
      <CyberButton type="warning" @click="onScrollToSection('#api-tester')">
        <i class="fas fa-flask" />
        <span>{{ $t('docs.intro.actions.test') }}</span>
      </CyberButton>
    </div>
  </div>
</template>

<style scoped>
  .docs-intro {
    position: relative;
    background: linear-gradient(
      135deg,
      rgba(var(--color-brand-500-rgb), 0.08) 0%,
      rgba(var(--color-badge-accent-text-rgb), 0.06) 50%,
      rgba(var(--color-success-rgb), 0.04) 100%
    );
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.2);
    border-radius: var(--radius-sm);
    padding: 2rem 1.8rem;
    margin-bottom: 2rem;
    overflow: hidden;
    box-shadow:
      0 12px 24px rgba(var(--color-background-700-rgb), 0.8),
      0 0 20px rgba(var(--color-brand-500-rgb), 0.08),
      inset 0 1px 0 rgba(var(--color-content-rgb), 0.1);
    z-index: 1;
  }

  .intro-background {
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    pointer-events: none;
    z-index: 0;
  }

  .intro-grid {
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background-image:
      linear-gradient(rgba(var(--color-brand-500-rgb), 0.03) 1px, transparent 1px),
      linear-gradient(90deg, rgba(var(--color-brand-500-rgb), 0.03) 1px, transparent 1px);
    background-size: 30px 30px;
    animation: grid-move 20s linear infinite;
  }

  @keyframes grid-move {
    0% {
      transform: translate(0, 0);
    }
    100% {
      transform: translate(30px, 30px);
    }
  }

  .intro-particles {
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
  }

  .intro-particles::before,
  .intro-particles::after {
    content: '';
    position: absolute;
    width: 4px;
    height: 4px;
    background: var(--color-brand-500);
    border-radius: var(--radius-full);
    box-shadow: 0 0 10px rgba(var(--color-brand-500-rgb), 0.8);
  }

  .intro-particles::after {
    top: 70%;
    right: 25%;
    animation: particle-float 8s ease-in-out infinite reverse;
  }

  @keyframes particle-float {
    0%,
    100% {
      transform: translateY(0px) scale(1);
      opacity: 0.6;
    }
    50% {
      transform: translateY(-20px) scale(1.2);
      opacity: 1;
    }
  }

  .title-section {
    position: relative;
    z-index: 1;
    text-align: center;
    margin-bottom: 1.8rem;
  }

  .title-wrapper {
    position: relative;
    display: inline-block;
  }

  .page-title {
    font-size: 2.5rem;
    font-weight: 800;
    margin: 0 0 0.8rem 0;
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 0.8rem;
    position: relative;
    z-index: 1;
  }

  .title-icon {
    color: var(--color-brand-500);
    text-shadow: 0 0 20px rgba(var(--color-brand-500-rgb), 0.6);
    animation: icon-pulse 3s ease-in-out infinite;
  }

  @keyframes icon-pulse {
    0%,
    100% {
      transform: scale(1);
      filter: drop-shadow(0 0 10px rgba(var(--color-brand-500-rgb), 0.4));
    }
    50% {
      transform: scale(1.1);
      filter: drop-shadow(0 0 20px rgba(var(--color-brand-500-rgb), 0.8));
    }
  }

  .title-text {
    background: linear-gradient(135deg, var(--color-brand-500) 0%, var(--color-error-500) 50%, var(--color-success-500) 100%);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    background-clip: text;
    animation: title-glow 4s ease-in-out infinite;
  }

  @keyframes title-glow {
    0%,
    100% {
      filter: drop-shadow(0 0 10px rgba(var(--color-brand-500-rgb), 0.3));
    }
    50% {
      filter: drop-shadow(0 0 20px rgba(var(--color-badge-accent-text-rgb), 0.5));
    }
  }

  .title-glow {
    position: absolute;
    top: 50%;
    left: 50%;
    width: 200%;
    height: 200%;
    background: radial-gradient(circle, rgba(var(--color-brand-500-rgb), 0.2) 0%, transparent 70%);
    transform: translate(-50%, -50%);
    z-index: -1;
    animation: glow-pulse 3s ease-in-out infinite;
  }

  @keyframes glow-pulse {
    0%,
    100% {
      opacity: 0.3;
      transform: translate(-50%, -50%) scale(1);
    }
    50% {
      opacity: 0.6;
      transform: translate(-50%, -50%) scale(1.2);
    }
  }

  .page-subtitle {
    font-size: 1.1rem;
    color: var(--color-gray-200);
    margin: 0;
    line-height: 1.5;
    font-weight: 500;
    opacity: 0.9;
  }

  .status-indicators {
    display: flex;
    justify-content: center;
    gap: 2rem;
    margin-top: 1.5rem;
    flex-wrap: wrap;
  }

  .status-item {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    font-size: 0.85rem;
    color: var(--color-content-muted);
    font-weight: 500;
  }

  .status-dot {
    width: 8px;
    height: 8px;
    border-radius: var(--radius-full);
    background: var(--color-content-disabled);
    animation: dot-pulse 2s infinite;
  }

  .status-dot.online {
    background: var(--color-success-500);
    box-shadow: 0 0 10px rgba(var(--color-success-rgb), 0.6);
  }

  .status-dot.perf {
    background: var(--color-error-500);
    box-shadow: 0 0 10px rgba(var(--color-badge-accent-text-rgb), 0.6);
    animation-delay: 0.5s;
  }

  @keyframes dot-pulse {
    0%,
    100% {
      transform: scale(1);
      opacity: 1;
    }
    50% {
      transform: scale(1.3);
      opacity: 0.7;
    }
  }

  .quick-actions {
    display: flex;
    justify-content: center;
    gap: 1rem;
    flex-wrap: wrap;
    position: relative;
    z-index: 1;
  }

  @media (max-width: 768px) {
    .page-title {
      font-size: 2rem;
      gap: 0.6rem;
    }

    .status-indicators {
      gap: 1rem;
    }

    .quick-actions {
      gap: 0.75rem;
    }
  }
</style>
