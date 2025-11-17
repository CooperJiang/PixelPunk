<script setup lang="ts">
  import { computed, onMounted, ref } from 'vue'
  import { getDashboardShareStats } from '@/api/admin/dashboard'
  import DashboardCard from './DashboardCard.vue'
  import { useTexts } from '@/composables/useTexts'
  const { $t } = useTexts()

  defineOptions({
    name: 'ShareStatsCard',
  })

  /* 数据状态 */
  const totalShares = ref(0)
  const publicShares = ref(0)
  const privateShares = ref(0)
  const totalVisits = ref(0)
  const totalDownloads = ref(0)
  const growthRate = ref(0)
  const newSharesToday = ref(0)
  const lastUpdated = ref('')
  const isRefreshing = ref(false)

  const _publicPercentage = computed(() =>
    totalShares.value > 0 ? Math.round((publicShares.value / totalShares.value) * 100) : 0
  )
  const _privatePercentage = computed(() =>
    totalShares.value > 0 ? Math.round((privateShares.value / totalShares.value) * 100) : 0
  )

  const safeGrowthRate = computed(() => {
    const rate = growthRate.value
    if (rate == null || isNaN(rate) || !isFinite(rate)) {
      return 0
    }
    return Math.round(rate)
  })

  /* 格式化数字 */
  const formatNumber = (num: number) => {
    if (num == null || isNaN(num) || !isFinite(num)) {
      return '0'
    }
    if (num >= 1000000) {
      return `${(num / 1000000).toFixed(1)}M`
    } else if (num >= 1000) {
      return `${(num / 1000).toFixed(1)}K`
    }
    return Math.round(num).toString()
  }

  const fetchData = async () => {
    const startTime = Date.now()
    isRefreshing.value = true
    try {
      const result = await getDashboardShareStats()

      if (result.success && result.data) {
        totalShares.value = result.data.total_shares || 0
        publicShares.value = result.data.public_shares || 0
        privateShares.value = result.data.private_shares || 0
        totalVisits.value = result.data.total_visits || 0
        totalDownloads.value = result.data.total_downloads || 0
        growthRate.value = result.data.growth_rate || 0
        newSharesToday.value = result.data.new_shares_today || 0

        const now = new Date()
        lastUpdated.value = `${now.getHours().toString().padStart(2, '0')}:${now.getMinutes().toString().padStart(2, '0')}`
      } else {
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
    :title="$t('admin.dashboard.shareStats.title')"
    :badge="`${safeGrowthRate > 0 ? '+' : ''}${safeGrowthRate}%`"
    :badge-type="safeGrowthRate > 0 ? 'success' : 'normal'"
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
          d="M8.684 13.342C8.886 12.938 9 12.482 9 12c0-.482-.114-.938-.316-1.342m0 2.684a3 3 0 110-2.684m0 2.684l6.632 3.316m-6.632-6l6.632-3.316m0 0a3 3 0 105.367-2.684 3 3 0 00-5.367 2.684zm0 9.316a3 3 0 105.367 2.684 3 3 0 00-5.367-2.684z"
        />
      </svg>
    </template>

    <div class="share-content">
      <div class="main-metric">
        <div class="metric-value">
          <span class="value-number">{{ formatNumber(totalShares) }}</span>
          <span class="value-label">{{ $t('admin.dashboard.shareStats.totalShares') }}</span>
        </div>
        <div class="today-upload">
          <span class="upload-count">{{ formatNumber(newSharesToday) }}</span>
          <span class="upload-label">{{ $t('admin.dashboard.shareStats.newToday') }}</span>
        </div>
      </div>

      <div class="sub-metrics">
        <div class="metric-item">
          <span class="metric-label">{{ $t('admin.dashboard.shareStats.publicShares') }}</span>
          <div class="metric-value-row">
            <div class="metric-icon public">
              <svg class="h-4 w-4" fill="currentColor" viewBox="0 0 20 20">
                <path d="M10 12a2 2 0 100-4 2 2 0 000 4z" />
                <path
                  fill-rule="evenodd"
                  d="M.458 10C1.732 5.943 5.522 3 10 3s8.268 2.943 9.542 7c-1.274 4.057-5.064 7-9.542 7S1.732 14.057.458 10zM14 10a4 4 0 11-8 0 4 4 0 018 0z"
                  clip-rule="evenodd"
                />
              </svg>
            </div>
            <span class="metric-number">{{ formatNumber(publicShares) }}</span>
          </div>
        </div>

        <div class="metric-item">
          <span class="metric-label">{{ $t('admin.dashboard.shareStats.privateShares') }}</span>
          <div class="metric-value-row">
            <div class="metric-icon private">
              <svg class="h-4 w-4" fill="currentColor" viewBox="0 0 20 20">
                <path
                  fill-rule="evenodd"
                  d="M5 9V7a5 5 0 0110 0v2a2 2 0 012 2v5a2 2 0 01-2 2H5a2 2 0 01-2-2v-5a2 2 0 012-2zm8-2v2H7V7a3 3 0 016 0z"
                  clip-rule="evenodd"
                />
              </svg>
            </div>
            <span class="metric-number">{{ formatNumber(privateShares) }}</span>
          </div>
        </div>

        <div class="metric-item">
          <span class="metric-label">{{ $t('admin.dashboard.shareStats.totalDownloads') }}</span>
          <div class="metric-value-row">
            <div class="metric-icon downloads">
              <svg class="h-4 w-4" fill="currentColor" viewBox="0 0 20 20">
                <path
                  fill-rule="evenodd"
                  d="M3 17a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zm3.293-7.707a1 1 0 011.414 0L9 10.586V3a1 1 0 112 0v7.586l1.293-1.293a1 1 0 111.414 1.414l-3 3a1 1 0 01-1.414 0l-3-3a1 1 0 010-1.414z"
                  clip-rule="evenodd"
                />
              </svg>
            </div>
            <span class="metric-number">{{ formatNumber(totalDownloads) }}</span>
          </div>
        </div>
      </div>
    </div>
  </DashboardCard>
</template>

<style scoped>
  .share-content {
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

  @media (max-width: 768px) {
    .value-number {
      font-size: 1.25rem;
    }
  }

  @media (max-width: 640px) {
    .value-number {
      font-size: 1.1rem;
    }
  }

  .value-label {
    color: var(--color-slate-400);
    font-size: 0.75rem;
    font-weight: 500;
  }

  @media (max-width: 640px) {
    .value-label {
      font-size: 0.65rem;
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
    border-radius: var(--radius-sm);
    gap: 0.4rem;
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

  .metric-icon.public {
    background: rgba(34, 197, 94, 0.2);
    color: var(--color-success);
  }

  .metric-icon.private {
    background: rgba(251, 146, 60, 0.2);
    color: var(--color-warning-500);
  }

  .metric-icon.downloads {
    background: rgba(59, 130, 246, 0.2);
    color: var(--color-cyan-500);
  }

  .metric-number {
    color: var(--color-content-heading);
    font-size: 0.85rem;
    font-weight: bold;
    line-height: 1.2;
    white-space: nowrap;
  }
</style>
