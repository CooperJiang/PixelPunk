import type { Composer } from '@/composables/useTexts'

/**
 * 导航选项常量
 * 用于菜单、导航、页面配置等
 */

/* ==================== 管理员菜单项 ==================== */

/**
 * 获取管理员菜单项
 */
export const getAdminMenuItems = ($t: Composer['$t']) => [
  /* 核心概览 */
  { path: '/admin/dashboard', label: $t('constants.navigation.admin.dashboard'), icon: 'fa-tachometer-alt' },

  /* 内容管理 */
  { path: '/admin/files', label: $t('constants.navigation.admin.files'), icon: 'fa-file' },
  { path: '/admin/tags', label: $t('constants.navigation.admin.tags'), icon: 'fa-tags' },
  { path: '/admin/categories', label: $t('constants.navigation.admin.categories'), icon: 'fa-folder-open' },
  { path: '/admin/shares', label: $t('constants.navigation.admin.shares'), icon: 'fa-share-alt' },

  /* AI智能功能 */
  { path: '/admin/ai', label: $t('constants.navigation.admin.ai'), icon: 'fa-brain' },
  { path: '/admin/content-review', label: $t('constants.navigation.admin.contentReview'), icon: 'fa-gavel' },
  { path: '/admin/tagging', label: $t('constants.navigation.admin.tagging'), icon: 'fa-robot' },
  { path: '/admin/vectors', label: $t('constants.navigation.admin.vectors'), icon: 'fa-vector-square' },

  /* 用户与渠道 */
  { path: '/admin/users', label: $t('constants.navigation.admin.users'), icon: 'fa-users' },
  { path: '/admin/channels', label: $t('constants.navigation.admin.channels'), icon: 'fa-server' },

  /* 系统配置 */
  { path: '/admin/construction', label: $t('constants.navigation.admin.construction'), icon: 'fa-hammer' },
  { path: '/admin/announcements', label: $t('constants.navigation.admin.announcements'), icon: 'fa-bullhorn' },
  { path: '/admin/settings', label: $t('constants.navigation.admin.settings'), icon: 'fa-cog' },
]

export type AdminMenuItemI18n = ReturnType<typeof getAdminMenuItems>[number]

/* ==================== 文档导航项 ==================== */

/**
 * 获取文档导航项
 */
export const getDocsNavItems = ($t: Composer['$t']) => [
  { id: 'api-overview', label: $t('constants.navigation.docs.overview'), iconClass: 'fas fa-book' },
  { id: 'api-authentication', label: $t('constants.navigation.docs.authentication'), iconClass: 'fas fa-key' },
  { id: 'api-upload', label: $t('constants.navigation.docs.upload'), iconClass: 'fas fa-upload' },
  { id: 'api-limits', label: $t('constants.navigation.docs.limits'), iconClass: 'fas fa-exclamation-circle' },
  { id: 'api-examples', label: $t('constants.navigation.docs.examples'), iconClass: 'fas fa-code' },
  { id: 'api-tester', label: $t('constants.navigation.docs.tester'), iconClass: 'fas fa-flask' },
  { id: 'api-faq', label: $t('constants.navigation.docs.faq'), iconClass: 'fas fa-question-circle' },
]

export type DocsNavItemI18n = ReturnType<typeof getDocsNavItems>[number]

/* ==================== 代码标签选项（编程语言名称，不需要翻译） ==================== */

export const CODE_TAB_OPTIONS = [
  { label: 'cURL', value: 'curl' },
  { label: 'JavaScript', value: 'javascript' },
  { label: 'Node.js', value: 'nodejs' },
  { label: 'Python', value: 'python' },
  { label: 'Java', value: 'java' },
  { label: 'Go', value: 'go' },
] as const

export type CodeTabOption = (typeof CODE_TAB_OPTIONS)[number]

/* ==================== 主导航菜单项 ==================== */

/**
 * 获取主导航菜单项
 */
export const getMainNavItems = ($t: Composer['$t']) => [
  { path: '/', label: $t('constants.navigation.main.home'), icon: 'fa-home' },
  { path: '/explore', label: $t('constants.navigation.main.explore'), icon: 'fa-images' },
  { path: '/upload', label: $t('constants.navigation.main.upload'), icon: 'fa-upload' },
  { path: '/random', label: $t('constants.navigation.main.random'), icon: 'fa-random' },
  { path: '/docs', label: $t('constants.navigation.main.docs'), icon: 'fa-code' },
]

export type MainNavItemI18n = ReturnType<typeof getMainNavItems>[number]

/* ==================== 用户菜单项 ==================== */

/**
 * 获取用户菜单项
 */
export const getUserMenuItems = ($t: Composer['$t']) => [
  { path: '/dashboard', label: $t('constants.navigation.user.dashboard'), icon: 'fa-tachometer-alt' },
  { path: '/my-files', label: $t('constants.navigation.user.myFiles'), icon: 'fa-file' },
  { path: '/folders', label: $t('constants.navigation.user.folders'), icon: 'fa-folder' },
  { path: '/tag-manage', label: $t('constants.navigation.user.tagManage'), icon: 'fa-tags' },
  { path: '/category-manage', label: $t('constants.navigation.user.categoryManage'), icon: 'fa-folder-open' },
  { path: '/shares', label: $t('constants.navigation.user.shares'), icon: 'fa-share-alt' },
  { path: '/settings', label: $t('constants.navigation.user.settings'), icon: 'fa-cog' },
]

export type UserMenuItemI18n = ReturnType<typeof getUserMenuItems>[number]

/* ==================== 设置导航项 ==================== */

/**
 * 获取设置导航项
 */
export const getSettingsNavItems = ($t: Composer['$t']) => [
  { key: 'profile', label: $t('constants.navigation.settings.profile'), icon: 'fa-user' },
  { key: 'security', label: $t('constants.navigation.settings.security'), icon: 'fa-shield-alt' },
  { key: 'apikey', label: $t('constants.navigation.settings.apikey'), icon: 'fa-key' },
  { key: 'storage', label: $t('constants.navigation.settings.storage'), icon: 'fa-database' },
]

export type SettingsNavItemI18n = ReturnType<typeof getSettingsNavItems>[number]

/* ==================== 管理员设置导航项 ==================== */

/**
 * 获取管理员设置导航项
 */
export const getAdminSettingsNavItems = ($t: Composer['$t']) => [
  { key: 'website', label: $t('constants.navigation.adminSettings.website'), icon: 'fa-globe' },
  { key: 'upload', label: $t('constants.navigation.adminSettings.upload'), icon: 'fa-upload' },
  { key: 'security', label: $t('constants.navigation.adminSettings.security'), icon: 'fa-shield-alt' },
  { key: 'registration', label: $t('constants.navigation.adminSettings.registration'), icon: 'fa-user-plus' },
  { key: 'mail', label: $t('constants.navigation.adminSettings.mail'), icon: 'fa-envelope' },
  { key: 'ai', label: $t('constants.navigation.adminSettings.ai'), icon: 'fa-robot' },
  { key: 'vector', label: $t('constants.navigation.adminSettings.vector'), icon: 'fa-vector-square' },
]

export type AdminSettingsNavItemI18n = ReturnType<typeof getAdminSettingsNavItems>[number]
