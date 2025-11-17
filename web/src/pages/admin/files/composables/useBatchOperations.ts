import { ref } from 'vue'
import { useToast } from '@/components/Toast/useToast'
import { useTexts } from '@/composables/useTexts'
import { deleteAdminFile, setFileRecommendation, batchSetFileRecommendation, batchDeleteAdminFiles } from '@/api/admin/files'

export function useBatchOperations() {
  const toast = useToast()
  const { $t } = useTexts()
  const isDeleting = ref(false)
  const isBatchProcessing = ref(false)
  const showBatchDeleteConfirm = ref(false)

  const batchSetRecommendation = async (
    selectedImages: string[],
    recommended: boolean,
    images: ImageInfo[],
    onSuccess?: () => void
  ) => {
    if (selectedImages.length === 0) {
      return
    }

    isBatchProcessing.value = true

    try {
      const result = await batchSetFileRecommendation(selectedImages, recommended)

      if (!result.success) {
        toast.error($t('admin.files.batch.operationFailed'))
        return
      }

      const { data } = result

      if (data.success_count > 0) {
        const action = recommended ? $t('admin.files.actions.recommend') : $t('admin.files.actions.unrecommend')
        const failPart = data.fail_count > 0 ? $t('admin.files.batch.failCount', { count: data.fail_count }) : ''
        toast.success(
          $t('admin.files.batch.recommendSuccess', {
            action,
            successCount: data.success_count,
            failPart,
          })
        )

        data.updated_files.forEach((updatedFile) => {
          const imageIndex = images.findIndex((img) => img.id === updatedFile.id)
          if (imageIndex !== -1) {
            images[imageIndex].is_recommended = updatedFile.is_recommended
          }
        })
      }

      if (data.fail_count > 0 && data.success_count === 0) {
        const action = recommended ? $t('admin.files.actions.recommend') : $t('admin.files.actions.unrecommend')
        toast.error($t('admin.files.batch.recommendFailed', { action }))
      }

      onSuccess?.()
    } catch (_error) {
      toast.error($t('admin.files.batch.operationError'))
    } finally {
      isBatchProcessing.value = false
    }
  }

  const confirmBatchDelete = (selectedImages: string[]) => {
    if (selectedImages.length === 0) {
      return
    }
    showBatchDeleteConfirm.value = true
  }

  const batchDelete = async (selectedImages: string[], images: ImageInfo[], onSuccess?: (deletedIds: string[]) => void) => {
    if (selectedImages.length === 0) {
      return
    }

    isBatchProcessing.value = true

    try {
      const result = await batchDeleteAdminFiles(selectedImages)

      if (!result.success) {
        toast.error($t('admin.files.batch.deleteFailed'))
        return
      }

      const { data } = result

      if (data.success_count > 0) {
        const failPart = data.fail_count > 0 ? $t('admin.files.batch.failCount', { count: data.fail_count }) : ''
        toast.success(
          $t('admin.files.batch.deleteSuccess', {
            successCount: data.success_count,
            failPart,
          })
        )

        const _newImages = images.filter((img) => !data.success_ids.includes(img.id))

        onSuccess?.(data.success_ids)
      }

      if (data.fail_count > 0 && data.success_count === 0) {
        toast.error($t('admin.files.batch.deleteFailed'))
      }

      showBatchDeleteConfirm.value = false
    } catch (_error) {
      toast.error($t('admin.files.batch.deleteError'))
    } finally {
      isBatchProcessing.value = false
    }
  }

  const deleteImage = async (fileId: string, onSuccess?: () => void) => {
    if (!fileId) {
      return
    }

    isDeleting.value = true

    try {
      const result = await deleteAdminFile(fileId)

      if (!result.success) {
        toast.error($t('admin.files.single.deleteFailed'))
        return
      }

      toast.success($t('admin.files.single.deleteSuccess'))
      onSuccess?.()
    } catch (_error) {
      toast.error($t('admin.files.single.deleteError'))
    } finally {
      isDeleting.value = false
    }
  }

  const toggleRecommendation = async (id: string, recommended: boolean, images: ImageInfo[]) => {
    try {
      const result = await setFileRecommendation(id)

      if (!result.success) {
        const action = recommended ? $t('admin.files.actions.recommend') : $t('admin.files.actions.unrecommend')
        toast.error($t('admin.files.single.toggleFailed', { action }))
        return
      }

      const { data } = result
      const action = data.is_recommended ? $t('admin.files.actions.recommend') : $t('admin.files.actions.unrecommend')
      toast.success($t('admin.files.single.toggleSuccess', { action }))

      const imageIndex = images.findIndex((img) => img.id === id)
      if (imageIndex !== -1) {
        images[imageIndex].is_recommended = data.is_recommended
      }
    } catch (_error) {
      toast.error($t('admin.files.single.toggleError'))
    }
  }

  return {
    isDeleting,
    isBatchProcessing,
    showBatchDeleteConfirm,
    batchSetRecommendation,
    confirmBatchDelete,
    batchDelete,
    deleteImage,
    toggleRecommendation,
  }
}
