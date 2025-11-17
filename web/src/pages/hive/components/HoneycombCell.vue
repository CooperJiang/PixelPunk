<script setup lang="ts">
  import { computed, onMounted, ref } from 'vue'
  import type { ImageInfo } from '@/types/global'
  import { useTexts } from '@/composables/useTexts'

  const { $t } = useTexts()

  const props = defineProps<{
    imageData: ImageInfo
    position: { x: number; y: number }
    transform: string
    currentTransform?: string
    row: number
    col: number
  }>()

  const emit = defineEmits<{
    click: [imageData: ImageInfo]
    mouseenter: [row: number, col: number]
    mouseleave: [event: MouseEvent]
  }>()

  const isLoading = ref(false)
  const isLoaded = ref(false)
  const hasError = ref(false)
  const retryCount = ref(0)
  const maxRetries = 3

  const imageUrl = computed(() => {
    const baseUrl = props.imageData.full_thumb_url || props.imageData.thumb_url || ''
    if (!baseUrl) {
      return ''
    }

    if (retryCount.value > 0) {
      const separator = baseUrl.includes('?') ? '&' : '?'
      return `${baseUrl}${separator}retry=${retryCount.value}&t=${Date.now()}`
    }

    return baseUrl
  })

  const description = computed(() => {
    const aiDescription = props.imageData.ai_info?.description || ''
    const desc = aiDescription || $t('hive.cell.noDescription')
    return desc.length > 60 ? `${desc.substring(0, 60)}...` : desc
  })

  const cellStyle = computed(() => {
    const finalTransform = props.currentTransform || props.transform
    const hasHoverEffect = props.currentTransform && props.currentTransform !== props.transform

    return {
      position: 'absolute',
      left: `${props.position.x}px`,
      top: `${props.position.y}px`,
      transform: finalTransform,
      transformOrigin: 'center center',
      zIndex: hasHoverEffect ? 10 : 1,
    }
  })

  const getSpinnerStyle = (index: number) => {
    const angle = (index - 1) * 60
    const radius = 12
    const x = Math.cos((angle * Math.PI) / 180) * radius
    const y = Math.sin((angle * Math.PI) / 180) * radius

    return {
      '--x': `${x}px`,
      '--y': `${y}px`,
      animationDelay: `${(index - 1) * 0.1}s`,
    }
  }

  const formatFileSize = (bytes: number): string => {
    if (bytes === 0) {
      return '0 B'
    }
    const k = 1024
    const sizes = ['B', 'KB', 'MB', 'GB']
    const i = Math.floor(Math.log(bytes) / Math.log(k))
    return `${parseFloat((bytes / Math.pow(k, i)).toFixed(1))} ${sizes[i]}`
  }

  const handleImageLoad = () => {
    isLoading.value = false
    isLoaded.value = true
    hasError.value = false
    retryCount.value = 0
  }

  const handleImageError = () => {
    isLoading.value = false
    isLoaded.value = false
    hasError.value = true

    if (retryCount.value < maxRetries) {
      setTimeout(
        () => {
          if (hasError.value && retryCount.value < maxRetries) {
            retryLoad()
          }
        },
        1000 * (retryCount.value + 1)
      )
    }
  }

  const retryLoad = () => {
    if (retryCount.value >= maxRetries) {
      retryCount.value = 0
    }

    isLoading.value = true
    isLoaded.value = false
    hasError.value = false

    const currentUrl = imageUrl.value
    if (currentUrl) {
      retryCount.value++
    }
  }

  const handleImageClick = () => {
    if (isLoaded.value && !hasError.value) {
      emit('click', props.imageData)
    }
  }

  const handleClick = () => {
    if (hasError.value) {
      retryLoad()
    } else if (isLoaded.value) {
      emit('click', props.imageData)
    }
  }

  const handleMouseEnter = () => {
    emit('mouseenter', props.row, props.col)
  }

  const handleMouseLeave = (e: MouseEvent) => {
    emit('mouseleave', e)
  }

  onMounted(() => {
    if (!imageUrl.value) {
      isLoading.value = false
      hasError.value = true
    }
  })
</script>

<template>
  <div
    class="honeycomb-item"
    :style="cellStyle"
    @click="handleClick"
    @mouseenter="handleMouseEnter"
    @mouseleave="handleMouseLeave"
  >
    <div class="honeycomb-content">
      <div v-if="isLoading" class="honeycomb-loading">
        <div class="loading-spinner">
          <div v-for="i in 6" :key="i" class="hex-spinner" :style="getSpinnerStyle(i)" />
        </div>
        <div class="loading-text">{{ $t('hive.cell.reloading') }}</div>
      </div>

      <div v-else-if="hasError" class="honeycomb-error">
        <div class="error-icon">
          <i class="fas fa-exclamation-triangle" />
        </div>
        <div class="error-text">{{ $t('hive.cell.loadFailed') }}</div>
        <div class="error-retry" @click.stop="retryLoad">
          <i class="fas fa-redo" />
          <span>{{ $t('hive.cell.retry') }}</span>
        </div>
      </div>

      <img
        v-else-if="imageUrl"
        :src="imageUrl"
        :alt="$t('hive.cell.imageAlt')"
        class="honeycomb-image"
        :class="{ 'image-loaded': isLoaded }"
        @load="handleImageLoad"
        @error="handleImageError"
        @click.stop="handleImageClick"
      />

      <div v-if="isLoaded && !hasError" class="honeycomb-overlay">
        <div class="overlay-content">
          <div class="overlay-description">{{ description }}</div>
          <div class="overlay-meta">
            <span class="meta-item">
              <i class="fas fa-eye" />
              {{ imageData.views || 0 }}
            </span>
            <span v-if="imageData.size" class="meta-item">
              <i class="fas fa-file" />
              {{ formatFileSize(imageData.size) }}
            </span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
  .honeycomb-item {
    width: var(--item-size, 180px);
    height: calc(var(--item-size, 180px) * 1.1547);
    position: relative;
    cursor:
      url('data:image/svg+xml;utf8,<svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 20 20"><defs><linearGradient id="grad" x1="0%" y1="0%" x2="100%" y2="100%"><stop offset="0%" style="stop-color:%2300d4ff;stop-opacity:1" /><stop offset="100%" style="stop-color:%23ff1493;stop-opacity:1" /></linearGradient></defs><polygon points="10,1 17,5 17,15 10,19 3,15 3,5" fill="url(%23grad)" stroke="%23ffffff" stroke-width="1" opacity="0.8"/></svg>')
        10 10,
      pointer;
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    flex-shrink: 0;
    z-index: 1;
    opacity: 0;
    animation: fadeIn 0.8s ease forwards;
  }

  .honeycomb-loading {
    position: absolute;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    width: 100%;
    height: 100%;
    background: linear-gradient(135deg, rgba(var(--color-info-rgb), 0.1) 0%, rgba(var(--color-brand-500-rgb), 0.1) 100%);
    clip-path: polygon(
      50% calc(var(--honeycomb-clip-gap, 0) * 1%),
      calc(100% - var(--honeycomb-clip-gap, 0) * 1%) 25%,
      calc(100% - var(--honeycomb-clip-gap, 0) * 1%) 75%,
      50% calc(100% - var(--honeycomb-clip-gap, 0) * 1%),
      calc(var(--honeycomb-clip-gap, 0) * 1%) 75%,
      calc(var(--honeycomb-clip-gap, 0) * 1%) 25%
    );
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    border: 1px solid rgba(var(--color-info-rgb), 0.3);
    z-index: 2;
  }

  .loading-spinner {
    position: relative;
    width: 40px;
    height: 40px;
    margin-bottom: 8px;
  }

  .hex-spinner {
    position: absolute;
    width: 8px;
    height: 8px;
    background: var(--color-info-400);
    clip-path: polygon(50% 0%, 100% 25%, 100% 75%, 50% 100%, 0% 75%, 0% 25%);
    animation: hexSpinner 1.5s ease-in-out infinite;
    left: 50%;
    top: 50%;
    margin-left: -4px;
    margin-top: -4px;
  }

  @keyframes hexSpinner {
    0%,
    100% {
      opacity: 0.3;
      transform: translate(var(--x, 0), var(--y, 0)) scale(0.8);
    }
    50% {
      opacity: 1;
      transform: translate(var(--x, 0), var(--y, 0)) scale(1.2);
    }
  }

  .loading-text {
    font-size: 0.7rem;
    color: var(--color-info-400);
    font-weight: bold;
    text-shadow: 0 1px 2px rgba(0, 0, 0, 0.8);
  }

  .honeycomb-error {
    position: absolute;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    width: 100%;
    height: 100%;
    background: linear-gradient(135deg, rgba(var(--color-error-rgb), 0.1) 0%, rgba(var(--color-error-rgb), 0.1) 100%);
    clip-path: polygon(
      50% calc(var(--honeycomb-clip-gap, 0) * 1%),
      calc(100% - var(--honeycomb-clip-gap, 0) * 1%) 25%,
      calc(100% - var(--honeycomb-clip-gap, 0) * 1%) 75%,
      50% calc(100% - var(--honeycomb-clip-gap, 0) * 1%),
      calc(var(--honeycomb-clip-gap, 0) * 1%) 75%,
      calc(var(--honeycomb-clip-gap, 0) * 1%) 25%
    );
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    border: 1px solid rgba(var(--color-error-rgb), 0.4);
    cursor: pointer;
    z-index: 2;
    transition: all 0.3s ease;
  }

  .honeycomb-error {
    background: linear-gradient(135deg, rgba(var(--color-error-rgb), 0.15) 0%, rgba(var(--color-error-rgb), 0.15) 100%);
    border-color: rgba(var(--color-error-rgb), 0.6);
    box-shadow: 0 0 15px rgba(var(--color-error-rgb), 0.3);
  }

  .error-icon {
    font-size: 1.2rem;
    color: var(--color-red-400);
    margin-bottom: 4px;
    animation: errorPulse 2s ease-in-out infinite;
  }

  @keyframes errorPulse {
    0%,
    100% {
      opacity: 0.7;
      transform: scale(1);
    }
    50% {
      opacity: 1;
      transform: scale(1.1);
    }
  }

  .error-text {
    font-size: 0.7rem;
    color: var(--color-red-400);
    font-weight: bold;
    text-shadow: 0 1px 2px rgba(0, 0, 0, 0.8);
    margin-bottom: 6px;
  }

  .error-retry {
    display: flex;
    align-items: center;
    gap: 4px;
    font-size: 0.6rem;
    color: var(--color-error-300);
    background: rgba(var(--color-error-rgb), 0.2);
    padding: 4px 8px;
    border-radius: var(--radius-sm);
    border: 1px solid rgba(var(--color-error-rgb), 0.3);
    cursor: pointer;
    transition: all 0.3s ease;
    font-weight: bold;
  }

  .error-retry:hover {
    background: rgba(var(--color-error-rgb), 0.3);
    border-color: rgba(var(--color-error-rgb), 0.5);
    color: var(--color-error-200);
    transform: scale(1.05);
  }

  .error-retry i {
    font-size: 0.5rem;
    animation: retryRotate 2s linear infinite;
  }

  @keyframes retryRotate {
    from {
      transform: rotate(0deg);
    }
    to {
      transform: rotate(360deg);
    }
  }

  .honeycomb-content {
    width: 100%;
    height: 100%;
    position: relative;
  }

  .honeycomb-image {
    width: 100%;
    height: 100%;
    object-fit: cover;
    object-position: center;
    clip-path: polygon(
      50% calc(var(--honeycomb-clip-gap, 0) * 1%),
      calc(100% - var(--honeycomb-clip-gap, 0) * 1%) 25%,
      calc(100% - var(--honeycomb-clip-gap, 0) * 1%) 75%,
      50% calc(100% - var(--honeycomb-clip-gap, 0) * 1%),
      calc(var(--honeycomb-clip-gap, 0) * 1%) 75%,
      calc(var(--honeycomb-clip-gap, 0) * 1%) 25%
    );
    transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
    position: absolute;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    border: 1px solid rgba(var(--color-info-rgb), 0.2);
    opacity: 0;
  }

  .honeycomb-image.image-loaded {
    opacity: 1;
    animation: imageSlideIn 0.6s ease-out forwards;
  }

  @keyframes imageSlideIn {
    from {
      opacity: 0;
      transform: translate(-50%, -50%) scale(0.9);
    }
    to {
      opacity: 1;
      transform: translate(-50%, -50%) scale(1);
    }
  }

  .honeycomb-item:hover .honeycomb-image {
    filter: brightness(1.05) saturate(1.1) contrast(1.05);
    border-color: rgba(var(--color-info-rgb), 0.5);
    box-shadow: 0 0 20px rgba(var(--color-info-rgb), 0.3);
  }

  .honeycomb-overlay {
    position: absolute;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%) scale(0.9);
    width: 100%;
    height: 100%;
    background: rgba(0, 0, 0, 0.85);
    clip-path: polygon(
      50% calc(var(--honeycomb-clip-gap, 0) * 1%),
      calc(100% - var(--honeycomb-clip-gap, 0) * 1%) 25%,
      calc(100% - var(--honeycomb-clip-gap, 0) * 1%) 75%,
      50% calc(100% - var(--honeycomb-clip-gap, 0) * 1%),
      calc(var(--honeycomb-clip-gap, 0) * 1%) 75%,
      calc(var(--honeycomb-clip-gap, 0) * 1%) 25%
    );
    opacity: 0;
    transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
    display: flex;
    align-items: center;
    justify-content: center;
    color: var(--color-content-heading);
    font-weight: bold;
    text-align: center;
    pointer-events: none;
  }

  .honeycomb-item:hover .honeycomb-overlay {
    opacity: 1;
    transform: translate(-50%, -50%) scale(1);
  }

  .overlay-content {
    width: 85%;
    height: 80%;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 8px;
    box-sizing: border-box;
  }

  .overlay-description {
    font-size: 0.7rem;
    line-height: 1.3;
    color: var(--color-content);
    text-shadow: 0 1px 2px rgba(0, 0, 0, 0.8);
    margin-bottom: 8px;
    overflow: hidden;
    display: -webkit-box;
    -webkit-line-clamp: 4;
    -webkit-box-orient: vertical;
    word-break: break-word;
    text-align: center;
  }

  .overlay-meta {
    display: flex;
    gap: 8px;
    font-size: 0.6rem;
    color: var(--color-error-500);
  }

  .meta-item {
    display: flex;
    align-items: center;
    gap: 2px;
  }

  .meta-item i {
    font-size: 0.5rem;
  }

  @keyframes fadeIn {
    from {
      opacity: 0;
    }
    to {
      opacity: 1;
    }
  }

  @media (max-width: 768px) {
    .overlay-description {
      font-size: 0.65rem;
      -webkit-line-clamp: 3;
    }

    .overlay-meta {
      font-size: 0.55rem;
    }
  }
</style>
