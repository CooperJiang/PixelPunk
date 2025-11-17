/**
 * Message Translations - English Normal Theme
 */
export const message = {
  system: {
    maintenance: {
      title: 'System Maintenance Notice',
      content: 'System will undergo maintenance at {time}, expected duration: {duration}. Some features may be unavailable during maintenance, please prepare in advance.',
    },
    announcement: {
      title: 'System Announcement',
      content: '{content}',
    },
  },
  account: {
    register: {
      title: 'Welcome to PixelPunk',
      content: 'Thank you for registering PixelPunk! Your account has been created successfully. We have prepared {storage} of storage space for you. Start your image management journey!',
    },
    storageGranted: {
      title: 'Storage Space Increased',
      content: 'Congratulations! Your storage space has been increased by {size}, current total capacity is {total}.',
    },
    bandwidthGranted: {
      title: 'Bandwidth Increased',
      content: 'Congratulations! Your monthly bandwidth has been increased by {size}, current monthly bandwidth is {total}.',
    },
  },
  content: {
    reviewApproved: {
      title: 'Content Review Approved',
      content: 'Your content「{content_name}」has been approved and is now publicly visible.',
    },
    reviewRejected: {
      title: 'Content Review Rejected',
      content: 'Sorry, your content「{content_name}」was not approved. Reason: {reason}',
    },
    reviewPending: {
      title: 'Content Under Review',
      content: 'Your content「{content_name}」has been submitted for review, we will process it soon.',
    },
  },
  storage: {
    quotaWarning: {
      title: 'Storage Space Low',
      content: 'You have used {used_percent}% ({used}/{total}) of your storage space, please clean up or upgrade capacity.',
    },
    quotaIncreased: {
      title: 'Storage Space Upgraded',
      content: 'Your storage space has been upgraded from {old_size} to {new_size}.',
    },
    quotaDecreased: {
      title: 'Storage Space Adjusted',
      content: 'Your storage space has been adjusted from {old_size} to {new_size}.',
    },
  },
  file: {
    deletedByAdmin: {
      title: 'File Deleted',
      content: 'Your file「{file_name}」has been deleted by administrator for violating community rules.',
    },
    batchDeletedByAdmin: {
      title: 'Files Batch Deleted',
      content: 'You have {count} files deleted by administrator for violating community rules.',
    },
    hardDeletedByAdmin: {
      title: 'File Permanently Deleted',
      content: 'Your file「{file_name}」has been permanently deleted by administrator and cannot be recovered.',
    },
    expiryWarning: {
      title: 'You have {count} file(s) about to expire',
      content: 'The following files are about to expire, please backup in time:\n{file_list}\n\nExpired files will be automatically deleted and cannot be recovered.',
    },
    thumbnailFailed: {
      title: 'Thumbnail Generation Failed',
      content: 'Thumbnail generation for file「{file_name}」failed, but the file itself is not affected.',
    },
  },
  security: {
    loginAlert: {
      title: 'New Device Login Alert',
      content: 'Your account was logged in at {time} from {ip} ({location}). If this was not you, please change your password immediately.',
    },
  },
  apikey: {
    created: {
      title: 'API Key Created',
      content: 'You created a new API key「{key_name}」. Please keep your key safe.',
    },
    deleted: {
      title: 'API Key Deleted',
      content: 'API key「{key_name}」has been deleted.',
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
      title: 'API Key Enabled',
      content: 'API key「{key_name}」has been enabled.',
    },
  },
  randomApi: {
    created: {
      title: 'Random Image API Created',
      content: 'You created random image API「{api_name}」.',
    },
    deleted: {
      title: 'Random Image API Deleted',
      content: 'Random image API「{api_name}」has been deleted.',
    },
    disabled: {
      title: 'Random Image API Disabled',
      content: 'Random image API「{api_name}」has been disabled.',
    },
    enabled: {
      title: 'Random Image API Enabled',
      content: 'Random image API「{api_name}」has been enabled.',
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
