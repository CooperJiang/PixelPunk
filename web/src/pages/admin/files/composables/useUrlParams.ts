import { useRoute, useRouter } from 'vue-router'
import type { ImageListParams } from '@/api/admin/files'
import { useToast } from '@/components/Toast/useToast'
import { useTexts } from '@/composables/useTexts'

export function useUrlParams(pageSizeRef: { value: number }, currentPageRef: { value: number }, filtersRef: ImageListParams) {
  const route = useRoute()
  const router = useRouter()
  const toast = useToast()
  const { $t } = useTexts()

  const syncUrlParams = (currentPage: number, pageSize: number, filters: ImageListParams) => {
    const query: Record<string, string | number> = {}

    if (route.query.tag_id) {
      query.tag_id = route.query.tag_id
    }
    if (route.query.tag_name) {
      query.tag_name = route.query.tag_name
    }

    if (currentPage > 1) {
      query.page = currentPage
    }
    if (pageSize !== 24) {
      query.size = pageSize
    }
    if (filters.sort !== 'newest') {
      query.sort = filters.sort
    }

    if (filters.keyword) {
      query.keyword = filters.keyword
    }
    if (filters.tags) {
      query.tags = filters.tags
    }
    if (filters.dominant_color) {
      query.dominant_color = filters.dominant_color
    }
    if (filters.resolution) {
      query.resolution = filters.resolution
    }
    if (filters.is_nsfw !== undefined) {
      query.is_nsfw = filters.is_nsfw
    }
    if (filters.is_recommended !== undefined) {
      query.is_recommended = filters.is_recommended
    }
    if (filters.storage_type) {
      query.storage_type = filters.storage_type
    }
    if (filters.min_width) {
      query.min_width = filters.min_width
    }
    if (filters.max_width) {
      query.max_width = filters.max_width
    }
    if (filters.min_height) {
      query.min_height = filters.min_height
    }
    if (filters.max_height) {
      query.max_height = filters.max_height
    }
    if (filters.user_id) {
      query.user_id = filters.user_id
    }

    router.replace({ query })
  }

  const clearUrlParams = () => {
    router.replace({ query: {} })
  }

  const restoreFromUrlParams = () => {
    const { query } = route

    if (query.size) {
      const newSize = Number(query.size)
      pageSizeRef.value = newSize
      filtersRef.size = newSize
    }

    if (query.page) {
      const pageNum = Number(query.page)
      currentPageRef.value = pageNum
      filtersRef.page = pageNum
    }

    if (query.keyword) {
      filtersRef.keyword = String(query.keyword)
    }
    if (query.tags) {
      filtersRef.tags = String(query.tags)
    }
    if (query.dominant_color) {
      filtersRef.dominant_color = String(query.dominant_color)
    }
    if (query.resolution) {
      filtersRef.resolution = String(query.resolution)
    }
    if (query.is_nsfw !== undefined) {
      filtersRef.is_nsfw = query.is_nsfw === 'true'
    }
    if (query.is_recommended !== undefined) {
      filtersRef.is_recommended = query.is_recommended === 'true'
    }
    if (query.storage_type) {
      filtersRef.storage_type = String(query.storage_type)
    }
    if (query.min_width) {
      filtersRef.min_width = Number(query.min_width)
    }
    if (query.max_width) {
      filtersRef.max_width = Number(query.max_width)
    }
    if (query.min_height) {
      filtersRef.min_height = Number(query.min_height)
    }
    if (query.max_height) {
      filtersRef.max_height = Number(query.max_height)
    }
    if (query.user_id) {
      filtersRef.user_id = Number(query.user_id)
    }
    if (query.sort) {
      filtersRef.sort = String(query.sort) as 'newest' | 'oldest' | 'largest' | 'smallest'
    }

    const tagId = query.tag_id
    const tagName = query.tag_name

    if (tagId) {
      filtersRef.tags = Array.isArray(tagId) ? tagId.join(',') : String(tagId)

      if (tagName) {
        const displayName = Array.isArray(tagName) ? tagName.join(', ') : tagName
        toast.info($t('admin.files.urlParams.tagFiltered').replace('{name}', displayName))
      } else {
        toast.info($t('admin.files.urlParams.tagIdFiltered').replace('{id}', filtersRef.tags))
      }
    }
  }

  const checkRouteParams = () => {
    const tagId = route.query.tag_id
    return Boolean(tagId)
  }

  return {
    syncUrlParams,
    clearUrlParams,
    restoreFromUrlParams,
    checkRouteParams,
  }
}
