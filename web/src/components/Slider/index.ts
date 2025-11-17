import type { App } from 'vue'
import Slider from './index.vue'
import type { SliderEmits, SliderProps } from './types'

export { Slider }
export type { SliderProps, SliderEmits }

export default {
  install(app: App) {
    app.component('CyberSlider', Slider)
  },
}
