<script setup lang="ts">
  import { computed, onMounted, onUnmounted, ref, watch } from 'vue'
  import { useRoute, useRouter } from 'vue-router'
  import { useBreadcrumbStore } from '@/store/breadcrumb'
  import CreateFolderDialog from './components/CreateFolderDialog.vue'
  import CreateShareDialog from './components/CreateShareDialog.vue'
  import FolderContentView from './components/FolderContentView/index.vue'
  import PageHeader from './components/PageHeader.vue'
  import ShareSelectTip from './components/ShareSelectTip.vue'
  import DeleteConfirmDialogs from './components/DeleteConfirmDialogs.vue'
  import { useKeyboardEvents } from './composables/useKeyboardEvents'
  import {
    useDragSort,
    useFileNavigation,
    useFolderManagement,
    useFileManagement,
    useSelectionMode,
    useViewMode,
  } from '@/hooks/pages/folders'

  defineOptions({
    name: 'FoldersPage',
  })

  const route = useRoute()
  const router = useRouter()
  const breadcrumbStore = useBreadcrumbStore()

  const viewMode = useViewMode()
  const navigation = useFileNavigation()
  const selection = useSelectionMode()
  const dragSort = useDragSort()
  const folderManagement = useFolderManagement()
  const imageManagement = useFileManagement()

  const shareDialogVisible = ref(false)
  const pageHeaderRef = ref()
  const shareButton = computed(() => pageHeaderRef.value?.shareButtonRef)

  const keyboardEvents = useKeyboardEvents(selection.selectMode, shareButton)
  const { setupKeyboardEvents, cleanupKeyboardEvents } = keyboardEvents

  /* 跳转到上传页面并携带当前文件夹ID */
  const navigateToUpload = () => {
    const currentFolderId = navigation.currentFolderId.value

    if (currentFolderId) {
      router.push({
        path: '/upload',
        query: { folderId: currentFolderId },
      })
    } else {
      router.push('/upload')
    }
  }

  const handleToggleSelectAllImages = () => {
    if (selection.selectMode.value) {
      shareButton.value?.toggleSelectAllImages?.(folderManagement.images.value)
    } else {
      selection.toggleSelectAllFiles(folderManagement.images.value)
    }
  }

  const handleInvertImageSelection = () => {
    if (selection.selectMode.value) {
      shareButton.value?.invertImageSelection?.(folderManagement.images.value)
    } else {
      selection.invertFileSelection(folderManagement.images.value)
    }
  }

  const handleShareCreated = () => {
    shareButton.value?.cancelSelectMode()
    shareDialogVisible.value = false
  }

  const handleSelectModeChange = (value: boolean) => {
    selection.selectMode.value = value
  }

  const confirmDelete = async () => {
    if (!folderManagement.deleteTarget.value) {
      return
    }

    if (folderManagement.deleteType.value === 'folder') {
      await folderManagement.confirmDelete(navigation.currentFolderId.value)
    } else {
      folderManagement.isDeleting.value = true
      try {
        await imageManagement.confirmDeleteImage(folderManagement.deleteTarget.value, () =>
          folderManagement.loadFolders(navigation.currentFolderId.value)
        )
      } finally {
        folderManagement.isDeleting.value = false
      }
      folderManagement.deleteDialogVisible.value = false
      folderManagement.deleteTarget.value = null
    }
  }

  const executeBatchDelete = async () => {
    await imageManagement.executeBatchDelete(
      selection.selectedBatchFiles.value,
      () => folderManagement.loadFolders(navigation.currentFolderId.value),
      selection.cancelBatchMode,
      (processing: boolean) => (selection.isBatchProcessing.value = processing)
    )
  }

  watch(
    () => navigation.breadcrumbItems.value,
    (newItems) => {
      breadcrumbStore.setItems(
        newItems.map((item) => ({
          id: item.id,
          name: item.name,
          path: item.id ? `/folders/${item.id}` : '/folders',
        }))
      )
    },
    { immediate: true, deep: true }
  )

  watch(
    () => route.params.folderPath,
    async (newPath, oldPath) => {
      if (newPath !== oldPath) {
        const pathString = Array.isArray(newPath) ? newPath.join('/') : newPath || ''
        await navigation.buildBreadcrumbFromPath(pathString)
        folderManagement.loadFolders(navigation.currentFolderId.value)
      }
    },
    { immediate: false }
  )

  onMounted(async () => {
    const { folderPath } = route.params
    const pathString = Array.isArray(folderPath) ? folderPath.join('/') : folderPath || ''

    if (pathString) {
      await navigation.buildBreadcrumbFromPath(pathString)
    }

    folderManagement.loadFolders(navigation.currentFolderId.value)
    setupKeyboardEvents()
  })

  onUnmounted(() => {
    cleanupKeyboardEvents()
    dragSort.cleanup()
  })
</script>

<template>
  <div class="folders-page">
    <CyberFileViewer
      v-model="imageManagement.previewDialogVisible.value"
      :file="imageManagement.previewImage.value"
      @close="imageManagement.handlePreviewClose"
    />

    <PageHeader
      ref="pageHeaderRef"
      :breadcrumb-items="navigation.breadcrumbItems.value"
      :folders="folderManagement.folders.value"
      :images="folderManagement.images.value"
      :current-folder-id="navigation.currentFolderId.value || null"
      :select-mode="selection.selectMode.value"
      @breadcrumb-click="navigation.handleBreadcrumbClick"
      @navigate-to-upload="navigateToUpload"
      @show-create-dialog="folderManagement.showCreateDialog"
      @update:select-mode="handleSelectModeChange"
      @update:share-dialog-visible="shareDialogVisible = $event"
    />

    <ShareSelectTip :visible="selection.selectMode.value" />

    <div class="content-area">
      <FolderContentView
        :folders="folderManagement.folders.value"
        :images="folderManagement.images.value"
        :select-mode="selection.selectMode.value"
        :batch-mode="selection.batchMode.value"
        :view-mode="viewMode.viewMode.value"
        :selected-batch-images="selection.selectedBatchFiles.value"
        :is-all-folders-selected="selection.isAllFoldersSelected.value(folderManagement.folders.value)"
        :is-all-images-selected="selection.isAllFilesSelected.value(folderManagement.images.value)"
        :is-folder-selected="shareButton?.isFolderSelected"
        :is-image-selected="shareButton?.isImageSelected"
        :is-batch-image-selected="selection.isBatchFileSelected"
        @folder-click="navigation.navigateToFolder"
        @folder-moved="() => folderManagement.loadFolders(navigation.currentFolderId.value)"
        @folder-double-click="navigation.navigateToFolder"
        @start-batch-mode="selection.startBatchMode"
        @cancel-batch-mode="selection.cancelBatchMode"
        @batch-delete="imageManagement.confirmBatchDeleteImages(selection.selectedBatchFiles.value)"
        @toggle-select-all-images="handleToggleSelectAllImages"
        @toggle-select-all-folders="() => shareButton?.toggleSelectAllFolders?.(folderManagement.folders.value)"
        @invert-folder-selection="() => shareButton?.invertFolderSelection?.(folderManagement.folders.value)"
        @invert-image-selection="handleInvertImageSelection"
        @toggle-batch-image-select="selection.toggleBatchFileSelect"
        @toggle-image-select="(image) => shareButton?.toggleImageSelect(image)"
        @toggle-folder-select="(folder) => shareButton?.toggleFolderSelect(folder)"
        @change-view-mode="viewMode.toggleViewMode"
        @preview-image="(image) => imageManagement.handlePreview(image, folderManagement.images.value)"
        @view-details="imageManagement.handleViewDetails"
        @copy-link="imageManagement.handleCopyLink"
        @download-image="imageManagement.handleDownload"
        @delete-image="(image) => imageManagement.handleDeleteImage(image, folderManagement.showDeleteDialog)"
        @delete-folder="(folder) => folderManagement.handleDeleteFolder(folder)"
        @edit-folder="(folder) => folderManagement.handleEditFolder(folder)"
        @toggle-image-visibility="
          (image) =>
            imageManagement.handleToggleVisibility(image, () => folderManagement.loadFolders(navigation.currentFolderId.value))
        "
        @toggle-folder-visibility="
          (folder, event) => folderManagement.handleToggleFolderVisibility(folder, event, navigation.currentFolderId.value)
        "
        @create-folder="folderManagement.showCreateDialog"
        @folder-drag-start="dragSort.onFolderDragStart"
        @folder-drag-end="
          (data) =>
            dragSort.onFolderDragEnd(data.event, data.sortedFolders, navigation.currentFolderId.value, () =>
              folderManagement.loadFolders(navigation.currentFolderId.value)
            )
        "
        @image-drag-start="dragSort.onFileDragStart"
        @image-drag-end="
          (data) =>
            dragSort.onFileDragEnd(data.event, data.sortedFiles, navigation.currentFolderId.value, () =>
              folderManagement.loadFolders(navigation.currentFolderId.value)
            )
        "
        @image-moved="
          (fileId, targetFolderId) => {
            folderManagement.loadFolders(navigation.currentFolderId.value)
          }
        "
        @image-deleted="
          (fileId) => {
            folderManagement.loadFolders(navigation.currentFolderId.value)
          }
        "
      />
    </div>

    <CreateFolderDialog
      v-model="folderManagement.createDialogVisible.value"
      :parent-id="navigation.currentFolderId.value"
      :mode="folderManagement.dialogMode.value"
      :folder="folderManagement.currentFolder.value"
      @created="(folder) => folderManagement.handleFolderCreated(folder, navigation.currentFolderId.value)"
      @updated="(folder) => folderManagement.handleFolderUpdated(folder, navigation.currentFolderId.value)"
    />

    <CreateShareDialog
      v-model="shareDialogVisible"
      :selected-folders="shareButton?.selectedFolders || []"
      :selected-images="shareButton?.selectedImages || []"
      @created="handleShareCreated"
    />

    <DeleteConfirmDialogs
      :delete-dialog-visible="folderManagement.deleteDialogVisible.value"
      :delete-type="folderManagement.deleteType.value"
      :delete-target="folderManagement.deleteTarget.value"
      :is-deleting="folderManagement.isDeleting.value"
      :batch-delete-dialog-visible="imageManagement.batchDeleteDialogVisible.value"
      :selected-batch-images-count="selection.selectedBatchFiles.value.length"
      :is-batch-deleting="imageManagement.isBatchDeleting.value"
      @update:delete-dialog-visible="folderManagement.deleteDialogVisible.value = $event"
      @update:batch-delete-dialog-visible="imageManagement.batchDeleteDialogVisible.value = $event"
      @confirm-delete="confirmDelete"
      @confirm-batch-delete="executeBatchDelete"
      @cancel-delete="folderManagement.cancelDelete"
      @cancel-batch-delete="imageManagement.cancelBatchDelete"
    />

    <CyberFileDetailModal
      :model-value="imageManagement.detailDialogVisible.value"
      :fileInfo="imageManagement.detailImage.value"
      :is-admin="false"
      @update:model-value="imageManagement.detailDialogVisible.value = $event"
      @download="imageManagement.handleDownload(imageManagement.detailImage.value)"
    />
  </div>
</template>

<style scoped>
  .folders-page {
    padding: 0;
    display: flex;
    flex-direction: column;
    position: relative;
  }

  .content-area {
    display: flex;
    flex-direction: column;
    gap: 1.5rem;
    flex: 1;
    padding-right: 0.5rem;
    margin-right: -0.5rem;
  }

  .content-area::-webkit-scrollbar {
    display: none;
  }

  .section-container {
    background: rgba(var(--color-background-800-rgb), 0.95);
    backdrop-filter: blur(10px);
    border: 1.5px solid rgba(var(--color-brand-500-rgb), 0.2);
    border-radius: var(--radius-sm);
    overflow: hidden;
    padding: 1.25rem;
    box-shadow:
      0 4px 12px var(--color-overlay-medium),
      0 2px 6px var(--color-overlay-light);
  }

  .section-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 1rem;
    padding-bottom: 0.6rem;
    border-bottom: 1.5px solid rgba(var(--color-brand-500-rgb), 0.25);
  }

  .section-title {
    margin: 0;
    font-size: 1.1rem;
    font-weight: 600;
    color: var(--color-content-heading);
    display: flex;
    align-items: center;
    gap: 0.4rem;
  }

  .section-title i {
    color: var(--color-brand-500);
    font-size: 1rem;
  }

  .item-count {
    background: rgba(var(--color-brand-500-rgb), 0.2);
    color: var(--color-brand-500);
    padding: 0.15rem 0.5rem;
    border-radius: var(--radius-sm);
    font-size: 0.75rem;
    font-weight: 600;
    margin-left: 0.4rem;
  }

  .drag-tip {
    background: rgba(var(--color-warning-rgb), 0.1);
    color: var(--color-warning-500);
    padding: 0.1rem 0.4rem;
    border-radius: var(--radius-sm);
    font-size: 0.7rem;
    margin-left: 0.5rem;
    display: flex;
    align-items: center;
    gap: 0.2rem;
    opacity: 0.8;
    transition: all 0.2s ease;
  }

  .drag-tip:hover {
    opacity: 1;
    background: rgba(var(--color-warning-rgb), 0.15);
  }

  .drag-tip i {
    font-size: 0.6rem;
  }

  .section-actions {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    flex-wrap: wrap;
  }

  .action-group {
    display: flex;
    gap: 0.4rem;
    align-items: center;
  }

  .action-divider {
    width: 1.5px;
    height: 1.2rem;
    background: rgba(var(--color-brand-500-rgb), 0.3);
    margin: 0 0.2rem;
  }

  .empty-state {
    text-align: center;
    padding: 3rem 2rem;
    color: rgba(var(--color-content-muted-rgb), 0.6);
  }

  .empty-state i {
    font-size: 4rem;
    margin-bottom: 1rem;
    color: rgba(var(--color-brand-500-rgb), 0.3);
  }

  .empty-state h3 {
    margin: 0 0 0.5rem 0;
    font-size: 1.5rem;
    color: rgba(var(--color-content-heading-rgb), 0.8);
  }

  .empty-state p {
    margin: 0;
    font-size: 1rem;
  }

  @media (max-width: 768px) {
    .section-container {
      padding: 1rem;
    }

    .section-header {
      flex-direction: column;
      align-items: flex-start;
      gap: 0.75rem;
    }

    .section-actions {
      width: 100%;
      justify-content: flex-start;
    }

    .action-group {
      flex-wrap: wrap;
    }
  }
</style>
<style src="@/styles/selection.css"></style>
