import { computed, ref } from 'vue'
import { useToast } from '@/components/Toast/useToast'
import { useTexts } from '@/composables/useTexts'
import { StorageUtil } from '@/utils/storage/storage'
import { TOKEN_KEY } from '@/constants'
import type { FileInfo } from '@/api/types/file'
import type { FolderInfo } from '@/api/types/folder'

/**
 * 选择模式管理 Hook
 * 负责分享选择模式、批量操作模式的状态管理
 */
export function useSelectionMode() {
  const toast = useToast()
  const { $t } = useTexts()

  const selectMode = ref(false)

  const batchMode = ref(false)
  const selectedBatchFiles = ref<string[]>([])
  const isBatchProcessing = ref(false)

  const startBatchMode = () => {
    batchMode.value = true
    selectedBatchFiles.value = []
  }

  const cancelBatchMode = () => {
    batchMode.value = false
    selectedBatchFiles.value = []
    isBatchProcessing.value = false
  }

  const isBatchFileSelected = (file: FileInfo) => selectedBatchFiles.value.includes(file.id)

  const checkLoginStatus = (requireLogin: boolean = true) => {
    if (requireLogin) {
      const token = StorageUtil.get<string>(TOKEN_KEY)
      if (!token) {
        toast.warning($t('common.auth.loginRequired'))
        return false
      }
    }
    return true
  }

  const toggleBatchFileSelect = (file: FileInfo) => {
    const fileId = file.id
    const index = selectedBatchFiles.value.indexOf(fileId)
    if (index > -1) {
      selectedBatchFiles.value.splice(index, 1)
    } else {
      selectedBatchFiles.value.push(fileId)
    }
  }

  const isAllFoldersSelected = computed(() => (_folders: FolderInfo[]) => false)

  const isAllFilesSelected = computed(() => (files: FileInfo[]) => {
    if (batchMode.value) {
      return files.length > 0 && selectedBatchFiles.value.length === files.length
    }
    return false
  })

  const toggleSelectAllFiles = (files: FileInfo[], _requireLogin: boolean = true) => {
    if (batchMode.value) {
      if (selectedBatchFiles.value.length === files.length) {
        selectedBatchFiles.value = []
      } else {
        selectedBatchFiles.value = files.map((file) => file.id)
      }
    }
  }

  const invertFileSelection = (files: FileInfo[], _requireLogin: boolean = true) => {
    if (batchMode.value) {
      const allIds = files.map((file) => file.id)
      selectedBatchFiles.value = allIds.filter((id) => !selectedBatchFiles.value.includes(id))
    }
  }

  return {
    selectMode,

    batchMode,
    selectedBatchFiles,
    isBatchProcessing,

    startBatchMode,
    cancelBatchMode,
    isBatchFileSelected,
    toggleBatchFileSelect,

    isAllFoldersSelected,
    isAllFilesSelected,

    toggleSelectAllFiles,
    invertFileSelection,

    checkLoginStatus,
  }
}
