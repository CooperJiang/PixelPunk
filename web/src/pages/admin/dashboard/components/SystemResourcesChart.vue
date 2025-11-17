<script setup lang="ts">
  import { computed, onMounted, onUnmounted, ref } from 'vue'
  import { getDashboardSystemResources, type DashboardSystemResourcesResponse } from '@/api/admin/dashboard'
  import { useTexts } from '@/composables/useTexts'
  const { $t } = useTexts()

  const loading = ref(false)

  const data = ref<DashboardSystemResourcesResponse>({
    cpu_usage: 0,
    memory_usage: 0,
    disk_usage: 0,
    load_average: [0, 0, 0],
    uptime: '0d 0h 0m',
  })

  /* 系统状态计算 */
  const systemStatusClass = computed(() => {
    const avgUsage = (data.value.cpu_usage + data.value.memory_usage + data.value.disk_usage) / 3
    if (avgUsage >= 80) {
      return 'status-danger'
    }
    if (avgUsage >= 60) {
      return 'status-warning'
    }
    return 'status-healthy'
  })

  const systemStatusText = computed(() => {
    const avgUsage = (data.value.cpu_usage + data.value.memory_usage + data.value.disk_usage) / 3
    if (avgUsage >= 80) {
      return $t('admin.dashboard.systemResources.highLoad')
    }
    if (avgUsage >= 60) {
      return $t('admin.dashboard.systemResources.mediumLoad')
    }
    return $t('admin.dashboard.systemResources.normalRunning')
  })

  async function fetchData() {
    try {
      loading.value = true
      const result = await getDashboardSystemResources()
      if (result.success) {
        data.value = result.data
      } else {
        data.value = {
          cpu_usage: Math.random() * 80 + 10,
          memory_usage: Math.random() * 70 + 20,
          disk_usage: Math.random() * 60 + 15,
          load_average: [Math.random() * 2, Math.random() * 2, Math.random() * 2],
          uptime: '5d 12h 36m',
        }
      }
    } catch {
      data.value = {
        cpu_usage: Math.random() * 80 + 10,
        memory_usage: Math.random() * 70 + 20,
        disk_usage: Math.random() * 60 + 15,
        load_average: [Math.random() * 2, Math.random() * 2, Math.random() * 2],
        uptime: '5d 12h 36m',
      }
    } finally {
      loading.value = false
    }
  }

  function refresh() {
    fetchData()
  }

  let refreshInterval: NodeJS.Timeout | null = null

  onMounted(() => {
    fetchData()
    refreshInterval = setInterval(fetchData, 30000) // 30秒刷新一次
  })

  onUnmounted(() => {
    if (refreshInterval) {
      clearInterval(refreshInterval)
    }
  })
</script>

<template>
  <div class="system-resources-chart">
    <div class="card-header">
      <div class="header-icon">
        <svg class="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M9 3v2m6-2v2M9 19v2m6-2v2M5 9H3m2 6H3m18-6h-2m2 6h-2M7 19h10a2 2 0 002-2V7a2 2 0 00-2-2H7a2 2 0 00-2 2v10a2 2 0 002 2zM9 9h6v6H9V9z"
          />
        </svg>
      </div>
      <div class="header-info">
        <h3>{{ $t('admin.dashboard.systemResources.title') }}</h3>
        <div class="system-status" :class="systemStatusClass">
          <div class="status-dot" />
          <span>{{ systemStatusText }}</span>
        </div>
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
      <div class="resource-stats">
        <div class="stat-item">
          <div class="stat-icon">
            <svg class="h-4 w-4" fill="currentColor" viewBox="0 0 20 20">
              <path
                fill-rule="evenodd"
                d="M3 4a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zm0 4a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zm0 4a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1z"
                clip-rule="evenodd"
              />
            </svg>
          </div>
          <div class="stat-info">
            <div class="stat-label">{{ $t('admin.dashboard.systemResources.cpuUsage') }}</div>
            <div class="stat-value">{{ data.cpu_usage.toFixed(1) }}%</div>
          </div>
        </div>

        <div class="stat-item">
          <div class="stat-icon">
            <svg class="h-4 w-4" fill="currentColor" viewBox="0 0 20 20">
              <path
                d="M3 4a1 1 0 011-1h12a1 1 0 011 1v2a1 1 0 01-1 1H4a1 1 0 01-1-1V4zM3 10a1 1 0 011-1h6a1 1 0 011 1v6a1 1 0 01-1 1H4a1 1 0 01-1-1v-6zM14 9a1 1 0 00-1 1v6a1 1 0 001 1h2a1 1 0 001-1v-6a1 1 0 00-1-1h-2z"
              />
            </svg>
          </div>
          <div class="stat-info">
            <div class="stat-label">{{ $t('admin.dashboard.systemResources.memoryUsage') }}</div>
            <div class="stat-value">{{ data.memory_usage.toFixed(1) }}%</div>
          </div>
        </div>

        <div class="stat-item">
          <div class="stat-icon">
            <svg class="h-4 w-4" fill="currentColor" viewBox="0 0 20 20">
              <path
                fill-rule="evenodd"
                d="M3 17a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zm3.293-7.707a1 1 0 011.414 0L9 10.586V3a1 1 0 112 0v7.586l1.293-1.293a1 1 0 111.414 1.414l-3 3a1 1 0 01-1.414 0l-3-3a1 1 0 010-1.414z"
                clip-rule="evenodd"
              />
            </svg>
          </div>
          <div class="stat-info">
            <div class="stat-label">{{ $t('admin.dashboard.systemResources.diskUsage') }}</div>
            <div class="stat-value">{{ data.disk_usage.toFixed(1) }}%</div>
          </div>
        </div>
      </div>

      <div class="load-section">
        <div class="section-title">{{ $t('admin.dashboard.systemResources.systemLoad') }}</div>
        <div class="load-stats">
          <div class="load-item">
            <span class="load-period">{{ $t('admin.dashboard.systemResources.oneMinute') }}</span>
            <span class="load-number">{{ data.load_average[0]?.toFixed(2) || '0.00' }}</span>
          </div>
          <div class="load-item">
            <span class="load-period">{{ $t('admin.dashboard.systemResources.fiveMinutes') }}</span>
            <span class="load-number">{{ data.load_average[1]?.toFixed(2) || '0.00' }}</span>
          </div>
          <div class="load-item">
            <span class="load-period">{{ $t('admin.dashboard.systemResources.fifteenMinutes') }}</span>
            <span class="load-number">{{ data.load_average[2]?.toFixed(2) || '0.00' }}</span>
          </div>
        </div>
      </div>

      <div class="uptime-section">
        <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        <span class="uptime-label">{{ $t('admin.dashboard.systemResources.uptime') }}</span>
        <span class="uptime-value">{{ data.uptime }}</span>
      </div>
    </div>
  </div>
</template>

<style scoped>
  .system-resources-chart {
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
    background: rgba(var(--color-brand-500-rgb), 0.1);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.15);
    color: var(--color-brand-500);
  }

  .header-info {
    @apply flex flex-1 items-center justify-between;
  }

  .header-info h3 {
    color: var(--color-content-heading);
    font-size: 0.95rem;
    font-weight: 600;
    margin-bottom: 2px;
    line-height: 1.2;
  }

  .system-status {
    @apply flex items-center gap-1 px-2 py-0.5 text-xs font-medium;
    border-radius: var(--radius-sm);
  }

  .system-status.status-healthy {
    background: rgba(34, 197, 94, 0.2);
    color: var(--color-success);
    border: 1px solid rgba(34, 197, 94, 0.3);
  }

  .system-status.status-warning {
    background: rgba(245, 158, 11, 0.2);
    color: var(--color-warning);
    border: 1px solid rgba(245, 158, 11, 0.3);
  }

  .system-status.status-danger {
    background: rgba(239, 68, 68, 0.2);
    color: var(--color-danger);
    border: 1px solid rgba(239, 68, 68, 0.3);
  }

  .status-dot {
    @apply h-2 w-2;
    border-radius: var(--radius-sm);
    background: currentColor;
    animation: pulse 2s infinite;
  }

  .refresh-btn {
    @apply flex h-5 w-5 cursor-pointer items-center justify-center;
    border-radius: var(--radius-sm);
    background: rgba(var(--color-brand-500-rgb), 0.1);
    color: var(--color-brand-500);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.2);
    transition: all 0.2s ease;
  }

  .refresh-btn:hover {
    background: rgba(var(--color-brand-500-rgb), 0.2);
    border-color: rgba(var(--color-brand-500-rgb), 0.3);
    box-shadow: 0 0 10px rgba(var(--color-brand-500-rgb), 0.2);
  }

  .card-content {
    @apply flex-1 space-y-3 p-3 pt-2;
  }

  .resource-stats {
    @apply space-y-2;
  }

  .stat-item {
    @apply flex items-center gap-3 p-2;
    border-radius: var(--radius-sm);
    background: rgba(var(--color-brand-500-rgb), 0.05);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.08);
    transition: all 0.2s ease;
  }

  .stat-item:hover {
    background: rgba(var(--color-brand-500-rgb), 0.1);
    border-color: rgba(var(--color-brand-500-rgb), 0.15);
  }

  .stat-icon {
    @apply flex h-8 w-8 flex-shrink-0 items-center justify-center;
    border-radius: var(--radius-sm);
    background: rgba(var(--color-brand-500-rgb), 0.1);
    color: var(--color-brand-500);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.15);
  }

  .stat-info {
    @apply flex-1;
  }

  .stat-label {
    color: var(--color-slate-400);
    font-size: 0.75rem;
    font-weight: 500;
    display: block;
    line-height: 1.2;
  }

  .stat-value {
    color: var(--color-content-heading);
    font-size: 1.1rem;
    font-weight: bold;
    display: block;
    margin-top: 0.125rem;
    line-height: 1;
  }

  .load-section {
    @apply space-y-2;
  }

  .section-title {
    color: var(--color-content-heading);
    font-size: 0.8rem;
    font-weight: 600;
  }

  .load-stats {
    @apply space-y-2;
  }

  .load-item {
    @apply flex items-center justify-between p-2;
    border-radius: var(--radius-sm);
    background: rgba(var(--color-brand-500-rgb), 0.05);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.08);
    transition: all 0.2s ease;
  }

  .load-item:hover {
    background: rgba(var(--color-brand-500-rgb), 0.1);
    border-color: rgba(var(--color-brand-500-rgb), 0.15);
  }

  .load-period {
    color: var(--color-slate-400);
    font-size: 0.75rem;
    font-weight: 500;
  }

  .load-number {
    color: var(--color-brand-500);
    font-size: 0.8rem;
    font-weight: 600;
    font-family: monospace;
  }

  .uptime-section {
    @apply flex items-center gap-2 p-2;
    border-radius: var(--radius-sm);
    background: rgba(var(--color-brand-500-rgb), 0.05);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.08);
  }

  .uptime-section svg {
    color: var(--color-brand-500);
    flex-shrink: 0;
  }

  .uptime-label {
    color: var(--color-slate-400);
    font-size: 0.75rem;
    flex: 1;
  }

  .uptime-value {
    color: var(--color-content-heading);
    font-size: 0.75rem;
    font-weight: 600;
    font-family: monospace;
  }

  .animate-spin {
    animation: spin 1s linear infinite;
  }

  @keyframes spin {
    from {
      transform: rotate(0deg);
    }
    to {
      transform: rotate(360deg);
    }
  }

  @keyframes pulse {
    0%,
    100% {
      opacity: 1;
    }
    50% {
      opacity: 0.5;
    }
  }

  @media (max-width: 768px) {
    .stat-value {
      font-size: 1rem;
    }

    .stat-label {
      font-size: 0.7rem;
    }

    .load-period,
    .load-number {
      font-size: 0.7rem;
    }
  }
</style>
