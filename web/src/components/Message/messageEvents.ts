/* 消息事件类型 */
type MessageEventType = 'new-message' | 'message-read' | 'message-deleted'

/* 消息事件管理器 */
class MessageEventManager {
  private listeners = new Map<MessageEventType, Set<(data: any) => void>>()

  on(type: MessageEventType, callback: (data: any) => void): void {
    if (!this.listeners.has(type)) {
      this.listeners.set(type, new Set())
    }
    this.listeners.get(type)?.add(callback)
  }

  off(type: MessageEventType, callback: (data: any) => void): void {
    const callbacks = this.listeners.get(type)
    if (callbacks) {
      callbacks.delete(callback)
      if (callbacks.size === 0) {
        this.listeners.delete(type)
      }
    }
  }

  emit(type: MessageEventType, data: any): void {
    const callbacks = this.listeners.get(type)
    if (callbacks) {
      callbacks.forEach((callback) => {
        try {
          callback(data)
        } catch {}
      })
    }
  }

  clear(): void {
    this.listeners.clear()
  }
}

export const messageEvents = new MessageEventManager()

export function useMessageEvents() {
  return {
    on: messageEvents.on.bind(messageEvents),
    off: messageEvents.off.bind(messageEvents),
    emit: messageEvents.emit.bind(messageEvents),
    clear: messageEvents.clear.bind(messageEvents),
  }
}

export function setupMessageEventHandlers() {}

export default messageEvents
