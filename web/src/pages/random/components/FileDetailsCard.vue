<script setup lang="ts">
  import { getCurrentLocale } from '@/utils/locale'
  import type { ImageInfo } from '@/api/types'
  import { useTexts } from '@/composables/useTexts'

  const { $t } = useTexts()

  interface Props {
    imageData: ImageInfo
  }

  defineProps<Props>()

  const formatFileSize = (bytes: number): string => {
    if (bytes === 0) {
      return '0 B'
    }
    const k = 1024
    const sizes = ['B', 'KB', 'MB', 'GB']
    const i = Math.floor(Math.log(bytes) / Math.log(k))
    return `${parseFloat((bytes / Math.pow(k, i)).toFixed(2))} ${sizes[i]}`
  }

  const formatDate = (dateString: string): string => {
    const date = new Date(dateString)
    return date.toLocaleDateString(getCurrentLocale(), {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
    })
  }
</script>

<template>
  <div class="image-details-card">
    <h3 class="section-title"><i class="fas fa-info-circle mr-2" />{{ $t('random.fileDetails.title') }}</h3>
    <div class="details-content">
      <div class="detail-row">
        <span class="detail-label">{{ $t('random.fileDetails.dimensions') }}</span>
        <span class="detail-value">{{ imageData.width }} × {{ imageData.height }}</span>
      </div>
      <div class="detail-row">
        <span class="detail-label">{{ $t('random.fileDetails.fileSize') }}</span>
        <span class="detail-value">{{ formatFileSize(imageData.size) }}</span>
      </div>
      <div class="detail-row">
        <span class="detail-label">{{ $t('random.fileDetails.format') }}</span>
        <span class="detail-value">{{ imageData.format?.toUpperCase() }}</span>
      </div>
      <div class="detail-row">
        <span class="detail-label">{{ $t('random.fileDetails.uploadTime') }}</span>
        <span class="detail-value">{{ formatDate(imageData.created_at) }}</span>
      </div>
      <div v-if="imageData.ai_info?.dominant_color" class="detail-row">
        <span class="detail-label">{{ $t('random.fileDetails.dominantColor') }}</span>
        <div class="color-display">
          <div class="color-preview" :style="{ backgroundColor: imageData.ai_info.dominant_color }" />
          <span class="color-code">{{ imageData.ai_info.dominant_color }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
  .image-details-card {
    background: rgba(var(--color-background-700-rgb), 0.35);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.2);
    border-radius: var(--radius-sm);
    padding: 1.25rem;
    backdrop-filter: blur(12px);
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
    transition: all 0.3s ease;
  }

  .image-details-card:hover {
    background: rgba(var(--color-background-700-rgb), 0.45);
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.12);
  }

  .section-title {
    font-size: 1rem;
    font-weight: 600;
    color: var(--color-brand-500);
    margin-bottom: 1rem;
    text-shadow: 0 0 10px rgba(var(--color-brand-500-rgb), 0.3);
  }

  .details-content {
    display: flex;
    flex-direction: column;
    gap: 0.75rem;
  }

  .detail-row {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 0.5rem 0;
    border-bottom: 1px solid rgba(var(--color-brand-500-rgb), 0.1);
  }

  .detail-row:last-child {
    border-bottom: none;
  }

  .detail-label {
    color: rgba(var(--color-content-rgb), 0.7);
    font-size: 0.8rem;
  }

  .detail-value {
    color: var(--color-content-heading);
    font-weight: 500;
    font-size: 0.8rem;
  }

  .color-display {
    display: flex;
    align-items: center;
    gap: 0.5rem;
  }

  .color-preview {
    width: 18px;
    height: 18px;
    border-radius: var(--radius-full);
    border: 2px solid rgba(var(--color-content-rgb), 0.3);
  }

  .color-code {
    font-family: monospace;
    font-size: 0.7rem;
    color: rgba(var(--color-content-rgb), 0.7);
  }
</style>
