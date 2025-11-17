/**
 * Author Page Text - Cyber Style
 */
export const author = {
  loading: 'Loading data...',

  header: {
    joinedLabel: 'Connected Time',
    websiteLabel: 'Network Node',
    stats: {
      views: 'Access Count',
      shares: 'Share Count',
      images: 'Data Units',
    },
  },

  folders: {
    rootTitle: 'All Directories',
    rootSubtitle: 'Browse all data directories by author',
    childTitle: 'Subdirectories',
    childSubtitle: 'Subdirectories in current directory',
    countLabel: ' units',
    createdLabel: 'Created At',
    unknownTime: 'Unknown',
    rootFolder: 'Root Directory',
  },

  shares: {
    sectionTitle: 'Public Shares',
    sectionSubtitle: 'Author\'s public data shares',
    statusBadge: 'Status',
    meta: {
      viewsLabel: 'Access',
      createdLabel: 'Created At',
      keyLabel: 'Share Key',
    },
    empty: {
      title: 'No Shares',
      description: 'This author has not created any public shares yet',
    },
  },

  toolbar: {
    searchPlaceholder: 'Search directories...',
    refresh: 'Refresh',
    fullscreenEnter: 'Fullscreen Mode',
    fullscreenExit: 'Exit Fullscreen',
  },

  empty: {
    title: 'No Data',
    description: 'This author has not uploaded any data yet',
    suggestion: 'Please check back later',
    retry: 'Retry',
    folderEmpty: {
      title: 'Directory Empty',
      description: 'This directory has no content at the moment',
      suggestion: 'Go back or browse other content',
    },
  },

  error: {
    title: 'Load Failed',
    retry: 'Retry',
    invalidAuthor: 'Invalid author ID',
    loadFailed: 'Failed to load author info',
    loadFolderFailed: 'Failed to load directory',
  },

  toast: {
    copySuccess: 'Copy successful',
    copyFailed: 'Copy failed',
  },
}
