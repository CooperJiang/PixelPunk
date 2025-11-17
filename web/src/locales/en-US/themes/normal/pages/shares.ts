/**
 * Share management page texts
 */
export const shares = {
  title: 'My Shares',
  search: 'Search Shares',
  page: {
    title: 'My Shares',
    subtitle: 'Manage and view all share links you created',
  },
  filter: {
    all: 'All',
    normal: 'Normal',
    expired: 'Expired',
    deleted: 'Deleted',
    disabled: 'Disabled',
  },
  status: {
    normal: 'Normal',
    expired: 'Expired',
    deleted: 'Deleted',
    disabled: 'Disabled',
    unknown: 'Unknown',
  },
  card: {
    viewCount: '{count} views',
    maxViews: 'Limit {max} views',
    passwordProtected: 'Password Protected',
    collectVisitor: 'Collect Visitors',
    accessNotification: 'Access Notification',
    copyLink: 'Copy Link',
    folderCount: '{count} folders',
    fileCount: '{count} resources',
    viewShare: 'View Share',
    visitorInfo: 'Visitor Info',
    delete: 'Delete',
  },
  empty: {
    title: 'No Share Content',
    subtitle: 'You have not created any shares yet, you can create shares from the folders page',
    goToFolders: 'Go to Folders',
  },
  delete: {
    confirm: 'Are you sure you want to delete share {name}?',
    forceDelete: 'Permanently Delete (Cannot be recovered)',
    cancel: 'Cancel',
    confirmButton: 'Delete',
  },
  unnamed: 'Unnamed Share',
  toast: {
    loadFailed: 'Failed to load share list',
    linkCopied: 'Share link copied to clipboard',
    deleteSuccess: 'Share deleted successfully',
    deleteFailed: 'Failed to delete share',
  },
  dialog: {
    deleteTitle: 'Confirm Delete',
    deleteMessage: 'Are you sure you want to delete share "{name}"?',
    deleteHint: 'Cannot be recovered after deletion',
  },
  visitor: {
    dialogTitle: 'Visitor Information List',
    search: 'Search visitor name, email or IP address',
    deleteConfirm: 'Are you sure you want to delete this visitor information?',
    deleteHint: 'Cannot be recovered after deletion',
    loadFailed: 'Failed to load visitor information, please try again',
    deleteSuccess: 'Visitor information deleted successfully',
    deleteFailed: 'Failed to delete visitor information, please try again',
    table: {
      name: 'Visitor Name',
      email: 'Email',
      ip: 'IP Address',
      visits: 'Visit Count',
      lastVisit: 'Last Visit Time',
      actions: 'Actions',
    },
    pagination: {
      info: 'Total {total} items, showing {from}-{to}',
    },
    empty: {
      title: 'No Visitor Information',
      desc: 'No visitors have submitted information yet, or your filter conditions have no matching results',
    },
  },
}
