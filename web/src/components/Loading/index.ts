import type { App } from 'vue'
import Loading from './index.vue'
import type { LoadingProps } from './types'

export { Loading }
export type { LoadingProps }

export default {
  install(app: App) {
    app.component('CyberLoading', Loading)
  },
}
