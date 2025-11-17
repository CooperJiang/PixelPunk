<script setup lang="ts">
  import { computed, onMounted, onUnmounted, ref, watch } from 'vue'
  import { useFileNavigation } from './composables/useFileNavigation'
  import { useFileViewer } from './composables/useFileViewer'
  import { useKeyboardShortcuts } from './composables/useKeyboardShortcuts'
  import NavigationBar from './components/NavigationBar.vue'
  import FileCanvas from './components/FileCanvas.vue'
  import SimilarFilesDrawer from './components/SimilarFilesDrawer.vue'
  import FileDetailDrawer from './components/FileDetailDrawer.vue'
  import InfoPanel from './components/InfoPanel.vue'
  import ControlBar from './components/ControlBar.vue'
  import { downloadFileQuick } from '@/utils/file/downloader'
  import { adminSimilarImages, gallerySimilarImages, searchSimilarImages, userSimilarImages } from '@/api/search'
  import { useToast } from '@/components/Toast/useToast'
  import type { FileTag, FileViewerProps, ViewerEmits } from './types'
  import { logger } from '@/utils/system/logger'
  import { useTexts } from '@/composables/useTexts'

  /* Props定义 */
  interface Props extends FileViewerProps {
    modelValue: boolean
  }

  const props = withDefaults(defineProps<Props>(), {
    file: null,
    files: () => [],
    initialIndex: 0,
    searchScope: 'default',
    showSideNav: false,
    showKeyboardTips: false,
    controlsHideTimeout: 5000,
  })

  /* 文件列表管理 */
  const safeFiles = ref(props.files || [])

  /* 监听文件变化 */
  watch(
    () => props.files,
    (newFiles) => {
      safeFiles.value = newFiles || []
    },
    { immediate: true }
  )

  const emit = defineEmits<ViewerEmits>()

  const toast = useToast()
  const { $t } = useTexts()

  /* 组件引用 */
  const viewerContent = ref<HTMLDivElement>()
  const fileCanvasRef = ref<InstanceType<typeof FileCanvas>>()

  /* 使用composables */
  const {
    currentIndex,
    totalFiles,
    hasMultipleFiles,
    hasPreviousFile,
    hasNextFile,
    currentFile,
    fileUrl: _fileUrl,
    showPreviousFile,
    showNextFile,
    goToFile: _goToFile,
    handleKeyboardNavigation,
  } = useFileNavigation(safeFiles, props.initialIndex)

  const {
    isLoading,
    hasError,
    isFirstView,
    isAllContentVisible,
    isImmersiveMode,
    isFullscreen,
    isFillMode: _isFillMode,
    isHorizontalImage: _isHorizontalImage,
    showModeIndicator,
    showZoomIndicator,
    isControlsHidden: _isControlsHidden,
    loadingProgress,
    loadingText,
    shouldUseFillMode,
    toggleControls,
    toggleFullscreen,
    toggleFitMode,
    showModeIndicatorTemporary: _showModeIndicatorTemporary,
    showZoomIndicatorTemporary,
    resetControlsTimer,
    handleFileLoad: handleImageLoadInternal,
    handleFileError: handleImageErrorInternal,
    retryLoading,
    startRealImageLoading,
    startLoadingProgressSimulation,
    stopLoadingProgressSimulation: _stopLoadingProgressSimulation,
    handleFullscreenChange,
  } = useFileViewer(props.file)

  /* 可见性计算 */
  const visible = computed({
    get: () => props.modelValue,
    set: (value) => emit('update:modelValue', value),
  })

  /* 当前缩放比例 */
  const currentScale = ref(1)

  /* 相似文件相关状态 */
  const showSimilarDrawer = ref(false)
  const similarLoading = ref(false)
  const similarError = ref<Error | null>(null)
  const similarImages = ref<any[]>([])

  /* 文件详情相关状态 */
  const showDetailDrawer = ref(false)
  const lastSimilarQueryFileId = ref<string | null>(null) // 记录上次查询相似文件的源文件ID

  /* 临时预览模式，用于预览相似文件 */
  const tempPreviewImage = ref<any>(null)
  const isInTempPreview = ref(false)

  const displayImage = computed(() => {
    if (isInTempPreview.value && tempPreviewImage.value) {
      return tempPreviewImage.value
    }
    const current = currentFile.value
    const fallback = props.file
    return current || fallback || null
  })

  const hasAIInfo = computed(() => {
    const image = displayImage.value
    return image && image.ai_info && image.ai_info.description
  })

  const imageBrightness = ref(0.5) // 默认中等亮度
  const isLightBackground = computed(() => imageBrightness.value > 0.55)

  const detectImageBrightness = (imageUrl: string) => {
    const img = new Image()
    img.crossOrigin = 'anonymous'
    img.onload = () => {
      const canvas = document.createElement('canvas')
      const ctx = canvas.getContext('2d')

      if (!ctx) return

      canvas.width = img.width
      canvas.height = img.height
      ctx.drawImage(img, 0, 0)

      try {
        const imageData = ctx.getImageData(0, 0, canvas.width, canvas.height)
        const data = imageData.data

        let total = 0
        let count = 0

        for (let i = 0; i < data.length; i += 12) {
          const r = data[i]
          const g = data[i + 1]
          const b = data[i + 2]

          const brightness = (0.299 * r + 0.587 * g + 0.114 * b) / 255
          total += brightness
          count++
        }

        imageBrightness.value = total / count
      } catch {
        imageBrightness.value = 0.5
      }
    }

    img.onerror = () => {
      imageBrightness.value = 0.5
    }

    img.src = imageUrl
  }

  watch(
    () => {
      const image = displayImage.value
      return image?.full_url || image?.url || null
    },
    (newUrl) => {
      if (newUrl) {
        detectImageBrightness(newUrl)
      }
    },
    { immediate: true }
  )

  const loadImageWithProgress = async () => {
    const imageUrl = displayImage.value?.full_url || displayImage.value?.url
    if (!imageUrl) return

    isLoading.value = true
    hasError.value = false

    try {
      const blob = await startRealImageLoading(imageUrl)
      if (blob) {
        URL.createObjectURL(blob)
      } else {
      }
    } catch (error) {
      logger.error('Error loading image with progress:', error)
      startLoadingProgressSimulation()
    }
  }

  const handleClose = () => {
    visible.value = false
    showSimilarDrawer.value = false
    if (isInTempPreview.value) {
      exitTempPreview()
    }
    emit('close')
  }

  const handleBackdropClick = () => {
    handleClose()
  }

  const handleMouseMove = () => {
    resetControlsTimer()
  }

  const handlePreviousClick = (e?: Event) => {
    if (e) {
      e.stopPropagation()
    }
    if (isInTempPreview.value) {
      exitTempPreview()
    }
    showPreviousFile()
  }

  const handleNextClick = (e?: Event) => {
    if (e) {
      e.stopPropagation()
    }
    if (isInTempPreview.value) {
      exitTempPreview()
    }
    showNextFile()
  }

  const exitTempPreview = () => {
    isInTempPreview.value = false
    tempPreviewImage.value = null
    loadImage() // 重新加载原来的文件
  }

  const handleImageLoad = (imageRef: HTMLImageElement) => {
    handleImageLoadInternal(imageRef)
    emit('load', displayImage.value)
  }

  const handleImageError = () => {
    handleImageErrorInternal()
    emit('error', { file: displayImage.value, event: new Event('error') })
  }

  const handleDownload = () => {
    const image = displayImage.value
    if (image) {
      downloadFileQuick(image.id, image.display_name || 'file')
    }
  }

  const handleTagClick = (tag: FileTag) => {
    emit('tag-click', tag)
  }

  const zoomIn = () => {
    if (fileCanvasRef.value) {
      fileCanvasRef.value.zoomIn()
      currentScale.value = fileCanvasRef.value.getScale()
      showZoomIndicatorTemporary()
    }
  }

  const zoomOut = () => {
    if (fileCanvasRef.value) {
      fileCanvasRef.value.zoomOut()
      currentScale.value = fileCanvasRef.value.getScale()
      showZoomIndicatorTemporary()
    }
  }

  const rotateLeft = () => {
    if (fileCanvasRef.value) {
      fileCanvasRef.value.rotateLeft()
    }
  }

  const rotateRight = () => {
    if (fileCanvasRef.value) {
      fileCanvasRef.value.rotateRight()
    }
  }

  const resetTransform = () => {
    if (fileCanvasRef.value) {
      fileCanvasRef.value.resetTransform()
      currentScale.value = 1
    }
  }

  const handleScaleChange = (newScale: number) => {
    currentScale.value = newScale
    showZoomIndicatorTemporary()
  }

  const handleFindSimilar = async () => {
    if (isInTempPreview.value) {
      showSimilarDrawer.value = true
      return
    }

    const currentFileId = currentFile.value?.id || props.file?.id

    if (currentFileId && lastSimilarQueryFileId.value === String(currentFileId)) {
      showSimilarDrawer.value = true
      return
    }

    showSimilarDrawer.value = true
    await findSimilarImages()
  }

  const findSimilarImages = async () => {
    const image = isInTempPreview.value ? currentFile.value || props.file : displayImage.value
    if (!image) {
      similarError.value = new Error($t('components.fileViewer.errors.noFileInfo'))
      return
    }

    lastSimilarQueryFileId.value = String(image.id)

    similarLoading.value = true
    similarError.value = null
    similarImages.value = []

    try {
      let result
      const fileId = String(image.id)

      switch (props.searchScope) {
        case 'gallery':
          result = await gallerySimilarImages(fileId)
          break
        case 'user':
          result = await userSimilarImages(fileId)
          break
        case 'admin':
          result = await adminSimilarImages(fileId)
          break
        default:
          result = await searchSimilarImages(fileId)
          break
      }

      if (result && result.data && result.data.items && result.data.items.length > 0) {
        similarImages.value = result.data.items

        similarImages.value = similarImages.value.filter((item: any) => String(item.id) !== String(image.id))

        if (similarImages.value.length === 0) {
          similarError.value = new Error($t('components.fileViewer.errors.noSimilarFiles'))
        }
      } else {
        similarError.value = new Error($t('components.fileViewer.errors.noSimilarFiles'))
      }
    } catch (error: unknown) {
      const errorObj = error as any

      if (errorObj.response?.status === 404) {
        similarError.value = new Error($t('components.fileViewer.errors.noVectorData'))
      } else if (errorObj.code === 110) {
        similarError.value = new Error(errorObj.message || $t('components.fileViewer.errors.vectorSearchUnavailable'))
      } else if (errorObj.message) {
        similarError.value = new Error(errorObj.message)
      } else {
        similarError.value = new Error($t('components.fileViewer.errors.searchFailed'))
      }
    } finally {
      similarLoading.value = false
    }
  }

  const viewSimilarImage = (image: Event) => {
    tempPreviewImage.value = image
    isInTempPreview.value = true

    showSimilarDrawer.value = false

    loadImage()
  }

  const downloadSimilarImage = async (image: Event) => {
    await downloadFileQuick(image.id, image.display_name || image.original_name || 'file')
    toast.success($t('components.fileViewer.messages.downloadStarted'))
  }

  useKeyboardShortcuts(visible, {
    onEscape: handleClose,
    onSpace: () => {
      if (currentScale.value !== 1) {
        resetTransform()
      }
      toggleFitMode()
    },
    onArrowLeft: () => handleKeyboardNavigation('ArrowLeft'),
    onArrowRight: () => handleKeyboardNavigation('ArrowRight'),
    onPlus: zoomIn,
    onMinus: zoomOut,
    onR: rotateRight,
    onL: rotateLeft,
    onF: toggleFullscreen,
    onReset: resetTransform,
  })

  watch(
    () => props.file,
    (newFile, oldFile) => {
      if (newFile && (!oldFile || newFile.id !== oldFile.id)) {
        isLoading.value = true
        hasError.value = false
        startLoadingProgressSimulation()
        resetTransform()
        if (isInTempPreview.value) {
          exitTempPreview()
        }
      }
    }
  )

  watch(
    () => props.files,
    () => {
      if (currentIndex.value >= (props.files?.length || 0)) {
        currentIndex.value = Math.max(0, (props.files?.length || 0) - 1)
      }
    }
  )

  watch(
    () => props.initialIndex,
    (newIndex) => {
      if (newIndex && newIndex >= 0 && newIndex < (props.files?.length || 0)) {
        currentIndex.value = newIndex
      }
    }
  )

  watch(currentIndex, (newIndex) => {
    if (props.files && props.files[newIndex]) {
      emit('change', props.files[newIndex], newIndex)
      loadImageWithProgress()
    }
  })

  onMounted(() => {
    document.addEventListener('fullscreenchange', handleFullscreenChange)

    if (props.file || (props.files && props.files.length > 0 && props.files[currentIndex.value])) {
      loadImageWithProgress()
    }
  })

  onUnmounted(() => {
    document.removeEventListener('fullscreenchange', handleFullscreenChange)
  })
</script>

<template>
  <Teleport to="body">
    <transition name="cyber-viewer-fade">
      <div v-if="visible" class="image-viewer-wrapper">
        <div class="cyber-image-viewer" :class="{ 'is-immersive': isImmersiveMode }">
          <div class="viewer-backdrop" @click="handleBackdropClick" />

          <div class="viewer-container">
            <NavigationBar
              :is-visible="isAllContentVisible"
              :current-file="displayImage"
              :is-fullscreen="isFullscreen"
              :has-ai-info="hasAIInfo"
              :is-light-background="isLightBackground"
              @toggle-fullscreen="toggleFullscreen"
              @download="handleDownload"
              @close="handleClose"
              @tag-click="handleTagClick"
              @find-similar="handleFindSimilar"
            />

            <div ref="viewerContent" class="viewer-content" @mousemove="handleMouseMove" @mouseleave="resetControlsTimer">
              <CyberFileLoading
                v-if="isLoading"
                :progress="loadingProgress"
                :loading-text="String(loadingText)"
                :show-data-stream="true"
              />

              <FileCanvas
                ref="fileCanvasRef"
                :current-file="displayImage"
                :file-url="displayImage?.full_url || ''"
                :is-loading="isLoading"
                :has-error="hasError"
                :show-side-nav="hasMultipleFiles"
                :has-multiple-files="hasMultipleFiles"
                :has-previous-file="hasPreviousFile || isInTempPreview"
                :has-next-file="hasNextFile || isInTempPreview"
                :should-use-fill-mode="shouldUseFillMode"
                :show-mode-indicator="showModeIndicator"
                :show-zoom-indicator="showZoomIndicator"
                :is-light-background="isLightBackground"
                @canvas-click="toggleControls"
                @prev-image="handlePreviousClick"
                @next-image="handleNextClick"
                @reset-transform="resetTransform"
                @retry-loading="retryLoading"
                @image-load="handleImageLoad"
                @image-error="handleImageError"
                @scale-change="handleScaleChange"
              />

              <InfoPanel
                :current-file="displayImage"
                :is-visible="isFirstView || isAllContentVisible"
                :show-keyboard-tips="showKeyboardTips"
                :should-use-fill-mode="shouldUseFillMode"
                :is-light-background="isLightBackground"
              />
            </div>

            <ControlBar
              :is-visible="isAllContentVisible"
              :scale="currentScale"
              :should-use-fill-mode="shouldUseFillMode"
              :has-multiple-files="hasMultipleFiles || isInTempPreview"
              :has-previous-file="hasPreviousFile || isInTempPreview"
              :has-next-file="hasNextFile || isInTempPreview"
              :current-index="isInTempPreview ? 0 : currentIndex"
              :total-files="isInTempPreview ? 1 : totalFiles"
              :is-light-background="isLightBackground"
              @zoom-in="zoomIn"
              @zoom-out="zoomOut"
              @rotate-left="rotateLeft"
              @rotate-right="rotateRight"
              @reset-transform="resetTransform"
              @toggle-fit-mode="toggleFitMode"
              @prev-image="handlePreviousClick"
              @next-image="handleNextClick"
            />

            <div class="fixed-actions-wrapper">
              <cyberTooltip
                :content="$t('components.fileViewer.actions.fileDetail')"
                placement="top"
                :offset="[-8, 20]"
                :show-delay="0"
                theme="dark"
              >
                <div class="fixed-action-btn detail-btn" @click="showDetailDrawer = true">
                  <i class="fas fa-info-circle" />
                </div>
              </cyberTooltip>

              <cyberTooltip
                v-if="hasAIInfo"
                :content="$t('components.fileViewer.actions.similarFiles')"
                placement="top"
                :offset="[-8, 20]"
                :show-delay="0"
                theme="dark"
              >
                <div class="fixed-action-btn similar-btn" @click="handleFindSimilar">
                  <i class="fas fa-layer-group" />
                </div>
              </cyberTooltip>
            </div>
          </div>
        </div>
      </div>
    </transition>
  </Teleport>

  <FileDetailDrawer v-model="showDetailDrawer" :file="displayImage" />

  <SimilarFilesDrawer
    v-model="showSimilarDrawer"
    :loading="similarLoading"
    :error="similarError"
    :similar-files="similarImages"
    @retry="findSimilarImages"
    @select="viewSimilarImage"
    @download="downloadSimilarImage"
    @view="viewSimilarImage"
  />
</template>

<style scoped>
  .image-viewer-wrapper {
    position: fixed;
    top: 0;
    left: 0;
    width: 100vw;
    height: 100vh;
    z-index: 9000;
  }

  .cyber-image-viewer {
    position: relative;
    width: 100%;
    height: 100%;
    background: rgba(var(--color-background-900-rgb), 0.95);
    z-index: 1;
    display: flex;
    align-items: center;
    justify-content: center;
    backdrop-filter: blur(10px);

    --bg-overlay: rgba(var(--color-background-900-rgb), calc(0.95 - var(--viewer-is-light, 0) * 0.15));
    --text-shadow-intensity: calc(1 + var(--viewer-is-light, 0) * 2);
    --border-opacity: calc(0.3 + var(--viewer-is-light, 0) * 0.4);
    --shadow-intensity: calc(0.2 + var(--viewer-is-light, 0) * 0.3);

    background: var(--bg-overlay);
  }

  .cyber-image-viewer.is-immersive {
    background: rgba(var(--color-background-900-rgb), 1);
  }

  .viewer-backdrop {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background: transparent;
    cursor: pointer;
    z-index: 1;
  }

  .viewer-container {
    position: relative;
    width: 100%;
    height: 100%;
    display: flex;
    flex-direction: column;
    z-index: 2;
  }

  .viewer-content {
    flex: 1;
    position: relative;
    overflow: hidden;
  }

  .cyber-viewer-fade-enter-active,
  .cyber-viewer-fade-leave-active {
    transition: all 0.3s ease;
  }

  .cyber-viewer-fade-enter-from,
  .cyber-viewer-fade-leave-to {
    opacity: 0;
    transform: scale(0.9);
  }

  .fixed-actions-wrapper {
    position: absolute;
    bottom: 60px;
    right: 24px;
    z-index: 1000;
    display: flex;
    flex-direction: column;
    gap: 12px;
  }

  .fixed-action-btn {
    position: relative;
    width: 44px;
    height: 44px;
    display: flex;
    align-items: center;
    justify-content: center;
    background: rgba(var(--color-background-900-rgb), 0.9);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.3);
    border-radius: var(--radius-full);
    color: rgba(var(--color-white-rgb), 0.8);
    cursor: pointer;
    user-select: none;
    backdrop-filter: blur(12px);
    box-shadow:
      0 4px 12px rgba(var(--color-background-900-rgb), 0.3),
      0 2px 6px rgba(var(--color-brand-500-rgb), 0.2),
      inset 0 1px 0 rgba(255, 255, 255, 0.1);
    transition: all 0.3s ease;
    font-size: 16px;
  }

  .fixed-action-btn:hover {
    color: var(--color-brand-400);
    border-color: rgba(var(--color-brand-500-rgb), 0.6);
    background: rgba(var(--color-background-900-rgb), 0.98);
    transform: translateY(-2px);
    box-shadow:
      0 6px 20px rgba(var(--color-background-900-rgb), 0.4),
      0 3px 12px rgba(var(--color-brand-500-rgb), 0.3),
      inset 0 1px 0 rgba(255, 255, 255, 0.2);
  }

  .fixed-action-btn:active {
    transform: translateY(-1px);
    box-shadow:
      0 4px 12px rgba(var(--color-background-900-rgb), 0.3),
      0 2px 6px rgba(var(--color-brand-500-rgb), 0.4),
      inset 0 2px 4px rgba(var(--color-background-900-rgb), 0.2);
  }

  .fixed-action-btn.detail-btn::after {
    content: '';
    position: absolute;
    inset: 0;
    border-radius: var(--radius-full);
    background: rgba(var(--color-info-rgb), 0.15);
    animation: pulse 3s ease-in-out infinite;
    opacity: 0;
  }

  .fixed-action-btn.similar-btn::after {
    content: '';
    position: absolute;
    inset: 0;
    border-radius: var(--radius-full);
    background: rgba(var(--color-brand-500-rgb), 0.2);
    animation: pulse 3s ease-in-out infinite;
    opacity: 0;
  }

  @keyframes pulse {
    0%,
    100% {
      transform: scale(1);
      opacity: 0;
    }
    50% {
      transform: scale(1.2);
      opacity: 0.3;
    }
  }

  @media (max-width: 768px) {
    .cyber-image-viewer {
      background: rgba(var(--color-background-900-rgb), 0.98);
    }

    .fixed-actions-wrapper {
      bottom: 50px;
      right: 20px;
      gap: 10px;
    }

    .fixed-action-btn {
      width: 40px;
      height: 40px;
      font-size: 14px;
    }
  }

  :deep(.cyber-tooltip) {
    z-index: 99999 !important;
  }

  :global(.cyber-tooltip) {
    z-index: 99999 !important;
  }
</style>
