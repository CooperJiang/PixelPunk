/**
 * Constants 常量文案
 */
export const constants = {
  fileTypes: {
    image: '文件',
    video: '视频',
    document: '文档',
    archive: '压缩包',
    audio: '音频',
    other: '其他',
  },

  accessLevels: {
    private: '私密',
    public: '公开',
    protected: '受保护',
  },

  sortOptions: {
    created_at: '创建时间',
    updated_at: '更新时间',
    size: '文件大小',
    name: '文件名称',
    views: '访问次数',
    downloads: '下载次数',
  },

  categorySortBy: {
    sort_order: '按排序',
    name: '按名称',
    usage_count: '按使用次数',
    created_at: '按创建时间',
  },

  sortOrder: {
    asc: '升序',
    desc: '降序',
  },

  categoryType: {
    all: '全部',
    hot: '热门',
    normal: '普通',
  },

  batchOperations: {
    delete: '删除',
    move: '移动',
    update_access: '修改权限',
    update_folder: '移动到文件夹',
    download: '下载',
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
      block: '将向访问者返回403错误，完全拒绝访问文件资源，可能影响用户体验。',
      redirect: '将访问者重定向到您指定的URL，请确保目标URL可正常访问。',
      thumbnail: '仅向未授权访问者展示文件的缩略图版本，保护原图不被盗用。',
      watermark: '在文件上自动添加水印，可标识您的所有权但会改变原图外观。',
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
    unlimited: '不限量',
    storageTypes: {
      local: '本地存储',
      s3: 'S3存储',
      aliyun: '阿里云OSS',
      tencent: '腾讯云COS',
    },
    channelTypes: {
      oss: '阿里云 OSS',
      cos: '腾讯云COS',
      rainyun: '雨云 RainYun',
      local: '本地存储',
      unknown: '未知',
    },
    storageClasses: {
      standard: '标准存储',
      standardMultiAz: '标准存储(多AZ)',
      infrequent: '低频存储',
      infrequentMultiAz: '低频存储(多AZ)',
      infrequentAccess: '低频访问',
      intelligentTiering: '智能分层存储',
      intelligentTieringMultiAz: '智能分层存储(多AZ)',
      archive: '归档存储',
      deepArchive: '深度归档存储',
      coldArchive: '冷归档存储',
      glacier: 'Glacier存储',
    },
  },

  filterOptions: {
    sortOptions: {
      newest: '最新上传',
      oldest: '最早上传',
      name: '名称',
      size: '文件大小',
      width: '宽度',
      height: '高度',
      quality: '文件质量',
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
      all: '所有文件类型',
      jpeg: 'JPEG文件',
      png: 'PNG文件',
      gif: 'GIF文件',
      webp: 'WebP文件',
      svg: 'SVG文件',
      bmp: 'BMP文件',
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
      all: '全部动作',
      approve: '批准',
      reject: '拒绝',
    },
    sort: {
      newest: '最新上传',
      oldest: '最旧上传',
      size: '文件大小',
      nsfw: 'NSFW评分',
    },
    nsfwLevels: {
      safe: '安全',
      low: '轻微',
      medium: '中等',
      high: '严重',
    },
    columns: {
      file: '文件',
      fileInfo: '文件信息',
      action: '审核动作',
      softDelete: '删除状态',
      auditor: '审核员',
      auditTime: '审核时间',
      reason: '审核原因',
      nsfw: 'NSFW评分',
      operations: '操作',
    },
    defaultReasons: {
      approve: '文件审核通过',
      reject: '不符合平台规范',
      batchApprove: '批量批准',
      batchReject: '批量拒绝',
    },
    messages: {
      approveSuccess: '文件审批通过',
      rejectSuccess: '文件已拒绝',
      batchApproveSuccess: '批量批准成功，共处理 {count} 个资源',
      batchRejectSuccess: '批量拒绝成功,共处理 {count} 个资源',
      restoreSuccess: '文件已恢复',
      hardDeleteSuccess: '文件已彻底删除',
      batchRestoreSuccess: '批量恢复完成：成功 {success} 个',
      batchRestoreWithFail: '批量恢复完成：成功 {success} 个，失败 {fail} 个',
      batchHardDeleteSuccess: '批量删除完成：成功 {success} 个',
      batchHardDeleteWithFail: '批量删除完成：成功 {success} 个，失败 {fail} 个',
      loadError: '加载数据失败',
      operationError: '操作失败',
      noSelection: '请先选择要操作的项目',
      queueEmpty: '队列为空',
      queueEmptyDesc: '所有文件都已审核完成',
      logsEmpty: '暂无审核记录',
      logsEmptyDesc: '还没有任何审核操作记录',
    },
    warnings: {
      hardDelete: '警告：此操作不可逆',
      hardDeleteDesc: '文件文件将从服务器彻底删除，无法恢复。',
      batchOperation: '此操作将影响多个项目',
    },
    descriptions: {
      approveEffects: ['文件将可以在公共画廊中显示', '文件状态改为"已审核"', '用户可以正常分享和传播'],
      restoreEffects: ['文件将从软删除状态恢复为正常状态'],
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
        description: '根据文件高度自动排列',
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
      dashboard: '仪表盘',
      files: '文件管理',
      tags: '标签管理',
      categories: '分类管理',
      shares: '分享管理',
      ai: '智能管理',
      contentReview: '内容审核',
      tagging: '打标管理',
      vectors: '向量管理',
      users: '用户管理',
      channels: '渠道管理',
      construction: '网站建设',
      announcements: '公告管理',
      settings: '全局设置',
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
      home: '首页',
      explore: '图库',
      upload: '上传',
      random: '随机',
      docs: 'API',
    },
    user: {
      dashboard: '仪表盘',
      myFiles: '我的文件',
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
      badRequest: '请求参数错误',
      unauthorized: '未授权访问',
      forbidden: '权限不足',
      notFound: '资源不存在',
      methodNotAllowed: '请求方法不允许',
      conflict: '资源冲突',
      payloadTooLarge: '请求体过大',
      unprocessableEntity: '数据验证失败',
      tooManyRequests: '请求过于频繁',
      unavailableForLegalReasons: 'IP地址被禁用',
      internalServerError: '服务器内部错误',
      badGateway: '网关错误',
      serviceUnavailable: '服务不可用',
      gatewayTimeout: '网关超时',
    },
  },

  systemOptions: {
    sensitiveContent: {
      markOnly: '仅标记',
      pendingReview: '管理员审核',
      autoDelete: '自动删除',
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
      admin: '管理员',
      superAdmin: '超级管理员',
    },
    userStatus: {
      normal: '正常',
      disabled: '禁用',
      pending: '待审核',
    },
    imageStatus: {
      normal: '正常',
      hidden: '隐藏',
      pending: '待审核',
      deleted: '已删除',
    },
    shareStatus: {
      normal: '正常',
      disabled: '禁用',
      expired: '过期',
    },
  },
}
