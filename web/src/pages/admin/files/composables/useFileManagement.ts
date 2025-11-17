import { reactive, ref } from 'vue'
import { useToast } from '@/components/Toast/useToast'
import {
  type ImageInfo,
  type ImageListParams,
  batchDeleteAdminImages,
  batchSetImageRecommendation,
  deleteAdminFile,
  getAdminFileList,
  setFileRecommendation,
} from '@/api/admin/files'
import { processParams } from '@/utils/formatting/format'
import { useTexts } from '@/composables/useTexts'

export function useFileManagement() {
  const toast = useToast()
  const { $t } = useTexts()

  const images = ref<ImageInfo[]>([])
  const totalImages = ref(0)
  const loading = ref(false)

  const filters = reactive<ImageListParams>({
    page: 1,
    page_size: 24,
    keyword: '',
    uploader_id: null,
    channel: '',
    is_recommended: null,
    start_date: '',
    end_date: '',
    sort: 'created_at',
    order: 'desc',
  })

  const loadImages = async () => {
    try {
      loading.value = true
      const params = processParams(filters)
      const response = await getAdminFileList(params)

      images.value = response.images || []
      totalImages.value = response.total || 0
    } catch (_error) {
    } finally {
      loading.value = false
    }
  }

  const deleteImage = async (image: ImageInfo) => {
    try {
      await deleteAdminFile(image.id)
      toast.success($t('admin.files.management.deleteSuccess'))
      await loadImages()
    } catch (_error) {}
  }

  const batchDelete = async (fileIds: number[]) => {
    try {
      await batchDeleteAdminImages(fileIds)
      toast.success($t('admin.files.management.batchDeleteSuccess').replace('{count}', fileIds.length.toString()))
      await loadImages()
    } catch (_error) {}
  }

  const toggleRecommendation = async (image: ImageInfo) => {
    try {
      const newStatus = !image.is_recommended
      await setFileRecommendation(image.id, newStatus)

      const index = images.value.findIndex((img) => img.id === image.id)
      if (index !== -1) {
        images.value[index].is_recommended = newStatus
      }

      toast.success(newStatus ? $t('admin.files.management.setRecommended') : $t('admin.files.management.cancelRecommended'))
    } catch (_error) {}
  }

  const batchSetRecommendation = async (fileIds: number[], isRecommended: boolean) => {
    try {
      await batchSetImageRecommendation(fileIds, isRecommended)

      images.value.forEach((image) => {
        if (fileIds.includes(image.id)) {
          image.is_recommended = isRecommended
        }
      })

      const action = isRecommended ? $t('admin.files.actions.recommend') : $t('admin.files.actions.unrecommend')
      toast.success(
        $t('admin.files.management.batchOperationSuccess')
          .replace('{action}', action)
          .replace('{count}', fileIds.length.toString())
      )
    } catch (_error) {}
  }

  return {
    images,
    totalImages,
    loading,
    filters,

    loadImages,
    deleteImage,
    batchDelete,
    toggleRecommendation,
    batchSetRecommendation,
  }
}
