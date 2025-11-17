import { type Ref, computed, ref } from 'vue'
import type { FileInfo } from '../types'

export function useFileNavigation(files: Ref<FileInfo[]>, initialIndex: number = 0) {
  const currentIndex = ref(initialIndex)

  const totalFiles = computed(() => files.value.length)
  const hasMultipleFiles = computed(() => totalFiles.value > 1)
  const hasPreviousFile = computed(() => currentIndex.value > 0)
  const hasNextFile = computed(() => currentIndex.value < totalFiles.value - 1)
  const currentFile = computed(() => files.value[currentIndex.value] || null)
  const fileUrl = computed(() => currentFile.value?.full_url || '')

  const showPreviousFile = () => {
    if (hasPreviousFile.value) {
      currentIndex.value--
      return true
    }
    return false
  }

  const showNextFile = () => {
    if (hasNextFile.value) {
      currentIndex.value++
      return true
    }
    return false
  }

  const goToFile = (index: number) => {
    if (index >= 0 && index < totalFiles.value) {
      currentIndex.value = index
      return true
    }
    return false
  }

  const handleKeyboardNavigation = (key: string) => {
    switch (key) {
      case 'ArrowLeft':
        return showPreviousFile()
      case 'ArrowRight':
        return showNextFile()
      default:
        return false
    }
  }

  return {
    currentIndex,

    totalFiles,
    hasMultipleFiles,
    hasPreviousFile,
    hasNextFile,
    currentFile,
    fileUrl,

    showPreviousFile,
    showNextFile,
    goToFile,
    handleKeyboardNavigation,
  }
}
