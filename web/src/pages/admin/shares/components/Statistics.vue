<script setup lang="ts">
  import { computed, onMounted, reactive, ref } from 'vue'
  import { useRouter } from 'vue-router'
  import { shareApi } from '@/api'
  import { useToast } from '@/components/Toast/useToast'
  import { useTexts } from '@/composables/useTexts'
  import type { ShareStats } from '@/api/share/types'

  const router = useRouter()
  const toast = useToast()
  const { $t } = useTexts()

  const stats = reactive<ShareStats>({
    total_shares: 0,
    active_shares: 0,
    views_today: 0,
    popular_shares: [],
    stats_by_date: [],
  })

  const trendPeriod = ref('7')

  const periodOptions = computed(() => [
    { label: $t('admin.shares.statistics.periodOptions.last7Days'), value: '7' },
    { label: $t('admin.shares.statistics.periodOptions.last30Days'), value: '30' },
    { label: $t('admin.shares.statistics.periodOptions.last90Days'), value: '90' },
  ])

  const loading = ref(true)

  const hasData = computed(() => stats.stats_by_date && stats.stats_by_date.length > 0)

  const loadStatistics = async () => {
    try {
      loading.value = true

      const result = await shareApi.getShareStats()

      if (result.success && result.data) {
        stats.total_shares = result.data.total_shares || 0
        stats.active_shares = result.data.active_shares || 0
        stats.views_today = result.data.views_today || 0
        stats.popular_shares = result.data.popular_shares || []
        stats.stats_by_date = result.data.stats_by_date || []
      }
    } catch (error: unknown) {
      toast.error((error as any)?.message || $t('admin.shares.statistics.messages.fetchFailed'))
    } finally {
      loading.value = false
    }
  }

  const refreshStats = () => {
    loadStatistics()
  }

  const viewShareDetail = async (share: { id: string; share_key: string }) => {
    try {
      loading.value = true

      try {
        const result = await shareApi.adminGetShareAccessToken(share.id)

        if (result.success && result.data && result.data.access_token) {
          const accessToken = result.data.access_token
          router.push(`/share/${share.share_key}?access_token=${accessToken}`)
        } else {
          router.push(`/share/${share.share_key}`)
          toast.warning($t('admin.shares.shareList.messages.getAccessFailed'))
        }
      } catch {
        router.push(`/share/${share.share_key}`)
        toast.warning($t('admin.shares.shareList.messages.getAccessFailed'))
      }
    } catch {
      router.push(`/share/${share.share_key}`)
      toast.warning($t('admin.shares.shareList.messages.getAccessFailed'))
    } finally {
      loading.value = false
    }
  }

  const calculateBarHeight = (value: number) => {
    if (!value) {
      return 0
    }

    const maxValue = Math.max(...stats.stats_by_date.map((item) => Math.max(item.total_views, item.new_shares * 10)))
    if (maxValue === 0) {
      return 0
    }

    return Math.min(Math.round((value / maxValue) * 100), 100)
  }

  const formatDateLabel = (dateString: string) => {
    if (!dateString) {
      return ''
    }
    const date = new Date(dateString)
    return `${date.getMonth() + 1}/${date.getDate()}`
  }

  onMounted(() => {
    loadStatistics()
  })
</script>

<template>
  <div class="statistics-component">
    <div class="stat-cards">
      <div class="stat-card">
        <div class="card-icon">
          <i class="fas fa-share-alt" />
        </div>
        <div class="card-content">
          <div class="card-title">{{ $t('admin.shares.statistics.cards.totalShares') }}</div>
          <div class="card-value">{{ stats.total_shares || 0 }}</div>
        </div>
      </div>

      <div class="stat-card">
        <div class="card-icon active-icon">
          <i class="fas fa-check-circle" />
        </div>
        <div class="card-content">
          <div class="card-title">{{ $t('admin.shares.statistics.cards.activeShares') }}</div>
          <div class="card-value">{{ stats.active_shares || 0 }}</div>
        </div>
      </div>

      <div class="stat-card">
        <div class="card-icon views-icon">
          <i class="fas fa-eye" />
        </div>
        <div class="card-content">
          <div class="card-title">{{ $t('admin.shares.statistics.cards.viewsToday') }}</div>
          <div class="card-value">{{ stats.views_today || 0 }}</div>
        </div>
      </div>
    </div>

    <div class="charts-grid">
      <div class="chart-container">
        <div class="chart-header">
          <h3 class="chart-title">{{ $t('admin.shares.statistics.charts.trend.title') }}</h3>
          <div class="chart-actions">
            <CyberDropdown v-model="trendPeriod" :options="periodOptions" />
          </div>
        </div>
        <div class="chart-body">
          <div v-if="loading" class="chart-placeholder">
            <div class="loading-spinner">
              <i class="fas fa-circle-notch fa-spin" />
              <span>{{ $t('admin.shares.statistics.loading.text') }}</span>
            </div>
          </div>
          <div v-else-if="!hasData" class="chart-placeholder">
            <div class="no-data">
              <i class="fas fa-chart-line" />
              <span>{{ $t('admin.shares.statistics.empty.noData') }}</span>
            </div>
          </div>
          <div v-else class="trend-chart">
            <div class="chart-bars">
              <div v-for="(item, index) in stats.stats_by_date" :key="index" class="chart-bar-group">
                <div class="date-label">{{ formatDateLabel(item.date) }}</div>
                <div class="bar-container">
                  <div
                    class="views-bar"
                    :style="{ height: calculateBarHeight(item.total_views) + '%' }"
                    :title="$t('admin.shares.statistics.charts.trend.viewsTooltip').replace('{count}', item.total_views)"
                  />
                  <div
                    class="shares-bar"
                    :style="{ height: calculateBarHeight(item.new_shares * 10) + '%' }"
                    :title="$t('admin.shares.statistics.charts.trend.sharesTooltip').replace('{count}', item.new_shares)"
                  />
                </div>
                <div class="bar-labels">
                  <div class="views-label">{{ item.total_views }}</div>
                  <div class="shares-label">{{ item.new_shares }}</div>
                </div>
              </div>
            </div>
            <div class="chart-legend">
              <div class="legend-item">
                <div class="legend-color views-color" />
                <div class="legend-text">{{ $t('admin.shares.statistics.legend.views') }}</div>
              </div>
              <div class="legend-item">
                <div class="legend-color shares-color" />
                <div class="legend-text">{{ $t('admin.shares.statistics.legend.shares') }}</div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div class="chart-container">
        <div class="chart-header">
          <h3 class="chart-title">{{ $t('admin.shares.statistics.charts.popular.title') }}</h3>
          <div class="chart-actions">
            <CyberButton size="mini" type="primary" @click="refreshStats">
              <i class="fas fa-sync-alt" />
            </CyberButton>
          </div>
        </div>
        <div class="chart-body">
          <div class="popular-shares">
            <div v-if="stats.popular_shares && stats.popular_shares.length > 0" class="shares-list">
              <div
                v-for="(share, index) in stats.popular_shares"
                :key="share.id"
                class="popular-share-item"
                @click="viewShareDetail(share)"
              >
                <div class="share-rank">{{ index + 1 }}</div>
                <div class="share-info">
                  <div class="share-name" :title="share.name">{{ share.name }}</div>
                  <div class="share-creator">{{ share.username }}</div>
                </div>
                <div class="share-views">
                  <span class="total-views">{{ share.total_views }}</span>
                  <span v-if="share.recent_views > 0" class="recent-badge">+{{ share.recent_views }}</span>
                </div>
              </div>
            </div>
            <div v-else class="chart-placeholder">
              <div class="no-data">
                <i class="fas fa-chart-pie" />
                <span>{{ $t('admin.shares.statistics.empty.noPopular') }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped lang="scss">
  .statistics-component {
    position: relative;
  }

  .stat-cards {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
    gap: var(--space-xl);
    margin-bottom: var(--space-2xl);
  }

  .stat-card {
    background: var(--color-background-700);
    border-radius: var(--radius-lg);
    padding: var(--space-xl);
    display: flex;
    align-items: center;
    border: 1px solid var(--color-border-default);
    transition: all var(--transition-normal) var(--ease-out);
    position: relative;
    overflow: hidden;
  }

  .stat-card:hover {
    border-color: var(--color-hover-border);
    transform: translateY(calc(var(--space-xs) * -1));
    box-shadow: var(--shadow-cyber-lg);
  }

  .card-icon {
    width: var(--space-4xl);
    height: var(--space-4xl);
    border-radius: var(--radius-md);
    background: rgba(var(--color-brand-500-rgb), 0.1);
    display: flex;
    align-items: center;
    justify-content: center;
    margin-right: var(--space-md);
    position: relative;
    z-index: 1;
  }

  .card-icon i {
    font-size: var(--text-xl);
    color: var(--color-brand-500);
    filter: drop-shadow(var(--shadow-glow-sm));
  }

  .active-icon {
    background: rgba(52, 199, 89, 0.1);
  }

  .active-icon i {
    color: rgb(52, 199, 89);
    filter: drop-shadow(0 0 5px rgba(52, 199, 89, 0.5));
  }

  .views-icon {
    background: rgba(255, 149, 0, 0.1);
  }

  .views-icon i {
    color: rgb(255, 149, 0);
    filter: drop-shadow(0 0 5px rgba(255, 149, 0, 0.5));
  }

  .card-content {
    flex: 1;
    position: relative;
    z-index: 1;
  }

  .card-title {
    color: var(--color-content-default);
    font-size: var(--text-sm);
    margin-bottom: var(--space-xs);
  }

  .card-value {
    color: var(--color-content-heading);
    font-size: var(--text-xl);
    font-weight: var(--font-semibold);
  }

  .charts-grid {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: var(--space-xl);
  }

  .chart-container {
    background: var(--color-background-700);
    border-radius: var(--radius-lg);
    border: 1px solid var(--color-border-default);
    overflow: hidden;
  }

  .chart-header {
    padding: var(--space-md) var(--space-lg);
    display: flex;
    justify-content: space-between;
    align-items: center;
    border-bottom: 1px solid rgba(var(--color-brand-500-rgb), 0.1);
  }

  .chart-title {
    margin: 0;
    color: var(--color-content-heading);
    font-size: var(--text-lg);
    font-weight: var(--font-semibold);
  }

  .chart-body {
    padding: var(--space-xl);
    min-height: 300px;
    position: relative;
  }

  .chart-placeholder {
    display: flex;
    justify-content: center;
    align-items: center;
    height: 250px;
  }

  .loading-spinner,
  .no-data {
    display: flex;
    flex-direction: column;
    align-items: center;
    color: var(--color-content-muted);
  }

  .loading-spinner i,
  .no-data i {
    font-size: var(--text-3xl);
    margin-bottom: var(--space-md);
    color: rgba(var(--color-brand-500-rgb), 0.3);
  }

  .trend-chart {
    height: 250px;
    display: flex;
    flex-direction: column;
  }

  .chart-bars {
    flex: 1;
    display: flex;
    gap: var(--space-xl);
    align-items: flex-end;
    padding-bottom: var(--space-lg);
    position: relative;
  }

  .chart-bar-group {
    display: flex;
    flex-direction: column;
    align-items: center;
    flex: 1;
  }

  .date-label {
    font-size: var(--text-xs);
    color: var(--color-content-muted);
    margin-top: var(--space-xs);
    position: absolute;
    bottom: 0;
  }

  .bar-container {
    height: 200px;
    display: flex;
    align-items: flex-end;
    justify-content: center;
    gap: var(--space-xs);
    width: 100%;
  }

  .views-bar,
  .shares-bar {
    width: 45%;
    background: linear-gradient(to top, rgba(var(--color-brand-500-rgb), 0.2), rgba(var(--color-brand-500-rgb), 0.7));
    border-radius: var(--radius-sm) var(--radius-sm) 0 0;
    transition: height var(--transition-normal) var(--ease-out);
    position: relative;
    min-height: 1px;
  }

  .shares-bar {
    background: linear-gradient(to top, rgba(255, 149, 0, 0.2), rgba(255, 149, 0, 0.7));
  }

  .bar-labels {
    display: flex;
    justify-content: center;
    gap: var(--space-sm);
    font-size: var(--text-xs);
    position: absolute;
    bottom: var(--space-md);
    opacity: 0;
    transition: opacity var(--transition-normal);
  }

  .chart-bar-group:hover .bar-labels {
    opacity: 1;
  }

  .views-label {
    color: rgba(var(--color-brand-500-rgb), 0.9);
  }

  .shares-label {
    color: rgba(255, 149, 0, 0.9);
  }

  .chart-legend {
    display: flex;
    justify-content: center;
    gap: var(--space-xl);
    margin-top: var(--space-md);
  }

  .legend-item {
    display: flex;
    align-items: center;
    gap: var(--space-sm);
  }

  .legend-color {
    width: var(--space-md);
    height: var(--space-sm);
    border-radius: var(--radius-xs);
  }

  .views-color {
    background: rgba(var(--color-brand-500-rgb), 0.7);
  }

  .shares-color {
    background: rgba(255, 149, 0, 0.7);
  }

  .legend-text {
    font-size: var(--text-sm);
    color: var(--color-content-default);
  }

  .popular-shares {
    height: 250px;
    overflow-y: auto;
  }

  .shares-list {
    padding-right: var(--space-sm);
  }

  .popular-share-item {
    display: flex;
    align-items: center;
    padding: var(--space-sm);
    border-radius: var(--radius-md);
    margin-bottom: var(--space-sm);
    background: rgba(var(--color-brand-500-rgb), 0.05);
    cursor: pointer;
    transition: all var(--transition-normal);
  }

  .popular-share-item:hover {
    background: rgba(var(--color-brand-500-rgb), 0.1);
  }

  .popular-share-item:last-child {
    margin-bottom: 0;
  }

  .share-rank {
    width: var(--space-xl);
    height: var(--space-xl);
    border-radius: var(--radius-full);
    background: rgba(var(--color-brand-500-rgb), 0.1);
    color: var(--color-brand-500);
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: var(--text-sm);
    font-weight: var(--font-bold);
    margin-right: var(--space-sm);
  }

  .share-info {
    flex: 1;
    overflow: hidden;
  }

  .share-name {
    font-weight: var(--font-medium);
    color: var(--color-content-heading);
    margin-bottom: var(--space-xs);
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .share-creator {
    font-size: var(--text-sm);
    color: var(--color-content-muted);
  }

  .share-views {
    display: flex;
    align-items: center;
    gap: var(--space-sm);
  }

  .total-views {
    color: var(--color-content-default);
    font-weight: var(--font-medium);
  }

  .recent-badge {
    background: rgba(52, 199, 89, 0.1);
    color: rgb(52, 199, 89);
    padding: var(--space-xs) var(--space-xs);
    border-radius: var(--radius-sm);
    font-size: var(--text-xs);
  }

  @media (max-width: 991px) {
    .charts-grid {
      grid-template-columns: 1fr;
    }
  }
</style>
