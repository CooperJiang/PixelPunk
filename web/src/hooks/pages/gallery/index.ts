/* Gallery页面Hooks统一导出 */

export * from './types'
export { useGalleryCore } from './useGalleryCore'
export { useFileActions } from './useFileActions'
export { useFileModals } from './useFileModals'
export { useFileSelection } from './useFileSelection'

/* 默认导出：完整的Gallery功能Hook */
import { useGalleryCore } from './useGalleryCore'
import { useFileActions } from './useFileActions'
import { useFileModals } from './useFileModals'
import { useFileSelection } from './useFileSelection'
import type { GalleryOptions } from './types'

/**
 * 完整的Gallery功能Hook
 * 集成所有Gallery相关功能
 */
export function useGallery(options: GalleryOptions = {}) {
  const galleryCore = useGalleryCore(options)

  const imageActions = useFileActions()

  const imageModals = useFileModals()

  const imageSelection = useFileSelection()

  const syncImageList = () => {
    imageModals.setImageList(galleryCore.images.value)
    imageSelection.setImageList(galleryCore.images.value)
  }

  const loadImages = async (resetPage = false) => {
    await galleryCore.loadImages(resetPage)
    syncImageList()
  }

  const refreshImages = async () => {
    await galleryCore.refreshImages()
    syncImageList()
  }

  const deleteImageWithRefresh = async (id: string) => {
    const success = await imageActions.deleteFile(id)
    if (success) {
      await refreshImages()
    }
    return success
  }

  const deleteMultipleImagesWithRefresh = async (fileIds: string[]) => {
    const success = await imageActions.deleteMultipleImages(fileIds)
    if (success) {
      await refreshImages()
      imageSelection.clearSelection()
    }
    return success
  }

  const toggleImageVisibilityWithRefresh = async (id: string, currentAccessLevel: string) => {
    const result = await imageActions.toggleImageVisibility(id, currentAccessLevel)
    if (result.success) {
      const image = galleryCore.findImageById(id)
      if (image && result.newAccessLevel) {
        image.access_level = result.newAccessLevel
        syncImageList()
      }
    }
    return result
  }

  return {
    ...galleryCore,

    ...imageActions,

    ...imageModals,

    ...imageSelection,

    loadImages,
    refreshImages,
    deleteImage: deleteImageWithRefresh,
    deleteMultipleImages: deleteMultipleImagesWithRefresh,
    toggleImageVisibility: toggleImageVisibilityWithRefresh,

    syncImageList,
  }
}
