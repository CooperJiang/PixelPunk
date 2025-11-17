/**
 * 文件夹相关类型定义
 */
import type { AccessLevel, TimeStamps } from './common'
import type { AuthorFileInfo } from './file'

/* ==================== 文件夹基础信息类型 ==================== */
export interface FolderInfo extends TimeStamps {
  id: string
  name: string
  parent_id?: string
  permission?: AccessLevel
  description?: string
  file_count?: number
  total_size?: number
  has_children?: boolean
  sort_order?: number
  level?: number // 文件夹层级，从1开始
}

/* ==================== 文件夹创建请求类型 ==================== */
export interface CreateFolderRequest {
  name: string
  parent_id?: string
  permission?: AccessLevel
  description?: string
}

/* ==================== 文件夹更新请求类型 ==================== */
export interface UpdateFolderRequest {
  id?: string
  name?: string
  parent_id?: string
  permission?: AccessLevel
  description?: string
}

/* ==================== 文件夹列表响应类型 ==================== */
export interface FolderListResponse {
  items: FolderInfo[]
}

/* ==================== 文件夹树节点类型 ==================== */
export interface FolderTreeNode {
  id: string
  label: string
  icon: string
  count?: number
  level?: number // 树节点层级，从1开始
  children?: FolderTreeNode[]
  data?: FolderInfo
}

/* ==================== 文件夹路径链类型 ==================== */
export interface FolderPathChain {
  id: string
  name: string
  parent_id?: string
  level: number
}

export interface FolderPathChainResponse {
  folder_id: string
  full_path: string
  path_chain: FolderPathChain[]
  total_levels: number
}

/* ==================== 作者文件夹相关类型 ==================== */
export interface AuthorFolderInfo {
  id: string
  name: string
  fileCount: number
  coverFile: string
  coverFileFullPath: string
  coverFileThumbURL: string
  coverFileFullThumbURL: string
  createdAt: string
  totalSize: number
  totalSizeFormatted: string
}

export interface AuthorFolderContent {
  folder: AuthorFolderInfo
  files: AuthorFileInfo[]
  subFolders: AuthorFolderInfo[]
  pagination: {
    currentPage: number
    perPage: number
    total: number
    lastPage: number
  }
}
