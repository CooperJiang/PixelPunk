/**
 * API类型统一导出
 * 新的分层类型系统入口文件

/* ==================== 通用类型 ==================== */
export * from './common'

/* ==================== 用户模块类型 ==================== */
export * from './user'

/* ==================== 文件模块类型 ==================== */
export * from './file'

/* ==================== 文件夹模块类型 ==================== */
export * from './folder'

/* ==================== 系统模块类型 ==================== */
export * from './system'

/* ==================== 标签模块类型 ==================== */
export * from './tag'

/* ==================== 分类模块类型 ==================== */
export * from './category'

/* ==================== 消息模块类型 ==================== */
export * from './message'

/* ==================== 分享模块类型 ==================== */
export * from './share'

/* ==================== 通用导出 ==================== */
export type { ApiResponse } from './common'
export type { UserInfo, LoginResponse } from './user'
export type { FileInfo, FileListResponse } from './file'
export type { FolderInfo, FolderListResponse } from './folder'

/* 向量管理模块类型 */
export type {
  VectorListParams,
  VectorItem,
  VectorStats,
  VectorLogItem,
  VectorLogParams,
  PaginationInfo,
} from '../admin/vectors/index'
