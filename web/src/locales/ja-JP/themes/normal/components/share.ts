/**
 * Share component
 */
export const share = {
  createShare: '共有作成',
  shareSelected: '選択したものを共有',
  cancel: 'キャンセル',
  passwordDialog: {
    title: 'この共有にはパスワード確認が必要です',
    message: 'この共有コンテンツはパスワード保護されています。続行するにはパスワードを入力してください',
    placeholder: '共有パスワードを入力してください',
    verify: 'パスワード確認',
  },
  selectionToolbar: {
    selectedPrefix: '選択済み',
    selectedSuffix: '件',
    selectAll: 'すべて選択',
    batchDownload: 'パックダウンロード',
  },
  selectionInfo: {
    folders: 'フォルダ',
    images: 'リソース',
  },
  messages: {
    loginRequired: '共有する前にログインしてください',
    noItemsSelected: 'まず共有するフォルダまたはファイルを選択してください',
  },
  image: {
    alt: 'ファイル',
    downloadError: 'ダウンロード失敗',
    loadTimeout: 'ファイル読み込みタイムアウト',
  },
  errors: {
    invalidLink: '無効な共有リンク',
    invalidLinkDesc: 'リンクが完全かどうかを確認してください。または、この共有が削除された可能性があります',
    loadFailed: '読み込み失敗',
    loadFailedDesc: '共有コンテンツを読み込めません。後でもう一度お試しください',
    maxViewsReached: '最大閲覧数到達',
    maxViewsReachedDesc: 'この共有は作成者が設定した最大閲覧数に達しました',
    shareExpired: '共有期限切れ',
    shareExpiredDesc: 'この共有リンクは期限切れです。共有者に連絡して再共有してください',
    accessRestricted: 'アクセス制限',
    accessRestrictedDesc: 'この共有にアクセスする権限がありません',
    shareNotFound: '共有が見つかりません',
    shareNotFoundDesc: 'この共有は削除されたか、リンクが無効です',
    shareInvalid: '共有が無効です',
    shareInvalidDesc: 'この共有は期限切れです。共有者に連絡して再共有してください',
    passwordVerifySuccess: 'パスワード確認成功、コンテンツを読み込み中...',
    passwordVerifyNoToken: 'パスワード確認成功、ただしアクセストークンの取得に失敗しました',
    passwordVerifyFailed: 'パスワード確認失敗',
    passwordVerifyFailedRetry: 'パスワード確認失敗、もう一度お試しください',
    visitorSubmitSuccess: '情報送信成功、ご協力ありがとうございます',
    visitorSubmitFailed: '情報送信失敗、もう一度お試しください',
  },
  defaults: {
    shareName: '共有',
    shareTitle: 'ファイル共有',
    shareDescription: '共有には{count}件のファイルが含まれています',
  },
}

