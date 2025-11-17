/**
 * Network Utils Text - Cyber Style
 */
export const network = {
  http: {
    cancelDuplicate: '重複リクエストを中止: {method} {url}',
    requestTimeout: 'リクエストタイムアウトクリーンアップ',
    requestCancelled: 'リクエスト中止済み',
    pageSwitch: 'ページ切り替え、すべての保留中のリクエストを中止',
    cancelSpecificUrl: '特定のURLのリクエストを中止',

    operationSuccess: '操作完了',
    operationFailed: '操作失敗',
    requestFailed: '転送失敗',
    requestFailedWithStatus: '転送失敗（{status}）',

    networkError: 'ニューラルリンク切断、ネットワークを確認してください',
    configError: 'リクエスト設定エラー',
    unknown: '不明',

    ipNotInWhitelist: 'あなたのIPはアクセスホワイトリストにありません',
    ipInBlacklist: 'あなたのIPはアクセスブラックリストにあります',
    domainNotInWhitelist: 'あなたのドメインはアクセスホワイトリストにありません',
    domainInBlacklist: 'あなたのドメインはアクセスブラックリストにあります',
    userAccountDisabled: 'ユーザーアカウントが無効化されています',

    envVarReadError: '環境変数VITE_API_BASE_URLを読み取れません:',
    toastLoadError: 'トーストコンポーネントを読み込めません',
  },
}
