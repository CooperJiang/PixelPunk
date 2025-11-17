import { computed, ref } from 'vue'
import type { ImageDetail } from './types'

/**
 * 文件模态框管理Hook
 * 提供文件预览、详情模态框的状态管理和导航功能
 */
export function useFileModals() {
  const showFilePreview = ref(false)
  const showDetailModal = ref(false)
  const selectedImageDetails = ref<ImageDetail | null>(null)
  const currentPreviewIndex = ref(0)

  const imageList = ref<ImageDetail[]>([])

  const hasSelectedImage = computed(() => selectedImageDetails.value !== null)
  const hasPrevImage = computed(() => currentPreviewIndex.value > 0)
  const hasNextImage = computed(() => imageList.value.length > 0 && currentPreviewIndex.value < imageList.value.length - 1)
  const currentImageInfo = computed(() => {
    if (!imageList.value.length || currentPreviewIndex.value < 0) {
      return null
    }
    return imageList.value[currentPreviewIndex.value]
  })

  const setImageList = (images: ImageDetail[]) => {
    imageList.value = images
  }

  const openFilePreview = (fileId: string, images?: ImageDetail[]) => {
    if (images) {
      setImageList(images)
    }

    const imageIndex = imageList.value.findIndex((img) => img.id === fileId)
    if (imageIndex !== -1) {
      currentPreviewIndex.value = imageIndex
      selectedImageDetails.value = imageList.value[imageIndex]
      showFilePreview.value = true
    }
  }

  const closeFilePreview = () => {
    showFilePreview.value = false
    selectedImageDetails.value = null
  }

  const openImageDetails = (fileId: string, images?: ImageDetail[]) => {
    if (images) {
      setImageList(images)
    }

    const image = imageList.value.find((img) => img.id === fileId)
    if (image) {
      selectedImageDetails.value = image
      showDetailModal.value = true
    }
  }

  const closeImageDetails = () => {
    showDetailModal.value = false
    selectedImageDetails.value = null
  }

  const switchToImageAtIndex = (index: number) => {
    if (index < 0 || index >= imageList.value.length) {
      return false
    }

    currentPreviewIndex.value = index
    selectedImageDetails.value = imageList.value[index]
    return true
  }

  const goToPrevImage = () => {
    if (hasPrevImage.value) {
      return switchToImageAtIndex(currentPreviewIndex.value - 1)
    }
    return false
  }

  const goToNextImage = () => {
    if (hasNextImage.value) {
      return switchToImageAtIndex(currentPreviewIndex.value + 1)
    }
    return false
  }

  const setCurrentImage = (image: ImageDetail) => {
    selectedImageDetails.value = image
    const index = imageList.value.findIndex((img) => img.id === image.id)
    if (index !== -1) {
      currentPreviewIndex.value = index
    }
  }

  const handleKeyboardNavigation = (event: KeyboardEvent) => {
    if (!showFilePreview.value) {
      return
    }

    switch (event.key) {
      case 'ArrowLeft':
        event.preventDefault()
        goToPrevImage()
        break
      case 'ArrowRight':
        event.preventDefault()
        goToNextImage()
        break
      case 'Escape':
        event.preventDefault()
        closeFilePreview()
        break
    }
  }

  const closeAllModals = () => {
    closeFilePreview()
    closeImageDetails()
  }

  const resetModalState = () => {
    closeAllModals()
    currentPreviewIndex.value = 0
    selectedImageDetails.value = null
    imageList.value = []
  }

  return {
    showFilePreview,
    showDetailModal,
    selectedImageDetails,
    currentPreviewIndex,
    imageList,

    hasSelectedImage,
    hasPrevImage,
    hasNextImage,
    currentImageInfo,

    setImageList,
    openFilePreview,
    closeFilePreview,
    openImageDetails,
    closeImageDetails,
    switchToImageAtIndex,
    goToPrevImage,
    goToNextImage,
    setCurrentImage,
    handleKeyboardNavigation,
    closeAllModals,
    resetModalState,
  }
}
