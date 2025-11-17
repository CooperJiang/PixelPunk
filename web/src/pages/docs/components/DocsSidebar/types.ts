export interface NavItem {
  id: string
  label: string
  iconClass: string
}

export interface DocsSidebarProps {
  navItems: NavItem[]
  activeSection: string
  readingProgress: number
  scrollToSection: (event: Event) => void
}
