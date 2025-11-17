/**
 * 全局WebSocket管理组合函数
 * 用于管理后台管理系统的WebSocket连接生命周期

 */
import { onMounted, onUnmounted } from 'vue'
import { useRoute } from 'vue-router'
import { useWebSocketStore } from '@/store/websocket'
import { useAuthStore } from '@/store/auth'

/**
 * 全局WebSocket管理
 * 在管理员路由下自动初始化WebSocket连接
 */
export function useGlobalWebSocket() {
  const route = useRoute()
  const webSocketStore = useWebSocketStore()
  const authStore = useAuthStore()

  const shouldConnectWebSocket = () => {
    if (!authStore.isAdmin) {
      return false
    }

    if (!route.path.startsWith('/admin')) {
      return false
    }

    return true
  }

  const initWebSocketIfNeeded = () => {
    if (shouldConnectWebSocket()) {
      webSocketStore.initWebSocket()
    }
  }

  onMounted(() => {
    initWebSocketIfNeeded()
  })

  onUnmounted(() => {})

  return {
    webSocketStore,
    shouldConnectWebSocket,
    initWebSocketIfNeeded,
  }
}

export function useWebSocketRouteGuard() {
  const webSocketStore = useWebSocketStore()
  const authStore = useAuthStore()

  const onRouteEnter = (to: any) => {
    if (authStore.isAdmin && to.path.startsWith('/admin')) {
      webSocketStore.initWebSocket()
    }
  }

  const onRouteLeave = (from: any, to: any) => {
    if (from.path.startsWith('/admin') && !to.path.startsWith('/admin')) {
      webSocketStore.disconnectWebSocket()
    }
  }

  const onUserLogout = () => {
    webSocketStore.destroyWebSocket()
  }

  return {
    onRouteEnter,
    onRouteLeave,
    onUserLogout,
  }
}
