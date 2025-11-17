<script setup lang="ts">
  import { computed, onUnmounted, ref } from 'vue'
  import { type ActivityLog, getUserActivities } from '@/api/activity'
  import { formatActivityDescription } from '@/utils/activity/formatter'
  import { useTexts } from '@/composables/useTexts'

  defineOptions({
    name: 'ActivityMonitor',
  })

  const { $t } = useTexts()
  const monitorStatus = ref('active')

  const statusText = computed(() => {
    switch (monitorStatus.value) {
      case 'active':
        return $t('dashboard.activityMonitor.status.active')
      case 'warning':
        return $t('dashboard.activityMonitor.status.warning')
      case 'error':
        return $t('dashboard.activityMonitor.status.error')
      default:
        return $t('dashboard.activityMonitor.status.offline')
    }
  })

  const realtimeStats = ref({
    uploadsPerHour: 12,
    viewsPerHour: 847,
  })

  const activities = ref<
    Array<{
      id: number
      type: string
      icon: string
      class: string
      description: string
      timestamp: number
    }>
  >([])

  const currentPage = ref(1)
  const pageSize = ref(15)
  const totalActivities = ref(0)
  const hasMore = ref(true)
  const initialLoading = ref(false)
  const loadingMore = ref(false)

  const scrollContainer = ref<HTMLElement>()

  const getActivityIcon = (type: string): { icon: string; class: string } => {
    switch (type) {
      case 'user_login':
        return { icon: 'fas fa-sign-in-alt', class: 'login' }
      case 'user_register':
        return { icon: 'fas fa-user-plus', class: 'register' }
      case 'batch_image_upload':
        return { icon: 'fas fa-upload', class: 'upload' }
      case 'file_delete':
        return { icon: 'fas fa-file-times', class: 'delete' }
      case 'file_rename':
        return { icon: 'fas fa-file-signature', class: 'edit' }
      case 'file_move':
        return { icon: 'fas fa-file-export', class: 'move' }
      case 'share_create':
        return { icon: 'fas fa-share-alt', class: 'share' }
      case 'share_delete':
        return { icon: 'fas fa-share-square', class: 'delete' }
      case 'share_milestone':
        return { icon: 'fas fa-eye', class: 'view' }
      case 'apikey_create':
        return { icon: 'fas fa-key', class: 'create' }
      case 'apikey_delete':
        return { icon: 'fas fa-key', class: 'delete' }
      case 'apikey_toggle_status':
        return { icon: 'fas fa-toggle-on', class: 'edit' }
      case 'apikey_regenerate':
        return { icon: 'fas fa-sync-alt', class: 'security' }
      case 'random_api_create':
        return { icon: 'fas fa-random', class: 'create' }
      case 'random_api_delete':
        return { icon: 'fas fa-random', class: 'delete' }
      case 'random_api_toggle_status':
        return { icon: 'fas fa-toggle-on', class: 'edit' }
      case 'batch_delete':
      case 'admin_delete':
        return { icon: 'fas fa-trash', class: 'delete' }
      case 'system_cleanup':
        return { icon: 'fas fa-cog', class: 'system' }
      case 'folder_create':
        return { icon: 'fas fa-folder-plus', class: 'folder' }
      case 'folder_rename':
        return { icon: 'fas fa-edit', class: 'folder' }
      case 'folder_delete':
        return { icon: 'fas fa-folder-minus', class: 'delete' }
      case 'folder_access_level_change':
        return { icon: 'fas fa-lock', class: 'permission' }
      case 'image_access_level_change':
        return { icon: 'fas fa-eye-slash', class: 'permission' }
      case 'profile_update':
        return { icon: 'fas fa-user-edit', class: 'profile' }
      case 'password_change':
        return { icon: 'fas fa-key', class: 'security' }
      case 'email_change':
        return { icon: 'fas fa-envelope', class: 'security' }
      case 'hotlink_protection_change':
        return { icon: 'fas fa-shield-alt', class: 'security' }
      case 'image_expired':
        return { icon: 'fas fa-clock', class: 'expired' }
      case 'guest_image_expired':
        return { icon: 'fas fa-user-clock', class: 'system' }
      default:
        return { icon: 'fas fa-circle', class: 'default' }
    }
  }

  const loadActivities = async (reset = true) => {
    if (initialLoading.value || loadingMore.value) {
      return
    }

    try {
      if (reset) {
        initialLoading.value = true
        currentPage.value = 1
        activities.value = []
        hasMore.value = true
      } else {
        loadingMore.value = true
      }

      monitorStatus.value = 'active'

      const response = await getUserActivities({
        page: currentPage.value,
        size: pageSize.value,
      })

      if (response.code === 200 && response.data) {
        const newActivities = response.data.list.map((activity: ActivityLog) => {
          const iconInfo = getActivityIcon(activity.type)

          return {
            id: activity.id,
            type: activity.type,
            icon: iconInfo.icon,
            class: iconInfo.class,
            description: formatActivityDescription(activity, $t),
            timestamp: new Date(activity.created_at).getTime(),
          }
        })

        if (reset) {
          activities.value = newActivities

          if (response.data.today_stats) {
            realtimeStats.value.uploadsPerHour = response.data.today_stats.today_uploads
            realtimeStats.value.viewsPerHour = response.data.today_stats.total_views
          }
        } else {
          const existingIds = new Set(activities.value.map((a) => a.id))
          const uniqueNewActivities = newActivities.filter((a) => !existingIds.has(a.id))
          activities.value.push(...uniqueNewActivities)
        }

        totalActivities.value = response.data.total || 0
        hasMore.value = response.data.list.length === pageSize.value && activities.value.length < totalActivities.value
      } else {
        console.warn('Failed to fetch activity logs:', response.message)
        monitorStatus.value = 'warning'
      }
    } catch (_error) {
      monitorStatus.value = 'error'

      if (reset) {
        activities.value = [
          {
            id: 1,
            type: 'upload',
            icon: 'fas fa-upload',
            class: 'upload',
            description: $t('dashboard.activityMonitor.toast.noData'),
            timestamp: Date.now(),
          },
        ]
        totalActivities.value = 1
        hasMore.value = false
      }
    } finally {
      initialLoading.value = false
      loadingMore.value = false
    }
  }

  const loadMoreActivities = async () => {
    if (!hasMore.value || loadingMore.value) {
      return
    }

    currentPage.value++
    await loadActivities(false)
  }

  const handleScroll = async (event: Event) => {
    const target = event.target as HTMLElement
    if (!target || !hasMore.value || loadingMore.value) {
      return
    }

    const { scrollTop } = target
    const { scrollHeight } = target
    const { clientHeight } = target

    if (scrollTop + clientHeight >= scrollHeight - 50) {
      await loadMoreActivities()
    }
  }

  const formatTime = (timestamp: number): string => {
    const now = Date.now()
    const diff = now - timestamp
    const minutes = Math.floor(diff / (1000 * 60))
    const hours = Math.floor(diff / (1000 * 60 * 60))

    if (hours > 0) {
      return $t('dashboard.activityMonitor.timeAgo.hoursAgo', { hours })
    }
    if (minutes > 0) {
      return $t('dashboard.activityMonitor.timeAgo.minutesAgo', { minutes })
    }
    return $t('dashboard.activityMonitor.timeAgo.justNow')
  }

  onMounted(() => {
    loadActivities(true)

    const activityInterval = setInterval(() => {
      if (!initialLoading.value && !loadingMore.value) {
        loadActivities(true) // 这会同时刷新今日统计数据
      }
    }, 120000) // 2分钟

    onUnmounted(() => {
      clearInterval(activityInterval)
    })
  })
</script>

<template>
  <div class="activity-monitor cyber-card">
    <div class="monitor-content">
      <div class="monitor-header">
        <h3 class="section-title">
          <i class="fas fa-chart-line" />
          {{ $t('dashboard.activityMonitor.title') }}
        </h3>
        <div class="monitor-status">
          <div class="status-dot" :class="monitorStatus" />
          <span class="status-text">{{ statusText }}</span>
        </div>
      </div>

      <div class="stats-overview">
        <div class="stat-item">
          <div class="stat-icon upload">
            <i class="fas fa-upload" />
          </div>
          <div class="stat-content">
            <div class="stat-value">{{ realtimeStats.uploadsPerHour }}</div>
            <div class="stat-label">{{ $t('dashboard.activityMonitor.todayUploads') }}</div>
          </div>
        </div>

        <div class="stat-item">
          <div class="stat-icon view">
            <i class="fas fa-eye" />
          </div>
          <div class="stat-content">
            <div class="stat-value">{{ realtimeStats.viewsPerHour }}</div>
            <div class="stat-label">{{ $t('dashboard.activityMonitor.totalViews') }}</div>
          </div>
        </div>
      </div>

      <div class="recent-activities">
        <div class="activities-header">
          <h4 class="section-subtitle">{{ $t('dashboard.activityMonitor.recentActivity') }}</h4>
          <div v-if="activities.length > 0" class="activities-counter">
            {{ $t('dashboard.activityMonitor.recordCount', { count: totalActivities }) }}
          </div>
        </div>

        <div ref="scrollContainer" class="activities-scroll-container" @scroll="handleScroll">
          <div class="activities-list">
            <div v-if="initialLoading" class="activity-item loading">
              <div class="activity-icon loading">
                <i class="fas fa-spinner fa-spin" />
              </div>
              <div class="activity-content">
                <div class="activity-description">{{ $t('dashboard.activityMonitor.loadingActivities') }}</div>
                <div class="activity-time">{{ $t('dashboard.activityMonitor.loadingHint') }}</div>
              </div>
            </div>

            <div v-for="activity in activities" v-else :key="activity.id" class="activity-item">
              <div class="activity-icon" :class="activity.class || activity.type">
                <i :class="activity.icon" />
              </div>
              <div class="activity-content">
                <div class="activity-description">{{ activity.description }}</div>
                <div class="activity-time">{{ formatTime(activity.timestamp) }}</div>
              </div>
            </div>

            <div v-if="loadingMore" class="activity-item loading-more">
              <div class="activity-icon loading">
                <i class="fas fa-spinner fa-spin" />
              </div>
              <div class="activity-content">
                <div class="activity-description">{{ $t('dashboard.activityMonitor.loadingMore') }}</div>
                <div class="activity-time">{{ $t('dashboard.activityMonitor.loadingHint') }}</div>
              </div>
            </div>

            <div v-if="!hasMore && activities.length > 0 && !initialLoading" class="activity-item no-more">
              <div class="activity-icon default">
                <i class="fas fa-check-circle" />
              </div>
              <div class="activity-content">
                <div class="activity-description">{{ $t('dashboard.activityMonitor.allLoaded') }}</div>
                <div class="activity-time">{{ $t('dashboard.activityMonitor.recordCount', { count: activities.length }) }}</div>
              </div>
            </div>

            <div v-if="!initialLoading && activities.length === 0" class="activity-item empty">
              <div class="activity-icon default">
                <i class="fas fa-info-circle" />
              </div>
              <div class="activity-content">
                <div class="activity-description">{{ $t('dashboard.activityMonitor.empty') }}</div>
                <div class="activity-time">{{ $t('dashboard.activityMonitor.emptyHint') }}</div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
  .activity-monitor {
    padding: 20px;
  }

  .monitor-content {
    display: flex;
    flex-direction: column;
    gap: 20px;
    height: 100%;
    min-height: 500px;
  }

  .monitor-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .section-title {
    margin: 0;
    font-size: 18px;
    font-weight: 600;
    color: var(--color-content-heading);
    display: flex;
    align-items: center;
    gap: 10px;
  }

  .section-title i {
    color: var(--color-brand-500);
    font-size: 16px;
  }

  .monitor-status {
    display: flex;
    align-items: center;
    gap: 8px;
  }

  .status-dot {
    width: 6px;
    height: 6px;
    border-radius: var(--radius-full);
  }

  .status-dot.active {
    background: var(--color-success-500);
    box-shadow: 0 0 6px rgba(var(--color-success-rgb), 0.6);
  }

  .status-dot.warning {
    background: var(--color-warning-500);
    box-shadow: 0 0 6px rgba(255, 167, 38, 0.6);
  }

  .status-dot.error {
    background: var(--color-error-500);
    box-shadow: 0 0 6px rgba(255, 0, 102, 0.6);
  }

  .status-text {
    font-size: 12px;
    color: rgba(var(--color-content-rgb), 0.8);
    font-weight: 500;
  }

  .section-subtitle {
    font-size: 14px;
    color: rgba(var(--color-content-rgb), 0.8);
    font-weight: 600;
    margin-bottom: 12px;
  }

  .stats-overview {
    display: grid;
    grid-template-columns: repeat(2, 1fr);
    gap: 12px;
    margin-bottom: 4px;
  }

  .stat-item {
    display: flex;
    align-items: center;
    gap: 12px;
    background: var(--color-background-800);
    border: 1px solid var(--color-border-subtle);
    border-radius: var(--radius-sm);
    padding: 16px;
    transition: all 0.2s ease;
  }

  .stat-item:hover {
    border-color: var(--color-hover-border);
    background: var(--color-hover-bg);
  }

  .stat-icon {
    width: 36px;
    height: 36px;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: var(--radius-sm);
    font-size: 16px;
    flex-shrink: 0;
  }

  .stat-icon.upload {
    background: rgba(var(--color-brand-500-rgb), 0.15);
    color: var(--color-brand-500);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.2);
  }

  .stat-icon.view {
    background: rgba(var(--color-success-rgb), 0.15);
    color: var(--color-success-500);
    border: 1px solid rgba(var(--color-success-rgb), 0.2);
  }

  .stat-content {
    flex: 1;
  }

  .stat-value {
    font-size: 18px;
    font-weight: 700;
    color: var(--color-content-heading);
    font-family: 'SF Mono', monospace;
    margin-bottom: 2px;
  }

  .stat-label {
    font-size: 12px;
    color: rgba(var(--color-content-rgb), 0.6);
    font-weight: 500;
  }

  .recent-activities {
    flex: 1;
    display: flex;
    flex-direction: column;
    min-height: 0;
  }

  .activities-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 12px;
    flex-shrink: 0;
  }

  .activities-counter {
    font-size: 12px;
    color: rgba(var(--color-content-rgb), 0.6);
    font-weight: 500;
  }

  .activities-scroll-container {
    flex: 1;
    max-height: 360px;
    overflow-y: auto;
    overflow-x: hidden;
    padding-right: 4px;
    margin-right: -4px;
  }

  .activities-scroll-container::-webkit-scrollbar {
    width: 6px;
  }

  .activities-scroll-container::-webkit-scrollbar-track {
    background: var(--color-background-900);
    border-radius: var(--radius-sm);
  }

  .activities-scroll-container::-webkit-scrollbar-thumb {
    background: var(--color-brand-500);
    border-radius: var(--radius-sm);
    transition: background 0.2s ease;
  }

  .activities-scroll-container::-webkit-scrollbar-thumb:hover {
    background: var(--color-brand-600);
  }

  .activities-list {
    display: flex;
    flex-direction: column;
    gap: 8px;
    padding-bottom: 8px;
  }

  .activity-item {
    display: flex;
    align-items: center;
    gap: 12px;
    background: var(--color-background-800);
    border: 1px solid var(--color-border-subtle);
    border-radius: var(--radius-sm);
    padding: 12px;
  }

  .activity-icon {
    width: 28px;
    height: 28px;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: var(--radius-sm);
    font-size: 12px;
    flex-shrink: 0;
  }

  .activity-icon.upload {
    background: rgba(var(--color-brand-500-rgb), 0.15);
    color: var(--color-brand-500);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.2);
  }

  .activity-icon.view {
    background: rgba(var(--color-success-rgb), 0.15);
    color: var(--color-success-500);
    border: 1px solid rgba(var(--color-success-rgb), 0.2);
  }

  .activity-icon.share {
    background: rgba(var(--color-badge-accent-text-rgb), 0.15);
    color: var(--color-pink-400);
    border: 1px solid rgba(var(--color-badge-accent-text-rgb), 0.2);
  }

  .activity-icon.api {
    background: rgba(var(--color-badge-accent-text-rgb), 0.15);
    color: var(--color-brand-500);
    border: 1px solid rgba(var(--color-badge-accent-text-rgb), 0.2);
  }

  .activity-icon.system {
    background: rgba(255, 167, 38, 0.15);
    color: var(--color-warning-500);
    border: 1px solid rgba(255, 167, 38, 0.2);
  }

  .activity-icon.login {
    background: rgba(var(--color-success-rgb), 0.15);
    color: var(--color-success-500);
    border: 1px solid rgba(var(--color-success-rgb), 0.2);
  }

  .activity-icon.register {
    background: rgba(var(--color-brand-500-rgb), 0.15);
    color: var(--color-brand-500);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.2);
  }

  .activity-icon.edit {
    background: rgba(103, 194, 255, 0.15);
    color: var(--color-info-500);
    border: 1px solid rgba(103, 194, 255, 0.2);
  }

  .activity-icon.move {
    background: rgba(156, 39, 176, 0.15);
    color: var(--color-brand-600);
    border: 1px solid rgba(156, 39, 176, 0.2);
  }

  .activity-icon.create {
    background: rgba(var(--color-success-rgb), 0.15);
    color: var(--color-success-500);
    border: 1px solid rgba(var(--color-success-rgb), 0.2);
  }

  .activity-icon.delete {
    background: rgba(255, 107, 107, 0.15);
    color: var(--color-red-400);
    border: 1px solid rgba(255, 107, 107, 0.2);
  }

  .activity-icon.folder {
    background: rgba(255, 193, 7, 0.15);
    color: var(--color-warning);
    border: 1px solid rgba(255, 193, 7, 0.2);
  }

  .activity-icon.profile {
    background: rgba(156, 39, 176, 0.15);
    color: var(--color-brand-600);
    border: 1px solid rgba(156, 39, 176, 0.2);
  }

  .activity-icon.security {
    background: rgba(255, 87, 34, 0.15);
    color: var(--color-error-600);
    border: 1px solid rgba(255, 87, 34, 0.2);
  }

  .activity-icon.permission {
    background: rgba(63, 81, 181, 0.15);
    color: var(--color-info-600);
    border: 1px solid rgba(63, 81, 181, 0.2);
  }

  .activity-icon.expired {
    background: rgba(255, 152, 0, 0.15);
    color: var(--color-warning-500);
    border: 1px solid rgba(255, 152, 0, 0.2);
  }

  .activity-icon.default {
    background: rgba(148, 163, 184, 0.15);
    color: var(--color-slate-400);
    border: 1px solid rgba(148, 163, 184, 0.2);
  }

  .activity-icon.loading {
    background: rgba(var(--color-brand-500-rgb), 0.15);
    color: var(--color-brand-500);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.2);
  }

  .activity-item.loading {
    opacity: 0.6;
  }

  .activity-item.empty {
    opacity: 0.8;
  }

  .activity-item.loading-more {
    opacity: 0.7;
    background: rgba(var(--color-brand-500-rgb), 0.05);
    border-color: rgba(var(--color-brand-500-rgb), 0.1);
  }

  .activity-item.no-more {
    opacity: 0.6;
    background: rgba(var(--color-success-rgb), 0.05);
    border-color: rgba(var(--color-success-rgb), 0.1);
  }

  .activity-content {
    flex: 1;
  }

  .activity-description {
    font-size: 12px;
    color: var(--color-content-default);
    margin-bottom: 2px;
  }

  .activity-time {
    font-size: 10px;
    color: var(--color-content-muted);
  }

  @media (max-width: 768px) {
    .stats-overview {
      grid-template-columns: 1fr;
    }

    .monitor-header {
      flex-direction: column;
      gap: 12px;
      align-items: flex-start;
    }

    .activity-monitor {
      padding: 16px;
    }
  }
</style>
