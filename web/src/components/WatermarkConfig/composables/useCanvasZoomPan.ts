import { ref, type Ref } from 'vue'
import { CANVAS_CONFIG } from './constants'

export function useCanvasZoomPan(canvasZoom: Ref<number>, canvasPan: Ref<{ x: number; y: number }>) {
  const isSpacePressed = ref(false)
  const isPanningCanvas = ref(false)
  const panStartPos = ref({ x: 0, y: 0 })

  const handleWheel = (e: WheelEvent, mouseX: number, mouseY: number) => {
    e.preventDefault()

    const canvasXBefore = (mouseX - canvasPan.value.x) / canvasZoom.value
    const canvasYBefore = (mouseY - canvasPan.value.y) / canvasZoom.value

    const zoomDelta = e.deltaY > 0 ? CANVAS_CONFIG.zoomSpeed : CANVAS_CONFIG.inverseZoomSpeed
    const newZoom = Math.max(CANVAS_CONFIG.minZoom, Math.min(CANVAS_CONFIG.maxZoom, canvasZoom.value * zoomDelta))

    const canvasXAfter = (mouseX - canvasPan.value.x) / newZoom
    const canvasYAfter = (mouseY - canvasPan.value.y) / newZoom

    canvasPan.value = {
      x: canvasPan.value.x + (canvasXAfter - canvasXBefore) * newZoom,
      y: canvasPan.value.y + (canvasYAfter - canvasYBefore) * newZoom,
    }

    canvasZoom.value = newZoom
  }

  const startPan = (clientX: number, clientY: number) => {
    if (isSpacePressed.value) {
      isPanningCanvas.value = true
      panStartPos.value = { x: clientX, y: clientY }
      return true
    }
    return false
  }

  const updatePan = (clientX: number, clientY: number) => {
    if (!isPanningCanvas.value) return

    const deltaX = clientX - panStartPos.value.x
    const deltaY = clientY - panStartPos.value.y
    canvasPan.value = {
      x: canvasPan.value.x + deltaX,
      y: canvasPan.value.y + deltaY,
    }
    panStartPos.value = { x: clientX, y: clientY }
  }

  const endPan = () => {
    isPanningCanvas.value = false
  }

  const handleKeyDown = (e: KeyboardEvent) => {
    if (e.code === 'Space' && !isSpacePressed.value) {
      isSpacePressed.value = true
      e.preventDefault()
    }
  }

  const handleKeyUp = (e: KeyboardEvent) => {
    if (e.code === 'Space') {
      isSpacePressed.value = false
      isPanningCanvas.value = false
    }
  }

  return {
    isSpacePressed,
    isPanningCanvas,
    handleWheel,
    startPan,
    updatePan,
    endPan,
    handleKeyDown,
    handleKeyUp,
  }
}
