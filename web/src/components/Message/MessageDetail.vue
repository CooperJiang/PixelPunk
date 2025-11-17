<script setup lang="ts">
  import { getCurrentLocale } from '@/utils/locale'
  import { computed, ref, watch } from 'vue'
  import { useRouter } from 'vue-router'
  import { deleteMessage as deleteMessageApi, markMessageRead, formatRelativeTime } from '@/api/message'
  import {
    MessagePriority,
    MessageStatus,
    getMessageTypeConfig,
    getPriorityConfig,
    getStatusConfig,
    type Message,
  } from '@/api/message/types'
  import { showConfirm } from '@/utils/dialog'
  import { useTexts } from '@/composables/useTexts'

  interface Props {
    message: Message
    show: boolean
  }

  const props = defineProps<Props>()

  const emit = defineEmits<{
    close: []
    updated: [message: Message]
    deleted: [messageId: number]
  }>()

  const router = useRouter()

  const { $t } = useTexts()

  const messageTypeConfig = getMessageTypeConfig($t)
  const priorityConfig = getPriorityConfig($t)
  const statusConfig = getStatusConfig($t)

  const dialogVisible = ref(false)

  const messageIcon = computed(() => messageTypeConfig[props.message.type as keyof typeof messageTypeConfig]?.icon || '📧')

  const messageTypeLabel = computed(
    () => messageTypeConfig[props.message.type as keyof typeof messageTypeConfig]?.label || props.message.type
  )

  const statusInfo = computed(
    () => statusConfig[props.message.status as keyof typeof statusConfig] || statusConfig[MessageStatus.UNREAD]
  )

  const priorityInfo = computed(
    () => priorityConfig[props.message.priority as keyof typeof priorityConfig] || priorityConfig[MessagePriority.NORMAL]
  )

  const hasMetaData = computed(() => props.message.meta_data && Object.keys(props.message.meta_data).length > 0)

  const handleClose = () => {
    dialogVisible.value = false
    emit('close')
  }

  const _markRead = async () => {
    if (props.message.status === MessageStatus.read) {
      return
    }

    try {
      await markMessageRead(props.message.id.toString())
      const updatedMessage = {
        ...props.message,
        status: MessageStatus.read,
        read_at: new Date().toISOString(),
      }
      emit('updated', updatedMessage)
    } catch {}
  }

  const _deleteMessage = async () => {
    if (!showConfirm($t('components.messageDetail.confirmDelete'))) {
      return
    }

    try {
      await deleteMessageApi(props.message.id.toString())
      emit('deleted', props.message.id)
      handleClose()
    } catch {}
  }

  const _handleAction = () => {
    if (props.message.action_url) {
      if (props.message.action_url.startsWith('http')) {
        window.open(props.message.action_url, '_blank')
      } else {
        router.push(props.message.action_url)
      }
      handleClose()
    }
  }

  const _getActionButtonType = (style?: string): string => {
    switch (style) {
      case 'primary':
        return 'primary'
      case 'success':
        return 'success'
      case 'warning':
        return 'warning'
      case 'danger':
        return 'danger'
      case 'secondary':
      default:
        return 'outlined'
    }
  }

  const formatMetaKey = (key: string): string => {
    const keyMap: Record<string, string> = {
      storage_mb: $t('components.messageDetail.metaKeys.storageMb'),
      bandwidth_mb: $t('components.messageDetail.metaKeys.bandwidthMb'),
      file_name: $t('components.messageDetail.metaKeys.fileName'),
      file_id: $t('components.messageDetail.metaKeys.fileId'),
      reason: $t('components.messageDetail.metaKeys.reason'),
      used_percent: $t('components.messageDetail.metaKeys.usedPercent'),
      used_mb: $t('components.messageDetail.metaKeys.usedMb'),
      total_mb: $t('components.messageDetail.metaKeys.totalMb'),
      granted_mb: $t('components.messageDetail.metaKeys.grantedMb'),
      start_time: $t('components.messageDetail.metaKeys.startTime'),
      end_time: $t('components.messageDetail.metaKeys.endTime'),
      notice_id: $t('components.messageDetail.metaKeys.noticeId'),
      review_id: $t('components.messageDetail.metaKeys.reviewId'),
    }
    return keyMap[key] || key
  }

  const formatMetaValue = (value: any): string => {
    if (value === null || value === undefined) {
      return '-'
    }
    if (typeof value === 'boolean') {
      return value ? $t('components.messageDetail.yes') : $t('components.messageDetail.no')
    }
    return String(value)
  }

  const formatDateTime = (dateString: string): string => {
    const date = new Date(dateString)
    return date.toLocaleString(getCurrentLocale(), {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit',
    })
  }

  watch(
    () => props.show,
    (_newShow) => {},
    { immediate: true }
  )

  watch(dialogVisible, (newValue) => {
    if (!newValue) {
      emit('close')
    }
  })
</script>

<template>
  <CyberDialog
    v-if="false"
    v-model="dialogVisible"
    :title="message.title"
    width="600px"
    max-width="95vw"
    max-height="90vh"
    :show-default-footer="false"
    @close="handleClose"
  >
    <div class="message-header mb-4 rounded-lg border border-cyan-800/30 bg-gradient-to-r from-cyan-900/20 to-purple-900/20 p-4">
      <div class="mb-2 flex items-center space-x-3">
        <span class="text-2xl">{{ messageIcon }}</span>
        <div class="flex-1">
          <h3 class="text-lg font-medium text-cyan-400">{{ message.title }}</h3>
          <p class="text-sm text-content-muted">{{ messageTypeLabel }} · {{ formatRelativeTime(message.created_at, $t) }}</p>
        </div>
      </div>

      <div class="mt-3 flex items-center space-x-2">
        <span class="inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium" :class="statusInfo.class">
          {{ statusInfo.label }}
        </span>

        <span
          v-if="message.priority === MessagePriority.HIGH"
          class="inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium"
          :class="priorityInfo.class"
        >
          {{ priorityInfo.label }}{{ $t('components.messageDetail.priority') }}
        </span>

        <span
          class="inline-flex items-center rounded-full bg-background-700 px-2.5 py-0.5 text-xs font-medium text-content-muted"
        >
          {{ messageTypeLabel }}
        </span>
      </div>
    </div>

    <div class="message-content mb-4">
      <div class="border-default-light rounded-lg border bg-background-800 p-4">
        <p class="whitespace-pre-wrap leading-relaxed text-content-muted">
          {{ message.content }}
        </p>
      </div>
    </div>

    <div v-if="message.related_type || message.related_id" class="mb-4">
      <h4 class="mb-2 text-sm font-medium text-cyan-400">{{ $t('components.messageDetail.relatedInfo') }}</h4>
      <div class="rounded-lg border border-cyan-800/30 bg-cyan-900/10 p-3">
        <div class="space-y-1 text-sm text-content-muted">
          <p v-if="message.related_type">
            {{ $t('components.messageDetail.type') }}: <span class="text-content-muted">{{ message.related_type }}</span>
          </p>
          <p v-if="message.related_id">
            {{ $t('components.messageDetail.id') }}: <span class="text-content-muted">{{ message.related_id }}</span>
          </p>
        </div>
      </div>
    </div>

    <div v-if="hasMetaData" class="mb-4">
      <h4 class="mb-2 text-sm font-medium" style="color: var(--color-brand-500)">
        {{ $t('components.messageDetail.detailedInfo') }}
      </h4>
      <div
        class="rounded-lg border p-3"
        style="border-color: var(--color-border-subtle); background-color: var(--color-background-700)"
      >
        <div class="text-sm text-content-muted">
          <div v-for="(value, key) in message.meta_data" :key="key" class="flex justify-between py-1">
            <span class="font-medium">{{ formatMetaKey(key) }}:</span>
            <span class="text-content-muted">{{ formatMetaValue(value) }}</span>
          </div>
        </div>
      </div>
    </div>

    <div v-if="message.expires_at" class="mb-4">
      <div class="flex items-center text-xs text-content-disabled">
        <i class="fas fa-clock mr-2" />
        {{ $t('components.messageDetail.expiresAt') }}: {{ formatDateTime(message.expires_at) }}
      </div>
    </div>

    <template #footer>
      <div class="flex justify-end space-x-3 border-t border-cyan-800/30 p-4">
        <CyberButton type="outlined" @click="handleClose"> {{ $t('components.messageDetail.close') }} </CyberButton>
      </div>
    </template>
  </CyberDialog>
</template>

<style scoped>
  .message-header {
    background: linear-gradient(135deg, rgba(6, 182, 212, 0.1) 0%, rgba(147, 51, 234, 0.1) 100%);
  }

  .cyber-button.text-red-400 {
    color: var(--color-danger);
    border-color: var(--color-danger);
  }

  .cyber-button.text-red-400:hover {
    background-color: var(--color-danger);
    color: var(--color-content-default);
  }
</style>
