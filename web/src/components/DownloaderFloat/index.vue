<script setup lang="ts">
  import { computed, onMounted, onUnmounted, ref } from 'vue'
  import { useToast } from '@/components/Toast/useToast'
  import { downloadFile, downloadSharedFile } from '@/api/file'
  import { onQueueChange } from '@/utils/file/downloader'
  import { StorageUtil } from '@/utils/storage/storage'
  import { useTexts } from '@/composables/useTexts'
  import type { DownloadOptions } from './types'

  const props = withDefaults(
    defineProps<{
      options: DownloadOptions
    }>(),
    {
      options: () => ({
        fileId: '',
        position: { side: 'right', y: '50%' },
        autoClose: true,
        closeDelay: 3000,
        queueSize: 0,
      }),
    }
  )

  const emit = defineEmits<{
    complete: [success: boolean]
    close: []
  }>()

  const toast = useToast()
  const { $t } = useTexts()

  /* 状态管理 */
  const visible = ref(true)
  const status = ref<'preparing' | 'downloading' | 'success' | 'error'>('preparing')
  const progress = ref(0)
  const loaded = ref(0)
  const total = ref(0)
  const fileName = ref(props.options.fileName || '')
  const showFileInfo = ref(false)
  const downloadStartTime = ref(0)

  /* 拖动相关状态 */
  const isDragging = ref(false)
  const dragOffset = ref({ x: 0, y: 0 })
  const currentPosition = ref({ x: 0, y: 0 })
  const dragStartTime = ref(0)

  /* 队列状态 */
  const queueSize = ref(props.options.queueSize || 0)

  /* 下载控制 */
  let downloadController: AbortController | null = null
  let autoCloseTimer: number | null = null

  const getSavedPosition = () => StorageUtil.get<{ x: number; y: number }>('cyber-downloader-position')

  const savePosition = (x: number, y: number) => {
    StorageUtil.set('cyber-downloader-position', { x, y })
  }

  const initializePosition = () => {
    const savedPos = getSavedPosition()
    const pos = props.options.position || {}
    const { innerWidth: windowWidth, innerHeight: windowHeight } = window
    const elementSize = 120

    let x, y

    if (savedPos && savedPos.x !== 0 && savedPos.y !== 0) {
      const isInViewport =
        savedPos.x >= 0 && savedPos.y >= 0 && savedPos.x + elementSize <= windowWidth && savedPos.y + elementSize <= windowHeight

      if (isInViewport) {
        x = savedPos.x
        y = savedPos.y
      } else {
        x = windowWidth - elementSize - 30
        y = (windowHeight - elementSize) / 2
      }
    } else {
      x = pos.side === 'left' ? 30 : pos.side === 'center' ? (windowWidth - elementSize) / 2 : windowWidth - elementSize - 30
      y = (windowHeight - elementSize) / 2

      if (pos.y) {
        y =
          typeof pos.y === 'string' && pos.y.includes('%')
            ? (windowHeight - elementSize) * (parseInt(pos.y) / 100)
            : typeof pos.y === 'string'
              ? parseInt(pos.y)
              : pos.y
      }
      if (pos.x) {
        x = typeof pos.x === 'string' ? parseInt(pos.x) : pos.x
      }
    }

    currentPosition.value.x = Math.max(0, Math.min(x, windowWidth - elementSize))
    currentPosition.value.y = Math.max(0, Math.min(y, windowHeight - elementSize))
  }

  const floatStyle = computed(() => ({
    position: 'fixed',
    zIndex: 10000,
    pointerEvents: 'auto',
    userSelect: 'none',
    cursor: isDragging.value ? 'grabbing' : 'grab',
    left: `${currentPosition.value.x}px`,
    top: `${currentPosition.value.y}px`,
    transform: 'none',
  }))

  const progressStyle = computed(() => {
    const circumference = 2 * Math.PI * 36
    return {
      strokeDasharray: circumference,
      strokeDashoffset: circumference - (progress.value / 100) * circumference,
      transition: 'stroke-dashoffset 0.3s ease',
      filter: 'url(#glow)',
    }
  })

  const formatBytes = (bytes: number): string => {
    if (bytes === 0) {
      return '0 B'
    }
    const k = 1024
    const sizes = ['B', 'KB', 'MB', 'GB']
    const i = Math.floor(Math.log(bytes) / Math.log(k))
    return `${parseFloat((bytes / Math.pow(k, i)).toFixed(1))} ${sizes[i]}`
  }

  const formatSpeed = (): string => {
    if (downloadStartTime.value === 0 || loaded.value === 0) {
      return ''
    }
    const elapsed = (Date.now() - downloadStartTime.value) / 1000
    return `${formatBytes(loaded.value / elapsed)}/s`
  }

  const startDownload = async () => {
    try {
      status.value = 'downloading'
      downloadStartTime.value = Date.now()
      downloadController = new AbortController()

      let downloadResult: { blob: Blob; filename?: string }

      if (props.options.shareKey) {
        downloadResult = await downloadSharedFile(
          props.options.fileId,
          props.options.shareKey,
          props.options.accessToken,
          {
            quality: props.options.quality,
            format: props.options.format,
          },
          (progressInfo) => {
            progress.value = progressInfo.percent
            loaded.value = progressInfo.loaded
            total.value = progressInfo.total
          }
        )
      } else {
        downloadResult = await downloadFile(
          props.options.fileId,
          {
            quality: props.options.quality,
            format: props.options.format,
          },
          (progressInfo) => {
            progress.value = progressInfo.percent
            loaded.value = progressInfo.loaded
            total.value = progressInfo.total
          }
        )
      }

      if (downloadController?.signal.aborted) {
        return
      }

      let downloadFileName = downloadResult.filename || fileName.value || `file_${props.options.fileId}`

      if (!downloadFileName.includes('.')) {
        downloadFileName += '.jpg'
      }

      const url = URL.createObjectURL(downloadResult.blob)
      const link = document.createElement('a')
      link.href = url
      link.download = downloadFileName
      document.body.appendChild(link)
      link.click()
      document.body.removeChild(link)
      URL.revokeObjectURL(url)

      status.value = 'success'
      toast.success($t('components.downloader.completed'))
      emit('complete', true)

      if (props.options.autoClose) {
        scheduleAutoClose()
      }
    } catch (error) {
      if (downloadController?.signal.aborted) {
        return
      }
      toast.error(error.message)

      console.error($t('components.downloader.downloadFailed'), error)
      status.value = 'error'
      emit('complete', false)

      scheduleAutoClose(5000)
    }
  }

  const handleMouseDown = (e: MouseEvent) => {
    clearAutoCloseTimer()

    dragStartTime.value = Date.now()
    const element = e.currentTarget as HTMLElement
    const rect = element.getBoundingClientRect()

    dragOffset.value = { x: e.clientX - rect.left, y: e.clientY - rect.top }
    currentPosition.value = { x: rect.left, y: rect.top }
    isDragging.value = true

    document.addEventListener('mousemove', handleMouseMove, { passive: false })
    document.addEventListener('mouseup', handleMouseUp, { passive: false })

    e.preventDefault()
    e.stopPropagation()
  }

  const handleMouseMove = (e: MouseEvent) => {
    if (!isDragging.value) {
      return
    }

    const elementSize = 120
    const maxX = window.innerWidth - elementSize
    const maxY = window.innerHeight - elementSize

    currentPosition.value.x = Math.max(0, Math.min(e.clientX - dragOffset.value.x, maxX))
    currentPosition.value.y = Math.max(0, Math.min(e.clientY - dragOffset.value.y, maxY))

    e.preventDefault()
  }

  const handleMouseUp = (e: MouseEvent) => {
    if (!isDragging.value) {
      return
    }

    isDragging.value = false
    savePosition(currentPosition.value.x, currentPosition.value.y)

    document.removeEventListener('mousemove', handleMouseMove)
    document.removeEventListener('mouseup', handleMouseUp)

    if ((status.value === 'success' || status.value === 'error') && props.options.autoClose) {
      const delay = status.value === 'success' ? 3000 : 5000
      scheduleAutoClose(delay)
    }

    e.preventDefault()
    e.stopPropagation()
  }

  const _handleClick = () => {
    const dragDuration = Date.now() - dragStartTime.value
    if (dragDuration < 100) {
      return
    }

    if (status.value === 'error') {
      progress.value = loaded.value = total.value = downloadStartTime.value = 0
      startDownload()
    } else if (status.value === 'success') {
      close()
    }
  }

  const handleMouseEnter = () => (showFileInfo.value = true)
  const handleMouseLeave = () => (showFileInfo.value = false)

  const close = () => {
    clearAutoCloseTimer()
    downloadController?.abort()
    visible.value = false
    emit('close')
  }

  let unsubscribeQueueChange: (() => void) | null = null

  const scheduleAutoClose = (delay: number = 3000) => {
    clearAutoCloseTimer()
    if (!isDragging.value) {
      autoCloseTimer = window.setTimeout(() => {
        close()
      }, delay)
    }
  }

  const clearAutoCloseTimer = () => {
    if (autoCloseTimer) {
      clearTimeout(autoCloseTimer)
      autoCloseTimer = null
    }
  }

  onMounted(() => {
    initializePosition()
    queueSize.value = props.options.queueSize || 0

    unsubscribeQueueChange = onQueueChange((newSize) => {
      queueSize.value = newSize
    })

    startDownload()
  })

  onUnmounted(() => {
    clearAutoCloseTimer()
    downloadController?.abort()
    document.removeEventListener('mousemove', handleMouseMove)
    document.removeEventListener('mouseup', handleMouseUp)

    if (unsubscribeQueueChange) {
      unsubscribeQueueChange()
    }
  })
</script>

<template>
  <Teleport to="body">
    <Transition name="cyber-float">
      <div
        v-if="visible"
        class="cyber-downloader-float"
        :style="floatStyle"
        @mousedown="handleMouseDown"
        @mouseenter="handleMouseEnter"
        @mouseleave="handleMouseLeave"
      >
        <div class="cyber-ring-outer">
          <div class="cyber-ring-main">
            <svg viewBox="0 0 86 86" class="progress-svg">
              <circle cx="43" cy="43" r="36" stroke="var(--color-border-subtle)" stroke-width="3" fill="none" class="ring-bg" />

              <circle
                cx="43"
                cy="43"
                r="36"
                stroke="url(#cyberGradient)"
                stroke-width="3"
                fill="none"
                stroke-linecap="round"
                class="ring-progress"
                :style="progressStyle"
              />

              <defs>
                <linearGradient id="cyberGradient" x1="0%" y1="0%" x2="100%" y2="100%">
                  <stop offset="0%" style="stop-color: var(--color-brand-500); stop-opacity: 1" />
                  <stop offset="50%" style="stop-color: var(--color-error-500); stop-opacity: 1" />
                  <stop offset="100%" style="stop-color: var(--color-brand-300); stop-opacity: 1" />
                </linearGradient>
                <filter id="glow">
                  <feGaussianBlur stdDeviation="3" result="coloredBlur" />
                  <feMerge>
                    <feMergeNode in="coloredBlur" />
                    <feMergeNode in="SourceGraphic" />
                  </feMerge>
                </filter>
              </defs>
            </svg>

            <div class="ring-center">
              <div v-if="status === 'downloading'" class="center-content">
                <div class="percent-line">
                  <span class="percent-display">{{ Math.round(progress) }}</span>
                  <span class="percent-symbol">%</span>
                </div>
                <div class="speed-display">{{ formatSpeed() }}</div>
                <div v-if="queueSize > 0" class="queue-display">
                  {{ $t('components.downloader.queueRemaining').replace('{count}', queueSize.toString()) }}
                </div>
              </div>

              <div v-else-if="status === 'success'" class="center-content success">
                <div class="status-icon">✓</div>
                <div class="status-text">{{ $t('components.downloader.completed') }}</div>
                <div v-if="queueSize > 0" class="queue-display">
                  {{ $t('components.downloader.queueRemaining').replace('{count}', queueSize.toString()) }}
                </div>
              </div>

              <div v-else-if="status === 'error'" class="center-content error">
                <div class="status-icon">✗</div>
                <div class="status-text">{{ $t('components.downloader.failed') }}</div>
                <div v-if="queueSize > 0" class="queue-display">
                  {{ $t('components.downloader.queueRemaining').replace('{count}', queueSize.toString()) }}
                </div>
              </div>

              <div v-else class="center-content preparing">
                <div class="loading-dots">
                  <span />
                  <span />
                  <span />
                </div>
                <div class="status-text">{{ $t('components.downloader.preparing') }}</div>
                <div v-if="queueSize > 0" class="queue-display">
                  {{ $t('components.downloader.queueRemaining').replace('{count}', queueSize.toString()) }}
                </div>
              </div>
            </div>
          </div>

          <div class="pulse-ring" :class="{ active: status === 'downloading' }" />
        </div>

        <div v-if="showFileInfo" class="file-tooltip">
          <div class="tooltip-content">
            <div class="file-name">{{ fileName || $t('components.downloader.downloading') }}</div>
            <div class="file-progress">{{ formatBytes(loaded) }} / {{ formatBytes(total) }}</div>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
  .cyber-downloader-float {
    @apply select-none transition-transform duration-200 ease-in-out;
  }

  .cyber-downloader-float:active {
    @apply scale-[0.98];
  }

  .cyber-ring-outer {
    @apply relative h-[90px] w-[90px];
  }

  .cyber-ring-main {
    @apply relative h-[90px] w-[90px] rounded-full border-2 border-brand-400 bg-background-800 transition-all duration-300 ease-in-out;
    box-shadow: var(--shadow-glow-md);
    backdrop-filter: blur(10px);
  }

  .cyber-ring-main:hover {
    @apply scale-105;
    box-shadow: var(--shadow-glow-lg);
  }

  .progress-svg {
    @apply absolute left-1/2 top-1/2 -translate-x-1/2 -translate-y-1/2 -rotate-90 transform;
    width: 86px;
    height: 86px;
    filter: drop-shadow(0 0 8px var(--color-brand-500));
  }

  .ring-bg {
    @apply opacity-30;
  }

  .ring-progress {
    stroke-width: 4;
    transition: stroke-dashoffset 0.3s ease;
    animation: pulse-glow 2s ease-in-out infinite alternate;
  }

  .ring-center {
    @apply absolute left-1/2 top-1/2 z-[2] -translate-x-1/2 -translate-y-1/2 transform text-center;
  }

  .center-content {
    @apply flex flex-col items-center gap-0.5;
  }

  .percent-line {
    @apply flex items-baseline gap-0.5 leading-none;
  }

  .percent-display {
    @apply font-mono text-lg font-bold text-brand-400;
    text-shadow: 0 0 8px var(--color-brand-500);
  }

  .percent-symbol {
    @apply ml-px text-xs font-medium text-error-500;
    text-shadow: 0 0 4px var(--color-error-500);
  }

  .speed-display {
    @apply mt-px font-mono text-[7px] text-content;
    text-shadow: 0 0 4px var(--color-warning-500);
  }

  .queue-display {
    @apply mt-px rounded-sm border border-badge-accent-border bg-badge-accent-bg px-1 py-px font-mono text-[7px] font-semibold text-error-500;
    text-shadow: 0 0 6px var(--color-error-500);
  }

  .center-content.success .status-icon {
    @apply text-2xl text-success-500;
    text-shadow: 0 0 12px var(--color-success-500);
    animation: scaleIn var(--transition-fast) var(--ease-out);
  }

  .center-content.error .status-icon {
    @apply text-xl text-error-500;
    text-shadow: 0 0 15px var(--color-error-500);
    animation: error-shake 0.6s ease-out;
  }

  .status-text {
    @apply mt-1 text-[10px] font-medium text-content;
    text-shadow: 0 0 5px var(--color-content-default);
  }

  .loading-dots {
    @apply mb-1 flex gap-1;
  }

  .loading-dots span {
    @apply h-1 w-1 rounded-full bg-brand-500;
    animation: loading-bounce 1.4s ease-in-out infinite both;
    box-shadow: 0 0 5px var(--color-brand-500);
  }

  .loading-dots span:nth-child(1) {
    animation-delay: -0.32s;
  }
  .loading-dots span:nth-child(2) {
    animation-delay: -0.16s;
  }
  .loading-dots span:nth-child(3) {
    animation-delay: 0;
  }

  .pulse-ring {
    @apply pointer-events-none absolute -left-2.5 -top-2.5 h-[110px] w-[110px] rounded-full border-2 border-brand-400 opacity-0;
  }

  .pulse-ring.active {
    animation: cyber-pulse 2s ease-out infinite;
  }

  .file-tooltip {
    @apply pointer-events-none absolute -top-16 left-1/2 -translate-x-1/2 opacity-0 transition-opacity duration-300;
  }

  .cyber-downloader-float:hover .file-tooltip {
    @apply opacity-100;
  }

  .tooltip-content {
    @apply whitespace-nowrap rounded-lg border border-brand-400 bg-background-900 px-3 py-2;
    box-shadow: 0 0 15px var(--color-brand-500);
  }

  .file-name {
    @apply mb-0.5 text-xs font-medium text-brand-400;
    text-shadow: 0 0 5px var(--color-brand-500);
  }

  .file-progress {
    @apply font-mono text-[10px] text-content;
    text-shadow: 0 0 3px var(--color-warning-500);
  }

  @keyframes pulse-glow {
    0% {
      filter: drop-shadow(0 0 5px var(--color-brand-500));
    }
    100% {
      filter: drop-shadow(0 0 15px var(--color-brand-400));
    }
  }

  @keyframes cyber-pulse {
    0% {
      transform: scale(1);
      opacity: 0.6;
    }
    50% {
      opacity: 0.2;
    }
    100% {
      transform: scale(1.3);
      opacity: 0;
    }
  }

  @keyframes loading-bounce {
    0%,
    80%,
    100% {
      transform: scale(0.8);
      opacity: 0.5;
    }
    40% {
      transform: scale(1.2);
      opacity: 1;
    }
  }

  @keyframes error-shake {
    0%,
    100% {
      transform: translateX(0);
    }
    25% {
      transform: translateX(-5px);
    }
    75% {
      transform: translateX(5px);
    }
  }

  .cyber-float-enter-active,
  .cyber-float-leave-active {
    transition: all 0.4s cubic-bezier(0.25, 0.8, 0.25, 1);
  }

  .cyber-float-enter-from {
    opacity: 0;
    transform: scale(0.5) rotate(180deg);
  }

  .cyber-float-leave-to {
    opacity: 0;
    transform: scale(0.5) rotate(-180deg);
  }

  @media (max-width: 768px) {
    .cyber-ring-outer,
    .cyber-ring-main {
      @apply h-20 w-20;
    }

    .progress-svg {
      width: 76px;
      height: 76px;
    }

    .percent-display {
      @apply text-base;
    }

    .pulse-ring {
      @apply -left-2.5 -top-2.5 h-[100px] w-[100px];
    }
  }
</style>
