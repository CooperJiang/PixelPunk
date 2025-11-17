<script setup lang="ts">
  import { computed, onMounted, ref } from 'vue'
  import { getDashboardImageStats } from '@/api/admin/dashboard'
  import { useTexts } from '@/composables/useTexts'
  import DashboardCard from './DashboardCard.vue'

  defineOptions({
    name: 'ImageStatsCard',
  })

  const { $t } = useTexts()

  /* 数据状态 */
  const totalImages = ref(0)
  const uploadedToday = ref(0)
  const aiTaggedPercentage = ref(0)
  const pendingReview = ref(0)
  const nsfwDetected = ref(0)
  const untaggedCount = ref(0)
  const lastUpdated = ref('')
  const isRefreshing = ref(false)

  const safeAiTaggedPercentage = computed(() => {
    const percentage = aiTaggedPercentage.value
    if (percentage == null || isNaN(percentage) || !isFinite(percentage)) {
      return 0
    }
    return Math.round(percentage)
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
      const result = await getDashboardImageStats()

      if (result.success && result.data) {
        totalImages.value = result.data.total_files || 0
        uploadedToday.value = result.data.uploaded_today || 0
        aiTaggedPercentage.value = Math.round(result.data.ai_tagged_percentage || 0)
        pendingReview.value = result.data.pending_review || 0
        nsfwDetected.value = result.data.nsfw_detected || 0
        untaggedCount.value = result.data.untagged_count || 0

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
    :title="$t('admin.dashboard.fileStats.title')"
    :badge="`AI ${safeAiTaggedPercentage}%`"
    :badge-type="safeAiTaggedPercentage >= 80 ? 'success' : safeAiTaggedPercentage >= 60 ? 'warning' : 'danger'"
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
          d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"
        />
      </svg>
    </template>

    <div class="image-content">
      <div class="metric-main">
        <div class="metric-value">
          <span class="value-number">{{ formatNumber(totalImages) }}</span>
          <span class="value-label">{{ $t('admin.dashboard.fileStats.totalFiles') }}</span>
        </div>
        <div class="today-upload">
          <span class="upload-count">+{{ uploadedToday }}</span>
          <span class="upload-label">{{ $t('admin.dashboard.fileStats.uploadedToday') }}</span>
        </div>
      </div>

      <div class="metrics-grid">
        <div class="metric-item">
          <span class="metric-label">{{ $t('admin.dashboard.fileStats.pendingReview') }}</span>
          <div class="metric-value-row">
            <div class="metric-icon pending">
              <svg class="h-4 w-4" fill="currentColor" viewBox="0 0 20 20">
                <path
                  fill-rule="evenodd"
                  d="M10 18a8 8 0 100-16 8 8 0 000 16zm1-12a1 1 0 10-2 0v4a1 1 0 00.293.707l2.828 2.829a1 1 0 101.415-1.415L11 9.586V6z"
                  clip-rule="evenodd"
                />
              </svg>
            </div>
            <span class="metric-number">{{ formatNumber(pendingReview) }}</span>
          </div>
        </div>

        <div class="metric-item">
          <span class="metric-label">{{ $t('admin.dashboard.fileStats.untagged') }}</span>
          <div class="metric-value-row">
            <div class="metric-icon untagged">
              <svg class="h-4 w-4" fill="currentColor" viewBox="0 0 20 20">
                <path
                  fill-rule="evenodd"
                  d="M17.707 9.293a1 1 0 010 1.414l-7 7a1 1 0 01-1.414 0l-7-7A.997.997 0 012 10V5a3 3 0 013-3h5c.256 0 .512.098.707.293l7 7zM5 6a1 1 0 100-2 1 1 0 000 2z"
                  clip-rule="evenodd"
                />
              </svg>
            </div>
            <span class="metric-number">{{ formatNumber(untaggedCount) }}</span>
          </div>
        </div>

        <div class="metric-item">
          <span class="metric-label">NSFW</span>
          <div class="metric-value-row">
            <div class="metric-icon danger">
              <svg class="h-4 w-4" fill="currentColor" viewBox="0 0 20 20">
                <path
                  fill-rule="evenodd"
                  d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z"
                  clip-rule="evenodd"
                />
              </svg>
            </div>
            <span class="metric-number">{{ formatNumber(nsfwDetected) }}</span>
          </div>
        </div>
      </div>
    </div>
  </DashboardCard>
</template>

<style scoped>
  .image-content {
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

  .metric-value {
    @apply flex flex-col;
  }

  .value-number {
    color: var(--color-brand-500);
    font-size: 1.5rem;
    font-weight: var(--font-bold);
    text-shadow: var(--shadow-glow-sm);
    margin-bottom: var(--space-xs);
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
    font-size: var(--text-xs);
    font-weight: var(--font-medium);
  }

  .today-upload {
    @apply flex flex-col items-end;
  }

  .upload-count {
    color: var(--color-success);
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
    background: rgba(var(--color-brand-500-rgb), 0.05);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.08);
    transition: all var(--transition-normal) var(--ease-out);
    min-height: 58px;
  }

  .metric-item:hover {
    background: rgba(var(--color-brand-500-rgb), 0.1);
    border-color: rgba(var(--color-brand-500-rgb), 0.15);
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

  .metric-icon.pending {
    background: rgba(251, 146, 60, 0.2);
    color: var(--color-warning-500);
  }

  .metric-icon.untagged {
    background: rgba(250, 204, 21, 0.2);
    color: var(--color-warning-400);
  }

  .metric-icon.danger {
    background: rgba(239, 68, 68, 0.2);
    color: var(--color-danger);
  }

  .metric-number {
    color: var(--color-content-heading);
    font-size: 0.85rem;
    font-weight: var(--font-bold);
    white-space: nowrap;
  }
</style>
