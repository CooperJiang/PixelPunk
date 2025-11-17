<script setup lang="ts">
  import { computed, onMounted, ref } from 'vue'
  import { getTaggingStats, type TaggingStatsResponse } from '@/api/admin/tagging'
  import { getVectorStats, type VectorStats } from '@/api/admin/vectors'
  import { useTexts } from '@/composables/useTexts'
  const { $t } = useTexts()

  defineOptions({
    name: 'AIServicesCard',
  })

  const loading = ref(false)
  const taggingData = ref<TaggingStatsResponse | null>(null)
  const vectorData = ref<VectorStats | null>(null)

  const data = computed(() => {
    const tagging = {
      total_count: 0,
      none_count: 0,
      pending_count: 0,
      done_count: 0,
      failed_count: 0,
      skipped_count: 0,
      ignored_count: 0,
    }

    if (taggingData.value?.status_stats) {
      taggingData.value.status_stats.forEach((item) => {
        tagging.total_count += item.count
        switch (item.status) {
          case 'none':
            tagging.none_count = item.count
            break
          case 'pending':
            tagging.pending_count = item.count
            break
          case 'done':
            tagging.done_count = item.count
            break
          case 'failed':
            tagging.failed_count = item.count
            break
          case 'ignored':
            tagging.ignored_count = item.count
            break
        }
      })
    }

    const vectors = {
      total_count: vectorData.value?.total_count || 0,
      pending_count: vectorData.value?.pending_count || 0,
      processing_count: vectorData.value?.running_count || 0,
      completed_count: vectorData.value?.completed_count || 0,
      failed_count: vectorData.value?.failed_count || 0,
      reset_count: vectorData.value?.reset_count || 0,
      active_workers: vectorData.value?.active_workers || 0,
      max_workers: vectorData.value?.max_workers || 0,
    }

    return {
      tagging,
      vectors,
      health_status: {
        tagging_service: taggingData.value ? 'healthy' : 'down',
        vector_service: vectorData.value ? 'healthy' : 'down',
      },
    }
  })

  const hasData = computed(
    () =>
      data.value.tagging.total_count > 0 ||
      data.value.vectors.total_count > 0 ||
      data.value.tagging.pending_count > 0 ||
      data.value.vectors.processing_count > 0
  )

  const formatNumber = (num: number) => {
    if (num >= 1000000) {
      return `${(num / 1000000).toFixed(1)}M`
    } else if (num >= 1000) {
      return `${(num / 1000).toFixed(1)}K`
    }
    return num.toString()
  }

  const getCompletionRate = (completed: number, total: number) => {
    if (total === 0) {
      return '0%'
    }
    return `${Math.round((completed / total) * 100)}%`
  }

  function getStatusClass(status: string): string {
    switch (status) {
      case 'healthy':
        return 'status-healthy'
      case 'degraded':
        return 'status-warning'
      case 'down':
        return 'status-danger'
      default:
        return 'status-unknown'
    }
  }

  function getStatusText(status: string): string {
    switch (status) {
      case 'healthy':
        return $t('admin.dashboard.aiServices.statusNormal')
      case 'degraded':
        return $t('admin.dashboard.aiServices.statusDegraded')
      case 'down':
        return $t('admin.dashboard.aiServices.statusDown')
      default:
        return $t('admin.dashboard.aiServices.statusUnknown')
    }
  }

  async function fetchData() {
    try {
      loading.value = true

      const [taggingResult, vectorResult] = await Promise.all([
        getTaggingStats().catch(() => ({ success: false, data: null })),
        getVectorStats().catch(() => ({ success: false, data: null })),
      ])

      if (taggingResult?.success && taggingResult.data) {
        taggingData.value = taggingResult.data
      }

      if (vectorResult?.success && vectorResult.data) {
        vectorData.value = vectorResult.data
      }
    } catch (_error) {
    } finally {
      loading.value = false
    }
  }

  function refresh() {
    fetchData()
  }

  onMounted(() => {
    fetchData()
  })
</script>

<template>
  <div class="ai-services-card">
    <div class="card-header">
      <div class="header-icon">
        <svg class="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M9.75 17L9 20l-1 1h8l-1-1-.75-3M3 13h18M5 17h14a2 2 0 002-2V5a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z"
          />
        </svg>
      </div>
      <div class="header-title">
        <h3>{{ $t('admin.dashboard.aiServices.title') }}</h3>
      </div>
      <div class="refresh-btn" @click="refresh">
        <svg class="h-4 w-4" :class="{ 'animate-spin': loading }" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"
          />
        </svg>
      </div>
    </div>

    <div class="card-content">
      <div v-if="hasData" class="services-columns">
        <div class="service-column">
          <div class="service-header">
            <div class="service-icon tagging">
              <svg class="h-5 w-5" fill="currentColor" viewBox="0 0 20 20">
                <path
                  fill-rule="evenodd"
                  d="M17.707 9.293a1 1 0 010 1.414l-7 7a1 1 0 01-1.414 0l-7-7A.997.997 0 012 10V5a3 3 0 013-3h5c.256 0 .512.098.707.293l7 7zM5 6a1 1 0 100-2 1 1 0 000 2z"
                  clip-rule="evenodd"
                />
              </svg>
            </div>
            <div class="service-info">
              <div class="service-title">
                <span class="service-name">{{ $t('admin.dashboard.aiServices.taggingService') }}</span>
                <span class="completion-rate">{{ getCompletionRate(data.tagging.done_count, data.tagging.total_count) }}</span>
              </div>
              <div class="service-status" :class="getStatusClass(data.health_status.tagging_service)">
                <span class="status-dot" />
                {{ getStatusText(data.health_status.tagging_service) }}
              </div>
            </div>
          </div>

          <div class="service-stats">
            <div class="stat-item primary">
              <div class="stat-info">
                <div class="stat-value">{{ formatNumber(data.tagging.total_count) }}</div>
                <div class="stat-label">{{ $t('admin.dashboard.aiServices.totalCount') }}</div>
              </div>
              <div class="stat-icon"><i class="fas fa-database" /></div>
            </div>
            <div class="stat-item success">
              <div class="stat-info">
                <div class="stat-value">{{ formatNumber(data.tagging.done_count) }}</div>
                <div class="stat-label">{{ $t('admin.dashboard.aiServices.completed') }}</div>
              </div>
              <div class="stat-icon"><i class="fas fa-check" /></div>
            </div>
            <div class="stat-item warning">
              <div class="stat-info">
                <div class="stat-value">{{ formatNumber(data.tagging.none_count) }}</div>
                <div class="stat-label">{{ $t('admin.dashboard.aiServices.unprocessed') }}</div>
              </div>
              <div class="stat-icon"><i class="fas fa-clock" /></div>
            </div>
            <div class="stat-item processing">
              <div class="stat-info">
                <div class="stat-value">{{ formatNumber(data.tagging.pending_count) }}</div>
                <div class="stat-label">{{ $t('admin.dashboard.aiServices.processing') }}</div>
              </div>
              <div class="stat-icon"><i class="fas fa-spinner fa-spin" /></div>
            </div>
            <div class="stat-item danger">
              <div class="stat-info">
                <div class="stat-value">{{ formatNumber(data.tagging.failed_count) }}</div>
                <div class="stat-label">{{ $t('admin.dashboard.aiServices.failed') }}</div>
              </div>
              <div class="stat-icon"><i class="fas fa-times" /></div>
            </div>
            <div class="stat-item muted">
              <div class="stat-info">
                <div class="stat-value">{{ formatNumber(data.tagging.skipped_count) }}</div>
                <div class="stat-label">{{ $t('admin.dashboard.aiServices.skipped') }}</div>
              </div>
              <div class="stat-icon"><i class="fas fa-forward" /></div>
            </div>
          </div>
        </div>

        <div class="service-column">
          <div class="service-header">
            <div class="service-icon vector">
              <svg class="h-5 w-5" fill="currentColor" viewBox="0 0 20 20">
                <path
                  d="M3 4a1 1 0 011-1h12a1 1 0 011 1v2a1 1 0 01-1 1H4a1 1 0 01-1-1V4zM3 10a1 1 0 011-1h6a1 1 0 011 1v6a1 1 0 01-1 1H4a1 1 0 01-1-1v-6zM14 9a1 1 0 00-1 1v6a1 1 0 001 1h2a1 1 0 001-1v-6a1 1 0 00-1-1h-2z"
                />
              </svg>
            </div>
            <div class="service-info">
              <div class="service-title">
                <span class="service-name">{{ $t('admin.dashboard.aiServices.vectorService') }}</span>
                <span class="completion-rate">{{
                  getCompletionRate(data.vectors.completed_count, data.vectors.total_count)
                }}</span>
              </div>
              <div class="service-status" :class="getStatusClass(data.health_status.vector_service)">
                <span class="status-dot" />
                {{ getStatusText(data.health_status.vector_service) }}
              </div>
            </div>
          </div>

          <div class="service-stats">
            <div class="stat-item primary">
              <div class="stat-info">
                <div class="stat-value">{{ formatNumber(data.vectors.total_count) }}</div>
                <div class="stat-label">{{ $t('admin.dashboard.aiServices.totalCount') }}</div>
              </div>
              <div class="stat-icon"><i class="fas fa-database" /></div>
            </div>
            <div class="stat-item success">
              <div class="stat-info">
                <div class="stat-value">{{ formatNumber(data.vectors.completed_count) }}</div>
                <div class="stat-label">{{ $t('admin.dashboard.aiServices.completed') }}</div>
              </div>
              <div class="stat-icon"><i class="fas fa-check" /></div>
            </div>
            <div class="stat-item warning">
              <div class="stat-info">
                <div class="stat-value">{{ formatNumber(data.vectors.pending_count) }}</div>
                <div class="stat-label">{{ $t('admin.dashboard.aiServices.pending') }}</div>
              </div>
              <div class="stat-icon"><i class="fas fa-clock" /></div>
            </div>
            <div class="stat-item processing">
              <div class="stat-info">
                <div class="stat-value">{{ formatNumber(data.vectors.processing_count) }}</div>
                <div class="stat-label">{{ $t('admin.dashboard.aiServices.processing') }}</div>
              </div>
              <div class="stat-icon"><i class="fas fa-cog fa-spin" /></div>
            </div>
            <div class="stat-item danger">
              <div class="stat-info">
                <div class="stat-value">{{ formatNumber(data.vectors.failed_count) }}</div>
                <div class="stat-label">{{ $t('admin.dashboard.aiServices.failed') }}</div>
              </div>
              <div class="stat-icon"><i class="fas fa-times" /></div>
            </div>
            <div class="stat-item muted">
              <div class="stat-info">
                <div class="stat-value">{{ data.vectors.active_workers }}/{{ data.vectors.max_workers }}</div>
                <div class="stat-label">{{ $t('admin.dashboard.aiServices.workers') }}</div>
              </div>
              <div class="stat-icon"><i class="fas fa-server" /></div>
            </div>
          </div>
        </div>
      </div>
      <div v-else class="empty-state">
        <div class="empty-icon">🤖</div>
        <div class="empty-title">{{ $t('admin.dashboard.aiServices.emptyTitle') }}</div>
        <div class="empty-subtitle">{{ $t('admin.dashboard.aiServices.emptySubtitle') }}</div>
      </div>
    </div>
  </div>
</template>

<style scoped>
  .ai-services-card {
    @apply flex h-full flex-col;
    background: linear-gradient(
      135deg,
      rgba(var(--color-background-800-rgb), 0.6) 0%,
      rgba(var(--color-background-900-rgb), 0.4) 100%
    );
    backdrop-filter: blur(10px);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.1);
    border-radius: var(--radius-sm);
    box-shadow: var(--shadow-cyber-sm);
    overflow: hidden;
  }

  .card-header {
    @apply flex items-center gap-2 p-3 pb-2;
    border-bottom: 1px solid rgba(var(--color-brand-500-rgb), 0.1);
  }

  .header-icon {
    @apply flex h-10 w-10 flex-shrink-0 items-center justify-center;
    border-radius: var(--radius-sm);
    background: rgba(var(--color-brand-500-rgb), 0.15);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.25);
    color: var(--color-brand-500);
    box-shadow: 0 0 10px rgba(var(--color-brand-500-rgb), 0.2);
  }

  .header-title {
    @apply flex flex-1 flex-col gap-1;
  }

  .card-header h3 {
    color: var(--color-content-heading);
    font-size: 0.95rem;
    font-weight: 600;
    margin: 0;
    line-height: 1.2;
  }

  .refresh-btn {
    @apply flex h-5 w-5 cursor-pointer items-center justify-center;
    border-radius: var(--radius-full);
    background: rgba(var(--color-brand-500-rgb), 0.1);
    color: var(--color-brand-500);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.2);
    transition: all 0.2s ease;
  }

  .refresh-btn:hover {
    background: rgba(var(--color-brand-500-rgb), 0.2);
    border-color: rgba(var(--color-brand-500-rgb), 0.3);
    box-shadow: 0 0 15px rgba(var(--color-brand-500-rgb), 0.3);
  }

  .card-content {
    @apply flex-1 p-3 pt-2;
  }

  .empty-state {
    @apply flex flex-col items-center justify-center gap-2 py-8;
  }

  .empty-icon {
    font-size: 1.5rem;
  }
  .empty-title {
    color: var(--color-content-heading);
    font-weight: 600;
  }
  .empty-subtitle {
    color: var(--color-content-muted);
    font-size: 0.85rem;
  }

  .services-columns {
    @apply grid h-full grid-cols-2 gap-4;
  }

  .service-column {
    @apply flex flex-col;
    background: linear-gradient(
      135deg,
      rgba(var(--color-brand-500-rgb), 0.02) 0%,
      rgba(var(--color-background-800-rgb), 0.3) 100%
    );
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.08);
    border-radius: var(--radius-sm);
    padding: 1rem;
  }

  .service-header {
    @apply mb-3 flex items-center p-2.5;
    border-radius: var(--radius-sm);
    background: rgba(var(--color-brand-500-rgb), 0.04);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.08);
    gap: 0.75rem;
    height: 68px;
    min-height: 68px;
    max-height: 68px;
  }

  .service-icon {
    @apply flex h-12 w-12 flex-shrink-0 items-center justify-center;
    border-radius: var(--radius-sm);
    background: rgba(var(--color-brand-500-rgb), 0.15);
    color: var(--color-brand-500);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.25);
  }

  .service-icon svg {
    width: 1.5rem;
    height: 1.5rem;
  }

  .service-info {
    @apply flex-1;
    min-width: 0;
    display: flex;
    flex-direction: column;
    justify-content: center;
  }

  .service-title {
    @apply mb-1;
    display: flex;
    flex-direction: row;
    align-items: center;
    justify-content: space-between;
    gap: 0.5rem;
    flex-wrap: wrap;
  }

  .service-name {
    color: var(--color-content-heading);
    font-size: 0.95rem;
    font-weight: 700;
    line-height: 1.1;
    letter-spacing: 0.01em;
  }

  .completion-rate {
    color: var(--color-brand-500);
    font-size: 1.2rem;
    font-weight: 800;
    line-height: 1;
    letter-spacing: 0.01em;
    text-shadow: 0 0 15px rgba(var(--color-brand-500-rgb), 0.4);
  }

  .service-status {
    font-size: 0.7rem;
    font-weight: 600;
    display: inline-flex;
    align-items: center;
    gap: 0.4rem;
    padding: 1px 6px;
    border-radius: var(--radius-full);
    border: 1px solid rgba(148, 163, 184, 0.25);
    background: rgba(148, 163, 184, 0.06);
    align-self: flex-start;
  }

  .status-dot {
    width: 8px;
    height: 8px;
    border-radius: var(--radius-full);
    display: inline-block;
  }

  .service-status.status-healthy .status-dot {
    background: var(--color-success-500);
    box-shadow: 0 0 12px rgba(0, 255, 136, 0.5);
  }

  .service-status.status-warning .status-dot {
    background: var(--color-warning-500);
    box-shadow: 0 0 12px rgba(255, 136, 0, 0.5);
  }

  .service-status.status-danger .status-dot {
    background: var(--color-error-500);
    box-shadow: 0 0 12px rgba(255, 0, 102, 0.5);
  }

  .service-status.status-healthy {
    color: var(--color-success-500);
    background: rgba(0, 255, 136, 0.1);
    border-color: rgba(0, 255, 136, 0.3);
    box-shadow: 0 0 10px rgba(0, 255, 136, 0.2);
  }

  .service-status.status-warning {
    color: var(--color-warning-500);
    background: rgba(255, 136, 0, 0.1);
    border-color: rgba(255, 136, 0, 0.3);
    box-shadow: 0 0 10px rgba(255, 136, 0, 0.2);
  }

  .service-status.status-danger {
    color: var(--color-error-500);
    background: rgba(255, 0, 102, 0.1);
    border-color: rgba(255, 0, 102, 0.3);
    box-shadow: 0 0 10px rgba(255, 0, 102, 0.2);
  }

  .service-status.status-unknown {
    color: var(--color-content-muted);
    background: rgba(var(--color-content-rgb), 0.08);
    border-color: var(--color-border-default);
  }

  .service-stats {
    display: grid;
    grid-template-columns: 1fr 1fr 1fr;
    gap: 0.375rem;
    align-items: stretch;
    justify-items: stretch;
  }

  .stat-item {
    @apply text-left;
    border-radius: var(--radius-sm);
    background: rgba(var(--color-brand-500-rgb), 0.04);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.08);
    transition: all 0.3s ease;
    height: 65px !important;
    min-height: 65px !important;
    max-height: 65px !important;
    display: flex !important;
    flex-direction: row !important;
    align-items: center !important;
    justify-content: flex-start !important;
    gap: 0.25rem !important;
    box-sizing: border-box !important;
    padding: 0.5rem 0.5rem !important;
    margin: 0 !important;
    flex-grow: 0 !important;
    flex-shrink: 0 !important;
    position: relative;
    overflow: hidden;
  }

  .stat-item:hover {
    background: rgba(var(--color-brand-500-rgb), 0.08);
    border-color: rgba(var(--color-brand-500-rgb), 0.15);
  }

  .stat-info {
    @apply flex flex-col items-start justify-center;
    flex: 1;
    min-height: 0;
  }

  .stat-icon {
    @apply flex h-4 w-4 items-center justify-center rounded;
    font-size: 0.7rem;
    flex-shrink: 0;
    min-height: 16px;
    max-height: 16px;
    margin-top: 0;
    margin-right: 0.375rem;
    order: -1;
  }

  .stat-item.primary .stat-value,
  .stat-item.primary .stat-icon {
    color: var(--color-brand-500);
    text-shadow: 0 0 8px rgba(var(--color-brand-500-rgb), 0.3);
  }

  .stat-item.success .stat-value,
  .stat-item.success .stat-icon {
    color: var(--color-success-500);
    text-shadow: 0 0 8px rgba(0, 255, 136, 0.3);
  }

  .stat-item.warning .stat-value,
  .stat-item.warning .stat-icon {
    color: var(--color-warning-500);
    text-shadow: 0 0 8px rgba(255, 136, 0, 0.3);
  }

  .stat-item.processing .stat-value,
  .stat-item.processing .stat-icon {
    color: var(--color-brand-500);
    text-shadow: 0 0 8px rgba(var(--color-brand-500-rgb), 0.3);
  }

  .stat-item.danger .stat-value,
  .stat-item.danger .stat-icon {
    color: var(--color-error-500);
    text-shadow: 0 0 8px rgba(255, 0, 102, 0.3);
  }

  .stat-item.muted .stat-value,
  .stat-item.muted .stat-icon {
    color: var(--color-content-muted);
  }

  .stat-value {
    font-size: 1rem;
    font-weight: 700;
    line-height: 1;
    margin-bottom: 0.25rem;
  }

  .stat-label {
    color: var(--color-content-muted);
    font-size: 0.65rem;
    font-weight: 500;
    line-height: 1;
    margin: 0;
    padding: 0;
    white-space: nowrap;
  }

  @media (max-width: 768px) {
    .services-columns {
      @apply grid-cols-1 gap-2;
    }

    .service-stats {
      gap: 0.25rem;
    }

    .stat-item {
      gap: 0.125rem !important;
      padding: 0.375rem 0.375rem !important;
    }

    .stat-icon {
      @apply h-3 w-3;
      font-size: 0.6rem;
      min-height: 12px;
      max-height: 12px;
      margin-right: 0.25rem;
    }

    .stat-value {
      font-size: 0.75rem;
    }

    .stat-label {
      font-size: 0.6rem;
    }

    .service-icon {
      width: 2.25rem;
      height: 2.25rem;
    }
    .service-icon svg {
      width: 1.125rem;
      height: 1.125rem;
    }
    .service-header {
      gap: 0.5rem;
    }
    .completion-rate {
      font-size: 1.05rem;
    }
    .service-name {
      font-size: 0.9rem;
    }
  }
</style>
