import { ref, computed, onMounted, onUnmounted } from 'vue'
import type { Position } from '../types'
import { StorageUtil } from '@/utils/storage/storage'

export function useDrag(initialPosition?: Position, onPositionChange?: () => void) {
  const isDragging = ref(false)
  const hasDragged = ref(false)
  const position = ref<Position>(initialPosition || { x: window.innerWidth - 100, y: window.innerHeight - 100 })
  const dragOffset = ref<Position>({ x: 0, y: 0 })

  const floatStyle = computed(() => ({
    left: `${position.value.x}px`,
    top: `${position.value.y}px`,
    cursor: isDragging.value ? 'grabbing' : 'grab',
  }))

  const savePosition = () => {
    StorageUtil.set('upload-float-position', position.value)
    onPositionChange?.()
  }

  const loadPosition = () => {
    const saved = StorageUtil.get<{ x: number; y: number }>('upload-float-position')
    const elementSize = 80
    const margin = 20

    if (saved) {
      const windowWidth = window.innerWidth
      const windowHeight = window.innerHeight
      const maxX = windowWidth - elementSize - margin
      const maxY = windowHeight - elementSize - margin

      let { x, y } = saved

      if (x < margin || x > maxX || y < margin || y > maxY) {
        x = maxX
        y = maxY
      }

      position.value = {
        x: Math.max(margin, Math.min(x, maxX)),
        y: Math.max(margin, Math.min(y, maxY)),
      }
    } else {
      const maxX = window.innerWidth - elementSize - margin
      const maxY = window.innerHeight - elementSize - margin
      position.value = { x: maxX, y: maxY }
    }

    savePosition()
  }

  const handleMouseDown = (e: MouseEvent) => {
    isDragging.value = true
    hasDragged.value = false
    dragOffset.value = {
      x: e.clientX - position.value.x,
      y: e.clientY - position.value.y,
    }

    document.addEventListener('mousemove', handleMouseMove)
    document.addEventListener('mouseup', handleMouseUp)
    e.preventDefault()
  }

  const handleMouseMove = (e: MouseEvent) => {
    if (!isDragging.value) return

    hasDragged.value = true

    const newX = e.clientX - dragOffset.value.x
    const newY = e.clientY - dragOffset.value.y

    const maxX = window.innerWidth - 80
    const maxY = window.innerHeight - 80

    position.value = {
      x: Math.max(0, Math.min(newX, maxX)),
      y: Math.max(0, Math.min(newY, maxY)),
    }

    savePosition()
  }

  const handleMouseUp = (e: MouseEvent) => {
    isDragging.value = false

    document.removeEventListener('mousemove', handleMouseMove)
    document.removeEventListener('mouseup', handleMouseUp)
    e.preventDefault()

    return !hasDragged.value /* Return true if it was a click, not a drag */
  }

  const handleResize = () => {
    const elementSize = 80
    const margin = 20
    const maxX = window.innerWidth - elementSize - margin
    const maxY = window.innerHeight - elementSize - margin

    let needsUpdate = false

    if (position.value.x < margin) {
      position.value.x = margin
      needsUpdate = true
    } else if (position.value.x > maxX) {
      position.value.x = maxX
      needsUpdate = true
    }

    if (position.value.y < margin) {
      position.value.y = margin
      needsUpdate = true
    } else if (position.value.y > maxY) {
      position.value.y = maxY
      needsUpdate = true
    }

    if (needsUpdate) {
      savePosition()
    }
  }

  onMounted(() => {
    loadPosition()
    window.addEventListener('resize', handleResize)

    setTimeout(() => {
      handleResize()
    }, 100)
  })

  onUnmounted(() => {
    window.removeEventListener('resize', handleResize)
    document.removeEventListener('mousemove', handleMouseMove)
    document.removeEventListener('mouseup', handleMouseUp)
  })

  return {
    isDragging,
    hasDragged,
    position,
    floatStyle,
    handleMouseDown,
    handleMouseUp,
    loadPosition,
    savePosition,
  }
}
