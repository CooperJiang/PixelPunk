/**
 * Shares 分享管理文案 - 赛博朋克风格
 */
export const shares = {
  title: '分享网络管理',
  subtitle: '数据分享协议管理和节点访问监控',
  tabs: {
    list: '协议列表',
    visitors: '节点访问记录',
    stats: '网络数据统计',
  },
  shareList: {
    buttons: {
      refresh: '重新同步',
      filter: '数据筛选',
    },
    columns: {
      name: '协议名称',
      shareKey: '访问密钥',
      expiredAt: '协议到期',
      views: '连接计数',
      createdAt: '协议建立',
      status: '协议状态',
      actions: '指令操作',
    },
    actions: {
      copyLink: '复制访问链接',
      viewContent: '查看协议内容',
      disable: '禁用协议',
      enable: '启用协议',
      delete: '终止协议',
      viewIpInfo: '查看节点信息',
    },
    status: {
      normal: '在线',
      expired: '已过期',
      deleted: '已终止',
      disabled: '已禁用',
      unknown: '未知状态',
    },
    format: {
      unlimited: '永久有效',
      infiniteViews: '∞',
      viewsFormat: '{current} / {max}',
    },
    messages: {
      fetchFailed: '协议列表同步失败',
      updateStatusSuccess: '协议状态更新成功',
      updateStatusFailed: '协议状态更新失败',
      deleteSuccess: '协议终止成功',
      deleteFailed: '协议终止失败',
      copySuccess: '访问链接已复制至数据缓冲区',
      copySuccessNormal: '访问链接已复制至数据缓冲区（需要验证密钥）',
      copyFailed: '数据复制失败，请手动执行',
      getAccessFailed: '管理员访问令牌获取失败，需要验证密钥',
    },
    loading: {
      text: '协议数据同步中...',
    },
    empty: {
      title: '协议列表为空',
      description: '当前没有任何分享协议，节点建立协议后将显示在这里',
      resetButton: '重置筛选指令',
    },
    statusDialog: {
      titleEnable: '启用协议',
      titleUpdate: '更新协议状态',
      message: '您确定要将此协议的状态更改为',
      messageHighlight: '{status}',
      reasonLabel: '状态变更原因（可选）',
      reasonPlaceholder: '请输入状态变更原因...',
      cancel: '中止操作',
      confirm: '确认执行',
    },
    deleteDialog: {
      title: '确认终止',
      message: '您确定要终止分享协议',
      messageHighlight: '{name}',
      forceDelete: '强制清除（不可恢复）',
      cancel: '中止操作',
      confirm: '确认终止',
    },
  },
  visitorList: {
    ipInfoTitle: '查看IP信息',
    filter: {
      shareIdPlaceholder: '协议ID筛选',
      keywordPlaceholder: '检索节点名称、通信邮箱或节点地址...',
      searchButton: '数据检索',
    },
    columns: {
      id: '节点ID',
      visitorName: '节点名称',
      visitorEmail: '通信邮箱',
      createdAt: '首次连接',
      ipAddress: '节点地址',
      shareKey: '所属协议',
      visitCount: '连接次数',
      actions: '指令操作',
    },
    actions: {
      delete: '数据清除',
    },
    format: {
      unknown: '未知节点',
      empty: '-',
    },
    messages: {
      fetchFailed: '节点访问记录同步失败',
      deleteSuccess: '访问记录清除成功',
      deleteFailed: '访问记录清除失败',
    },
    loading: {
      text: '节点访问数据同步中...',
    },
    empty: {
      title: '访问记录为空',
      description: '当前没有任何节点访问记录，节点建立协议连接后将显示在这里',
    },
    deleteDialog: {
      title: '确认清除',
      message: '您确定要清除此节点访问记录吗？',
      visitorName: '节点名称：',
      email: '通信邮箱：',
      visitTime: '连接时间：',
      cancel: '中止操作',
      confirm: '确认清除',
    },
  },
  statistics: {
    cards: {
      totalShares: '协议总量',
      activeShares: '在线协议',
      viewsToday: '今日连接',
    },
    charts: {
      trend: {
        title: '连接趋势分析',
        viewsBar: '连接次数',
        sharesBar: '新建协议',
        viewsTooltip: '{count} 次连接',
        sharesTooltip: '{count} 个新协议',
      },
      popular: {
        title: '热门协议排行',
        refreshButton: '重新同步',
      },
    },
    periodOptions: {
      last7Days: '近7日周期',
      last30Days: '近30日周期',
      last90Days: '近90日周期',
    },
    legend: {
      views: '连接次数',
      shares: '新建协议',
    },
    messages: {
      fetchFailed: '统计数据同步失败',
    },
    loading: {
      text: '数据同步中...',
    },
    empty: {
      noData: '数据空缺',
      noPopular: '热门协议数据空缺',
    },
  },
}
