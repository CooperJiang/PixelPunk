import { ref, type Ref } from 'vue'
import { CANVAS_CONFIG } from './constants'
import { getHoveredResizeHandle } from '../utils/canvasUtils'

export function usePreviewResize(
  previewSize: Ref<{ width: number; height: number }>,
  canvasZoom: Ref<number>,
  canvasPan: Ref<{ x: number; y: number }>
) {
  const isResizingPreview = ref(false)
  const resizeCorner = ref<string | null>(null)
  const resizeStartSize = ref({ width: 0, height: 0 })
  const dragStartPos = ref({ x: 0, y: 0 })
  const hoveredCorner = ref<string | null>(null)

  const checkHover = (screenX: number, screenY: number) => {
    const zoom = canvasZoom.value
    const { x: panX, y: panY } = canvasPan.value
    const previewW = previewSize.value.width
    const previewH = previewSize.value.height
    const previewScreenX = 0 * zoom + panX
    const previewScreenY = 0 * zoom + panY
    const previewScreenW = previewW * zoom
    const previewScreenH = previewH * zoom

    const hovered = getHoveredResizeHandle(screenX, screenY, previewScreenX, previewScreenY, previewScreenW, previewScreenH)
    hoveredCorner.value = hovered
    return hovered
  }

  const startResize = (corner: string, screenX: number, screenY: number) => {
    isResizingPreview.value = true
    resizeCorner.value = corner
    resizeStartSize.value = { width: previewSize.value.width, height: previewSize.value.height }
    dragStartPos.value = { x: screenX, y: screenY }
  }

  const updateResize = (screenX: number, screenY: number) => {
    if (!isResizingPreview.value || !resizeCorner.value) return

    const zoom = canvasZoom.value
    const deltaScreenX = screenX - dragStartPos.value.x
    const deltaScreenY = screenY - dragStartPos.value.y
    const deltaCanvasX = deltaScreenX / zoom
    const deltaCanvasY = deltaScreenY / zoom

    let newWidth = resizeStartSize.value.width
    let newHeight = resizeStartSize.value.height

    switch (resizeCorner.value) {
      case 'top-left':
        newWidth = resizeStartSize.value.width - deltaCanvasX
        newHeight = resizeStartSize.value.height - deltaCanvasY
        break
      case 'top-right':
        newWidth = resizeStartSize.value.width + deltaCanvasX
        newHeight = resizeStartSize.value.height - deltaCanvasY
        break
      case 'bottom-left':
        newWidth = resizeStartSize.value.width - deltaCanvasX
        newHeight = resizeStartSize.value.height + deltaCanvasY
        break
      case 'bottom-right':
        newWidth = resizeStartSize.value.width + deltaCanvasX
        newHeight = resizeStartSize.value.height + deltaCanvasY
        break
    }

    newWidth = Math.max(CANVAS_CONFIG.minPreviewSize, Math.min(CANVAS_CONFIG.maxPreviewSize, newWidth))
    newHeight = Math.max(CANVAS_CONFIG.minPreviewSize, Math.min(CANVAS_CONFIG.maxPreviewSize, newHeight))

    previewSize.value = {
      width: Math.round(newWidth),
      height: Math.round(newHeight),
    }
  }

  const endResize = () => {
    isResizingPreview.value = false
    resizeCorner.value = null
  }

  return {
    isResizingPreview,
    hoveredCorner,
    checkHover,
    startResize,
    updateResize,
    endResize,
  }
}
