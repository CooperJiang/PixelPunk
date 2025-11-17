<script setup lang="ts">
  import { nextTick, onMounted, onUnmounted, ref, watch, computed } from 'vue'
  import type { WatermarkConfig, WatermarkPreviewProps } from './types'
  import './styles/preview.css'
  import {
    calculateWatermarkPosition,
    drawPreviewAreaBorder,
    getAnchorPosition,
    drawAnchorPoint,
    drawAnchorLine,
    drawResizeHandles,
    getImageRect,
    drawGuideLines,
  } from './utils/canvasUtils'
  import { useCanvasState } from './composables/useCanvasState'
  import { useImageLoader } from './composables/useImageLoader'
  import { useWatermarkDrag } from './composables/useWatermarkDrag'
  import { usePreviewResize } from './composables/usePreviewResize'
  import { useCanvasZoomPan } from './composables/useCanvasZoomPan'
  import { useRulers } from './composables/useRulers'
  import { useFileDrop } from './composables/useFileDrop'
  import { PRESET_SIZES, CANVAS_CONFIG } from './composables/constants'
  import { getBrandColor, getSuccessColor } from './utils/canvas/themeColors'
  import { useTexts } from '@/composables/useTexts'

  defineOptions({
    name: 'WatermarkPreview',
  })

  const { $t } = useTexts()

  const props = defineProps<WatermarkPreviewProps>()

  const emit = defineEmits<{
    (e: 'config-change', config: WatermarkConfig): void
  }>()

  const canvasRef = ref<HTMLCanvasElement>()
  const animationFrameId = ref<number>()

  const watermarkRect = ref({ x: 0, y: 0, width: 0, height: 0 })

  const localConfig = ref<WatermarkConfig>({ ...props.config })

  /* 监听 props.config 变化，同步到本地配置 */
  watch(
    () => props.config,
    (newConfig) => {
      localConfig.value = { ...newConfig }
    },
    { deep: true, immediate: true }
  )

  const { canvasSize, canvasZoom, canvasPan, previewSize, setPreviewSize, fitToView, resetZoom } = useCanvasState()

  const { previewImageData, watermarkImageData, droppedImageUrl, loadPreviewImage, loadWatermarkImage } = useImageLoader(
    props.previewImage,
    localConfig,
    drawCanvas
  )

  const { isDraggingWatermark, isPointInWatermark, startDrag, updateDrag, endDrag } = useWatermarkDrag(
    watermarkRect,
    previewImageData,
    previewSize,
    canvasZoom,
    canvasPan,
    localConfig,
    emit
  )

  const { isResizingPreview, hoveredCorner, checkHover, startResize, updateResize, endResize } = usePreviewResize(
    previewSize,
    canvasZoom,
    canvasPan
  )

  const { isSpacePressed, isPanningCanvas, handleWheel, startPan, updatePan, endPan, handleKeyDown, handleKeyUp } =
    useCanvasZoomPan(canvasZoom, canvasPan)

  const { horizontalRulerMarks, verticalRulerMarks } = useRulers(canvasZoom, canvasPan, canvasSize)

  const { isDragOver, handleDragOver, handleDragLeave, handleDrop } = useFileDrop(droppedImageUrl, loadPreviewImage)

  const initCanvas = () => {
    if (!canvasRef.value) return

    const canvas = canvasRef.value
    const ctx = canvas.getContext('2d')
    if (!ctx) return

    const container = canvas.parentElement
    if (container) {
      const width = container.clientWidth
      const height = container.clientHeight

      canvas.width = width
      canvas.height = height
      canvas.style.width = `${width}px`
      canvas.style.height = `${height}px`

      canvasSize.value = { width, height }

      const margin = CANVAS_CONFIG.fitMargin
      const fitZoomX = (width - margin) / previewSize.value.width
      const fitZoomY = (height - margin) / previewSize.value.height
      const fitZoom = Math.min(fitZoomX, fitZoomY, 1.0)

      const centerPanX = (width - previewSize.value.width * fitZoom) / 2
      const centerPanY = (height - previewSize.value.height * fitZoom) / 2

      canvasZoom.value = fitZoom
      canvasPan.value = { x: centerPanX, y: centerPanY }
    }

    drawCanvas()
  }

  const drawInfiniteGrid = (ctx: CanvasRenderingContext2D, width: number, height: number) => {
    const zoom = canvasZoom.value
    const { x: panX, y: panY } = canvasPan.value

    const baseGridSize = CANVAS_CONFIG.baseGridSize
    const gridSize = baseGridSize * zoom
    const startX = (panX % gridSize) - gridSize
    const startY = (panY % gridSize) - gridSize

    ctx.strokeStyle = getBrandColor(0.15)
    ctx.lineWidth = 1

    for (let x = startX; x < width + gridSize; x += gridSize) {
      ctx.beginPath()
      ctx.moveTo(Math.round(x), 0)
      ctx.lineTo(Math.round(x), height)
      ctx.stroke()
    }

    for (let y = startY; y < height + gridSize; y += gridSize) {
      ctx.beginPath()
      ctx.moveTo(0, Math.round(y))
      ctx.lineTo(width, Math.round(y))
      ctx.stroke()
    }

    const originX = panX
    const originY = panY
    if (originX >= 0 && originX <= width) {
      ctx.strokeStyle = getBrandColor(0.4)
      ctx.lineWidth = 2
      ctx.beginPath()
      ctx.moveTo(Math.round(originX), 0)
      ctx.lineTo(Math.round(originX), height)
      ctx.stroke()
    }
    if (originY >= 0 && originY <= height) {
      ctx.strokeStyle = getBrandColor(0.4)
      ctx.lineWidth = 2
      ctx.beginPath()
      ctx.moveTo(0, Math.round(originY))
      ctx.lineTo(width, Math.round(originY))
      ctx.stroke()
    }
  }

  function drawCanvas() {
    const canvas = canvasRef.value
    if (!canvas) return

    const ctx = canvas.getContext('2d')
    if (!ctx) return

    const { width, height } = canvasSize.value

    ctx.clearRect(0, 0, width, height)
    drawInfiniteGrid(ctx, width, height)

    const zoom = canvasZoom.value
    const { x: panX, y: panY } = canvasPan.value

    const previewW = previewSize.value.width
    const previewH = previewSize.value.height
    const previewCanvasX = 0
    const previewCanvasY = 0
    const previewScreenX = previewCanvasX * zoom + panX
    const previewScreenY = previewCanvasY * zoom + panY
    const previewScreenW = previewW * zoom
    const previewScreenH = previewH * zoom

    ctx.save()
    if (previewImageData.value) {
      const imgRect = getImageRect(previewImageData.value, previewW, previewH)
      const imgScreenX = imgRect.x * zoom + panX
      const imgScreenY = imgRect.y * zoom + panY
      const imgScreenW = imgRect.width * zoom
      const imgScreenH = imgRect.height * zoom
      ctx.drawImage(previewImageData.value, imgScreenX, imgScreenY, imgScreenW, imgScreenH)
    }
    ctx.restore()

    let wmWidth = 100
    let wmHeight = 30
    if (localConfig.value.type === 'text') {
      const tempCanvas = document.createElement('canvas')
      const tempCtx = tempCanvas.getContext('2d')
      if (tempCtx) {
        tempCtx.font = `${localConfig.value.fontWeight} ${localConfig.value.fontSize}px ${localConfig.value.fontFamily}`
        const metrics = tempCtx.measureText(localConfig.value.text)
        wmWidth = metrics.width * localConfig.value.scale
        wmHeight = localConfig.value.fontSize * localConfig.value.scale
      }
    } else if (watermarkImageData.value) {
      wmWidth = watermarkImageData.value.width * localConfig.value.scale
      wmHeight = watermarkImageData.value.height * localConfig.value.scale
    }

    if (localConfig.value.rotation !== 0) {
      const rad = (Math.abs(localConfig.value.rotation) * Math.PI) / 180
      const cos = Math.abs(Math.cos(rad))
      const sin = Math.abs(Math.sin(rad))
      const rotatedW = wmWidth * cos + wmHeight * sin
      const rotatedH = wmWidth * sin + wmHeight * cos
      wmWidth = rotatedW
      wmHeight = rotatedH
    }

    const position = calculateWatermarkPosition(
      localConfig.value,
      watermarkImageData.value,
      previewW,
      previewH,
      previewImageData.value
    )

    const watermarkScreenX = position.x * zoom + panX
    const watermarkScreenY = position.y * zoom + panY
    const watermarkScreenW = wmWidth * zoom
    const watermarkScreenH = wmHeight * zoom

    ctx.save()
    ctx.globalAlpha = localConfig.value.opacity
    if (localConfig.value.shadow) {
      ctx.shadowColor = localConfig.value.shadowColor
      ctx.shadowBlur = localConfig.value.shadowBlur * zoom
      ctx.shadowOffsetX = localConfig.value.shadowOffsetX * zoom
      ctx.shadowOffsetY = localConfig.value.shadowOffsetY * zoom
    }

    const centerX = watermarkScreenX + watermarkScreenW / 2
    const centerY = watermarkScreenY + watermarkScreenH / 2
    ctx.translate(centerX, centerY)

    if (localConfig.value.rotation !== 0) {
      ctx.rotate((localConfig.value.rotation * Math.PI) / 180)
    }

    if (localConfig.value.type === 'text') {
      ctx.font = `${localConfig.value.fontWeight} ${localConfig.value.fontSize * localConfig.value.scale * zoom}px ${localConfig.value.fontFamily}`
      ctx.fillStyle = localConfig.value.fontColor
      ctx.textAlign = 'center'
      ctx.textBaseline = 'middle'
      ctx.fillText(localConfig.value.text, 0, 0)
    } else if (localConfig.value.type === 'image' && watermarkImageData.value) {
      ctx.drawImage(watermarkImageData.value, -watermarkScreenW / 2, -watermarkScreenH / 2, watermarkScreenW, watermarkScreenH)
    }

    ctx.restore()

    const imgRect = getImageRect(previewImageData.value, previewW, previewH)
    const imgScreenX = imgRect.x * zoom + panX
    const imgScreenY = imgRect.y * zoom + panY
    const imgScreenW = imgRect.width * zoom
    const imgScreenH = imgRect.height * zoom

    if (isDraggingWatermark.value) {
      const wmScreenRect = { x: watermarkScreenX, y: watermarkScreenY, width: watermarkScreenW, height: watermarkScreenH }
      const imgScreenRect = { x: imgScreenX, y: imgScreenY, width: imgScreenW, height: imgScreenH }
      drawGuideLines(ctx, wmScreenRect, imgScreenRect, localConfig.value.position, true)
    }

    const previewAreaLabel = $t('components.watermark.preview.previewArea', {
      width: previewW,
      height: previewH,
    })
    drawPreviewAreaBorder(
      ctx,
      previewScreenX,
      previewScreenY,
      previewScreenW,
      previewScreenH,
      previewW,
      previewH,
      previewAreaLabel
    )

    if (previewImageData.value) {
      ctx.save()
      ctx.strokeStyle = getSuccessColor(0.6)
      ctx.lineWidth = 2
      ctx.setLineDash([6, 4])
      ctx.strokeRect(imgScreenX, imgScreenY, imgScreenW, imgScreenH)
      ctx.setLineDash([])

      ctx.fillStyle = getSuccessColor(0.9)
      ctx.font = 'bold 12px -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif'
      ctx.textAlign = 'left'
      ctx.textBaseline = 'top'

      const labelText = $t('components.watermark.preview.imageArea', {
        width: Math.round(imgRect.width),
        height: Math.round(imgRect.height),
      })
      const labelPadding = 6
      const labelMetrics = ctx.measureText(labelText)
      const labelWidth = labelMetrics.width + labelPadding * 2
      const labelHeight = 20
      ctx.fillStyle = getSuccessColor(0.9)
      ctx.fillRect(imgScreenX + 10, imgScreenY + 10, labelWidth, labelHeight)

      ctx.fillStyle = '#ffffff'
      ctx.fillText(labelText, imgScreenX + 10 + labelPadding, imgScreenY + 10 + 4)

      ctx.restore()
    }

    const anchorPos = getAnchorPosition(localConfig.value.position, imgRect)
    const anchorScreenX = anchorPos.x * zoom + panX
    const anchorScreenY = anchorPos.y * zoom + panY
    drawAnchorPoint(ctx, anchorScreenX, anchorScreenY)

    const watermarkCenterScreenX = watermarkScreenX + watermarkScreenW / 2
    const watermarkCenterScreenY = watermarkScreenY + watermarkScreenH / 2
    drawAnchorLine(ctx, anchorScreenX, anchorScreenY, watermarkCenterScreenX, watermarkCenterScreenY)

    drawResizeHandles(ctx, previewScreenX, previewScreenY, previewScreenW, previewScreenH, hoveredCorner.value)

    watermarkRect.value = {
      x: position.x,
      y: position.y,
      width: wmWidth,
      height: wmHeight,
    }
  }

  const handleCanvasMouseDown = (e: MouseEvent) => {
    const canvas = canvasRef.value
    if (!canvas) return

    const rect = canvas.getBoundingClientRect()
    const screenX = e.clientX - rect.left
    const screenY = e.clientY - rect.top

    if (startPan(e.clientX, e.clientY)) {
      canvas.style.cursor = 'grabbing'
      e.preventDefault()
      return
    }

    const zoom = canvasZoom.value
    const { x: panX, y: panY } = canvasPan.value
    const canvasX = (screenX - panX) / zoom
    const canvasY = (screenY - panY) / zoom

    const clickedHandle = checkHover(screenX, screenY)

    if (clickedHandle) {
      startResize(clickedHandle, screenX, screenY)
      e.preventDefault()
      return
    }

    if (startDrag(canvasX, canvasY)) {
      canvas.style.cursor = 'grabbing'
      e.preventDefault()
    }
  }

  const handleCanvasMouseMove = (e: MouseEvent) => {
    const canvas = canvasRef.value
    if (!canvas) return

    const rect = canvas.getBoundingClientRect()
    const screenX = e.clientX - rect.left
    const screenY = e.clientY - rect.top

    if (animationFrameId.value !== undefined) {
      cancelAnimationFrame(animationFrameId.value)
    }

    animationFrameId.value = requestAnimationFrame(() => {
      if (isPanningCanvas.value) {
        updatePan(e.clientX, e.clientY)
        return
      }

      const zoom = canvasZoom.value
      const { x: panX, y: panY } = canvasPan.value
      const canvasX = (screenX - panX) / zoom
      const canvasY = (screenY - panY) / zoom

      if (isResizingPreview.value) {
        updateResize(screenX, screenY)
        return
      }

      if (isDraggingWatermark.value) {
        updateDrag(canvasX, canvasY)
      } else {
        const hovered = checkHover(screenX, screenY)

        if (isSpacePressed.value) {
          canvas.style.cursor = 'grab'
        } else if (hovered) {
          if (hovered === 'top-left' || hovered === 'bottom-right') {
            canvas.style.cursor = 'nwse-resize'
          } else if (hovered === 'top-right' || hovered === 'bottom-left') {
            canvas.style.cursor = 'nesw-resize'
          }
        } else {
          canvas.style.cursor = isPointInWatermark(canvasX, canvasY) ? 'grab' : 'default'
        }
      }
    })
  }

  const handleCanvasMouseUp = () => {
    if (isDraggingWatermark.value) {
      endDrag()
      const canvas = canvasRef.value
      if (canvas) {
        canvas.style.cursor = 'default'
      }
      drawCanvas()
    }

    if (isResizingPreview.value) {
      endResize()
      const canvas = canvasRef.value
      if (canvas) {
        canvas.style.cursor = 'default'
      }
    }

    if (isPanningCanvas.value) {
      endPan()
      const canvas = canvasRef.value
      if (canvas) {
        canvas.style.cursor = isSpacePressed.value ? 'grab' : 'default'
      }
    }
  }

  const handleCanvasMouseLeave = () => {
    handleCanvasMouseUp()
  }

  const positionMap = computed<Record<string, string>>(() => ({
    'top-left': $t('components.watermark.position.topLeft'),
    'top-center': $t('components.watermark.position.topCenter'),
    'top-right': $t('components.watermark.position.topRight'),
    'middle-left': $t('components.watermark.position.middleLeft'),
    'middle-center': $t('components.watermark.position.center'),
    'middle-right': $t('components.watermark.position.middleRight'),
    'bottom-left': $t('components.watermark.position.bottomLeft'),
    'bottom-center': $t('components.watermark.position.bottomCenter'),
    'bottom-right': $t('components.watermark.position.bottomRight'),
  }))

  const getPositionName = (position: string): string => {
    return positionMap.value[position] || position
  }

  const handleCanvasWheel = (e: WheelEvent) => {
    const canvas = canvasRef.value
    if (!canvas) return

    const rect = canvas.getBoundingClientRect()
    const mouseX = e.clientX - rect.left
    const mouseY = e.clientY - rect.top

    handleWheel(e, mouseX, mouseY)
  }

  const handleResize = () => {
    initCanvas()
  }

  const handleSetPreviewSize = (width: number, height: number) => {
    setPreviewSize(width, height)
    nextTick(() => {
      fitToView()
    })
  }

  const resetView = () => {
    fitToView()
  }

  onMounted(() => {
    initCanvas()
    loadPreviewImage()

    document.addEventListener('mouseup', handleCanvasMouseUp)
    document.addEventListener('keydown', handleKeyDown)
    document.addEventListener('keyup', handleKeyUp)
    window.addEventListener('resize', handleResize)
  })

  onUnmounted(() => {
    document.removeEventListener('mouseup', handleCanvasMouseUp)
    document.removeEventListener('keydown', handleKeyDown)
    document.removeEventListener('keyup', handleKeyUp)
    window.removeEventListener('resize', handleResize)

    if (animationFrameId.value !== undefined) {
      cancelAnimationFrame(animationFrameId.value)
    }
  })

  watch(
    localConfig,
    () => {
      nextTick(() => {
        drawCanvas()
      })
    },
    { deep: true }
  )

  watch(
    () => props.previewImage,
    () => {
      loadPreviewImage()
    }
  )

  watch(
    () => [localConfig.value.type, localConfig.value.fileBase64],
    () => {
      loadWatermarkImage()
    },
    { immediate: true }
  )

  watch(
    () => [canvasZoom.value, canvasPan.value.x, canvasPan.value.y],
    () => {
      nextTick(() => {
        drawCanvas()
      })
    }
  )

  watch(
    () => [previewSize.value.width, previewSize.value.height],
    () => {
      nextTick(() => {
        drawCanvas()
      })
    }
  )
</script>

<template>
  <div class="watermark-preview">
    <div class="preview-toolbar">
      <div class="toolbar-left">
        <div class="zoom-indicator">
          <i class="fas fa-search-plus mr-1.5 text-xs" />
          {{ Math.round(canvasZoom * 100) }}%
        </div>

        <div class="size-selector-group">
          <label
            v-for="preset in PRESET_SIZES"
            :key="preset.name"
            class="size-radio-label"
            :class="{ active: previewSize.width === preset.width && previewSize.height === preset.height }"
          >
            <input
              type="radio"
              :value="`${preset.width}x${preset.height}`"
              :checked="previewSize.width === preset.width && previewSize.height === preset.height"
              @change="handleSetPreviewSize(preset.width, preset.height)"
              class="size-radio-input"
            />
            <span class="size-radio-text">{{ preset.name }}</span>
          </label>
        </div>

        <div class="current-size-display">
          <i class="fas fa-ruler-combined mr-1 text-xs" />
          <span class="size-text">{{ previewSize.width }} × {{ previewSize.height }}</span>
        </div>
      </div>

      <div class="toolbar-right">
        <div v-if="isDraggingWatermark" class="dragging-status">
          <i class="fas fa-hand-rock mr-1.5 text-xs" />
          {{ $t('components.watermark.preview.dragging') }}
        </div>

        <div class="view-controls-group">
          <button class="view-control-btn" :title="$t('components.watermark.preview.fitWindow')" @click="fitToView">
            <i class="fas fa-expand mr-1 text-xs" />
            <span class="btn-text">{{ $t('components.watermark.preview.fit') }}</span>
          </button>
          <button class="view-control-btn" :title="$t('components.watermark.preview.actualSize')" @click="resetZoom">
            <i class="fas fa-search mr-1 text-xs" />
            <span class="btn-text">1:1</span>
          </button>
          <button class="view-control-btn" :title="$t('components.watermark.preview.resetView')" @click="resetView">
            <i class="fas fa-redo mr-1 text-xs" />
            <span class="btn-text">{{ $t('components.watermark.preview.reset') }}</span>
          </button>
        </div>
      </div>
    </div>

    <div
      class="preview-main"
      :class="{ 'drag-over': isDragOver }"
      @dragover="handleDragOver"
      @dragleave="handleDragLeave"
      @drop="handleDrop"
    >
      <div class="ruler-vertical">
        <div v-for="mark in verticalRulerMarks" :key="mark.value" class="ruler-mark" :style="{ top: `${mark.position}px` }">
          <span class="ruler-label">{{ mark.value }}</span>
          <div class="ruler-tick" :class="{ major: mark.value % 100 === 0 }" />
        </div>
      </div>

      <div class="canvas-wrapper">
        <canvas
          ref="canvasRef"
          class="preview-canvas"
          :class="{ 'drag-active': isDragOver, panning: isSpacePressed }"
          @mousedown="handleCanvasMouseDown"
          @mousemove="handleCanvasMouseMove"
          @mouseup="handleCanvasMouseUp"
          @mouseleave="handleCanvasMouseLeave"
          @wheel="handleCanvasWheel"
        />

        <div v-if="isDragOver" class="drag-overlay">
          <i class="fas fa-image text-4xl" />
          <p class="drag-text">{{ $t('components.watermark.preview.dropToPreview') }}</p>
        </div>

        <div v-if="previewImageData && !isDraggingWatermark && !isDragOver" class="canvas-hints-panel">
          <div class="hint-item-row hint-info">
            <i class="fas fa-info-circle text-xs" />
            <span>
              <span class="hint-color-green">{{ $t('components.watermark.preview.greenFrame') }}</span
              >={{ $t('components.watermark.preview.imageArea2') }}
              <span class="hint-color-purple">{{ $t('components.watermark.preview.purpleFrame') }}</span
              >={{ $t('components.watermark.preview.previewCanvas') }}
            </span>
          </div>
        </div>

        <div v-if="isDraggingWatermark" class="offset-info-panel">
          <div class="info-header">
            <i class="fas fa-info-circle mr-1.5 text-xs" />
            {{ $t('components.watermark.preview.positionInfo') }}
          </div>
          <div class="info-content">
            <div class="info-item">
              <span class="info-label">{{ $t('components.watermark.preview.anchor') }}:</span>
              <span class="info-value">{{ getPositionName(config.position) }}</span>
            </div>
            <div class="info-item">
              <span class="info-label">{{ $t('components.watermark.preview.horizontal') }}:</span>
              <span class="info-value">{{ config.offsetX ?? 20 }}px</span>
            </div>
            <div class="info-item">
              <span class="info-label">{{ $t('components.watermark.preview.vertical') }}:</span>
              <span class="info-value">{{ config.offsetY ?? 20 }}px</span>
            </div>
          </div>
        </div>
      </div>

      <div class="ruler-horizontal">
        <div v-for="mark in horizontalRulerMarks" :key="mark.value" class="ruler-mark" :style="{ left: `${mark.position}px` }">
          <span class="ruler-label">{{ mark.value }}</span>
          <div class="ruler-tick" :class="{ major: mark.value % 100 === 0 }" />
        </div>
      </div>
    </div>
  </div>
</template>
