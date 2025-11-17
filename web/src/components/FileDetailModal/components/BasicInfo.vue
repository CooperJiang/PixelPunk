<script setup lang="ts">
  import { getCurrentLocale } from '@/utils/locale'
  import { useTexts } from '@/composables/useTexts'
  import type { FileInfo } from '@/api/admin/files'
  import { VisibilityLabelMap } from '../types'

  defineOptions({
    name: 'BasicInfo',
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

  const formatDateTime = (dateString: string) => {
    const date = new Date(dateString)
    return date.toLocaleString(getCurrentLocale(), {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit',
    })
  }

  const translateVisibility = (visibility: string) => {
    return (
      $t(`components.fileDetailModal.visibility.${visibility}`) ||
      VisibilityLabelMap[visibility as keyof typeof VisibilityLabelMap] ||
      visibility
    )
  }
</script>

<template>
  <div class="basic-info-area">
    <div class="info-section h-full">
      <h3 class="section-title-compact flex items-center">
        <i class="fas fa-info-circle mr-2" />{{ $t('components.fileDetailModal.basicInfo.title') }}
      </h3>
      <div class="info-grid-compact mt-1 space-y-1">
        <div class="info-item-compact">
          <span class="info-label-compact">{{ $t('components.fileDetailModal.basicInfo.fileId') }}</span>
          <span class="info-value-compact font-mono text-xs">{{ fileInfo.id }}</span>
        </div>

        <div class="info-item-compact">
          <span class="info-label-compact">{{ $t('components.fileDetailModal.basicInfo.originalName') }}</span>
          <span class="info-value-compact truncate" :title="fileInfo.original_name">{{ fileInfo.original_name }}</span>
        </div>

        <div class="info-item-compact">
          <span class="info-label-compact">{{ $t('components.fileDetailModal.basicInfo.fileFormat') }}</span>
          <CyberTag variant="primary" size="small">
            {{ fileInfo.format.toUpperCase() }}
          </CyberTag>
        </div>

        <div class="info-item-compact">
          <span class="info-label-compact">{{ $t('components.fileDetailModal.basicInfo.fileSize') }}</span>
          <span class="info-value-compact">{{ formatFileSize(fileInfo.size) }}</span>
        </div>

        <div class="info-item-compact">
          <span class="info-label-compact">{{ $t('components.fileDetailModal.basicInfo.imageDimensions') }}</span>
          <span class="info-value-compact">{{ fileInfo.width }} × {{ fileInfo.height }}</span>
        </div>

        <div class="info-item-compact">
          <span class="info-label-compact">{{ $t('components.fileDetailModal.basicInfo.visibility') }}</span>
          <CyberTag variant="success" size="small">
            {{ translateVisibility(fileInfo.access_level) }}
          </CyberTag>
        </div>

        <div class="info-item-compact">
          <span class="info-label-compact">{{ $t('components.fileDetailModal.basicInfo.uploader') }}</span>
          <span class="info-value-compact truncate" :title="fileInfo.user_name">{{ fileInfo.user_name }}</span>
        </div>

        <div class="info-item-compact">
          <span class="info-label-compact">{{ $t('components.fileDetailModal.basicInfo.uploadTime') }}</span>
          <span class="info-value-compact text-xs">{{ formatDateTime(fileInfo.created_at) }}</span>
        </div>

        <div v-if="fileInfo.md5_hash" class="info-item-compact">
          <span class="info-label-compact">{{ $t('components.fileDetailModal.basicInfo.fileHash') }}</span>
          <span class="info-value-compact break-all font-mono text-xs">{{ fileInfo.md5_hash.substring(0, 16) }}...</span>
        </div>

        <div v-if="fileInfo.views !== undefined" class="info-item-compact">
          <span class="info-label-compact">{{ $t('components.fileDetailModal.basicInfo.viewCount') }}</span>
          <span class="info-value-compact">{{ fileInfo.views }} {{ $t('components.fileDetailModal.basicInfo.times') }}</span>
        </div>

        <div v-if="fileInfo.is_duplicate || fileInfo.is_recommended || fileInfo.ai_info?.is_nsfw" class="info-item-compact">
          <span class="info-label-compact">{{ $t('components.fileDetailModal.basicInfo.status') }}</span>
          <div class="flex flex-wrap items-center gap-1">
            <CyberTag v-if="fileInfo.is_duplicate" variant="warning" size="small">
              {{ $t('components.fileDetailModal.statusLabels.duplicate') }}
            </CyberTag>
            <CyberTag v-if="fileInfo.is_recommended" variant="primary" size="small">
              {{ $t('components.fileDetailModal.statusLabels.recommended') }}
            </CyberTag>
            <CyberTag v-if="fileInfo.ai_info?.is_nsfw" variant="error" size="small">
              {{ $t('components.fileDetailModal.statusLabels.nsfw') }}
            </CyberTag>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
  .info-section {
    @apply rounded-xl border border-default bg-background-800 p-3;
  }

  .section-title-compact {
    @apply mb-0 pb-1 text-sm font-semibold text-brand-500;
  }

  .info-item-compact {
    @apply flex items-start justify-between gap-2 rounded border border-subtle bg-background-800 p-1.5;
  }

  .info-label-compact {
    @apply min-w-16 flex-shrink-0 text-xs font-medium uppercase tracking-wide text-content-muted;
    letter-spacing: 0.3px;
  }

  .info-value-compact {
    @apply flex-1 break-words text-right text-xs leading-tight text-content;
    min-width: 0; /* 允许 flex 子元素收缩到内容宽度以下 */
  }

  .info-value-compact.truncate {
    @apply overflow-hidden text-ellipsis whitespace-nowrap;
    max-width: 200px; /* 固定最大宽度，确保 truncate 生效 */
  }
</style>
