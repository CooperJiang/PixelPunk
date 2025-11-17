/**
 * Announcement Management Module
 */
export const announcements = {
  title: '█ Announcement Management █',
  subtitle: 'Manage system announcements, supports Markdown & HTML',

  actions: {
    create: 'Create announcement',
    refresh: 'Refresh',
    filter: 'Filter',
    edit: 'Edit',
    delete: 'Delete',
    pin: 'Pin',
    unpin: 'Unpin',
    cancel: 'Cancel',
    update: 'Update',
  },

  status: {
    draft: 'Draft',
    published: 'Published',
    archived: 'Archived',
  },

  table: {
    title: 'Title',
    summary: 'Summary',
    status: 'Status',
    viewCount: 'Views',
    publishedAt: 'Published',
    createdAt: 'Created',
    actions: 'Actions',
    noSummary: 'No summary',
    notSet: 'Not set',
    loadingText: '>>> Loading...',
  },

  empty: {
    title: 'No announcements',
    description: 'Click "Create announcement" to create your first announcement',
  },

  form: {
    title: {
      label: 'Title',
      placeholder: 'Enter announcement title',
    },
    summary: {
      label: 'Summary',
      placeholder: 'Enter summary (optional)',
    },
    content: {
      label: 'Content',
      placeholder: 'Enter Markdown content...',
      hint: 'Supports Markdown with live preview. Drag or paste images to upload',
    },
    status: {
      label: 'Status',
    },
    options: {
      label: 'Options',
      pinned: 'Pin to top',
      pinnedHint: 'Auto popup on first visit',
    },
  },

  dialog: {
    createTitle: '◢ Create Announcement',
    editTitle: '◢ Edit Announcement',
  },

  messages: {
    fetchError: 'Failed to load announcements: ',
    deleteConfirm: 'Delete this announcement? This action cannot be undone.',
    deleteSuccess: '✓ Deleted successfully',
    deleteError: '✗ Delete failed: ',
    pinSuccess: '✓ {action} successful',
    pinError: '✗ {action} failed: ',
    createSuccess: '✓ Announcement created',
    createError: '✗ Create failed: ',
    updateSuccess: '✓ Announcement updated',
    updateError: '✗ Update failed: ',
    operationError: '✗ Operation failed: ',
    uploadSuccess: '✓ {count} images uploaded',
    uploadError: '✗ Image upload failed: ',
    uploadFailed: 'Upload failed',
  },

  validation: {
    titleRequired: 'Title required',
    contentRequired: 'Content required',
  },

  filter: {
    status: {
      label: 'Status',
      all: 'All',
    },
    pinned: {
      label: 'Pinned',
      all: 'All',
      onlyPinned: 'Pinned only',
      notPinned: 'Not pinned',
    },
    search: {
      label: 'Search',
      placeholder: 'Search title or summary...',
    },
    actions: {
      query: 'Query',
      reset: 'Reset',
    },
  },

  settings: {
    systemToggle: 'System switch',
    drawerPosition: {
      label: 'Drawer position',
      left: 'Left',
      right: 'Right',
    },
    displayLimit: 'Display limit',
    autoShowDelay: 'Popup delay',
    saveButton: 'Save config',
    messages: {
      fetchError: 'Failed to fetch config:',
      limitRange: 'Display limit must be between 1-20',
      saveSuccess: '✓ Config saved',
      saveError: '✗ Save failed: ',
    },
    units: {
      second: 'sec',
    },
  },
}
