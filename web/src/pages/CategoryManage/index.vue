<script setup lang="ts">
  import { h, onMounted } from 'vue'
  import { useRouter } from 'vue-router'
  import { formatDate } from '@/utils/formatting/format'
  import { useConfirmDialog } from '@/composables/useConfirmDialog'
  import { useTexts } from '@/composables/useTexts'
  import type { ImageCategory } from '@/api/types/category'
  import type { CyberTableColumn } from '@/components/Table/types'
  import { useCategoryManagement } from './composables/useCategoryManagement'
  import { getSourceInfo, getStatusStyle } from './composables/useCategoryHelpers'
  import CategoryForm from './components/CategoryForm.vue'
  import CategoryFilterPanel from './components/CategoryFilterPanel.vue'
  import CreateCategoryDialog from './components/CreateCategoryDialog.vue'
  import BatchOperations from './components/BatchOperations.vue'

  defineOptions({
    name: 'CategoryManage',
  })

  const router = useRouter()
  const { $t } = useTexts()

  /* 使用 confirm dialog */
  const { dialogState, showConfirmDialog, handleConfirm, handleCancel } = useConfirmDialog()

  /* 使用分类管理逻辑 */
  const {
    categories,
    loading,
    saveLoading,
    totalCategories,
    currentPage,
    pageSize,
    showFilter,
    currentFilters,
    showEditDialog,
    currentEditCategory,
    showCreateDialog,
    selectedCategoryIds,
    fetchCategoryList,
    handleFilter,
    handlePageChange,
    handlePageSizeChange,
    editCategory,
    saveCategoryChanges,
    toggleCategoryStatus,
    deleteCategoryItem,
    handleBatchDelete,
    handleRefresh,
    handleCategoryCreated,
    handleBatchCompleted,
    handleSelectionChange,
  } = useCategoryManagement()

  /* 跳转到资源页面并按分类筛选 */
  const goToResourceByCategory = (category: ImageCategory) => {
    router.push({
      path: '/resource',
      query: {
        categoryId: String(category.id),
        category_name: category.name,
      },
    })
  }

  const dataTableColumns: CyberTableColumn<ImageCategory>[] = [
    {
      key: 'name',
      title: $t('category.table.columns.name'),
      width: 200,
      align: 'left',
      render: (value, record) =>
        h('div', { class: 'flex flex-col' }, [
          h(
            'div',
            {
              class:
                'font-medium text-content border-b border-dashed border-content-muted/30 cursor-pointer hover:text-brand-400 hover:border-brand-400/60 transition-colors w-fit',
              title: $t('category.table.clickToView'),
              onClick: (e: MouseEvent) => {
                e.stopPropagation()
                goToResourceByCategory(record)
              },
            },
            record.name
          ),
          h('div', { class: 'text-xs text-content-muted' }, `ID: ${record.id}`),
        ]),
    },
    {
      key: 'description',
      title: $t('category.table.columns.description'),
      dataIndex: 'description',
      width: 250,
      align: 'left',
      ellipsis: true,
      render: (value) => h('span', { class: 'text-content-muted' }, value || $t('category.table.noDescription')),
    },
    {
      key: 'source',
      title: $t('category.table.columns.source'),
      width: 110,
      align: 'center',
      render: (value, record) => {
        const source = getSourceInfo(record.source)
        return h(
          'span',
          {
            class:
              'inline-flex items-center gap-1 whitespace-nowrap rounded-md px-2 py-1 text-xs font-medium transition-all border',
            style: source.style,
          },
          [h('i', { class: source.icon }), h('span', source.text)]
        )
      },
    },
    {
      key: 'file_count',
      title: $t('category.table.columns.fileCount'),
      width: 100,
      align: 'center',
      render: (value) => h('span', { class: 'font-mono text-brand-400' }, value || 0),
    },
    {
      key: 'sort_order',
      title: $t('category.table.columns.sortOrder'),
      width: 80,
      align: 'center',
      render: (value) => h('span', { class: 'font-mono text-content-muted' }, value),
    },
    {
      key: 'status',
      title: $t('category.table.columns.status'),
      width: 90,
      align: 'center',
      render: (value, record) => {
        const statusInfo = getStatusStyle(record.status)
        return h(
          'span',
          {
            class:
              'inline-flex items-center gap-1.5 whitespace-nowrap rounded-md px-2.5 py-1 text-xs font-medium transition-all border',
            style: statusInfo.containerStyle,
          },
          [
            h('span', {
              class: 'w-1.5 h-1.5 rounded-full',
              style: statusInfo.dotStyle,
            }),
            h('span', statusInfo.text),
          ]
        )
      },
    },
    {
      key: 'created_at',
      title: $t('category.table.columns.createdAt'),
      width: 160,
      align: 'center',
      render: (value) => h('span', { class: 'text-content-muted' }, formatDate(value)),
    },
    {
      key: 'actions',
      title: $t('category.table.columns.actions'),
      width: 160,
      align: 'center',
      slot: 'actions',
    },
  ]

  const handleRowClick = (_row: ImageCategory, _index: number) => {}

  const handleDeleteCategory = async (category: ImageCategory) => {
    const confirmed = await showConfirmDialog({
      title: $t('category.dialog.delete.title'),
      message: $t('category.dialog.delete.message'),
      type: 'danger',
    })
    if (confirmed) {
      deleteCategoryItem(category.id)
    }
  }

  onMounted(() => {
    fetchCategoryList()
  })
</script>

<template>
  <div class="category-manage-page">
    <div class="page-header mb-4 p-4">
      <div class="flex items-center justify-between gap-4">
        <div class="flex items-center gap-3">
          <div class="page-header-icon flex h-9 w-9 items-center justify-center rounded-lg border">
            <i class="fas fa-folder-tree text-lg text-brand-400" />
          </div>
          <div>
            <div class="flex items-center gap-2">
              <h1 class="text-lg font-bold text-content-heading">{{ $t('category.page.title') }}</h1>
              <span
                class="inline-flex items-center rounded-full px-2 py-0.5 text-xs font-medium text-brand-400"
                :style="{
                  backgroundColor: 'rgba(var(--color-brand-500-rgb), 0.1)',
                }"
              >
                {{ totalCategories }}
              </span>
            </div>
            <p class="mt-0.5 text-xs text-content-muted">{{ $t('category.page.subtitle') }}</p>
          </div>
        </div>

        <div class="flex gap-2">
          <CyberButton type="primary" size="small" @click="showCreateDialog = true">
            <i class="fas fa-plus mr-1.5" />
            {{ $t('category.page.create') }}
          </CyberButton>

          <CyberButton type="outlined" size="small" :loading="loading" @click="handleRefresh">
            <i class="fas fa-sync-alt mr-1.5" />
            {{ $t('category.page.refresh') }}
          </CyberButton>

          <CyberButton
            type="outlined"
            size="small"
            :class="{ 'bg-brand-500/5 border-brand-400 text-brand-400': showFilter }"
            @click="showFilter = !showFilter"
          >
            <i class="fas fa-filter mr-1.5" />
            {{ showFilter ? $t('category.page.hideFilter') : $t('category.page.filter') }}
          </CyberButton>
        </div>
      </div>
    </div>

    <transition name="filter-slide">
      <CategoryFilterPanel v-if="showFilter" :initial-filters="currentFilters" @filter="handleFilter" />
    </transition>

    <BatchOperations
      v-if="selectedCategoryIds.length > 0"
      :selected-categories="selectedCategoryIds"
      :all-categories="categories"
      @clear-selection="() => (selectedCategoryIds = [])"
      @batch-delete="handleBatchDelete"
      @batch-completed="handleBatchCompleted"
    />

    <div class="table-container mt-4 overflow-hidden">
      <div class="table-wrapper">
        <CyberSkeleton type="table" :count="pageSize" :loading="loading" />

        <div v-if="!loading && categories.length === 0" class="category-empty-state">
          <div class="category-icon-wrapper">
            <i class="fas fa-folder-open text-content-muted/30 text-5xl" />
          </div>
          <h3 class="mt-4 text-lg font-medium text-content">{{ $t('category.empty.title') }}</h3>
          <p class="mt-2 max-w-md text-center text-sm text-content-muted">
            {{ $t('category.empty.description') }}
          </p>
          <CyberButton type="primary" size="medium" class="mt-4" @click="showCreateDialog = true">
            <i class="fas fa-plus mr-2" />
            {{ $t('category.empty.action') }}
          </CyberButton>
        </div>

        <CyberDataTable
          v-if="!loading && categories.length > 0"
          :data="categories"
          :columns="dataTableColumns"
          :loading="loading"
          :selectable="true"
          :hoverable="true"
          :striped="true"
          :bordered="true"
          size="small"
          row-key="id"
          max-height="calc(100vh - 200px)"
          :selected-row-keys="selectedCategoryIds"
          :loading-text="$t('category.table.loading')"
          @update:selected-row-keys="handleSelectionChange"
          @row-click="handleRowClick"
        >
          <template #actions="{ record }">
            <div class="action-buttons">
              <button class="action-btn action-btn-edit" :title="$t('category.actions.edit')" @click.stop="editCategory(record)">
                <i class="fas fa-edit" />
              </button>
              <button
                class="action-btn action-btn-toggle"
                :title="record.status === 'active' ? $t('category.actions.archive') : $t('category.actions.activate')"
                @click.stop="toggleCategoryStatus(record)"
              >
                <i :class="record.status === 'active' ? 'fas fa-archive' : 'fas fa-undo'" />
              </button>
              <button
                class="action-btn action-btn-delete"
                :title="$t('category.actions.delete')"
                @click.stop="handleDeleteCategory(record)"
              >
                <i class="fas fa-trash" />
              </button>
            </div>
          </template>
        </CyberDataTable>
      </div>

      <div v-if="totalCategories > 0" class="pagination-wrapper bg-background-800/50 border-t border-subtle p-4 backdrop-blur-sm">
        <CyberPagination
          v-model:current-page="currentPage"
          :total="totalCategories"
          :page-size="pageSize"
          :show-page-size-selector="true"
          :show-quick-jumper="true"
          @update:current-page="handlePageChange"
          @update:page-size="handlePageSizeChange"
        />
      </div>
    </div>

    <CyberDialog v-model="showEditDialog" :title="$t('category.dialog.edit.title')" width="500px">
      <CategoryForm v-model="currentEditCategory" />

      <template #footer>
        <div class="flex justify-end gap-2 p-4">
          <CyberButton type="secondary" @click="showEditDialog = false">{{ $t('category.dialog.edit.cancel') }}</CyberButton>
          <CyberButton type="primary" :loading="saveLoading" @click="saveCategoryChanges">{{
            $t('category.dialog.edit.save')
          }}</CyberButton>
        </div>
      </template>
    </CyberDialog>

    <CreateCategoryDialog v-model="showCreateDialog" @category-created="handleCategoryCreated" />

    <!-- Confirm Dialog -->
    <CyberConfirmDialog
      v-if="dialogState.isVisible"
      v-model="dialogState.isVisible"
      :title="dialogState.config.title"
      :message="dialogState.config.message"
      :type="dialogState.config.type"
      :confirm-text="dialogState.config.confirmText"
      :cancel-text="dialogState.config.cancelText"
      :show-input="dialogState.config.showInput"
      :input-placeholder="dialogState.config.inputPlaceholder"
      :input-value="dialogState.config.inputValue"
      :input-validator="dialogState.config.inputValidator"
      @confirm="handleConfirm"
      @cancel="handleCancel"
    />
  </div>
</template>

<style scoped lang="scss">
  .category-manage-page {
    .page-header {
      position: relative;
      overflow: hidden;
      background: rgba(var(--color-background-800-rgb), 0.6);
      border: 1px solid rgba(var(--color-brand-500-rgb), 0.15);
      border-radius: var(--radius-sm);
      box-shadow: var(--shadow-sm);
      backdrop-filter: var(--backdrop-blur-md);

      &::before {
        content: '';
        position: absolute;
        top: 0;
        left: 0;
        right: 0;
        bottom: 0;
        background: linear-gradient(135deg, rgba(var(--color-brand-500-rgb), 0.05), transparent);
        pointer-events: none;
      }

      .page-header-icon {
        background: linear-gradient(135deg, rgba(var(--color-brand-500-rgb), 0.12), rgba(var(--color-brand-500-rgb), 0.06));
        border: 1px solid rgba(var(--color-brand-500-rgb), 0.2);
        border-radius: var(--radius-sm);
      }
    }
  }

  .table-container {
    display: flex;
    flex-direction: column;
    position: relative;
    background: rgba(var(--color-background-800-rgb), 0.6);
    border: 1.5px solid rgba(var(--color-brand-500-rgb), 0.15);
    border-radius: var(--radius-sm);
    box-shadow: var(--shadow-sm);
    backdrop-filter: var(--backdrop-blur-md);
    padding: 1rem;
    overflow: hidden;

    :deep(.cyber-data-table) {
      border: none !important;
      border-radius: var(--radius-sm) !important;
      overflow: hidden !important;
    }

    :deep(.cyber-table-wrapper) {
      border-radius: var(--radius-sm) !important;
      overflow: hidden !important;
    }

    :deep(.cyber-table) {
      border-radius: var(--radius-sm) !important;
      overflow: hidden !important;
    }

    :deep(.cyber-table-thead) {
      border-radius: var(--radius-sm) var(--radius-sm) 0 0 !important;
    }

    :deep(.cyber-table-tbody tr:last-child td:first-child) {
      border-bottom-left-radius: var(--radius-sm) !important;
    }

    :deep(.cyber-table-tbody tr:last-child td:last-child) {
      border-bottom-right-radius: var(--radius-sm) !important;
    }
  }

  .table-wrapper {
    flex: 1;
    min-height: 0;
    overflow: hidden;
  }

  .pagination-wrapper {
    flex-shrink: 0;
    display: flex;
    justify-content: center;
  }

  .action-buttons {
    display: flex;
    gap: 0.5rem;
    justify-content: center;
    align-items: center;

    .action-btn {
      display: inline-flex;
      align-items: center;
      justify-content: center;
      width: 2rem;
      height: 2rem;
      padding: 0.375rem;
      font-size: 0.875rem;
      border-radius: var(--radius-md);
      transition: all var(--transition-normal) var(--ease-in-out);
      cursor: pointer;

      &:hover {
        transform: translateY(-2px) scale(1.08);
      }

      &:active {
        transform: scale(0.92);
      }
    }

    .action-btn-edit {
      background: rgba(var(--color-brand-500-rgb), 0.1);
      color: var(--color-brand-400);
      border: 1px solid rgba(var(--color-brand-500-rgb), 0.3);

      &:hover {
        background: rgba(var(--color-brand-500-rgb), 0.2);
        border-color: rgba(var(--color-brand-500-rgb), 0.5);
        box-shadow: 0 0 12px rgba(var(--color-brand-500-rgb), 0.3);
      }
    }

    .action-btn-toggle {
      background: rgba(var(--color-warning-rgb), 0.1);
      color: var(--color-warning-400);
      border: 1px solid rgba(var(--color-warning-rgb), 0.3);

      &:hover {
        background: rgba(var(--color-warning-rgb), 0.2);
        border-color: rgba(var(--color-warning-rgb), 0.5);
        box-shadow: 0 0 12px rgba(var(--color-warning-rgb), 0.3);
      }
    }

    .action-btn-delete {
      background: rgba(var(--color-error-rgb), 0.1);
      color: var(--color-error-400);
      border: 1px solid rgba(var(--color-error-rgb), 0.3);

      &:hover {
        background: rgba(var(--color-error-rgb), 0.2);
        border-color: rgba(var(--color-error-rgb), 0.5);
        box-shadow: 0 0 12px rgba(var(--color-error-rgb), 0.3);
      }
    }
  }

  .category-empty-state {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 4rem 2rem;
    text-align: center;
  }

  .category-icon-wrapper {
    width: 5rem;
    height: 5rem;
    display: flex;
    align-items: center;
    justify-content: center;
    background: linear-gradient(135deg, rgba(var(--color-brand-500-rgb), 0.12), rgba(var(--color-brand-500-rgb), 0.06));
    border: 2px solid rgba(var(--color-brand-500-rgb), 0.2);
    border-radius: var(--radius-2xl);
  }
</style>
