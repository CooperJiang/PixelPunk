/**
 * WebSocket 状态管理 Store
 * 全局管理WebSocket连接状态和实时数据

 */
import { defineStore } from 'pinia'
import { computed, readonly, ref } from 'vue'
import { useAuthStore } from './auth'
import { ConnectionStatus, type WebSocketManager, destroyWebSocketManager, getWebSocketManager } from '@/utils/websocket'
import { useTexts } from '@/composables/useTexts'

/* 队列状态数据类型 */
export interface QueueStatsData {
  queue_stats: {
    total_count: number
    none_count: number
    pending_count: number
    done_count: number
    failed_count: number
    ignored_count: number
  }
  config: {
    max_concurrency: number
    current_concurrency: number
  }
  performance?: {
    avg_processing_time: number
    total_processed_today: number
  }
}

/* 向量统计数据类型 */
export interface VectorStatsData {
  runtime?: {
    active_workers: number
    configured_concurrency: number
    paused: boolean
  }
  queue_stats_ext?: {
    queued: number
    processing: number
    delayed: number
    dlq: number
  }
  queue_stats: {
    completed_count: number
    failed_count: number
    pending_count: number
    reset_count: number
    running_count: number
    total_count: number
  }
}

export const useWebSocketStore = defineStore('websocket', () => {
  const { $t } = useTexts()

  const isConnected = ref(false)
  const connectionStatus = ref<ConnectionStatus>(ConnectionStatus.DISCONNECTED)
  const lastError = ref<Error | null>(null)
  const lastUpdated = ref<Date>(new Date())

  const queueStats = ref<QueueStatsData | null>(null)
  const vectorStats = ref<VectorStatsData | null>(null)
  const announcements = ref<any[]>([])
  const systemLogs = ref<any[]>([])

  let wsManager: WebSocketManager | null = null
  let unsubscribeFunctions: Array<() => void> = []

  const connectionStatusText = computed(() => {
    switch (connectionStatus.value) {
      case ConnectionStatus.CONNECTED:
        return $t('store.websocket.status.connected')
      case ConnectionStatus.CONNECTING:
        return $t('store.websocket.status.connecting')
      case ConnectionStatus.RECONNECTING:
        return $t('store.websocket.status.reconnecting')
      case ConnectionStatus.DISCONNECTED:
        return $t('store.websocket.status.disconnected')
      case ConnectionStatus.ERROR:
        return $t('store.websocket.status.error')
      default:
        return $t('store.websocket.status.unknown')
    }
  })

  const lastUpdatedText = computed(() => {
    const now = new Date()
    const diff = now.getTime() - lastUpdated.value.getTime()
    const seconds = Math.floor(diff / 1000)

    if (seconds < 60) {
      return $t('store.websocket.lastUpdated.seconds', { seconds: String(seconds) })
    }
    const minutes = Math.floor(seconds / 60)
    return $t('store.websocket.lastUpdated.minutes', { minutes: String(minutes) })
  })

  const queueCompletionRate = computed(() => {
    if (!queueStats.value || queueStats.value.queue_stats.total_count === 0) {
      return 0
    }
    return Math.round((queueStats.value.queue_stats.done_count / queueStats.value.queue_stats.total_count) * 100)
  })

  const initWebSocket = () => {
    const authStore = useAuthStore()

    if (!authStore.isAdmin) {
      return
    }

    if (wsManager) {
      return
    }

    try {
      wsManager = getWebSocketManager(authStore.token || undefined)

      const statusUnsubscribe = wsManager.onStatusChange((status) => {
        connectionStatus.value = status
        isConnected.value = status === ConnectionStatus.CONNECTED
      })

      const errorUnsubscribe = wsManager.onError((error) => {
        lastError.value = error
      })

      const queueUnsubscribe = wsManager.subscribeQueueStats((data) => {
        queueStats.value = data
        lastUpdated.value = new Date()
      })

      const vectorUnsubscribe = wsManager.subscribeVectorStats((data) => {
        vectorStats.value = data
        lastUpdated.value = new Date()
      })

      const announcementUnsubscribe = wsManager.subscribeAnnouncement((data) => {
        announcements.value.unshift(data)
        if (announcements.value.length > 50) {
          announcements.value = announcements.value.slice(0, 50)
        }
      })

      const logsUnsubscribe = wsManager.subscribeLogs((data) => {
        systemLogs.value.unshift(data)
        if (systemLogs.value.length > 100) {
          systemLogs.value = systemLogs.value.slice(0, 100)
        }
      })

      unsubscribeFunctions = [
        statusUnsubscribe,
        errorUnsubscribe,
        queueUnsubscribe,
        vectorUnsubscribe,
        announcementUnsubscribe,
        logsUnsubscribe,
      ]

      wsManager.connect()
    } catch (error) {
      lastError.value = error as Error
    }
  }

  const disconnectWebSocket = () => {
    if (wsManager) {
      unsubscribeFunctions.forEach((unsubscribe) => unsubscribe())
      unsubscribeFunctions = []

      wsManager.disconnect()
      wsManager = null
    }

    isConnected.value = false
    connectionStatus.value = ConnectionStatus.DISCONNECTED
    lastError.value = null
  }

  const destroyWebSocket = () => {
    disconnectWebSocket()
    destroyWebSocketManager()

    queueStats.value = null
    vectorStats.value = null
    announcements.value = []
    systemLogs.value = []
  }

  const reconnectWebSocket = () => {
    disconnectWebSocket()
    setTimeout(() => {
      initWebSocket()
    }, 1000)
  }

  const sendMessage = (message: Event) => {
    if (wsManager && isConnected.value) {
      wsManager.send(message)
    }
  }

  return {
    isConnected: readonly(isConnected),
    connectionStatus: readonly(connectionStatus),
    connectionStatusText,
    lastError: readonly(lastError),
    lastUpdated: readonly(lastUpdated),
    lastUpdatedText,

    queueStats: readonly(queueStats),
    vectorStats: readonly(vectorStats),
    announcements: readonly(announcements),
    systemLogs: readonly(systemLogs),

    queueCompletionRate,

    initWebSocket,
    disconnectWebSocket,
    destroyWebSocket,
    reconnectWebSocket,
    sendMessage,
  }
})
