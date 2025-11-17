<script setup lang="ts">
  import GalleryToolbarControls from './GalleryToolbarControls.vue'
  import CompactSearchBar from './CompactSearchBar.vue'
  import ShareButton from '@/components/Share/ShareButton.vue'
  import type { FileInfo } from '@/api/common'

  type Density = 'compact' | 'normal' | 'comfortable'

  const _props = defineProps<{
    images: FileInfo[]
    hasContent: boolean
    setShareButtonRef: (el: any) => void

    filtersSort: string
    sortOptions: { value: string; label: string; icon: string }[]
    currentDensity: Density
    densityOptions: { value: Density; label: string; icon: string }[]

    searchKeyword: string
    searchMode: 'normal' | 'vector'
    showFilter: boolean
    filterCount: number
    isSearching: boolean

    isLeftLayout: boolean
  }>()

  const emit = defineEmits<{
    (e: 'update:select-mode', value: boolean): void
    (e: 'update:share-dialog-visible', value: boolean): void

    (e: 'sort', value: string): void
    (e: 'change-density', value: Density): void

    (e: 'update:search-keyword', value: string): void
    (e: 'update:search-mode', value: 'normal' | 'vector'): void
    (e: 'update:show-filter', value: boolean): void
    (e: 'search', payload: { keyword: string; mode: 'normal' | 'vector' }): void
  }>()
</script>

<template>
  <div
    class="gallery-header-aligned"
    :class="{
      'px-4 sm:px-6 lg:px-8 xl:px-12 2xl:px-16': !isLeftLayout,
      'px-0': isLeftLayout,
    }"
  >
    <div
      class="gallery-header-inner flex flex-col md:flex-row md:items-center"
      :class="{
        'md:justify-between': !isLeftLayout,
        'md:justify-end': isLeftLayout,
        'gallery-header-transparent': !hasContent,
      }"
    >
      <div v-if="!isLeftLayout && hasContent">
        <h1 class="text-lg font-bold text-content md:text-xl">{{ $t('explore.header.title') }}</h1>
        <p class="whitespace-nowrap text-xs text-content-muted">{{ $t('explore.header.subtitle') }}</p>
      </div>

      <div v-else-if="!isLeftLayout && !hasContent"></div>

      <div class="mt-3 flex items-center gap-3 md:mt-0">
        <div class="toolbar-section flex items-center gap-3">
          <ShareButton
            :ref="setShareButtonRef"
            :images="images"
            :require-login="true"
            @update:select-mode="(v: boolean) => emit('update:select-mode', v)"
            @update:share-dialog-visible="(v: boolean) => emit('update:share-dialog-visible', v)"
          />

          <GalleryToolbarControls
            class="toolbar-group"
            :filters-sort="filtersSort"
            :sort-options="sortOptions"
            :current-density="currentDensity"
            :density-options="densityOptions"
            @sort="(v: string) => emit('sort', v)"
            @change-density="(v: Density) => emit('change-density', v)"
          />
        </div>

        <div class="search-section flex-shrink-0">
          <CompactSearchBar
            :model-value="searchKeyword"
            :search-mode="searchMode"
            :show-filter="showFilter"
            :filter-count="filterCount"
            :is-searching="isSearching"
            @update:modelValue="(v: string) => emit('update:search-keyword', v)"
            @update:search-mode="(v: 'normal' | 'vector') => emit('update:search-mode', v)"
            @update:show-filter="(v: boolean) => emit('update:show-filter', v)"
            @search="(p) => emit('search', p)"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped></style>
