<script setup lang="ts">
  import { getCurrentLocale } from '@/utils/locale'
  import { computed, ref, watch } from 'vue'
  import type { TaggingLogItem } from '@/api/admin/tagging'
  import { useTexts } from '@/composables/useTexts'
  import { formatTaggingLog } from '@/utils/tagging/formatter'

  const { $t } = useTexts()

  interface Props {
    logs: TaggingLogItem[]
    pagination: {
      page: number
      limit: number
      total: number
    }
    selectedFileId?: string
  }

  interface Emits {
    (e: 'load-more'): void
    (e: 'refresh'): void
    (e: 'clear-filter'): void
  }

  const props = defineProps<Props>()
  const emit = defineEmits<Emits>()

  const logContainer = ref<HTMLElement>()
  const isLoadingMore = ref(false)

  const hasMore = computed(() => {
    const currentCount = props.logs.length
    return currentCount < props.pagination.total
  })

  const getActionIcon = (action: string): string => {
    const actionIcons: Record<string, string> = {
      manual: 'fas fa-user-cog',
      auto: 'fas fa-robot',
      system: 'fas fa-server',
      retry: 'fas fa-redo',
      reset: 'fas fa-sync-alt',
      trigger: 'fas fa-play',
      default: 'fas fa-info-circle',
    }

    return actionIcons[action] || actionIcons.default
  }

  const getActionColor = (action: string): string => {
    const actionColors: Record<string, string> = {
      manual: 'text-content',
      auto: 'text-green-400',
      system: 'text-yellow-400',
      retry: 'text-error-500',
      reset: 'text-orange-400',
      trigger: 'text-indigo-400',
      default: 'text-content-content-muted',
    }

    return actionColors[action] || actionColors.default
  }

  const getActionDescription = (action: string, status: string): string => {
    const actionDesc = $t(`admin.tagging.logAction.${action}`) || action
    const statusDesc = $t(`admin.tagging.logStatus.${status}`) || status

    return `${actionDesc} ${statusDesc} `
  }

  const formatTime = (timeString: string): string => {
    try {
      const date = new Date(timeString)
      return date.toLocaleString(getCurrentLocale(), {
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit',
        second: '2-digit',
        hour12: false,
      })
    } catch {
      return timeString
    }
  }

  const handleScroll = async () => {
    if (!hasMore.value || isLoadingMore.value || !logContainer.value) return

    const { scrollTop, scrollHeight, clientHeight } = logContainer.value
    const threshold = 50 // 距离底部50px时触发加载

    if (scrollTop + clientHeight >= scrollHeight - threshold) {
      isLoadingMore.value = true
      emit('load-more')
    }
  }

  watch(
    () => props.logs,
    () => {
      isLoadingMore.value = false
    }
  )

  const handleRefresh = () => {
    emit('refresh')
  }

  const handleClearFilter = () => {
    emit('clear-filter')
  }
</script>

<template>
  <div class="tagging-log-list">
    <div class="log-header">
      <div class="log-info">
        <div v-if="pagination.total > 0" class="total-info">
          {{ $t('admin.tagging.logs.loaded', { current: logs.length, total: pagination.total }) }}
        </div>
        <span v-if="selectedFileId" class="log-filter-info">
          <i class="fas fa-filter mr-1"></i>
          {{ $t('admin.tagging.logs.fileFilter', { id: selectedFileId }) }}
        </span>
      </div>
      <div class="log-actions">
        <CyberButton v-if="selectedFileId" type="warning" @click="handleClearFilter">
          <i class="fas fa-times mr-1"></i>{{ $t('admin.tagging.buttons.showAll') }}
        </CyberButton>
        <CyberButton type="secondary" @click="handleRefresh">
          <i class="fas fa-refresh mr-1"></i>{{ $t('admin.tagging.buttons.refresh') }}
        </CyberButton>
      </div>
    </div>

    <div ref="logContainer" class="log-content" @scroll="handleScroll">
      <div v-if="logs.length === 0" class="empty-logs">
        <div class="empty-icon">
          <i class="fas fa-file-alt"></i>
        </div>
        <div class="empty-text">{{ $t('admin.tagging.logs.noRecords') }}</div>
      </div>

      <div v-else class="log-timeline">
        <div
          v-for="log in logs"
          :key="log.id"
          class="log-item mb-2 border-b border-default pb-2 last:mb-0 last:border-b-0 last:pb-0"
        >
          <div class="flex items-center justify-between">
            <div class="flex items-center">
              <div :class="getActionColor(log.action)" class="w-6 text-center">
                <i :class="getActionIcon(log.action)" />
              </div>
              <div class="ml-2">
                <div class="text-xs text-content-heading">
                  {{ getActionDescription(log.action, log.status) }}
                  <span class="text-content">{{ log.file_id }}</span>
                </div>
              </div>
            </div>
            <div class="text-content-content-disabled text-xs">
              {{ formatTime(log.created_at) }}
            </div>
          </div>
          <div class="ml-8 mt-1 text-xs text-content-muted">
            {{ formatTaggingLog(log, $t) }}
          </div>
        </div>

        <div v-if="isLoadingMore" class="loading-indicator">
          <div class="loading-spinner">
            <i class="fas fa-spinner fa-spin"></i>
          </div>
          <span class="loading-text">{{ $t('admin.tagging.logs.loading') }}</span>
        </div>

        <div v-else-if="!hasMore && logs.length > 0" class="no-more-hint">
          <i class="fas fa-check-circle mr-2"></i>
          {{ $t('admin.tagging.logs.allLoaded') }}
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped lang="scss">
  .tagging-log-list {
    background: linear-gradient(
      135deg,
      rgba(var(--color-background-800-rgb), 0.8) 0%,
      rgba(var(--color-background-700-rgb), 0.6) 100%
    );
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.2);
    border-radius: var(--radius-sm);
    backdrop-filter: blur(10px);
    box-shadow: var(--shadow-cyber-md);
    overflow: hidden;
    margin-bottom: 20px;
  }

  .log-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 20px 24px;
    border-bottom: 1px solid rgba(255, 255, 255, 0.1);

    .log-info {
      display: flex;
      align-items: center;
      gap: 12px;

      .total-info {
        color: var(--color-slate-400);
        font-size: 0.875rem;
        white-space: nowrap;
      }

      .log-filter-info {
        padding: 4px 8px;
        background: rgba(255, 193, 7, 0.2);
        border: 1px solid rgba(255, 193, 7, 0.4);
        border-radius: var(--radius-sm);
        font-size: 0.8rem;
        color: var(--color-warning);
        font-weight: 500;
      }
    }

    .log-actions {
      display: flex;
      align-items: center;
      gap: 12px;
    }
  }

  .log-content {
    padding: 20px 24px;
    max-height: 400px;
    overflow-y: auto;

    &::-webkit-scrollbar {
      width: 6px;
    }

    &::-webkit-scrollbar-track {
      background: rgba(0, 0, 0, 0.1);
      border-radius: var(--radius-sm);
    }

    &::-webkit-scrollbar-thumb {
      background: rgba(var(--color-brand-500-rgb), 0.3);
      border-radius: var(--radius-sm);

      &:hover {
        background: rgba(var(--color-brand-500-rgb), 0.5);
      }
    }
  }

  .empty-logs {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 40px 20px;
    color: var(--color-slate-400);

    .empty-icon {
      font-size: 3rem;
      margin-bottom: 12px;
      opacity: 0.5;
    }

    .empty-text {
      font-size: 0.9rem;
    }
  }

  .log-timeline {
    display: flex;
    flex-direction: column;
    gap: 8px;
  }

  .log-item {
    transition: background-color 0.2s ease;
    padding: 8px 0;
  }

  .log-item:hover {
    background-color: rgba(var(--color-brand-500-rgb), 0.05);
  }

  .loading-indicator {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 8px;
    padding: 16px;
    color: var(--color-brand-500);
    font-size: 0.9rem;

    .loading-spinner i {
      animation: spin 1s linear infinite;
    }
  }

  .no-more-hint {
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 16px;
    color: var(--color-slate-400);
    font-size: 0.85rem;
  }

  @keyframes spin {
    from {
      transform: rotate(0deg);
    }
    to {
      transform: rotate(360deg);
    }
  }

  .tagging-log-list::-webkit-scrollbar-thumb:hover {
    background: rgba(var(--color-brand-500-rgb), 0.5);
  }
</style>
