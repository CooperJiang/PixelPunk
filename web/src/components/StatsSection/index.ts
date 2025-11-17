import type { App } from 'vue'
import StatsSection from './index.vue'

/* Vue插件安装（按规范注册为 CyberStatsSection） */
StatsSection.install = (app: App) => {
  app.component('CyberStatsSection', StatsSection)
}

export default StatsSection
export { StatsSection }
