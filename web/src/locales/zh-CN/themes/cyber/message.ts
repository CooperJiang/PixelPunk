/**
 * Message Translations - Chinese Cyber Theme
 */
export const message = {
  system: {
    maintenance: {
      title: '系统维护协议',
      content: '神经网络将于 {time} 进入维护模式，预计持续 {duration}。维护期间部分节点可能离线，请提前同步数据。',
    },
    announcement: {
      title: '系统广播',
      content: '{content}',
    },
  },
  account: {
    register: {
      title: '欢迎接入 PixelPunk 神经网络',
      content: '连接已建立！您的节点已成功注册。系统为您分配了 {storage} 的数据空间，开始您的数据管理之旅吧！',
    },
    storageGranted: {
      title: '数据空间已扩展',
      content: '系统升级！您的数据空间已扩展 {size}，当前总容量 {total}。',
    },
    bandwidthGranted: {
      title: '带宽已增强',
      content: '网络增强！您的每月带宽已增加 {size}，当前每月带宽 {total}。',
    },
  },
  content: {
    reviewApproved: {
      title: '数据审核通过',
      content: '您的数据「{content_name}」已通过审核，现已在网络中可见。',
    },
    reviewRejected: {
      title: '数据审核失败',
      content: '抱歉，您的数据「{content_name}」未通过审核。原因：{reason}',
    },
    reviewPending: {
      title: '数据审核中',
      content: '您的数据「{content_name}」已提交审核，系统正在处理。',
    },
  },
  storage: {
    quotaWarning: {
      title: '数据空间告警',
      content: '您的数据空间已使用 {used_percent}%（{used}/{total}），建议清理或升级容量。',
    },
    quotaIncreased: {
      title: '数据空间已升级',
      content: '您的数据空间已从 {old_size} 升级至 {new_size}。',
    },
    quotaDecreased: {
      title: '数据空间已调整',
      content: '您的数据空间已从 {old_size} 调整至 {new_size}。',
    },
  },
  file: {
    deletedByAdmin: {
      title: '数据已被移除',
      content: '您的数据单元「{file_name}」因违反网络协议被管理员移除。',
    },
    batchDeletedByAdmin: {
      title: '数据已被批量移除',
      content: '您有 {count} 个数据单元因违反网络协议被管理员移除。',
    },
    hardDeletedByAdmin: {
      title: '数据已被永久删除',
      content: '您的数据单元「{file_name}」已被管理员永久删除且无法恢复。',
    },
    expiryWarning: {
      title: '您有 {count} 个数据单元即将过期',
      content: '以下数据单元即将过期，请及时备份：\n{file_list}\n\n过期数据将被自动清除且无法恢复。',
    },
    thumbnailFailed: {
      title: '预览生成失败',
      content: '数据单元「{file_name}」的预览生成失败，但数据本身未受影响。',
    },
  },
  security: {
    loginAlert: {
      title: '新节点接入提醒',
      content: '检测到您的账号于 {time} 从 {ip} ({location}) 建立连接。如非本人操作，请立即修改访问密钥。',
    },
  },
  apikey: {
    created: {
      title: 'API 密钥已生成',
      content: '您生成了新的 API 密钥「{key_name}」。请妥善保管您的访问凭证。',
    },
    deleted: {
      title: 'API 密钥已销毁',
      content: 'API 密钥「{key_name}」已被销毁。',
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
      title: 'API 密钥已激活',
      content: 'API 密钥「{key_name}」已被激活。',
    },
  },
  randomApi: {
    created: {
      title: '随机数据 API 已创建',
      content: '您创建了随机数据 API「{api_name}」。',
    },
    deleted: {
      title: '随机数据 API 已删除',
      content: '随机数据 API「{api_name}」已被删除。',
    },
    disabled: {
      title: '随机数据 API 已停用',
      content: '随机数据 API「{api_name}」已被停用。',
    },
    enabled: {
      title: '随机数据 API 已激活',
      content: '随机数据 API「{api_name}」已被激活。',
    },
  },
  share: {
    expiryWarning: {
      title: '共享链接即将过期',
      content: '您的共享链接「{share_name}」将于 {expires_at} 过期。',
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
