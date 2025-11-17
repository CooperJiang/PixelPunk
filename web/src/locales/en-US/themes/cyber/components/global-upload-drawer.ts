/**
 * Global Upload Drawer Component
 */
export const globalUploadDrawer = {
  stats: {
    uploading: 'Transferring',
    pending: 'In Queue',
    completed: 'Completed',
    failed: 'Failed',
    batchActions: {
      pauseAll: 'Pause All',
      resumeAll: 'Resume All',
      clearCompleted: 'Clear Completed',
      clearAll: 'Clear Queue',
      statusMessages: {
        uploading: 'Transferring data',
        completed: 'Transfer complete',
        failed: 'Transfer failed',
        paused: 'Transfer paused',
        pending: 'Waiting for transfer',
        preparing: 'Preparing transfer',
        actions: {
          pause: 'Pause',
          resume: 'Resume',
          retry: 'Retry',
          copyLink: 'Copy Link',
          remove: 'Remove',
          messages: {
            pausedAll: 'All transfers paused',
            resumedAll: 'All transfers resumed',
            clearedCompleted: 'Cleared {count} completed tasks',
            clearedAll: 'Transfer queue cleared',
            linkCopied: 'Link copied to clipboard',
            copyFailed: 'Copy operation failed',
            misc: {
              chunkedUpload: 'Chunked Transfer',
              noTasks: 'No transfer tasks',
              goToUpload: 'Go to Transfer Console',
              viewFullPage: 'View Full Console',
            },
          },
        },
      },
    },
  },
}
