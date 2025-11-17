<script setup lang="ts">
  import { computed, h, onMounted } from 'vue'
  import { useTaggingHistory } from '../composables/useTaggingHistory'
  import { useTexts } from '@/composables/useTexts'
  import type { TaggingTaskItem } from '@/api/user/automation'
  import type { CyberTableColumn } from '@/components/Table/types'

  defineOptions({
    name: 'TaggingHistoryTable',
  })

  const { $t } = useTexts()

  const { isLoading, data, total, query, hasLoaded, loadData, handlePageChange, handleSizeChange, handleStatusFilter } =
    useTaggingHistory()

  /* 组件挂载时加载数据（仅首次） */
  onMounted(() => {
    if (!hasLoaded.value) {
      loadData()
    }
  })

  const statusMap = computed(() => ({
    none: {
      label: $t('automation.tagging.status.none'),
      color: 'var(--color-content-muted)',
      bgColor: 'var(--color-badge-neutral-bg)',
      icon: 'fas fa-circle',
    },
    pending: {
      label: $t('automation.tagging.status.pending'),
      color: 'var(--color-brand-400)',
      bgColor: 'var(--color-badge-primary-bg)',
      icon: 'fas fa-clock',
    },
    processing: {
      label: $t('automation.tagging.status.processing'),
      color: 'var(--color-warning-500)',
      bgColor: 'var(--color-badge-warning-bg)',
      icon: 'fas fa-spinner fa-spin',
    },
    done: {
      label: $t('automation.tagging.status.done'),
      color: 'var(--color-success-400)',
      bgColor: 'var(--color-badge-success-bg)',
      icon: 'fas fa-check-circle',
    },
    failed: {
      label: $t('automation.tagging.status.failed'),
      color: 'var(--color-error-400)',
      bgColor: 'var(--color-badge-error-bg)',
      icon: 'fas fa-times-circle',
    },
    ignored: {
      label: $t('automation.tagging.status.ignored'),
      color: 'var(--color-content-muted)',
      bgColor: 'var(--color-badge-neutral-bg)',
      icon: 'fas fa-ban',
    },
  }))

  const statusOptions = computed(() => [
    { label: $t('common.all'), value: '' },
    { label: $t('automation.tagging.status.none'), value: 'none' },
    { label: $t('automation.tagging.status.pending'), value: 'pending' },
    { label: $t('automation.tagging.status.processing'), value: 'processing' },
    { label: $t('automation.tagging.status.done'), value: 'done' },
    { label: $t('automation.tagging.status.failed'), value: 'failed' },
    { label: $t('automation.tagging.status.ignored'), value: 'ignored' },
  ])

  const columns: CyberTableColumn<TaggingTaskItem>[] = [
    {
      key: 'thumbnail',
      title: $t('automation.history.columns.preview'),
      width: 70,
      align: 'center',
      render: (value, record) => {
        if (!record.thumbnail_url) {
          return h(
            'div',
            {
              class: 'flex h-10 w-10 items-center justify-center rounded border border-subtle bg-background-700 mx-auto',
            },
            [h('i', { class: 'fas fa-image text-content-muted text-xs' })]
          )
        }
        return h('img', {
          src: record.thumbnail_url,
          class: 'h-10 w-10 rounded border border-subtle object-cover mx-auto',
        })
      },
    },
    {
      key: 'file_name',
      title: $t('automation.history.columns.filename'),
      dataIndex: 'file_name',
      minWidth: 180,
      ellipsis: true,
    },
    {
      key: 'status',
      title: $t('automation.history.columns.status'),
      width: 100,
      align: 'center',
      render: (value, record) => {
        const status = statusMap.value[record.status] || statusMap.value.none
        return h(
          'span',
          {
            class: 'inline-flex items-center gap-1.5 rounded-full px-2.5 py-1 text-xs font-medium',
            style: {
              backgroundColor: status.bgColor,
              color: status.color,
            },
          },
          [h('i', { class: status.icon }), status.label]
        )
      },
    },
    {
      key: 'format',
      title: $t('automation.history.columns.format'),
      dataIndex: 'format',
      width: 70,
      align: 'center',
      render: (value, record) => {
        return h('span', { class: 'text-xs font-mono text-brand-400' }, record.format?.toUpperCase() || '-')
      },
    },
    {
      key: 'size',
      title: $t('automation.history.columns.size'),
      dataIndex: 'size_formatted',
      width: 90,
      align: 'center',
      render: (value, record) => {
        return h('span', { class: 'text-xs text-content-muted' }, record.size_formatted || '-')
      },
    },
    {
      key: 'resolution',
      title: $t('automation.history.columns.resolution'),
      dataIndex: 'resolution',
      width: 100,
      align: 'center',
      render: (value, record) => {
        return h('span', { class: 'text-xs text-content-muted' }, record.resolution || '-')
      },
    },
    {
      key: 'created_at',
      title: $t('automation.history.columns.createdAt'),
      dataIndex: 'created_at',
      width: 160,
      align: 'center',
    },
    {
      key: 'updated_at',
      title: $t('automation.history.columns.updatedAt'),
      dataIndex: 'updated_at',
      width: 160,
      align: 'center',
    },
    {
      key: 'error_message',
      title: $t('automation.history.columns.error'),
      minWidth: 200,
      ellipsis: true,
      render: (value, record) => {
        if (!record.error_message) {
          return h('span', { class: 'text-content-muted' }, '-')
        }
        return h('span', { class: 'text-error-400 text-xs' }, record.error_message)
      },
    },
  ]

  const pageSizeOptions = [10, 20, 50, 100]

  const currentPageData = computed(() => data.value)
</script>

<template>
  <div class="tagging-history-table">
    <div class="table-header">
      <div class="flex items-center gap-2">
        <i class="fas fa-history" />
        <h3 class="table-title">{{ $t('automation.history.title') }}</h3>
        <span class="record-count">{{ $t('automation.history.recordCount', { count: total }) }}</span>
      </div>

      <div class="filter-controls">
        <div class="flex items-center gap-2">
          <span class="filter-label">{{ $t('actions.filter') }}</span>
          <CyberDropdown
            v-model="query.status"
            :options="statusOptions"
            :placeholder="$t('automation.taggingHistory.filter.statusPlaceholder')"
            @change="handleStatusFilter(query.status)"
          />
        </div>

        <CyberButton type="outlined" :loading="isLoading" @click="loadData">
          <i class="fas fa-sync-alt mr-1.5" :class="{ 'fa-spin': isLoading }" />
          {{ $t('automation.history.refresh') }}
        </CyberButton>
      </div>
    </div>

    <div class="table-container">
      <CyberTable
        :columns="columns"
        :data="currentPageData"
        :loading="isLoading"
        :skeleton="true"
        :emptyText="$t('automation.history.empty')"
      />
    </div>

    <CyberPagination
      v-if="total > 0"
      class="pagination-section"
      :current-page="query.page"
      :total="total"
      :page-size="query.limit"
      :page-size-options="pageSizeOptions"
      @update:current-page="handlePageChange"
      @update:page-size="handleSizeChange"
    />
  </div>
</template>

<style scoped lang="scss">
  .tagging-history-table {
    padding: var(--space-md);
    background: rgba(var(--color-background-800-rgb), 0.6);
    border: 1.5px solid rgba(var(--color-brand-500-rgb), 0.15);
    border-radius: var(--radius-sm);
    box-shadow: var(--shadow-sm);
    backdrop-filter: var(--backdrop-blur-md);
  }

  .table-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: var(--space-sm);

    i {
      color: var(--color-brand-400);
    }

    .table-title {
      font-size: var(--text-base);
      font-weight: 600;
      color: var(--text-heading);
    }

    .record-count {
      font-size: var(--text-xs);
      color: var(--text-muted);
    }
  }

  .filter-controls {
    display: flex;
    align-items: center;
    gap: var(--space-sm);

    .filter-label {
      white-space: nowrap;
      font-size: var(--text-sm);
      color: var(--text-muted);
    }
  }

  .table-container {
    padding: var(--space-md);
    background: rgba(var(--color-background-900-rgb), 0.3);
    border-radius: var(--radius-sm);
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
      border: none !important;
      border-radius: var(--radius-sm) !important;
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

  .pagination-section {
    margin-top: var(--space-sm);
  }
</style>
