<script setup lang="ts">
  import { getCurrentLocale } from '@/utils/locale'
  import { ref } from 'vue'

  const props = defineProps<{
    imageData: any
    selectMode?: boolean
    isSelected?: boolean
    isVectorSearch?: boolean
  }>()

  const emit = defineEmits(['view', 'copy', 'download', 'delete', 'visibility', 'details', 'select'])

  const isImageLoading = ref(false)

  const handleClick = (_event: Event) => {
    if (props.selectMode) {
      emit('select', props.imageData)
    } else {
      emit('view', props.imageData.id)
    }
  }

  const handleToggleVisibility = () => {
    emit('visibility', props.imageData.id, props.imageData.access_level)
  }

  const formatFileSize = (bytes: number) => {
    if (bytes === 0) {
      return '0 B'
    }
    const k = 1024
    const sizes = ['B', 'KB', 'MB', 'GB']
    const i = Math.floor(Math.log(bytes) / Math.log(k))
    return `${parseFloat((bytes / Math.pow(k, i)).toFixed(2))} ${sizes[i]}`
  }

  const _formatDate = (dateString: string) => {
    const date = new Date(dateString)
    return date.toLocaleDateString(getCurrentLocale(), {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
    })
  }
</script>

<template>
  <div class="image-card" :class="{ selected: isSelected, 'select-mode': selectMode }" @click="handleClick">
    <CyberBadge
      v-if="isVectorSearch && imageData.similarity !== undefined"
      :text="imageData.similarity?.toFixed(2)"
      color="var(--color-badge-primary-bg)"
      position="top-left"
    />
    <div v-else-if="selectMode && isSelected" class="selected-indicator">
      <i class="fas fa-check"></i>
    </div>
    <div class="image-thumbnail" v-loading="isImageLoading">
      <CyberFile
        :src="imageData.full_thumb_url"
        :alt="imageData.display_name"
        fit-mode="contain"
        :maintain-aspect-ratio="true"
        :retry-count="2"
        :is-nsfw="imageData.is_nsfw"
        :error-text="$t('resource.fileCard.loadError')"
        class="image-thumbnail-img"
        @loading="isImageLoading = $event"
      />
      <CyberAccessLevelToggle
        v-if="!selectMode"
        :access-level="imageData.access_level"
        type="image"
        size="medium"
        @toggle="handleToggleVisibility"
      />
      <DuplicateBadge
        v-if="imageData.is_duplicate"
        position="top-left"
        icon="fas fa-clone"
        :text="$t('resource.fileCard.duplicate')"
      />
      <div v-if="!selectMode" class="image-hover-overlay">
        <div class="hover-actions">
          <button
            v-if="imageData.ai_info"
            class="cyber-action-btn info-btn"
            :title="$t('resource.fileCard.aiInfo')"
            @click.stop="$emit('details', imageData.id)"
          >
            <i class="fas fa-info-circle" />
            <div class="cyber-btn-glow" />
          </button>
          <button
            class="cyber-action-btn preview-btn"
            :title="$t('resource.fileCard.preview')"
            @click.stop="$emit('view', imageData.id)"
          >
            <i class="fas fa-eye" />
            <div class="cyber-btn-glow" />
          </button>
          <button
            class="cyber-action-btn copy-btn"
            :title="$t('resource.fileCard.copyLink')"
            @click.stop="$emit('copy', imageData)"
          >
            <i class="fas fa-link" />
            <div class="cyber-btn-glow" />
          </button>
          <button
            class="cyber-action-btn download-btn"
            :title="$t('resource.fileCard.download')"
            @click.stop="$emit('download', imageData.id, imageData)"
          >
            <i class="fas fa-download" />
            <div class="cyber-btn-glow" />
          </button>
        </div>
      </div>
    </div>
    <div class="image-info">
      <div class="image-name" :title="imageData.display_name">{{ imageData.display_name }}</div>
      <div class="image-meta">
        <span class="image-size">{{ formatFileSize(imageData.size) }}</span>
        <span v-if="imageData.width && imageData.height" class="image-dimensions">
          {{ imageData.width }}x{{ imageData.height }}
        </span>
      </div>
    </div>
  </div>
</template>

<style scoped>
  .image-card {
    position: relative;
    background: rgba(var(--color-background-700-rgb), 0.7);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.15);
    border-radius: var(--radius-sm);
    overflow: visible;
    transition: all 0.3s ease;
  }

  .image-card.select-mode {
    cursor: pointer;
  }

  .image-card.select-mode:hover {
    transform: translateY(-2px);
    border-color: rgba(var(--color-brand-500-rgb), 0.3);
    box-shadow: 0 4px 15px rgba(var(--color-brand-500-rgb), 0.15);
  }

  .image-thumbnail {
    position: relative;
    cursor: pointer;
    overflow: hidden;
    aspect-ratio: 16/10;
    background: rgba(var(--color-background-800-rgb), 0.5);
    border-radius: var(--radius-sm) var(--radius-sm) 0 0;
    z-index: 1;
    perspective: 1000px;
  }

  .image-card.select-mode .image-thumbnail {
    cursor: pointer;
  }

  .image-thumbnail-img {
    width: 100%;
    height: 100%;
  }

  .image-thumbnail-img :deep(.cyber-file-container),
  .image-thumbnail-img :deep(img) {
    transition: transform 0.5s cubic-bezier(0.34, 1.56, 0.64, 1);
    transform-style: preserve-3d;
    transform-origin: center center;
  }

  .image-card:hover .image-thumbnail-img :deep(.cyber-file-container),
  .image-card:hover .image-thumbnail-img :deep(img) {
    transform: scale(1.15) translateZ(30px);
  }

  .image-hover-overlay {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background: var(--color-overlay-medium);
    display: flex;
    align-items: center;
    justify-content: center;
    opacity: 0;
    transition: opacity 0.3s ease;
  }

  .image-card:hover .image-hover-overlay {
    opacity: 1;
  }

  .hover-actions {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 0.8rem;
  }

  .cyber-action-btn {
    position: relative;
    width: 32px;
    height: 32px;
    display: flex;
    align-items: center;
    justify-content: center;
    background: rgba(var(--color-background-700-rgb), 0.95);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.4);
    border-radius: var(--radius-sm);
    color: var(--color-content-default);
    font-size: 0.9rem;
    cursor: pointer;
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    overflow: hidden;
    backdrop-filter: blur(8px);
    z-index: 1;
  }

  .cyber-btn-glow {
    position: absolute;
    top: -2px;
    left: -2px;
    right: -2px;
    bottom: -2px;
    border-radius: var(--radius-sm);
    opacity: 0;
    transition: all 0.3s ease;
    pointer-events: none;
    z-index: -1;
  }

  .cyber-action-btn:hover {
    transform: translateY(-3px) scale(1.05);
    box-shadow: 0 8px 25px rgba(var(--color-brand-500-rgb), 0.4);
    border-color: rgba(var(--color-brand-500-rgb), 0.8);
  }

  .cyber-action-btn:hover .cyber-btn-glow {
    opacity: 1;
    box-shadow:
      0 0 30px rgba(var(--color-brand-500-rgb), 0.6),
      inset 0 0 30px rgba(var(--color-brand-500-rgb), 0.15);
  }

  .cyber-action-btn:active {
    transform: translateY(-1px) scale(1.02);
    box-shadow: 0 4px 15px rgba(var(--color-brand-500-rgb), 0.3);
  }

  .cyber-action-btn.info-btn:hover {
    color: var(--color-info-500);
    border-color: rgba(var(--color-info-rgb), 0.8);
    box-shadow: 0 8px 25px rgba(var(--color-info-rgb), 0.4);
  }

  .cyber-action-btn.info-btn:hover .cyber-btn-glow {
    box-shadow:
      0 0 30px rgba(var(--color-info-rgb), 0.6),
      inset 0 0 30px rgba(var(--color-info-rgb), 0.15);
  }

  .cyber-action-btn.preview-btn:hover {
    color: var(--color-brand-500);
    border-color: rgba(var(--color-brand-500-rgb), 0.8);
    box-shadow: 0 8px 25px rgba(var(--color-brand-500-rgb), 0.4);
  }

  .cyber-action-btn.preview-btn:hover .cyber-btn-glow {
    box-shadow:
      0 0 30px rgba(var(--color-brand-500-rgb), 0.6),
      inset 0 0 30px rgba(var(--color-brand-500-rgb), 0.15);
  }

  .cyber-action-btn.copy-btn:hover {
    color: var(--color-warning-500);
    border-color: rgba(var(--color-warning-rgb), 0.8);
    box-shadow: 0 8px 25px rgba(var(--color-warning-rgb), 0.4);
  }

  .cyber-action-btn.copy-btn:hover .cyber-btn-glow {
    box-shadow:
      0 0 30px rgba(var(--color-warning-rgb), 0.6),
      inset 0 0 30px rgba(var(--color-warning-rgb), 0.15);
  }

  .cyber-action-btn.download-btn:hover {
    color: var(--color-success-500);
    border-color: rgba(var(--color-success-rgb), 0.8);
    box-shadow: 0 8px 25px rgba(var(--color-success-rgb), 0.4);
  }

  .cyber-action-btn.download-btn:hover .cyber-btn-glow {
    box-shadow:
      0 0 30px rgba(var(--color-success-rgb), 0.6),
      inset 0 0 30px rgba(var(--color-success-rgb), 0.15);
  }

  .image-info {
    position: relative;
    padding: 10px;
    flex-grow: 1;
    flex-direction: column;
    justify-content: space-between;
    z-index: 1;
    background: rgba(var(--color-background-700-rgb), 0.7);
    border-radius: 0 0 var(--radius-sm) var(--radius-sm);
  }

  .image-name {
    font-size: 0.8rem;
    font-weight: 500;
    color: var(--color-content-default);
    margin-bottom: 5px;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .image-meta {
    font-size: 0.7rem;
    color: var(--color-content-muted);
  }

  .image-size {
    margin-right: 5px;
  }

  .image-dimensions {
    margin-left: 5px;
  }
</style>
<style src="@/styles/selection.css"></style>
