/**
 * 上传页面文案 - 赛博风格
 */
export const upload = {
  title: '数据注入',

  settings: {
    title: '注入设置',
    advanced: '高级配置',
    folder: '数据文件夹',
    folderHint: '支持多级目录，可搜索文件夹名称',
    guestFolderInfo: '游客注入的数据单元将存放在系统文件夹中，默认{level}访问',
    guestLimitInfo: '游客模式仅支持单数据单元注入（并发限制为1）',
    guestDisabled: '游客注入功能已禁用，请联系管理员',
    accessLevelLabel: '访问权限设置',
    accessLevelPlaceholder: '选择访问权限',
    accessLevel: {
      public: '🌍 公开',
      publicDesc: '可被推荐展示，任何人可通过链接访问',
      private: '🔗 私有',
      privateDesc: '不会被推荐，但可正常分享和外部引用',
      protected: '🔒 受保护',
      protectedDesc: '只有您登录后可见，外部链接无法访问',
    },
    storageDurationLabel: '存储时长',
    storageDurationPlaceholder: '请选择存储时长',
    storageDurationHint: '游客注入必须选择存储时长，数据单元到期后将自动删除',
    required: '(必选)',
    autoOptimize: '自动优化数据单元（推荐）',
    autoOptimizeHint: '自动压缩和优化缩略图，减小数据单元大小而不明显降低质量',
    autoRemove: '注入成功后自动移除',
    autoRemoveHint: '注入成功的数据单元将自动从列表中移除，适合批量注入时保持列表简洁',
    duration: {
      permanent: '永久',
    },
  },

  uploadProgress: {
    status: {
      uploadSuccess: '注入成功',
      instantSuccess: '瞬时注入成功',
      uploadFailed: '注入失败',
      analyzing: '正在分析数据单元...',
      paused: '发现未完成的注入，可以继续',
      retrying: '正在重试注入...',
      pending: '等待注入',
    },
    time: {
      seconds: '{count}秒',
      minutes: '{count}分钟',
      hoursMinutes: '{hours}小时{minutes}分钟',
    },
    remaining: '剩余',
    actions: {
      openInNewTab: '新窗口打开',
      copyLink: '复制数据单元链接',
      resume: '继续注入',
      retry: '重新注入',
      remove: '移除数据单元',
    },
    toast: {
      linkCopied: '链接已复制',
      copyFailed: '复制失败',
      uploadSuccessWithRemainingCount: '数据注入成功！剩余 {count} 次注入机会',
      watermarkApplyFailed: '数据单元已注入，但水印应用失败: {reason}',
      thumbnailGenerationFailed: '数据单元已注入，但缩略图生成失败: {reason}',
      fileTypeNotSupported: '数据类型不支持: {name}',
      fileSizeExceeded: '数据大小超限: {name} ({size})',
      instantUploadSuccess: '秒传成功: {name}',
      cacheCleared: '缓存已清除',
    },
  },

  dropZone: {
    title: '拖拽数据单元到此处注入',
    orText: '或者',
    clickToSelect: '点击选择数据',
    supportFolder: '支持拖拽整个数据目录',
    features: {
      batch: '批量注入',
      autoOptimize: '自动优化',
      secureStorage: '安全存储',
    },
    shortcut: {
      label: '快捷键',
      paste: '粘贴注入',
    },
  },

  actions: {
    startUpload: '启动注入',
    selectFiles: '选择数据',
    selectFolder: '选择数据目录',
    copyLinks: '复制链接',
    clearQueue: '清空队列',
    cancelUpload: '取消注入',
    clearCache: '清除缓存',
  },

  queue: {
    title: '注入队列',
    filesCount: '{count} 个数据单元',
    concurrent: '{current}/{max} 并发',
    stats: {
      pending: '待注入',
      uploading: '注入中',
      success: '成功',
      failed: '失败',
      total: '总计',
      totalSize: '总大小',
    },
    empty: {
      title: '队列为空',
      desc: '拖拽数据、点击选择或粘贴截图即可开始注入',
      action: '选择数据',
    },
  },

  guide: {
    title: '注入提示',
    multipleWays: {
      title: '多种注入方式',
      desc: '支持拖拽、点击选择或使用 {key}+V 粘贴剪贴板数据',
    },
    smartChunking: {
      title: '智能分片',
      desc: '大数据单元自动分片传输，支持断点续传',
    },
    oneCopy: {
      title: '一键复制',
      desc: '批量注入后可一键复制所有数据链接',
    },
  },

  advanced: {
    title: '高级配置',
    storageDuration: '存储时长',
    storageDurationDesc: '数据单元在服务器上的保存时间',
    watermark: '添加数据水印',
    watermarkDesc: '为注入的数据添加水印标识',
    watermarkInfo: {
      type: '类型:',
      content: '内容:',
      position: '位置:',
    },
    watermarkType: {
      text: '文字',
      image: '图片',
    },
    loginPrompt: {
      title: '更多高级功能',
      desc: '接入系统后可使用永久存储、自定义水印等功能',
      action: '立即接入',
    },
  },

  folder: {
    browserNotSupported: '您的浏览器不支持数据目录注入功能',
    noValidFiles: '数据目录中没有可注入的数据单元',
    parsing: '正在解析数据目录...',
    addSuccess: '已从数据目录中添加 {count} 个数据单元',
    dialog: {
      title: '数据目录注入确认',
      totalFiles: '数据单元总数',
      validFiles: '可注入单元',
      invalidFiles: '不支持单元',
      confirmMessage: '检测到 <strong class="highlight-text">{count}</strong> 个可注入数据单元，是否加入注入队列？',
      warningMessage: '已过滤 {count} 个不支持的数据单元',
      cancel: '取消',
      confirm: '确认注入',
    },
  },

  toast: {
    folderSelected: '已选择数据目录: {path}',
    folderPreselected: '数据目录已预选，可直接开始数据注入',
    sensitiveAutoDelete: '违规数据将被自动清除，请确保数据单元符合协议',
    sensitivePendingReview: '敏感数据将进入审核队列，审核期间暂不可访问',
    defaultHint: '请确保注入数据符合平台协议',
    guestDurationRequired: '访客注入必须选择存储时长',
    uploadCancelled: '数据注入已取消',
    noFilesInQueue: '队列中暂无需要注入的数据单元',
    allLinksCopied: '所有链接已复制到剪贴板',
    noLinksToCopy: '暂无可复制的数据链接',
  },

  duration: {
    permanent: '永久',
  },

  watermark: {
    position: {
      'top-left': '左上角',
      'top-center': '上方居中',
      'top-right': '右上角',
      'middle-left': '左侧居中',
      'middle-center': '正中央',
      'middle-right': '右侧居中',
      'bottom-left': '左下角',
      'bottom-center': '下方居中',
      'bottom-right': '右下角',
      custom: '自定义位置',
    },
    configure: '配置数据水印',
  },

  form: {
    storageDuration: {
      placeholder: '选择存储时长',
    },
  },

  smartUpload: {
    error: {
      fileTypeNotSupported: '数据类型不支持: {name} ({type})',
      fileSizeExceeded: '数据单元超限: {name} ({size})',
      chunkedUploadNotReady: '分片传输未初始化，请稍后重试',
      uploadFailed: '传输失败',
      analysisFailed: '数据分析失败',
      analysisTimeout: '数据分析超时',
    },
    status: {
      preparing: '准备注入',
      uploading: '注入中...',
      uploadingWithProgress: '注入中 {progress}% - {speed} - 剩余{remaining}',
      completed: '注入完成',
      failed: '注入失败: {error}',
      paused: '已暂停',
      cancelled: '已取消',
      retrying: '准备重试',
      userPaused: '用户暂停',
      userCancelled: '用户取消',
      analyzing: '正在分析数据...',
      analyzingProgress: '数据分析中 {progress}%',
      checkingInstant: '检查秒传...',
      checkingInstantProgress: '检查秒传 {progress}%',
      instantCompleted: '秒传完成',
      foundResumable: '发现未完成的传输，可继续',
      analyzingCompleted: '✓ 分析完毕，可注入',
      analyzingFailed: '分析失败: {error}',
      retryingUpload: '正在重试注入...',
      resumingUpload: '正在恢复注入...',
      continuingUpload: '继续注入中...',
      restartingUpload: '重新开始注入...',
      preparingSession: '准备传输会话...',
      verifyingIntegrity: '验证数据完整性...',
      completingUpload: '正在完成注入...',
      chunkUploadFailed: '分片 {chunk} 传输失败: {error}',
      chunkRetrying: '分片 {chunk} 重试中... ({attempt}/{maxRetry}) - {delay}秒后重试',
    },
    time: {
      minutes: '{minutes}分{seconds}秒',
      seconds: '{seconds}秒',
    },
  },
}
