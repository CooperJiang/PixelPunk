<script setup lang="ts">
  import { computed, nextTick, onMounted, onUnmounted, ref } from 'vue'
  import { useRouter } from 'vue-router'
  import { getUnreadCount, getUserMessages, markAllMessagesRead, markMessageRead, formatRelativeTime } from '@/api/message'
  import { MessageStatus, getMessageTypeConfig, type Message } from '@/api/message/types'
  import { useMessageToast } from './useMessageToast'
  import { useMessageModal } from './useMessageModal'
  import { useTexts } from '@/composables/useTexts'
  import { formatMessage } from '@/utils/message/formatter'

  const props = withDefaults(defineProps<{ offsetX?: number; offsetY?: number }>(), {
    offsetX: 12,
    offsetY: 8,
  })

  const router = useRouter()
  const messageToast = useMessageToast()
  const messageModal = useMessageModal()
  const { $t } = useTexts()

  const messageTypeConfig = getMessageTypeConfig($t)

  const showDropdown = ref(false)
  const loading = ref(false)
  const unreadCount = ref(0)
  const recentMessages = ref<Message[]>([])
  const dropdownRef = ref<HTMLElement>()
  const messageButtonRef = ref<HTMLElement>()
  const dropdownStyle = ref<Record<string, string>>({})

  let pollTimer: NodeJS.Timeout | null = null
  let positionTimer: number | null = null

  const hasUnreadMessages = computed(() => unreadCount.value > 0)

  const toggleDropdown = async () => {
    showDropdown.value = !showDropdown.value
    if (showDropdown.value) {
      loadRecentMessages()
      await nextTick()
      updateDropdownPosition()
    }
  }

  const updateDropdownPosition = () => {
    const iconButtonComponent = messageButtonRef.value
    const button = iconButtonComponent?.$el || iconButtonComponent
    if (!button) {
      return
    }
    const rect = button.getBoundingClientRect()
    const vw = window.innerWidth
    const scrollX = window.pageXOffset || document.documentElement.scrollLeft || 0
    const scrollY = window.pageYOffset || document.documentElement.scrollTop || 0
    const dropdownWidth = 380

    let left = rect.left + scrollX + rect.width / 2 - dropdownWidth / 2 + (props.offsetX || 0)
    if (left < 8) {
      left = 8
    }
    if (left + dropdownWidth > scrollX + vw - 8) {
      left = scrollX + vw - dropdownWidth - 8
    }
    const top = rect.bottom + scrollY + (props.offsetY || 0)

    const arrowLeft = rect.left + scrollX + rect.width / 2 - left - 5

    dropdownStyle.value = {
      left: `${left}px`,
      top: `${top}px`,
      width: `${dropdownWidth}px`,
      '--arrow-left': `${Math.max(12, Math.min(dropdownWidth - 12, arrowLeft))}px`,
    } as Record<string, string>
  }

  const handleViewportChange = () => {
    if (!showDropdown.value) {
      return
    }
    if (positionTimer) {
      window.clearTimeout(positionTimer)
    }
    positionTimer = window.setTimeout(() => updateDropdownPosition(), 120)
  }

  const closeDropdown = () => {
    showDropdown.value = false
  }

  const loadUnreadCount = async () => {
    try {
      const response: Event = await getUnreadCount()
      unreadCount.value = response?.data?.count ?? 0
    } catch {
      unreadCount.value = 0
    }
  }

  const loadRecentMessages = async () => {
    loading.value = true
    try {
      const response: Event = await getUserMessages({
        page: 1,
        page_size: 50,
        status: MessageStatus.UNREAD,
      })
      recentMessages.value = response?.data?.items || []
    } finally {
      loading.value = false
    }
  }

  const markRead = async (messageId: number) => {
    try {
      await markMessageRead(messageId.toString())
      recentMessages.value = recentMessages.value.filter((m) => m.id !== messageId)
      unreadCount.value = Math.max(0, unreadCount.value - 1)
    } catch {}
  }

  const markAllRead = async () => {
    try {
      await markAllMessagesRead()
      recentMessages.value = []
      unreadCount.value = 0
    } catch {}
  }

  // 格式化消息显示
  const getFormattedMessage = (message: Message) => {
    return formatMessage(message, $t)
  }

  const handleMessageClick = (message: Message) => {
    if (message.status === MessageStatus.UNREAD) {
      markRead(message.id)
    }
    closeDropdown()
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

  const goToMessageCenter = () => {
    closeDropdown()
    router.push('/messages')
  }
  const getMessageIcon = (type: string): string => {
    const iconMap: Record<string, string> = {
      'system.maintenance': 'wrench',
      'system.announcement': 'bullhorn',
      'account.register': 'party-horn',
      'account.storage_granted': 'hard-drive',
      'account.bandwidth_granted': 'rocket',
      'content.review_approved': 'check-circle',
      'content.review_rejected': 'times-circle',
      'content.review_pending': 'hourglass-half',
      'storage.quota_warning': 'exclamation-triangle',
      'storage.quota_increased': 'chart-line',
      'storage.quota_decreased': 'arrow-down',
      'file.deleted_by_admin': 'trash-alt',
      'file.batch_deleted_by_admin': 'folder-minus',
      'file.hard_deleted_by_admin': 'exclamation-triangle',
      'security.login_alert': 'shield-alt',
      'file.expiry_warning': 'clock',
    }
    return iconMap[type] || 'envelope'
  }
  const getMessageTypeLabel = (type: string): string => messageTypeConfig[type as keyof typeof messageTypeConfig]?.label || type
  const getMessageColor = (type: string): string => {
    const colorMap: Record<string, string> = {
      'system.maintenance': 'warning',
      'system.announcement': 'info',
      'account.register': 'success',
      'account.storage_granted': 'brand',
      'account.bandwidth_granted': 'info',
      'content.review_approved': 'success',
      'content.review_rejected': 'error',
      'content.review_pending': 'warning',
      'storage.quota_warning': 'warning',
      'storage.quota_increased': 'success',
      'storage.quota_decreased': 'warning',
      'file.deleted_by_admin': 'error',
      'file.batch_deleted_by_admin': 'error',
      'file.hard_deleted_by_admin': 'error',
      'security.login_alert': 'error',
      'file.expiry_warning': 'warning',
    }
    return colorMap[type] || 'brand'
  }

  const handleKeydown = (e: KeyboardEvent) => {
    if (e.key === 'Escape') {
      closeDropdown()
    }
  }

  const handleDocumentClick = (e: MouseEvent) => {
    if (!showDropdown.value) {
      return
    }
    const target = e.target as Node
    if (dropdownRef.value?.contains(target)) {
      return
    }
    const iconButtonComponent = messageButtonRef.value
    const button = iconButtonComponent?.$el || iconButtonComponent
    if (button?.contains(target)) {
      return
    }
    closeDropdown()
  }

  onMounted(() => {
    loadUnreadCount()
    pollTimer = setInterval(loadUnreadCount, 30000)
    window.addEventListener('scroll', handleViewportChange, true)
    window.addEventListener('resize', handleViewportChange)
    document.addEventListener('mousedown', handleDocumentClick, true)
    document.addEventListener('keydown', handleKeydown)
  })

  onUnmounted(() => {
    if (pollTimer) {
      clearInterval(pollTimer)
    }
    window.removeEventListener('scroll', handleViewportChange, true)
    window.removeEventListener('resize', handleViewportChange)
    document.removeEventListener('mousedown', handleDocumentClick, true)
    document.removeEventListener('keydown', handleKeydown)
  })
</script>

<template>
  <div class="message-notification relative">
    <button
      ref="messageButtonRef"
      class="message-button"
      :class="{ 'has-unread': hasUnreadMessages }"
      aria-haspopup="listbox"
      :aria-expanded="String(showDropdown)"
      @click="toggleDropdown"
    >
      <i class="fas fa-bell message-icon" />

      <span v-if="unreadCount > 0" class="message-badge">
        {{ unreadCount > 99 ? '99+' : unreadCount }}
      </span>
    </button>

    <Teleport to="body">
      <Transition
        enter-active-class="transition ease-out duration-300"
        enter-from-class="transform opacity-0 scale-95 translate-y-4"
        enter-to-class="transform opacity-100 scale-100 translate-y-0"
        leave-active-class="transition ease-in duration-200"
        leave-from-class="transform opacity-100 scale-100 translate-y-0"
        leave-to-class="transform opacity-0 scale-95 translate-y-4"
      >
        <div
          v-if="showDropdown"
          ref="dropdownRef"
          class="cyber-dropdown z-[9999]"
          :style="dropdownStyle"
          role="listbox"
          @click.stop
        >
          <div class="cyber-dropdown-header px-4 py-3">
            <div class="flex items-center justify-between">
              <h3 id="notify-title" class="dropdown-title">{{ $t('components.messageNotification.title') }}</h3>
              <div class="flex items-center space-x-2">
                <span class="text-xs text-content-muted"
                  >{{ unreadCount }} {{ $t('components.messageNotification.unreadCount') }}</span
                >
                <button
                  v-if="unreadCount > 0"
                  class="cyber-mark-all-btn rounded px-2 py-1 text-xs font-medium transition-all duration-200"
                  @click="markAllRead"
                >
                  {{ $t('components.messageNotification.markAllRead') }}
                </button>
              </div>
            </div>
          </div>

          <div class="cyber-dropdown-body max-h-96 overflow-y-auto" aria-labelledby="notify-title">
            <div v-if="loading" class="px-4 py-8 text-center">
              <div class="cyber-spinner mx-auto mb-3" />
              <p class="text-sm text-content-muted">{{ $t('components.messageNotification.loading') }}</p>
            </div>

            <div
              v-else-if="recentMessages.length === 0"
              class="empty-state flex min-h-[120px] flex-col items-center justify-center px-4 py-8 text-center"
            >
              <div class="empty-icon mb-3">
                <i class="fas fa-bell-slash text-3xl" />
              </div>
              <h3 class="empty-title mb-2">{{ $t('components.messageNotification.emptyTitle') }}</h3>
              <p class="empty-description">{{ $t('components.messageNotification.emptyDescription') }}</p>
            </div>

            <div v-else class="py-1">
              <div
                v-for="message in recentMessages"
                :key="message.id"
                class="cyber-message-item cursor-pointer transition-all duration-300"
                :class="{ 'is-unread': message.status === MessageStatus.UNREAD }"
                @click="handleMessageClick(message)"
              >
                <div class="flex items-start space-x-3">
                  <div class="flex-shrink-0">
                    <div class="cyber-message-type-icon" :data-color="getMessageColor(message.type)">
                      <i :class="['fas', `fa-${getMessageIcon(message.type)}`]"></i>
                    </div>
                  </div>

                  <div class="min-w-0 flex-1">
                    <div class="mb-1 flex items-center justify-between">
                      <div class="flex min-w-0 flex-1 items-center space-x-2">
                        <h4 class="message-title truncate text-sm font-medium">
                          {{ getFormattedMessage(message).title }}
                        </h4>
                        <div v-if="message.status === MessageStatus.UNREAD" class="cyber-unread-dot flex-shrink-0" />
                      </div>
                      <span class="ml-3 flex-shrink-0 text-xs text-content-disabled">{{
                        formatRelativeTime(message.created_at, $t)
                      }}</span>
                    </div>

                    <p class="truncate text-xs text-content-muted">{{ getFormattedMessage(message).content }}</p>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <div class="cyber-dropdown-footer px-4 py-1">
            <div
              class="cyber-view-all-link group flex cursor-pointer items-center justify-center py-1 transition-all duration-300"
              @click="goToMessageCenter"
            >
              <span class="text-xs font-medium">{{ $t('components.messageNotification.viewAll') }}</span>
              <i class="fas fa-arrow-right ml-1 transform text-xs transition-transform duration-300 group-hover:translate-x-1" />
            </div>
          </div>
        </div>
      </Transition>
    </Teleport>
  </div>
</template>

<style scoped>
  .message-notification {
    flex-shrink: 0;
    position: relative;
    display: inline-block;
  }

  .message-button {
    position: relative;
    display: inline-flex;
    align-items: center;
    justify-content: center;
    width: var(--size-button-sm);
    height: var(--size-button-sm);
    background: rgba(var(--color-background-700-rgb), 0.1);
    border: 1px solid var(--color-border-default);
    border-radius: var(--radius-sm);
    color: var(--color-content-muted);
    cursor: pointer;
    transition: all var(--transition-normal) var(--ease-in-out);
    backdrop-filter: var(--backdrop-blur-md) saturate(120%);
    box-shadow: var(--shadow-sm);
    outline: none;
  }

  .message-button::after {
    content: '';
    position: absolute;
    inset: -1px;
    border-radius: inherit;
    background: linear-gradient(45deg, transparent, rgba(var(--color-brand-500-rgb), 0.1), transparent);
    opacity: 0;
    transition: opacity var(--transition-normal) ease;
    z-index: -1;
  }

  .message-button:hover {
    background: rgba(var(--color-brand-500-rgb), 0.1);
    border-color: var(--color-hover-border);
    color: var(--color-brand-500);
    transform: translateY(-2px);
    box-shadow:
      var(--shadow-lg),
      0 0 20px rgba(var(--color-brand-500-rgb), 0.15);
  }

  .message-button:hover::after {
    opacity: 1;
  }

  .message-button:active {
    transform: translateY(1px) scale(0.98);
  }

  .message-button:focus-visible {
    outline: none;
    box-shadow:
      0 0 0 2px var(--color-focus-ring),
      0 0 25px rgba(var(--color-brand-500-rgb), 0.2);
  }

  .message-button.has-unread {
    background: rgba(var(--color-brand-500-rgb), 0.15);
    border-color: var(--color-active-border);
    color: var(--color-brand-500);
    box-shadow: 0 4px 12px rgba(var(--color-brand-500-rgb), 0.15);
  }

  .message-button.has-unread:hover {
    background: rgba(var(--color-brand-500-rgb), 0.2);
    border-color: var(--color-hover-border);
    box-shadow:
      0 8px 25px rgba(var(--color-brand-500-rgb), 0.25),
      0 0 20px rgba(var(--color-brand-500-rgb), 0.15);
  }

  .message-icon {
    font-size: var(--size-icon-sm);
    line-height: 1;
  }

  .message-badge {
    position: absolute;
    top: -8px;
    right: -8px;
    display: flex;
    align-items: center;
    justify-content: center;
    min-width: 20px;
    height: 20px;
    padding: 0 var(--space-xs);
    background: linear-gradient(135deg, var(--color-error-500), var(--color-error-600));
    color: var(--color-text-on-error);
    border-radius: var(--radius-full);
    font-size: 11px;
    font-weight: var(--font-semibold);
    line-height: 1;
    box-shadow:
      0 4px 12px rgba(var(--color-error-rgb), 0.4),
      var(--shadow-sm),
      0 0 0 2px var(--color-background-900);
    z-index: 10;
    border: 1px solid rgba(var(--color-error-rgb), 0.5);
    backdrop-filter: var(--backdrop-blur-md) saturate(120%);
    transition: all var(--transition-normal) var(--ease-in-out);
  }

  .message-badge::after {
    content: '';
    position: absolute;
    inset: 0;
    border-radius: inherit;
    background: linear-gradient(135deg, rgba(255, 255, 255, 0.15), rgba(255, 255, 255, 0.05));
    opacity: 0.8;
  }

  .cyber-dropdown {
    background: linear-gradient(145deg, rgba(var(--color-background-700-rgb), 0.98), rgba(var(--color-background-800-rgb), 0.96));
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.3);
    border-radius: var(--radius-sm);
    backdrop-filter: var(--backdrop-blur-xl) saturate(180%);
    box-shadow:
      var(--shadow-xl),
      0 0 30px rgba(var(--color-brand-500-rgb), 0.15),
      inset 0 1px 0 rgba(255, 255, 255, 0.08);
    overflow: hidden;
    position: absolute;
  }

  .cyber-dropdown-header {
    padding: var(--space-md);
    background: linear-gradient(135deg, rgba(var(--color-brand-500-rgb), 0.08), rgba(var(--color-background-700-rgb), 0.8));
    border-bottom: 1px solid var(--color-border-default);
  }

  .dropdown-title {
    font-size: var(--text-sm);
    font-weight: var(--font-medium);
    color: var(--color-brand-500);
  }

  .cyber-mark-all-btn {
    background: var(--color-hover-bg);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.3);
    color: var(--color-brand-500);
    border-radius: var(--radius-md);
    padding: var(--space-xs) var(--space-sm);
    font-size: var(--text-xs);
    font-weight: var(--font-medium);
    transition: all var(--transition-normal) ease;
  }

  .cyber-mark-all-btn:hover {
    background: var(--color-active-bg);
    border-color: var(--color-hover-border);
    box-shadow: 0 0 15px rgba(var(--color-brand-500-rgb), 0.3);
  }

  .cyber-dropdown-body {
    background: linear-gradient(
      180deg,
      rgba(var(--color-brand-500-rgb), 0.02) 0%,
      transparent 10%,
      transparent 90%,
      rgba(var(--color-brand-500-rgb), 0.02) 100%
    );
    padding: var(--space-sm);
  }

  .cyber-message-item {
    position: relative;
    border-radius: var(--radius-sm);
    margin: 0 var(--space-sm) var(--space-sm) var(--space-sm);
    padding: var(--space-sm);
    transition: all var(--transition-normal) var(--ease-in-out);
  }

  .cyber-message-item:hover {
    background: var(--color-hover-bg-neutral);
    transform: translateY(-1px);
    box-shadow: 0 0 15px rgba(var(--color-brand-500-rgb), 0.15);
  }

  .cyber-message-item.is-unread {
    background: rgba(var(--color-brand-500-rgb), 0.04);
    border-left: 3px solid var(--color-brand-500);
    box-shadow: 0 0 0 1px rgba(var(--color-brand-500-rgb), 0.1);
  }

  .cyber-message-type-icon {
    width: 32px;
    height: 32px;
    border-radius: var(--radius-sm);
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: var(--text-sm);
    transition: all var(--transition-normal) ease;
    position: relative;
    overflow: hidden;
  }

  .cyber-message-type-icon {
    background: linear-gradient(135deg, rgba(var(--color-brand-500-rgb), 0.15), rgba(var(--color-background-900-rgb), 0.6));
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.3);
    color: var(--color-brand-500);
  }

  .cyber-message-type-icon[data-color='success'] {
    background: linear-gradient(135deg, rgba(var(--color-success-rgb), 0.15), rgba(var(--color-background-900-rgb), 0.6));
    border-color: rgba(var(--color-success-rgb), 0.3);
    color: var(--color-success-500);
  }

  .cyber-message-type-icon[data-color='error'] {
    background: linear-gradient(135deg, rgba(var(--color-error-rgb), 0.15), rgba(var(--color-background-900-rgb), 0.6));
    border-color: rgba(var(--color-error-rgb), 0.3);
    color: var(--color-error-500);
  }

  .cyber-message-type-icon[data-color='warning'] {
    background: linear-gradient(135deg, rgba(var(--color-warning-rgb), 0.15), rgba(var(--color-background-900-rgb), 0.6));
    border-color: rgba(var(--color-warning-rgb), 0.3);
    color: var(--color-warning-500);
  }

  .cyber-message-type-icon[data-color='info'] {
    background: linear-gradient(135deg, rgba(var(--color-info-rgb), 0.15), rgba(var(--color-background-900-rgb), 0.6));
    border-color: rgba(var(--color-info-rgb), 0.3);
    color: var(--color-info-500);
  }

  .cyber-message-item:hover .cyber-message-type-icon {
    transform: scale(1.05);
    box-shadow: 0 0 15px currentColor;
  }

  .cyber-message-item:hover .cyber-message-type-icon[data-color='brand'] {
    background: linear-gradient(135deg, rgba(var(--color-brand-500-rgb), 0.25), rgba(var(--color-background-700-rgb), 0.8));
    border-color: var(--color-hover-border);
  }

  .cyber-message-item:hover .cyber-message-type-icon[data-color='success'] {
    background: linear-gradient(135deg, rgba(var(--color-success-rgb), 0.25), rgba(var(--color-background-700-rgb), 0.8));
    border-color: rgba(var(--color-success-rgb), 0.5);
  }

  .cyber-message-item:hover .cyber-message-type-icon[data-color='error'] {
    background: linear-gradient(135deg, rgba(var(--color-error-rgb), 0.25), rgba(var(--color-background-700-rgb), 0.8));
    border-color: rgba(var(--color-error-rgb), 0.5);
  }

  .cyber-message-item:hover .cyber-message-type-icon[data-color='warning'] {
    background: linear-gradient(135deg, rgba(var(--color-warning-rgb), 0.25), rgba(var(--color-background-700-rgb), 0.8));
    border-color: rgba(var(--color-warning-rgb), 0.5);
  }

  .cyber-message-item:hover .cyber-message-type-icon[data-color='info'] {
    background: linear-gradient(135deg, rgba(var(--color-info-rgb), 0.25), rgba(var(--color-background-700-rgb), 0.8));
    border-color: rgba(var(--color-info-rgb), 0.5);
  }

  .cyber-unread-dot {
    width: 6px;
    height: 6px;
    background: var(--color-brand-500);
    border-radius: var(--radius-full);
    animation: dot-pulse 2s infinite;
    box-shadow: 0 0 4px rgba(var(--color-brand-500-rgb), 0.5);
  }

  @keyframes dot-pulse {
    0%,
    100% {
      opacity: 1;
    }
    50% {
      opacity: 0.6;
    }
  }

  .message-title {
    color: var(--color-content-heading);
  }

  .cyber-message-tag-small {
    display: inline-flex;
    align-items: center;
    padding: var(--space-xs) var(--space-sm);
    border-radius: var(--radius-md);
    font-size: 10px;
    font-weight: var(--font-medium);
    background: linear-gradient(135deg, rgba(var(--color-brand-500-rgb), 0.15), rgba(var(--color-brand-500-rgb), 0.1));
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.3);
    color: var(--color-brand-500);
    transition: all var(--transition-normal) ease;
    flex-shrink: 0;
  }

  .cyber-message-item:hover .cyber-message-tag-small {
    background: linear-gradient(135deg, rgba(var(--color-brand-500-rgb), 0.25), rgba(var(--color-brand-500-rgb), 0.15));
    border-color: var(--color-hover-border);
    transform: scale(1.05);
  }

  .cyber-spinner {
    width: 24px;
    height: 24px;
    border: 2px solid rgba(var(--color-brand-500-rgb), 0.2);
    border-top: 2px solid var(--color-brand-500);
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

  .cyber-dropdown-footer {
    background: linear-gradient(135deg, rgba(var(--color-background-700-rgb), 0.8), rgba(var(--color-brand-500-rgb), 0.05));
    padding: var(--space-md);
    border-top: 1px solid var(--color-border-default);
  }

  .cyber-view-all-link {
    color: var(--color-brand-500);
    padding: var(--space-sm) var(--space-md);
    border-radius: var(--radius-sm);
    position: relative;
    overflow: hidden;
    background: var(--color-hover-bg-neutral);
    border: 1px solid var(--color-border-subtle);
    transition: all var(--transition-normal) ease;
    font-weight: var(--font-medium);
  }

  .cyber-view-all-link:hover {
    color: var(--color-hover-text);
    background: var(--color-hover-bg);
    border-color: rgba(var(--color-brand-500-rgb), 0.3);
    box-shadow: 0 0 15px rgba(var(--color-brand-500-rgb), 0.15);
    transform: translateY(-1px);
  }

  .line-clamp-2 {
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
  }

  .empty-state {
    background: linear-gradient(135deg, rgba(var(--color-brand-500-rgb), 0.02), rgba(var(--color-background-900-rgb), 0.3));
    border: 1px dashed var(--color-border-default);
    border-radius: var(--radius-sm);
    margin: var(--space-sm);
    position: relative;
    overflow: hidden;
  }

  .empty-icon {
    position: relative;
    z-index: 1;
    width: 48px;
    height: 48px;
    border-radius: var(--radius-full);
    background: linear-gradient(135deg, rgba(var(--color-brand-500-rgb), 0.1), rgba(var(--color-background-900-rgb), 0.6));
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.2);
    display: flex;
    align-items: center;
    justify-content: center;
    color: var(--color-content-muted);
    transition: all var(--transition-normal) ease;
  }

  .empty-state:hover .empty-icon {
    transform: scale(1.05);
    border-color: rgba(var(--color-brand-500-rgb), 0.4);
    color: var(--color-brand-500);
    box-shadow: 0 0 20px rgba(var(--color-brand-500-rgb), 0.2);
  }

  .empty-title {
    position: relative;
    z-index: 1;
    font-size: var(--text-base);
    font-weight: var(--font-semibold);
    color: var(--color-content-heading);
    margin: 0;
  }

  .empty-description {
    position: relative;
    z-index: 1;
    font-size: var(--text-sm);
    color: var(--color-content-muted);
    margin: 0;
    line-height: var(--leading-normal);
  }

  .cyber-dropdown {
    position: absolute;
  }

  .cyber-dropdown::after,
  .cyber-dropdown::after {
    border-bottom: 8px solid rgba(var(--color-brand-500-rgb), 0.5);
    filter: drop-shadow(0 1px 2px rgba(0, 0, 0, 0.3));
    z-index: 11;
  }
</style>
