import { ref, type Ref } from 'vue'
import type { WatermarkConfig } from '../types'
import { getImageRect } from '../utils/canvasUtils'

export function useWatermarkDrag(
  watermarkRect: Ref<{ x: number; y: number; width: number; height: number }>,
  previewImageData: Ref<HTMLImageElement | undefined>,
  previewSize: Ref<{ width: number; height: number }>,
  canvasZoom: Ref<number>,
  canvasPan: Ref<{ x: number; y: number }>,
  config: Ref<WatermarkConfig>,
  emit: (event: 'config-change', config: WatermarkConfig) => void
) {
  const isDraggingWatermark = ref(false)
  const dragStartPos = ref({ x: 0, y: 0 })
  const dragStartWatermarkPos = ref({ x: 0, y: 0 })

  const isPointInWatermark = (x: number, y: number) => {
    const rect = watermarkRect.value
    return x >= rect.x && x <= rect.x + rect.width && y >= rect.y && y <= rect.y + rect.height
  }

  const startDrag = (canvasX: number, canvasY: number) => {
    if (isPointInWatermark(canvasX, canvasY)) {
      isDraggingWatermark.value = true
      dragStartPos.value = { x: canvasX, y: canvasY }
      dragStartWatermarkPos.value = { x: watermarkRect.value.x, y: watermarkRect.value.y }
      return true
    }
    return false
  }

  const updateDrag = (canvasX: number, canvasY: number) => {
    if (!isDraggingWatermark.value) return

    const deltaX = canvasX - dragStartPos.value.x
    const deltaY = canvasY - dragStartPos.value.y

    let newWatermarkX = dragStartWatermarkPos.value.x + deltaX
    let newWatermarkY = dragStartWatermarkPos.value.y + deltaY
    const wmWidth = watermarkRect.value.width
    const wmHeight = watermarkRect.value.height

    const previewW = previewSize.value.width
    const previewH = previewSize.value.height

    const imgRect = getImageRect(previewImageData.value, previewW, previewH)
    newWatermarkX = Math.max(imgRect.x, Math.min(newWatermarkX, imgRect.x + imgRect.width - wmWidth))
    newWatermarkY = Math.max(imgRect.y, Math.min(newWatermarkY, imgRect.y + imgRect.height - wmHeight))

    const currentPosition = config.value.position
    let offsetX = 0,
      offsetY = 0

    if (currentPosition.includes('left')) {
      offsetX = newWatermarkX - imgRect.x
    } else if (currentPosition.includes('right')) {
      offsetX = imgRect.x + imgRect.width - (newWatermarkX + wmWidth)
    } else {
      offsetX = 0
    }

    if (currentPosition.includes('top')) {
      offsetY = newWatermarkY - imgRect.y
    } else if (currentPosition.includes('bottom')) {
      offsetY = imgRect.y + imgRect.height - (newWatermarkY + wmHeight)
    } else if (currentPosition.includes('middle')) {
      offsetY = 0
    }

    emit('config-change', {
      ...config.value,
      offsetX: Math.round(offsetX),
      offsetY: Math.round(offsetY),
      offsetUnit: 'px',
    })
  }

  const endDrag = () => {
    isDraggingWatermark.value = false
  }

  return {
    isDraggingWatermark,
    isPointInWatermark,
    startDrag,
    updateDrag,
    endDrag,
  }
}
