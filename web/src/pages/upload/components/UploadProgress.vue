<script setup lang="ts">
  import { nextTick, onMounted, ref, watch } from 'vue'
  import { useToast } from '@/components/Toast/useToast'
  import { useTexts } from '@/composables/useTexts'
  import { RecycleScroller } from 'vue-virtual-scroller'
  import { SIZE_LIMITS, TIMING } from '@/constants'
  import { createFileDataURL } from '@/utils/file/heicConverter'
  import CyberTooltip from '@/components/Tooltip/index.vue'
  import 'vue-virtual-scroller/dist/vue-virtual-scroller.css'

  export interface UploadFile {
    file: File
    name: string
    size: number
    progress: number
    status: 'pending' | 'uploading' | 'success' | 'error' | 'analyzing' | 'retrying' | 'paused' | 'completed' | 'instant'
    url?: string
    full_url?: string
    preview?: string
    error?: string
    id?: string
    uniqueId?: string // 用于虚拟列表的唯一标识
    lastStatusChange?: number // 上次状态变化的时间戳
    previewLoaded?: boolean // 标记预览是否已加载
    itemLoaded?: boolean // 标记项目是否已加载完成
    statusMessage?: string // 状态消息
    speed?: number // 上传速度
    remainingTime?: number // 剩余时间
    is_duplicate?: boolean // 是否是秒传
    type?: 'regular' | 'chunked' // 上传类型
  }

  const props = defineProps<{
    files: UploadFile[]
  }>()

  const emit = defineEmits<{
    (e: 'remove', index: number): void
    (e: 'upload'): void
    (e: 'cancel'): void
    (e: 'retry', index: number): void
    (e: 'resume', index: number): void
    (e: 'update:files', files: UploadFile[]): void
  }>()

  const toast = useToast()
  const { $t } = useTexts()
  const listKey = ref<number>(0) // 添加一个key，在数组变化时使虚拟列表重新渲染
  const recycleScroller = ref<any>(null)
  const itemRefs = ref<Map<string, HTMLElement>>(new Map())
  const _lastScrollTime = ref<number>(0)
  const _scrollDebounceTime = TIMING.SCROLL.UPLOAD_DEBOUNCE_TIME // 滚动防抖时间
  const scrollQueue = ref<string[]>([])
  const processingScroll = ref<boolean>(false)
  const previewQueue = ref<UploadFile[]>([]) // 预览加载队列
  const isPreloading = ref<boolean>(false) // 是否正在预加载
  const preloadBatchSize = SIZE_LIMITS.COMPONENT.PRELOAD_BATCH_SIZE // 每批预加载的数量
  const preloadDelay = TIMING.ANIMATION.PRELOAD_DELAY // 预加载延迟
  const loadingDelay = TIMING.ANIMATION.LOADING_DELAY // 显示骨架屏的最小时间

  const setItemRef = (el: any, uniqueId: string | undefined) => {
    if (el && uniqueId) {
      itemRefs.value.set(uniqueId, el)
    }
  }

  const isItemLoaded = (file: UploadFile): boolean => {
    if (file.itemLoaded !== undefined) {
      return file.itemLoaded
    }

    if (file.status !== 'pending' || file.progress > 0) {
      return true
    }

    return false
  }

  const onImageLoaded = (file: UploadFile) => {
    file.previewLoaded = true

    setTimeout(() => {
      file.itemLoaded = true
    }, loadingDelay)
  }

  const processScrollQueue = async () => {
    if (processingScroll.value || scrollQueue.value.length === 0) {
      return
    }

    processingScroll.value = true

    const uniqueId = scrollQueue.value.shift()
    if (!uniqueId) {
      processingScroll.value = false
      return
    }

    const fileIndex = props.files.findIndex((f) => f.uniqueId === uniqueId)
    if (fileIndex === -1 || !props.files[fileIndex]) {
      processingScroll.value = false
      return
    }

    if (recycleScroller.value && typeof recycleScroller.value.scrollToItem === 'function') {
      await nextTick()
      recycleScroller.value.scrollToItem(fileIndex)
    }

    processingScroll.value = false

    setTimeout(() => {
      if (scrollQueue.value.length > 0) {
        processScrollQueue()
      }
    }, 200)
  }

  watch(
    () => [...props.files],
    (newFiles, oldFiles) => {
      newFiles.forEach((file) => {
        if (!file.uniqueId) {
          return
        }

        const oldFile = oldFiles.find((f) => f.uniqueId === file.uniqueId)

        if (oldFile && oldFile.status === 'uploading' && (file.status === 'success' || file.status === 'error')) {
          file.lastStatusChange = Date.now()

          scrollQueue.value.push(file.uniqueId)

          if (!processingScroll.value) {
            processScrollQueue()
          }
        }

        if (!file.previewLoaded && !previewQueue.value.includes(file)) {
          previewQueue.value.push(file)
          if (!isPreloading.value) {
            preloadPreviews()
          }
        }
      })
    },
    { deep: true }
  )

  const preloadPreviews = async () => {
    if (isPreloading.value || previewQueue.value.length === 0) {
      return
    }

    isPreloading.value = true

    const batch = previewQueue.value.splice(0, preloadBatchSize)

    await Promise.all(
      batch.map(async (file) => {
        if (!file || file.preview || file.previewLoaded) {
          return
        }

        if (file.file instanceof File) {
          try {
            const preview = await createFilePreview(file.file)

            const fileIndex = props.files.findIndex((f) => f.uniqueId === file.uniqueId)
            if (fileIndex !== -1 && props.files[fileIndex]) {
              const updatedFiles = [...props.files]
              updatedFiles[fileIndex] = { ...updatedFiles[fileIndex], preview }
              emit('update:files', updatedFiles)
            }
          } catch (error) {
            console.error(`预览图加载失败: ${file.name}`, error)

            if (file.uniqueId) {
              const fileIndex = props.files.findIndex((f) => f.uniqueId === file.uniqueId)
              if (fileIndex !== -1 && props.files[fileIndex]) {
                setTimeout(() => {
                  if (props.files[fileIndex]) {
                    const updatedFiles = [...props.files]
                    updatedFiles[fileIndex] = { ...updatedFiles[fileIndex], itemLoaded: true }
                    emit('update:files', updatedFiles)
                  }
                }, loadingDelay)
              }
            }
          }
        }
      })
    )

    await new Promise((resolve) => setTimeout(resolve, preloadDelay))

    isPreloading.value = false

    if (previewQueue.value.length > 0) {
      preloadPreviews()
    }
  }

  const createFilePreview = async (file: File): Promise<string> => {
    try {
      return await createFileDataURL(file)
    } catch (error) {
      console.error('预览图生成失败:', file.name, error)
      return 'data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iMTAwIiBoZWlnaHQ9IjEwMCIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj48cmVjdCB3aWR0aD0iMTAwIiBoZWlnaHQ9IjEwMCIgZmlsbD0iIzMzMyIvPjx0ZXh0IHg9IjUwJSIgeT0iNTAlIiBkb21pbmFudC1iYXNlbGluZT0ibWlkZGxlIiB0ZXh0LWFuY2hvcj0ibWlkZGxlIiBmaWxsPSIjZmZmIiBmb250LXNpemU9IjE2Ij7mianlh4Y8L3RleHQ+PC9zdmc+'
    }
  }

  watch(
    () => props.files.length,
    () => {
      listKey.value++
    }
  )

  const ensureUniqueIds = () => {
    props.files.forEach((file, _index) => {
      if (!file.uniqueId) {
        const timestamp = new Date().getTime()
        file.uniqueId = `file-${index}-${timestamp}-${file.size}`

        file.itemLoaded = false
        file.previewLoaded = false
      }
    })
  }

  onMounted(() => {
    ensureUniqueIds()

    const filesToPreload = props.files.filter((file) => !file.preview && !file.previewLoaded)

    if (filesToPreload.length > 0) {
      previewQueue.value.push(...filesToPreload)
      preloadPreviews()
    }

    setTimeout(() => {
      props.files.forEach((file, _index) => {
        if (!file.itemLoaded) {
          setTimeout(
            () => {
              file.itemLoaded = true
            },
            300 + Math.random() * 200
          ) // 300ms到500ms的随机延迟
        }
      })
    }, 100)
  })

  watch(
    () => props.files,
    () => {
      ensureUniqueIds()
    },
    { deep: true }
  )

  const formatFileSize = (bytes: number | undefined): string => {
    if (bytes === undefined || bytes === null || isNaN(bytes)) {
      return '0 B'
    }
    if (bytes < 1024) {
      return `${bytes} B`
    }
    if (bytes < 1024 * 1024) {
      return `${(bytes / 1024).toFixed(1)} KB`
    }
    return `${(bytes / (1024 * 1024)).toFixed(1)} MB`
  }

  const removeFile = (index: number) => {
    emit('remove', index)
    listKey.value++
  }

  const retryUpload = (index: number) => {
    emit('retry', index)
  }

  const resumeUpload = (index: number) => {
    emit('resume', index)
  }

  const openInNewTab = (url: string) => {
    if (!url) {
      return
    }
    window.open(url, '_blank', 'noopener,noreferrer')
  }

  const copyImageUrl = async (url: string) => {
    try {
      await navigator.clipboard.writeText(url)
      toast.success($t('upload.uploadProgress.toast.linkCopied'))
    } catch {
      toast.error($t('upload.uploadProgress.toast.copyFailed'))
    }
  }

  const copyMarkdownUrl = async (url: string, name: string) => {
    try {
      const markdownText = `![${name}](${url})`
      await navigator.clipboard.writeText(markdownText)
      toast.success($t('upload.uploadProgress.toast.markdownLinkCopied'))
    } catch {
      toast.error($t('upload.uploadProgress.toast.copyFailed'))
    }
  }

  const copyHtmlUrl = async (url: string, name: string) => {
    try {
      const htmlTag = `<img src="${url}" alt="${name}" />`
      await navigator.clipboard.writeText(htmlTag)
      toast.success($t('upload.uploadProgress.toast.htmlLinkCopied'))
    } catch {
      toast.error($t('upload.uploadProgress.toast.copyFailed'))
    }
  }

  const copyThumbnailUrl = async (url: string) => {
    try {
      await navigator.clipboard.writeText(url)
      toast.success($t('upload.uploadProgress.toast.thumbnailLinkCopied'))
    } catch {
      toast.error($t('upload.uploadProgress.toast.copyFailed'))
    }
  }

  const formatRemainingTime = (seconds: number): string => {
    if (seconds < 60) {
      return $t('upload.uploadProgress.time.seconds', { count: Math.ceil(seconds) })
    } else if (seconds < 3600) {
      const minutes = Math.ceil(seconds / 60)
      return $t('upload.uploadProgress.time.minutes', { count: minutes })
    }
    const hours = Math.floor(seconds / 3600)
    const minutes = Math.ceil((seconds % 3600) / 60)
    return hours > 0
      ? $t('upload.uploadProgress.time.hoursMinutes', { hours, minutes })
      : $t('upload.uploadProgress.time.minutes', { count: minutes })
  }

  const isInstantUpload = (file: UploadFile): boolean => {
    if (file.is_duplicate) {
      return true
    }

    if (file.statusMessage) {
      // 获取翻译后的秒传相关文本用于比较
      const instantCompleted = $t('upload.smartUpload.status.instantCompleted')
      const instantSuccess = $t('upload.uploadProgress.status.instantSuccess')
      const checkingInstant = $t('upload.smartUpload.status.checkingInstant')

      // 检查翻译后的文本
      if (
        file.statusMessage.includes(instantCompleted) ||
        file.statusMessage.includes(instantSuccess) ||
        file.statusMessage.includes(checkingInstant) ||
        file.statusMessage.includes('instant')
      ) {
        return true
      }
    }

    return false
  }

  const getSuccessText = (file: UploadFile): string =>
    isInstantUpload(file) ? $t('upload.uploadProgress.status.instantSuccess') : $t('upload.uploadProgress.status.uploadSuccess')

  const getSuccessIconClass = (file: UploadFile): string => {
    if (isInstantUpload(file)) {
      return 'fas fa-bolt text-cyan-400'
    }
    return 'fas fa-check-circle text-green-400'
  }
</script>

<template>
  <div class="upload-progress-container">
    <RecycleScroller
      ref="recycleScroller"
      :key="listKey"
      v-slot="{ item: file, index }"
      class="scroller"
      :items="files"
      :item-size="82"
      :buffer="500"
      :prerender="10"
      key-field="uniqueId"
    >
      <div
        :id="`upload-item-${file.uniqueId}`"
        :ref="(el) => setItemRef(el, file.uniqueId)"
        class="upload-item"
        :data-status="file.status"
      >
        <div class="upload-card">
          <div class="upload-preview">
            <div v-if="!file.preview || !file.previewLoaded" class="skeleton absolute inset-0" />
            <img
              v-if="file.preview"
              :src="file.preview"
              class="h-full w-full object-contain"
              loading="lazy"
              :class="{
                'opacity-0': !file.previewLoaded,
                'opacity-100 transition-opacity duration-300': file.previewLoaded,
              }"
              @load="onImageLoaded(file)"
            />
            <i v-else class="fas fa-file-image text-muted z-10 text-lg" />
          </div>

          <div class="upload-body">
            <div class="flex items-center justify-between">
              <div class="relative w-2/3">
                <div v-if="!isItemLoaded(file)" class="skeleton h-4 w-full rounded" />
                <p v-else class="truncate text-sm font-medium text-content">{{ file.name }}</p>
              </div>
              <div class="relative">
                <div v-if="!isItemLoaded(file)" class="skeleton h-3 w-10 rounded" />
                <span v-else class="text-muted text-xs">{{ formatFileSize(file.size) }}</span>
              </div>
            </div>

            <div v-if="!isItemLoaded(file)" class="relative mt-1">
              <div class="skeleton h-3 w-2/3 rounded" />
            </div>
            <div v-else-if="file.status === 'error'" class="status-badge status-badge--error">
              <i class="fas fa-exclamation-circle mr-1.5 text-error-400" />
              <span class="font-medium">{{ file.error || $t('upload.uploadProgress.status.uploadFailed') }}</span>
            </div>
            <div
              v-else-if="file.status === 'success' || file.status === 'completed'"
              :class="['status-badge', isInstantUpload(file) ? 'status-badge--instant' : 'status-badge--success']"
            >
              <i :class="getSuccessIconClass(file)" class="mr-1.5" />
              <span class="font-medium">{{ getSuccessText(file) }}</span>
            </div>
            <div v-else-if="file.status === 'uploading'" class="mt-0.5">
              <div class="mb-0.5 flex items-center">
                <div class="relative h-1.5 w-full overflow-hidden rounded-full bg-background-800 shadow-inner">
                  <div
                    class="relative h-full rounded-full bg-gradient-to-r from-brand-500 via-info-400 to-error-500 transition-all duration-300 ease-out"
                    :style="{ width: `${file.progress}%` }"
                  >
                    <div
                      class="absolute inset-0 animate-pulse bg-gradient-to-r from-transparent via-white to-transparent opacity-20"
                    />
                  </div>
                </div>
                <span class="ml-2 min-w-10 text-right text-xs font-medium text-content">{{ file.progress }}%</span>
              </div>
              <div v-if="file.statusMessage || file.remainingTime" class="text-muted flex items-center justify-between text-xs">
                <span v-if="file.statusMessage" class="truncate">{{ file.statusMessage }}</span>
                <span v-if="file.remainingTime && file.remainingTime > 0" class="text-2xs text-muted ml-2 flex-shrink-0">
                  {{ $t('upload.uploadProgress.remaining') }} {{ formatRemainingTime(file.remainingTime) }}
                </span>
              </div>
            </div>
            <div v-else-if="file.status === 'analyzing'" class="status-badge status-badge--analyzing">
              <i class="fas fa-spinner fa-spin mr-1.5 text-content" />
              <span class="font-medium">{{ $t('upload.uploadProgress.status.analyzing') }}</span>
            </div>
            <div v-else-if="file.status === 'paused'" class="status-badge status-badge--paused">
              <i class="fas fa-pause-circle mr-1.5 text-warning-400" />
              <span class="font-medium">{{ file.statusMessage || $t('upload.uploadProgress.status.paused') }}</span>
            </div>
            <div v-else-if="file.status === 'retrying'" class="status-badge status-badge--retry">
              <i class="fas fa-redo-alt fa-spin mr-1.5 text-warning-400" />
              <span class="font-medium">{{ $t('upload.uploadProgress.status.retrying') }}</span>
            </div>
            <div v-else class="status-badge status-badge--pending">
              <i class="fas fa-clock text-muted mr-1.5" />
              <span class="font-medium">{{ file.statusMessage || $t('upload.uploadProgress.status.pending') }}</span>
            </div>
          </div>

          <div class="upload-actions">
            <div v-if="!isItemLoaded(file)" class="flex space-x-1">
              <div class="skeleton h-6 w-6 rounded" />
              <div class="skeleton h-6 w-6 rounded" />
              <div class="skeleton h-6 w-6 rounded" />
            </div>
            <template v-else>
              <button
                v-if="(file.status === 'success' || file.status === 'completed') && file.full_url"
                class="upload-action upload-action--success"
                :title="$t('upload.uploadProgress.actions.openInNewTab')"
                @click="openInNewTab(file.full_url)"
              >
                <i class="fas fa-external-link-alt" />
              </button>

              <CyberTooltip
                v-if="(file.status === 'success' || file.status === 'completed') && file.full_url"
                :content="$t('upload.uploadProgress.actions.copyLinkTooltip')"
                placement="top"
              >
                <button class="upload-action upload-action--brand" @click="copyImageUrl(file.full_url)">
                  <i class="fas fa-link" />
                </button>
              </CyberTooltip>

              <CyberTooltip
                v-if="(file.status === 'success' || file.status === 'completed') && file.full_url"
                :content="$t('upload.uploadProgress.actions.copyMarkdownTooltip')"
                placement="top"
              >
                <button class="upload-action upload-action--brand" @click="copyMarkdownUrl(file.full_url, file.name)">
                  <i class="fab fa-markdown" />
                </button>
              </CyberTooltip>

              <CyberTooltip
                v-if="(file.status === 'success' || file.status === 'completed') && file.full_url"
                :content="$t('upload.uploadProgress.actions.copyHtmlTooltip')"
                placement="top"
              >
                <button class="upload-action upload-action--warning" @click="copyHtmlUrl(file.full_url, file.name)">
                  <i class="fab fa-html5" />
                </button>
              </CyberTooltip>

              <CyberTooltip
                v-if="(file.status === 'success' || file.status === 'completed') && file.full_thumb_url"
                :content="$t('upload.uploadProgress.actions.copyThumbnailTooltip')"
                placement="top"
              >
                <button class="upload-action upload-action--info" @click="copyThumbnailUrl(file.full_thumb_url)">
                  <i class="fas fa-image" />
                </button>
              </CyberTooltip>

              <button
                v-if="file.status === 'paused'"
                class="upload-action upload-action--success"
                :title="$t('upload.uploadProgress.actions.resume')"
                @click="resumeUpload(index)"
              >
                <i class="fas fa-play" />
              </button>

              <button
                v-if="file.status === 'error'"
                class="upload-action upload-action--warning"
                :title="$t('upload.uploadProgress.actions.retry')"
                @click="retryUpload(index)"
              >
                <i class="fas fa-redo-alt" />
              </button>

              <button
                class="upload-action upload-action--danger"
                :disabled="file.status === 'uploading'"
                :class="{
                  'cursor-not-allowed opacity-50 hover:scale-100 hover:bg-transparent hover:text-error-500':
                    file.status === 'uploading',
                }"
                :title="$t('upload.uploadProgress.actions.remove')"
                @click="removeFile(index)"
              >
                <i class="fas fa-times" />
              </button>
            </template>
          </div>
        </div>
      </div>
    </RecycleScroller>
  </div>
</template>

<style scoped>
  .upload-item {
    transition: all 0.3s ease;
  }

  .upload-card {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    border-radius: var(--radius-sm);
    border: 1px solid var(--color-border-default);
    background-color: rgba(var(--color-background-700-rgb), 0.85);
    padding: 0.625rem 0.75rem;
    margin-bottom: 1rem;
    transition:
      border-color 0.3s ease,
      background-color 0.3s ease,
      box-shadow 0.3s ease;
  }

  .upload-card:hover {
    border-color: rgba(var(--color-brand-500-rgb), 0.6);
    background-color: rgba(var(--color-background-700-rgb), 0.95);
    box-shadow:
      0 0 0 1px rgba(var(--color-brand-500-rgb), 0.3),
      0 0 20px rgba(var(--color-brand-500-rgb), 0.15),
      0 4px 12px rgba(0, 0, 0, 0.16);
  }

  .upload-preview {
    position: relative;
    flex: 0 0 auto;
    height: 2.5rem;
    width: 2.5rem;
    display: flex;
    align-items: center;
    justify-content: center;
    overflow: hidden;
    border-radius: var(--radius-sm);
    border: 1px solid var(--color-border-default);
    background-color: rgba(var(--color-background-800-rgb), 0.95);
  }

  .upload-body {
    min-width: 0;
    flex: 1;
    padding-right: 0.5rem;
  }

  .upload-actions {
    display: flex;
    align-items: center;
    gap: 0.375rem;
  }

  .status-badge {
    margin-top: 0.25rem;
    display: inline-flex;
    align-items: center;
    gap: 0.35rem;
    border-radius: var(--radius-sm);
    border: 1px solid transparent;
    padding: 0.25rem 0.5rem;
    font-size: 0.6875rem;
    font-weight: 600;
  }

  .status-badge--error {
    background-color: rgba(var(--color-error-rgb), 0.16);
    border-color: rgba(var(--color-error-rgb), 0.32);
    color: var(--color-error-400);
  }

  .status-badge--success {
    background-color: rgba(var(--color-success-rgb), 0.16);
    border-color: rgba(var(--color-success-rgb), 0.32);
    color: var(--color-success-400);
  }

  .status-badge--instant {
    background-color: rgba(var(--color-brand-500-rgb), 0.18);
    border-color: rgba(var(--color-brand-500-rgb), 0.36);
    color: var(--color-brand-400);
  }

  .status-badge--analyzing {
    background-color: rgba(var(--color-brand-500-rgb), 0.14);
    border-color: rgba(var(--color-brand-500-rgb), 0.3);
    color: var(--color-brand-400);
  }

  .status-badge--paused,
  .status-badge--retry {
    background-color: rgba(var(--color-warning-rgb), 0.16);
    border-color: rgba(var(--color-warning-rgb), 0.3);
    color: var(--color-warning-400);
  }

  .status-badge--pending {
    background-color: rgba(var(--color-background-700-rgb), 0.6);
    border-color: var(--color-border-subtle);
    color: var(--color-content-muted);
  }

  .upload-action {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    border-radius: var(--radius-sm);
    height: 2rem;
    width: 2rem;
    border: 1px solid var(--color-border-default);
    background-color: rgba(var(--color-background-700-rgb), 0.7);
    color: var(--color-content-muted);
    font-size: 0.875rem;
    transition:
      transform 0.2s ease,
      border-color 0.2s ease,
      background-color 0.2s ease,
      color 0.2s ease,
      box-shadow 0.2s ease;
  }

  .upload-action:hover {
    box-shadow:
      0 0 0 2px currentColor,
      0 2px 8px rgba(0, 0, 0, 0.15);
  }

  .upload-action--brand {
    color: var(--color-brand-400);
    border-color: rgba(var(--color-brand-500-rgb), 0.36);
    background-color: rgba(var(--color-brand-500-rgb), 0.12);
  }

  .upload-action--brand:hover {
    background-color: rgba(var(--color-brand-500-rgb), 0.18);
  }

  .upload-action--success {
    color: var(--color-success-400);
    border-color: rgba(var(--color-success-rgb), 0.32);
    background-color: rgba(var(--color-success-rgb), 0.12);
  }

  .upload-action--success:hover {
    background-color: rgba(var(--color-success-rgb), 0.18);
  }

  .upload-action--warning {
    color: var(--color-warning-400);
    border-color: rgba(var(--color-warning-rgb), 0.28);
    background-color: rgba(var(--color-warning-rgb), 0.12);
  }

  .upload-action--warning:hover {
    background-color: rgba(var(--color-warning-rgb), 0.18);
  }

  .upload-action--danger {
    color: var(--color-error-400);
    border-color: rgba(var(--color-error-rgb), 0.28);
    background-color: rgba(var(--color-error-rgb), 0.12);
  }

  .upload-action--danger:hover {
    background-color: rgba(var(--color-error-rgb), 0.18);
  }

  .upload-action[disabled] {
    opacity: 0.5;
    cursor: not-allowed;
  }

  .scroller {
    height: var(--upload-scroll-height);
    overflow-y: auto;
  }

  .upload-item.highlight {
    animation: highlight-animation var(--animation-slow) ease;
  }

  @keyframes highlight-animation {
    0% {
      background-color: rgba(var(--color-brand-500-rgb), 0.2);
    }
    100% {
      background-color: transparent;
    }
  }
</style>
