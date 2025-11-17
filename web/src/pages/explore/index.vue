<script setup lang="ts">
  import { computed, onMounted, onUnmounted, reactive, ref } from 'vue'
  import { useToast } from '@/components/Toast/useToast'
  import { useTexts } from '@/composables/useTexts'
  import type { FileInfo, FileListParams } from '@/api/types'
  import { useRoute } from 'vue-router'
  import { downloadFileQuick } from '@/utils/file/downloader'
  import { MAX_GALLERY_ITEMS, DEFAULT_PAGE_SIZE } from '@/constants/pagination'
  import { useLayoutStore } from '@/store/layout'
  import GalleryFilterSection from './components/GalleryFilterSection.vue'
  import GalleryMasonry from './components/GalleryMasonry.vue'
  import GalleryModals from './components/GalleryModals.vue'
  import SelectionToolbar from './components/SelectionToolbar.vue'
  import GalleryHeaderBar from './components/GalleryHeaderBar.vue'
  import GalleryFooterStatus from './components/GalleryFooterStatus.vue'
  import EmptyState from './components/EmptyState.vue'
  import CreateShareDialog from '@/pages/folders/components/CreateShareDialog.vue'
  import { useGalleryFilters } from './composables/useGalleryFilters'
  import { useGallerySearch } from './composables/useGallerySearch'
  import { useGalleryList } from './composables/useGalleryList'
  import { useGalleryPreview } from './composables/useGalleryPreview'
  import { useGallerySelection } from './composables/useGallerySelection'
  import { useDisplayDensity } from './composables/useDisplayDensity'
  import { SORT_OPTIONS, DENSITY_OPTIONS } from './constants'
  import type { SearchMode } from './types'

  defineOptions({
    name: 'ExplorePage',
  })

  const { $t } = useTexts()
  const toast = useToast()
  const route = useRoute()
  const layoutStore = useLayoutStore()

  const isLeftLayout = computed(() => layoutStore.isLeftLayout)

  const showFilter = ref(false)
  const images = ref<FileInfo[]>([])
  const searchKeyword = ref('')
  const searchMode = ref<SearchMode>('vector')
  const isVectorSearchMode = ref(false)
  const lastVectorSearchQuery = ref('')

  const { currentDensity, changeDensity, loadSavedDensity } = useDisplayDensity('compact')

  const filters = reactive<FileListParams>({
    page: 1,
    size: DEFAULT_PAGE_SIZE,
    sort: 'newest',
    dominant_color: '',
  })

  const {
    showDetailModal,
    showFilePreview,
    selectedImage,
    previewImageUrl,
    previewImageName,
    currentPreviewIndex,
    closeFilePreview,
    handlePrevImage,
    handleNextImage,
    viewImageDetails,
    viewImage,
  } = useGalleryPreview(images)

  const {
    selectMode,
    shareButton,
    isAllImagesSelected,
    showShareDialog,
    handleSelectModeChange,
    setShareButtonRef,
    isFileSelected,
    handleFileSelect,
    toggleSelectAllImages,
    invertImageSelection,
    handleShareCreated,
    attachEscListener,
    detachEscListener,
  } = useGallerySelection(images)

  const { currentPage, hasMore, hasReachedLimit, isLoading, isLoadingMore, loadImages, withPageLoading } = useGalleryList({
    images,
    filters,
    searchKeyword,
    isVectorSearchMode,
    pageSize: DEFAULT_PAGE_SIZE,
    maxImages: MAX_GALLERY_ITEMS,
  })

  const { getFilterCount, quickSort, applyFilter, resetFilter } = useGalleryFilters({
    filters,
    loadFiles: loadImages,
    defaultPageSize: DEFAULT_PAGE_SIZE,
    searchKeyword,
  })

  const { handleAdvancedSearch } = useGallerySearch({
    images,
    currentPage,
    hasMore,
    hasReachedLimit,
    withPageLoading,
    loadImages,
    searchKeyword,
    searchMode,
    isVectorSearchMode,
    lastVectorSearchQuery,
    toast,
  })
  const onFilter = (newFilters: FileListParams) => {
    applyFilter(newFilters)
    showFilter.value = false
  }
  const downloadFile = async (file: FileInfo) => {
    await downloadFileQuick(file.id, file.display_name || file.original_name || 'file')
  }
  const handleShowDetailModalChange = (value: boolean) => {
    showDetailModal.value = value
  }

  const handleShowFilePreviewChange = (value: boolean) => {
    showFilePreview.value = value
  }

  onMounted(() => {
    loadSavedDensity()

    const urlTags = route.query.tags
    const fromShare = route.query.from === 'share'

    if (urlTags && typeof urlTags === 'string') {
      filters.tags = urlTags
      searchKeyword.value = ''

      if (fromShare) {
        toast.info($t('explore.toast.showingTagResults', { tags: urlTags }))
      }

      loadImages(false)
    } else {
      loadImages(false)
    }

    attachEscListener()
  })

  onUnmounted(() => {
    detachEscListener()
  })
</script>

<template>
  <div class="gallery-page text-sm">
    <CyberLoading :visible="isLoading" :full-screen="true" :text="$t('explore.loadingText')" />

    <div class="gallery-container w-full pb-4">
      <GalleryHeaderBar
        :images="images"
        :has-content="images.length > 0"
        :set-share-button-ref="setShareButtonRef"
        :filters-sort="filters.sort"
        :sort-options="SORT_OPTIONS"
        :current-density="currentDensity"
        :density-options="DENSITY_OPTIONS"
        :search-keyword="searchKeyword"
        :search-mode="searchMode"
        :show-filter="showFilter"
        :filter-count="getFilterCount()"
        :is-searching="false"
        :is-left-layout="isLeftLayout"
        @update:select-mode="handleSelectModeChange"
        @update:share-dialog-visible="showShareDialog = $event"
        @sort="quickSort"
        @change-density="changeDensity"
        @update:search-keyword="searchKeyword = $event"
        @update:search-mode="searchMode = $event"
        @update:show-filter="showFilter = $event"
        @search="handleAdvancedSearch"
      />

      <GalleryFilterSection :show="showFilter" :initial-filters="filters" @filter="onFilter" />

      <div
        v-if="hasReachedLimit"
        :class="{
          'px-4 sm:px-6 lg:px-8 xl:px-12 2xl:px-16': !isLeftLayout,
          'px-0': isLeftLayout,
        }"
      >
        <div
          class="mb-3 px-4 py-2 text-center text-sm"
          style="
            background: rgba(var(--color-background-900-rgb), 0.5);
            backdrop-filter: blur(20px);
            border: 1px solid rgba(var(--color-brand-500-rgb), 0.1);
            border-radius: var(--radius-sm);
            color: var(--color-content-muted);
            box-shadow: 0 2px 8px var(--color-overlay-light);
          "
        >
          <i class="fas fa-info-circle mr-2" />
          {{ $t('explore.footer.narrowSearch') }}
        </div>
      </div>

      <div
        v-if="selectMode"
        :class="{
          'px-4 sm:px-6 lg:px-8 xl:px-12 2xl:px-16': !isLeftLayout,
          'px-0': isLeftLayout,
        }"
      >
        <div class="share-select-tip-container">
          <div class="share-select-tip">
            <i class="fas fa-hand-pointer mr-2" />
            {{ $t('explore.selection.shareHint') }}
          </div>
          <SelectionToolbar
            :is-all-selected="isAllImagesSelected"
            :selected-count="shareButton?.selectedFiles?.length || 0"
            @toggle-all="toggleSelectAllImages"
            @invert="invertImageSelection"
          />
        </div>
      </div>

      <div
        v-if="images.length > 0"
        class="gallery-images-fullscreen"
        :class="[
          `density-${currentDensity}`,
          {
            'px-4 sm:px-6 lg:px-8 xl:px-12 2xl:px-16': !isLeftLayout,
            'px-0': isLeftLayout,
          },
        ]"
      >
        <GalleryMasonry
          :files="images"
          :select-mode="selectMode"
          :is-file-selected="isFileSelected"
          :is-vector-search="isVectorSearchMode"
          :show-tags="true"
          @view="viewImage"
          @download="downloadFile"
          @details="viewImageDetails"
          @select="handleFileSelect"
        />

        <GalleryFooterStatus
          :files-length="images.length"
          :is-loading-more="isLoadingMore"
          :has-reached-limit="hasReachedLimit"
          :has-more="hasMore"
        />
      </div>

      <div v-if="!isLoading && images.length === 0" class="empty-state-fullscreen" :class="{ 'left-layout': isLeftLayout }">
        <EmptyState @reset="resetFilter" />
      </div>

      <GalleryModals
        :show-detail-modal="showDetailModal"
        :show-file-preview="showFilePreview"
        :selected-file="selectedImage"
        :preview-file-url="previewImageUrl"
        :preview-file-name="previewImageName"
        :files="images"
        :current-preview-index="currentPreviewIndex"
        @update:show-detail-modal="handleShowDetailModalChange"
        @update:show-file-preview="handleShowFilePreviewChange"
        @download="downloadFile"
        @preview-close="closeFilePreview"
        @prev="handlePrevImage"
        @next="handleNextImage"
      />

      <CreateShareDialog
        v-model="showShareDialog"
        :selected-folders="shareButton?.selectedFolders || []"
        :selected-images="shareButton?.selectedFiles || []"
        @created="handleShareCreated"
      />
    </div>
  </div>
</template>

<style src="./styles/gallery.css"></style>
<style src="@/styles/selection.css"></style>
