/**
 * Resumable uploads component
 */
export const resumableUploads = {
  title: 'アップロード再開',
  relativeTime: {
    justNow: 'たった今',
    minutesAgo: '{count}分前',
    hoursAgo: '{count}時間前',
    daysAgo: '{count}日前',
    empty: {
      title: '再開可能なアップロードがありません',
      subtitle: '現在中断されたアップロードタスクがありません',
      labels: {
        chunks: 'チャンク',
        actions: {
          selectAll: 'すべて選択',
          clearSelection: '選択クリア',
          deleteSelected: '選択したものを削除',
          footer: {
            cancel: 'キャンセル',
            resume: 'アップロード再開',
          },
        },
      },
    },
  },
}

