/**
 * Shares 分享管理文案
 */
export const shares = {
  title: '分享管理',
  subtitle: '用户分享链接管理和访客信息统计',
  tabs: {
    list: '分享列表',
    visitors: '访客信息',
    stats: '数据统计',
  },
  shareList: {
    buttons: {
      refresh: '刷新数据',
      filter: '筛选',
    },
    columns: {
      name: '分享名称',
      shareKey: '分享Key',
      expiredAt: '过期时间',
      views: '访问次数',
      createdAt: '创建时间',
      status: '状态',
      actions: '操作',
    },
    actions: {
      copyLink: '复制分享链接',
      viewContent: '查看分享内容',
      disable: '禁用分享',
      enable: '启用分享',
      delete: '删除分享',
      viewIpInfo: '查看IP信息',
    },
    status: {
      normal: '正常',
      expired: '已过期',
      deleted: '已删除',
      disabled: '已禁用',
      unknown: '未知',
    },
    format: {
      unlimited: '无限期',
      infiniteViews: '∞',
      viewsFormat: '{current} / {max}',
    },
    messages: {
      fetchFailed: '获取分享列表失败',
      updateStatusSuccess: '更新分享状态成功',
      updateStatusFailed: '更新分享状态失败',
      deleteSuccess: '删除分享成功',
      deleteFailed: '删除分享失败',
      copySuccess: '分享链接已复制到剪贴板',
      copySuccessNormal: '分享链接已复制到剪贴板（普通链接，可能需要密码）',
      copyFailed: '复制失败，请手动复制',
      getAccessFailed: '获取管理员访问权限失败，可能需要输入密码',
    },
    loading: {
      text: '正在加载分享数据...',
    },
    empty: {
      title: '暂无分享数据',
      description: '当前没有任何分享链接，用户创建分享后将显示在这里',
      resetButton: '重置筛选条件',
    },
    statusDialog: {
      titleEnable: '启用分享',
      titleUpdate: '更新分享状态',
      message: '您确定要将此分享的状态更改为',
      messageHighlight: '{status}',
      reasonLabel: '状态变更原因（可选）',
      reasonPlaceholder: '请输入状态变更原因...',
      cancel: '取消',
      confirm: '确定',
    },
    deleteDialog: {
      title: '确认删除',
      message: '您确定要删除分享',
      messageHighlight: '{name}',
      forceDelete: '彻底删除（不可恢复）',
      cancel: '取消',
      confirm: '删除',
    },
  },
  visitorList: {
    ipInfoTitle: '查看IP信息',
    filter: {
      shareIdPlaceholder: '分享ID过滤',
      keywordPlaceholder: '搜索访客姓名、邮箱或IP...',
      searchButton: '搜索',
    },
    columns: {
      id: 'ID',
      visitorName: '访客姓名',
      visitorEmail: '邮箱',
      createdAt: '访问时间',
      ipAddress: 'IP地址',
      shareKey: '所属分享',
      visitCount: '访问次数',
      actions: '操作',
    },
    actions: {
      delete: '删除',
    },
    format: {
      unknown: '未知',
      empty: '-',
    },
    messages: {
      fetchFailed: '获取访客信息列表失败',
      deleteSuccess: '删除访客信息成功',
      deleteFailed: '删除访客信息失败',
    },
    loading: {
      text: '正在加载访客信息...',
    },
    empty: {
      title: '暂无访客信息',
      description: '当前没有任何访客记录，用户访问分享链接后将显示在这里',
    },
    deleteDialog: {
      title: '确认删除',
      message: '您确定要删除此访客信息记录吗？',
      visitorName: '访客姓名：',
      email: '邮箱：',
      visitTime: '访问时间：',
      cancel: '取消',
      confirm: '删除',
    },
  },
  statistics: {
    cards: {
      totalShares: '总分享数',
      activeShares: '活跃分享',
      viewsToday: '今日访问',
    },
    charts: {
      trend: {
        title: '访问趋势',
        viewsBar: '访问次数',
        sharesBar: '新分享数',
        viewsTooltip: '{count} 次访问',
        sharesTooltip: '{count} 个新分享',
      },
      popular: {
        title: '热门分享',
        refreshButton: '刷新',
      },
    },
    periodOptions: {
      last7Days: '最近7天',
      last30Days: '最近30天',
      last90Days: '最近90天',
    },
    legend: {
      views: '访问次数',
      shares: '新分享数',
    },
    messages: {
      fetchFailed: '获取统计数据失败',
    },
    loading: {
      text: '加载中...',
    },
    empty: {
      noData: '暂无数据',
      noPopular: '暂无热门分享数据',
    },
  },
}
