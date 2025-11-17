import { defineStore } from 'pinia'
import { StorageUtil } from '@/utils/storage/storage'
import { LAYOUT_MODE_KEY, SIDEBAR_COLLAPSED_KEY } from '@/constants/storage'
import { useSettingsStore } from './settings'

export type LayoutMode = 'top' | 'left'

interface LayoutConfig {
  sidebarWidth: number
  collapsedWidth: number
  breakpoints: {
    mobile: number
    tablet: number
    desktop: number
  }
}

interface LayoutState {
  mode: LayoutMode
  sidebarCollapsed: boolean
  config: LayoutConfig
}

const DEFAULT_CONFIG: LayoutConfig = {
  sidebarWidth: 210,
  collapsedWidth: 50,
  breakpoints: {
    mobile: 768,
    tablet: 1024,
    desktop: 1280,
  },
}

export const useLayoutStore = defineStore('layout', {
  state: (): LayoutState => {
    const savedMode = StorageUtil.get<LayoutMode>(LAYOUT_MODE_KEY)
    const validMode: LayoutMode = savedMode === 'top' || savedMode === 'left' ? savedMode : 'top'

    return {
      mode: validMode,
      sidebarCollapsed: StorageUtil.get<boolean>(SIDEBAR_COLLAPSED_KEY) || false,
      config: DEFAULT_CONFIG,
    }
  },

  getters: {
    isTopLayout: (state) => state.mode === 'top',
    isLeftLayout: (state) => state.mode === 'left',

    canToggleLayout: () => {
      const settingsStore = useSettingsStore()
      return settingsStore.enableMultiLayout ?? false
    },

    currentSidebarWidth: (state) => {
      if (state.mode === 'top') return 0
      return state.sidebarCollapsed ? state.config.collapsedWidth : state.config.sidebarWidth
    },

    mainContentMargin: (state) => {
      if (state.mode === 'top') return 0
      return state.sidebarCollapsed ? state.config.collapsedWidth : state.config.sidebarWidth
    },

    layoutClasses: (state) => ({
      'layout-top': state.mode === 'top',
      'layout-left': state.mode === 'left',
      'sidebar-collapsed': state.sidebarCollapsed,
    }),

    containerStyle: (state) => {
      if (state.mode === 'top') return {}
      return {
        marginLeft: `${state.sidebarCollapsed ? state.config.collapsedWidth : state.config.sidebarWidth}px`,
      }
    },
  },

  actions: {
    async initializeLayout() {
      const settingsStore = useSettingsStore()

      if (!settingsStore.isLoaded) {
        await settingsStore.loadGlobalSettings()
      }

      const multiLayoutEnabled = settingsStore.enableMultiLayout
      const rawDefaultLayout = settingsStore.defaultLayout
      const rawUserSelectedLayout = StorageUtil.get<LayoutMode>(LAYOUT_MODE_KEY)

      const validateLayout = (layout: any): LayoutMode => {
        return layout === 'top' || layout === 'left' ? layout : 'top'
      }

      const defaultLayout = validateLayout(rawDefaultLayout)
      const userSelectedLayout =
        rawUserSelectedLayout === 'top' || rawUserSelectedLayout === 'left' ? rawUserSelectedLayout : null

      if (multiLayoutEnabled) {
        const targetLayout = userSelectedLayout || defaultLayout
        this.setLayoutMode(targetLayout)
      } else {
        const targetLayout = defaultLayout
        this.setLayoutMode(targetLayout)
      }
    },

    setLayoutMode(mode: LayoutMode) {
      this.mode = mode
      StorageUtil.set(LAYOUT_MODE_KEY, mode)

      if (mode === 'top') {
        this._resetSidebarState()
      }
    },

    toggleLayoutMode() {
      const newMode = this.mode === 'left' ? 'top' : 'left'
      this.setLayoutMode(newMode)
      return newMode
    },

    toggleSidebar() {
      if (this.mode === 'left') {
        this.sidebarCollapsed = !this.sidebarCollapsed
        this._saveSidebarState()
      }
    },

    setSidebarCollapsed(collapsed: boolean) {
      if (this.mode === 'left') {
        this.sidebarCollapsed = collapsed
        this._saveSidebarState()
      }
    },

    handleResize(width: number) {
      if (this.mode === 'left') {
        const { tablet, desktop } = this.config.breakpoints

        if (width < tablet) {
          this.setSidebarCollapsed(true)
        } else if (width >= desktop) {
          this.setSidebarCollapsed(false)
        }
      }
    },

    _saveSidebarState() {
      StorageUtil.set(SIDEBAR_COLLAPSED_KEY, this.sidebarCollapsed)
    },

    _resetSidebarState() {
      this.sidebarCollapsed = false
      this._saveSidebarState()
    },
  },
})
