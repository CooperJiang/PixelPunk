import { computed, type Ref } from 'vue'
import { RULER_CONFIG } from './constants'

export function useRulers(
  canvasZoom: Ref<number>,
  canvasPan: Ref<{ x: number; y: number }>,
  canvasSize: Ref<{ width: number; height: number }>
) {
  const getInterval = (zoom: number) => {
    if (zoom >= RULER_CONFIG.zoomThresholds.high) return RULER_CONFIG.intervals.zoom2x
    if (zoom >= RULER_CONFIG.zoomThresholds.medium) return RULER_CONFIG.intervals.zoom1x
    if (zoom >= RULER_CONFIG.zoomThresholds.low) return RULER_CONFIG.intervals.zoom05x
    return RULER_CONFIG.intervals.default
  }

  const horizontalRulerMarks = computed(() => {
    const marks: Array<{ value: number; position: number }> = []
    const zoom = canvasZoom.value
    const panX = canvasPan.value.x
    const viewportWidth = canvasSize.value.width

    const interval = getInterval(zoom)

    const canvasStart = -panX / zoom
    const canvasEnd = (viewportWidth - panX) / zoom
    const firstMark = Math.floor(canvasStart / interval) * interval

    for (let i = firstMark; i <= canvasEnd; i += interval) {
      const screenPos = i * zoom + panX
      if (screenPos >= 0 && screenPos <= viewportWidth) {
        marks.push({
          value: i,
          position: screenPos + 30,
        })
      }
    }
    return marks
  })

  const verticalRulerMarks = computed(() => {
    const marks: Array<{ value: number; position: number }> = []
    const zoom = canvasZoom.value
    const panY = canvasPan.value.y
    const viewportHeight = canvasSize.value.height

    const interval = getInterval(zoom)

    const canvasStart = -panY / zoom
    const canvasEnd = (viewportHeight - panY) / zoom
    const firstMark = Math.floor(canvasStart / interval) * interval

    for (let i = firstMark; i <= canvasEnd; i += interval) {
      const screenPos = i * zoom + panY
      if (screenPos >= 0 && screenPos <= viewportHeight) {
        marks.push({
          value: i,
          position: screenPos + 40,
        })
      }
    }
    return marks
  })

  return {
    horizontalRulerMarks,
    verticalRulerMarks,
  }
}
