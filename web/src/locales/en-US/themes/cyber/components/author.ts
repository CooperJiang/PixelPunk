/**
 * Author Info Component
 */
export const author = {
  loading: 'Loading data...',
  header: {
    joinedLabel: 'Connected {days} days',
    websiteLabel: 'Network Address',
    stats: {
      views: 'Access Count',
      shares: 'Share Data',
      images: 'File Storage',
    },
  },
  folders: {
    rootTitle: 'Data Container',
    rootSubtitle: 'Browse user data storage structure',
    childTitle: 'Sub-container',
    childSubtitle: 'Sub-level storage units in current container',
    countLabel: ' data units',
    createdLabel: '{time}',
  },
  images: {
    sectionTitle: 'File Database',
    sectionSubtitle: 'Browse user file dataset',
    countLabel: ' resource data',
    loadingText: 'Loading file data...',
    paginationInfo: 'Section {current} of {total}',
    actions: {
      preview: {
        label: 'Data Preview',
        tooltip: 'Launch file preview mode',
      },
      download: {
        label: 'Data Download',
        tooltip: 'Get original file data',
      },
      detail: {
        label: 'Details',
        tooltip: 'Copy data link',
      },
    },
  },
  shares: {
    sectionTitle: 'Share Network',
    sectionSubtitle: 'Browse user public data shares',
    statusBadge: 'Connection Status',
    meta: {
      viewsLabel: 'Access Count',
      createdLabel: '{time}',
      keyLabel: 'Access Key',
    },
    empty: {
      title: 'Network Space Empty',
      description: 'User has not established any public data connections',
    },
  },
  empty: {
    title: 'Database Empty',
    description: 'User has not uploaded any data content',
    suggestion: 'Rescan later',
    retry: 'Reconnect',
  },
  error: {
    title: 'Connection Failed',
    retry: 'Retry',
    invalidAuthor: 'Invalid user ID',
    loadFailed: 'User profile load failed',
    loadFolderFailed: 'Data container load failed',
  },
  toolbar: {
    searchPlaceholder: 'Data retrieval...',
    refresh: 'Resync',
    gridView: 'Matrix View',
    listView: 'List View',
    fullscreenEnter: 'Fullscreen Mode',
    fullscreenExit: 'Exit Fullscreen',
  },
  breadcrumb: {
    root: 'Root Directory',
  },
  toast: {
    copySuccess: 'Data link copied to buffer',
    copyFailed: 'Data copy operation failed',
  },
}
