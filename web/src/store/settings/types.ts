/**
 * Settings Store 类型定义
 * 统一导出所有设置相关的类型
 */

/* 网站设置接口 */
export interface WebsiteSettings {
  admin_email?: string
  site_base_url?: string
}

/* 网站信息设置接口 */
export interface WebsiteInfoSettings {
  site_name?: string
  site_description?: string
  site_keywords?: string
  icp_number?: string
  show_file_count?: boolean
  show_storage_usage?: boolean
  site_logo_url?: string
  favicon_url?: string
  copyright_text?: string
  contact_email?: string
  footer_custom_text?: string
  site_hero_title?: string
  site_features_text?: string
  site_slogan?: string
}

/* 上传设置接口 */
export interface UploadSettings {
  allowed_image_formats?: string[]
  max_file_size?: number
  max_batch_size?: number
  thumbnail_quality?: number
  thumbnail_max_width?: number
  thumbnail_max_height?: number
  preserve_exif?: boolean
  daily_upload_limit?: number
}

/* 注册设置接口 */
export interface RegistrationSettings {
  enable_registration?: boolean
  email_verification?: boolean
}

/* 版本设置接口 */
export interface VersionSettings {
  current_version?: string
  build_time?: string
  update_available?: boolean
  last_update_check?: string
}

/* AI设置接口 */
export interface AISettings {
  ai_enabled?: boolean
}

/* 向量设置接口 */
export interface VectorSettings {
  vector_enabled?: boolean
}

/* 外观设置接口 */
export interface AppearanceSettings {
  show_official_site?: boolean
  official_site_url?: string
  show_github_link?: boolean
  github_url?: string
  show_wechat_group?: boolean
  wechat_qr_image_url?: string
  wechat_contact_account?: string
  show_qq_group?: boolean
  qq_qr_image_url?: string
  qq_group_number?: string
  enable_multi_layout?: boolean
  default_layout?: string
  enable_multi_language?: boolean
  default_language?: string
  default_theme?: string
}

/* 埋点统计设置接口 */
export interface AnalyticsSettings {
  baidu_analytics_enabled?: boolean | string
  baidu_analytics_site_id?: string
  google_analytics_enabled?: boolean | string
  google_analytics_measurement_id?: string
}
