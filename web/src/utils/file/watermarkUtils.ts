import type { WatermarkConfig } from '@/components/WatermarkConfig/types'

const PREVIEW_WIDTH = 2000
const PREVIEW_HEIGHT = 1200

export function convertWatermarkConfigForActualFile(
  config: WatermarkConfig,
  actualWidth: number,
  actualHeight: number
): WatermarkConfig {
  if (config.offsetUnit) {
    const converted = JSON.parse(JSON.stringify(config)) as any
    if (config.offsetUnit === 'px') {
      const canvasRatio = PREVIEW_WIDTH / PREVIEW_HEIGHT
      const imgRatio = actualWidth / actualHeight
      let previewDrawnWidth: number
      let previewDrawnHeight: number
      if (imgRatio > canvasRatio) {
        previewDrawnWidth = PREVIEW_WIDTH
        previewDrawnHeight = PREVIEW_WIDTH / imgRatio
      } else {
        previewDrawnHeight = PREVIEW_HEIGHT
        previewDrawnWidth = PREVIEW_HEIGHT * imgRatio
      }
      const scaleX = actualWidth / previewDrawnWidth
      const scaleY = actualHeight / previewDrawnHeight
      converted.offsetX = Math.round((config.offsetX || 0) * scaleX)
      converted.offsetY = Math.round((config.offsetY || 0) * scaleY)
      if (config.shadow) {
        converted.shadowOffsetX = Math.round((config.shadowOffsetX || 0) * scaleX)
        converted.shadowOffsetY = Math.round((config.shadowOffsetY || 0) * scaleY)
        converted.shadowBlur = Math.round((config.shadowBlur || 0) * Math.max(scaleX, scaleY))
      }
    }
    return converted
  }

  const canvasRatio = PREVIEW_WIDTH / PREVIEW_HEIGHT
  const imgRatio = actualWidth / actualHeight
  let previewDrawnWidth: number
  let _previewDrawnHeight: number
  if (imgRatio > canvasRatio) {
    previewDrawnWidth = PREVIEW_WIDTH
    _previewDrawnHeight = PREVIEW_WIDTH / imgRatio
  } else {
    _previewDrawnHeight = PREVIEW_HEIGHT
    previewDrawnWidth = PREVIEW_HEIGHT * imgRatio
  }
  const scaleRatio = actualWidth / previewDrawnWidth

  const converted = JSON.parse(JSON.stringify(config))

  if ((config as any).generateMode === 'frontend' && (config as any).fileBase64) {
    return converted
  }

  converted.scale = config.scale * scaleRatio
  if (config.position !== 'custom') {
    Object.assign(converted, convertFixedPositionToCustom(converted, actualWidth, actualHeight))
  }
  if (config.type === 'text') {
    converted.fontSize = Math.round((config as any).fontSize * scaleRatio)
  }
  if (converted.shadow) {
    converted.shadowBlur = Math.round((converted as any).shadowBlur * scaleRatio)
    converted.shadowOffsetX = Math.round((converted as any).shadowOffsetX * scaleRatio)
    converted.shadowOffsetY = Math.round((converted as any).shadowOffsetY * scaleRatio)
  }
  return converted
}

function convertFixedPositionToCustom(config: WatermarkConfig, actualWidth: number, actualHeight: number): WatermarkConfig {
  const canvasRatio = PREVIEW_WIDTH / PREVIEW_HEIGHT
  const imgRatio = actualWidth / actualHeight
  let previewDrawnWidth: number
  let _previewDrawnHeight: number
  if (imgRatio > canvasRatio) {
    previewDrawnWidth = PREVIEW_WIDTH
    _previewDrawnHeight = PREVIEW_WIDTH / imgRatio
  } else {
    _previewDrawnHeight = PREVIEW_HEIGHT
    previewDrawnWidth = PREVIEW_HEIGHT * imgRatio
  }
  const scaleRatio = actualWidth / previewDrawnWidth

  const margin = 20 * scaleRatio
  let customX = 0.5
  let customY = 0.5

  let watermarkWidth = 100
  let watermarkHeight = 30

  if (config.type === 'text') {
    watermarkWidth = config.text.length * config.fontSize * 0.6 * config.scale
    watermarkHeight = config.fontSize * config.scale
  } else {
    watermarkWidth = 100 * config.scale
    watermarkHeight = 30 * config.scale
  }

  switch (config.position) {
    case 'top-left':
      customX = (margin + watermarkWidth / 2) / actualWidth
      customY = (margin + watermarkHeight / 2) / actualHeight
      break
    case 'top-center':
      customX = 0.5
      customY = (margin + watermarkHeight / 2) / actualHeight
      break
    case 'top-right':
      customX = (actualWidth - margin - watermarkWidth / 2) / actualWidth
      customY = (margin + watermarkHeight / 2) / actualHeight
      break
    case 'middle-left':
      customX = (margin + watermarkWidth / 2) / actualWidth
      customY = 0.5
      break
    case 'middle-center':
      customX = 0.5
      customY = 0.5
      break
    case 'middle-right':
      customX = (actualWidth - margin - watermarkWidth / 2) / actualWidth
      customY = 0.5
      break
    case 'bottom-left':
      customX = (margin + watermarkWidth / 2) / actualWidth
      customY = (actualHeight - margin - watermarkHeight / 2) / actualHeight
      break
    case 'bottom-center':
      customX = 0.5
      customY = (actualHeight - margin - watermarkHeight / 2) / actualHeight
      break
    case 'bottom-right':
      customX = (actualWidth - margin - watermarkWidth / 2) / actualWidth
      customY = (actualHeight - margin - watermarkHeight / 2) / actualHeight
      break
  }

  customX = Math.max(0.05, Math.min(0.95, customX))
  customY = Math.max(0.05, Math.min(0.95, customY))

  return {
    ...config,
    position: 'custom',
    customX,
    customY,
  }
}

export async function getFileDimensions(file: File): Promise<{ width: number; height: number }> {
  const { createFileObjectURL } = await import('./heicConverter')
  const { useTextThemeStore } = await import('@/store/textTheme')
  const store = useTextThemeStore()

  return new Promise((resolve, reject) => {
    const img = new Image()

    createFileObjectURL(file)
      .then((objectUrl) => {
        img.onload = () => {
          resolve({
            width: img.naturalWidth,
            height: img.naturalHeight,
          })
          URL.revokeObjectURL(objectUrl)
        }

        img.onerror = () => {
          reject(new Error(store.getText('utils.file.watermark.errors.cannotLoadImage')))
          URL.revokeObjectURL(objectUrl)
        }

        img.src = objectUrl
      })
      .catch((error) => {
        reject(error)
      })
  })
}
