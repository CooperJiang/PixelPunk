import { computed, ref } from 'vue'
import type { ImageDetail } from './types'

/**
 * 文件选择模式管理Hook
 * 提供多选、全选、反选等选择功能
 */
export function useFileSelection() {
  const selectMode = ref(false)
  const selectedFiles = ref<Set<string>>(new Set())

  const fileList = ref<ImageDetail[]>([])

  const selectedCount = computed(() => selectedFiles.value.size)
  const hasSelectedFiles = computed(() => selectedCount.value > 0)
  const isAllSelected = computed(() => fileList.value.length > 0 && selectedCount.value === fileList.value.length)
  const isPartiallySelected = computed(() => selectedCount.value > 0 && selectedCount.value < fileList.value.length)
  const selectedFileList = computed(() => fileList.value.filter((file) => selectedFiles.value.has(file.id)))

  const setFileList = (files: ImageDetail[]) => {
    fileList.value = files
    const existingIds = new Set(files.map((file) => file.id))
    selectedFiles.value.forEach((id) => {
      if (!existingIds.has(id)) {
        selectedFiles.value.delete(id)
      }
    })
  }

  const enableSelectMode = () => {
    selectMode.value = true
  }

  const disableSelectMode = () => {
    selectMode.value = false
    clearSelection()
  }

  const toggleSelectMode = () => {
    if (selectMode.value) {
      disableSelectMode()
    } else {
      enableSelectMode()
    }
  }

  const isFileSelected = (fileId: string): boolean => selectedFiles.value.has(fileId)

  const selectFile = (fileId: string) => {
    selectedFiles.value.add(fileId)
  }

  const unselectFile = (fileId: string) => {
    selectedFiles.value.delete(fileId)
  }

  const toggleFileSelection = (fileId: string) => {
    if (isFileSelected(fileId)) {
      unselectFile(fileId)
    } else {
      selectFile(fileId)
    }
  }

  const selectFiles = (fileIds: string[]) => {
    fileIds.forEach((id) => selectFile(id))
  }

  const unselectFiles = (fileIds: string[]) => {
    fileIds.forEach((id) => unselectFile(id))
  }

  const selectAll = () => {
    fileList.value.forEach((file) => selectFile(file.id))
  }

  const unselectAll = () => {
    selectedFiles.value.clear()
  }

  const toggleSelectAll = () => {
    if (isAllSelected.value) {
      unselectAll()
    } else {
      selectAll()
    }
  }

  const invertSelection = () => {
    fileList.value.forEach((file) => {
      toggleFileSelection(file.id)
    })
  }

  const clearSelection = () => {
    selectedFiles.value.clear()
  }

  const selectFilesByCondition = (condition: (file: ImageDetail) => boolean) => {
    fileList.value.forEach((file) => {
      if (condition(file)) {
        selectFile(file.id)
      }
    })
  }

  const selectFilesByAccessLevel = (accessLevel: string) => {
    selectFilesByCondition((file) => file.access_level === accessLevel)
  }

  const selectFilesByFormat = (format: string) => {
    selectFilesByCondition((file) => file.format.toLowerCase() === format.toLowerCase())
  }

  const getSelectionStats = () => {
    const selected = selectedFileList.value
    const totalSize = selected.reduce((sum, file) => sum + file.size, 0)
    const accessLevels = selected.reduce(
      (acc, file) => {
        acc[file.access_level] = (acc[file.access_level] || 0) + 1
        return acc
      },
      {} as Record<string, number>
    )

    return {
      count: selected.length,
      totalSize,
      accessLevels,
      files: selected,
    }
  }

  const exportSelectedFiles = () => ({
    files: selectedFileList.value,
    stats: getSelectionStats(),
  })

  return {
    selectMode,
    selectedFiles,
    fileList,

    selectedCount,
    hasSelectedFiles,
    isAllSelected,
    isPartiallySelected,
    selectedFileList,

    setFileList,
    enableSelectMode,
    disableSelectMode,
    toggleSelectMode,
    isFileSelected,
    selectFile,
    unselectFile,
    toggleFileSelection,
    selectFiles,
    unselectFiles,
    selectAll,
    unselectAll,
    toggleSelectAll,
    invertSelection,
    clearSelection,
    selectFilesByCondition,
    selectFilesByAccessLevel,
    selectFilesByFormat,
    getSelectionStats,
    exportSelectedFiles,
  }
}
