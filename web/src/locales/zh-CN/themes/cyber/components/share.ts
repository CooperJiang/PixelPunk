/**
 */
export const share = {
  createShare: '建立分享链路',
  shareSelected: '分享选中项',
  cancel: '中止操作',
  passwordDialog: {
    title: '访问受限：需要密码验证',
    message: '该分享已启用访问加密，请输入访问密钥继续',
    placeholder: '输入访问密钥',
    verify: '验证密钥',
  },
  selectionToolbar: {
    selectedPrefix: '已选中',
    selectedSuffix: '项',
    selectAll: '全选',
    batchDownload: '批量下载',
  },
  selectionInfo: {
    folders: '个数据容器',
    images: '个数据文件',
  },
  messages: {
    loginRequired: '请先建立神经链接后再进行分享',
    noItemsSelected: '请先选择要分享的数据容器或数据文件',
  },
  image: {
    alt: '数据文件',
    downloadError: '数据提取失败',
    loadTimeout: '数据加载超时',
  },
  errors: {
    invalidLink: '无效的分享链接',
    invalidLinkDesc: '链接格式错误或该分享已被清除',
    loadFailed: '数据加载失败',
    loadFailedDesc: '无法建立链接，请稍后重试',
    maxViewsReached: '访问次数耗尽',
    maxViewsReachedDesc: '该分享已达到最大访问次数限制',
    shareExpired: '分享已过期',
    shareExpiredDesc: '该链接已过有效期，请联系源节点重新分享',
    accessRestricted: '访问受限',
    accessRestrictedDesc: '权限不足，无法访问此分享',
    shareNotFound: '分享不存在',
    shareNotFoundDesc: '该分享已被清除或链接无效',
    shareInvalid: '分享已失效',
    shareInvalidDesc: '该分享已失效，请联系源节点重新分享',
    passwordVerifySuccess: '密钥验证成功，正在加载数据...',
    passwordVerifyNoToken: '密钥验证成功，但未获取到访问令牌',
    passwordVerifyFailed: '密钥验证失败',
    passwordVerifyFailedRetry: '密钥验证失败，请重试',
    visitorSubmitSuccess: '信息提交成功，感谢配合',
    visitorSubmitFailed: '信息提交失败，请重试',
  },
  defaults: {
    shareName: '分享',
    shareTitle: '数据分享',
    shareDescription: '分享包含 {count} 个数据文件',
  },
}
