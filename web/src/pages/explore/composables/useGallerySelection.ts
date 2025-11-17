import { ref, type Ref } from 'vue'
import type { FileInfo } from '@/api/types'
import { useToast } from '@/components/Toast/useToast'
import type { ShareButtonRef } from '../types'
import { useTexts } from '@/composables/useTexts'

export function useGallerySelection(images: Ref<FileInfo[]>) {
  const toast = useToast()
  const { $t } = useTexts()

  const selectMode = ref(false)
  const shareButton = ref<ShareButtonRef>()
  const isAllImagesSelected = ref(false)
  const showShareDialog = ref(false)

  const handleSelectModeChange = (value: boolean) => {
    selectMode.value = value
  }

  const isFileSelected = (file: FileInfo) => shareButton.value?.isFileSelected?.(file)

  const handleFileSelect = (file: FileInfo) => {
    shareButton.value?.toggleFileSelect(file)
  }

  const toggleSelectAllImages = () => {
    if (!selectMode.value) return

    isAllImagesSelected.value = !isAllImagesSelected.value
    images.value.forEach((image) => {
      const selected = isFileSelected(image)
      if (isAllImagesSelected.value && !selected) {
        handleFileSelect(image)
      } else if (!isAllImagesSelected.value && selected) {
        handleFileSelect(image)
      }
    })
  }

  const invertImageSelection = () => {
    if (!selectMode.value) return
    images.value.forEach((image) => handleFileSelect(image))
  }

  const handleShareCreated = (_share: unknown) => {
    shareButton.value?.cancelSelectMode?.()
    showShareDialog.value = false
    toast.success($t('explore.explore.toast.shareCreated'))
  }

  const setShareButtonRef = (el: ShareButtonRef) => {
    shareButton.value = el || undefined
  }

  const handleKeydown = (event: KeyboardEvent) => {
    if (event.key === 'Escape' && selectMode.value) {
      shareButton.value?.cancelSelectMode?.()
      event.preventDefault()
    }
  }

  const attachEscListener = () => {
    window.addEventListener('keydown', handleKeydown)
  }

  const detachEscListener = () => {
    window.removeEventListener('keydown', handleKeydown)
  }

  return {
    selectMode,
    shareButton,
    isAllImagesSelected,
    showShareDialog,

    handleSelectModeChange,
    isFileSelected,
    handleFileSelect,
    toggleSelectAllImages,
    invertImageSelection,
    handleShareCreated,
    setShareButtonRef,
    attachEscListener,
    detachEscListener,
  }
}
