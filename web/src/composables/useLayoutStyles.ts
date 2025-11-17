import { computed } from 'vue'
import { useLayoutStore } from '@/store/layout'
import { useRoute } from 'vue-router'

/**
 * 布局感知的样式组合式函数
 * 为不同布局模式提供自适应的样式配置
 */
export function useLayoutStyles() {
  const layoutStore = useLayoutStore()
  const route = useRoute()

  const isTopLayout = computed(() => layoutStore.isTopLayout)
  const isLeftLayout = computed(() => layoutStore.isLeftLayout)

  const pageType = computed(() => {
    const path = route.path
    if (path === '/') return 'home'
    if (path.startsWith('/admin')) return 'admin'
    if (path.startsWith('/docs')) return 'docs'
    if (path.includes('gallery') || path.includes('explore')) return 'gallery'
    if (path.includes('share')) return 'share'
    return 'default'
  })

  const containerClasses = computed(() => ({
    'layout-container': true,
    'layout-top': isTopLayout.value,
    'layout-left': isLeftLayout.value,
    [`page-${pageType.value}`]: true,
  }))

  const pageContainerStyle = computed(() => {
    const baseStyle = {
      width: 'var(--layout-container-width)',
      maxWidth: 'var(--layout-content-max-width)',
      padding: 'var(--layout-content-padding)',
      margin: '0 auto',
    }

    return baseStyle
  })

  const gridContainerStyle = computed(() => ({
    display: 'grid',
    gridTemplateColumns: 'var(--layout-grid-cols)',
    gap: 'var(--layout-grid-gap)',
  }))

  const cardStyle = computed(() => ({
    padding: 'var(--layout-card-padding)',
    borderRadius: 'var(--layout-card-radius)',
    transition: `all var(--layout-animation-duration) ease`,
  }))

  const titleStyle = computed(() => ({
    fontSize: 'var(--layout-title-size)',
    lineHeight: 1.2,
  }))

  const getGridCols = (minWidth: number = 300) => {
    if (isLeftLayout.value) {
      return `repeat(auto-fill, minmax(${Math.max(minWidth - 50, 200)}px, 1fr))`
    }
    return `repeat(auto-fill, minmax(${minWidth}px, 1fr))`
  }

  const getSpacing = (type: 'section' | 'card' | 'element' = 'element') => {
    return `var(--layout-${type}-gap)`
  }

  const getPageSpecificStyles = () => {
    switch (pageType.value) {
      case 'admin':
        return {
          '--layout-card-gap': isLeftLayout.value ? '0.75rem' : '1.5rem',
          '--layout-table-density': isLeftLayout.value ? 'compact' : 'comfortable',
        }
      case 'gallery':
        return {
          '--layout-grid-gap': isLeftLayout.value ? '0.5rem' : '1rem',
          '--layout-image-size': isLeftLayout.value ? '150px' : '200px',
        }
      case 'home':
        return isTopLayout.value
          ? {
              '--layout-hero-height': '100vh',
              '--layout-section-gap': '5rem',
            }
          : {}
      default:
        return {}
    }
  }

  return {
    isTopLayout,
    isLeftLayout,
    pageType,

    containerClasses,

    pageContainerStyle,
    gridContainerStyle,
    cardStyle,
    titleStyle,

    getGridCols,
    getSpacing,
    getPageSpecificStyles,
  }
}

export function useAdminLayoutStyles() {
  const { isLeftLayout } = useLayoutStyles()

  const adminTableStyle = computed(() => ({
    '--table-row-height': isLeftLayout.value ? '2.5rem' : '3rem',
    '--table-cell-padding': isLeftLayout.value ? '0.5rem' : '0.75rem',
    '--table-font-size': isLeftLayout.value ? '0.8125rem' : '0.875rem',
  }))

  const adminFormStyle = computed(() => ({
    '--form-field-height': isLeftLayout.value ? '2.25rem' : '2.5rem',
    '--form-input-padding': isLeftLayout.value ? '0.5rem' : '0.75rem',
    '--form-label-size': isLeftLayout.value ? '0.8125rem' : '0.875rem',
  }))

  const adminCardGrid = computed(() => {
    if (isLeftLayout.value) {
      return 'repeat(auto-fill, minmax(280px, 1fr))'
    }
    return 'repeat(auto-fill, minmax(350px, 1fr))'
  })

  return {
    adminTableStyle,
    adminFormStyle,
    adminCardGrid,
  }
}

export function useGalleryLayoutStyles() {
  const { isLeftLayout } = useLayoutStyles()

  const fileGridStyle = computed(() => {
    const baseStyle = {
      display: 'grid',
      gap: 'var(--layout-grid-gap)',
    }

    if (isLeftLayout.value) {
      return {
        ...baseStyle,
        gridTemplateColumns: 'repeat(auto-fill, minmax(180px, 1fr))',
        '--image-aspect-ratio': '1 / 1',
      }
    }

    return {
      ...baseStyle,
      gridTemplateColumns: 'repeat(auto-fill, minmax(250px, 1fr))',
      '--image-aspect-ratio': '4 / 3',
    }
  })

  const fileCardStyle = computed(() => ({
    borderRadius: 'var(--layout-card-radius)',
    overflow: 'hidden',
    transition: `transform var(--layout-animation-duration) ease`,
  }))

  return {
    fileGridStyle,
    fileCardStyle,
  }
}
