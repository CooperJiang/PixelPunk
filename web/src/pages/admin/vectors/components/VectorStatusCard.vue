<template>
  <div class="vector-status-card" v-loading="props.isLoading">
    <div class="card-header">
      <div class="worker-info">
        <span class="text-content-content-muted text-xs">
          {{ $t('admin.vectors.qdrant.workerThreads') }}: {{ workerInfo.current }}/{{ workerInfo.max }}
        </span>
        <div class="connection-status">
          <i
            :class="{
              'fas fa-wifi': isConnected && useWebSocketMode,
              'fas fa-wifi text-green-500': isConnected && useWebSocketMode,
              'fas fa-exclamation-triangle text-yellow-500': connectionStatus === 'reconnecting',
              'fas fa-times-circle text-red-500': connectionStatus === 'error' || connectionStatus === 'disconnected',
              'fas fa-clock': !useWebSocketMode,
            }"
          />
          <span class="status-text text-xs">{{ connectionStatusText }}</span>
        </div>
        <div class="last-update">
          <i class="fas fa-clock mr-1" />
          <span class="update-time">{{ lastUpdatedText }}</span>
        </div>
      </div>
    </div>

    <div class="status-grid">
      <div
        v-for="status in statusList"
        :key="status.key"
        :class="getStatusItemClass(status.key)"
        @click="selectStatus(status.key)"
      >
        <div class="status-icon" :class="status.iconClass">
          <i :class="status.icon"></i>
        </div>
        <div class="status-content">
          <div class="status-label">{{ status.label }}</div>
          <div class="status-count">{{ getStatusCount(status.key) }}</div>
        </div>
      </div>
    </div>

    <div class="summary-info">
      <div class="summary-item">
        <span class="summary-label">{{ $t('admin.vectors.pagination.total') }}:</span>
        <span class="summary-value">{{ stats?.total_count || 0 }}</span>
      </div>
      <div class="summary-item">
        <span class="summary-label">{{ $t('admin.vectors.qdrant.throughput') }}:</span>
        <span class="summary-value">{{ completionRate }}%</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
  import { computed, onMounted, ref } from 'vue'
  import { useTexts } from '@/composables/useTexts'
  import { getVectorStats, type VectorStats } from '@/api/admin/vectors'
  import { useWebSocketStore } from '@/store/websocket'

  interface Props {
    selectedStatus?: string
    isLoading?: boolean
  }

  interface Emits {
    (e: 'select', status: string): void
  }

  const props = defineProps<Props>()
  const emit = defineEmits<Emits>()
  const { $t } = useTexts()

  const webSocketStore = useWebSocketStore()

  const fallbackVectorStats = ref<VectorStats | null>(null)
  const lastUpdate = ref<Date>(new Date())
  const statsLoading = ref<boolean>(false)
  const useWebSocketMode = ref<boolean>(true)

  const ws = computed(() => webSocketStore.vectorStats)
  const stats = computed(() => ws.value?.queue_stats || fallbackVectorStats.value)
  const _coverage = computed(() => ws.value?.coverage || null)
  const _divergence = computed(() => ws.value?.divergence || null)

  const lastUpdatedText = computed(() => {
    const updateTime = useWebSocketMode.value ? webSocketStore.lastUpdated : lastUpdate.value
    const now = new Date()
    const diff = now.getTime() - updateTime.getTime()
    const seconds = Math.floor(diff / 1000)

    if (seconds < 60) {
      return $t('admin.vectors.status.updatedSecondsAgo', { seconds })
    }
    const minutes = Math.floor(seconds / 60)
    return $t('admin.vectors.status.updatedMinutesAgo', { minutes })
  })

  const isConnected = computed(() => webSocketStore.isConnected)
  const connectionStatus = computed(() => webSocketStore.connectionStatus)

  const connectionStatusText = computed(() => {
    if (useWebSocketMode.value) {
      return webSocketStore.connectionStatusText
    }
    return $t('admin.vectors.status.scheduledUpdate')
  })

  const workerInfo = computed(() => {
    if (ws.value?.runtime) {
      return {
        current: ws.value.runtime.active_workers || 0,
        max: ws.value.runtime.configured_concurrency || 0,
      }
    }
    return { current: 0, max: 0 }
  })

  const loadVectorStats = async () => {
    if (statsLoading.value) {
      return
    }

    statsLoading.value = true
    try {
      const result = await getVectorStats()

      if (result && (result.success === true || result.code === 200)) {
        if (result.data && typeof result.data === 'object') {
          const data = result.data as Record<string, unknown>
          const q = (data.queue_stats ?? data) as Record<string, unknown>
          const rt = (data.runtime ?? {}) as Record<string, unknown>
          fallbackVectorStats.value = {
            total_count: q.total_count || 0,
            pending_count: q.pending_count || 0,
            running_count: (q.processing_count ?? q.running_count ?? 0) as number,
            completed_count: q.completed_count || 0,
            failed_count: q.failed_count || 0,
            reset_count: q.reset_count || 0,
            active_workers: rt.active_workers || data.active_workers || 0,
            max_workers: rt.configured_concurrency || data.max_workers || 0,
          }
          if (!useWebSocketMode.value) {
            lastUpdate.value = new Date()
          }
        }
      }
    } catch {
      fallbackVectorStats.value = null
    } finally {
      statsLoading.value = false
    }
  }

  onMounted(async () => {
    await loadVectorStats()
  })

  const completionRate = computed(() => {
    if (!stats.value || stats.value.total_count === 0) {
      return 0
    }
    return Math.round((stats.value.completed_count / stats.value.total_count) * 100)
  })

  const hasQueueActivity = computed(() => {
    if (!stats.value) {
      return false
    }

    const statsValue = stats.value as Record<string, unknown>
    const running = statsValue.running_count ?? statsValue.processing_count ?? 0
    return workerInfo.value.current > 0 || (stats.value.pending_count || 0) > 0 || running > 0
  })

  defineExpose({
    refresh: loadVectorStats,
    hasActivity: hasQueueActivity,
    isConnected: webSocketStore.isConnected,
    connectionStatus: webSocketStore.connectionStatusText,
  })

  const statusList = computed(() => [
    {
      key: 'all',
      label: $t('admin.vectors.filter.options.all'),
      icon: 'fas fa-list',
      iconClass: 'icon-all',
    },
    {
      key: 'pending',
      label: $t('admin.vectors.status.indexing'),
      icon: 'fas fa-clock',
      iconClass: 'icon-pending',
    },
    {
      key: 'running',
      label: $t('admin.vectors.status.processing'),
      icon: 'fas fa-spinner animate-spin',
      iconClass: 'icon-running',
    },
    {
      key: 'completed',
      label: $t('admin.vectors.status.active'),
      icon: 'fas fa-check-circle',
      iconClass: 'icon-completed',
    },
    {
      key: 'failed',
      label: $t('admin.vectors.status.error'),
      icon: 'fas fa-exclamation-circle',
      iconClass: 'icon-failed',
    },
  ])

  const getStatusCount = (status: string) => {
    if (!stats.value) return 0

    switch (status) {
      case 'all':
        return stats.value.total_count || 0
      case 'pending':
        return stats.value.pending_count || 0
      case 'running': {
        const statsValue = stats.value as Record<string, unknown>
        return statsValue.running_count ?? statsValue.processing_count ?? 0
      }
      case 'completed':
        return stats.value.completed_count || 0
      case 'failed':
        return stats.value.failed_count || 0
      default:
        return 0
    }
  }

  const selectStatus = (status: string) => {
    const statusMap: Record<string, string> = {
      all: '',
      pending: 'pending',
      running: 'processing',
      completed: 'completed',
      failed: 'failed',
    }

    const mappedStatus = statusMap[status] || status
    emit('select', mappedStatus === props.selectedStatus ? '' : mappedStatus)
  }

  const getStatusItemClass = (status: string): string => {
    const statusMap: Record<string, string> = {
      all: '',
      pending: 'pending',
      running: 'processing',
      completed: 'completed',
      failed: 'failed',
      reset: 'reset',
    }

    const mappedStatus = statusMap[status]
    const isSelected = props.selectedStatus === mappedStatus
    const baseClass = 'status-item'

    return isSelected ? `${baseClass} status-item-selected` : baseClass
  }

  const _handleSelect = (status: string) => {
    selectStatus(status)
  }
</script>

<style scoped lang="scss">
  .vector-status-card {
    position: relative;
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

    .card-header {
      display: flex;
      justify-content: flex-end;
      align-items: center;
      margin-bottom: var(--space-lg);

      .worker-info {
        display: flex;
        align-items: center;
        gap: var(--space-sm);

        .connection-status {
          display: flex;
          align-items: center;
          gap: var(--space-xs);
          font-size: var(--text-xs);
          color: var(--color-content-default-content-muted);

          .status-text {
            white-space: nowrap;
          }

          i {
            font-size: var(--text-sm);
          }
        }

        .last-update {
          display: flex;
          align-items: center;
          font-size: var(--text-xs);
          color: var(--color-content-default-content-muted);

          .update-time {
            margin-left: var(--space-xs);
          }
        }
      }
    }

    .status-grid {
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(160px, 1fr));
      gap: var(--space-sm);
      margin-bottom: var(--space-lg);
    }

    .status-item {
      display: flex;
      align-items: center;
      padding: var(--space-sm) var(--space-md);
      background: rgba(var(--color-background-800-rgb), 0.4);
      border: 1px solid rgba(255, 255, 255, 0.1);
      border-radius: var(--radius-sm);
      cursor: pointer;
      transition: all var(--transition-normal) var(--ease-out);

      &:hover {
        background: rgba(var(--color-brand-500-rgb), 0.1);
        border-color: rgba(var(--color-brand-500-rgb), 0.3);
        transform: translateY(-2px);
      }

      &.active {
        background: rgba(var(--color-brand-500-rgb), 0.2);
        border-color: var(--color-brand-500);
        box-shadow: var(--shadow-glow-md);
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
        background: rgba(108, 117, 125, 0.2);
        color: var(--color-content-muted);
      }

      &.icon-pending {
        background: rgba(255, 193, 7, 0.2);
        color: var(--color-warning);
      }

      &.icon-running {
        background: rgba(13, 202, 240, 0.2);
        color: var(--color-info-500);
      }

      &.icon-completed {
        background: rgba(25, 135, 84, 0.2);
        color: var(--color-success-500);
      }

      &.icon-failed {
        background: rgba(220, 53, 69, 0.2);
        color: var(--color-danger);
      }

      &.icon-reset {
        background: rgba(255, 107, 1, 0.2);
        color: var(--color-warning-500);
      }
    }

    .status-content {
      flex: 1;

      .status-label {
        font-size: var(--text-sm);
        color: var(--color-slate-400);
        margin-bottom: var(--space-xs);
      }

      .status-count {
        font-size: var(--text-base);
        font-weight: var(--font-semibold);
        color: var(--color-white);
      }
    }

    .summary-info {
      display: flex;
      justify-content: space-between;
      padding-top: var(--space-md);
      border-top: 1px solid rgba(255, 255, 255, 0.1);

      .summary-item {
        display: flex;
        align-items: center;
        gap: var(--space-sm);

        .summary-label {
          font-size: var(--text-sm);
          color: var(--color-slate-400);
        }

        .summary-value {
          font-size: var(--text-sm);
          font-weight: var(--font-semibold);
          color: var(--color-brand-500);
        }
      }
    }
  }
</style>
