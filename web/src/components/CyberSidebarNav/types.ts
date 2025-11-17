export interface TabItem {
  key: string
  name: string
  icon?: string
  badge?: string | number
}

export interface CyberSidebarNavProps {
  tabs: TabItem[]
  activeTab: string
}

export interface CyberSidebarNavEmits {
  (e: 'tab-change', key: string): void
}
