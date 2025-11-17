/**
 * Utils File i18n configuration - Normal theme
 */
export const file = {
  loadFailed: 'File load failed',
  nsfwWarning: 'This content may contain sensitive material',
  fileHash: {
    errors: {
      emptyResult: 'Read result is empty',
      md5Failed: 'MD5 calculation failed: {error}',
      readFailedAfterRetries: 'File read failed, retried {maxRetries} times',
      readAborted: 'File read aborted',
      chunkIndexOutOfRange: 'Chunk index exceeds file size',
      chunkCreationFailed: 'Cannot create file chunk: {error}',
      chunkReadFailed: 'Chunk read failed',
    },
  },
  watermark: {
    errors: {
      canvasContextFailed: 'Cannot create 2D rendering context',
      unsupportedWatermarkType: 'Unsupported watermark type',
      watermarkGenerationFailed: 'Watermark generation failed',
      tempCanvasContextFailed: 'Cannot create temporary 2D rendering context',
      noWatermarkFile: 'Watermark file not provided',
      fileLoadFailed: 'File load failed',
      imageLoadFailed: 'File load failed: {url}',
      cannotLoadImage: 'Cannot load file',
    },
    status: {
      completed: '✅ Watermark processing completed ({processedCount}/{totalFiles})',
      failed: '❌ Watermark processing failed ({failedCount}/{totalFiles})',
      partialSuccess: '⚠️ Watermark processing partially successful (Success: {processedCount}, Failed: {failedCount})',
    },
  },
  colorUtils: {
    errors: {
      invalidImageUrl: 'Invalid file URL',
      cannotCreateCanvasContext: 'Cannot create canvas context',
      cannotGetImageData: 'Cannot get file data (may be CORS issue)',
      imageLoadFailed: 'Failed to load file, cannot extract color',
      extractColorFailed: 'Color extraction failed: {error}',
    },
  },
}
