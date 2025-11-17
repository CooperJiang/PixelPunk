/**
 * 网络工具文案
 */
export const network = {
  http: {
    cancelDuplicate: '取消重复请求: {method} {url}',
    requestTimeout: '请求超时清理',
    requestCancelled: '请求被取消',
    pageSwitch: '页面切换，取消所有待处理请求',
    cancelSpecificUrl: '取消特定URL的请求',

    operationSuccess: '操作成功',
    operationFailed: '操作失败',
    requestFailed: '请求失败',
    requestFailedWithStatus: '请求失败 ({status})',

    networkError: '网络连接失败，请检查网络',
    configError: '请求配置错误',
    unknown: '未知',

    ipNotInWhitelist: '您的IP不在访问白名单中',
    ipInBlacklist: '您的IP在访问黑名单中',
    domainNotInWhitelist: '您的域名不在访问白名单中',
    domainInBlacklist: '您的域名在访问黑名单中',
    userAccountDisabled: '用户账号已被禁用',

    envVarReadError: '无法读取环境变量VITE_API_BASE_URL:',
    toastLoadError: '无法加载toast组件',
  },
}
