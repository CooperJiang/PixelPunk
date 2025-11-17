<script setup lang="ts">
  import { getCurrentLocale } from '@/utils/locale'
  import { ref, onMounted } from 'vue'
  import { getPublicAnnouncementList, getPublicAnnouncementDetail } from '@/api/announcement'
  import type { AnnouncementDetail } from '@/api/types/announcement'
  import { useToast } from '@/components/Toast/useToast'
  import { useTexts } from '@/composables/useTexts'
  import AnnouncementDialog from './AnnouncementDialog.vue'

  defineOptions({
    name: 'AnnouncementButton',
  })

  const toast = useToast()
  const { $t } = useTexts()

  const showDrawer = ref(false)
  const showDetailDialog = ref(false)
  const loading = ref(false)
  const detailLoading = ref(false)
  const announcements = ref<AnnouncementDetail[]>([])
  const selectedAnnouncement = ref<AnnouncementDetail | null>(null)

  /* 公告系统配置 - 使用默认值 */
  const settings = ref({
    announcement_enabled: true,
    announcement_drawer_position: 'right' as 'left' | 'right',
    announcement_drawer_width: '450px',
    announcement_display_limit: 10,
    announcement_auto_show_delay: 3, // 自动弹窗延迟（秒）
  })

  /* localStorage 常量 */
  const DISMISSED_ANNOUNCEMENTS_KEY = 'dismissed_announcements'
  const TEMPORARY_DISMISSED_KEY = 'temporary_dismissed_announcements'
  const DISMISS_DURATION = 5 * 60 * 1000 // 5分钟

  /* 从 localStorage 加载已永久忽略的公告 ID */
  const loadDismissedAnnouncements = (): number[] => {
    try {
      const stored = localStorage.getItem(DISMISSED_ANNOUNCEMENTS_KEY)
      return stored ? JSON.parse(stored) : []
    } catch (error) {
      console.error('Failed to load dismissed announcements:', error)
      return []
    }
  }

  const saveDismissedAnnouncement = (announcementId: number) => {
    try {
      const dismissed = loadDismissedAnnouncements()
      if (!dismissed.includes(announcementId)) {
        dismissed.push(announcementId)
        localStorage.setItem(DISMISSED_ANNOUNCEMENTS_KEY, JSON.stringify(dismissed))
      }
    } catch (error) {
      console.error('Failed to save dismissed announcement:', error)
    }
  }

  const loadTemporaryDismissed = (): Record<number, number> => {
    try {
      const stored = localStorage.getItem(TEMPORARY_DISMISSED_KEY)
      return stored ? JSON.parse(stored) : {}
    } catch (error) {
      console.error('Failed to load temporary dismissed announcements:', error)
      return {}
    }
  }

  const saveTemporaryDismissed = (announcementId: number) => {
    try {
      const dismissed = loadTemporaryDismissed()
      dismissed[announcementId] = Date.now()
      localStorage.setItem(TEMPORARY_DISMISSED_KEY, JSON.stringify(dismissed))
    } catch (error) {
      console.error('Failed to save temporary dismissed announcement:', error)
    }
  }

  const isTemporarilyDismissed = (announcementId: number): boolean => {
    const dismissed = loadTemporaryDismissed()
    const dismissTime = dismissed[announcementId]
    if (!dismissTime) return false

    const elapsed = Date.now() - dismissTime
    if (elapsed > DISMISS_DURATION) {
      delete dismissed[announcementId]
      localStorage.setItem(TEMPORARY_DISMISSED_KEY, JSON.stringify(dismissed))
      return false
    }
    return true
  }

  const shouldAutoShow = (announcement: AnnouncementDetail): boolean => {
    if (!announcement.is_pinned) return false

    const dismissed = loadDismissedAnnouncements()
    if (dismissed.includes(announcement.id)) return false

    if (isTemporarilyDismissed(announcement.id)) return false

    return true
  }

  const fetchAnnouncements = async () => {
    if (!settings.value.announcement_enabled) {
      return
    }

    loading.value = true
    try {
      const result = await getPublicAnnouncementList()
      if (result.code === 200 && result.data) {
        announcements.value = result.data.announcements || []

        if (result.data.config) {
          settings.value = {
            announcement_enabled: result.data.config.announcement_enabled ?? true,
            announcement_drawer_position: (result.data.config.announcement_drawer_position ?? 'right') as 'left' | 'right',
            announcement_drawer_width: result.data.config.announcement_drawer_width ?? '450px',
            announcement_display_limit: result.data.config.announcement_display_limit ?? 10,
            announcement_auto_show_delay: result.data.config.announcement_auto_show_delay ?? 3,
          }
        }
      }
    } catch (error: any) {
      console.error('Failed to fetch announcements:', error)
      toast.error(error?.message || $t('components.announcement.messages.loadListFailed'))
    } finally {
      loading.value = false
    }
  }

  const toggleDrawer = () => {
    if (!settings.value.announcement_enabled) {
      toast.info($t('components.announcement.messages.systemDisabled'))
      return
    }

    if (!showDrawer.value) {
      fetchAnnouncements()
    }
    showDrawer.value = !showDrawer.value
  }

  const viewAnnouncement = async (announcement: AnnouncementDetail) => {
    showDrawer.value = false

    detailLoading.value = true
    showDetailDialog.value = true
    selectedAnnouncement.value = announcement

    try {
      const result = await getPublicAnnouncementDetail(announcement.id)
      if (result.code === 200 && result.data) {
        selectedAnnouncement.value = result.data
      }
    } catch (error: any) {
      console.error('Failed to fetch announcement detail:', error)
      toast.error(error?.message || $t('components.announcement.messages.loadDetailFailed'))
      showDetailDialog.value = false
    } finally {
      detailLoading.value = false
    }
  }

  const closeDetailDialog = () => {
    showDetailDialog.value = false
    selectedAnnouncement.value = null
  }

  const handleNeverShowAgain = (announcementId: number) => {
    saveDismissedAnnouncement(announcementId)
    toast.success($t('components.announcement.messages.neverShowAgainSuccess'))
  }

  const handleDismissTemporarily = (announcementId: number) => {
    saveTemporaryDismissed(announcementId)
  }

  const handleViewMore = () => {
    showDrawer.value = true
  }

  const formatDate = (dateStr: string) => {
    const date = new Date(dateStr)
    return date.toLocaleDateString(getCurrentLocale(), {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit',
    })
  }

  onMounted(async () => {
    if (!settings.value.announcement_enabled) {
      return
    }

    await fetchAnnouncements()

    const pinnedToShow = announcements.value.find(shouldAutoShow)
    if (pinnedToShow) {
      const delay = (settings.value.announcement_auto_show_delay || 3) * 1000
      setTimeout(async () => {
        await viewAnnouncement(pinnedToShow)
      }, delay)
    }
  })
</script>

<template>
  <div class="announcement-container">
    <button
      class="announcement-button"
      :class="{ disabled: !settings.announcement_enabled }"
      :aria-label="$t('components.announcement.button.ariaLabel')"
      :disabled="!settings.announcement_enabled"
      @click="toggleDrawer"
    >
      <i class="fas fa-bullhorn announcement-icon" />
    </button>

    <CyberDrawer
      v-model="showDrawer"
      :title="$t('components.announcement.drawer.title')"
      :width="settings.announcement_drawer_width"
      :position="settings.announcement_drawer_position"
    >
      <template #title-icon>
        <i class="fas fa-bullhorn text-brand-500" />
      </template>

      <div class="announcement-list">
        <!-- 加载状态 -->
        <div v-if="loading" class="loading-state">
          <div class="spinner" />
          <p class="text-sm text-content-muted">{{ $t('components.announcement.status.loading') }}</p>
        </div>

        <!-- 空状态 -->
        <div v-else-if="announcements.length === 0" class="empty-state">
          <i class="fas fa-inbox empty-icon" />
          <h3 class="empty-title">{{ $t('components.announcement.empty.title') }}</h3>
          <p class="empty-description">{{ $t('components.announcement.empty.description') }}</p>
        </div>

        <!-- 公告列表 -->
        <div v-else class="announcements">
          <div
            v-for="announcement in announcements"
            :key="announcement.id"
            class="announcement-card"
            :class="{ pinned: announcement.is_pinned }"
            @click="viewAnnouncement(announcement)"
          >
            <!-- 置顶标签 -->
            <div v-if="announcement.is_pinned" class="pinned-badge">
              <i class="fas fa-thumbtack" />
              {{ $t('components.announcement.list.pinned') }}
            </div>

            <!-- 标题 -->
            <h3 class="announcement-title">
              {{ announcement.title }}
            </h3>

            <!-- 摘要 -->
            <p v-if="announcement.summary" class="announcement-summary">
              {{ announcement.summary }}
            </p>

            <!-- 底部信息 -->
            <div class="announcement-footer">
              <span class="announcement-date">
                <i class="fas fa-clock" />
                {{ formatDate(announcement.created_at) }}
              </span>
              <span class="view-detail">
                {{ $t('components.announcement.list.viewDetail') }}
                <i class="fas fa-chevron-right" />
              </span>
            </div>
          </div>
        </div>
      </div>
    </CyberDrawer>

    <!-- 公告详情弹窗 -->
    <AnnouncementDialog
      v-model="showDetailDialog"
      :announcement="selectedAnnouncement"
      :loading="detailLoading"
      @close="closeDetailDialog"
      @never-show-again="handleNeverShowAgain"
      @dismiss-temporarily="handleDismissTemporarily"
      @view-more="handleViewMore"
    />
  </div>
</template>

<style scoped>
  .announcement-container {
    flex-shrink: 0;
    position: relative;
    display: inline-block;
  }

  .announcement-button {
    position: relative;
    display: inline-flex;
    align-items: center;
    justify-content: center;
    width: var(--size-button-sm);
    height: var(--size-button-sm);
    background: rgba(var(--color-background-700-rgb), 0.1);
    border: 1px solid var(--color-border-default);
    border-radius: var(--radius-sm);
    color: var(--color-content-muted);
    cursor: pointer;
    transition: all var(--transition-normal) var(--ease-in-out);
    backdrop-filter: var(--backdrop-blur-md) saturate(120%);
    box-shadow: var(--shadow-sm);
    outline: none;
  }

  .announcement-button::after {
    content: '';
    position: absolute;
    inset: -1px;
    border-radius: inherit;
    background: linear-gradient(45deg, transparent, rgba(var(--color-brand-500-rgb), 0.1), transparent);
    opacity: 0;
    transition: opacity var(--transition-normal) ease;
    z-index: -1;
  }

  .announcement-button:hover:not(.disabled) {
    background: rgba(var(--color-brand-500-rgb), 0.1);
    border-color: var(--color-hover-border);
    color: var(--color-brand-500);
    transform: translateY(-2px);
    box-shadow:
      var(--shadow-lg),
      0 0 20px rgba(var(--color-brand-500-rgb), 0.15);
  }

  .announcement-button:hover:not(.disabled)::after {
    opacity: 1;
  }

  .announcement-button:active:not(.disabled) {
    transform: translateY(1px) scale(0.98);
  }

  .announcement-button.disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  .announcement-icon {
    font-size: var(--size-icon-sm);
    line-height: 1;
  }

  .announcement-list {
    padding: var(--space-md);
  }

  .loading-state,
  .empty-state {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: var(--space-2xl);
    text-align: center;
  }

  .spinner {
    width: 32px;
    height: 32px;
    border: 3px solid rgba(var(--color-brand-500-rgb), 0.2);
    border-top: 3px solid var(--color-brand-500);
    border-radius: var(--radius-full);
    animation: spin 1s linear infinite;
    margin-bottom: var(--space-md);
  }

  @keyframes spin {
    0% {
      transform: rotate(0deg);
    }
    100% {
      transform: rotate(360deg);
    }
  }

  .empty-icon {
    font-size: 48px;
    color: var(--color-content-muted);
    margin-bottom: var(--space-md);
  }

  .empty-title {
    font-size: var(--text-lg);
    font-weight: var(--font-semibold);
    color: var(--color-content-heading);
    margin-bottom: var(--space-xs);
  }

  .empty-description {
    font-size: var(--text-sm);
    color: var(--color-content-muted);
  }

  .announcements {
    display: flex;
    flex-direction: column;
    gap: var(--space-sm);
  }

  .announcement-card {
    position: relative;
    padding: var(--space-md);
    background: rgba(var(--color-background-700-rgb), 0.6);
    border: 1px solid var(--color-border-default);
    border-radius: var(--radius-md);
    cursor: pointer;
    transition: all var(--transition-normal) ease;
  }

  .announcement-card:hover {
    background: rgba(var(--color-background-600-rgb), 0.8);
    border-color: rgba(var(--color-brand-500-rgb), 0.4);
    transform: translateY(-2px);
    box-shadow: 0 4px 12px rgba(var(--color-brand-500-rgb), 0.15);
  }

  .announcement-card.pinned {
    background: linear-gradient(135deg, rgba(var(--color-warning-rgb), 0.08), rgba(var(--color-background-700-rgb), 0.6));
    border-color: rgba(var(--color-warning-rgb), 0.3);
  }

  .pinned-badge {
    display: inline-flex;
    align-items: center;
    gap: var(--space-xs);
    padding: 2px var(--space-sm);
    background: linear-gradient(135deg, rgba(var(--color-warning-rgb), 0.2), rgba(var(--color-warning-rgb), 0.1));
    border: 1px solid rgba(var(--color-warning-rgb), 0.4);
    border-radius: var(--radius-sm);
    color: var(--color-warning-500);
    font-size: 11px;
    font-weight: var(--font-semibold);
    margin-bottom: var(--space-sm);
  }

  .pinned-badge.large {
    padding: var(--space-xs) var(--space-md);
    font-size: var(--text-sm);
  }

  .announcement-title {
    font-size: var(--text-sm);
    font-weight: var(--font-semibold);
    color: var(--color-content-heading);
    margin-bottom: var(--space-xs);
    line-height: var(--leading-normal);
  }

  .announcement-summary {
    font-size: var(--text-xs);
    color: var(--color-content-muted);
    line-height: var(--leading-relaxed);
    margin-bottom: var(--space-sm);
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
  }

  .announcement-footer {
    display: flex;
    align-items: center;
    justify-content: space-between;
    font-size: 11px;
    color: var(--color-content-muted);
  }

  .announcement-date {
    display: flex;
    align-items: center;
    gap: var(--space-xs);
  }

  .view-detail {
    display: flex;
    align-items: center;
    gap: var(--space-xs);
    color: var(--color-brand-500);
    font-weight: var(--font-medium);
    transition: all var(--transition-fast) ease;
  }

  .announcement-card:hover .view-detail {
    color: var(--color-brand-400);
  }
</style>
