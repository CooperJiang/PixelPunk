<template>
  <div class="admin-vectors-page admin-page-container">
    <CyberAdminWrapper :title="$t('admin.vectors.title')" :subtitle="$t('admin.vectors.subtitle')" icon="fas fa-vector-square">
      <template #actions>
        <div class="flex items-center gap-2">
          <span class="text-xs">{{ $t('admin.vectors.qdrant.systemInfo') }}</span>
          <CyberSwitch
            v-model="autoEnabled"
            :disabled="updatingAuto"
            @change="
              async (v: boolean) => {
                try {
                  updatingAuto = true
                  await batchUpsertSettings([
                    {
                      key: 'vector_auto_processing_enabled',
                      value: v,
                      type: 'boolean',
                      group: 'vector',
                      description: $t('admin.vectors.settings.autoProcessingDesc'),
                    },
                  ])
                  await setAutoProcessing(v)
                } finally {
                  updatingAuto = false
                }
              }
            "
          />
        </div>

        <div class="bg-cyber-border mx-3 h-4 w-px"></div>

        <div class="flex items-center gap-2">
          <span class="whitespace-nowrap text-xs">{{ $t('admin.vectors.qdrant.throughput') }}</span>
          <CyberSlider
            v-model="concurrency"
            :min="1"
            :max="30"
            :step="1"
            class="w-20"
            @change="
              async () => {
                try {
                  await batchUpsertSettings([
                    {
                      key: 'vector_concurrency',
                      value: concurrency,
                      type: 'number',
                      group: 'vector',
                      description: $t('admin.vectors.settings.concurrencyDesc'),
                    },
                  ])
                  await setConcurrency(concurrency)
                } catch {}
              }
            "
          />
          <span class="w-6 text-center text-xs">{{ concurrency }}</span>
        </div>
      </template>
      <template #default>
        <div class="module-wrapper">
          <div class="module-header">
            <h3 class="module-title">
              <i class="fas fa-chart-bar mr-2" />
              {{ $t('admin.vectors.overview.title') }}
            </h3>
          </div>
          <div class="module-content">
            <VectorStatusCard
              :selected-status="filters.status as string"
              :is-loading="isStatusSelecting"
              @select="handleStatusSelect"
            />
          </div>
        </div>

        <div class="module-wrapper">
          <div class="module-header">
            <h3 class="module-title">
              <span class="icon">🔍</span>
              {{ $t('admin.vectors.qdrant.title') }}
            </h3>
          </div>
          <div class="module-content">
            <QdrantStatsCard :qdrant-stats="qdrantRealStats" />
          </div>
        </div>

        <div class="module-wrapper">
          <div class="module-header">
            <h3 class="module-title">
              <i class="fas fa-cogs mr-2" />
              {{ $t('admin.vectors.vectors.vectorActions') }}
            </h3>
          </div>
          <div class="module-content">
            <VectorActionPanelSimplified
              :processing="isProcessing"
              @regenerate-all="handleRegenerateAll"
              @retry-all-failed="handleRetryAllFailed"
              @refresh="handleRefresh"
              @verification-completed="handleVerificationCompleted"
            />
          </div>
        </div>

        <div class="module-wrapper">
          <div class="module-header">
            <h3 class="module-title">
              <i class="fas fa-list mr-2" />
              {{ $t('admin.vectors.vectors.title') }}
            </h3>
          </div>
          <div class="module-content">
            <VectorList
              :vectors="vectors"
              :is-loading="isLoading"
              :pagination="pagination"
              :processing="isProcessing"
              :show-filter="showFilter"
              :available-models="availableModels"
              :filters="filters"
              @retry="handleRetryVector"
              @batch-action="handleBatchAction"
              @page-change="handlePageChange"
              @page-size-change="handlePageSizeChange"
              @view-detail="handleViewDetail"
              @toggle-filter="showFilter = !showFilter"
              @filter="handleFilter"
            />
          </div>
        </div>

        <div v-if="showLogs" class="module-wrapper module-wrapper-last">
          <div class="module-header">
            <h3 class="module-title">
              <i class="fas fa-history mr-2" />
              {{ $t('admin.vectors.logs.title') }}
            </h3>
          </div>
          <div class="module-content">
            <VectorLogList
              :logs="logs"
              :pagination="logPagination"
              :selected-file-id="selectedFileId"
              @load-more="handleLogLoadMore"
              @refresh="handleLogRefresh"
              @clear-filter="handleClearLogFilter"
            />
          </div>
        </div>
      </template>
    </CyberAdminWrapper>

    <CyberConfirmDialog
      v-model="showDeleteConfirm"
      type="danger"
      :title="$t('admin.vectors.confirm.delete.title')"
      :message="deleteConfirmMessage"
      :require-input="true"
      :input-label="deleteInputLabel"
      :input-placeholder="deleteConfirmData.expectedText"
      :expected-input="deleteConfirmData.expectedText"
      :confirm-text="$t('admin.vectors.confirm.delete.confirm')"
      :cancel-text="$t('admin.vectors.confirm.delete.cancel')"
      @confirm="handleDeleteConfirm"
      @cancel="handleDialogCancel"
    />

    <CyberConfirmDialog
      v-model="showBatchConfirm"
      type="warning"
      :title="$t('admin.vectors.batchOperations.title')"
      :message="[batchConfirmData.message]"
      :confirm-text="$t('admin.vectors.buttons.confirm')"
      :cancel-text="$t('admin.vectors.buttons.cancel')"
      @confirm="handleBatchConfirm"
      @cancel="handleDialogCancel"
    />

    <CyberConfirmDialog
      v-model="showRegenerateConfirm"
      type="danger"
      :title="$t('admin.vectors.confirm.reindex.title')"
      :message="regenerateConfirmData.message"
      :confirm-text="$t('admin.vectors.confirm.reindex.confirm')"
      :cancel-text="$t('admin.vectors.confirm.reindex.cancel')"
      @confirm="handleRegenerateConfirm"
      @cancel="handleDialogCancel"
    />

    <CyberConfirmDialog
      v-model="showRetryFailedConfirm"
      type="warning"
      :title="$t('admin.vectors.confirm.reindex.title')"
      :message="[retryFailedConfirmData.message]"
      :confirm-text="$t('admin.vectors.confirm.reindex.confirm')"
      :cancel-text="$t('admin.vectors.confirm.reindex.cancel')"
      @confirm="handleRetryFailedConfirm"
      @cancel="handleDialogCancel"
    />
  </div>
</template>

<script setup lang="ts">
  import { ref, reactive, computed, onMounted, onUnmounted } from 'vue'
  import { useToast } from '@/components/Toast/useToast'
  import { useWebSocketStore } from '@/store/websocket'
  import { useTexts } from '@/composables/useTexts'
  import {
    getAutoProcessing,
    setAutoProcessing,
    setConcurrency,
    retryAllFailedVectors,
    regenerateAllVectors,
    type VectorListParams,
  } from '@/api/admin/vectors'
  import { getSettings, batchUpsertSettings, type Setting } from '@/api/admin/settings'

  import VectorStatusCard from './components/VectorStatusCard.vue'
  import QdrantStatsCard from './components/QdrantStatsCard.vue'
  import VectorList from './components/VectorList.vue'
  import VectorLogList from './components/VectorLogList.vue'
  import VectorActionPanelSimplified from './components/VectorActionPanelSimplified.vue'

  import { useVectorData, useVectorActions, useVectorDialogs } from './composables'

  defineOptions({
    name: 'VectorManagePage',
  })

  const toast = useToast()
  const { $t } = useTexts()

  const {
    vectors,
    qdrantRealStats,
    logs,
    availableModels,
    pagination,
    logPagination,
    isLoading,
    selectedFileId,
    loadVectors,
    loadQdrantStats,
    loadAvailableModels,
    loadVectorLogs,
    cleanup: cleanupData,
  } = useVectorData()

  const { isProcessing, handleRetryVector: retryVector, handleBatchActionExecute } = useVectorActions()

  const {
    showDeleteConfirm,
    showBatchConfirm,
    showRegenerateConfirm,
    showRetryFailedConfirm,
    deleteConfirmData,
    batchConfirmData,
    regenerateConfirmData,
    retryFailedConfirmData,
  } = useVectorDialogs()

  const showFilter = ref(false)
  const showLogs = ref(true)

  const autoEnabled = ref<boolean>(true)
  const updatingAuto = ref<boolean>(false)
  const concurrency = ref<number>(0)

  const filters = reactive<VectorListParams>({
    page: 1,
    page_size: 10,
    status: undefined,
    model: undefined,
    keyword: undefined,
  })

  const deleteInputLabel = computed(() => {
    return `${$t('admin.vectors.confirm.delete.message')} "${deleteConfirmData.value.expectedText}"`
  })

  const deleteConfirmMessage = computed(() => {
    return [
      `⚠️ ${$t('admin.vectors.error.title')}`,
      '',
      `${$t('admin.vectors.confirm.batchDelete.messagePrefix')} ${deleteConfirmData.value.count} ${$t('admin.vectors.vectors.title')}！`,
      $t('admin.vectors.confirm.delete.message'),
      '',
      $t('admin.vectors.form.placeholders.collectionName'),
    ]
  })

  const isStatusSelecting = ref(false)

  const handleStatusSelect = async (status: string) => {
    try {
      isStatusSelecting.value = true
      filters.status = status === 'all' ? undefined : (status as VectorListParams['status'])
      filters.page = 1
      await loadVectors(filters)
    } catch {
    } finally {
      isStatusSelecting.value = false
    }
  }

  const handleFilter = (newFilters: VectorListParams) => {
    Object.assign(filters, newFilters)
    filters.page = 1
    loadVectors(filters)
    showFilter.value = false
  }

  const handlePageChange = (page: number) => {
    filters.page = page
    loadVectors(filters)
  }

  const handlePageSizeChange = (size: number) => {
    filters.page_size = size
    filters.page = 1 // 重置到第一页
    loadVectors(filters)
  }

  const handleLogLoadMore = async () => {
    const nextPage = logPagination.value.page + 1
    await loadVectorLogs(selectedFileId.value, nextPage, true) // 传入append=true
  }

  const handleLogRefresh = () => {
    loadVectorLogs(selectedFileId.value)
  }

  const handleClearLogFilter = () => {
    selectedFileId.value = ''
    loadVectorLogs()
    toast.success($t('admin.vectors.logs.actions.clear'))
  }

  const handleViewDetail = (fileId: string) => {
    selectedFileId.value = fileId
    showLogs.value = true
    loadVectorLogs(fileId, 1)
    toast.info(`${$t('admin.vectors.logs.actions.viewDetails')} ${fileId}`)
  }

  const handleRefresh = () => {
    loadVectors(filters)
    loadQdrantStats()
    toast.info($t('admin.vectors.buttons.refresh'))
  }

  const handleVerificationCompleted = () => {
    handleRefresh()
  }

  const handleRetryVector = async (fileId: string) => {
    await retryVector(fileId, () => {
      handleRefresh()
    })
  }

  const handleBatchAction = async (fileIds: string[], action: 'reset' | 'retry' | 'delete') => {
    if (action === 'delete') {
      deleteConfirmData.value = {
        count: fileIds.length,
        expectedText: $t('admin.vectors.buttons.confirm'),
        fileIds: [...fileIds],
      }
      showDeleteConfirm.value = true
    } else {
      batchConfirmData.value = {
        action,
        message: `${$t('admin.vectors.batchOperations.confirm.delete')} ${fileIds.length} ${action === 'retry' ? $t('admin.vectors.actions.search') : $t('admin.vectors.actions.rebuild')} ${$t('admin.vectors.vectors.title')}。`,
        fileIds: [...fileIds],
      }
      showBatchConfirm.value = true
    }
  }

  const handleRegenerateAll = async () => {
    regenerateConfirmData.value = {
      message: [
        `⚠️ ${$t('admin.vectors.confirm.reindex.title')}`,
        '',
        $t('admin.vectors.confirm.reindex.message'),
        $t('admin.vectors.confirm.reindex.stepsIntro'),
        `1. ${$t('admin.vectors.confirm.delete.message')}`,
        `2. ${$t('admin.vectors.actions.rebuild')}${$t('admin.vectors.vectors.title')}`,
        `3. ${$t('admin.vectors.loading.processing')}`,
        '',
        `${$t('admin.vectors.confirm.reindex.message')}?`,
      ],
    }
    showRegenerateConfirm.value = true
  }

  const handleRetryAllFailed = async () => {
    const webSocketStore = useWebSocketStore()
    const failedCount = webSocketStore.vectorStats?.queue_stats.failed_count || 0

    if (failedCount === 0) {
      toast.info($t('admin.vectors.toast.noData'))
      return
    }

    retryFailedConfirmData.value = {
      message: `${$t('admin.vectors.batchOperations.confirm.reindex')} ${failedCount} ${$t('admin.vectors.status.error')} ${$t('admin.vectors.vectors.title')}。${$t('admin.vectors.confirm.reindex.message')}`,
      count: failedCount,
    }
    showRetryFailedConfirm.value = true
  }

  const handleDeleteConfirm = (inputValue: string) => {
    if (inputValue === deleteConfirmData.value.expectedText) {
      showDeleteConfirm.value = false
      handleBatchActionExecute(deleteConfirmData.value.fileIds, 'delete', handleRefresh)
    } else {
      toast.error($t('admin.vectors.toast.confirmFailed'))
    }
  }

  const handleBatchConfirm = () => {
    showBatchConfirm.value = false
    handleBatchActionExecute(batchConfirmData.value.fileIds, batchConfirmData.value.action, handleRefresh)
  }

  const handleRegenerateConfirm = async () => {
    try {
      showRegenerateConfirm.value = false
      const res = await regenerateAllVectors({})
      if (res.success) {
        toast.success($t('admin.vectors.toast.reindexed'))
        setTimeout(handleRefresh, 2000)
      } else {
        toast.error(res.message || $t('admin.vectors.toast.reindexFailed'))
      }
    } catch (e: any) {
      toast.error(e?.message || $t('admin.vectors.toast.reindexFailed'))
    }
  }

  const handleRetryFailedConfirm = async () => {
    try {
      showRetryFailedConfirm.value = false
      const res = await retryAllFailedVectors()
      if (res.success) {
        toast.success($t('admin.vectors.toast.batchReindexed'))
        setTimeout(handleRefresh, 1500)
      } else {
        toast.error(res.message || $t('admin.vectors.toast.reindexFailed'))
      }
    } catch (e: any) {
      toast.error(e?.message || $t('admin.vectors.toast.reindexFailed'))
    }
  }

  const handleDialogCancel = () => {
    showDeleteConfirm.value = false
    showBatchConfirm.value = false
    showRegenerateConfirm.value = false
    showRetryFailedConfirm.value = false
    toast.info($t('admin.vectors.buttons.cancel'))
  }

  onMounted(async () => {
    try {
      await Promise.allSettled([loadAvailableModels(), loadVectors(filters), loadQdrantStats(), loadVectorLogs()])

      try {
        const { data } = await getSettings({ group: 'vector' })
        const map = Object.fromEntries((data.settings || []).map((s: Setting) => [s.key, s.value])) as Record<string, unknown>
        if (typeof map.vector_auto_processing_enabled === 'boolean') autoEnabled.value = map.vector_auto_processing_enabled
        if (typeof map.vector_concurrency === 'number') concurrency.value = map.vector_concurrency
      } catch {}

      try {
        const st = await getAutoProcessing()
        if (st.success) {
          autoEnabled.value = !!st.data?.enabled
        }
        const ws = useWebSocketStore()
        const rt = ws.vectorStats?.runtime
        if (rt && typeof rt.configured_concurrency === 'number') {
          concurrency.value = rt.configured_concurrency
        }
      } catch {}
    } catch {
      toast.error($t('admin.vectors.error.networkError'))
    }
  })

  onUnmounted(() => {
    cleanupData()
  })
</script>

<style scoped lang="scss">
  .admin-vectors-page {
    color: var(--color-content);
  }

  .module-wrapper {
    display: flex;
    flex-direction: column;
    gap: var(--space-md);
    margin-bottom: var(--space-2xl); // 增加间距
  }

  .module-wrapper-last {
    margin-bottom: var(--space-4xl); // 处理日志底部额外空间
  }

  .module-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: var(--space-sm) 0;
  }

  .module-title {
    font-size: var(--text-lg);
    font-weight: var(--font-semibold);
    color: var(--color-brand-500);
    margin: 0;
    display: flex;
    align-items: center;
    text-shadow: 0 0 8px rgba(var(--color-brand-500-rgb), 0.3);
    letter-spacing: 0.5px;

    i {
      filter: drop-shadow(0 0 4px rgba(var(--color-brand-500-rgb), 0.5));
    }

    .icon {
      margin-right: 0.5rem;
      filter: drop-shadow(0 0 4px rgba(var(--color-brand-500-rgb), 0.5));
    }
  }

  .module-content {
    flex: 1;
  }

  .qdrant-status-section {
    margin-bottom: var(--space-md);
  }
</style>
