/**
 * API服务统一导出
 * 将所有API服务按模块组织并统一导出

/* ==================== 用户服务 ==================== */
export { default as userApi } from '../user'

/* ==================== 文件服务 ==================== */
export { default as fileApi } from '../file'

/* ==================== 文件夹服务 ==================== */
export { default as folderApi } from '../folder'

/* ==================== 消息服务 ==================== */
export { default as messageApi } from '../message'

/* ==================== 标签服务 ==================== */
export { default as tagApi } from '../tag'

/* ==================== 分享服务 ==================== */
export { default as shareApi } from '../share'

/* ==================== 存储服务 ==================== */
export { default as storageApi } from '../storage'

/* ==================== API密钥服务 ==================== */
export { default as apikeyApi } from '../apikey'

/* ==================== 统计服务 ==================== */
export { default as statsApi } from '../stats'

/* ==================== 搜索服务 ==================== */
export { default as searchApi } from '../search'

/* ==================== 公共服务 ==================== */
export { default as commonApi } from '../common'

/* ==================== 管理员服务 ==================== */
export { default as adminApi } from '../admin'

/* ==================== 作者页面服务 ==================== */
export { default as authorApi } from '../author'

/* ==================== 仪表盘服务 ==================== */
export { default as dashboardApi } from '../dashboard'

/* ==================== 上传服务 ==================== */
export * as uploadApi from '../upload'

/* ==================== 系统设置服务 ==================== */
export { default as setupApi } from '../setup'

/* ==================== 工作区服务 ==================== */
export * as workspaceApi from '../workspace'

/* ==================== 活动服务 ==================== */
export { default as activityApi } from '../activity'
