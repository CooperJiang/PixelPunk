/**
 * WebSocket 组合式函数
 * 提供Vue组件中使用WebSocket的便捷方法

 */
import { onMounted, onUnmounted, ref, type Ref } from 'vue'
import {
  ConnectionStatus,
  MessageType,
  getWebSocketManager,
  type MessageHandler,
  type WebSocketManager,
  type WebSocketMessage,
} from '@/utils/websocket'
import { useAuthStore } from '@/store/auth'

export interface UseWebSocketOptions {
  autoConnect?: boolean
  debug?: boolean
}

export function useWebSocket(options: UseWebSocketOptions = {}) {
  const { autoConnect = true, debug = false } = options

  const isConnected = ref(false)
  const connectionStatus = ref<ConnectionStatus>(ConnectionStatus.DISCONNECTED)
  const lastError = ref<Error | null>(null)
  const messageCount = ref(0)

  let wsManager: WebSocketManager | null = null

  const authStore = useAuthStore()

  const initWebSocket = () => {
    if (wsManager) {
      return wsManager
    }

    const { token } = authStore
    wsManager = getWebSocketManager(token)

    wsManager.onStatusChange((status) => {
      connectionStatus.value = status
      isConnected.value = status === ConnectionStatus.CONNECTED
    })

    wsManager.onError((error) => {
      lastError.value = error

      if (debug) {
      }
    })

    return wsManager
  }

  const connect = () => {
    const manager = initWebSocket()
    manager.connect()
  }

  const disconnect = () => {
    if (wsManager) {
      wsManager.disconnect()
      isConnected.value = false
      connectionStatus.value = ConnectionStatus.DISCONNECTED
    }
  }

  const send = (message: Partial<WebSocketMessage>) => {
    if (wsManager && wsManager.isConnected()) {
      wsManager.send(message)
    } else if (debug) {
      console.warn('[useWebSocket] Cannot send message: not connected')
    }
  }

  const subscribe = (type: MessageType, handler: MessageHandler) => {
    if (!wsManager) {
      initWebSocket()
    }

    wsManager?.on(type, (message) => {
      messageCount.value++
      handler(message)
    })

    return () => {
      wsManager?.off(type, handler)
    }
  }

  const subscribeQueueStats = (handler: (data: any) => void) =>
    subscribe(MessageType.QUEUE_STATS, (message) => {
      handler(message.data)
    })

  const subscribeVectorStats = (handler: (data: any) => void) =>
    subscribe(MessageType.VECTOR_STATS, (message) => {
      handler(message.data)
    })

  const subscribeAnnouncement = (handler: (data: any) => void) =>
    subscribe(MessageType.ANNOUNCEMENT, (message) => {
      handler(message.data)
    })

  const subscribeLogs = (handler: (data: any) => void) =>
    subscribe(MessageType.LOGS, (message) => {
      handler(message.data)
    })

  const reconnect = () => {
    disconnect()
    setTimeout(() => {
      connect()
    }, 1000)
  }

  onMounted(() => {
    if (autoConnect && authStore.isAdmin) {
      connect()
    }
  })

  onUnmounted(() => {})

  return {
    isConnected: isConnected as Readonly<Ref<boolean>>,
    connectionStatus: connectionStatus as Readonly<Ref<ConnectionStatus>>,
    lastError: lastError as Readonly<Ref<Error | null>>,
    messageCount: messageCount as Readonly<Ref<number>>,

    connect,
    disconnect,
    reconnect,
    send,
    subscribe,
    subscribeQueueStats,
    subscribeVectorStats,
    subscribeAnnouncement,
    subscribeLogs,

    wsManager: () => wsManager,
  }
}

export function useQueueWebSocket() {
  const ws = useWebSocket({ autoConnect: false })

  const queueStats = ref<any>(null)
  const lastUpdated = ref<Date>(new Date())

  let unsubscribe: (() => void) | null = null

  const startListening = () => {
    if (unsubscribe) {
      return
    } // 已经在监听

    ws.connect()

    unsubscribe = ws.subscribeQueueStats((data) => {
      queueStats.value = data
      lastUpdated.value = new Date()
    })
  }

  const stopListening = () => {
    if (unsubscribe) {
      unsubscribe()
      unsubscribe = null
    }
  }

  onUnmounted(() => {
    stopListening()
  })

  return {
    queueStats: queueStats as Readonly<Ref<any>>,
    lastUpdated: lastUpdated as Readonly<Ref<Date>>,
    isConnected: ws.isConnected,
    connectionStatus: ws.connectionStatus,

    startListening,
    stopListening,
    reconnect: ws.reconnect,
  }
}
