import { ref, type Ref } from 'vue'

export function useFileDrop(droppedImageUrl: Ref<string | undefined>, onDrop: () => void) {
  const isDragOver = ref(false)

  const handleDragOver = (e: DragEvent) => {
    e.preventDefault()
    isDragOver.value = true
  }

  const handleDragLeave = () => {
    isDragOver.value = false
  }

  const handleDrop = (e: DragEvent) => {
    e.preventDefault()
    isDragOver.value = false

    const files = e.dataTransfer?.files
    if (files && files.length > 0) {
      const file = files[0]
      if (file.type.startsWith('image/')) {
        const reader = new FileReader()
        reader.onload = (event) => {
          const result = event.target?.result as string
          if (result) {
            droppedImageUrl.value = result
            onDrop()
          }
        }
        reader.readAsDataURL(file)
      }
    }
  }

  return {
    isDragOver,
    handleDragOver,
    handleDragLeave,
    handleDrop,
  }
}
