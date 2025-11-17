<template>
  <div class="stats-card" :data-loading="loading">
    <div class="stats-container">
      <div class="stats-icon">
        <i :class="icon" />
      </div>
      <div class="stats-content">
        <div class="stats-value">{{ formattedValue }}</div>
        <div class="stats-label">{{ label }}</div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
  import { computed } from 'vue'
  import type { StatsCardProps } from './types'

  defineOptions({
    name: 'CyberStatsCard',
  })

  const props = withDefaults(defineProps<StatsCardProps>(), {
    loading: false,
  })

  const formattedValue = computed(() => {
    if (props.loading) return '---'

    if (typeof props.value === 'number') {
      return props.value.toLocaleString()
    }

    return props.value
  })
</script>

<style scoped lang="scss">
  .stats-card {
    background: linear-gradient(
      135deg,
      rgba(var(--color-brand-500-rgb), 0.05) 0%,
      rgba(var(--color-purple-500-rgb), 0.03) 50%,
      rgba(var(--color-error-rgb), 0.02) 100%
    );
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.15);
    border-radius: var(--radius-sm);
    padding: 1.25rem;
    backdrop-filter: blur(10px);
    position: relative;
    overflow: hidden;
    transition: all 0.3s ease;

    &::before {
      content: '';
      position: absolute;
      top: 0;
      left: 0;
      right: 0;
      bottom: 0;
      background: linear-gradient(
        135deg,
        rgba(var(--color-brand-500-rgb), 0.1) 0%,
        transparent 30%,
        transparent 70%,
        rgba(var(--color-purple-500-rgb), 0.05) 100%
      );
      opacity: 0;
      transition: opacity 0.3s ease;
      pointer-events: none;
    }

    &:hover {
      border-color: rgba(var(--color-brand-500-rgb), 0.25);
      transform: translateY(-2px);
      box-shadow:
        0 8px 25px rgba(var(--color-background-900-rgb), 0.15),
        0 0 20px rgba(var(--color-brand-500-rgb), 0.1);

      &::before {
        opacity: 1;
      }
    }
  }

  .stats-container {
    display: flex;
    align-items: center;
    gap: 1rem;
    position: relative;
    z-index: 1;
  }

  .stats-icon {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 3rem;
    height: 3rem;
    background: linear-gradient(135deg, rgba(var(--color-brand-500-rgb), 0.15) 0%, rgba(var(--color-purple-500-rgb), 0.1) 100%);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.2);
    border-radius: var(--radius-sm);
    position: relative;

    i {
      font-size: 1.25rem;
      color: var(--color-brand-500);
      filter: drop-shadow(0 0 8px rgba(var(--color-brand-500-rgb), 0.3));
    }

    &::after {
      content: '';
      position: absolute;
      inset: -1px;
      background: linear-gradient(45deg, var(--color-brand-500), var(--color-purple-500));
      border-radius: var(--radius-sm);
      opacity: 0;
      z-index: -1;
      transition: opacity 0.3s ease;
    }
  }

  .stats-card:hover .stats-icon::after {
    opacity: 0.1;
  }

  .stats-content {
    flex: 1;
    min-width: 0;
  }

  .stats-value {
    font-size: 1.75rem;
    font-weight: 700;
    color: var(--color-brand-500);
    line-height: 1.2;
    text-shadow: 0 0 10px rgba(var(--color-brand-500-rgb), 0.3);
    margin-bottom: 0.25rem;

    background: linear-gradient(135deg, var(--color-brand-500) 0%, rgba(var(--color-purple-500-rgb), 0.8) 100%);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    background-clip: text;
  }

  .stats-label {
    font-size: 0.875rem;
    font-weight: 500;
    color: var(--color-content-muted);
    line-height: 1.3;
  }

  .stats-card[data-loading='true'] {
    .stats-value {
      background: linear-gradient(
        90deg,
        rgba(var(--color-brand-500-rgb), 0.3) 0%,
        rgba(var(--color-brand-500-rgb), 0.6) 50%,
        rgba(var(--color-brand-500-rgb), 0.3) 100%
      );
      background-size: 200% 100%;
      animation: shimmer 2s infinite;
      -webkit-background-clip: text;
      -webkit-text-fill-color: transparent;
      background-clip: text;
    }
  }
  @media (max-width: 768px) {
    .stats-card {
      padding: 1rem;
    }

    .stats-container {
      gap: 0.75rem;
    }

    .stats-icon {
      width: 2.5rem;
      height: 2.5rem;

      i {
        font-size: 1.125rem;
      }
    }

    .stats-value {
      font-size: 1.5rem;
    }

    .stats-label {
      font-size: 0.8rem;
    }
  }
</style>
