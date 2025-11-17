import type { Ref } from 'vue'

export interface EventHandlerOptions {
  scale: Ref<number>
  translateX: Ref<number>
  translateY: Ref<number>
  isDragging: Ref<boolean>
  dragStartX: Ref<number>
  dragStartY: Ref<number>
  dragStartTranslateX: Ref<number>
  dragStartTranslateY: Ref<number>
  lastTouchDistance: Ref<number>
  initialTouchScale: Ref<number>
  showZoomIndicatorTemporary: () => void
  resetTransform: () => void
}

export const createEventHandlers = (options: EventHandlerOptions) => {
  const {
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
  } = options

  const handleWheel = (e: WheelEvent) => {
    e.preventDefault()
    const delta = e.deltaY > 0 ? -0.1 : 0.1
    const newScale = Math.max(0.1, Math.min(5, scale.value + delta))

    if (newScale !== scale.value) {
      scale.value = newScale
      showZoomIndicatorTemporary()
    }
  }

  const handleMouseDown = (e: MouseEvent) => {
    if (e.button !== 0) return

    isDragging.value = true
    dragStartX.value = e.clientX
    dragStartY.value = e.clientY
    dragStartTranslateX.value = translateX.value
    dragStartTranslateY.value = translateY.value

    document.body.style.cursor = 'grabbing'
    e.preventDefault()
  }

  const handleMouseMove = (e: MouseEvent) => {
    if (!isDragging.value) return

    const deltaX = e.clientX - dragStartX.value
    const deltaY = e.clientY - dragStartY.value

    translateX.value = dragStartTranslateX.value + deltaX
    translateY.value = dragStartTranslateY.value + deltaY
  }

  const handleMouseUp = () => {
    if (isDragging.value) {
      isDragging.value = false
      document.body.style.cursor = ''
    }
  }

  const handleTouchStart = (e: TouchEvent) => {
    if (e.touches.length === 1) {
      const touch = e.touches[0]
      dragStartX.value = touch.clientX
      dragStartY.value = touch.clientY
      dragStartTranslateX.value = translateX.value
      dragStartTranslateY.value = translateY.value
      isDragging.value = true
    } else if (e.touches.length === 2) {
      const touch1 = e.touches[0]
      const touch2 = e.touches[1]
      const distance = Math.sqrt(Math.pow(touch2.clientX - touch1.clientX, 2) + Math.pow(touch2.clientY - touch1.clientY, 2))
      lastTouchDistance.value = distance
      initialTouchScale.value = scale.value
      isDragging.value = false
    }
  }

  const handleTouchMove = (e: TouchEvent) => {
    e.preventDefault()

    if (e.touches.length === 1 && isDragging.value) {
      const touch = e.touches[0]
      const deltaX = touch.clientX - dragStartX.value
      const deltaY = touch.clientY - dragStartY.value

      translateX.value = dragStartTranslateX.value + deltaX
      translateY.value = dragStartTranslateY.value + deltaY
    } else if (e.touches.length === 2) {
      const touch1 = e.touches[0]
      const touch2 = e.touches[1]
      const distance = Math.sqrt(Math.pow(touch2.clientX - touch1.clientX, 2) + Math.pow(touch2.clientY - touch1.clientY, 2))

      if (lastTouchDistance.value > 0) {
        const scaleChange = distance / lastTouchDistance.value
        const newScale = Math.max(0.1, Math.min(5, initialTouchScale.value * scaleChange))

        if (newScale !== scale.value) {
          scale.value = newScale
          showZoomIndicatorTemporary()
        }
      }
    }
  }

  const handleTouchEnd = () => {
    isDragging.value = false
    lastTouchDistance.value = 0
  }

  const handleImageClick = () => {
    resetTransform()
  }

  return {
    handleWheel,
    handleMouseDown,
    handleMouseMove,
    handleMouseUp,
    handleTouchStart,
    handleTouchMove,
    handleTouchEnd,
    handleImageClick,
  }
}
