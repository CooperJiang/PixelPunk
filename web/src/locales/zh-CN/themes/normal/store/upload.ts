/**
 * Upload Store 多语言配置 - 常规主题
 */
export const upload = {
  control: {
    paused: '已暂停',
    preparingResume: '准备继续',
    userCancelled: '用户取消',
    cancelled: '已取消',
    preparingRetry: '准备重试',
  },
  queue: {
    preparingAnalysis: '准备分析...',
    analyzing: '分析文件...',
    checkingInstant: '检查秒传可能性...',
    calculatingMD5: '计算文件MD5... {progress}%',
    checkingServerDuplicate: '检查服务器重复文件...',
    preparingUpload: '执行秒传...',
    instantComplete: '秒传完成',
    analysisComplete: '分析完成',
    waitingUpload: '准备上传',
  },
  upload: {
    analyzingFile: '分析文件...',
    analyzingProgress: '分析文件 {progress}%',
    instantProgress: '秒传中 {progress}%',
    instantComplete: '秒传完成',
    uploading: '上传中...',
    mergingChunks: '合并分片中...',
    uploadComplete: '上传完成',
    uploadFailed: '上传失败',
    instantFailed: '秒传失败',
    preparingRetry: '准备重试',
    retrying: '重试中',
    postProcessing: '后期处理中...',
    analysisError: '分析错误',
    remainingTimeMinutes: '{minutes}分{seconds}秒',
    remainingTimeSeconds: '{seconds}秒',
    uploadingProgress: '上传中 {progress}% - {speed} - 剩余{remaining}',
    uploadFailedWithError: '上传失败: {error}',
  },
}
