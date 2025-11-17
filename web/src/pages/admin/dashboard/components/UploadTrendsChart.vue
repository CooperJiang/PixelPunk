<script setup lang="ts">
  import { computed, nextTick, onMounted, onUnmounted, ref, watch } from 'vue'
  import { getDashboardUploadTrends } from '@/api/admin/dashboard'
  import { useTheme } from '@/composables/useTheme'
  import * as echarts from 'echarts'
  import { useTexts } from '@/composables/useTexts'
  const { $t } = useTexts()

  defineOptions({
    name: 'UploadTrendsChart',
  })

  /* 主题监听 */
  const { appliedTheme } = useTheme()

  const selectedDays = ref('15')
  const isLoading = ref(false)
  const chartContainer = ref<HTMLDivElement>()
  const chartData = ref<{
    timeline: string[]
    uploadCounts: number[]
    storageGrowth: number[]
  }>({
    timeline: [],
    uploadCounts: [],
    storageGrowth: [],
  })

  /* ECharts 实例 */
  let chartInstance: echarts.ECharts | null = null

  /* 时间选项 */
  const timeOptions = computed(() => [
    { value: '7', label: $t('admin.dashboard.uploadTrends.last7Days') },
    { value: '15', label: $t('admin.dashboard.uploadTrends.last15Days') },
    { value: '30', label: $t('admin.dashboard.uploadTrends.last30Days') },
  ])

  const hasData = computed(() => chartData.value.timeline.length > 0 && chartData.value.uploadCounts.length > 0)

  const totalUploads = computed(() => chartData.value.uploadCounts.reduce((sum, count) => sum + count, 0))

  const avgDaily = computed(() => {
    const counts = chartData.value.uploadCounts
    return counts.length > 0 ? Math.round(totalUploads.value / counts.length) : 0
  })

  const trendPercentage = computed(() => {
    const counts = chartData.value.uploadCounts
    if (counts.length < 2) {
      return 0
    }

    const recent = counts.slice(-3).reduce((sum, count) => sum + count, 0) / 3
    const earlier = counts.slice(0, 3).reduce((sum, count) => sum + count, 0) / 3
    return earlier === 0 ? 0 : Math.round(((recent - earlier) / earlier) * 100)
  })

  const totalStorageGrowth = computed(() => {
    return chartData.value.storageGrowth.reduce((sum, size) => sum + size, 0)
  })

  const formattedStorageGrowth = computed(() => {
    const bytes = totalStorageGrowth.value
    if (bytes === 0 || isNaN(bytes) || !isFinite(bytes)) return '0 B'
    const k = 1024
    const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
    const i = Math.floor(Math.log(bytes) / Math.log(k))
    if (isNaN(i) || !isFinite(i) || i < 0) return '0 B'
    const sizeIndex = Math.min(i, sizes.length - 1)
    return `${(bytes / Math.pow(k, sizeIndex)).toFixed(2)} ${sizes[sizeIndex]}`
  })

  const avgDailyStorage = computed(() => {
    const growth = chartData.value.storageGrowth
    if (growth.length === 0) return '0 B'
    const avgBytes = totalStorageGrowth.value / growth.length
    if (avgBytes === 0 || isNaN(avgBytes) || !isFinite(avgBytes)) return '0 B'
    const k = 1024
    const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
    const i = Math.floor(Math.log(avgBytes) / Math.log(k))
    if (isNaN(i) || !isFinite(i) || i < 0) return '0 B'
    return `${(avgBytes / Math.pow(k, i)).toFixed(2)} ${sizes[i]}`
  })

  const fetchData = async () => {
    isLoading.value = true

    try {
      const days = parseInt(selectedDays.value)

      const result = await getDashboardUploadTrends(days)

      if (result?.success && result.data) {
        chartData.value = {
          timeline: result.data.timeline || [],
          uploadCounts: result.data.upload_counts || [],
          storageGrowth: result.data.storage_growth || [],
        }

        await nextTick()
        setTimeout(() => {
          renderChart()
        }, 100)
      } else {
        chartData.value = { timeline: [], uploadCounts: [], storageGrowth: [] }
      }
    } catch {
      chartData.value = { timeline: [], uploadCounts: [], storageGrowth: [] }
    } finally {
      isLoading.value = false
    }
  }

  const renderChart = () => {
    if (!chartContainer.value || !hasData.value) {
      return
    }

    if (chartInstance) {
      chartInstance.dispose()
      chartInstance = null
    }

    chartInstance = echarts.init(chartContainer.value)

    const rootStyles = getComputedStyle(document.documentElement)
    const contentHeading = rootStyles.getPropertyValue('--color-content-heading').trim()
    const contentDefault = rootStyles.getPropertyValue('--color-content-default').trim()
    const contentMuted = rootStyles.getPropertyValue('--color-content-muted').trim()
    const bg700 = rootStyles.getPropertyValue('--color-background-700').trim()
    const bg700Rgb = rootStyles.getPropertyValue('--color-background-700-rgb').trim()
    const brand500 = rootStyles.getPropertyValue('--color-brand-500').trim()
    const brand500Rgb = rootStyles.getPropertyValue('--color-brand-500-rgb').trim()
    const success500 = brand500 // rootStyles.getPropertyValue('--color-success-500').trim()
    const successRgb = brand500Rgb // rootStyles.getPropertyValue('--color-success-rgb').trim()

    const formatBytes = (bytes: number): string => {
      if (bytes === 0 || isNaN(bytes) || !isFinite(bytes) || bytes == null) return '0 B'
      const k = 1024
      const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
      const i = Math.floor(Math.log(bytes) / Math.log(k))
      if (isNaN(i) || !isFinite(i) || i < 0) return '0 B'
      const sizeIndex = Math.min(i, sizes.length - 1)
      return `${(bytes / Math.pow(k, sizeIndex)).toFixed(2)} ${sizes[sizeIndex]}`
    }

    let effectiveLength = chartData.value.uploadCounts.length
    let hasNonZeroData = false

    for (let i = chartData.value.uploadCounts.length - 1; i >= 0; i--) {
      if (chartData.value.uploadCounts[i] > 0) {
        effectiveLength = i + 1
        hasNonZeroData = true
        break
      }
    }

    if (!hasNonZeroData) {
      effectiveLength = chartData.value.uploadCounts.length
    }

    const effectiveTimeline = chartData.value.timeline.slice(0, effectiveLength)
    const effectiveCounts = chartData.value.uploadCounts.slice(0, effectiveLength)
    const effectiveStorage = chartData.value.storageGrowth.slice(0, effectiveLength)

    const xAxisData = effectiveTimeline.map((dateStr) => {
      const date = new Date(dateStr)
      return `${date.getMonth() + 1}/${date.getDate()}`
    })

    const option: echarts.EChartsOption = {
      backgroundColor: 'transparent',
      grid: {
        left: '1%',
        right: '1%',
        bottom: '3%',
        top: '10%',
        containLabel: true,
      },
      tooltip: {
        trigger: 'axis',
        backgroundColor: `rgba(${bg700Rgb}, 0.95)`,
        borderColor: brand500,
        borderWidth: 1,
        borderRadius: 8,
        textStyle: {
          color: contentHeading,
          fontSize: 12,
        },
        extraCssText: `backdrop-filter: blur(10px); box-shadow: 0 8px 32px rgba(${brand500Rgb}, 0.3), 0 0 16px rgba(${brand500Rgb}, 0.1);`,
        formatter(params: any) {
          if (!params || params.length === 0) {
            return ''
          }
          const uploadParam = params[0]
          const storageParam = params[1]
          const { name } = uploadParam
          return `
          <div class="radius-sm" style="padding: 12px;">
            <div style="color: ${brand500}; font-weight: 600; margin-bottom: 10px; font-size: 14px; text-shadow: 0 0 8px rgba(${brand500Rgb}, 0.5);">
              📅 ${name}
            </div>
            <div style="display: flex; flex-direction: column; gap: 6px;">
              <div style="display: flex; align-items: center; gap: 8px;">
                <span style="color: ${contentDefault}; font-size: 12px;">${$t('admin.dashboard.uploadTrends.tooltip.uploadCountLabel')}</span>
                <span style="color: ${brand500}; font-weight: 700; font-size: 15px; text-shadow: 0 0 8px rgba(${brand500Rgb}, 0.5);">${uploadParam.value}</span>
                <span style="color: ${contentMuted}; font-size: 11px; font-weight: 500;">${$t('admin.dashboard.uploadTrends.filesUnit')}</span>
              </div>
              <div style="display: flex; align-items: center; gap: 8px;">
                <span style="color: ${contentDefault}; font-size: 12px;">${$t('admin.dashboard.uploadTrends.tooltip.storageGrowthLabel')}</span>
                <span style="color: ${success500}; font-weight: 700; font-size: 15px; text-shadow: 0 0 8px rgba(${successRgb}, 0.5);">${formatBytes(storageParam.value)}</span>
              </div>
            </div>
          </div>
        `
        },
      },
      xAxis: {
        type: 'category',
        data: xAxisData,
        boundaryGap: false,
        axisLine: {
          show: false,
        },
        axisTick: {
          show: false,
        },
        axisLabel: {
          color: contentMuted,
          fontSize: 11,
          margin: 10,
        },
        splitLine: {
          show: true,
          lineStyle: {
            color: `rgba(${brand500Rgb}, 0.15)`,
            type: 'dashed',
          },
        },
      },
      yAxis: [
        {
          type: 'value',
          name: $t('admin.dashboard.uploadTrends.uploadCount'),
          nameTextStyle: {
            color: brand500,
            fontSize: 11,
            padding: [0, 0, 0, 10],
          },
          position: 'left',
          axisLine: {
            show: false,
          },
          axisTick: {
            show: false,
          },
          axisLabel: {
            color: contentMuted,
            fontSize: 11,
            formatter: '{value}',
          },
          splitLine: {
            show: true,
            lineStyle: {
              color: `rgba(${brand500Rgb}, 0.15)`,
              type: 'dashed',
            },
          },
        },
        {
          type: 'value',
          name: $t('admin.dashboard.uploadTrends.storageGrowth'),
          nameTextStyle: {
            color: brand500,
            fontSize: 11,
            padding: [0, 10, 0, 0],
          },
          position: 'right',
          axisLine: {
            show: false,
          },
          axisTick: {
            show: false,
          },
          axisLabel: {
            color: contentMuted,
            fontSize: 11,
            formatter: (value: number) => formatBytes(value),
          },
          splitLine: {
            show: false,
          },
        },
      ],
      series: [
        {
          name: $t('admin.dashboard.uploadTrends.uploadCount'),
          type: 'line',
          yAxisIndex: 0,
          data: effectiveCounts,
          smooth: true,
          symbol: 'circle',
          symbolSize: 8,
          showSymbol: true,
          lineStyle: {
            color: brand500,
            width: 3,
            shadowColor: `rgba(${brand500Rgb}, 0.5)`,
            shadowBlur: 10,
          },
          itemStyle: {
            color: brand500,
            borderColor: bg700,
            borderWidth: 2,
          },
          areaStyle: {
            color: {
              type: 'linear',
              x: 0,
              y: 0,
              x2: 0,
              y2: 1,
              colorStops: [
                {
                  offset: 0,
                  color: `rgba(${brand500Rgb}, 0.3)`,
                },
                {
                  offset: 1,
                  color: `rgba(${brand500Rgb}, 0.05)`,
                },
              ],
            },
          },
          emphasis: {
            scale: true,
            focus: 'series',
          },
        },
        {
          name: $t('admin.dashboard.uploadTrends.storageGrowth'),
          type: 'line',
          yAxisIndex: 1,
          data: effectiveStorage,
          smooth: true,
          symbol: 'circle',
          symbolSize: 8,
          showSymbol: true,
          lineStyle: {
            color: success500,
            width: 3,
            shadowColor: `rgba(${successRgb}, 0.5)`,
            shadowBlur: 10,
          },
          itemStyle: {
            color: success500,
            borderColor: bg700,
            borderWidth: 2,
          },
          areaStyle: {
            color: {
              type: 'linear',
              x: 0,
              y: 0,
              x2: 0,
              y2: 1,
              colorStops: [
                {
                  offset: 0,
                  color: `rgba(${successRgb}, 0.25)`,
                },
                {
                  offset: 1,
                  color: `rgba(${successRgb}, 0.03)`,
                },
              ],
            },
          },
          emphasis: {
            scale: true,
            focus: 'series',
          },
        },
      ],
    }

    chartInstance.setOption(option)
  }

  const handleResize = () => {
    if (chartInstance) {
      chartInstance.resize()
    }
  }

  watch(appliedTheme, (_newTheme) => {
    if (hasData.value && chartContainer.value) {
      setTimeout(() => {
        renderChart()
      }, 50)
    }
  })

  onMounted(() => {
    fetchData()
    window.addEventListener('resize', handleResize)
  })

  onUnmounted(() => {
    window.removeEventListener('resize', handleResize)
    if (chartInstance) {
      chartInstance.dispose()
      chartInstance = null
    }
  })
</script>

<template>
  <div class="chart-card upload-trends">
    <div class="card-header">
      <div class="header-icon">
        <svg class="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"
          />
        </svg>
      </div>
      <div class="header-info">
        <h3>{{ $t('admin.dashboard.uploadTrends.title') }}</h3>
        <div class="time-selector">
          <CyberDropdown
            v-model="selectedDays"
            :options="timeOptions"
            :clearable="false"
            :searchable="false"
            :placeholder="$t('admin.dashboard.uploadTrends.selectTimeRange')"
            @change="fetchData"
          />
        </div>
      </div>
    </div>

    <div class="chart-content">
      <div v-if="isLoading" class="chart-loading">
        <div class="loading-spinner" />
        <span>{{ $t('admin.dashboard.uploadTrends.loading') }}</span>
      </div>

      <div v-else-if="hasData" class="chart-wrapper">
        <div ref="chartContainer" class="chart-canvas" />
      </div>

      <div v-else-if="!isLoading && !hasData" class="empty-state">
        <div class="empty-icon">📈</div>
        <div class="empty-title">{{ $t('admin.dashboard.uploadTrends.noData') }}</div>
        <div class="empty-subtitle">{{ $t('admin.dashboard.uploadTrends.tryOtherRange') }}</div>
      </div>

      <div v-if="hasData && !isLoading" class="chart-summary">
        <div class="summary-item">
          <div class="summary-value">{{ totalUploads }}</div>
          <div class="summary-label">{{ $t('admin.dashboard.uploadTrends.totalUploads') }}</div>
        </div>
        <div class="summary-item">
          <div class="summary-value">{{ avgDaily }}</div>
          <div class="summary-label">{{ $t('admin.dashboard.uploadTrends.avgDaily') }}</div>
        </div>
        <div class="summary-item storage">
          <div class="summary-value">{{ formattedStorageGrowth }}</div>
          <div class="summary-label">{{ $t('admin.dashboard.uploadTrends.totalStorageGrowth') }}</div>
        </div>
        <div class="summary-item storage">
          <div class="summary-value">{{ avgDailyStorage }}</div>
          <div class="summary-label">{{ $t('admin.dashboard.uploadTrends.avgDailyStorage') }}</div>
        </div>
        <div class="summary-item">
          <div class="summary-value" :class="{ positive: trendPercentage > 0, negative: trendPercentage < 0 }">
            {{ trendPercentage > 0 ? '+' : '' }}{{ trendPercentage }}%
          </div>
          <div class="summary-label">{{ $t('admin.dashboard.uploadTrends.trend') }}</div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
  .chart-card {
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
    overflow: visible;
  }

  .card-header {
    @apply flex items-center gap-3 p-4 pb-3;
    border-bottom: 1px solid rgba(var(--color-brand-500-rgb), 0.1);
  }

  .header-icon {
    @apply flex h-12 w-12 flex-shrink-0 items-center justify-center;
    border-radius: var(--radius-sm);
    background: linear-gradient(135deg, rgba(var(--color-brand-500-rgb), 0.15) 0%, rgba(124, 77, 255, 0.15) 100%);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.25);
    color: var(--color-brand-500);
    box-shadow: 0 0 15px rgba(var(--color-brand-500-rgb), 0.2);
  }

  .header-info {
    @apply flex flex-1 items-center justify-between;
  }

  .header-info h3 {
    color: var(--color-content-heading);
    font-size: 1.1rem;
    font-weight: 600;
  }

  .time-selector {
    min-width: 100px;
  }

  .time-selector :deep(.cyber-dropdown-header) {
    height: 28px;
    min-height: 28px;
    font-size: 0.8rem;
    padding: 0 0.5rem;
  }

  .time-selector :deep(.cyber-dropdown-selected) {
    font-size: 0.8rem;
  }

  .chart-content {
    @apply flex-1 p-4;
    position: relative;
    overflow: visible;
  }

  .chart-container {
    position: relative;
    overflow: visible;
  }

  .chart-wrapper {
    @apply mb-2;
    height: 340px;
    position: relative;
    overflow: visible;
  }

  .chart-canvas {
    @apply h-full w-full;
    position: relative;
    overflow: visible;
  }

  .chart-summary {
    @apply flex justify-around py-2;
    border-top: 1px solid rgba(var(--color-brand-500-rgb), 0.1);
  }

  .summary-item {
    @apply text-center;
  }

  .summary-value {
    color: var(--color-brand-500);
    font-size: 1.6rem;
    font-weight: bold;
    text-shadow: 0 0 12px rgba(var(--color-brand-500-rgb), 0.5);
    margin-bottom: 0.15rem;
  }

  .summary-item.storage .summary-value {
    color: var(--color-brand-500);
    text-shadow: 0 0 12px rgba(var(--color-brand-500-rgb), 0.5);
  }

  .summary-value.positive {
    color: var(--color-success-500);
    text-shadow: 0 0 12px rgba(var(--color-success-rgb), 0.5);
  }

  .summary-value.negative {
    color: var(--color-error-500);
    text-shadow: 0 0 12px rgba(var(--color-error-rgb), 0.5);
  }

  .summary-label {
    color: var(--color-content-muted) !important;
    font-size: 0.8rem;
    font-weight: 500;
  }

  .chart-loading {
    @apply flex flex-col items-center justify-center;
    height: 340px;
    gap: 1rem;
  }

  .loading-spinner {
    @apply h-8 w-8 animate-spin;
    border-radius: var(--radius-full);
    border: 4px solid rgba(var(--color-brand-500-rgb), 0.2);
    border-top: 4px solid var(--color-brand-500);
    box-shadow: 0 0 20px rgba(var(--color-brand-500-rgb), 0.3);
  }

  .chart-loading span {
    color: var(--color-content-muted) !important;
    font-size: 0.9rem;
  }

  .empty-state {
    @apply flex flex-col items-center justify-center gap-2;
    height: 340px;
    border-top: 1px solid rgba(var(--color-brand-500-rgb), 0.1);
  }

  .empty-icon {
    font-size: 1.75rem;
  }

  .empty-title {
    color: var(--color-info-200) !important;
    font-weight: 600;
  }

  .empty-subtitle {
    color: var(--color-content-muted) !important;
    font-size: 0.85rem;
  }

  @media (max-width: 768px) {
    .chart-wrapper {
      height: 300px;
    }

    .summary-value {
      font-size: 1.25rem;
    }

    .chart-summary {
      @apply flex-col gap-2;
    }

    .summary-item {
      @apply flex items-center justify-between;
    }
  }
</style>
