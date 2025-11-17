<template>
  <div class="vector-list">
    <div class="list-header">
      <div class="list-title">
        <span class="list-count"
          >{{ $t('admin.vectors.pagination.total') }} {{ pagination.total }} {{ $t('admin.vectors.pagination.items') }}</span
        >
      </div>

      <div class="list-actions">
        <CyberButton type="secondary" @click="$emit('toggle-filter')" :class="{ 'border-error-500 text-error-500': showFilter }">
          <i class="fas fa-filter mr-1"></i>{{ $t('admin.vectors.filter.title') }}
        </CyberButton>
        <CyberButton
          v-if="selectedVectors.length > 0"
          type="warning"
          @click="showBatchActions = !showBatchActions"
          class="batch-btn ml-2"
        >
          <i class="fas fa-layer-group mr-1"></i>
          {{ $t('admin.vectors.batchOperations.title') }} ({{ selectedVectors.length }})
        </CyberButton>
      </div>
    </div>

    <transition
      name="filter-panel"
      enter-active-class="transition duration-300 ease-out"
      leave-active-class="transition duration-200 ease-in"
      enter-from-class="opacity-0 transform -translate-y-8"
      enter-to-class="opacity-100 transform translate-y-0"
      leave-from-class="opacity-100 transform translate-y-0"
      leave-to-class="opacity-0 transform -translate-y-8"
    >
      <div v-show="showFilter" class="filter-panel-container">
        <VectorFilterPanel :initial-filters="filters" :available-models="availableModels" @filter="$emit('filter', $event)" />
      </div>
    </transition>

    <transition
      name="batch-panel"
      enter-active-class="transition duration-300 ease-out"
      leave-active-class="transition duration-200 ease-in"
      enter-from-class="opacity-0 transform -translate-y-4"
      enter-to-class="opacity-100 transform translate-y-0"
      leave-from-class="opacity-100 transform translate-y-0"
      leave-to-class="opacity-0 transform -translate-y-4"
    >
      <div v-show="showBatchActions" class="batch-actions-panel">
        <div class="batch-info">
          {{ $t('admin.vectors.batchOperations.selected') }} {{ selectedVectors.length }}
          {{ $t('admin.vectors.batchOperations.items') }}
        </div>
        <div class="batch-buttons">
          <CyberButton type="primary" :loading="processing" @click="handleBatchAction('retry')">
            <i class="fas fa-redo mr-1"></i>{{ $t('admin.vectors.actions.search') }}
          </CyberButton>
          <CyberButton type="warning" :loading="processing" @click="handleBatchAction('reset')">
            <i class="fas fa-sync-alt mr-1"></i>{{ $t('admin.vectors.buttons.reset') }}
          </CyberButton>
          <CyberButton type="danger" :loading="processing" @click="handleBatchAction('delete')">
            <i class="fas fa-trash mr-1"></i>{{ $t('admin.vectors.actions.delete') }}
          </CyberButton>
          <CyberButton type="secondary" @click="clearSelection">
            <i class="fas fa-times mr-1"></i>{{ $t('admin.vectors.batchOperations.deselectAll') }}
          </CyberButton>
        </div>
      </div>
    </transition>

    <div class="vector-table-container">
      <CyberSkeleton type="table" :count="pagination.page_size" :loading="isLoading" />

      <div v-if="!isLoading && vectors.length === 0" class="empty-state">
        <div class="empty-icon">
          <i class="fas fa-vector-square"></i>
        </div>
        <div class="empty-text">{{ $t('admin.vectors.empty.title') }}</div>
      </div>

      <div v-if="!isLoading && vectors.length > 0" class="vector-table">
        <div class="table-header">
          <div class="cell-checkbox table-cell">
            <CyberCheckbox :model-value="isAllSelected" :indeterminate="isIndeterminate" @update:model-value="handleSelectAll" />
          </div>
          <div class="cell-image table-cell">{{ $t('admin.vectors.columns.filename') }}</div>
          <div class="cell-status table-cell">{{ $t('admin.vectors.columns.status') }}</div>
          <div class="cell-model table-cell">{{ $t('admin.vectors.qdrant.version') }}</div>
          <div class="cell-dimension table-cell">{{ $t('admin.vectors.columns.dimensions') }}</div>
          <div class="cell-vector table-cell">{{ $t('admin.vectors.form.description') }}</div>
          <div class="cell-retry table-cell">{{ $t('admin.vectors.actions.rebuild') }}</div>
          <div class="cell-duration table-cell">{{ $t('admin.vectors.logs.duration') }}</div>
          <div class="cell-time table-cell">{{ $t('admin.vectors.columns.updatedAt') }}</div>
          <div class="cell-actions table-cell">{{ $t('admin.vectors.columns.actions') }}</div>
        </div>

        <div class="table-body">
          <div
            v-for="vector in vectors"
            :key="vector.id"
            class="table-row"
            :class="{ selected: selectedVectors.includes(vector.file_id) }"
          >
            <div class="cell-checkbox table-cell">
              <CyberCheckbox
                :model-value="selectedVectors.includes(vector.file_id)"
                @update:model-value="handleSelectVector(vector.file_id, $event)"
              />
            </div>

            <div class="cell-image table-cell">
              <div class="image-info">
                <div v-if="vector.thumbnail_url" class="image-preview" @click="handleViewDetail(vector.file_id)">
                  <img :src="vector.thumbnail_url" :alt="vector.file_name || $t('admin.vectors.columns.filename')" />
                </div>
                <div v-else class="image-placeholder">
                  <i class="fas fa-image"></i>
                </div>
                <div class="image-details">
                  <div class="file-name">{{ vector.file_name || 'Unknown' }}</div>
                  <div class="image-id">{{ vector.file_id }}</div>
                </div>
              </div>
            </div>

            <div class="cell-status table-cell">
              <div class="status-badge" :class="`status-${vector.status}`">
                <i :class="getStatusIcon(vector.status)"></i>
                <span>{{ getStatusText(vector.status) }}</span>
              </div>
              <div v-if="vector.error_message" class="error-message">
                {{ vector.error_message }}
              </div>
            </div>

            <div class="cell-model table-cell">
              <div class="model-info">
                <div class="model-name">{{ vector.model }}</div>
              </div>
            </div>

            <div class="cell-dimension table-cell">
              <div class="dimension-info">
                {{ vector.dimension || '-' }}
              </div>
            </div>

            <div class="cell-vector table-cell">
              <div class="vector-content">
                <CyberTooltip v-if="vector.description" :content="vector.description" placement="top" :max-width="400">
                  <div class="vector-description">
                    {{ vector.description }}
                  </div>
                </CyberTooltip>
                <div v-else class="vector-empty">
                  <span>-</span>
                </div>
              </div>
            </div>

            <div class="cell-retry table-cell">
              <div class="retry-info">
                <span class="retry-count">{{ vector.retry_count }}</span>
                <span v-if="vector.last_retry_at" class="retry-time">
                  {{ formatTime(vector.last_retry_at) }}
                </span>
              </div>
            </div>

            <div class="cell-duration table-cell">
              <div class="duration-info">
                {{ formatDuration(vector.processing_duration) }}
              </div>
            </div>

            <div class="cell-time table-cell">
              <div class="time-info">
                <div class="update-time">{{ formatTime(vector.updated_at) }}</div>
                <div class="create-time">{{ formatTime(vector.created_at) }}</div>
              </div>
            </div>

            <div class="cell-actions table-cell">
              <div class="action-buttons">
                <CyberButton
                  v-if="canRetry(vector)"
                  type="primary"
                  size="mini"
                  :loading="processing"
                  @click="handleRetry(vector.file_id)"
                  :title="$t('admin.vectors.actions.search')"
                >
                  <i class="fas fa-redo"></i>
                </CyberButton>

                <CyberButton
                  type="info"
                  size="mini"
                  @click="handleViewDetail(vector.file_id)"
                  :title="$t('admin.vectors.actions.view')"
                >
                  <i class="fas fa-eye"></i>
                </CyberButton>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div v-if="pagination.total > 0" class="pagination-container">
      <CyberPagination
        :currentPage="pagination.page"
        :total="pagination.total"
        :pageSize="pagination.page_size"
        @page-change="handlePageChange"
        @page-size-change="handlePageSizeChange"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
  import { ref, computed, watch } from 'vue'
  import { useTexts } from '@/composables/useTexts'
  import type { VectorItem, PaginationInfo, VectorListParams } from '@/api/admin/vectors'
  import { formatDistanceToNow } from '@/utils/formatting/dateUtils'
  import VectorFilterPanel from './VectorFilterPanel.vue'

  interface Props {
    vectors: VectorItem[]
    isLoading?: boolean
    pagination: PaginationInfo
    processing?: boolean
    showFilter?: boolean
    availableModels?: string[]
    filters?: VectorListParams
  }

  interface Emits {
    (e: 'retry', fileId: string): void
    (e: 'batch-action', fileIds: string[], action: 'reset' | 'retry' | 'delete'): void
    (e: 'page-change', page: number): void
    (e: 'page-size-change', size: number): void
    (e: 'view-detail', fileId: string): void
    (e: 'toggle-filter'): void
    (e: 'filter', filters: VectorListParams): void
  }

  const props = withDefaults(defineProps<Props>(), {
    isLoading: false,
    processing: false,
    showFilter: false,
    availableModels: () => [],
    filters: () => ({ page: 1, page_size: 10 }),
  })

  const emit = defineEmits<Emits>()
  const { $t } = useTexts()

  const selectedVectors = ref<string[]>([])
  const showBatchActions = ref(false)

  watch(
    () => props.vectors,
    () => {
      selectedVectors.value = []
    },
    { deep: true, immediate: true }
  )

  const isAllSelected = computed(() => {
    const result = props.vectors.length > 0 && selectedVectors.value.length === props.vectors.length
    return result
  })

  const isIndeterminate = computed(() => {
    const result = selectedVectors.value.length > 0 && selectedVectors.value.length < props.vectors.length
    return result
  })

  const handleSelectAll = (selected: boolean) => {
    if (selected) {
      selectedVectors.value = props.vectors.map((v) => v.file_id)
    } else {
      selectedVectors.value = []
    }
  }

  const handleSelectVector = (fileId: string, selected: boolean) => {
    if (selected) {
      if (!selectedVectors.value.includes(fileId)) {
        selectedVectors.value.push(fileId)
      }
    } else {
      const index = selectedVectors.value.indexOf(fileId)
      if (index > -1) {
        selectedVectors.value.splice(index, 1)
      }
    }
  }

  const clearSelection = () => {
    selectedVectors.value = []
    showBatchActions.value = false
  }

  const handleRetry = (fileId: string) => {
    emit('retry', fileId)
  }

  const handleBatchAction = (action: 'reset' | 'retry' | 'delete') => {
    if (selectedVectors.value.length === 0) return
    emit('batch-action', selectedVectors.value, action)
    clearSelection()
  }

  const handlePageChange = (page: number) => {
    emit('page-change', page)
    clearSelection()
  }

  const handlePageSizeChange = (size: number) => {
    emit('page-size-change', size)
    clearSelection()
  }

  const handleViewDetail = (fileId: string) => {
    emit('view-detail', fileId)
  }

  const getStatusText = (status: string): string => {
    const statusMap: Record<string, string> = {
      pending: $t('admin.vectors.status.indexing'),
      running: $t('admin.vectors.status.processing'),
      completed: $t('admin.vectors.status.active'),
      failed: $t('admin.vectors.status.error'),
      reset: $t('admin.vectors.buttons.reset'),
    }
    return statusMap[status] || status
  }

  const getStatusIcon = (status: string): string => {
    const iconMap: Record<string, string> = {
      pending: 'fas fa-clock',
      running: 'fas fa-spinner fa-spin',
      completed: 'fas fa-check-circle',
      failed: 'fas fa-exclamation-circle',
      reset: 'fas fa-redo',
    }
    return iconMap[status] || 'fas fa-question-circle'
  }

  const formatTime = (time: string): string => {
    return formatDistanceToNow(time)
  }

  const formatDuration = (duration: number): string => {
    if (!duration || duration === 0) return '-'
    if (duration < 60) return `${duration}s`
    const minutes = Math.floor(duration / 60)
    const seconds = duration % 60
    return seconds > 0 ? `${minutes}m ${seconds}s` : `${minutes}m`
  }

  const canRetry = (vector: VectorItem): boolean => {
    return vector.status === 'failed' || vector.status === 'reset'
  }
</script>

<style scoped lang="scss">
  .vector-list {
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
  }

  .list-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 20px 24px;
    border-bottom: 1px solid rgba(255, 255, 255, 0.1);

    .list-title {
      display: flex;
      align-items: center;
      gap: 12px;

      h3 {
        font-size: 1.1rem;
        font-weight: 600;
        color: var(--color-brand-500);
        margin: 0;
      }

      .list-count {
        font-size: 0.85rem;
        color: var(--color-slate-400);
      }
    }

    .list-actions {
      display: flex;
      gap: 8px;
    }
  }

  .batch-actions-panel {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 16px 24px;
    background: rgba(255, 193, 7, 0.1);
    border-bottom: 1px solid rgba(255, 193, 7, 0.3);

    .batch-info {
      color: var(--color-warning);
      font-size: 0.9rem;
      font-weight: 500;
    }

    .batch-buttons {
      display: flex;
      gap: 8px;
    }
  }

  .vector-table-container {
    position: relative;
    overflow-x: auto;
    min-height: 300px;

    &::-webkit-scrollbar {
      height: 8px;
      background: transparent;
    }

    &::-webkit-scrollbar-track {
      background: transparent;
    }

    &::-webkit-scrollbar-thumb {
      background: rgba(var(--color-brand-500-rgb), 0.3);
      border-radius: var(--radius-sm);
      transition: background 0.3s ease;

      &:hover {
        background: rgba(var(--color-brand-500-rgb), 0.5);
      }
    }

    scrollbar-width: thin;
    scrollbar-color: rgba(var(--color-brand-500-rgb), 0.3) transparent;
  }

  .empty-state {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 60px 20px;

    .empty-icon {
      font-size: 3rem;
      color: var(--color-content-muted);
      margin-bottom: 16px;
    }

    .empty-text {
      color: var(--color-slate-400);
      font-size: 1rem;
    }
  }

  .vector-table {
    min-width: 1200px; // 设置最小宽度，确保内容不会被压缩

    .table-header {
      display: grid;
      grid-template-columns:
        50px minmax(220px, 1fr) 120px minmax(180px, 200px) 80px minmax(200px, 300px)
        100px 100px 120px 100px;
      gap: 16px;
      padding: 16px 24px;
      background: rgba(var(--color-background-800-rgb), 0.6);
      border-bottom: 1px solid rgba(255, 255, 255, 0.1);
      font-weight: 600;
      color: var(--color-slate-400);
      font-size: 0.85rem;
    }

    .table-body {
      .table-row {
        display: grid;
        grid-template-columns:
          50px minmax(220px, 1fr) 120px minmax(180px, 200px) 80px minmax(200px, 300px)
          100px 100px 120px 100px;
        gap: 16px;
        padding: 16px 24px;
        border-bottom: 1px solid rgba(255, 255, 255, 0.05);
        transition: all 0.2s ease;

        &:hover {
          background: rgba(var(--color-brand-500-rgb), 0.05);
        }

        &.selected {
          background: rgba(var(--color-brand-500-rgb), 0.1);
          border-color: rgba(var(--color-brand-500-rgb), 0.3);
        }
      }
    }

    .table-cell {
      display: flex;
      align-items: center;
      font-size: 0.85rem;
      min-width: 0; // 允许flex布局正确处理溢出
      overflow: hidden; // 防止内容溢出

      &.cell-image {
        align-items: flex-start;
        overflow: hidden; // 确保文件列不会溢出
      }

      &.cell-model {
        min-width: 150px; // 确保模型列有足够宽度
      }

      &.cell-status {
        min-width: 100px; // 确保状态列有足够宽度
      }
    }
  }

  .image-info {
    display: flex;
    gap: 12px;
    width: 100%;
    min-width: 0; // 允许容器收缩
    overflow: hidden;

    .image-preview,
    .image-placeholder {
      width: 48px;
      height: 48px;
      border-radius: var(--radius-sm);
      overflow: hidden;
      cursor: pointer;
      transition: transform 0.2s ease;
      flex-shrink: 0; // 防止文件被压缩

      &:hover {
        transform: scale(1.05);
      }

      img {
        width: 100%;
        height: 100%;
        object-fit: cover;
      }
    }

    .image-placeholder {
      background: rgba(255, 255, 255, 0.1);
      display: flex;
      align-items: center;
      justify-content: center;
      color: var(--color-content-muted);
      font-size: 1.2rem;
    }

    .image-details {
      flex: 1;
      min-width: 0; // 允许flex子项收缩
      overflow: hidden; // 确保内容不会溢出

      .image-name {
        color: var(--color-white);
        font-weight: 500;
        margin-bottom: 4px;
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
        width: 100%; // 确保宽度约束
      }

      .image-id {
        color: var(--color-slate-400);
        font-size: 0.75rem;
        font-family: monospace;
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
        width: 100%; // 确保宽度约束
      }
    }
  }

  .status-badge {
    display: flex;
    align-items: center;
    gap: 6px;
    padding: 4px 8px;
    border-radius: var(--radius-sm);
    font-size: 0.8rem;
    font-weight: 500;

    &.status-pending {
      background: rgba(255, 193, 7, 0.2);
      color: var(--color-warning);
    }

    &.status-running {
      background: rgba(13, 202, 240, 0.2);
      color: var(--color-info-500);
    }

    &.status-completed {
      background: rgba(25, 135, 84, 0.2);
      color: var(--color-success-500);
    }

    &.status-failed {
      background: rgba(220, 53, 69, 0.2);
      color: var(--color-danger);
    }

    &.status-reset {
      background: rgba(255, 107, 1, 0.2);
      color: var(--color-warning-500);
    }
  }

  .error-message {
    color: var(--color-danger);
    font-size: 0.75rem;
    margin-top: 4px;
    word-break: break-all;
  }

  .model-info {
    width: 100%;
    min-width: 0;

    .model-name {
      color: var(--color-white);
      font-weight: 500;
      white-space: nowrap;
      overflow: hidden;
      text-overflow: ellipsis;
      max-width: 100%;
      font-size: 0.85rem;
    }
  }

  .dimension-info {
    color: var(--color-slate-400);
    font-family: monospace;
  }

  .retry-info {
    display: flex;
    flex-direction: column;
    gap: 2px;

    .retry-count {
      color: var(--color-white);
      font-weight: 500;
    }

    .retry-time {
      color: var(--color-content-muted);
      font-size: 0.75rem;
    }
  }

  .duration-info {
    color: var(--color-slate-400);
    font-family: monospace;
  }

  .time-info {
    display: flex;
    flex-direction: column;
    gap: 2px;

    .update-time {
      color: var(--color-white);
      font-size: 0.8rem;
    }

    .create-time {
      color: var(--color-content-muted);
      font-size: 0.75rem;
    }
  }

  .action-buttons {
    display: flex;
    gap: 4px;
  }

  .vector-content {
    width: 100%;
    min-width: 0;

    .vector-description {
      color: var(--color-gray-200);
      font-size: 0.8rem;
      line-height: 1.4;
      cursor: pointer;
      width: 100%;
      max-width: 100%;

      display: -webkit-box;
      -webkit-box-orient: vertical;
      -webkit-line-clamp: 2;
      overflow: hidden;
      text-overflow: ellipsis;
      word-wrap: break-word;
      word-break: break-word;
      max-height: 2.8em;
      hyphens: auto;
    }

    .vector-empty {
      color: var(--color-content-muted);
      font-size: 0.8rem;
      text-align: center;
    }
  }

  .pagination-container {
    padding: 20px 24px;
    border-top: 1px solid rgba(255, 255, 255, 0.1);
    display: flex;
    justify-content: center;
  }

  @media (max-width: 1400px) {
    .vector-table {
      min-width: 1100px;

      .table-header,
      .table-body .table-row {
        grid-template-columns:
          50px minmax(180px, 1fr) 100px minmax(150px, 180px) 70px minmax(180px, 250px)
          90px 90px 100px 90px;
        gap: 12px;
      }
    }
  }

  @media (max-width: 1200px) {
    .vector-table {
      min-width: 1000px;

      .table-header,
      .table-body .table-row {
        grid-template-columns:
          40px minmax(150px, 200px) 90px minmax(120px, 150px) 60px minmax(160px, 220px)
          80px 80px 90px 80px;
        gap: 10px;
        padding: 12px 16px;
      }
    }

    .table-cell {
      font-size: 0.8rem;

      &.cell-model {
        min-width: 120px;
      }

      &.cell-image {
        max-width: 200px; // 限制文件列最大宽度
      }
    }

    .image-preview,
    .image-placeholder {
      width: 40px !important;
      height: 40px !important;
    }

    .image-info {
      gap: 8px; // 减小间距

      .image-details {
        .image-name {
          font-size: 0.8rem;
        }

        .image-id {
          font-size: 0.7rem;
        }
      }
    }
  }

  @media (max-width: 1000px) {
    .vector-table {
      .table-header,
      .table-body .table-row {
        grid-template-columns:
          35px minmax(130px, 160px) 80px minmax(100px, 130px) 50px minmax(140px, 200px)
          70px 70px 80px 70px;
      }
    }

    .table-cell {
      &.cell-image {
        max-width: 160px; // 进一步限制文件列最大宽度
      }
    }

    .image-info {
      gap: 6px;

      .image-details {
        .image-name {
          font-size: 0.75rem;
          line-height: 1.2;
        }

        .image-id {
          font-size: 0.65rem;
        }
      }
    }
  }
</style>
