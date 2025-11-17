<template>
  <div class="vector-log-list">
    <div class="log-header">
      <div class="log-info">
        <div v-if="pagination.total > 0" class="total-info">
          {{ $t('admin.vectors.logs.filters.keywords') }} {{ logs.length }} / {{ pagination.total }}
          {{ $t('admin.vectors.logs.title') }}
        </div>
        <span v-if="selectedFileId" class="log-filter-info">
          <i class="fas fa-filter mr-1"></i>
          {{ $t('admin.vectors.columns.filename') }}: {{ selectedFileId }}
        </span>
      </div>
      <div class="log-actions">
        <CyberButton v-if="selectedFileId" type="warning" @click="handleClearFilter">
          <i class="fas fa-times mr-1"></i>{{ $t('admin.vectors.filter.options.all') }}
        </CyberButton>
        <CyberButton type="secondary" @click="handleRefresh">
          <i class="fas fa-refresh mr-1"></i>{{ $t('admin.vectors.buttons.refresh') }}
        </CyberButton>
      </div>
    </div>

    <div ref="logContainer" class="log-content" @scroll="handleScroll">
      <div v-if="logs.length === 0" class="empty-logs">
        <div class="empty-icon">
          <i class="fas fa-file-alt"></i>
        </div>
        <div class="empty-text">{{ $t('admin.vectors.logs.title') }}</div>
      </div>

      <div v-else class="log-timeline">
        <div v-for="log in logs" :key="log.id" class="log-item" :class="`log-${log.action}`">
          <div class="log-icon">
            <i :class="getActionIcon(log.action)"></i>
          </div>

          <div class="log-content-wrapper">
            <div class="log-main">
              <div class="log-header-info">
                <div class="log-primary">
                  <span class="action-badge" :class="`action-${log.action}`">
                    {{ getActionText(log.action) }}
                  </span>
                  <span class="log-message">{{ formatVectorLog(log, $t) }}</span>
                </div>
                <span class="log-time">{{ formatTime(log.created_at) }}</span>
              </div>

              <div class="log-meta">
                <span v-if="log.model" class="meta-item"> <i class="fas fa-brain"></i>{{ log.model }} </span>
                <span v-if="log.duration > 0" class="meta-item">
                  <i class="fas fa-clock"></i>{{ formatDuration(log.duration) }}
                </span>
                <span v-if="log.task_id" class="meta-item"> <i class="fas fa-tag"></i>{{ log.task_id.slice(-8) }} </span>
                <span class="meta-item"> <i class="fas fa-image"></i>{{ log.file_id }} </span>
              </div>

              <div v-if="log.error_code" class="log-error">
                <i class="fas fa-exclamation-triangle"></i>
                {{ log.error_code }}
              </div>
            </div>
          </div>
        </div>

        <div v-if="isLoadingMore" class="loading-indicator">
          <div class="loading-spinner">
            <i class="fas fa-spinner fa-spin"></i>
          </div>
          <span class="loading-text">{{ $t('admin.vectors.loading.title') }}...</span>
        </div>

        <div v-else-if="!hasMore && logs.length > 0" class="no-more-hint">
          <i class="fas fa-check-circle mr-2"></i>
          {{ $t('admin.vectors.logs.filters.all') }} {{ $t('admin.vectors.logs.title') }}
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
  import { computed, ref, watch } from 'vue'
  import { useTexts } from '@/composables/useTexts'
  import type { VectorLogItem, PaginationInfo } from '@/api/admin/vectors'
  import { formatDistanceToNow } from '@/utils/formatting/dateUtils'
  import { formatVectorLog } from '@/utils/vector/formatter'

  interface Props {
    logs: VectorLogItem[]
    pagination: PaginationInfo
    selectedFileId?: string
  }

  interface Emits {
    (e: 'load-more'): void
    (e: 'refresh'): void
    (e: 'clear-filter'): void
  }

  const props = defineProps<Props>()
  const emit = defineEmits<Emits>()
  const { $t } = useTexts()

  const logContainer = ref<HTMLElement>()

  const isLoadingMore = ref(false)

  const hasMore = computed(() => {
    const currentCount = props.logs.length
    return currentCount < props.pagination.total
  })

  const getActionText = (action: string): string => {
    const actionMap: Record<string, string> = {
      start: $t('admin.vectors.status.indexing'),
      success: $t('admin.vectors.status.active'),
      error: $t('admin.vectors.status.error'),
      retry: $t('admin.vectors.actions.search'),
      reset: $t('admin.vectors.buttons.reset'),
      delete: $t('admin.vectors.actions.delete'),
      skip: $t('admin.vectors.status.processing'),
    }
    return actionMap[action] || action
  }

  const getActionIcon = (action: string): string => {
    const iconMap: Record<string, string> = {
      start: 'fas fa-play-circle',
      success: 'fas fa-check-circle',
      error: 'fas fa-exclamation-circle',
      retry: 'fas fa-redo',
      reset: 'fas fa-sync-alt',
      delete: 'fas fa-trash-alt',
      skip: 'fas fa-forward',
    }
    return iconMap[action] || 'fas fa-info-circle'
  }

  const formatTime = (time: string): string => {
    return formatDistanceToNow(time)
  }

  const formatDuration = (duration: number): string => {
    if (!duration || duration === 0) return '-'
    if (duration < 60) return $t('admin.vectors.logs.durationSeconds', { seconds: duration })
    const minutes = Math.floor(duration / 60)
    const seconds = duration % 60
    return seconds > 0
      ? $t('admin.vectors.logs.durationMinutesSeconds', { minutes, seconds })
      : $t('admin.vectors.logs.durationMinutes', { minutes })
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

<style scoped lang="scss">
  .vector-log-list {
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
    max-height: 600px;
    overflow-y: auto;

    &::-webkit-scrollbar {
      width: 6px;
    }

    &::-webkit-scrollbar-track {
      background: rgba(255, 255, 255, 0.05);
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

    .empty-icon {
      font-size: 2.5rem;
      color: var(--color-content-muted);
      margin-bottom: 12px;
    }

    .empty-text {
      color: var(--color-slate-400);
      font-size: 0.9rem;
    }
  }

  .log-timeline {
    position: relative;

    &::before {
      content: '';
      position: absolute;
      left: 20px;
      top: 0;
      bottom: 0;
      width: 2px;
      background: linear-gradient(
        to bottom,
        transparent 0%,
        rgba(var(--color-brand-500-rgb), 0.3) 20%,
        rgba(var(--color-brand-500-rgb), 0.3) 80%,
        transparent 100%
      );
    }
  }

  .log-item {
    position: relative;
    display: flex;
    margin-bottom: 16px;

    &:last-child {
      margin-bottom: 0;
    }

    .log-icon {
      position: relative;
      z-index: 2;
      width: 32px;
      height: 32px;
      border-radius: var(--radius-full);
      display: flex;
      align-items: center;
      justify-content: center;
      font-size: 0.8rem;
      margin-right: 12px;
      flex-shrink: 0;
    }

    .log-content-wrapper {
      flex: 1;
      background: rgba(var(--color-background-800-rgb), 0.4);
      border: 1px solid rgba(255, 255, 255, 0.1);
      border-radius: var(--radius-sm);
      padding: 12px 16px;
      margin-top: 2px;
      transition: all 0.2s ease;

      &:hover {
        background: rgba(var(--color-background-800-rgb), 0.6);
        border-color: rgba(var(--color-brand-500-rgb), 0.3);
      }
    }

    &.log-start .log-icon {
      background: rgba(13, 202, 240, 0.2);
      color: var(--color-info-500);
      border: 2px solid rgba(13, 202, 240, 0.4);
    }

    &.log-success .log-icon {
      background: rgba(25, 135, 84, 0.2);
      color: var(--color-success-500);
      border: 2px solid rgba(25, 135, 84, 0.4);
    }

    &.log-error .log-icon {
      background: rgba(220, 53, 69, 0.2);
      color: var(--color-danger);
      border: 2px solid rgba(220, 53, 69, 0.4);
    }

    &.log-retry .log-icon {
      background: rgba(255, 193, 7, 0.2);
      color: var(--color-warning);
      border: 2px solid rgba(255, 193, 7, 0.4);
    }

    &.log-reset .log-icon {
      background: rgba(255, 107, 1, 0.2);
      color: var(--color-warning-500);
      border: 2px solid rgba(255, 107, 1, 0.4);
    }

    &.log-delete .log-icon {
      background: rgba(220, 53, 69, 0.2);
      color: var(--color-danger);
      border: 2px solid rgba(220, 53, 69, 0.4);
    }

    &.log-skip .log-icon {
      background: rgba(108, 117, 125, 0.2);
      color: var(--color-content-muted);
      border: 2px solid rgba(108, 117, 125, 0.4);
    }
  }

  .log-main {
    .log-header-info {
      display: flex;
      justify-content: space-between;
      align-items: flex-start;
      margin-bottom: 8px;

      .log-primary {
        display: flex;
        align-items: center;
        gap: 8px;
        flex: 1;

        .action-badge {
          padding: 3px 6px;
          border-radius: var(--radius-sm);
          font-size: 0.7rem;
          font-weight: 600;
          white-space: nowrap;

          &.action-start {
            background: rgba(13, 202, 240, 0.2);
            color: var(--color-info-500);
          }

          &.action-success {
            background: rgba(25, 135, 84, 0.2);
            color: var(--color-success-500);
          }

          &.action-error {
            background: rgba(220, 53, 69, 0.2);
            color: var(--color-danger);
          }

          &.action-retry {
            background: rgba(255, 193, 7, 0.2);
            color: var(--color-warning);
          }

          &.action-reset {
            background: rgba(255, 107, 1, 0.2);
            color: var(--color-warning-500);
          }

          &.action-delete {
            background: rgba(220, 53, 69, 0.2);
            color: var(--color-danger);
          }

          &.action-skip {
            background: rgba(108, 117, 125, 0.2);
            color: var(--color-content-muted);
          }
        }

        .log-message {
          color: var(--color-gray-200);
          font-size: 0.8rem;
          line-height: 1.3;
        }
      }

      .log-time {
        color: var(--color-content-muted);
        font-size: 0.7rem;
        white-space: nowrap;
        margin-left: 8px;
      }
    }

    .log-meta {
      display: flex;
      flex-wrap: wrap;
      gap: 6px;
      margin-bottom: 6px;

      .meta-item {
        display: flex;
        align-items: center;
        gap: 3px;
        background: rgba(255, 255, 255, 0.05);
        padding: 2px 5px;
        border-radius: var(--radius-sm);
        font-size: 0.65rem;
        color: var(--color-slate-400);

        i {
          color: var(--color-brand-500);
          font-size: 0.6rem;
        }
      }
    }

    .log-error {
      display: flex;
      align-items: center;
      gap: 4px;
      color: var(--color-red-400);
      background: rgba(255, 107, 107, 0.1);
      padding: 3px 6px;
      border-radius: var(--radius-sm);
      font-size: 0.7rem;
      border: 1px solid rgba(255, 107, 107, 0.2);

      i {
        font-size: 0.65rem;
      }
    }
  }

  .loading-indicator {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 20px;
    gap: 8px;

    .loading-spinner {
      color: var(--color-brand-500);
      font-size: 1.2rem;
    }

    .loading-text {
      color: var(--color-slate-400);
      font-size: 0.875rem;
    }
  }

  .no-more-hint {
    display: flex;
    align-items: center;
    justify-content: center;
    color: var(--color-slate-400);
    font-size: 0.875rem;
    opacity: 0.8;
    padding: 20px;
    border-top: 1px solid rgba(255, 255, 255, 0.05);
    margin-top: 16px;
  }

  @media (max-width: 768px) {
    .log-header-info {
      flex-direction: column;
      gap: 8px;

      .log-meta {
        justify-content: flex-start;
      }
    }
  }
</style>
