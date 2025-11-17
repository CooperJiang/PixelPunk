import { ref } from 'vue'
import { folderApi } from '@/api'
import { reorderFiles } from '@/api/file'
import { useToast } from '@/components/Toast/useToast'
import { useTexts } from '@/composables/useTexts'

/**
 * 拖拽排序管理 Hook
 * 负责文件夹和文件的拖拽排序功能
 */
export function useDragSort() {
  const toast = useToast()
  const { $t } = useTexts()

  const isDragging = ref(false)
  const dragTimeout = ref<number | null>(null)

  const isFileDragging = ref(false)
  const fileDragTimeout = ref<number | null>(null)

  const onFolderDragStart = (_event: Event) => {
    isDragging.value = true
  }

  const onFolderDragEnd = async (
    event: Event,
    folders: any[],
    currentFolderId: string | undefined,
    onReloadData: () => Promise<void>
  ) => {
    const { oldIndex, newIndex } = event
    isDragging.value = false

    if (oldIndex === newIndex) {
      return
    }

    if (dragTimeout.value) {
      clearTimeout(dragTimeout.value)
    }

    dragTimeout.value = setTimeout(async () => {
      try {
        const orderedFolderIds = folders.map((folder) => folder.id)

        await folderApi.reorderFolders({
          parent_id: currentFolderId || '',
          folder_ids: orderedFolderIds,
        })

        await onReloadData()

        toast.success($t('folders.dragSort.folderSortUpdated'))
      } catch (error) {
        console.error('Folder sort failed:', error)

        await onReloadData()
      }
    }, 500) as unknown as number
  }

  const onFileDragStart = (_event: Event) => {
    isFileDragging.value = true
  }

  const onFileDragEnd = async (
    event: Event,
    files: any[],
    currentFolderId: string | undefined,
    onReloadData: () => Promise<void>
  ) => {
    const { oldIndex, newIndex } = event
    isFileDragging.value = false

    if (oldIndex === newIndex) {
      return
    }

    if (fileDragTimeout.value) {
      clearTimeout(fileDragTimeout.value)
    }

    fileDragTimeout.value = setTimeout(async () => {
      try {
        const orderedFileIds = files.map((file) => file.id)

        await reorderFiles({
          folder_id: currentFolderId || '',
          file_ids: orderedFileIds,
        })

        await onReloadData()

        toast.success($t('folders.dragSort.fileSortUpdated'))
      } catch (error) {
        console.error('File sort failed:', error)

        await onReloadData()
      }
    }, 500) as unknown as number
  }

  const cleanup = () => {
    if (dragTimeout.value) {
      clearTimeout(dragTimeout.value)
      dragTimeout.value = null
    }
    if (fileDragTimeout.value) {
      clearTimeout(fileDragTimeout.value)
      fileDragTimeout.value = null
    }
  }

  return {
    isDragging,

    isFileDragging,

    onFolderDragStart,
    onFolderDragEnd,

    onFileDragStart,
    onFileDragEnd,

    cleanup,
  }
}
