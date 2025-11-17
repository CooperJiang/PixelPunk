/**
 * 所有设置的默认配置

 */
import type { TranslationFunction } from '@/composables/useTexts'
import type { Setting, SettingGroup } from './types'
import { mailDefaults } from './mail'
import { aiDefaults } from './ai'
import { vectorDefaults } from './vector'

/* 网站后端功能设置（带翻译） */
export function getWebsiteDefaults($t: TranslationFunction): Setting[] {
  return [
    {
      key: 'admin_email',
      value: 'admin@example.com',
      type: 'string',
      group: 'website',
      description: $t('api.settingsDefaults.website.admin_email'),
      is_system: true,
    },
    {
      key: 'site_base_url',
      value: '',
      type: 'string',
      group: 'website',
      description: $t('api.settingsDefaults.website.site_base_url'),
      is_system: true,
    },
  ]
}

/* 向后兼容 - 默认中文 */
export const websiteDefaults: Setting[] = [
  {
    key: 'admin_email',
    value: 'admin@example.com',
    type: 'string',
    group: 'website',
    description: '管理员邮箱',
    is_system: true,
  },
  {
    key: 'site_base_url',
    value: '',
    type: 'string',
    group: 'website',
    description: '网站基础URL',
    is_system: true,
  },
]

/* 网站前端显示配置 - 不再使用默认值，完全依赖后端接口 */
export const websiteInfoDefaults: Setting[] = []

/* 注册设置 */
export const registrationDefaults: Setting[] = [
  {
    key: 'enable_registration',
    value: true,
    type: 'boolean',
    group: 'registration',
    description: '开放注册',
    is_system: true,
  },
  {
    key: 'email_verification',
    value: true,
    type: 'boolean',
    group: 'registration',
    description: '邮箱验证',
    is_system: true,
  },
  {
    key: 'user_initial_storage',
    value: 100,
    type: 'number',
    group: 'registration',
    description: '新用户默认存储空间(MB)',
    is_system: true,
  },
  {
    key: 'user_initial_bandwidth',
    value: 1024,
    type: 'number',
    group: 'registration',
    description: '新用户默认带宽流量(MB)',
    is_system: true,
  },
]

/* 安全设置 */
export const securityDefaults: Setting[] = [
  {
    key: 'max_login_attempts',
    value: 5,
    type: 'number',
    group: 'security',
    description: '最大登录尝试次数',
    is_system: true,
  },
  {
    key: 'account_lockout_minutes',
    value: 30,
    type: 'number',
    group: 'security',
    description: '账户锁定分钟数',
    is_system: true,
  },
  {
    key: 'login_expire_hours',
    value: 240,
    type: 'number',
    group: 'security',
    description: '登录有效期',
    is_system: true,
  },
  {
    key: 'jwt_secret',
    value: '',
    type: 'string',
    group: 'security',
    description: '登录安全密钥',
    is_system: true,
  },
  {
    key: 'ip_whitelist',
    value: '',
    type: 'string',
    group: 'security',
    description: 'IP白名单',
    is_system: true,
  },
  {
    key: 'ip_blacklist',
    value: '',
    type: 'string',
    group: 'security',
    description: 'IP黑名单',
    is_system: true,
  },
  {
    key: 'domain_whitelist',
    value: '',
    type: 'string',
    group: 'security',
    description: '域名白名单',
    is_system: true,
  },
  {
    key: 'domain_blacklist',
    value: '',
    type: 'string',
    group: 'security',
    description: '域名黑名单',
    is_system: true,
  },
  {
    key: 'hide_remote_url',
    value: true,
    type: 'boolean',
    group: 'security',
    description: '隐藏三方存储地址（全局优先，渠道未设置时回退）',
    is_system: true,
  },
]

/* 上传设置 */
export const uploadDefaults: Setting[] = [
  {
    key: 'allowed_image_formats',
    value: ['jpg', 'jpeg', 'png', 'gif', 'webp', 'bmp', 'svg', 'ico', 'apng', 'tiff', 'tif', 'heic', 'heif'],
    type: 'array',
    group: 'upload',
    description: '允许上传的文件格式',
    is_system: true,
  },
  {
    key: 'max_file_size',
    value: 20,
    type: 'number',
    group: 'upload',
    description: '单个文件最大大小(MB)',
    is_system: true,
  },
  {
    key: 'max_batch_size',
    value: 100,
    type: 'number',
    group: 'upload',
    description: '批量上传总大小限制(MB)',
    is_system: true,
  },
  {
    key: 'thumbnail_quality',
    value: 50,
    type: 'number',
    group: 'upload',
    description: '缩略图质量设置(0-100)',
    is_system: true,
  },
  {
    key: 'thumbnail_max_width',
    value: 200,
    type: 'number',
    group: 'upload',
    description: '缩略图最大宽度',
    is_system: true,
  },
  {
    key: 'thumbnail_max_height',
    value: 200,
    type: 'number',
    group: 'upload',
    description: '缩略图最大高度',
    is_system: true,
  },
  {
    key: 'preserve_exif',
    value: true,
    type: 'boolean',
    group: 'upload',
    description: '是否保留EXIF信息',
    is_system: true,
  },
  {
    key: 'daily_upload_limit',
    value: 50,
    type: 'number',
    group: 'upload',
    description: '用户每日上传数量限制',
    is_system: true,
  },
  {
    key: 'client_max_concurrent_uploads',
    value: 3,
    type: 'number',
    group: 'upload',
    description: '客户端最大并发上传数',
    is_system: true,
  },
  /* 分片上传配置 */
  {
    key: 'chunked_upload_enabled',
    value: true,
    type: 'boolean',
    group: 'upload',
    description: '分片上传功能开关',
    is_system: true,
  },
  {
    key: 'chunked_threshold',
    value: 20,
    type: 'number',
    group: 'upload',
    description: '分片上传阈值(MB)',
    is_system: true,
  },
  {
    key: 'chunk_size',
    value: 2,
    type: 'number',
    group: 'upload',
    description: '分片大小(MB)',
    is_system: true,
  },
  {
    key: 'max_concurrency',
    value: 5,
    type: 'number',
    group: 'upload',
    description: '分片上传最大并发数',
    is_system: true,
  },
  {
    key: 'session_timeout',
    value: 24,
    type: 'number',
    group: 'upload',
    description: '分片上传会话超时时间(小时)',
    is_system: true,
  },
  {
    key: 'retry_count',
    value: 3,
    type: 'number',
    group: 'upload',
    description: '分片上传重试次数',
    is_system: true,
  },
  {
    key: 'cleanup_interval',
    value: 60,
    type: 'number',
    group: 'upload',
    description: '分片上传清理间隔(分钟)',
    is_system: true,
  },
  {
    key: 'user_allowed_storage_durations',
    value: ['1h', '3d', '7d', '30d', 'permanent'],
    type: 'array',
    group: 'upload',
    description: '已登录用户可选择的存储时长选项（permanent为内置选项）',
    is_system: true,
  },
  {
    key: 'user_default_storage_duration',
    value: 'permanent',
    type: 'string',
    group: 'upload',
    description: '已登录用户默认存储时长',
    is_system: true,
  },
]

/* 网站建设设置 */
export const constructionDefaults: Setting[] = [
  /* 公告相关配置 */
  {
    key: 'announcement_enabled',
    value: false,
    type: 'boolean',
    group: 'construction',
    description: '公告开启状态',
    is_system: true,
  },
  {
    key: 'announcement_title',
    value: '',
    type: 'string',
    group: 'construction',
    description: '公告标题',
    is_system: true,
  },
  {
    key: 'announcement_content',
    value: '',
    type: 'string',
    group: 'construction',
    description: '公告内容（支持HTML）',
    is_system: true,
  },
  {
    key: 'announcement_delay_before_show',
    value: 5000,
    type: 'number',
    group: 'construction',
    description: '公告显示延迟时间(毫秒)',
    is_system: true,
  },
  {
    key: 'announcement_reopen_delay',
    value: 300000,
    type: 'number',
    group: 'construction',
    description: '公告重新打开间隔时间(毫秒)',
    is_system: true,
  },
  {
    key: 'announcement_width',
    value: 600,
    type: 'number',
    group: 'construction',
    description: '公告弹窗宽度(像素)',
    is_system: true,
  },
  {
    key: 'announcement_icon',
    value: 'fas fa-bullhorn',
    type: 'string',
    group: 'construction',
    description: '公告图标',
    is_system: true,
  },
]

/* 网站装修设置 */
export const themeDefaults: Setting[] = [
  {
    key: 'site_mode',
    value: 'website',
    type: 'string',
    group: 'theme',
    description: '网站显示模式(固定为website:传统网站模式)',
    is_system: true,
  },
]

/* 访客控制设置 */
export const guestDefaults: Setting[] = [
  {
    key: 'enable_guest_upload',
    value: true,
    type: 'boolean',
    group: 'guest',
    description: '是否开放游客上传',
    is_system: true,
  },
  {
    key: 'guest_daily_limit',
    value: 10,
    type: 'number',
    group: 'guest',
    description: '游客每日上传次数限制',
    is_system: true,
  },
  {
    key: 'guest_default_access_level',
    value: 'public',
    type: 'string',
    group: 'guest',
    description: '默认访问级别',
    is_system: true,
  },
  {
    key: 'guest_allowed_storage_durations',
    value: ['3d', '7d', '30d'],
    type: 'array',
    group: 'guest',
    description: '允许的存储时长选项',
    is_system: true,
  },
  {
    key: 'guest_default_storage_duration',
    value: '7d',
    type: 'string',
    group: 'guest',
    description: '默认存储时长',
    is_system: true,
  },
  {
    key: 'guest_ip_daily_limit',
    value: 50,
    type: 'number',
    group: 'guest',
    description: 'IP每日限制，防止通过刷新指纹盗刷',
    is_system: true,
  },
]

/* 默认设置配置 - 统一导出 */
export const defaultSettings: Record<SettingGroup, Setting[]> = {
  website: websiteDefaults,
  website_info: websiteInfoDefaults,
  registration: registrationDefaults,
  security: securityDefaults,
  mail: mailDefaults,
  ai: aiDefaults,
  vector: vectorDefaults,
  upload: uploadDefaults,
  construction: constructionDefaults,
  theme: themeDefaults,
  guest: guestDefaults,
  appearance: [], // 预留外观设置
}

/* 翻译设置描述的辅助函数 */
function translateSettingDescription($t: TranslationFunction, setting: Setting): Setting {
  const { group, key } = setting
  try {
    return {
      ...setting,
      description: $t(`api.settingsDefaults.${group}.${key}`),
    }
  } catch {
    return setting // 如果翻译键不存在，返回原始设置
  }
}

/* 获取带翻译的默认设置（带翻译） */
export function getDefaultSettings($t: TranslationFunction): Record<SettingGroup, Setting[]> {
  return {
    website: getWebsiteDefaults($t),
    website_info: websiteInfoDefaults,
    registration: registrationDefaults.map((s) => translateSettingDescription($t, s)),
    security: securityDefaults.map((s) => translateSettingDescription($t, s)),
    mail: mailDefaults.map((s) => translateSettingDescription($t, s)),
    ai: aiDefaults.map((s) => translateSettingDescription($t, s)),
    vector: vectorDefaults.map((s) => translateSettingDescription($t, s)),
    upload: uploadDefaults.map((s) => translateSettingDescription($t, s)),
    construction: constructionDefaults.map((s) => translateSettingDescription($t, s)),
    theme: themeDefaults.map((s) => translateSettingDescription($t, s)),
    guest: guestDefaults.map((s) => translateSettingDescription($t, s)),
    appearance: [],
  }
}
