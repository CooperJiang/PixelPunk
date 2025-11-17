/**
 * 初始化相关操作

 */
import { computed } from 'vue'
import { useChunkedUpload } from '@/composables/useChunkedUpload'
import { useUploadConfig } from '@/composables/useUploadConfig'
import { chunkedUploadInstance, isChunkedReady, globalOptions } from '../state'
import type { GlobalUploadOptions } from '../types'

export function initChunkedUpload() {
  if (chunkedUploadInstance.value) {
    return
  }

  const uploadConfig = useUploadConfig()
  chunkedUploadInstance.value = useChunkedUpload({
    chunkSize: uploadConfig.getChunkSize(),
    concurrency: uploadConfig.getChunkedMaxConcurrency(),
    retryCount: uploadConfig.getRetryCount(),
    folderId: computed(() => globalOptions.value.folderId),
    accessLevel: computed(() => globalOptions.value.accessLevel),
    optimize: computed(() => globalOptions.value.optimize),
    watermarkEnabled: computed(() => globalOptions.value.watermarkEnabled),
    watermarkConfig: computed(() => globalOptions.value.watermarkConfig),
  })

  isChunkedReady.value = true
}

export function setGlobalOptions(options: Partial<GlobalUploadOptions>) {
  globalOptions.value = { ...globalOptions.value, ...options }
}

export function classifyFile(file: File): boolean {
  const uploadConfig = useUploadConfig()
  const shouldUseChunked = uploadConfig.shouldUseChunkedUpload(file.size)
  return shouldUseChunked
}

export async function initialize() {
  const uploadConfig = useUploadConfig()

  if (!uploadConfig.isConfigLoaded.value) {
    await uploadConfig.loadConfig()
  }
  initChunkedUpload()
}
