/**
 * Message System i18n - Cyber Theme
 */
export const message = {
  types: {
    system: {
      maintenance: 'System maintenance',
      update: 'System update',
      announcement: 'System announcement',
    },
    account: {
      register: 'Registration successful',
      storage_granted: 'Storage quota granted',
      bandwidth_granted: 'Bandwidth quota granted',
    },
    content: {
      review_approved: 'Content review approved',
      review_rejected: 'Content review rejected',
      review_pending: 'Pending review',
    },
    storage: {
      quota_warning: 'Storage quota warning',
      quota_exceeded: 'Storage quota exceeded',
      quota_increased: 'Storage quota increased',
      quota_decreased: 'Storage quota adjusted',
    },
    file: {
      deleted_by_admin: 'Data deletion',
      batch_deleted_by_admin: 'Batch data deletion',
      hard_deleted_by_admin: 'Permanent data deletion',
      expiry_warning: 'Data expiry reminder',
      thumbnail_failed: 'Thumbnail generation failed',
    },
    security: {
      login_alert: 'Login alert',
      password_changed: 'Password changed',
    },
    apikey: {
      created: 'API key created',
      deleted: 'API key deleted',
      regenerated: 'API key regenerated',
      disabled: 'API key disabled',
      enabled: 'API key enabled',
    },
    random_api: {
      created: 'Random API created',
      deleted: 'Random API deleted',
      disabled: 'Random API disabled',
      enabled: 'Random API enabled',
    },
    share: {
      expiry_warning: 'Share expiring soon',
    },
  },
  priority: {
    high: 'High',
    normal: 'Normal',
    low: 'Low',
  },
  status: {
    unread: 'Unread',
    read: 'Read',
    deleted: 'Deleted',
    unknown: 'Unknown',
  },
  time: {
    justNow: 'Just now',
    minutesAgo: '{n} minutes ago',
    hoursAgo: '{n} hours ago',
    daysAgo: '{n} days ago',
  },
}
