/**
 * Upload Store i18n configuration - Normal theme
 */
export const upload = {
  control: {
    paused: '一時停止',
    preparingResume: '再開準備中',
    userCancelled: 'ユーザーキャンセル',
    cancelled: 'キャンセル済み',
    preparingRetry: '再試行準備中',
  },
  queue: {
    preparingAnalysis: '分析準備中...',
    analyzing: 'ファイル分析中...',
    checkingInstant: 'インスタントアップロード可能性を確認中...',
    calculatingMD5: 'ファイルMD5計算中... {progress}%',
    checkingServerDuplicate: 'サーバー重複ファイルを確認中...',
    preparingUpload: 'インスタントアップロード実行中...',
    instantComplete: 'インスタントアップロード完了',
    analysisComplete: '分析完了',
    waitingUpload: 'アップロード準備中',
  },
  upload: {
    analyzingFile: 'ファイル分析中...',
    analyzingProgress: 'ファイル分析中 {progress}%',
    instantProgress: 'インスタントアップロード {progress}%',
    instantComplete: 'インスタントアップロード完了',
    uploading: 'アップロード中...',
    mergingChunks: 'チャンク統合中...',
    uploadComplete: 'アップロード完了',
    uploadFailed: 'アップロード失敗',
    instantFailed: 'インスタントアップロード失敗',
    preparingRetry: '再試行準備中',
    retrying: '再試行中',
    postProcessing: '後処理中...',
    analysisError: '分析エラー',
    remainingTimeMinutes: '{minutes}分{seconds}秒',
    remainingTimeSeconds: '{seconds}秒',
    uploadingProgress: 'アップロード中 {progress}% - {speed} - 残り{remaining}',
    uploadFailedWithError: 'アップロード失敗：{error}',
  },
}
