<script setup lang="ts">
  import { nextTick, onBeforeUnmount, onMounted, reactive, ref, watch } from 'vue'
  import { debounce } from '@/utils/common'
  import { LAYOUT, TIMING } from '@/constants'
  import WaterfallItem from './components/WaterfallItem.vue'
  import type { WaterfallLayoutEmits, WaterfallLayoutProps } from './types'

  const props = withDefaults(defineProps<WaterfallLayoutProps>(), {
    items: () => [],
    columnCount: LAYOUT.WATERFALL.DEFAULT_COLUMNS,
    columnWidth: LAYOUT.WATERFALL.DEFAULT_COLUMN_WIDTH,
    selectable: false,
    selectedIds: () => [],
    lazyLoad: true,
    initialLoadCount: LAYOUT.WATERFALL.INITIAL_LOAD_COUNT,
    gap: LAYOUT.WATERFALL.DEFAULT_GAP,
  })

  const emit = defineEmits<WaterfallLayoutEmits>()

  const containerRef = ref(null)
  const columns = ref<any[][]>([])
  const columnHeights = ref<number[]>([])
  const loadedItems = reactive(new Set<string>())
  const expandedTags = reactive(new Set<string>())
  const currentVisibleCount = ref(props.initialLoadCount)

  const initColumns = () => {
    columns.value = Array.from({ length: props.columnCount }, () => [])
    columnHeights.value = Array.from({ length: props.columnCount }, () => 0)
  }

  const calculateLayout = () => {
    initColumns()
    const effectiveColumnWidth = Math.max(props.columnWidth * 1.2, 180)

    props.items.forEach((item, index) => {
      const processedItem = {
        ...item,
        loaded: loadedItems.has(item.id),
        visible: index < currentVisibleCount.value || loadedItems.has(item.id),
        showAllTags: expandedTags.has(item.id),
      }

      const minHeightIndex = columnHeights.value.indexOf(Math.min(...columnHeights.value))
      columns.value[minHeightIndex].push(processedItem)

      if (processedItem.width && processedItem.height) {
        const aspectRatio = processedItem.width / processedItem.height
        const calculatedHeight = effectiveColumnWidth / aspectRatio + LAYOUT.WATERFALL.INFO_AREA_HEIGHT
        columnHeights.value[minHeightIndex] += calculatedHeight
      } else {
        columnHeights.value[minHeightIndex] += LAYOUT.WATERFALL.DEFAULT_IMAGE_HEIGHT
      }
    })
  }

  const handleResize = () => {
    const width = window.innerWidth
    const newColumns = Math.floor(width / LAYOUT.WATERFALL.DEFAULT_COLUMN_WIDTH) || 1
    if (newColumns !== columns.value.length) {
      calculateLayout()
    }
  }

  const handleImageLoad = (item) => {
    loadedItems.add(item.id)
    item.loaded = true
    emit('image-load', item)
  }

  const handleImageError = (item) => {
    item.hasError = true
  }

  const handleImageClick = (item) => {
    if (props.selectable) {
      emit('select', item.id)
    } else {
      emit('image-click', item)
    }
  }

  const handleToggleTags = (id) => {
    if (expandedTags.has(id)) {
      expandedTags.delete(id)
    } else {
      expandedTags.add(id)
    }
    calculateLayout()
  }

  const checkLoadMore = () => {
    if (!props.lazyLoad || currentVisibleCount.value >= props.items.length) return

    const container = containerRef.value
    if (!container) return

    const containerBottom = container.getBoundingClientRect().bottom
    const windowHeight = window.innerHeight
    const threshold = windowHeight * 2

    if (containerBottom - windowHeight < threshold) {
      const loadIncrement = Math.min(LAYOUT.WATERFALL.MIN_LOAD_COUNT, props.items.length - currentVisibleCount.value)
      currentVisibleCount.value = Math.min(currentVisibleCount.value + loadIncrement, props.items.length)
      calculateLayout()
    }
  }

  const handleScroll = debounce(() => {
    checkLoadMore()
  }, TIMING.SCROLL.DEBOUNCE_TIME)

  watch(() => props.items, calculateLayout, { deep: true })
  watch(() => props.columnCount, calculateLayout)

  onMounted(() => {
    calculateLayout()
    window.addEventListener('resize', handleResize)
    if (props.lazyLoad) {
      window.addEventListener('scroll', handleScroll)
    }
    nextTick(() => checkLoadMore())
  })

  onBeforeUnmount(() => {
    window.removeEventListener('resize', handleResize)
    if (props.lazyLoad) {
      window.removeEventListener('scroll', handleScroll)
    }
  })
</script>

<template>
  <div ref="containerRef" class="waterfall-layout" :style="{ gap: `${gap / 2}px` }">
    <div
      v-for="(column, colIndex) in columns"
      :key="colIndex"
      class="waterfall-column"
      :style="{ width: `${100 / columnCount}%`, padding: `0 ${gap / 4}px` }"
    >
      <WaterfallItem
        v-for="item in column"
        :key="item.id"
        :item="item"
        :selectable="selectable"
        :selected-ids="selectedIds"
        :column-width="columnWidth"
        :gap="gap"
        @image-click="handleImageClick"
        @image-load="handleImageLoad"
        @image-error="handleImageError"
        @toggle-tags="handleToggleTags"
      />
    </div>
  </div>
</template>

<style scoped>
  .waterfall-layout {
    @apply flex w-full;
    align-items: flex-start;
  }

  .waterfall-column {
    @apply flex flex-col;
  }
</style>
