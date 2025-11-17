<script setup lang="ts">
  import { computed, ref } from 'vue'
  import { useFileTransform } from '../composables/useFileTransform'
  import type { FileInfo } from '../types'
  import { useTexts } from '@/composables/useTexts'

  interface Props {
    currentFile: FileInfo | null
    fileUrl: string
    isLoading: boolean
    hasError: boolean
    showSideNav: boolean
    hasMultipleFiles: boolean
    hasPreviousFile: boolean
    hasNextFile: boolean
    shouldUseFillMode: boolean
    showModeIndicator: boolean
    showZoomIndicator: boolean
    isLightBackground?: boolean
  }

  const props = defineProps<Props>()

  const { $t } = useTexts()

  const emit = defineEmits<{
    'canvas-click': []
    'prev-image': []
    'next-image': []
    'reset-transform': []
    'retry-loading': []
    'image-load': [imageRef: HTMLImageElement]
    'image-error': []
    'scale-change': [scale: number]
  }>()

  const canvasRef = ref<HTMLDivElement>()
  const imageRef = ref<HTMLImageElement>()

  /* 使用文件变换composable */
  const {
    scale,
    rotation,
    translateX,
    translateY,
    imageStyle,
    containerStyle,
    handleMouseDown,
    handleMouseUp,
    handleTouchStart,
    handleTouchMove,
    handleTouchEnd,
    handleWheel,
  } = useFileTransform((newScale) => {
    emit('scale-change', newScale)
  })

  const computedImageStyle = computed(() => {
    const baseStyle = imageStyle.value

    if (
      props.shouldUseFillMode &&
      scale.value === 1 &&
      rotation.value === 0 &&
      translateX.value === 0 &&
      translateY.value === 0
    ) {
      return {
        ...baseStyle,
        width: '100vw',
        height: '100vh',
        objectFit: 'cover' as const,
        maxWidth: 'none',
        maxHeight: 'none',
      }
    }
    return {
      ...baseStyle,
      width: 'auto',
      height: 'auto',
      objectFit: 'contain' as const,
      maxWidth: '100%',
      maxHeight: '100%',
    }
  })

  /* 文件加载处理 */
  const handleImageLoad = () => {
    if (imageRef.value) {
      emit('image-load', imageRef.value)
    }
  }

  const handleImageError = () => {
    emit('image-error')
  }

  const handlePrevClick = (e?: Event) => {
    if (e) {
      e.stopPropagation()
    }
    emit('prev-image')
  }

  const handleNextClick = (e?: Event) => {
    if (e) {
      e.stopPropagation()
    }
    emit('next-image')
  }

  defineExpose({
    zoomIn: (amount = 0.25) => {
      const newScale = scale.value + amount
      scale.value = Math.min(newScale, 5)
    },
    zoomOut: (amount = 0.25) => {
      const newScale = scale.value - amount
      scale.value = Math.max(newScale, 0.1)
    },
    rotateLeft: () => {
      rotation.value -= 90
    },
    rotateRight: () => {
      rotation.value += 90
    },
    resetTransform: () => {
      scale.value = 1
      rotation.value = 0
      translateX.value = 0
      translateY.value = 0
    },
    getScale: () => scale.value,
  })
</script>

<template>
  <div ref="canvasRef" class="image-canvas" @wheel="handleWheel" @click="$emit('canvas-click')">
    <button
      v-if="hasMultipleFiles"
      class="side-nav-btn prev-side-btn"
      :class="{ 'is-light-bg': props.isLightBackground }"
      :disabled="!hasPreviousFile"
      :title="$t('components.fileCanvas.previous')"
      @click="handlePrevClick"
    >
      <div class="nav-btn-bg" />
      <i class="fas fa-chevron-left" />
    </button>

    <button
      v-if="hasMultipleFiles"
      class="side-nav-btn next-side-btn"
      :class="{ 'is-light-bg': props.isLightBackground }"
      :disabled="!hasNextFile"
      :title="$t('components.fileCanvas.next')"
      @click="handleNextClick"
    >
      <div class="nav-btn-bg" />
      <i class="fas fa-chevron-right" />
    </button>

    <div
      class="image-container"
      :style="containerStyle"
      @mousedown="handleMouseDown"
      @mouseup="handleMouseUp"
      @dblclick="$emit('reset-transform')"
      @touchstart="handleTouchStart"
      @touchmove="handleTouchMove"
      @touchend="handleTouchEnd"
    >
      <img
        v-show="!isLoading && !hasError"
        ref="imageRef"
        :src="fileUrl"
        :alt="currentFile?.display_name || $t('components.fileCanvas.file')"
        class="main-image"
        :style="computedImageStyle"
        draggable="false"
        @load="handleImageLoad"
        @error="handleImageError"
      />

      <div v-if="hasError" class="error-container">
        <i class="fas fa-exclamation-triangle" />
        <span>{{ $t('components.fileCanvas.loadFailed') }}</span>
        <button class="retry-btn" @click="$emit('retry-loading')">{{ $t('components.fileCanvas.retry') }}</button>
      </div>
    </div>

    <div v-if="showModeIndicator" class="mode-indicator">
      {{ $t(shouldUseFillMode ? 'components.fileCanvas.fillMode' : 'components.fileCanvas.fitMode') }}
    </div>

    <div v-if="showZoomIndicator && scale !== 1" class="zoom-indicator">{{ Math.round((isNaN(scale) ? 1 : scale) * 100) }}%</div>
  </div>
</template>

<style scoped>
  .image-canvas {
    position: relative;
    width: 100%;
    height: 100%;
    overflow: hidden;
  }

  .image-container {
    position: relative;
    width: 100%;
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
    overflow: hidden;
  }

  .main-image {
    display: block;
    user-select: none;
    pointer-events: none;
  }

  .side-nav-btn {
    position: absolute;
    top: 50%;
    transform: translateY(-50%);
    width: 45px;
    height: 45px;
    border-radius: var(--radius-lg);
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    transition: all var(--transition-normal) var(--ease-in-out);
    z-index: 10;
    background: transparent;
    color: var(--color-content-heading);
    font-size: var(--text-sm);
    overflow: hidden;
    backdrop-filter: var(--backdrop-blur-md);
    border: 1.5px solid rgba(var(--color-brand-500-rgb), 0.4);
    box-shadow:
      0 4px 20px rgba(var(--color-background-900-rgb), 0.35),
      0 1px 4px rgba(var(--color-brand-500-rgb), 0.2),
      inset 0 1px 0 rgba(255, 255, 255, 0.08);
  }

  .nav-btn-bg {
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: linear-gradient(
      135deg,
      rgba(var(--color-background-800-rgb), calc(0.7 + var(--viewer-is-light, 0) * 0.15)) 0%,
      rgba(var(--color-background-900-rgb), calc(0.8 + var(--viewer-is-light, 0) * 0.1)) 50%,
      rgba(var(--color-background-800-rgb), calc(0.7 + var(--viewer-is-light, 0) * 0.15)) 100%
    );
    transition: all var(--transition-normal) var(--ease-in-out);
    opacity: calc(0.7 + var(--viewer-is-light, 0) * 0.2);
    pointer-events: none;
  }

  .side-nav-btn:hover {
    transform: translateY(-50%) scale(1.05);
    border-color: rgba(var(--color-brand-500-rgb), 0.6);
    box-shadow:
      0 6px 25px rgba(var(--color-background-900-rgb), 0.4),
      0 2px 10px rgba(var(--color-brand-500-rgb), 0.3),
      0 0 20px rgba(var(--color-brand-500-rgb), 0.25),
      inset 0 1px 0 rgba(255, 255, 255, 0.12);
    color: var(--color-brand-400);
  }

  .side-nav-btn:hover .nav-btn-bg {
    opacity: 1;
    background: linear-gradient(
      135deg,
      rgba(var(--color-background-700-rgb), calc(0.8 + var(--viewer-is-light, 0) * 0.1)) 0%,
      rgba(var(--color-background-800-rgb), calc(0.9 + var(--viewer-is-light, 0) * 0.08)) 50%,
      rgba(var(--color-background-700-rgb), calc(0.8 + var(--viewer-is-light, 0) * 0.1)) 100%
    );
  }

  .side-nav-btn:active {
    transform: translateY(-50%) scale(0.95);
  }

  .side-nav-btn:disabled {
    opacity: 0.3;
    cursor: not-allowed;
    filter: grayscale(0.8);
  }

  .side-nav-btn:disabled:hover {
    transform: translateY(-50%);
    border-color: rgba(var(--color-border-subtle-rgb), calc(0.3 + var(--viewer-is-light, 0) * 0.4));
    box-shadow:
      0 4px 20px rgba(var(--color-background-900-rgb), calc(0.3 + var(--viewer-is-light, 0) * 0.3)),
      0 1px 4px rgba(var(--color-border-subtle-rgb), calc(0.1 + var(--viewer-is-light, 0) * 0.1)),
      inset 0 1px 0 rgba(255, 255, 255, 0.08);
  }

  .side-nav-btn i {
    position: relative;
    z-index: 1;
    text-shadow: 0 0 6px rgba(var(--color-content-rgb), 0.3);
  }

  .prev-side-btn {
    left: 30px;
  }

  .next-side-btn {
    right: 30px;
  }

  .error-container {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    color: var(--color-red-400);
    font-size: 16px;
    gap: 15px;
    text-shadow: 0 0 10px rgba(var(--color-error-rgb), 0.5);
  }

  .error-container i {
    font-size: 48px;
    color: var(--color-danger);
    text-shadow: 0 0 20px rgba(var(--color-error-rgb), 0.6);
    animation: errorPulse 2s ease-in-out infinite;
  }

  .retry-btn {
    background: linear-gradient(135deg, var(--color-error-500) 0%, var(--color-error-600) 100%);
    color: var(--color-white);
    border: 1px solid rgba(var(--color-error-rgb), 0.5);
    padding: var(--space-sm) var(--space-xl);
    border-radius: var(--radius-md);
    cursor: pointer;
    transition: all var(--transition-normal) var(--ease-in-out);
    font-weight: 600;
    text-transform: uppercase;
    letter-spacing: 1px;
    backdrop-filter: var(--backdrop-blur-md);
    box-shadow:
      0 4px 15px rgba(var(--color-error-rgb), 0.3),
      inset 0 1px 0 rgba(255, 255, 255, 0.2);
  }

  .retry-btn:hover {
    background: linear-gradient(135deg, var(--color-error-400) 0%, var(--color-error-500) 100%);
    box-shadow:
      0 6px 20px rgba(var(--color-error-rgb), 0.4),
      0 0 30px rgba(var(--color-error-rgb), 0.3);
    transform: translateY(-2px);
  }

  .mode-indicator,
  .zoom-indicator {
    position: absolute;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    background: linear-gradient(
      135deg,
      rgba(var(--color-background-800-rgb), calc(0.9 + var(--viewer-is-light, 0) * 0.08)) 0%,
      rgba(var(--color-background-900-rgb), calc(0.92 + var(--viewer-is-light, 0) * 0.06)) 50%,
      rgba(var(--color-background-800-rgb), calc(0.9 + var(--viewer-is-light, 0) * 0.08)) 100%
    );
    backdrop-filter: var(--backdrop-blur-lg);
    color: var(--color-content-default);
    padding: var(--space-sm) var(--space-lg);
    border-radius: var(--radius-full);
    font-size: var(--text-sm);
    font-weight: 600;
    z-index: 20;
    border: 1px solid rgba(var(--color-border-subtle-rgb), calc(0.4 + var(--viewer-is-light, 0) * 0.3));
    box-shadow:
      var(--shadow-lg),
      0 0 0 1px rgba(255, 255, 255, 0.05),
      inset 0 1px 0 rgba(255, 255, 255, 0.08);
    animation: cyberFadeInOut 2s ease-in-out;
    text-shadow: 0 0 8px rgba(var(--color-content-rgb), 0.4);
    letter-spacing: 0.5px;
  }

  .zoom-indicator {
    top: 60%;
    background: linear-gradient(
      135deg,
      rgba(var(--color-background-800-rgb), calc(0.9 + var(--viewer-is-light, 0) * 0.08)) 0%,
      rgba(var(--color-background-900-rgb), calc(0.92 + var(--viewer-is-light, 0) * 0.06)) 50%,
      rgba(var(--color-background-800-rgb), calc(0.9 + var(--viewer-is-light, 0) * 0.08)) 100%
    );
    box-shadow:
      var(--shadow-lg),
      0 0 0 1px rgba(255, 255, 255, 0.05),
      inset 0 1px 0 rgba(255, 255, 255, 0.08);
  }

  @keyframes cyberFadeInOut {
    0% {
      opacity: 0;
      transform: translate(-50%, -50%) scale(0.8) rotateX(90deg);
    }
    15% {
      opacity: 1;
      transform: translate(-50%, -50%) scale(1.05) rotateX(0deg);
    }
    85% {
      opacity: 1;
      transform: translate(-50%, -50%) scale(1) rotateX(0deg);
    }
    100% {
      opacity: 0;
      transform: translate(-50%, -50%) scale(0.8) rotateX(-90deg);
    }
  }

  @keyframes errorPulse {
    0%,
    100% {
      transform: scale(1);
      opacity: 1;
    }
    50% {
      transform: scale(1.1);
      opacity: 0.8;
    }
  }

  .side-nav-btn.is-light-bg {
    background: linear-gradient(
      135deg,
      rgba(var(--color-background-900-rgb), 0.75) 0%,
      rgba(var(--color-background-900-rgb), 0.7) 50%,
      rgba(var(--color-background-900-rgb), 0.78) 100%
    );
    border: 1.5px solid rgba(var(--color-brand-500-rgb), 0.4);
    color: rgba(255, 255, 255, 0.95);
    backdrop-filter: blur(25px) saturate(110%);
    box-shadow:
      0 12px 32px rgba(var(--color-background-900-rgb), 0.6),
      0 6px 16px rgba(var(--color-background-900-rgb), 0.4),
      0 2px 8px rgba(var(--color-brand-500-rgb), 0.3),
      0 0 0 1px rgba(255, 255, 255, 0.08),
      inset 0 1px 0 rgba(var(--color-brand-500-rgb), 0.15);
    text-shadow:
      0 0 12px rgba(var(--color-brand-500-rgb), 0.4),
      0 0 6px rgba(255, 255, 255, 0.3),
      2px 2px 4px rgba(var(--color-background-900-rgb), 0.8);
  }

  .side-nav-btn.is-light-bg:hover:not(:disabled) {
    background: linear-gradient(
      135deg,
      rgba(var(--color-background-900-rgb), 0.8) 0%,
      rgba(var(--color-background-900-rgb), 0.75) 50%,
      rgba(var(--color-background-900-rgb), 0.85) 100%
    );
    border-color: rgba(var(--color-brand-500-rgb), 0.6);
    color: rgba(255, 255, 255, 0.98);
    box-shadow:
      0 16px 40px rgba(var(--color-background-900-rgb), 0.7),
      0 8px 20px rgba(var(--color-background-900-rgb), 0.5),
      0 4px 12px rgba(var(--color-brand-500-rgb), 0.4),
      0 0 0 1px rgba(255, 255, 255, 0.12),
      inset 0 1px 0 rgba(var(--color-brand-500-rgb), 0.25);
  }

  .side-nav-btn.is-light-bg .nav-btn-bg {
    background: linear-gradient(
      135deg,
      rgba(var(--color-background-900-rgb), 0.75) 0%,
      rgba(var(--color-background-900-rgb), 0.7) 50%,
      rgba(var(--color-background-900-rgb), 0.78) 100%
    );
    backdrop-filter: blur(25px) saturate(110%);
  }

  @media (max-width: 768px) {
    .side-nav-btn {
      width: 40px;
      height: 40px;
      font-size: 12px;
    }

    .prev-side-btn {
      left: 20px;
    }

    .next-side-btn {
      right: 20px;
    }

    .mode-indicator,
    .zoom-indicator {
      padding: 10px 16px;
      font-size: 13px;
    }
  }

  @media (max-width: 480px) {
    .side-nav-btn {
      width: 35px;
      height: 35px;
      font-size: 11px;
    }

    .prev-side-btn {
      left: 15px;
    }

    .next-side-btn {
      right: 15px;
    }

    .mode-indicator,
    .zoom-indicator {
      padding: 8px 14px;
      font-size: 12px;
    }
  }
</style>
