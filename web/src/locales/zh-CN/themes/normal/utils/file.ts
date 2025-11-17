/**
 * Utils File 多语言配置 - 常规主题
 */
export const file = {
  loadFailed: '文件加载失败',
  nsfwWarning: '此内容可能包含敏感信息',
  fileHash: {
    errors: {
      emptyResult: '读取结果为空',
      md5Failed: 'MD5计算失败: {error}',
      readFailedAfterRetries: '文件读取失败，已重试 {maxRetries} 次',
      readAborted: '文件读取被中止',
      chunkIndexOutOfRange: '分片索引超出文件大小',
      chunkCreationFailed: '无法创建文件分片: {error}',
      chunkReadFailed: '分片读取失败',
    },
  },
  watermark: {
    errors: {
      canvasContextFailed: '无法创建2D渲染上下文',
      unsupportedWatermarkType: '不支持的水印类型',
      watermarkGenerationFailed: '生成水印失败',
      tempCanvasContextFailed: '无法创建临时2D渲染上下文',
      noWatermarkFile: '未提供水印文件',
      fileLoadFailed: '文件加载失败',
      imageLoadFailed: '文件加载失败: {url}',
      cannotLoadImage: '无法加载文件',
    },
    status: {
      completed: '✅ 水印处理完成 ({processedCount}/{totalFiles})',
      failed: '❌ 水印处理失败 ({failedCount}/{totalFiles})',
      partialSuccess: '⚠️ 水印处理部分成功 (成功: {processedCount}, 失败: {failedCount})',
    },
  },
  colorUtils: {
    errors: {
      invalidImageUrl: '无效的文件URL',
      cannotCreateCanvasContext: '无法创建canvas上下文',
      cannotGetImageData: '无法获取文件数据（可能是跨域问题）',
      imageLoadFailed: '加载文件失败，无法提取颜色',
      extractColorFailed: '提取颜色失败: {error}',
    },
  },
}
