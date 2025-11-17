<script setup lang="ts">
  defineEmits<{
    (e: 'reset'): void
  }>()
</script>

<template>
  <div class="empty-state-container">
    <div class="empty-content">
      <div class="empty-bg-decoration">
        <div class="glow-orb glow-1" />
        <div class="glow-orb glow-2" />
        <div class="glow-orb glow-3" />
      </div>

      <div class="empty-icon-wrapper">
        <div class="icon-circle">
          <i class="fas fa-images" />
        </div>
        <div class="icon-ring ring-1" />
        <div class="icon-ring ring-2" />
        <div class="icon-ring ring-3" />
      </div>

      <div class="empty-text-content">
        <h3 class="empty-title">{{ $t('explore.empty.title') }}</h3>
        <p class="empty-description">{{ $t('explore.empty.description') }}</p>
      </div>
    </div>
  </div>
</template>

<style scoped>
  .empty-state-container {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 100%;
    max-width: 480px;
    margin: 0 auto;
    padding: 2rem 1rem;
  }

  .empty-content {
    position: relative;
    width: 100%;
    display: flex;
    flex-direction: column;
    align-items: center;
    background: linear-gradient(145deg, rgba(var(--color-background-800-rgb), 0.3), rgba(var(--color-background-700-rgb), 0.2));
    border-radius: var(--radius-sm);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.2);
    padding: 2.5rem 2rem;
    backdrop-filter: blur(20px);
    box-shadow:
      0 8px 32px rgba(0, 0, 0, 0.1),
      0 0 0 1px rgba(var(--color-brand-500-rgb), 0.05) inset;
    animation: fadeInScale 0.6s cubic-bezier(0.4, 0, 0.2, 1);
    overflow: hidden;
  }

  .empty-bg-decoration {
    position: absolute;
    inset: 0;
    pointer-events: none;
    overflow: hidden;
    border-radius: var(--radius-sm);
  }

  .glow-orb {
    position: absolute;
    border-radius: var(--radius-full);
    filter: blur(60px);
    opacity: 0.2;
    animation: orbFloat 8s ease-in-out infinite;
  }

  .glow-1 {
    width: 300px;
    height: 300px;
    background: radial-gradient(circle, rgba(var(--color-brand-500-rgb), 0.2) 0%, transparent 70%);
    top: -100px;
    right: -100px;
  }

  .glow-2 {
    width: 250px;
    height: 250px;
    background: radial-gradient(circle, rgba(var(--color-info-rgb), 0.15) 0%, transparent 70%);
    bottom: -80px;
    left: -80px;
    animation-delay: 3s;
  }

  .glow-3 {
    display: none;
  }

  .empty-icon-wrapper {
    position: relative;
    width: 100px;
    height: 100px;
    margin-bottom: 1.5rem;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .icon-circle {
    position: relative;
    width: 70px;
    height: 70px;
    display: flex;
    align-items: center;
    justify-content: center;
    background: linear-gradient(145deg, rgba(var(--color-brand-500-rgb), 0.15), rgba(var(--color-brand-500-rgb), 0.08));
    border-radius: var(--radius-full);
    border: 2px solid rgba(var(--color-brand-500-rgb), 0.25);
    box-shadow:
      0 0 40px rgba(var(--color-brand-500-rgb), 0.2),
      0 0 0 1px rgba(var(--color-brand-500-rgb), 0.1) inset;
    animation: iconPulse 3s ease-in-out infinite;
    z-index: 2;
  }

  .icon-circle i {
    font-size: 32px;
    color: var(--color-brand-500);
    filter: drop-shadow(0 2px 8px rgba(var(--color-brand-500-rgb), 0.6));
  }

  .icon-ring {
    position: absolute;
    border-radius: var(--radius-full);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.15);
    animation: ringExpand 3s ease-in-out infinite;
  }

  .ring-1 {
    width: 100%;
    height: 100%;
    animation-delay: 0s;
  }

  .ring-2 {
    width: 90%;
    height: 90%;
    animation-delay: 1s;
    border-color: rgba(var(--color-info-rgb), 0.12);
  }

  .ring-3 {
    width: 110%;
    height: 110%;
    animation-delay: 2s;
    border-color: rgba(var(--color-brand-500-rgb), 0.08);
  }

  .empty-text-content {
    position: relative;
    z-index: 1;
    text-align: center;
  }

  .empty-title {
    margin: 0 0 0.75rem 0;
    font-size: 1.25rem;
    font-weight: 700;
    color: var(--color-content-heading);
    background: linear-gradient(135deg, var(--color-brand-500) 0%, var(--color-info-500) 100%);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    background-clip: text;
    letter-spacing: 0.5px;
  }

  .empty-description {
    margin: 0;
    font-size: 0.875rem;
    line-height: 1.6;
    color: var(--color-content-muted);
  }

  @keyframes fadeInScale {
    from {
      opacity: 0;
      transform: scale(0.95);
    }
    to {
      opacity: 1;
      transform: scale(1);
    }
  }

  @keyframes orbFloat {
    0%,
    100% {
      transform: translate(0, 0);
    }
    50% {
      transform: translate(10px, -10px);
    }
  }

  @keyframes iconPulse {
    0%,
    100% {
      transform: scale(1);
      box-shadow:
        0 0 40px rgba(var(--color-brand-500-rgb), 0.2),
        0 0 0 1px rgba(var(--color-brand-500-rgb), 0.1) inset;
    }
    50% {
      transform: scale(1.05);
      box-shadow:
        0 0 50px rgba(var(--color-brand-500-rgb), 0.3),
        0 0 0 1px rgba(var(--color-brand-500-rgb), 0.15) inset;
    }
  }

  @keyframes ringExpand {
    0%,
    100% {
      transform: scale(1);
      opacity: 0.6;
    }
    50% {
      transform: scale(1.1);
      opacity: 0.3;
    }
  }

  @media (max-width: 768px) {
    .empty-state-container {
      padding: 1.5rem 1rem;
    }

    .empty-content {
      padding: 2rem 1.5rem;
    }

    .empty-icon-wrapper {
      width: 90px;
      height: 90px;
      margin-bottom: 1.25rem;
    }

    .icon-circle {
      width: 60px;
      height: 60px;
    }

    .icon-circle i {
      font-size: 28px;
    }

    .empty-title {
      font-size: 1.125rem;
    }

    .empty-description {
      font-size: 0.8125rem;
    }
  }

  @media (max-width: 480px) {
    .empty-content {
      padding: 1.75rem 1.25rem;
      border-radius: var(--radius-sm);
    }

    .empty-icon-wrapper {
      width: 80px;
      height: 80px;
      margin-bottom: 1rem;
    }

    .icon-circle {
      width: 55px;
      height: 55px;
    }

    .icon-circle i {
      font-size: 24px;
    }

    .empty-title {
      font-size: 1rem;
    }

    .empty-description {
      font-size: 0.75rem;
    }
  }
</style>
