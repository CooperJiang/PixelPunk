/**
 * Message system i18n configuration - Normal theme
 */
export const message = {
  types: {
    system: {
      maintenance: 'システムメンテナンス',
      update: 'システム更新',
      announcement: 'システムお知らせ',
    },
    account: {
      register: '登録ウェルカム',
      storage_granted: 'ストレージ付与',
      bandwidth_granted: '帯域幅付与',
    },
    content: {
      review_approved: 'レビュー承認',
      review_rejected: 'レビュー拒否',
      review_pending: 'レビュー待ち',
    },
    storage: {
      quota_warning: 'ストレージ警告',
      quota_exceeded: 'ストレージ超過',
      quota_increased: 'ストレージ増加',
      quota_decreased: 'ストレージ調整',
    },
    file: {
      deleted_by_admin: 'ファイル削除',
      batch_deleted_by_admin: '一括削除',
      hard_deleted_by_admin: '完全削除',
      expiry_warning: 'ファイル期限警告',
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
      created: 'ランダム画像API作成',
      deleted: 'ランダム画像API削除',
      disabled: 'ランダム画像API無効化',
      enabled: 'ランダム画像API有効化',
    },
    share: {
      expiry_warning: '共有期限間近',
    },
  },
  priority: {
    high: '高',
    normal: '中',
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
