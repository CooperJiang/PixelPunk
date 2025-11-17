/**
 * Composables 文案
 */
export const performance = {
  suggestions: {
    slowPageLoad: '页面加载时间过长，建议优化文件大小和代码分割',
    slowRender: '组件渲染时间过长，建议使用Vue.memo或优化计算属性',
    highMemory: '内存使用过高，建议检查内存泄漏或优化数据结构',
  },
}

export const chunkedUpload = {
  errors: {
    sessionNotFound: '会话ID不存在，无法验证上传状态',
    sessionDataMissing: '会话ID或文件MD5缺失，无法创建会话数据',
    sessionIdMissing: '会话ID不存在，无法上传分片',
    unknownError: '未知错误',
    networkError: ' (网络错误，请检查网络连接)',
    serverError: ' (服务器错误，请稍后重试)',
    rateLimited: ' (请求过于频繁，请稍后重试)',
  },
}
