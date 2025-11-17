/**
 * Message system i18n configuration - Normal theme
 */
export const message = {
  types: {
    system: {
      maintenance: 'System Maintenance',
      update: 'System Update',
      announcement: 'System Announcement',
    },
    account: {
      register: 'Registration Welcome',
      storage_granted: 'Storage Granted',
      bandwidth_granted: 'Bandwidth Granted',
    },
    content: {
      review_approved: 'Review Approved',
      review_rejected: 'Review Rejected',
      review_pending: 'Pending Review',
    },
    storage: {
      quota_warning: 'Storage Warning',
      quota_exceeded: 'Storage Exceeded',
      quota_increased: 'Storage Increased',
      quota_decreased: 'Storage Adjusted',
    },
    file: {
      deleted_by_admin: 'File Deleted',
      batch_deleted_by_admin: 'Batch Deleted',
      hard_deleted_by_admin: 'Permanently Deleted',
      expiry_warning: 'File Expiry Warning',
      thumbnail_failed: 'Thumbnail Generation Failed',
    },
    security: {
      login_alert: 'Login Alert',
      password_changed: 'Password Changed',
    },
    apikey: {
      created: 'API Key Created',
      deleted: 'API Key Deleted',
      regenerated: 'API Key Regenerated',
      disabled: 'API Key Disabled',
      enabled: 'API Key Enabled',
    },
    random_api: {
      created: 'Random Image API Created',
      deleted: 'Random Image API Deleted',
      disabled: 'Random Image API Disabled',
      enabled: 'Random Image API Enabled',
    },
    share: {
      expiry_warning: 'Share Expiring Soon',
    },
  },
  priority: {
    high: 'High',
    normal: 'Medium',
    low: 'Low',
  },
  status: {
    unread: 'Unread',
    read: 'Read',
    deleted: 'Deleted',
    unknown: 'Unknown',
  },
  time: {
    justNow: 'Just Now',
    minutesAgo: '{n} minutes ago',
    hoursAgo: '{n} hours ago',
    daysAgo: '{n} days ago',
  },
}
