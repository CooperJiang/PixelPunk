import type {
  VectorItem,
  VectorStats,
  QdrantRealStats,
  VectorLogItem,
  PaginationInfo,
  VectorListParams,
} from '@/api/admin/vectors'

export type { VectorItem, VectorStats, QdrantRealStats, VectorLogItem, PaginationInfo, VectorListParams }

export interface ConfirmDialogData {
  title: string
  message: string[]
  type: 'primary' | 'danger' | 'warning'
  confirmText: string
  action: string
}

export interface DeleteConfirmData {
  count: number
  expectedText: string
  fileIds: string[]
}

export interface BatchConfirmData {
  action: 'reset' | 'retry' | 'delete'
  message: string
  fileIds: string[]
}

export interface RegenerateConfirmData {
  message: string[]
  force?: boolean
}

export interface RetryFailedConfirmData {
  message: string
  count: number
}

export interface DialogStates {
  showDeleteConfirm: boolean
  showBatchConfirm: boolean
  showRegenerateConfirm: boolean
  showRetryFailedConfirm: boolean
}

export interface ProblemStats {
  failed: number
  stuck: number
  missing: number
  unknown: number
  total: number
}

export interface VectorPageState {
  isLoading: boolean
  isProcessing: boolean
  showFilter: boolean
  showLogs: boolean
  selectedFileId: string
}

export interface ConcurrencyState {
  isLoadingData: boolean
  isUpdatingStats: boolean
  requestAbortController: AbortController | null
}

export interface TimerState {
  refreshInterval: number | null
  statsRefreshInterval: number | null
}
