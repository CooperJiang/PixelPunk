/**
 * 仪表盘页面文案 - 赛博风格
 */
export const dashboard = {
  quickActions: {
    title: '快速操作',
    dragHint: '拖拽排序',
    actions: {
      myFiles: '数据单元',
      folders: '文件夹',
      tagManage: '标签管理',
      categoryManage: '分类管理',
      openApi: '开放API',
      automation: '自动任务',
    },
    settings: {
      divider: '设置中心',
      api: 'API',
      profile: '个人信息',
      security: '安全',
      accessControl: '防盗链',
      preferences: '偏好设置',
    },
  },
  stats: {
    title: '数据概览',
    items: {
      images: '资源数量',
      storage: '已用存储',
      views: '访问次数',
      shares: '分享数量',
    },
    quotas: {
      title: '额度概览',
      storage: {
        label: '空间额度',
        used: '已用',
        total: '总计',
      },
      bandwidth: {
        label: '带宽额度',
        used: '已用',
        total: '总计',
      },
      status: {
        healthy: '充裕',
        warning: '适中',
        danger: '紧张',
      },
      description: {
        healthy: '资源使用充裕，可以安心使用',
        warning: '资源使用适中，建议定期清理不需要的内容',
        danger: '资源使用紧张，请及时清理内容或联系管理员扩容',
      },
    },
  },
  messages: {
    title: '消息通知',
    unread: '未读',
    viewAll: '查看全部',
    loading: '加载消息中...',
    empty: '暂无消息',
    loadingMore: '正在加载更多...',
    loadingMoreDesc: '请稍候',
    noMore: '已加载全部消息',
    noMoreDesc: '共 {count} 条',
    scrollHint: '滑动查看更多',
    markReadFailed: '标记已读失败，请重新尝试',
  },
  errors: {
    fetchStatsFailed: '获取数据失败',
    networkFailed: '网络请求失败',
  },
  activityMonitor: {
    title: '活动监控',
    todayUploads: '今日上传',
    totalViews: '总浏览',
    recentActivity: '最近活动',
    recordCount: '共 {count} 条记录',
    loadingActivities: '正在加载活动日志...',
    loadingMore: '正在加载更多...',
    loadingHint: '请稍候',
    allLoaded: '已加载全部活动记录',
    empty: '暂无活动记录',
    emptyHint: '开始使用后这里会显示您的活动',
    status: {
      active: '正常监控',
      warning: '警告',
      error: '错误',
      offline: '离线',
    },
    toast: {
      fetchFailed: '获取活动日志失败',
      noData: '暂时无法获取实时活动数据',
    },
    timeAgo: {
      justNow: '刚刚',
      hoursAgo: '{hours}小时前',
      minutesAgo: '{minutes}分钟前',
      daysAgo: '{days}天前',
    },
  },
  recentUploads: {
    title: '最近上传',
    loading: '加载中...',
    empty: '暂无上传记录',
    toast: {
      noImages: '暂无图片可复制',
      copySuccess: '已复制 {count} 个数据单元链接',
      copyFailed: '复制失败',
      fetchFailed: '获取数据失败',
    },
    actions: {
      copyAll: '复制全部链接',
      viewAll: '查看全部',
    },
  },
  uploadQueue: {
    title: '上传队列',
    statsTitle: '上传统计',
    totalFiles: '总计 {count} 个文件',
    overallProgress: '总体进度',
    queueTotal: '队列总数',
    uploading: '上传中',
    pending: '待上传',
    completed: '已上传',
    failed: '失败',
    actions: {
      goToUpload: '前往上传页面',
    },
  },
  uploadEntry: {
    title: {
      default: '上传文件',
      withCount: '队列中有 {count} 个文件',
    },
  },
  folderDistribution: {
    title: '文件夹分布',
    description: '块大小表示存储空间,颜色深浅表示占比',
    totalFolders: '总文件夹数: {count}',
    empty: '暂无文件夹',
    tooltipSize: '大小: {size}',
    tooltipFiles: '文件: {count}张',
    tooltipPercent: '占比: {percent}%',
    unnamed: '未命名',
  },
}
