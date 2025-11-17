<script setup lang="ts">
  import { computed, nextTick, onBeforeUnmount, onMounted, ref } from 'vue'
  import { useRoute, useRouter } from 'vue-router'
  import { useToast } from '@/components/Toast/useToast'
  import { useDevice } from '@/hooks/useDevice'
  import { useImagePreload } from '@/composables/useFilePreload'
  import { estimateImageHeight } from '@/utils/file/fileUtils'
  import { useTexts } from '@/composables/useTexts'

  import { INITIAL_LOAD_COUNT, getLayoutOptions, PAGE_SIZE } from '@/constants/share'

  const { $t } = useTexts()

  const LAYOUT_OPTIONS = computed(() => getLayoutOptions($t))

  import PasswordDialog from './components/PasswordDialog.vue'
  import VisitorInfoDialog from './components/VisitorInfoDialog.vue'
  import ShareHeader from './components/ShareHeader.vue'
  import ErrorState from './components/ErrorState.vue'
  import MobileShareInfo from './components/MobileShareInfo.vue'
  import FolderSection from './components/FolderSection.vue'
  import FileSection from './components/FileSection.vue'
  import ShareEmptyState from './components/ShareEmptyState.vue'
  import ShareFooter from './components/ShareFooter.vue'
  import FullscreenControls from './components/FullscreenControls.vue'
  import LayoutSelector from './components/LayoutSelector.vue'
  import GridSizeSelector from './components/GridSizeSelector.vue'
  import FolderBreadcrumb from '@/pages/folders/components/FolderBreadcrumb.vue'

  import { useShareData } from '@/composables/useShareData'
  import { useShareNavigation } from '@/composables/useShareNavigation'
  import { useShareLayout } from '@/composables/useShareLayout'
  import { useShareSorting } from '@/composables/useShareSorting'
  import { useSharePagination } from '@/composables/useSharePagination'
  import { useSharePreview } from '@/composables/useSharePreview'
  import { useShareSelection } from '@/composables/useShareSelection'

  const route = useRoute()
  const router = useRouter()
  const toast = useToast()

  const { isMobile } = useDevice()

  const { preloadImage } = useImagePreload()

  const showError = (title: string, message: string, type: string = 'error') => {
    error.value = true
    errorTitle.value = title
    errorMessage.value = message
    errorType.value = type
    initialLoading.value = false
    folderLoading.value = false
  }

  const {
    initialLoading,
    folderLoading,
    error,
    errorTitle,
    errorMessage,
    errorType,
    shareData,
    shareKey,
    accessToken,
    shareName,
    passwordDialogVisible,
    visitorDialogVisible,
    folders,
    images,
    rootImages,
    loadShareData,
    handlePasswordVerify: doPasswordVerify,
    handleVisitorSubmit,
    handleVisitorDismiss,
  } = useShareData(route, showError)

  const {
    currentFolderId,
    breadcrumbItems,
    restoreCurrentScrollPosition,
    handleFolderClick: doFolderClick,
    handleBreadcrumbClick: doBreadcrumbClick,
  } = useShareNavigation()

  const handleFolderClick = (folder: Record<string, unknown>) => {
    doFolderClick(folder, () => {
      loadShareData(currentFolderId.value)
    })
  }

  const handleBreadcrumbClick = (item: Record<string, unknown> | null) => {
    if (item && item.id === null) {
      doBreadcrumbClick(null, () => {
        loadShareData(currentFolderId.value)
      })
      return
    }

    doBreadcrumbClick(item, () => {
      loadShareData(currentFolderId.value)
    })
  }

  const handlePasswordVerify = (password: string) => {
    doPasswordVerify(
      password,
      (message) => toast.success(message),
      (message) => toast.error(message)
    )
  }

  const { currentSort, sortDirection, handleSortChange, sortedList } = useShareSorting(images)

  const {
    currentPage,
    paginatedList,
    loadMoreImages,
    showLoadMoreButton,
    isLoading: paginationLoading,
    canLoadMore,
  } = useSharePagination(sortedList, preloadImage)

  const allImages = computed(() => {
    if (currentFolderId.value) {
      return sortedList.value
    }
    return [...sortedList.value, ...rootImages.value].filter(
      (img, index, self) => index === self.findIndex((t) => t.id === img.id)
    )
  })

  const { previewDialogVisible, previewImage, handleImageClick, handlePreviewClose, handleKeydown } = useSharePreview(
    allImages,
    () => {},
    preloadImage
  )

  const {
    layoutMode,
    gridSizeValue,
    showGridSizePopover,
    showLayoutPopover,
    waterfallColumns,
    isImmersiveMode,
    loadLayoutSettings,
    openGridSizePopover,
    closeGridSizePopover,
    openLayoutPopover,
    closeLayoutPopover,
    toggleImmersiveMode,
    cleanup: cleanupLayout,
  } = useShareLayout(() => {
    nextTick(() => {
      restoreCurrentScrollPosition()
    })
  })

  const { isSelectionMode, selectedImages, toggleSelectionMode, toggleImageSelection, selectAll, deselectAll, downloadSelected } =
    useShareSelection(
      allImages,
      computed(() => accessToken.value),
      (message) => toast.error(message),
      (message) => toast.success(message),
      (message) => toast.info(message)
    )

  const fullscreenRef = ref<{ toggleFullScreen: () => void } | null>(null)

  const handleToggleFullscreen = () => {
    if (fullscreenRef.value) {
      fullscreenRef.value.toggleFullScreen()
    } else {
      toggleImmersiveMode()
    }
  }

  const getImageHeight = (item: Record<string, unknown>) => estimateImageHeight(item)

  const goHome = () => {
    router.push('/')
  }

  onMounted(() => {
    window.addEventListener('keydown', handleKeydown)

    loadLayoutSettings()

    document.addEventListener('click', closePopoversByOutsideClick)
  })

  onBeforeUnmount(() => {
    window.removeEventListener('keydown', handleKeydown)

    cleanupLayout()

    document.removeEventListener('click', closePopoversByOutsideClick)
  })

  const pageSize = PAGE_SIZE
  const initialLoadCount = INITIAL_LOAD_COUNT
  const layoutOptions = LAYOUT_OPTIONS

  const closePopoversByOutsideClick = (e: MouseEvent) => {
    if (showLayoutPopover.value || showGridSizePopover.value) {
      const layoutPopover = document.querySelector('.layout-popover')
      const gridSizePopover = document.querySelector('.grid-size-popover')

      const isClickInside = e.composedPath().some((el) => el === layoutPopover || el === gridSizePopover)

      if (!isClickInside) {
        closeLayoutPopover()
        closeGridSizePopover()
      }
    }
  }
</script>

<template>
  <div class="share-page" :class="{ 'is-mobile': isMobile }">
    <CyberLoading :visible="initialLoading" :text="$t('status.loading')" />

    <CyberMobileNavigation
      v-if="isMobile && !initialLoading && !error && shareData"
      :items="breadcrumbItems.length > 0 ? breadcrumbItems : []"
      :title="shareData?.share?.name"
      :has-password="shareData?.share?.has_password"
      :show-navigation="breadcrumbItems.length > 0"
      class="mobile-navigation"
      @click="handleBreadcrumbClick"
    />

    <cyberFileViewer v-model="previewDialogVisible" :file="previewImage" @close="handlePreviewClose" />

    <PasswordDialog v-model="passwordDialogVisible" :share-name="shareName" @verify="handlePasswordVerify" />

    <VisitorInfoDialog
      v-model="visitorDialogVisible"
      :share-key="shareKey"
      @submit="handleVisitorSubmit"
      @dismiss="handleVisitorDismiss"
    />

    <template v-if="!initialLoading && shareData">
      <template v-if="!isMobile">
        <ShareHeader :share="shareData.share" :user="shareData.user" />

        <div class="breadcrumb-container">
          <FolderBreadcrumb :items="breadcrumbItems" :show-icon="true" :auto-add-root="true" @click="handleBreadcrumbClick" />
        </div>
      </template>

      <FullscreenControls ref="fullscreenRef" :enabled="isImmersiveMode" @update:enabled="(value) => (isImmersiveMode = value)" />

      <div class="share-content">
        <MobileShareInfo
          v-if="isMobile && shareData"
          :show-back-button="breadcrumbItems.length > 0"
          :view-count="shareData.share.current_views"
          :created-at="shareData.share.created_at"
          @back-click="handleBreadcrumbClick(null)"
        />

        <FolderSection :folders="folders" @folder-click="handleFolderClick" />

        <FileSection
          v-if="images && images.length > 0"
          :images="images"
          :display-images="paginatedList"
          :is-root-images="!currentFolderId"
          :layout-mode="layoutMode"
          :is-selection-mode="isSelectionMode"
          :selected-ids="Array.from(selectedImages)"
          :initial-load-count="initialLoadCount"
          :current-sort="currentSort"
          :sort-direction="sortDirection"
          :waterfall-columns="waterfallColumns"
          :waterfall-column-width="Number(gridSizeValue)"
          :show-load-more-button="showLoadMoreButton"
          :current-page="currentPage"
          :page-size="pageSize"
          :get-image-height="getImageHeight"
          :is-full-screen="isImmersiveMode"
          :is-loading="paginationLoading"
          :can-load-more="canLoadMore"
          :share-key="shareKey"
          :access-token="accessToken"
          @image-click="handleImageClick"
          @toggle-selection="toggleImageSelection"
          @toggle-selection-mode="toggleSelectionMode"
          @select-all="selectAll"
          @deselect-all="deselectAll"
          @download-selected="downloadSelected"
          @sort-change="handleSortChange"
          @open-grid-size="openGridSizePopover"
          @toggle-immersive="toggleImmersiveMode"
          @open-layout="openLayoutPopover"
          @load-more="loadMoreImages"
          @toggle-fullscreen="handleToggleFullscreen"
        />

        <ShareEmptyState
          v-if="(!folders || folders.length === 0) && (!images || images.length === 0) && !initialLoading && !folderLoading"
        />

        <ShareFooter />

        <GridSizeSelector v-model="gridSizeValue" :visible="showGridSizePopover" @close="closeGridSizePopover" />

        <LayoutSelector
          v-model="layoutMode"
          :visible="showLayoutPopover"
          :layout-options="layoutOptions"
          @close="closeLayoutPopover"
        />
      </div>
    </template>

    <ErrorState
      v-if="error"
      :type="errorType"
      :title="errorTitle"
      :message="errorMessage"
      :action-text="$t('share.error.backHome')"
      @home-click="goHome"
    />
  </div>
</template>

<style scoped lang="scss">
  .share-page {
    min-height: 100vh;
    padding: var(--space-md);
    background: linear-gradient(
      to bottom,
      rgba(var(--color-background-900-rgb), 0.95),
      rgba(var(--color-background-900-rgb), 0.85)
    );

    &.is-mobile {
      padding: var(--space-sm);
    }
  }

  .share-content {
    width: 100%;
    max-width: 1920px;
    margin: 0 auto;
    position: relative;
    padding: 0 var(--space-md);
  }

  .breadcrumb-container {
    margin: var(--space-md) 0;
    padding: 0 var(--space-md);
  }

  @media (max-width: 768px) {
    .share-page {
      padding: var(--space-sm);
    }

    .share-content {
      padding: 0;
    }

    .breadcrumb-container {
      margin: var(--space-sm) 0;
      padding: 0;
      background: transparent;
      border: none;
    }
  }

  @media (min-width: 769px) and (max-width: 1024px) {
    .share-content {
      padding: 0 var(--space-sm);
    }
  }
</style>
