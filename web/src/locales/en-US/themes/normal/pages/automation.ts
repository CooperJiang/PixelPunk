/**
 * Automation tasks page texts
 */
export const automation = {
  page: {
    title: 'Automation Tasks',
    subtitle: 'Monitor your AI tagging and vectorization task progress in real-time',
    refresh: 'Refresh',
    tip: 'The queue will automatically process your uploaded images without manual intervention. The system will complete all tasks sequentially.',
  },
  tagging: {
    title: 'AI Tagging Tasks',
    subtitle: 'Automatically identify image content and generate tags',
    statusRunning: 'Running',
    statusPaused: 'Paused',
    progress: {
      label: 'Overall Progress',
      completed: '{done} / {total} completed',
    },
    queue: {
      position: 'Your image is at position {position} in the queue',
      pending: '{count} images waiting in queue',
    },
    status: {
      none: 'Not Processed',
      pending: 'In Queue',
      processing: 'Processing',
      done: 'Completed',
      failed: 'Failed',
      ignored: 'Ignored',
    },
    empty: {
      title: 'No AI Tagging Tasks',
      subtitle: 'AI analysis will start automatically after uploading images',
    },
  },
  vector: {
    title: 'Vectorization Tasks',
    subtitle: 'Convert images to vectors for similarity search',
    statusRunning: 'Running',
    statusPaused: 'Paused',
    progress: {
      label: 'Overall Progress',
      completed: '{done} / {total} completed',
    },
    queue: {
      position: 'Your vector task is at position {position} in the queue',
      pending: '{count} vector tasks waiting in queue',
    },
    status: {
      pending: 'Pending',
      reset: 'Reset',
      processing: 'Processing',
      completed: 'Completed',
      failed: 'Failed',
    },
    empty: {
      title: 'No Vectorization Tasks',
      subtitle: 'Vectorization will start automatically after AI tagging is completed',
    },
  },
  taggingHistory: {
    filter: {
      statusPlaceholder: 'Select Status',
    },
  },
  history: {
    title: 'Tagging History',
    recordCount: '(Total {count} records)',
    refresh: 'Refresh',
    empty: 'No tagging records',
    columns: {
      preview: 'Preview',
      filename: 'Filename',
      status: 'Status',
      format: 'Format',
      size: 'Size',
      resolution: 'Resolution',
      createdAt: 'Created At',
      updatedAt: 'Updated At',
      error: 'Error Message',
    },
  },
  toast: {
    loadHistoryError: 'Failed to load tagging records',
    loadDataError: 'Failed to load task data',
  },
}
