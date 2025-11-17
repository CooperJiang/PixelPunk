import type { App } from 'vue'
import lazyDirective from './lazy'
import loadingDirective from './loading'
/* 布局自适应相关指令（全局注册） */
import { vLayoutAdaptive, vResponsiveGrid, vResponsiveSpacing } from './layout-adaptive'
/* 右键菜单指令 */
import { vContextMenu } from './contextMenu'

/**
 * 注册所有自定义指令
 */
export default {
  install(app: App) {
    app.use(lazyDirective)
    app.use(loadingDirective)
    app.directive('layout-adaptive', vLayoutAdaptive)
    app.directive('responsive-grid', vResponsiveGrid)
    app.directive('responsive-spacing', vResponsiveSpacing)
    app.directive('context-menu', vContextMenu)
  },
}

export { lazyDirective, loadingDirective }
