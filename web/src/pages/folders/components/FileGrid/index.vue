<script setup lang="ts">
  import { computed } from 'vue'
  import type { FileGridEmits, FileGridProps, ImageGridEmits, ImageGridProps } from './types'
  import FilesGridView from './FilesGridView.vue'
  import FilesListView from './FilesListView.vue'
  import { useTexts } from '@/composables/useTexts'

  defineOptions({
    name: 'FileGrid',
  })

  const { $t } = useTexts()

  interface DualProps extends Partial<FileGridProps>, Partial<ImageGridProps> {
    files?: FileGridProps['files']
    selectedBatchFiles?: FileGridProps['selectedBatchFiles']
    isAllFilesSelected?: FileGridProps['isAllFilesSelected']
    isFileSelected?: FileGridProps['isFileSelected']
    isBatchFileSelected?: FileGridProps['isBatchFileSelected']
    images?: ImageGridProps['images']
    selectedBatchImages?: ImageGridProps['selectedBatchImages']
    isAllImagesSelected?: ImageGridProps['isAllImagesSelected']
    isImageSelected?: ImageGridProps['isImageSelected']
    isBatchImageSelected?: ImageGridProps['isBatchImageSelected']
    pagination?: { total: number; page: number; limit: number }
    selectMode: boolean
    batchMode: boolean
    viewMode: 'grid' | 'list'
    preview?: boolean
  }

  const props = defineProps<DualProps>()
  const _emit = defineEmits<FileGridEmits & ImageGridEmits>()

  const images = computed(() => props.files || props.images || [])
  const selectedBatchImages = computed(() => props.selectedBatchFiles || props.selectedBatchImages || [])
  const isAllImagesSelected = computed(() => props.isAllFilesSelected ?? props.isAllImagesSelected ?? false)
  const isImageSelected = computed(() => props.isFileSelected || props.isImageSelected)
  const isBatchImageSelected = computed(() => props.isBatchFileSelected || props.isBatchImageSelected)
</script>

<template>
  <div v-if="images.length > 0" class="section-container">
    <div class="section-header">
      <h2 class="section-title">
        <i class="fas fa-images" /> {{ $t('folders.fileGrid.sectionTitle') }}
        <span class="item-count">{{ pagination?.total || images.length }}</span>
        <span v-if="!selectMode && !batchMode && !preview && images.length > 1 && viewMode === 'grid'" class="drag-tip">
          <i class="fas fa-arrows-alt" /> {{ $t('folders.fileGrid.dragTip') }}
        </span>
      </h2>
      <div class="section-actions">
        <div v-if="!batchMode && !selectMode && !preview" class="batch-actions">
          <button class="batch-btn" :title="$t('folders.fileGrid.actions.batchMode')" @click="$emit('start-batch-mode')">
            <i class="fas fa-tasks" />
            {{ $t('folders.fileGrid.actions.batchMode') }}
          </button>
        </div>

        <div v-else-if="batchMode && !preview" class="batch-mode-controls">
          <div class="batch-selection-info">
            <span class="batch-count bg-error-200 text-error-500">
              <i class="fas fa-image" />{{ $t('folders.fileGrid.selectedCount', { count: selectedBatchImages.length }) }}
            </span>
          </div>

          <button
            class="batch-btn danger-btn"
            :disabled="selectedBatchImages.length === 0"
            :title="$t('folders.fileGrid.actions.batchDelete')"
            @click="$emit('batch-delete')"
          >
            <i class="fas fa-trash" />
            {{ $t('folders.fileGrid.actions.batchDelete') }}
          </button>
          <button
            class="batch-btn cancel-btn"
            :title="$t('folders.fileGrid.actions.cancelBatch')"
            @click="$emit('cancel-batch-mode')"
          >
            <i class="fas fa-times" />
            {{ $t('actions.cancel') }}
          </button>
        </div>

        <div v-if="(selectMode || batchMode) && !preview" class="select-actions">
          <CyberButton
            type="ghost"
            :icon="isAllImagesSelected ? 'check-double' : 'check'"
            :title="
              isAllImagesSelected
                ? $t('folders.fileGrid.actions.deselectAllTitle')
                : $t('folders.fileGrid.actions.selectAllTitle')
            "
            :class="{ active: isAllImagesSelected }"
            @click="$emit('toggle-select-all-images')"
          >
            {{ isAllImagesSelected ? $t('folders.fileGrid.actions.deselectAll') : $t('folders.fileGrid.actions.selectAll') }}
          </CyberButton>
          <CyberButton
            type="ghost"
            icon="exchange-alt"
            :title="$t('folders.fileGrid.actions.invertTitle')"
            @click="$emit('invert-image-selection')"
          >
            {{ $t('folders.fileGrid.actions.invert') }}
          </CyberButton>
        </div>

        <div class="view-actions">
          <button
            class="view-toggle-btn"
            :class="{ active: viewMode === 'grid' }"
            :title="$t('folders.fileGrid.actions.gridView')"
            @click="$emit('change-view-mode', 'grid')"
          >
            <i class="fas fa-th" />
          </button>
          <button
            class="view-toggle-btn"
            :class="{ active: viewMode === 'list' }"
            :title="$t('folders.fileGrid.actions.listView')"
            @click="$emit('change-view-mode', 'list')"
          >
            <i class="fas fa-list" />
          </button>
        </div>
      </div>
    </div>

    <FilesGridView
      v-if="viewMode === 'grid'"
      :images="images"
      :select-mode="selectMode"
      :batch-mode="batchMode"
      :is-image-selected="isImageSelected"
      :is-batch-image-selected="isBatchImageSelected"
      :preview="preview"
      @preview-image="$emit('preview-image', $event)"
      @view-details="$emit('view-details', $event)"
      @copy-link="$emit('copy-link', $event)"
      @download-image="$emit('download-image', $event)"
      @delete-image="$emit('delete-image', $event)"
      @toggle-image-visibility="$emit('toggle-image-visibility', $event)"
      @toggle-image-select="$emit('toggle-image-select', $event)"
      @toggle-batch-image-select="$emit('toggle-batch-image-select', $event)"
      @image-drag-start="$emit('image-drag-start')"
      @image-drag-end="$emit('image-drag-end', $event)"
      @image-moved="$emit('image-moved', $event[0], $event[1])"
      @image-deleted="$emit('image-deleted', $event)"
    />

    <FilesListView
      v-else
      :images="images"
      :select-mode="selectMode"
      :batch-mode="batchMode"
      :is-all-images-selected="isAllImagesSelected"
      :is-image-selected="isImageSelected"
      :is-batch-image-selected="isBatchImageSelected"
      :preview="preview"
      @preview-image="$emit('preview-image', $event)"
      @view-details="$emit('view-details', $event)"
      @copy-link="$emit('copy-link', $event)"
      @download-image="$emit('download-image', $event)"
      @delete-image="$emit('delete-image', $event)"
      @toggle-image-select="$emit('toggle-image-select', $event)"
      @toggle-batch-image-select="$emit('toggle-batch-image-select', $event)"
      @toggle-select-all-images="$emit('toggle-select-all-images')"
    />
  </div>
</template>

<style scoped>
  .section-container {
    background: rgba(var(--color-background-800-rgb), 0.95);
    backdrop-filter: blur(10px);
    border: 1.5px solid rgba(var(--color-brand-500-rgb), 0.2);
    border-radius: var(--radius-sm);
    overflow: visible;
    padding: 1.25rem;
    box-shadow:
      0 4px 12px var(--color-overlay-medium),
      0 2px 6px var(--color-overlay-light);
  }

  .section-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 1rem;
    padding-bottom: 0.6rem;
    border-bottom: 1.5px solid rgba(var(--color-brand-500-rgb), 0.25);
    flex-wrap: wrap;
    gap: 0.75rem;
  }

  .section-title {
    margin: 0;
    font-size: 1.1rem;
    font-weight: 600;
    color: var(--color-content-heading);
    display: flex;
    align-items: center;
    gap: 0.4rem;
  }

  .section-title i {
    color: var(--color-brand-500);
    font-size: 1rem;
  }

  .item-count {
    background: rgba(var(--color-brand-500-rgb), 0.2);
    color: var(--color-brand-500);
    padding: 0.15rem 0.5rem;
    border-radius: var(--radius-sm);
    font-size: 0.75rem;
    font-weight: 600;
    margin-left: 0.4rem;
  }

  .drag-tip {
    background: rgba(var(--color-warning-rgb), 0.1);
    color: var(--color-warning-500);
    padding: 0.1rem 0.4rem;
    border-radius: var(--radius-sm);
    font-size: 0.7rem;
    margin-left: 0.5rem;
    display: flex;
    align-items: center;
    gap: 0.2rem;
    opacity: 0.8;
    transition: all 0.2s ease;
  }

  .drag-tip:hover {
    opacity: 1;
    background: rgba(var(--color-warning-rgb), 0.15);
  }

  .drag-tip i {
    font-size: 0.65rem;
  }

  .section-actions {
    display: flex;
    gap: 0.4rem;
    align-items: center;
    flex-wrap: wrap;
  }

  .batch-actions {
    display: flex;
    gap: 0.5rem;
  }

  .batch-btn {
    display: flex;
    align-items: center;
    gap: 0.3rem;
    padding: 0.4rem 0.8rem;
    border-radius: var(--radius-sm);
    font-size: 0.85rem;
    cursor: pointer;
    transition: all 0.2s ease;
    border: 1px solid transparent;
    background: rgba(var(--color-brand-500-rgb), 0.1);
    color: var(--color-content-default);
  }

  .batch-btn:hover {
    background: rgba(var(--color-brand-500-rgb), 0.2);
    transform: translateY(-1px);
    color: var(--color-brand-500);
  }

  .batch-btn:disabled {
    opacity: 0.5;
    cursor: not-allowed;
    transform: none;
  }

  .batch-btn.danger-btn {
    background: rgba(var(--color-error-rgb), 0.1);
    color: var(--color-error-500);
  }

  .batch-btn.danger-btn:hover:not(:disabled) {
    background: rgba(var(--color-error-rgb), 0.2);
    color: var(--color-error-600);
  }

  .batch-btn.cancel-btn {
    background: rgba(var(--color-content-rgb), 0.1);
    color: var(--color-content-default);
  }

  .batch-btn.cancel-btn:hover {
    background: rgba(var(--color-content-rgb), 0.2);
    color: var(--color-content-heading);
  }

  .batch-btn i {
    font-size: 0.9rem;
  }

  .batch-mode-controls {
    display: flex;
    align-items: center;
    gap: 0.5rem;
  }

  .batch-selection-info {
    display: flex;
    align-items: center;
    gap: 0.3rem;
  }

  .batch-count {
    font-size: 0.85rem;
    padding: 0.2rem 0.5rem;
    border-radius: var(--radius-sm);
    display: flex;
    align-items: center;
    gap: 0.3rem;
    background: rgba(var(--color-error-rgb), 0.15);
    color: var(--color-error-500);
  }

  .batch-count i {
    font-size: 0.8rem;
  }

  .select-actions {
    display: flex;
    gap: 0.5rem;
    margin-right: 0.5rem;
    align-items: center;
  }

  .select-mode-tip {
    display: flex;
    align-items: center;
    padding: 0.4rem 0.8rem;
    background: rgba(var(--color-brand-500-rgb), 0.1);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.3);
    border-radius: var(--radius-sm);
    color: rgba(var(--color-brand-500-rgb), 0.9);
    font-size: 0.85rem;
    font-weight: 500;
    white-space: nowrap;
    margin-right: 0.8rem;
    animation: tipPulse 2s ease-in-out infinite;
  }

  .select-mode-tip i {
    font-size: 0.9rem;
    color: rgba(var(--color-brand-500-rgb), 0.8);
  }

  @keyframes tipPulse {
    0%,
    100% {
      opacity: 0.8;
      transform: scale(1);
    }
    50% {
      opacity: 1;
      transform: scale(1.02);
    }
  }

  .view-actions {
    display: flex;
    gap: 0.4rem;
    margin-left: 0.5rem;
  }

  .view-toggle-btn {
    background: rgba(var(--color-brand-500-rgb), 0.1);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.2);
    color: var(--color-content-muted);
    width: 32px;
    height: 32px;
    border-radius: var(--radius-sm);
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    transition: all 0.2s ease;
    flex-shrink: 0;
  }

  .view-toggle-btn:hover {
    background: rgba(var(--color-brand-500-rgb), 0.15);
    color: var(--color-brand-500);
    border-color: rgba(var(--color-brand-500-rgb), 0.3);
    transform: translateY(-1px);
  }

  .view-toggle-btn.active {
    background: rgba(var(--color-brand-500-rgb), 0.2);
    color: var(--color-brand-500);
    border-color: rgba(var(--color-brand-500-rgb), 0.4);
    box-shadow: 0 0 8px rgba(var(--color-brand-500-rgb), 0.3);
  }
</style>
