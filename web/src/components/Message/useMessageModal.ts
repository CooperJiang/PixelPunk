import { createApp, type App } from 'vue'
import MessageDetail from './MessageDetail.vue'
import type { Message } from '@/api/message/types'
import router from '@/router'
import CyberComponentsPlugin from '@/components'

interface MessageModalInstance {
  close: () => void
}

class MessageModalManager {
  private currentModal: { app: App; mountPoint: HTMLElement } | null = null

  showMessageDetail(message: Message): MessageModalInstance {
    if (this.currentModal) {
      this.closeModal()
    }

    const mountPoint = document.createElement('div')
    document.body.appendChild(mountPoint)

    const app = createApp(MessageDetail, {
      message,
      show: true,
      onClose: () => {
        this.closeModal()
      },
      onUpdated: (_updatedMessage: Message) => {},
      onDeleted: (_messageId: number) => {
        this.closeModal()
      },
    })

    app.use(router)
    app.use(CyberComponentsPlugin)

    app.mount(mountPoint)

    this.currentModal = { app, mountPoint }

    return {
      close: () => {
        this.closeModal()
      },
    }
  }

  private closeModal(): void {
    if (!this.currentModal) {
      return
    }

    const { app, mountPoint } = this.currentModal

    app.unmount()

    if (mountPoint.parentElement) {
      mountPoint.parentElement.removeChild(mountPoint)
    }

    this.currentModal = null
  }

  hasOpenModal(): boolean {
    return this.currentModal !== null
  }
}

export const messageModal = new MessageModalManager()

export function useMessageModal() {
  return {
    showMessageDetail: messageModal.showMessageDetail.bind(messageModal),
    hasOpenModal: messageModal.hasOpenModal.bind(messageModal),
  }
}

export type { MessageModalInstance }
export default messageModal
