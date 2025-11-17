<script setup lang="ts">
  import { computed, nextTick, onMounted, onUnmounted, ref, watch } from 'vue'
  import type { ImageInfo } from '@/types/global'
  import HoneycombCell from './HoneycombCell.vue'

  interface CellData {
    imageData: ImageInfo
    position: { x: number; y: number }
    transform: string
    currentTransform?: string
    row: number
    col: number
  }

  interface LayoutConfig {
    itemSize: number
    itemsPerOddRow: number
    itemsPerEvenRow: number
    rowHeight: number
    horizontalSpacing: number
  }

  const props = defineProps<{
    images: ImageInfo[]
    horizontalGap: number
    verticalGap: number
    strength3D: number
  }>()

  const emit = defineEmits<{
    loadMore: []
    imageClick: [image: ImageInfo]
  }>()

  const canvasContainer = ref<HTMLElement>()
  const honeycombCanvas = ref<HTMLElement>()

  const isDragging = ref(false)
  const lastMouseX = ref(0)
  const lastMouseY = ref(0)
  const canvasOffsetX = ref(0)
  const canvasOffsetY = ref(0)
  const dragDistance = ref(0)

  const visibleCells = ref<CellData[]>([])
  const cellsMap = new Map<string, number>()

  const canvasStyle = computed(() => ({
    transform: `translate(${canvasOffsetX.value}px, ${canvasOffsetY.value}px)`,
  }))

  const calculateLayout = (): LayoutConfig => {
    const minItemSize = 180
    const viewportWidth = window.innerWidth
    const itemsPerOddRow = Math.floor(viewportWidth / minItemSize)
    const itemsPerEvenRow = itemsPerOddRow + 1
    const actualItemSize = Math.floor(viewportWidth / itemsPerOddRow)

    const hexHeight = actualItemSize * 1.1547
    const verticalSpacing = actualItemSize * props.verticalGap
    const horizontalSpacing = actualItemSize * (1 + props.horizontalGap)

    return {
      itemSize: actualItemSize,
      itemsPerOddRow,
      itemsPerEvenRow,
      rowHeight: hexHeight * 0.75 + verticalSpacing,
      horizontalSpacing,
    }
  }

  const getCellKey = (row: number, col: number): string => `${row}_${col}`

  const usedImagesInViewport = new Map<string, Set<number>>()
  const imageUsageCount = new Map<number, number>()

  const getRegionKey = (row: number, col: number): string => {
    const regionSize = 6
    const regionRow = Math.floor(row / regionSize)
    const regionCol = Math.floor(col / regionSize)
    return `${regionRow}_${regionCol}`
  }

  const cleanupDistantRegions = (currentRow: number, currentCol: number) => {
    const currentRegionRow = Math.floor(currentRow / 6)
    const currentRegionCol = Math.floor(currentCol / 6)
    const cleanupDistance = 10

    for (const [regionKey] of usedImagesInViewport) {
      const [regionRow, regionCol] = regionKey.split('_').map(Number)
      const distance = Math.max(Math.abs(regionRow - currentRegionRow), Math.abs(regionCol - currentRegionCol))

      if (distance > cleanupDistance) {
        const usedSet = usedImagesInViewport.get(regionKey)
        if (usedSet) {
          for (const imageIndex of usedSet) {
            const currentCount = imageUsageCount.get(imageIndex) || 0
            if (currentCount <= 1) {
              imageUsageCount.delete(imageIndex)
            } else {
              imageUsageCount.set(imageIndex, currentCount - usedSet.size)
            }
          }
        }
        usedImagesInViewport.delete(regionKey)
      }
    }
  }

  const getImageFromPool = (row: number, col: number): ImageInfo | null => {
    if (props.images.length === 0) {
      return null
    }

    const regionKey = getRegionKey(row, col)

    if (Math.random() < 0.01) {
      cleanupDistantRegions(row, col)
    }

    if (!usedImagesInViewport.has(regionKey)) {
      usedImagesInViewport.set(regionKey, new Set())
    }
    const usedInRegion = usedImagesInViewport.get(regionKey) as Set<number>

    const adjacentUsed = new Set<number>()
    for (let dr = -1; dr <= 1; dr++) {
      for (let dc = -1; dc <= 1; dc++) {
        const adjRegionRow = Math.floor(row / 6) + dr
        const adjRegionCol = Math.floor(col / 6) + dc
        const adjRegionKey = `${adjRegionRow}_${adjRegionCol}`
        const adjUsed = usedImagesInViewport.get(adjRegionKey)
        if (adjUsed) {
          for (const imageIndex of adjUsed) {
            adjacentUsed.add(imageIndex)
          }
        }
      }
    }

    const candidates: number[] = []
    const fallbackCandidates: number[] = []

    for (let i = 0; i < props.images.length; i++) {
      if (!usedInRegion.has(i) && !adjacentUsed.has(i)) {
        candidates.push(i)
      } else if (!usedInRegion.has(i)) {
        fallbackCandidates.push(i)
      }
    }

    let selectedIndex: number

    if (candidates.length > 0) {
      candidates.sort((a, b) => {
        const usageA = imageUsageCount.get(a) || 0
        const usageB = imageUsageCount.get(b) || 0
        return usageA - usageB
      })

      const topCandidates = candidates.slice(0, Math.min(5, candidates.length))
      selectedIndex = topCandidates[Math.floor(Math.random() * topCandidates.length)]
    } else if (fallbackCandidates.length > 0) {
      fallbackCandidates.sort((a, b) => {
        const usageA = imageUsageCount.get(a) || 0
        const usageB = imageUsageCount.get(b) || 0
        return usageA - usageB
      })

      const topFallbacks = fallbackCandidates.slice(0, Math.min(3, fallbackCandidates.length))
      selectedIndex = topFallbacks[Math.floor(Math.random() * topFallbacks.length)]
    } else {
      const allIndices = Array.from({ length: props.images.length }, (_, i) => i)
      allIndices.sort((a, b) => {
        const usageA = imageUsageCount.get(a) || 0
        const usageB = imageUsageCount.get(b) || 0
        return usageA - usageB
      })

      selectedIndex = allIndices[0]
    }

    usedInRegion.add(selectedIndex)
    const currentUsage = imageUsageCount.get(selectedIndex) || 0
    imageUsageCount.set(selectedIndex, currentUsage + 1)

    return props.images[selectedIndex]
  }

  const calculate3DTransform = (x: number, y: number): string => {
    const centerX = window.innerWidth / 2
    const centerY = window.innerHeight / 2
    const distanceFromCenterX = (x - centerX + canvasOffsetX.value) / window.innerWidth
    const distanceFromCenterY = (y - centerY + canvasOffsetY.value) / window.innerHeight

    const rotateY = distanceFromCenterX * (25 * props.strength3D)
    const rotateX = distanceFromCenterY * (12 * props.strength3D)
    const translateZ = -Math.abs(distanceFromCenterX) * (150 * props.strength3D) - 50 * props.strength3D

    return `rotateY(${rotateY}deg) rotateX(${rotateX}deg) translateZ(${translateZ}px)`
  }

  const getVisibleRange = (layout: LayoutConfig) => {
    const viewportLeft = -canvasOffsetX.value
    const viewportTop = -canvasOffsetY.value
    const viewportRight = viewportLeft + window.innerWidth
    const viewportBottom = viewportTop + window.innerHeight

    const startRow = Math.floor(viewportTop / layout.rowHeight) - 2
    const endRow = Math.ceil(viewportBottom / layout.rowHeight) + 2
    const startCol = Math.floor(viewportLeft / layout.horizontalSpacing) - 2
    const endCol = Math.ceil(viewportRight / layout.horizontalSpacing) + 2

    return { startRow, endRow, startCol, endCol }
  }

  const updateVisibleContent = () => {
    if (props.images.length === 0) {
      return
    }

    const layout = calculateLayout()
    const range = getVisibleRange(layout)
    const cellsToKeep = new Set<string>()

    for (let row = range.startRow; row <= range.endRow; row++) {
      const isEvenRow = row % 2 === 0

      for (let col = range.startCol; col <= range.endCol; col++) {
        const key = getCellKey(row, col)
        cellsToKeep.add(key)

        if (!cellsMap.has(key)) {
          const imageData = getImageFromPool(row, col)
          if (imageData) {
            const x = col * layout.horizontalSpacing + (isEvenRow ? 0 : layout.horizontalSpacing * 0.5)
            const y = row * layout.rowHeight
            const baseTransform = calculate3DTransform(x, y)

            visibleCells.value.push({
              imageData,
              position: { x, y },
              transform: baseTransform,
              currentTransform: baseTransform,
              row,
              col,
            })
            cellsMap.set(key, visibleCells.value.length - 1)
          }
        } else {
          const index = cellsMap.get(key) as number
          const cell = visibleCells.value[index]
          const x = col * layout.horizontalSpacing + (isEvenRow ? 0 : layout.horizontalSpacing * 0.5)
          const y = row * layout.rowHeight
          const newBaseTransform = calculate3DTransform(x, y)

          const hasHoverTransform = cell.currentTransform && cell.currentTransform !== cell.transform

          const updatedCell = {
            ...cell,
            position: { x, y },
            transform: newBaseTransform,
          }

          if (!hasHoverTransform) {
            updatedCell.currentTransform = newBaseTransform
          }

          visibleCells.value[index] = updatedCell
        }
      }
    }

    for (let i = visibleCells.value.length - 1; i >= 0; i--) {
      const key = getCellKey(visibleCells.value[i].row, visibleCells.value[i].col)
      if (!cellsToKeep.has(key)) {
        visibleCells.value.splice(i, 1)
        cellsMap.delete(key)

        cellsMap.clear()
        visibleCells.value.forEach((cell, index) => {
          const cellKey = getCellKey(cell.row, cell.col)
          cellsMap.set(cellKey, index)
        })
      }
    }
  }

  const resetImageTracking = () => {
    usedImagesInViewport.clear()
    imageUsageCount.clear()
  }

  const updateLayout = () => {
    const layout = calculateLayout()
    document.documentElement.style.setProperty('--item-size', `${layout.itemSize}px`)

    visibleCells.value = []
    cellsMap.clear()

    resetImageTracking()

    updateVisibleContent()
  }

  const handleCellClick = (imageData: ImageInfo) => {
    if (dragDistance.value < 10) {
      emit('imageClick', imageData)
    }
  }

  const handleCellHover = (row: number, col: number) => {
    if (isDragging.value) {
      return
    }

    clearHoverEffects()

    const currentKey = getCellKey(row, col)
    const currentIndex = cellsMap.get(currentKey)
    if (currentIndex === undefined) {
      return
    }

    applyHoverScale(visibleCells.value[currentIndex], 1.15)

    const neighbors = getNeighborCells(row, col)
    neighbors.forEach((neighbor) => {
      applyHoverScale(neighbor, 0.85)
    })
  }

  const handleCellLeave = (e: MouseEvent) => {
    const relatedTarget = e.relatedTarget as HTMLElement
    if (relatedTarget && relatedTarget.closest('.honeycomb-item')) {
      return
    }

    clearHoverEffects()
  }

  const getNeighborCells = (row: number, col: number) => {
    const neighbors: CellData[] = []
    const isEvenRow = row % 2 === 0

    const neighborOffsets = !isEvenRow
      ? [
          { row: 0, col: -1 },
          { row: 0, col: 1 },
          { row: -1, col: 0 },
          { row: -1, col: 1 },
          { row: 1, col: 0 },
          { row: 1, col: 1 },
        ]
      : [
          { row: 0, col: -1 },
          { row: 0, col: 1 },
          { row: -1, col: -1 },
          { row: -1, col: 0 },
          { row: 1, col: -1 },
          { row: 1, col: 0 },
        ]

    neighborOffsets.forEach((offset) => {
      const neighborRow = row + offset.row
      const neighborCol = col + offset.col
      const key = getCellKey(neighborRow, neighborCol)
      const neighborIndex = cellsMap.get(key)
      if (neighborIndex !== undefined) {
        neighbors.push(visibleCells.value[neighborIndex])
      }
    })

    return neighbors
  }

  const clearHoverEffects = () => {
    for (const cell of visibleCells.value) {
      removeHoverScale(cell)
    }
  }

  const applyHoverScale = (cell: CellData, scale: number) => {
    const baseTransform = cell.transform
    const newTransform = `${baseTransform} scale(${scale})`

    const key = getCellKey(cell.row, cell.col)
    const index = cellsMap.get(key)
    if (index !== undefined) {
      visibleCells.value[index] = {
        ...cell,
        currentTransform: newTransform,
      }
    }

    nextTick(() => {})
  }

  const removeHoverScale = (cell: CellData) => {
    const key = getCellKey(cell.row, cell.col)
    const index = cellsMap.get(key)
    if (index !== undefined) {
      visibleCells.value[index] = {
        ...cell,
        currentTransform: cell.transform,
      }
    }

    nextTick(() => {})
  }

  const handleMouseDown = (e: MouseEvent) => {
    isDragging.value = true
    lastMouseX.value = e.clientX
    lastMouseY.value = e.clientY
    dragDistance.value = 0
    e.preventDefault()
  }

  const handleMouseMove = (e: MouseEvent) => {
    if (!isDragging.value) {
      return
    }

    const deltaX = e.clientX - lastMouseX.value
    const deltaY = e.clientY - lastMouseY.value

    canvasOffsetX.value += deltaX
    canvasOffsetY.value += deltaY

    dragDistance.value += Math.abs(deltaX) + Math.abs(deltaY)

    updateVisibleContent()

    checkLoadMoreNeeded()

    lastMouseX.value = e.clientX
    lastMouseY.value = e.clientY
  }

  const handleMouseUp = () => {
    isDragging.value = false
  }

  const handleTouchStart = (e: TouchEvent) => {
    e.preventDefault()
    isDragging.value = true
    const touch = e.touches[0]
    lastMouseX.value = touch.clientX
    lastMouseY.value = touch.clientY
    dragDistance.value = 0
  }

  const handleTouchMove = (e: TouchEvent) => {
    if (!isDragging.value) {
      return
    }
    e.preventDefault()

    const touch = e.touches[0]
    const deltaX = touch.clientX - lastMouseX.value
    const deltaY = touch.clientY - lastMouseY.value

    canvasOffsetX.value += deltaX
    canvasOffsetY.value += deltaY

    dragDistance.value += Math.abs(deltaX) + Math.abs(deltaY)

    updateVisibleContent()

    checkLoadMoreNeeded()

    lastMouseX.value = touch.clientX
    lastMouseY.value = touch.clientY
  }

  const handleTouchEnd = (e: TouchEvent) => {
    e.preventDefault()
    isDragging.value = false
  }

  const checkLoadMoreNeeded = () => {
    if (props.images.length === 0) {
      return
    }

    const layout = calculateLayout()
    const range = getVisibleRange(layout)

    const visibleRows = range.endRow - range.startRow + 1
    const visibleCols = range.endCol - range.startCol + 1
    const avgItemsPerRow = (layout.itemsPerOddRow + layout.itemsPerEvenRow) / 2
    const neededImages = Math.ceil(visibleRows * avgItemsPerRow)

    const currentUsageRatio = neededImages / props.images.length
    const horizontalExpansion = visibleCols > avgItemsPerRow * 1.5

    if (currentUsageRatio > 0.6 || (horizontalExpansion && currentUsageRatio > 0.4)) {
      emit('loadMore')
    }
  }

  const handleResize = () => {
    updateLayout()
  }

  watch(
    () => props.images.length,
    (newLength, oldLength) => {
      if (newLength !== oldLength) {
        resetImageTracking()
        if (newLength > 0) {
          nextTick(() => {
            updateVisibleContent()
          })
        }
      }
    },
    { immediate: false }
  )

  onMounted(() => {
    document.addEventListener('mousemove', handleMouseMove)
    document.addEventListener('mouseup', handleMouseUp)
    document.addEventListener('touchmove', handleTouchMove)
    document.addEventListener('touchend', handleTouchEnd)
    window.addEventListener('resize', handleResize)

    updateLayout()
  })

  onUnmounted(() => {
    document.removeEventListener('mousemove', handleMouseMove)
    document.removeEventListener('mouseup', handleMouseUp)
    document.removeEventListener('touchmove', handleTouchMove)
    document.removeEventListener('touchend', handleTouchEnd)
    window.removeEventListener('resize', handleResize)
  })

  defineExpose({
    updateVisibleContent,
    updateLayout,
  })
</script>

<template>
  <div
    ref="canvasContainer"
    class="canvas-container"
    :class="{ dragging: isDragging }"
    @mousedown="handleMouseDown"
    @touchstart="handleTouchStart"
  >
    <div ref="honeycombCanvas" class="honeycomb-canvas" :style="canvasStyle">
      <HoneycombCell
        v-for="cell in visibleCells"
        :key="`${cell.row}_${cell.col}`"
        :image-data="cell.imageData"
        :position="cell.position"
        :transform="cell.transform"
        :current-transform="cell.currentTransform"
        :row="cell.row"
        :col="cell.col"
        @click="handleCellClick"
        @mouseenter="handleCellHover"
        @mouseleave="handleCellLeave"
      />
    </div>
  </div>
</template>

<style scoped>
  .canvas-container {
    width: 100vw;
    height: 100vh;
    position: relative;
    overflow: hidden;
    cursor:
      url('data:image/svg+xml;utf8,<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"><defs><linearGradient id="grad" x1="0%" y1="0%" x2="100%" y2="100%"><stop offset="0%" style="stop-color:%2300d4ff;stop-opacity:1" /><stop offset="100%" style="stop-color:%23ff1493;stop-opacity:1" /></linearGradient></defs><path d="M12 2L15.09 8.26L22 9L17 14L18.18 21L12 17.77L5.82 21L7 14L2 9L8.91 8.26L12 2Z" fill="url(%23grad)" stroke="%23ffffff" stroke-width="1"/></svg>')
        12 12,
      grab;
    perspective: 1200px;
    perspective-origin: center center;
  }

  .canvas-container.dragging {
    cursor:
      url('data:image/svg+xml;utf8,<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"><defs><linearGradient id="grad" x1="0%" y1="0%" x2="100%" y2="100%"><stop offset="0%" style="stop-color:%23ff1493;stop-opacity:1" /><stop offset="100%" style="stop-color:%2300d4ff;stop-opacity:1" /></linearGradient></defs><path d="M12 2L15.09 8.26L22 9L17 14L18.18 21L12 17.77L5.82 21L7 14L2 9L8.91 8.26L12 2Z" fill="url(%23grad)" stroke="%23ffffff" stroke-width="1"/><circle cx="12" cy="12" r="8" fill="none" stroke="url(%23grad)" stroke-width="2" opacity="0.5"><animate attributeName="r" values="8;12;8" dur="1s" repeatCount="indefinite"/></circle></svg>')
        12 12,
      grabbing;
  }

  .honeycomb-canvas {
    position: relative;
    min-width: 100vw;
    min-height: 100vh;
    transform-style: preserve-3d;
    transform: rotateY(0deg) rotateX(0deg);
  }
</style>
