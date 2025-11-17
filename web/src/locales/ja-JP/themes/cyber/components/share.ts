/**
 * Share Component
 */
export const share = {
  createShare: '共有リンクを確立',
  shareSelected: '選択項目を共有',
  cancel: '操作を中止',
  passwordDialog: {
    title: 'アクセス制限: パスワード必須',
    message: 'この共有はアクセス暗号化が有効です。アクセスキーを入力して続行してください',
    placeholder: 'アクセスキーを入力',
    verify: 'キーを確認',
  },
  selectionToolbar: {
    selectedPrefix: '選択済み',
    selectedSuffix: '件',
    selectAll: 'すべて選択',
    batchDownload: '一括ダウンロード',
  },
  selectionInfo: {
    folders: '個のデータコンテナ',
    images: '個のデータファイル',
  },
  messages: {
    loginRequired: '共有する前にニューラルリンクを確立してください',
    noItemsSelected: '共有するデータコンテナまたはデータファイルを選択してください',
  },
  image: {
    alt: 'データファイル',
    downloadError: 'データ抽出失敗',
    loadTimeout: 'データ読み込みタイムアウト',
  },
  errors: {
    invalidLink: '無効な共有リンク',
    invalidLinkDesc: 'リンクフォーマットエラーまたは共有がクリアされました',
    loadFailed: 'データ読み込み失敗',
    loadFailedDesc: 'リンクを確立できません。後でもう一度お試しください',
    maxViewsReached: 'アクセス数が上限に達しました',
    maxViewsReachedDesc: 'この共有は最大アクセス数制限に達しました',
    shareExpired: '共有期限切れ',
    shareExpiredDesc: 'このリンクは期限切れです。ソースノードに連絡して再共有してください',
    accessRestricted: 'アクセス制限',
    accessRestrictedDesc: '権限が不足しています。この共有にアクセスできません',
    shareNotFound: '共有が見つかりません',
    shareNotFoundDesc: 'この共有はクリアされたか、リンクが無効です',
    shareInvalid: '共有が無効です',
    shareInvalidDesc: 'この共有は期限切れです。ソースノードに連絡して再共有してください',
    passwordVerifySuccess: 'キー確認成功、データを読み込み中...',
    passwordVerifyNoToken: 'キー確認成功しましたが、アクセストークンを取得できませんでした',
    passwordVerifyFailed: 'キー確認失敗',
    passwordVerifyFailedRetry: 'キー確認失敗、再試行してください',
    visitorSubmitSuccess: '情報送信成功、ご協力ありがとうございます',
    visitorSubmitFailed: '情報送信失敗、再試行してください',
  },
  defaults: {
    shareName: '共有',
    shareTitle: 'データ共有',
    shareDescription: '共有には{count}個のデータファイルが含まれています',
  },
}
