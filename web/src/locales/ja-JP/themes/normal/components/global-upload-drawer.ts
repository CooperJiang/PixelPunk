/**
 * Global upload drawer component
 */
export const globalUploadDrawer = {
  stats: {
    uploading: 'アップロード中',
    pending: '保留中',
    completed: '完了',
    failed: '失敗',
    batchActions: {
      pauseAll: 'すべて一時停止',
      resumeAll: 'すべて再開',
      clearCompleted: '完了をクリア',
      clearAll: 'すべてクリア',
      statusMessages: {
        uploading: 'アップロード中',
        completed: 'アップロード完了',
        failed: 'アップロード失敗',
        paused: '一時停止',
        pending: 'アップロード待ち',
        preparing: '準備中',
        actions: {
          pause: '一時停止',
          resume: '再開',
          retry: '再試行',
          copyLink: 'リンクコピー',
          remove: '削除',
          messages: {
            pausedAll: 'すべてのアップロードを一時停止',
            resumedAll: 'すべてのアップロードを再開',
            clearedCompleted: '{count}件の完了タスクをクリア',
            clearedAll: 'すべてのタスクをクリア',
            linkCopied: 'リンクコピー済み',
            copyFailed: 'コピー失敗',
            misc: {
              chunkedUpload: 'チャンクアップロード',
              noTasks: 'アップロードタスクがありません',
              goToUpload: 'アップロードページに移動',
              viewFullPage: 'フルページを表示',
            },
          },
        },
      },
    },
  },
}

