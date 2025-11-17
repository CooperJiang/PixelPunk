<script setup lang="ts">
  import { getCurrentLocale } from '@/utils/locale'
  import { computed, h, onMounted, reactive, ref } from 'vue'
  import { useRouter } from 'vue-router'
  import { shareApi } from '@/api'
  import { useToast } from '@/components/Toast/useToast'
  import { useTexts } from '@/composables/useTexts'
  import type { VisitorInfoDetail } from '@/api/share/types'

  const router = useRouter()
  const toast = useToast()
  const { $t } = useTexts()

  /* 表格列定义 */
  const tableColumns = computed(() => [
    {
      key: 'id',
      dataIndex: 'id',
      title: $t('admin.shares.visitorList.columns.id'),
      width: 120,
      align: 'center',
      render: (value: string) => shortId(value),
    },
    {
      key: 'visitor_name',
      dataIndex: 'visitor_name',
      title: $t('admin.shares.visitorList.columns.visitorName'),
      width: 150,
      align: 'center',
      render: (value: string) => value || $t('admin.shares.visitorList.format.unknown'),
    },
    {
      key: 'visitor_email',
      dataIndex: 'visitor_email',
      title: $t('admin.shares.visitorList.columns.visitorEmail'),
      width: 200,
      ellipsis: true,
      align: 'center',
      render: (value: string) => value || $t('admin.shares.visitorList.format.empty'),
    },
    {
      key: 'created_at',
      dataIndex: 'created_at',
      title: $t('admin.shares.visitorList.columns.createdAt'),
      width: 160,
      align: 'center',
      render: (value: string) => formatDate(value),
    },
    {
      key: 'ip_address',
      dataIndex: 'ip_address',
      title: $t('admin.shares.visitorList.columns.ipAddress'),
      width: 160,
      align: 'center',
      render: (value: string, record: VisitorInfoDetail) => {
        if (record.ip_address) {
          return h('div', { class: 'ip-info' }, [
            h('span', { class: 'ip-text' }, record.ip_address),
            h('i', {
              class: 'fas fa-info-circle ml-1',
              title: $t('admin.shares.visitorList.ipInfoTitle'),
              style: 'cursor: pointer; opacity: 0.6; transition: all 0.2s; color: var(--color-brand-500);',
            }),
          ])
        }
        return $t('admin.shares.visitorList.format.empty')
      },
    },
    {
      key: 'share_key',
      dataIndex: 'share_key',
      title: $t('admin.shares.visitorList.columns.shareKey'),
      width: 200,
      ellipsis: true,
      align: 'center',
      render: (value: string, record: VisitorInfoDetail) =>
        h(
          'div',
          {
            class: 'share-link',
            onClick: () => viewShareContent(record.share_id, record.share_key),
            style: 'color: var(--color-brand-500); cursor: pointer; font-family: monospace; transition: all 0.2s;',
          },
          record.share_key
        ),
    },
    {
      key: 'visit_count',
      dataIndex: 'visit_count',
      title: $t('admin.shares.visitorList.columns.visitCount'),
      width: 100,
      align: 'center',
    },
    {
      key: 'actions',
      title: $t('admin.shares.visitorList.columns.actions'),
      width: 100,
      align: 'center',
      render: (_value: unknown, record: VisitorInfoDetail) =>
        h('div', { class: 'action-buttons flex gap-2 justify-center' }, [
          h(
            'button',
            {
              class:
                'inline-flex items-center justify-center px-3 py-1 text-xs bg-red-500/20 hover:bg-red-500/30 border border-red-500/30 text-red-400 rounded transition-all duration-200',
              onClick: (e: Event) => {
                e.stopPropagation()
                confirmDelete(record)
              },
            },
            $t('admin.shares.visitorList.actions.delete')
          ),
        ]),
    },
  ])

  const searchParams = reactive({
    page: 1,
    size: 10,
    keyword: '',
    share_id: '',
    order_by: 'created_at DESC',
  })

  const totalItems = ref(0)
  const visitorList = ref<VisitorInfoDetail[]>([])
  const loading = ref(false)

  const deleteDialog = reactive({
    visible: false,
    visitorId: '',
    visitorName: '',
    visitorEmail: '',
    createdAt: '',
    loading: false,
  })

  const loadVisitorList = async () => {
    try {
      loading.value = true
      const result = await shareApi.adminGetVisitorList(searchParams)

      if (result.success && result.data) {
        visitorList.value = result.data.list || []
        totalItems.value = result.data.total || 0
      }
    } catch (error: unknown) {
      toast.error((error as any)?.message || $t('admin.shares.visitorList.messages.fetchFailed'))
    } finally {
      loading.value = false
    }
  }

  const searchVisitors = () => {
    searchParams.page = 1
    loadVisitorList()
  }

  const handlePageChange = (page: number) => {
    searchParams.page = page
    loadVisitorList()
  }

  const handlePageSizeChange = (size: number) => {
    searchParams.size = size
    loadVisitorList()
  }

  const viewShareContent = async (shareId: string, shareKey: string) => {
    try {
      loading.value = true

      try {
        const data = await shareApi.adminGetShareAccessToken(shareId)

        if (data && data.access_token) {
          const accessToken = data.access_token
          router.push(`/share/${shareKey}?access_token=${accessToken}`)
        } else {
          router.push(`/share/${shareKey}`)
          toast.warning($t('admin.shares.shareList.messages.getAccessFailed'))
        }
      } catch {
        router.push(`/share/${shareKey}`)
        toast.warning($t('admin.shares.shareList.messages.getAccessFailed'))
      }
    } catch {
      router.push(`/share/${shareKey}`)
      toast.warning($t('admin.shares.shareList.messages.getAccessFailed'))
    } finally {
      loading.value = false
    }
  }

  const confirmDelete = (visitor: VisitorInfoDetail) => {
    deleteDialog.visitorId = visitor.id
    deleteDialog.visitorName = visitor.visitor_name || ''
    deleteDialog.visitorEmail = visitor.visitor_email || ''
    deleteDialog.createdAt = visitor.created_at
    deleteDialog.visible = true
  }

  const confirmDeleteVisitor = async () => {
    try {
      deleteDialog.loading = true
      await shareApi.adminDeleteVisitor(deleteDialog.visitorId)

      toast.success($t('admin.shares.visitorList.messages.deleteSuccess'))
      deleteDialog.visible = false
      loadVisitorList()
    } catch (error: unknown) {
      toast.error((error as any)?.message || $t('admin.shares.visitorList.messages.deleteFailed'))
    } finally {
      deleteDialog.loading = false
    }
  }

  const formatDate = (dateString?: string) => {
    if (!dateString) {
      return $t('admin.shares.visitorList.format.empty')
    }

    const date = new Date(dateString)
    return date.toLocaleString(getCurrentLocale(), {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit',
    })
  }

  const shortId = (id: string) => {
    if (!id) {
      return $t('admin.shares.visitorList.format.empty')
    }
    return id.replace(/^visitor_/, '').substring(0, 8)
  }

  onMounted(() => {
    loadVisitorList()
  })
</script>

<template>
  <div class="visitor-list-component">
    <div class="filter-bar">
      <div class="filters">
        <div class="filter-item">
          <CyberInput
            v-model="searchParams.share_id"
            :placeholder="$t('admin.shares.visitorList.filter.shareIdPlaceholder')"
            @keyup.enter="searchVisitors"
          />
        </div>

        <div class="filter-item search-box">
          <CyberInput
            v-model="searchParams.keyword"
            :placeholder="$t('admin.shares.visitorList.filter.keywordPlaceholder')"
            @keyup.enter="searchVisitors"
          >
            <template #suffix>
              <i class="fas fa-search cursor-pointer" @click="searchVisitors" />
            </template>
          </CyberInput>
        </div>
      </div>

      <div class="right-actions">
        <CyberButton type="primary" @click="searchVisitors">
          <i class="fas fa-search mr-1" /> {{ $t('admin.shares.visitorList.filter.searchButton') }}
        </CyberButton>
      </div>
    </div>

    <!-- 加载状态 -->
    <div v-if="loading" class="admin-loading">
      <CyberLoading :visible="true" :text="$t('admin.shares.visitorList.loading.text')" />
    </div>

    <!-- 空状态 -->
    <div v-else-if="!loading && visitorList.length === 0" class="empty-state">
      <div class="empty-content animate-slide-up">
        <div class="empty-icon-wrapper">
          <i class="fas fa-users empty-icon" />
        </div>
        <h3 class="empty-title">{{ $t('admin.shares.visitorList.empty.title') }}</h3>
        <p class="empty-desc">{{ $t('admin.shares.visitorList.empty.description') }}</p>
      </div>
    </div>

    <!-- 表格容器 -->
    <div v-else class="table-container">
      <div class="table-wrapper">
        <CyberDataTable
          :columns="tableColumns"
          :data="visitorList"
          :loading="false"
          :hoverable="true"
          :striped="true"
          :bordered="true"
          class="visitor-data-table"
        />
      </div>
    </div>

    <!-- 分页 -->
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

    <CyberDialog v-model="deleteDialog.visible" :title="$t('admin.shares.visitorList.deleteDialog.title')" width="400px">
      <div class="delete-dialog-content">
        <p>{{ $t('admin.shares.visitorList.deleteDialog.message') }}</p>
        <div class="visitor-info">
          <div class="info-item">
            <span class="label">{{ $t('admin.shares.visitorList.deleteDialog.visitorName') }}</span>
            <span class="value">{{ deleteDialog.visitorName || $t('admin.shares.visitorList.format.unknown') }}</span>
          </div>
          <div v-if="deleteDialog.visitorEmail" class="info-item">
            <span class="label">{{ $t('admin.shares.visitorList.deleteDialog.email') }}</span>
            <span class="value">{{ deleteDialog.visitorEmail }}</span>
          </div>
          <div class="info-item">
            <span class="label">{{ $t('admin.shares.visitorList.deleteDialog.visitTime') }}</span>
            <span class="value">{{ formatDate(deleteDialog.createdAt) }}</span>
          </div>
        </div>
      </div>

      <template #footer>
        <div class="dialog-footer">
          <CyberButton @click="deleteDialog.visible = false">{{
            $t('admin.shares.visitorList.deleteDialog.cancel')
          }}</CyberButton>
          <CyberButton type="danger" :loading="deleteDialog.loading" @click="confirmDeleteVisitor">
            {{ $t('admin.shares.visitorList.deleteDialog.confirm') }}
          </CyberButton>
        </div>
      </template>
    </CyberDialog>
  </div>
</template>

<style scoped lang="scss">
  .visitor-list-component {
    display: flex;
    flex-direction: column;
    gap: var(--space-md);
    flex: 1;
    min-height: 0;
  }

  .filter-bar {
    display: flex;
    justify-content: space-between;
    margin-bottom: var(--space-md);
    flex-wrap: wrap;
    gap: var(--space-md);
  }

  .right-actions {
    display: flex;
    gap: var(--space-sm);
  }

  .filters {
    display: flex;
    flex-wrap: wrap;
    gap: var(--space-sm);
    flex: 1;
    justify-content: flex-start;
  }

  .filter-item {
    min-width: 120px;
  }

  .search-box {
    min-width: 260px;
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

  .ip-info {
    display: flex;
    align-items: center;
    font-family: monospace;
    justify-content: center;
  }

  .ip-text {
    margin-right: var(--space-xs);
  }

  .share-link:hover {
    text-decoration: underline;
  }

  .action-buttons {
    display: flex;
    gap: var(--space-xs);
  }

  .pagination-container {
    margin-top: var(--space-xl);
    display: flex;
    justify-content: flex-end;
  }

  .delete-dialog-content {
    padding: var(--space-sm) 0;
  }

  .visitor-info {
    margin-top: var(--space-md);
    padding: var(--space-md);
    background: rgba(var(--color-brand-500-rgb), 0.05);
    border-radius: var(--radius-md);
    border-left: 2px solid rgba(var(--color-brand-500-rgb), 0.3);
  }

  .info-item {
    margin-bottom: var(--space-sm);
  }

  .info-item:last-child {
    margin-bottom: 0;
  }

  .label {
    color: var(--color-content-muted);
    margin-right: var(--space-sm);
  }

  .value {
    color: var(--color-content-default);
    font-weight: var(--font-medium);
  }

  .dialog-footer {
    display: flex;
    justify-content: flex-end;
    gap: var(--space-sm);
    padding-top: var(--space-md);
  }

  @media (max-width: 991px) {
    .filter-bar {
      flex-direction: column;
    }

    .filters {
      justify-content: flex-start;
    }
  }
</style>
