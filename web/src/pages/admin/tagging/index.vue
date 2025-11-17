<script setup lang="ts">
  import { computed, onMounted, onUnmounted, reactive, ref } from 'vue'
  import { useToast } from '@/components/Toast/useToast'
  import { useWebSocketStore } from '@/store/websocket'
  import { useTexts } from '@/composables/useTexts'

  const { $t } = useTexts()
  import {
    getAutoProcessingStatus,
    getTaggingStats,
    getTaggingLogs,
    getTaggingStatusList,
    ignoreTagging,
    resetStuckTagging,
    retryFailedAll,
    retryTagging,
    setAutoProcessing,
    setConcurrency,
    triggerTagging,
    unignoreTagging,
    type StatusStatsItem,
    type TaggingLogItem,
  } from '@/api/admin/tagging'
  import { batchUpsertSettings, getSettings, type Setting } from '@/api/admin/settings'
  import TaggingLogList from './components/TaggingLogList.vue'
  import TaggingFileList from './components/TaggingFileList.vue'
  import TaggingQueueStatus from './components/TaggingQueueStatus.vue'
  import TaggingFilterPanel from './components/TaggingFilterPanel.vue'

  const toast = useToast()

  const webSocketStore = useWebSocketStore()

  const isLoading = ref(false)
  const isProcessing = ref(false)
  const showSkeleton = ref(false)
  const logs = ref<TaggingLogItem[]>([])
  const logPagination = ref({
    page: 1,
    limit: 20,
    total: 0,
  })
  const selectedFileId = ref<string>('')
  const images = ref<Record<string, unknown>[]>([])
  const total = ref(0)
  const currentPage = ref(1)
  const pageSize = ref(10)
  const selectedIds = ref<string[]>([])
  const _selectedCount = computed(() => selectedIds.value.length)
  const showFilter = ref(false)
  const filters = reactive({
    status: '',
    page: 1,
    limit: 10,
    order_by: 'updated_at',
    order: 'desc',
  })
  const autoEnabled = ref<boolean>(true)
  const updatingAuto = ref<boolean>(false)
  const concurrency = ref<number>(0)

  const fallbackStats = ref<StatusStatsItem[]>([])

  const stats = computed<StatusStatsItem[]>(() => {
    if (webSocketStore.queueStats) {
      const queueStats = webSocketStore.queueStats.queue_stats
      return [
        { status: 'none', count: queueStats.none_count || 0 },
        { status: 'pending', count: queueStats.pending_count || 0 },
        { status: 'done', count: queueStats.done_count || 0 },
        { status: 'failed', count: queueStats.failed_count || 0 },
        { status: 'ignored', count: queueStats.ignored_count || 0 },
        { status: 'skipped', count: 0 }, // WebSocket数据中没有skipped
      ]
    }

    return fallbackStats.value
  })

  const loadTaggingStatus = async () => {
    try {
      isLoading.value = true
      showSkeleton.value = true
      const result = await getTaggingStatusList({ ...filters, page: currentPage.value })
      if (result.success) {
        images.value = result.data?.data || []
        total.value = result.data?.total || 0
      }
    } catch (e: unknown) {
      toast.error((e as Error)?.message || $t('admin.tagging.messages.loadFailed'))
    } finally {
      isLoading.value = false
      showSkeleton.value = false
    }
  }

  const loadFallbackTaggingStats = async () => {
    try {
      const result = await getTaggingStats()

      if (result.success) {
        const { data } = result
        if (!webSocketStore.isConnected || !webSocketStore.queueStats) {
          fallbackStats.value = data?.status_stats || []
        }
      }
    } catch (error: unknown) {
      toast.error(error.message)
    }
  }

  const loadTaggingLogs = async (fileID = '', page = 1, append = false) => {
    try {
      const result = await getTaggingLogs({
        file_id: fileID,
        page,
        limit: logPagination.value.limit,
      })

      if (result.success && result.data) {
        const { logs: newLogs, pagination } = result.data

        if (append) {
          logs.value = [...logs.value, ...newLogs]
        } else {
          logs.value = newLogs
        }

        logPagination.value = {
          page: pagination.page,
          limit: pagination.limit,
          total: pagination.total,
        }
      }
    } catch (error: unknown) {
      toast.error((error as Error)?.message || $t('admin.tagging.messages.loadLogsFailed'))
    }
  }

  const handleLogLoadMore = async () => {
    const nextPage = logPagination.value.page + 1
    await loadTaggingLogs(selectedFileId.value, nextPage, true)
  }

  const handleLogRefresh = async () => {
    logPagination.value.page = 1
    await loadTaggingLogs(selectedFileId.value, 1, false)
  }

  const handleClearLogFilter = async () => {
    selectedFileId.value = ''
    logPagination.value.page = 1
    await loadTaggingLogs('', 1, false)
  }

  const reloadData = async () => {
    await Promise.all([loadTaggingStatus(), loadFallbackTaggingStats(), loadTaggingLogs()])
  }

  const loadAiSettings = async () => {
    try {
      const { data } = await getSettings({ group: 'ai' })
      const map = Object.fromEntries((data.settings || []).map((s: Setting) => [s.key, s.value])) as Record<string, unknown>
      if (typeof map.ai_concurrency === 'number') concurrency.value = map.ai_concurrency
      if (typeof map.ai_auto_processing_enabled === 'boolean') autoEnabled.value = map.ai_auto_processing_enabled
    } catch {}
  }

  const handleStatusSelect = (status: string) => {
    filters.status = status
    currentPage.value = 1
    selectedIds.value = []
    loadTaggingStatus()
  }

  const handleFilter = (newFilters: Record<string, unknown>) => {
    Object.assign(filters, newFilters)
    currentPage.value = 1
    selectedIds.value = []
    loadTaggingStatus()
    showFilter.value = false
  }

  const handlePageChange = (page: number) => {
    currentPage.value = page
    loadTaggingStatus()
  }

  const handlePageSizeChange = (size: number) => {
    pageSize.value = size
    currentPage.value = 1
    filters.page = 1
    filters.limit = size
    selectedIds.value = []
    loadTaggingStatus()
  }

  const handleRetryTagging = async (fileIds: string[]) => {
    if (fileIds.length === 0) return
    try {
      isProcessing.value = true
      const res = await retryTagging(fileIds)
      if (res.success) {
        toast.success(
          $t('admin.tagging.messages.retrySuccess')
            .replace('{selected}', String(res.data?.requested ?? 0))
            .replace('{enqueued}', String(res.data?.enqueued ?? 0))
            .replace('{skipped}', String(res.data?.skipped ?? 0))
        )
        await reloadData()
      }
    } catch (e: unknown) {
      toast.error((e as Error)?.message || $t('admin.tagging.messages.retryFailed'))
    } finally {
      isProcessing.value = false
    }
  }

  const handleResetStuck = async (timeThreshold: number) => {
    try {
      isProcessing.value = true
      const result = await resetStuckTagging(timeThreshold)

      if (result.success) {
        const { data } = result
        if (data?.reset_count > 0) {
          toast.success($t('admin.tagging.messages.resetSuccess').replace('{count}', String(data.reset_count)))
        } else {
          toast.info($t('admin.tagging.messages.resetEmpty'))
        }
        await reloadData()
      }
    } catch (error: unknown) {
      toast.error(error.message)
    } finally {
      isProcessing.value = false
    }
  }

  const handleRetryAllFailed = async () => {
    try {
      isProcessing.value = true
      const res = await retryFailedAll(500)
      if (res.success) {
        toast.success(
          $t('admin.tagging.messages.retrySuccess')
            .replace('{selected}', String(res.data?.selected ?? 0))
            .replace('{enqueued}', String(res.data?.enqueued ?? 0))
            .replace('{skipped}', String(res.data?.skipped ?? 0))
        )
        await reloadData()
      }
    } catch (error: unknown) {
      toast.error(error.message)
    } finally {
      isProcessing.value = false
    }
  }

  const handleTriggerTagging = async () => {
    try {
      isProcessing.value = true
      const res = await triggerTagging(200)
      if (res.success) {
        const count = res.data?.submitted_count ?? 0
        if (count > 0) {
          toast.success($t('admin.tagging.messages.triggerSuccess').replace('{count}', String(count)))
        } else {
          toast.info($t('admin.tagging.messages.triggerEmpty'))
        }
        await reloadData()
      }
    } catch (error: unknown) {
      toast.error(error.message)
    } finally {
      isProcessing.value = false
    }
  }

  const handleToggleSelect = (id: string) => {
    const set = new Set(selectedIds.value)
    if (set.has(id)) set.delete(id)
    else set.add(id)
    selectedIds.value = Array.from(set)
  }

  const handleToggleSelectAll = (checked: boolean) => {
    const pageIds = images.value.map((i) => i.id as string)
    const set = new Set(selectedIds.value)
    if (checked) pageIds.forEach((id) => set.add(id))
    else pageIds.forEach((id) => set.delete(id))
    selectedIds.value = Array.from(set)
  }

  const handleIgnoreSelected = async (reason = '') => {
    if (selectedIds.value.length === 0) return
    const ids = selectedIds.value.slice(0, 500)
    try {
      isProcessing.value = true
      const res = await ignoreTagging(ids, reason)
      if (res.success) {
        toast.success($t('admin.tagging.messages.ignoreSuccess').replace('{count}', String(res.data?.count ?? ids.length)))
        selectedIds.value = []
        await reloadData()
      }
    } catch (e: unknown) {
      toast.error((e as Error)?.message || $t('admin.tagging.messages.ignoreFailed'))
    } finally {
      isProcessing.value = false
    }
  }

  const handleUnignoreSelected = async () => {
    if (selectedIds.value.length === 0) return
    const ids = selectedIds.value.slice(0, 500)
    try {
      isProcessing.value = true
      const res = await unignoreTagging(ids)
      if (res.success) {
        if (res.data?.enqueued !== undefined && res.data?.skipped !== undefined) {
          toast.success(
            $t('admin.tagging.messages.unignoreSuccessDetail')
              .replace('{count}', String(res.data?.count ?? ids.length))
              .replace('{enqueued}', String(res.data?.enqueued ?? 0))
              .replace('{skipped}', String(res.data?.skipped ?? 0))
          )
        } else {
          toast.success($t('admin.tagging.messages.unignoreSuccess').replace('{count}', String(res.data?.count ?? ids.length)))
        }
        selectedIds.value = []
        await reloadData()
      }
    } catch (e: unknown) {
      toast.error((e as Error)?.message || $t('admin.tagging.messages.unignoreFailed'))
    } finally {
      isProcessing.value = false
    }
  }

  onMounted(async () => {
    await Promise.all([reloadData(), loadAiSettings()])
    try {
      const st = await getAutoProcessingStatus()
      if (st.success) {
        autoEnabled.value = !!st.data?.enabled
        const cfg = webSocketStore.queueStats?.runtime || webSocketStore.queueStats?.config
        if (cfg) {
          const v = (cfg as Record<string, unknown>).configured_concurrency || (cfg as Record<string, unknown>).max_concurrency
          if (typeof v === 'number') concurrency.value = v
        }
      }
    } catch {}
  })

  onUnmounted(() => {})
</script>

<template>
  <div class="admin-tagging-page admin-page-container">
    <CyberAdminWrapper :title="$t('admin.tagging.title')" :subtitle="$t('admin.tagging.subtitle')" icon="fas fa-tags">
      <template #stats>
        <span v-if="stats.length > 0">
          {{ $t('admin.tagging.stats.total').replace('{count}', String(stats.reduce((sum, item) => sum + item.count, 0))) }}
        </span>
      </template>

      <template #actions>
        <div class="flex items-center gap-2">
          <span class="text-xs">{{ $t('admin.tagging.actions.autoProcessing') }}</span>
          <CyberSwitch
            v-model="autoEnabled"
            :disabled="updatingAuto"
            @change="
              async (v: boolean) => {
                try {
                  updatingAuto = true
                  await batchUpsertSettings([
                    {
                      key: 'ai_auto_processing_enabled',
                      value: v,
                      type: 'boolean',
                      group: 'ai',
                      description: $t('admin.tagging.settingsDescriptions.autoProcessing'),
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

        <div class="bg-border-default mx-3 h-4 w-px"></div>

        <div class="flex items-center gap-2">
          <span class="whitespace-nowrap text-xs">{{ $t('admin.tagging.concurrency') }}</span>
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
                      key: 'ai_concurrency',
                      value: concurrency,
                      type: 'number',
                      group: 'ai',
                      description: $t('admin.tagging.settingsDescriptions.concurrency'),
                    },
                  ])
                  await setConcurrency(concurrency)
                } catch {}
              }
            "
          />
          <span class="w-6 text-center text-xs">{{ concurrency }}</span>
        </div>
        <div class="bg-border-default mx-3 h-4 w-px"></div>

        <CyberButton type="warning" class="text-xs" :disabled="isProcessing || !autoEnabled" @click="handleTriggerTagging()">
          <i class="fas fa-wrench mr-1.5" />{{ $t('admin.tagging.buttons.fixQueue') }}
        </CyberButton>
        <CyberButton type="primary" class="ml-2 text-xs" :disabled="isProcessing || !autoEnabled" @click="handleRetryAllFailed()">
          <i class="fas fa-redo mr-1.5" />{{ $t('admin.tagging.buttons.retryFailed') }}
        </CyberButton>
        <CyberButton type="secondary" class="ml-2 text-xs" :disabled="isProcessing" @click="handleResetStuck(60)">
          <i class="fas fa-broom mr-1.5" />{{ $t('admin.tagging.buttons.resetStuck') }}
        </CyberButton>
        <CyberButton
          type="outlined"
          class="ml-2 text-xs"
          @click="showFilter = !showFilter"
          :class="{ 'border-brand-500 text-brand-500': showFilter }"
        >
          <i class="fas fa-filter mr-1.5"></i>{{ $t('admin.tagging.buttons.filter') }}
        </CyberButton>
      </template>

      <template #content>
        <div class="tagging-content admin-content-wrapper">
          <div class="module-wrapper">
            <div class="module-header">
              <h3 class="module-title">
                <i class="fas fa-chart-bar mr-2" />
                {{ $t('admin.tagging.queueStatus.title') }}
              </h3>
            </div>
            <div class="module-content">
              <TaggingQueueStatus :selected-status="filters.status" @select="handleStatusSelect" />
            </div>
          </div>

          <div v-if="showFilter" class="filter-section">
            <transition
              name="filter-panel"
              enter-active-class="transition duration-300 ease-out"
              leave-active-class="transition duration-200 ease-in"
              enter-from-class="opacity-0 transform -translate-y-8"
              enter-to-class="opacity-100 transform translate-y-0"
              leave-from-class="opacity-100 transform translate-y-0"
              leave-to-class="opacity-0 transform -translate-y-8"
            >
              <TaggingFilterPanel :initial-filters="filters" @filter="handleFilter" />
            </transition>
          </div>

          <div class="image-list-section">
            <TaggingFileList
              :images="images"
              :is-loading="showSkeleton || isLoading"
              :total="total"
              :current-page="currentPage"
              :page-size="pageSize"
              :processing="isProcessing"
              :selected-ids="selectedIds"
              @retry="handleRetryTagging"
              @page-change="handlePageChange"
              @page-size-change="handlePageSizeChange"
              @toggle-select="handleToggleSelect"
              @toggle-select-all="handleToggleSelectAll"
              @ignore-selected="handleIgnoreSelected"
              @unignore-selected="handleUnignoreSelected"
            />
          </div>

          <div class="module-wrapper module-wrapper-last">
            <div class="module-header">
              <h3 class="module-title">
                <i class="fas fa-history mr-2" />
                {{ $t('admin.tagging.recentLogs.title') }}
              </h3>
            </div>
            <div class="module-content">
              <TaggingLogList
                :logs="logs"
                :pagination="logPagination"
                :selected-file-id="selectedFileId"
                @load-more="handleLogLoadMore"
                @refresh="handleLogRefresh"
                @clear-filter="handleClearLogFilter"
              />
            </div>
          </div>
        </div>
      </template>
    </CyberAdminWrapper>
  </div>
</template>

<style scoped lang="scss">
  .page-title {
    font-size: var(--text-lg);
    color: var(--color-content-heading);
    margin-bottom: var(--space-md);
    padding-bottom: var(--space-sm);
  }

  .tagging-content {
    display: flex;
    flex-direction: column;
    gap: var(--space-xl);
    flex: 1;
    min-height: 0;
  }

  .module-wrapper {
    display: flex;
    flex-direction: column;
    gap: var(--space-md);
  }

  .module-wrapper-last {
    margin-bottom: var(--space-4xl); // 最后一个模块底部额外空间
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

  .module-subtitle {
    font-size: var(--text-sm);
    color: var(--color-content-muted);
    font-weight: var(--font-medium);
  }

  .module-content {
    flex: 1;
  }

  :deep(.form-container) {
    background: var(--color-background-800);
    border: 1px solid var(--color-border-default);
    border-radius: var(--radius-sm);
    padding: var(--space-xl);
    margin-bottom: var(--space-xl);
    box-shadow: var(--shadow-lg);
    backdrop-filter: var(--backdrop-blur-md);
    transition: all var(--transition-fast) var(--ease-out);
  }

  :deep(.form-container):hover {
    border-color: var(--color-border-hover);
    box-shadow: var(--shadow-xl);
  }

  .image-list-section {
    flex: 1;
    min-height: auto;
    display: flex;
    flex-direction: column;
  }

  .stats-text {
    color: var(--color-content-muted);
    font-size: var(--text-sm);
  }
</style>
