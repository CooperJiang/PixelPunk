/**
 * Share Management Page Text - Cyber Style
 */
export const shares = {
  title: 'マイ共有リンク',
  search: '共有リンクを検索',
  page: {
    title: 'マイ共有リンク',
    subtitle: '作成したすべての共有リンクを管理および表示',
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
    viewCount: '{count}回のアクセス',
    maxViews: '制限{max}回',
    passwordProtected: 'パスワード保護',
    collectVisitor: '訪問者収集',
    accessNotification: 'アクセス通知',
    copyLink: 'リンクをコピー',
    folderCount: '{count}個のデータフォルダ',
    fileCount: '{count}個のデータユニット',
    viewShare: '共有を表示',
    visitorInfo: '訪問者情報',
    delete: '削除',
  },
  empty: {
    title: '共有コンテンツなし',
    subtitle: 'まだ共有リンクを作成していません。フォルダページから共有を作成できます',
    goToFolders: 'フォルダに移動',
  },
  delete: {
    confirm: '共有リンク{name}を削除してもよろしいですか？',
    forceDelete: '完全削除（元に戻せません）',
    cancel: 'キャンセル',
    confirmButton: '削除',
  },
  unnamed: '無名共有',
  toast: {
    loadFailed: '共有リンクリストの読み込みに失敗しました',
    linkCopied: '共有リンクをクリップボードにコピーしました',
    deleteSuccess: '共有リンクを削除しました',
    deleteFailed: '共有リンクの削除に失敗しました',
  },
  dialog: {
    deleteTitle: '削除確認',
    deleteMessage: '共有リンク「{name}」を削除してもよろしいですか？',
    deleteHint: '削除後は回復できません',
  },
  visitor: {
    dialogTitle: '訪問者情報リスト',
    search: '訪問者名、メールまたはIPノードを検索',
    deleteConfirm: 'この訪問者情報を削除してもよろしいですか？',
    deleteHint: '削除後は回復できません',
    loadFailed: '訪問者情報の読み込みに失敗しました。再試行してください',
    deleteSuccess: '訪問者情報を削除しました',
    deleteFailed: '訪問者情報の削除に失敗しました。再試行してください',
    table: {
      name: '訪問者名',
      email: 'メール',
      ip: 'IPアドレス',
      visits: 'アクセス数',
      lastVisit: '最終訪問時間',
      actions: 'アクション',
    },
    pagination: {
      info: '合計{total}件、現在{from}-{to}を表示',
    },
    empty: {
      title: '訪問者情報なし',
      desc: 'まだ訪問者が情報を送信していないか、フィルター条件に一致する結果がありませんでした',
    },
  },
}
