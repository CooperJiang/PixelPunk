/**
 * 系统相关类型定义
 */
import type { PaginatedResponse, PaginationParams, StatusInfo, TimeStamps } from './common'

/* ==================== 系统安装相关类型 ==================== */
export interface InstallStatus {
  installed: boolean
  database_ok: boolean
  redis_ok: boolean
  message: string
  deploy_mode?: string // standalone | docker | compose
  setup_level?: number // 0=已安装 1=简化配置 2=完整配置
  config_preset?: boolean // 配置是否预设
}

export interface DatabaseTestRequest {
  type: 'mysql' | 'sqlite'
  host?: string
  port?: number
  username?: string
  password?: string
  name?: string
  path?: string
}

export interface InstallRequest {
  database: {
    type: 'mysql' | 'sqlite'
    host?: string
    port?: number
    username?: string
    password?: string
    name?: string
    path?: string
  }
  redis: {
    host: string
    port: number
    password: string
    db: number
  }
  vector?: {
    qdrant_url: string
    qdrant_timeout: number
    use_builtin?: boolean
    http_port?: number
    grpc_port?: number
  }
  admin_username: string
  admin_password: string
}

export interface QdrantTestRequest {
  qdrant_url: string
  qdrant_timeout?: number
}

/* ==================== 存储渠道相关类型 ==================== */
/* 支持后端动态扩展的渠道类型（避免前端硬编码） */
export type StorageChannelType = string

export interface StorageChannel extends StatusInfo, TimeStamps {
  id: string
  name: string
  type: StorageChannelType
  is_default: boolean
  remark?: string
  upload_limit?: number // 上传限制 (字节)
  file_count?: number // 文件数量
}

export interface StorageConfigItem extends TimeStamps {
  id: string
  channel_id: string
  name: string
  key_name: string
  value: string
  type: 'string' | 'int' | 'bool' | 'password'
  is_secret: boolean
  required: boolean
  description: string
}

export interface CreateStorageChannelRequest {
  name: string
  type: StorageChannelType
  status?: number
  is_default?: boolean
  remark?: string
  upload_limit?: number
}

export interface UpdateStorageChannelRequest {
  name?: string
  is_default?: boolean
  status?: number
  remark?: string
  upload_limit?: number
}

export interface UpdateStorageConfigRequest {
  [key: string]: string | number | boolean
}

/* ==================== 存储高级配置类型 ==================== */
export interface StorageAdvancedConfig {
  storage_class?: string
  custom_domain?: string
  access_control?: string
  use_https?: boolean
  hide_remote_url?: boolean
  thumb_max_width?: number
  thumb_max_height?: number
  thumb_quality?: number
  upload_limit_mb?: number
}

export interface StorageChannelExport {
  channel: StorageChannel
  configs: StorageConfigItem[]
}

/* 批量导出渠道配置（排除本地存储） */
export interface StorageChannelsBatchExport {
  channels: StorageChannelExport[] | any[]
  total_count: number
  export_type: string
  version: string
  export_time: string
}

/* ==================== API密钥相关类型 ==================== */
export interface ApiKeyInfo extends StatusInfo, TimeStamps {
  id: string
  name: string
  key?: string
  status_text?: string
  storage_limit: number
  storage_used?: number
  single_file_limit: number
  upload_count_limit: number
  upload_count_used?: number
  allowed_types: string[]
  folder_id?: string
  folder_path?: string
  is_expired?: boolean
  expires_in_days?: number
  last_used_at?: string
  expires_at?: string
}

export interface CreateApiKeyRequest {
  name: string
  storage_limit?: number
  single_file_limit?: number
  upload_count_limit?: number
  allowed_types?: string[]
  folder_id?: string
  expires_in_days?: number
}

export interface UpdateApiKeyRequest {
  name?: string
  status?: number
  storage_limit?: number
  single_file_limit?: number
  upload_count_limit?: number
  allowed_types?: string[]
  folder_id?: string
  expires_in_days?: number
}

export interface ApiKeyListParams extends PaginationParams {
  status?: number
  search?: string
}

export type ApiKeyListResponse = PaginatedResponse<ApiKeyInfo>

export interface ApiKeyStatsResponse {
  storage_used: number
  upload_count: number
  last_used_at?: string
}

export interface RegenerateApiKeyResponse {
  id: string
  key: string
}

export interface ToggleApiKeyResponse {
  id: string
  status: number
  status_text: string
  is_active: boolean
}
