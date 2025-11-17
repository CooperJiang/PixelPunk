/**
 * 图库浏览文案
 */
export const gallery = {
  pages: {
    gallery: '数据库',
    upload: '数据上传中心',
    login: '身份验证',
    register: '系统注册',
    settings: '系统配置',
    profile: '用户档案',
    folders: '数据容器',
    dashboard: '控制面板',
    admin: '管理终端',
  },
  viewer: {
    previous: '前一帧',
    next: '后一帧',
    fitMode: '适配视窗',
    fillMode: '填充视窗',
    fullscreen: '沉浸模式',
    exitFullscreen: '退出沉浸',
    originalSize: '原始分辨率',
    imageInfo: '数据属性',
    similarImages: '相似数据',
  },
  waterfallLayout: {
    loadFailed: '数据损坏',
    collapse: '收回',
  },
  file: {
    loadFailed: '数据加载失败',
    nsfwWarning: '受限内容',
  },
  enhancedFilePreview: {
    defaultImageName: '数据文件',
    fillMode: '填充模式',
    fitMode: '适应模式',
    switchToFit: '切换适应',
    switchToFill: '切换填充',
    wheelZoom: '滚轮缩放',
    dragMove: '拖拽移动',
    spaceKey: '空格',
    exitPreview: '退出预览',
    exitFullscreen: '退出全屏',
    exitPreviewEsc: '退出预览 (ESC)',
    exitFullscreenEsc: '退出全屏 (ESC)',
    enterFullscreen: '进入浏览器全屏',
    fullscreen: '全屏',
    preview: '预览',
    mode: '模式',
  },
  fileDetailModal: {
    title: '数据详情',
    untitled: '未命名数据',

    basicInfo: {
      title: '基本信息',
      fileId: '数据ID',
      originalName: '原始名称',
      fileFormat: '数据格式',
      fileSize: '数据大小',
      imageDimensions: '数据尺寸',
      visibility: '可见性',
      uploader: '上传者',
      uploadTime: '上传时间',
      updateTime: '更新时间',
      fileHash: '数据哈希',
      viewCount: '查看次数',
      status: '状态',
      times: '次',
      count: '个',
      unknownResolution: '未知分辨率',
      none: '无',
    },

    statusLabels: {
      duplicate: '重复数据',
      recommended: '推荐',
      nsfw: 'NSFW',
    },

    visibility: {
      public: '公开协议',
      private: '私密档案',
      protected: '安全防护',
      link: '链接可见',
    },

    qualityLevels: {
      high: '高级',
      medium: '中级',
      low: '低级',
      excellent: '优秀',
      normal: '正常',
      poor: '较差',
    },

    contentSafety: {
      title: '内容安全',
      nsfwStatus: 'NSFW状态',
      inappropriate: '不适宜',
      safe: '安全',
      nsfwScore: 'NSFW分数',
      high: '高风险',
      medium: '中等风险',
      low: '低风险',
    },

    visualAnalysis: {
      title: '视觉分析',
      dominantColor: '主色调',
      imageDimensions: '数据尺寸',
      aspectRatio: '宽高比',
      resolution: '分辨率',
      colorPalette: '色彩板',
      composition: '构图',
      objectsCount: '物体数量',
    },

    ai: {
      tags: 'AI 标签',
      description: 'AI 描述',
    },

    actions: {
      deleteImage: '清除数据',
      cancelRecommend: '取消推荐',
      setRecommend: '设为推荐',
      close: '断开连接',
    },
  },
  fileExpiryTag: {
    statuses: {
      expired: '数据过期',
      expiringSoon: '即将过期',
      active: '限时存储',
    },
    labels: {
      storageDuration: '存储时长',
      expiresAt: '过期时间',
      expiredAtPrefix: '已于',
      expiredAtSuffix: '过期',
    },
    units: {
      day: '天',
      hour: '小时',
      minute: '分钟',
    },
  },
  fileLoading: {
    loading: '数据传输中',
  },
  fileActions: {
    copyLink: {
      protectedNotSupported: '受保护数据不支持外部传输，请在终端内查看',
      publicLinkCopied: '公开链接已复制到剪贴板',
      privateLinkCopied: '私有链接已复制到剪贴板',
      copyFailed: '链接复制失败',
    },
    delete: {
      success: '数据删除成功',
      failed: '删除失败，请重试',
    },
    accessLevel: {
      public: '公开协议',
      private: '私密档案',
      protected: '安全防护',
      switched: '数据已设为{level}',
      switchFailed: '切换失败，请重试',
    },
    batch: {
      noPublicFiles: '所选数据中没有可复制的公开资源',
      filteredFiles: '已过滤掉 {count} 个私有/受保护数据',
      linksCopied: '已复制 {count} 个数据链接',
      copyFailed: '批量复制失败',
      deleteSuccess: '成功删除 {count} 个资源',
      deletePartial: '成功删除 {successCount}/{totalCount} 个资源',
      deleteFailed: '删除失败，请重试',
    },
  },
}
