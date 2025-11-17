import { computed, ref } from 'vue'
import type { ImageInfo } from '@/api/admin/files'

export function useFileSelection() {
  const selectMode = ref(false)
  const selectedImages = ref<string[]>([])

  const startSelectMode = () => {
    selectMode.value = true
  }

  const cancelSelectMode = () => {
    selectMode.value = false
    selectedImages.value = []
  }

  const toggleImageSelection = (id: string) => {
    const index = selectedImages.value.indexOf(id)
    if (index > -1) {
      selectedImages.value.splice(index, 1)
    } else {
      selectedImages.value.push(id)
    }
  }

  const createSelectionHelpers = (images: () => ImageInfo[]) => {
    const isAllSelected = computed(() => {
      if (images().length === 0) {
        return false
      }
      const currentIds = new Set(images().map((i) => i.id))
      for (const id of currentIds) {
        if (!selectedImages.value.includes(id)) {
          return false
        }
      }
      return true
    })

    const hasCrossPageSelection = computed(() => {
      if (selectedImages.value.length === 0) {
        return false
      }
      const currentIds = new Set(images().map((i) => i.id))
      return selectedImages.value.some((id) => !currentIds.has(id))
    })

    const canInvertSelection = computed(() => !hasCrossPageSelection.value)

    const toggleSelectAll = () => {
      const currentIds = images().map((img) => img.id)
      if (isAllSelected.value) {
        selectedImages.value = selectedImages.value.filter((id) => !currentIds.includes(id))
      } else {
        const set = new Set(selectedImages.value)
        currentIds.forEach((id) => set.add(id))
        selectedImages.value = Array.from(set)
      }
    }

    const selectInvert = () => {
      if (!canInvertSelection.value) {
        return
      }
      const allIds = images().map((img) => img.id)
      const set = new Set(selectedImages.value)
      const invertedCurrent: string[] = []
      for (const id of allIds) {
        if (!set.has(id)) {
          invertedCurrent.push(id)
        }
      }
      const otherPageSelected = selectedImages.value.filter((id) => !allIds.includes(id))
      selectedImages.value = [...otherPageSelected, ...invertedCurrent]
    }

    return {
      isAllSelected,
      hasCrossPageSelection,
      canInvertSelection,
      toggleSelectAll,
      selectInvert,
    }
  }

  return {
    selectMode,
    selectedImages,
    startSelectMode,
    cancelSelectMode,
    toggleImageSelection,
    createSelectionHelpers,
  }
}
