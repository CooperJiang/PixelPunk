/**
 * 内容审核模块常量
 */

import type { Composer } from '@/composables/useTexts'

/* 审核动作 */
export const REVIEW_ACTIONS = {
  APPROVE: 'approve',
  REJECT: 'reject',
} as const

export type ReviewAction = (typeof REVIEW_ACTIONS)[keyof typeof REVIEW_ACTIONS]

export const DELETE_TYPES = {
  SOFT: 'soft',
  HARD: 'hard',
} as const

export type DeleteType = (typeof DELETE_TYPES)[keyof typeof DELETE_TYPES]

/* 文件状态 */
export const FILE_STATUS = {
  ACTIVE: 'active',
  DELETED: 'deleted',
  PENDING: 'pending',
} as const

export type FileStatus = (typeof FILE_STATUS)[keyof typeof FILE_STATUS]

/* 审核状态选项 */
export const getReviewActionOptions = ($t: Composer['$t']) => [
  { label: $t('constants.contentReview.actions.all'), value: '' },
  { label: $t('constants.contentReview.actions.approve'), value: REVIEW_ACTIONS.APPROVE },
  { label: $t('constants.contentReview.actions.reject'), value: REVIEW_ACTIONS.REJECT },
]

/* 排序选项 */
export const getReviewSortOptions = ($t: Composer['$t']) => [
  { label: $t('constants.contentReview.sort.newest'), value: 'newest' },
  { label: $t('constants.contentReview.sort.oldest'), value: 'oldest' },
  { label: $t('constants.contentReview.sort.size'), value: 'size' },
  { label: $t('constants.contentReview.sort.nsfw'), value: 'nsfw' },
]

/* NSFW评分等级 */
export const NSFW_SCORE_LEVELS = {
  SAFE: 0.3,
  LOW: 0.6,
  MEDIUM: 0.8,
  HIGH: 1.0,
} as const

/* NSFW评分显示配置 */
export const getNsfwScoreConfig = ($t: Composer['$t']) => ({
  SAFE: {
    threshold: NSFW_SCORE_LEVELS.SAFE,
    class: 'nsfw-safe',
    color: 'text-green-400',
    label: $t('constants.contentReview.nsfwLevels.safe'),
  },
  LOW: {
    threshold: NSFW_SCORE_LEVELS.LOW,
    class: 'nsfw-low',
    color: 'text-yellow-400',
    label: $t('constants.contentReview.nsfwLevels.low'),
  },
  MEDIUM: {
    threshold: NSFW_SCORE_LEVELS.MEDIUM,
    class: 'nsfw-medium',
    color: 'text-orange-400',
    label: $t('constants.contentReview.nsfwLevels.medium'),
  },
  HIGH: {
    threshold: NSFW_SCORE_LEVELS.HIGH,
    class: 'nsfw-high',
    color: 'text-red-400',
    label: $t('constants.contentReview.nsfwLevels.high'),
  },
})

/* 分页配置 */
export const REVIEW_PAGINATION = {
  PAGE_SIZE: 20,
  PAGE_SIZE_OPTIONS: [10, 20, 50, 100],
  MAX_PAGE_SIZE: 100,
} as const

/* 表格列配置 */
export const getReviewTableColumns = ($t: Composer['$t']) => ({
  FILE: {
    key: 'file',
    title: $t('constants.contentReview.columns.file'),
    width: 80,
    align: 'center' as const,
  },
  FILE_INFO: {
    key: 'fileInfo',
    title: $t('constants.contentReview.columns.fileInfo'),
    width: 200,
    ellipsis: true,
  },
  ACTION: {
    key: 'action',
    title: $t('constants.contentReview.columns.action'),
    width: 120,
    align: 'center' as const,
    sortable: true,
  },
  SOFT_DELETE: {
    key: 'softDelete',
    title: $t('constants.contentReview.columns.softDelete'),
    width: 120,
    align: 'center' as const,
    sortable: true,
  },
  AUDITOR: {
    key: 'auditor',
    title: $t('constants.contentReview.columns.auditor'),
    width: 120,
    sortable: true,
  },
  AUDIT_TIME: {
    key: 'auditTime',
    title: $t('constants.contentReview.columns.auditTime'),
    width: 140,
    sortable: true,
  },
  REASON: {
    key: 'reason',
    title: $t('constants.contentReview.columns.reason'),
    minWidth: 150,
    ellipsis: true,
  },
  NSFW: {
    key: 'nsfw',
    title: $t('constants.contentReview.columns.nsfw'),
    width: 100,
    align: 'center' as const,
    sortable: true,
  },
  OPERATIONS: {
    key: 'operations',
    title: $t('constants.contentReview.columns.operations'),
    width: 120,
    align: 'center' as const,
    fixed: 'right' as const,
  },
})

/* 默认审核原因 */
export const getDefaultReviewReasons = ($t: Composer['$t']) => ({
  APPROVE: $t('constants.contentReview.defaultReasons.approve'),
  REJECT: $t('constants.contentReview.defaultReasons.reject'),
  BATCH_APPROVE: $t('constants.contentReview.defaultReasons.batchApprove'),
  BATCH_REJECT: $t('constants.contentReview.defaultReasons.batchReject'),
})

/* 文本截断长度 */
export const TEXT_TRUNCATE = {
  REASON: 40,
  FILENAME: 20,
  USERNAME: 15,
} as const

/* 文件网格配置 */
export const FILE_GRID_CONFIG = {
  MIN_CARD_WIDTH: {
    DESKTOP: 200,
    TABLET: 180,
    MOBILE: 160,
    SMALL: 140,
    TINY: 120,
  },
  ASPECT_RATIO: '4/3',
  GAP: {
    DESKTOP: '1rem',
    TABLET: '0.75rem',
    MOBILE: '0.5rem',
  },
} as const

/* 审核提示信息 */
export const getReviewMessages = ($t: Composer['$t']) => ({
  APPROVE_SUCCESS: $t('constants.contentReview.messages.approveSuccess'),
  REJECT_SUCCESS: $t('constants.contentReview.messages.rejectSuccess'),
  BATCH_APPROVE_SUCCESS: (count: number) => $t('constants.contentReview.messages.batchApproveSuccess', { count }),
  BATCH_REJECT_SUCCESS: (count: number) => $t('constants.contentReview.messages.batchRejectSuccess', { count }),
  RESTORE_SUCCESS: $t('constants.contentReview.messages.restoreSuccess'),
  HARD_DELETE_SUCCESS: $t('constants.contentReview.messages.hardDeleteSuccess'),
  BATCH_RESTORE_SUCCESS: (success: number, fail: number) =>
    fail > 0
      ? $t('constants.contentReview.messages.batchRestoreWithFail', { success, fail })
      : $t('constants.contentReview.messages.batchRestoreSuccess', { success }),
  BATCH_HARD_DELETE_SUCCESS: (success: number, fail: number) =>
    fail > 0
      ? $t('constants.contentReview.messages.batchHardDeleteWithFail', { success, fail })
      : $t('constants.contentReview.messages.batchHardDeleteSuccess', { success }),
  LOAD_ERROR: $t('constants.contentReview.messages.loadError'),
  OPERATION_ERROR: $t('constants.contentReview.messages.operationError'),
  NO_SELECTION: $t('constants.contentReview.messages.noSelection'),
  QUEUE_EMPTY: $t('constants.contentReview.messages.queueEmpty'),
  QUEUE_EMPTY_DESC: $t('constants.contentReview.messages.queueEmptyDesc'),
  LOGS_EMPTY: $t('constants.contentReview.messages.logsEmpty'),
  LOGS_EMPTY_DESC: $t('constants.contentReview.messages.logsEmptyDesc'),
})

/* 审核警告信息 */
export const getReviewWarnings = ($t: Composer['$t']) => ({
  HARD_DELETE: $t('constants.contentReview.warnings.hardDelete'),
  HARD_DELETE_DESC: $t('constants.contentReview.warnings.hardDeleteDesc'),
  BATCH_OPERATION: $t('constants.contentReview.warnings.batchOperation'),
})

/* 审核说明信息 */
export const getReviewDescriptions = ($t: Composer['$t']) => ({
  APPROVE_EFFECTS: $t('constants.contentReview.descriptions.approveEffects') as any,
  RESTORE_EFFECTS: $t('constants.contentReview.descriptions.restoreEffects') as any,
})
