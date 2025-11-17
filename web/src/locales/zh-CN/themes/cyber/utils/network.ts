/**
 * 网络工具文案 - 赛博风格
 */
export const network = {
  http: {
    cancelDuplicate: '中止重复请求: {method} {url}',
    requestTimeout: '请求超时清理',
    requestCancelled: '请求已中止',
    pageSwitch: '页面切换，中止所有待处理请求',
    cancelSpecificUrl: '中止特定URL的请求',

    operationSuccess: '操作完成',
    operationFailed: '操作失败',
    requestFailed: '传输失败',
    requestFailedWithStatus: '传输失败 ({status})',

    networkError: '神经链接断开，请检查网络',
    configError: '请求配置错误',
    unknown: '未知',

    ipNotInWhitelist: '您的IP不在访问白名单中',
    ipInBlacklist: '您的IP在访问黑名单中',
    domainNotInWhitelist: '您的域名不在访问白名单中',
    domainInBlacklist: '您的域名在访问黑名单中',
    userAccountDisabled: '用户账号已被禁用',

    envVarReadError: 'Unable to read environment variable VITE_API_BASE_URL:',
    toastLoadError: 'Unable to load toast component',
  },
}
