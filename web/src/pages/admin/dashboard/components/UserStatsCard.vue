<script setup lang="ts">
  import { computed, onMounted, ref } from 'vue'
  import { getDashboardUserStats } from '@/api/admin/dashboard'
  import { useTexts } from '@/composables/useTexts'
  import DashboardCard from './DashboardCard.vue'

  defineOptions({
    name: 'UserStatsCard',
  })

  const { $t } = useTexts()

  /* 数据状态 */
  const totalUsers = ref(0)
  const activeUsersToday = ref(0)
  const newUsersToday = ref(0)
  const growthRate = ref(0)
  const bannedUsers = ref(0)
  const lastUpdated = ref('')
  const isRefreshing = ref(false)
  const _trendClass = computed(() => {
    if (growthRate.value > 0) {
      return 'trend-positive'
    }
    if (growthRate.value < 0) {
      return 'trend-negative'
    }
    return 'trend-neutral'
  })

  const _trendIcon = computed(() => {
    if (growthRate.value > 0) {
      return 'M7 14l3-3 3 3m-3-3v12'
    }
    if (growthRate.value < 0) {
      return 'M17 10l-3 3-3-3m3 3V2'
    }
    return 'M5 12h14'
  })

  const safeGrowthRate = computed(() => {
    const rate = growthRate.value
    if (rate == null || isNaN(rate) || !isFinite(rate)) {
      return 0
    }
    return Math.round(rate)
  })

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
      const result = await getDashboardUserStats()

      if (result.success && result.data) {
        totalUsers.value = result.data.total_users || 0
        activeUsersToday.value = result.data.active_users_today || 0
        newUsersToday.value = result.data.new_users_today || 0
        growthRate.value = result.data.growth_rate || 0
        bannedUsers.value = result.data.banned_users || 0

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
        lastUpdated.value = updated
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
    :title="$t('admin.dashboard.userStats.title')"
    :badge="`${safeGrowthRate > 0 ? '+' : ''}${safeGrowthRate}%`"
    :badge-type="safeGrowthRate > 0 ? 'success' : safeGrowthRate < 0 ? 'danger' : 'normal'"
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
          d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z"
        />
      </svg>
    </template>

    <div class="user-content">
      <div class="metric-main">
        <div class="metric-left">
          <div class="metric-value">{{ formatNumber(totalUsers) }}</div>
          <div class="metric-label">{{ $t('admin.dashboard.userStats.totalUsers') }}</div>
        </div>
        <div class="metric-right">
          <div class="today-upload">
            <span class="upload-count">+{{ formatNumber(newUsersToday) }}</span>
            <span class="upload-label">{{ $t('admin.dashboard.userStats.newToday') }}</span>
          </div>
        </div>
      </div>

      <div class="metrics-grid">
        <div class="metric-item">
          <span class="metric-label">{{ $t('admin.dashboard.userStats.activeToday') }}</span>
          <div class="metric-value-row">
            <div class="metric-icon active">
              <svg class="h-4 w-4" fill="currentColor" viewBox="0 0 20 20">
                <path d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
            </div>
            <span class="metric-number">{{ formatNumber(activeUsersToday) }}</span>
          </div>
        </div>

        <div class="metric-item">
          <span class="metric-label">{{ $t('admin.dashboard.userStats.newToday') }}</span>
          <div class="metric-value-row">
            <div class="metric-icon new">
              <svg class="h-4 w-4" fill="currentColor" viewBox="0 0 20 20">
                <path
                  d="M8 9a3 3 0 100-6 3 3 0 000 6zM8 11a6 6 0 0 1 6 6H2a6 6 0 0 1 6-6zM16 7a1 1 0 10-2 0v1h-1a1 1 0 100 2h1v1a1 1 0 102 0v-1h1a1 1 0 100-2h-1V7z"
                />
              </svg>
            </div>
            <span class="metric-number">{{ formatNumber(newUsersToday) }}</span>
          </div>
        </div>

        <div class="metric-item">
          <span class="metric-label">{{ $t('admin.dashboard.userStats.banned') }}</span>
          <div class="metric-value-row">
            <div class="metric-icon danger">
              <svg class="h-4 w-4" fill="currentColor" viewBox="0 0 20 20">
                <path
                  fill-rule="evenodd"
                  d="M13.477 14.89A6 6 0 015.11 6.524l8.367 8.368zm1.414-1.414L6.524 5.11a6 6 0 018.367 8.367zM18 10a8 8 0 11-16 0 8 8 0 0116 0z"
                  clip-rule="evenodd"
                />
              </svg>
            </div>
            <span class="metric-number">{{ formatNumber(bannedUsers) }}</span>
          </div>
        </div>
      </div>
    </div>
  </DashboardCard>
</template>

<style scoped>
  .user-content {
    @apply flex flex-col;
    min-height: 145px;
    height: auto;
  }

  .metric-main {
    @apply mb-2.5 flex items-center justify-between;
    height: 56px;
    min-height: 56px;
    max-height: 56px;
  }

  .metric-left {
    @apply flex flex-col;
  }

  .metric-right {
    @apply flex items-center;
  }

  .metric-value {
    color: var(--color-brand-500);
    font-size: 1.5rem;
    font-weight: var(--font-bold);
    text-shadow: var(--shadow-glow-md);
    margin-bottom: var(--space-xs);
    line-height: 1.1;
    white-space: nowrap;
  }

  .metric-label {
    color: var(--color-slate-400);
    font-size: var(--text-xs);
    font-weight: var(--font-medium);
  }

  @media (max-width: 768px) {
    .metric-value {
      font-size: 1.25rem;
    }
  }

  @media (max-width: 640px) {
    .metric-value {
      font-size: 1.1rem;
    }

    .metric-label {
      font-size: 0.65rem;
    }
  }

  .today-upload {
    @apply flex flex-col items-end;
  }

  .upload-count {
    color: var(--color-success-500);
    font-size: var(--text-base);
    font-weight: var(--font-bold);
    text-shadow: var(--shadow-glow-sm);
    line-height: 1;
  }

  .upload-label {
    color: var(--color-slate-400);
    font-size: var(--text-xs);
    margin-top: var(--space-xs);
  }

  .metrics-grid {
    @apply grid grid-cols-3;
    gap: 0.5rem;
  }

  .metric-item {
    @apply flex flex-col justify-center p-2;
    border-radius: var(--radius-sm);
    gap: 0.4rem;
    background: rgba(var(--color-brand-500-rgb), 0.08);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.15);
    backdrop-filter: var(--backdrop-blur-md);
    min-height: 58px;
  }

  .metric-label {
    color: var(--color-slate-400);
    font-size: 0.65rem;
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
    @apply flex h-4 w-4 flex-shrink-0 items-center justify-center;
    border-radius: var(--radius-md);
  }

  .metric-icon.active {
    background: rgba(0, 255, 136, 0.2);
    color: var(--color-success-500);
    box-shadow: var(--shadow-glow-sm);
  }

  .metric-icon.new {
    background: rgba(var(--color-brand-500-rgb), 0.2);
    color: var(--color-brand-500);
    box-shadow: var(--shadow-glow-sm);
  }

  .metric-icon.danger {
    background: rgba(255, 0, 102, 0.2);
    color: var(--color-error-500);
    box-shadow: var(--shadow-glow-sm);
  }

  .metric-number {
    color: var(--color-content-heading);
    font-size: 0.85rem;
    font-weight: var(--font-bold);
    white-space: nowrap;
  }
</style>
