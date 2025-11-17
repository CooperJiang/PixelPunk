import { ref, type Ref } from 'vue'
import type { WatermarkConfig } from '../types'

export function useImageLoader(previewImageProp: string | undefined, config: Ref<WatermarkConfig>, onLoad: () => void) {
  const previewImageData = ref<HTMLImageElement>()
  const watermarkImageData = ref<HTMLImageElement>()
  const droppedImageUrl = ref<string>()

  const loadPreviewImage = () => {
    const imageUrl = droppedImageUrl.value || previewImageProp

    if (imageUrl) {
      const img = new Image()
      img.crossOrigin = 'anonymous'
      img.onload = () => {
        previewImageData.value = img
        onLoad()
      }
      img.onerror = () => {
        previewImageData.value = undefined
        onLoad()
      }
      img.src = imageUrl
    } else {
      previewImageData.value = undefined
      onLoad()
    }
  }

  const loadWatermarkImage = () => {
    if (config.value.type === 'image' && config.value.fileBase64) {
      const img = new Image()
      img.crossOrigin = 'anonymous'
      img.onload = () => {
        watermarkImageData.value = img
        onLoad()
      }
      img.onerror = () => {
        watermarkImageData.value = undefined
        onLoad()
      }
      img.src = `data:image/png;base64,${config.value.fileBase64}`
    } else {
      watermarkImageData.value = undefined
    }
  }

  return {
    previewImageData,
    watermarkImageData,
    droppedImageUrl,
    loadPreviewImage,
    loadWatermarkImage,
  }
}
