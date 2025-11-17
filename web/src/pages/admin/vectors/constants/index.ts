import { useTexts } from '@/composables/useTexts'

export const VECTOR_STATUS = {
  PENDING: 'pending',
  RUNNING: 'running',
  COMPLETED: 'completed',
  FAILED: 'failed',
  RESET: 'reset',
} as const

export const getStatusTextMap = (): Record<string, string> => {
  const { $t } = useTexts()
  return {
    [VECTOR_STATUS.PENDING]: $t('admin.vectors.status.pending'),
    [VECTOR_STATUS.RUNNING]: $t('admin.vectors.status.running'),
    [VECTOR_STATUS.COMPLETED]: $t('admin.vectors.status.completed'),
    [VECTOR_STATUS.FAILED]: $t('admin.vectors.status.failed'),
    [VECTOR_STATUS.RESET]: $t('admin.vectors.status.reset'),
  }
}

/* 保持向后兼容 */
export const STATUS_TEXT_MAP: Record<string, string> = getStatusTextMap()

export const STATUS_ICON_MAP: Record<string, string> = {
  [VECTOR_STATUS.PENDING]: 'fas fa-clock',
  [VECTOR_STATUS.RUNNING]: 'fas fa-spinner fa-spin',
  [VECTOR_STATUS.COMPLETED]: 'fas fa-check-circle',
  [VECTOR_STATUS.FAILED]: 'fas fa-exclamation-circle',
  [VECTOR_STATUS.RESET]: 'fas fa-redo',
}

export const BATCH_ACTIONS = {
  RETRY: 'retry',
  RESET: 'reset',
  DELETE: 'delete',
} as const

export const getBatchActionTextMap = (): Record<string, string> => {
  const { $t } = useTexts()
  return {
    [BATCH_ACTIONS.RETRY]: $t('admin.vectors.actions.retry'),
    [BATCH_ACTIONS.RESET]: $t('admin.vectors.actions.reset'),
    [BATCH_ACTIONS.DELETE]: $t('admin.vectors.actions.delete'),
  }
}

export const BATCH_ACTION_TEXT_MAP: Record<string, string> = getBatchActionTextMap()

export const getCollectionStatusTextMap = (): Record<string, string> => {
  const { $t } = useTexts()
  return {
    healthy: $t('admin.vectors.collectionStatus.healthy'),
    empty_normal: $t('admin.vectors.collectionStatus.emptyNormal'),
    empty_but_should_have_data: $t('admin.vectors.collectionStatus.dataMissing'),
    partially_synced: $t('admin.vectors.collectionStatus.partiallySynced'),
    needs_resync: $t('admin.vectors.collectionStatus.needsResync'),
    engine_not_initialized: $t('admin.vectors.collectionStatus.engineNotInitialized'),
    qdrant_connection_failed: $t('admin.vectors.collectionStatus.qdrantConnectionFailed'),
    error: $t('admin.vectors.collectionStatus.error'),
    unknown: $t('admin.vectors.collectionStatus.unknown'),
  }
}

export const COLLECTION_STATUS_TEXT_MAP: Record<string, string> = getCollectionStatusTextMap()

export const getModelDescriptions = (): Record<string, string> => {
  const { $t } = useTexts()
  return {
    'text-embedding-3-small': $t('admin.vectors.modelDesc.small'),
    'text-embedding-3-large': $t('admin.vectors.modelDesc.large'),
    'text-embedding-ada-002': $t('admin.vectors.modelDesc.ada'),
  }
}

export const MODEL_DESCRIPTIONS: Record<string, string> = getModelDescriptions()

export const REFRESH_INTERVALS = {
  STATS_FAST: 2000, // 快速更新统计信息（有活跃任务时）
  STATS_SLOW: 30000, // 慢速更新统计信息（无活跃任务时）
  DATA_DEBOUNCE: 300, // 数据加载防抖延迟
} as const

export const PAGINATION_CONFIG = {
  DEFAULT_PAGE_SIZE: 10,
  LOG_PAGE_SIZE: 20,
  MAX_PAGE_SIZE: 100,
} as const
