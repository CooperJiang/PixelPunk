export const shares = {
  title: 'Share Management',
  subtitle: 'Manage user share links and visitor analytics',
  tabs: {
    list: 'Share List',
    visitors: 'Visitors',
    stats: 'Statistics',
  },
  shareList: {
    buttons: {
      refresh: 'Refresh',
      filter: 'Filter',
    },
    columns: {
      name: 'Name',
      shareKey: 'Share Key',
      expiredAt: 'Expiry',
      views: 'Views',
      createdAt: 'Created At',
      status: 'Status',
      actions: 'Actions',
    },
    actions: {
      copyLink: 'Copy Link',
      viewContent: 'View Content',
      disable: 'Disable',
      enable: 'Enable',
      delete: 'Delete',
      viewIpInfo: 'View IP Info',
    },
    status: {
      normal: 'Normal',
      expired: 'Expired',
      deleted: 'Deleted',
      disabled: 'Disabled',
      unknown: 'Unknown',
    },
    format: {
      unlimited: 'Unlimited',
      infiniteViews: '∞',
      viewsFormat: '{current} / {max}',
    },
    messages: {
      fetchFailed: 'Failed to fetch shares',
      updateStatusSuccess: 'Share status updated',
      updateStatusFailed: 'Failed to update status',
      deleteSuccess: 'Share deleted',
      deleteFailed: 'Delete failed',
      copySuccess: 'Share link copied',
      copySuccessNormal: 'Copied (normal link, may require password)',
      copyFailed: 'Copy failed, please copy manually',
      getAccessFailed: 'Failed to get admin access, password may be required',
    },
    loading: {
      text: 'Loading shares...',
    },
    empty: {
      title: 'No shares',
      description: 'No share links yet. Newly created shares will appear here.',
      resetButton: 'Reset Filters',
    },
    statusDialog: {
      titleEnable: 'Enable Share',
      titleUpdate: 'Update Status',
      message: 'Change status to',
      messageHighlight: '{status}',
      reasonLabel: 'Reason (optional)',
      reasonPlaceholder: 'Enter reason...',
      cancel: 'Cancel',
      confirm: 'Confirm',
    },
    deleteDialog: {
      title: 'Confirm Delete',
      message: 'Delete share',
      messageHighlight: '{name}',
      forceDelete: 'Hard delete (irreversible)',
      cancel: 'Cancel',
      confirm: 'Delete',
    },
  },
  visitorList: {
    ipInfoTitle: 'View IP Info',
    filter: {
      shareIdPlaceholder: 'Filter by Share ID',
      keywordPlaceholder: 'Search name, email or IP...',
      searchButton: 'Search',
    },
    columns: {
      id: 'ID',
      visitorName: 'Visitor',
      visitorEmail: 'Email',
      createdAt: 'Visited At',
      ipAddress: 'IP Address',
      shareKey: 'Share',
      visitCount: 'Visits',
      actions: 'Actions',
    },
    actions: {
      delete: 'Delete',
    },
    format: {
      unknown: 'Unknown',
      empty: '-',
    },
    messages: {
      fetchFailed: 'Failed to fetch visitors',
      deleteSuccess: 'Visitor deleted',
      deleteFailed: 'Delete failed',
    },
    loading: {
      text: 'Loading visitors...',
    },
    empty: {
      title: 'No visitors',
      description: 'No visitor records yet. Visits will appear here.',
    },
    deleteDialog: {
      title: 'Confirm Delete',
      message: 'Delete this visitor record?',
      visitorName: 'Name:',
      email: 'Email:',
      visitTime: 'Visited at:',
      cancel: 'Cancel',
      confirm: 'Delete',
    },
  },
  statistics: {
    cards: {
      totalShares: 'Total Shares',
      activeShares: 'Active Shares',
      viewsToday: 'Views Today',
    },
    charts: {
      trend: {
        title: 'Traffic Trend',
        viewsBar: 'Views',
        sharesBar: 'New Shares',
        viewsTooltip: '{count} views',
        sharesTooltip: '{count} new shares',
      },
      popular: {
        title: 'Popular Shares',
        refreshButton: 'Refresh',
      },
    },
    periodOptions: {
      last7Days: 'Last 7 days',
      last30Days: 'Last 30 days',
      last90Days: 'Last 90 days',
    },
    legend: {
      views: 'Views',
      shares: 'New Shares',
    },
    messages: {
      fetchFailed: 'Failed to fetch statistics',
    },
    loading: {
      text: 'Loading...',
    },
    empty: {
      noData: 'No data',
      noPopular: 'No popular shares yet',
    },
  },
}
