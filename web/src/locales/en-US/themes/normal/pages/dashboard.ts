/**
 * Dashboard page texts
 */
export const dashboard = {
  quickActions: {
    title: 'Quick Actions',
    dragHint: 'Drag to reorder',
    actions: {
      myFiles: 'My Files',
      folders: 'Folders',
      tagManage: 'Tag Management',
      categoryManage: 'Category Management',
      openApi: 'Open API',
      automation: 'Automation Tasks',
    },
    settings: {
      divider: 'Settings Center',
      api: 'API',
      profile: 'Profile',
      security: 'Security',
      accessControl: 'Hotlink Protection',
      preferences: 'Preferences',
    },
  },
  stats: {
    title: 'Data Overview',
    items: {
      images: 'Resources',
      storage: 'Storage Used',
      views: 'Views',
      shares: 'Shares',
    },
    quotas: {
      title: 'Quota Overview',
      storage: {
        label: 'Storage Quota',
        used: 'Used',
        total: 'Total',
      },
      bandwidth: {
        label: 'Bandwidth Quota',
        used: 'Used',
        total: 'Total',
      },
      status: {
        healthy: 'Sufficient',
        warning: 'Moderate',
        danger: 'Critical',
      },
      description: {
        healthy: 'Resources are sufficient, you can use with peace of mind',
        warning: 'Resource usage is moderate, consider cleaning up unnecessary content regularly',
        danger: 'Resource usage is critical, please clean up content or contact administrator for expansion',
      },
    },
  },
  messages: {
    title: 'Notifications',
    unread: 'Unread',
    viewAll: 'View All',
    loading: 'Loading messages...',
    empty: 'No messages',
    loadingMore: 'Loading more...',
    loadingMoreDesc: 'Please wait',
    noMore: 'All messages loaded',
    noMoreDesc: 'Total {count} messages',
    scrollHint: 'Scroll for more',
    markReadFailed: 'Failed to mark as read, please try again',
  },
  errors: {
    fetchStatsFailed: 'Failed to fetch data',
    networkFailed: 'Network request failed',
  },
  activityMonitor: {
    title: 'Activity Monitor',
    todayUploads: "Today's Uploads",
    totalViews: 'Total Views',
    recentActivity: 'Recent Activity',
    recordCount: 'Total {count} records',
    loadingActivities: 'Loading activity logs...',
    loadingMore: 'Loading more...',
    loadingHint: 'Please wait',
    allLoaded: 'All activity records loaded',
    empty: 'No activity records',
    emptyHint: 'Your activities will appear here once you start using',
    status: {
      active: 'Monitoring',
      warning: 'Warning',
      error: 'Error',
      offline: 'Offline',
    },
    toast: {
      fetchFailed: 'Failed to fetch activity logs',
      noData: 'Unable to get real-time activity data temporarily',
    },
    timeAgo: {
      justNow: 'Just now',
      hoursAgo: '{hours} hours ago',
      minutesAgo: '{minutes} minutes ago',
      daysAgo: '{days} days ago',
    },
  },
  recentUploads: {
    title: 'Recent Uploads',
    loading: 'Loading...',
    empty: 'No upload records',
    toast: {
      noImages: 'No images to copy',
      copySuccess: 'Copied {count} image links',
      copyFailed: 'Copy failed',
      fetchFailed: 'Failed to fetch data',
    },
    actions: {
      copyAll: 'Copy All Links',
      viewAll: 'View All',
    },
  },
  uploadQueue: {
    title: 'Upload Queue',
    statsTitle: 'Upload Statistics',
    totalFiles: 'Total {count} files',
    overallProgress: 'Overall Progress',
    queueTotal: 'Queue Total',
    uploading: 'Uploading',
    pending: 'Pending',
    completed: 'Completed',
    failed: 'Failed',
    actions: {
      goToUpload: 'Go to Upload Page',
    },
  },
  uploadEntry: {
    title: {
      default: 'Upload Files',
      withCount: '{count} files in queue',
    },
  },
  folderDistribution: {
    title: 'Folder Distribution',
    description: 'Block size represents storage space, color depth represents proportion',
    totalFolders: 'Total folders: {count}',
    empty: 'No folders',
    tooltipSize: 'Size: {size}',
    tooltipFiles: 'Files: {count}',
    tooltipPercent: 'Proportion: {percent}%',
    unnamed: 'Unnamed',
  },
}
