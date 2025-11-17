/**
 * Global upload drawer component
 */
export const globalUploadDrawer = {
  stats: {
    uploading: 'Uploading',
    pending: 'Pending',
    completed: 'Completed',
    failed: 'Failed',
    batchActions: {
      pauseAll: 'Pause All',
      resumeAll: 'Resume All',
      clearCompleted: 'Clear Completed',
      clearAll: 'Clear All',
      statusMessages: {
        uploading: 'Uploading',
        completed: 'Upload Completed',
        failed: 'Upload Failed',
        paused: 'Paused',
        pending: 'Waiting to Upload',
        preparing: 'Preparing',
        actions: {
          pause: 'Pause',
          resume: 'Resume',
          retry: 'Retry',
          copyLink: 'Copy Link',
          remove: 'Remove',
          messages: {
            pausedAll: 'Paused all uploads',
            resumedAll: 'Resumed all uploads',
            clearedCompleted: 'Cleared {count} completed tasks',
            clearedAll: 'Cleared all tasks',
            linkCopied: 'Link copied',
            copyFailed: 'Copy failed',
            misc: {
              chunkedUpload: 'Chunked Upload',
              noTasks: 'No upload tasks',
              goToUpload: 'Go to Upload Page',
              viewFullPage: 'View Full Page',
            },
          },
        },
      },
    },
  },
}
