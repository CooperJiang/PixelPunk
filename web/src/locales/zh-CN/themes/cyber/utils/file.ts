/**
 * Utils File 多语言配置 - 赛博主题
 */
export const file = {
  loadFailed: '数据加载失败',
  nsfwWarning: '此内容可能包含敏感信息',
  fileHash: {
    errors: {
      emptyResult: '数据读取异常',
      md5Failed: 'MD5计算异常: {error}',
      readFailedAfterRetries: '数据读取失败，已重试 {maxRetries} 次',
      readAborted: '数据读取已终止',
      chunkIndexOutOfRange: '分片索引越界',
      chunkCreationFailed: '无法创建数据分片: {error}',
      chunkReadFailed: '分片读取异常',
    },
  },
  watermark: {
    errors: {
      canvasContextFailed: '无法初始化2D渲染上下文',
      unsupportedWatermarkType: '不支持的水印类型',
      watermarkGenerationFailed: '水印生成异常',
      tempCanvasContextFailed: '无法创建临时渲染上下文',
      noWatermarkFile: '未提供水印数据',
      fileLoadFailed: '数据加载失败',
      imageLoadFailed: '数据加载失败: {url}',
      cannotLoadImage: '无法加载数据',
    },
    status: {
      completed: '✅ 水印处理完成 ({processedCount}/{totalFiles})',
      failed: '❌ 水印处理失败 ({failedCount}/{totalFiles})',
      partialSuccess: '⚠️ 水印处理部分成功 (成功: {processedCount}, 失败: {failedCount})',
    },
  },
  colorUtils: {
    errors: {
      invalidImageUrl: '无效的数据URL',
      cannotCreateCanvasContext: '无法创建渲染上下文',
      cannotGetImageData: '无法获取数据（可能存在访问限制）',
      imageLoadFailed: '数据加载失败，无法提取颜色',
      extractColorFailed: '颜色提取异常: {error}',
    },
  },
}
