<script setup lang="ts">
  import { computed, onMounted, onUnmounted, ref } from 'vue'
  import type { FileInfo } from '@/api/common'
  import { useRouter } from 'vue-router'
  import { useTexts } from '@/composables/useTexts'

  const { $t } = useTexts()

  const props = defineProps<{
    files: FileInfo[]
    selectMode?: boolean
    isFileSelected?: (file: FileInfo) => boolean
    isVectorSearch?: boolean
    columnCount?: number
    showTags?: boolean
  }>()

  const emit = defineEmits<{
    (e: 'view', id: string): void
    (e: 'download', file: FileInfo): void
    (e: 'details', id: string): void
    (e: 'select', file: FileInfo): void
  }>()

  const router = useRouter()
  const masonryGrid = ref<HTMLElement>()

  /* 动态计算列数 - 针对全屏优化 */
  const dynamicColumnCount = ref(4)

  const calculateColumns = () => {
    if (!masonryGrid.value) {
      return
    }
    const width = masonryGrid.value.offsetWidth

    const density = getCurrentDensity()

    const baseColumns = getBaseColumns(width)

    if (density === 'compact') {
      dynamicColumnCount.value = Math.min(baseColumns + 2, 12) // 紧凑模式：增加2列，最多12列
    } else if (density === 'comfortable') {
      dynamicColumnCount.value = Math.max(baseColumns - 2, 2) // 舒适模式：减少2列，最少2列
    } else {
      dynamicColumnCount.value = baseColumns // 标准模式：使用基础列数
    }
  }

  const getCurrentDensity = () => {
    const galleryElement = document.querySelector('.gallery-files-fullscreen')
    if (!galleryElement) {
      return 'normal'
    }

    if (galleryElement.classList.contains('density-compact')) {
      return 'compact'
    }
    if (galleryElement.classList.contains('density-comfortable')) {
      return 'comfortable'
    }
    return 'normal'
  }

  const getBaseColumns = (width: number) => {
    if (width < 640) {
      return 2
    } // 移动端
    if (width < 768) {
      return 3
    } // 小平板
    if (width < 1024) {
      return 4
    } // 平板
    if (width < 1280) {
      return 5
    } // 小屏笔记本
    if (width < 1536) {
      return 6
    } // 笔记本
    if (width < 1920) {
      return 7
    } // 桌面
    if (width < 2560) {
      return 8
    } // 大屏幕
    return 10 // 超宽屏
  }

  const columns = computed(() => {
    const cols = props.columnCount || dynamicColumnCount.value
    const result: FileInfo[][] = Array.from({ length: cols }, () => [])

    if (!props.files || !Array.isArray(props.files)) {
      return result
    }

    const columnHeights = new Array(cols).fill(0)

    props.files.forEach((file) => {
      const minHeightIndex = columnHeights.indexOf(Math.min(...columnHeights))
      result[minHeightIndex].push(file)

      const ratio = (file.height || 100) / (file.width || 100)
      columnHeights[minHeightIndex] += ratio
    })

    return result
  })

  const handleFileClick = (file: FileInfo) => {
    if (props.selectMode) {
      emit('select', file)
    } else {
      emit('view', file.id)
    }
  }

  const handleTagClick = (tag: string) => {
    router.push({
      path: '/gallery',
      query: { tags: tag, from: 'gallery' },
    })
  }

  const openInNewWindow = (file: FileInfo) => {
    const url = file.full_url || file.url
    if (url) {
      window.open(url, '_blank', 'noopener,noreferrer')
    }
  }

  const getMaxVisibleTags = (tags: string[]) => {
    if (!tags || tags.length === 0) return 0

    const currentDensity = getCurrentDensity()
    let maxTags = 2 // 默认最多显示2个标签

    if (currentDensity === 'compact') {
      maxTags = 1 // 紧凑模式只显示1个
    } else if (currentDensity === 'comfortable') {
      maxTags = 2 // 舒适模式也显示2个，确保能显示+xx
    }

    return Math.min(maxTags, tags.length)
  }

  const formatFileSize = (bytes: number) => {
    if (bytes < 1024) {
      return `${bytes} B`
    }
    if (bytes < 1024 * 1024) {
      return `${(bytes / 1024).toFixed(1)} KB`
    }
    return `${(bytes / (1024 * 1024)).toFixed(1)} MB`
  }

  const handleImageLoad = (event: Event) => {
    const img = event.target as HTMLImageElement
    img.classList.add('loaded')
  }

  const handleImageError = (event: Event) => {
    const img = event.target as HTMLImageElement
    img.src = '/placeholder-file.png' // 备用文件
  }

  let resizeObserver: ResizeObserver

  onMounted(() => {
    calculateColumns()

    if (masonryGrid.value) {
      resizeObserver = new ResizeObserver(() => {
        calculateColumns()
      })
      resizeObserver.observe(masonryGrid.value)
    }

    window.addEventListener('resize', calculateColumns)
  })

  onUnmounted(() => {
    if (resizeObserver && masonryGrid.value) {
      resizeObserver.unobserve(masonryGrid.value)
    }
    window.removeEventListener('resize', calculateColumns)
  })
</script>

<template>
  <div class="gallery-masonry-container">
    <div ref="masonryGrid" class="masonry-grid">
      <div v-for="(column, colIndex) in columns" :key="`col-${colIndex}`" class="masonry-column">
        <div
          v-for="file in column"
          :key="file.id"
          class="masonry-item selectable-item"
          :class="{
            selected: selectMode && isFileSelected(file),
            'is-selected': selectMode && isFileSelected(file),
            selectable: selectMode,
            'has-tags': showTags && file.ai_info?.tags?.length,
          }"
          @click="handleFileClick(file)"
        >
          <div v-if="selectMode && isFileSelected(file)" class="selection-indicator file-selected-indicator">
            <i class="fas fa-check" />
          </div>

          <div class="file-container">
            <img
              :src="file.full_thumb_url || file.thumb_url"
              :alt="file.display_name || file.original_name"
              class="masonry-file"
              :style="{ aspectRatio: `${file.width}/${file.height}` }"
              loading="lazy"
              @load="handleImageLoad"
              @error="handleImageError"
            />

            <div v-if="!selectMode" class="image-actions">
              <button
                v-if="file.ai_info"
                class="image-action-icon"
                :title="$t('explore.actions.details')"
                @click.stop="$emit('details', file.id)"
              >
                <i class="fas fa-info-circle" />
              </button>
              <button class="image-action-icon" :title="$t('explore.actions.preview')" @click.stop="openInNewWindow(file)">
                <i class="fas fa-external-link-alt" />
              </button>
              <button class="image-action-icon" :title="$t('explore.actions.download')" @click.stop="$emit('download', file)">
                <i class="fas fa-download" />
              </button>
            </div>

            <div
              v-if="isVectorSearch && file.similarity !== undefined"
              class="similarity-badge"
              :class="{ 'with-actions': !selectMode }"
            >
              <i class="fas fa-percentage" />
              {{ Math.round(file.similarity * 100) }}
            </div>

            <div class="hover-overlay">
              <div class="file-info">
                <h3 class="file-title">{{ file.display_name || file.original_name }}</h3>
                <div class="file-meta">
                  <span>{{ file.width }} × {{ file.height }}</span>
                  <span>{{ formatFileSize(file.size) }}</span>
                </div>
              </div>
            </div>
          </div>

          <div v-if="showTags && file.ai_info?.tags?.length" class="file-tags">
            <CyberTag
              v-for="(tag, idx) in file.ai_info.tags.slice(0, getMaxVisibleTags(file.ai_info.tags))"
              :key="idx"
              variant="primary"
              size="small"
              @click.stop="handleTagClick(tag)"
            >
              {{ tag }}
            </CyberTag>
            <CyberTag v-if="file.ai_info.tags.length > getMaxVisibleTags(file.ai_info.tags)" variant="secondary" size="small">
              +{{ file.ai_info.tags.length - getMaxVisibleTags(file.ai_info.tags) }}
            </CyberTag>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
  .gallery-masonry-container {
    width: 100%;
    padding: 0.5rem;
  }

  .masonry-grid {
    display: flex;
    gap: var(--file-gap, 1rem);
    align-items: flex-start;
  }

  @media (min-width: 1920px) {
    .masonry-grid {
      gap: var(--file-gap, 1.25rem);
    }
  }

  @media (min-width: 2560px) {
    .masonry-grid {
      gap: var(--file-gap, 1.5rem);
    }
  }

  .masonry-column {
    flex: 1;
    display: flex;
    flex-direction: column;
    gap: var(--file-gap, 1rem);
  }

  @media (min-width: 1920px) {
    .masonry-column {
      gap: var(--file-gap, 1.25rem);
    }
  }

  @media (min-width: 2560px) {
    .masonry-column {
      gap: var(--file-gap, 1.5rem);
    }
  }

  .masonry-item {
    position: relative;
    border-radius: var(--radius-md);
    overflow: hidden;
    background: var(--color-background-700);
    border: 1px solid var(--color-border-default);
    transition: all 0.3s ease;
    cursor: pointer;
  }

  .masonry-item:hover {
    box-shadow: var(--shadow-cyber-lg);
    border-color: var(--color-hover-border);
  }

  .masonry-item.selectable {
    cursor: pointer;
  }

  .selection-indicator {
    position: absolute;
    top: 8px;
    right: 8px;
    width: 26px;
    height: 26px;
    background: linear-gradient(135deg, var(--color-badge-accent-text), var(--color-brand-500));
    border-radius: var(--radius-full);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 10;
    box-shadow: var(--shadow-glow-md);
    border: 2px solid var(--color-background-700);
    animation: scaleIn var(--transition-fast) var(--ease-out);
  }

  .selection-indicator i {
    color: var(--color-content-heading);
    font-size: 12px;
    text-shadow: 0 1px 2px rgba(0, 0, 0, 0.5);
  }

  .file-container {
    position: relative;
    width: 100%;
    overflow: hidden;
    perspective: 1000px;
  }

  .masonry-file {
    width: 100%;
    height: auto;
    display: block;
    opacity: 0;
    transition:
      opacity 0.3s ease,
      transform 0.5s cubic-bezier(0.34, 1.56, 0.64, 1);
    transform-style: preserve-3d;
    transform-origin: center center;
  }

  .masonry-file.loaded {
    opacity: 1;
  }

  .masonry-item:hover .masonry-file {
    transform: scale(1.15) translateZ(30px);
  }

  .image-actions {
    position: absolute;
    top: 8px;
    left: 8px;
    display: flex;
    gap: 6px;
    z-index: 5;
    opacity: 0;
    transform: translateY(-4px);
    transition: all 0.3s ease;
    pointer-events: none;
  }

  .masonry-item:hover .image-actions {
    opacity: 1;
    transform: translateY(0);
    pointer-events: auto;
  }

  .image-action-icon {
    width: 32px;
    height: 32px;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: var(--radius-sm);
    background: rgba(var(--color-background-900-rgb), 0.9);
    backdrop-filter: blur(12px);
    border: 1.5px solid rgba(var(--color-brand-500-rgb), 0.4);
    color: var(--color-content);
    transition: all 0.2s ease;
    cursor: pointer;
    box-shadow: 0 2px 8px var(--color-overlay-medium);
  }

  .image-action-icon:hover {
    border-color: var(--color-brand-500);
    color: var(--color-brand-500);
    box-shadow: 0 4px 12px rgba(var(--color-brand-500-rgb), 0.5);
  }

  .image-action-icon i {
    font-size: 14px;
    transition: all 0.2s ease;
  }

  .similarity-badge {
    position: absolute;
    top: 8px;
    right: 8px;
    background: linear-gradient(135deg, rgba(var(--color-brand-500-rgb), 0.9), rgba(var(--color-brand-500-rgb), 0.7));
    color: var(--color-text-on-brand);
    padding: 4px 8px;
    border-radius: var(--radius-lg);
    font-size: 11px;
    font-weight: 600;
    display: flex;
    align-items: center;
    gap: 4px;
    backdrop-filter: blur(10px);
    box-shadow: var(--shadow-cyber-sm);
  }

  .similarity-badge.with-actions {
    top: 8px;
  }

  .hover-overlay {
    position: absolute;
    inset: 0;
    background: linear-gradient(to bottom, transparent 60%, rgba(0, 0, 0, 0.9) 100%);
    opacity: 0;
    transition: opacity 0.3s ease;
    display: flex;
    flex-direction: column;
    justify-content: flex-end;
    padding: 0.75rem;
  }

  .masonry-item:hover .hover-overlay {
    opacity: 1;
  }

  .file-info {
    color: var(--color-content-heading);
    position: absolute;
    bottom: 0;
    left: 0;
    right: 0;
    padding: 0.75rem;
    background: transparent;
    z-index: 1;
  }

  .file-title {
    font-size: 13px;
    font-weight: 500;
    margin-bottom: 4px;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    text-shadow: 0 1px 3px rgba(0, 0, 0, 0.5);
  }

  .file-meta {
    display: flex;
    gap: 0.75rem;
    font-size: 11px;
    color: var(--color-content-default);
  }

  .file-meta span {
    display: flex;
    align-items: center;
    gap: 4px;
  }

  .file-tags {
    padding: 6px;
    display: flex;
    flex-wrap: nowrap;
    gap: 4px;
    background: rgba(var(--color-background-900-rgb), 0.5);
    backdrop-filter: blur(8px);
    overflow: hidden;
  }

  @media (max-width: 640px) {
    .masonry-grid {
      gap: var(--file-gap, 0.5rem);
    }

    .masonry-column {
      gap: var(--file-gap, 0.5rem);
    }

    .image-action-icon {
      width: 28px;
      height: 28px;
    }

    .image-action-icon i {
      font-size: 12px;
    }
  }
</style>
