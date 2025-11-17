export interface Tab {
  label: string
  value: string
}

export interface CodeTabsProps {
  tabs: Tab[]
  default?: string
}
