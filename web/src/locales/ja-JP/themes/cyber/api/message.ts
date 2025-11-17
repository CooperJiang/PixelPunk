/**
 * Message System i18n - Cyber Theme
 */
export const message = {
  types: {
    system: {
      maintenance: 'システムメンテナンス',
      update: 'システム更新',
      announcement: 'システムお知らせ',
    },
    account: {
      register: '登録成功',
      storage_granted: 'ストレージクォータ付与',
      bandwidth_granted: '帯域幅クォータ付与',
    },
    content: {
      review_approved: 'コンテンツレビュー承認',
      review_rejected: 'コンテンツレビュー拒否',
      review_pending: 'レビュー待ち',
    },
    storage: {
      quota_warning: 'ストレージクォータ警告',
      quota_exceeded: 'ストレージクォータ超過',
      quota_increased: 'ストレージクォータ増加',
      quota_decreased: 'ストレージクォータ調整',
    },
    file: {
      deleted_by_admin: 'データ削除',
      batch_deleted_by_admin: '一括データ削除',
      hard_deleted_by_admin: '永続データ削除',
      expiry_warning: 'データ期限切れリマインダー',
      thumbnail_failed: 'サムネイル生成失敗',
    },
    security: {
      login_alert: 'ログインアラート',
      password_changed: 'パスワード変更',
    },
    apikey: {
      created: 'APIキー作成',
      deleted: 'APIキー削除',
      regenerated: 'APIキー再生成',
      disabled: 'APIキー無効化',
      enabled: 'APIキー有効化',
    },
    random_api: {
      created: 'ランダムAPI作成',
      deleted: 'ランダムAPI削除',
      disabled: 'ランダムAPI無効化',
      enabled: 'ランダムAPI有効化',
    },
    share: {
      expiry_warning: 'シェア期限切れ間近',
    },
  },
  priority: {
    high: '高',
    normal: '通常',
    low: '低',
  },
  status: {
    unread: '未読',
    read: '既読',
    deleted: '削除済み',
    unknown: '不明',
  },
  time: {
    justNow: 'たった今',
    minutesAgo: '{n}分前',
    hoursAgo: '{n}時間前',
    daysAgo: '{n}日前',
  },
}
