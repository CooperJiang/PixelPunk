/**
 * Announcement Component
 */
export const announcement = {
  button: {
    ariaLabel: 'システムお知らせ',
  },

  drawer: {
    title: 'システムお知らせ',
  },

  status: {
    loading: 'データを読み込み中...',
    disabled: 'お知らせシステム切断',
  },

  empty: {
    title: 'お知らせなし',
    description: '利用可能なお知らせがありません',
  },

  list: {
    pinned: 'ピン留め',
    viewDetail: '詳細を表示',
  },

  messages: {
    systemDisabled: 'お知らせシステム無効',
    loadListFailed: 'お知らせリストの読み込みに失敗しました',
    loadDetailFailed: 'お知らせ詳細の読み込みに失敗しました',
    neverShowAgainSuccess: 'このお知らせを再度表示しないように設定しました',
  },

  dialog: {
    title: 'お知らせ詳細',
    loading: 'データを読み込み中...',
    meta: {
      published: '公開日時',
      updated: '更新日時',
    },
    actions: {
      neverShowAgain: '再度表示しない',
      viewMore: 'もっと見る',
      gotIt: '了解',
    },
  },
}
