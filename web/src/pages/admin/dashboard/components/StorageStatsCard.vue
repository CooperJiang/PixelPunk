<script setup lang="ts">
  import { computed, onMounted, ref } from 'vue'
  import { getDashboardStorageStats } from '@/api/admin/dashboard'
  import DashboardCard from './DashboardCard.vue'
  import { useTexts } from '@/composables/useTexts'
  const { $t } = useTexts()

  defineOptions({
    name: 'StorageStatsCard',
  })

  /* 数据状态 */
  const totalStorage = ref(0)
  const formattedStorage = ref('0 B')
  const averageFileSize = ref(0)
  const usagePercentage = ref(0)
  const growthRate = ref(0)
  const totalBandwidth = ref(0)
  const formattedBandwidth = ref('0 B')
  const newStorageToday = ref(0)
  const formattedNewStorage = ref('0 B')
  const lastUpdated = ref('')
  const isRefreshing = ref(false)

  const dailyChangeRate = computed(() => {
    const yesterday = totalStorage.value - newStorageToday.value
    if (yesterday <= 0 || newStorageToday.value === 0) return 0
    const rate = (newStorageToday.value / yesterday) * 100
    return isNaN(rate) || !isFinite(rate) ? 0 : rate
  })

  /* 格式化字节数 */
  const formatBytes = (bytes: number) => {
    if (bytes == null || isNaN(bytes) || !isFinite(bytes) || bytes <= 0) {
      return '0 B'
    }
    if (bytes >= 1024 * 1024 * 1024) {
      return `${(bytes / (1024 * 1024 * 1024)).toFixed(1)} GB`
    } else if (bytes >= 1024 * 1024) {
      return `${(bytes / (1024 * 1024)).toFixed(1)} MB`
    } else if (bytes >= 1024) {
      return `${(bytes / 1024).toFixed(1)} KB`
    }
    return `${bytes.toFixed(0)} B`
  }

  const fetchData = async () => {
    const startTime = Date.now()
    isRefreshing.value = true
    try {
      const result = await getDashboardStorageStats()

      if (result.success && result.data) {
        totalStorage.value = result.data.total_storage || 0
        formattedStorage.value = result.data.formatted_storage || '0 B'
        averageFileSize.value = result.data.average_file_size || 0
        usagePercentage.value = Math.round(result.data.usage_percentage || 0)
        growthRate.value = result.data.growth_rate || 0
        totalBandwidth.value = result.data.total_bandwidth || 0
        formattedBandwidth.value = result.data.formatted_bandwidth || '0 B'
        newStorageToday.value = result.data.new_storage_today || 0
        formattedNewStorage.value = result.data.formatted_new_storage || '0 B'

        const now = new Date()
        lastUpdated.value = `${now.getHours().toString().padStart(2, '0')}:${now.getMinutes().toString().padStart(2, '0')}`
      }
    } catch {
    } finally {
      const endTime = Date.now()
      const elapsed = endTime - startTime
      const minDuration = 1000
      const delay = Math.max(0, minDuration - elapsed)
      const now = new Date()
      const updated = `${now.getHours().toString().padStart(2, '0')}:${now.getMinutes().toString().padStart(2, '0')}`
      setTimeout(() => {
        isRefreshing.value = false
        if (!lastUpdated.value) {
          lastUpdated.value = updated
        } else {
          lastUpdated.value = updated
        }
      }, delay)
    }
  }

  const refresh = () => {
    if (!isRefreshing.value) {
      fetchData()
    }
  }

  onMounted(() => {
    fetchData()
  })
</script>

<template>
  <DashboardCard
    :title="$t('admin.dashboard.storageStats.title')"
    :badge="`${dailyChangeRate > 0 ? '+' : ''}${dailyChangeRate.toFixed(1)}%`"
    :badge-type="dailyChangeRate > 0 ? 'success' : dailyChangeRate < 0 ? 'danger' : 'normal'"
    :refreshable="true"
    :is-refreshing="isRefreshing"
    :last-updated="lastUpdated"
    @refresh="refresh"
  >
    <template #icon>
      <svg class="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-width="2"
          d="M4 7v10c0 2.21 3.582 4 8 4s8-1.79 8-4V7M4 7c0 2.21 3.582 4 8 4s8-1.79 8-4M4 7c0-2.21 3.582-4 8-4s8 1.79 8 4m0 5c0 2.21-3.582 4-8 4s-8-1.79-8-4"
        />
      </svg>
    </template>

    <div class="storage-content">
      <div class="main-metric">
        <div class="metric-value">
          <span class="value-number">{{ formattedStorage }}</span>
          <span class="value-label">{{ $t('admin.dashboard.storageStats.totalStorage') }}</span>
        </div>
        <div class="today-upload">
          <span class="upload-count">{{ formattedNewStorage }}</span>
          <span class="upload-label">{{ $t('admin.dashboard.storageStats.newToday') }}</span>
        </div>
      </div>

      <div class="sub-metrics">
        <div class="metric-item">
          <span class="metric-label">{{ $t('admin.dashboard.storageStats.averageSize') }}</span>
          <div class="metric-value-row">
            <div class="metric-icon size">
              <svg class="h-4 w-4" fill="currentColor" viewBox="0 0 20 20">
                <path
                  fill-rule="evenodd"
                  d="M3 4a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zm0 4a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zm0 4a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zm0 4a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1z"
                  clip-rule="evenodd"
                />
              </svg>
            </div>
            <span class="metric-number">{{ formatBytes(averageFileSize) }}</span>
          </div>
        </div>

        <div class="metric-item">
          <span class="metric-label">{{ $t('admin.dashboard.storageStats.growthRate') }}</span>
          <div class="metric-value-row">
            <div class="metric-icon growth" :class="{ negative: growthRate < 0 }">
              <svg v-if="growthRate >= 0" class="h-4 w-4" fill="currentColor" viewBox="0 0 20 20">
                <path
                  fill-rule="evenodd"
                  d="M3.293 9.707a1 1 0 010-1.414l6-6a1 1 0 011.414 0l6 6a1 1 0 01-1.414 1.414L11 5.414V17a1 1 0 11-2 0V5.414L4.707 9.707a1 1 0 01-1.414 0z"
                  clip-rule="evenodd"
                />
              </svg>
              <svg v-else class="h-4 w-4" fill="currentColor" viewBox="0 0 20 20">
                <path
                  fill-rule="evenodd"
                  d="M16.707 10.293a1 1 0 010 1.414l-6 6a1 1 0 01-1.414 0l-6-6a1 1 0 111.414-1.414L9 14.586V3a1 1 0 012 0v11.586l4.293-4.293a1 1 0 011.414 0z"
                  clip-rule="evenodd"
                />
              </svg>
            </div>
            <span class="metric-number">{{ Math.abs(growthRate).toFixed(1) }}%</span>
          </div>
        </div>

        <div class="metric-item">
          <span class="metric-label">{{ $t('admin.dashboard.storageStats.totalBandwidth') }}</span>
          <div class="metric-value-row">
            <div class="metric-icon size">
              <svg class="h-4 w-4" fill="currentColor" viewBox="0 0 20 20">
                <path
                  fill-rule="evenodd"
                  d="M4 3h12a1 1 0 011 1v12a1 1 0 01-1 1H4a1 1 0 01-1-1V4a1 1 0 011-1zm3 4a1 1 0 000 2h6a1 1 0 100-2H7zm0 4a1 1 0 000 2h6a1 1 0 100-2H7z"
                  clip-rule="evenodd"
                />
              </svg>
            </div>
            <span class="metric-number">{{ formattedBandwidth }}</span>
          </div>
        </div>
      </div>
    </div>
  </DashboardCard>
</template>

<style scoped>
  .storage-content {
    @apply flex flex-col;
    min-height: 145px;
    height: auto;
  }

  .main-metric {
    @apply mb-2.5 flex items-center justify-between;
    height: 56px;
    min-height: 56px;
    max-height: 56px;
  }

  .metric-value {
    @apply flex flex-col;
  }

  .value-number {
    color: var(--color-brand-500);
    font-size: 1.5rem;
    font-weight: bold;
    text-shadow: 0 0 8px rgba(var(--color-brand-500-rgb), 0.3);
    margin-bottom: 0.25rem;
    line-height: 1.1;
    white-space: nowrap;
  }

  .value-label {
    color: var(--color-slate-400);
    font-size: 0.7rem;
    font-weight: 500;
  }

  @media (max-width: 768px) {
    .value-number {
      font-size: 1.25rem;
    }

    .value-label {
      font-size: 0.65rem;
    }
  }

  @media (max-width: 640px) {
    .value-number {
      font-size: 1.1rem;
    }

    .value-label {
      font-size: 0.6rem;
    }
  }

  .today-upload {
    @apply flex flex-col items-end;
  }
  .upload-count {
    color: var(--color-success);
    font-size: 1rem;
    font-weight: bold;
    text-shadow: 0 0 5px rgba(34, 197, 94, 0.3);
    line-height: 1;
  }
  .upload-label {
    color: var(--color-slate-400);
    font-size: 0.65rem;
    margin-top: 0.125rem;
  }

  .sub-metrics {
    @apply grid grid-cols-3 gap-2;
  }

  .metric-item {
    @apply flex flex-col justify-center p-2;
    gap: 0.4rem;
    border-radius: var(--radius-sm);
    background: rgba(var(--color-brand-500-rgb), 0.05);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.08);
    transition: all 0.2s ease;
    min-height: 58px;
  }

  .metric-item:hover {
    background: rgba(var(--color-brand-500-rgb), 0.1);
    border-color: rgba(var(--color-brand-500-rgb), 0.15);
  }

  .metric-label {
    color: var(--color-slate-400);
    font-size: 0.65rem;
    display: block;
    text-align: center;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .metric-value-row {
    @apply flex items-center justify-center;
    gap: 0.4rem;
  }

  .metric-icon {
    @apply flex h-5 w-5 flex-shrink-0 items-center justify-center;
    border-radius: var(--radius-sm);
  }

  .metric-icon.size {
    background: rgba(59, 130, 246, 0.2);
    color: var(--color-cyan-500);
  }

  .metric-icon.growth {
    background: rgba(34, 197, 94, 0.2);
    color: var(--color-success);
  }

  .metric-icon.growth.negative {
    background: rgba(239, 68, 68, 0.2);
    color: var(--color-danger);
  }

  .metric-number {
    color: var(--color-content-heading);
    font-size: 0.85rem;
    font-weight: bold;
    line-height: 1.2;
    white-space: nowrap;
  }

  @media (max-width: 1536px) {
    .metric-icon {
      display: none;
    }

    .metric-value-row {
      gap: 0;
    }
  }

  @media (max-width: 1024px) {
    .metric-number {
      font-size: 0.8rem;
    }
  }

  @media (max-width: 768px) {
    .metric-number {
      font-size: 0.75rem;
    }

    .metric-label {
      font-size: 0.6rem;
    }
  }
</style>
