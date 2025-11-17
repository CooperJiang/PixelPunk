/**
 * Message Translations - English Cyber Theme
 */
export const message = {
  system: {
    maintenance: {
      title: 'System Maintenance Protocol',
      content: 'Neural network will enter maintenance mode at {time}, expected duration: {duration}. Some nodes may go offline during maintenance, please sync data in advance.',
    },
    announcement: {
      title: 'System Broadcast',
      content: '{content}',
    },
  },
  account: {
    register: {
      title: 'Welcome to PixelPunk Neural Network',
      content: 'Connection established! Your node has been successfully registered. System allocated {storage} of data space for you. Start your data management journey!',
    },
    storageGranted: {
      title: 'Data Space Expanded',
      content: 'System upgrade! Your data space has been expanded by {size}, current total capacity {total}.',
    },
    bandwidthGranted: {
      title: 'Bandwidth Enhanced',
      content: 'Network enhancement! Your monthly bandwidth has been increased by {size}, current monthly bandwidth {total}.',
    },
  },
  content: {
    reviewApproved: {
      title: 'Data Review Approved',
      content: 'Your data「{content_name}」has been approved and is now visible on the network.',
    },
    reviewRejected: {
      title: 'Data Review Failed',
      content: 'Sorry, your data「{content_name}」was not approved. Reason: {reason}',
    },
    reviewPending: {
      title: 'Data Under Review',
      content: 'Your data「{content_name}」has been submitted for review, system is processing.',
    },
  },
  storage: {
    quotaWarning: {
      title: 'Data Space Alert',
      content: 'You have used {used_percent}% ({used}/{total}) of your data space, recommend cleanup or upgrade capacity.',
    },
    quotaIncreased: {
      title: 'Data Space Upgraded',
      content: 'Your data space has been upgraded from {old_size} to {new_size}.',
    },
    quotaDecreased: {
      title: 'Data Space Adjusted',
      content: 'Your data space has been adjusted from {old_size} to {new_size}.',
    },
  },
  file: {
    deletedByAdmin: {
      title: 'Data Removed',
      content: 'Your data unit「{file_name}」has been removed by administrator for violating network protocol.',
    },
    batchDeletedByAdmin: {
      title: 'Data Batch Removed',
      content: 'You have {count} data units removed by administrator for violating network protocol.',
    },
    hardDeletedByAdmin: {
      title: 'Data Permanently Deleted',
      content: 'Your data unit「{file_name}」has been permanently deleted by administrator and cannot be recovered.',
    },
    expiryWarning: {
      title: 'You have {count} data unit(s) about to expire',
      content: 'The following data units are about to expire, please backup in time:\n{file_list}\n\nExpired data will be automatically purged and cannot be recovered.',
    },
    thumbnailFailed: {
      title: 'Preview Generation Failed',
      content: 'Preview generation for data unit「{file_name}」failed, but the data itself is not affected.',
    },
  },
  security: {
    loginAlert: {
      title: 'New Node Connection Alert',
      content: 'Your account established connection at {time} from {ip} ({location}). If this was not you, please change your access key immediately.',
    },
  },
  apikey: {
    created: {
      title: 'API Key Generated',
      content: 'You generated a new API key「{key_name}」. Please keep your access credential safe.',
    },
    deleted: {
      title: 'API Key Destroyed',
      content: 'API key「{key_name}」has been destroyed.',
    },
    regenerated: {
      title: 'API Key Regenerated',
      content: 'API key「{key_name}」has been regenerated, the old key will become invalid immediately.',
    },
    disabled: {
      title: 'API Key Disabled',
      content: 'API key「{key_name}」has been disabled.',
    },
    enabled: {
      title: 'API Key Activated',
      content: 'API key「{key_name}」has been activated.',
    },
  },
  randomApi: {
    created: {
      title: 'Random Data API Created',
      content: 'You created random data API「{api_name}」.',
    },
    deleted: {
      title: 'Random Data API Deleted',
      content: 'Random data API「{api_name}」has been deleted.',
    },
    disabled: {
      title: 'Random Data API Disabled',
      content: 'Random data API「{api_name}」has been disabled.',
    },
    enabled: {
      title: 'Random Data API Activated',
      content: 'Random data API「{api_name}」has been activated.',
    },
  },
  share: {
    expiryWarning: {
      title: 'Share Link About to Expire',
      content: 'Your share link「{share_name}」will expire at {expires_at}.',
    },
  },
  messageList: {
    title: 'Message Center',
    total: 'Total',
    unread: 'Unread',
    allStatus: 'All Status',
    allTypes: 'All Types',
    markAllRead: 'Mark All as Read',
    emptyTitle: 'No Messages',
    emptyDescription: 'You have not received any messages yet',
    loadingText: 'Loading...',
    columns: {
      type: 'Type',
      title: 'Title',
      content: 'Content',
      status: 'Status',
      priority: 'Priority',
      time: 'Time',
      actions: 'Actions',
    },
    statusLabels: {
      unread: 'Unread',
      read: 'Read',
    },
    priorityLabels: {
      high: 'High',
      normal: 'Normal',
      low: 'Low',
    },
    actions: {
      markRead: 'Mark as Read',
      delete: 'Delete',
    },
    dialog: {
      deleteTitle: 'Confirm Deletion',
      deleteMessage: 'Are you sure you want to delete this message? This action cannot be undone.',
      cancel: 'Cancel',
      confirmDelete: 'Delete',
    },
    toasts: {
      deleted: 'Message deleted',
      deleteFailed: 'Failed to delete, please try again',
      markReadSuccess: 'Marked as read',
      markReadFailed: 'Failed to mark as read, please try again',
      markAllReadSuccess: 'All messages marked as read',
      markAllReadFailed: 'Failed to mark all as read, please try again',
    },
  },
}
