import { computed, nextTick, onMounted, ref, watch } from 'vue'
import { StorageUtil } from '@/utils/storage/storage'
import { setCssVariable, toggleGlobalClass } from '@/utils/ui/domUtils'
import { DEFAULTS, LAYOUT_MODES, STORAGE_KEYS } from '@/constants'
import { debounce } from '@/utils/common'
import { useElementSize } from '@vueuse/core'

/**
 */
export type LayoutMode = 'grid' | 'waterfall' | 'masonry' | 'large'

/**
 * 分享布局功能
 * 处理网格大小、布局方式等
 */
export function useShareLayout(onLayoutChange?: (layoutMode: LayoutMode) => void) {
  const layoutMode = ref<LayoutMode>(LAYOUT_MODES.WATERFALL)

  const gridSizeValue = ref(DEFAULTS.GRID_SIZE)

  const showGridSizePopover = ref(false)
  const gridSizePopoverStyle = ref({
    top: '40px',
    left: 'auto',
    right: '0',
  })

  const showLayoutPopover = ref(false)

  const settingsLoaded = ref(false)

  const containerRef = ref(null)
  const { width: containerWidth } = useElementSize(containerRef)

  const waterfallColumns = computed(() => {
    if (containerWidth.value) {
      const columnWidth = Math.max(Math.round(gridSizeValue.value * 1.2), 180)
      const gap = 20
      const columnsCount = Math.max(1, Math.floor(containerWidth.value / (columnWidth + gap)))
      return columnsCount
    }

    const width = window.innerWidth
    if (width < 768) {
      return 2 // 移动端默认2列
    } else if (width < 1024) {
      return 3
    } else if (width < 1440) {
      return 4
    } else if (width < 1600) {
      return 5
    }
    return 6 // 超宽屏幕使用6列
  })

  const largeColumns = computed(() => {
    const width = window.innerWidth
    if (width < 768) {
      return 1 // 移动端1列
    } else if (width < 1024) {
      return 2
    }
    return 3
  })

  const loadLayoutSettings = () => {
    const savedLayoutMode = StorageUtil.get<LayoutMode>(STORAGE_KEYS.LAYOUT_MODE)
    if (savedLayoutMode) {
      layoutMode.value = savedLayoutMode
    }

    const savedGridSize = StorageUtil.get<number>(STORAGE_KEYS.GRID_SIZE)
    if (savedGridSize) {
      gridSizeValue.value = savedGridSize
      updateGridSize(false)
    } else {
      updateGridSize(false)
    }

    settingsLoaded.value = true

    nextTick(() => {
      initContainer()
    })
  }

  const debouncedSaveGridSize = debounce((size: number) => {
    if (settingsLoaded.value) {
      StorageUtil.set(STORAGE_KEYS.GRID_SIZE, size)
    }
  }, 300)

  const updateGridSize = (shouldSave = true) => {
    setCssVariable('grid-item-size', `${gridSizeValue.value}px`, document.documentElement, true)

    if (layoutMode.value === 'waterfall') {
      setTimeout(() => {
        window.dispatchEvent(new Event('resize'))
      }, 50)
    }

    if (shouldSave) {
      debouncedSaveGridSize(gridSizeValue.value)
    }
  }

  watch(gridSizeValue, () => {
    if (settingsLoaded.value) {
      updateGridSize()
    }
  })

  const openGridSizePopover = (e?: MouseEvent | Event) => {
    if (e && typeof e.stopPropagation === 'function') {
      e.stopPropagation()
    }

    showGridSizePopover.value = true

    if (e && e.target) {
      const button = (e.target as HTMLElement).closest('button')
      if (button) {
        const rect = button.getBoundingClientRect()

        setTimeout(() => {
          const popup = document.querySelector('.grid-size-popover') as HTMLElement
          if (popup) {
            const top = rect.bottom + 10
            const right = window.innerWidth - rect.right

            const maxTop = window.innerHeight - popup.offsetHeight - 20
            const finalTop = Math.min(top, maxTop)

            popup.style.transformOrigin = `top ${window.innerWidth - right}px`
            popup.style.opacity = '0'
            popup.style.transform = 'scale(0.9) translateY(-10px)'

            popup.style.top = `${finalTop}px`
            popup.style.right = `${right}px`

            requestAnimationFrame(() => {
              popup.style.opacity = '1'
              popup.style.transform = 'scale(1) translateY(0)'
            })
          }
        }, 0)
      }
    }
  }

  const closeGridSizePopover = () => {
    showGridSizePopover.value = false
  }

  const openLayoutPopover = (e?: MouseEvent | Event) => {
    if (e && typeof e.stopPropagation === 'function') {
      e.stopPropagation()
    }

    showLayoutPopover.value = true
  }

  const closeLayoutPopover = () => {
    showLayoutPopover.value = false
  }

  const changeLayout = (value: LayoutMode) => {
    layoutMode.value = value
    showLayoutPopover.value = false

    nextTick(() => {
      updateGridSize(false)
    })

    StorageUtil.set(STORAGE_KEYS.LAYOUT_MODE, value)

    if (onLayoutChange) {
      onLayoutChange(value)
    }
  }

  const isImmersiveMode = ref(false)

  const toggleImmersiveMode = (value?: boolean) => {
    const newValue = value !== undefined ? value : !isImmersiveMode.value
    isImmersiveMode.value = newValue

    toggleGlobalClass('immersive-mode', newValue)
  }

  const cleanup = () => {
    if (isImmersiveMode.value) {
      toggleImmersiveMode(false)
    }

    containerRef.value = null
  }

  const initContainer = () => {
    if (!containerRef.value) {
      const contentEl = document.querySelector('.share-content') || document.body
      containerRef.value = contentEl
    }
  }

  onMounted(() => {
    nextTick(() => {
      initContainer()
    })
  })

  return {
    layoutMode,
    gridSizeValue,
    showGridSizePopover,
    gridSizePopoverStyle,
    showLayoutPopover,
    waterfallColumns,
    largeColumns,
    isImmersiveMode,
    settingsLoaded,
    containerRef,
    loadLayoutSettings,
    updateGridSize,
    openGridSizePopover,
    closeGridSizePopover,
    openLayoutPopover,
    closeLayoutPopover,
    changeLayout,
    toggleImmersiveMode,
    cleanup,
  }
}
