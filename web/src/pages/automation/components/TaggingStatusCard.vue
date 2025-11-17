<script setup lang="ts">
  import { computed } from 'vue'
  import { useTexts } from '@/composables/useTexts'
  import type { TaggingTaskStats } from '@/api/user/automation'

  defineOptions({
    name: 'TaggingStatusCard',
  })

  const { $t } = useTexts()

  interface Props {
    stats: TaggingTaskStats | null
    systemEnabled: boolean
  }

  const props = defineProps<Props>()

  const progress = computed(() => {
    if (!props.stats || props.stats.total_count === 0) return 0
    return Math.round((props.stats.done_count / props.stats.total_count) * 100)
  })

  /* 状态项配置 */
  const statusItems = computed(() => {
    if (!props.stats) return []

    return [
      {
        label: $t('automation.tagging.status.none'),
        count: props.stats.none_count,
        icon: 'fas fa-circle',
        statusClass: 'status-none',
      },
      {
        label: $t('automation.tagging.status.pending'),
        count: props.stats.pending_count,
        icon: 'fas fa-clock',
        statusClass: 'status-pending',
      },
      {
        label: $t('automation.tagging.status.processing'),
        count: props.stats.processing_count,
        icon: 'fas fa-spinner fa-spin',
        statusClass: 'status-processing',
      },
      {
        label: $t('automation.tagging.status.done'),
        count: props.stats.done_count,
        icon: 'fas fa-check-circle',
        statusClass: 'status-done',
      },
      {
        label: $t('automation.tagging.status.failed'),
        count: props.stats.failed_count,
        icon: 'fas fa-times-circle',
        statusClass: 'status-failed',
      },
    ]
  })

  const queueInfo = computed(() => {
    if (!props.stats) return null
    const { queue_position, pending_count } = props.stats
    if (queue_position !== undefined && queue_position > 0) {
      return $t('automation.tagging.queue.position', { position: queue_position })
    }
    if (pending_count > 0) {
      return $t('automation.tagging.queue.pending', { count: pending_count })
    }
    return null
  })
</script>

<template>
  <div class="tagging-status-card">
    <div class="card-header">
      <div class="flex items-center gap-3">
        <div class="card-icon">
          <i class="fas fa-robot" />
        </div>
        <div>
          <h3 class="card-title">{{ $t('automation.tagging.title') }}</h3>
          <p class="card-subtitle">{{ $t('automation.tagging.subtitle') }}</p>
        </div>
      </div>

      <div class="flex items-center gap-2">
        <span class="system-status-badge" :class="systemEnabled ? 'status-enabled' : 'status-disabled'">
          <span class="status-dot" />
          {{ systemEnabled ? $t('automation.tagging.statusRunning') : $t('automation.tagging.statusPaused') }}
        </span>
      </div>
    </div>

    <div v-if="stats && stats.total_count > 0" class="progress-section">
      <div class="progress-header">
        <span class="progress-label">{{ $t('automation.tagging.progress.label') }}</span>
        <span class="progress-value">{{ progress }}%</span>
      </div>
      <div class="progress-bar-container">
        <div class="progress-bar-bg">
          <div class="progress-bar-fill" :style="{ width: progress + '%' }">
            <div class="progress-bar-shine" />
          </div>
        </div>
      </div>
      <div class="progress-footer">
        {{ $t('automation.tagging.progress.completed', { done: stats.done_count, total: stats.total_count }) }}
      </div>
    </div>

    <div v-if="queueInfo" class="queue-info">
      <i class="fas fa-info-circle" />
      {{ queueInfo }}
    </div>

    <div class="status-grid">
      <div v-for="item in statusItems" :key="item.label" class="status-item" :class="item.statusClass">
        <i :class="[item.icon, 'status-icon']" />
        <div class="status-content">
          <span class="status-count">{{ item.count }}</span>
          <span class="status-label">{{ item.label }}</span>
        </div>
      </div>
    </div>

    <div v-if="!stats || stats.total_count === 0" class="empty-state">
      <div class="empty-icon-wrapper">
        <i class="fas fa-inbox" />
      </div>
      <p class="empty-title">{{ $t('automation.tagging.empty.title') }}</p>
      <p class="empty-subtitle">{{ $t('automation.tagging.empty.subtitle') }}</p>
    </div>
  </div>
</template>

<style scoped lang="scss">
  .tagging-status-card {
    position: relative;
    padding: 1.25rem;
    background: linear-gradient(145deg, rgba(var(--color-background-800-rgb), 0.8), rgba(var(--color-background-800-rgb), 0.6));
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.2);
    border-radius: var(--radius-sm);
    box-shadow:
      0 2px 0 rgba(var(--color-brand-500-rgb), 0.15) inset,
      0 -1px 0 rgba(0, 0, 0, 0.15) inset,
      0 10px 20px rgba(0, 0, 0, 0.2),
      0 5px 10px rgba(var(--color-brand-500-rgb), 0.15),
      0 2px 5px rgba(0, 0, 0, 0.1);
    backdrop-filter: var(--backdrop-blur-md);
    overflow: hidden;
    transition: border-color 0.3s ease;
    transform: translateZ(0);

    &::before {
      content: '';
      position: absolute;
      top: 0;
      left: 0;
      right: 0;
      height: 4px;
      background: linear-gradient(90deg, transparent, var(--color-brand-500), transparent);
      opacity: 0.8;
      box-shadow: 0 2px 6px rgba(var(--color-brand-500-rgb), 0.6);
    }

    &::after {
      content: '';
      position: absolute;
      top: 0;
      left: 0;
      right: 0;
      bottom: 0;
      background: radial-gradient(circle at 10% 20%, rgba(var(--color-brand-500-rgb), 0.03) 0%, transparent 50%);
      pointer-events: none;
    }

    &:hover {
      border-color: rgba(var(--color-brand-500-rgb), 0.35);
    }
  }

  .card-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: var(--space-md);
  }

  .card-icon {
    position: relative;
    width: 2.25rem;
    height: 2.25rem;
    display: flex;
    align-items: center;
    justify-content: center;
    background: linear-gradient(135deg, rgba(var(--color-brand-500-rgb), 0.15), rgba(var(--color-brand-500-rgb), 0.05));
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.3);
    border-radius: var(--radius-sm);
    color: var(--color-brand-500);
    transition: all 0.3s ease;

    &::after {
      content: '';
      position: absolute;
      inset: -2px;
      background: radial-gradient(circle, rgba(var(--color-brand-500-rgb), 0.2) 0%, transparent 70%);
      border-radius: var(--radius-sm);
      opacity: 0;
      transition: opacity 0.3s ease;
    }

    .tagging-status-card:hover & {
      transform: scale(1.05);
      box-shadow: 0 0 15px rgba(var(--color-brand-500-rgb), 0.3);

      &::after {
        opacity: 1;
      }
    }
  }

  .card-title {
    font-size: var(--text-base);
    font-weight: 600;
    color: var(--color-content-heading);
  }

  .card-subtitle {
    font-size: var(--text-xs);
    color: var(--color-content-muted);
  }

  .system-status-badge {
    display: inline-flex;
    align-items: center;
    gap: 0.375rem;
    padding: 0.25rem 0.625rem;
    border-radius: var(--radius-md);
    border: 1.5px solid;
    font-size: var(--text-xs);
    font-weight: 500;
    transition: all var(--transition-normal) var(--ease-in-out);

    .status-dot {
      width: 0.375rem;
      height: 0.375rem;
      border-radius: var(--radius-full);
    }

    &.status-enabled {
      background: rgba(var(--color-success-rgb), 0.1);
      border-color: rgba(var(--color-success-rgb), 0.3);
      color: var(--color-success-500);

      .status-dot {
        background: var(--color-success-500);
        box-shadow: 0 0 8px rgba(var(--color-success-rgb), 0.5);
      }
    }

    &.status-disabled {
      background: rgba(var(--color-background-700-rgb), 0.5);
      border-color: var(--color-border-subtle);
      color: var(--color-content-muted);

      .status-dot {
        background: var(--color-content-muted);
      }
    }
  }

  .progress-section {
    margin-bottom: var(--space-md);
    padding: 1rem;
    background: linear-gradient(145deg, rgba(var(--color-background-700-rgb), 0.5), rgba(var(--color-background-700-rgb), 0.3));
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.1);
    border-radius: var(--radius-sm);
  }

  .progress-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 0.75rem;
    font-size: var(--text-xs);

    .progress-label {
      font-weight: 600;
      color: var(--color-content-default);
    }

    .progress-value {
      font-family: var(--font-mono);
      font-weight: 700;
      font-size: 1.125rem;
      color: var(--color-brand-500);
    }
  }

  .progress-bar-container {
    margin-bottom: 0.75rem;
  }

  .progress-bar-bg {
    position: relative;
    height: 12px;
    background: rgba(var(--color-content-rgb), 0.08);
    border-radius: var(--radius-sm);
    overflow: hidden;
    box-shadow:
      0 1px 0 rgba(0, 0, 0, 0.1) inset,
      0 -1px 0 rgba(255, 255, 255, 0.05) inset;
  }

  .progress-bar-fill {
    position: relative;
    height: 100%;
    background: linear-gradient(90deg, var(--color-brand-500), var(--color-brand-600));
    border-radius: var(--radius-sm);
    transition: width 0.6s cubic-bezier(0.4, 0, 0.2, 1);
    box-shadow:
      0 1px 0 rgba(255, 255, 255, 0.3) inset,
      0 2px 4px rgba(var(--color-brand-500-rgb), 0.4);
    overflow: hidden;
  }

  .progress-bar-shine {
    position: absolute;
    top: 0;
    left: -100%;
    width: 100%;
    height: 100%;
    background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.3), transparent);
    animation: progressShine 2s ease-in-out infinite;
  }

  @keyframes progressShine {
    0% {
      left: -100%;
    }
    100% {
      left: 200%;
    }
  }

  .progress-footer {
    text-align: center;
    font-size: var(--text-xs);
    color: var(--color-content-muted);
    font-weight: 500;
  }

  .queue-info {
    margin-bottom: var(--space-sm);
    padding: 0.375rem 0.75rem;
    background: rgba(var(--color-brand-500-rgb), 0.1);
    border: 1.5px solid rgba(var(--color-brand-500-rgb), 0.3);
    border-radius: var(--radius-sm);
    font-size: var(--text-xs);
    color: var(--color-brand-500);

    i {
      margin-right: 0.375rem;
    }
  }

  .status-grid {
    display: grid;
    grid-template-columns: repeat(5, 1fr);
    gap: 0.75rem;
  }

  .status-item {
    position: relative;
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 0.625rem;
    padding: 0.75rem;
    border-radius: var(--radius-sm);
    border: 1px solid;
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    overflow: hidden;

    &::before {
      content: '';
      position: absolute;
      top: 0;
      left: -100%;
      width: 100%;
      height: 100%;
      background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.1), transparent);
      transition: left 0.5s ease;
    }

    &:hover {
      transform: translateY(-3px) scale(1.02);
      box-shadow: 0 8px 20px rgba(0, 0, 0, 0.15);

      &::before {
        left: 100%;
      }
    }

    .status-icon {
      font-size: var(--text-lg);
    }

    .status-content {
      display: flex;
      flex-direction: column;
    }

    .status-count {
      font-size: 1.25rem;
      font-weight: 700;
      line-height: 1.2;
    }

    .status-label {
      font-size: var(--text-xs);
      color: var(--color-content-muted);
    }

    &.status-none {
      background: rgba(var(--color-content-rgb), 0.05);
      border-color: var(--color-border-subtle);

      .status-icon,
      .status-count {
        color: var(--color-content-muted);
      }
    }

    &.status-pending {
      background: linear-gradient(145deg, rgba(var(--color-brand-500-rgb), 0.15), rgba(var(--color-brand-500-rgb), 0.08));
      border-color: rgba(var(--color-brand-500-rgb), 0.3);
      box-shadow:
        0 1px 0 rgba(var(--color-brand-500-rgb), 0.2) inset,
        0 -1px 0 rgba(0, 0, 0, 0.1) inset,
        0 2px 4px rgba(var(--color-brand-500-rgb), 0.2);

      .status-icon,
      .status-count {
        color: var(--color-brand-500);
      }
    }

    &.status-processing {
      background: linear-gradient(145deg, rgba(var(--color-warning-rgb), 0.15), rgba(var(--color-warning-rgb), 0.08));
      border-color: rgba(var(--color-warning-rgb), 0.3);
      box-shadow:
        0 1px 0 rgba(var(--color-warning-rgb), 0.2) inset,
        0 -1px 0 rgba(0, 0, 0, 0.1) inset,
        0 2px 4px rgba(var(--color-warning-rgb), 0.2);

      .status-icon,
      .status-count {
        color: var(--color-warning-500);
      }
    }

    &.status-done {
      background: linear-gradient(145deg, rgba(var(--color-success-rgb), 0.15), rgba(var(--color-success-rgb), 0.08));
      border-color: rgba(var(--color-success-rgb), 0.3);
      box-shadow:
        0 1px 0 rgba(var(--color-success-rgb), 0.2) inset,
        0 -1px 0 rgba(0, 0, 0, 0.1) inset,
        0 2px 4px rgba(var(--color-success-rgb), 0.2);

      .status-icon,
      .status-count {
        color: var(--color-success-500);
      }
    }

    &.status-failed {
      background: linear-gradient(145deg, rgba(var(--color-error-rgb), 0.15), rgba(var(--color-error-rgb), 0.08));
      border-color: rgba(var(--color-error-rgb), 0.3);
      box-shadow:
        0 1px 0 rgba(var(--color-error-rgb), 0.2) inset,
        0 -1px 0 rgba(0, 0, 0, 0.1) inset,
        0 2px 4px rgba(var(--color-error-rgb), 0.2);

      .status-icon,
      .status-count {
        color: var(--color-error-500);
      }
    }
  }

  .empty-state {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: var(--space-xl) 0;
    text-align: center;
  }

  .empty-icon-wrapper {
    width: 4rem;
    height: 4rem;
    display: flex;
    align-items: center;
    justify-content: center;
    background: linear-gradient(135deg, rgba(var(--color-brand-500-rgb), 0.1), rgba(var(--color-brand-500-rgb), 0.05));
    border: 2px solid rgba(var(--color-brand-500-rgb), 0.2);
    border-radius: var(--radius-2xl);
    margin-bottom: var(--space-sm);

    i {
      font-size: var(--text-3xl);
      color: var(--text-muted);
      opacity: 0.3;
    }
  }

  .empty-title {
    font-size: var(--text-sm);
    color: var(--color-content-heading);
  }

  .empty-subtitle {
    margin-top: var(--space-xs);
    font-size: var(--text-xs);
    color: var(--color-content-muted);
  }
</style>
