/**
 * Reset password page texts
 */
export const resetPassword = {
  title: 'パスワードをリセット',
  verifying: {
    title: 'リセットトークンを確認中',
    loading: 'リセットトークンを確認中...',
  },
  error: {
    missingToken: 'リセットトークンがありません',
    invalidToken: 'トークンが無効です',
    invalidOrExpired: 'トークンが無効または期限切れです',
    hint: 'もう一度パスワードリセットをリクエストしてください',
  },
  form: {
    subtitle: 'アカウント {email} の新しいパスワードを設定',
    newPassword: {
      label: '新しいパスワード',
      placeholder: '新しいパスワードを入力してください（6-50文字）',
    },
    confirmPassword: {
      label: 'パスワードを確認',
      placeholder: '新しいパスワードを再度入力してください',
    },
    submit: 'パスワードをリセット',
    submitting: 'リセット中...',
  },
  validation: {
    passwordRequired: '新しいパスワードを入力してください',
    passwordTooShort: 'パスワードは6文字以上である必要があります',
    passwordTooLong: 'パスワードは50文字を超えることはできません',
    confirmRequired: '新しいパスワードを確認してください',
    passwordMismatch: 'パスワードが一致しません',
  },
  actions: {
    backToLogin: 'ログインページに戻る',
  },
  protocol: 'システムセキュリティプロトコル有効化 / SYSTEM SECURITY PROTOCOL',
}

