/**
 */
export const share = {
  createShare: '创建分享',
  shareSelected: '分享选中项',
  cancel: '取消',
  passwordDialog: {
    title: '该分享需要密码验证',
    message: '该分享内容已被密码保护，请输入密码后继续访问',
    placeholder: '请输入分享密码',
    verify: '验证密码',
  },
  selectionToolbar: {
    selectedPrefix: '已选择',
    selectedSuffix: '项',
    selectAll: '全选',
    batchDownload: '打包下载',
  },
  selectionInfo: {
    folders: '个文件夹',
    images: '个资源',
  },
  messages: {
    loginRequired: '请先登录后再进行分享',
    noItemsSelected: '请先选择要分享的文件夹或文件',
  },
  image: {
    alt: '文件',
    downloadError: '下载失败',
    loadTimeout: '文件加载超时',
  },
  errors: {
    invalidLink: '无效的分享链接',
    invalidLinkDesc: '请检查您的链接是否完整，或者该分享可能已被删除',
    loadFailed: '加载失败',
    loadFailedDesc: '无法加载分享内容，请稍后重试',
    maxViewsReached: '访问次数已达上限',
    maxViewsReachedDesc: '该分享已达到创建者设置的最大访问次数',
    shareExpired: '分享已过期',
    shareExpiredDesc: '该分享链接已过期，请联系分享者重新分享',
    accessRestricted: '访问受限',
    accessRestrictedDesc: '您没有权限访问此分享',
    shareNotFound: '分享不存在',
    shareNotFoundDesc: '该分享已被删除或链接无效',
    shareInvalid: '分享已失效',
    shareInvalidDesc: '该分享已失效，请联系分享者重新分享',
    passwordVerifySuccess: '密码验证成功，正在加载内容...',
    passwordVerifyNoToken: '密码验证成功，但未获取到访问令牌',
    passwordVerifyFailed: '密码验证失败',
    passwordVerifyFailedRetry: '密码验证失败，请重试',
    visitorSubmitSuccess: '信息提交成功，感谢您的配合',
    visitorSubmitFailed: '提交信息失败，请重试',
  },
  defaults: {
    shareName: '分享',
    shareTitle: '文件分享',
    shareDescription: '分享包含 {count} 个文件',
  },
}
