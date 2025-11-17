<template>
  <div v-if="qdrantStats" class="qdrant-stats-card">
    <div class="card-header">
      <div
        class="health-indicator"
        :class="{
          healthy: qdrantStats.is_healthy,
          unhealthy: !qdrantStats.is_healthy,
        }"
      >
        {{
          qdrantStats.is_healthy ? $t('admin.vectors.stats.systemStatus.healthy') : $t('admin.vectors.stats.systemStatus.error')
        }}
      </div>
    </div>

    <div class="stats-grid">
      <div class="stat-item">
        <span class="label">{{ $t('admin.vectors.qdrant.vectors') }}</span>
        <span class="value primary">{{ qdrantStats.qdrant_vector_count }}</span>
      </div>
      <div class="stat-item">
        <span class="label">{{ $t('admin.vectors.stats.indexStatus.title') }}</span>
        <span class="value">{{ qdrantStats.qdrant_indexed_count }}</span>
      </div>
      <div class="stat-item">
        <span class="label">{{ $t('admin.vectors.qdrant.collections') }}</span>
        <span class="value">{{ qdrantStats.mysql_total_count }}</span>
      </div>
      <div class="stat-item">
        <span class="label">{{ $t('admin.vectors.stats.indexStatus.completed') }}</span>
        <span class="value">{{ qdrantStats.mysql_completed_count }}</span>
      </div>
      <div class="stat-item">
        <span class="label">{{ $t('admin.vectors.qdrant.throughput') }}</span>
        <span
          class="value"
          :class="{
            success: qdrantStats.sync_ratio >= 95,
            warning: qdrantStats.sync_ratio >= 70 && qdrantStats.sync_ratio < 95,
            error: qdrantStats.sync_ratio < 70,
          }"
        >
          {{ qdrantStats.sync_ratio.toFixed(1) }}%
        </span>
      </div>
      <div class="stat-item">
        <span class="label">{{ $t('admin.vectors.columns.collection') }}</span>
        <span class="value status" :class="qdrantStats.collection_status">
          {{ getCollectionStatusText(qdrantStats.collection_status) }}
        </span>
      </div>
    </div>

    <div class="card-footer">
      <span class="last-checked">{{ $t('admin.vectors.qdrant.uptime') }}: {{ qdrantStats.last_checked }}</span>
    </div>
  </div>
</template>

<script setup lang="ts">
  import { useTexts } from '@/composables/useTexts'
  import type { QdrantRealStats } from '@/api/admin/vectors'
  import { getCollectionStatusText } from '../utils'

  interface Props {
    qdrantStats: QdrantRealStats | null
  }

  defineOptions({
    name: 'QdrantStatsCard',
  })

  const { $t } = useTexts()
  defineProps<Props>()
</script>

<style scoped lang="scss">
  .qdrant-stats-card {
    background: linear-gradient(
      135deg,
      rgba(var(--color-background-800-rgb), 0.8) 0%,
      rgba(var(--color-background-700-rgb), 0.6) 100%
    );
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.2);
    border-radius: var(--radius-sm);
    padding: var(--space-xl);
    backdrop-filter: var(--backdrop-blur-md);
    margin-top: var(--space-md);
    box-shadow: var(--shadow-cyber-md);
  }

  .card-header {
    display: flex;
    justify-content: flex-end;
    align-items: center;
    margin-bottom: var(--space-lg);
  }

  .health-indicator {
    padding: var(--space-xs) var(--space-sm);
    border-radius: var(--radius-sm);
    font-size: var(--text-sm);
    font-weight: var(--font-medium);
    text-transform: uppercase;
    letter-spacing: 0.5px;
  }

  .health-indicator.healthy {
    background: rgba(34, 197, 94, 0.2);
    color: var(--color-success);
    border: 1px solid rgba(34, 197, 94, 0.3);
  }

  .health-indicator.unhealthy {
    background: rgba(239, 68, 68, 0.2);
    color: var(--color-danger);
    border: 1px solid rgba(239, 68, 68, 0.3);
  }

  .stats-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
    gap: var(--space-sm);
    margin-bottom: var(--space-lg);
  }

  .stat-item {
    display: flex;
    flex-direction: column;
    gap: var(--space-sm);
  }

  .stat-item .label {
    font-size: var(--text-sm);
    color: var(--color-slate-400);
    font-weight: var(--font-medium);
  }

  .stat-item .value {
    font-size: var(--text-xl);
    font-weight: var(--font-semibold);
    color: var(--color-white);
  }

  .stat-item .value.primary {
    color: var(--color-brand-500);
  }

  .stat-item .value.success {
    color: var(--color-success);
  }

  .stat-item .value.warning {
    color: var(--color-warning);
  }

  .stat-item .value.error {
    color: var(--color-danger);
  }

  .stat-item .value.status {
    font-size: var(--text-base);
    padding: var(--space-xs) var(--space-sm);
    border-radius: var(--radius-sm);
    text-align: center;
  }

  .stat-item .value.status.healthy {
    background: rgba(34, 197, 94, 0.2);
    color: var(--color-success);
    border: 1px solid rgba(34, 197, 94, 0.3);
  }

  .stat-item .value.status.empty_normal {
    background: rgba(59, 130, 246, 0.2);
    color: var(--color-cyan-500);
    border: 1px solid rgba(59, 130, 246, 0.3);
  }

  .stat-item .value.status.empty_but_should_have_data {
    background: rgba(239, 68, 68, 0.2);
    color: var(--color-danger);
    border: 1px solid rgba(239, 68, 68, 0.3);
  }

  .stat-item .value.status.partially_synced {
    background: rgba(245, 158, 11, 0.2);
    color: var(--color-warning);
    border: 1px solid rgba(245, 158, 11, 0.3);
  }

  .stat-item .value.status.needs_resync {
    background: rgba(239, 68, 68, 0.2);
    color: var(--color-danger);
    border: 1px solid rgba(239, 68, 68, 0.3);
  }

  .stat-item .value.status.engine_not_initialized,
  .stat-item .value.status.qdrant_connection_failed {
    background: rgba(239, 68, 68, 0.2);
    color: var(--color-danger);
    border: 1px solid rgba(239, 68, 68, 0.3);
  }

  .stat-item .value.status.error,
  .stat-item .value.status.unknown {
    background: rgba(107, 114, 128, 0.2);
    color: var(--color-text-quaternary);
    border: 1px solid rgba(107, 114, 128, 0.3);
  }

  .card-footer {
    display: flex;
    justify-content: flex-end;
    padding-top: var(--space-md);
    border-top: 1px solid rgba(var(--color-brand-500-rgb), 0.1);
  }

  .last-checked {
    font-size: var(--text-xs);
    color: var(--color-content-muted);
    font-style: italic;
  }

  @media (max-width: 768px) {
    .stats-grid {
      grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
      gap: var(--space-sm);
    }

    .stat-item .value {
      font-size: var(--text-lg);
    }

    .card-header {
      justify-content: flex-end;
    }
  }
</style>
