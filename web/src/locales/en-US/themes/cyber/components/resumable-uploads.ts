/**
 * Resumable Uploads Component
 */
export const resumableUploads = {
  title: 'Resume Data Transfer',
  relativeTime: {
    justNow: 'Just Now',
    minutesAgo: '{count} minutes ago',
    hoursAgo: '{count} hours ago',
    daysAgo: '{count} days ago',
    empty: {
      title: 'No Interrupted Transfer Tasks',
      subtitle: 'No data transfers need to be resumed',
      labels: {
        chunks: 'Data Chunks',
        actions: {
          selectAll: 'Select All',
          clearSelection: 'Clear Selection',
          deleteSelected: 'Delete Selected',
          footer: {
            cancel: 'Cancel',
            resume: 'Resume Transfer',
          },
        },
      },
    },
  },
}
