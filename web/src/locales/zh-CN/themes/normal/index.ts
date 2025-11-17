import { activity } from './activity'
import { message } from './message'
import { common } from './common'
import { admin } from './admin/index'
import { api } from './api/index'
import { components } from './components'
import { composables } from './composables/index'
import { constants } from './constants'
import { layouts } from './layouts'
import { pages } from './pages/index'
import { store } from './store'
import { theme } from './theme'
import { network, validation, formatting, errorHandler, file, common as utilsCommon } from './utils'

/**
 * 常规主题完整导出
 */
export const normalTheme = {
  activity,
  ...common,
  message, // 放在 common 后面，这样 message 不会被 common.message 覆盖
  admin,
  api,
  components,
  composables,
  constants,
  layouts,
  ...pages,
  store,
  theme,
  utils: {
    formatting,
    errorHandler,
    file,
    common: utilsCommon,
    ...network,
    ...validation,
  },
}
