/**
 * Message Translations - Chinese Normal Theme
 */
export const message = {
  system: {
    maintenance: {
      title: '系统维护通知',
      content: '系统将于 {time} 进行维护，预计耗时 {duration}。维护期间部分功能可能无法使用，请提前做好准备。',
    },
    announcement: {
      title: '系统公告',
      content: '{content}',
    },
  },
  account: {
    register: {
      title: '欢迎加入 PixelPunk',
      content: '感谢您注册 PixelPunk！您的账号已创建成功。我们为您准备了 {storage} 的存储空间，开始您的图片管理之旅吧！',
    },
    storageGranted: {
      title: '存储空间已增加',
      content: '恭喜！您的存储空间已增加 {size}，当前总容量为 {total}。',
    },
    bandwidthGranted: {
      title: '流量已增加',
      content: '恭喜！您的每月流量已增加 {size}，当前每月流量为 {total}。',
    },
  },
  content: {
    reviewApproved: {
      title: '内容审核通过',
      content: '您的内容「{content_name}」已通过审核，现已公开显示。',
    },
    reviewRejected: {
      title: '内容审核未通过',
      content: '很抱歉，您的内容「{content_name}」未通过审核。原因：{reason}',
    },
    reviewPending: {
      title: '内容待审核',
      content: '您的内容「{content_name}」已提交审核，我们会尽快处理。',
    },
  },
  storage: {
    quotaWarning: {
      title: '存储空间不足',
      content: '您的存储空间已使用 {used_percent}%（{used}/{total}），请及时清理或升级容量。',
    },
    quotaIncreased: {
      title: '存储空间已升级',
      content: '您的存储空间已从 {old_size} 升级到 {new_size}。',
    },
    quotaDecreased: {
      title: '存储空间已调整',
      content: '您的存储空间已从 {old_size} 调整到 {new_size}。',
    },
  },
  file: {
    deletedByAdmin: {
      title: '文件已被删除',
      content: '您的文件「{file_name}」因违反社区规则被管理员删除。',
    },
    batchDeletedByAdmin: {
      title: '文件已被批量删除',
      content: '您有 {count} 个文件因违反社区规则被管理员删除。',
    },
    hardDeletedByAdmin: {
      title: '文件已被永久删除',
      content: '您的文件「{file_name}」已被管理员永久删除且无法恢复。',
    },
    expiryWarning: {
      title: '您有 {count} 个文件即将过期',
      content: '以下文件即将过期，请及时备份：\n{file_list}\n\n过期文件将被自动删除且无法恢复。',
    },
    thumbnailFailed: {
      title: '缩略图生成失败',
      content: '文件「{file_name}」的缩略图生成失败，但文件本身未受影响。',
    },
  },
  security: {
    loginAlert: {
      title: '新设备登录提醒',
      content: '检测到您的账号于 {time} 从 {ip} ({location}) 登录。如非本人操作，请立即修改密码。',
    },
  },
  apikey: {
    created: {
      title: 'API 密钥已创建',
      content: '您创建了新的 API 密钥「{key_name}」。请妥善保管您的密钥。',
    },
    deleted: {
      title: 'API 密钥已删除',
      content: 'API 密钥「{key_name}」已被删除。',
    },
    regenerated: {
      title: 'API 密钥已重新生成',
      content: 'API 密钥「{key_name}」已重新生成，旧密钥将立即失效。',
    },
    disabled: {
      title: 'API 密钥已停用',
      content: 'API 密钥「{key_name}」已被停用。',
    },
    enabled: {
      title: 'API 密钥已启用',
      content: 'API 密钥「{key_name}」已被启用。',
    },
  },
  randomApi: {
    created: {
      title: '随机图片 API 已创建',
      content: '您创建了随机图片 API「{api_name}」。',
    },
    deleted: {
      title: '随机图片 API 已删除',
      content: '随机图片 API「{api_name}」已被删除。',
    },
    disabled: {
      title: '随机图片 API 已停用',
      content: '随机图片 API「{api_name}」已被停用。',
    },
    enabled: {
      title: '随机图片 API 已启用',
      content: '随机图片 API「{api_name}」已被启用。',
    },
  },
  share: {
    expiryWarning: {
      title: '分享链接即将过期',
      content: '您的分享链接「{share_name}」将于 {expires_at} 过期。',
    },
  },
  messageList: {
    title: '消息中心',
    total: '总数',
    unread: '未读',
    allStatus: '全部状态',
    allTypes: '全部类型',
    markAllRead: '全部标为已读',
    emptyTitle: '暂无消息',
    emptyDescription: '您还没有收到任何消息',
    loadingText: '加载中...',
    columns: {
      type: '类型',
      title: '标题',
      content: '内容',
      status: '状态',
      priority: '优先级',
      time: '时间',
      actions: '操作',
    },
    statusLabels: {
      unread: '未读',
      read: '已读',
    },
    priorityLabels: {
      high: '高',
      normal: '普通',
      low: '低',
    },
    actions: {
      markRead: '标为已读',
      delete: '删除',
    },
    dialog: {
      deleteTitle: '确认删除',
      deleteMessage: '确定要删除这条消息吗？此操作无法撤销。',
      cancel: '取消',
      confirmDelete: '删除',
    },
    toasts: {
      deleted: '消息已删除',
      deleteFailed: '删除失败，请重试',
      markReadSuccess: '已标记为已读',
      markReadFailed: '标记失败，请重试',
      markAllReadSuccess: '全部消息已标记为已读',
      markAllReadFailed: '批量标记失败，请重试',
    },
  },
}
