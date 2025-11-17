/**
 * Global Upload Drawer Component
 */
export const globalUploadDrawer = {
  stats: {
    uploading: '転送中',
    pending: 'キュー待機中',
    completed: '完了',
    failed: '失敗',
    batchActions: {
      pauseAll: 'すべて一時停止',
      resumeAll: 'すべて再開',
      clearCompleted: '完了をクリア',
      clearAll: 'キューをクリア',
      statusMessages: {
        uploading: 'データ転送中',
        completed: '転送完了',
        failed: '転送失敗',
        paused: '転送一時停止',
        pending: '転送待機中',
        preparing: '転送準備中',
        actions: {
          pause: '一時停止',
          resume: '再開',
          retry: '再試行',
          copyLink: 'リンクをコピー',
          remove: '削除',
          messages: {
            pausedAll: 'すべての転送を一時停止しました',
            resumedAll: 'すべての転送を再開しました',
            clearedCompleted: '{count}件の完了タスクをクリアしました',
            clearedAll: '転送キューをクリアしました',
            linkCopied: 'リンクをクリップボードにコピーしました',
            copyFailed: 'コピー操作失敗',
            misc: {
              chunkedUpload: 'チャンク転送',
              noTasks: '転送タスクなし',
              goToUpload: '転送コンソールに移動',
              viewFullPage: 'フルコンソールを表示',
            },
          },
        },
      },
    },
  },
}
