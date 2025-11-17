/**
 * Share Management Page Text - Cyber Style
 */
export const shares = {
  title: 'My Share Links',
  search: 'Search Share Links',
  page: {
    title: 'My Share Links',
    subtitle: 'Manage and view all share links you have created',
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
    viewCount: '{count} accesses',
    maxViews: 'Limit {max} times',
    passwordProtected: 'Password Protected',
    collectVisitor: 'Collect Visitor',
    accessNotification: 'Access Notification',
    copyLink: 'Copy Link',
    folderCount: '{count} data folders',
    fileCount: '{count} data units',
    viewShare: 'View Share',
    visitorInfo: 'Visitor Info',
    delete: 'Delete',
  },
  empty: {
    title: 'No Share Content',
    subtitle: 'You have not created any share links yet, you can create shares from folders page',
    goToFolders: 'Go to Folders',
  },
  delete: {
    confirm: 'Are you sure you want to delete share link {name}?',
    forceDelete: 'Permanently Delete (irreversible)',
    cancel: 'Cancel',
    confirmButton: 'Delete',
  },
  unnamed: 'Unnamed Share',
  toast: {
    loadFailed: 'Failed to load share link list',
    linkCopied: 'Share link copied to clipboard',
    deleteSuccess: 'Share link deleted successfully',
    deleteFailed: 'Failed to delete share link',
  },
  dialog: {
    deleteTitle: 'Confirm Deletion',
    deleteMessage: 'Are you sure you want to delete share link "{name}"?',
    deleteHint: 'Cannot be recovered after deletion',
  },
  visitor: {
    dialogTitle: 'Visitor Info List',
    search: 'Search visitor name, email or IP node',
    deleteConfirm: 'Are you sure you want to delete this visitor info?',
    deleteHint: 'Cannot be recovered after deletion',
    loadFailed: 'Failed to load visitor info, please retry',
    deleteSuccess: 'Visitor info deleted successfully',
    deleteFailed: 'Failed to delete visitor info, please retry',
    table: {
      name: 'Visitor Name',
      email: 'Email',
      ip: 'IP Address',
      visits: 'Access Count',
      lastVisit: 'Last Visit Time',
      actions: 'Actions',
    },
    pagination: {
      info: 'Total {total} items, currently showing {from}-{to}',
    },
    empty: {
      title: 'No Visitor Info',
      desc: 'No visitors have submitted info yet, or your filter criteria did not match any results',
    },
  },
}
