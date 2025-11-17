import { ref, type Ref } from 'vue'
import type { FileInfo } from '@/api/types'

export function useGalleryPreview(images: Ref<FileInfo[]>) {
  const showDetailModal = ref(false)
  const showFilePreview = ref(false)
  const selectedImage = ref<FileInfo | null>(null)
  const previewImageUrl = ref('')
  const previewImageName = ref('')
  const currentPreviewIndex = ref(0)

  function closeFilePreview() {
    showFilePreview.value = false
  }

  const handlePreviewChange = (index: number) => {
    if (index < 0 || index >= images.value.length) return
    currentPreviewIndex.value = index
    selectedImage.value = images.value[index]
  }

  const handlePrevImage = () => {
    if (currentPreviewIndex.value > 0) handlePreviewChange(currentPreviewIndex.value - 1)
  }

  const handleNextImage = () => {
    if (currentPreviewIndex.value < images.value.length - 1) handlePreviewChange(currentPreviewIndex.value + 1)
  }

  const viewImageDetails = (id: string) => {
    const image = images.value.find((img) => img.id === id)
    if (image) {
      selectedImage.value = image
      showDetailModal.value = true
    }
  }

  const viewImage = (id: string) => {
    const imageIndex = images.value.findIndex((img) => img.id === id)
    if (imageIndex !== -1) {
      currentPreviewIndex.value = imageIndex
      selectedImage.value = images.value[imageIndex]
      showFilePreview.value = true
    }
  }

  return {
    showDetailModal,
    showFilePreview,
    selectedImage,
    previewImageUrl,
    previewImageName,
    currentPreviewIndex,
    closeFilePreview,
    handlePreviewChange,
    handlePrevImage,
    handleNextImage,
    viewImageDetails,
    viewImage,
  }
}
