/**
 * Utils File i18n - Cyber Theme
 */
export const file = {
  loadFailed: 'データ読み込み失敗',
  nsfwWarning: 'このコンテンツには機密情報が含まれている可能性があります',
  fileHash: {
    errors: {
      emptyResult: 'データ読み取り例外',
      md5Failed: 'MD5計算例外: {error}',
      readFailedAfterRetries: 'データ読み取り失敗、{maxRetries}回再試行しました',
      readAborted: 'データ読み取り終了',
      chunkIndexOutOfRange: 'チャンクインデックスが範囲外',
      chunkCreationFailed: 'データチャンクを作成できません: {error}',
      chunkReadFailed: 'チャンク読み取り例外',
    },
  },
  watermark: {
    errors: {
      canvasContextFailed: '2Dレンダリングコンテキストを初期化できません',
      unsupportedWatermarkType: 'サポートされていない透かしタイプ',
      watermarkGenerationFailed: '透かし生成例外',
      tempCanvasContextFailed: '一時レンダリングコンテキストを作成できません',
      noWatermarkFile: '透かしデータが提供されていません',
      fileLoadFailed: 'データ読み込み失敗',
      imageLoadFailed: 'データ読み込み失敗: {url}',
      cannotLoadImage: 'データを読み込めません',
    },
    status: {
      completed: '✅ 透かし処理完了（{processedCount}/{totalFiles}）',
      failed: '❌ 透かし処理失敗（{failedCount}/{totalFiles}）',
      partialSuccess: '⚠️ 透かし処理部分成功（成功: {processedCount}、失敗: {failedCount}）',
    },
  },
  colorUtils: {
    errors: {
      invalidImageUrl: '無効なデータURL',
      cannotCreateCanvasContext: 'レンダリングコンテキストを作成できません',
      cannotGetImageData: 'データを取得できません（アクセス制限の可能性）',
      imageLoadFailed: 'データ読み込み失敗、色を抽出できません',
      extractColorFailed: '色抽出例外: {error}',
    },
  },
}
