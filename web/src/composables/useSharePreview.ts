import { ref, type ComputedRef } from 'vue'

/**
 * 分享文件预览功能
 * 处理文件预览、翻页等操作
 */
export function useSharePreview<T>(
  allImages: ComputedRef<T[]>,
  onPreviewChange?: (index: number) => void,
  preloadFileFn?: (url: string) => void
) {
  const previewDialogVisible = ref(false)

  const previewImage = ref<any>({
    id: '',
    url: '',
    display_name: '',
  })

  const currentPreviewIndex = ref(0)

  const handleImageClick = (image: Event) => {
    previewImage.value = { ...image }

    const index = allImages.value.findIndex((img: any) => img.id === image.id)
    if (index !== -1) {
      currentPreviewIndex.value = index

      preloadAdjacentImages(index)

      if (onPreviewChange) {
        onPreviewChange(index)
      }
    } else {
      currentPreviewIndex.value = 0
    }

    previewDialogVisible.value = true
  }

  const preloadAdjacentImages = (index: number) => {
    if (!preloadFileFn) {
      return
    }

    for (let i = -3; i <= 3; i++) {
      if (i === 0) {
        continue
      } // 跳过当前文件

      const preloadIndex = index + i
      if (preloadIndex >= 0 && preloadIndex < allImages.value.length) {
        const preloadFile = allImages.value[preloadIndex] as any
        const preloadFileUrl = preloadFile.full_thumb_url || preloadFile.thumb_url || preloadFile.full_url || preloadFile.url
        if (preloadFileUrl) {
          preloadFileFn(preloadFileUrl)
        }
      }
    }
  }

  const handlePreviewChange = (index: number) => {
    if (index < 0 || index >= allImages.value.length) {
      return
    }

    currentPreviewIndex.value = index

    previewImage.value = allImages.value[index]

    if (preloadFileFn) {
      if (index + 1 < allImages.value.length) {
        const nextFile = allImages.value[index + 1] as any
        if (nextFile) {
          const fileUrl = nextFile.full_thumb_url || nextFile.thumb_url || nextFile.full_url || nextFile.url
          if (fileUrl) {
            preloadFileFn(fileUrl)
          }
        }
      }
    }

    if (onPreviewChange) {
      onPreviewChange(index)
    }
  }

  const handlePreviewClose = () => {
    previewDialogVisible.value = false
  }

  const showPreviousImage = () => {
    if (currentPreviewIndex.value > 0) {
      handlePreviewChange(currentPreviewIndex.value - 1)
    }
  }

  const showNextImage = () => {
    if (currentPreviewIndex.value < allImages.value.length - 1) {
      handlePreviewChange(currentPreviewIndex.value + 1)
    }
  }

  const handleKeydown = (e: KeyboardEvent) => {
    if (!previewDialogVisible.value) {
      return
    }

    if (e.key === 'ArrowLeft') {
      showPreviousImage()
    } else if (e.key === 'ArrowRight') {
      showNextImage()
    } else if (e.key === 'Escape') {
      handlePreviewClose()
    }
  }

  const handlePrevImage = (_index: number) => {}

  const handleNextImage = (_index: number) => {}

  return {
    previewDialogVisible,
    previewImage,
    currentPreviewIndex,
    handleImageClick,
    handlePreviewChange,
    handlePreviewClose,
    showPreviousImage,
    showNextImage,
    handleKeydown,
    handlePrevImage,
    handleNextImage,
  }
}
