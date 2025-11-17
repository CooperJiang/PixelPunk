/**
 * Resumable Uploads Component
 */
export const resumableUploads = {
  title: 'データ転送を再開',
  relativeTime: {
    justNow: 'たった今',
    minutesAgo: '{count}分前',
    hoursAgo: '{count}時間前',
    daysAgo: '{count}日前',
    empty: {
      title: '中断された転送タスクなし',
      subtitle: '再開が必要なデータ転送なし',
      labels: {
        chunks: 'データチャンク',
        actions: {
          selectAll: 'すべて選択',
          clearSelection: '選択をクリア',
          deleteSelected: '選択を削除',
          footer: {
            cancel: 'キャンセル',
            resume: '転送を再開',
          },
        },
      },
    },
  },
}
