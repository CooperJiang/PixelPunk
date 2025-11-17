/**
 * Constants 常量文案 - 赛博风格
 */
export const constants = {
  fileTypes: {
    image: '数据文件',
    video: '视频流',
    document: '文档数据',
    archive: '压缩包',
    audio: '音频流',
    other: '未知类型',
  },

  accessLevels: {
    private: '私密模式',
    public: '公开访问',
    protected: '受保护',
  },

  sortOptions: {
    created_at: '创建时间',
    updated_at: '更新时间',
    size: '数据大小',
    name: '文件名称',
    views: '访问次数',
    downloads: '下载次数',
  },

  categorySortBy: {
    sort_order: '按优先级',
    name: '按名称',
    usage_count: '按使用频率',
    created_at: '按创建时间',
  },

  sortOrder: {
    asc: '升序排列',
    desc: '降序排列',
  },

  categoryType: {
    all: '全部分类',
    hot: '热门标签',
    normal: '普通标签',
  },

  batchOperations: {
    delete: '数据清除',
    move: '数据转移',
    update_access: '权限变更',
    update_folder: '移动至文件夹',
    download: '数据提取',
  },

  accessControl: {
    modes: {
      ipWhitelist: 'IP白名单',
      ipBlacklist: 'IP黑名单',
      domainWhitelist: '域名白名单',
      domainBlacklist: '域名黑名单',
    },
    restrictionModes: {
      strict: '严格模式',
      moderate: '中等模式',
      loose: '宽松模式',
    },
    blockActions: {
      block: '完全阻止',
      redirect: '重定向',
      thumbnail: '返回缩略图',
      watermark: '添加水印',
    },
    warnings: {
      block: '将向访问者返回403错误，完全拒绝访问数据资源，可能影响用户体验。',
      redirect: '将访问者重定向到您指定的URL，请确保目标URL可正常访问。',
      thumbnail: '仅向未授权访问者展示文件的缩略图版本，保护原始数据不被盗用。',
      watermark: '在文件上自动添加水印，可标识您的所有权但会改变原始外观。',
    },
    descriptions: {
      ipWhitelist: '白名单模式：仅允许列表中的IP访问',
      ipBlacklist: '黑名单模式：阻止列表中的IP访问',
      domainWhitelist: '白名单模式：仅允许列表中的域名引用',
      domainBlacklist: '黑名单模式：阻止列表中的域名引用',
      strictMode: '严格模式：必须完全匹配规则',
      moderateMode: '中等模式：部分匹配即可',
      looseMode: '宽松模式：最少限制',
    },
  },

  storageOptions: {
    unlimited: '无限制',
    storageTypes: {
      local: '本地存储',
      s3: 'S3云端',
      aliyun: '阿里云OSS',
      tencent: '腾讯云COS',
    },
    channelTypes: {
      oss: '阿里云 OSS',
      cos: '腾讯云COS',
      rainyun: '雨云 RainYun',
      local: '本地存储',
      unknown: '未知节点',
    },
    storageClasses: {
      standard: '标准存储',
      standardMultiAz: '标准存储(多AZ)',
      infrequent: '低频存储',
      infrequentMultiAz: '低频存储(多AZ)',
      infrequentAccess: '低频访问',
      intelligentTiering: '智能分层',
      intelligentTieringMultiAz: '智能分层(多AZ)',
      archive: '归档存储',
      deepArchive: '深度归档',
      coldArchive: '冷归档',
      glacier: 'Glacier存储',
    },
  },

  filterOptions: {
    sortOptions: {
      newest: '最新注入',
      oldest: '最早注入',
      name: '名称',
      size: '数据大小',
      width: '宽度',
      height: '高度',
      quality: '数据质量',
      nsfw_score: '敏感度评分',
    },
    resolutionOptions: {
      low: '低分辨率',
      '720p': '720p',
      '1080p': '1080p',
      '2k': '2k',
      '4k': '4K',
      '8k': '8K',
    },
    imageFormatOptions: {
      jpg: 'JPG',
      jpeg: 'JPEG',
      png: 'PNG',
      gif: 'GIF',
      webp: 'WebP',
      bmp: 'BMP',
      svg: 'SVG',
      ico: 'ICO',
    },
    fileTypeOptions: {
      all: '全部数据类型',
      jpeg: 'JPEG数据',
      png: 'PNG数据',
      gif: 'GIF数据',
      webp: 'WebP数据',
      svg: 'SVG数据',
      bmp: 'BMP数据',
      ico: 'ICO图标',
    },
  },

  category: {
    sortByOptions: {
      sort_order: '按排序',
      name: '按名称',
      usage_count: '按使用次数',
      created_at: '按创建时间',
    },
    sortOrderOptions: {
      asc: '升序',
      desc: '降序',
    },
    typeOptions: {
      all: '全部',
      popular: '热门',
      normal: '普通',
    },
  },

  contentReview: {
    actions: {
      all: '全部行动',
      approve: '通过验证',
      reject: '数据拒绝',
    },
    sort: {
      newest: '最新注入',
      oldest: '最早注入',
      size: '数据大小',
      nsfw: '风险评分',
    },
    nsfwLevels: {
      safe: '安全级别',
      low: '低风险',
      medium: '中风险',
      high: '高风险',
    },
    columns: {
      file: '数据文件',
      fileInfo: '数据信息',
      action: '审查动作',
      softDelete: '标记状态',
      auditor: '审查员',
      auditTime: '审查时间',
      reason: '审查原因',
      nsfw: '风险评分',
      operations: '系统操作',
    },
    defaultReasons: {
      approve: '数据验证通过',
      reject: '不符合系统规范',
      batchApprove: '批量验证通过',
      batchReject: '批量数据拒绝',
    },
    messages: {
      approveSuccess: '数据验证通过',
      rejectSuccess: '数据已拒绝',
      batchApproveSuccess: '批量验证成功，已处理 {count} 个数据',
      batchRejectSuccess: '批量拒绝成功，已处理 {count} 个数据',
      restoreSuccess: '数据已恢复',
      hardDeleteSuccess: '数据已彻底清除',
      batchRestoreSuccess: '批量恢复完成：成功 {success} 个',
      batchRestoreWithFail: '批量恢复完成：成功 {success} 个，失败 {fail} 个',
      batchHardDeleteSuccess: '批量清除完成：成功 {success} 个',
      batchHardDeleteWithFail: '批量清除完成：成功 {success} 个，失败 {fail} 个',
      loadError: '数据加载失败',
      operationError: '操作执行失败',
      noSelection: '请先选择要操作的数据项',
      queueEmpty: '队列清空',
      queueEmptyDesc: '所有数据均已完成审查',
      logsEmpty: '暂无审查日志',
      logsEmptyDesc: '系统中还没有任何审查记录',
    },
    warnings: {
      hardDelete: '警告：此操作不可逆',
      hardDeleteDesc: '数据将从服务器彻底清除，无法恢复。',
      batchOperation: '此操作将影响多个数据项',
    },
    descriptions: {
      approveEffects: ['数据将在公共数据库中显示', '数据状态更改为"已验证"', '用户可以正常分享和传播数据'],
      restoreEffects: ['数据将从标记删除状态恢复为正常状态'],
    },
  },

  share: {
    sortOptions: {
      date: '日期',
      name: '名称',
      size: '大小',
    },
    layoutOptions: {
      grid: {
        label: '网格布局',
        description: '整齐的网格展示',
      },
      waterfall: {
        label: '瀑布流',
        description: '根据数据高度自动排列',
      },
      masonry: {
        label: '砖墙布局',
        description: '左对齐的瀑布流',
      },
      large: {
        label: '大图模式',
        description: '查看更大的缩略图',
      },
    },
  },

  navigation: {
    admin: {
      dashboard: '控制面板',
      files: '数据管理',
      tags: '标签系统',
      categories: '分类系统',
      shares: '分享管理',
      ai: '智能系统',
      contentReview: '内容审核',
      tagging: '自动标注',
      vectors: '向量系统',
      users: '用户管理',
      channels: '渠道管理',
      construction: '站点建设',
      announcements: '公告系统',
      settings: '系统设置',
    },
    docs: {
      overview: 'API概述',
      authentication: '认证机制',
      upload: '上传API',
      limits: '限制说明',
      examples: '代码示例',
      tester: 'API测试',
      faq: '常见问题',
    },
    main: {
      home: '主页',
      explore: '数据库',
      upload: '上传',
      random: '随机',
      docs: 'API',
    },
    user: {
      dashboard: '控制面板',
      myFiles: '我的数据',
      folders: '文件夹',
      tagManage: '标签管理',
      categoryManage: '分类管理',
      shares: '我的分享',
      settings: '设置',
    },
    settings: {
      profile: '个人信息',
      security: '安全设置',
      apikey: 'API密钥',
      storage: '存储配置',
    },
    adminSettings: {
      website: '站点设置',
      upload: '上传设置',
      security: '安全设置',
      registration: '注册设置',
      mail: '邮件设置',
      ai: 'AI设置',
      vector: '向量设置',
    },
  },

  api: {
    errors: {
      badRequest: '请求参数异常',
      unauthorized: '未授权访问',
      forbidden: '访问被拒绝',
      notFound: '资源未找到',
      methodNotAllowed: '方法不允许',
      conflict: '数据冲突',
      payloadTooLarge: '数据包过大',
      unprocessableEntity: '数据验证失败',
      tooManyRequests: '请求频率过高',
      unavailableForLegalReasons: 'IP地址已封禁',
      internalServerError: '系统内部错误',
      badGateway: '网关异常',
      serviceUnavailable: '服务不可用',
      gatewayTimeout: '网关超时',
    },
  },

  systemOptions: {
    sensitiveContent: {
      markOnly: '仅标记',
      pendingReview: '系统审核',
      autoDelete: '自动清除',
    },
    sizeUnits: {
      bytes: 'Bytes (B)',
      kilobytes: 'Kilobytes (KB)',
      megabytes: 'Megabytes (MB)',
      gigabytes: 'Gigabytes (GB)',
      terabytes: 'Terabytes (TB)',
    },
    userRoles: {
      normal: '普通用户',
      advanced: '高级用户',
      admin: '系统管理员',
      superAdmin: '超级管理员',
    },
    userStatus: {
      normal: '正常',
      disabled: '已禁用',
      pending: '待审核',
    },
    imageStatus: {
      normal: '正常',
      hidden: '已隐藏',
      pending: '待审核',
      deleted: '已删除',
    },
    shareStatus: {
      normal: '正常',
      disabled: '已禁用',
      expired: '已过期',
    },
  },
}
