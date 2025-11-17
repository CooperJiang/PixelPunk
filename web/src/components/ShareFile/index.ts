import type { App } from 'vue'
import ShareFile from './index.vue'
import type { ShareFileProps, ShareFileEmits, ShareFileData, ShareFileTag } from './types'

export { ShareFile }
export type { ShareFileProps, ShareFileEmits, ShareFileData, ShareFileTag }

export default {
  install(app: App) {
    app.component('CyberShareFile', ShareFile)
  },
}
