<script setup lang="ts">
  import { useTexts } from '@/composables/useTexts'
  import type { FileInfo } from '@/api/admin/files'
  import { File } from '@/components/File'

  defineOptions({
    name: 'FilePreview',
  })

  const { $t } = useTexts()

  defineProps<{
    fileInfo: FileInfo
  }>()

  /* 格式化文件大小 */
  const formatFileSize = (bytes: number) => {
    if (bytes === 0) {
      return '0 B'
    }

    const k = 1024
    const sizes = ['B', 'KB', 'MB', 'GB']
    const i = Math.floor(Math.log(bytes) / Math.log(k))

    return `${parseFloat((bytes / Math.pow(k, i)).toFixed(2))} ${sizes[i]}`
  }

  const _handlePreviewError = (e: Event) => {
    const img = e.target as HTMLImageElement
    img.style.display = 'none'
  }
</script>

<template>
  <div class="preview-area flex h-full flex-col">
    <div class="image-wrapper relative h-[340px] flex-shrink-0 overflow-hidden rounded-lg border border-default">
      <File
        :src="fileInfo.full_url || fileInfo.url"
        :alt="fileInfo.display_name || fileInfo.original_name"
        fit-mode="contain"
        :maintain-aspect-ratio="true"
        class="h-full w-full"
      />
      <div class="image-info-card absolute bottom-2 left-2 right-2 rounded px-3 py-2">
        <div class="flex justify-between text-xs text-content">
          <span>{{ fileInfo.width }} × {{ fileInfo.height }}</span>
          <span>{{ formatFileSize(fileInfo.size) }}</span>
        </div>
        <div class="mt-1 text-xs text-content-muted">
          {{ fileInfo.format?.toUpperCase() || 'UNKNOWN' }} ·
          {{ fileInfo.ai_info?.resolution || $t('components.fileDetailModal.basicInfo.unknownResolution') }}
        </div>
      </div>
    </div>

    <div v-if="fileInfo.ai_info?.tags?.length" class="tags-section mt-2">
      <h3 class="section-title-compact flex items-center"><i class="fas fa-tags mr-2" />{{ $t('components.fileDetailModal.ai.tags') }}</h3>
      <div class="tags-container tags-scroll-area mt-1 flex flex-wrap content-start gap-1.5">
        <CyberTag v-for="(tag, index) in fileInfo.ai_info.tags" :key="index" variant="primary" size="small" truncate>
          {{ tag }}
        </CyberTag>
      </div>
    </div>
  </div>
</template>

<style scoped>
  .section-title-compact {
    @apply mb-0 pb-1 text-sm font-semibold text-brand-500;
  }

  .tags-section {
    @apply flex flex-col rounded-xl border border-default bg-background-800 p-3;
  }

  .tags-scroll-area {
    @apply overflow-y-auto;
    scrollbar-width: thin;
    max-height: 120px;
  }

  .tags-scroll-area::-webkit-scrollbar {
    @apply w-1;
  }

  .tags-scroll-area::-webkit-scrollbar-track {
    @apply bg-transparent;
  }

  .tags-scroll-area::-webkit-scrollbar-thumb {
    @apply rounded bg-brand-500;
  }

  .tags-scroll-area::-webkit-scrollbar-thumb:hover {
    @apply bg-brand-400;
  }

  .image-info-card {
    background: rgba(var(--color-background-800-rgb), 0.7);
    backdrop-filter: blur(8px);
    border: 1px solid rgba(var(--color-border-default-rgb), 0.5);
  }
</style>
