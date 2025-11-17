/**
 * 消息系统多语言配置 - 常规主题
 */
export const message = {
  types: {
    system: {
      maintenance: '系统维护',
      update: '系统更新',
      announcement: '系统公告',
    },
    account: {
      register: '注册欢迎',
      storage_granted: '存储赠送',
      bandwidth_granted: '流量赠送',
    },
    content: {
      review_approved: '审核通过',
      review_rejected: '审核拒绝',
      review_pending: '待审核',
    },
    storage: {
      quota_warning: '存储预警',
      quota_exceeded: '存储超限',
      quota_increased: '存储增加',
      quota_decreased: '存储调整',
    },
    file: {
      deleted_by_admin: '文件删除',
      batch_deleted_by_admin: '批量删除',
      hard_deleted_by_admin: '永久删除',
      expiry_warning: '文件过期提醒',
      thumbnail_failed: '缩略图生成失败',
    },
    security: {
      login_alert: '登录提醒',
      password_changed: '密码变更',
    },
    apikey: {
      created: 'API密钥创建',
      deleted: 'API密钥删除',
      regenerated: 'API密钥重新生成',
      disabled: 'API密钥禁用',
      enabled: 'API密钥启用',
    },
    random_api: {
      created: '随机图API创建',
      deleted: '随机图API删除',
      disabled: '随机图API禁用',
      enabled: '随机图API启用',
    },
    share: {
      expiry_warning: '分享即将过期',
    },
  },
  priority: {
    high: '高',
    normal: '中',
    low: '低',
  },
  status: {
    unread: '未读',
    read: '已读',
    deleted: '已删除',
    unknown: '未知',
  },
  time: {
    justNow: '刚刚',
    minutesAgo: '{n}分钟前',
    hoursAgo: '{n}小时前',
    daysAgo: '{n}天前',
  },
}
