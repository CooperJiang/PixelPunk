/**
 * 上传Store主入口
 * 使用Pinia defineStore统一导出

 */
import { defineStore } from 'pinia'
import { computed } from 'vue'
import { useUploadConfig } from '@/composables/useUploadConfig'
import { useTexts } from '@/composables/useTexts'

/* 导出类型 */
export type { UploadStatus, UploadItem, GlobalUploadOptions } from './types'

import { globalOptions, runningUploads, $t as stateT } from './state'

import { allUploads, globalProgress, globalSpeed, isUploading, statistics, maxConcurrentUploads } from './getters'

/* 操作方法 */
import { initialize, setGlobalOptions } from './actions/init'
import { addFiles, removeUploadItem } from './actions/queue'
import { startUpload, pauseUpload, resumeUpload, cancelUpload, retryUpload } from './actions/control'
import { clearQueue, clearAllSessions, cleanup } from './actions/cleanup'

export const useUploadStore = defineStore('upload', () => {
  const uploadConfig = useUploadConfig()
  const { $t } = useTexts()

  // 初始化翻译函数
  stateT.value = $t

  initialize()

  return {
    allUploads,
    globalProgress,
    globalSpeed,
    isUploading,
    statistics,
    maxConcurrentUploads,
    globalOptions,
    runningUploads: computed(() => runningUploads.value),

    setGlobalOptions,
    addFiles,
    startUpload,
    pauseUpload,
    resumeUpload,
    cancelUpload,
    retryUpload,
    removeUploadItem,
    clearQueue,
    clearAllSessions,
    cleanup,

    formatFileSize: uploadConfig.formatFileSize,
    isValidFileType: (file: File) => uploadConfig.isAllowedFileType(file.type),
  }
})
