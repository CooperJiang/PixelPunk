/**
 * Network utility texts
 */
export const network = {
  http: {
    cancelDuplicate: '重複リクエストをキャンセル：{method} {url}',
    requestTimeout: 'リクエストタイムアウトクリーンアップ',
    requestCancelled: 'リクエストキャンセル済み',
    pageSwitch: 'ページ切り替え、すべての保留中のリクエストをキャンセル',
    cancelSpecificUrl: '特定のURLのリクエストをキャンセル',

    operationSuccess: '操作成功',
    operationFailed: '操作失敗',
    requestFailed: 'リクエスト失敗',
    requestFailedWithStatus: 'リクエスト失敗（{status}）',

    networkError: 'ネットワーク接続失敗、ネットワークを確認してください',
    configError: 'リクエスト設定エラー',
    unknown: '不明',

    ipNotInWhitelist: 'あなたのIPはアクセスホワイトリストにありません',
    ipInBlacklist: 'あなたのIPはアクセスブラックリストにあります',
    domainNotInWhitelist: 'あなたのドメインはアクセスホワイトリストにありません',
    domainInBlacklist: 'あなたのドメインはアクセスブラックリストにあります',
    userAccountDisabled: 'ユーザーアカウントが無効化されています',

    envVarReadError: '環境変数VITE_API_BASE_URLを読み取れません：',
    toastLoadError: 'トーストコンポーネントを読み込めません',
  },
}

