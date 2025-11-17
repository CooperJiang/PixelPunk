<script setup lang="ts">
  import { useTexts } from '@/composables/useTexts'
  import type { FileInfo } from '@/api/admin/files'
  defineOptions({
    name: 'VisualAnalysis',
  })

  const { $t } = useTexts()

  defineProps<{
    fileInfo: FileInfo
  }>()
</script>

<template>
  <div v-if="fileInfo.ai_info" class="ai-info-area">
    <div class="info-section h-full">
      <h3 class="section-title-compact flex items-center">
        <i class="fas fa-eye mr-2" />{{ $t('components.fileDetailModal.visualAnalysis.title') }}
      </h3>
      <div class="info-grid-compact mt-1 space-y-1">
        <div class="info-item-compact">
          <span class="info-label-compact">{{ $t('components.fileDetailModal.visualAnalysis.dominantColor') }}</span>
          <div class="flex items-center">
            <div
              class="mr-2 h-4 w-4 rounded border border-default"
              :style="{ backgroundColor: fileInfo.ai_info.dominant_color }"
            />
            <span class="font-mono text-xs">{{ fileInfo.ai_info.dominant_color }}</span>
          </div>
        </div>

        <div class="info-item-compact">
          <span class="info-label-compact">{{ $t('components.fileDetailModal.visualAnalysis.imageDimensions') }}</span>
          <span class="info-value-compact">{{ fileInfo.width }} × {{ fileInfo.height }}</span>
        </div>

        <div v-if="fileInfo.ai_info.aspect_ratio" class="info-item-compact">
          <span class="info-label-compact">{{ $t('components.fileDetailModal.visualAnalysis.aspectRatio') }}</span>
          <span class="info-value-compact">{{ fileInfo.ai_info.aspect_ratio.toFixed(2) }}</span>
        </div>

        <div class="info-item-compact">
          <span class="info-label-compact">{{ $t('components.fileDetailModal.visualAnalysis.resolution') }}</span>
          <span class="info-value-compact">{{ fileInfo.ai_info.resolution }}</span>
        </div>

        <div v-if="fileInfo.ai_info.color_palette?.length" class="info-item-compact">
          <span class="info-label-compact">{{ $t('components.fileDetailModal.visualAnalysis.colorPalette') }}</span>
          <div class="flex flex-wrap gap-1">
            <div
              v-for="(color, index) in fileInfo.ai_info.color_palette"
              :key="index"
              class="h-3 w-3 rounded border border-default"
              :style="{ backgroundColor: color }"
              :title="color"
            />
          </div>
        </div>

        <div v-if="fileInfo.ai_info.composition" class="info-item-compact">
          <span class="info-label-compact">{{ $t('components.fileDetailModal.visualAnalysis.composition') }}</span>
          <span class="info-value-compact">{{ fileInfo.ai_info.composition }}</span>
        </div>

        <div v-if="fileInfo.ai_info.objects_count" class="info-item-compact">
          <span class="info-label-compact">{{ $t('components.fileDetailModal.visualAnalysis.objectsCount') }}</span>
          <span class="info-value-compact">{{ fileInfo.ai_info.objects_count }} {{ $t('components.fileDetailModal.basicInfo.count') }}</span>
        </div>
      </div>

      <div class="mt-3 border-t border-subtle pt-3">
        <h4 class="mb-1 flex items-center text-sm font-semibold text-content">
          <i class="fas fa-shield-alt mr-2" />{{ $t('components.fileDetailModal.contentSafety.title') }}
        </h4>
        <div class="space-y-1">
          <div class="info-item-compact">
            <span class="info-label-compact">{{ $t('components.fileDetailModal.contentSafety.nsfwStatus') }}</span>
            <CyberTag :variant="fileInfo.ai_info.is_nsfw ? 'error' : 'success'" size="small">
              {{
                fileInfo.ai_info.is_nsfw
                  ? $t('components.fileDetailModal.contentSafety.inappropriate')
                  : $t('components.fileDetailModal.contentSafety.safe')
              }}
            </CyberTag>
          </div>
          <div class="info-item-compact">
            <span class="info-label-compact">{{ $t('components.fileDetailModal.contentSafety.nsfwScore') }}</span>
            <div class="flex-1">
              <div class="mb-1 flex items-center justify-between">
                <span class="text-xs text-content">{{ (fileInfo.ai_info.nsfw_score * 100).toFixed(1) }}%</span>
                <span
                  class="text-xs"
                  :class="
                    fileInfo.ai_info.nsfw_score > 0.5
                      ? 'text-status-error-text'
                      : fileInfo.ai_info.nsfw_score > 0.3
                        ? 'text-badge-warning-text'
                        : 'text-status-success-text'
                  "
                >
                  {{
                    fileInfo.ai_info.nsfw_score > 0.5
                      ? $t('components.fileDetailModal.contentSafety.high')
                      : fileInfo.ai_info.nsfw_score > 0.3
                        ? $t('components.fileDetailModal.contentSafety.medium')
                        : $t('components.fileDetailModal.contentSafety.low')
                  }}
                </span>
              </div>
              <div class="relative h-2 overflow-hidden rounded-full bg-background-700">
                <div
                  class="h-full rounded-full transition-all"
                  :class="
                    fileInfo.ai_info.nsfw_score > 0.5
                      ? 'bg-status-error-border'
                      : fileInfo.ai_info.nsfw_score > 0.3
                        ? 'bg-badge-warning-border'
                        : 'bg-status-success-border'
                  "
                  :style="{ width: `${fileInfo.ai_info.nsfw_score * 100}%` }"
                />
              </div>
            </div>
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
    @apply flex items-start justify-between gap-2 rounded border border-default bg-background-700 p-1.5;
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

  .quality-badge {
    @apply rounded px-1.5 py-0.5 text-xs font-medium capitalize;
  }
</style>
