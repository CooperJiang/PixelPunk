<script setup lang="ts">
  import { computed, onMounted, ref } from 'vue'
  import { getDashboardTagStats, type DashboardTagStatsResponse } from '@/api/admin/dashboard'
  import { useTexts } from '@/composables/useTexts'
  const { $t } = useTexts()

  const loading = ref(false)

  const data = ref<DashboardTagStatsResponse>({
    ai_tags_count: 0,
    tagged_images: 0,
    untagged_images: 0,
    popular_tags: [],
    manual_tags_count: 0,
  })

  function formatNumber(num: number | undefined | null): string {
    if (num === null || num === undefined || Number.isNaN(num)) {
      return '0'
    }
    if (num >= 1000000) {
      return `${(num / 1000000).toFixed(1)}M`
    } else if (num >= 1000) {
      return `${(num / 1000).toFixed(1)}K`
    }
    return num.toString()
  }

  function getTagGridClass(index: number): string {
    if (index === 0) {
      return 'tag-gold'
    }
    if (index === 1) {
      return 'tag-silver'
    }
    if (index === 2) {
      return 'tag-bronze'
    }
    return 'tag-normal'
  }

  async function fetchData() {
    try {
      loading.value = true
      const result = await getDashboardTagStats()

      if (result.success && result.data) {
        data.value = result.data
      }
    } catch {
    } finally {
      loading.value = false
    }
  }

  const hasData = computed(
    () =>
      data.value.ai_tags_count > 0 ||
      data.value.tagged_images > 0 ||
      data.value.untagged_images > 0 ||
      (data.value.popular_tags && data.value.popular_tags.length > 0)
  )

  function refresh() {
    fetchData()
  }

  onMounted(() => {
    fetchData()
  })
</script>

<template>
  <div class="tag-stats-card">
    <div class="card-header">
      <div class="header-icon">
        <svg class="h-6 w-6" fill="currentColor" viewBox="0 0 20 20">
          <path
            fill-rule="evenodd"
            d="M17.707 9.293a1 1 0 010 1.414l-7 7a1 1 0 01-1.414 0l-7-7A.997.997 0 012 10V5a3 3 0 013-3h5c.256 0 .512.098.707.293l7 7zM5 6a1 1 0 100-2 1 1 0 000 2z"
            clip-rule="evenodd"
          />
        </svg>
      </div>
      <h3>{{ $t('admin.dashboard.tagStats.title') }}</h3>
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
      <div v-if="hasData" class="stats-wrapper">
        <div class="overview-stats">
          <div class="stat-item primary">
            <div class="stat-info">
              <div class="stat-value">{{ formatNumber(data.ai_tags_count) }}</div>
              <div class="stat-label">{{ $t('admin.dashboard.tagStats.aiTags') }}</div>
            </div>
            <div class="stat-icon"><i class="fas fa-robot" /></div>
          </div>

          <div class="stat-item success">
            <div class="stat-info">
              <div class="stat-value">{{ formatNumber(data.tagged_images) }}</div>
              <div class="stat-label">{{ $t('admin.dashboard.tagStats.taggedFiles') }}</div>
            </div>
            <div class="stat-icon"><i class="fas fa-check-circle" /></div>
          </div>

          <div class="stat-item warning">
            <div class="stat-info">
              <div class="stat-value">{{ formatNumber(data.untagged_images) }}</div>
              <div class="stat-label">{{ $t('admin.dashboard.tagStats.untaggedFiles') }}</div>
            </div>
            <div class="stat-icon"><i class="fas fa-exclamation-triangle" /></div>
          </div>
        </div>

        <div class="tags-grid">
          <div
            v-for="(tag, index) in data.popular_tags.slice(0, 12)"
            :key="tag.name"
            class="stat-item tag-item"
            :class="getTagGridClass(index)"
          >
            <div class="stat-info">
              <div class="stat-value">{{ formatNumber(tag.count) }}</div>
              <div class="stat-label">{{ tag.name }}</div>
            </div>
            <div class="tag-rank-badge">{{ index + 1 }}</div>
          </div>

          <div v-for="i in Math.max(0, 12 - data.popular_tags.length)" :key="`empty-tag-${i}`" class="stat-item empty-item">
            <div class="stat-info">
              <div class="stat-value">-</div>
              <div class="stat-label">{{ $t('admin.dashboard.tagStats.none') }}</div>
            </div>
          </div>
        </div>
      </div>
      <div v-else class="empty-state">
        <div class="empty-icon">🏷️</div>
        <div class="empty-title">{{ $t('admin.dashboard.tagStats.noData') }}</div>
        <div class="empty-subtitle">{{ $t('admin.dashboard.tagStats.waitingForData') }}</div>
      </div>
    </div>
  </div>
</template>

<style scoped>
  .tag-stats-card {
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

  .card-header h3 {
    color: var(--color-content-heading);
    font-size: 0.95rem;
    font-weight: 600;
    margin-bottom: 2px;
    line-height: 1.2;
    flex: 1;
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
    box-shadow: 0 0 10px rgba(var(--color-brand-500-rgb), 0.2);
  }

  .card-content {
    @apply flex-1 p-3 pt-2;
  }

  .stats-wrapper {
    background: linear-gradient(
      135deg,
      rgba(var(--color-brand-500-rgb), 0.02) 0%,
      rgba(var(--color-background-800-rgb), 0.3) 100%
    );
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.08);
    border-radius: var(--radius-sm);
    padding: 1rem;
  }

  .empty-state {
    @apply flex flex-col items-center justify-center gap-2 py-10;
  }

  .empty-icon {
    font-size: 1.5rem;
  }
  .empty-title {
    color: var(--color-slate-200);
    font-weight: 600;
  }
  .empty-subtitle {
    color: var(--color-slate-400);
    font-size: 0.85rem;
  }

  .overview-stats {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 0.75rem;
    align-items: stretch;
    justify-items: stretch;
    margin-bottom: 0.75rem;
  }

  .overview-stats .stat-item {
    height: 68px !important;
    min-height: 68px !important;
    max-height: 68px !important;
    display: flex !important;
    flex-direction: row !important;
    align-items: center !important;
    justify-content: flex-start !important;
    text-align: left !important;
    padding: 0.5rem 0.75rem !important;
  }

  .overview-stats .stat-info {
    @apply flex flex-col items-start justify-center;
  }

  .overview-stats .stat-icon {
    @apply flex h-6 w-6 items-center justify-center;
    border-radius: var(--radius-full);
    margin-top: 0 !important;
    margin-right: 0.5rem !important;
    flex-shrink: 0;
    order: -1;

    background: rgba(148, 163, 184, 0.15);
    border: 1px solid rgba(148, 163, 184, 0.25);
  }

  .overview-stats .stat-item.primary .stat-icon {
    background: rgba(var(--color-brand-500-rgb), 0.15);
    border-color: rgba(var(--color-brand-500-rgb), 0.25);
  }

  .overview-stats .stat-item.success .stat-icon {
    background: rgba(34, 197, 94, 0.15);
    border-color: rgba(34, 197, 94, 0.25);
  }

  .overview-stats .stat-item.warning .stat-icon {
    background: rgba(245, 158, 11, 0.15);
    border-color: rgba(245, 158, 11, 0.25);
  }

  .tags-grid {
    display: grid;
    grid-template-columns: repeat(6, 1fr);
    gap: 0.75rem;
    align-items: stretch;
    justify-items: stretch;
  }

  .stat-item {
    @apply text-center;
    border-radius: var(--radius-sm);
    background: rgba(var(--color-brand-500-rgb), 0.04);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.08);
    transition: all 0.3s ease;
    height: 65px !important;
    min-height: 65px !important;
    max-height: 65px !important;
    display: flex !important;
    flex-direction: column !important;
    align-items: center !important;
    justify-content: center !important;
    box-sizing: border-box !important;
    padding: 0.5rem 0.75rem !important;
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
    @apply flex flex-col items-center justify-center;
    flex: 1;
    min-height: 0;
  }

  .stat-icon {
    @apply flex h-6 w-6 items-center justify-center rounded;
    font-size: 0.8rem;
    flex-shrink: 0;
    min-height: 24px;
    max-height: 24px;
    margin-top: 0.5rem;
  }

  .stat-item.primary .stat-value,
  .stat-item.primary .stat-icon {
    color: var(--color-brand-500);
  }

  .stat-item.success .stat-value,
  .stat-item.success .stat-icon {
    color: var(--color-success);
  }

  .stat-item.warning .stat-value,
  .stat-item.warning .stat-icon {
    color: var(--color-warning);
  }

  .stat-value {
    font-size: 1.1rem;
    font-weight: 700;
    line-height: 1;
    margin-bottom: 0.25rem;
  }

  .stat-label {
    color: var(--color-slate-400);
    font-size: 0.7rem;
    font-weight: 500;
    line-height: 1;
    margin: 0;
    padding: 0;
  }

  .stat-item.tag-item {
    position: relative;
  }

  .tag-rank-badge {
    position: absolute;
    top: 0.375rem;
    right: 0.375rem;
    width: 1.25rem;
    height: 1.25rem;
    border-radius: var(--radius-full);
    background: rgba(var(--color-brand-500-rgb), 0.2);
    color: var(--color-brand-500);
    font-size: 0.65rem;
    font-weight: 700;
    display: flex;
    align-items: center;
    justify-content: center;
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.3);
    z-index: 1;
  }

  .stat-item.tag-gold .tag-rank-badge {
    background: linear-gradient(135deg, var(--color-warning), #d97706);
    color: var(--color-black);
    border-color: var(--color-warning);
    box-shadow: 0 0 8px rgba(245, 158, 11, 0.3);
  }

  .stat-item.tag-silver .tag-rank-badge {
    background: linear-gradient(135deg, var(--color-cyan-500), #2563eb);
    color: var(--color-content-heading);
    border-color: var(--color-cyan-500);
    box-shadow: 0 0 8px rgba(59, 130, 246, 0.3);
  }

  .stat-item.tag-bronze .tag-rank-badge {
    background: linear-gradient(135deg, var(--color-brand-500), #9333ea);
    color: var(--color-content-heading);
    border-color: var(--color-brand-500);
    box-shadow: 0 0 8px rgba(168, 85, 247, 0.3);
  }

  .stat-item.tag-gold {
    background: linear-gradient(135deg, rgba(245, 158, 11, 0.08), rgba(var(--color-brand-500-rgb), 0.04));
    border-color: rgba(245, 158, 11, 0.2);
  }

  .stat-item.tag-silver {
    background: linear-gradient(135deg, rgba(59, 130, 246, 0.08), rgba(var(--color-brand-500-rgb), 0.04));
    border-color: rgba(59, 130, 246, 0.2);
  }

  .stat-item.tag-bronze {
    background: linear-gradient(135deg, rgba(168, 85, 247, 0.08), rgba(var(--color-brand-500-rgb), 0.04));
    border-color: rgba(168, 85, 247, 0.2);
  }

  .stat-item.empty-item {
    background: rgba(0, 0, 0, 0.1);
    border: 1px dashed rgba(var(--color-brand-500-rgb), 0.2);
  }

  .stat-item.empty-item .stat-value,
  .stat-item.empty-item .stat-label {
    color: var(--color-content-disabled);
  }

  .stat-item.empty-item:hover {
    transform: none;
    box-shadow: none;
    background: rgba(0, 0, 0, 0.15);
  }

  .no-tags {
    @apply flex items-center justify-center gap-2 py-8 text-content-muted;
  }

  @media (max-width: 1200px) {
    .tags-grid {
      grid-template-columns: repeat(4, 1fr);
      gap: 0.5rem;
    }
  }

  @media (max-width: 768px) {
    .overview-stats {
      grid-template-columns: 1fr;
      gap: 0.5rem;
      margin-bottom: 0.5rem;
    }

    .tags-grid {
      grid-template-columns: repeat(3, 1fr);
      gap: 0.5rem;
    }

    .stat-item {
      min-height: 55px;
      padding: 0.5rem;
    }

    .stat-value {
      font-size: 0.9rem;
    }

    .stat-label {
      font-size: 0.65rem;
    }

    .tag-rank-badge {
      width: 1rem;
      height: 1rem;
      font-size: 0.6rem;
    }
  }
</style>
