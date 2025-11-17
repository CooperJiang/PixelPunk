<script setup lang="ts">
  import { getCurrentLocale } from '@/utils/locale'
  import { computed, h, onMounted, reactive, ref } from 'vue'
  import { shareApi } from '@/api'
  import { useToast } from '@/components/Toast/useToast'
  import { useTexts } from '@/composables/useTexts'
  import type { AdminShareListParams, ShareInfo, ShareStatus } from '@/api/share/types'
  import type { CyberTableColumn } from '@/components/Table/types'
  import ShareFilterPanel from './ShareFilterPanel.vue'

  const toast = useToast()
  const { $t } = useTexts()

  /* 筛选相关状态 */
  const showFilter = ref(false)
  const currentFilters = reactive<AdminShareListParams>({
    page: 1,
    size: 10,
    status: undefined,
    keyword: '',
    start_date: '',
    end_date: '',
    order_by: 'created_at DESC',
  })

  const hasActiveFilters = computed(() =>
    Boolean(currentFilters.status !== undefined || currentFilters.keyword || currentFilters.start_date || currentFilters.end_date)
  )

  const handleFilter = () => {
    showFilter.value = !showFilter.value
  }

  const applyFilter = (filters: AdminShareListParams) => {
    Object.assign(currentFilters, filters)
    Object.assign(searchParams, filters)
    loadShareList()
    showFilter.value = false
  }

  const resetFilters = () => {
    Object.assign(currentFilters, {
      page: 1,
      size: 10,
      status: undefined,
      keyword: '',
      start_date: '',
      end_date: '',
      order_by: 'created_at DESC',
    })
    Object.assign(searchParams, currentFilters)
    loadShareList()
  }

  const tableColumns = computed<CyberTableColumn<ShareInfo>[]>(() => [
    {
      key: 'name',
      dataIndex: 'name',
      title: $t('admin.shares.shareList.columns.name'),
      width: 150,
      ellipsis: true,
      align: 'center',
      render: (value: string, record: ShareInfo) =>
        h(
          'div',
          {
            class: 'share-name',
            title: record.name,
          },
          record.name
        ),
    },
    {
      key: 'share_key',
      dataIndex: 'share_key',
      title: $t('admin.shares.shareList.columns.shareKey'),
      width: 200,
      ellipsis: true,
      align: 'center',
      render: (value: string, record: ShareInfo) =>
        h('div', { class: 'share-key', title: record.share_key }, [
          h('span', { class: 'key-text' }, record.share_key),
          h(
            'button',
            {
              class:
                'inline-flex items-center justify-center w-6 h-6 ml-2 text-xs bg-cyan-500/20 hover:bg-cyan-500/30 border border-cyan-500/30 text-cyan-400 rounded transition-all duration-200',
              title: $t('admin.shares.shareList.actions.copyLink'),
              onClick: (e: Event) => {
                e.stopPropagation()
                copyShareUrl(record)
              },
            },
            [h('i', { class: 'fas fa-copy' })]
          ),
        ]),
    },
    {
      key: 'expired_at',
      dataIndex: 'expired_at',
      title: $t('admin.shares.shareList.columns.expiredAt'),
      width: 140,
      align: 'center',
      render: (value: string) => formatDate(value),
    },
    {
      key: 'views',
      dataIndex: 'current_views',
      title: $t('admin.shares.shareList.columns.views'),
      width: 110,
      align: 'center',
      render: (value: number, record: ShareInfo) =>
        `${record.current_views} / ${record.max_views || $t('admin.shares.shareList.format.infiniteViews')}`,
    },
    {
      key: 'created_at',
      dataIndex: 'created_at',
      title: $t('admin.shares.shareList.columns.createdAt'),
      width: 140,
      align: 'center',
      render: (value: string) => formatDate(value),
    },
    {
      key: 'status',
      dataIndex: 'status',
      title: $t('admin.shares.shareList.columns.status'),
      width: 80,
      align: 'center',
      render: (value: ShareStatus, record: ShareInfo) =>
        h(
          'div',
          {
            class: ['status-tag', getStatusClass(record.status)],
          },
          getStatusText(record.status)
        ),
    },
    {
      key: 'actions',
      title: $t('admin.shares.shareList.columns.actions'),
      width: 180,
      align: 'center',
      render: (_value: unknown, record: ShareInfo) =>
        h('div', { class: 'action-buttons flex gap-1 justify-center items-center' }, [
          h(
            'button',
            {
              class:
                'inline-flex items-center justify-center w-7 h-7 text-xs bg-blue-500/20 hover:bg-blue-500/30 border border-blue-500/30 text-blue-400 rounded transition-all duration-200',
              title: $t('admin.shares.shareList.actions.viewContent'),
              onClick: (e: Event) => {
                e.stopPropagation()
                viewShareContent(record)
              },
            },
            [h('i', { class: 'fas fa-eye' })]
          ),

          ...(record.status === 1
            ? [
                h(
                  'button',
                  {
                    class:
                      'inline-flex items-center justify-center w-7 h-7 text-xs bg-yellow-500/20 hover:bg-yellow-500/30 border border-yellow-500/30 text-yellow-400 rounded transition-all duration-200',
                    title: $t('admin.shares.shareList.actions.disable'),
                    onClick: (e: Event) => {
                      e.stopPropagation()
                      updateStatus(record.id, 4)
                    },
                  },
                  [h('i', { class: 'fas fa-ban' })]
                ),
              ]
            : []),

          ...(record.status === 4
            ? [
                h(
                  'button',
                  {
                    class:
                      'inline-flex items-center justify-center w-7 h-7 text-xs bg-green-500/20 hover:bg-green-500/30 border border-green-500/30 text-green-400 rounded transition-all duration-200',
                    title: $t('admin.shares.shareList.actions.enable'),
                    onClick: (e: Event) => {
                      e.stopPropagation()
                      updateStatus(record.id, 1)
                    },
                  },
                  [h('i', { class: 'fas fa-check' })]
                ),
              ]
            : []),

          h(
            'button',
            {
              class:
                'inline-flex items-center justify-center w-7 h-7 text-xs bg-red-500/20 hover:bg-red-500/30 border border-red-500/30 text-red-400 rounded transition-all duration-200',
              title: $t('admin.shares.shareList.actions.delete'),
              onClick: (e: Event) => {
                e.stopPropagation()
                confirmDelete(record)
              },
            },
            [h('i', { class: 'fas fa-trash' })]
          ),
        ]),
    },
  ])

  const searchParams = reactive<AdminShareListParams>({
    page: 1,
    size: 10,
    status: undefined,
    keyword: '',
    start_date: '',
    end_date: '',
    order_by: 'created_at DESC',
  })

  const totalItems = ref(0)

  const shareList = ref<ShareInfo[]>([])

  const loading = ref(false)

  const statusDialog = reactive({
    visible: false,
    shareId: '',
    status: 1 as ShareStatus,
    reason: '',
    title: $t('admin.shares.statusDialog.titleUpdate'),
    loading: false,
  })

  const deleteDialog = reactive({
    visible: false,
    shareId: '',
    shareName: '',
    force: false,
    loading: false,
  })

  const loadShareList = async () => {
    try {
      loading.value = true

      const queryParams: AdminShareListParams = {
        page: searchParams.page,
        size: searchParams.size,
        order_by: searchParams.order_by,
      }

      if (searchParams.status !== undefined) {
        queryParams.status = searchParams.status
      }

      if (searchParams.keyword) {
        queryParams.keyword = searchParams.keyword
      }

      if (searchParams.start_date) {
        queryParams.start_date = searchParams.start_date
      }

      if (searchParams.end_date) {
        queryParams.end_date = searchParams.end_date
      }

      const result = await shareApi.adminGetShareList(queryParams)

      if (result.success && result.data) {
        shareList.value = result.data.list || []
        totalItems.value = result.data.total || 0
      }
    } catch (error: unknown) {
      toast.error((error as any)?.message || $t('admin.shares.shareList.messages.fetchFailed'))
    } finally {
      loading.value = false
    }
  }

  const refreshData = () => {
    loadShareList()
  }

  const handlePageChange = (page: number) => {
    searchParams.page = page
    loadShareList()
  }

  const handlePageSizeChange = (size: number) => {
    searchParams.size = size
    searchParams.page = 1 // 重置到第一页
    loadShareList()
  }

  const viewShareContent = async (share: ShareInfo) => {
    try {
      loading.value = true

      const result = await shareApi.adminGetShareAccessToken(share.id)

      const baseUrl = window.location.origin
      let shareUrl = `${baseUrl}/share/${share.share_key}`

      if (result.success && result.data && result.data.access_token) {
        const accessToken = result.data.access_token
        shareUrl += `?access_token=${accessToken}`
      }

      window.open(shareUrl, '_blank')
    } catch {
      const baseUrl = window.location.origin
      const shareUrl = `${baseUrl}/share/${share.share_key}`
      window.open(shareUrl, '_blank')
      toast.warning($t('admin.shares.shareList.messages.getAccessFailed'))
    } finally {
      loading.value = false
    }
  }

  const updateStatus = (id: string, status: ShareStatus) => {
    statusDialog.shareId = id
    statusDialog.status = status
    statusDialog.reason = ''
    statusDialog.title =
      status === 1 ? $t('admin.shares.shareList.statusDialog.titleEnable') : $t('admin.shares.shareList.statusDialog.titleUpdate')
    statusDialog.visible = true
  }

  const confirmUpdateStatus = async () => {
    try {
      statusDialog.loading = true

      const params = {
        status: statusDialog.status,
        ...(statusDialog.reason ? { reason: statusDialog.reason } : {}),
      }

      await shareApi.adminUpdateShareStatus(statusDialog.shareId, params)

      toast.success($t('admin.shares.shareList.messages.updateStatusSuccess'))
      statusDialog.visible = false
      loadShareList()
    } catch (error: unknown) {
      toast.error((error as any)?.message || $t('admin.shares.shareList.messages.updateStatusFailed'))
    } finally {
      statusDialog.loading = false
    }
  }

  const confirmDelete = (share: ShareInfo) => {
    deleteDialog.shareId = share.id
    deleteDialog.shareName = share.name
    deleteDialog.force = false
    deleteDialog.visible = true
  }

  const confirmDeleteShare = async () => {
    try {
      deleteDialog.loading = true
      await shareApi.adminDeleteShare(deleteDialog.shareId, deleteDialog.force)

      toast.success($t('admin.shares.shareList.messages.deleteSuccess'))
      deleteDialog.visible = false
      loadShareList()
    } catch (error: unknown) {
      toast.error((error as any)?.message || $t('admin.shares.shareList.messages.deleteFailed'))
    } finally {
      deleteDialog.loading = false
    }
  }

  const copyShareUrl = async (share: ShareInfo) => {
    try {
      const result = await shareApi.adminGetShareAccessToken(share.id)

      const baseUrl = window.location.origin
      let shareUrl = `${baseUrl}/share/${share.share_key}`

      if (result.success && result.data && result.data.access_token) {
        shareUrl += `?access_token=${result.data.access_token}`
      }

      navigator.clipboard
        .writeText(shareUrl)
        .then(() => {
          toast.success($t('admin.shares.shareList.messages.copySuccess'))
        })
        .catch(() => {
          toast.error($t('admin.shares.shareList.messages.copyFailed'))
        })
    } catch {
      const baseUrl = window.location.origin
      const shareUrl = `${baseUrl}/share/${share.share_key}`

      navigator.clipboard
        .writeText(shareUrl)
        .then(() => {
          toast.success($t('admin.shares.shareList.messages.copySuccessNormal'))
        })
        .catch(() => {
          toast.error($t('admin.shares.shareList.messages.copyFailed'))
        })
    }
  }

  const formatDate = (dateString?: string) => {
    if (!dateString) {
      return $t('admin.shares.shareList.format.unlimited')
    }

    const date = new Date(dateString)
    return date.toLocaleString(getCurrentLocale(), {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit',
    })
  }

  const getStatusText = (status?: ShareStatus) => {
    switch (status) {
      case 1:
        return $t('admin.shares.shareList.status.normal')
      case 2:
        return $t('admin.shares.shareList.status.expired')
      case 3:
        return $t('admin.shares.shareList.status.deleted')
      case 4:
        return $t('admin.shares.shareList.status.disabled')
      default:
        return $t('admin.shares.shareList.status.unknown')
    }
  }

  const getStatusClass = (status?: ShareStatus) => {
    switch (status) {
      case 1:
        return 'status-normal'
      case 2:
        return 'status-expired'
      case 3:
        return 'status-deleted'
      case 4:
        return 'status-disabled'
      default:
        return ''
    }
  }

  const handleRowClick = (_row: ShareInfo, _index: number) => {}

  onMounted(() => {
    loadShareList()
  })
</script>

<template>
  <div class="share-list-component">
    <div class="action-bar">
      <div class="action-buttons">
        <CyberIconButton type="primary" size="small" :title="$t('admin.shares.shareList.buttons.refresh')" @click="refreshData">
          <i class="fas fa-sync-alt" />
        </CyberIconButton>

        <CyberIconButton
          type="secondary"
          size="small"
          :title="$t('admin.shares.shareList.buttons.filter')"
          :class="{ 'filter-active': showFilter || hasActiveFilters }"
          @click="handleFilter"
        >
          <i class="fas fa-filter" />
        </CyberIconButton>
      </div>
    </div>

    <ShareFilterPanel v-show="showFilter" :initial-filters="currentFilters" @filter="applyFilter" />

    <div v-if="loading" class="admin-loading">
      <CyberLoading :visible="true" :text="$t('admin.shares.shareList.loading.text')" />
    </div>

    <div v-else-if="!loading && shareList.length === 0" class="empty-state">
      <div class="empty-content animate-slide-up">
        <div class="empty-icon-wrapper">
          <i class="fas fa-share-alt empty-icon" />
        </div>
        <h3 class="empty-title">{{ $t('admin.shares.shareList.empty.title') }}</h3>
        <p class="empty-desc">{{ $t('admin.shares.shareList.empty.description') }}</p>
        <CyberButton v-if="hasActiveFilters" type="secondary" @click="resetFilters">
          <i class="fas fa-redo mr-2" />
          {{ $t('admin.shares.shareList.empty.resetButton') }}
        </CyberButton>
      </div>
    </div>

    <div v-else class="table-container">
      <div class="table-wrapper">
        <CyberDataTable
          :columns="tableColumns"
          :data="shareList"
          :loading="false"
          :hoverable="true"
          :striped="true"
          :bordered="true"
          row-key="id"
          @row-click="handleRowClick"
        />
      </div>
    </div>

    <div v-if="!loading && totalItems > 0" class="pagination-container">
      <CyberPagination
        v-model:current-page="searchParams.page"
        :total="totalItems"
        :page-size="searchParams.size"
        :show-page-size-selector="true"
        :show-quick-jumper="true"
        @update:current-page="handlePageChange"
        @update:page-size="handlePageSizeChange"
      />
    </div>

    <CyberDialog v-model="statusDialog.visible" :title="statusDialog.title" width="400px">
      <div class="status-dialog-content">
        <p>
          {{ $t('admin.shares.shareList.statusDialog.messagePrefix') }}
          <span class="highlight">{{ getStatusText(statusDialog.status) }}</span
          >{{ $t('admin.shares.shareList.statusDialog.messageSuffix') }}
        </p>

        <div v-if="statusDialog.status !== 1" class="form-item">
          <label>{{ $t('admin.shares.shareList.statusDialog.reasonLabel') }}</label>
          <CyberInput
            v-model="statusDialog.reason"
            type="textarea"
            :placeholder="$t('admin.shares.shareList.statusDialog.reasonPlaceholder')"
            :rows="3"
          />
        </div>
      </div>

      <template #footer>
        <div class="dialog-footer">
          <CyberButton @click="statusDialog.visible = false">{{ $t('admin.shares.shareList.statusDialog.cancel') }}</CyberButton>
          <CyberButton type="primary" :loading="statusDialog.loading" @click="confirmUpdateStatus">{{
            $t('admin.shares.shareList.statusDialog.confirm')
          }}</CyberButton>
        </div>
      </template>
    </CyberDialog>

    <CyberDialog v-model="deleteDialog.visible" :title="$t('admin.shares.shareList.deleteDialog.title')" width="400px">
      <div class="delete-dialog-content">
        <p>
          {{ $t('admin.shares.shareList.deleteDialog.messagePrefix') }} <span class="highlight">{{ deleteDialog.shareName }}</span
          >{{ $t('admin.shares.shareList.deleteDialog.messageSuffix') }}
        </p>

        <div class="form-item force-delete">
          <CyberCheckbox v-model="deleteDialog.force">{{ $t('admin.shares.shareList.deleteDialog.forceDelete') }}</CyberCheckbox>
        </div>
      </div>

      <template #footer>
        <div class="dialog-footer p-4">
          <CyberButton type="warning" @click="deleteDialog.visible = false">{{
            $t('admin.shares.shareList.deleteDialog.cancel')
          }}</CyberButton>
          <CyberButton type="danger" :loading="deleteDialog.loading" @click="confirmDeleteShare">{{
            $t('admin.shares.shareList.deleteDialog.confirm')
          }}</CyberButton>
        </div>
      </template>
    </CyberDialog>
  </div>
</template>

<style scoped lang="scss">
  .share-list-component {
    display: flex;
    flex-direction: column;
    gap: var(--space-md);
    flex: 1;
    min-height: 0;
  }

  .action-bar {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .action-buttons {
    display: flex;
    gap: var(--space-sm);
  }

  .filter-active {
    background-color: rgba(var(--color-brand-500-rgb), 0.2) !important;
    color: var(--color-brand-500) !important;
    box-shadow: var(--shadow-glow-sm) !important;
  }

  .admin-loading {
    display: flex;
    align-items: center;
    justify-content: center;
    padding: var(--space-3xl) var(--space-md);
    min-height: 100%;
    height: 100%;
    flex: 1;
  }

  .empty-state {
    flex: 1;
    display: flex;
    align-items: center;
    justify-content: center;
    min-height: 100%;
    height: 100%;
    padding: var(--space-3xl) var(--space-2xl);
  }

  .empty-content {
    text-align: center;
    max-width: 400px;
  }

  .empty-icon-wrapper {
    margin-bottom: var(--space-xl);
    position: relative;
    display: inline-block;
  }

  .empty-icon {
    font-size: var(--text-4xl);
    color: var(--color-brand-500);
    filter: drop-shadow(var(--shadow-glow-md));
    animation: float 3s var(--ease-in-out) infinite;
  }

  .empty-title {
    font-size: var(--text-lg);
    font-weight: var(--font-semibold);
    color: var(--color-content-default);
    margin-bottom: var(--space-sm);
    letter-spacing: var(--tracking-tight);
  }

  .empty-desc {
    font-size: var(--text-sm);
    color: var(--color-content-muted);
    line-height: var(--leading-relaxed);
    margin-bottom: var(--space-xl);
  }

  .table-container {
    flex: 1;
    min-height: 0;
    display: flex;
    flex-direction: column;
    border-radius: var(--radius-sm);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.1);
    overflow: hidden;
  }

  .table-wrapper {
    flex: 1;
    min-height: 0;
    background: rgba(var(--color-background-800-rgb), 0.3);
    overflow-y: auto;
    overflow-x: hidden;
  }

  .table-wrapper::-webkit-scrollbar {
    width: var(--space-sm);
    background: transparent;
  }

  .table-wrapper::-webkit-scrollbar-track {
    background: rgba(var(--color-brand-500-rgb), 0.1);
    border-radius: var(--radius-sm);
  }

  .table-wrapper::-webkit-scrollbar-thumb {
    background: rgba(var(--color-brand-500-rgb), 0.3);
    border-radius: var(--radius-sm);
    transition: all var(--transition-normal);
  }

  .table-wrapper::-webkit-scrollbar-thumb:hover {
    background: rgba(var(--color-brand-500-rgb), 0.5);
  }

  .share-name {
    max-width: 180px;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .share-key {
    display: flex;
    align-items: center;
    gap: var(--space-sm);
    font-family: var(--font-mono);
    color: var(--color-brand-500);
    min-width: 0;
  }

  .key-text {
    flex: 1;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    min-width: 0;
  }

  .copy-button {
    flex-shrink: 0;
    margin-left: var(--space-sm);
  }

  .status-tag {
    display: inline-block;
    padding: var(--space-xs) var(--space-sm);
    border-radius: var(--radius-sm);
    font-size: var(--text-sm);
    font-weight: var(--font-medium);
  }

  .status-normal {
    background: rgba(var(--color-success-rgb), 0.15);
    color: var(--color-success-400);
  }

  .status-expired {
    background: rgba(var(--color-warning-rgb), 0.15);
    color: var(--color-warning-400);
  }

  .status-deleted {
    background: rgba(var(--color-error-rgb), 0.15);
    color: var(--color-error-400);
  }

  .status-disabled {
    background: rgba(var(--color-content-muted-rgb), 0.15);
    color: var(--color-content-muted);
  }

  .pagination-container {
    margin-top: var(--space-xl);
    display: flex;
    justify-content: flex-end;
  }

  .status-dialog-content,
  .delete-dialog-content {
    padding: var(--space-sm) 0;
  }

  .highlight {
    color: var(--color-brand-500);
    font-weight: var(--font-semibold);
  }

  .form-item {
    margin-top: var(--space-md);
  }

  .force-delete {
    margin-top: var(--space-xl);
    color: var(--color-error-400);
  }

  .dialog-footer {
    display: flex;
    justify-content: flex-end;
    gap: var(--space-sm);
    padding-top: var(--space-md);
  }
</style>
