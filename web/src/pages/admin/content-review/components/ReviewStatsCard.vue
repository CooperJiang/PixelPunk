<script setup lang="ts">
  import type { ReviewStats } from '@/api/admin/content-review'
  import { useTexts } from '@/composables/useTexts'

  defineOptions({
    name: 'ReviewStatsCard',
  })

  defineProps<Props>()

  const { $t } = useTexts()

  interface Props {
    stats?: ReviewStats | null
  }
</script>

<template>
  <div class="mb-6 grid grid-cols-1 gap-6 md:grid-cols-3">
    <div class="stats-card pending">
      <div class="stats-content">
        <div class="stats-icon">
          <i class="fas fa-clock" />
        </div>
        <div class="stats-info">
          <h3 class="stats-title">{{ $t('admin.contentReview.stats.pending.title') }}</h3>
          <p class="stats-value">{{ stats?.pending_count || 0 }}</p>
          <p class="stats-desc">{{ $t('admin.contentReview.stats.pending.description') }}</p>
        </div>
      </div>
    </div>

    <div class="stats-card approved">
      <div class="stats-content">
        <div class="stats-icon">
          <i class="fas fa-check-circle" />
        </div>
        <div class="stats-info">
          <h3 class="stats-title">{{ $t('admin.contentReview.stats.approved.title') }}</h3>
          <p class="stats-value">{{ stats?.approved_today || 0 }}</p>
          <p class="stats-desc">{{ $t('admin.contentReview.stats.approved.description') }}</p>
        </div>
      </div>
    </div>

    <div class="stats-card rejected">
      <div class="stats-content">
        <div class="stats-icon">
          <i class="fas fa-times-circle" />
        </div>
        <div class="stats-info">
          <h3 class="stats-title">{{ $t('admin.contentReview.stats.rejected.title') }}</h3>
          <p class="stats-value">{{ stats?.rejected_today || 0 }}</p>
          <p class="stats-desc">{{ $t('admin.contentReview.stats.rejected.description') }}</p>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped lang="scss">
  .stats-card {
    border-radius: var(--radius-sm);
    border: 1px solid var(--color-border-default);
    background: var(--color-background-600);
    padding: var(--space-xl);
    transition: all var(--transition-normal) var(--ease-out);

    &:hover {
      box-shadow: var(--shadow-md);
    }

    &.pending {
      border-color: rgba(var(--color-warning-rgb), 0.3);

      &:hover {
        border-color: rgba(var(--color-warning-rgb), 0.5);
      }
    }

    &.approved {
      border-color: rgba(var(--color-success-rgb), 0.3);

      &:hover {
        border-color: rgba(var(--color-success-rgb), 0.5);
      }
    }

    &.rejected {
      border-color: rgba(var(--color-error-rgb), 0.3);

      &:hover {
        border-color: rgba(var(--color-error-rgb), 0.5);
      }
    }
  }

  .stats-content {
    display: flex;
    align-items: flex-start;
    gap: var(--space-md);
  }

  .stats-icon {
    display: flex;
    height: var(--space-4xl);
    width: var(--space-4xl);
    flex-shrink: 0;
    align-items: center;
    justify-content: center;
    border-radius: var(--radius-sm);
    font-size: var(--text-xl);
  }

  .pending .stats-icon {
    background: rgba(var(--color-warning-rgb), 0.2);
    color: var(--color-warning-400);
  }

  .approved .stats-icon {
    background: rgba(var(--color-success-rgb), 0.2);
    color: var(--color-success-400);
  }

  .rejected .stats-icon {
    background: rgba(var(--color-error-rgb), 0.2);
    color: var(--color-error-400);
  }

  .stats-info {
    flex: 1;
  }

  .stats-title {
    margin-bottom: var(--space-xs);
    font-size: var(--text-sm);
    font-weight: var(--font-medium);
    color: rgba(255, 255, 255, 0.8);
  }

  .stats-value {
    margin-bottom: var(--space-xs);
    font-size: var(--text-2xl);
    font-weight: var(--font-bold);
    color: var(--color-white);
  }

  .stats-desc {
    font-size: var(--text-xs);
    color: var(--color-content-disabled);
  }
</style>
