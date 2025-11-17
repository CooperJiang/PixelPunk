/**
 * 消息系统多语言配置 - 赛博主题
 */
export const message = {
  types: {
    system: {
      maintenance: '系统维护',
      update: '系统更新',
      announcement: '系统公告',
    },
    account: {
      register: '注册成功',
      storage_granted: '存储空间授予',
      bandwidth_granted: '流量额度授予',
    },
    content: {
      review_approved: '内容审核通过',
      review_rejected: '内容审核拒绝',
      review_pending: '等待审核',
    },
    storage: {
      quota_warning: '存储空间预警',
      quota_exceeded: '存储空间超限',
      quota_increased: '存储空间增加',
      quota_decreased: '存储空间调整',
    },
    file: {
      deleted_by_admin: '数据删除',
      batch_deleted_by_admin: '批量数据删除',
      hard_deleted_by_admin: '永久删除数据',
      expiry_warning: '数据过期提醒',
      thumbnail_failed: '缩略图生成失败',
    },
    security: {
      login_alert: '登录警报',
      password_changed: '密码已变更',
    },
    apikey: {
      created: 'API密钥已创建',
      deleted: 'API密钥已删除',
      regenerated: 'API密钥已重新生成',
      disabled: 'API密钥已禁用',
      enabled: 'API密钥已启用',
    },
    random_api: {
      created: '随机图API已创建',
      deleted: '随机图API已删除',
      disabled: '随机图API已禁用',
      enabled: '随机图API已启用',
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
