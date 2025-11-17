/**
 * 用户界面文本和消息常量
 *
 * 使用方式：
 * import { getCommonActions } from '@/constants/messages'
 * import { useTexts } from '@/composables/useTexts'
 *
 * const { $t } = useTexts()
 * const actions = getCommonActions($t)
 */

import type { Composer } from '@/composables/useTexts'

/* ===== 通用操作文本 ===== */
export const getCommonActions = ($t: Composer['$t']) =>
  ({
    CONFIRM: $t('actions.confirm'),
    CANCEL: $t('actions.cancel'),
    SAVE: $t('actions.save'),
    EDIT: $t('actions.edit'),
    DELETE: $t('actions.delete'),
    SUBMIT: $t('actions.submit'),
    RESET: $t('actions.reset'),
    SEARCH: $t('actions.search'),
    FILTER: $t('actions.filter'),
    UPLOAD: $t('actions.upload'),
    DOWNLOAD: $t('actions.download'),
    COPY: $t('actions.copy'),
    SHARE: $t('actions.share'),
    BACK: $t('actions.back'),
    CLOSE: $t('actions.close'),
    REFRESH: $t('actions.refresh'),
    RETRY: $t('actions.retry'),
    VIEW: $t('actions.view'),
    PREVIEW: $t('actions.preview'),
    SELECT: $t('actions.select'),
    CLEAR: $t('actions.clear'),
  }) as const

/* ===== 状态文本 ===== */
export const getStatusTexts = ($t: Composer['$t']) =>
  ({
    SUCCESS: $t('status.success'),
    FAILED: $t('status.failed'),
    ERROR: $t('status.error'),
    WARNING: $t('status.warning'),
    INFO: $t('status.info'),
    LOADING: $t('status.loading'),
    UPLOADING: $t('status.uploading'),
    DELETING: $t('status.deleting'),
    PROCESSING: $t('status.processing'),
    SAVING: $t('status.saving'),
    WAITING: $t('status.waiting'),
    COMPLETED: $t('status.completed'),
    PENDING: $t('status.pending'),
    CANCELLED: $t('status.cancelled'),
  }) as const

/* ===== 表单字段标签 ===== */
export const getFormLabels = ($t: Composer['$t']) =>
  ({
    USERNAME: $t('form.username'),
    PASSWORD: $t('form.password'),
    EMAIL: $t('form.email'),
    VERIFY_CODE: $t('form.verifyCode'),
    FILENAME: $t('form.filename'),
    DESCRIPTION: $t('form.description'),
    TITLE: $t('form.title'),
    TAGS: $t('form.tags'),
    CATEGORY: $t('form.category'),
    ACCESS_LEVEL: $t('form.accessLevel'),
    FOLDER: $t('form.folder'),
    SIZE: $t('form.size'),
    FORMAT: $t('form.format'),
    DATE: $t('form.date'),
    AUTHOR: $t('form.author'),
    REMARKS: $t('form.remarks'),
  }) as const

/* ===== 空状态提示 ===== */
export const getEmptyStateMessages = ($t: Composer['$t']) =>
  ({
    NO_DATA: $t('empty.noData'),
    NO_CONTENT: $t('empty.noContent'),
    NO_IMAGES: $t('empty.noImages'),
    NO_FILES: $t('empty.noFiles'),
    NO_RESULTS: $t('empty.noResults'),
    NO_MORE: $t('empty.noMore'),
    PLEASE_SELECT: $t('empty.pleaseSelect'),
    PLEASE_UPLOAD: $t('empty.pleaseUpload'),
    NO_PERMISSION: $t('empty.noPermission'),
    NETWORK_ERROR: $t('empty.networkError'),
  }) as const

/* ===== 成功消息 ===== */
export const getSuccessMessages = ($t: Composer['$t']) =>
  ({
    UPLOAD_SUCCESS: $t('messages.success.upload'),
    DELETE_SUCCESS: $t('messages.success.delete'),
    SAVE_SUCCESS: $t('messages.success.save'),
    UPDATE_SUCCESS: $t('messages.success.update'),
    COPY_SUCCESS: $t('messages.success.copy'),
    SHARE_SUCCESS: $t('messages.success.share'),
    LOGIN_SUCCESS: $t('messages.success.login'),
    LOGOUT_SUCCESS: $t('messages.success.logout'),
    REGISTER_SUCCESS: $t('messages.success.register'),
    PASSWORD_RESET: $t('messages.success.passwordReset'),
    EMAIL_SENT: $t('messages.success.emailSent'),
    OPERATION_SUCCESS: $t('messages.success.operation'),
  }) as const

/* ===== 错误消息 ===== */
export const getErrorMessages = ($t: Composer['$t']) =>
  ({
    UPLOAD_FAILED: $t('messages.error.upload'),
    DELETE_FAILED: $t('messages.error.delete'),
    SAVE_FAILED: $t('messages.error.save'),
    UPDATE_FAILED: $t('messages.error.update'),
    COPY_FAILED: $t('messages.error.copy'),
    SHARE_FAILED: $t('messages.error.share'),
    LOGIN_FAILED: $t('messages.error.login'),
    REGISTER_FAILED: $t('messages.error.register'),
    NETWORK_ERROR: $t('messages.error.network'),
    SERVER_ERROR: $t('messages.error.server'),
    PERMISSION_DENIED: $t('messages.error.permissionDenied'),
    SESSION_EXPIRED: $t('messages.error.sessionExpired'),
    FILE_TOO_LARGE: $t('messages.error.fileTooLarge'),
    INVALID_FORMAT: $t('messages.error.invalidFormat'),
    REQUIRED_FIELD: $t('messages.error.requiredField'),
    INVALID_EMAIL: $t('messages.error.invalidEmail'),
    PASSWORD_TOO_SHORT: $t('messages.error.passwordTooShort'),
    CONFIRM_PASSWORD_MISMATCH: $t('messages.error.confirmPasswordMismatch'),
  }) as const

/* ===== 确认对话框消息 ===== */
export const getConfirmMessages = ($t: Composer['$t']) =>
  ({
    DELETE_CONFIRM: $t('confirm.delete'),
    CLEAR_CONFIRM: $t('confirm.clear'),
    LOGOUT_CONFIRM: $t('confirm.logout'),
    RESET_CONFIRM: $t('confirm.reset'),
    CANCEL_CONFIRM: $t('confirm.cancel'),
    OVERWRITE_CONFIRM: $t('confirm.overwrite'),
    BULK_DELETE_CONFIRM: $t('confirm.bulkDelete'),
    PERMANENT_DELETE: $t('confirm.permanentDelete'),
  }) as const

/* ===== 加载状态消息 ===== */
export const getLoadingMessages = ($t: Composer['$t']) =>
  ({
    LOADING_DATA: $t('messages.loading.data'),
    LOADING_IMAGES: $t('messages.loading.images'),
    UPLOADING_FILES: $t('messages.loading.uploadingFiles'),
    PROCESSING_REQUEST: $t('messages.loading.processingRequest'),
    SAVING_CHANGES: $t('messages.loading.savingChanges'),
    DELETING_ITEMS: $t('messages.loading.deletingItems'),
    GENERATING_LINK: $t('messages.loading.generatingLink'),
    COPYING_LINK: $t('messages.loading.copyingLink'),
  }) as const

/* ===== 提示信息 ===== */
export const getTipMessages = ($t: Composer['$t']) =>
  ({
    SELECT_ALL: $t('tips.selectAll'),
    DESELECT_ALL: $t('tips.deselectAll'),
    DRAG_TO_UPLOAD: $t('tips.dragToUpload'),
    CLICK_TO_UPLOAD: $t('tips.clickToUpload'),
    SUPPORT_FORMATS: $t('tips.supportFormats'),
    MAX_FILE_SIZE: $t('tips.maxFileSize'),
    BATCH_OPERATION: $t('tips.batchOperation'),
    KEYBOARD_SHORTCUTS: $t('tips.keyboardShortcuts'),
    RIGHT_CLICK_MENU: $t('tips.rightClickMenu'),
  }) as const

/* ===== 页面标题 ===== */
export const getPageTitles = ($t: Composer['$t']) =>
  ({
    HOME: $t('navigation.home'),
    GALLERY: $t('navigation.explore'),
    UPLOAD: $t('navigation.upload'),
    SETTINGS: $t('navigation.settings'),
    PROFILE: $t('navigation.settings'),
    SHARE: $t('routes.shares'),
    DASHBOARD: $t('navigation.dashboard'),
    LOGIN: $t('routes.auth'),
    REGISTER: $t('routes.auth'),
    RESET_PASSWORD: $t('routes.resetPassword'),
    ADMIN: $t('routes.admin'),
    STATISTICS: $t('routes.adminDashboard'),
  }) as const
