<script setup lang="ts">
  import { h, onMounted, reactive, ref } from 'vue'
  import { formatDate } from '@/utils/formatting/format'
  import { useToast } from '@/components/Toast/useToast'
  import { useTexts } from '@/composables/useTexts'
  import type { CategoryTemplate } from '@/api/types/category'
  import { getTemplateList, updateTemplate, deleteTemplate } from '@/api/admin/category'
  import { showConfirm } from '@/utils/dialog'
  import TemplateForm from './components/TemplateForm.vue'
  import TemplateFilterPanel from './components/TemplateFilterPanel.vue'
  import CreateTemplateDialog from './components/CreateTemplateDialog.vue'
  import BatchOperations from './components/BatchOperations.vue'
  import type { CyberTableColumn } from '@/components/Table/types'
  import { DEFAULT_CATEGORY_SORT } from '@/constants'

  defineOptions({
    name: 'AdminCategories',
  })

  /* 工具 */
  const { $t } = useTexts()
  const toast = useToast()

  /* 表格列配置 */
  const dataTableColumns: CyberTableColumn<CategoryTemplate>[] = [
    {
      key: 'name',
      title: $t('admin.categories.columns.name'),
      width: 200,
      align: 'left',
      render: (value, record) =>
        h('div', { class: 'flex items-center' }, [
          record.icon ? h('i', { class: `${record.icon} mr-2`, style: { color: 'var(--color-brand-500)' } }) : null,
          h('div', {}, [
            h('div', { class: 'font-medium', style: { color: 'var(--color-content-default)' } }, record.name),
            h('div', { class: 'text-xs', style: { color: 'var(--color-content-muted)' } }, `ID: ${record.id}`),
          ]),
        ]),
    },
    {
      key: 'description',
      title: $t('admin.categories.columns.description'),
      dataIndex: 'description',
      width: 250,
      align: 'left',
      ellipsis: true,
      render: (value) => value || $t('admin.categories.status.noDescription'),
    },
    {
      key: 'sort_order',
      title: $t('admin.categories.columns.sortOrder'),
      width: 80,
      align: 'center',
      render: (value) => h('span', { class: 'font-mono', style: { color: 'var(--color-brand-500)' } }, value),
    },
    {
      key: 'is_popular',
      title: $t('admin.categories.columns.isPopular'),
      width: 80,
      align: 'center',
      render: (value, record) =>
        h(
          'span',
          {
            class: 'inline-block whitespace-nowrap rounded px-2 py-1 text-xs font-medium',
            style: record.is_popular
              ? {
                  background: 'rgba(var(--color-error-rgb), 0.2)',
                  color: 'var(--color-error-400)',
                }
              : {
                  background: 'rgba(var(--color-content-muted-rgb), 0.2)',
                  color: 'var(--color-content-muted)',
                },
          },
          record.is_popular ? $t('admin.categories.status.popular') : $t('admin.categories.status.normal')
        ),
    },
    {
      key: 'usage_count',
      title: $t('admin.categories.columns.usageCount'),
      width: 100,
      align: 'center',
      render: (value) => h('span', { class: 'font-mono', style: { color: 'var(--color-success-400)' } }, value),
    },
    {
      key: 'created_at',
      title: $t('admin.categories.columns.createdAt'),
      width: 160,
      align: 'center',
      render: (value) => formatDate(value),
    },
    {
      key: 'actions',
      title: $t('admin.categories.columns.actions'),
      width: 120,
      align: 'center',
      render: (value, record) =>
        h('div', { class: 'action-buttons flex gap-1 justify-center items-center' }, [
          h(
            'button',
            {
              class: 'inline-flex items-center justify-center w-7 h-7 text-xs rounded transition-colors focus:outline-none',
              style: {
                background: 'rgba(var(--color-brand-500-rgb), 0.2)',
                color: 'var(--color-brand-500)',
                border: '1px solid rgba(var(--color-brand-500-rgb), 0.4)',
              },
              title: $t('admin.categories.actions.editTemplate'),
              onClick: () => editTemplate(record),
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
              title: $t('admin.categories.actions.deleteTemplate'),
              onClick: () => {
                if (showConfirm($t('admin.categories.confirm.deleteOne'))) {
                  deleteTemplateItem(record.id)
                }
              },
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
  ]

  const templates = ref<CategoryTemplate[]>([])
  const loading = ref(false)
  const saveLoading = ref(false)

  /* 分页 */
  const totalTemplates = ref(0)
  const currentPage = ref(1)
  const pageSize = ref(10)

  /* 筛选相关 */
  const showFilter = ref(false)
  const currentFilters = reactive({
    keyword: '',
    is_popular: undefined as boolean | undefined,
    sort_by: DEFAULT_CATEGORY_SORT.sortBy,
    sort_order: DEFAULT_CATEGORY_SORT.sortOrder,
  })

  /* 编辑模板 */
  const showEditDialog = ref(false)
  const currentEditTemplate = ref<CategoryTemplate | null>(null)

  /* 新增模板 */
  const showCreateDialog = ref(false)

  /* 批量选择 */
  const selectedTemplateIds = ref<number[]>([])

  const fetchTemplateList = async () => {
    loading.value = true

    const params = {
      page: currentPage.value,
      size: pageSize.value,
      keyword: currentFilters.keyword || undefined,
      is_popular: currentFilters.is_popular,
      sort_by: currentFilters.sort_by || DEFAULT_CATEGORY_SORT.sortBy,
      sort_order: currentFilters.sort_order || DEFAULT_CATEGORY_SORT.sortOrder,
    }

    try {
      const result = await getTemplateList(params)
      if (result.success) {
        const { data } = result
        templates.value = data.templates
        totalTemplates.value = data.total

        selectedTemplateIds.value = selectedTemplateIds.value.filter((id) =>
          templates.value.some((template) => template.id === id)
        )
      }
    } catch {
      toast.error($t('admin.categories.messages.fetchFailed'))
    }

    loading.value = false
  }

  const handleFilter = (filters) => {
    Object.assign(currentFilters, filters)
    currentPage.value = 1
    selectedTemplateIds.value = [] // 清空选择
    fetchTemplateList()
    showFilter.value = false
  }

  const handlePageChange = (page) => {
    currentPage.value = page
    fetchTemplateList()
  }

  const handlePageSizeChange = (size) => {
    pageSize.value = size
    currentPage.value = 1 // 重置到第一页
    fetchTemplateList()
  }

  const editTemplate = (template) => {
    currentEditTemplate.value = JSON.parse(JSON.stringify(template))
    showEditDialog.value = true
  }

  const saveTemplateChanges = async () => {
    if (!currentEditTemplate.value) {
      return
    }

    saveLoading.value = true

    try {
      const result = await updateTemplate({
        id: currentEditTemplate.value.id,
        name: currentEditTemplate.value.name,
        description: currentEditTemplate.value.description,
        icon: currentEditTemplate.value.icon,
        is_popular: currentEditTemplate.value.is_popular,
        sort_order: currentEditTemplate.value.sort_order,
      })

      if (result.success) {
        toast.success($t('admin.categories.messages.updated'))
        showEditDialog.value = false
        fetchTemplateList()
      }
    } catch {
      toast.error($t('admin.categories.messages.updateFailed'))
    }
    saveLoading.value = false
  }

  const deleteTemplateItem = async (id: number) => {
    try {
      const result = await deleteTemplate({ id })
      if (result.success) {
        toast.success($t('admin.categories.messages.deleted'))
        fetchTemplateList()
      }
    } catch {
      toast.error($t('admin.categories.messages.deleteFailed'))
    }
  }

  const handleRefresh = () => {
    selectedTemplateIds.value = []
    fetchTemplateList()
  }

  const handleTemplateCreated = (_newTemplate: CategoryTemplate) => {
    // 不需要重复显示提示，CreateTemplateDialog 组件内部已经显示
    fetchTemplateList()
  }

  const handleBatchCompleted = () => {
    fetchTemplateList() // 重新获取列表
  }

  const handleSelectionChange = (selectedRows: (string | number)[]) => {
    selectedTemplateIds.value = selectedRows.map((id) => Number(id))
  }

  const handleRowClick = (_row: CategoryTemplate, _index: number) => {}

  onMounted(() => {
    fetchTemplateList()
  })
</script>

<template>
  <div class="admin-categories-page admin-page-container">
    <CyberAdminWrapper
      :title="$t('admin.categories.title')"
      :subtitle="$t('admin.categories.subtitle')"
      icon="fas fa-tags"
      :show-topbar="showFilter"
    >
      <template #actions>
        <CyberButton type="primary" @click="showCreateDialog = true">
          <i class="fas fa-plus mr-2" />
          {{ $t('admin.categories.buttons.create') }}
        </CyberButton>

        <CyberButton type="secondary" :loading="loading" @click="handleRefresh">
          <i class="fas fa-sync-alt mr-2" />
          {{ $t('admin.categories.buttons.refresh') }}
        </CyberButton>

        <CyberButton
          type="secondary"
          :class="{ 'border-brand-400 text-brand-400': showFilter }"
          @click="showFilter = !showFilter"
        >
          <i class="fas fa-filter mr-1.5" />{{ $t('admin.categories.buttons.filter') }}
        </CyberButton>
      </template>

      <template #topbar>
        <TemplateFilterPanel :initial-filters="currentFilters" @filter="handleFilter" />
      </template>

      <template #content>
        <div class="admin-content-wrapper">
          <BatchOperations
            :selected-templates="selectedTemplateIds"
            :all-templates="templates"
            @clear-selection="() => (selectedTemplateIds = [])"
            @batch-completed="handleBatchCompleted"
          />

          <div class="table-container">
            <div class="table-wrapper">
              <CyberSkeleton type="table" :count="pageSize" :loading="loading" />

              <CyberDataTable
                v-if="!loading && templates.length > 0"
                :data="templates"
                :columns="dataTableColumns"
                :loading="loading"
                :selectable="true"
                :hoverable="true"
                :striped="true"
                :bordered="true"
                size="small"
                row-key="id"
                :selected-row-keys="selectedTemplateIds"
                :loading-text="$t('admin.categories.loading.table')"
                @update:selected-row-keys="handleSelectionChange"
                @row-click="handleRowClick"
              />

              <div v-if="!loading && templates.length === 0" class="empty-state">
                <div class="empty-content animate-slide-up">
                  <div class="empty-icon-wrapper">
                    <i class="fas fa-tags empty-icon" />
                  </div>
                  <h3 class="empty-title">{{ $t('admin.categories.empty.title') }}</h3>
                  <p class="empty-desc">{{ $t('admin.categories.empty.description') }}</p>
                </div>
              </div>
            </div>

            <div v-if="totalTemplates > 0" class="pagination-wrapper">
              <CyberPagination
                v-model:current-page="currentPage"
                :total="totalTemplates"
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

    <CyberDialog v-model="showEditDialog" :title="$t('admin.categories.dialogs.edit')" width="500px">
      <TemplateForm v-model="currentEditTemplate" />

      <template #footer>
        <div class="flex justify-end gap-2 p-4">
          <CyberButton type="secondary" @click="showEditDialog = false">{{ $t('admin.categories.buttons.cancel') }}</CyberButton>
          <CyberButton type="primary" :loading="saveLoading" @click="saveTemplateChanges">{{
            $t('admin.categories.buttons.save')
          }}</CyberButton>
        </div>
      </template>
    </CyberDialog>

    <CreateTemplateDialog v-model="showCreateDialog" @template-created="handleTemplateCreated" />
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

    :deep(.cyber-table) {
      tbody tr.selected {
        background-color: transparent !important;
      }

      tbody tr.selected:hover {
        background-color: rgba(var(--color-brand-500-rgb), 0.05) !important;
      }
    }
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

  .action-buttons {
    display: flex;
    gap: var(--space-xs);
    justify-content: center;
    align-items: center;
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
</style>
