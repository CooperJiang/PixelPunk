<script setup lang="ts">
  import { ref, watch } from 'vue'
  import { VueDraggable } from 'vue-draggable-plus'
  import type { ImageInfo } from '@/api/types/index'
  import { useFileContextMenu } from '../../composables/useFileContextMenu'

  defineOptions({
    name: 'FilesGridView',
  })

  interface Props {
    images: ImageInfo[]
    selectMode: boolean
    batchMode: boolean
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
    'toggle-image-visibility': [image: ImageInfo]
    'toggle-image-select': [image: ImageInfo]
    'toggle-batch-image-select': [image: ImageInfo]
    'image-drag-start': []
    'image-drag-end': [{ event: Event; sortedFiles: ImageInfo[] }]
    'image-moved': [fileId: string, targetFolderId?: string]
    'image-deleted': [fileId: string]
  }

  const props = defineProps<Props>()
  const emit = defineEmits<Emits>()
  const imageLoadingStates = ref<Record<string, boolean>>({})

  const mutableImages = ref([...props.images])

  watch(
    () => props.images,
    (newImages) => {
      mutableImages.value = [...newImages]
    },
    { deep: true }
  )

  const handleImageClick = (image: ImageInfo) => {
    if (props.selectMode && !props.preview) {
      emit('toggle-image-select', image)
    } else if (props.batchMode && !props.preview) {
      emit('toggle-batch-image-select', image)
    } else {
      emit('preview-image', image)
    }
  }

  const { showContextMenu, contextMenuPosition, contextMenuItems, showImageContextMenu, hideContextMenu } = useFileContextMenu({
    onImageMoved: (fileId, targetFolderId) => {
      emit('image-moved', fileId, targetFolderId)
    },
    onImageDeleted: (fileId) => {
      emit('image-deleted', fileId)
    },
    onPreview: (image) => emit('preview-image', image),
    onCopyLink: (image) => emit('copy-link', image),
    onDownload: (image) => emit('download-image', image),
    onDelete: (image) => emit('delete-image', image),
  })

  const getImageUrl = (url: string, _accessLevel?: string) => url

  const formatFileSize = (size: number) => {
    if (size < 1024) {
      return `${size} B`
    }
    if (size < 1024 * 1024) {
      return `${(size / 1024).toFixed(1)} KB`
    }
    return `${(size / (1024 * 1024)).toFixed(1)} MB`
  }
</script>

<template>
  <VueDraggable
    v-model="mutableImages"
    animation="150"
    ghost-class="image-ghost"
    chosen-class="image-chosen"
    drag-class="image-drag"
    :disabled="selectMode || batchMode || preview"
    class="images-grid"
    @start="$emit('image-drag-start')"
    @end="(event) => $emit('image-drag-end', { event, sortedFiles: mutableImages })"
  >
    <div
      v-for="image in mutableImages"
      :key="image.id"
      class="image-item selectable-item"
      :class="{
        selected: (selectMode && isImageSelected?.(image)) || (batchMode && isBatchImageSelected?.(image)),
        'is-selected': (selectMode && isImageSelected?.(image)) || (batchMode && isBatchImageSelected?.(image)),
        selectable: selectMode || batchMode,
      }"
      @click="handleImageClick(image)"
      @contextmenu="!preview ? showImageContextMenu($event, image) : undefined"
    >
      <div
        v-if="(selectMode && isImageSelected?.(image)) || (batchMode && isBatchImageSelected?.(image))"
        class="image-selected-indicator selectable-check"
      >
        <i class="fas fa-check" />
      </div>
      <div class="image-thumbnail" v-loading="imageLoadingStates[image.id]">
        <CyberBackground pattern="cyber">
          <CyberFile
            :src="getImageUrl(image.full_thumb_url, image.access_level)"
            :alt="image.display_name"
            fit-mode="contain"
            :maintain-aspect-ratio="true"
            :retry-count="2"
            :is-nsfw="image.is_nsfw"
            @loading="imageLoadingStates[image.id] = $event"
          />
        </CyberBackground>
        <CyberAccessLevelToggle
          v-if="!preview && !selectMode && !batchMode"
          :access-level="image.access_level"
          type="image"
          size="medium"
          @toggle="$emit('toggle-image-visibility', image)"
        />

        <cyber-file-expiry-tag
          v-if="image.is_time_limited"
          :expires-at="image.expires_at"
          :storage-duration="image.storage_duration"
          :is-time-limited="image.is_time_limited"
          position="top-left"
          mode="both"
          :show-icon="true"
        />

        <DuplicateBadge
          v-if="image.is_duplicate"
          position="top-left"
          icon="fas fa-clone"
          :text="$t('folders.fileGrid.badge.duplicate')"
        />
        <div v-if="!selectMode && !batchMode" class="image-hover-overlay">
          <div class="hover-actions">
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
      <div class="image-info">
        <div class="image-name" :title="image.display_name">{{ image.display_name }}</div>
        <div class="image-meta">
          <span class="image-size" :title="$t('folders.fileGrid.meta.size')"
            ><i class="fas fa-file-image" /> {{ formatFileSize(image.size) }}</span
          >
          <span v-if="image.width && image.height" class="image-dimensions" :title="$t('folders.fileGrid.meta.dimensions')">
            <i class="fas fa-expand" /> {{ image.width }}x{{ image.height }}
          </span>
        </div>
      </div>
    </div>
  </VueDraggable>

  <CyberContextMenu
    v-model="showContextMenu"
    :items="contextMenuItems"
    :x="contextMenuPosition.x"
    :y="contextMenuPosition.y"
    @close="hideContextMenu"
  />
</template>

<style scoped>
  .images-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(230px, 1fr));
    gap: 15px;
    perspective: 1000px;
  }

  .image-item {
    background: rgba(var(--color-background-800-rgb), 0.7);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.15);
    border-radius: var(--radius-sm);
    position: relative;
    cursor: grab;
    width: 100%;
    min-width: 0;
    transform-style: preserve-3d;
    overflow: hidden;
    transition: all 0.3s ease;
  }

  .image-item:hover {
    box-shadow:
      var(--shadow-cyber-lg),
      0 0 12px rgba(var(--color-brand-500-rgb), 0.3);
    border-color: rgba(var(--color-brand-500-rgb), 0.4);
    background: rgba(var(--color-background-800-rgb), 0.85);
  }

  .image-item:active {
    cursor: grabbing;
  }

  .image-item.selectable {
    cursor: pointer;
  }

  .image-selected-indicator {
    position: absolute;
    top: 8px;
    right: 8px;
    width: 26px;
    height: 26px;
    background: linear-gradient(135deg, var(--color-pink-400), var(--color-brand-500));
    border-radius: var(--radius-full);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 10;
    box-shadow: 0 2px 8px rgba(var(--color-badge-accent-text-rgb), 0.4);
    animation: checkmark-appear 0.3s ease-out;
    border: 2px solid rgba(var(--color-content-default-rgb), 0.9);
  }

  .image-selected-indicator i {
    color: var(--color-content-heading);
    font-size: 12px;
    text-shadow: var(--shadow-cyber-sm);
  }

  @keyframes checkmark-appear {
    0% {
      transform: scale(0);
      opacity: 0;
    }
    100% {
      transform: scale(1);
      opacity: 1;
    }
  }

  .image-thumbnail {
    aspect-ratio: 16/10;
    position: relative;
    overflow: hidden;
    perspective: 1000px;
  }

  .image-thumbnail :deep(.cyber-background) {
    width: 100%;
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .image-thumbnail :deep(.cyber-file-container) {
    width: 100%;
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .image-thumbnail :deep(.cyber-file-container),
  .image-thumbnail :deep(img) {
    transition: transform 0.5s cubic-bezier(0.34, 1.56, 0.64, 1);
    transform-style: preserve-3d;
    transform-origin: center center;
  }

  .image-thumbnail :deep(img) {
    max-width: 100%;
    max-height: 100%;
    width: auto;
    height: auto;
    object-fit: contain;
  }

  .image-item:hover .image-thumbnail :deep(.cyber-file-container),
  .image-item:hover .image-thumbnail :deep(img) {
    transform: scale(1.15) translateZ(30px);
  }

  .image-hover-overlay {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background: var(--color-overlay-medium);
    opacity: 0;
    transition: opacity 0.3s ease;
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 1;
  }

  .image-item:hover .image-hover-overlay {
    opacity: 1;
  }

  .hover-actions {
    display: flex;
    gap: 0.8rem;
  }

  .hover-actions .action-btn {
    background: rgba(var(--color-background-700-rgb), 0.8);
    width: 30px;
    height: 30px;
    color: var(--color-content-default);
    border-radius: var(--radius-sm);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.2);
    transition: all 0.2s ease;
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
  }

  .hover-actions .action-btn:hover {
    background: rgba(var(--color-background-600-rgb), 0.9);
    border-color: rgba(var(--color-brand-500-rgb), 0.4);
    transform: translateY(-2px);
    box-shadow: 0 2px 8px rgba(var(--color-brand-500-rgb), 0.3);
  }

  .info-btn:hover {
    color: var(--color-info-500);
    border-color: rgba(var(--color-info-rgb), 0.4);
  }
  .copy-btn:hover {
    color: var(--color-warning-500);
    border-color: rgba(var(--color-warning-rgb), 0.4);
  }
  .download-btn:hover {
    color: var(--color-success-500);
    border-color: rgba(var(--color-success-rgb), 0.4);
  }
  .delete-btn:hover {
    color: var(--color-pink-400);
    border-color: rgba(var(--color-badge-accent-text-rgb), 0.4);
  }

  .image-info {
    padding: 0.7rem;
  }

  .image-name {
    font-weight: 500;
    font-size: 0.85rem;
    color: rgba(var(--color-content-default-rgb), 0.9);
    margin-bottom: 0.4rem;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .image-meta {
    display: flex;
    gap: 0.6rem;
    font-size: 0.75rem;
    color: rgba(var(--color-content-muted-rgb), 0.6);
  }

  .image-meta i {
    margin-right: 0.2rem;
    color: rgba(var(--color-brand-500-rgb), 0.7);
  }

  .image-ghost {
    opacity: 0.4 !important;
    background: rgba(var(--color-brand-500-rgb), 0.15) !important;
    border: 2px dashed rgba(var(--color-brand-500-rgb), 0.6) !important;
    border-radius: var(--radius-sm) !important;
    transform: scale(0.92) !important;
    box-shadow: 0 0 20px rgba(var(--color-brand-500-rgb), 0.3) !important;
  }

  .image-chosen {
    opacity: 1 !important;
    background: rgba(var(--color-brand-500-rgb), 0.1) !important;
    border-color: rgba(var(--color-brand-500-rgb), 0.5) !important;
    transform: scale(1.03) !important;
    box-shadow: 0 5px 15px rgba(var(--color-brand-500-rgb), 0.2) !important;
    z-index: 10 !important;
  }
</style>
