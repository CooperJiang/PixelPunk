/**
 * 图库浏览文案
 */
export const gallery = {
  pages: {
    gallery: '图库',
    upload: '上传文件',
    login: '登录',
    register: '注册',
    settings: '设置',
    profile: '个人资料',
    folders: '文件夹',
    dashboard: '仪表板',
    admin: '管理后台',
  },
  viewer: {
    previous: '上一张',
    next: '下一张',
    fitMode: '适应模式',
    fillMode: '填充模式',
    fullscreen: '全屏',
    exitFullscreen: '退出全屏',
    originalSize: '原始尺寸',
    imageInfo: '文件信息',
    similarImages: '相似文件',
  },
  waterfallLayout: {
    loadFailed: '加载失败',
    collapse: '收起',
  },
  file: {
    loadFailed: '加载失败',
    nsfwWarning: '违规内容',
  },
  enhancedFilePreview: {
    defaultImageName: '文件',
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
    title: '文件详情',
    untitled: '未命名',

    basicInfo: {
      title: '基本信息',
      fileId: '文件ID',
      originalName: '原始名称',
      fileFormat: '文件格式',
      fileSize: '文件大小',
      imageDimensions: '文件尺寸',
      visibility: '可见性',
      uploader: '上传者',
      uploadTime: '上传时间',
      updateTime: '更新时间',
      fileHash: '文件哈希',
      viewCount: '查看次数',
      status: '状态',
      times: '次',
      count: '个',
      unknownResolution: '未知分辨率',
      none: '无',
    },

    statusLabels: {
      duplicate: '重复',
      recommended: '推荐',
      nsfw: 'NSFW',
    },

    visibility: {
      public: '公开',
      private: '私有',
      protected: '受保护',
      link: '链接可见',
    },

    qualityLevels: {
      high: '高',
      medium: '中',
      low: '低',
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
      high: '高',
      medium: '中',
      low: '低',
    },

    visualAnalysis: {
      title: '视觉分析',
      dominantColor: '主色调',
      imageDimensions: '文件尺寸',
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
      deleteImage: '删除文件',
      cancelRecommend: '取消推荐',
      setRecommend: '设为推荐',
      close: '关闭',
    },
  },
  fileExpiryTag: {
    statuses: {
      expired: '已过期',
      expiringSoon: '即将过期',
      active: '限时',
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
    loading: '加载中',
  },
  fileActions: {
    copyLink: {
      protectedNotSupported: '受保护文件不支持外部分享，请在应用内查看',
      publicLinkCopied: '公开链接已复制到剪贴板',
      privateLinkCopied: '私有链接已复制到剪贴板',
      copyFailed: '复制链接失败',
    },
    delete: {
      success: '文件删除成功',
      failed: '删除失败，请重试',
    },
    accessLevel: {
      public: '公开',
      private: '私有',
      protected: '受保护',
      switched: '文件已设为{level}',
      switchFailed: '切换失败，请重试',
    },
    batch: {
      noPublicFiles: '所选文件中没有可复制的公开文件',
      filteredFiles: '已过滤掉 {count} 张私有/受保护文件',
      linksCopied: '已复制 {count} 个文件链接',
      copyFailed: '批量复制失败',
      deleteSuccess: '成功删除 {count} 个资源',
      deletePartial: '成功删除 {successCount}/{totalCount} 个资源',
      deleteFailed: '删除失败，请重试',
    },
  },
}
