import { onBeforeUnmount, onMounted, ref } from 'vue'

interface Options {
  containerSelector?: string
  gridSelector?: string
  childSelector?: string
  itemMinWidth?: number // 单卡片的最小宽度（含gap）
  columnGap?: number // 列间距，用于更精确计算
  rowMultiple?: number // 行数基数（例如默认取3行填满）
  debounceMs?: number
  defaultSize?: number
  debug?: boolean
  preferCssColumns?: boolean // 优先使用 CSS 列定义推断列数（默认开启）
  minDeltaPx?: number // 触发重算的最小宽度变化阈值，避免频繁抖动
  mode?: 'once' | 'observe'
}

export function useResponsivePageSize(options: Options = {}) {
  const {
    containerSelector = '.images-grid',
    gridSelector = containerSelector,
    childSelector = undefined,
    itemMinWidth = 230,
    columnGap = 15,
    rowMultiple = 3,
    debounceMs = 120,
    defaultSize = 24,
    debug: _debug = false,
    preferCssColumns = true,
    minDeltaPx = 8,
    mode = 'observe',
  } = options

  const pageSize = ref<number>(defaultSize)
  const columns = ref<number>(1)
  const containerWidth = ref<number>(0)

  let resizeObserver: ResizeObserver | null = null
  let debounceTimer: number | null = null
  let lastWidth = 0
  let stopped = false

  const debug = _debug
  const log = (..._args: any[]) => {
    if (debug) {
    }
  }

  const parsePx = (val: string | null | undefined): number => {
    if (!val) {
      return 0
    }
    const n = parseFloat(val.toString())
    return isNaN(n) ? 0 : n
  }

  const cleanup = () => {
    if (resizeObserver) {
      resizeObserver.disconnect()
    }
    window.removeEventListener('resize', handleResize)
    stopped = true
  }

  const calc = (force = false) => {
    if (stopped) {
      log('❌ Calculation stopped (cleanup called)')
      return
    }

    log('🔄 Starting calculation...', force ? '(FORCED)' : '')
    const container = document.querySelector(containerSelector) as HTMLElement | null
    const grid = document.querySelector(gridSelector) as HTMLElement | null

    if (!container) {
      log('❌ Container not found:', containerSelector)
      return
    }
    log('✅ Container found:', containerSelector)

    const width = container.clientWidth
    log('📏 Container width:', width, 'px (last:', lastWidth, 'px)')

    if (!force && Math.abs(width - lastWidth) < minDeltaPx) {
      log('⏭️  Width change too small, skipping (delta:', Math.abs(width - lastWidth), 'px < minDelta:', minDeltaPx, 'px)')
      return
    }
    lastWidth = width
    containerWidth.value = width

    let colsFromCss: number | null = null
    let gapFromCss = columnGap
    let colsFromChild: number | null = null

    if (grid) {
      log('✅ Grid element found:', gridSelector)
      const style = window.getComputedStyle(grid)
      const gridTemplateColumns = style.getPropertyValue('grid-template-columns') || ''
      const colDefs = gridTemplateColumns.trim().split(/\s+/).filter(Boolean)
      if (colDefs.length > 0 && gridTemplateColumns !== 'none') {
        colsFromCss = colDefs.length
        log('🎯 Columns from CSS grid-template-columns:', colsFromCss, '(value:', gridTemplateColumns, ')')
      } else {
        log('⚠️  No valid grid-template-columns found')
      }
      const cssGap = parsePx(style.getPropertyValue('column-gap'))
      gapFromCss = cssGap > 0 ? cssGap : columnGap
      log('📐 Column gap:', gapFromCss, 'px (from CSS:', cssGap, ', fallback:', columnGap, ')')

      if (!preferCssColumns) {
        let childEl: HTMLElement | null = null
        if (childSelector) {
          childEl = grid.querySelector(childSelector) as HTMLElement | null
        } else {
          childEl = grid.firstElementChild as HTMLElement | null
        }
        const childWidth = childEl?.getBoundingClientRect().width || 0
        if (childWidth > 0) {
          colsFromChild = Math.max(1, Math.round((width + gapFromCss) / (childWidth + gapFromCss)))
          log('👶 Columns from child element width:', colsFromChild, '(child width:', childWidth, 'px)')
        }
      }
    } else {
      log('❌ Grid element not found:', gridSelector)
    }

    const effectiveItem = itemMinWidth + gapFromCss
    const colsFromMin = Math.max(1, Math.floor((width + gapFromCss) / effectiveItem))
    log(
      '🔢 Columns from min width calculation:',
      colsFromMin,
      '(width:',
      width,
      ', itemMinWidth:',
      itemMinWidth,
      ', gap:',
      gapFromCss,
      ')'
    )

    const finalCols = colsFromCss || colsFromChild || colsFromMin
    columns.value = finalCols
    log(
      '✨ Final columns selected:',
      finalCols,
      '(priority: CSS:',
      colsFromCss,
      '→ Child:',
      colsFromChild,
      '→ Min:',
      colsFromMin,
      ')'
    )

    pageSize.value = finalCols * rowMultiple
    log('📦 Final pageSize:', pageSize.value, '=', finalCols, '×', rowMultiple, 'rows')

    if (mode === 'once' && colsFromCss !== null) {
      log('⏹️  Mode is "once" and CSS columns detected, cleaning up observers')
      cleanup()
    } else if (mode === 'once') {
      log('⚠️  Mode is "once" but CSS columns not detected yet, keeping hook active for retry')
    }
  }

  const handleResize = () => {
    if (debounceTimer) {
      window.clearTimeout(debounceTimer)
    }
    debounceTimer = window.setTimeout(calc, debounceMs)
  }

  onMounted(() => {
    log('🚀 Hook initialized with config:', {
      containerSelector,
      gridSelector,
      childSelector,
      itemMinWidth,
      columnGap,
      rowMultiple,
      defaultSize,
      preferCssColumns,
      mode,
      minDeltaPx,
      debounceMs,
    })

    if (mode === 'observe') {
      log('👀 Setting up observers (mode: observe)...')
      const container = document.querySelector(containerSelector)
      if (container && 'ResizeObserver' in window) {
        resizeObserver = new ResizeObserver(() => {
          log('🔔 ResizeObserver triggered')
          handleResize()
        })
        resizeObserver.observe(container as Element)
        log('✅ ResizeObserver attached to container')
      } else {
        log('⚠️  ResizeObserver not available or container not found')
      }
      window.addEventListener('resize', handleResize, { passive: true })
      log('✅ Window resize listener attached')
    } else {
      log('ℹ️  Mode is "once", no observers will be set up. Call recalc() manually.')
    }
  })

  onBeforeUnmount(() => cleanup())

  return { pageSize, columns, containerWidth, recalc: calc, stop: cleanup }
}
