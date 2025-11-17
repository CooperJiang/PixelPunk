export interface CyberAdminWrapperProps {
  title?: string
  subtitle?: string
  icon?: string
  compact?: boolean
  contentPadding?: string
  sidebarLayout?: boolean
  showTopbar?: boolean
}

export interface CyberAdminWrapperSlots {
  icon?: () => unknown
  title?: () => unknown
  subtitle?: () => unknown
  stats?: () => unknown
  actions?: () => unknown
  topbar?: () => unknown
  toolbar?: () => unknown
  default?: () => unknown
  content?: () => unknown
}

export interface CyberAdminWrapperEmits {
  [key: string]: never
}

/* Layout configuration options */
export interface LayoutConfig {
  fixedHeader?: boolean
  showScrollbar?: boolean
  contentMaxWidth?: string
  breakpoints?: {
    mobile?: number
    tablet?: number
    desktop?: number
  }
}

/* Theme configuration */
export interface ThemeConfig {
  primaryColor?: string
  accentColor?: string
  backgroundColor?: string
  borderColor?: string
}
