export const users = {
  title: 'Node Management',
  subtitle: 'Manage and monitor all connected nodes',

  columns: {
    user: 'Node',
    email: 'Communication address',
    registeredAt: 'Access time',
    lastActivity: 'Last activity',
    status: 'Status',
    role: 'Permissions',
    storage: 'Storage/Bandwidth',
    actions: 'Operations',
    userId: 'ID',
    neverActive: 'Never active',
    storageLabel: 'Storage',
    bandwidthLabel: 'Bandwidth',
  },

  status: {
    all: 'All',
    active: 'Online',
    disabled: 'Disabled',
    suspended: 'Suspended',
  },

  role: {
    all: 'All',
    superAdmin: 'Super admin',
    admin: 'Admin node',
    user: 'Normal node',
    guest: 'Guest node',
    unknown: 'Unknown',
  },

  roles: {
    all: 'All',
    superAdmin: 'Super admin',
    admin: 'Admin node',
    normal: 'Normal node',
    guest: 'Guest node',
    unknown: 'Unknown',
  },

  actions: {
    view: 'View details',
    edit: 'Edit node',
    disable: 'Disable node',
    enable: 'Enable node',
    delete: 'Delete node',
    export: 'Export',
    import: 'Import',
    refresh: 'Refresh',
  },

  buttons: {
    createUser: 'New node',
    create: 'Create node',
    batchOperation: 'Batch operation',
    filter: 'Filter',
    resetFilter: 'Reset filter',
    applyFilter: 'Apply filter',
    cancel: 'Cancel',
    confirm: 'Confirm',
    save: 'Save',
    close: 'Close',
  },

  toast: {
    userCreated: 'Node created successfully',
    userUpdated: 'Node info updated',
    userDeleted: 'Node deleted successfully',
    statusChanged: 'Status changed',
    batchOperationSuccess: 'Batch operation complete',
    operationFailed: 'Operation failed',
    pleaseSelectUsers: 'Please select nodes to operate',
    userEnabled: 'Node enabled',
    userDisabled: 'Node disabled',
    userDataNotFound: 'Node data not found',
    success: '{operation} successful',
    error: '{operation} failed',
    exportSuccess: 'Node data exported successfully',
    exportError: 'Node data export failed',
  },

  confirm: {
    enableUser: 'Enable this node?',
    disableUser: 'Disable this node?',
    deleteUser: 'Delete this node?',
    delete: 'Delete selected nodes?',
    batchDelete: 'Delete the selected {count} nodes?',
    batchDisable: 'Disable selected nodes?',
    batchEnable: 'Enable the selected {count} nodes?',
    disableTitle: 'About to disable {count} nodes',
    disableMessage1: 'After disabling, these nodes will be unable to access the system',
    disableMessage2: 'Currently online node sessions will terminate immediately',
    disableMessage3: 'All related protocol links will become invalid',
  },

  disableDialog: {
    title: 'Disable node confirmation',
    warning: 'You are about to disable node',
    consequences: 'This operation will result in:',
    consequence1: 'Node cannot login to system',
    consequence2: 'All node sessions will terminate immediately',
    consequence3: 'All node protocol links will become invalid',
    consequence4: 'Node cannot access its data',
    inputLabel: 'Enter node name to confirm',
    inputPlaceholder: 'Enter node name',
    inputHint: 'Enter "{username}" to confirm disable',
    usernameNotMatch: 'Node name does not match',
    confirmButton: 'Confirm disable',
    cancelButton: 'Cancel',
  },

  filter: {
    title: 'Filter nodes',
    keyword: 'Keyword',
    keywordPlaceholder: 'Search node name or communication address',
    status: {
      label: 'Node status',
      all: 'All nodes',
      normal: 'Normal nodes',
      disabled: 'Disabled nodes',
      deleted: 'Deleted nodes',
    },
    role: {
      label: 'Node permissions',
      all: 'All permissions',
      superAdmin: 'Super admin nodes',
      admin: 'Admin nodes',
      user: 'Normal nodes',
    },
    registrationTime: 'Access time',
    registrationDate: 'Access date',
    dateRangePlaceholder: 'Select time range',
    startDatePlaceholder: 'Start date',
    endDatePlaceholder: 'End date',
    startDate: 'Start date',
    endDate: 'End date',
    sortBy: 'Sort by',
    sort: {
      label: 'Sort by',
      placeholder: 'Select sort method',
      default: 'Default',
      newest: 'Newest access',
      oldest: 'Oldest access',
      lastActivityNewest: 'Recently active',
      lastActivityOldest: 'Earliest active',
      usernameAZ: 'Node name A-Z',
      usernameZA: 'Node name Z-A',
    },
  },

  batch: {
    title: 'Batch operations',
    selected: 'Selected {count} nodes',
    selectAll: 'Select all',
    deselectAll: 'Deselect all',
    clearSelection: 'Clear selection',
    enable: 'Enable',
    disable: 'Disable',
    delete: 'Delete',
    export: 'Export',
    sendEmail: 'Send communication',
    role: 'Set permissions',
    apply: 'Apply',
    cancel: 'Cancel',
    noSelection: 'Please select nodes to operate',
    setRole: 'Set permissions',
    selectRole: 'Please select permissions',
    selectUsersToExport: 'Please select nodes to export',
    confirmDisable: 'Confirm disable',
  },

  form: {
    createTitle: 'Create new node',
    editTitle: 'Edit node',
    username: 'Node name',
    usernamePlaceholder: 'Enter node name',
    usernameHint: '3-20 characters, only letters, numbers and underscores',
    email: 'Communication address',
    emailPlaceholder: 'Enter communication address',
    emailHint: 'Used for login and receiving notifications',
    password: 'Password',
    passwordPlaceholder: 'Enter password',
    passwordHint: 'At least 8 characters',
    confirmPassword: 'Confirm password',
    confirmPasswordPlaceholder: 'Enter password again',
    role: 'Permissions',
    rolePlaceholder: 'Select node permissions',
    roleUser: 'Normal node',
    roleAdmin: 'Admin node',
    roleSuperAdmin: 'Super admin node',
    status: 'Status',
    statusActive: 'Normal node',
    statusDisabled: 'Disabled node',
    storageConfig: 'Storage config',
    storageLimit: 'Storage limit',
    storageLimitPlaceholder: 'Enter storage limit',
    storageLimitDefault: 'Default 5GB',
    storageLimitHint: 'Unit: GB, 0 means unlimited',
    bandwidthLimit: 'Bandwidth limit',
    bandwidthLimitPlaceholder: 'Enter bandwidth limit',
    bandwidthLimitDefault: 'Default 100GB',
    bandwidthLimitHint: 'Unit: GB, 0 means unlimited',
    avatar: 'Avatar',
    avatarUpload: 'Upload avatar',
    usernameRequired: 'Node name cannot be empty',
    usernameMinLength: 'Node name at least 2 characters',
    usernameMaxLength: 'Node name cannot exceed 20 characters',
    emailRequired: 'Communication address cannot be empty',
    emailInvalid: 'Communication address format incorrect',
    passwordRequired: 'Password cannot be empty',
    passwordMinLength: 'Password at least 6 characters',
    passwordMaxLength: 'Password cannot exceed 20 characters',
    required: 'This field is required',
    invalidEmail: 'Communication address format incorrect',
    passwordTooShort: 'Password needs at least 8 characters',
    passwordNotMatch: 'Passwords do not match',
  },

  detail: {
    title: 'Node details',
    loading: 'Loading data...',
    noData: 'No data',
    basicInfo: 'Basic info',
    accountInfo: 'Account info',
    statisticsInfo: 'Statistics info',
    activityInfo: 'Activity info',
    usageStats: 'Usage statistics',
    storageManagement: 'Storage management',
    quickActions: 'Quick actions',
    userId: 'Node ID',
    username: 'Node name',
    email: 'Communication address',
    role: 'Permissions',
    registeredAt: 'Access time',
    lastLogin: 'Last login',
    lastLoginAt: 'Last login',
    lastActivity: 'Last activity',
    lastActivityAt: 'Last activity',
    lastActivityIp: 'IP address',
    uploadedFiles: 'Uploaded data count',
    shareCount: 'Protocol count',
    totalViews: 'Total views',
    usedStorage: 'Used storage',
    storageLimit: 'Storage limit',
    monthlyBandwidth: 'Monthly bandwidth',
    bandwidthLimit: 'Bandwidth limit',
    fileCount: 'Data count',
    unlimited: 'Unlimited',
    never: 'Never',
    currentSetting: 'Current config',
    enterValue: 'Enter value',
    unit: 'Unit',
    updateStorageSettings: 'Update storage config',
    status: {
      label: 'Status',
      active: 'Online',
      disabled: 'Disabled',
      deleted: 'Deleted',
      unknown: 'Unknown',
    },
    roles: {
      superAdmin: 'Super admin node',
      admin: 'Admin node',
      user: 'Normal node',
      unknown: 'Unknown',
    },
    units: {
      mb: 'MB',
      gb: 'GB',
      tb: 'TB',
    },
    actions: {
      resetPassword: 'Reset password',
      resetPasswordDesc: 'Generate new password for this node',
      sendEmail: 'Send email',
      sendEmailDesc: 'Send email to this node',
      enableUser: 'Enable node',
      enableUserDesc: 'Allow this node to login',
      disableUser: 'Disable node',
      disableUserDesc: 'Prevent this node from login',
    },
    messages: {
      storageUpdateSuccess: 'Storage config updated',
      passwordResetSuccess: 'Password reset to: {password}',
      emailSentSuccess: 'Email sent successfully',
      userEnabled: 'Node enabled',
      userDisabled: 'Node disabled',
    },
    errors: {
      fetchFailed: 'Get node info failed',
      storageUpdateFailed: 'Update storage config failed',
      passwordResetFailed: 'Reset password failed',
      emailSentFailed: 'Send email failed',
      operationFailed: 'Operation failed',
    },
  },

  search: {
    placeholder: 'Search node name or communication address...',
    button: 'Search',
  },

  sendEmail: {
    title: 'Send communication',
    recipients: 'Recipient nodes',
    recipientsHint: 'Selected {count} nodes',
    subject: 'Subject',
    subjectPlaceholder: 'Enter communication subject',
    content: 'Content',
    contentPlaceholder: 'Enter communication content',
    send: 'Send',
    cancel: 'Cancel',
    sending: 'Sending...',
    sent: 'Communication sent',
    failed: 'Send failed',
    templates: 'Quick templates',
    templateButtons: {
      welcome: 'Welcome template',
      warning: 'Warning template',
      notice: 'Notice template',
      custom: 'Custom template',
    },
    errors: {
      subjectRequired: 'Communication subject cannot be empty',
      subjectTooLong: 'Communication subject length cannot exceed 100 characters',
      contentRequired: 'Communication content cannot be empty',
      contentTooLong: 'Communication content length cannot exceed 5000 characters',
    },
    messages: {
      sent: 'Communication send request submitted',
      failed: 'Submission failed, please retry',
    },
    templatesContent: {
      welcome: {
        subject: 'Welcome to our platform',
        content: `Dear User,

Welcome to our image hosting platform!

We are pleased that you chose our service. If you encounter any issues during usage, please feel free to contact our customer service team.

Best wishes!

Sincerely,
Management Team`,
      },
      warning: {
        subject: 'Account security reminder',
        content: `Dear User,

We noticed some security issues with your account that need attention:

1. Please ensure your password is secure enough
2. Do not share your account information with others
3. If abnormal activity is detected, please contact us immediately

If you have questions, please contact customer service.

Sincerely,
Security Team`,
      },
      notice: {
        subject: 'System notice',
        content: `Dear User,

We want to notify you of the following important information:

[Please fill in specific notice content here]

Thank you for your attention.

Sincerely,
Management Team`,
      },
      custom: {
        content: 'Dear User,\n\n[Please fill in your email content here]\n\nSincerely,\nManagement Team',
      },
    },
  },

  empty: {
    noUsers: 'No nodes',
    noSearchResults: 'No nodes matching conditions found',
    createFirst: 'Create first node',
    title: 'No nodes matching conditions found',
    description: 'Try adjusting filter conditions or create new nodes',
    resetFilter: 'Reset filter conditions',
  },

  pagination: {
    total: 'Total {total} nodes',
    pageSize: '{size} per page',
  },

  statusText: {
    normal: 'Online',
    disabled: 'Disabled',
    deleted: 'Deleted',
    unknown: 'Unknown',
  },

  loading: {
    users: 'Loading node data...',
  },

  disableWarning: {
    title: 'Disable node confirmation',
    warning: 'Are you sure you want to disable node "{username}"?',
    consequence1: 'After disabling, this node will immediately lose all operation permissions',
    consequence2: 'If the node is currently online, it will be forced to exit the system',
    consequence3: 'Node needs to wait for re-enabling before login again',
    confirmButton: 'Confirm disable',
  },

  editDialog: {
    title: 'Edit node',
    saveChanges: 'Save changes',
  },

  export: {
    id: 'Node ID',
    username: 'Node name',
    email: 'Communication address',
    status: 'Status',
    role: 'Permissions',
    createdAt: 'Access time',
    updatedAt: 'Last updated',
    filename: 'Node data',
  },
}
