import { ref, type Ref } from 'vue'

export function useFileUpload(selectedFiles: Ref<File[]>, showValidation: Ref<boolean>, isDragOver: Ref<boolean>) {
  const fileInput = ref<HTMLInputElement | null>(null)

  const formatFileSize = (bytes: number) => {
    if (bytes === 0) return '0 B'
    const k = 1024
    const sizes = ['B', 'KB', 'MB', 'GB']
    const i = Math.floor(Math.log(bytes) / Math.log(k))
    return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + ' ' + sizes[i]
  }

  const triggerFileInput = () => {
    fileInput.value?.click()
  }

  const handleFileSelect = (event: Event) => {
    const target = event.target as HTMLInputElement
    if (target.files) {
      const files = Array.from(target.files)
      selectedFiles.value = [...selectedFiles.value, ...files]
      showValidation.value = false
    }
  }

  const handleDrop = (event: DragEvent) => {
    event.preventDefault()
    isDragOver.value = false
    if (event.dataTransfer?.files) {
      const files = Array.from(event.dataTransfer.files)
      const imageFiles = files.filter((file) => file.type.startsWith('image/'))
      selectedFiles.value = [...selectedFiles.value, ...imageFiles]
      showValidation.value = false
    }
  }

  const handleDragOver = () => {
    isDragOver.value = true
  }

  const handleDragLeave = () => {
    isDragOver.value = false
  }

  const removeFile = (index: number) => {
    selectedFiles.value.splice(index, 1)
  }

  const handleImageError = (event: Event) => {
    const target = event.target as HTMLImageElement
    target.style.display = 'none'
  }

  return {
    fileInput,
    formatFileSize,
    triggerFileInput,
    handleFileSelect,
    handleDrop,
    handleDragOver,
    handleDragLeave,
    removeFile,
    handleImageError,
  }
}
