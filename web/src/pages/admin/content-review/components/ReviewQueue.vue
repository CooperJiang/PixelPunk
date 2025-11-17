<script setup lang="ts">
  import { nextTick, reactive, ref, watch } from 'vue'
  import { useToast } from '@/components/Toast/useToast'
  import { useTexts } from '@/composables/useTexts'
  import { useContentReview, useReviewSelection } from '../composables/useContentReview'
  import {
    batchReview,
    getReviewQueue,
    getReviewStats,
    reviewImage,
    type ReviewImage,
    type ReviewQueueQuery,
  } from '@/api/admin/content-review'
  import { getDefaultReviewReasons, getReviewMessages } from '@/constants'
  import { useResponsivePageSize } from '@/hooks/useResponsivePageSize'

  const props = withDefaults(
    defineProps<{
      isActive?: boolean
    }>(),
    {
      isActive: false,
    }
  )

  const emit = defineEmits<{
    (e: 'refresh-stats', stats: unknown): void
  }>()

  const toast = useToast()
  const { $t } = useTexts()

  /* 使用工厂函数生成i18n常量 */
  const REVIEW_MESSAGES = getReviewMessages($t)
  const DEFAULT_REVIEW_REASONS = getDefaultReviewReasons($t)
  const { formatDate, handleImageError } = useContentReview()
  const {
    selectMode,
    selectedItems: selectedImages,
    startSelectMode,
    exitSelectMode,
    toggleItemSelection: toggleImageSelection,
    toggleSelectAll,
    selectInvert,
    isAllSelected,
  } = useReviewSelection()

  const isLoading = ref(false)
  const showSkeleton = ref(false)
  const images = ref<ReviewImage[]>([])
  const totalImages = ref(0)
  const currentPage = ref(1)
  const pageSize = ref(20)
  const sortBy = ref('newest')
  const processingImages = reactive<Record<string, boolean>>({})
  const isBatchProcessing = ref(false)
  const isBatchApproving = ref(false)
  const isBatchRejecting = ref(false)

  /* 动态页面大小计算 - 基于网格布局自动计算合理的文件数量 */
  const { pageSize: autoPageSize, recalc } = useResponsivePageSize({
    containerSelector: '.review-queue-container',
    gridSelector: '.images-grid',
    itemMinWidth: 200,
    columnGap: 16,
    rowMultiple: 3,
    defaultSize: 20,
    debug: false,
    preferCssColumns: true,
    mode: 'once',
  })

  const showRejectDialog = ref(false)
  const showBatchRejectDialog = ref(false)
  const showBatchApproveDialog = ref(false)
  const showApproveDialog = ref(false)
  const selectedRejectImage = ref<ReviewImage | null>(null)
  const selectedApproveImage = ref<ReviewImage | null>(null)
  const rejectReason = ref('')
  const batchRejectReason = ref('')
  const batchApproveReason = ref('')
  const approveReason = ref('')
  const hardDeleteOnReject = ref(false)
  const batchHardDelete = ref(false)
  const isProcessingReject = ref(false)
  const isProcessingApprove = ref(false)
  const showPreview = ref(false)
  const previewImage = ref<ReviewImage | null>(null)

  const loadData = async (page = 1) => {
    try {
      isLoading.value = true
      showSkeleton.value = true

      const params: ReviewQueueQuery = {
        page,
        size: pageSize.value,
        sort: sortBy.value,
      }

      const response = await getReviewQueue(params)

      if (response.success && response.data) {
        images.value = response.data.data || []
        totalImages.value = response.data.pagination?.total || 0
        currentPage.value = response.data.pagination?.page || 1
      } else {
        images.value = []
        totalImages.value = 0
        currentPage.value = 1
      }
    } catch (error) {
      const isRequestCancelled = error && typeof error === 'object' && 'code' in error && error.code === -1
      if (isRequestCancelled) {
      } else {
        console.error('[ReviewQueue] 加载数据失败:', error)
        toast.error(REVIEW_MESSAGES.LOAD_ERROR)
        images.value = []
        totalImages.value = 0
        currentPage.value = 1
      }
    } finally {
      isLoading.value = false
      showSkeleton.value = false
    }

    await loadStats()
  }

  const loadStats = async () => {
    try {
      const statsResponse = await getReviewStats()
      if (statsResponse.success) {
        emit('refresh-stats', statsResponse.data)
      }
    } catch {}
  }

  const refreshData = () => {
    loadData(currentPage.value)
  }

  const handleImageClick = (image: ReviewImage) => {
    if (selectMode.value) {
      toggleImageSelection(image.id)
    } else {
      showFilePreview(image)
    }
  }

  const _approveImage = async (fileId: string) => {
    try {
      processingImages[fileId] = true

      const response = await reviewImage({
        file_id: fileId,
        action: 'approve',
        reason: DEFAULT_REVIEW_REASONS.APPROVE,
      })

      if (response.success) {
        toast.success(REVIEW_MESSAGES.APPROVE_SUCCESS)
        images.value = images.value.filter((img) => img.id !== fileId)
        totalImages.value--
        await loadStats()
      } else {
        toast.error(response.message || REVIEW_MESSAGES.OPERATION_ERROR)
      }
    } catch {
      toast.error(REVIEW_MESSAGES.OPERATION_ERROR)
    } finally {
      processingImages[fileId] = false
    }
  }

  const confirmRejectImage = (image: ReviewImage) => {
    selectedRejectImage.value = image
    rejectReason.value = ''
    hardDeleteOnReject.value = false
    showRejectDialog.value = true
  }

  const executeReject = async () => {
    if (!selectedRejectImage.value) {
      return
    }

    try {
      isProcessingReject.value = true

      const response = await reviewImage({
        file_id: selectedRejectImage.value.id,
        action: 'reject',
        reason: rejectReason.value || DEFAULT_REVIEW_REASONS.REJECT,
        hard_delete: hardDeleteOnReject.value,
      })

      if (response.success) {
        toast.success(REVIEW_MESSAGES.REJECT_SUCCESS)
        images.value = images.value.filter((img) => img.id !== selectedRejectImage.value?.id)
        totalImages.value--
        showRejectDialog.value = false
        selectedRejectImage.value = null
        await loadStats()
      } else {
        toast.error(response.message || REVIEW_MESSAGES.OPERATION_ERROR)
      }
    } catch {
      toast.error(REVIEW_MESSAGES.OPERATION_ERROR)
    } finally {
      isProcessingReject.value = false
    }
  }

  const confirmBatchApprove = () => {
    batchApproveReason.value = ''
    showBatchApproveDialog.value = true
  }

  const executeBatchApprove = async () => {
    if (selectedImages.value.length === 0) {
      return
    }

    try {
      isBatchApproving.value = true
      isBatchProcessing.value = true

      const response = await batchReview({
        file_ids: selectedImages.value,
        action: 'approve',
        reason: batchApproveReason.value || DEFAULT_REVIEW_REASONS.BATCH_APPROVE,
      })

      if (response.success) {
        const successCount = response.data?.success_count || 0
        toast.success(REVIEW_MESSAGES.BATCH_APPROVE_SUCCESS(successCount))

        images.value = images.value.filter((img) => !selectedImages.value.includes(img.id))
        totalImages.value -= successCount
        showBatchApproveDialog.value = false
        exitSelectMode()
        await loadStats()
      } else {
        toast.error(response.message || REVIEW_MESSAGES.OPERATION_ERROR)
      }
    } catch {
      toast.error(REVIEW_MESSAGES.OPERATION_ERROR)
    } finally {
      isBatchApproving.value = false
      isBatchProcessing.value = false
    }
  }

  const confirmBatchReject = () => {
    batchRejectReason.value = ''
    batchHardDelete.value = false
    showBatchRejectDialog.value = true
  }

  const executeBatchReject = async () => {
    if (selectedImages.value.length === 0) {
      return
    }

    try {
      isBatchRejecting.value = true
      isBatchProcessing.value = true

      const response = await batchReview({
        file_ids: selectedImages.value,
        action: 'reject',
        reason: batchRejectReason.value || DEFAULT_REVIEW_REASONS.BATCH_REJECT,
        hard_delete: batchHardDelete.value,
      })

      if (response.success) {
        const successCount = response.data?.success_count || 0
        toast.success(REVIEW_MESSAGES.BATCH_REJECT_SUCCESS(successCount))

        images.value = images.value.filter((img) => !selectedImages.value.includes(img.id))
        totalImages.value -= successCount
        showBatchRejectDialog.value = false
        exitSelectMode()
        await loadStats()
      } else {
        toast.error(response.message || REVIEW_MESSAGES.OPERATION_ERROR)
      }
    } catch {
      toast.error(REVIEW_MESSAGES.OPERATION_ERROR)
    } finally {
      isBatchRejecting.value = false
      isBatchProcessing.value = false
    }
  }

  const confirmApproveImage = (image: ReviewImage) => {
    selectedApproveImage.value = image
    approveReason.value = ''
    showApproveDialog.value = true
  }

  const executeApprove = async () => {
    if (!selectedApproveImage.value) {
      return
    }

    try {
      isProcessingApprove.value = true
      processingImages[selectedApproveImage.value.id] = true

      const response = await reviewImage({
        file_id: selectedApproveImage.value.id,
        action: 'approve',
        reason: approveReason.value || DEFAULT_REVIEW_REASONS.APPROVE,
      })

      if (response.success) {
        toast.success(REVIEW_MESSAGES.APPROVE_SUCCESS)
        images.value = images.value.filter((img) => img.id !== selectedApproveImage.value?.id)
        totalImages.value -= 1
        showApproveDialog.value = false
        await loadStats()
      } else {
        toast.error(response.message || REVIEW_MESSAGES.OPERATION_ERROR)
      }
    } catch {
      toast.error(REVIEW_MESSAGES.OPERATION_ERROR)
    } finally {
      if (selectedApproveImage.value) {
        delete processingImages[selectedApproveImage.value.id]
      }
      isProcessingApprove.value = false
    }
  }

  const showFilePreview = (image: ReviewImage) => {
    previewImage.value = image
    showPreview.value = true
  }

  const isInitialized = ref(false)

  watch(
    () => props.isActive,
    async (newVal) => {
      if (newVal) {
        if (!isInitialized.value) {
          isLoading.value = true
          showSkeleton.value = true

          await nextTick()
          setTimeout(() => {
            recalc(true)
            if (autoPageSize.value && autoPageSize.value !== pageSize.value) {
              pageSize.value = autoPageSize.value
            }
            isInitialized.value = true
            loadData(1)
          }, 200)
        } else {
          loadData(1)
        }
      }
    },
    { immediate: true }
  )

  watch(currentPage, (newPage) => {
    if (newPage && newPage !== 1) {
      loadData(newPage)
    }
  })

  watch(pageSize, (newSize, oldSize) => {
    if (newSize && newSize !== oldSize && isInitialized.value) {
      currentPage.value = 1
      loadData(1)
    }
  })

  defineExpose({
    loadData,
    refreshData,
  })
</script>

<template>
  <div class="review-queue-container" :data-select-mode="selectMode">
    <div class="queue-header">
      <div class="header-info">
        <span v-if="!isLoading && totalImages > 0" class="queue-count">
          {{ $t('admin.contentReview.queue.selection.selectedCount', { count: totalImages, total: totalImages }) }}
        </span>
      </div>

      <div class="header-actions">
        <Button type="secondary" :loading="isLoading" @click="refreshData">
          <i class="fas fa-sync-alt mr-1.5" />{{ $t('admin.contentReview.queue.buttons.refresh') }}
        </Button>

        <Button v-if="!selectMode" type="secondary" :disabled="images.length === 0" @click="startSelectMode">
          <i class="fas fa-tasks mr-1.5" />{{ $t('admin.contentReview.queue.buttons.batchOperation') }}
        </Button>

        <div v-else class="batch-controls">
          <div class="selection-info">
            <span class="selection-count">
              {{
                $t('admin.contentReview.queue.selection.selectedCount', { count: selectedImages.length, total: images.length })
              }}
            </span>
          </div>

          <Button
            type="success"
            size="small"
            :disabled="selectedImages.length === 0 || isBatchProcessing"
            :loading="isBatchApproving"
            @click="confirmBatchApprove"
          >
            <i class="fas fa-check mr-1" />{{
              $t('admin.contentReview.queue.batchActions.approve', { count: selectedImages.length })
            }}
          </Button>
          <Button
            type="danger"
            size="small"
            :disabled="selectedImages.length === 0 || isBatchProcessing"
            :loading="isBatchRejecting"
            @click="confirmBatchReject"
          >
            <i class="fas fa-times mr-1" />{{
              $t('admin.contentReview.queue.batchActions.reject', { count: selectedImages.length })
            }}
          </Button>
          <Button type="secondary" size="small" @click="exitSelectMode">
            <i class="fas fa-times mr-1" />{{ $t('admin.contentReview.queue.buttons.cancel') }}
          </Button>
        </div>
      </div>
    </div>

    <div v-if="selectMode" class="toolbar">
      <div class="toolbar-left">
        <div class="selection-controls">
          <button
            class="select-btn select-all-btn"
            :class="{ active: isAllSelected }"
            :title="
              isAllSelected
                ? $t('admin.contentReview.queue.selection.deselectAll')
                : $t('admin.contentReview.queue.selection.selectAll')
            "
            @click="() => toggleSelectAll(images)"
          >
            <i class="fas" :class="isAllSelected ? 'fa-check-double' : 'fa-check'" />
            {{
              isAllSelected
                ? $t('admin.contentReview.queue.selection.deselectAll')
                : $t('admin.contentReview.queue.selection.selectAll')
            }}
          </button>
          <button
            class="select-btn invert-btn"
            :title="$t('admin.contentReview.queue.selection.invertSelection')"
            @click="() => selectInvert(images)"
          >
            <i class="fas fa-exchange-alt" />
            {{ $t('admin.contentReview.queue.selection.invertSelection') }}
          </button>
        </div>
      </div>
    </div>

    <div class="content-area">
      <div class="images-grid">
        <CyberSkeleton type="card" :count="autoPageSize || pageSize" :loading="showSkeleton" simple />
        <template v-if="!showSkeleton && images.length > 0">
          <div
            v-for="image in images"
            :key="image.id"
            class="review-image-card"
            :class="{
              selected: selectedImages.includes(image.id),
              processing: processingImages[image.id],
            }"
          >
            <div v-if="selectMode" class="image-select" @click.stop>
              <cyberCheckbox
                :model-value="selectedImages.includes(image.id)"
                @update:model-value="toggleImageSelection(image.id)"
              />
            </div>

            <div class="image-preview" @click="handleImageClick(image)">
              <img
                :src="image.full_thumb_url"
                :alt="image.display_name || image.original_name"
                class="preview-image"
                @error="handleImageError"
              />

              <div v-if="image.ai_info?.is_nsfw || image.nsfw" class="nsfw-badge">
                <i class="fas fa-exclamation-triangle mr-1" />
                <span class="nsfw-label">NSFW</span>
                <span v-if="image.ai_info?.nsfw_score" class="nsfw-score">
                  {{ Math.round(image.ai_info.nsfw_score * 100) }}%
                </span>
              </div>

              <div class="date-overlay">
                {{ formatDate(image.created_at, 'MM-DD HH:mm') }}
              </div>

              <div v-if="!selectMode" class="action-overlay">
                <Button
                  type="success"
                  size="small"
                  :disabled="processingImages[image.id]"
                  @click.stop="confirmApproveImage(image)"
                >
                  <i class="fas fa-check mr-1" />
                  {{ $t('admin.contentReview.queue.buttons.pass') }}
                </Button>
                <Button type="danger" size="small" :disabled="processingImages[image.id]" @click.stop="confirmRejectImage(image)">
                  <i class="fas fa-times mr-1" />
                  {{ $t('admin.contentReview.queue.buttons.reject') }}
                </Button>
              </div>
            </div>

            <div class="image-info">
              <div class="image-title" :title="image.display_name || image.original_name">
                {{ image.display_name || image.original_name }}
              </div>

              <div class="image-meta">
                <span class="meta-item">
                  <i class="fas fa-image" />
                  {{ image.width }}×{{ image.height }}
                </span>
                <span class="meta-item">
                  <i class="fas fa-hdd" />
                  {{ image.size_formatted }}
                </span>
                <span v-if="image.uploader" class="meta-item uploader-info">
                  <i class="fas fa-user" />
                  {{ image.uploader.username }}
                </span>
              </div>
            </div>
          </div>
        </template>
        <template v-else-if="!showSkeleton && images.length === 0">
          <div class="empty-state-container">
            <div class="empty-state">
              <i class="fas fa-check-circle mb-4 text-6xl text-green-400/30" />
              <h3 class="text-content-content-muted mb-2 text-lg">{{ $t('admin.contentReview.queue.empty.title') }}</h3>
              <p class="text-content-disabled">{{ $t('admin.contentReview.queue.empty.description') }}</p>
            </div>
          </div>
        </template>
      </div>

      <div v-if="totalImages > 0" class="pagination-wrapper">
        <CyberPagination v-model:current-page="currentPage" v-model:page-size="pageSize" :total="totalImages" />
      </div>
    </div>

    <Dialog
      v-model="showRejectDialog"
      :title="$t('admin.contentReview.dialogs.rejectSingle.title')"
      :loading="isProcessingReject"
      @confirm="executeReject"
      @cancel="showRejectDialog = false"
    >
      <div class="space-y-4">
        <p class="text-content">
          {{ $t('admin.contentReview.dialogs.rejectSingle.message', { name: selectedRejectImage?.original_name }) }}
        </p>
        <div class="form-item">
          <label class="mb-2 block text-sm text-content">{{ $t('admin.contentReview.dialogs.rejectSingle.reasonLabel') }}</label>
          <Input
            v-model="rejectReason"
            :placeholder="$t('admin.contentReview.dialogs.rejectSingle.reasonPlaceholder')"
            type="textarea"
            rows="3"
          />
        </div>
        <div class="form-item">
          <cyberCheckbox v-model="hardDeleteOnReject">{{
            $t('admin.contentReview.dialogs.rejectSingle.hardDelete')
          }}</cyberCheckbox>
        </div>
      </div>
    </Dialog>

    <Dialog
      v-model="showBatchApproveDialog"
      :title="$t('admin.contentReview.dialogs.batchApprove.title')"
      :loading="isBatchApproving"
      @confirm="executeBatchApprove"
      @cancel="showBatchApproveDialog = false"
    >
      <div class="space-y-4">
        <p class="text-content">{{ $t('admin.contentReview.dialogs.batchApprove.message', { count: selectedImages.length }) }}</p>
        <div class="form-item">
          <label class="mb-2 block text-sm text-content">{{ $t('admin.contentReview.dialogs.batchApprove.reasonLabel') }}</label>
          <Input
            v-model="batchApproveReason"
            :placeholder="$t('admin.contentReview.dialogs.batchApprove.reasonPlaceholder')"
            type="textarea"
            rows="3"
          />
        </div>
        <div class="rounded border border-green-500/20 bg-green-500/10 p-3">
          <p class="text-xs text-green-400">
            <i class="fas fa-info-circle mr-1" />
            {{ $t('admin.contentReview.dialogs.batchApprove.infoTitle') }}：
          </p>
          <ul class="ml-4 mt-2 space-y-1 text-xs text-content-muted">
            <li v-for="(item, index) in $t('admin.contentReview.dialogs.batchApprove.infoItems')" :key="index">• {{ item }}</li>
          </ul>
        </div>
      </div>
    </Dialog>

    <Dialog
      v-model="showBatchRejectDialog"
      :title="$t('admin.contentReview.dialogs.batchReject.title')"
      :loading="isBatchRejecting"
      @confirm="executeBatchReject"
      @cancel="showBatchRejectDialog = false"
    >
      <div class="space-y-4">
        <p class="text-content">{{ $t('admin.contentReview.dialogs.batchReject.message', { count: selectedImages.length }) }}</p>
        <div class="form-item">
          <label class="mb-2 block text-sm text-content">{{ $t('admin.contentReview.dialogs.batchReject.reasonLabel') }}</label>
          <Input
            v-model="batchRejectReason"
            :placeholder="$t('admin.contentReview.dialogs.batchReject.reasonPlaceholder')"
            type="textarea"
            rows="3"
          />
        </div>
        <div class="form-item">
          <cyberCheckbox v-model="batchHardDelete">{{ $t('admin.contentReview.dialogs.batchReject.hardDelete') }}</cyberCheckbox>
        </div>
      </div>
    </Dialog>

    <Dialog
      v-model="showApproveDialog"
      :title="$t('admin.contentReview.dialogs.approveSingle.title')"
      :loading="isProcessingApprove"
      @confirm="executeApprove"
      @cancel="showApproveDialog = false"
    >
      <div class="space-y-4">
        <p class="text-content">{{ $t('admin.contentReview.dialogs.approveSingle.message') }}</p>
        <div class="form-item">
          <label class="mb-2 block text-sm text-content">{{ $t('admin.contentReview.dialogs.approveSingle.reasonLabel') }}</label>
          <Input
            v-model="approveReason"
            :placeholder="$t('admin.contentReview.dialogs.approveSingle.reasonPlaceholder')"
            type="textarea"
            rows="3"
          />
        </div>
        <div class="rounded border border-green-500/20 bg-green-500/10 p-3">
          <p class="text-xs text-green-400">
            <i class="fas fa-info-circle mr-1" />
            {{ $t('admin.contentReview.dialogs.approveSingle.infoTitle') }}：
          </p>
          <ul class="ml-4 mt-2 space-y-1 text-xs text-content-muted">
            <li v-for="(item, index) in $t('admin.contentReview.dialogs.approveSingle.infoItems')" :key="index">• {{ item }}</li>
          </ul>
        </div>
      </div>
    </Dialog>

    <cyberFileViewer v-model="showPreview" :file="previewImage" />
  </div>
</template>

<style scoped lang="scss">
  .review-queue-container {
    @apply flex flex-col;
    flex: 1;
    min-height: 0;
  }

  .queue-header {
    @apply mb-4 flex flex-shrink-0 items-start justify-between;
  }

  .header-info {
    @apply space-y-1;
  }

  .queue-count {
    @apply text-sm text-content-muted;
  }

  .header-actions {
    @apply flex items-center gap-3;
  }

  .batch-controls {
    @apply flex items-center gap-2;
  }

  .selection-info {
    @apply mr-2 flex items-center gap-2;
  }

  .selection-count {
    @apply rounded px-3 py-1.5 text-sm font-medium;
    background: linear-gradient(135deg, rgba(var(--color-brand-500-rgb), 0.15) 0%, rgba(var(--color-brand-600-rgb), 0.1) 100%);
    color: var(--color-brand-400);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.25);
    box-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
    backdrop-filter: blur(4px);
  }

  .toolbar {
    @apply mb-4 flex flex-shrink-0 items-center justify-between rounded-lg bg-background-400 px-4 py-3;
  }

  .toolbar-left {
    @apply flex items-center gap-2;
  }

  .selection-controls {
    display: flex;
    align-items: center;
    gap: 0.4rem;
  }

  .select-btn {
    background: rgba(var(--color-brand-500-rgb), 0.1);
    color: var(--color-brand-500);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.2);
    border-radius: var(--radius-sm);
    padding: 0.25rem 0.75rem;
    font-size: 0.75rem;
    cursor: pointer;
    transition: all 0.3s ease;
    display: flex;
    align-items: center;
    gap: 0.25rem;
  }

  .select-btn:hover {
    background: rgba(var(--color-brand-500-rgb), 0.2);
    border-color: rgba(var(--color-brand-500-rgb), 0.4);
  }

  .select-btn.active {
    background: rgba(var(--color-brand-500-rgb), 0.3);
    border-color: rgba(var(--color-brand-500-rgb), 0.5);
  }

  .select-btn i {
    font-size: 0.75rem;
  }

  .content-area {
    @apply flex flex-1 flex-col overflow-hidden;
  }

  .images-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
    gap: 1rem;
    align-content: start;
    contain: layout paint;
    will-change: contents;
    padding: 1rem;
    flex: 1;
    overflow-y: auto;
  }

  @media (max-width: 1536px) {
    .images-grid {
      grid-template-columns: repeat(auto-fill, minmax(180px, 1fr));
    }
  }

  @media (max-width: 1280px) {
    .images-grid {
      grid-template-columns: repeat(auto-fill, minmax(160px, 1fr));
    }
  }

  @media (max-width: 1024px) {
    .images-grid {
      grid-template-columns: repeat(auto-fill, minmax(140px, 1fr));
    }
  }

  @media (max-width: 768px) {
    .images-grid {
      grid-template-columns: repeat(auto-fill, minmax(120px, 1fr));
      gap: 0.75rem;
    }
  }

  @media (max-width: 640px) {
    .images-grid {
      grid-template-columns: repeat(auto-fill, minmax(100px, 1fr));
      gap: 0.5rem;
    }
  }

  .skeleton-card {
    background: linear-gradient(
      135deg,
      rgba(var(--color-background-800-rgb), 0.95) 0%,
      rgba(var(--color-background-900-rgb), 0.9) 100%
    ) !important;
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.25) !important;
    border-radius: var(--radius-sm) !important;
    backdrop-filter: blur(8px) !important;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.2) !important;

    margin: 0 !important;
    padding: 0;
    box-sizing: border-box !important;

    min-width: 0;
    width: 100%;
  }

  .review-image-card {
    background: linear-gradient(
      135deg,
      rgba(var(--color-background-800-rgb), 0.95) 0%,
      rgba(var(--color-background-900-rgb), 0.9) 100%
    );
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.25);
    border-radius: var(--radius-sm);
    overflow: hidden;
    transition: all 0.3s ease;
    position: relative;
    backdrop-filter: blur(8px);

    box-sizing: border-box;
    margin: 0;

    &:hover {
      transform: translateY(-6px);
      border-color: rgba(var(--color-brand-500-rgb), 0.5);
      box-shadow:
        0 0 25px rgba(var(--color-brand-500-rgb), 0.2),
        0 10px 40px rgba(0, 0, 0, 0.3);
    }

    &.selected {
      border-color: var(--color-brand-500);
      background: linear-gradient(135deg, rgba(var(--color-brand-500-rgb), 0.15) 0%, rgba(var(--color-error-rgb), 0.1) 100%);
      box-shadow:
        0 0 0 2px rgba(var(--color-brand-500-rgb), 0.4),
        0 0 30px rgba(var(--color-brand-500-rgb), 0.2);
    }

    &.processing {
      opacity: 0.7;
      filter: grayscale(0.3);
    }
  }

  .review-queue-container[data-select-mode='true'] .review-image-card {
    cursor: pointer;

    &:hover {
      border-color: rgba(var(--color-brand-500-rgb), 0.6);
      transform: translateY(-2px);
    }

    &:not(.selected):hover {
      background: linear-gradient(135deg, rgba(var(--color-brand-500-rgb), 0.08) 0%, rgba(var(--color-error-rgb), 0.05) 100%);
    }
  }

  .image-select {
    position: absolute;
    top: 8px;
    left: 8px;
    z-index: 20;
    cursor: pointer;
  }

  .image-preview {
    position: relative;
    cursor: pointer;
    overflow: hidden;
    aspect-ratio: 4/3;
  }

  .review-queue-container[data-select-mode='true'] .image-preview {
    cursor: pointer;
  }

  .review-queue-container[data-select-mode='true'] .image-preview:hover {
    opacity: 0.9;
  }

  .preview-image {
    width: 100%;
    height: 100%;
    object-fit: cover;
    transition: transform 0.3s ease;
  }

  .image-preview:hover .preview-image {
    transform: scale(1.05);
  }

  .nsfw-badge {
    position: absolute;
    top: var(--space-sm);
    right: var(--space-sm);
    display: inline-flex;
    align-items: center;
    gap: 0.35rem;
    background: linear-gradient(135deg, rgba(var(--color-error-rgb), 0.92) 0%, rgba(var(--color-error-rgb), 0.85) 100%);
    color: var(--color-text-on-brand);
    padding: 0.35rem 0.6rem;
    border-radius: var(--radius-md);
    font-size: var(--text-xs);
    font-weight: var(--font-bold);
    line-height: 1;
    z-index: 6;
    box-shadow:
      0 2px 6px rgba(var(--color-error-rgb), 0.35),
      0 4px 12px rgba(0, 0, 0, 0.15),
      inset 0 1px 0 rgba(255, 255, 255, 0.15);
    border: 1px solid rgba(var(--color-error-rgb), 0.7);
    text-shadow: 0 1px 2px rgba(0, 0, 0, 0.4);
    backdrop-filter: blur(6px);
    transition: all var(--transition-fast) var(--ease-out);
  }

  .nsfw-badge:hover {
    transform: translateY(-1px) scale(1.02);
    box-shadow:
      0 3px 8px rgba(var(--color-error-rgb), 0.45),
      0 6px 16px rgba(0, 0, 0, 0.2),
      inset 0 1px 0 rgba(255, 255, 255, 0.2);
  }

  .nsfw-badge i {
    font-size: 11px;
    opacity: 0.95;
  }

  .nsfw-label {
    font-weight: var(--font-extrabold);
    letter-spacing: 0.3px;
  }

  .nsfw-score {
    padding: 0.15rem 0.35rem;
    background: rgba(0, 0, 0, 0.25);
    border-radius: var(--radius-sm);
    font-size: 10px;
    font-weight: var(--font-bold);
    border: 1px solid rgba(0, 0, 0, 0.15);
    line-height: 1;
  }

  .date-overlay {
    position: absolute;
    bottom: var(--space-sm);
    left: var(--space-sm);
    background: rgba(var(--color-background-900-rgb), 0.7);
    color: var(--color-content-heading);
    padding: var(--space-xs) var(--space-sm);
    border-radius: var(--radius-sm);
    font-size: var(--text-xs);
    font-weight: var(--font-medium);
    z-index: 5;
    text-shadow: var(--shadow-glow-sm);
    backdrop-filter: var(--backdrop-blur-sm);
  }

  .action-overlay {
    position: absolute;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    display: flex;
    gap: var(--space-md);
    opacity: 0;
    transition: opacity var(--transition-normal) var(--ease-out);
    z-index: 10;
    pointer-events: none;
  }

  .image-preview:hover .action-overlay {
    opacity: 1;
    pointer-events: auto;
  }

  .image-info {
    padding: var(--space-sm) var(--space-md);
    background: linear-gradient(
      135deg,
      rgba(var(--color-background-800-rgb), 0.95) 0%,
      rgba(var(--color-background-900-rgb), 0.85) 100%
    );
    backdrop-filter: var(--backdrop-blur-sm);
    border-top: 1px solid rgba(var(--color-brand-500-rgb), 0.1);
  }

  .image-title {
    color: var(--color-content-heading);
    font-weight: var(--font-medium);
    font-size: var(--text-sm);
    margin-bottom: var(--space-xs);
    display: -webkit-box;
    -webkit-line-clamp: 1;
    -webkit-box-orient: vertical;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .image-meta {
    display: flex;
    flex-wrap: wrap;
    gap: var(--space-sm);
    align-items: center;
  }

  .meta-item {
    display: inline-flex;
    align-items: center;
    gap: 0.3rem;
    color: var(--color-content-muted);
    font-size: 0.7rem;
    line-height: 1.2;
    white-space: nowrap;

    i {
      font-size: 0.65rem;
      color: var(--color-brand-500);
      opacity: 0.6;
    }

    &.uploader-info {
      i {
        color: var(--color-success-500);
        opacity: 0.7;
      }
    }
  }

  .empty-state-container {
    @apply col-span-full;
  }

  .empty-state {
    @apply flex flex-col items-center justify-center py-20;
  }

  .pagination-wrapper {
    @apply flex-shrink-0 border-t border-subtle p-4;
  }

  .form-item {
    @apply space-y-2;
  }
</style>
