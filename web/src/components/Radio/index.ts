import type { App } from 'vue'
import Radio from './index.vue'
import RadioGroup from './RadioGroup.vue'
import type { RadioProps, RadioEmits, RadioGroupProps, RadioGroupEmits, RadioGroupOption } from './types'

export { Radio, RadioGroup }
export type { RadioProps, RadioEmits, RadioGroupProps, RadioGroupEmits, RadioGroupOption }

export default {
  install(app: App) {
    app.component('CyberRadio', Radio)
    app.component('CyberRadioGroup', RadioGroup)
  },
}
