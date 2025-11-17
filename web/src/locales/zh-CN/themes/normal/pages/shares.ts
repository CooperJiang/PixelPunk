/**
 * 分享管理页面文案
 */
export const shares = {
  title: '我的分享',
  search: '搜索分享',
  page: {
    title: '我的分享',
    subtitle: '管理和查看您创建的所有分享链接',
  },
  filter: {
    all: '全部',
    normal: '正常',
    expired: '已过期',
    deleted: '已删除',
    disabled: '已禁用',
  },
  status: {
    normal: '正常',
    expired: '已过期',
    deleted: '已删除',
    disabled: '已禁用',
    unknown: '未知',
  },
  card: {
    viewCount: '{count} 次访问',
    maxViews: '限制 {max} 次',
    passwordProtected: '密码保护',
    collectVisitor: '收集访客',
    accessNotification: '访问通知',
    copyLink: '复制链接',
    folderCount: '{count} 个文件夹',
    fileCount: '{count} 个资源',
    viewShare: '查看分享',
    visitorInfo: '访客信息',
    delete: '删除',
  },
  empty: {
    title: '暂无分享内容',
    subtitle: '您还没有创建任何分享，可以从文件夹页面创建分享',
    goToFolders: '前往文件夹',
  },
  delete: {
    confirm: '您确定要删除分享 {name} 吗？',
    forceDelete: '彻底删除（不可恢复）',
    cancel: '取消',
    confirmButton: '删除',
  },
  unnamed: '未命名分享',
  toast: {
    loadFailed: '加载分享列表失败',
    linkCopied: '分享链接已复制到剪贴板',
    deleteSuccess: '删除分享成功',
    deleteFailed: '删除分享失败',
  },
  dialog: {
    deleteTitle: '确认删除',
    deleteMessage: '确定要删除分享 "{name}" 吗？',
    deleteHint: '删除后将无法恢复',
  },
  visitor: {
    dialogTitle: '访客信息列表',
    search: '搜索访客姓名、邮箱或IP地址',
    deleteConfirm: '确定要删除这条访客信息吗？',
    deleteHint: '删除后将无法恢复',
    loadFailed: '加载访客信息失败，请重试',
    deleteSuccess: '删除访客信息成功',
    deleteFailed: '删除访客信息失败，请重试',
    table: {
      name: '访客姓名',
      email: '邮箱',
      ip: 'IP地址',
      visits: '访问次数',
      lastVisit: '最近访问时间',
      actions: '操作',
    },
    pagination: {
      info: '共 {total} 项，当前显示 {from}-{to}',
    },
    empty: {
      title: '暂无访客信息',
      desc: '还没有访客提交过信息，或者您的筛选条件没有匹配结果',
    },
  },
}
