export interface CyberTableColumn<T = any> {
  key: string
  title: string
  dataIndex?: string
  width?: number | string
  minWidth?: number
  maxWidth?: number
  align?: 'left' | 'center' | 'right'
  sortable?: boolean
  filterable?: boolean
  fixed?: 'left' | 'right'
  ellipsis?: boolean
  slot?: string
  render?: (value: any, record: T, index: number) => any
}

export interface CyberTableProps<T = any> {
  data: T[]
  columns: CyberTableColumn<T>[]
  loading?: boolean
  loadingText?: string
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

export interface CyberTableEmits<T = any> {
  'update:selectedRowKeys': [keys: (string | number)[]]
  'row-click': [record: T, index: number]
  'row-select': [record: T, selected: boolean]
  'select-all': [selected: boolean]
  'sort-change': [column: string, order: 'asc' | 'desc' | null]
  'page-change': [page: number]
  'page-size-change': [size: number]
}

export interface CyberTableInstance {
  scrollTo: (options: { top?: number; left?: number }) => void
  clearSelection: () => void
  toggleRowSelection: (row: any, selected?: boolean) => void
}
