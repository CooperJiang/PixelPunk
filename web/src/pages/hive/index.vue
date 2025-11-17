<script setup lang="ts">
  import { nextTick, onMounted, ref, watch } from 'vue'
  import { useTexts } from '@/composables/useTexts'

  import HoneycombCanvas from './components/HoneycombCanvas.vue'
  import FilePreview from './components/FilePreview.vue'

  import { useEventListener } from './composables/useEventListener'
  import { useHiveData } from './composables/useHiveData'
  import { useHiveUI } from './composables/useHiveUI'
  import { useHiveInit } from './composables/useHiveInit'
  import { getHexagonStyle, getLoaderHexStyle } from './utils/hexagonHelpers'
  import {
    DEFAULT_HORIZONTAL_GAP,
    DEFAULT_VERTICAL_GAP,
    DEFAULT_STRENGTH_3D,
    DEFAULT_CLIP_GAP,
    BACKGROUND_HEX_COUNT,
    LOADER_HEX_COUNT,
  } from './constants'

  const { $t } = useTexts()

  defineOptions({
    name: 'HivePage',
  })

  const canvasRef = ref<InstanceType<typeof HoneycombCanvas>>()

  const horizontalGap = ref(DEFAULT_HORIZONTAL_GAP)
  const verticalGap = ref(DEFAULT_VERTICAL_GAP)
  const strength3D = ref(DEFAULT_STRENGTH_3D)
  const _clipGap = ref(DEFAULT_CLIP_GAP)

  const {
    images,
    isLoading: _isLoading,
    isInitialLoading,
    isLoadingFadingOut,
    error,
    loadingProgress,
    loadingStartTime,
    loadFiles,
    loadMoreImages,
    preloadVisibleFiles,
    setCanvasUpdateCallback,
  } = useHiveData()

  setCanvasUpdateCallback(() => {
    canvasRef.value?.updateVisibleContent()
  })

  const {
    isFullscreen,
    tipsHidden,
    showPreview,
    previewFile,
    updateViewportDimensions,
    toggleTips,
    handleFullscreenChange,
    handleFileClick,
    closePreview,
    handleKeydown,
  } = useHiveUI()

  const { initializeHoneycomb } = useHiveInit(
    loadFiles,
    preloadVisibleFiles,
    images,
    isInitialLoading,
    isLoadingFadingOut,
    loadingStartTime
  )

  watch(
    [horizontalGap, verticalGap],
    () => {
      nextTick(() => {
        canvasRef.value?.updateLayout()
      })
    },
    { flush: 'post' }
  )

  useEventListener('document', 'keydown', handleKeydown)
  useEventListener('document', 'fullscreenchange', handleFullscreenChange)
  useEventListener('window', 'resize', updateViewportDimensions)

  onMounted(() => {
    updateViewportDimensions()
    initializeHoneycomb()
  })
</script>

<template>
  <div class="hive-container" :class="{ 'is-fullscreen': isFullscreen }">
    <div v-if="isInitialLoading" class="loading-overlay" :class="{ 'fading-out': isLoadingFadingOut }">
      <div class="loading-content">
        <div class="honeycomb-loader">
          <div v-for="i in LOADER_HEX_COUNT" :key="i" class="hex-loader" :style="getLoaderHexStyle(i)" />
        </div>
        <div class="loading-text">{{ $t('hive.loading.text') }}</div>
        <div class="loading-progress-container">
          <div class="progress-bar">
            <div class="progress-fill" :style="{ width: `${loadingProgress}%` }" />
          </div>
          <div class="progress-text">{{ loadingProgress.toFixed(0) }}%</div>
        </div>
      </div>
    </div>

    <div v-if="error" class="error-message">
      {{ error }}
    </div>

    <div class="user-tips" :class="{ 'tips-hidden': tipsHidden }">
      <div class="tip-item">
        <i class="fas fa-images" />
        {{ $t('hive.tips.loaded', { count: images.length }) }}
      </div>
      <div class="tip-item">
        <i class="fas fa-expand" />
        {{ $t('hive.tips.fullscreen') }}
      </div>
      <div class="tip-item">
        <i class="fas fa-hand-rock" />
        {{ $t('hive.tips.drag') }}
      </div>
      <button class="tips-toggle" @click="toggleTips">
        <i class="fas" :class="tipsHidden ? 'fa-chevron-down' : 'fa-chevron-up'" />
      </button>
    </div>

    <HoneycombCanvas
      v-if="!isInitialLoading"
      ref="canvasRef"
      :images="images"
      :horizontal-gap="horizontalGap"
      :vertical-gap="verticalGap"
      :strength3-d="strength3D"
      @load-more="loadMoreImages"
      @image-click="handleFileClick"
    />

    <FilePreview :show="showPreview" :image-data="previewFile" @close="closePreview" />

    <div class="hive-background">
      <div v-for="i in BACKGROUND_HEX_COUNT" :key="i" class="hive-hexagon" :style="getHexagonStyle(i)" />
    </div>
  </div>
</template>

<style scoped>
  .hive-container {
    min-height: 100vh;
    background: linear-gradient(
      135deg,
      var(--color-background-darkest) 0%,
      var(--color-background-dark) 50%,
      var(--color-background-secondary) 100%
    );
    color: var(--color-content);
    position: relative;
    overflow: hidden;
  }

  .hive-container.is-fullscreen {
    background: linear-gradient(
      135deg,
      var(--color-black) 0%,
      var(--color-background-darkest) 50%,
      var(--color-background-dark) 100%
    );
  }

  .loading-overlay {
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background: linear-gradient(
      135deg,
      var(--color-background-darkest) 0%,
      var(--color-background-dark) 50%,
      var(--color-background-secondary) 100%
    );
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: var(--z-index-dropdown);
    opacity: 1;
    transition: opacity 0.5s ease-out;
  }

  .loading-overlay.fading-out {
    opacity: 0;
  }

  .loading-content {
    text-align: center;
  }

  .honeycomb-loader {
    position: relative;
    width: 120px;
    height: 120px;
    margin: 0 auto 30px;
  }

  .hex-loader {
    position: absolute;
    width: 20px;
    height: 20px;
    background: linear-gradient(135deg, var(--color-info-400) 0%, #0099cc 100%);
    clip-path: polygon(50% 0%, 100% 25%, 100% 75%, 50% 100%, 0% 75%, 0% 25%);
    animation: hexPulse 2s ease-in-out infinite;
    left: 50%;
    top: 50%;
    margin-left: -10px;
    margin-top: -10px;
    box-shadow: 0 0 15px rgba(0, 212, 255, 0.5);
  }

  @keyframes hexPulse {
    0%,
    100% {
      opacity: 0.4;
      transform: translate(var(--x, 0), var(--y, 0)) scale(0.8) rotate(0deg);
      box-shadow: 0 0 10px rgba(0, 212, 255, 0.3);
    }

    50% {
      opacity: 1;
      transform: translate(var(--x, 0), var(--y, 0)) scale(1.2) rotate(180deg);
      box-shadow: 0 0 25px rgba(0, 212, 255, 0.8);
    }
  }

  .loading-text {
    font-size: 1.5rem;
    color: var(--color-info-400);
    margin-bottom: 10px;
    font-weight: bold;
  }

  .loading-progress-container {
    width: 300px;
    margin: 20px auto 0;
    position: relative;
    display: flex;
    flex-direction: column;
    gap: 20px;
  }

  .progress-bar {
    width: 100%;
    height: 8px;
    background: rgba(255, 255, 255, 0.2);
    border-radius: var(--radius-sm);
    overflow: hidden;
  }

  .progress-fill {
    height: 100%;
    background: linear-gradient(90deg, var(--color-info-400) 0%, var(--color-error-500) 100%);
    width: 0%;
    transition: width 0.2s ease-out;
    border-radius: var(--radius-sm);
  }

  .progress-text {
    text-align: center;
    color: var(--color-error-500);
    font-size: 1.2rem;
    font-weight: bold;
    margin-top: 20px !important;
  }

  .user-tips {
    position: fixed;
    top: 20px;
    left: 20px;
    background: rgba(0, 0, 0, 0.8);
    backdrop-filter: blur(10px);
    border: 1px solid rgba(0, 212, 255, 0.3);
    border-radius: var(--radius-sm);
    padding: 15px;
    z-index: 100;
    transition: all 0.3s ease;
    max-width: 300px;
  }

  .user-tips.tips-hidden {
    transform: translateY(-80%);
  }

  .tip-item {
    display: flex;
    align-items: center;
    margin-bottom: 8px;
    font-size: 0.9rem;
    color: var(--color-content);
  }

  .tip-item:last-child {
    margin-bottom: 0;
  }

  .tip-item i {
    margin-right: 8px;
    color: var(--color-info-400);
    width: 16px;
  }

  .tip-item kbd {
    background: var(--color-error-500);
    color: var(--color-text-on-brand);
    padding: 2px 6px;
    border-radius: var(--radius-sm);
    font-size: 0.8rem;
    font-weight: bold;
    margin: 0 4px;
  }

  .highlight {
    color: var(--color-error-500);
    font-weight: bold;
  }

  .tips-toggle {
    position: absolute;
    bottom: -15px;
    left: 50%;
    transform: translateX(-50%);
    background: rgba(var(--color-info-rgb), 0.2);
    border: 1px solid var(--color-info-400);
    color: var(--color-info-400);
    width: 30px;
    height: 30px;
    border-radius: var(--radius-full);
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    transition: all 0.3s ease;
  }

  .tips-toggle:hover {
    background: rgba(var(--color-info-rgb), 0.4);
    box-shadow: 0 0 10px rgba(var(--color-info-rgb), 0.5);
  }

  .error-message {
    position: fixed;
    top: 20px;
    left: 50%;
    transform: translateX(-50%);
    background: rgba(var(--color-error-rgb), 0.8);
    color: var(--color-content-heading);
    padding: 15px 20px;
    border-radius: var(--radius-sm);
    z-index: 100;
    backdrop-filter: blur(10px);
  }

  .hive-background {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    pointer-events: none;
    z-index: 1;
  }

  .hive-hexagon {
    position: absolute;
    background: rgba(0, 212, 255, 0.1);
    clip-path: polygon(30% 0%, 70% 0%, 100% 50%, 70% 100%, 30% 100%, 0% 50%);
    animation: float 6s ease-in-out infinite;
  }

  @keyframes float {
    0%,
    100% {
      transform: translateY(0px) rotate(0deg);
    }

    50% {
      transform: translateY(-20px) rotate(180deg);
    }
  }

  @media (max-width: 768px) {
    .hive-container {
      font-size: 0.9rem;
    }

    .user-tips {
      left: 10px;
      right: 10px;
      max-width: none;
    }

    .loading-text {
      font-size: 1.2rem;
    }
  }
</style>
