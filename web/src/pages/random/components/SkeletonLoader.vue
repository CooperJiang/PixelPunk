<template>
  <div class="skeleton-loader">
    <div class="skeleton-image" />
    <div class="skeleton-details">
      <div class="skeleton-card">
        <div class="skeleton-title" />
        <div class="skeleton-line" />
        <div class="skeleton-line short" />
      </div>
      <div class="skeleton-card">
        <div class="skeleton-title" />
        <div class="skeleton-line" />
        <div class="skeleton-line" />
        <div class="skeleton-line short" />
      </div>
      <div class="skeleton-card">
        <div class="skeleton-title" />
        <div class="skeleton-line" />
        <div class="skeleton-line short" />
      </div>
    </div>
  </div>
</template>

<style scoped>
  .skeleton-loader {
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    z-index: 5;
  }

  .skeleton-image {
    width: 100%;
    height: calc(100vh - 4rem);
    background: linear-gradient(
      90deg,
      rgba(var(--color-background-700-rgb), 0.3) 25%,
      rgba(var(--color-background-600-rgb), 0.3) 50%,
      rgba(var(--color-background-700-rgb), 0.3) 75%
    );
    background-size: 200% 100%;
    animation: shimmer 2s infinite linear;
  }

  .skeleton-details {
    position: absolute;
    bottom: 0;
    left: 0;
    right: 0;
    background: rgba(var(--color-background-800-rgb), 0.7);
    backdrop-filter: blur(24px);
    border-radius: var(--radius-sm) var(--radius-sm) 0 0;
    padding: 2rem;
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 1.5rem;
    min-height: 25vh;
  }

  .skeleton-card {
    background: rgba(var(--color-background-700-rgb), 0.35);
    border-radius: var(--radius-sm);
    padding: 1.25rem;
    display: flex;
    flex-direction: column;
    gap: 0.75rem;
  }

  .skeleton-title {
    width: 40%;
    height: 1.25rem;
    background: linear-gradient(
      90deg,
      rgba(var(--color-background-600-rgb), 0.6) 25%,
      rgba(var(--color-background-500-rgb), 0.6) 50%,
      rgba(var(--color-background-600-rgb), 0.6) 75%
    );
    background-size: 200% 100%;
    animation: shimmer 2s infinite linear;
    border-radius: var(--radius-sm);
  }

  .skeleton-line {
    width: 100%;
    height: 0.875rem;
    background: linear-gradient(
      90deg,
      rgba(var(--color-background-600-rgb), 0.4) 25%,
      rgba(var(--color-background-500-rgb), 0.4) 50%,
      rgba(var(--color-background-600-rgb), 0.4) 75%
    );
    background-size: 200% 100%;
    animation: shimmer 2s infinite linear;
    border-radius: var(--radius-sm);
  }

  .skeleton-line.short {
    width: 60%;
  }

  @keyframes shimmer {
    0% {
      background-position: 200% 0;
    }
    100% {
      background-position: -200% 0;
    }
  }

  @media (max-width: 1024px) {
    .skeleton-details {
      grid-template-columns: repeat(2, 1fr);
    }
  }

  @media (max-width: 768px) {
    .skeleton-image {
      height: calc(100vh - 60px);
    }

    .skeleton-details {
      grid-template-columns: 1fr;
      min-height: 35vh;
      padding: 1.25rem 1rem;
      border-radius: var(--radius-sm) var(--radius-sm) 0 0;
    }
  }
</style>
