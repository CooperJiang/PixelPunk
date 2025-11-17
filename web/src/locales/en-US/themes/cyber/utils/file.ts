/**
 * Utils File i18n - Cyber Theme
 */
export const file = {
  loadFailed: 'Data load failed',
  nsfwWarning: 'This content may contain sensitive material',
  fileHash: {
    errors: {
      emptyResult: 'Data read exception',
      md5Failed: 'MD5 calculation exception: {error}',
      readFailedAfterRetries: 'Data read failed, retried {maxRetries} times',
      readAborted: 'Data read terminated',
      chunkIndexOutOfRange: 'Chunk index out of bounds',
      chunkCreationFailed: 'Cannot create data chunk: {error}',
      chunkReadFailed: 'Chunk read exception',
    },
  },
  watermark: {
    errors: {
      canvasContextFailed: 'Cannot initialize 2D rendering context',
      unsupportedWatermarkType: 'Unsupported watermark type',
      watermarkGenerationFailed: 'Watermark generation exception',
      tempCanvasContextFailed: 'Cannot create temporary rendering context',
      noWatermarkFile: 'No watermark data provided',
      fileLoadFailed: 'Data load failed',
      imageLoadFailed: 'Data load failed: {url}',
      cannotLoadImage: 'Cannot load data',
    },
    status: {
      completed: '✅ Watermark processing complete ({processedCount}/{totalFiles})',
      failed: '❌ Watermark processing failed ({failedCount}/{totalFiles})',
      partialSuccess: '⚠️ Watermark processing partial success (Success: {processedCount}, Failed: {failedCount})',
    },
  },
  colorUtils: {
    errors: {
      invalidImageUrl: 'Invalid data URL',
      cannotCreateCanvasContext: 'Cannot create rendering context',
      cannotGetImageData: 'Cannot get data (possible access restriction)',
      imageLoadFailed: 'Data load failed, cannot extract color',
      extractColorFailed: 'Color extraction exception: {error}',
    },
  },
}
