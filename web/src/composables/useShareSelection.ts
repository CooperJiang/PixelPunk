import { ref, type ComputedRef } from 'vue'
import { getFileDownloadUrl } from '@/utils/file/fileUtils'
import { useTexts } from '@/composables/useTexts'
import JSZip from 'jszip'

/**
 * 分享多选功能
 * 处理文件多选、批量下载等操作
 */
export function useShareSelection(
  allImages: ComputedRef<any[]>,
  accessToken?: ComputedRef<string>,
  onError?: (message: string) => void,
  onSuccess?: (message: string) => void,
  onInfo?: (message: string) => void
) {
  const { $t } = useTexts()
  const isSelectionMode = ref(false)

  const selectedImagesSet = ref(new Set<string>())

  const selectedImages = ref<string[]>([])

  const toggleSelectionMode = (value?: boolean) => {
    isSelectionMode.value = value !== undefined ? value : !isSelectionMode.value

    if (!isSelectionMode.value) {
      selectedImagesSet.value.clear()
      selectedImages.value = []
    }
  }

  const toggleImageSelection = (fileId: string | any) => {
    if (!isSelectionMode.value) {
      return
    }

    const id = typeof fileId === 'string' ? fileId : fileId && fileId.id ? fileId.id : null

    if (!id) {
      return
    }

    const wasSelected = selectedImagesSet.value.has(id)

    if (wasSelected) {
      selectedImagesSet.value.delete(id)
    } else {
      selectedImagesSet.value.add(id)
    }

    selectedImages.value = Array.from(selectedImagesSet.value)
  }

  const isImageSelected = (fileId: string) => selectedImagesSet.value.has(fileId)

  const selectAll = () => {
    selectedImagesSet.value = new Set(allImages.value.map((image) => image.id))
    selectedImages.value = Array.from(selectedImagesSet.value)
  }

  const deselectAll = () => {
    selectedImagesSet.value.clear()
    selectedImages.value = []
  }

  const getSelectedImagesData = () => {
    const selectedIds = selectedImagesSet.value
    return allImages.value.filter((img) => selectedIds.has(img.id))
  }

  const downloadSelected = async () => {
    if (selectedImages.value.length === 0) {
      if (onInfo) {
        onInfo($t('pages.share.selection.noFilesSelected'))
      }
      return
    }

    try {
      const selectedImagesData = getSelectedImagesData()

      if (selectedImagesData.length === 1) {
        const accessTokenValue = accessToken?.value || ''
        const fileUrl = getFileDownloadUrl(selectedImagesData[0], accessTokenValue)
        downloadSingleFile(fileUrl, selectedImagesData[0].display_name || 'file')
        return
      }

      if (onInfo) {
        onInfo($t('pages.share.selection.preparing', { count: selectedImagesData.length.toString() }))
      }

      const zip = new JSZip()

      const imagesFolder = zip.folder('images')
      if (!imagesFolder) {
        if (onError) {
          onError($t('pages.share.selection.createZipFailed'))
        }
        return
      }

      let completedCount = 0
      const totalCount = selectedImagesData.length

      const downloadPromises = selectedImagesData.map(async (image, index) => {
        try {
          const accessTokenValue = accessToken?.value || ''
          const fileUrl = getFileDownloadUrl(image, accessTokenValue)

          let filename = image.display_name || `image-${index + 1}`

          if (!filename.match(/\.(jpg|jpeg|png|gif|webp|bmp)$/i)) {
            const extension = fileUrl.split('.').pop()?.split('?')[0] || 'jpg'
            filename += `.${extension}`
          }

          const response = await fetch(fileUrl)
          if (!response.ok) {
            throw new Error(`Failed to fetch ${filename}`)
          }

          const blob = await response.blob()

          imagesFolder.file(filename, blob)

          completedCount++

          if (completedCount % 3 === 0 || completedCount === totalCount) {
            if (onInfo) {
              onInfo($t('pages.share.selection.progress', { completed: completedCount.toString(), total: totalCount.toString() }))
            }
          }

          return true
        } catch {
          return false
        }
      })

      await Promise.all(downloadPromises)

      const zipContent = await zip.generateAsync({
        type: 'blob',
        compression: 'DEFLATE',
        compressionOptions: {
          level: 6,
        },
      })

      const zipUrl = URL.createObjectURL(zipContent)
      const a = document.createElement('a')
      a.href = zipUrl
      a.download = `selected_images_${new Date().getTime()}.zip`
      document.body.appendChild(a)
      a.click()

      setTimeout(() => {
        document.body.removeChild(a)
        URL.revokeObjectURL(zipUrl)
      }, 100)

      if (onSuccess) {
        onSuccess($t('pages.share.selection.successSummary', { count: selectedImagesData.length.toString() }))
      }
    } catch {
      if (onError) {
        onError($t('pages.share.selection.packFailedRetry'))
      }
    }
  }

  const downloadSingleFile = (url: string, filename: string) => {
    if (!url) {
      return
    }

    fetch(url)
      .then((response) => response.blob())
      .then((blob) => {
        const blobUrl = URL.createObjectURL(blob)

        const a = document.createElement('a')
        a.href = blobUrl
        a.download = filename || 'image'
        a.style.display = 'none'

        document.body.appendChild(a)
        a.click()

        setTimeout(() => {
          document.body.removeChild(a)
          URL.revokeObjectURL(blobUrl)
        }, 100)

        if (onSuccess) {
          onSuccess($t('pages.share.selection.downloadSuccess'))
        }
      })
      .catch(() => {
        if (onError) {
          onError($t('pages.share.selection.downloadFailed'))
        }

        try {
          const a = document.createElement('a')
          a.href = url
          a.download = filename || 'image'
          a.target = '_blank' // 确保在新窗口打开，避免导航离开当前页面
          a.rel = 'noopener noreferrer'
          document.body.appendChild(a)
          a.click()
          document.body.removeChild(a)
        } catch {}
      })
  }

  return {
    isSelectionMode,
    selectedImages,
    toggleSelectionMode,
    toggleImageSelection,
    isImageSelected,
    selectAll,
    deselectAll,
    getSelectedImagesData,
    downloadSelected,
  }
}
