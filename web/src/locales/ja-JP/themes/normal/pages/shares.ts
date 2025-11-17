/**
 * Share management page texts
 */
export const shares = {
  title: 'マイシェア',
  search: 'シェアを検索',
  page: {
    title: 'マイシェア',
    subtitle: '作成したすべてのシェアリンクを管理および表示',
  },
  filter: {
    all: 'すべて',
    normal: '通常',
    expired: '期限切れ',
    deleted: '削除済み',
    disabled: '無効',
  },
  status: {
    normal: '通常',
    expired: '期限切れ',
    deleted: '削除済み',
    disabled: '無効',
    unknown: '不明',
  },
  card: {
    viewCount: '{count}閲覧',
    maxViews: '制限{max}閲覧',
    passwordProtected: 'パスワード保護',
    collectVisitor: '訪問者を収集',
    accessNotification: 'アクセス通知',
    copyLink: 'リンクをコピー',
    folderCount: '{count}フォルダ',
    fileCount: '{count}リソース',
    viewShare: 'シェアを表示',
    visitorInfo: '訪問者情報',
    delete: '削除',
  },
  empty: {
    title: 'シェアコンテンツなし',
    subtitle: 'まだシェアを作成していません。フォルダページからシェアを作成できます',
    goToFolders: 'フォルダに移動',
  },
  delete: {
    confirm: 'シェア{name}を削除してもよろしいですか？',
    forceDelete: '完全に削除（復元不可）',
    cancel: 'キャンセル',
    confirmButton: '削除',
  },
  unnamed: '無名シェア',
  toast: {
    loadFailed: 'シェアリストの読み込みに失敗しました',
    linkCopied: 'シェアリンクをクリップボードにコピーしました',
    deleteSuccess: 'シェアが正常に削除されました',
    deleteFailed: 'シェアの削除に失敗しました',
  },
  dialog: {
    deleteTitle: '削除を確認',
    deleteMessage: 'シェア「{name}」を削除してもよろしいですか？',
    deleteHint: '削除後は復元できません',
  },
  visitor: {
    dialogTitle: '訪問者情報リスト',
    search: '訪問者名、メールまたはIPアドレスを検索',
    deleteConfirm: 'この訪問者情報を削除してもよろしいですか？',
    deleteHint: '削除後は復元できません',
    loadFailed: '訪問者情報の読み込みに失敗しました。もう一度お試しください',
    deleteSuccess: '訪問者情報が正常に削除されました',
    deleteFailed: '訪問者情報の削除に失敗しました。もう一度お試しください',
    table: {
      name: '訪問者名',
      email: 'メール',
      ip: 'IPアドレス',
      visits: '訪問数',
      lastVisit: '最終訪問日時',
      actions: 'アクション',
    },
    pagination: {
      info: '合計{total}項目、{from}-{to}を表示',
    },
    empty: {
      title: '訪問者情報なし',
      desc: 'まだ訪問者が情報を送信していないか、フィルター条件に一致する結果がありません',
    },
  },
}
