import { computed, reactive, ref } from 'vue'
import { getFileList } from '@/api/file'
import { StorageUtil } from '@/utils/storage/storage'
import { TOKEN_KEY } from '@/constants'
import type { FilterParams, GalleryOptions, ImageDetail, Pagination } from './types'

/**
 * Gallery核心功能Hook - 文件列表管理
 * 提供文件加载、筛选、分页等通用功能
 */
export function useGalleryCore(options: GalleryOptions = {}) {
  const defaultOptions = {
    defaultPageSize: 18,
    enableInfiniteScroll: false,
    apiEndpoint: '/api/files',
    enableSelection: true,
    ...options,
  }

  const images = ref<ImageDetail[]>([])
  const pagination = ref<Pagination>({
    total: 0,
    size: defaultOptions.defaultPageSize,
    current_page: 1,
    last_page: 1,
  })

  const isLoading = ref(false)
  const currentPage = ref(1)
  const showFilter = ref(false)

  const filterParams = reactive<FilterParams>({
    page: 1,
    size: defaultOptions.defaultPageSize,
    access_level: '',
    sort: 'newest',
    keyword: '',
    tags: '',
    dominant_color: '',
    resolution: '',
    min_width: undefined,
    min_height: undefined,
  })

  const hasImages = computed(() => images.value.length > 0)
  const hasMorePages = computed(() => pagination.value.current_page < pagination.value.last_page)
  const totalImages = computed(() => pagination.value.total)

  const processImageUrl = (image: ImageDetail): ImageDetail => {
    const token = StorageUtil.get<string>(TOKEN_KEY)

    if ((image.access_level === 'private' || image.access_level === 'protected') && token) {
      const separator = image.full_url.includes('?') ? '&' : '?'
      return {
        ...image,
        full_url: `${image.full_url}${separator}token=${token}`,
        full_thumb_url: `${image.full_thumb_url}${separator}token=${token}`,
      }
    }

    return image
  }

  const processImageData = (imageList: ImageDetail[]): ImageDetail[] =>
    imageList.map((image) => {
      const processedImage = processImageUrl(image)

      if (processedImage.ai_info?.description) {
        processedImage.description = processedImage.ai_info.description
      }

      return processedImage
    })

  const loadImages = async (resetPage = false) => {
    if (resetPage) {
      currentPage.value = 1
      filterParams.page = 1
    }

    isLoading.value = true

    const params = {
      page: currentPage.value,
      size: pagination.value.size,
      access_level: filterParams.access_level === 'all' ? '' : filterParams.access_level,
      sort: filterParams.sort,
      keyword: filterParams.keyword,
      tags: filterParams.tags,
      dominant_color: filterParams.dominant_color,
      resolution: filterParams.resolution,
      min_width: filterParams.min_width,
      min_height: filterParams.min_height,
    }

    try {
      const result = await getFileList(params)
      if (result.success && result.data?.items) {
        const processedImages = processImageData(result.data.items)

        if (defaultOptions.enableInfiniteScroll && !resetPage) {
          images.value.push(...processedImages)
        } else {
          images.value = processedImages
        }

        pagination.value = result.data.pagination
      }
    } catch (error) {
      console.error('加载文件失败:', error)
    } finally {
      isLoading.value = false
    }
  }

  const handleFilter = async (newFilters: Partial<FilterParams>) => {
    Object.assign(filterParams, newFilters)
    await loadImages(true) // 重置到第一页
    showFilter.value = false
  }

  const resetFilter = async () => {
    Object.assign(filterParams, {
      page: 1,
      size: defaultOptions.defaultPageSize,
      access_level: '',
      sort: 'newest',
      keyword: '',
      tags: '',
      dominant_color: '',
      resolution: '',
      min_width: undefined,
      min_height: undefined,
    })
    await loadImages(true)
  }

  const handlePageChange = async (page: number) => {
    currentPage.value = page
    filterParams.page = page
    await loadImages()
  }

  const loadMore = async () => {
    if (!hasMorePages.value || isLoading.value) {
      return
    }

    currentPage.value++
    filterParams.page = currentPage.value
    await loadImages()
  }

  const refreshImages = async () => {
    await loadImages(true)
  }

  const toggleFilter = () => {
    showFilter.value = !showFilter.value
  }

  const findImageById = (id: string): ImageDetail | null => images.value.find((img) => img.id === id) || null

  const findImageIndex = (id: string): number => images.value.findIndex((img) => img.id === id)

  return {
    images,
    pagination,
    isLoading,
    currentPage,
    showFilter,
    filterParams,

    hasImages,
    hasMorePages,
    totalImages,

    loadImages,
    handleFilter,
    resetFilter,
    handlePageChange,
    loadMore,
    refreshImages,
    toggleFilter,
    findImageById,
    findImageIndex,

    processImageUrl,
    processImageData,
  }
}
