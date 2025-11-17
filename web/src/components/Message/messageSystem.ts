import { useMessageToast } from './useMessageToast'
import { useMessageModal } from './useMessageModal'
import { useMessageEvents } from './messageEvents'
import type { Message } from '@/api/message/types'

/* 全局消息系统管理 */
class MessageSystemManager {
  private toast = useMessageToast()
  private modal = useMessageModal()
  private events = useMessageEvents()
  private isInitialized = false

  init(): void {
    if (this.isInitialized) {
      return
    }

    this.events.on('new-message', (message: Message) => {
      this.handleNewMessage(message)
    })

    this.events.on('message-read', (messageId: number) => {
      this.handleMessageRead(messageId)
    })

    this.events.on('message-deleted', (messageId: number) => {
      this.handleMessageDeleted(messageId)
    })

    this.isInitialized = true
  }

  private handleNewMessage(message: Message): void {
    if (this.shouldShowToastForMessage(message)) {
      this.toast.showMessage({
        title: message.title,
        content: message.content,
        type: message.type,
        action_url: message.action_url,
        action_text: message.action_text,
        action_style: message.action_style,
      })
    }

    if (this.shouldPlaySound(message)) {
      this.playNotificationSound()
    }
  }

  private handleMessageRead(_messageId: number): void {}

  private handleMessageDeleted(_messageId: number): void {}

  private shouldShowToastForMessage(message: Message): boolean {
    if (message.priority === 1) {
      return true
    }

    const importantTypes = [
      'system.maintenance',
      'system.announcement',
      'content.review_approved',
      'content.review_rejected',
      'storage.quota_warning',
      'security.login_alert',
      'account.storage_granted',
      'account.bandwidth_granted',
    ]

    return importantTypes.includes(message.type)
  }

  private shouldPlaySound(message: Message): boolean {
    return message.priority === 1 || message.type === 'security.login_alert'
  }

  private playNotificationSound(): void {
    try {
      const audioContext = new (window.AudioContext || (window as any).webkitAudioContext)()
      const oscillator = audioContext.createOscillator()
      const gainNode = audioContext.createGain()

      oscillator.connect(gainNode)
      gainNode.connect(audioContext.destination)

      oscillator.frequency.setValueAtTime(800, audioContext.currentTime)
      oscillator.frequency.setValueAtTime(600, audioContext.currentTime + 0.1)

      gainNode.gain.setValueAtTime(0.1, audioContext.currentTime)
      gainNode.gain.exponentialRampToValueAtTime(0.01, audioContext.currentTime + 0.3)

      oscillator.start(audioContext.currentTime)
      oscillator.stop(audioContext.currentTime + 0.3)
    } catch {}
  }

  showMessageToast(message: Message): void {
    this.toast.showMessage({
      title: message.title,
      content: message.content,
      type: message.type,
      action_url: message.action_url,
      action_text: message.action_text,
      action_style: message.action_style,
    })
  }

  showMessageDetail(_message: Message): void {}

  simulateNewMessage(message: Message): void {
    this.events.emit('new-message', message)
  }

  markMessageAsRead(messageId: number): void {
    this.events.emit('message-read', messageId)
  }

  deleteMessage(messageId: number): void {
    this.events.emit('message-deleted', messageId)
  }

  destroy(): void {
    this.events.clear()
    this.toast.clear()
    this.isInitialized = false
  }
}

export const messageSystem = new MessageSystemManager()

export function useMessageSystem() {
  return {
    init: messageSystem.init.bind(messageSystem),
    showMessageToast: messageSystem.showMessageToast.bind(messageSystem),
    showMessageDetail: messageSystem.showMessageDetail.bind(messageSystem),
    simulateNewMessage: messageSystem.simulateNewMessage.bind(messageSystem),
    markMessageAsRead: messageSystem.markMessageAsRead.bind(messageSystem),
    deleteMessage: messageSystem.deleteMessage.bind(messageSystem),
    destroy: messageSystem.destroy.bind(messageSystem),
  }
}

export default messageSystem
