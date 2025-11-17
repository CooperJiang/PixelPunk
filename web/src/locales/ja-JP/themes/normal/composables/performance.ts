/**
 * Composables texts
 */
export const performance = {
  suggestions: {
    slowPageLoad: 'ページ読み込み時間が長すぎます。ファイルサイズとコード分割の最適化を推奨します',
    slowRender: 'コンポーネントレンダリング時間が長すぎます。Vue.memoの使用またはcomputedプロパティの最適化を推奨します',
    highMemory: 'メモリ使用量が高すぎます。メモリリークの確認またはデータ構造の最適化を推奨します',
  },
}

export const chunkedUpload = {
  errors: {
    sessionNotFound: 'セッションIDが存在しません。アップロードステータスを確認できません',
    sessionDataMissing: 'セッションIDまたはファイルMD5が不足しています。セッションデータを作成できません',
    sessionIdMissing: 'セッションIDが存在しません。チャンクをアップロードできません',
    unknownError: '不明なエラー',
    networkError: '（ネットワークエラー、ネットワーク接続を確認してください）',
    serverError: '（サーバーエラー、後でもう一度お試しください）',
    rateLimited: '（リクエストが頻繁すぎます、後でもう一度お試しください）',
  },
}
