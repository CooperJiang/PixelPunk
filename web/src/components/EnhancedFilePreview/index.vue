<script setup lang="ts">
  import { computed, nextTick, onMounted, onUnmounted, ref, watch } from 'vue'
  import type { EnhancedFilePreviewProps, EnhancedFilePreviewEmits } from './types'
  import { useTexts } from '@/composables/useTexts'
  import { createEventHandlers } from './utils/eventHandlers'
  import FileContainer from './components/FileContainer.vue'
  import ControlButtons from './components/ControlButtons.vue'
  import ControlsHint from './components/ControlsHint.vue'
  import Indicators from './components/Indicators.vue'

  defineOptions({
    name: 'CyberEnhancedFilePreview',
  })

  const props = withDefaults(defineProps<EnhancedFilePreviewProps>(), {
    fileWidth: 0,
    fileHeight: 0,
  })

  const emit = defineEmits<EnhancedFilePreviewEmits>()
  const { $t } = useTexts()

  const overlayRef = ref<HTMLElement>()
  const containerRef = ref<any>()
  const imageRef = computed(() => containerRef.value?.imageRef)

  const scale = ref(1)
  const translateX = ref(0)
  const translateY = ref(0)
  const showZoomIndicator = ref(false)
  const showModeIndicator = ref(false)
  const showControlsHint = ref(true)
  const isImageLoaded = ref(false)
  const isFillMode = ref(true)
  const isBrowserFullscreen = ref(false)
  const isHorizontalImage = ref(true)

  const isDragging = ref(false)
  const dragStartX = ref(0)
  const dragStartY = ref(0)
  const dragStartTranslateX = ref(0)
  const dragStartTranslateY = ref(0)

  const lastTouchDistance = ref(0)
  const initialTouchScale = ref(1)
  let zoomIndicatorTimer: number | null = null
  let modeIndicatorTimer: number | null = null
  let controlsHintTimer: number | null = null

  const shouldUseFillMode = computed(() => isFillMode.value)

  const imageStyle = computed(() => {
    const baseTransform = `translate(${translateX.value}px, ${translateY.value}px) scale(${scale.value})`

    if (shouldUseFillMode.value && scale.value === 1) {
      return {
        transform: baseTransform,
        transformOrigin: 'center center',
        transition: isDragging.value ? 'none' : 'transform 0.2s ease-out',
        width: '100vw',
        height: '100vh',
        objectFit: 'cover' as const,
      }
    }

    if (!shouldUseFillMode.value && !isHorizontalImage.value && scale.value === 1) {
      return {
        transform: baseTransform,
        transformOrigin: 'center center',
        transition: isDragging.value ? 'none' : 'transform 0.2s ease-out',
        width: 'auto',
        height: '100vh',
        objectFit: 'contain' as const,
        maxWidth: '100vw',
        maxHeight: '100vh',
      }
    }

    return {
      transform: baseTransform,
      transformOrigin: 'center center',
      transition: isDragging.value ? 'none' : 'transform 0.2s ease-out',
      width: 'auto',
      height: 'auto',
      objectFit: 'contain' as const,
      maxWidth: 'none',
      maxHeight: 'none',
    }
  })

  const showZoomIndicatorTemporary = () => {
    showZoomIndicator.value = true
    if (zoomIndicatorTimer) {
      clearTimeout(zoomIndicatorTimer)
    }
    zoomIndicatorTimer = window.setTimeout(() => {
      showZoomIndicator.value = false
    }, 1000)
  }

  const showModeIndicatorTemporary = () => {
    showModeIndicator.value = true
    if (modeIndicatorTimer) {
      clearTimeout(modeIndicatorTimer)
    }
    modeIndicatorTimer = window.setTimeout(() => {
      showModeIndicator.value = false
    }, 1500)
  }

  const startControlsHintTimer = () => {
    showControlsHint.value = true
    if (controlsHintTimer) {
      clearTimeout(controlsHintTimer)
    }
    controlsHintTimer = window.setTimeout(() => {
      showControlsHint.value = false
    }, 5000)
  }

  const resetTransform = () => {
    const img = imageRef.value
    const container = containerRef.value?.$el || containerRef.value

    if (!img || !container) {
      return
    }

    const containerWidth = container.clientWidth
    const containerHeight = container.clientHeight
    const imgWidth = img.naturalWidth
    const imgHeight = img.naturalHeight

    if (imgWidth === 0 || imgHeight === 0) {
      return
    }

    if (shouldUseFillMode.value || !isHorizontalImage.value) {
      scale.value = 1
    } else {
      const scaleX = containerWidth / imgWidth
      const scaleY = containerHeight / imgHeight
      scale.value = Math.min(scaleX, scaleY)
    }

    translateX.value = 0
    translateY.value = 0
  }

  const toggleDisplayMode = () => {
    isFillMode.value = !isFillMode.value
    showModeIndicatorTemporary()
    resetTransform()
  }

  const enterBrowserFullscreen = async () => {
    try {
      if (overlayRef.value && !isBrowserFullscreen.value && document.fullscreenElement === null) {
        await overlayRef.value.requestFullscreen()
        isBrowserFullscreen.value = true
      }
    } catch {}
  }

  const exitBrowserFullscreen = async () => {
    try {
      if (isBrowserFullscreen.value && document.fullscreenElement) {
        await document.exitFullscreen()
        isBrowserFullscreen.value = false
      }
    } catch (_err) {}
  }

  const onImageLoad = () => {
    isImageLoaded.value = true

    const img = imageRef.value
    if (img) {
      isHorizontalImage.value = img.naturalWidth >= img.naturalHeight
    }

    resetTransform()
    startControlsHintTimer()
  }

  const onImageError = () => {}

  const {
    handleWheel,
    handleMouseDown,
    handleMouseMove,
    handleMouseUp,
    handleTouchStart,
    handleTouchMove,
    handleTouchEnd,
    handleImageClick,
  } = createEventHandlers({
    scale,
    translateX,
    translateY,
    isDragging,
    dragStartX,
    dragStartY,
    dragStartTranslateX,
    dragStartTranslateY,
    lastTouchDistance,
    initialTouchScale,
    showZoomIndicatorTemporary,
    resetTransform,
  })

  const handleKeydown = (e: KeyboardEvent) => {
    switch (e.key) {
      case 'Escape':
        if (isBrowserFullscreen.value) {
          exitBrowserFullscreen()
        } else {
          handleClose()
        }
        break
      case 'f':
      case 'F':
        if (!isBrowserFullscreen.value) {
          enterBrowserFullscreen()
        }
        break
      case ' ':
        e.preventDefault()
        toggleDisplayMode()
        break
    }
  }

  const handleOverlayClick = (e: MouseEvent) => {
    if (e.target === overlayRef.value) {
      handleClose()
    }
  }

  const handleClose = () => {
    emit('update:visible', false)
    emit('close')
  }

  const handleFullscreenChange = () => {
    isBrowserFullscreen.value = Boolean(document.fullscreenElement)
  }

  watch(
    () => props.visible,
    (newVal: boolean) => {
      if (newVal) {
        document.body.style.overflow = 'hidden'
        nextTick(() => {
          overlayRef.value?.focus()
        })
      } else {
        document.body.style.overflow = ''
        if (isBrowserFullscreen.value) {
          exitBrowserFullscreen()
        }
        isDragging.value = false
        document.body.style.cursor = ''
        if (zoomIndicatorTimer) {
          clearTimeout(zoomIndicatorTimer)
          zoomIndicatorTimer = null
        }
        if (modeIndicatorTimer) {
          clearTimeout(modeIndicatorTimer)
          modeIndicatorTimer = null
        }
        if (controlsHintTimer) {
          clearTimeout(controlsHintTimer)
          controlsHintTimer = null
        }
      }
    }
  )

  watch(
    () => props.fileUrl,
    () => {
      isImageLoaded.value = false
      resetTransform()
    }
  )

  watch(isFillMode, () => {
    if (isImageLoaded.value) {
      resetTransform()
    }
  })

  onMounted(() => {
    document.addEventListener('fullscreenchange', handleFullscreenChange)
  })

  onUnmounted(() => {
    document.removeEventListener('fullscreenchange', handleFullscreenChange)
    document.body.style.overflow = ''
    document.body.style.cursor = ''
    if (zoomIndicatorTimer) {
      clearTimeout(zoomIndicatorTimer)
    }
    if (modeIndicatorTimer) {
      clearTimeout(modeIndicatorTimer)
    }
    if (controlsHintTimer) {
      clearTimeout(controlsHintTimer)
    }
  })
</script>

<template>
  <Teleport to="body">
    <Transition name="enhanced-fullscreen">
      <div
        v-if="visible"
        ref="overlayRef"
        class="enhanced-fullscreen-overlay"
        tabindex="0"
        @click="handleOverlayClick"
        @keydown="handleKeydown"
      >
        <FileContainer
          ref="containerRef"
          :file-url="props.fileUrl"
          :file-name="props.fileName || $t('components.enhancedFilePreview.defaultImageName')"
          :should-use-fill-mode="shouldUseFillMode"
          :file-style="imageStyle"
          @wheel="handleWheel"
          @mousedown="handleMouseDown"
          @mousemove="handleMouseMove"
          @mouseup="handleMouseUp"
          @mouseleave="handleMouseUp"
          @touchstart="handleTouchStart"
          @touchmove="handleTouchMove"
          @touchend="handleTouchEnd"
          @file-click="handleImageClick"
          @file-load="onImageLoad"
          @file-error="onImageError"
        />

        <ControlButtons
          :is-browser-fullscreen="isBrowserFullscreen"
          @close="handleClose"
          @exit-fullscreen="exitBrowserFullscreen"
        />

        <ControlsHint
          :show-controls-hint="showControlsHint"
          :is-browser-fullscreen="isBrowserFullscreen"
          :should-use-fill-mode="shouldUseFillMode"
        />

        <Indicators
          :show-zoom-indicator="showZoomIndicator"
          :show-mode-indicator="showModeIndicator"
          :scale="scale"
          :should-use-fill-mode="shouldUseFillMode"
        />
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
  .enhanced-fullscreen-overlay {
    @apply fixed inset-0 z-[9999] outline-none backdrop-blur-sm;
    background: rgba(var(--color-background-900-rgb), 0.95);
  }

  .enhanced-fullscreen-enter-active,
  .enhanced-fullscreen-leave-active {
    @apply transition-opacity duration-300;
  }

  .enhanced-fullscreen-enter-from,
  .enhanced-fullscreen-leave-to {
    @apply opacity-0;
  }
</style>
