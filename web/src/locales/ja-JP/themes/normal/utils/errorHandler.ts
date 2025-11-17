/**
 * Error handler utility i18n configuration - Normal theme
 */
export const errorHandler = {
  httpStatus: {
    400: 'リクエストパラメータエラー',
    401: 'ログイン期限切れ、再度ログインしてください',
    403: 'アクセス権限がありません',
    404: 'リクエストされたリソースが見つかりません',
    429: 'リクエストが多すぎます。後でもう一度お試しください',
    500: '内部サーバーエラー',
    502: '不正なゲートウェイ',
    503: 'サービス一時利用不可',
  },
  fallback: {
    unknown: '不明なエラー',
    toastUnavailable: 'トーストが利用できません。エラーメッセージ：',
  },
}
