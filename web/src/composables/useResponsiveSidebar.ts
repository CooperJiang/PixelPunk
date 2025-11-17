/**
 * 响应式侧边栏组合式函数
 * 根据屏幕尺寸动态调整侧边栏宽度和行为

 */
import { computed, onMounted, onUnmounted, ref } from 'vue'
import { StorageUtil } from '@/utils/storage/storage'
import { RESPONSIVE_BREAKPOINTS } from '@/constants/responsive'

export interface SidebarConfig {
  isCollapsed: boolean
  isMobile: boolean
  isTablet: boolean
  screenWidth: number
  sidebarWidth: number
  collapsedWidth: number
  shouldAutoCollapse: boolean
}

export function useResponsiveSidebar() {
  const STORAGE_KEY = 'admin_sidebar_collapsed'
  const MOBILE_COLLAPSED_KEY = 'admin_sidebar_mobile_collapsed'

  const screenWidth = ref(typeof window !== 'undefined' ? window.innerWidth : 1024)

  const savedState = typeof window !== 'undefined' ? StorageUtil.get<boolean>(STORAGE_KEY) : null
  const isCollapsed = ref(savedState !== null ? savedState : false) // 默认展开状态

  const isMobileOverlayOpen = ref(false)
  const isTransitioning = ref(false)

  const isMobile = computed(() => screenWidth.value < RESPONSIVE_BREAKPOINTS.mobile)
  const isTablet = computed(
    () => screenWidth.value >= RESPONSIVE_BREAKPOINTS.mobile && screenWidth.value < RESPONSIVE_BREAKPOINTS.desktop
  )
  const isDesktop = computed(() => screenWidth.value >= RESPONSIVE_BREAKPOINTS.desktop)

  const currentSidebarWidth = computed(() => (isCollapsed.value ? 60 : 220))

  const collapsedWidth = computed(() => 60)

  const shouldAutoCollapse = computed(() => false)

  const sidebarConfig = computed<SidebarConfig>(() => ({
    isCollapsed: isCollapsed.value,
    isMobile: isMobile.value,
    isTablet: isTablet.value,
    screenWidth: screenWidth.value,
    sidebarWidth: currentSidebarWidth.value,
    collapsedWidth: collapsedWidth.value,
    shouldAutoCollapse: shouldAutoCollapse.value,
  }))

  const updateScreenWidth = () => {
    screenWidth.value = window.innerWidth
  }

  const initializeSidebar = () => {
    if (isMobile.value) {
      isCollapsed.value = true
      const mobileState = StorageUtil.get<boolean>(MOBILE_COLLAPSED_KEY)
      if (mobileState !== null) {
        isMobileOverlayOpen.value = !mobileState
      }
    }
  }

  const toggleSidebar = () => {
    isTransitioning.value = true

    if (isMobile.value) {
      isMobileOverlayOpen.value = !isMobileOverlayOpen.value
      StorageUtil.set(MOBILE_COLLAPSED_KEY, !isMobileOverlayOpen.value)
    } else {
      isCollapsed.value = !isCollapsed.value
      StorageUtil.set(STORAGE_KEY, isCollapsed.value)
    }

    setTimeout(() => {
      isTransitioning.value = false
      window.dispatchEvent(new Event('resize'))
    }, 300)
  }

  const closeMobileSidebar = () => {
    if (isMobile.value && isMobileOverlayOpen.value) {
      isMobileOverlayOpen.value = false
      StorageUtil.set(MOBILE_COLLAPSED_KEY, true)
    }
  }

  const expandSidebar = () => {
    if (isCollapsed.value && !isMobile.value) {
      isCollapsed.value = false
      StorageUtil.set(STORAGE_KEY, false)
    }
  }

  const collapseSidebar = () => {
    if (!isCollapsed.value && !isMobile.value) {
      isCollapsed.value = true
      StorageUtil.set(STORAGE_KEY, true)
    }
  }

  const handleResize = () => {
    const oldIsMobile = isMobile.value
    const _oldScreenWidth = screenWidth.value

    updateScreenWidth()

    if (oldIsMobile !== isMobile.value) {
      initializeSidebar()
    }
  }

  const handleKeydown = (event: KeyboardEvent) => {
    if ((event.ctrlKey || event.metaKey) && event.key === 'b') {
      event.preventDefault()
      toggleSidebar()
    }

    if (event.key === 'Escape' && isMobile.value && isMobileOverlayOpen.value) {
      event.preventDefault()
      closeMobileSidebar()
    }
  }

  const handleRouteChange = () => {
    if (isMobile.value && isMobileOverlayOpen.value) {
      closeMobileSidebar()
    }
  }

  onMounted(() => {
    updateScreenWidth()
    initializeSidebar()
    window.addEventListener('resize', handleResize)
    window.addEventListener('keydown', handleKeydown)

    document.body.style.overflowX = 'hidden'
  })

  onUnmounted(() => {
    window.removeEventListener('resize', handleResize)
    window.removeEventListener('keydown', handleKeydown)
    document.body.style.overflowX = ''
  })

  return {
    isCollapsed,
    isMobileOverlayOpen,
    isTransitioning,
    screenWidth,

    isMobile,
    isTablet,
    isDesktop,
    currentSidebarWidth,
    collapsedWidth,
    shouldAutoCollapse,
    sidebarConfig,

    toggleSidebar,
    closeMobileSidebar,
    expandSidebar,
    collapseSidebar,
    handleRouteChange,

    updateScreenWidth,
    initializeSidebar,
  }
}
