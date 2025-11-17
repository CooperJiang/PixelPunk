<script setup lang="ts">
  import { computed, ref, watch } from 'vue'
  import { type UploadSessionData, UploadStorageManager } from '@/utils/storage/uploadStorage'
  import type { ResumableUploadsProps, ResumableUploadsEmits } from './types'
  import { useTexts } from '@/composables/useTexts'

  const props = defineProps<ResumableUploadsProps>()
  const emit = defineEmits<ResumableUploadsEmits>()
  const { $t } = useTexts()

  /* 对话框显示状态 */
  const dialogVisible = computed({
    get: () => props.modelValue,
    set: (value) => emit('update:modelValue', value),
  })

  /* 可恢复的上传会话 */
  const resumableSessions = ref<UploadSessionData[]>([])
  const selectedSessions = ref<Set<string>>(new Set())

  const loadResumableSessions = () => {
    resumableSessions.value = UploadStorageManager.getResumableSessions()
    selectedSessions.value.clear()
  }

  watch(
    () => props.modelValue,
    (visible) => {
      if (visible) {
        loadResumableSessions()
      }
    }
  )

  const toggleSession = (sessionId: string) => {
    if (selectedSessions.value.has(sessionId)) {
      selectedSessions.value.delete(sessionId)
    } else {
      selectedSessions.value.add(sessionId)
    }
  }

  const selectAll = () => {
    resumableSessions.value.forEach((session) => {
      selectedSessions.value.add(session.id)
    })
  }

  const clearSelection = () => {
    selectedSessions.value.clear()
  }

  const handleResumeSelected = () => {
    const selected = resumableSessions.value.filter((session) => selectedSessions.value.has(session.id))

    if (selected.length > 0) {
      emit('resume', selected)
      dialogVisible.value = false
    }
  }

  const handleDeleteSelected = () => {
    selectedSessions.value.forEach((sessionId) => {
      UploadStorageManager.removeSession(sessionId)
    })

    loadResumableSessions()
  }

  const handleClose = () => {
    dialogVisible.value = false
  }

  const formatFileSize = (bytes: number): string => {
    if (bytes === 0) {
      return '0 B'
    }
    const k = 1024
    const sizes = ['B', 'KB', 'MB', 'GB']
    const i = Math.floor(Math.log(bytes) / Math.log(k))
    return `${parseFloat((bytes / Math.pow(k, i)).toFixed(2))} ${sizes[i]}`
  }

  const formatTime = (timestamp: number): string => {
    const now = Date.now()
    const diff = now - timestamp

    if (diff < 60000) {
      return $t('resumableUploads.relativeTime.justNow') as string
    } else if (diff < 3600000) {
      return ($t('resumableUploads.relativeTime.minutesAgo') as string).replace('{count}', String(Math.floor(diff / 60000)))
    } else if (diff < 86400000) {
      return ($t('resumableUploads.relativeTime.hoursAgo') as string).replace('{count}', String(Math.floor(diff / 3600000)))
    }
    return ($t('resumableUploads.relativeTime.daysAgo') as string).replace('{count}', String(Math.floor(diff / 86400000)))
  }
</script>

<template>
  <CyberDialog
    v-model="dialogVisible"
    :title="$t('resumableUploads.title')"
    width="600px"
    @confirm="handleResumeSelected"
    @close="handleClose"
  >
    <div class="resumable-uploads">
      <div v-if="resumableSessions.length === 0" class="empty-state py-8 text-center">
        <i class="fas fa-cloud-upload-alt mb-4 text-4xl text-brand-400 opacity-40" />
        <p class="text-brand-500">{{ $t('resumableUploads.empty.title') }}</p>
        <p class="mt-2 text-sm text-content-muted">{{ $t('resumableUploads.empty.subtitle') }}</p>
      </div>

      <div v-else class="space-y-3">
        <div
          v-for="session in resumableSessions"
          :key="session.id"
          class="upload-session rounded-lg border border-default bg-background-200 p-4 transition-all"
          :class="{ 'border-transparent ring-2 ring-brand-300': selectedSessions.has(session.id) }"
          @click="toggleSession(session.id)"
        >
          <div class="flex items-center justify-between">
            <div class="flex items-center space-x-3">
              <cyberCheckbox :model-value="selectedSessions.has(session.id)" @update:model-value="toggleSession(session.id)" />

              <div class="file-info">
                <h4 class="font-medium text-content">{{ session.fileName }}</h4>
                <div class="flex items-center space-x-4 text-sm text-content-muted">
                  <span>{{ formatFileSize(session.fileSize) }}</span>
                  <span
                    >{{ session.uploadedChunks.length }}/{{ session.totalChunks }}
                    {{ $t('resumableUploads.labels.chunks') }}</span
                  >
                  <span>{{ Math.round((session.uploadedChunks.length / session.totalChunks) * 100) }}%</span>
                </div>
              </div>
            </div>

            <div class="text-xs text-content-muted">
              {{ formatTime(session.lastActivity) }}
            </div>
          </div>

          <div class="mt-3 h-2 w-full rounded-full bg-background-500">
            <div
              class="h-2 rounded-full bg-brand-500 transition-all duration-300"
              :style="{
                width: Math.round((session.uploadedChunks.length / session.totalChunks) * 100) + '%',
              }"
            />
          </div>
        </div>
      </div>

      <div v-if="resumableSessions.length > 0" class="mt-6 flex items-center justify-between">
        <div class="flex space-x-2">
          <button
            class="rounded bg-brand-100 px-3 py-1 text-sm text-brand-500 transition-colors hover:bg-brand-200"
            @click.stop.prevent="selectAll"
          >
            {{ $t('resumableUploads.actions.selectAll') }}
          </button>
          <button
            class="rounded bg-background-400 px-3 py-1 text-sm text-content transition-colors hover:bg-background-300"
            @click.stop.prevent="clearSelection"
          >
            {{ $t('resumableUploads.actions.clearSelection') }}
          </button>
        </div>

        <div class="flex space-x-2">
          <button
            :disabled="selectedSessions.size === 0"
            class="rounded bg-error-100 px-3 py-1 text-sm text-error-500 transition-colors hover:bg-error-200 disabled:opacity-50"
            @click.stop.prevent="handleDeleteSelected"
          >
            {{ $t('resumableUploads.actions.deleteSelected') }} ({{ selectedSessions.size }})
          </button>
        </div>
      </div>
    </div>

    <template #footer>
      <div class="flex justify-end space-x-3">
        <CyberButton variant="secondary" @click="handleClose"> {{ $t('resumableUploads.footer.cancel') }} </CyberButton>
        <CyberButton type="primary" :disabled="selectedSessions.size === 0" @click="handleResumeSelected">
          {{ $t('resumableUploads.footer.resume') }} ({{ selectedSessions.size }})
        </CyberButton>
      </div>
    </template>
  </CyberDialog>
</template>

<style scoped>
  .upload-session {
    cursor: pointer;
  }

  .upload-session:hover {
    background-color: rgba(var(--color-brand-500-rgb), 0.1);
  }
</style>
