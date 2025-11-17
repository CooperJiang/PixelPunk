<script setup lang="ts">
  import { computed, h, onMounted, ref, resolveComponent, watch } from 'vue'
  import { useRouter } from 'vue-router'
  import { useTexts } from '@/composables/useTexts'
  import { useToast } from '@/components/Toast/useToast'
  import {
    deleteMessage as deleteMessageApi,
    getUnreadCount,
    getUserMessages,
    markAllMessagesRead,
    markMessageRead,
    formatRelativeTime,
  } from '@/api/message'
  import { MessagePriority, MessageStatus, getMessageTypeConfig, getPriorityConfig, type Message } from '@/api/message/types'
  import { formatMessage } from '@/utils/message/formatter'
  import type { MessageListProps } from './types'

  defineOptions({
    name: 'MessageList',
  })

  const CyberButton = resolveComponent('CyberButton')
  const toast = useToast()

  const props = withDefaults(defineProps<MessageListProps>(), {
    pageSize: 12,
  })

  const router = useRouter()
  const { $t } = useTexts()

  const messageTypeConfig = getMessageTypeConfig($t)
  const priorityConfig = getPriorityConfig($t)

  const loading = ref(false)
  const messages = ref<Message[]>([])
  const pageSize = ref(props.pageSize)
  const showSkeleton = ref(false)

  const statusOptions = ref([
    { label: $t('message.messageList.allStatus'), value: '' },
    { label: $t('message.messageList.statusLabels.unread'), value: MessageStatus.UNREAD },
    { label: $t('message.messageList.statusLabels.read'), value: MessageStatus.read },
  ])

  const typeOptions = ref([
    { label: $t('message.messageList.allTypes'), value: '' },
    ...Object.entries(messageTypeConfig).map(([type, config]) => ({
      label: config.label,
      value: type,
    })),
  ])

  const tableColumns = [
    {
      key: 'type',
      dataIndex: 'type',
      title: $t('message.messageList.columns.type'),
      width: 180,
      align: 'center',
      ellipsis: true,
      render: (value: string, record: Message) =>
        h('div', { class: 'flex items-center justify-center' }, [
          h('span', { class: 'mr-2 text-base' }, getMessageIcon(record.type)),
          h('span', { class: 'truncate', title: getMessageTypeLabel(record.type) }, getMessageTypeLabel(record.type)),
        ]),
    },
    {
      key: 'title',
      dataIndex: 'title',
      title: $t('message.messageList.columns.title'),
      width: 150,
      ellipsis: true,
      align: 'center',
      render: (value: string, record: Message) => {
        const formatted = formatMessage(record, $t)
        return h(
          'span',
          {
            title: formatted.title,
            class: 'block truncate',
          },
          formatted.title
        )
      },
    },
    {
      key: 'content',
      dataIndex: 'content',
      title: $t('message.messageList.columns.content'),
      width: 300,
      ellipsis: true,
      align: 'center',
      render: (value: string, record: Message) => {
        const formatted = formatMessage(record, $t)
        const content = formatted.content
        return h(
          'span',
          {
            title: content,
          },
          content.length > 50 ? `${content.substring(0, 50)}...` : content
        )
      },
    },
    {
      key: 'status',
      dataIndex: 'status',
      title: $t('message.messageList.columns.status'),
      width: 100,
      align: 'center',
      render: (value: number, record: Message) => {
        if (record.status === MessageStatus.UNREAD) {
          return h('span', { class: 'text-brand-400' }, $t('message.messageList.statusLabels.unread'))
        }
        return h('span', { class: 'text-content-muted' }, $t('message.messageList.statusLabels.read'))
      },
    },
    {
      key: 'priority',
      dataIndex: 'priority',
      title: $t('message.messageList.columns.priority'),
      width: 100,
      align: 'center',
      render: (value: number, record: Message) =>
        h('div', { class: 'flex justify-center' }, [
          h(
            'span',
            {
              class: `inline-flex items-center px-2 py-1 rounded text-xs font-medium ${getPriorityBadgeClass(record.priority)}`,
            },
            getPriorityLabel(record.priority)
          ),
        ]),
    },
    {
      key: 'created_at',
      dataIndex: 'created_at',
      title: $t('message.messageList.columns.time'),
      width: 120,
      align: 'center',
      render: (value: string, record: Message) => formatRelativeTime(record.created_at, $t),
    },
    {
      key: 'actions',
      title: $t('message.messageList.columns.actions'),
      width: 160,
      align: 'center',
      render: (value: any, record: Message) =>
        h('div', { class: 'flex gap-2 justify-center items-center' }, [
          ...(record.status === MessageStatus.UNREAD
            ? [
                h(
                  CyberButton,
                  {
                    type: 'outlined',
                    size: 'small',
                    icon: 'check',
                    onClick: (e: Event) => {
                      e.stopPropagation()
                      markRead(record.id)
                    },
                  },
                  () => $t('message.messageList.actions.markRead')
                ),
              ]
            : []),

          h(
            CyberButton,
            {
              type: 'outlined',
              size: 'small',
              icon: 'trash',
              onClick: (e: Event) => {
                e.stopPropagation()
                confirmDeleteMessage(record)
              },
            },
            () => $t('message.messageList.actions.delete')
          ),
        ]),
    },
  ]
  const unreadCount = ref(0)
  const totalMessages = ref(0)
  const totalPages = ref(0)
  const currentPage = ref(1)

  const filterStatus = ref('')
  const filterType = ref('')

  const showDeleteDialog = ref(false)
  const isDeleting = ref(false)
  const selectedMessage = ref<Message | null>(null)

  const _hasUnreadMessages = computed(() => unreadCount.value > 0)

  const loadMessages = async () => {
    loading.value = true
    showSkeleton.value = true
    try {
      const params: any = {
        page: currentPage.value,
        pageSize: pageSize.value,
      }

      if (filterStatus.value) {
        params.status = filterStatus.value
      }
      if (filterType.value) {
        params.type = filterType.value
      }

      const response = await getUserMessages(params)
      const data = response?.data
      if (data) {
        messages.value = data.items || []
        totalMessages.value = data.pagination?.total || 0
        totalPages.value = data.pagination?.last_page || Math.ceil((data.pagination?.total || 0) / pageSize.value)
      } else {
        messages.value = []
        totalMessages.value = 0
        totalPages.value = 1
      }
    } catch {
      messages.value = []
      totalMessages.value = 0
      totalPages.value = 1
    } finally {
      loading.value = false
      showSkeleton.value = false
    }
  }

  const loadUnreadCount = async () => {
    try {
      const response = await getUnreadCount()
      unreadCount.value = (response as any)?.data?.count || 0
    } catch {
      unreadCount.value = 0
    }
  }

  const markRead = async (messageId: number) => {
    try {
      await markMessageRead(messageId.toString())
      const message = messages.value.find((m) => m.id === messageId)
      if (message && message.status === MessageStatus.UNREAD) {
        message.status = MessageStatus.read
        message.read_at = new Date().toISOString()
        unreadCount.value = Math.max(0, unreadCount.value - 1)
      }
      toast.success($t('message.messageList.toasts.markReadSuccess'))
    } catch (error) {
      toast.error($t('message.messageList.toasts.markReadFailed'))
    }
  }

  const markAllRead = async () => {
    try {
      await markAllMessagesRead()
      messages.value.forEach((message) => {
        if (message.status === MessageStatus.UNREAD) {
          message.status = MessageStatus.read
          message.read_at = new Date().toISOString()
        }
      })
      unreadCount.value = 0
      toast.success($t('message.messageList.toasts.markAllReadSuccess'))
    } catch (error) {
      toast.error($t('message.messageList.toasts.markAllReadFailed'))
    }
  }

  const confirmDeleteMessage = (message: Message) => {
    selectedMessage.value = message
    showDeleteDialog.value = true
  }

  const closeDeleteDialog = () => {
    showDeleteDialog.value = false
    selectedMessage.value = null
  }

  const deleteMessage = async () => {
    if (!selectedMessage.value) {
      return
    }

    try {
      isDeleting.value = true
      await deleteMessageApi(selectedMessage.value.id.toString())

      const index = messages.value.findIndex((m) => m.id === selectedMessage.value!.id)
      if (index !== -1) {
        const message = messages.value[index]
        messages.value.splice(index, 1)
        totalMessages.value--
        if (message.status === MessageStatus.UNREAD) {
          unreadCount.value = Math.max(0, unreadCount.value - 1)
        }
      }

      toast.success($t('message.messageList.toasts.deleted'))
      closeDeleteDialog()
    } catch (error) {
      toast.error($t('message.messageList.toasts.deleteFailed'))
    } finally {
      isDeleting.value = false
    }
  }

  const handleMessageClick = (message: Message) => {
    if (message.status === MessageStatus.UNREAD) {
      markRead(message.id)
    }
  }

  const _handleAction = (message: Message) => {
    if (message.action_url) {
      if (message.action_url.startsWith('http')) {
        window.open(message.action_url, '_blank')
      } else {
        router.push(message.action_url)
      }
    }
  }

  const handleFilterChange = () => {
    currentPage.value = 1
    loadMessages()
  }

  const goToPage = (page: number) => {
    currentPage.value = page
    loadMessages()
  }

  const handlePageSizeChange = (size: number) => {
    pageSize.value = size
    currentPage.value = 1
    loadMessages()
  }

  const getMessageIcon = (type: string): string => messageTypeConfig[type as keyof typeof messageTypeConfig]?.icon || '📧'

  const getMessageTypeLabel = (type: string): string => messageTypeConfig[type as keyof typeof messageTypeConfig]?.label || type

  const _getActionButtonType = (style?: string): string => {
    switch (style) {
      case 'primary':
        return 'primary'
      case 'success':
        return 'secondary'
      case 'warning':
        return 'outlined'
      case 'danger':
        return 'danger'
      default:
        return 'outlined'
    }
  }

  const _getPriorityBadgeType = (priority: number | null | undefined): string => {
    const actualPriority = priority ?? MessagePriority.NORMAL
    switch (actualPriority) {
      case MessagePriority.HIGH:
        return 'danger'
      case MessagePriority.NORMAL:
        return 'primary'
      case MessagePriority.LOW:
        return 'info'
      default:
        return 'primary'
    }
  }

  const getPriorityLabel = (priority: number | null | undefined): string => {
    const actualPriority = priority ?? MessagePriority.NORMAL
    return priorityConfig[actualPriority as keyof typeof priorityConfig]?.label || $t('message.messageList.priorityLabels.normal')
  }

  const getPriorityBadgeClass = (priority: number | null | undefined): string => {
    const actualPriority = priority ?? MessagePriority.NORMAL
    switch (actualPriority) {
      case MessagePriority.HIGH:
        return 'priority-badge priority-badge--high'
      case MessagePriority.NORMAL:
        return 'priority-badge priority-badge--normal'
      case MessagePriority.LOW:
        return 'priority-badge priority-badge--low'
      default:
        return 'priority-badge priority-badge--normal'
    }
  }

  watch(
    [filterStatus, filterType],
    () => {
      handleFilterChange()
    },
    { deep: true }
  )

  onMounted(() => {
    loadMessages()
    loadUnreadCount()
  })
</script>

<template>
  <div class="cyber-message-list">
    <div class="cyber-header">
      <div class="cyber-header-left">
        <h2 class="cyber-title">{{ $t('message.messageList.title') }}</h2>
        <div class="cyber-stats">
          <span class="cyber-stat-item">
            <i class="fas fa-envelope mr-1 text-brand-400" />
            {{ $t('message.messageList.total') }}: {{ totalMessages }}
          </span>
          <span class="cyber-stat-item">
            <i class="fas fa-circle mr-1 text-error-400" />
            {{ $t('message.messageList.unread') }}: {{ unreadCount }}
          </span>
        </div>
      </div>
      <div class="cyber-header-right">
        <div class="cyber-filters">
          <CyberDropdown
            :model-value="filterStatus"
            :placeholder="$t('message.messageList.allStatus')"
            :options="statusOptions"
            class="cyber-filter-dropdown"
            @update:model-value="
              (val) => {
                filterStatus = val
                handleFilterChange()
              }
            "
          />
          <CyberDropdown
            :model-value="filterType"
            :placeholder="$t('message.messageList.allTypes')"
            :options="typeOptions"
            class="cyber-filter-dropdown"
            @update:model-value="
              (val) => {
                filterType = val
                handleFilterChange()
              }
            "
          />
        </div>
        <CyberButton type="secondary" icon="check" :disabled="unreadCount === 0 || loading" @click="markAllRead">
          {{ $t('message.messageList.markAllRead') }}
        </CyberButton>
      </div>
    </div>

    <div class="cyber-content">
      <CyberSkeleton type="table" :count="8" :loading="showSkeleton" />

      <div v-if="!showSkeleton && !loading && (!messages || messages.length === 0)" class="cyber-empty">
        <i class="fas fa-inbox mb-4 text-6xl text-content-disabled" />
        <h3 class="cyber-empty-title">{{ $t('message.messageList.emptyTitle') }}</h3>
        <p class="cyber-empty-desc">{{ $t('message.messageList.emptyDescription') }}</p>
      </div>

      <CyberDataTable
        v-if="!showSkeleton && messages && messages.length > 0"
        :columns="tableColumns"
        :data="messages"
        :loading="false"
        :hoverable="true"
        :striped="true"
        :bordered="true"
        :empty-text="$t('message.messageList.emptyTitle')"
        :loading-text="$t('message.messageList.loadingText')"
        class="cyber-message-table"
        @row-click="handleMessageClick"
      />
    </div>

    <div v-if="totalPages > 1" class="cyber-pagination-wrapper">
      <CyberPagination
        v-model:current-page="currentPage"
        v-model:page-size="pageSize"
        :total="totalMessages"
        :show-page-size-selector="true"
        :show-quick-jumper="true"
        :show-total="true"
        @update:current-page="goToPage"
        @update:size="handlePageSizeChange"
      />
    </div>

    <CyberDialog
      id="delete-message-dialog"
      v-model="showDeleteDialog"
      :title="$t('message.messageList.dialog.deleteTitle')"
      width="450px"
      :show-default-footer="false"
      @cancel="closeDeleteDialog"
    >
      <p class="mb-6">
        {{ $t('message.messageList.dialog.deleteMessage') }}
      </p>

      <template #footer>
        <div class="flex justify-end gap-3 p-4">
          <CyberButton type="ghost" @click="closeDeleteDialog">
            {{ $t('message.messageList.dialog.cancel') }}
          </CyberButton>
          <CyberButton type="danger" icon="trash" :loading="isDeleting" @click="deleteMessage">
            {{ isDeleting ? $t('status.deleting') : $t('message.messageList.dialog.confirmDelete') }}
          </CyberButton>
        </div>
      </template>
    </CyberDialog>
  </div>
</template>

<style scoped>
  .cyber-message-list {
    @apply relative z-[1] w-full;
  }

  .cyber-header {
    @apply relative mb-6 flex flex-wrap items-center justify-between gap-4 rounded-xl px-5 py-4;
    background: linear-gradient(145deg, rgba(var(--color-background-800-rgb), 0.8), rgba(var(--color-background-800-rgb), 0.6));
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.2);
    box-shadow:
      0 2px 0 rgba(var(--color-brand-500-rgb), 0.15) inset,
      0 -1px 0 rgba(0, 0, 0, 0.15) inset,
      0 10px 20px rgba(0, 0, 0, 0.2),
      0 5px 10px rgba(var(--color-brand-500-rgb), 0.15),
      0 2px 5px rgba(0, 0, 0, 0.1);
    backdrop-filter: var(--backdrop-blur-md);
    overflow: hidden;
  }

  .cyber-header-left {
    @apply flex flex-wrap items-center gap-4;
  }

  .cyber-title {
    @apply m-0 text-xl font-semibold text-content-heading;
  }

  .cyber-stats {
    @apply flex flex-wrap gap-2 text-xs text-content-muted;
  }

  .cyber-stat-item {
    @apply inline-flex items-center gap-1 whitespace-nowrap rounded-lg border px-2.5 py-1.5 font-medium;
    background-color: rgba(var(--color-background-700-rgb), 0.72);
    border-color: var(--color-border-subtle);
  }

  .cyber-header-right {
    @apply flex flex-nowrap items-center gap-3;
  }

  .cyber-filters {
    @apply flex flex-nowrap gap-2 whitespace-nowrap;
  }

  .cyber-filter-dropdown {
    @apply min-w-[132px];
  }

  .cyber-content {
    @apply relative min-h-[400px] rounded-2xl p-6;
    background: linear-gradient(145deg, rgba(var(--color-background-800-rgb), 0.8), rgba(var(--color-background-800-rgb), 0.6));
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.2);
    box-shadow:
      0 2px 0 rgba(var(--color-brand-500-rgb), 0.15) inset,
      0 -1px 0 rgba(0, 0, 0, 0.15) inset,
      0 10px 20px rgba(0, 0, 0, 0.2),
      0 5px 10px rgba(var(--color-brand-500-rgb), 0.15),
      0 2px 5px rgba(0, 0, 0, 0.1);
    backdrop-filter: var(--backdrop-blur-md);
    overflow: hidden;
  }

  .cyber-empty {
    @apply px-8 py-16 text-center;
    color: var(--color-content);
  }

  .cyber-empty i {
    color: var(--color-brand-400);
    text-shadow: 0 0 16px rgba(var(--color-brand-500-rgb), 0.25);
  }

  .cyber-empty-title {
    @apply m-0 mb-2 text-xl font-semibold text-content-heading;
  }

  .cyber-empty-desc {
    @apply m-0 text-sm text-content-muted;
  }

  .cyber-message-table {
    @apply mt-0;
  }

  .cyber-pagination-wrapper {
    @apply mt-8 flex justify-center rounded-2xl py-6;
    background: linear-gradient(145deg, rgba(var(--color-background-800-rgb), 0.8), rgba(var(--color-background-800-rgb), 0.6));
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.2);
    box-shadow:
      0 2px 0 rgba(var(--color-brand-500-rgb), 0.15) inset,
      0 -1px 0 rgba(0, 0, 0, 0.15) inset,
      0 10px 20px rgba(0, 0, 0, 0.2),
      0 5px 10px rgba(var(--color-brand-500-rgb), 0.15),
      0 2px 5px rgba(0, 0, 0, 0.1);
    backdrop-filter: var(--backdrop-blur-md);
  }

  .priority-badge {
    @apply inline-flex items-center rounded border px-2 py-1 text-xs font-medium;
  }

  .priority-badge--high {
    background: linear-gradient(145deg, rgba(var(--color-error-rgb), 0.2), rgba(var(--color-error-rgb), 0.1));
    border-color: rgba(var(--color-error-rgb), 0.3);
    color: var(--color-error-500);
    box-shadow:
      0 1px 0 rgba(var(--color-error-rgb), 0.3) inset,
      0 1px 2px rgba(0, 0, 0, 0.15);
  }

  .priority-badge--normal {
    background: linear-gradient(145deg, rgba(var(--color-brand-500-rgb), 0.2), rgba(var(--color-brand-500-rgb), 0.1));
    border-color: rgba(var(--color-brand-500-rgb), 0.3);
    color: var(--color-brand-500);
    box-shadow:
      0 1px 0 rgba(var(--color-brand-500-rgb), 0.3) inset,
      0 1px 2px rgba(0, 0, 0, 0.15);
  }

  .priority-badge--low {
    background: linear-gradient(145deg, rgba(var(--color-content-rgb), 0.15), rgba(var(--color-content-rgb), 0.08));
    border-color: rgba(var(--color-content-rgb), 0.2);
    color: var(--color-content-muted);
    box-shadow:
      0 1px 0 rgba(255, 255, 255, 0.1) inset,
      0 1px 2px rgba(0, 0, 0, 0.15);
  }

  @media (max-width: 768px) {
    .cyber-header {
      @apply flex-col items-stretch gap-3 p-5;
    }

    .cyber-header-left {
      @apply flex-col items-center gap-4;
    }

    .cyber-header-right {
      @apply flex-col items-stretch gap-4;
    }

    .cyber-filters {
      @apply w-full flex-col;
    }

    .cyber-filter-dropdown {
      @apply w-full min-w-0;
    }

    .cyber-stats {
      @apply w-full flex-col gap-3;
    }

    .cyber-stat-item {
      @apply w-full justify-center;
    }

    .cyber-content {
      @apply p-5;
    }

    .cyber-empty {
      @apply px-4 py-8;
    }
  }

  @media (max-width: 480px) {
    .cyber-header {
      @apply mb-4 p-4;
    }

    .cyber-content {
      @apply p-4;
    }

    .cyber-title {
      @apply text-xl;
    }

    .cyber-stat-item {
      @apply px-3 py-2 text-xs;
    }
  }
</style>
