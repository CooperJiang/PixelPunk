<script setup lang="ts">
  import { computed, onMounted, ref } from 'vue'
  import { type QueueStatsResponse, getQueueStats, getTaggingDiagnosis } from '@/api/admin/tagging'
  import { useWebSocketStore } from '@/store/websocket'
  import { useTexts } from '@/composables/useTexts'

  const { $t } = useTexts()

  const props = defineProps<{
    selectedStatus?: string
  }>()

  const emit = defineEmits<{
    (e: 'select', status: string): void
  }>()

  const webSocketStore = useWebSocketStore()

  const fallbackQueueStats = ref<QueueStatsResponse | null>(null)
  const lastUpdate = ref<Date>(new Date())
  const isLoading = ref<boolean>(false)
  const useWebSocketMode = ref<boolean>(true) // 是否使用WebSocket模式
  const diagnosis = ref<{
    queue_initialized: boolean
    service_running: boolean
    paused: boolean
    config_issues: string[]
    recommendations: string[]
    worker_count: number
    active_workers: number
  } | null>(null)

  const queueStats = computed(() => webSocketStore.queueStats || fallbackQueueStats.value)

  const ext = computed(() => {
    const data = queueStats.value as Record<string, unknown>
    if (!data || !data.queue_stats_ext) {
      return { queued: 0, processing: 0, delayed: 0, dlq: 0 }
    }
    return data.queue_stats_ext as { queued: number; processing: number; delayed: number; dlq: number }
  })

  const runtime = computed(() => {
    const data = queueStats.value as Record<string, unknown>
    return (data && data.runtime) || {}
  })

  const _utilization = computed(() => {
    const active =
      (runtime.value as Record<string, unknown>)?.active_workers ??
      (queueStats.value as Record<string, unknown>)?.config?.current_concurrency ??
      0
    const max =
      (runtime.value as Record<string, unknown>)?.configured_concurrency ??
      (queueStats.value as Record<string, unknown>)?.config?.max_concurrency ??
      0
    if (!max) return 0
    return Math.min(100, Math.round((active / max) * 100))
  })

  const lastUpdatedText = computed(() => {
    const updateTime = useWebSocketMode.value ? webSocketStore.lastUpdated : lastUpdate.value
    const now = new Date()
    const diff = now.getTime() - updateTime.getTime()
    const seconds = Math.floor(diff / 1000)

    if (seconds < 60) {
      return $t('admin.tagging.queue.secondsAgo').replace('{seconds}', String(seconds))
    }
    const minutes = Math.floor(seconds / 60)
    return $t('admin.tagging.queue.minutesAgo').replace('{minutes}', String(minutes))
  })

  const isConnected = computed(() => webSocketStore.isConnected)
  const connectionStatus = computed(() => webSocketStore.connectionStatus)

  const connectionStatusText = computed(() => {
    if (useWebSocketMode.value) {
      return webSocketStore.connectionStatusText
    }
    return $t('admin.tagging.queue.pollingMode')
  })

  const _completionRate = computed(() => {
    if (useWebSocketMode.value && webSocketStore.queueStats) {
      return webSocketStore.queueCompletionRate
    }

    if (!queueStats.value || queueStats.value.queue_stats.total_count === 0) {
      return 0
    }
    return Math.round((queueStats.value.queue_stats.done_count / queueStats.value.queue_stats.total_count) * 100)
  })

  const hasQueueActivity = computed(() => {
    if (!queueStats.value) {
      return false
    }
    const stats = queueStats.value.queue_stats
    const { config } = queueStats.value

    return config.current_concurrency > 0 || stats.pending_count > 0 || stats.none_count > 0
  })

  const paused = computed(() => {
    if (diagnosis.value) {
      return diagnosis.value.paused
    }
    if (useWebSocketMode.value && webSocketStore.queueStats) {
      return !!webSocketStore.queueStats?.config?.paused
    }
    return !!queueStats.value?.config?.paused
  })

  const queueInitialized = computed(() => {
    return diagnosis.value?.queue_initialized ?? true
  })

  const getStatusText = () => {
    if (!queueInitialized.value) {
      return $t('admin.dashboard.aiServices.statusDown')
    }
    if (paused.value) {
      return $t('admin.dashboard.aiServices.statusDegraded')
    }
    return $t('admin.dashboard.aiServices.statusNormal')
  }

  const getStatusColor = () => {
    if (!queueInitialized.value) {
      return 'text-error-500'
    }
    if (paused.value) {
      return 'text-warning-500'
    }
    return 'text-success-500'
  }

  const _recentFailures = computed(() => {
    if (useWebSocketMode.value && webSocketStore.queueStats) {
      return webSocketStore.queueStats?.config?.recent_failures || 0
    }
    return queueStats.value?.config?.recent_failures || 0
  })

  const _queueLength = computed(() => {
    if (useWebSocketMode.value && webSocketStore.queueStats) {
      return webSocketStore.queueStats?.config?.queue_length || 0
    }
    return queueStats.value?.config?.queue_length || 0
  })

  const loadQueueStats = async () => {
    if (isLoading.value) {
      return
    }

    isLoading.value = true
    try {
      const result = await getQueueStats()

      if (result.success) {
        fallbackQueueStats.value = result.data
        if (!useWebSocketMode.value) {
          lastUpdate.value = new Date()
        }
      }
    } catch {
      fallbackQueueStats.value = null
    } finally {
      isLoading.value = false
    }
  }

  const loadDiagnosis = async () => {
    try {
      const result = await getTaggingDiagnosis()
      if (result.success) {
        diagnosis.value = result.data
      }
    } catch (_error) {}
  }

  onMounted(async () => {
    await Promise.all([loadQueueStats(), loadDiagnosis()])
  })

  const selectStatus = (status: string) => {
    const statusMap: Record<string, string> = {
      total: '',
      none: 'none',
      pending: 'pending',
      done: 'done',
      failed: 'failed',
      ignored: 'ignored',
    }

    const mappedStatus = statusMap[status] || status

    if (status === 'total') {
      emit('select', '')
    } else {
      emit('select', mappedStatus === props.selectedStatus ? '' : mappedStatus)
    }
  }

  const getStatusItemClass = (status: string): string => {
    const statusMap: Record<string, string> = {
      total: '',
      none: 'none',
      pending: 'pending',
      done: 'done',
      failed: 'failed',
      ignored: 'ignored',
    }

    const mappedStatus = statusMap[status]
    const isSelected = props.selectedStatus === mappedStatus
    const baseClass = 'status-item'

    return isSelected ? `${baseClass} status-item-selected` : baseClass
  }

  defineExpose({
    refresh: loadQueueStats,
    hasActivity: hasQueueActivity,
    isConnected: webSocketStore.isConnected,
    connectionStatus: webSocketStore.connectionStatusText,
  })
</script>

<template>
  <div class="ai-queue-status-card">
    <div v-if="diagnosis && diagnosis.config_issues.length > 0" class="diagnosis-banner">
      <div class="diagnosis-header">
        <div class="diagnosis-icon">
          <i class="fas fa-exclamation-triangle"></i>
        </div>
        <div class="diagnosis-title">{{ $t('admin.tagging.queueDiagnosis.title') }}</div>
      </div>

      <div class="diagnosis-content">
        <div class="diagnosis-section">
          <div class="section-label">
            <i class="fas fa-bug mr-1.5"></i>
            {{ $t('admin.tagging.queueDiagnosis.issuesFound') }}
          </div>
          <ul class="issue-list">
            <li v-for="issue in diagnosis.config_issues" :key="issue">
              <i class="fas fa-circle mr-2 text-xs"></i>
              <span>{{ issue }}</span>
            </li>
          </ul>
        </div>

        <div class="diagnosis-section">
          <div class="section-label">
            <i class="fas fa-lightbulb mr-1.5"></i>
            {{ $t('admin.tagging.queueDiagnosis.recommendations') }}
          </div>
          <ul class="recommendation-list">
            <li v-for="rec in diagnosis.recommendations" :key="rec">
              <i class="fas fa-arrow-right mr-2 text-xs"></i>
              <span>{{ rec }}</span>
            </li>
          </ul>
        </div>
      </div>

      <div class="diagnosis-action-hint">
        <i class="fas fa-info-circle mr-2"></i>
        <span>{{ $t('admin.tagging.queueDiagnosis.hint') }}</span>
      </div>
    </div>

    <div v-if="diagnosis && !diagnosis.queue_initialized" class="warning-banner">
      <i class="fas fa-exclamation-circle mr-2"></i>
      <span>{{ $t('admin.dashboard.aiServices.statusDown') }}</span>
    </div>

    <div class="card-info-bar">
      <div class="info-items">
        <span class="info-item" :class="getStatusColor()">
          {{ getStatusText() }}
        </span>
        <div class="info-item connection-status">
          <i
            :class="{
              'fas fa-wifi': isConnected && useWebSocketMode,
              'fas fa-wifi text-green-500': isConnected && useWebSocketMode,
              'fas fa-exclamation-triangle text-yellow-500': connectionStatus === 'reconnecting',
              'fas fa-times-circle text-red-500': connectionStatus === 'error' || connectionStatus === 'disconnected',
              'fas fa-clock': !useWebSocketMode,
            }"
          />
          <span class="status-text">{{ connectionStatusText }}</span>
        </div>
        <div class="info-item last-update">
          <i class="fas fa-clock mr-1" />
          <span class="update-time">{{ lastUpdatedText }}</span>
        </div>
      </div>
    </div>

    <div class="status-grid">
      <div v-if="paused" class="paused-banner">
        <i class="fas fa-pause-circle mr-2"></i>
        {{ $t('admin.tagging.hints.queuePaused') }}
      </div>
      <div :class="getStatusItemClass('total')" @click="selectStatus('total')">
        <div class="status-icon icon-all">
          <i class="fas fa-list" />
        </div>
        <div class="status-content">
          <div class="status-label">{{ $t('admin.tagging.status.all') }}</div>
          <div class="status-count">{{ queueStats?.queue_stats.total_count || 0 }}</div>
        </div>
      </div>

      <div :class="getStatusItemClass('none')" @click="selectStatus('none')">
        <div class="status-icon icon-pending">
          <i class="fas fa-clock" />
        </div>
        <div class="status-content">
          <div class="status-label">{{ $t('admin.tagging.status.none') }}</div>
          <div class="status-count">{{ queueStats?.queue_stats.none_count || 0 }}</div>
        </div>
      </div>

      <div :class="getStatusItemClass('pending')" @click="selectStatus('pending')">
        <div class="status-icon icon-running">
          <i class="fas fa-spinner" />
        </div>
        <div class="status-content">
          <div class="status-label">{{ $t('admin.tagging.status.pending') }}</div>
          <div class="status-count">{{ queueStats?.queue_stats.pending_count || 0 }}</div>
        </div>
      </div>

      <div :class="getStatusItemClass('done')" @click="selectStatus('done')">
        <div class="status-icon icon-completed">
          <i class="fas fa-check-circle" />
        </div>
        <div class="status-content">
          <div class="status-label">{{ $t('admin.tagging.status.done') }}</div>
          <div class="status-count">{{ queueStats?.queue_stats.done_count || 0 }}</div>
        </div>
      </div>

      <div :class="getStatusItemClass('failed')" @click="selectStatus('failed')">
        <div class="status-icon icon-failed">
          <i class="fas fa-exclamation-circle" />
        </div>
        <div class="status-content">
          <div class="status-label">{{ $t('admin.tagging.status.failed') }}</div>
          <div class="status-count">{{ queueStats?.queue_stats.failed_count || 0 }}</div>
        </div>
      </div>

      <div :class="getStatusItemClass('ignored')" @click="selectStatus('ignored')">
        <div class="status-icon icon-ignored">
          <i class="fas fa-ban" />
        </div>
        <div class="status-content">
          <div class="status-label">{{ $t('admin.tagging.status.ignored') }}</div>
          <div class="status-count">{{ queueStats?.queue_stats.ignored_count || 0 }}</div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped lang="scss">
  .ai-queue-status-card {
    background: linear-gradient(
      135deg,
      rgba(var(--color-background-800-rgb), 0.8) 0%,
      rgba(var(--color-background-700-rgb), 0.6) 100%
    );
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.2);
    border-radius: var(--radius-sm);
    padding: var(--space-xl);
    backdrop-filter: var(--backdrop-blur-md);
    box-shadow: var(--shadow-cyber-md);

    .card-info-bar {
      display: flex;
      justify-content: flex-end;
      align-items: center;
      margin-bottom: var(--space-lg);
      padding: var(--space-sm) var(--space-md);
      background: rgba(var(--color-background-900-rgb), 0.3);
      border-radius: var(--radius-sm);
      border: 1px solid rgba(var(--color-brand-500-rgb), 0.1);

      .info-items {
        display: flex;
        align-items: center;
        gap: var(--space-lg);
      }

      .info-item {
        display: flex;
        align-items: center;
        font-size: var(--text-xs);
        color: var(--color-content-muted);
        white-space: nowrap;

        &.connection-status {
          gap: var(--space-xs);

          i {
            font-size: var(--text-sm);
          }
        }

        &.last-update {
          i {
            font-size: var(--text-xs);
          }
        }
      }
    }

    .diagnosis-banner {
      margin-bottom: var(--space-lg);
      background: linear-gradient(135deg, rgba(var(--color-warning-rgb), 0.1) 0%, rgba(var(--color-warning-rgb), 0.05) 100%);
      border: 1px solid rgba(var(--color-warning-rgb), 0.3);
      border-radius: var(--radius-sm);
      overflow: hidden;
      box-shadow: var(--shadow-md);

      .diagnosis-header {
        display: flex;
        align-items: center;
        padding: var(--space-md) var(--space-lg);
        background: rgba(var(--color-warning-rgb), 0.15);
        border-bottom: 1px solid rgba(var(--color-warning-rgb), 0.2);

        .diagnosis-icon {
          display: flex;
          align-items: center;
          justify-content: center;
          width: 32px;
          height: 32px;
          background: rgba(var(--color-warning-rgb), 0.2);
          border-radius: var(--radius-md);
          color: var(--color-warning-500);
          font-size: var(--text-lg);
          margin-right: var(--space-md);
        }

        .diagnosis-title {
          font-size: var(--text-base);
          font-weight: var(--font-semibold);
          color: var(--color-warning-500);
        }
      }

      .diagnosis-content {
        padding: var(--space-lg);
        display: flex;
        flex-direction: column;
        gap: var(--space-lg);

        .diagnosis-section {
          .section-label {
            display: flex;
            align-items: center;
            font-size: var(--text-sm);
            font-weight: var(--font-semibold);
            color: var(--color-content-heading);
            margin-bottom: var(--space-sm);

            i {
              color: var(--color-warning-500);
            }
          }

          .issue-list,
          .recommendation-list {
            list-style: none;
            padding: 0;
            margin: 0;
            display: flex;
            flex-direction: column;
            gap: var(--space-xs);

            li {
              display: flex;
              align-items: flex-start;
              font-size: var(--text-sm);
              color: var(--color-content-default);
              line-height: 1.6;
              padding: var(--space-xs) var(--space-sm);
              background: rgba(var(--color-background-700-rgb), 0.5);
              border-radius: var(--radius-md);
              transition: all var(--transition-fast) var(--ease-out);

              &:hover {
                background: rgba(var(--color-background-700-rgb), 0.8);
              }

              i {
                margin-top: 4px;
                flex-shrink: 0;
              }

              span {
                flex: 1;
              }
            }
          }

          .issue-list li i {
            color: var(--color-warning-500);
          }

          .recommendation-list li i {
            color: var(--color-brand-500);
          }
        }
      }

      .diagnosis-action-hint {
        display: flex;
        align-items: center;
        padding: var(--space-md) var(--space-lg);
        background: rgba(var(--color-brand-500-rgb), 0.08);
        border-top: 1px solid rgba(var(--color-brand-500-rgb), 0.15);
        color: var(--color-brand-500);
        font-size: var(--text-sm);

        i {
          flex-shrink: 0;
        }
      }
    }

    .warning-banner {
      display: flex;
      align-items: center;
      padding: var(--space-md);
      margin-bottom: var(--space-lg);
      border: 1px solid rgba(var(--color-error-rgb), 0.5);
      background: rgba(var(--color-error-rgb), 0.1);
      color: var(--color-error-500);
      border-radius: var(--radius-sm);
      font-size: var(--text-sm);

      i {
        color: var(--color-error-500);
      }
    }

    .status-grid {
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(160px, 1fr));
      gap: var(--space-sm);
      margin-bottom: var(--space-lg);

      .paused-banner {
        grid-column: 1 / -1;
        display: flex;
        align-items: center;
        padding: var(--space-sm) var(--space-sm);
        border: 1px dashed rgba(var(--color-error-rgb), 0.5);
        background: rgba(var(--color-error-rgb), 0.06);
        color: var(--color-error-500);
        border-radius: var(--radius-sm);
        font-size: var(--text-sm);
      }
    }

    .status-item {
      display: flex;
      align-items: center;
      padding: var(--space-sm) var(--space-md);
      background: rgba(var(--color-background-800-rgb), 0.4);
      border: 1px solid var(--color-border-subtle);
      border-radius: var(--radius-sm);
      transition: all var(--transition-normal) var(--ease-out);
      cursor: pointer;

      &:hover {
        background: rgba(var(--color-brand-500-rgb), 0.1);
        border-color: rgba(var(--color-brand-500-rgb), 0.3);
        transform: translateY(-2px);
      }

      &.status-item-selected {
        background: rgba(var(--color-error-rgb), 0.15);
        border-color: rgba(var(--color-error-rgb), 0.5);
        box-shadow: var(--shadow-md);

        &:hover {
          background: rgba(var(--color-error-rgb), 0.2);
          border-color: rgba(var(--color-error-rgb), 0.6);
        }
      }
    }

    .status-icon {
      width: var(--space-3xl);
      height: var(--space-3xl);
      border-radius: var(--radius-md);
      display: flex;
      align-items: center;
      justify-content: center;
      margin-right: var(--space-sm);
      font-size: var(--text-sm);

      &.icon-all {
        background: rgba(var(--color-brand-500-rgb), 0.15);
        color: var(--color-brand-500);
      }

      &.icon-pending {
        background: rgba(var(--color-warning-rgb), 0.15);
        color: var(--color-warning-400);
      }

      &.icon-running {
        background: rgba(var(--color-error-rgb), 0.15);
        color: var(--color-error-500);

        i {
          @apply animate-spin;
        }
      }

      &.icon-completed {
        background: rgba(var(--color-success-rgb), 0.15);
        color: var(--color-success-400);
      }

      &.icon-failed {
        background: rgba(var(--color-error-rgb), 0.15);
        color: var(--color-error-400);
      }

      &.icon-ignored {
        background: rgba(var(--color-content-muted-rgb), 0.15);
        color: var(--color-content-muted);
      }
    }

    .status-content {
      flex: 1;

      .status-label {
        font-size: var(--text-sm);
        color: var(--color-content-muted);
        margin-bottom: var(--space-xs);
      }

      .status-count {
        font-size: var(--text-base);
        font-weight: var(--font-semibold);
        color: var(--color-content-default);
      }
    }

    .summary-info {
      display: flex;
      justify-content: space-between;
      padding-top: var(--space-md);
      border-top: 1px solid var(--color-border-subtle);

      .summary-item {
        display: flex;
        align-items: center;
        gap: var(--space-sm);

        .summary-label {
          font-size: var(--text-sm);
          color: var(--color-content-muted);
        }

        .summary-value {
          font-size: var(--text-sm);
          font-weight: var(--font-semibold);
          color: var(--color-brand-500);
        }
      }
    }
  }

  @media (max-width: 768px) {
    .ai-queue-status-card {
      .status-grid {
        grid-template-columns: repeat(auto-fit, minmax(140px, 1fr));
        gap: var(--space-sm);
      }

      .card-info-bar {
        .info-items {
          flex-wrap: wrap;
          gap: var(--space-sm);
        }
      }
    }
  }

  @media (max-width: 640px) {
    .ai-queue-status-card {
      .status-grid {
        grid-template-columns: repeat(2, 1fr);
      }
    }
  }
</style>
