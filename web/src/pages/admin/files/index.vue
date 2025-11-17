<script setup lang="ts">
  import { computed, onBeforeUnmount, onMounted, reactive, ref, watch } from 'vue'
  import { useRoute, useRouter } from 'vue-router'
  import { useToast } from '@/components/Toast/useToast'
  import { useTexts } from '@/composables/useTexts'
  import {
    getAdminFileList,
    setFileRecommendation,
    deleteAdminFile,
    type ImageInfo,
    type ImageListParams,
  } from '@/api/admin/files'
  import { processParams } from '@/utils/formatting/format'
  import FileFilter from './components/FileFilterPanel.vue'
  import FileCard from './components/FileCard.vue'
  import AdminModals from './components/AdminModals.vue'
  import { useResponsivePageSize } from '@/hooks/useResponsivePageSize'
  import { useFileSelection } from './composables/useFileSelection'
  import { useUrlParams } from './composables/useUrlParams'
  import { useBatchOperations } from './composables/useBatchOperations'

  defineOptions({
    name: 'AdminFiles',
  })
  /* 工具 */
  const { $t } = useTexts()
  const toast = useToast()
  const route = useRoute()
  const router = useRouter()

  /* 状态变量 */
  const showFilter = ref(false)
  const isLoading = ref(false)
  const showSkeleton = ref(true) // 初始显示骨架屏
  const currentPage = ref(1)
  const pageSize = ref(24)
  /* 基于真实网格容器 `.images-grid` 自动推断列数，并在窗口/容器尺寸变化时自动更新 */
  const { pageSize: autoPageSize, recalc } = useResponsivePageSize({
    containerSelector: '.images-grid',
    gridSelector: '.images-grid',
    childSelector: '.image-card',
    itemMinWidth: 187,
    columnGap: 16,
    rowMultiple: 4,
    defaultSize: 24,
    debug: false,
    preferCssColumns: true,
    mode: 'once',
  })

  /* 基于自适应页大小的动态分页选项（1x/2x/3x/20x） */
  const pageSizeOptions = computed<number[]>(() => {
    const base = autoPageSize.value || pageSize.value || 24
    const list = [base, base * 2, base * 3, base * 20].map((n) => Math.max(1, Math.floor(n)))
    return Array.from(new Set(list)).sort((a, b) => a - b)
  })

  const totalFiles = ref(0)
  const files = ref<ImageInfo[]>([])
  const filters = reactive<ImageListParams>({
    page: 1,
    size: 24,
    sort: 'newest',
    tags: '',
  })

  const isInitialized = ref(false)

  /* 使用composables */
  const { selectMode, selectedImages, startSelectMode, cancelSelectMode, toggleImageSelection, createSelectionHelpers } =
    useFileSelection()

  const { syncUrlParams, clearUrlParams, restoreFromUrlParams } = useUrlParams(pageSize, currentPage, filters)

  const {
    isDeleting,
    isBatchProcessing,
    showBatchDeleteConfirm,
    batchSetRecommendation: batchSetRecommendationAction,
    confirmBatchDelete,
    batchDelete: batchDeleteAction,
    deleteImage: _deleteImageAction,
    toggleRecommendation: _toggleRecommendationAction,
  } = useBatchOperations()

  const {
    isAllSelected,
    hasCrossPageSelection: _hasCrossPageSelection,
    canInvertSelection,
    toggleSelectAll,
    selectInvert,
  } = createSelectionHelpers(() => files.value)

  /* 当前预览文件的索引 */
  const currentPreviewIndex = computed(() => {
    if (!selectedFile.value) {
      return 0
    }
    return files.value.findIndex((img) => img.id === selectedFile.value?.id) || 0
  })

  const handlePrevImage = () => {
    const currentIndex = currentPreviewIndex.value
    if (currentIndex > 0) {
      selectedFile.value = files.value[currentIndex - 1]
    }
  }

  const handleNextImage = () => {
    const currentIndex = currentPreviewIndex.value
    if (currentIndex < files.value.length - 1) {
      selectedFile.value = files.value[currentIndex + 1]
    }
  }

  const showDetailModal = ref(false)
  const showDeleteConfirm = ref(false)
  const showFilePreview = ref(false)
  const selectedFile = ref<ImageInfo | null>(null)
  const selectedFileId = ref<string | null>(null)

  const batchSetRecommendation = async (recommended: boolean) => {
    await batchSetRecommendationAction(selectedImages.value, recommended, files.value, () => {
      selectedImages.value = []
      selectMode.value = false
    })
  }

  const batchDelete = async () => {
    await batchDeleteAction(selectedImages.value, files.value, (deletedIds: string[]) => {
      files.value = files.value.filter((img) => !deletedIds.includes(img.id))

      if (files.value.length === 0 && currentPage.value > 1) {
        currentPage.value--
        syncUrlParams(currentPage.value, pageSize.value, filters)
        loadFiles()
      } else if (files.value.length === 0) {
        loadFiles()
      }

      selectedImages.value = []
      selectMode.value = false
    })
  }

  const handlePageSizeChange = (newSize: number) => {
    if (!isInitialized.value) {
      return
    }

    pageSize.value = newSize
    filters.size = newSize
    filters.page = 1 // 重置页码
    currentPage.value = 1
    selectedImages.value = [] // 切换每页数量时，清空选择
    syncUrlParams(currentPage.value, pageSize.value, filters)
    loadFiles()
  }

  const loadFiles = async () => {
    if (isLoading.value) {
      return
    }
    isLoading.value = true
    showSkeleton.value = true

    try {
      const params = processParams(
        {
          ...filters,
          page: currentPage.value,
          size: pageSize.value,
        },
        {
          booleanFields: ['is_recommended', 'is_nsfw'],
          removeUndefined: true,
          removeEmpty: true,
        }
      )

      const result = await getAdminFileList(params)

      if (result.success) {
        const { data } = result
        files.value = data?.items || []
        totalFiles.value = data?.pagination?.total || 0
      }
    } finally {
      isLoading.value = false
      showSkeleton.value = false
    }
  }

  const goToUpload = () => {
    router.push('/upload')
  }

  const handleFilter = (newFilters: ImageListParams) => {
    Object.assign(filters, newFilters)

    if (isInitialized.value) {
      filters.page = 1 // 重置页码
      currentPage.value = 1
    } else {
      filters.page = currentPage.value
    }

    selectedImages.value = [] // 筛选变更清空选择
    syncUrlParams(currentPage.value, pageSize.value, filters)
    loadFiles()
    showFilter.value = false
  }

  const resetFilter = () => {
    Object.assign(filters, {
      page: 1,
      size: pageSize.value,
      sort: 'newest',
      keyword: '',
      tags: '',
      dominant_color: '',
      resolution: '',
      nsfw_max_score: undefined,
      is_nsfw: undefined,
      is_recommended: undefined, // 改为undefined，避免设置默认值
      storage_type: '',
      min_width: undefined,
      max_width: undefined,
      min_height: undefined,
      max_height: undefined,
      user_id: undefined,
    })
    currentPage.value = 1
    selectedImages.value = [] // 清空选择
    syncUrlParams(currentPage.value, pageSize.value, filters)
    loadFiles()
  }

  const handlePageChange = (page: number) => {
    if (!isInitialized.value) {
      return
    }

    currentPage.value = page
    filters.page = page // 确保filters中的page也被更新
    syncUrlParams(currentPage.value, pageSize.value, filters)
    loadFiles()
  }

  const viewFileDetails = (id: string) => {
    const file = files.value.find((img) => img.id === id)
    if (file) {
      selectedFile.value = file
      showDetailModal.value = true
    }
  }

  const viewFile = (id: string) => {
    const file = files.value.find((img) => img.id === id)
    if (file) {
      selectedFile.value = file
      showFilePreview.value = true
    }
  }

  const confirmDeleteFile = (id?: string) => {
    selectedFileId.value = id || selectedFile.value?.id || null
    showDeleteConfirm.value = true
    showDetailModal.value = false
  }

  const deleteFile = async () => {
    if (!selectedFileId.value) {
      return
    }

    isDeleting.value = true

    try {
      const result = await deleteAdminFile(selectedFileId.value)

      if (!result.success) {
        toast.error($t('admin.files.messages.deleteFailed'))
        return
      }

      toast.success($t('admin.files.messages.deleteSuccess'))
      showDeleteConfirm.value = false

      files.value = files.value.filter((img) => img.id !== selectedFileId.value)

      selectedImages.value = selectedImages.value.filter((id) => id !== selectedFileId.value)

      if (files.value.length === 0 && currentPage.value > 1) {
        currentPage.value--
        filters.page = currentPage.value // 确保filters.page也同步更新
        syncUrlParams(currentPage.value, pageSize.value, filters)
        loadFiles()
      } else if (files.value.length === 0) {
        loadFiles()
      }
    } catch (_error) {
      toast.error($t('admin.files.messages.deleteError'))
    } finally {
      isDeleting.value = false
      selectedFileId.value = null
    }
  }

  const toggleRecommendation = async (id: string, _recommended: boolean) => {
    try {
      const result = await setFileRecommendation(id)

      if (!result.success) {
        toast.error($t('admin.files.messages.operationFailed'))
        return
      }

      const { data } = result

      toast.success(
        data.is_recommended ? $t('admin.files.messages.recommendSuccess') : $t('admin.files.messages.cancelRecommendSuccess')
      )

      const index = files.value.findIndex((img) => img.id === id)
      if (index !== -1) {
        files.value[index].is_recommended = data.is_recommended
      }

      if (selectedFile.value && selectedFile.value.id === id) {
        selectedFile.value.is_recommended = data.is_recommended
      }
    } catch (_error) {
      toast.error($t('admin.files.messages.operationError'))
    }
  }

  watch(
    [() => route.query, () => route.params],
    ([newQuery, newParams], [oldQuery, oldParams]) => {
      if (
        isInitialized.value &&
        (route.name === 'adminImages' || route.name === 'adminImagesByTags') &&
        (JSON.stringify(newQuery) !== JSON.stringify(oldQuery) || JSON.stringify(newParams) !== JSON.stringify(oldParams))
      ) {
        restoreFromUrlParams()
        loadImages()
      }
    },
    { immediate: false }
  )

  watch(
    () => route.name,
    (newName, oldName) => {
      if (
        (oldName === 'adminImages' || oldName === 'adminImagesByTags') &&
        newName !== 'adminImages' &&
        newName !== 'adminImagesByTags'
      ) {
        clearUrlParams()
        isInitialized.value = false
      }
    },
    { immediate: false }
  )

  onMounted(() => {
    restoreFromUrlParams()

    recalc()
    if (autoPageSize.value && autoPageSize.value !== pageSize.value) {
      pageSize.value = autoPageSize.value
      filters.size = autoPageSize.value
    }

    loadFiles()
    isInitialized.value = true
  })

  onBeforeUnmount(() => {
    clearUrlParams()
    isInitialized.value = false
  })
</script>

<template>
  <div class="admin-files-page admin-page-container admin-files">
    <CyberAdminWrapper
      :title="$t('admin.files.title')"
      :subtitle="$t('admin.files.subtitle')"
      icon="fas fa-file"
      :show-topbar="showFilter"
    >
      <template #actions>
        <CyberButton type="primary" @click="goToUpload">
          <i class="fas fa-upload mr-1.5" />{{ $t('admin.files.buttons.upload') }}
        </CyberButton>

        <CyberButton v-if="!selectMode" type="outlined" @click="startSelectMode">
          <i class="fas fa-tasks mr-1.5" />{{ $t('admin.files.buttons.batchOperation') }}
        </CyberButton>

        <template v-else>
          <div class="selection-info">
            <span class="selection-count">
              <i class="fas fa-file mr-1" />{{ $t('admin.files.selection.selected') }} {{ selectedImages.length }}
              {{ $t('admin.files.selection.selectedCount') }}
            </span>
          </div>

          <CyberButton
            type="success"
            :disabled="selectedImages.length === 0 || isBatchProcessing"
            @click="batchSetRecommendation(true)"
          >
            <i class="fas fa-star mr-1" />{{ $t('admin.files.batchActions.recommend') }}
          </CyberButton>
          <CyberButton
            type="warning"
            :disabled="selectedImages.length === 0 || isBatchProcessing"
            @click="batchSetRecommendation(false)"
          >
            <i class="far fa-star mr-1" />{{ $t('admin.files.batchActions.cancelRecommend') }}
          </CyberButton>
          <CyberButton
            type="danger"
            :disabled="selectedImages.length === 0 || isBatchProcessing"
            @click="confirmBatchDelete(selectedImages)"
          >
            <i class="fas fa-trash mr-1" />{{ $t('admin.files.batchActions.batchDelete') }}
          </CyberButton>
          <CyberButton type="outlined" @click="cancelSelectMode">
            <i class="fas fa-times mr-1" />{{ $t('admin.files.buttons.cancel') }}
          </CyberButton>
        </template>

        <CyberButton type="outlined" :class="{ 'border-error-500 text-error-500': showFilter }" @click="showFilter = !showFilter">
          <i class="fas fa-filter mr-1.5" />{{ $t('admin.files.buttons.filter') }}
        </CyberButton>
      </template>

      <template #topbar>
        <FileFilter :initial-filters="filters" @filter="handleFilter" />
      </template>

      <template #content>
        <div class="admin-content-wrapper">
          <div class="images-content">
            <div v-if="(files.length > 0 || isLoading) && selectMode" class="selection-toolbar">
              <button
                class="select-btn select-all-btn"
                :class="{ active: isAllSelected }"
                :title="isAllSelected ? $t('admin.files.selection.deselectAll') : $t('admin.files.selection.selectAll')"
                @click="toggleSelectAll"
              >
                <i class="fas" :class="isAllSelected ? 'fa-check-double' : 'fa-check'" />
                {{ isAllSelected ? $t('admin.files.selection.deselectAll') : $t('admin.files.selection.selectAll') }}
              </button>
              <button
                class="select-btn invert-btn"
                :disabled="!canInvertSelection"
                :title="
                  canInvertSelection ? $t('admin.files.selection.invertSelection') : $t('admin.files.selection.invertDisabledTip')
                "
                @click="selectInvert"
              >
                <i class="fas fa-exchange-alt" />
                {{ $t('admin.files.selection.invertSelection') }}
              </button>
            </div>

            <div class="images-grid-container">
              <div class="images-grid">
                <CyberSkeleton type="card" :count="24" :loading="showSkeleton" />
                <template v-if="!isLoading && files.length > 0">
                  <FileCard
                    v-for="file in files"
                    :key="file.id"
                    :file-data="file"
                    :select-mode="selectMode"
                    :selected="selectedImages.includes(file.id)"
                    class="image-card"
                    @view="viewFile"
                    @delete="confirmDeleteFile"
                    @recommend="toggleRecommendation"
                    @details="viewFileDetails"
                    @select="toggleImageSelection"
                  />
                </template>
              </div>
            </div>

            <div v-if="!showSkeleton && !isLoading && files.length > 0" class="pagination-wrapper">
              <CyberPagination
                v-model:current-page="currentPage"
                v-model:page-size="pageSize"
                :total="totalFiles"
                :page-size-options="pageSizeOptions"
                :show-page-size-selector="true"
                :show-quick-jumper="true"
                :show-total="true"
                @update:current-page="handlePageChange"
                @update:page-size="handlePageSizeChange"
              />
            </div>

            <div v-if="!showSkeleton && !isLoading && files.length === 0" class="empty-state">
              <div class="empty-content">
                <div class="empty-icon-wrapper">
                  <i class="fas fa-images empty-icon" />
                </div>
                <h3 class="empty-title">{{ $t('admin.files.empty.title') }}</h3>
                <p class="empty-desc">{{ $t('admin.files.empty.description') }}</p>
                <CyberButton type="secondary" @click="resetFilter">
                  <i class="fas fa-redo mr-2" />
                  {{ $t('admin.files.buttons.resetFilters') }}
                </CyberButton>
              </div>
            </div>
          </div>
        </div>
      </template>
    </CyberAdminWrapper>

    <cyberFileViewer
      v-model="showFilePreview"
      :file="selectedFile"
      :files="files"
      :initial-index="currentPreviewIndex"
      :show-side-nav="true"
      :show-keyboard-tips="true"
      @close="showFilePreview = false"
      @prev="handlePrevImage"
      @next="handleNextImage"
    />

    <AdminModals
      :show-detail-modal="showDetailModal"
      :show-delete-confirm="showDeleteConfirm"
      :show-batch-delete-confirm="showBatchDeleteConfirm"
      :selected-file="selectedFile"
      :selected-count="selectedImages.length"
      :is-deleting="isDeleting"
      :is-batch-processing="isBatchProcessing"
      @update:show-detail-modal="showDetailModal = $event"
      @update:show-delete-confirm="showDeleteConfirm = $event"
      @update:show-batch-delete-confirm="showBatchDeleteConfirm = $event"
      @delete="deleteFile"
      @batch-delete="batchDelete"
      @recommend="toggleRecommendation"
      @confirm-delete="confirmDeleteFile"
    />
  </div>
</template>

<style scoped>
  @import './styles/admin-files.css';

  .admin-content-wrapper {
    height: 100%;
    background: transparent;
    padding: 0;
    display: flex;
    flex-direction: column;
    min-height: 0;
  }

  .images-content {
    flex: 1;
    display: flex;
    flex-direction: column;
    min-height: 0;
  }

  .title-count {
    margin-left: var(--space-xs);
    font-size: var(--text-sm);
    color: var(--color-content-muted);
  }

  .pagination-wrapper {
    margin-top: var(--space-xl);
    padding-top: var(--space-md);
    border-top: 1px solid var(--color-border-subtle);
  }
</style>
<style src="@/styles/selection.css"></style>
