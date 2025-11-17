<script setup lang="ts">
  import { getCurrentLocale } from '@/utils/locale'
  import { ref, watch } from 'vue'
  import type { ImageInfo } from '@/types/global'

  const props = defineProps<{
    show: boolean
    imageData: ImageInfo | null
  }>()

  const emit = defineEmits<{
    close: []
  }>()

  const loading = ref(false)
  const imageLoaded = ref(false)

  /* 监听文件数据变化，预加载文件 */
  watch(
    () => props.imageData,
    (newImageData) => {
      if (newImageData && props.show) {
        loadImage(newImageData.full_url)
      }
    },
    { immediate: true }
  )

  watch(
    () => props.show,
    (newShow) => {
      if (newShow && props.imageData) {
        loadImage(props.imageData.full_url)
      } else {
        resetState()
      }
    }
  )

  const loadImage = (url: string) => {
    loading.value = true
    imageLoaded.value = false

    const img = new Image()
    img.onload = () => {
      loading.value = false
      imageLoaded.value = true
    }
    img.onerror = () => {
      loading.value = false
      imageLoaded.value = false
    }
    img.src = url
  }

  const resetState = () => {
    loading.value = false
    imageLoaded.value = false
  }

  const closePreview = () => {
    emit('close')
  }

  const handleOverlayClick = () => {
    closePreview()
  }

  const downloadImage = async () => {
    if (!props.imageData || !imageLoaded.value) {
      return
    }

    try {
      const response = await fetch(props.imageData.full_url)
      const blob = await response.blob()
      const url = window.URL.createObjectURL(blob)

      const link = document.createElement('a')
      link.href = url
      link.download = props.imageData.display_name || props.imageData.original_name || 'image'
      document.body.appendChild(link)
      link.click()
      document.body.removeChild(link)

      window.URL.revokeObjectURL(url)
    } catch {
      return
    }
  }

  const formatFileSize = (bytes: number): string => {
    if (bytes === 0) {
      return '0 Bytes'
    }
    const k = 1024
    const sizes = ['Bytes', 'KB', 'MB', 'GB']
    const i = Math.floor(Math.log(bytes) / Math.log(k))
    return `${parseFloat((bytes / Math.pow(k, i)).toFixed(2))} ${sizes[i]}`
  }

  const formatDate = (dateString: string): string => {
    const date = new Date(dateString)
    return date.toLocaleString(getCurrentLocale(), {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit',
    })
  }
</script>

<template>
  <div v-if="show" class="image-preview-overlay" @click="handleOverlayClick">
    <div class="preview-container" @click.stop>
      <div class="preview-header">
        <div class="preview-title">
          <div class="title-icon">
            <i class="fas fa-image" />
          </div>
          <span>{{ $t('hive.preview.title') }}</span>
        </div>
        <div class="header-actions">
          <button class="action-btn download-btn" :disabled="!imageLoaded" @click="downloadImage">
            <i class="fas fa-download" />
          </button>
          <button class="action-btn close-btn" @click="closePreview">
            <i class="fas fa-times" />
          </button>
        </div>
      </div>

      <div class="preview-content">
        <div class="preview-image-section">
          <div class="preview-image-wrapper">
            <div v-if="loading" class="image-loading">
              <div class="cyber-loader">
                <div class="loader-ring" />
                <div class="loader-ring" />
                <div class="loader-ring" />
                <div class="loader-core" />
              </div>
              <div class="loading-text">
                <span class="loading-dots">{{ $t('hive.preview.loading') }}</span>
                <span class="dots">...</span>
              </div>
            </div>

            <img
              v-else-if="imageLoaded"
              :src="imageData?.full_url"
              :alt="imageData?.display_name || $t('hive.preview.imageAlt')"
              class="preview-image"
            />

            <div v-else class="image-error">
              <div class="error-icon">
                <i class="fas fa-exclamation-triangle" />
              </div>
              <div class="error-text">{{ $t('hive.preview.error') }}</div>
            </div>
          </div>

          <div v-if="imageData" class="preview-info">
            <div class="info-grid">
              <div class="info-item">
                <span class="info-label">{{ $t('hive.preview.info.name') }}</span>
                <span class="info-value">{{
                  imageData.display_name || imageData.original_name || $t('hive.preview.info.unnamed')
                }}</span>
              </div>
              <div v-if="imageData.ai_info?.description" class="info-item">
                <span class="info-label">{{ $t('hive.preview.info.description') }}</span>
                <span class="info-value">{{ imageData.ai_info.description }}</span>
              </div>
              <div v-if="imageData.width && imageData.height" class="info-item">
                <span class="info-label">{{ $t('hive.preview.info.size') }}</span>
                <span class="info-value">{{ imageData.width }}×{{ imageData.height }}</span>
              </div>
              <div class="info-item">
                <span class="info-label">{{ $t('hive.preview.info.fileSize') }}</span>
                <span class="info-value">{{ formatFileSize(imageData.size) }}</span>
              </div>
              <div v-if="imageData.created_at" class="info-item">
                <span class="info-label">{{ $t('hive.preview.info.createTime') }}</span>
                <span class="info-value">{{ formatDate(imageData.created_at) }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
  .image-preview-overlay {
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background: rgba(0, 0, 0, 0.6);
    backdrop-filter: blur(12px);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: var(--z-index-modal);
    animation: fadeIn 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  }

  @keyframes fadeIn {
    from {
      opacity: 0;
      backdrop-filter: blur(0px);
    }
    to {
      opacity: 1;
      backdrop-filter: blur(12px);
    }
  }

  .preview-container {
    background: rgba(var(--color-background-900-rgb), 0.75);
    backdrop-filter: blur(20px);
    border: 1px solid rgba(var(--color-info-rgb), 0.3);
    border-radius: var(--radius-sm);
    width: 90vw;
    max-width: 1000px;
    max-height: 90vh;
    overflow: hidden;
    box-shadow:
      0 0 30px rgba(var(--color-info-rgb), 0.2),
      0 0 60px rgba(var(--color-brand-500-rgb), 0.1),
      inset 0 1px 0 rgba(var(--color-content-rgb), 0.05);
    animation: slideIn 0.4s cubic-bezier(0.34, 1.56, 0.64, 1);
    display: flex;
    flex-direction: column;
  }

  @keyframes slideIn {
    from {
      transform: scale(0.9) translateY(30px);
      opacity: 0;
    }
    to {
      transform: scale(1) translateY(0);
      opacity: 1;
    }
  }

  .preview-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 15px 20px;
    background: rgba(var(--color-info-rgb), 0.08);
    border-bottom: 1px solid rgba(var(--color-info-rgb), 0.15);
    backdrop-filter: blur(10px);
    flex-shrink: 0;
  }

  .preview-title {
    display: flex;
    align-items: center;
    gap: 10px;
    font-size: 1.1rem;
    font-weight: 600;
    color: var(--color-white);
    text-shadow: 0 0 8px rgba(var(--color-info-rgb), 0.4);
  }

  .title-icon {
    width: 28px;
    height: 28px;
    background: linear-gradient(45deg, var(--color-info-400), var(--color-error-500));
    border-radius: var(--radius-sm);
    display: flex;
    align-items: center;
    justify-content: center;
    box-shadow: 0 0 15px rgba(var(--color-info-rgb), 0.3);
  }

  .title-icon i {
    color: var(--color-content-heading);
    font-size: 0.9rem;
  }

  .header-actions {
    display: flex;
    gap: 10px;
  }

  .action-btn {
    background: var(--color-hover-bg-neutral);
    border: 1px solid var(--color-border-default);
    color: var(--color-content);
    padding: 8px 10px;
    border-radius: var(--radius-sm);
    cursor: pointer;
    font-size: 0.9rem;
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    backdrop-filter: blur(10px);
    min-width: 36px;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .action-btn:hover:not(:disabled) {
    background: rgba(var(--color-info-rgb), 0.2);
    border-color: rgba(var(--color-info-rgb), 0.4);
    color: var(--color-info-400);
    box-shadow: 0 0 15px rgba(var(--color-info-rgb), 0.3);
    transform: translateY(-1px);
  }

  .action-btn:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  .download-btn:hover:not(:disabled) {
    background: rgba(var(--color-success-rgb), 0.2);
    border-color: rgba(var(--color-success-rgb), 0.4);
    color: var(--color-success-500);
    box-shadow: 0 0 15px rgba(var(--color-success-rgb), 0.3);
  }

  .close-btn:hover {
    background: rgba(var(--color-error-rgb), 0.2);
    border-color: rgba(var(--color-error-rgb), 0.4);
    color: var(--color-error-500);
    box-shadow: 0 0 15px rgba(var(--color-error-rgb), 0.3);
  }

  .preview-content {
    flex: 1;
    overflow-y: auto;
    overflow-x: hidden;
    padding: 0;
    scrollbar-width: thin;
    scrollbar-color: rgba(0, 212, 255, 0.3) transparent;
  }

  .preview-content::-webkit-scrollbar {
    width: 6px;
  }

  .preview-content::-webkit-scrollbar-track {
    background: transparent;
  }

  .preview-content::-webkit-scrollbar-thumb {
    background: rgba(0, 212, 255, 0.3);
    border-radius: var(--radius-sm);
  }

  .preview-content::-webkit-scrollbar-thumb:hover {
    background: rgba(0, 212, 255, 0.5);
  }

  .preview-image-section {
    display: flex;
    flex-direction: column;
  }

  .preview-image-wrapper {
    padding: 30px;
    display: flex;
    justify-content: center;
    align-items: center;
    background: rgba(0, 0, 0, 0.2);
    min-height: 400px;
    position: relative;
  }

  .preview-image {
    max-width: 100%;
    max-height: 60vh;
    border-radius: var(--radius-sm);
    box-shadow:
      0 0 20px rgba(0, 212, 255, 0.2),
      0 0 40px rgba(255, 20, 147, 0.1);
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    border: 1px solid rgba(255, 255, 255, 0.1);
  }

  .preview-image:hover {
    box-shadow:
      0 0 25px rgba(0, 212, 255, 0.3),
      0 0 50px rgba(255, 20, 147, 0.15);
    transform: scale(1.01);
  }

  .image-loading {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: 25px;
    width: 100%;
    height: 100%;
    position: absolute;
    top: 0;
    left: 0;
  }

  .cyber-loader {
    position: relative;
    width: 60px;
    height: 60px;
  }

  .loader-ring {
    position: absolute;
    width: 100%;
    height: 100%;
    border: 2px solid transparent;
    border-radius: var(--radius-full);
    animation: spin 3s linear infinite;
  }

  .loader-ring:nth-child(1) {
    border-top-color: var(--color-info-400);
    animation-delay: 0s;
    animation-duration: 2s;
  }

  .loader-ring:nth-child(2) {
    border-right-color: var(--color-error-500);
    animation-delay: 0.5s;
    animation-duration: 2.5s;
  }

  .loader-ring:nth-child(3) {
    border-bottom-color: var(--color-info-400);
    animation-delay: 1s;
    animation-duration: 3s;
  }

  .loader-core {
    position: absolute;
    top: 50%;
    left: 50%;
    width: 16px;
    height: 16px;
    background: linear-gradient(45deg, var(--color-info-400), var(--color-error-500));
    border-radius: var(--radius-full);
    transform: translate(-50%, -50%);
    animation: pulse 2s ease-in-out infinite;
  }

  @keyframes spin {
    0% {
      transform: rotate(0deg);
    }
    100% {
      transform: rotate(360deg);
    }
  }

  @keyframes pulse {
    0%,
    100% {
      transform: translate(-50%, -50%) scale(1);
      box-shadow: 0 0 15px rgba(0, 212, 255, 0.4);
    }
    50% {
      transform: translate(-50%, -50%) scale(1.2);
      box-shadow: 0 0 25px rgba(255, 20, 147, 0.6);
    }
  }

  .loading-text {
    font-size: 1.1rem;
    color: var(--color-info-400);
    font-weight: 500;
    text-shadow: 0 0 10px rgba(0, 212, 255, 0.5);
  }

  .dots {
    animation: dots 1.5s infinite;
  }

  @keyframes dots {
    0%,
    20% {
      opacity: 0;
    }
    50% {
      opacity: 1;
    }
    100% {
      opacity: 0;
    }
  }

  .image-error {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: 15px;
    color: var(--color-red-400);
    width: 100%;
    height: 100%;
    position: absolute;
    top: 0;
    left: 0;
  }

  .error-icon {
    font-size: 2.5rem;
    opacity: 0.7;
  }

  .error-text {
    font-size: 1rem;
    font-weight: 500;
  }

  .preview-info {
    padding: 25px 30px 30px;
    background: rgba(0, 0, 0, 0.3);
    border-top: 1px solid rgba(0, 212, 255, 0.15);
  }

  .info-grid {
    display: grid;
    gap: 15px;
    max-width: 100%;
  }

  .info-item {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    padding: 12px 16px;
    background: rgba(255, 255, 255, 0.03);
    border: 1px solid rgba(255, 255, 255, 0.08);
    border-radius: var(--radius-sm);
    transition: all 0.3s ease;
  }

  .info-item:hover {
    background: rgba(0, 212, 255, 0.05);
    border-color: rgba(0, 212, 255, 0.2);
    transform: translateY(-1px);
  }

  .info-label {
    font-size: 0.9rem;
    font-weight: 600;
    color: var(--color-info-400);
    min-width: 70px;
    text-shadow: 0 0 5px rgba(0, 212, 255, 0.3);
  }

  .info-value {
    font-size: 0.9rem;
    color: var(--color-content);
    text-align: right;
    flex: 1;
    margin-left: 15px;
    word-break: break-word;
    line-height: 1.4;
  }

  @media (max-width: 768px) {
    .preview-container {
      width: 95vw;
      max-height: 95vh;
      border-radius: var(--radius-sm);
    }

    .preview-header {
      padding: 12px 16px;
    }

    .preview-title {
      font-size: 1rem;
    }

    .title-icon {
      width: 24px;
      height: 24px;
    }

    .preview-image-wrapper {
      padding: 20px;
      min-height: 300px;
    }

    .preview-info {
      padding: 20px;
    }

    .info-item {
      flex-direction: column;
      gap: 6px;
      text-align: left;
    }

    .info-value {
      text-align: left;
      margin-left: 0;
    }

    .action-btn {
      padding: 6px 8px;
      min-width: 32px;
    }
  }
</style>
