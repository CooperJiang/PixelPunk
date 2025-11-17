import { computed, ref } from 'vue'

export function useFileTransform(onScaleChange?: (scale: number) => void) {
  const scale = ref(1)
  const rotation = ref(0)
  const translateX = ref(0)
  const translateY = ref(0)
  const isDragging = ref(false)

  const imageStyle = computed(() => {
    const safeScale = isNaN(scale.value) ? 1 : scale.value
    const safeRotation = isNaN(rotation.value) ? 0 : rotation.value
    const safeTranslateX = isNaN(translateX.value) ? 0 : translateX.value
    const safeTranslateY = isNaN(translateY.value) ? 0 : translateY.value

    const baseTransform = `translate(${safeTranslateX}px, ${safeTranslateY}px) scale(${safeScale}) rotate(${safeRotation}deg)`

    return {
      transform: baseTransform,
      transformOrigin: 'center center',
      transition: isDragging.value ? 'none' : 'transform 0.3s ease',
    }
  })

  const containerStyle = computed(() => ({
    cursor: isDragging.value ? 'grabbing' : 'grab',
  }))

  const zoomIn = (amount = 0.25) => {
    const newScale = scale.value + amount
    scale.value = Math.min(newScale, 5) // 最大5倍
  }

  const zoomOut = (amount = 0.25) => {
    const newScale = scale.value - amount
    scale.value = Math.max(newScale, 0.1) // 最小0.1倍
  }

  const rotateLeft = () => {
    rotation.value -= 90
  }

  const rotateRight = () => {
    rotation.value += 90
  }

  const resetTransform = () => {
    scale.value = 1
    rotation.value = 0
    translateX.value = 0
    translateY.value = 0
  }

  let startX = 0
  let startY = 0
  let startTranslateX = 0
  let startTranslateY = 0

  const handleMouseDown = (e: MouseEvent) => {
    if (e.button !== 0) {
      return
    } // 只处理左键

    isDragging.value = true
    startX = e.clientX
    startY = e.clientY
    startTranslateX = translateX.value
    startTranslateY = translateY.value

    document.addEventListener('mousemove', handleMouseMove)
    document.addEventListener('mouseup', handleMouseUp)
  }

  const handleMouseMove = (e: MouseEvent) => {
    if (!isDragging.value) {
      return
    }

    const deltaX = e.clientX - startX
    const deltaY = e.clientY - startY

    translateX.value = startTranslateX + deltaX
    translateY.value = startTranslateY + deltaY
  }

  const handleMouseUp = () => {
    isDragging.value = false
    document.removeEventListener('mousemove', handleMouseMove)
    document.removeEventListener('mouseup', handleMouseUp)
  }

  let touchStartX = 0
  let touchStartY = 0
  let touchStartTranslateX = 0
  let touchStartTranslateY = 0

  const handleTouchStart = (e: TouchEvent) => {
    if (e.touches.length === 1) {
      const touch = e.touches[0]
      touchStartX = touch.clientX
      touchStartY = touch.clientY
      touchStartTranslateX = translateX.value
      touchStartTranslateY = translateY.value
    }
  }

  const handleTouchMove = (e: TouchEvent) => {
    if (e.touches.length === 1) {
      const touch = e.touches[0]
      const deltaX = touch.clientX - touchStartX
      const deltaY = touch.clientY - touchStartY

      translateX.value = touchStartTranslateX + deltaX
      translateY.value = touchStartTranslateY + deltaY
    }
  }

  const handleTouchEnd = () => {}

  const handleWheel = (e: WheelEvent) => {
    e.preventDefault()
    const delta = e.deltaY > 0 ? -0.1 : 0.1
    const newScale = scale.value + delta
    scale.value = Math.max(0.1, Math.min(5, newScale))

    if (onScaleChange) {
      onScaleChange(scale.value)
    }
  }

  return {
    scale,
    rotation,
    translateX,
    translateY,
    isDragging,

    imageStyle,
    containerStyle,

    zoomIn,
    zoomOut,
    rotateLeft,
    rotateRight,
    resetTransform,

    handleMouseDown,
    handleMouseMove,
    handleMouseUp,
    handleTouchStart,
    handleTouchMove,
    handleTouchEnd,
    handleWheel,
  }
}
