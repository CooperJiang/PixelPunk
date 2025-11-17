import CyberContextMenu from './index.vue'
import type { App } from 'vue'

export type { ContextMenuItem, ContextMenuOptions, CyberContextMenuProps } from './types'

CyberContextMenu.install = (app: App) => {
  app.component('CyberContextMenu', CyberContextMenu)
}

export default CyberContextMenu
