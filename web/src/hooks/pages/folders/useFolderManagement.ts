import { ref } from 'vue'
import { folderApi } from '@/api'
import { visibilityFolder } from '@/api/folder'
import type { FolderInfo } from '@/api/types/index'
import { useToast } from '@/components/Toast/useToast'
import { useTexts } from '@/composables/useTexts'

/**
 * 文件夹管理 Hook
 * 负责文件夹的CRUD操作、可见性切换等功能
 */
export function useFolderManagement() {
  const toast = useToast()
  const { $t } = useTexts()

  const folders = ref<FolderInfo[]>([])
  const images = ref<any[]>([])
  const pagination = ref<any>(null)
  const loading = ref(false)

  const createDialogVisible = ref(false)
  const deleteDialogVisible = ref(false)
  const deleteTarget = ref<any>(null)
  const deleteType = ref<'folder' | 'image'>('folder')
  const dialogMode = ref<'create' | 'edit'>('create')
  const currentFolder = ref<FolderInfo | undefined>()
  const isDeleting = ref(false)

  const loadFolders = async (currentFolderId?: string) => {
    try {
      loading.value = true
      const result = await folderApi.getFolderListWithImages(currentFolderId)

      if (result.success) {
        const { data } = result

        folders.value = data.folders || []

        const files = data.files || data.images || data.items || []
        images.value = files
        pagination.value = data.pagination
      }
    } catch (error) {
      console.error('Load folders failed:', error)
    } finally {
      setTimeout(() => {
        loading.value = false
      }, 300)
    }
  }

  const showCreateDialog = () => {
    dialogMode.value = 'create'
    currentFolder.value = undefined
    createDialogVisible.value = true
  }

  const handleEdit = (folder: FolderInfo) => {
    dialogMode.value = 'edit'
    currentFolder.value = folder
    createDialogVisible.value = true
  }

  const handleFolderCreated = (folder: FolderInfo, currentFolderId?: string) => {
    loadFolders(currentFolderId)
  }

  const handleFolderUpdated = (folder: FolderInfo, currentFolderId?: string) => {
    loadFolders(currentFolderId)
  }

  const showDeleteDialog = (target: any, type: 'folder' | 'image') => {
    deleteTarget.value = target
    deleteType.value = type
    deleteDialogVisible.value = true
  }

  const confirmDelete = async (currentFolderId?: string) => {
    if (!deleteTarget.value) {
      return
    }

    isDeleting.value = true
    try {
      if (deleteType.value === 'folder') {
        await folderApi.deleteFolder(deleteTarget.value.id)
        toast.success($t('folders.folderManagement.deleteSuccess', { name: deleteTarget.value.name }))
        loadFolders(currentFolderId)
      } else {
      }
    } catch (_error) {
    } finally {
      isDeleting.value = false
    }

    deleteDialogVisible.value = false
    deleteTarget.value = null
  }

  const cancelDelete = () => {
    deleteDialogVisible.value = false
    deleteTarget.value = null
  }

  const handleDelete = async (folder: FolderInfo) => {
    showDeleteDialog(folder, 'folder')
  }

  const handleToggleFolderVisibility = async (folder: any, event: Event, currentFolderId?: string) => {
    event.stopPropagation()

    try {
      await visibilityFolder(folder.id)
      const newVisibility =
        folder.permission === 'public'
          ? $t('folders.folderManagement.visibility.private')
          : $t('folders.folderManagement.visibility.public')
      toast.success(
        $t('folders.folderManagement.visibility.switched', {
          name: folder.name,
          level: newVisibility,
        })
      )
      loadFolders(currentFolderId)
    } catch (_error: unknown) {}
  }

  return {
    folders,
    images,
    pagination,
    loading,

    createDialogVisible,
    deleteDialogVisible,
    deleteTarget,
    deleteType,
    dialogMode,
    currentFolder,
    isDeleting,

    loadFolders,

    showCreateDialog,
    handleEditFolder: handleEdit,
    handleFolderCreated,
    handleFolderUpdated,

    showDeleteDialog,
    confirmDelete,
    cancelDelete,
    handleDeleteFolder: handleDelete,

    handleToggleFolderVisibility,
  }
}
