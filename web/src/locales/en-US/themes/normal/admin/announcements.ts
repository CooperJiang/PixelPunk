export const announcements = {
  title: 'Announcement Management',
  subtitle: 'Manage system announcements. Supports Markdown and HTML.',

  actions: {
    create: 'Create Announcement',
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
    publishedAt: 'Published At',
    createdAt: 'Created At',
    actions: 'Actions',
    noSummary: 'No summary',
    notSet: 'Not set',
    loadingText: 'Loading...',
  },

  empty: {
    title: 'No announcements',
    description: 'Click "Create Announcement" to add your first one',
  },

  form: {
    title: {
      label: 'Title',
      placeholder: 'Enter announcement title',
    },
    summary: {
      label: 'Summary',
      placeholder: 'Enter announcement summary (optional)',
    },
    content: {
      label: 'Content',
      placeholder: 'Enter content in Markdown...',
      hint: 'Markdown supported with live preview. Drag/paste images to upload.',
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
    createTitle: 'Create Announcement',
    editTitle: 'Edit Announcement',
  },

  messages: {
    fetchError: 'Failed to load announcements: ',
    deleteConfirm: 'Delete this announcement? This cannot be undone.',
    deleteSuccess: 'Deleted successfully',
    deleteError: 'Delete failed: ',
    pinSuccess: '{action} succeeded',
    pinError: '{action} failed: ',
    createSuccess: 'Announcement created',
    createError: 'Create failed: ',
    updateSuccess: 'Announcement updated',
    updateError: 'Update failed: ',
    operationError: 'Operation failed: ',
    uploadSuccess: 'Uploaded {count} image(s)',
    uploadError: 'Image upload failed: ',
    uploadFailed: 'Upload failed',
  },

  validation: {
    titleRequired: 'Please enter a title',
    contentRequired: 'Please enter content',
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
    systemToggle: 'System Toggle',
    drawerPosition: {
      label: 'Drawer Position',
      left: 'Left',
      right: 'Right',
    },
    displayLimit: 'Display Count',
    autoShowDelay: 'Popup Delay',
    saveButton: 'Save Settings',
    messages: {
      fetchError: 'Failed to load settings:',
      limitRange: 'Display count must be between 1 and 20',
      saveSuccess: 'Settings saved',
      saveError: 'Save failed: ',
    },
    units: {
      second: 's',
    },
  },
}
