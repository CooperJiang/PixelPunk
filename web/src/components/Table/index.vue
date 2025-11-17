<script setup lang="ts" generic="T extends Record<string, any>">
  import { computed, h, ref, watch, useSlots } from 'vue'
  import {
    type ColumnDef,
    FlexRender,
    type RowSelectionState,
    type SortingState,
    getCoreRowModel,
    getFilteredRowModel,
    getPaginationRowModel,
    getSortedRowModel,
    useVueTable,
  } from '@tanstack/vue-table'
  import { Checkbox } from '../Checkbox'
  import { useTexts } from '@/composables/useTexts'
  import type { CyberTableColumn } from './types'

  defineOptions({
    name: 'CyberTable',
  })

  /* Props 接口定义 */
  interface Props {
    data: T[]
    columns: CyberTableColumn<T>[]
    loading?: boolean
    loadingText?: string
    skeleton?: boolean
    skeletonRows?: number
    rowKey?: string | ((record: T) => string | number)
    selectable?: boolean
    selectedRowKeys?: (string | number)[]
    pagination?:
      | {
          current: number
          pageSize: number
          total: number
          showSizeChanger?: boolean
          showQuickJumper?: boolean
          pageSizeOptions?: number[]
        }
      | false
    bordered?: boolean
    striped?: boolean
    hoverable?: boolean
    size?: 'small' | 'medium' | 'large'
    maxHeight?: number | string
    emptyText?: string
    sortBy?: string
    sortOrder?: 'asc' | 'desc'
  }

  const props = withDefaults(defineProps<Props>(), {
    loading: false,
    loadingText: '',
    skeleton: false,
    skeletonRows: 5,
    rowKey: 'id',
    selectable: false,
    selectedRowKeys: () => [],
    bordered: true,
    striped: true,
    hoverable: true,
    size: 'medium',
    emptyText: '',
  })

  const emit = defineEmits<{
    'update:selectedRowKeys': [keys: (string | number)[]]
    'row-click': [record: T, index: number]
    'row-select': [record: T, selected: boolean]
    'select-all': [selected: boolean]
    'sort-change': [column: string, order: 'asc' | 'desc' | null]
    'page-change': [page: number]
    'page-size-change': [size: number]
  }>()

  const tableContainer = ref<HTMLElement>()

  const sorting = ref<SortingState>([])
  const rowSelection = ref<RowSelectionState>({})

  const { $t } = useTexts()
  const slots = useSlots()
  const resolvedLoadingText = computed(() => props.loadingText || ($t('table.loading') as string))
  const resolvedEmptyText = computed(() => props.emptyText || ($t('table.noData') as string))

  const tableClasses = computed(() => [
    `cyber-table-${props.size}`,
    {
      'cyber-table-bordered': props.bordered,
      'cyber-table-striped': props.striped,
      'cyber-table-hoverable': props.hoverable,
      'cyber-table-loading': props.loading,
    },
  ])

  const wrapperStyle = computed(() => ({
    maxHeight: typeof props.maxHeight === 'number' ? `${props.maxHeight}px` : props.maxHeight,
  }))

  /* 转换列定义 */
  const columns = computed((): ColumnDef<T>[] => {
    const cols: ColumnDef<T>[] = []

    if (props.selectable) {
      cols.push({
        id: 'select',
        size: 50,
        minSize: 50,
        maxSize: 50,
        header: ({ table }) =>
          h(
            'div',
            {
              class: 'flex items-center justify-center h-full',
            },
            [
              h(Checkbox, {
                modelValue: table.getIsAllRowsSelected(),
                indeterminate: table.getIsSomeRowsSelected() && !table.getIsAllRowsSelected(),
                onChange: (value) => {
                  table.getToggleAllRowsSelectedHandler()({ target: { checked: value } })
                },
              }),
            ]
          ),
        cell: ({ row }) =>
          h(
            'div',
            {
              class: 'flex items-center justify-center h-full',
            },
            [
              h(Checkbox, {
                modelValue: row.getIsSelected(),
                disabled: !row.getCanSelect(),
                onChange: (value) => {
                  row.getToggleSelectedHandler()({ target: { checked: value } })
                },
              }),
            ]
          ),
        meta: {
          align: 'center',
        },
      })
    }

    props.columns.forEach((column) => {
      const col: ColumnDef<T> = {
        id: column.key,
        accessorKey: column.dataIndex || column.key,
        header: column.title,
        size: typeof column.width === 'number' ? column.width : 150,
        minSize: typeof column.width === 'number' ? column.width : column.minWidth || 50,
        maxSize: typeof column.width === 'number' ? column.width : column.maxWidth || 1000,
        enableSorting: column.sortable ?? false,
        enableResizing: false,
        meta: {
          align: column.align || 'left',
          className: column.ellipsis ? 'ellipsis' : undefined,
        },
      }

      if (column.slot) {
        col.cell = ({ getValue, row }) => {
          const slotFn = (slots as any)[column.slot]
          if (typeof slotFn === 'function') {
            return h('div', { class: 'w-full' }, slotFn({ value: getValue(), record: row.original, index: row.index }))
          }
          return h('span', String(getValue()))
        }
      } else if (column.render) {
        col.cell = ({ getValue, row }) => column.render!(getValue(), row.original, row.index)
      }

      cols.push(col)
    })

    return cols
  })

  const table = useVueTable({
    get data() {
      return props.data
    },
    get columns() {
      return columns.value
    },
    state: {
      get sorting() {
        return sorting.value
      },
      get rowSelection() {
        return rowSelection.value
      },
    },
    onSortingChange: (updater) => {
      const newSorting = typeof updater === 'function' ? updater(sorting.value) : updater
      sorting.value = newSorting

      if (newSorting.length > 0) {
        const { id, desc } = newSorting[0]
        emit('sort-change', id, desc ? 'desc' : 'asc')
      } else {
        emit('sort-change', '', null)
      }
    },
    onRowSelectionChange: (updater) => {
      const newSelection = typeof updater === 'function' ? updater(rowSelection.value) : updater
      rowSelection.value = newSelection

      const selectedKeys = Object.keys(newSelection).filter((key) => newSelection[key])
      const selectedRowKeys = selectedKeys.map((key) => {
        const rowData = props.data[Number(key)]
        return typeof props.rowKey === 'string' ? rowData[props.rowKey] : props.rowKey(rowData)
      })
      emit('update:selectedRowKeys', selectedRowKeys)
    },
    getCoreRowModel: getCoreRowModel(),
    getSortedRowModel: getSortedRowModel(),
    getFilteredRowModel: getFilteredRowModel(),
    getPaginationRowModel: getPaginationRowModel(),
    enableRowSelection: props.selectable,
    enableSorting: true,
    enableFiltering: true,
  })

  const handleRowClick = (record: T, index: number) => {
    emit('row-click', record, index)
  }

  const handlePageChange = (page: number) => {
    emit('page-change', page)
  }

  const handlePageSizeChange = (size: number) => {
    emit('page-size-change', size)
  }

  watch(
    () => props.selectedRowKeys,
    (newKeys) => {
      const newSelection: RowSelectionState = {}
      newKeys.forEach((key) => {
        const index = props.data.findIndex((item) =>
          typeof props.rowKey === 'string' ? item[props.rowKey] === key : props.rowKey(item) === key
        )
        if (index !== -1) {
          newSelection[index] = true
        }
      })
      rowSelection.value = newSelection
    },
    { immediate: true }
  )

  defineExpose({
    scrollTo: (options: { top?: number; left?: number }) => {
      if (tableContainer.value) {
        tableContainer.value.scrollTo(options)
      }
    },
    clearSelection: () => {
      rowSelection.value = {}
      emit('update:selectedRowKeys', [])
    },
    toggleRowSelection: (row: T, selected?: boolean) => {
      const index = props.data.findIndex((item) =>
        typeof props.rowKey === 'string'
          ? item[props.rowKey] === (typeof props.rowKey === 'string' ? row[props.rowKey] : props.rowKey(row))
          : props.rowKey(item) === props.rowKey(row)
      )
      if (index !== -1) {
        const newSelection = { ...rowSelection.value }
        if (selected !== undefined) {
          newSelection[index] = selected
        } else {
          newSelection[index] = !newSelection[index]
        }
        rowSelection.value = newSelection
      }
    },
  })
</script>

<template>
  <div class="cyber-table" :class="tableClasses">
    <div class="cyber-table-wrapper" :style="wrapperStyle">
      <div v-if="loading" class="cyber-table-loading">
        <div class="loading-content">
          <div class="cyber-loading-spinner" />
          <p class="loading-text">{{ resolvedLoadingText }}</p>
        </div>
      </div>
      <div ref="tableContainer" class="cyber-table-container">
        <table class="cyber-table-element">
          <thead class="cyber-table-thead">
            <tr class="cyber-table-header-row">
              <th
                v-for="header in table.getFlatHeaders()"
                :key="header.id"
                class="cyber-table-th"
                :class="[
                  {
                    sortable: header.column.getCanSort(),
                    sorted: header.column.getIsSorted(),
                    'sort-asc': header.column.getIsSorted() === 'asc',
                    'sort-desc': header.column.getIsSorted() === 'desc',
                  },
                ]"
                :style="{
                  width: `${header.getSize()}px`,
                  minWidth: `${header.getSize()}px`,
                  maxWidth: `${header.getSize()}px`,
                  textAlign: (header.column.columnDef.meta as any)?.align || 'left',
                }"
                @click="header.column.getCanSort() ? header.column.getToggleSortingHandler()?.() : undefined"
              >
                <div class="cyber-table-th-content">
                  <FlexRender :render="header.column.columnDef.header" :props="header.getContext()" />
                  <div v-if="header.column.getCanSort()" class="sort-indicator">
                    <i v-if="!header.column.getIsSorted()" class="fas fa-sort" />
                    <i v-else-if="header.column.getIsSorted() === 'asc'" class="fas fa-sort-up" />
                    <i v-else-if="header.column.getIsSorted() === 'desc'" class="fas fa-sort-down" />
                  </div>
                </div>
              </th>
            </tr>
          </thead>

          <tbody class="cyber-table-tbody">
            <template v-if="loading && skeleton">
              <tr v-for="rowIndex in skeletonRows" :key="`skeleton-${rowIndex}`" class="cyber-table-row">
                <td
                  v-for="header in table.getFlatHeaders()"
                  :key="header.id"
                  class="cyber-table-td"
                  :style="{
                    width: `${header.getSize()}px`,
                    minWidth: `${header.getSize()}px`,
                    maxWidth: `${header.getSize()}px`,
                  }"
                >
                  <div class="p-2">
                    <div class="skeleton h-4 rounded" :style="{ width: `${60 + Math.random() * 30}%` }" />
                  </div>
                </td>
              </tr>
            </template>

            <!-- 正常数据显示 -->
            <template v-else-if="table.getRowModel().rows.length > 0">
              <tr
                v-for="row in table.getRowModel().rows"
                :key="row.id"
                class="cyber-table-row"
                :class="[
                  {
                    selected: row.getIsSelected(),
                    clickable: !!$attrs['onRowClick'],
                  },
                ]"
                @click="handleRowClick(row.original, row.index)"
              >
                <td
                  v-for="cell in row.getVisibleCells()"
                  :key="cell.id"
                  class="cyber-table-td"
                  :class="[(cell.column.columnDef.meta as any)?.className]"
                  :style="{
                    width: `${cell.column.getSize()}px`,
                    minWidth: `${cell.column.getSize()}px`,
                    maxWidth: `${cell.column.getSize()}px`,
                    textAlign: (cell.column.columnDef.meta as any)?.align || 'left',
                  }"
                >
                  <FlexRender :render="cell.column.columnDef.cell" :props="cell.getContext()" />
                </td>
              </tr>
            </template>

            <!-- 空状态 -->
            <tr v-else class="cyber-table-empty">
              <td :colspan="table.getAllColumns().length" class="cyber-table-empty-cell">
                <div class="empty-content">
                  <slot name="empty">
                    <div class="default-empty">
                      <i class="fas fa-inbox empty-icon" />
                      <p class="empty-text">{{ resolvedEmptyText }}</p>
                    </div>
                  </slot>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
      <div v-if="pagination && pagination.total > 0" class="cyber-table-pagination">
        <CyberPagination
          :total="pagination.total"
          :size="pagination.pageSize"
          :current-page="pagination.current"
          :show-page-size-selector="pagination.showSizeChanger"
          :show-quick-jumper="pagination.showQuickJumper"
          :page-size-options="pagination.pageSizeOptions"
          @update:current-page="handlePageChange"
          @update:size="handlePageSizeChange"
        />
      </div>
    </div>
  </div>
</template>

<style scoped lang="scss">
  @use './styles.scss';
</style>
