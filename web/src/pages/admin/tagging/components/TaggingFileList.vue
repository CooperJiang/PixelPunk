<script setup lang="ts">
  import { getCurrentLocale } from '@/utils/locale'
  import { computed, onMounted, onUnmounted, ref } from 'vue'
  import type { ImageTaggingInfo } from '@/api/admin/tagging'
  import { useTexts } from '@/composables/useTexts'

  const { $t } = useTexts()

  const props = defineProps<{
    images: ImageTaggingInfo[]
    isLoading: boolean
    total: number
    currentPage: number
    pageSize: number
    processing: boolean
    selectedIds: string[]
  }>()

  const imageLoadingStates = ref<Record<string, boolean>>({})

  const emit = defineEmits<{
    (e: 'page-change', page: number): void
    (e: 'page-size-change', size: number): void
    (e: 'toggle-select', fileId: string): void
    (e: 'toggle-select-all', checked: boolean): void
    (e: 'ignore-selected', reason: string): void
    (e: 'unignore-selected'): void
    (e: 'retry', fileIds: string[]): void
  }>()

  function debounce<T extends unknown[]>(fn: (...args: T) => void, delay: number) {
    let timer: number | null = null
    return function (this: unknown, ...args: T) {
      if (timer) {
        clearTimeout(timer)
      }
      timer = setTimeout(() => {
        fn.apply(this, args)
        timer = null
      }, delay) as unknown as number
    }
  }

  const isPageAllSelected = computed(() => {
    if (!props.images || props.images.length === 0) {
      return false
    }
    return props.images.every((img) => props.selectedIds.includes(img.id))
  })

  const toggleSelectAll = (checked: boolean) => {
    emit('toggle-select-all', checked)
  }

  const tableContainer = ref<HTMLElement | null>(null)
  const columnWidths = ref({
    select: 32, // 新增：选择列固定宽度
    index: 40, // 固定宽度
    thumb: 70, // 固定宽度
    filename: 0, // 动态计算 - 权重大
    id: 0, // 动态计算 - 权重中
    dimension: 0, // 动态计算 - 权重小
    size: 0, // 动态计算 - 权重小
    format: 0, // 动态计算 - 权重小
    uploader: 0, // 动态计算 - 权重小
    status: 0, // 动态计算 - 权重小
    tries: 0, // 动态计算 - 权重小
    duration: 0, // 动态计算 - 权重中 - 耗时列
    time: 0, // 动态计算 - 权重中
    actions: 40, // 固定宽度 - 操作按钮
  })

  const columnWeights = {
    filename: 4, // 文件名更重要，给最多空间
    id: 2.5, // ID列需要较多空间
    dimension: 1.5,
    size: 1,
    format: 1,
    uploader: 1.2,
    status: 1.2,
    tries: 1,
    duration: 1.5, // 耗时列需要一定空间
    time: 1.8, // 时间列需要一定空间
  }

  const calculateColumnWidths = () => {
    if (!tableContainer.value) {
      return
    }

    const totalWidth = tableContainer.value.clientWidth
    if (totalWidth <= 0) {
      return
    }

    const paddingWidth = 32 // px-4 = 16px * 2

    const fixedWidth = columnWidths.value.select + columnWidths.value.index + columnWidths.value.thumb + columnWidths.value.actions

    const gapWidth = 13 * 16 // gap-4 = 16px (13 columns now)

    const availableWidth = totalWidth - fixedWidth - gapWidth - paddingWidth
    if (availableWidth <= 0) {
      return
    }

    const totalWeight = Object.values(columnWeights).reduce((sum, weight) => sum + weight, 0)

    const tempWidths: Record<string, number> = {}

    tempWidths.filename = Math.floor(availableWidth * (columnWeights.filename / totalWeight))
    tempWidths.id = Math.floor(availableWidth * (columnWeights.id / totalWeight))
    tempWidths.dimension = Math.floor(availableWidth * (columnWeights.dimension / totalWeight))
    tempWidths.size = Math.floor(availableWidth * (columnWeights.size / totalWeight))
    tempWidths.format = Math.floor(availableWidth * (columnWeights.format / totalWeight))
    tempWidths.uploader = Math.floor(availableWidth * (columnWeights.uploader / totalWeight))
    tempWidths.status = Math.floor(availableWidth * (columnWeights.status / totalWeight))
    tempWidths.tries = Math.floor(availableWidth * (columnWeights.tries / totalWeight))
    tempWidths.duration = Math.floor(availableWidth * (columnWeights.duration / totalWeight))
    tempWidths.time = Math.floor(availableWidth * (columnWeights.time / totalWeight))

    const allocatedWidth = Object.values(tempWidths).reduce((sum, width) => sum + width, 0)

    if (allocatedWidth > availableWidth) {
      const scaleFactor = availableWidth / allocatedWidth
      Object.keys(tempWidths).forEach((key) => {
        tempWidths[key] = Math.floor(tempWidths[key] * scaleFactor)
      })
    }

    const newAllocatedWidth = Object.values(tempWidths).reduce((sum, width) => sum + width, 0)
    const remainingWidth = availableWidth - newAllocatedWidth
    if (remainingWidth > 0) {
      tempWidths.filename += remainingWidth
    }

    Object.keys(tempWidths).forEach((key) => {
      columnWidths.value[key as keyof typeof columnWidths.value] = tempWidths[key]
    })

    if (tableContainer.value) {
      Object.entries(columnWidths.value).forEach(([key, value]) => {
        tableContainer.value?.style.setProperty(`--col-${key}-width`, `${value}px`)
      })
    }
  }

  const handleResize = debounce(() => {
    calculateColumnWidths()
  }, 200)

  onMounted(() => {
    window.addEventListener('resize', handleResize)
    setTimeout(() => {
      calculateColumnWidths()
      setTimeout(() => {
        calculateColumnWidths()
      }, 100)
    }, 0)
  })

  onUnmounted(() => {
    window.removeEventListener('resize', handleResize)
  })

  const formatTime = (timeString: string): string => {
    try {
      const date = new Date(timeString)
      return date.toLocaleString(getCurrentLocale(), {
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit',
        hour12: false,
      })
    } catch {
      return timeString
    }
  }

  const formatFileSize = (bytes?: number): string => {
    if (!bytes) {
      return $t('admin.tagging.fileList.unknown')
    }

    const units = ['B', 'KB', 'MB', 'GB']
    let size = bytes
    let unitIndex = 0

    while (size >= 1024 && unitIndex < units.length - 1) {
      size /= 1024
      unitIndex++
    }

    return `${size.toFixed(1)} ${units[unitIndex]}`
  }

  const getStatusText = (status: string): string => {
    return $t(`admin.tagging.status.${status}`) || status
  }

  const getStatusClass = (status: string): string => {
    const statusClassMap: Record<string, string> = {
      none: 'bg-background-600 text-content-heading border border-default',
      pending: 'bg-error-200 text-error-500 border border-error-500',
      done: 'bg-green-900/30 text-green-400 border border-green-500/50',
      failed: 'bg-red-900/30 text-red-400 border border-red-500/50',
      skipped: 'bg-background-500 text-content-muted border border-default',
      ignored: 'bg-background-700 text-content-content-muted border border-subtle',
    }
    return statusClassMap[status] || ''
  }

  const timeString = (timeString: string): string => {
    try {
      const date = new Date(timeString)
      return date.toLocaleString(getCurrentLocale(), {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit',
        second: '2-digit',
        hour12: false,
      })
    } catch {
      return timeString
    }
  }
</script>

<template>
  <div ref="tableContainer" class="tagging-image-list">
    <div class="mb-3 flex items-center justify-between">
      <h2 class="module-title"><i class="fas fa-clipboard-list mr-2"></i>{{ $t('admin.tagging.fileList.title') }}</h2>
      <div v-if="selectedIds.length > 0" class="flex gap-2">
        <CyberButton type="secondary" size="small" @click="$emit('ignore-selected', $t('admin.tagging.actions.ignore'))">
          <i class="fas fa-ban mr-1.5" />{{ $t('admin.tagging.actions.ignore') }}（{{ selectedIds.length }}）
        </CyberButton>
        <CyberButton type="primary" size="small" @click="$emit('unignore-selected')">
          <i class="fas fa-undo mr-1.5" />{{ $t('admin.tagging.actions.unignore') }}（{{ selectedIds.length }}）
        </CyberButton>
      </div>
    </div>

    <div
      class="image-list-header grid-cols-table mb-2 grid gap-4 rounded border border-subtle px-4 py-2 text-xs font-medium"
      style="background: rgba(var(--color-background-800-rgb), 0.5); color: var(--color-content-muted)"
    >
      <div class="col-select flex items-center justify-center overflow-hidden">
        <CyberCheckbox :model-value="isPageAllSelected" @update:model-value="toggleSelectAll" />
      </div>
      <div class="col-index flex items-center justify-center overflow-hidden">#</div>
      <div class="col-thumb flex items-center justify-center overflow-hidden">
        {{ $t('admin.tagging.fileList.headers.thumbnail') }}
      </div>
      <div class="col-filename flex items-center overflow-hidden truncate">
        {{ $t('admin.tagging.fileList.headers.filename') }}
      </div>
      <div class="col-id flex items-center overflow-hidden truncate">{{ $t('admin.tagging.fileList.headers.id') }}</div>
      <div class="col-dimension flex items-center overflow-hidden truncate">
        {{ $t('admin.tagging.fileList.headers.dimension') }}
      </div>
      <div class="col-size flex items-center overflow-hidden truncate">{{ $t('admin.tagging.fileList.headers.size') }}</div>
      <div class="col-format flex items-center overflow-hidden truncate">{{ $t('admin.tagging.fileList.headers.format') }}</div>
      <div class="col-uploader flex items-center overflow-hidden truncate">
        {{ $t('admin.tagging.fileList.headers.uploader') }}
      </div>
      <div class="col-status flex items-center overflow-hidden truncate">{{ $t('admin.tagging.fileList.headers.status') }}</div>
      <div class="col-tries flex items-center overflow-hidden truncate">{{ $t('admin.tagging.fileList.headers.tries') }}</div>
      <div class="col-duration flex items-center overflow-hidden truncate">{{ $t('admin.tagging.fileList.headers.duration') }}</div>
      <div class="col-time flex items-center overflow-hidden truncate">{{ $t('admin.tagging.fileList.headers.time') }}</div>
      <div class="col-actions flex items-center justify-center overflow-hidden">{{ $t('admin.tagging.fileList.headers.actions') }}</div>
    </div>

    <CyberSkeleton type="table" :count="6" :loading="isLoading" />

    <div
      v-if="!isLoading && images.length === 0"
      class="rounded border border-subtle p-3 py-6 text-center"
      style="background: rgba(var(--color-background-800-rgb), 0.3); color: var(--color-content-muted)"
    >
      <i class="fas fa-search mb-2 text-2xl text-brand-500" />
      <p class="text-xs">{{ $t('admin.tagging.fileList.noFiles') }}</p>
    </div>

    <div v-if="!isLoading && images.length > 0" class="space-y-1">
      <div
        v-for="(image, index) in images"
        :key="image.id"
        class="image-list-item grid-cols-table grid gap-4 rounded border border-subtle px-4 py-2 transition-all duration-300"
        style="background: rgba(var(--color-background-800-rgb), 0.3)"
        @mouseenter="$event.currentTarget.style.background = 'rgba(var(--color-background-800-rgb), 0.5)'"
        @mouseleave="$event.currentTarget.style.background = 'rgba(var(--color-background-800-rgb), 0.3)'"
      >
        <div class="col-select flex items-center justify-center">
          <CyberCheckbox
            :model-value="selectedIds.includes(image.id)"
            @update:model-value="() => $emit('toggle-select', image.id)"
          />
        </div>

        <div class="col-index text-2xs text-content-content-muted flex items-center justify-center">
          {{ (currentPage - 1) * pageSize + index + 1 }}
        </div>

        <div class="col-thumb flex items-center justify-center">
          <div class="thumb-container overflow-hidden rounded border border-brand-300" v-loading="imageLoadingStates[image.id]">
            <CyberFile
              :src="image.full_thumb_url"
              :alt="image.original_name"
              class="h-auto max-h-10 w-full"
              fit-mode="cover"
              :retry-count="2"
              :is-nsfw="image.is_nsfw"
              @loading="imageLoadingStates[image.id] = $event"
            />
          </div>
        </div>

        <div class="col-filename flex items-center">
          <div class="text-2xs w-full truncate text-content-heading" :title="image.original_name">
            {{ image.original_name }}
          </div>
        </div>

        <div class="col-id flex items-center">
          <div class="text-2xs w-full truncate text-content" :title="image.id">
            {{ image.id }}
          </div>
        </div>

        <div class="col-dimension flex items-center">
          <div class="text-2xs w-full truncate text-content-heading" :title="`${image.width || 0}x${image.height || 0}`">
            {{ image.width || 0 }}x{{ image.height || 0 }}
          </div>
        </div>

        <div class="col-size flex items-center">
          <div class="text-2xs w-full truncate text-content" :title="formatFileSize(image.size)">
            {{ formatFileSize(image.size) }}
          </div>
        </div>

        <div class="col-format flex items-center">
          <div class="text-2xs w-full truncate text-content" :title="image.format || $t('admin.tagging.fileList.unknown')">
            {{ image.format || $t('admin.tagging.fileList.unknown') }}
          </div>
        </div>

        <div class="col-uploader flex items-center">
          <div class="text-2xs w-full truncate text-content" :title="`ID ${image.user_id}`">ID {{ image.user_id }}</div>
        </div>

        <div class="col-status flex items-center">
          <span
            class="status-badge text-2xs whitespace-nowrap rounded-full px-1.5 py-0.5"
            :class="getStatusClass(image.ai_tagging_status)"
          >
            {{ getStatusText(image.ai_tagging_status) }}
          </span>
        </div>

        <div class="col-tries flex items-center">
          <div class="text-2xs w-full truncate text-content" :title="`${image.ai_tagging_tries}`">
            {{ image.ai_tagging_tries }}
          </div>
        </div>

        <div class="col-duration flex items-center">
          <div class="text-2xs w-full truncate text-content" :title="image.ai_http_duration_formatted || '-'">
            <template v-if="image.ai_http_duration_formatted">
              <span class="text-brand-400">HTTP:</span> {{ image.ai_http_duration_formatted }}
            </template>
            <template v-else>-</template>
          </div>
        </div>

        <div class="col-time flex items-center">
          <div class="text-2xs w-full truncate text-content" :title="timeString(image.updated_at)">
            {{ formatTime(image.updated_at) }}
          </div>
        </div>

        <div class="col-actions flex items-center justify-center">
          <button
            v-if="image.ai_tagging_status === 'failed'"
            class="retry-button"
            :title="$t('admin.tagging.fileList.retryTooltip')"
            @click="$emit('retry', [image.id])"
          >
            <i class="fas fa-redo"></i>
          </button>
        </div>
      </div>
    </div>

    <div v-if="!isLoading && Math.ceil(total / pageSize) > 1" class="mt-4 flex justify-center">
      <CyberPagination
        :current-page="currentPage"
        :total-pages="Math.ceil(total / pageSize)"
        :total="total"
        :page-size="pageSize"
        :show-page-size-selector="true"
        :show-quick-jumper="true"
        class="text-xs"
        @update:current-page="$emit('page-change', $event)"
        @update:page-size="$emit('page-size-change', $event)"
      />
    </div>
  </div>
</template>

<style scoped lang="scss">
  .module-title {
    font-size: var(--text-lg);
    font-weight: var(--font-semibold);
    color: var(--color-brand-500);
    margin: 0;
    display: flex;
    align-items: center;
    text-shadow: 0 0 8px rgba(var(--color-brand-500-rgb), 0.3);
    letter-spacing: 0.5px;

    i {
      filter: drop-shadow(0 0 4px rgba(var(--color-brand-500-rgb), 0.5));
    }
  }

  .tagging-image-list {
    overflow-x: auto;
    --col-select-width: 32px;
    --col-index-width: 40px;
    --col-thumb-width: 70px;
    --col-filename-width: 120px;
    --col-id-width: 100px;
    --col-dimension-width: 80px;
    --col-size-width: 60px;
    --col-format-width: 50px;
    --col-uploader-width: 70px;
    --col-status-width: 70px;
    --col-tries-width: 70px;
    --col-duration-width: 90px;
    --col-time-width: 80px;
    --col-actions-width: 40px;
    width: 100%;
  }

  .image-list-item {
    box-shadow: 0 2px 6px rgba(0, 0, 0, 0.2);
    backdrop-filter: blur(10px);
  }

  .image-list-item:hover {
    box-shadow: 0 3px 8px rgba(0, 0, 0, 0.3);
  }

  .status-badge {
    text-shadow: 0 0 3px rgba(0, 0, 0, 0.5);
    box-shadow: 0 1px 2px rgba(0, 0, 0, 0.3);
  }

  .thumb-container {
    width: 100%;
    max-width: 60px;
    height: 50px;
    display: flex;
    align-items: center;
    justify-content: center;
    background-color: var(--color-background-800);
    margin: 0 auto;
    position: relative;
  }

  .text-2xs {
    font-size: 0.72rem;
    line-height: 1rem;
  }

  .grid-cols-table {
    display: grid;
    grid-template-columns:
      var(--col-select-width)
      var(--col-index-width)
      var(--col-thumb-width)
      var(--col-filename-width)
      var(--col-id-width)
      var(--col-dimension-width)
      var(--col-size-width)
      var(--col-format-width)
      var(--col-uploader-width)
      var(--col-status-width)
      var(--col-tries-width)
      var(--col-duration-width)
      var(--col-time-width)
      var(--col-actions-width);
  }

  .col-select {
    width: var(--col-select-width);
  }
  .col-index {
    width: var(--col-index-width);
  }
  .col-thumb {
    width: var(--col-thumb-width);
  }
  .col-filename {
    width: var(--col-filename-width);
  }
  .col-id {
    width: var(--col-id-width);
  }
  .col-dimension {
    width: var(--col-dimension-width);
  }
  .col-size {
    width: var(--col-size-width);
  }
  .col-format {
    width: var(--col-format-width);
  }
  .col-uploader {
    width: var(--col-uploader-width);
  }
  .col-status {
    width: var(--col-status-width);
  }
  .col-tries {
    width: var(--col-tries-width);
  }
  .col-duration {
    width: var(--col-duration-width);
  }
  .col-time {
    width: var(--col-time-width);
  }
  .col-actions {
    width: var(--col-actions-width);
  }

  .retry-button {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 28px;
    height: 28px;
    padding: 0;
    border: 1px solid var(--color-brand-500);
    background: rgba(var(--color-brand-500-rgb), 0.1);
    color: var(--color-brand-500);
    border-radius: 4px;
    cursor: pointer;
    transition: all 0.2s ease;
    font-size: 0.7rem;
  }

  .retry-button:hover {
    background: rgba(var(--color-brand-500-rgb), 0.2);
    border-color: var(--color-brand-400);
    color: var(--color-brand-400);
    box-shadow: 0 0 8px rgba(var(--color-brand-500-rgb), 0.3);
  }

  .retry-button:active {
    transform: scale(0.95);
  }

  @keyframes wave {
    0% {
      background-position: -200% 0;
    }
    100% {
      background-position: 200% 0;
    }
  }
</style>
