/**
 * 虚拟滚动工具
 * 用于处理大量数据的高性能滚动

 */
interface VirtualScrollOptions {
  itemHeight: number | ((index: number) => number)
  containerHeight: number
  buffer?: number
  overscan?: number
  scrollContainer?: HTMLElement
}

interface VirtualScrollItem {
  index: number
  top: number
  height: number
  visible: boolean
}

interface VirtualScrollState {
  scrollTop: number
  startIndex: number
  endIndex: number
  visibleItems: VirtualScrollItem[]
  totalHeight: number
}

export class VirtualScrollManager {
  private options: Required<VirtualScrollOptions>
  private state: VirtualScrollState
  private itemHeights: Map<number, number> = new Map()
  private totalItemCount = 0
  private listeners: Array<(state: VirtualScrollState) => void> = []

  constructor(options: VirtualScrollOptions) {
    this.options = {
      buffer: 5,
      overscan: 3,
      scrollContainer: document.documentElement,
      ...options,
    }

    this.state = {
      scrollTop: 0,
      startIndex: 0,
      endIndex: 0,
      visibleItems: [],
      totalHeight: 0,
    }

    this.bindScrollEvent()
  }

  private bindScrollEvent() {
    const handleScroll = () => {
      this.updateScrollTop()
      this.calculateVisibleRange()
      this.notifyListeners()
    }

    this.options.scrollContainer.addEventListener('scroll', handleScroll, { passive: true })
  }

  private updateScrollTop() {
    this.state.scrollTop = this.options.scrollContainer.scrollTop
  }

  private getItemHeight(index: number): number {
    if (typeof this.options.itemHeight === 'function') {
      if (!this.itemHeights.has(index)) {
        const height = this.options.itemHeight(index)
        this.itemHeights.set(index, height)
      }
      return this.itemHeights.get(index) ?? 0
    }
    return this.options.itemHeight as number
  }

  private calculateItemTop(index: number): number {
    if (typeof this.options.itemHeight === 'number') {
      return index * this.options.itemHeight
    }

    let top = 0
    for (let i = 0; i < index; i++) {
      top += this.getItemHeight(i)
    }
    return top
  }

  private calculateVisibleRange() {
    const { scrollTop } = this.state
    const { containerHeight, buffer, overscan } = this.options

    let startIndex = 0
    let endIndex = 0

    if (typeof this.options.itemHeight === 'number') {
      startIndex = Math.floor(scrollTop / this.options.itemHeight)
      endIndex = Math.min(this.totalItemCount - 1, Math.ceil((scrollTop + containerHeight) / this.options.itemHeight))
    } else {
      let currentTop = 0

      for (let i = 0; i < this.totalItemCount; i++) {
        const itemHeight = this.getItemHeight(i)
        if (currentTop + itemHeight > scrollTop) {
          startIndex = i
          break
        }
        currentTop += itemHeight
      }

      currentTop = this.calculateItemTop(startIndex)
      for (let i = startIndex; i < this.totalItemCount; i++) {
        const itemHeight = this.getItemHeight(i)
        if (currentTop > scrollTop + containerHeight) {
          endIndex = i - 1
          break
        }
        currentTop += itemHeight
        endIndex = i
      }
    }

    startIndex = Math.max(0, startIndex - buffer - overscan)
    endIndex = Math.min(this.totalItemCount - 1, endIndex + buffer + overscan)

    this.state.startIndex = startIndex
    this.state.endIndex = endIndex

    this.generateVisibleItems()
  }

  private generateVisibleItems() {
    const visibleItems: VirtualScrollItem[] = []

    for (let i = this.state.startIndex; i <= this.state.endIndex; i++) {
      const top = this.calculateItemTop(i)
      const height = this.getItemHeight(i)

      visibleItems.push({
        index: i,
        top,
        height,
        visible: true,
      })
    }

    this.state.visibleItems = visibleItems
  }

  private calculateTotalHeight() {
    if (typeof this.options.itemHeight === 'number') {
      this.state.totalHeight = this.totalItemCount * this.options.itemHeight
    } else {
      let totalHeight = 0
      for (let i = 0; i < this.totalItemCount; i++) {
        totalHeight += this.getItemHeight(i)
      }
      this.state.totalHeight = totalHeight
    }
  }

  private notifyListeners() {
    this.listeners.forEach((listener) => listener(this.state))
  }

  setItemCount(count: number) {
    this.totalItemCount = count
    this.calculateTotalHeight()
    this.calculateVisibleRange()
    this.notifyListeners()
  }

  updateItemHeight(index: number, height: number) {
    if (typeof this.options.itemHeight === 'function') {
      this.itemHeights.set(index, height)
      this.calculateTotalHeight()
      this.calculateVisibleRange()
      this.notifyListeners()
    }
  }

  scrollToIndex(index: number, alignment: 'start' | 'center' | 'end' = 'start') {
    const itemTop = this.calculateItemTop(index)
    const itemHeight = this.getItemHeight(index)
    let scrollTop = itemTop

    switch (alignment) {
      case 'center':
        scrollTop = itemTop + itemHeight / 2 - this.options.containerHeight / 2
        break
      case 'end':
        scrollTop = itemTop + itemHeight - this.options.containerHeight
        break
    }

    scrollTop = Math.max(0, Math.min(scrollTop, this.state.totalHeight - this.options.containerHeight))
    this.options.scrollContainer.scrollTop = scrollTop
  }

  getState(): VirtualScrollState {
    return { ...this.state }
  }

  subscribe(listener: (state: VirtualScrollState) => void) {
    this.listeners.push(listener)
    return () => {
      const index = this.listeners.indexOf(listener)
      if (index > -1) {
        this.listeners.splice(index, 1)
      }
    }
  }

  destroy() {
    this.listeners.length = 0
    this.itemHeights.clear()
  }
}

export function useVirtualScroll(options: VirtualScrollOptions) {
  const scrollManager = ref<VirtualScrollManager | null>(null)
  const state = ref<VirtualScrollState>({
    scrollTop: 0,
    startIndex: 0,
    endIndex: 0,
    visibleItems: [],
    totalHeight: 0,
  })

  onMounted(() => {
    scrollManager.value = new VirtualScrollManager(options)

    const unsubscribe = scrollManager.value.subscribe((newState) => {
      state.value = newState
    })

    onUnmounted(() => {
      unsubscribe()
      scrollManager.value?.destroy()
    })
  })

  const setItemCount = (count: number) => {
    scrollManager.value?.setItemCount(count)
  }

  const updateItemHeight = (index: number, height: number) => {
    scrollManager.value?.updateItemHeight(index, height)
  }

  const scrollToIndex = (index: number, alignment?: 'start' | 'center' | 'end') => {
    scrollManager.value?.scrollToIndex(index, alignment)
  }

  return {
    state: readonly(state),
    setItemCount,
    updateItemHeight,
    scrollToIndex,
    scrollManager: readonly(scrollManager),
  }
}

export function createVirtualScrollProps() {
  return {
    items: {
      type: Array,
      required: true,
    },
    itemHeight: {
      type: [Number, Function],
      required: true,
    },
    containerHeight: {
      type: Number,
      required: true,
    },
    buffer: {
      type: Number,
      default: 5,
    },
    overscan: {
      type: Number,
      default: 3,
    },
  }
}
