<script setup lang="ts">
  import { getCurrentLocale } from '@/utils/locale'
  import type { ImageInfo } from '@/api/file'
  import { useTexts } from '@/composables/useTexts'

  defineOptions({
    name: 'RecentFileItem',
  })

  const props = defineProps<{
    image: ImageInfo
  }>()

  const emit = defineEmits<{
    click: [image: ImageInfo]
  }>()

  const getImageUrl = (): string => props.image.full_thumb_url || props.image.thumb_url || getPlaceholderImage()

  const handleImageError = (event: Event) => {
    const img = event.target as HTMLImageElement
    img.src = getPlaceholderImage()
  }

  const getPlaceholderImage = (): string =>
    'data:image/svg+xml;charset=UTF-8,%3csvg xmlns="http://www.w3.org/2000/svg" width="80" height="80" viewBox="0 0 80 80"%3e%3crect width="80" height="80" fill="%23161b22"/%3e%3cpath d="M25 35h30v10H25z" fill="%23484f58"/%3e%3ccircle cx="35" cy="30" r="3" fill="%23484f58"/%3e%3c/svg%3e'

  const { $t } = useTexts()

  const formatSize = (bytes: number): string => {
    if (bytes === 0) {
      return '0 B'
    }
    const k = 1024
    const sizes = ['B', 'KB', 'MB', 'GB']
    const i = Math.floor(Math.log(bytes) / Math.log(k))
    return `${parseFloat((bytes / Math.pow(k, i)).toFixed(1))} ${sizes[i]}`
  }

  const formatTime = (dateString: string): string => {
    const now = new Date()
    const date = new Date(dateString)
    const diff = now.getTime() - date.getTime()
    const minutes = Math.floor(diff / (1000 * 60))
    const hours = Math.floor(diff / (1000 * 60 * 60))
    const days = Math.floor(diff / (1000 * 60 * 60 * 24))

    if (minutes < 60) {
      return $t('dashboard.activityMonitor.timeAgo.minutesAgo', { minutes: String(minutes) })
    } else if (hours < 24) {
      return $t('dashboard.activityMonitor.timeAgo.hoursAgo', { hours: String(hours) })
    } else if (days < 7) {
      return $t('dashboard.activityMonitor.timeAgo.daysAgo', { days: String(days) })
    }
    return date.toLocaleDateString(getCurrentLocale(), { month: 'short', day: 'numeric' })
  }

  const handleClick = () => {
    const imageUrl = props.image.full_url || props.image.url || getImageUrl()
    window.open(imageUrl, '_blank', 'noopener,noreferrer')

    emit('click', props.image)
  }
</script>

<template>
  <div class="recent-image-item" :title="`${image.display_name} (${formatSize(image.size)})`" @click="handleClick">
    <div class="image-thumbnail">
      <CyberBackground pattern="cyber">
        <img :src="getImageUrl()" :alt="image.display_name" @error="handleImageError" />
      </CyberBackground>

      <cyber-file-expiry-tag
        v-if="image.is_time_limited"
        :expires-at="image.expires_at"
        :storage-duration="image.storage_duration"
        :is-time-limited="image.is_time_limited"
        position="top-left"
        mode="both"
        :show-icon="true"
      />

      <div class="hover-overlay">
        <div class="image-info">
          <div class="image-name">{{ image.display_name }}</div>
          <div class="image-details">
            <span class="size">{{ formatSize(image.size) }}</span>
            <span class="time">{{ formatTime(image.created_at) }}</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
  .recent-image-item {
    position: relative;
    background: var(--color-background-700);
    border: 1px solid var(--color-border-default);
    border-radius: var(--radius-md);
    overflow: hidden;
    cursor: pointer;
    transition: all 0.3s ease;
    width: 100%;
    height: 100px;
  }

  .recent-image-item:hover {
    box-shadow: var(--shadow-cyber-lg);
    z-index: 10;
  }

  .image-thumbnail {
    position: relative;
    cursor: pointer;
    overflow: hidden;
    width: 100%;
    height: 100%;
    background: var(--color-background-800);
  }

  .image-thumbnail :deep(.cyber-background) {
    width: 100%;
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .image-thumbnail img {
    max-width: 100%;
    max-height: 100%;
    width: auto;
    height: auto;
    object-fit: contain;
    transition: transform 0.3s ease;
    position: relative;
    z-index: 1;
  }

  .recent-image-item:hover .image-thumbnail img {
  }

  .hover-overlay {
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(0, 0, 0, 0.85);
    padding: 8px;
    opacity: 0;
    transition: all 0.3s ease;
    pointer-events: none;
    backdrop-filter: blur(8px);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 2;
    border-radius: var(--radius-sm);
  }

  .recent-image-item:hover .hover-overlay {
    opacity: 1;
    transform: translateY(0);
  }

  .image-info {
    display: flex;
    flex-direction: column;
    gap: 3px;
  }

  .image-name {
    font-size: 11px;
    font-weight: 700;
    color: var(--color-content-heading);
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    text-shadow:
      0 2px 4px rgba(0, 0, 0, 1),
      0 0 8px rgba(0, 0, 0, 0.8);
    line-height: 1.3;
    margin-bottom: 4px;
    text-align: center;
  }

  .image-details {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 2px;
    font-size: 9px;
    color: var(--color-content-default);
    text-shadow: 0 1px 2px rgba(0, 0, 0, 0.8);
    text-align: center;
  }

  .size {
    color: var(--color-brand-500);
    font-weight: 600;
    text-shadow: 0 1px 2px rgba(0, 0, 0, 0.8);
  }

  .time {
    color: var(--color-content-default);
    font-weight: 500;
  }

  :deep(.cyber-file-expiry-tag) {
    z-index: 5;
  }
</style>
