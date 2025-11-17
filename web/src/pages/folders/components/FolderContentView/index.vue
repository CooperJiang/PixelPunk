<script setup lang="ts">
  import type { FolderContentViewEmits, FolderContentViewProps } from './types'
  import FolderGrid from '../FolderGrid/index.vue'
  import FileGrid from '../FileGrid/index.vue'
  import FolderEmptyState from '../FolderEmptyState/index.vue'

  defineOptions({
    name: 'FolderContentView',
  })

  const _props = defineProps<FolderContentViewProps>()
  const _emit = defineEmits<FolderContentViewEmits>()
</script>

<template>
  <div class="mixed-content-panel">
    <div class="content-body">
      <FolderGrid
        v-if="folders.length > 0"
        :folders="folders"
        :select-mode="selectMode && !preview"
        :is-all-folders-selected="isAllFoldersSelected"
        :is-folder-selected="isFolderSelected"
        :preview="preview"
        @folder-click="$emit('folder-click', $event)"
        @folder-double-click="$emit('folder-double-click', $event)"
        @toggle-folder-select="$emit('toggle-folder-select', $event)"
        @toggle-select-all-folders="$emit('toggle-select-all-folders')"
        @invert-folder-selection="$emit('invert-folder-selection')"
        @delete-folder="$emit('delete-folder', $event)"
        @edit-folder="$emit('edit-folder', $event)"
        @toggle-folder-visibility="(...args) => $emit('toggle-folder-visibility', ...args)"
        @folder-drag-start="$emit('folder-drag-start')"
        @folder-drag-end="$emit('folder-drag-end', $event)"
        @folder-moved="$emit('folder-moved', $event[0], $event[1])"
      />

      <FileGrid
        v-if="images.length > 0"
        :images="images"
        :pagination="pagination"
        :select-mode="selectMode && !preview"
        :batch-mode="batchMode && !preview"
        :view-mode="viewMode"
        :selected-batch-images="selectedBatchImages"
        :is-all-images-selected="isAllImagesSelected"
        :is-image-selected="isImageSelected"
        :is-batch-image-selected="isBatchImageSelected"
        :preview="preview"
        @preview-image="$emit('preview-image', $event)"
        @view-details="$emit('view-details', $event)"
        @copy-link="$emit('copy-link', $event)"
        @download-image="$emit('download-image', $event)"
        @delete-image="$emit('delete-image', $event)"
        @toggle-image-visibility="$emit('toggle-image-visibility', $event)"
        @toggle-image-select="$emit('toggle-image-select', $event)"
        @toggle-batch-image-select="$emit('toggle-batch-image-select', $event)"
        @toggle-select-all-images="$emit('toggle-select-all-images')"
        @invert-image-selection="$emit('invert-image-selection')"
        @start-batch-mode="$emit('start-batch-mode')"
        @cancel-batch-mode="$emit('cancel-batch-mode')"
        @batch-delete="$emit('batch-delete')"
        @change-view-mode="$emit('change-view-mode', $event)"
        @image-drag-start="$emit('image-drag-start')"
        @image-drag-end="$emit('image-drag-end', $event)"
        @image-moved="$emit('image-moved', $event[0], $event[1])"
        @image-deleted="$emit('image-deleted', $event)"
      />

      <FolderEmptyState v-if="!folders.length && !images.length" @create-folder="$emit('create-folder')" />
    </div>
  </div>
</template>

<style scoped>
  .mixed-content-panel {
    display: flex;
    flex-direction: column;
    height: 100%;
    min-height: 0;
    background: transparent;
    position: relative;
  }

  .content-body {
    position: relative;
    display: flex;
    flex-direction: column;
    gap: 1.5rem;
    flex: 1;
    min-height: 0;
    max-height: 100%;
    overflow-y: auto;
    scrollbar-width: thin;
    scrollbar-color: var(--color-divider-default) transparent;
  }

  .content-body::-webkit-scrollbar {
    width: 6px;
  }

  .content-body::-webkit-scrollbar-track {
    background: transparent;
  }

  .content-body::-webkit-scrollbar-thumb {
    background: var(--color-divider-default);
    border-radius: var(--radius-sm);
  }

  .content-body::-webkit-scrollbar-thumb:hover {
    background: var(--color-divider-subtle);
  }
</style>
