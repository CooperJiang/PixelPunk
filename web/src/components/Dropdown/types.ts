export interface DropdownOption {
  label: string
  value: string | number | boolean
  disabled?: boolean
  color?: string
  isCreateOption?: boolean
  originalValue?: string
  [key: string]: unknown
}

export interface DropdownProps {
  modelValue: string | number | boolean | null | undefined | Array<string | number | boolean>
  options: DropdownOption[]
  placeholder?: string
  disabled?: boolean
  multiple?: boolean
  searchable?: boolean
  maxHeight?: string
  borderColor?: string
  textColor?: string
  bgColor?: string
  isColorMode?: boolean
  clearable?: boolean
  allowCreate?: boolean
  createLabel?: string
  pagination?: boolean
  loading?: boolean
  loadingMore?: boolean
  hasMore?: boolean
  pageSize?: number
  width?: string
  height?: string
}

export interface DropdownEmits {
  (e: 'update:modelValue', value: string | number | boolean | null | undefined | Array<string | number | boolean>): void
  (e: 'change', value: string | number | boolean | null | undefined | Array<string | number | boolean>): void
  (e: 'create', option: DropdownOption): void
  (e: 'load-more'): void
  (e: 'search', query: string): void
}
