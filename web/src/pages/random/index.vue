<script setup lang="ts">
  import { onMounted, onUnmounted, ref } from 'vue'
  import { getRandomRecommendedFile } from '@/api/file'
  import type { ImageInfo } from '@/api/types'
  import { useToast } from '@/components/Toast/useToast'
  import { useTexts } from '@/composables/useTexts'

  import ErrorState from './components/ErrorState.vue'
  import FileDisplay from './components/FileDisplay.vue'
  import CreatorCard from './components/CreatorCard.vue'
  import AIAnalysisCard from './components/AIAnalysisCard.vue'
  import FileDetailsCard from './components/FileDetailsCard.vue'
  import ActionButtons from './components/ActionButtons.vue'
  import RefreshButton from './components/RefreshButton.vue'
  import SkeletonLoader from './components/SkeletonLoader.vue'

  const loading = ref(false)
  const error = ref<string>('')
  const imageData = ref<ImageInfo | null>(null)

  const toast = useToast()
  const { $t } = useTexts()

  const _siteDomain = import.meta.env.VITE_SITE_DOMAIN || 'http://localhost:5173'

  const goToAuthorPage = () => {
    if (imageData.value?.user_info?.id) {
      window.open(`/author/${imageData.value.user_info.id}`, '_blank')
    }
  }

  const viewWork = (work: any) => {
    window.open(work.full_url, '_blank')
  }

  const fetchRandomImage = async () => {
    loading.value = true
    error.value = ''
    imageData.value = null

    try {
      const result = await getRandomRecommendedFile()

      if (result.success) {
        const newImageData = result.data

        if (!newImageData) {
          error.value = result.message || $t('random.error.noFiles')
          loading.value = false
          return
        }

        imageData.value = newImageData
      } else {
        error.value = result.message || $t('random.error.fetchFailed')
        loading.value = false
      }
    } catch (err: any) {
      const is404 = err?.response?.status === 404

      if (is404) {
        error.value = $t('random.error.noFiles')
      } else {
        error.value = $t('random.error.fetchFailed')
      }

      loading.value = false
    }
  }

  const refreshImage = () => {
    fetchRandomImage()
  }

  const handleKeydown = (e: KeyboardEvent) => {
    const activeElement = document.querySelector(':focus') as HTMLElement
    if (
      activeElement &&
      (activeElement.tagName === 'INPUT' || activeElement.tagName === 'TEXTAREA' || activeElement.contentEditable === 'true')
    ) {
      return
    }

    const fullscreenOverlay = document.querySelector('.enhanced-fullscreen-overlay')
    if (fullscreenOverlay) {
      return // 全屏模式下不处理这些快捷键
    }

    switch (e.key) {
      case ' ':
      case 'Enter':
        e.preventDefault()
        fetchRandomImage()
        break
      case 'f':
      case 'F':
        e.preventDefault()
        break
    }
  }

  const onImageLoad = () => {
    loading.value = false
  }

  const onImageError = () => {
    toast.error($t('random.error.fileLoadFailed'))
    loading.value = false
  }

  onMounted(() => {
    fetchRandomImage()
    document.addEventListener('keydown', handleKeydown)
  })

  onUnmounted(() => {
    document.removeEventListener('keydown', handleKeydown)
  })
</script>

<template>
  <div class="random-page-container">
    <RefreshButton :loading="loading" @refresh="refreshImage" />

    <div class="image-hero-section" v-loading="loading">
      <ErrorState v-if="error" :error="error" @retry="refreshImage" />

      <FileDisplay
        v-else-if="imageData"
        :image-data="imageData"
        @image-load="onImageLoad"
        @image-error="onImageError"
      />
    </div>

    <div v-if="imageData" class="details-section">
      <div class="scroll-indicator">
        <i class="fas fa-chevron-up" />
        <span>{{ $t('random.scrollIndicator') }}</span>
      </div>

      <div class="w-full px-6 py-8">
        <div class="grid grid-cols-1 gap-6 lg:grid-cols-3">
          <CreatorCard :image-data="imageData" @go-to-author="goToAuthorPage" @view-work="viewWork" />

          <AIAnalysisCard :image-data="imageData" />

          <FileDetailsCard :image-data="imageData" />
        </div>

        <ActionButtons :image-data="imageData" @refresh="refreshImage" />
      </div>
    </div>

    <div class="keyboard-hint">
      <div class="hint-content">
        <span class="hint-text">{{ $t('random.keyboardHint.title') }}</span>
        <div class="key-hints">
          <span class="key-hint"> <kbd>Space</kbd> {{ $t('random.keyboardHint.refresh') }} </span>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
  .random-page-container {
    position: relative;
    height: calc(100vh - 4rem);
    margin-top: 4rem;
    overflow-y: auto;
    overflow-x: hidden;
    scroll-behavior: smooth;
  }

  .image-hero-section {
    height: calc(100vh - 4rem);
    position: sticky;
    top: 0;
    left: 0;
    right: 0;
    overflow: hidden;
    z-index: 0;
  }

  .details-section {
    position: relative;
    z-index: 10;
    backdrop-filter: blur(24px);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.15);
    border-radius: var(--radius-sm) var(--radius-sm) 0 0;
    box-shadow:
      0 -8px 32px rgba(0, 0, 0, 0.12),
      0 -4px 16px rgba(0, 0, 0, 0.08);
    min-height: 25vh;
    max-width: 1600px;
    margin: 0 auto;
    padding: 1.5rem 1.5rem 1.5rem;
  }

  .scroll-indicator {
    position: absolute;
    top: -3rem;
    left: 50%;
    transform: translateX(-50%);
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 0.25rem;
    color: rgba(255, 255, 255, 0.5);
    font-size: 0.75rem;
    font-weight: 500;
    animation: bounce 2s ease-in-out infinite;
    pointer-events: none;
    z-index: 1;
  }

  .scroll-indicator i {
    font-size: 1.25rem;
    opacity: 0.8;
  }

  @keyframes bounce {
    0%,
    100% {
      transform: translateX(-50%) translateY(0);
      opacity: 0.5;
    }
    50% {
      transform: translateX(-50%) translateY(-10px);
      opacity: 0.8;
    }
  }

  .keyboard-hint {
    position: fixed;
    bottom: 2rem;
    left: 2rem;
    background: rgba(var(--color-background-900-rgb), 0.5);
    backdrop-filter: blur(20px);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.2);
    border-radius: var(--radius-sm);
    padding: 0.75rem 1rem;
    z-index: 100;
    box-shadow: 0 2px 8px var(--color-overlay-light);
    transition: all 0.3s ease;
  }

  .keyboard-hint:hover {
    background: rgba(var(--color-background-900-rgb), 0.6);
    border-color: rgba(var(--color-brand-500-rgb), 0.3);
  }

  .hint-content {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }

  .hint-text {
    color: rgba(var(--color-white-rgb), 0.8);
    font-size: 0.75rem;
    font-weight: 500;
  }

  .key-hints {
    display: flex;
    flex-direction: column;
    gap: 0.25rem;
  }

  .key-hint {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    font-size: 0.7rem;
    color: rgba(var(--color-content-rgb), 0.7);
  }

  kbd {
    display: inline-block;
    background: rgba(var(--color-brand-500-rgb), 0.1);
    color: var(--color-brand-500);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.3);
    border-radius: var(--radius-sm);
    padding: 0.125rem 0.25rem;
    font-size: 0.65rem;
    font-family: 'Courier New', monospace;
    min-width: 1.2rem;
    text-align: center;
    margin-right: 0.125rem;
  }

  kbd:last-child {
    margin-right: 0;
  }

  @media (max-width: 768px) {
    .random-page-container {
      height: calc(100vh - 60px);
      margin-top: 60px;
    }

    .image-hero-section {
      height: calc(100vh - 60px);
    }

    .details-section {
      min-height: 35vh;
      padding: 1.25rem 1rem 1.25rem;
      border-radius: var(--radius-sm) var(--radius-sm) 0 0;
    }

    .keyboard-hint {
      bottom: 1rem;
      left: 1rem;
      right: 1rem;
      padding: 0.5rem 0.75rem;
    }

    .hint-content {
      flex-direction: row;
      align-items: center;
      justify-content: space-between;
    }

    .key-hints {
      flex-direction: row;
      gap: 0.5rem;
    }

    .hint-text {
      font-size: 0.7rem;
    }

    .key-hint {
      font-size: 0.65rem;
    }

    kbd {
      font-size: 0.6rem;
      padding: 0.1rem 0.2rem;
      min-width: 1rem;
    }
  }
</style>
