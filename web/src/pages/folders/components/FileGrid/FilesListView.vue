<script setup lang="ts">
  import { getCurrentLocale } from '@/utils/locale'
  import { ref } from 'vue'
  import { useTexts } from '@/composables/useTexts'
  import type { ImageInfo } from '@/api/types/index'

  defineOptions({
    name: 'FilesListView',
  })

  const { $t } = useTexts()

  interface Props {
    images: ImageInfo[]
    selectMode: boolean
    batchMode: boolean
    isAllImagesSelected: boolean
    isImageSelected?: (image: ImageInfo) => boolean
    isBatchImageSelected?: (image: ImageInfo) => boolean
    preview?: boolean
  }

  interface Emits {
    'preview-image': [image: ImageInfo]
    'view-details': [image: ImageInfo]
    'copy-link': [image: ImageInfo]
    'download-image': [image: ImageInfo]
    'delete-image': [image: ImageInfo]
    'toggle-image-select': [image: ImageInfo]
    'toggle-batch-image-select': [image: ImageInfo]
    'toggle-select-all-images': []
  }

  const props = defineProps<Props>()
  const emit = defineEmits<Emits>()
  const imageLoadingStates = ref<Record<string, boolean>>({})

  const handleItemClick = (image: ImageInfo) => {
    if ((props.selectMode || props.batchMode) && !props.preview) {
      if (props.batchMode) {
        emit('toggle-batch-image-select', image)
      } else if (props.selectMode) {
        emit('toggle-image-select', image)
      }
    } else {
      emit('preview-image', image)
    }
  }

  const getImageUrl = (url: string, _accessLevel?: string) => url

  const getAccessIcon = (accessLevel: string) => {
    switch (accessLevel) {
      case 'private':
        return 'fas fa-lock'
      case 'protected':
        return 'fas fa-shield-alt'
      default:
        return 'fas fa-globe'
    }
  }

  const formatFileSize = (size: number) => {
    if (!size) {
      return '0 B'
    }
    if (size < 1024) {
      return `${size} B`
    }
    if (size < 1024 * 1024) {
      return `${(size / 1024).toFixed(1)} KB`
    }
    return `${(size / (1024 * 1024)).toFixed(1)} MB`
  }

  const formatDate = (dateStr: string) => {
    const date = new Date(dateStr)
    return date.toLocaleDateString(getCurrentLocale(), {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
    })
  }
</script>

<template>
  <div class="images-list">
    <div class="list-header">
      <div class="col-preview">{{ $t('folders.fileGrid.listView.columns.preview') }}</div>
      <div class="col-filename">{{ $t('folders.fileGrid.listView.columns.filename') }}</div>
      <div class="col-size">{{ $t('folders.fileGrid.listView.columns.size') }}</div>
      <div class="col-dimensions">{{ $t('folders.fileGrid.listView.columns.dimensions') }}</div>
      <div class="col-date">{{ $t('folders.fileGrid.listView.columns.date') }}</div>
      <div class="col-actions">{{ $t('folders.fileGrid.listView.columns.actions') }}</div>
    </div>

    <div class="list-body">
      <div
        v-for="image in images"
        :key="image.id"
        class="list-item"
        :class="{
          'is-selected': (selectMode && isImageSelected?.(image)) || (batchMode && isBatchImageSelected?.(image)),
        }"
        @click="handleItemClick(image)"
      >
        <div class="col-preview" v-loading="imageLoadingStates[image.id]">
          <div class="item-preview">
            <CyberFile
              :src="getImageUrl(image.full_thumb_url, image.access_level)"
              :alt="image.display_name"
              fit-mode="cover"
              :retry-count="2"
              :is-nsfw="image.is_nsfw"
              @loading="imageLoadingStates[image.id] = $event"
              @click="$emit('preview-image', image)"
            />

            <div
              v-if="image.access_level && image.access_level !== 'public'"
              class="access-badge"
              :class="`access-${image.access_level}`"
            >
              <i :class="getAccessIcon(image.access_level)" />
            </div>
          </div>
        </div>

        <div class="col-filename">
          <div class="filename-content">
            <div class="file-name" :title="image.display_name">
              {{ image.display_name }}
            </div>
            <div class="file-meta">
              <span class="file-format">.{{ image.format?.toUpperCase() || 'UNKNOWN' }}</span>
              <span v-if="image.is_duplicate" class="duplicate-indicator" :title="$t('folders.fileGrid.meta.duplicateFile')">
                <i class="fas fa-clone" />
              </span>
            </div>
          </div>
        </div>

        <div class="col-size">
          {{ formatFileSize(image.size) }}
        </div>

        <div class="col-dimensions">
          <span v-if="image.width && image.height"> {{ image.width }} × {{ image.height }} </span>
          <span v-else class="text-content-muted">-</span>
        </div>

        <div class="col-date">
          {{ formatDate(image.created_at) }}
        </div>

        <div class="col-actions">
          <div
            v-if="
              (selectMode || batchMode) &&
              ((selectMode && isImageSelected?.(image)) || (batchMode && isBatchImageSelected?.(image)))
            "
            class="selected-indicator"
          >
            <i class="fas fa-check" />
          </div>
          <div v-else-if="!selectMode && !batchMode" class="action-buttons">
            <button
              v-if="image.ai_info"
              class="action-btn info-btn"
              :title="$t('folders.fileGrid.actions.aiInfo')"
              @click.stop="$emit('view-details', image)"
            >
              <i class="fas fa-info-circle" />
            </button>
            <button
              class="action-btn copy-btn"
              :title="$t('folders.fileGrid.actions.copyLink')"
              @click.stop="$emit('copy-link', image)"
            >
              <i class="fas fa-link" />
            </button>
            <button
              class="action-btn download-btn"
              :title="$t('folders.fileGrid.actions.download')"
              @click.stop="$emit('download-image', image)"
            >
              <i class="fas fa-download" />
            </button>
            <button
              v-if="!preview"
              class="action-btn delete-btn"
              :title="$t('folders.fileGrid.actions.delete')"
              @click.stop="$emit('delete-image', image)"
            >
              <i class="fas fa-trash" />
            </button>
          </div>
        </div>
      </div>

      <div v-if="images.length === 0" class="list-empty">
        <i class="fas fa-images" />
        <p>{{ $t('folders.fileGrid.listView.empty') }}</p>
      </div>
    </div>
  </div>
</template>

<style scoped>
  .images-list {
    background: rgba(var(--color-background-900-rgb), 0.5);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.1);
    border-radius: var(--radius-sm);
    overflow: hidden;
    box-shadow: var(--shadow-md);
  }

  .list-header {
    display: grid;
    grid-template-columns:
      [preview] 60px
      [filename] 1fr
      [size] 80px
      [dimensions] 100px
      [date] 120px
      [actions] 140px;
    gap: 16px;
    padding: 12px 16px;
    background: linear-gradient(135deg, rgba(var(--color-brand-500-rgb), 0.15) 0%, rgba(var(--color-brand-500-rgb), 0.05) 100%);
    border-bottom: 1px solid rgba(var(--color-brand-500-rgb), 0.2);
    font-size: 12px;
    font-weight: 600;
    color: rgba(var(--color-brand-500-rgb), 0.9);
    letter-spacing: 0.5px;
    text-transform: uppercase;
  }

  .list-header > div {
    display: flex;
    align-items: center;
  }

  .list-body {
    max-height: 600px;
    overflow-y: auto;
    scrollbar-width: thin;
    scrollbar-color: rgba(var(--color-brand-500-rgb), 0.3) transparent;
  }

  .list-body::-webkit-scrollbar {
    width: 6px;
  }

  .list-body::-webkit-scrollbar-track {
    background: transparent;
  }

  .list-body::-webkit-scrollbar-thumb {
    background: rgba(var(--color-brand-500-rgb), 0.3);
    border-radius: var(--radius-sm);
  }

  .list-body::-webkit-scrollbar-thumb:hover {
    background: rgba(var(--color-brand-500-rgb), 0.5);
  }

  .list-item {
    display: grid;
    grid-template-columns:
      [preview] 60px
      [filename] 1fr
      [size] 80px
      [dimensions] 100px
      [date] 120px
      [actions] 140px;
    gap: 16px;
    padding: 12px 16px;
    align-items: center;
    background: rgba(var(--color-background-800-rgb), 0.7);
    border-bottom: 1px solid rgba(var(--color-brand-500-rgb), 0.05);
    transition: all 0.2s ease;
    cursor: pointer;
    position: relative;
  }

  .list-item:hover {
    background: rgba(var(--color-background-800-rgb), 0.9);
    border-left: 3px solid rgba(var(--color-brand-500-rgb), 0.6);
    transform: translateX(2px);
    box-shadow: var(--shadow-cyber-md);
  }

  .list-item.is-selected {
  }

  .list-item.is-selected:hover {
    background: rgba(var(--color-badge-accent-text-rgb), 0.15);
    transform: translateX(2px);
  }

  .col-preview {
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .item-preview {
    position: relative;
    width: 44px;
    height: 44px;
    border-radius: var(--radius-sm);
    overflow: hidden;
    background: rgba(var(--color-background-800-rgb), 0.6);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.2);
    cursor: pointer;
    transition: all 0.2s ease;
  }

  .item-preview:hover {
    border-color: rgba(var(--color-brand-500-rgb), 0.5);
    transform: scale(1.05);
    box-shadow: var(--shadow-cyber-lg);
  }

  .item-preview img {
    width: 100%;
    height: 100%;
    object-fit: cover;
    transition: transform 0.2s ease;
  }

  .item-preview:hover img {
    transform: scale(1.1);
  }

  .access-badge {
    position: absolute;
    top: 2px;
    right: 2px;
    width: 14px;
    height: 14px;
    border-radius: var(--radius-full);
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 8px;
    color: var(--color-content-heading);
  }

  .access-badge.access-private {
    background: rgba(var(--color-error-rgb), 0.9);
  }

  .access-badge.access-protected {
    background: rgba(var(--color-warning-rgb), 0.9);
  }

  .col-filename {
    min-width: 0;
  }

  .filename-content {
    display: flex;
    flex-direction: column;
    gap: 2px;
  }

  .file-name {
    font-weight: 500;
    color: var(--color-content-heading);
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    font-size: 14px;
  }

  .file-meta {
    display: flex;
    align-items: center;
    gap: 6px;
  }

  .file-format {
    font-size: 11px;
    color: rgba(var(--color-brand-500-rgb), 0.7);
    background: rgba(var(--color-brand-500-rgb), 0.1);
    padding: 1px 4px;
    border-radius: var(--radius-sm);
    font-weight: 500;
    letter-spacing: 0.5px;
  }

  .duplicate-indicator {
    color: rgba(var(--color-info-rgb), 0.95);
    font-size: 11px;
  }

  .col-size,
  .col-dimensions,
  .col-date {
    font-size: 12px;
    color: var(--color-content-default);
    font-weight: 400;
  }

  .col-dimensions {
    font-family: 'Monaco', 'Consolas', monospace;
  }

  .text-content-muted {
    color: var(--color-content-muted);
  }

  .col-actions {
    display: flex;
    justify-content: flex-end;
  }

  .action-buttons {
    display: flex;
    gap: 6px;
    align-items: center;
  }

  .action-btn {
    width: 28px;
    height: 28px;
    border: none;
    border-radius: var(--radius-sm);
    background: var(--color-hover-bg-neutral);
    color: var(--color-content-default);
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    transition: all 0.2s ease;
    font-size: 11px;
  }

  .action-btn:hover {
    background: var(--color-hover-bg);
    transform: translateY(-1px);
    box-shadow: var(--shadow-cyber-md);
  }

  .info-btn:hover {
    color: var(--color-info-500);
    background: rgba(var(--color-info-rgb), 0.15);
  }

  .copy-btn:hover {
    color: var(--color-warning-500);
    background: rgba(var(--color-warning-rgb), 0.15);
  }

  .download-btn:hover {
    color: var(--color-success-500);
    background: rgba(var(--color-success-rgb), 0.15);
  }

  .delete-btn:hover {
    color: var(--color-pink-400);
    background: rgba(var(--color-badge-accent-text-rgb), 0.15);
  }

  .list-empty {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 60px 20px;
    color: var(--color-content-muted);
    text-align: center;
  }

  .list-empty i {
    font-size: 48px;
    color: rgba(var(--color-brand-500-rgb), 0.3);
    margin-bottom: 16px;
  }

  .list-empty p {
    font-size: 14px;
    margin: 0;
  }

  @media (max-width: 1024px) {
    .list-header,
    .list-item {
      grid-template-columns:
        [preview] 50px
        [filename] 1fr
        [size] 70px
        [dimensions] 90px
        [date] 100px
        [actions] 120px;
      gap: 12px;
    }

    .list-header,
    .list-item {
      padding: 10px 12px;
    }

    .item-preview {
      width: 36px;
      height: 36px;
    }
  }

  @media (max-width: 768px) {
    .list-header,
    .list-item {
      grid-template-columns:
        [preview] 40px
        [filename] 1fr
        [actions] 100px;
      gap: 8px;
    }

    .col-size,
    .col-dimensions,
    .col-date {
      display: none;
    }

    .action-buttons {
      gap: 4px;
    }

    .action-btn {
      width: 24px;
      height: 24px;
      font-size: 10px;
    }
  }
</style>
