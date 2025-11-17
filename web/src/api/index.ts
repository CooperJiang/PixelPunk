/**
 * API模块统一入口
 * 导出所有API服务和类型定义

/* ==================== 类型定义导出 ==================== */
export * from './types'

/* ==================== 传统模块导出 (保持兼容) ==================== */
import * as folderApi from './folder'
import * as fileApi from './file'
import * as userApi from './user'
import * as tagApi from './tag'
import * as storageApi from './storage'
import * as apikeyApi from './apikey'
import * as statsApi from './stats'
import * as shareApi from './share'
import * as authorApi from './author'

export { folderApi, fileApi, userApi, tagApi, storageApi, apikeyApi, statsApi, shareApi, authorApi }

export default {
  user: userApi,
  file: fileApi,
  folder: folderApi,
  stats: statsApi,
  apikey: apikeyApi,
  storage: storageApi,
  tag: tagApi,
  share: shareApi,
  author: authorApi,
}
