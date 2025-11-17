/**
 * Automation Tasks Page Text - Cyber Style
 */
export const automation = {
  page: {
    title: 'Automation Commands',
    subtitle: 'Monitor AI tagging and vectorization task execution status in real-time',
    refresh: 'Refresh',
    tip: 'Queue will automatically process injected data units without manual intervention. System executes all commands sequentially.',
  },
  tagging: {
    title: 'AI Tagging Tasks',
    subtitle: 'Auto identify data unit content and generate tags',
    statusRunning: 'Running',
    statusPaused: 'Paused',
    progress: {
      label: 'Execution Progress',
      completed: '{done} / {total} completed',
    },
    queue: {
      position: 'Your data unit is at position {position} in queue',
      pending: '{count} data units waiting in queue',
    },
    status: {
      none: 'Not Processed',
      pending: 'In Queue',
      processing: 'Executing',
      done: 'Completed',
      failed: 'Failed',
      ignored: 'Ignored',
    },
    empty: {
      title: 'No AI Tagging Tasks',
      subtitle: 'AI analysis will start automatically after injecting data units',
    },
  },
  vector: {
    title: 'Vectorization Tasks',
    subtitle: 'Convert data units to vectors for similarity search',
    statusRunning: 'Running',
    statusPaused: 'Paused',
    progress: {
      label: 'Execution Progress',
      completed: '{done} / {total} completed',
    },
    queue: {
      position: 'Your vector task is at position {position} in queue',
      pending: '{count} vector tasks waiting in queue',
    },
    status: {
      pending: 'Pending',
      reset: 'Reset',
      processing: 'Executing',
      completed: 'Completed',
      failed: 'Failed',
    },
    empty: {
      title: 'No Vectorization Tasks',
      subtitle: 'Vectorization will start automatically after AI tagging completes',
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
