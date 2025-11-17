<template>
  <div v-if="loading" class="cyber-skeleton-container" :class="[`type-${type}`, `size-${size}`, `animation-${animation}`]">
    <template v-if="type === 'card'">
      <div v-for="i in actualCount" :key="i" class="skeleton-card">
        <div class="skeleton skeleton-image" />
        <div class="skeleton-info">
          <div class="skeleton skeleton-text title" />
          <div v-if="!simple" class="skeleton skeleton-text subtitle" />
          <div v-if="!simple" class="skeleton-tags">
            <div class="skeleton skeleton-tag" />
            <div class="skeleton skeleton-tag" />
            <div class="skeleton skeleton-tag small" />
          </div>
          <div v-if="!simple" class="skeleton-actions">
            <div class="skeleton skeleton-text date" />
            <div class="skeleton-buttons">
              <div class="skeleton skeleton-btn" />
              <div class="skeleton skeleton-btn" />
              <div class="skeleton skeleton-btn" />
            </div>
          </div>
        </div>
      </div>
    </template>

    <template v-if="type === 'table'">
      <div class="skeleton-table">
        <div v-for="i in actualCount" :key="i" class="skeleton-table-row">
          <div class="skeleton skeleton-cell" />
          <div class="skeleton skeleton-cell" />
          <div class="skeleton skeleton-cell" />
          <div class="skeleton skeleton-cell" />
          <div class="skeleton skeleton-cell small" />
        </div>
      </div>
    </template>

    <template v-if="type === 'list'">
      <div v-for="i in actualCount" :key="i" class="skeleton-list-item">
        <div class="skeleton skeleton-avatar" />
        <div class="skeleton-content">
          <div class="skeleton skeleton-text title" />
          <div class="skeleton skeleton-text subtitle" />
        </div>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
  import { computed } from 'vue'
  import type { CyberSkeletonProps } from './types'

  defineOptions({
    name: 'CyberSkeleton',
  })

  const props = withDefaults(defineProps<CyberSkeletonProps>(), {
    size: 'normal',
    loading: true,
    animation: 'wave',
    simple: false,
  })

  const actualCount = computed(() => {
    if (props.count) return props.count

    switch (props.type) {
      case 'card':
        return 6
      case 'table':
        return 5
      case 'list':
        return 3
      default:
        return 3
    }
  })
</script>

<style scoped>
  .cyber-skeleton-container {
    @apply w-full;
  }

  .skeleton-image,
  .skeleton-text,
  .skeleton-tag,
  .skeleton-btn,
  .skeleton-cell,
  .skeleton-avatar {
    @apply rounded;
  }

  .type-card {
    display: contents;
  }

  .skeleton-card {
    @apply flex flex-col overflow-hidden rounded-lg border border-subtle bg-background-800;
  }

  .size-small .skeleton-card {
    @apply h-[180px];
  }
  .size-normal .skeleton-card {
    @apply h-[220px];
  }
  .size-large .skeleton-card {
    @apply h-[260px];
  }

  .skeleton-image {
    @apply bg-skeleton-highlight;
  }

  .size-small .skeleton-image {
    @apply h-[100px];
  }
  .size-normal .skeleton-image {
    @apply h-[140px];
  }
  .size-large .skeleton-image {
    @apply h-[180px];
  }

  .skeleton-info {
    @apply flex flex-1 flex-col gap-2 p-3;
  }

  .skeleton-text.title {
    @apply h-4 w-[85%];
  }

  .skeleton-text.subtitle {
    @apply h-3 w-[70%];
  }

  .skeleton-text.date {
    @apply h-3 w-1/2;
  }

  .skeleton-tags {
    @apply flex gap-1.5;
  }

  .skeleton-tag {
    @apply h-5 w-10;
  }

  .skeleton-tag.small {
    @apply w-6;
  }

  .skeleton-actions {
    @apply mt-auto flex items-center justify-between;
  }

  .skeleton-buttons {
    @apply flex gap-1;
  }

  .skeleton-btn {
    @apply h-7 w-7 rounded-md;
  }

  .skeleton-table {
    @apply flex w-full flex-col gap-2;
  }

  .skeleton-table-row {
    @apply flex gap-4 rounded-lg border border-subtle bg-background-800 p-4;
  }

  .skeleton-cell {
    @apply h-4 flex-1;
  }

  .skeleton-cell.small {
    @apply flex-[0.5];
  }

  .skeleton-list-item {
    @apply mb-2 flex gap-4 rounded-lg border border-subtle bg-background-800 p-4;
  }

  .skeleton-avatar {
    @apply h-12 w-12 shrink-0 rounded-full;
  }

  .skeleton-content {
    @apply flex flex-1 flex-col gap-2;
  }

  .animation-pulse .skeleton-image,
  .animation-pulse .skeleton-text,
  .animation-pulse .skeleton-tag,
  .animation-pulse .skeleton-btn,
  .animation-pulse .skeleton-cell,
  .animation-pulse .skeleton-avatar {
    animation: pulse 1.5s ease-in-out infinite;
  }

  @media (max-width: 768px) {
    .type-card {
      grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
      gap: 0.75rem;
    }

    .size-small .skeleton-card {
      @apply h-[160px];
    }
    .size-normal .skeleton-card {
      @apply h-[200px];
    }
    .size-large .skeleton-card {
      @apply h-[240px];
    }

    .size-small .skeleton-image {
      @apply h-[80px];
    }
    .size-normal .skeleton-image {
      @apply h-[120px];
    }
    .size-large .skeleton-image {
      @apply h-[160px];
    }

    .skeleton-info {
      @apply gap-1.5 p-2;
    }
  }
</style>
