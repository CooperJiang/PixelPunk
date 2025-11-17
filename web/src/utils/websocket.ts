/**
 * WebSocket 管理器
 * 实现与后端WebSocket连接的实时通信

/* WebSocket 消息类型 */
export enum MessageType {
  QUEUE_STATS = 'queue_stats',
  VECTOR_STATS = 'vector_stats',
  LOGS = 'logs',
  ANNOUNCEMENT = 'announcement',
  SYSTEM_STATUS = 'system_status',
  ERROR = 'error',
  PING = 'ping',
  PONG = 'pong',
}

/* 消息优先级 */
export enum MessagePriority {
  HIGH = 'high',
  NORMAL = 'normal',
  LOW = 'low',
}

/* WebSocket 消息接口 */
export interface WebSocketMessage {
  id: string
  type: MessageType
  priority: MessagePriority
  timestamp: number
  source?: string
  data?: any
  require_ack?: boolean
}

/* WebSocket 连接状态 */
export enum ConnectionStatus {
  CONNECTING = 'connecting',
  CONNECTED = 'connected',
  DISCONNECTED = 'disconnected',
  RECONNECTING = 'reconnecting',
  ERROR = 'error',
}

/* WebSocket 管理器选项 */
export interface WebSocketOptions {
  url: string
  token?: string
  reconnectDelay?: number
  maxReconnectAttempts?: number
  pingInterval?: number
  debug?: boolean
}

/* WebSocket 事件回调 */
export type MessageHandler = (message: WebSocketMessage) => void
export type StatusChangeHandler = (status: ConnectionStatus) => void
export type ErrorHandler = (_error: Error) => void

/**
 * WebSocket 管理器类
 */
export class WebSocketManager {
  private ws: WebSocket | null = null
  private options: Required<WebSocketOptions>
  private status: ConnectionStatus = ConnectionStatus.DISCONNECTED
  private reconnectAttempts = 0
  private reconnectTimeout: number | null = null
  private pingInterval: number | null = null

  private messageHandlers: Map<MessageType, MessageHandler[]> = new Map()
  private statusChangeHandlers: StatusChangeHandler[] = []
  private errorHandlers: ErrorHandler[] = []

  constructor(options: WebSocketOptions) {
    this.options = {
      reconnectDelay: 5000,
      maxReconnectAttempts: 5,
      pingInterval: 30000,
      debug: false,
      ...options,
    }

    Object.values(MessageType).forEach((type) => {
      this.messageHandlers.set(type, [])
    })
  }

  connect(): void {
    if (this.ws && this.ws.readyState === WebSocket.OPEN) {
      this.log('WebSocket already connected')
      return
    }

    this.setStatus(ConnectionStatus.CONNECTING)

    try {
      const wsUrl = this.buildWebSocketUrl()
      this.log(`Connecting to WebSocket: ${wsUrl}`)

      this.ws = new WebSocket(wsUrl)
      this.setupEventListeners()
    } catch (error) {
      this.handleError(new Error(`Failed to create WebSocket connection: ${error}`))
    }
  }

  disconnect(): void {
    this.log('Disconnecting WebSocket')

    if (this.reconnectTimeout) {
      clearTimeout(this.reconnectTimeout)
      this.reconnectTimeout = null
    }

    if (this.pingInterval) {
      clearInterval(this.pingInterval)
      this.pingInterval = null
    }

    if (this.ws) {
      this.ws.close(1000, 'Manual disconnect')
      this.ws = null
    }

    this.setStatus(ConnectionStatus.DISCONNECTED)
    this.reconnectAttempts = 0
  }

  send(message: Partial<WebSocketMessage>): void {
    if (!this.ws || this.ws.readyState !== WebSocket.OPEN) {
      this.log('WebSocket not connected, cannot send message')
      return
    }

    const fullMessage: WebSocketMessage = {
      id: this.generateMessageId(),
      type: MessageType.PING,
      priority: MessagePriority.NORMAL,
      timestamp: Date.now(),
      ...message,
    }

    try {
      this.ws.send(JSON.stringify(fullMessage))
      this.log('Message sent:', fullMessage)
    } catch (error) {
      this.handleError(new Error(`Failed to send message: ${error}`))
    }
  }

  on(type: MessageType, handler: MessageHandler): void {
    const handlers = this.messageHandlers.get(type) || []
    handlers.push(handler)
    this.messageHandlers.set(type, handlers)
  }

  off(type: MessageType, handler: MessageHandler): void {
    const handlers = this.messageHandlers.get(type) || []
    const index = handlers.indexOf(handler)
    if (index > -1) {
      handlers.splice(index, 1)
      this.messageHandlers.set(type, handlers)
    }
  }

  onStatusChange(handler: StatusChangeHandler): () => void {
    this.statusChangeHandlers.push(handler)

    return () => {
      const index = this.statusChangeHandlers.indexOf(handler)
      if (index > -1) {
        this.statusChangeHandlers.splice(index, 1)
      }
    }
  }

  onError(handler: ErrorHandler): () => void {
    this.errorHandlers.push(handler)

    return () => {
      const index = this.errorHandlers.indexOf(handler)
      if (index > -1) {
        this.errorHandlers.splice(index, 1)
      }
    }
  }

  getStatus(): ConnectionStatus {
    return this.status
  }

  isConnected(): boolean {
    return this.status === ConnectionStatus.CONNECTED && this.ws?.readyState === WebSocket.OPEN
  }

  subscribeQueueStats(handler: (data: any) => void): () => void {
    const messageHandler: MessageHandler = (message) => {
      handler(message.data)
    }

    this.on(MessageType.QUEUE_STATS, messageHandler)

    return () => {
      this.off(MessageType.QUEUE_STATS, messageHandler)
    }
  }

  subscribeVectorStats(handler: (data: any) => void): () => void {
    const messageHandler: MessageHandler = (message) => {
      handler(message.data)
    }

    this.on(MessageType.VECTOR_STATS, messageHandler)

    return () => {
      this.off(MessageType.VECTOR_STATS, messageHandler)
    }
  }

  subscribeAnnouncement(handler: (data: any) => void): () => void {
    const messageHandler: MessageHandler = (message) => {
      handler(message.data)
    }

    this.on(MessageType.ANNOUNCEMENT, messageHandler)

    return () => {
      this.off(MessageType.ANNOUNCEMENT, messageHandler)
    }
  }

  subscribeLogs(handler: (data: any) => void): () => void {
    const messageHandler: MessageHandler = (message) => {
      handler(message.data)
    }

    this.on(MessageType.LOGS, messageHandler)

    return () => {
      this.off(MessageType.LOGS, messageHandler)
    }
  }

  private buildWebSocketUrl(): string {
    const { url, token } = this.options
    const wsUrl = url.replace(/^http/, 'ws')

    if (token) {
      const separator = wsUrl.includes('?') ? '&' : '?'
      return `${wsUrl}${separator}token=${token}`
    }

    return wsUrl
  }

  private setupEventListeners(): void {
    if (!this.ws) {
      return
    }

    this.ws.onopen = () => {
      this.log('WebSocket connected')
      this.setStatus(ConnectionStatus.CONNECTED)
      this.reconnectAttempts = 0
      this.startPing()
    }

    this.ws.onmessage = (event) => {
      try {
        const message: WebSocketMessage = JSON.parse(event.data)
        this.handleMessage(message)
      } catch (error) {
        this.handleError(new Error(`Failed to parse message: ${error}`))
      }
    }

    this.ws.onclose = (event) => {
      this.log(`WebSocket closed: ${event.code} - ${event.reason}`)
      this.setStatus(ConnectionStatus.DISCONNECTED)

      if (this.pingInterval) {
        clearInterval(this.pingInterval)
        this.pingInterval = null
      }

      if (event.code !== 1000 && this.reconnectAttempts < this.options.maxReconnectAttempts) {
        this.scheduleReconnect()
      }
    }

    this.ws.onerror = (error) => {
      this.log('WebSocket error:', error)
      this.handleError(new Error('WebSocket connection error'))
    }
  }

  private handleMessage(message: WebSocketMessage): void {
    this.log('Message received:', message)

    if (message.type === MessageType.PONG) {
      this.log('Received pong from server')
      return
    }

    const handlers = this.messageHandlers.get(message.type) || []
    handlers.forEach((handler) => {
      try {
        handler(message)
      } catch (error) {
        this.log('Error in message handler:', error)
      }
    })
  }

  private setStatus(status: ConnectionStatus): void {
    if (this.status !== status) {
      this.status = status
      this.log(`Status changed to: ${status}`)

      this.statusChangeHandlers.forEach((handler) => {
        try {
          handler(status)
        } catch (error) {
          this.log('Error in status change handler:', error)
        }
      })
    }
  }

  private handleError(_error: Error): void {
    this.log('Error:', _error.message)
    this.setStatus(ConnectionStatus.ERROR)

    this.errorHandlers.forEach((handler) => {
      try {
        handler(_error)
      } catch (err) {
        this.log('Error in error handler:', err)
      }
    })
  }

  private scheduleReconnect(): void {
    this.reconnectAttempts++
    this.setStatus(ConnectionStatus.RECONNECTING)

    const delay = this.options.reconnectDelay * this.reconnectAttempts
    this.log(`Scheduling reconnect attempt ${this.reconnectAttempts} in ${delay}ms`)

    this.reconnectTimeout = window.setTimeout(() => {
      this.connect()
    }, delay)
  }

  private startPing(): void {
    if (this.pingInterval) {
      clearInterval(this.pingInterval)
    }

    this.pingInterval = window.setInterval(() => {
      if (this.isConnected()) {
        this.send({
          type: MessageType.PING,
          data: { timestamp: Date.now() },
        })
      }
    }, this.options.pingInterval)
  }

  private generateMessageId(): string {
    return `msg_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`
  }

  private log(...args: any[]): void {
    if (this.options.debug && process.env.NODE_ENV !== 'production') {
      import('@/utils/system/logger').then(({ logger }) => {
        logger.debug('[WebSocketManager]', args.join(' '))
      })
    }
  }
}

let globalWebSocketManager: WebSocketManager | null = null

function getWebSocketURL(): string {
  const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
  const hostname = window.location.hostname

  let host: string
  if (process.env.NODE_ENV === 'production') {
    host = window.location.host
  } else {
    host = `${hostname}:9520`
  }

  return `${protocol}//${host}/api/v1/admin/ws/admin`
}

export function getWebSocketManager(token?: string): WebSocketManager {
  if (!globalWebSocketManager) {
    const url = getWebSocketURL()

    globalWebSocketManager = new WebSocketManager({
      url,
      token,
      debug: false,
    })
  }

  return globalWebSocketManager
}

export function destroyWebSocketManager(): void {
  if (globalWebSocketManager) {
    globalWebSocketManager.disconnect()
    globalWebSocketManager = null
  }
}
