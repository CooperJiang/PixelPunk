import type { WatermarkConfig } from '../../types'
import type { WatermarkPosition, ImageRect } from './types'

export const getImageRect = (img: HTMLImageElement | undefined, canvasWidth: number, canvasHeight: number): ImageRect => {
  if (!img) {
    return { x: 0, y: 0, width: canvasWidth, height: canvasHeight }
  }
  const imgRatio = img.width / img.height
  const canvasRatio = canvasWidth / canvasHeight
  let drawWidth, drawHeight, drawX, drawY
  if (imgRatio > canvasRatio) {
    drawWidth = canvasWidth
    drawHeight = canvasWidth / imgRatio
    drawX = 0
    drawY = (canvasHeight - drawHeight) / 2
  } else {
    drawHeight = canvasHeight
    drawWidth = canvasHeight * imgRatio
    drawX = (canvasWidth - drawWidth) / 2
    drawY = 0
  }
  return { x: drawX, y: drawY, width: drawWidth, height: drawHeight }
}

export const calculateWatermarkPosition = (
  config: WatermarkConfig,
  watermarkImageData: HTMLImageElement | undefined,
  width: number,
  height: number,
  previewImageData: HTMLImageElement | undefined
): WatermarkPosition => {
  let watermarkWidth = 0,
    watermarkHeight = 0

  if (config.type === 'text') {
    const canvas = document.createElement('canvas')
    const ctx = canvas.getContext('2d')
    if (ctx) {
      ctx.font = `${config.fontWeight} ${config.fontSize}px ${config.fontFamily}`
      const metrics = ctx.measureText(config.text)
      watermarkWidth = metrics.width * config.scale
      watermarkHeight = config.fontSize * config.scale
    }
  } else if (watermarkImageData) {
    watermarkWidth = watermarkImageData.width * config.scale
    watermarkHeight = watermarkImageData.height * config.scale
  }

  const minDragSize = 40
  watermarkWidth = Math.max(watermarkWidth, minDragSize)
  watermarkHeight = Math.max(watermarkHeight, minDragSize)

  if (config.rotation !== 0) {
    const rad = (Math.abs(config.rotation) * Math.PI) / 180
    const cos = Math.abs(Math.cos(rad))
    const sin = Math.abs(Math.sin(rad))
    const rotatedWidth = watermarkWidth * cos + watermarkHeight * sin
    const rotatedHeight = watermarkWidth * sin + watermarkHeight * cos
    watermarkWidth = rotatedWidth
    watermarkHeight = rotatedHeight
  }

  const imgRect = getImageRect(previewImageData, width, height)

  const offsetX = config.offsetX ?? 20
  const offsetY = config.offsetY ?? 20

  let x = 0,
    y = 0

  switch (config.position) {
    case 'top-left':
      x = imgRect.x + offsetX
      y = imgRect.y + offsetY
      break
    case 'top-center':
      x = imgRect.x + (imgRect.width - watermarkWidth) / 2 + offsetX
      y = imgRect.y + offsetY
      break
    case 'top-right':
      x = imgRect.x + imgRect.width - watermarkWidth - offsetX
      y = imgRect.y + offsetY
      break
    case 'middle-left':
      x = imgRect.x + offsetX
      y = imgRect.y + (imgRect.height - watermarkHeight) / 2 + offsetY
      break
    case 'middle-center':
      x = imgRect.x + (imgRect.width - watermarkWidth) / 2 + offsetX
      y = imgRect.y + (imgRect.height - watermarkHeight) / 2 + offsetY
      break
    case 'middle-right':
      x = imgRect.x + imgRect.width - watermarkWidth - offsetX
      y = imgRect.y + (imgRect.height - watermarkHeight) / 2 + offsetY
      break
    case 'bottom-left':
      x = imgRect.x + offsetX
      y = imgRect.y + imgRect.height - watermarkHeight - offsetY
      break
    case 'bottom-center':
      x = imgRect.x + (imgRect.width - watermarkWidth) / 2 + offsetX
      y = imgRect.y + imgRect.height - watermarkHeight - offsetY
      break
    case 'bottom-right':
      x = imgRect.x + imgRect.width - watermarkWidth - offsetX
      y = imgRect.y + imgRect.height - watermarkHeight - offsetY
      break
    case 'custom':
      x = imgRect.x + offsetX
      y = imgRect.y + offsetY
      break
    default:
      x = imgRect.x + imgRect.width - watermarkWidth - offsetX
      y = imgRect.y + imgRect.height - watermarkHeight - offsetY
  }

  return { x, y, width: watermarkWidth, height: watermarkHeight }
}

export const getAnchorPosition = (position: string, imgRect: ImageRect): { x: number; y: number } => {
  const { x, y, width, height } = imgRect

  switch (position) {
    case 'top-left':
      return { x, y }
    case 'top-center':
      return { x: x + width / 2, y }
    case 'top-right':
      return { x: x + width, y }
    case 'middle-left':
      return { x, y: y + height / 2 }
    case 'middle-center':
      return { x: x + width / 2, y: y + height / 2 }
    case 'middle-right':
      return { x: x + width, y: y + height / 2 }
    case 'bottom-left':
      return { x, y: y + height }
    case 'bottom-center':
      return { x: x + width / 2, y: y + height }
    case 'bottom-right':
      return { x: x + width, y: y + height }
    default:
      return { x: x + width, y: y + height }
  }
}

export const getHoveredResizeHandle = (
  mouseX: number,
  mouseY: number,
  previewX: number,
  previewY: number,
  previewWidth: number,
  previewHeight: number
): string | null => {
  const handleSize = 12
  const threshold = handleSize / 2 + 4

  const handles = [
    { name: 'top-left', x: previewX, y: previewY },
    { name: 'top-right', x: previewX + previewWidth, y: previewY },
    { name: 'bottom-left', x: previewX, y: previewY + previewHeight },
    { name: 'bottom-right', x: previewX + previewWidth, y: previewY + previewHeight },
  ]

  for (const handle of handles) {
    const dx = mouseX - handle.x
    const dy = mouseY - handle.y
    const distance = Math.sqrt(dx * dx + dy * dy)

    if (distance <= threshold) {
      return handle.name
    }
  }

  return null
}
