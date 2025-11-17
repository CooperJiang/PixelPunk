<script setup lang="ts">
  import { computed, onMounted, ref, watch } from 'vue'
  import { useRouter } from 'vue-router'
  import { formatRelativeTime, getUnreadCount, getUserMessages, markMessageRead } from '@/api/message'
  import type { Message } from '@/api/message/index'
  import { useToast } from '@/components/Toast/useToast'
  import { useMessageToast } from '@/components/Message/useMessageToast'
  import { useMessageModal } from '@/components/Message/useMessageModal'
  import { useTexts } from '@/composables/useTexts'
  import { formatMessage } from '@/utils/message/formatter'

  defineOptions({
    name: 'MessageCenter',
  })

  const props = defineProps<{
    refreshKey?: number
  }>()

  const router = useRouter()
  const toast = useToast()
  const { $t } = useTexts()
  const messageToast = useMessageToast()
  const messageModal = useMessageModal()

  const loading = ref(false)
  const messages = ref<Message[]>([])
  const currentPage = ref(1)
  const pageSize = ref(15)
  const hasMore = ref(true)
  const loadingMore = ref(false)
  const scrollContainer = ref<HTMLElement>()

  const unreadCount = computed(() => messages.value.filter((m) => m.status === 1).length)

  const sortMessages = (messageList: Message[], byStatus = false) => {
    return messageList.sort((a, b) => {
      if (byStatus && a.status !== b.status) {
        return a.status - b.status
      }
      return new Date(b.created_at).getTime() - new Date(a.created_at).getTime()
    })
  }

  const reorganizeMessages = () => {
    const unread = messages.value.filter((m) => m.status === 1)
    const read = messages.value.filter((m) => m.status === 2)
    sortMessages(read)
    messages.value = [...unread, ...read]
  }

  const fetchMessages = async (reset = true) => {
    if (loading.value || loadingMore.value) return

    try {
      if (reset) {
        loading.value = true
      } else {
        loadingMore.value = true
      }
      if (reset) {
        currentPage.value = 1
        hasMore.value = true
      }

      const [messagesResult] = await Promise.all([
        getUserMessages({ page: currentPage.value, pageSize: pageSize.value }),
        reset ? getUnreadCount() : Promise.resolve(null),
      ])

      if (messagesResult.success) {
        const messageList = messagesResult.data?.items || []

        if (reset) {
          messages.value = sortMessages(messageList, true)
        } else {
          const existingIds = new Set(messages.value.map((m) => m.id))
          const newMessages = messageList.filter((m) => !existingIds.has(m.id))
          const unread = messages.value.filter((m) => m.status === 1)
          const read = [...messages.value.filter((m) => m.status === 2), ...newMessages.filter((m) => m.status === 2)]
          const newUnread = newMessages.filter((m) => m.status === 1)
          sortMessages(read)
          messages.value = [...unread, ...newUnread, ...read]
        }

        hasMore.value = messageList.length === pageSize.value
      }
    } catch (_error: unknown) {
    } finally {
      loading.value = false
      loadingMore.value = false
    }
  }

  const loadMoreMessages = async () => {
    if (!hasMore.value || loadingMore.value) return
    currentPage.value++
    await fetchMessages(false)
  }

  const handleScroll = async (event: Event) => {
    const target = event.target as HTMLElement
    if (!target || !hasMore.value || loadingMore.value) return

    const { scrollTop, scrollHeight, clientHeight } = target
    if (scrollTop + clientHeight >= scrollHeight - 50) {
      await loadMoreMessages()
    }
  }

  const isUnread = (message: Message): boolean => message.status === 1

  const MESSAGE_TYPE_MAP: Record<string, string> = {
    'system.maintenance': 'system',
    'system.announcement': 'system',
    'account.register': 'account',
    'account.storage_granted': 'account',
    'account.bandwidth_granted': 'account',
    'content.review_approved': 'content',
    'content.review_rejected': 'content',
    'content.review_pending': 'content',
    'storage.quota_warning': 'storage',
    'storage.quota_increased': 'storage',
    'storage.quota_decreased': 'storage',
    'file.deleted_by_admin': 'file',
    'file.batch_deleted_by_admin': 'file',
    'file.hard_deleted_by_admin': 'file',
    'file.expiry_warning': 'file',
    'file.thumbnail_failed': 'file',
    'security.login_alert': 'security',
    'apikey.created': 'apikey',
    'apikey.deleted': 'apikey',
    'apikey.regenerated': 'apikey',
    'apikey.disabled': 'apikey',
    'apikey.enabled': 'apikey',
    'random_api.created': 'random_api',
    'random_api.deleted': 'random_api',
    'random_api.disabled': 'random_api',
    'random_api.enabled': 'random_api',
    'share.expiry_warning': 'share',
  }

  const MESSAGE_ICONS: Record<string, string> = {
    'system.maintenance': 'fas fa-tools',
    'system.announcement': 'fas fa-bullhorn',
    'account.register': 'fas fa-user-plus',
    'account.storage_granted': 'fas fa-gift',
    'account.bandwidth_granted': 'fas fa-rocket',
    'content.review_approved': 'fas fa-check-circle',
    'content.review_rejected': 'fas fa-times-circle',
    'content.review_pending': 'fas fa-clock',
    'storage.quota_warning': 'fas fa-exclamation-triangle',
    'storage.quota_increased': 'fas fa-arrow-up',
    'storage.quota_decreased': 'fas fa-arrow-down',
    'file.deleted_by_admin': 'fas fa-trash',
    'file.batch_deleted_by_admin': 'fas fa-trash-alt',
    'file.hard_deleted_by_admin': 'fas fa-skull-crossbones',
    'file.expiry_warning': 'fas fa-hourglass-half',
    'file.thumbnail_failed': 'fas fa-image',
    'security.login_alert': 'fas fa-shield-alt',
    'apikey.created': 'fas fa-key',
    'apikey.deleted': 'fas fa-key',
    'apikey.regenerated': 'fas fa-sync-alt',
    'apikey.disabled': 'fas fa-lock',
    'apikey.enabled': 'fas fa-unlock',
    'random_api.created': 'fas fa-random',
    'random_api.deleted': 'fas fa-random',
    'random_api.disabled': 'fas fa-ban',
    'random_api.enabled': 'fas fa-check-circle',
    'share.expiry_warning': 'fas fa-share-alt',
  }

  const getMessageTypeClass = (type: string): string => MESSAGE_TYPE_MAP[type] || 'default'
  const getMessageIcon = (type: string): string => MESSAGE_ICONS[type] || 'fas fa-bell'

  const handleMessageClick = async (message: Message) => {
    if (isUnread(message)) {
      try {
        const messageIndex = messages.value.findIndex((m) => m.id === message.id)
        if (messageIndex !== -1) {
          messages.value[messageIndex].status = 2
          messages.value[messageIndex].read_at = new Date().toISOString()
          reorganizeMessages()
        }

        markMessageRead(message.id.toString()).catch(() => {
          if (messageIndex !== -1) {
            messages.value[messageIndex].status = 1
            messages.value[messageIndex].read_at = undefined
            reorganizeMessages()
          }
          toast.error($t('dashboard.messages.markReadFailed'))
        })
      } catch (_error) {}
    }

    messageModal.showMessageDetail(message)

    // 使用 formatMessage 生成标题和内容
    const formatted = formatMessage(message, $t)

    messageToast.showMessage({
      title: formatted.title,
      content: formatted.content,
      type: message.type,
      action_url: message.action_url,
      action_text: message.action_text,
      action_style: message.action_style,
    })
  }

  const handleViewAllClick = () => router.push('/messages')

  onMounted(() => fetchMessages())

  watch(
    () => props.refreshKey,
    () => {
      if (props.refreshKey !== undefined) fetchMessages()
    }
  )
</script>

<template>
  <div class="message-center cyber-card">
    <div class="center-content">
      <div class="center-header">
        <h3 class="section-title">
          <i class="fas fa-bell" />
          {{ $t('dashboard.messages.title') }}
        </h3>
        <div class="message-stats">
          <span v-if="unreadCount > 0" class="unread-count"> {{ unreadCount }} {{ $t('dashboard.messages.unread') }} </span>
          <button class="view-all-btn" @click="handleViewAllClick">
            <i class="fas fa-arrow-right" />
            {{ $t('dashboard.messages.viewAll') }}
          </button>
        </div>
      </div>

      <div v-if="loading && messages.length === 0" class="loading-state">
        <div class="loading-spinner" />
        <span>{{ $t('dashboard.messages.loading') }}</span>
      </div>

      <div v-else-if="!loading && messages.length === 0" class="empty-state">
        <i class="fas fa-inbox" />
        <span>{{ $t('dashboard.messages.empty') }}</span>
      </div>

      <div v-else class="messages-wrapper">
        <div ref="scrollContainer" class="messages-container" @scroll="handleScroll">
          <div class="messages-list">
            <div
              v-for="message in messages"
              :key="message.id"
              class="message-item"
              :class="{ unread: isUnread(message) }"
              @click="handleMessageClick(message)"
            >
              <div class="message-icon" :class="getMessageTypeClass(message.type)">
                <i :class="getMessageIcon(message.type)" />
              </div>

              <div class="message-content">
                <div class="message-header">
                  <div class="message-title">{{ formatMessage(message, $t).title }}</div>
                  <div class="message-time">{{ formatRelativeTime(message.created_at, $t) }}</div>
                </div>
                <div class="message-body">{{ formatMessage(message, $t).content }}</div>
              </div>

              <div v-if="message.priority === 1" class="message-priority">
                <i class="fas fa-exclamation-triangle" />
              </div>

              <div v-if="isUnread(message)" class="unread-dot" />
            </div>

            <div v-if="loadingMore" class="message-item loading-more">
              <div class="message-icon loading">
                <i class="fas fa-spinner fa-spin" />
              </div>
              <div class="message-content">
                <div class="message-title">{{ $t('dashboard.messages.loadingMore') }}</div>
                <div class="message-body">{{ $t('dashboard.messages.loadingMoreDesc') }}</div>
              </div>
            </div>

            <div v-if="!hasMore && messages.length > 0" class="message-item no-more">
              <div class="message-icon default">
                <i class="fas fa-check-circle" />
              </div>
              <div class="message-content">
                <div class="message-title">{{ $t('dashboard.messages.noMore') }}</div>
                <div class="message-body">{{ $t('dashboard.messages.noMoreDesc', { count: messages.length }) }}</div>
              </div>
            </div>
          </div>
        </div>

        <div v-if="messages.length > 3" class="scroll-indicator">
          <i class="fas fa-chevron-down" />
          <span>{{ $t('dashboard.messages.scrollHint') }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
  .message-center {
    padding: 0;
    width: 100%;
    min-width: 0;
    box-sizing: border-box;
  }

  .center-content {
    display: flex;
    flex-direction: column;
    gap: 20px;
    width: 100%;
    min-width: 0;
    overflow: hidden;
  }

  .center-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 20px 20px 0;
  }

  .section-title {
    margin: 0;
    font-size: 16px;
    font-weight: 600;
    color: var(--color-content-heading);
    display: flex;
    align-items: center;
    gap: 10px;
  }

  .section-title i {
    color: rgba(var(--color-brand-500-rgb), 0.8);
    font-size: 14px;
  }

  .message-stats {
    display: flex;
    align-items: center;
    gap: 12px;
    font-size: 12px;
  }

  .unread-count {
    background: rgba(var(--color-brand-500-rgb), 0.15);
    color: rgba(var(--color-brand-500-rgb), 0.9);
    padding: 3px 8px;
    border-radius: var(--radius-sm);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.3);
    font-weight: 600;
    font-size: 10px;
    text-transform: uppercase;
    letter-spacing: 0.5px;
  }

  .view-all-btn {
    display: flex;
    align-items: center;
    gap: 6px;
    color: var(--color-content-muted);
    font-size: 11px;
    font-weight: 500;
    padding: 4px 8px;
    border-radius: var(--radius-sm);
    border: 1px solid var(--color-border-subtle);
    transition: all 0.3s ease;
  }

  .view-all-btn:hover {
    color: var(--color-brand-500);
    background: var(--color-hover-bg);
    border-color: var(--color-hover-border);
    transform: translateX(2px);
  }

  .view-all-btn i {
    font-size: 9px;
  }

  .loading-state,
  .empty-state {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 40px 20px 20px;
    gap: 12px;
    color: rgba(var(--color-content-rgb), 0.6);
    font-size: 13px;
  }

  .loading-spinner {
    width: 20px;
    height: 20px;
    border: 2px solid rgba(var(--color-brand-500-rgb), 0.2);
    border-top: 2px solid rgba(var(--color-brand-500-rgb), 0.8);
    border-radius: var(--radius-full);
    animation: spin 1s linear infinite;
  }

  @keyframes spin {
    0% {
      transform: rotate(0deg);
    }
    100% {
      transform: rotate(360deg);
    }
  }

  .empty-state i {
    font-size: 24px;
    color: rgba(var(--color-brand-500-rgb), 0.4);
  }

  .messages-wrapper {
    position: relative;
    padding: 0 20px;
  }

  .messages-container {
    max-height: 400px;
    overflow-y: auto;
    overflow-x: hidden;
    padding-bottom: 50px;
  }

  .messages-list {
    display: flex;
    flex-direction: column;
    gap: 8px;
  }

  .messages-container::-webkit-scrollbar {
    width: 6px;
  }

  .messages-container::-webkit-scrollbar-track {
    background: var(--color-background-900);
    border-radius: var(--radius-sm);
    margin: 4px 0;
  }

  .messages-container::-webkit-scrollbar-thumb {
    background: var(--color-brand-500);
    border-radius: var(--radius-sm);
    border: 1px solid var(--color-border-default);
    transition: all 0.3s ease;
  }

  .messages-container::-webkit-scrollbar-thumb:hover {
    background: var(--color-brand-600);
    border-color: var(--color-hover-border);
    box-shadow: var(--shadow-glow-primary-sm);
  }

  .scroll-indicator {
    position: absolute;
    bottom: 0;
    left: 0;
    right: 0;
    height: 40px;
    background: linear-gradient(
      transparent 0%,
      rgba(var(--color-background-800-rgb), 0.8) 30%,
      rgba(var(--color-background-800-rgb), 0.95) 100%
    );
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 6px;
    font-size: 10px;
    color: var(--color-content-muted);
    pointer-events: none;
    z-index: 2;
    backdrop-filter: blur(4px);
  }

  .scroll-indicator i {
    font-size: 8px;
    animation: bounce-down 2.5s ease-in-out infinite;
  }

  @keyframes bounce-down {
    0%,
    100% {
      transform: translateY(0);
      opacity: 0.5;
    }
    50% {
      transform: translateY(3px);
      opacity: 0.9;
    }
  }

  .message-item {
    display: flex;
    align-items: center;
    gap: 12px;
    background: var(--color-background-800);
    border: 1px solid var(--color-border-subtle);
    border-radius: var(--radius-sm);
    padding: 10px 12px;
    min-height: 48px;
    width: 100%;
    min-width: 0;
    cursor: pointer;
    position: relative;
    overflow: hidden;
    flex-shrink: 0;
    box-sizing: border-box;
    transition: all 0.3s ease;
  }

  .message-item.unread {
    border-left: 3px solid rgba(var(--color-brand-500-rgb), 0.8);
    background: rgba(var(--color-brand-500-rgb), 0.05);
  }

  .message-item:hover {
    transform: translateY(-1px);
    background: var(--color-hover-bg);
    border-top-color: transparent;
    border-right-color: transparent;
    border-bottom-color: transparent;
    box-shadow: var(--shadow-cyber-sm);
  }

  .message-item:not(.unread):hover {
    border-left-color: transparent;
  }

  .message-icon {
    width: 28px;
    height: 28px;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: var(--radius-sm);
    font-size: 12px;
    flex-shrink: 0;
    border: 1px solid;
  }

  .message-icon.system {
    background: rgba(var(--color-info-rgb), 0.12);
    color: rgba(var(--color-info-rgb), 0.9);
    border-color: rgba(var(--color-info-rgb), 0.25);
  }

  .message-icon.account {
    background: rgba(var(--color-brand-500-rgb), 0.12);
    color: rgba(var(--color-brand-500-rgb), 0.9);
    border-color: rgba(var(--color-brand-500-rgb), 0.25);
  }

  .message-icon.content {
    background: rgba(var(--color-success-rgb), 0.12);
    color: rgba(var(--color-success-rgb), 0.9);
    border-color: rgba(var(--color-success-rgb), 0.25);
  }

  .message-icon.storage {
    background: rgba(var(--color-brand-500-rgb), 0.12);
    color: rgba(var(--color-brand-500-rgb), 0.9);
    border-color: rgba(var(--color-brand-500-rgb), 0.25);
  }

  .message-icon.file {
    background: rgba(var(--color-error-rgb), 0.12);
    color: rgba(var(--color-error-rgb), 0.9);
    border-color: rgba(var(--color-error-rgb), 0.25);
  }

  .message-icon.security {
    background: rgba(var(--color-warning-rgb), 0.12);
    color: rgba(var(--color-warning-rgb), 0.9);
    border-color: rgba(var(--color-warning-rgb), 0.25);
  }

  .message-icon.apikey {
    background: rgba(103, 194, 255, 0.12);
    color: rgba(103, 194, 255, 0.9);
    border-color: rgba(103, 194, 255, 0.25);
  }

  .message-icon.random_api {
    background: rgba(156, 39, 176, 0.12);
    color: rgba(156, 39, 176, 0.9);
    border-color: rgba(156, 39, 176, 0.25);
  }

  .message-icon.default,
  .message-icon.loading {
    background: rgba(var(--color-content-rgb), 0.12);
    color: rgba(var(--color-content-rgb), 0.8);
    border-color: rgba(var(--color-content-rgb), 0.25);
  }

  .message-content {
    flex: 1;
    min-width: 0;
    overflow: hidden;
  }

  .message-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 2px;
    gap: 6px;
    width: 100%;
    min-width: 0;
  }

  .message-title {
    font-size: 12px;
    font-weight: 600;
    color: rgba(var(--color-content-rgb), 0.9);
    line-height: 1.2;
    flex: 1;
    min-width: 0;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .message-time {
    font-size: 9px;
    color: rgba(var(--color-content-rgb), 0.4);
    font-weight: 500;
    flex-shrink: 0;
  }

  .message-body {
    font-size: 11px;
    color: rgba(var(--color-content-rgb), 0.6);
    line-height: 1.2;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    width: 100%;
  }

  .message-priority {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 16px;
    height: 16px;
    color: rgba(var(--color-warning-rgb), 0.9);
    font-size: 9px;
  }

  .unread-dot {
    width: 6px;
    height: 6px;
    border-radius: var(--radius-full);
    background: rgba(var(--color-brand-500-rgb), 0.9);
    flex-shrink: 0;
    box-shadow: 0 0 6px rgba(var(--color-brand-500-rgb), 0.4);
    animation: pulse-glow 2s ease-in-out infinite;
  }

  @keyframes pulse-glow {
    0%,
    100% {
      box-shadow: 0 0 6px rgba(var(--color-brand-500-rgb), 0.4);
      transform: scale(1);
    }
    50% {
      box-shadow: 0 0 10px rgba(var(--color-brand-500-rgb), 0.7);
      transform: scale(1.1);
    }
  }

  .message-item.loading-more {
    cursor: default;
    opacity: 0.8;
    background: rgba(var(--color-info-rgb), 0.08);
    border-color: rgba(var(--color-info-rgb), 0.2);
  }

  .message-item.no-more {
    cursor: default;
    opacity: 0.7;
    background: rgba(var(--color-success-rgb), 0.08);
    border-color: rgba(var(--color-success-rgb), 0.2);
  }

  @media (max-width: 768px) {
    .center-header {
      padding: 16px 16px 0;
      flex-direction: column;
      gap: 12px;
      align-items: stretch;
    }

    .loading-state,
    .empty-state {
      padding: 40px 16px 16px;
    }

    .message-stats {
      justify-content: space-between;
    }

    .messages-wrapper {
      padding: 0 16px;
    }

    .messages-container {
      max-height: 350px;
    }

    .message-item {
      padding: 10px;
    }

    .message-icon {
      width: 28px;
      height: 28px;
      font-size: 12px;
    }
  }

  @media (max-width: 480px) {
    .center-header {
      padding: 12px 12px 0;
    }

    .loading-state,
    .empty-state {
      padding: 40px 12px 12px;
    }

    .message-stats {
      flex-direction: column;
      gap: 8px;
    }

    .view-all-btn {
      align-self: flex-end;
    }

    .messages-wrapper {
      padding: 0 12px;
    }

    .messages-container {
      max-height: 300px;
    }

    .scroll-indicator {
      font-size: 9px;
    }
  }
</style>
