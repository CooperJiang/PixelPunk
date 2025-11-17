import CyberTable from './index.vue'
import type { App } from 'vue'
import type { CyberTableColumn, CyberTableEmits, CyberTableInstance, CyberTableProps } from './types'

/* 组件安装函数 */
CyberTable.install = (app: App) => {
  app.component('CyberTable', CyberTable)
  app.component('CyberDataTable', CyberTable) /* 别名 */
}

/* 导出组件和类型 */
export default CyberTable
export type { CyberTableProps, CyberTableColumn, CyberTableEmits, CyberTableInstance }
