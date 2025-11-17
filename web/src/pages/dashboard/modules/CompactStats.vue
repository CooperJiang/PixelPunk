<script setup lang="ts">
  import { computed, onMounted, ref, watch } from 'vue'
  import { type WorkspaceStats, getWorkspaceStats } from '@/api/workspace'
  import { useToast } from '@/components/Toast/useToast'
  import { useTexts } from '@/composables/useTexts'

  defineOptions({
    name: 'CompactStats',
  })

  const props = defineProps<{
    refreshKey?: number
  }>()

  const toast = useToast()
  const { $t } = useTexts()
  const loading = ref(false)

  const stats = ref<WorkspaceStats>({
    total_files: 0,
    used_storage: 0,
    used_storage_formatted: '0 B',
    total_storage: 0,
    total_storage_formatted: '0 B',
    total_views: 0,
    total_shares: 0,
    used_bandwidth: 0,
    used_bandwidth_formatted: '0 B',
    total_bandwidth: 0,
    total_bandwidth_formatted: '0 B',
  })

  const fetchStats = async () => {
    try {
      loading.value = true
      const result = await getWorkspaceStats()

      if (result.success) {
        stats.value = result.data || stats.value
      } else {
        toast.error(result.message || $t('dashboard.errors.fetchStatsFailed'))
      }
    } catch (_error) {
      toast.error($t('dashboard.errors.networkFailed'))
    } finally {
      loading.value = false
    }
  }

  onMounted(() => {
    fetchStats()
  })

  watch(
    () => props.refreshKey,
    () => {
      if (props.refreshKey !== undefined) {
        fetchStats()
      }
    }
  )

  const storagePercent = computed(() => {
    if (stats.value.total_storage <= 0) {
      return 0
    }
    return Math.min(Math.round((stats.value.used_storage / stats.value.total_storage) * 100), 100)
  })

  const bandwidthPercent = computed(() => {
    if (stats.value.total_bandwidth <= 0) {
      return 0
    }
    return Math.min(Math.round((stats.value.used_bandwidth / stats.value.total_bandwidth) * 100), 100)
  })

  const getProgressColor = (percentage: number) => {
    if (percentage <= 60) {
      const ratio = percentage / 60
      const red = Math.floor(ratio * 150) // 0 -> 150
      const green = 255
      const blue = Math.floor(136 - ratio * 50) // 136 -> 86
      return {
        gradient: `linear-gradient(90deg, var(--color-success-500), #${red.toString(16).padStart(2, '0')}${green.toString(16)}${blue.toString(16).padStart(2, '0')})`,
        glow: `rgba(${red}, ${green}, ${blue}, 0.4)`,
        level: 'healthy',
      }
    } else if (percentage <= 85) {
      const ratio = (percentage - 60) / 25
      const red = 255
      const green = Math.floor(200 - ratio * 80) // 200 -> 120
      const blue = Math.floor(50 - ratio * 24) // 50 -> 26
      return {
        gradient: `linear-gradient(90deg, var(--color-warning-400), #${red.toString(16)}${green.toString(16).padStart(2, '0')}${blue.toString(16).padStart(2, '0')})`,
        glow: `rgba(${red}, ${green}, ${blue}, 0.4)`,
        level: 'warning',
      }
    }
    const ratio = Math.min((percentage - 85) / 15, 1) // 防止超过100%
    const red = 255
    const green = Math.floor(120 - ratio * 120) // 120 -> 0
    const blue = Math.floor(26 - ratio * 26) // 26 -> 0
    return {
      gradient: `linear-gradient(90deg, var(--color-warning-600), #${red.toString(16)}${green.toString(16).padStart(2, '0')}${blue.toString(16).padStart(2, '0')})`,
      glow: `rgba(${red}, ${Math.max(green, 50)}, ${Math.max(blue, 50)}, 0.5)`,
      level: 'danger',
    }
  }

  const storageProgressStyle = computed(() => {
    const color = getProgressColor(storagePercent.value)
    return {
      background: color.gradient,
      boxShadow: `0 0 12px ${color.glow}`,
    }
  })

  const bandwidthProgressStyle = computed(() => {
    const color = getProgressColor(bandwidthPercent.value)
    return {
      background: color.gradient,
      boxShadow: `0 0 12px ${color.glow}`,
    }
  })

  const getProgressLevel = (percentage: number) => getProgressColor(percentage).level

  const getStatusText = (percentage: number) => {
    if (percentage <= 60) {
      return $t('dashboard.stats.quotas.status.healthy')
    } else if (percentage <= 85) {
      return $t('dashboard.stats.quotas.status.warning')
    }
    return $t('dashboard.stats.quotas.status.danger')
  }

  const getStatusDescription = (percentage: number) => {
    if (percentage <= 60) {
      return $t('dashboard.stats.quotas.description.healthy')
    } else if (percentage <= 85) {
      return $t('dashboard.stats.quotas.description.warning')
    }
    return $t('dashboard.stats.quotas.description.danger')
  }

  const formatNumber = (num: number | undefined): string => {
    if (num === undefined || num === null || isNaN(num)) {
      return '0'
    }
    if (num >= 1000000) {
      return `${(num / 1000000).toFixed(1)}M`
    } else if (num >= 1000) {
      return `${(num / 1000).toFixed(1)}K`
    }
    return num.toString()
  }

  const statItems = ref([
    {
      id: 'images',
      icon: 'fas fa-images',
      label: $t('dashboard.stats.items.images'),
      getValue: () => formatNumber(stats.value.total_files),
      colorVar: 'primary',
    },
    {
      id: 'storage',
      icon: 'fas fa-hdd',
      label: $t('dashboard.stats.items.storage'),
      getValue: () => stats.value.used_storage_formatted,
      colorVar: 'secondary',
    },
    {
      id: 'views',
      icon: 'fas fa-eye',
      label: $t('dashboard.stats.items.views'),
      getValue: () => formatNumber(stats.value.total_views),
      colorVar: 'accent',
    },
    {
      id: 'shares',
      icon: 'fas fa-share-alt',
      label: $t('dashboard.stats.items.shares'),
      getValue: () => stats.value.total_shares,
      colorVar: 'info',
    },
  ])
</script>

<template>
  <div class="compact-stats cyber-card">
    <div class="stats-content">
      <div class="stats-header">
        <h3 class="section-title">
          <i class="fas fa-chart-pie" />
          {{ $t('dashboard.stats.title') }}
        </h3>
      </div>

      <div class="stats-grid">
        <div v-for="item in statItems" :key="item.id" class="stat-item" :data-color="item.colorVar">
          <div class="stat-icon">
            <i :class="item.icon" />
          </div>
          <div class="stat-content">
            <div class="stat-value">{{ loading ? '-' : item.getValue() }}</div>
            <div class="stat-label">{{ item.label }}</div>
          </div>
        </div>
      </div>

      <div class="quotas-section">
        <h4 class="section-subtitle">
          <i class="fas fa-tachometer-alt" />
          {{ $t('dashboard.stats.quotas.title') }}
        </h4>

        <div class="quota-item">
          <div class="quota-header">
            <span class="quota-name">
              <i class="fas fa-database" />
              {{ $t('dashboard.stats.quotas.storage.label') }}
              <span
                class="status-badge"
                :class="getProgressLevel(storagePercent)"
                :title="getStatusDescription(storagePercent)"
                >{{ getStatusText(storagePercent) }}</span
              >
            </span>
            <span class="quota-value" :class="getProgressLevel(storagePercent)">{{ storagePercent }}%</span>
          </div>
          <div class="quota-progress">
            <div class="progress-bar storage" :class="[getProgressLevel(storagePercent)]">
              <div class="progress-fill" :style="{ width: storagePercent + '%', ...storageProgressStyle }" />
            </div>
          </div>
          <div class="quota-details">
            <span class="used">{{ stats.used_storage_formatted }}</span>
            <span class="total">{{ stats.total_storage_formatted }}</span>
          </div>
        </div>

        <div class="quota-item">
          <div class="quota-header">
            <span class="quota-name">
              <i class="fas fa-wifi" />
              {{ $t('dashboard.stats.quotas.bandwidth.label') }}
              <span
                class="status-badge"
                :class="getProgressLevel(bandwidthPercent)"
                :title="getStatusDescription(bandwidthPercent)"
                >{{ getStatusText(bandwidthPercent) }}</span
              >
            </span>
            <span class="quota-value" :class="getProgressLevel(bandwidthPercent)">{{ bandwidthPercent }}%</span>
          </div>
          <div class="quota-progress">
            <div class="progress-bar bandwidth" :class="[getProgressLevel(bandwidthPercent)]">
              <div class="progress-fill" :style="{ width: bandwidthPercent + '%', ...bandwidthProgressStyle }" />
            </div>
          </div>
          <div class="quota-details">
            <span class="used">{{ stats.used_bandwidth_formatted }}</span>
            <span class="total">{{ stats.total_bandwidth_formatted }}</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
  .compact-stats {
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
  }

  .stats-content {
    display: flex;
    flex-direction: column;
    gap: 20px;
  }

  .stats-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .section-title {
    margin: 0;
    font-size: 16px;
    font-weight: 600;
    color: var(--color-content-heading);
    display: flex;
    align-items: center;
    gap: 10px;
  }

  .section-title i {
    color: var(--color-brand-500);
    font-size: 14px;
  }

  .section-subtitle {
    margin: 0 0 12px 0;
    font-size: 14px;
    font-weight: 600;
    color: var(--color-content-default);
    display: flex;
    align-items: center;
    gap: 8px;
  }

  .section-subtitle i {
    color: var(--color-warning-500);
    font-size: 12px;
  }

  .stats-grid {
    display: grid;
    grid-template-columns: repeat(4, 1fr);
    gap: 16px;
  }

  .stat-item {
    padding: 16px;
    background: var(--color-background-700);
    border-radius: var(--radius-sm);
    border: 1px solid var(--color-border-default);
    display: flex;
    align-items: center;
    gap: 12px;
    transition: all 0.3s ease;
    position: relative;
    overflow: hidden;
  }

  .stat-item:hover {
    transform: translateY(-2px);
    border-color: var(--color-hover-border);
    box-shadow: var(--shadow-cyber-md);
  }

  .stat-icon {
    width: 40px;
    height: 40px;
    border-radius: var(--radius-sm);
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 16px;
    flex-shrink: 0;
    position: relative;
    transition: all 0.3s ease;
  }

  .stat-item[data-color='primary'] .stat-icon {
    background: rgba(var(--color-func-primary-rgb), 0.12);
    color: var(--color-func-primary);
    border: 1px solid rgba(var(--color-func-primary-rgb), 0.25);
  }

  .stat-item[data-color='secondary'] .stat-icon {
    background: rgba(var(--color-func-secondary-rgb), 0.12);
    color: var(--color-func-secondary);
    border: 1px solid rgba(var(--color-func-secondary-rgb), 0.25);
  }

  .stat-item[data-color='accent'] .stat-icon {
    background: rgba(var(--color-brand-500-rgb), 0.12);
    color: var(--color-func-accent);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.25);
  }

  .stat-item[data-color='info'] .stat-icon {
    background: rgba(var(--color-func-info-rgb), 0.12);
    color: var(--color-func-info);
    border: 1px solid rgba(var(--color-func-info-rgb), 0.25);
  }

  .stat-item[data-color='primary']:hover .stat-icon {
    background: rgba(var(--color-func-primary-rgb), 0.2);
    border-color: rgba(var(--color-func-primary-rgb), 0.4);
    box-shadow: 0 0 15px rgba(var(--color-func-primary-rgb), 0.3);
  }

  .stat-item[data-color='secondary']:hover .stat-icon {
    background: rgba(var(--color-func-secondary-rgb), 0.2);
    border-color: rgba(var(--color-func-secondary-rgb), 0.4);
    box-shadow: 0 0 15px rgba(var(--color-func-secondary-rgb), 0.3);
  }

  .stat-item[data-color='accent']:hover .stat-icon {
    background: rgba(var(--color-brand-500-rgb), 0.2);
    border-color: rgba(var(--color-brand-500-rgb), 0.4);
    box-shadow: 0 0 15px rgba(var(--color-brand-500-rgb), 0.3);
  }

  .stat-item[data-color='info']:hover .stat-icon {
    background: rgba(var(--color-func-info-rgb), 0.2);
    border-color: rgba(var(--color-func-info-rgb), 0.4);
    box-shadow: 0 0 15px rgba(var(--color-func-info-rgb), 0.3);
  }

  .stat-content {
    flex: 1;
  }

  .stat-value {
    font-size: 18px;
    font-weight: 700;
    margin-bottom: 4px;
    font-family: 'SF Mono', monospace;
    line-height: 1;
    transition: color 0.3s ease;
  }

  .stat-item[data-color='primary'] .stat-value {
    color: var(--color-func-primary);
  }

  .stat-item[data-color='secondary'] .stat-value {
    color: var(--color-func-secondary);
  }

  .stat-item[data-color='accent'] .stat-value {
    color: var(--color-func-accent);
  }

  .stat-item[data-color='info'] .stat-value {
    color: var(--color-func-info);
  }

  .stat-label {
    font-size: 11px;
    color: var(--color-content-muted);
    font-weight: 500;
    line-height: 1;
  }

  .quotas-section {
    padding-top: 4px;
  }

  .quota-item {
    margin-bottom: 16px;
    padding: 16px;
    background: var(--color-background-700);
    border-radius: var(--radius-sm);
    border: 1px solid var(--color-border-default);
  }

  .quota-item:last-child {
    margin-bottom: 0;
  }

  .quota-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 8px;
  }

  .quota-name {
    font-size: 13px;
    font-weight: 600;
    color: var(--color-content-default);
    display: flex;
    align-items: center;
    gap: 8px;
  }

  .quota-name i {
    font-size: 12px;
    color: var(--color-brand-500);
  }

  .status-badge {
    margin-left: 8px;
    padding: 2px 8px;
    border-radius: var(--radius-sm);
    font-size: 10px;
    font-weight: 600;
    text-transform: uppercase;
    letter-spacing: 0.5px;
    transition: all 0.3s ease;
  }

  .status-badge.healthy {
    background: rgba(var(--color-success-rgb), 0.15);
    color: var(--color-success-500);
    border: 1px solid rgba(var(--color-success-rgb), 0.3);
    position: relative;
  }

  .status-badge.warning {
    background: rgba(var(--color-warning-rgb), 0.15);
    color: var(--color-warning-400);
    border: 1px solid rgba(var(--color-warning-rgb), 0.3);
    position: relative;
  }

  .status-badge.danger {
    background: rgba(var(--color-error-rgb), 0.15);
    color: var(--color-warning-600);
    border: 1px solid rgba(var(--color-error-rgb), 0.3);
    animation: pulse-tension 2s ease-in-out infinite;
    position: relative;
  }

  @keyframes pulse-tension {
    0%,
    100% {
      box-shadow: 0 0 8px rgba(var(--color-error-rgb), 0.3);
      transform: scale(1);
    }
    50% {
      box-shadow: 0 0 15px rgba(var(--color-error-rgb), 0.6);
      transform: scale(1.02);
    }
  }

  .quota-value {
    font-size: 13px;
    font-weight: 700;
    font-family: 'SF Mono', monospace;
    transition: color 0.3s ease;
  }

  .quota-value.healthy {
    color: var(--color-success-500);
    text-shadow: 0 0 8px rgba(var(--color-success-rgb), 0.3);
  }

  .quota-value.warning {
    color: var(--color-warning-500);
    text-shadow: 0 0 8px rgba(var(--color-warning-rgb), 0.3);
  }

  .quota-value.danger {
    color: var(--color-error-500);
    text-shadow: 0 0 8px rgba(var(--color-error-rgb), 0.4);
  }

  .quota-progress {
    margin-bottom: 8px;
  }

  .progress-bar {
    width: 100%;
    height: 8px;
    background: rgba(var(--color-content-rgb), 0.08);
    border-radius: var(--radius-sm);
    overflow: hidden;
    position: relative;
    box-shadow:
      0 1px 0 rgba(0, 0, 0, 0.1) inset,
      0 -1px 0 rgba(255, 255, 255, 0.05) inset;
  }

  .progress-bar.healthy {
    background: rgba(var(--color-success-rgb), 0.08);
  }

  .progress-bar.warning {
    background: rgba(var(--color-warning-rgb), 0.08);
  }

  .progress-bar.danger {
    background: rgba(var(--color-error-rgb), 0.08);
  }

  .progress-fill {
    height: 100%;
    border-radius: var(--radius-sm);
    transition: width 0.6s cubic-bezier(0.4, 0, 0.2, 1);
    position: relative;
    box-shadow:
      0 1px 0 rgba(255, 255, 255, 0.3) inset,
      0 2px 4px rgba(0, 0, 0, 0.2);
  }

  .progress-fill::after {
    content: '';
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

  .quota-details {
    display: flex;
    justify-content: space-between;
    align-items: center;
    font-size: 11px;
    font-weight: 500;
  }

  .quota-details .used {
    color: var(--color-content-default);
  }

  .quota-details .total {
    color: var(--color-content-muted);
  }

  @media (max-width: 768px) {
    .compact-stats {
      padding: 16px;
    }

    .stats-content {
      gap: 16px;
    }

    .stats-grid {
      grid-template-columns: repeat(2, 1fr);
      gap: 12px;
    }

    .stat-item {
      padding: 12px;
    }

    .stat-icon {
      width: 36px;
      height: 36px;
      font-size: 14px;
    }

    .stat-value {
      font-size: 16px;
    }
  }

  @media (max-width: 480px) {
    .compact-stats {
      padding: 12px;
    }

    .stats-grid {
      grid-template-columns: 1fr;
      gap: 8px;
    }

    .stat-item {
      padding: 12px;
    }

    .quota-item {
      padding: 12px;
    }

    .quota-name {
      flex-direction: column;
      align-items: flex-start;
      gap: 4px;
    }

    .status-badge {
      margin-left: 0;
      margin-top: 4px;
      font-size: 9px;
      padding: 2px 6px;
    }
  }
</style>
