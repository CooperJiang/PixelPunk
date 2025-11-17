export interface TabItem {
  id?: string
  key?: string
  name?: string
  label?: string
  icon?: string
  badge?: string | number
  disabled?: boolean
  visible?: boolean
}

export interface SideNavTabsProps {
  title: string
  subtitle?: string
  icon?: string
  tabs: TabItem[]
  activeTab?: string
  syncWithUrl?: boolean
  urlParamName?: string
}

export interface SideNavTabsEmits {
  (e: 'update:activeTab', key: string): void
  (e: 'tab-change', key: string, tab: TabItem): void
  (e: 'change', key: string): void
}
