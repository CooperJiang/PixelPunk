/**
 * 访问拒绝页面文案 - 赛博风格
 */
export const refuse = {
  defaultMessage: '您无权访问此系统',
  unknown: '未知',
  protocolHeader: '系统安全协议已启动 / SYSTEM SECURITY PROTOCOL',
  titles: {
    accountDisabled: '账号已被禁用',
    accessDenied: '访问被拒绝',
  },
  analysisTitle: {
    account: '账号状态说明:',
    system: '系统分析结果:',
  },
  messages: {
    ipNotInWhitelist: '您的IP不在访问白名单中',
    ipInBlacklist: '您的IP已被加入系统黑名单',
    domainNotInWhitelist: '您的域名不在访问白名单中',
    domainInBlacklist: '您的域名已被加入系统黑名单',
    accountDisabled: '您的账号已被禁用',
  },
  ip: {
    notInWhitelist: '您的IP地址不在系统访问白名单中',
    inBlacklist: '或者您的IP已被添加到系统黑名单',
    contactAdmin: '请联系系统管理员将您的IP加入白名单',
  },
  domain: {
    notInWhitelist: '您的访问域名不在系统白名单中',
    inBlacklist: '或者您的域名已被添加到系统黑名单',
    contactAdmin: '请联系系统管理员配置域名访问权限',
  },
  status: {
    accountStatus: '账号状态:',
    currentIp: '当前IP地址:',
    currentDomain: '当前域名:',
    disabled: '已禁用',
    detectionTime: '检测时间：',
  },
  accountDetails: {
    disabled: '您的账号已被系统管理员禁用',
    reason: '账号被禁用可能是由于违反平台规则或其他安全原因',
    contact: '如需解除限制，请联系系统管理员获取帮助',
  },
}
