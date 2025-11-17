/**
 * Composables - Cyberpunk Style
 */
export const performance = {
  suggestions: {
    slowPageLoad: 'ページ読み込み遅延が高すぎます。データパッケージサイズとコード分割の最適化を推奨します',
    slowRender: 'コンポーネントレンダリング遅延が高すぎます。Vue.memoの使用または計算プロパティの最適化を推奨します',
    highMemory: 'メモリ使用量が高すぎます。メモリリークの確認またはデータ構造の最適化を推奨します',
  },
}

export const chunkedUpload = {
  errors: {
    sessionNotFound: 'セッションIDが見つかりません。アップロードステータスを確認できません',
    sessionDataMissing: 'セッションIDまたはデータフィンガープリントが不足しています。セッションデータを作成できません',
    sessionIdMissing: 'セッションIDが見つかりません。データチャンクをアップロードできません',
    unknownError: '不明なエラー',
    networkError: '（ネットワーク接続エラー、ネットワークを確認してください）',
    serverError: '（サーバーレスポンスエラー、後でもう一度お試しください）',
    rateLimited: '（リクエスト頻度が高すぎます。後でもう一度お試しください）',
  },
}
