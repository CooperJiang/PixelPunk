/**
 * Resumable uploads component
 */
export const resumableUploads = {
  title: 'Resume Upload',
  relativeTime: {
    justNow: 'Just Now',
    minutesAgo: '{count} minutes ago',
    hoursAgo: '{count} hours ago',
    daysAgo: '{count} days ago',
    empty: {
      title: 'No Resumable Uploads',
      subtitle: 'Currently no interrupted upload tasks',
      labels: {
        chunks: 'Chunks',
        actions: {
          selectAll: 'Select All',
          clearSelection: 'Clear Selection',
          deleteSelected: 'Delete Selected',
          footer: {
            cancel: 'Cancel',
            resume: 'Resume Upload',
          },
        },
      },
    },
  },
}
