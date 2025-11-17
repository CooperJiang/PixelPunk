<script setup lang="ts">
  import { computed, nextTick, onMounted, reactive, ref, watch } from 'vue'
  import { useRoute, useRouter } from 'vue-router'
  import { downloadFileQuick } from '@/utils/file/downloader'
  import { useToast } from '@/components/Toast/useToast'
  import { getFileList, toggleAccessLevel } from '@/api/file'
  import { fileApi } from '@/api'
  import { userVectorSearch } from '@/api/search'
  import { getGlobalSettings } from '@/api/admin/settings/common'
  import { StorageUtil } from '@/utils/storage/storage'
  import { useLayoutStore } from '@/store/layout'
  import { TOKEN_KEY } from '@/constants'
  import { useResponsivePageSize } from '@/hooks/useResponsivePageSize'
  import { useTexts } from '@/composables/useTexts'

  import FileFilterPanel from './components/FileFilterPanel.vue'
  import FileCard from './components/FileCard.vue'
  import FileModals from './components/FileModals.vue'
  import ShareButton from '@/components/Share/ShareButton.vue'
  import CreateShareDialog from '@/pages/folders/components/CreateShareDialog.vue'
  import CompactSearchBar from '@/pages/explore/components/CompactSearchBar.vue'
  import EmptyState from '@/pages/explore/components/EmptyState.vue'

  interface ImageDetail {
    id: string
    url: string
    full_url: string
    full_thumb_url: string
    original_name: string
    display_name: string
    size: number
    width: number
    height: number
    format: string
    access_level: string
    created_at: string
    updated_at: string
    ai_info?: {
      description: string
      tags: string[]
      dominant_color: string
      resolution: string
      is_nsfw: boolean
      nsfw_score: number
      nsfw_evaluation: string
    }
    description?: string
  }

  interface Pagination {
    current_page: number
    last_page: number
    size: number
    total: number
  }

  const showFilter = ref(false)
  const route = useRoute()
  const router = useRouter()
  const layoutStore = useLayoutStore()
  const toast = useToast()
  const { $t } = useTexts()
  const isLoading = ref(false)
  const currentPage = ref(1)
  const showSkeleton = ref(false)
  const images = ref<ImageDetail[]>([])
  const pagination = ref<Pagination>({
    total: 0,
    size: 18,
    current_page: 1,
    last_page: 1,
  })

  const { pageSize: autoPageSize, recalc } = useResponsivePageSize({
    containerSelector: '.resource-container',
    gridSelector: '.images-grid',
    childSelector: undefined,
    itemMinWidth: 200,
    columnGap: 15,
    rowMultiple: 3,
    defaultSize: 18,
    preferCssColumns: true,
    debug: false,
    mode: 'once',
  })

  const pageSizeOptions = computed<number[]>(() => {
    const base = autoPageSize.value || pagination.value.size || 18
    const list = [base, base * 2, base * 3, base * 20].map((n) => Math.max(1, Math.floor(n)))
    return Array.from(new Set(list)).sort((a, b) => a - b)
  })

  watch(
    () => pagination.value.size,
    (newSize, oldSize) => {
      if (newSize !== oldSize) {
        currentPage.value = 1
        filterParams.size = newSize
        loadImages()
      }
    }
  )

  const searchKeyword = ref('')
  const searchMode = ref<'normal' | 'vector'>('vector')
  const isSearching = ref(false)
  const isVectorSearchMode = ref(false)
  const lastVectorSearchQuery = ref('')

  const filterParams = reactive({
    page: 1,
    size: 18,
    access_level: '',
    sort: 'newest',
    keyword: '',
    tags: '',
    categoryId: '',
    dominant_color: '',
    resolution: '',
    min_width: undefined,
    min_height: undefined,
  })

  const showFilePreview = ref(false)
  const showDetailModal = ref(false)
  const _selectedFileId = ref<string | null>(null)
  const selectedImageDetails = ref<ImageDetail | null>(null)

  const shareDialogVisible = ref(false)
  const shareButton = ref()
  const currentPreviewIndex = ref(0)

  const selectMode = ref(false)

  const handleSelectModeChange = (value: boolean) => {
    selectMode.value = value
  }

  const isImageSelected = (image: ImageDetail) => shareButton.value?.isImageSelected?.(image)

  const handleImageSelect = (image: ImageDetail) => {
    shareButton.value?.toggleImageSelect(image)
  }

  const loadImages = async () => {
    isLoading.value = true
    showSkeleton.value = true

    const params = {
      page: currentPage.value,
      size: pagination.value.size,
      access_level: filterParams.visibility === 'all' ? '' : filterParams.visibility,
      sort: filterParams.sort,
      keyword: filterParams.keyword,
      tags: filterParams.tags,
      categoryId: filterParams.categoryId,
      dominant_color: filterParams.dominant_color,
      resolution: filterParams.resolution,
      min_width: filterParams.min_width,
      min_height: filterParams.min_height,
    }

    try {
      const result = await getFileList(params)
      if (result.success) {
        const { data } = result
        const token = StorageUtil.get<string>(TOKEN_KEY)
        if (!data?.items) {
          return
        }

        images.value = data.items.map((image) => {
          let processedImage = image

          if (image.access_level === 'private' && token) {
            const separator = image.full_url.includes('?') ? '&' : '?'
            processedImage = {
              ...image,
              full_url: `${image.full_url}${separator}token=${token}`,
              full_thumb_url: `${image.full_thumb_url}${separator}token=${token}`,
            }
          }

          if (processedImage.ai_info?.description) {
            processedImage.description = processedImage.ai_info.description
          }

          return processedImage
        })

        pagination.value = data.pagination
      }
    } catch {
    } finally {
      isLoading.value = false
      showSkeleton.value = false
    }
  }

  const handleFilter = (newFilters) => {
    Object.assign(filterParams, newFilters)
    currentPage.value = 1

    if (isVectorSearchMode.value) {
      clearVectorSearch()
    } else {
      loadImages()
    }
    showFilter.value = false
  }

  const handleSearch = async (payload: { keyword: string; mode: 'normal' | 'vector' }) => {
    const { keyword, mode } = payload
    const query = keyword.trim()

    if (!query) {
      searchKeyword.value = ''
      searchMode.value = mode
      isVectorSearchMode.value = false
      resetFilter()
      return
    }

    searchKeyword.value = keyword
    searchMode.value = mode
    currentPage.value = 1

    if (mode === 'vector') {
      await handleVectorSearch(query)
    } else {
      await handleNormalSearch(query)
    }
  }

  const handleNormalSearch = async (query: string) => {
    isSearching.value = true
    isLoading.value = true
    showSkeleton.value = true
    isVectorSearchMode.value = false

    images.value = []

    filterParams.keyword = query
    await loadImages()

    isSearching.value = false
  }

  const handleVectorSearch = async (query: string) => {
    isSearching.value = true
    isLoading.value = true
    showSkeleton.value = true

    images.value = []

    try {
      const result = await userVectorSearch({
        query,
        page: currentPage.value,
      })

      if (result && result.data) {
        const { data } = result

        const token = StorageUtil.get<string>(TOKEN_KEY)

        if (data.items && Array.isArray(data.items) && data.items.length > 0) {
          images.value = data.items.map((item) => {
            const processedImage = { ...item }

            if (processedImage.access_level === 'private' && token) {
              const separator = processedImage.full_url?.includes('?') ? '&' : '?'
              if (processedImage.full_url) {
                processedImage.full_url = `${processedImage.full_url}${separator}token=${token}`
              }
              if (processedImage.full_thumb_url) {
                processedImage.full_thumb_url = `${processedImage.full_thumb_url}${separator}token=${token}`
              }
            }

            if (processedImage.ai_info?.description) {
              processedImage.description = processedImage.ai_info.description
            }

            if (item.similarity !== undefined) {
              processedImage.similarity = item.similarity
            }

            return processedImage
          })

          pagination.value = data.pagination
          isVectorSearchMode.value = true
          lastVectorSearchQuery.value = query

          toast.success($t('resource.toast.vectorSearchSuccess', { count: data.pagination.total }))
        } else {
          images.value = []
          pagination.value = {
            total: 0,
            size: pagination.value.size,
            current_page: 1,
            last_page: 1,
          }
          isVectorSearchMode.value = true
          lastVectorSearchQuery.value = query
          toast.info($t('resource.toast.vectorSearchEmpty'))
        }
      } else {
        images.value = []
        pagination.value = {
          total: 0,
          size: pagination.value.size,
          current_page: 1,
          last_page: 1,
        }
        toast.info($t('resource.toast.vectorSearchEmpty'))
      }
    } catch (_error) {
      images.value = []
      toast.error($t('resource.toast.vectorSearchFailed'))
    } finally {
      isSearching.value = false
      isLoading.value = false
      showSkeleton.value = false
    }
  }

  const handleVectorSearchPageChange = async (page: number) => {
    if (!isVectorSearchMode.value || !lastVectorSearchQuery.value) {
      return
    }

    currentPage.value = page
    isLoading.value = true
    showSkeleton.value = true

    try {
      const result = await userVectorSearch({
        query: lastVectorSearchQuery.value,
        page,
      })

      if (result && result.data && result.data.items) {
        const { data } = result
        const token = StorageUtil.get<string>(TOKEN_KEY)

        images.value = data.items.map((item) => {
          const processedImage = { ...item }

          if (processedImage.access_level === 'private' && token) {
            const separator = processedImage.full_url?.includes('?') ? '&' : '?'
            if (processedImage.full_url) {
              processedImage.full_url = `${processedImage.full_url}${separator}token=${token}`
            }
            if (processedImage.full_thumb_url) {
              processedImage.full_thumb_url = `${processedImage.full_thumb_url}${separator}token=${token}`
            }
          }

          if (processedImage.ai_info?.description) {
            processedImage.description = processedImage.ai_info.description
          }

          if (item.similarity !== undefined) {
            processedImage.similarity = item.similarity
          }

          return processedImage
        })

        pagination.value = data.pagination
      }
    } catch (_error) {
      toast.error($t('resource.toast.loadFailed'))
    } finally {
      isLoading.value = false
      showSkeleton.value = false
    }
  }

  const resetFilter = () => {
    Object.assign(filterParams, {
      page: 1,
      size: 18,
      visibility: '',
      sort: 'newest',
      keyword: '',
      tags: '',
      dominant_color: '',
      resolution: '',
      min_width: undefined,
      min_height: undefined,
    })
    currentPage.value = 1
    loadImages()
  }

  const handlePageChange = (page: number) => {
    currentPage.value = page
    if (isVectorSearchMode.value) {
      handleVectorSearchPageChange(page)
    } else {
      loadImages()
    }
  }

  const viewImage = (id: string) => {
    const imageIndex = images.value.findIndex((img) => img.id === id)
    if (imageIndex !== -1) {
      currentPreviewIndex.value = imageIndex
      selectedImageDetails.value = images.value[imageIndex]
      showFilePreview.value = true
    }
  }

  const handlePreviewChange = (index: number) => {
    if (index < 0 || index >= images.value.length) {
      return
    }

    currentPreviewIndex.value = index
    selectedImageDetails.value = images.value[index]
  }

  const handlePrevImage = () => {
    if (currentPreviewIndex.value > 0) {
      handlePreviewChange(currentPreviewIndex.value - 1)
    }
  }

  const handleNextImage = () => {
    if (currentPreviewIndex.value < images.value.length - 1) {
      handlePreviewChange(currentPreviewIndex.value + 1)
    }
  }

  const handleImageChange = (image: any, index: number) => {
    handlePreviewChange(index)
  }

  const handleShareCreated = () => {
    shareButton.value?.cancelSelectMode()
    shareDialogVisible.value = false
    toast.success($t('resource.toast.shareCreated'))
  }

  const handleCopyLink = async (image: ImageDetail) => {
    if (image.access_level === 'protected') {
      toast.warning($t('resource.toast.protectedFileWarning'))
      return
    }

    let url = image.full_url
    if (image.access_level === 'private') {
      const token = StorageUtil.get<string>(TOKEN_KEY)
      if (token) {
        url = `${url}?token=${token}`
      }
    }

    await navigator.clipboard.writeText(url)
    const message =
      image.access_level === 'public' ? $t('resource.toast.publicLinkCopied') : $t('resource.toast.privateLinkCopied')
    toast.success(message)
  }

  const handleDownload = async (fileId: string, image?: any) => {
    try {
      const fileName = image?.display_name || image?.original_name || `image_${fileId}`

      const success = await downloadFileQuick(fileId, fileName)

      if (success) {
        toast.success($t('resource.toast.downloadSuccess', { fileName }))
      } else {
        toast.error($t('resource.toast.downloadCancelled'))
      }
    } catch (_error) {
      toast.error($t('resource.toast.downloadFailed'))
    }
  }

  const deleteImage = async (id: string) => {
    isLoading.value = true

    try {
      const result = await fileApi.deleteFile(id)
      if (result.success) {
        toast.success($t('resource.toast.deleteSuccess'))
        loadImages()
      }
    } catch {
    } finally {
      isLoading.value = false
    }
  }

  const handleToggleVisibility = async (id: string, currentAccessLevel: string) => {
    const imageIndex = images.value.findIndex((img) => img.id === id)
    if (imageIndex === -1) {
      return
    }

    try {
      const result = await toggleAccessLevel(id)

      if (result.success) {
        let nextLevel = ''
        let nextAccessLevel = ''
        switch (currentAccessLevel) {
          case 'public':
            nextLevel = $t('resource.accessLevel.private')
            nextAccessLevel = 'private'
            break
          case 'private':
            nextLevel = $t('resource.accessLevel.protected')
            nextAccessLevel = 'protected'
            break
          case 'protected':
            nextLevel = $t('resource.accessLevel.public')
            nextAccessLevel = 'public'
            break
          default:
            nextLevel = $t('resource.accessLevel.public')
            nextAccessLevel = 'public'
        }

        images.value[imageIndex].access_level = nextAccessLevel

        if (nextAccessLevel === 'private') {
          const token = StorageUtil.get<string>(TOKEN_KEY)
          if (token) {
            const addTokenToUrl = (url: string) => {
              if (url && !url.includes('token=')) {
                const separator = url.includes('?') ? '&' : '?'
                return `${url}${separator}token=${token}`
              }
              return url
            }

            if (images.value[imageIndex].full_url) {
              images.value[imageIndex].full_url = addTokenToUrl(images.value[imageIndex].full_url)
            }
            if (images.value[imageIndex].full_thumb_url) {
              images.value[imageIndex].full_thumb_url = addTokenToUrl(images.value[imageIndex].full_thumb_url)
            }
          }
        } else if (currentAccessLevel === 'private' && nextAccessLevel !== 'private') {
          const removeTokenFromUrl = (url: string) => {
            if (url && url.includes('token=')) {
              return url.replace(/[?&]token=[^&]*/, '').replace(/\?$/, '')
            }
            return url
          }

          if (images.value[imageIndex].full_url) {
            images.value[imageIndex].full_url = removeTokenFromUrl(images.value[imageIndex].full_url)
          }
          if (images.value[imageIndex].full_thumb_url) {
            images.value[imageIndex].full_thumb_url = removeTokenFromUrl(images.value[imageIndex].full_thumb_url)
          }
        }

        toast.success($t('resource.toast.visibilityChanged', { level: nextLevel }))
      }
    } catch {
      toast.error($t('resource.toast.visibilityChangeFailed'))
    }
  }

  const viewImageDetails = (fileId: string) => {
    const image = images.value.find((img) => img.id === fileId)
    if (image) {
      selectedImageDetails.value = image
      showDetailModal.value = true
    }
  }

  const isAllImagesSelected = computed(() => {
    return shareButton.value?.selectedImages?.length === images.value.length && images.value.length > 0
  })

  const toggleSelectAllImages = () => {
    if (selectMode.value && shareButton.value) {
      shareButton.value.toggleSelectAllImages(images.value)
    }
  }

  const invertImageSelection = () => {
    if (selectMode.value && shareButton.value) {
      shareButton.value.invertImageSelection(images.value)
    }
  }

  const clearVectorSearch = () => {
    searchKeyword.value = ''
    searchMode.value = 'normal'
    isVectorSearchMode.value = false
    lastVectorSearchQuery.value = ''
    filterParams.keyword = ''
    currentPage.value = 1
    loadImages()
  }

  onMounted(async () => {
    isLoading.value = true
    showSkeleton.value = true

    try {
      const result = await getGlobalSettings()
      if (result.success && result.data) {
        const vectorConfig = result.data.vector as Record<string, unknown> | undefined
        if (vectorConfig && typeof vectorConfig.vector_enabled === 'boolean') {
          if (!vectorConfig.vector_enabled) {
            searchMode.value = 'normal'
          }
        }
      }
    } catch (_error) {
      console.error('Failed to fetch global settings:', error)
    }

    const tagsParam = route.query.tags
    const tagNamesParam = route.query.tag_names
    const categoryIdParam = route.query.categoryId
    const categoryNameParam = route.query.category_name

    if (tagsParam && typeof tagsParam === 'string') {
      filterParams.tags = tagsParam

      if (tagNamesParam && typeof tagNamesParam === 'string') {
        toast.info($t('resource.toast.tagFilter', { tagNames: tagNamesParam }))
      } else {
        toast.info($t('resource.toast.tagIdFilter', { tagIds: tagsParam }))
      }

      await router.replace({
        path: route.path,
        query: {},
      })
    }

    if (categoryIdParam && typeof categoryIdParam === 'string') {
      filterParams.categoryId = categoryIdParam

      if (categoryNameParam && typeof categoryNameParam === 'string') {
        toast.info($t('resource.toast.categoryFilter', { categoryName: categoryNameParam }))
      } else {
        toast.info($t('resource.toast.categoryIdFilter', { categoryId: categoryIdParam }))
      }

      await router.replace({
        path: route.path,
        query: {},
      })
    }

    await nextTick()
    setTimeout(() => {
      recalc(true)
      if (autoPageSize.value && autoPageSize.value !== pagination.value.size) {
        pagination.value.size = autoPageSize.value
        filterParams.size = autoPageSize.value
      }
      loadImages()
    }, 200)
  })
</script>

<template>
  <div
    :class="{
      'container mx-auto': layoutStore.isTopLayout,
    }"
  >
    <div class="resource-container">
      <div class="mb-3 flex items-center justify-between">
        <h1 class="page-title text-xl font-bold"><i class="fas fa-images mr-2 text-content" />{{ $t('resource.title') }}</h1>
        <div class="flex items-center space-x-2">
          <CompactSearchBar
            v-model="searchKeyword"
            :search-mode="searchMode"
            :show-filter="showFilter"
            :is-searching="isSearching"
            @update:search-keyword="searchKeyword = $event"
            @update:search-mode="searchMode = $event"
            @update:show-filter="showFilter = $event"
            @search="handleSearch"
          />

          <ShareButton
            ref="shareButton"
            :images="images"
            @update:select-mode="handleSelectModeChange"
            @update:share-dialog-visible="shareDialogVisible = $event"
          />
          <CyberButton
            type="secondary"
            class="text-xs"
            :class="{ 'border-error-500 text-error-500': showFilter }"
            @click="showFilter = !showFilter"
          >
            <i class="fas fa-filter" />
          </CyberButton>
        </div>
      </div>

      <transition
        name="filter-panel"
        enter-active-class="transition duration-300 ease-out"
        leave-active-class="transition duration-200 ease-in"
        enter-from-class="opacity-0 transform -translate-y-8"
        enter-to-class="opacity-100 transform translate-y-0"
        leave-from-class="opacity-100 transform translate-y-0"
        leave-to-class="opacity-0 transform -translate-y-8"
      >
        <FileFilterPanel v-show="showFilter" class="mb-4" @filter="handleFilter" />
      </transition>

      <div v-if="images.length > 0 || isLoading" class="mb-3 flex items-center justify-between px-1 text-xs">
        <div class="flex items-center gap-2">
          <div class="resource-count-text" v-html="$t('resource.count', { count: pagination.total })"></div>
          <div v-if="selectMode" class="flex items-center gap-2">
            <button
              class="select-btn select-all-btn"
              :class="{ active: isAllImagesSelected }"
              :title="isAllImagesSelected ? $t('resource.selection.deselectAll') : $t('resource.selection.selectAll')"
              @click="toggleSelectAllImages"
            >
              <i class="fas" :class="isAllImagesSelected ? 'fa-check-double' : 'fa-check'" />
              {{ isAllImagesSelected ? $t('resource.selection.deselectAll') : $t('resource.selection.selectAll') }}
            </button>
            <button class="select-btn invert-btn" :title="$t('resource.selection.invertTitle')" @click="invertImageSelection">
              <i class="fas fa-exchange-alt" />
              {{ $t('resource.selection.invert') }}
            </button>
          </div>
        </div>
      </div>

      <div v-if="!isLoading && images.length === 0" class="mb-4 flex items-center justify-center py-8">
        <EmptyState @reset="resetFilter" />
      </div>

      <div v-if="isLoading || images.length > 0" class="mb-4">
        <template v-if="showSkeleton">
          <div class="images-grid">
            <CyberSkeleton type="card" :count="18" :loading="showSkeleton" />
          </div>
        </template>

        <template v-else>
          <div class="images-grid relative">
            <template v-if="isLoading">
              <div
                class="pointer-events-none absolute inset-0 z-10 flex items-center justify-center bg-background-800 backdrop-blur-sm"
              >
                <div class="pointer-events-auto text-center">
                  <div class="mb-2">
                    <CyberSkeleton variant="circular" width="40px" height="40px" animation="wave" class="mx-auto" />
                  </div>
                  <CyberSkeleton variant="text" width="120px" height="14px" animation="wave" class="mx-auto" />
                </div>
              </div>
            </template>

            <FileCard
              v-for="image in images"
              :key="image.id"
              :image-data="image"
              :select-mode="selectMode"
              :is-selected="isImageSelected(image)"
              :is-vector-search="isVectorSearchMode"
              class="text-sm"
              @view="viewImage"
              @copy="handleCopyLink"
              @download="handleDownload"
              @delete="deleteImage"
              @visibility="handleToggleVisibility"
              @details="viewImageDetails"
              @select="handleImageSelect"
            />
          </div>
        </template>

        <div v-if="!isLoading && images.length > 0 && pagination.total > pagination.size" class="mt-6 flex justify-center">
          <CyberPagination
            v-model:page-size="pagination.size"
            v-model:current-page="currentPage"
            :total="pagination.total"
            :page-size-options="pageSizeOptions"
            :show-page-size-selector="true"
            class="text-xs"
            @update:current-page="handlePageChange"
            @update:page-size="
              (size) => {
                pagination.size = size
              }
            "
          />
        </div>
      </div>

      <cyberFileViewer
        v-model="showFilePreview"
        :file="selectedImageDetails"
        :files="images"
        :initial-index="currentPreviewIndex"
        :show-side-nav="true"
        :show-keyboard-tips="true"
        search-scope="user"
        @close="showFilePreview = false"
        @prev="handlePrevImage"
        @next="handleNextImage"
        @change="handleImageChange"
      />

      <FileModals
        :show-detail-modal="showDetailModal"
        :selected-image="selectedImageDetails"
        @update:show-detail-modal="showDetailModal = $event"
        @download="handleDownload"
      />

      <CreateShareDialog
        v-model="shareDialogVisible"
        :selected-folders="[]"
        :selected-images="shareButton?.selectedImages || []"
        @created="handleShareCreated"
      />
    </div>
  </div>
</template>

<style scoped>
  .resource-container {
    background: rgba(var(--color-background-800-rgb), 0.95);
    backdrop-filter: blur(10px);
    border: 1.5px solid rgba(var(--color-brand-500-rgb), 0.2);
    border-radius: var(--radius-sm);
    padding: 1.25rem;
    box-shadow:
      0 4px 12px var(--color-overlay-medium),
      0 2px 6px var(--color-overlay-light);
  }

  .page-title {
    color: var(--color-content-heading);
  }

  .resource-count-text {
    color: var(--color-content-muted);
  }

  .hover\:shadow-glow:hover {
    box-shadow: 0 0 12px rgba(var(--color-brand-500-rgb), 0.15);
  }

  .select-btn {
    display: flex;
    align-items: center;
    gap: 0.3rem;
    padding: 0.4rem 0.8rem;
    border-radius: var(--radius-sm);
    font-size: 0.85rem;
    cursor: pointer;
    transition: all 0.2s ease;
    border: 1px solid transparent;
    background: rgba(var(--color-brand-500-rgb), 0.1);
    color: var(--color-content-default);
  }

  .select-btn:hover {
    background: rgba(var(--color-brand-500-rgb), 0.2);
    transform: translateY(-1px);
  }

  .select-btn.active {
    background: rgba(var(--color-brand-500-rgb), 0.3);
    color: var(--color-brand-500);
    border-color: rgba(var(--color-brand-500-rgb), 0.4);
    box-shadow: 0 0 10px rgba(var(--color-brand-500-rgb), 0.2);
  }

  .select-btn i {
    font-size: 0.9rem;
  }

  .select-all-btn {
    background: linear-gradient(135deg, rgba(var(--color-brand-500-rgb), 0.15), rgba(var(--color-brand-500-rgb), 0.1));
  }

  .select-all-btn:hover {
    background: linear-gradient(135deg, rgba(var(--color-brand-500-rgb), 0.25), rgba(var(--color-brand-500-rgb), 0.2));
  }

  .select-all-btn.active {
    background: linear-gradient(135deg, rgba(var(--color-brand-500-rgb), 0.35), rgba(var(--color-brand-500-rgb), 0.3));
  }

  .invert-btn {
    background: linear-gradient(135deg, rgba(var(--color-error-rgb), 0.15), rgba(var(--color-error-rgb), 0.1));
    color: rgba(var(--color-error-rgb), 0.8);
  }

  .invert-btn:hover {
    background: linear-gradient(135deg, rgba(var(--color-error-rgb), 0.25), rgba(var(--color-error-rgb), 0.2));
    color: var(--color-error-500);
  }

  .invert-btn i {
    transform: rotate(90deg);
  }

  .images-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
    gap: 15px;
    perspective: 1000px;
  }
</style>
<style src="@/styles/selection.css"></style>
