import type { App } from 'vue'
import StatsCard from './index.vue'
export type { StatsCardProps } from './types'

StatsCard.install = (app: App) => {
  app.component('CyberStatsCard', StatsCard)
}

export default StatsCard
export { StatsCard }
