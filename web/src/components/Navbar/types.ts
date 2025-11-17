export interface NavbarEmits {
  (e: 'login'): void
  (e: 'toggleMobileMenu'): void
}

export interface NavItem {
  path: string
  name: string
  label: string
  external?: boolean
  icon?: string
}

export interface UserMenuState {
  mobileMenuOpen: boolean
  userMenuOpen: boolean
}

export interface ScrollState {
  isScrolled: boolean
  scrollY: number
}

export interface NavbarProps {
  transparent?: boolean
  fixed?: boolean
}

export interface MenuAction {
  icon: string
  label: string
  path?: string
  action?: () => void
  variant?: 'normal' | 'admin' | 'danger'
}
