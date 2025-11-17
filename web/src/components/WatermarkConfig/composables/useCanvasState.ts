import { ref } from 'vue'
import { CANVAS_CONFIG, DEFAULT_PREVIEW_SIZE } from './constants'

export function useCanvasState() {
  const canvasSize = ref({ width: 580, height: 326 })
  const canvasZoom = ref(1.0)
  const canvasPan = ref({ x: 0, y: 0 })
  const previewSize = ref({ ...DEFAULT_PREVIEW_SIZE })

  const setPreviewSize = (width: number, height: number) => {
    previewSize.value = { width, height }
  }

  const fitToView = () => {
    const { width, height } = canvasSize.value
    const margin = CANVAS_CONFIG.fitMargin

    const fitZoomX = (width - margin) / previewSize.value.width
    const fitZoomY = (height - margin) / previewSize.value.height
    const fitZoom = Math.min(fitZoomX, fitZoomY, 1.0)

    const centerPanX = (width - previewSize.value.width * fitZoom) / 2
    const centerPanY = (height - previewSize.value.height * fitZoom) / 2

    canvasZoom.value = fitZoom
    canvasPan.value = { x: centerPanX, y: centerPanY }
  }

  const resetZoom = () => {
    const { width, height } = canvasSize.value
    canvasZoom.value = 1.0

    const centerPanX = (width - previewSize.value.width) / 2
    const centerPanY = (height - previewSize.value.height) / 2

    canvasPan.value = { x: centerPanX, y: centerPanY }
  }

  return {
    canvasSize,
    canvasZoom,
    canvasPan,
    previewSize,
    setPreviewSize,
    fitToView,
    resetZoom,
  }
}
