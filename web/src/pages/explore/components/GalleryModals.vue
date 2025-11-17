<script setup lang="ts">
  /* 使用预览器的 FileInfo 类型，确保字段兼容（full_url/description 等） */
  import type { FileInfo as ViewerFileInfo } from '@/components/FileViewer/types'
  /* 直接使用全功能预览器，避免别名映射冲突 */
  import FileViewer from '@/components/FileViewer/index.vue'
  import { computed } from 'vue'

  const props = defineProps<{
    showDetailModal: boolean
    showFilePreview: boolean
    selectedFile: any | null
    previewFileUrl: string
    previewFileName: string
    files?: any[]
    currentPreviewIndex?: number
  }>()

  defineEmits<{
    (e: 'update:showDetailModal', value: boolean): void
    (e: 'update:showFilePreview', value: boolean): void
    (e: 'download', id?: string): void
    (e: 'prev'): void
    (e: 'next'): void
  }>()

  /* 将页面文件数据映射为预览器需要的字段（尽量不改变源数据结构） */
  const mapToViewerFile = (file: any): ViewerFileInfo => ({
    id: String(file.id),
    url: file.full_url || file.url || file.full_thumb_url || file.thumb_url || '',
    full_url: file.full_url || file.url || file.full_thumb_url || file.thumb_url || '',
    display_name: file.display_name || file.original_name || 'file',
    description: file.description || file.ai_info?.description || '',
    width: file.width,
    height: file.height,
    format: file.format,
    size_formatted: file.size_formatted,
    resolution: file.ai_info?.resolution,
    file_type: file.file_type,
    tags: file.ai_info?.tags?.map((tag: string) => ({ id: tag, name: tag })) || [],
    ai_info: file.ai_info
      ? {
          description: file.ai_info.description,
        }
      : undefined,
  })

  const viewerFiles = computed<ViewerFileInfo[]>(() => (props.files || []).map(mapToViewerFile))
  const viewerSelectedFile = computed<ViewerFileInfo | null>(() => {
    try {
      return props.selectedFile ? mapToViewerFile(props.selectedFile) : null
    } catch (_error) {
      return null
    }
  })
</script>

<template>
  <div>
    <CyberFileDetailModal
      :model-value="showDetailModal"
      :fileInfo="selectedFile"
      :is-admin="false"
      @update:model-value="$emit('update:showDetailModal', $event)"
      @download="$emit('download', selectedFile?.id)"
    />

    <FileViewer
      v-if="showFilePreview && viewerSelectedFile"
      :model-value="showFilePreview"
      :file="viewerSelectedFile"
      :files="viewerFiles"
      :initial-index="currentPreviewIndex"
      :show-side-nav="true"
      :show-keyboard-tips="true"
      :controls-hide-timeout="8000"
      search-scope="gallery"
      @update:model-value="$emit('update:showFilePreview', $event)"
      @close="$emit('update:showFilePreview', false)"
      @prev="$emit('prev')"
      @next="$emit('next')"
    />
  </div>
</template>
