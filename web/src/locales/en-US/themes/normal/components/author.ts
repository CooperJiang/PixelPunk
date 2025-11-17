/**
 * Author info component
 */
export const author = {
  loading: 'Loading...',
  header: {
    joinedLabel: 'Joined {days} days',
    websiteLabel: 'Website',
    stats: {
      views: 'Views',
      shares: 'Shares',
      images: 'Files',
    },
  },
  folders: {
    rootTitle: 'Folders',
    rootSubtitle: "Browse author's folders",
    childTitle: 'Subfolders',
    childSubtitle: 'Subfolders in current folder',
    countLabel: 'files',
    createdLabel: '{time}',
  },
  images: {
    sectionTitle: 'Files',
    sectionSubtitle: "Browse author's file works",
    countLabel: 'resources',
    loadingText: 'Loading files...',
    paginationInfo: 'Page {current} of {total}',
    actions: {
      preview: {
        label: 'Preview',
        tooltip: 'Click to preview file',
      },
      download: {
        label: 'Download',
        tooltip: 'Download original',
      },
      detail: {
        label: 'Details',
        tooltip: 'Copy file link',
      },
    },
  },
  shares: {
    sectionTitle: 'Shares',
    sectionSubtitle: "Browse author's public shares",
    statusBadge: 'Status',
    meta: {
      viewsLabel: 'Views',
      createdLabel: '{time}',
      keyLabel: 'Share Key',
    },
    empty: {
      title: 'No Shares',
      description: 'This author has not created any public shares yet',
    },
  },
  empty: {
    title: 'No Content',
    description: 'This author has not uploaded any content yet',
    suggestion: 'Please check back later',
    retry: 'Reload',
  },
  error: {
    title: 'Failed to Load',
    retry: 'Retry',
    invalidAuthor: 'Invalid author ID',
    loadFailed: 'Failed to load author information',
    loadFolderFailed: 'Failed to load folder content',
  },
  toolbar: {
    searchPlaceholder: 'Search content...',
    refresh: 'Refresh',
    gridView: 'Grid View',
    listView: 'List View',
    fullscreenEnter: 'Fullscreen',
    fullscreenExit: 'Exit Fullscreen',
  },
  breadcrumb: {
    root: 'Root',
  },
  toast: {
    copySuccess: 'Link copied to clipboard',
    copyFailed: 'Copy failed',
  },
}
