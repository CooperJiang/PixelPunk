<script setup lang="ts">
  import GalleryFilterPanel from './GalleryFilterPanel.vue'
  import type { ImageListParams } from '@/api/common'

  const _props = defineProps<{
    show: boolean
    initialFilters: ImageListParams
  }>()

  const emit = defineEmits<{
    (e: 'filter', payload: ImageListParams): void
  }>()
</script>

<template>
  <div class="gallery-header-aligned px-8 lg:px-[max(2rem,calc((100vw-1600px)/2))]">
    <transition
      name="filter-panel"
      enter-active-class="transition duration-300 ease-out"
      leave-active-class="transition duration-200 ease-in"
      enter-from-class="opacity-0 transform -translate-y-6"
      enter-to-class="opacity-100 transform translate-y-0"
      leave-from-class="opacity-100 transform translate-y-0"
      leave-to-class="opacity-0 transform -translate-y-6"
    >
      <GalleryFilterPanel
        v-show="show"
        :show="show"
        :initial-filters="initialFilters"
        class="mb-3"
        @filter="(p) => emit('filter', p)"
      />
    </transition>
  </div>
</template>

<style scoped></style>
