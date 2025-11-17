import { ref } from 'vue'
import { useToast } from '@/components/Toast/useToast'
import { downloadFileQuick } from '@/utils/file/downloader'
import { fileApi } from '@/api'
import { toggleAccessLevel } from '@/api/file'
import { StorageUtil } from '@/utils/storage/storage'
import { TOKEN_KEY } from '@/constants'
import { useTexts } from '@/composables/useTexts'
import type { ImageDetail } from './types'

/**
 * 文件操作功能Hook
 * 提供文件预览、下载、删除、访问级别切换等操作
 */
export function useFileActions() {
  const toast = useToast()
  const { $t } = useTexts()
  const isActionLoading = ref(false)

  const copyImageLink = async (image: ImageDetail) => {
    try {
      if (image.access_level === 'protected') {
        toast.warning($t('gallery.fileActions.copyLink.protectedNotSupported'))
        return false
      }

      let url = image.full_url
      if (image.access_level === 'private') {
        const token = StorageUtil.get<string>(TOKEN_KEY)
        if (token) {
          const separator = url.includes('?') ? '&' : '?'
          url = `${url}${separator}token=${token}`
        }
      }

      await navigator.clipboard.writeText(url)
      const message =
        image.access_level === 'public'
          ? $t('gallery.fileActions.copyLink.publicLinkCopied')
          : $t('gallery.fileActions.copyLink.privateLinkCopied')
      toast.success(message)
      return true
    } catch (error) {
      console.error('Copy link failed:', error)
      toast.error($t('gallery.fileActions.copyLink.copyFailed'))
      return false
    }
  }

  const downloadImage = async (fileId: string, image?: ImageDetail) => {
    try {
      const fileName = image?.display_name || image?.original_name || `image_${fileId}`
      const success = await downloadFileQuick(fileId, fileName)
      return success
    } catch (error) {
      console.error('Download failed:', error)
      return false
    }
  }

  const deleteImage = async (id: string): Promise<boolean> => {
    isActionLoading.value = true

    try {
      const result = await fileApi.deleteFile(id)
      if (result.success) {
        toast.success($t('gallery.fileActions.delete.success'))
        return true
      }
      toast.error($t('gallery.fileActions.delete.failed'))
      return false
    } catch (error) {
      console.error('Delete file failed:', error)
      return false
    } finally {
      isActionLoading.value = false
    }
  }

  const toggleImageVisibility = async (
    id: string,
    currentAccessLevel: string
  ): Promise<{ success: boolean; newAccessLevel?: string }> => {
    isActionLoading.value = true

    try {
      const result = await toggleAccessLevel(id)

      if (result.success) {
        let nextLevel = ''
        let nextAccessLevel = ''
        switch (currentAccessLevel) {
          case 'public':
            nextLevel = $t('gallery.fileActions.accessLevel.private')
            nextAccessLevel = 'private'
            break
          case 'private':
            nextLevel = $t('gallery.fileActions.accessLevel.protected')
            nextAccessLevel = 'protected'
            break
          case 'protected':
            nextLevel = $t('gallery.fileActions.accessLevel.public')
            nextAccessLevel = 'public'
            break
          default:
            nextLevel = $t('gallery.fileActions.accessLevel.public')
            nextAccessLevel = 'public'
        }

        toast.success($t('gallery.fileActions.accessLevel.switched', { level: nextLevel }))
        return { success: true, newAccessLevel: nextAccessLevel }
      }
      toast.error($t('gallery.fileActions.accessLevel.switchFailed'))
      return { success: false }
    } catch (error) {
      console.error('Toggle access level failed:', error)
      return { success: false }
    } finally {
      isActionLoading.value = false
    }
  }

  const copyMultipleLinks = async (images: ImageDetail[]): Promise<boolean> => {
    try {
      const publicImages = images.filter((img) => img.access_level === 'public')

      if (publicImages.length === 0) {
        toast.warning($t('gallery.fileActions.batch.noPublicFiles'))
        return false
      }

      if (publicImages.length < images.length) {
        toast.info(
          $t('gallery.fileActions.batch.filteredFiles', {
            count: images.length - publicImages.length,
          })
        )
      }

      const urls = publicImages.map((img) => img.full_url).join('\n')
      await navigator.clipboard.writeText(urls)
      toast.success($t('gallery.fileActions.batch.linksCopied', { count: publicImages.length }))
      return true
    } catch (error) {
      console.error('Batch copy links failed:', error)
      toast.error($t('gallery.fileActions.batch.copyFailed'))
      return false
    }
  }

  const deleteMultipleImages = async (fileIds: string[]): Promise<boolean> => {
    isActionLoading.value = true

    try {
      const promises = fileIds.map((id) => fileApi.deleteFile(id))
      const results = await Promise.allSettled(promises)

      const successCount = results.filter((result) => result.status === 'fulfilled' && result.value.success).length

      if (successCount === fileIds.length) {
        toast.success($t('gallery.fileActions.batch.deleteSuccess', { count: successCount }))
        return true
      } else if (successCount > 0) {
        toast.warning(
          $t('gallery.fileActions.batch.deletePartial', {
            successCount,
            totalCount: fileIds.length,
          })
        )
        return true
      }
      toast.error($t('gallery.fileActions.batch.deleteFailed'))
      return false
    } catch (error) {
      console.error('Batch delete failed:', error)
      return false
    } finally {
      isActionLoading.value = false
    }
  }

  return {
    isActionLoading,

    copyImageLink,
    downloadImage,
    deleteImage,
    toggleImageVisibility,

    copyMultipleLinks,
    deleteMultipleImages,
  }
}
