/**
 * 设置相关类型定义

 */
export type SettingType = 'string' | 'number' | 'boolean' | 'json' | 'array'

export type SettingGroup =
  | 'website'
  | 'website_info'
  | 'registration'
  | 'security'
  | 'mail'
  | 'ai'
  | 'vector'
  | 'upload'
  | 'construction'
  | 'theme'
  | 'guest'
  | 'appearance'
  | 'legal'
  | 'oauth'
  | 'analytics'

export interface Setting {
  id?: number
  key: string
  value: string | number | boolean
  type: SettingType
  group: SettingGroup
  description: string
  is_system?: boolean
  created_at?: string
  updated_at?: string
}

export interface SettingsListResponse {
  settings: Setting[]
}

export interface SettingsQuery {
  group?: SettingGroup
  key?: string
}

/* OAuth 提供商状态 */
export interface OAuthProviders {
  github_enabled: boolean // GitHub 登录是否启用
  google_enabled: boolean // Google 登录是否启用
  linuxdo_enabled: boolean // Linux DO 登录是否启用
}

/* 全局设置响应接口 - 优化为对象套对象格式 */
export interface GlobalSettingsResponse {
  construction?: Record<string, unknown> // 建设相关设置
  website?: Record<string, unknown> // 网站后端功能设置
  website_info?: Record<string, unknown> // 网站前端显示设置
  upload?: Record<string, unknown> & {
    is_allow_chunk_upload?: boolean // 是否允许分片上传（仅local存储允许）
  } // 上传相关设置
  theme?: Record<string, unknown> // 网站装修设置
  registration?: Record<string, unknown> // 注册相关设置
  version?: Record<string, unknown> // 版本信息设置
  ai?: Record<string, unknown> // AI相关设置
  vector?: Record<string, unknown> // 向量搜索相关设置
  guest?: Record<string, unknown> // 游客上传相关设置
  appearance?: Record<string, unknown> // 外观界面设置
  analytics?: Record<string, unknown> // 埋点统计相关设置
  oauth_providers: OAuthProviders // OAuth 提供商状态
  announcement_version?: string // 公告版本号
  deploy_mode?: string // 部署模式：standalone, docker, compose
}

/* 邮件服务器测试参数 */
export interface TestMailServerParams {
  email: string
  smtp_host: string
  smtp_port: number | string
  smtp_encryption: string
  smtp_username: string
  smtp_password: string
  smtp_from_address: string
  smtp_from_name: string
}

/* 向量配置测试参数 */
export interface VectorTestParams {
  provider: string
  model: string
  api_key: string
  base_url: string
  timeout: number
}

/* 向量配置测试结果 */
export interface VectorTestResult {
  success: boolean
  message: string
  details?: string
  config?: {
    provider: string
    model: string
    baseUrl: string
    timeout: number
    dimension?: number
  }
}

/* 重新生成向量参数 */
export interface RegenerateVectorsParams {
  file_id?: string
  limit?: number
}

/* 重新生成向量结果 */
export interface RegenerateVectorsResult {
  processed: number
  failed: number
  total: number
  message: string
}

/* AI配置测试结果类型 */
export interface AITestResult {
  success: boolean
  message: string
  model?: string
  api_endpoint?: string
  test_response?: string
  tokens_used?: number
  status?: string
  error?: string
  details?: string
}

/* AI配置测试API响应类型 */
export interface AITestResponse {
  code: number
  message: string
  data: AITestResult
  request_id: string
  timestamp: number
}

/* Qdrant连接测试参数 */
export interface QdrantTestParams {
  qdrant_url: string
  qdrant_timeout: number
}

/* Qdrant连接测试结果 */
export interface QdrantTestResult {
  success: boolean
  message: string
  details?: {
    url?: string
    version?: string
    collections?: number
    total_vectors?: number
    response_time?: number
  }
  error?: string
}

/* 法律文档响应接口 */
export interface LegalDocumentsResponse {
  privacy_policy: string // 隐私政策内容
  terms_of_service: string // 服务条款内容
}
