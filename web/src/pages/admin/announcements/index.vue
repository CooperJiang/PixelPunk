<script setup lang="ts">
  import { h, onMounted, reactive, ref, computed } from 'vue'
  import { formatDate } from '@/utils/formatting/format'
  import { useToast } from '@/components/Toast/useToast'
  import { useTexts } from '@/composables/useTexts'
  import type { Announcement, AnnouncementListQuery } from '@/api/types/announcement'
  import { getAnnouncementList, deleteAnnouncement, toggleAnnouncementPin } from '@/api/admin/announcement'
  import { showConfirm } from '@/utils/dialog'
  import type { CyberTableColumn } from '@/components/Table/types'
  import AnnouncementForm from './components/AnnouncementForm.vue'
  import AnnouncementFilterPanel from './components/AnnouncementFilterPanel.vue'
  import AnnouncementSettings from './components/AnnouncementSettings.vue'

  defineOptions({
    name: 'AdminAnnouncements',
  })

  /* 工具 */
  const { $t } = useTexts()
  const toast = useToast()

  /* 状态图标映射 */
  const getStatusConfig = (status: string) => {
    const statusMap = {
      draft: {
        text: $t('admin.announcements.status.draft'),
        color: 'var(--color-content-muted)',
        bgColor: 'rgba(var(--color-content-muted-rgb), 0.2)',
      },
      published: {
        text: $t('admin.announcements.status.published'),
        color: 'var(--color-success-400)',
        bgColor: 'rgba(var(--color-success-rgb), 0.2)',
      },
      archived: {
        text: $t('admin.announcements.status.archived'),
        color: 'var(--color-warning-400)',
        bgColor: 'rgba(var(--color-warning-rgb), 0.2)',
      },
    }
    return statusMap[status as keyof typeof statusMap] || statusMap.draft
  }

  const dataTableColumns = computed<CyberTableColumn<Announcement>[]>(() => [
    {
      key: 'title',
      title: $t('admin.announcements.table.title'),
      width: 280,
      align: 'left',
      render: (value, record) =>
        h('div', { class: 'flex items-center gap-2' }, [
          h('div', { class: 'flex flex-col gap-1' }, [
            h('div', { class: 'flex items-center gap-2' }, [
              h('span', { class: 'font-medium', style: { color: 'var(--color-content-default)' } }, record.title),
              record.is_pinned
                ? h(
                    'span',
                    {
                      class: 'px-1.5 py-0.5 text-xs rounded font-medium',
                      style: {
                        background: 'rgba(var(--color-error-rgb), 0.2)',
                        color: 'var(--color-error-400)',
                      },
                    },
                    $t('admin.announcements.actions.pin')
                  )
                : null,
            ]),
            h('div', { class: 'text-xs', style: { color: 'var(--color-content-muted)' } }, `ID: ${record.id}`),
          ]),
        ]),
    },
    {
      key: 'summary',
      title: $t('admin.announcements.table.summary'),
      dataIndex: 'summary',
      width: 300,
      align: 'left',
      ellipsis: true,
      render: (value) =>
        h('span', { style: { color: 'var(--color-content-muted)' } }, value || $t('admin.announcements.table.noSummary')),
    },
    {
      key: 'status',
      title: $t('admin.announcements.table.status'),
      width: 90,
      align: 'center',
      render: (value, record) => {
        const config = getStatusConfig(record.status)
        return h(
          'span',
          {
            class: 'inline-block px-2 py-1 text-xs rounded font-medium whitespace-nowrap',
            style: {
              background: config.bgColor,
              color: config.color,
            },
          },
          config.text
        )
      },
    },
    {
      key: 'view_count',
      title: $t('admin.announcements.table.viewCount'),
      width: 90,
      align: 'center',
      render: (value) => h('span', { class: 'font-mono', style: { color: 'var(--color-success-400)' } }, value || 0),
    },
    {
      key: 'published_at',
      title: $t('admin.announcements.table.publishedAt'),
      width: 160,
      align: 'center',
      render: (value) =>
        h(
          'span',
          { style: { color: 'var(--color-content-muted)' } },
          value ? formatDate(value) : $t('admin.announcements.table.notSet')
        ),
    },
    {
      key: 'created_at',
      title: $t('admin.announcements.table.createdAt'),
      width: 160,
      align: 'center',
      render: (value) => h('span', { style: { color: 'var(--color-content-muted)' } }, formatDate(value)),
    },
    {
      key: 'actions',
      title: $t('admin.announcements.table.actions'),
      width: 160,
      align: 'center',
      render: (value, record) =>
        h('div', { class: 'flex gap-1 justify-center items-center' }, [
          h(
            'button',
            {
              class: 'inline-flex items-center justify-center w-7 h-7 text-xs rounded transition-colors',
              style: {
                background: record.is_pinned ? 'rgba(var(--color-warning-rgb), 0.2)' : 'rgba(var(--color-success-rgb), 0.2)',
                color: record.is_pinned ? 'var(--color-warning-400)' : 'var(--color-success-400)',
                border: record.is_pinned
                  ? '1px solid rgba(var(--color-warning-rgb), 0.4)'
                  : '1px solid rgba(var(--color-success-rgb), 0.4)',
              },
              title: record.is_pinned ? $t('admin.announcements.actions.unpin') : $t('admin.announcements.actions.pin'),
              onClick: () => togglePin(record),
              onMouseenter: (e: MouseEvent) => {
                const target = e.target as HTMLElement
                if (record.is_pinned) {
                  target.style.background = 'rgba(var(--color-warning-rgb), 0.3)'
                } else {
                  target.style.background = 'rgba(var(--color-success-rgb), 0.3)'
                }
              },
              onMouseleave: (e: MouseEvent) => {
                const target = e.target as HTMLElement
                if (record.is_pinned) {
                  target.style.background = 'rgba(var(--color-warning-rgb), 0.2)'
                } else {
                  target.style.background = 'rgba(var(--color-success-rgb), 0.2)'
                }
              },
            },
            [h('i', { class: record.is_pinned ? 'fas fa-thumbtack' : 'fas fa-thumbtack' })]
          ),
          h(
            'button',
            {
              class: 'inline-flex items-center justify-center w-7 h-7 text-xs rounded transition-colors',
              style: {
                background: 'rgba(var(--color-brand-500-rgb), 0.2)',
                color: 'var(--color-brand-500)',
                border: '1px solid rgba(var(--color-brand-500-rgb), 0.4)',
              },
              title: $t('admin.announcements.actions.edit'),
              onClick: () => editAnnouncement(record),
              onMouseenter: (e: MouseEvent) => {
                const target = e.target as HTMLElement
                target.style.background = 'rgba(var(--color-brand-500-rgb), 0.3)'
              },
              onMouseleave: (e: MouseEvent) => {
                const target = e.target as HTMLElement
                target.style.background = 'rgba(var(--color-brand-500-rgb), 0.2)'
              },
            },
            [h('i', { class: 'fas fa-edit' })]
          ),
          h(
            'button',
            {
              class: 'inline-flex items-center justify-center w-7 h-7 text-xs rounded transition-colors',
              style: {
                background: 'rgba(var(--color-error-rgb), 0.2)',
                color: 'var(--color-error-400)',
                border: '1px solid rgba(var(--color-error-rgb), 0.3)',
              },
              title: $t('admin.announcements.actions.delete'),
              onClick: () => deleteAnnouncementItem(record.id),
              onMouseenter: (e: MouseEvent) => {
                const target = e.target as HTMLElement
                target.style.background = 'rgba(var(--color-error-rgb), 0.3)'
              },
              onMouseleave: (e: MouseEvent) => {
                const target = e.target as HTMLElement
                target.style.background = 'rgba(var(--color-error-rgb), 0.2)'
              },
            },
            [h('i', { class: 'fas fa-trash' })]
          ),
        ]),
    },
  ])

  const announcements = ref<Announcement[]>([])
  const loading = ref(false)
  const saveLoading = ref(false)

  const totalAnnouncements = ref(0)
  const currentPage = ref(1)
  const pageSize = ref(20)

  const showFilter = ref(false)
  const currentFilters = reactive<AnnouncementListQuery>({
    keyword: '',
    status: undefined,
    is_pinned: undefined,
  })

  const showEditDialog = ref(false)
  const currentEditAnnouncement = ref<Announcement | null>(null)
  const isCreating = ref(false)
  const announcementFormRef = ref<any>(null)

  const fetchAnnouncementList = async () => {
    loading.value = true

    const params: AnnouncementListQuery = {
      page: currentPage.value,
      page_size: pageSize.value,
      keyword: currentFilters.keyword || undefined,
      status: currentFilters.status,
      is_pinned: currentFilters.is_pinned,
    }

    try {
      const result = await getAnnouncementList(params)
      if (result.code === 200) {
        const { data } = result
        announcements.value = data.data
        totalAnnouncements.value = data.pagination.total
      } else {
        toast.error($t('admin.announcements.messages.fetchError') + result.message)
      }
    } catch (error: any) {
      toast.error($t('admin.announcements.messages.fetchError') + error.message)
    }

    loading.value = false
  }

  const handleFilter = (filters: Partial<AnnouncementListQuery>) => {
    Object.assign(currentFilters, filters)
    currentPage.value = 1
    fetchAnnouncementList()
    showFilter.value = false
  }

  const handlePageChange = (page: number) => {
    currentPage.value = page
    fetchAnnouncementList()
  }

  const handlePageSizeChange = (size: number) => {
    pageSize.value = size
    currentPage.value = 1
    fetchAnnouncementList()
  }

  const createAnnouncement = () => {
    isCreating.value = true
    currentEditAnnouncement.value = null
    showEditDialog.value = true
  }

  const editAnnouncement = (announcement: Announcement) => {
    isCreating.value = false
    currentEditAnnouncement.value = JSON.parse(JSON.stringify(announcement))
    showEditDialog.value = true
  }

  const deleteAnnouncementItem = async (id: number) => {
    if (!showConfirm($t('admin.announcements.messages.deleteConfirm'))) {
      return
    }

    try {
      const result = await deleteAnnouncement(id)
      if (result.code === 200) {
        toast.success($t('admin.announcements.messages.deleteSuccess'))
        fetchAnnouncementList()
      } else {
        toast.error($t('admin.announcements.messages.deleteError') + result.message)
      }
    } catch (error: any) {
      toast.error($t('admin.announcements.messages.deleteError') + error.message)
    }
  }

  const togglePin = async (announcement: Announcement) => {
    const action = announcement.is_pinned ? $t('admin.announcements.actions.unpin') : $t('admin.announcements.actions.pin')
    try {
      const result = await toggleAnnouncementPin(announcement.id)
      if (result.code === 200) {
        toast.success($t('admin.announcements.messages.pinSuccess').replace('{action}', action))
        fetchAnnouncementList()
      } else {
        toast.error($t('admin.announcements.messages.pinError').replace('{action}', action) + result.message)
      }
    } catch (error: any) {
      toast.error($t('admin.announcements.messages.pinError').replace('{action}', action) + error.message)
    }
  }

  const handleRefresh = () => {
    fetchAnnouncementList()
  }

  const handleFormSuccess = () => {
    showEditDialog.value = false
    fetchAnnouncementList()
  }

  const handleConfirm = async () => {
    if (announcementFormRef.value) {
      await announcementFormRef.value.submitForm()
    }
  }

  const handleCancel = () => {
    showEditDialog.value = false
  }

  onMounted(() => {
    fetchAnnouncementList()
  })
</script>

<template>
  <div class="admin-announcements-page admin-page-container">
    <CyberAdminWrapper
      :title="$t('admin.announcements.title')"
      :subtitle="$t('admin.announcements.subtitle')"
      icon="fas fa-bullhorn"
      :show-topbar="showFilter"
    >
      <template #actions>
        <CyberButton type="primary" @click="createAnnouncement">
          <i class="fas fa-plus mr-2" />
          {{ $t('admin.announcements.actions.create') }}
        </CyberButton>

        <CyberButton type="secondary" :loading="loading" @click="handleRefresh">
          <i class="fas fa-sync-alt mr-2" />
          {{ $t('admin.announcements.actions.refresh') }}
        </CyberButton>

        <CyberButton
          type="secondary"
          :class="{ 'border-brand-400 text-brand-400': showFilter }"
          @click="showFilter = !showFilter"
        >
          <i class="fas fa-filter mr-1.5" />{{ $t('admin.announcements.actions.filter') }}
        </CyberButton>
      </template>

      <template #topbar>
        <AnnouncementFilterPanel :initial-filters="currentFilters" @filter="handleFilter" />
      </template>

      <template #content>
        <div class="admin-content-wrapper">
          <AnnouncementSettings />

          <div class="table-container">
            <div class="table-wrapper">
              <CyberSkeleton type="table" :count="pageSize" :loading="loading" />

              <CyberDataTable
                v-if="!loading && announcements.length > 0"
                :data="announcements"
                :columns="dataTableColumns"
                :loading="loading"
                :selectable="false"
                :hoverable="true"
                :striped="true"
                :bordered="true"
                size="small"
                row-key="id"
                :loading-text="$t('admin.announcements.table.loadingText')"
              />

              <div v-if="!loading && announcements.length === 0" class="empty-state">
                <div class="empty-content animate-slide-up">
                  <div class="empty-icon-wrapper">
                    <i class="fas fa-bullhorn empty-icon" />
                  </div>
                  <h3 class="empty-title">{{ $t('admin.announcements.empty.title') }}</h3>
                  <p class="empty-desc">{{ $t('admin.announcements.empty.description') }}</p>
                </div>
              </div>
            </div>

            <div v-if="totalAnnouncements > 0" class="pagination-wrapper">
              <CyberPagination
                v-model:current-page="currentPage"
                :total="totalAnnouncements"
                :page-size="pageSize"
                :show-page-size-selector="true"
                :show-quick-jumper="true"
                class="compact-pagination"
                @update:current-page="handlePageChange"
                @update:page-size="handlePageSizeChange"
              />
            </div>
          </div>
        </div>
      </template>
    </CyberAdminWrapper>

    <!-- 编辑/创建弹窗 -->
    <CyberDialog
      v-model="showEditDialog"
      :title="isCreating ? $t('admin.announcements.dialog.createTitle') : $t('admin.announcements.dialog.editTitle')"
      width="1200px"
      max-width="95vw"
      :close-on-click-overlay="false"
      :loading="saveLoading"
    >
      <AnnouncementForm
        ref="announcementFormRef"
        v-model="currentEditAnnouncement"
        :is-creating="isCreating"
        @submit="handleFormSuccess"
      />

      <template #footer>
        <div class="flex justify-end gap-3 p-4">
          <CyberButton type="secondary" size="medium" @click="handleCancel">
            {{ $t('admin.announcements.actions.cancel') }}
          </CyberButton>
          <CyberButton type="primary" size="medium" :loading="saveLoading" @click="handleConfirm">
            {{ isCreating ? $t('admin.announcements.actions.create') : $t('admin.announcements.actions.update') }}
          </CyberButton>
        </div>
      </template>
    </CyberDialog>
  </div>
</template>

<style scoped lang="scss">
  .admin-content-wrapper {
    background: var(--color-background-900);
    padding: 0;
    display: flex;
    flex-direction: column;
    flex: 1;
    min-height: 0;
  }

  .table-container {
    flex: 1;
    min-height: 0;
    display: flex;
    flex-direction: column;
  }

  .table-container:has(.cyber-table) {
    border-radius: var(--radius-sm);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.1);
    overflow: hidden;
  }

  .table-wrapper {
    flex: 1;
    min-height: 0;
    overflow: hidden;
  }

  .table-wrapper:has(.cyber-table) {
    background: rgba(var(--color-background-800-rgb), 0.3);
  }

  .pagination-wrapper {
    flex-shrink: 0;
    display: flex;
    justify-content: center;
    padding: var(--space-sm) var(--space-md);
    background: rgba(var(--color-background-800-rgb), 0.25);
    border-top: 1px solid rgba(var(--color-brand-500-rgb), 0.08);
    margin: 0;
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
  }

  .compact-pagination {
    :deep(.pagination-wrapper) {
      padding: var(--space-sm) var(--space-sm);
      min-height: auto;
    }

    :deep(.pagination-controls) {
      gap: var(--space-xs);
    }

    :deep(.pagination-btn) {
      width: var(--space-2xl);
      height: var(--space-2xl);
      font-size: var(--text-xs);
      min-width: var(--space-2xl);
    }

    :deep(.page-size-selector) {
      gap: var(--space-sm);
    }

    :deep(.select-trigger) {
      min-height: var(--space-2xl);
      padding: 0 var(--space-sm);
      font-size: var(--text-xs);
    }

    :deep(.jump-input) {
      width: var(--space-4xl);
      height: var(--space-2xl);
      font-size: var(--text-xs);
    }

    :deep(.page-size-label),
    :deep(.jump-label),
    :deep(.total-info) {
      font-size: var(--text-xs);
    }
  }

  @keyframes float {
    0%,
    100% {
      transform: translateY(0px);
    }
    50% {
      transform: translateY(-10px);
    }
  }
</style>
