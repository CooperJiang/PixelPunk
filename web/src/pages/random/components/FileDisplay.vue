<script setup lang="ts">
  import { onMounted, onUnmounted, ref, computed, watch } from 'vue'
  import type { ImageInfo } from '@/api/types'
  import EnhancedFilePreview from '@/components/EnhancedFilePreview'
  import { useTexts } from '@/composables/useTexts'

  const { $t } = useTexts()

  interface Props {
    imageData: ImageInfo
  }

  interface Emits {
    (e: 'image-load'): void
    (e: 'image-error'): void
  }

  const props = defineProps<Props>()
  const emit = defineEmits<Emits>()

  const currentImageUrl = ref('')
  const isProgressiveLoading = ref(true)

  const imageFitMode = computed(() => {
    if (!props.imageData) return 'contain'

    const width = props.imageData.width || 0
    const height = props.imageData.height || 0
    const aspectRatio = width / height

    if (aspectRatio < 1) {
      return 'contain'
    }

    return 'cover'
  })

  const isPortrait = computed(() => {
    if (!props.imageData) return false
    return (props.imageData.width || 0) < (props.imageData.height || 0)
  })

  const loadProgressiveImage = () => {
    currentImageUrl.value = props.imageData.full_thumb_url || props.imageData.full_url
    isProgressiveLoading.value = true

    const img = new Image()
    img.onload = () => {
      currentImageUrl.value = props.imageData.full_url
      isProgressiveLoading.value = false
      emit('image-load')
    }
    img.onerror = () => {
      isProgressiveLoading.value = false
      emit('image-error')
    }
    img.src = props.imageData.full_url
  }

  watch(
    () => props.imageData,
    () => {
      if (props.imageData) {
        loadProgressiveImage()
      }
    },
    { immediate: true }
  )

  const isOverlayVisible = ref(false)
  const isFullscreen = ref(false)

  const formatFileSize = (bytes: number): string => {
    if (bytes === 0) {
      return '0 B'
    }
    const k = 1024
    const sizes = ['B', 'KB', 'MB', 'GB']
    const i = Math.floor(Math.log(bytes) / Math.log(k))
    return `${parseFloat((bytes / Math.pow(k, i)).toFixed(2))} ${sizes[i]}`
  }

  const handleMouseEnter = () => {
    isOverlayVisible.value = true
  }

  const handleMouseLeave = () => {
    isOverlayVisible.value = false
  }

  const handleOverlayEnter = () => {
    isOverlayVisible.value = true
  }

  const openFullscreen = () => {
    isFullscreen.value = true
  }

  const closeFullscreen = () => {
    isFullscreen.value = false
  }

  const handleKeydown = (event: KeyboardEvent) => {
    if (event.key === 'Escape' && isFullscreen.value) {
      closeFullscreen()
    }
  }

  onMounted(() => {
    document.addEventListener('keydown', handleKeydown)
  })

  onUnmounted(() => {
    document.removeEventListener('keydown', handleKeydown)
  })
</script>

<template>
  <div
    class="fullscreen-image-container"
    :class="{
      'progressive-loading': isProgressiveLoading,
      'portrait-mode': isPortrait,
    }"
    @mouseenter="handleMouseEnter"
    @mouseleave="handleMouseLeave"
  >
    <CyberFile
      :src="currentImageUrl"
      :alt="imageData.display_name"
      class="hero-image"
      width="100%"
      height="100%"
      :fit-mode="imageFitMode"
      :retry-count="3"
      :is-nsfw="imageData.is_nsfw"
      background-pattern="none"
      @click="openFullscreen"
    />

    <div v-if="isProgressiveLoading" class="progressive-indicator">
      <i class="fas fa-circle-notch fa-spin" />
      <span>{{ $t('random.fileDisplay.loadingHD') }}</span>
    </div>

    <div class="image-info-overlay" :class="{ show: isOverlayVisible }" @mouseenter="handleOverlayEnter">
      <div class="overlay-content">
        <h1 class="image-title">{{ imageData.display_name }}</h1>
        <div class="image-quick-info">
          <span class="info-badge">{{ imageData.width }} × {{ imageData.height }}</span>
          <span class="info-badge">{{ formatFileSize(imageData.size) }}</span>
          <span class="info-badge">{{ imageData.format?.toUpperCase() }}</span>
        </div>
        <div class="image-hint">
          <span class="hint-text">{{ $t('random.fileDisplay.clickHint') }}</span>
        </div>
      </div>
    </div>

    <EnhancedFilePreview
      :visible="isFullscreen"
      :file-url="imageData.full_url"
      :file-name="imageData.display_name"
      :file-width="imageData.width"
      :file-height="imageData.height"
      @update:visible="closeFullscreen"
      @close="closeFullscreen"
    />
  </div>
</template>

<style scoped>
  .fullscreen-image-container {
    width: 100%;
    height: 100%;
    position: relative;
    cursor: pointer;
    overflow: hidden;
  }

  .hero-image {
    width: 100%;
    height: 100%;
  }

  .hero-image :deep(.cyber-file-container) {
    width: 100% !important;
    height: 100% !important;
  }

  .hero-image :deep(img) {
    width: 100% !important;
    height: 100% !important;
    object-fit: cover !important;
  }

  .portrait-mode {
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .portrait-mode .hero-image :deep(.cyber-file-container) {
    background: none !important;
  }

  .portrait-mode .hero-image :deep(.cyber-file-container)::before,
  .portrait-mode .hero-image :deep(.cyber-file-container)::after {
    display: none !important;
  }

  .portrait-mode .hero-image :deep(img) {
    object-fit: contain !important;
    width: auto !important;
    max-width: 100% !important;
    height: 100% !important;
  }

  .fullscreen-image-container:hover .hero-image :deep(img) {
    transform: scale(1.01);
  }

  .portrait-mode:hover .hero-image :deep(img) {
    transform: none;
  }

  .image-info-overlay {
    position: absolute;
    bottom: 0;
    left: 0;
    right: 0;
    background: linear-gradient(transparent, rgba(var(--color-background-900-rgb), 0.8));
    padding: 2rem 1.5rem 1.5rem;
    transform: translateY(100%);
    transition: transform 0.3s ease;
    pointer-events: none;
  }

  .image-info-overlay.show {
    transform: translateY(0);
    pointer-events: auto;
  }

  .overlay-content {
    max-width: 1200px;
    margin: 0 auto;
  }

  .image-title {
    font-size: 1.2rem;
    font-weight: bold;
    color: var(--color-content-heading);
    margin-bottom: 0.75rem;
    text-shadow: 0 2px 10px rgba(var(--color-background-900-rgb), 0.5);
  }

  .image-quick-info {
    display: flex;
    flex-wrap: wrap;
    gap: 0.75rem;
    margin-bottom: 0.5rem;
  }

  .info-badge {
    background: rgba(var(--color-background-900-rgb), 0.7);
    color: var(--color-brand-500);
    padding: 0.25rem 0.75rem;
    border-radius: var(--radius-sm);
    font-size: 0.8rem;
    font-weight: 500;
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.3);
    backdrop-filter: blur(10px);
  }

  .quality-badge {
    background: rgba(var(--color-warning-rgb), 0.2);
    color: var(--color-warning-500);
    padding: 0.25rem 0.75rem;
    border-radius: var(--radius-sm);
    font-size: 0.8rem;
    font-weight: 500;
    border: 1px solid rgba(var(--color-warning-rgb), 0.4);
    backdrop-filter: blur(10px);
  }

  .image-hint {
    margin-top: 0.5rem;
  }

  .hint-text {
    color: rgba(var(--color-content-rgb), 0.7);
    font-size: 0.75rem;
    font-style: italic;
  }

  .progressive-indicator {
    position: absolute;
    bottom: 2rem;
    right: 2rem;
    display: flex;
    align-items: center;
    gap: 0.5rem;
    background: rgba(var(--color-background-900-rgb), 0.8);
    backdrop-filter: blur(10px);
    padding: 0.5rem 1rem;
    border-radius: var(--radius-sm);
    color: var(--color-brand-500);
    font-size: 0.85rem;
    font-weight: 500;
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.3);
    z-index: 5;
  }

  .progressive-indicator i {
    font-size: 1rem;
  }

  .progressive-loading .hero-image {
    filter: blur(2px);
    transition: filter 0.3s ease;
  }

  @media (max-width: 768px) {
    .fullscreen-image-container {
      height: calc(100vh - 3.5rem);
    }

    .hero-image {
      max-height: calc(100vh - 3.5rem);
    }

    .image-title {
      font-size: 1rem;
    }

    .image-quick-info {
      gap: 0.5rem;
    }
  }
</style>
