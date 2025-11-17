/**
 * Upload Store 多语言配置 - 赛博主题
 */
export const upload = {
  control: {
    paused: '已挂起',
    preparingResume: '准备恢复',
    userCancelled: '用户中止',
    cancelled: '已中止',
    preparingRetry: '准备重连',
  },
  queue: {
    preparingAnalysis: '预备扫描...',
    analyzing: '扫描数据...',
    checkingInstant: '检测瞬传可能性...',
    calculatingMD5: '计算数据签名... {progress}%',
    checkingServerDuplicate: '检查服务端数据副本...',
    preparingUpload: '执行瞬传...',
    instantComplete: '瞬传完毕',
    analysisComplete: '扫描完毕',
    waitingUpload: '等待传输',
  },
  upload: {
    analyzingFile: '扫描数据...',
    analyzingProgress: '扫描数据 {progress}%',
    instantProgress: '瞬传中 {progress}%',
    instantComplete: '瞬传完毕',
    uploading: '传输中...',
    mergingChunks: '合并数据块...',
    uploadComplete: '传输完毕',
    uploadFailed: '传输失败',
    instantFailed: '瞬传失败',
    preparingRetry: '准备重连',
    retrying: '重连中',
    postProcessing: '后期处理中...',
    analysisError: '扫描错误',
    remainingTimeMinutes: '{minutes}分{seconds}秒',
    remainingTimeSeconds: '{seconds}秒',
    uploadingProgress: '传输中 {progress}% - {speed} - 剩余{remaining}',
    uploadFailedWithError: '传输失败: {error}',
  },
}
