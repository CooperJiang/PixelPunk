/**
 * 上传页面文案
 */
export const upload = {
  title: '上传文件',

  settings: {
    title: '上传设置',
    advanced: '高级设置',
    folder: '文件夹',
    folderHint: '支持多级目录，可搜索文件夹名称',
    guestFolderInfo: '游客上传的文件将存放在系统文件夹中，默认{level}访问',
    guestLimitInfo: '游客模式仅支持单文件上传（并发限制为1）',
    guestDisabled: '游客上传功能已禁用，请联系管理员',
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
    storageDurationHint: '游客上传必须选择存储时长，文件到期后将自动删除',
    required: '(必选)',
    autoOptimize: '自动优化文件（推荐）',
    autoOptimizeHint: '自动压缩和优化缩略图，减小文件大小而不明显降低质量',
    autoRemove: '上传成功后自动移除',
    autoRemoveHint: '上传成功的文件将自动从列表中移除，适合批量上传时保持列表简洁',
    duration: {
      permanent: '永久',
    },
  },

  uploadProgress: {
    status: {
      uploadSuccess: '上传成功',
      instantSuccess: '秒传成功',
      uploadFailed: '上传失败',
      analyzing: '正在分析文件...',
      paused: '发现未完成的上传，可以继续',
      retrying: '正在重试...',
      pending: '等待上传',
    },
    time: {
      seconds: '{count}秒',
      minutes: '{count}分钟',
      hoursMinutes: '{hours}小时{minutes}分钟',
    },
    remaining: '剩余',
    actions: {
      openInNewTab: '新窗口打开',
      copyLink: '复制文件链接',
      resume: '继续上传',
      retry: '重新上传',
      remove: '移除文件',
    },
    toast: {
      linkCopied: '链接已复制',
      copyFailed: '复制失败',
      uploadSuccessWithRemainingCount: '上传成功！剩余 {count} 次上传机会',
      watermarkApplyFailed: '文件已上传，但水印应用失败: {reason}',
      thumbnailGenerationFailed: '文件已上传，但缩略图生成失败: {reason}',
      fileTypeNotSupported: '文件类型不支持: {name}',
      fileSizeExceeded: '文件大小超限: {name} ({size})',
      instantUploadSuccess: '秒传成功: {name}',
      cacheCleared: '缓存已清理',
    },
  },

  dropZone: {
    title: '拖拽文件到此处上传',
    orText: '或者',
    clickToSelect: '点击选择文件',
    features: {
      batch: '批量上传',
      autoOptimize: '自动优化',
      secureStorage: '安全存储',
    },
    shortcut: {
      label: '快捷键',
      paste: '粘贴上传',
    },
  },

  actions: {
    startUpload: '开始上传',
    selectFiles: '选择文件',
    copyLinks: '复制链接',
    clearQueue: '清空队列',
    cancelUpload: '取消上传',
    clearCache: '清理断点续传缓存，释放本地存储空间',
  },

  queue: {
    title: '上传队列',
    filesCount: '{count} 个文件',
    concurrent: '{current}/{max} 并发',
    stats: {
      pending: '等待中',
      uploading: '上传中',
      success: '成功',
      failed: '失败',
      total: '总计',
      totalSize: '总大小',
    },
    empty: {
      title: '队列为空',
      desc: '拖拽文件、点击选择或粘贴截图即可开始上传',
      action: '选择文件',
    },
  },

  guide: {
    title: '上传提示',
    multipleWays: {
      title: '多种上传方式',
      desc: '支持拖拽、点击选择或使用 {key}+V 粘贴剪贴板图片',
    },
    smartChunking: {
      title: '智能分片',
      desc: '大文件自动分片上传，支持断点续传',
    },
    oneCopy: {
      title: '一键复制',
      desc: '批量上传后可一键复制所有文件链接',
    },
  },

  advanced: {
    title: '高级设置',
    storageDuration: '存储时长',
    storageDurationDesc: '文件在服务器上的保存时间',
    watermark: '添加水印',
    watermarkDesc: '为上传的图片添加水印',
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
      desc: '登录后可使用永久存储、自定义水印等功能',
      action: '立即登录',
    },
  },

  toast: {
    folderSelected: '已选择文件夹: {path}',
    folderPreselected: '文件夹已预选，可直接开始上传',
    sensitiveAutoDelete: '违规内容将被自动删除，请确保文件符合规范',
    sensitivePendingReview: '敏感内容将进入审核队列，审核期间暂不可访问',
    defaultHint: '请确保上传内容符合平台规范',
    guestDurationRequired: '游客上传必须选择存储时长',
    uploadCancelled: '上传已取消',
    noFilesInQueue: '队列中暂无需要上传的文件',
    allLinksCopied: '所有链接已复制到剪贴板',
    noLinksToCopy: '暂无可复制的文件链接',
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
    configure: '配置水印',
  },

  form: {
    storageDuration: {
      placeholder: '选择存储时长',
    },
  },

  smartUpload: {
    error: {
      fileTypeNotSupported: '文件类型不支持: {name} ({type})',
      fileSizeExceeded: '文件大小超限: {name} ({size})',
      chunkedUploadNotReady: '分片上传未初始化，请稍后重试',
      uploadFailed: '上传失败',
      analysisFailed: '文件分析失败',
      analysisTimeout: '文件分析超时',
    },
    status: {
      preparing: '准备上传',
      uploading: '上传中...',
      uploadingWithProgress: '上传中 {progress}% - {speed} - 剩余{remaining}',
      completed: '上传完成',
      failed: '上传失败: {error}',
      paused: '已暂停',
      cancelled: '已取消',
      retrying: '准备重试',
      userPaused: '用户暂停',
      userCancelled: '用户取消',
      analyzing: '正在分析文件...',
      analyzingProgress: '分析中 {progress}%',
      checkingInstant: '检查秒传...',
      checkingInstantProgress: '检查秒传 {progress}%',
      instantCompleted: '秒传完成',
      foundResumable: '发现未完成的上传，可以继续',
      analyzingCompleted: '✓ 分析完毕，可上传',
      analyzingFailed: '分析失败: {error}',
      retryingUpload: '正在重试上传...',
      resumingUpload: '正在恢复上传...',
      continuingUpload: '继续上传中...',
      restartingUpload: '重新开始上传...',
      preparingSession: '准备上传会话...',
      verifyingIntegrity: '验证上传完整性...',
      completingUpload: '正在完成上传...',
      chunkUploadFailed: '分片 {chunk} 上传失败: {error}',
      chunkRetrying: '分片 {chunk} 重试中... ({attempt}/{maxRetry}) - {delay}秒后重试',
    },
    time: {
      minutes: '{minutes}分{seconds}秒',
      seconds: '{seconds}秒',
    },
  },
}
