/**
 * Composables 文案 - 赛博风格
 */
export const performance = {
  suggestions: {
    slowPageLoad: '页面加载延迟过高，建议优化数据包大小和代码分割',
    slowRender: '组件渲染延迟过高，建议使用Vue.memo或优化计算属性',
    highMemory: '内存占用过高，建议检查内存泄漏或优化数据结构',
  },
}

export const chunkedUpload = {
  errors: {
    sessionNotFound: '会话ID不存在，无法验证上传状态',
    sessionDataMissing: '会话ID或数据指纹缺失，无法创建会话数据',
    sessionIdMissing: '会话ID不存在，无法上传数据分片',
    unknownError: '未知错误',
    networkError: ' (网络连接异常，请检查网络)',
    serverError: ' (服务器响应异常，请稍后重试)',
    rateLimited: ' (请求频率过高，请稍后重试)',
  },
}
