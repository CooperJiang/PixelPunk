<script setup lang="ts">
  import { getCurrentLocale } from '@/utils/locale'
  import { computed } from 'vue'
  import { PLACEHOLDER_IMAGE } from '@/constants/files'
  import { useTexts } from '@/composables/useTexts'
  import type { ReviewImage } from '@/api/admin/content-review'

  const { $t } = useTexts()

  defineOptions({
    name: 'ReviewImageCard',
  })

  const props = defineProps<Props>()

  defineEmits<Emits>()

  interface Props {
    image: ReviewImage
    selected?: boolean
    selectMode?: boolean
    approving?: boolean
    rejecting?: boolean
  }

  interface Emits {
    (e: 'select', id: string): void
    (e: 'preview', image: ReviewImage): void
    (e: 'approve', id: string): void
    (e: 'reject', id: string): void
  }

  const operating = computed(() => props.approving || props.rejecting)

  const formatDate = (dateStr: string) => {
    try {
      return new Date(dateStr).toLocaleDateString(getCurrentLocale(), {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit',
      })
    } catch {
      return typeof dateStr === 'string' ? dateStr : '-'
    }
  }

  const handleImageError = (e: Event) => {
    const img = e.target as HTMLImageElement
    img.src = PLACEHOLDER_IMAGE
  }
</script>

<template>
  <div class="review-image-card" :class="{ selected: selected }">
    <div v-if="selectMode" class="select-checkbox">
      <cyberCheckbox :model-value="selected" @update:model-value="$emit('select', image.id)" />
    </div>

    <div class="image-preview" @click="$emit('preview', image)">
      <img
        :src="image.full_thumb_url || image.full_url || image.thumb_url || image.url"
        :alt="image.display_name"
        class="preview-img"
        @error="handleImageError"
      />
      <div class="image-overlay">
        <i class="fas fa-eye text-xl text-content-heading" />
      </div>
    </div>

    <div class="image-info">
      <h3 class="image-title" :title="image.display_name">
        {{ image.display_name }}
      </h3>
      <div class="image-meta">
        <span class="meta-item">
          <i class="fas fa-calendar mr-1 text-brand-600" />
          {{ formatDate(image.created_at) }}
        </span>
        <span class="meta-item">
          <i class="fas fa-ruler-combined mr-1 text-brand-600" />
          {{ image.width }}×{{ image.height }}
        </span>
        <span class="meta-item">
          <i class="fas fa-file-image mr-1 text-brand-600" />
          {{ image.format?.toUpperCase() }}
        </span>
        <span class="meta-item">
          <i class="fas fa-weight-hanging mr-1 text-brand-600" />
          {{ image.size_formatted }}
        </span>
      </div>

      <div class="nsfw-badge">
        <i class="fas fa-exclamation-triangle mr-1 text-amber-400" />
        <span class="text-amber-300">{{ $t('admin.contentReview.queue.card.sensitiveContent') }}</span>
      </div>
    </div>

    <div class="image-actions">
      <CyberButton type="success" size="small" :loading="approving" :disabled="operating" @click="$emit('approve', image.id)">
        <i class="fas fa-check mr-1" />
        {{ $t('admin.contentReview.queue.card.approve') }}
      </CyberButton>
      <CyberButton type="danger" size="small" :loading="rejecting" :disabled="operating" @click="$emit('reject', image.id)">
        <i class="fas fa-times mr-1" />
        {{ $t('admin.contentReview.queue.card.reject') }}
      </CyberButton>
    </div>
  </div>
</template>

<style scoped lang="scss">
  .review-image-card {
    position: relative;
    overflow: hidden;
    border-radius: var(--radius-sm);
    border: 1px solid var(--color-border-subtle);
    background: var(--color-background-600);
    transition: all var(--transition-normal) var(--ease-out);

    &:hover {
      border-color: var(--color-brand-400);
      box-shadow: var(--shadow-md);
    }

    &.selected {
      border-color: var(--color-brand-500);
      box-shadow: 0 0 0 2px rgba(var(--color-brand-500-rgb), 0.3);
    }
  }

  .select-checkbox {
    position: absolute;
    left: var(--space-sm);
    top: var(--space-sm);
    z-index: 10;
  }

  .cyber-checkbox {
    height: var(--space-lg);
    width: var(--space-lg);
    border-radius: var(--radius-sm);
    border: 2px solid var(--color-brand-600);
    background: transparent;

    &:checked {
      border-color: var(--color-brand-500);
      background: var(--color-brand-500);
    }

    &:focus {
      box-shadow: 0 0 0 2px rgba(var(--color-brand-500-rgb), 0.3);
    }
  }

  .image-preview {
    position: relative;
    aspect-ratio: 1/1;
    cursor: pointer;
    overflow: hidden;

    &:hover .preview-img {
      transform: scale(1.05);
    }
  }

  .preview-img {
    height: 100%;
    width: 100%;
    object-fit: cover;
    transition: transform var(--transition-normal) var(--ease-out);
  }

  .image-overlay {
    position: absolute;
    inset: 0;
    display: flex;
    align-items: center;
    justify-content: center;
    background: rgba(0, 0, 0, 0.5);
    opacity: 0;
    transition: opacity var(--transition-normal) var(--ease-out);
  }

  .image-preview:hover .image-overlay {
    opacity: 1;
  }

  .image-info {
    display: flex;
    flex-direction: column;
    gap: var(--space-sm);
    padding: var(--space-md);
  }

  .image-title {
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    font-size: var(--text-sm);
    font-weight: var(--font-medium);
    color: var(--color-white);
  }

  .image-meta {
    display: flex;
    flex-direction: column;
    gap: var(--space-xs);
  }

  .meta-item {
    display: block;
    font-size: var(--text-xs);
    color: var(--color-content-muted);
  }

  .nsfw-badge {
    display: inline-flex;
    align-items: center;
    border-radius: var(--radius-md);
    border: 1px solid rgba(var(--color-warning-rgb), 0.3);
    background: rgba(var(--color-warning-rgb), 0.2);
    padding: var(--space-xs) var(--space-sm);
    font-size: var(--text-xs);
  }

  .image-actions {
    display: flex;
    gap: var(--space-xs);
    padding: 0 var(--space-md) var(--space-md);
  }
</style>
