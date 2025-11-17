import { ref } from 'vue'
import { fileApi } from '@/api'
import { toggleAccessLevel } from '@/api/file'
import { useToast } from '@/components/Toast/useToast'
import { StorageUtil } from '@/utils/storage/storage'
import { TOKEN_KEY } from '@/constants'
import { copyToClipboard } from '@/utils/file/clipboard'
import { downloadFileQuick } from '@/utils/file/downloader'
import { useTexts } from '@/composables/useTexts'
import type { ImageInfo } from '@/api/types/file'

/**
 * 文件管理 Hook
 */
export function useFileManagement() {
  const toast = useToast()
  const { $t } = useTexts()

  const previewDialogVisible = ref(false)
  const previewImage = ref(null)
  const previewImageIndex = ref(0)

  const detailDialogVisible = ref(false)
  const detailImage = ref(null)

  const batchDeleteDialogVisible = ref(false)
  const isBatchDeleting = ref(false)

  const getImageUrl = (fullPath: string, accessLevel: string = 'private') => {
    if (accessLevel === 'private' || accessLevel === 'protected') {
      const token = StorageUtil.get<string>(TOKEN_KEY)
      if (token) {
        const separator = fullPath.includes('?') ? '&' : '?'
        return `${fullPath}${separator}token=${token}`
      }
    }

    return fullPath
  }

  const getOriginalImageUrl = (image: ImageInfo) => {
    if (image.full_url) {
      return getImageUrl(image.full_url, image.access_level)
    }

    return getImageUrl(image.full_thumb_url, image.access_level)
  }

  const formatFileSize = (bytes: number) => {
    if (bytes === 0) {
      return '0 B'
    }
    const k = 1024
    const sizes = ['B', 'KB', 'MB', 'GB']
    const i = Math.floor(Math.log(bytes) / Math.log(k))
    return `${parseFloat((bytes / Math.pow(k, i)).toFixed(2))} ${sizes[i]}`
  }

  const formatDate = (dateStr: string) => {
    try {
      const date = new Date(dateStr)
      const month = String(date.getMonth() + 1).padStart(2, '0')
      const day = String(date.getDate()).padStart(2, '0')
      const hours = String(date.getHours()).padStart(2, '0')
      const minutes = String(date.getMinutes()).padStart(2, '0')
      return `${month}-${day} ${hours}:${minutes}`
    } catch {
      return dateStr
    }
  }

  const handlePreview = (image: ImageInfo, images: ImageInfo[]) => {
    previewImage.value = image
    previewImageIndex.value = images.findIndex((img) => img.id === image.id)
    if (previewImageIndex.value === -1) {
      previewImageIndex.value = 0
    }
    previewDialogVisible.value = true
  }

  const handlePreviewClose = () => {
    previewDialogVisible.value = false
  }

  const handleViewDetails = (image: ImageInfo) => {
    detailImage.value = image
    detailDialogVisible.value = true
  }

  const handleDetailClose = () => {
    detailDialogVisible.value = false
  }

  const handleCopyLink = async (image: ImageInfo) => {
    if (image.access_level === 'protected') {
      toast.warning($t('folders.fileManagement.copyLink.protectedNotSupported'))
      return
    }

    const url = getOriginalImageUrl(image)
    await copyToClipboard(url)
    const message =
      image.access_level === 'public'
        ? $t('folders.fileManagement.copyLink.publicLinkCopied')
        : $t('folders.fileManagement.copyLink.privateLinkCopied')
    toast.success(message)
  }

  const handleDownload = async (image: ImageInfo) => {
    try {
      const fileName = image?.display_name || image?.original_name || `image_${image.id}`

      const success = await downloadFileQuick(image.id, fileName)

      if (success) {
        toast.success($t('folders.fileManagement.download.success', { name: fileName }))
      } else {
        toast.error($t('folders.fileManagement.download.cancelled'))
      }
    } catch (error) {
      console.error('Download failed:', error)
    }
  }

  const handleDeleteImage = async (image: ImageInfo, showDeleteDialog: (target: any, type: 'folder' | 'image') => void) => {
    showDeleteDialog(image, 'image')
  }

  const confirmDeleteImage = async (image: any, onReloadData: () => Promise<void>) => {
    try {
      await fileApi.deleteFile(image.id)
      toast.success($t('folders.fileManagement.delete.success', { name: image.display_name }))
      await onReloadData()
    } catch (_error) {}
  }

  const handleToggleVisibility = async (image: any, onReloadData: () => Promise<void>) => {
    try {
      await toggleAccessLevel(image.id)

      let nextLevel = ''
      switch (image.access_level) {
        case 'public':
          nextLevel = $t('folders.fileManagement.accessLevel.private')
          break
        case 'private':
          nextLevel = $t('folders.fileManagement.accessLevel.protected')
          break
        case 'protected':
          nextLevel = $t('folders.fileManagement.accessLevel.public')
          break
        default:
          nextLevel = $t('folders.fileManagement.accessLevel.public')
      }

      toast.success(
        $t('folders.fileManagement.accessLevel.switched', {
          name: image.display_name,
          level: nextLevel,
        })
      )
      await onReloadData()
    } catch (_error) {}
  }

  const confirmBatchDeleteImages = (selectedBatchImages: string[]) => {
    if (selectedBatchImages.length === 0) {
      toast.info($t('folders.fileManagement.delete.selectFirst'))
      return
    }
    batchDeleteDialogVisible.value = true
  }

  const executeBatchDelete = async (
    selectedBatchImages: string[],
    onReloadData: () => Promise<void>,
    cancelBatchMode: () => void,
    setProcessing: (processing: boolean) => void
  ) => {
    if (selectedBatchImages.length === 0) {
      return
    }

    isBatchDeleting.value = true
    setProcessing(true)

    try {
      const response = await fileApi.batchDeleteFiles(selectedBatchImages)
      const { success_count, fail_count } = response
      if (success_count > 0 && fail_count === 0) {
        toast.success($t('folders.fileManagement.delete.batchSuccess', { successCount: success_count }))
      } else if (success_count > 0 && fail_count > 0) {
        toast.success(
          $t('folders.fileManagement.delete.batchPartial', {
            successCount: success_count,
            failCount: fail_count,
          })
        )
      } else if (fail_count > 0 && success_count === 0) {
        toast.error($t('folders.fileManagement.delete.batchFailed', { failCount: fail_count }))
      }
      batchDeleteDialogVisible.value = false
      await onReloadData() // 重新加载数据
      cancelBatchMode() // 退出批量模式
    } catch (error) {
      console.error('Batch delete error:', error)
    } finally {
      isBatchDeleting.value = false
      setProcessing(false)
    }
  }

  const cancelBatchDelete = () => {
    batchDeleteDialogVisible.value = false
  }

  return {
    previewDialogVisible,
    previewImage,
    previewImageIndex,

    detailDialogVisible,
    detailImage,

    batchDeleteDialogVisible,
    isBatchDeleting,

    getImageUrl,
    getOriginalImageUrl,
    formatFileSize,
    formatDate,

    handlePreview,
    handlePreviewClose,

    handleViewDetails,
    handleDetailClose,

    handleCopyLink,
    handleDownload,
    handleDeleteImage,
    confirmDeleteImage,
    handleToggleVisibility,

    confirmBatchDeleteImages,
    executeBatchDelete,
    cancelBatchDelete,
  }
}
