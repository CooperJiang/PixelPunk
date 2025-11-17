<script setup lang="ts">
  import { computed } from 'vue'
  import { storeToRefs } from 'pinia'
  import { useUploadStore } from '@/store/upload'
  import { useToast } from '@/components/Toast/useToast'
  import { useRouter } from 'vue-router'
  import { useTexts } from '@/composables/useTexts'
  import type { GlobalUploadDrawerEmits } from './types'

  defineOptions({
    name: 'GlobalUploadDrawer',
  })

  const emit = defineEmits<GlobalUploadDrawerEmits>()

  const uploadStore = useUploadStore()
  const toast = useToast()
  const router = useRouter()
  const { $t } = useTexts()

  const { allUploads, statistics } = storeToRefs(uploadStore)

  const hasUploadingOrPending = computed(() => statistics.value.uploading > 0 || statistics.value.pending > 0)

  const hasPaused = computed(() => allUploads.value.some((item) => item.status === 'paused'))

  /* 格式化文件大小 */
  const formatFileSize = (size: number): string => {
    if (size < 1024) {
      return `${size} B`
    }
    if (size < 1024 * 1024) {
      return `${(size / 1024).toFixed(1)} KB`
    }
    if (size < 1024 * 1024 * 1024) {
      return `${(size / (1024 * 1024)).toFixed(1)} MB`
    }
    return `${(size / (1024 * 1024 * 1024)).toFixed(2)} GB`
  }

  const formatSpeed = (speed: number): string => {
    if (speed > 1024) {
      return `${(speed / 1024).toFixed(1)} MB/s`
    }
    return `${speed.toFixed(0)} KB/s`
  }

  const pauseUpload = (id: string) => {
    uploadStore.pauseUpload(id)
  }

  const resumeUpload = (id: string) => {
    uploadStore.resumeUpload(id)
  }

  const retryUpload = (id: string) => {
    uploadStore.retryUpload(id)
  }

  const removeUpload = (id: string) => {
    uploadStore.removeUploadItem(id)
  }

  const pauseAll = () => {
    uploadStore.pauseUpload()
    toast.success($t('globalUploadDrawer.messages.pausedAll'))
  }

  const resumeAll = () => {
    allUploads.value.filter((item) => item.status === 'paused').forEach((item) => uploadStore.resumeUpload(item.id))
    toast.success($t('globalUploadDrawer.messages.resumedAll'))
  }

  const clearCompleted = () => {
    const completed = allUploads.value.filter((item) => item.status === 'completed')
    completed.forEach((item) => uploadStore.removeUploadItem(item.id))
    toast.success($t('globalUploadDrawer.messages.clearedCompleted', { count: completed.length }))
  }

  const clearAll = () => {
    uploadStore.clearQueue()
    toast.success($t('globalUploadDrawer.messages.clearedAll'))
  }

  const copyUrl = async (url: string) => {
    try {
      await navigator.clipboard.writeText(url)
      toast.success($t('globalUploadDrawer.messages.linkCopied'))
    } catch {
      toast.error($t('globalUploadDrawer.messages.copyFailed'))
    }
  }

  const handlePreviewError = (e: Event) => {
    const img = e.target as HTMLImageElement
    img.style.display = 'none'
  }

  const goToUploadPage = () => {
    emit('close')
    router.push('/upload')
  }
</script>

<template>
  <div class="upload-drawer-container">
    <div class="upload-stats">
      <div class="stat-item">
        <i class="fas fa-cloud-upload-alt" />
        <div class="stat-value">{{ statistics.uploading }}</div>
        <div class="stat-label">{{ $t('globalUploadDrawer.stats.uploading') }}</div>
      </div>
      <div class="stat-item">
        <i class="fas fa-clock" />
        <div class="stat-value">{{ statistics.pending }}</div>
        <div class="stat-label">{{ $t('globalUploadDrawer.stats.pending') }}</div>
      </div>
      <div class="stat-item">
        <i class="fas fa-check-circle" />
        <div class="stat-value">{{ statistics.completed }}</div>
        <div class="stat-label">{{ $t('globalUploadDrawer.stats.completed') }}</div>
      </div>
      <div v-if="statistics.failed > 0" class="stat-item error">
        <i class="fas fa-exclamation-circle" />
        <div class="stat-value">{{ statistics.failed }}</div>
        <div class="stat-label">{{ $t('globalUploadDrawer.stats.failed') }}</div>
      </div>
    </div>

    <div v-if="allUploads.length > 0" class="batch-actions">
      <button
        v-if="hasUploadingOrPending"
        class="action-btn pause"
        :title="$t('globalUploadDrawer.batchActions.pauseAll')"
        @click="pauseAll"
      >
        <i class="fas fa-pause" />
        <span>{{ $t('globalUploadDrawer.batchActions.pauseAll') }}</span>
      </button>

      <button
        v-if="hasPaused"
        class="action-btn resume"
        :title="$t('globalUploadDrawer.batchActions.resumeAll')"
        @click="resumeAll"
      >
        <i class="fas fa-play" />
        <span>{{ $t('globalUploadDrawer.batchActions.resumeAll') }}</span>
      </button>

      <button
        class="action-btn clear"
        :title="$t('globalUploadDrawer.batchActions.clearCompleted')"
        :disabled="statistics.completed === 0"
        @click="clearCompleted"
      >
        <i class="fas fa-broom" />
        <span>{{ $t('globalUploadDrawer.batchActions.clearCompleted') }}</span>
      </button>

      <button class="action-btn danger" :title="$t('globalUploadDrawer.batchActions.clearAll')" @click="clearAll">
        <i class="fas fa-trash" />
        <span>{{ $t('globalUploadDrawer.batchActions.clearAll') }}</span>
      </button>
    </div>

    <div v-if="allUploads.length > 0" class="upload-list">
      <TransitionGroup name="list-item">
        <div v-for="item in allUploads" :key="item.id" class="upload-item" :class="[`status-${item.status}`]">
          <div class="file-preview">
            <img v-if="item.preview" :src="item.preview" alt="" class="preview-img" @error="handlePreviewError" />
            <div v-else class="preview-placeholder">
              <i class="fas fa-file-image" />
            </div>
          </div>

          <div class="file-info">
            <div class="file-name" :title="item.file.name">
              {{ item.file.name }}
            </div>

            <div class="file-meta">
              <span class="file-size">{{ formatFileSize(item.file.size) }}</span>
              <span v-if="item.type === 'chunked'" class="file-type">
                <i class="fas fa-puzzle-piece" /> {{ $t('globalUploadDrawer.misc.chunkedUpload') }}
              </span>
            </div>

            <div v-if="item.status === 'uploading' || item.status === 'preparing'" class="progress-bar">
              <div class="progress-track">
                <div class="progress-fill" :style="{ width: `${item.progress}%` }" />
              </div>
              <span class="progress-text">{{ item.progress }}%</span>
            </div>

            <div class="status-info">
              <span v-if="item.status === 'uploading'" class="status-text uploading">
                <i class="fas fa-spinner fa-spin" />
                {{ item.statusMessage || $t('globalUploadDrawer.statusMessages.uploading') }}
                <span v-if="item.speed > 0" class="speed">
                  {{ formatSpeed(item.speed) }}
                </span>
              </span>

              <span v-else-if="item.status === 'completed'" class="status-text success">
                <i class="fas fa-check" />
                {{ $t('globalUploadDrawer.statusMessages.completed') }}
              </span>

              <span v-else-if="item.status === 'failed'" class="status-text error">
                <i class="fas fa-times" />
                {{ item.error || $t('globalUploadDrawer.statusMessages.failed') }}
              </span>

              <span v-else-if="item.status === 'paused'" class="status-text paused">
                <i class="fas fa-pause" />
                {{ $t('globalUploadDrawer.statusMessages.paused') }}
              </span>

              <span v-else-if="item.status === 'pending'" class="status-text pending">
                <i class="fas fa-clock" />
                {{ $t('globalUploadDrawer.statusMessages.pending') }}
              </span>

              <span v-else class="status-text preparing">
                <i class="fas fa-hourglass-half" />
                {{ $t('globalUploadDrawer.statusMessages.preparing') }}
              </span>
            </div>
          </div>

          <div class="file-actions">
            <button
              v-if="item.status === 'uploading'"
              class="btn-action pause"
              :title="$t('globalUploadDrawer.actions.pause')"
              @click="pauseUpload(item.id)"
            >
              <i class="fas fa-pause" />
            </button>

            <button
              v-else-if="item.status === 'paused'"
              class="btn-action resume"
              :title="$t('globalUploadDrawer.actions.resume')"
              @click="resumeUpload(item.id)"
            >
              <i class="fas fa-play" />
            </button>

            <button
              v-else-if="item.status === 'failed'"
              class="btn-action retry"
              :title="$t('globalUploadDrawer.actions.retry')"
              @click="retryUpload(item.id)"
            >
              <i class="fas fa-redo" />
            </button>

            <button
              v-if="item.status === 'completed' && item.result?.full_url"
              class="btn-action copy"
              :title="$t('globalUploadDrawer.actions.copyLink')"
              @click="copyUrl(item.result.full_url)"
            >
              <i class="fas fa-copy" />
            </button>

            <button class="btn-action remove" :title="$t('globalUploadDrawer.actions.remove')" @click="removeUpload(item.id)">
              <i class="fas fa-times" />
            </button>
          </div>
        </div>
      </TransitionGroup>
    </div>

    <div v-else class="empty-state">
      <i class="fas fa-cloud-upload-alt" />
      <p>{{ $t('globalUploadDrawer.misc.noTasks') }}</p>
      <a class="go-upload" @click="goToUploadPage"> {{ $t('globalUploadDrawer.misc.goToUpload') }} </a>
    </div>

    <div class="drawer-footer">
      <a class="footer-link" @click="goToUploadPage">
        <span>{{ $t('globalUploadDrawer.misc.viewFullPage') }}</span>
        <i class="fas fa-arrow-right" />
      </a>
    </div>
  </div>
</template>

<style scoped>
  .upload-drawer-container {
    @apply flex h-full flex-col p-5;
  }

  .upload-stats {
    @apply mb-4 flex gap-4 rounded-lg border border-brand-500 bg-background-900 p-3;
  }

  .stat-item {
    @apply flex min-w-0 flex-1 flex-col items-center gap-1 text-brand-500;
  }

  .stat-item i {
    @apply mb-1 text-lg;
  }

  .stat-value {
    @apply text-xl font-bold leading-none;
  }

  .stat-label {
    @apply whitespace-nowrap text-xs opacity-80;
  }

  .stat-item.error {
    @apply text-error-500;
  }

  .batch-actions {
    @apply mb-4 flex gap-2;
  }

  .action-btn {
    @apply flex cursor-pointer items-center gap-1 rounded-md border border-subtle bg-background-800 px-2.5 py-1.5 text-xs text-content transition-all;
  }

  .action-btn:hover {
    @apply border-content-muted bg-background-700;
  }

  .action-btn:disabled {
    @apply cursor-not-allowed opacity-50;
  }

  .action-btn.pause {
    @apply border-warning-500 text-warning-500;
  }

  .action-btn.resume {
    @apply border-success-500 text-success-500;
  }

  .action-btn.clear {
    @apply border-brand-500 text-brand-500;
  }

  .action-btn.danger {
    @apply border-error-500 text-error-500;
  }

  .upload-list {
    @apply mb-5 flex-1 overflow-y-auto;
  }

  .upload-item {
    @apply mb-3 flex gap-4 rounded-lg border border-subtle bg-background-800 p-4 transition-all;
  }

  .upload-item:hover {
    @apply border-brand-500 bg-background-700;
  }

  .file-preview {
    @apply h-16 w-16 shrink-0 overflow-hidden rounded-md bg-background-900;
  }

  .preview-img {
    @apply h-full w-full object-cover;
  }

  .preview-placeholder {
    @apply flex h-full w-full items-center justify-center text-2xl text-content-muted;
  }

  .file-info {
    @apply flex min-w-0 flex-1 flex-col gap-2;
  }

  .file-name {
    @apply overflow-hidden text-ellipsis whitespace-nowrap text-sm font-medium text-content;
  }

  .file-meta {
    @apply flex gap-4 text-xs text-content-muted;
  }

  .file-type {
    @apply text-brand-600;
  }

  .progress-bar {
    @apply flex items-center gap-3;
  }

  .progress-track {
    @apply h-1 flex-1 overflow-hidden rounded-sm bg-background-700;
  }

  .progress-fill {
    @apply h-full rounded-sm bg-gradient-to-r from-brand-500 to-brand-600 transition-all;
    box-shadow: 0 0 10px var(--color-brand-500);
  }

  .progress-text {
    @apply min-w-10 text-right text-xs text-brand-500;
  }

  .status-info {
    @apply text-xs;
  }

  .status-text {
    @apply inline-flex items-center gap-1.5;
  }

  .status-text.uploading {
    @apply text-brand-500;
  }

  .status-text.success {
    @apply text-success-500;
  }

  .status-text.error {
    @apply text-error-500;
  }

  .status-text.paused {
    @apply text-warning-500;
  }

  .status-text.pending {
    @apply text-content-muted;
  }

  .status-text.preparing {
    @apply text-brand-600;
  }

  .speed {
    @apply ml-2 text-content;
  }

  .file-actions {
    @apply flex items-center gap-2;
  }

  .btn-action {
    @apply flex h-8 w-8 cursor-pointer items-center justify-center rounded-md border border-subtle bg-background-800 text-content-muted transition-all;
  }

  .btn-action:hover {
    @apply scale-110 bg-background-700 text-content;
  }

  .btn-action.pause {
    @apply border-warning-500 text-warning-500;
  }

  .btn-action.resume {
    @apply border-success-500 text-success-500;
  }

  .btn-action.retry {
    @apply border-warning-500 text-warning-500;
  }

  .btn-action.copy {
    @apply border-brand-500 text-brand-500;
  }

  .btn-action.remove {
    @apply border-error-500 text-error-500;
  }

  .empty-state {
    @apply flex flex-1 flex-col items-center justify-center gap-5 text-content-muted;
  }

  .empty-state i {
    @apply text-6xl text-brand-200;
  }

  .empty-state p {
    @apply text-base;
  }

  .go-upload {
    @apply inline-block cursor-pointer rounded-md bg-gradient-to-br from-brand-500 to-brand-600 px-5 py-2.5 text-content-heading no-underline transition-all;
  }

  .go-upload:hover {
    @apply -translate-y-0.5 shadow-glow-md;
  }

  .drawer-footer {
    @apply border-t border-subtle pt-2.5;
  }

  .footer-link {
    @apply flex cursor-pointer items-center justify-between py-2 text-xs text-content-muted no-underline transition-all;
  }

  .footer-link:hover {
    @apply text-brand-500;
  }

  .footer-link i {
    @apply text-xs transition-transform;
  }

  .footer-link:hover i {
    @apply translate-x-1;
  }

  .list-item-enter-active,
  .list-item-leave-active {
    @apply transition-all;
  }

  .list-item-enter-from {
    @apply -translate-x-8 opacity-0;
  }

  .list-item-leave-to {
    @apply translate-x-8 opacity-0;
  }

  .upload-list::-webkit-scrollbar {
    @apply w-1.5;
  }

  .upload-list::-webkit-scrollbar-track {
    @apply bg-background-900;
  }

  .upload-list::-webkit-scrollbar-thumb {
    @apply rounded bg-brand-300;
  }

  .upload-list::-webkit-scrollbar-thumb:hover {
    @apply bg-brand-500;
  }
</style>
