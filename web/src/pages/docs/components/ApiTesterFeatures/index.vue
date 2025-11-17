<template>
  <div class="tester-features">
    <div v-for="feature in features" :key="feature.id" class="feature-item">
      <div class="feature-icon">
        <i :class="feature.icon" />
      </div>
      <div class="feature-content">
        <h5>{{ feature.title }}</h5>
        <p>{{ feature.description }}</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
  import type { FeatureItem } from '../../types'

  defineOptions({
    name: 'ApiTesterFeatures',
  })

  interface Props {
    features: FeatureItem[]
  }

  withDefaults(defineProps<Props>(), {
    features: () => [],
  })
</script>

<style scoped>
  .tester-features {
    @apply mb-8 grid grid-cols-2 gap-4 rounded-xl border p-6;
    background: linear-gradient(
      135deg,
      rgba(var(--color-brand-500-rgb), 0.05) 0%,
      rgba(var(--color-badge-accent-text-rgb), 0.03) 50%,
      rgba(var(--color-success-rgb), 0.05) 100%
    );
    border-color: rgba(var(--color-brand-500-rgb), 0.2);
  }

  .feature-item {
    @apply flex items-center gap-3 rounded-lg border p-4 transition-all duration-300;
    background: rgba(var(--color-background-700-rgb), 0.6);
    border-color: rgba(var(--color-brand-500-rgb), 0.4);
  }

  .feature-item:hover {
    @apply -translate-y-0.5;
    border-color: rgba(var(--color-brand-500-rgb), 0.4);
  }

  .feature-icon {
    @apply flex h-12 w-12 shrink-0 items-center justify-center rounded-lg border text-xl;
    background: linear-gradient(
      135deg,
      rgba(var(--color-brand-500-rgb), 0.2) 0%,
      rgba(var(--color-badge-accent-text-rgb), 0.15) 50%,
      rgba(var(--color-success-rgb), 0.2) 100%
    );
    border-color: rgba(var(--color-brand-500-rgb), 0.4);
    color: var(--color-brand-500);
  }

  .feature-content h5 {
    @apply m-0 mb-1 text-sm font-bold;
    color: var(--color-brand-500);
  }

  .feature-content p {
    @apply m-0 text-xs leading-tight opacity-90;
    color: var(--color-gray-200);
  }

  @media (max-width: 768px) {
    .tester-features {
      @apply grid-cols-1 gap-3 p-4;
    }

    .feature-item {
      @apply gap-2 p-3;
    }

    .feature-icon {
      @apply h-10 w-10 text-base;
    }

    .feature-content h5 {
      @apply text-sm;
    }

    .feature-content p {
      @apply text-xs;
    }
  }
</style>
