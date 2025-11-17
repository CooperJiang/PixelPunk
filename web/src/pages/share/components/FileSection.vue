<script setup lang="ts">
  import { computed } from 'vue'
  import SortBar from './SortBar.vue'
  import SelectionToolbar from './SelectionToolbar.vue'
  import AutoLoader from './AutoLoader.vue'
  import { getShareSortOptions } from '@/constants/share'
  import { useTexts } from '@/composables/useTexts'

  const { $t } = useTexts()

  const SHARE_SORT_OPTIONS = computed(() => getShareSortOptions($t))

  const props = defineProps({
    images: {
      type: Array,
      default: () => [],
    },
    displayImages: {
      type: Array,
      default: () => [],
    },
    isRootImages: {
      type: Boolean,
      default: false,
    },
    layoutMode: {
      type: String,
      default: 'grid',
    },
    isSelectionMode: {
      type: Boolean,
      default: false,
    },
    selectedIds: {
      type: Array,
      default: () => [],
    },
    initialLoadCount: {
      type: Number,
      default: 20,
    },
    currentSort: {
      type: String,
      default: 'created_at',
    },
    sortDirection: {
      type: String,
      default: 'desc',
    },
    waterfallColumns: {
      type: Number,
      default: 4,
    },
    waterfallColumnWidth: {
      type: Number,
      default: 240,
    },
    showLoadMoreButton: {
      type: Boolean,
      default: false,
    },
    currentPage: {
      type: Number,
      default: 1,
    },
    pageSize: {
      type: Number,
      default: 50,
    },
    getImageHeight: {
      type: Function,
      default: () => 200,
    },
    isFullScreen: {
      type: Boolean,
      default: false,
    },
    isLoading: {
      type: Boolean,
      default: false,
    },
    canLoadMore: {
      type: Boolean,
      default: true,
    },
    shareKey: {
      type: String,
      required: true,
    },
    accessToken: {
      type: String,
      default: '',
    },
  })

  const _emit = defineEmits([
    'image-click',
    'toggle-selection',
    'toggle-selection-mode',
    'select-all',
    'deselect-all',
    'download-selected',
    'sort-change',
    'open-grid-size',
    'toggle-immersive',
    'open-layout',
    'load-more',
    'toggle-fullscreen',
  ])

  const selectedCount = computed(() => props.selectedIds.length)

  const isImageSelected = (fileId: string) => {
    const isSelected = props.selectedIds.includes(fileId)
    return isSelected
  }

  const sortOptions = SHARE_SORT_OPTIONS
</script>

<template>
  <div v-if="images && images.length > 0" class="share-section images-section" :class="{ 'root-images': isRootImages }">
    <div class="section-header">
      <div v-if="isRootImages" class="section-title">
        <i class="fas fa-images" />
        <span>{{ $t('share.file.rootFiles') }}</span>
        <span class="count-badge">{{ images.length }}</span>
      </div>

      <SortBar
        v-if="!isRootImages"
        :current-sort="currentSort"
        :sort-direction="sortDirection"
        :sort-options="sortOptions"
        :layout-mode="layoutMode"
        :is-selection-mode="isSelectionMode"
        :is-full-screen="isFullScreen"
        :title="$t('share.file.title')"
        :count="images.length"
        class="fixed-sort-bar"
        @sort-change="$emit('sort-change', $event)"
        @open-grid-size="$emit('open-grid-size', $event)"
        @toggle-immersive="$emit('toggle-immersive')"
        @open-layout="$emit('open-layout', $event)"
        @toggle-selection="$emit('toggle-selection-mode')"
        @toggle-fullscreen="$emit('toggle-fullscreen')"
      />
    </div>

    <SelectionToolbar
      v-if="!isRootImages && isSelectionMode"
      :visible="isSelectionMode"
      :selected-count="selectedCount"
      :total-count="images.length"
      @select-all="$emit('select-all')"
      @deselect-all="$emit('deselect-all')"
      @download="$emit('download-selected')"
    />

    <div v-if="layoutMode === 'grid'" class="images-grid">
      <CyberShareFile
        v-for="image in displayImages"
        :key="image.id"
        :file="image"
        :share-key="shareKey"
        :access-token="accessToken"
        :selectable="isSelectionMode"
        :selected="isImageSelected(image.id)"
        @click="$emit('image-click', image)"
        @select="$emit('toggle-selection', image.id)"
      />
    </div>

    <CyberWaterfallLayout
      v-else-if="layoutMode === 'waterfall'"
      :items="displayImages"
      :column-count="waterfallColumns"
      :column-width="Number(waterfallColumnWidth)"
      :selectable="isSelectionMode"
      :selected-ids="selectedIds"
      :initial-load-count="initialLoadCount"
      :lazy-load="true"
      :share-key="shareKey"
      :access-token="accessToken"
      class="images-waterfall"
      @image-click="$emit('image-click', $event)"
      @select="$emit('toggle-selection', $event)"
    />

    <div v-else-if="layoutMode === 'masonry'" class="images-grid masonry">
      <CyberShareFile
        v-for="image in displayImages"
        :key="image.id"
        :file="image"
        :share-key="shareKey"
        :access-token="accessToken"
        :selectable="isSelectionMode"
        :selected="isImageSelected(image.id)"
        @click="$emit('image-click', image)"
        @select="$emit('toggle-selection', image.id)"
      />
    </div>

    <div v-else-if="layoutMode === 'large'" class="images-grid large">
      <CyberShareFile
        v-for="image in displayImages"
        :key="image.id"
        :file="image"
        :share-key="shareKey"
        :access-token="accessToken"
        :selectable="isSelectionMode"
        :selected="isImageSelected(image.id)"
        @click="$emit('image-click', image)"
        @select="$emit('toggle-selection', image.id)"
      />
    </div>

    <AutoLoader
      v-if="!isRootImages && showLoadMoreButton"
      :current="Math.min(currentPage * pageSize, images.length)"
      :total="images.length"
      :external-loading="isLoading"
      :can-load-more="canLoadMore"
      @load-more="$emit('load-more')"
    />
  </div>
</template>

<style scoped>
  .fixed-sort-bar {
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    z-index: var(--z-index-dropdown);
  }
  .share-section {
    background: var(--color-background-800);
    border-radius: var(--radius-md);
    padding: var(--space-md);
    margin-bottom: var(--space-md);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.15);
    box-shadow: 0 2px 8px rgba(var(--color-background-900-rgb), 0.3);
    transition: all var(--transition-normal) var(--ease-in-out);
  }

  .share-section:hover {
    border-color: rgba(var(--color-brand-500-rgb), 0.25);
    box-shadow: 0 4px 12px rgba(var(--color-background-900-rgb), 0.4);
  }

  .section-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: var(--space-md);
    position: relative;
    padding: 0 0.5rem 0.75rem;
    border-bottom: 1px solid rgba(var(--color-brand-500-rgb), 0.15);
  }

  .section-title {
    display: flex;
    align-items: center;
    color: rgba(var(--color-content-default-rgb), 0.95);
    font-size: 1.1rem;
    font-weight: 600;
    gap: var(--space-sm);
  }

  .section-title i {
    color: var(--color-brand-500);
    font-size: 1.1rem;
  }

  .count-badge {
    background: rgba(var(--color-brand-500-rgb), 0.15);
    color: var(--color-brand-500);
    border-radius: var(--radius-md);
    padding: var(--space-md) 0.5rem;
    font-size: 0.75rem;
    font-weight: 500;
  }

  .images-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(var(--grid-item-size, 220px), 1fr));
    gap: var(--space-sm);
    width: 100%;
    transition: all var(--transition-normal) var(--ease-in-out);
  }

  .images-grid.masonry {
    display: flex;
    flex-flow: row wrap;
    align-items: flex-start;
    gap: var(--space-sm);
    justify-content: flex-start;
  }

  .images-grid.masonry > * {
    width: var(--grid-item-size, 180px);
    margin-bottom: var(--space-md);
    flex-grow: 0;
    flex-shrink: 0;
  }

  .images-grid.large {
    grid-template-columns: repeat(auto-fill, minmax(calc(var(--grid-item-size, 150px) * 1.8), 1fr));
    gap: var(--space-sm);
  }

  .images-waterfall {
    width: 100%;
    margin-top: var(--space-md);
    margin-bottom: var(--space-md);
  }

  .images-waterfall :deep(.cyber-share-file) {
    height: auto;
    overflow: hidden;
    margin-bottom: var(--space-md);
  }

  @media (min-width: 1600px) {
    .images-grid {
      gap: var(--space-sm);
    }

    .images-grid.masonry {
      gap: var(--space-sm);
    }

    .images-grid.masonry > * {
      margin-bottom: var(--space-md);
    }
  }

  @media (min-width: 1920px) {
    .images-grid.large {
      grid-template-columns: repeat(auto-fill, minmax(calc(var(--grid-item-size, 150px) * 2), 1fr));
    }
  }

  @media (max-width: 768px) {
    .share-section {
      padding: var(--space-md);
      margin-bottom: var(--space-md);
    }

    .section-header {
      flex-direction: column;
      align-items: flex-start;
    }

    .images-grid,
    .images-grid.masonry,
    .images-grid.large {
      gap: var(--space-sm);
    }

    .images-grid.masonry > * {
      width: calc(50% - 0.5rem);
    }
  }
</style>
