/**
 * Utils File i18n configuration - Normal theme
 */
export const file = {
  loadFailed: 'ファイル読み込み失敗',
  nsfwWarning: 'このコンテンツには機密情報が含まれている可能性があります',
  fileHash: {
    errors: {
      emptyResult: '読み取り結果が空です',
      md5Failed: 'MD5計算失敗：{error}',
      readFailedAfterRetries: 'ファイル読み取り失敗、{maxRetries}回再試行しました',
      readAborted: 'ファイル読み取りが中止されました',
      chunkIndexOutOfRange: 'チャンクインデックスがファイルサイズを超えています',
      chunkCreationFailed: 'ファイルチャンクを作成できません：{error}',
      chunkReadFailed: 'チャンク読み取り失敗',
    },
  },
  watermark: {
    errors: {
      canvasContextFailed: '2Dレンダリングコンテキストを作成できません',
      unsupportedWatermarkType: 'サポートされていないウォーターマークタイプ',
      watermarkGenerationFailed: 'ウォーターマーク生成失敗',
      tempCanvasContextFailed: '一時2Dレンダリングコンテキストを作成できません',
      noWatermarkFile: 'ウォーターマークファイルが提供されていません',
      fileLoadFailed: 'ファイル読み込み失敗',
      imageLoadFailed: 'ファイル読み込み失敗：{url}',
      cannotLoadImage: 'ファイルを読み込めません',
    },
    status: {
      completed: '✅ ウォーターマーク処理完了（{processedCount}/{totalFiles}）',
      failed: '❌ ウォーターマーク処理失敗（{failedCount}/{totalFiles}）',
      partialSuccess: '⚠️ ウォーターマーク処理が部分的に成功（成功：{processedCount}、失敗：{failedCount}）',
    },
  },
  colorUtils: {
    errors: {
      invalidImageUrl: '無効なファイルURL',
      cannotCreateCanvasContext: 'キャンバスコンテキストを作成できません',
      cannotGetImageData: 'ファイルデータを取得できません（CORSの問題の可能性があります）',
      imageLoadFailed: 'ファイルの読み込みに失敗しました。色を抽出できません',
      extractColorFailed: '色抽出失敗：{error}',
    },
  },
}
