import type { App } from 'vue'
import ImageExpiryTag from './index.vue'

/* 导出类型 */
export type { ImageExpiryTagProps, ImageExpiryData, ExpiryStatusConfig } from './types'

/* 导出组件 */
export { ImageExpiryTag }

/* 默认导出 */
export default ImageExpiryTag

/* Vue插件安装 */
export const installImageExpiryTag = (app: App) => {
  app.component('ImageExpiryTag', ImageExpiryTag)
}
