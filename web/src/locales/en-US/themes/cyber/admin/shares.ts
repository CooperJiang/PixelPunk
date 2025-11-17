/**
 * Shares Management - Cyberpunk Style
 */
export const shares = {
  title: 'Share Network Management',
  subtitle: 'Data sharing protocol management and node access monitoring',
  tabs: {
    list: 'Protocol List',
    visitors: 'Node Access Records',
    stats: 'Network Data Statistics',
  },
  shareList: {
    buttons: {
      refresh: 'Resync',
      filter: 'Data filter',
    },
    columns: {
      name: 'Protocol name',
      shareKey: 'Access key',
      expiredAt: 'Protocol expiry',
      views: 'Connection count',
      createdAt: 'Protocol established',
      status: 'Protocol status',
      actions: 'Command operations',
    },
    actions: {
      copyLink: 'Copy access link',
      viewContent: 'View protocol content',
      disable: 'Disable protocol',
      enable: 'Enable protocol',
      delete: 'Terminate protocol',
      viewIpInfo: 'View node info',
    },
    status: {
      normal: 'Online',
      expired: 'Expired',
      deleted: 'Terminated',
      disabled: 'Disabled',
      unknown: 'Unknown status',
    },
    format: {
      unlimited: 'Permanent valid',
      infiniteViews: '∞',
      viewsFormat: '{current} / {max}',
    },
    messages: {
      fetchFailed: 'Protocol list sync failed',
      updateStatusSuccess: 'Protocol status updated successfully',
      updateStatusFailed: 'Protocol status update failed',
      deleteSuccess: 'Protocol terminated successfully',
      deleteFailed: 'Protocol termination failed',
      copySuccess: 'Access link copied to data buffer',
      copySuccessNormal: 'Access link copied to data buffer (key verification required)',
      copyFailed: 'Data copy failed, please execute manually',
      getAccessFailed: 'Admin access token acquisition failed, key verification required',
    },
    loading: {
      text: 'Protocol data syncing...',
    },
    empty: {
      title: 'Protocol list empty',
      description: 'No sharing protocols currently, node-established protocols will display here',
      resetButton: 'Reset filter commands',
    },
    statusDialog: {
      titleEnable: 'Enable protocol',
      titleUpdate: 'Update protocol status',
      message: 'Are you sure you want to change this protocol status to',
      messageHighlight: '{status}',
      reasonLabel: 'Status change reason (optional)',
      reasonPlaceholder: 'Enter status change reason...',
      cancel: 'Abort operation',
      confirm: 'Confirm execute',
    },
    deleteDialog: {
      title: 'Confirm termination',
      message: 'Are you sure you want to terminate sharing protocol',
      messageHighlight: '{name}',
      forceDelete: 'Force purge (unrecoverable)',
      cancel: 'Abort operation',
      confirm: 'Confirm termination',
    },
  },
  visitorList: {
    ipInfoTitle: 'View IP info',
    filter: {
      shareIdPlaceholder: 'Protocol ID filter',
      keywordPlaceholder: 'Search node name, communication email or node address...',
      searchButton: 'Data search',
    },
    columns: {
      id: 'Node ID',
      visitorName: 'Node name',
      visitorEmail: 'Communication email',
      createdAt: 'First connection',
      ipAddress: 'Node address',
      shareKey: 'Belong to protocol',
      visitCount: 'Connection count',
      actions: 'Command operations',
    },
    actions: {
      delete: 'Data purge',
    },
    format: {
      unknown: 'Unknown node',
      empty: '-',
    },
    messages: {
      fetchFailed: 'Node access record sync failed',
      deleteSuccess: 'Access record purged successfully',
      deleteFailed: 'Access record purge failed',
    },
    loading: {
      text: 'Node access data syncing...',
    },
    empty: {
      title: 'Access records empty',
      description: 'No node access records currently, node protocol connections will display here',
    },
    deleteDialog: {
      title: 'Confirm purge',
      message: 'Are you sure you want to purge this node access record?',
      visitorName: 'Node name:',
      email: 'Communication email:',
      visitTime: 'Connection time:',
      cancel: 'Abort operation',
      confirm: 'Confirm purge',
    },
  },
  statistics: {
    cards: {
      totalShares: 'Total protocols',
      activeShares: 'Online protocols',
      viewsToday: 'Connections today',
    },
    charts: {
      trend: {
        title: 'Connection trend analysis',
        viewsBar: 'Connection count',
        sharesBar: 'New protocols',
        viewsTooltip: '{count} connections',
        sharesTooltip: '{count} new protocols',
      },
      popular: {
        title: 'Popular protocol ranking',
        refreshButton: 'Resync',
      },
    },
    periodOptions: {
      last7Days: 'Last 7 days cycle',
      last30Days: 'Last 30 days cycle',
      last90Days: 'Last 90 days cycle',
    },
    legend: {
      views: 'Connection count',
      shares: 'New protocols',
    },
    messages: {
      fetchFailed: 'Statistics data sync failed',
    },
    loading: {
      text: 'Data syncing...',
    },
    empty: {
      noData: 'Data gap',
      noPopular: 'Popular protocol data gap',
    },
  },
}
