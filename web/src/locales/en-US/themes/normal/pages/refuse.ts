/**
 * Access denied page texts
 */
export const refuse = {
  defaultMessage: 'You do not have permission to access this system',
  unknown: 'Unknown',
  protocolHeader: 'System Security Protocol Activated / SYSTEM SECURITY PROTOCOL',
  titles: {
    accountDisabled: 'Account Disabled',
    accessDenied: 'Access Denied',
  },
  analysisTitle: {
    account: 'Account Status:',
    system: 'System Analysis Result:',
  },
  messages: {
    ipNotInWhitelist: 'Your IP is not in the access whitelist',
    ipInBlacklist: 'Your IP has been added to the system blacklist',
    domainNotInWhitelist: 'Your domain is not in the access whitelist',
    domainInBlacklist: 'Your domain has been added to the system blacklist',
    accountDisabled: 'Your account has been disabled',
  },
  ip: {
    notInWhitelist: 'Your IP address is not in the system access whitelist',
    inBlacklist: 'Or your IP has been added to the system blacklist',
    contactAdmin: 'Please contact system administrator to add your IP to whitelist',
  },
  domain: {
    notInWhitelist: 'Your access domain is not in the system whitelist',
    inBlacklist: 'Or your domain has been added to the system blacklist',
    contactAdmin: 'Please contact system administrator to configure domain access permissions',
  },
  status: {
    accountStatus: 'Account Status:',
    currentIp: 'Current IP Address:',
    currentDomain: 'Current Domain:',
    disabled: 'Disabled',
    detectionTime: 'Detection Time:',
  },
  accountDetails: {
    disabled: 'Your account has been disabled by system administrator',
    reason: 'Account may be disabled due to violation of platform rules or other security reasons',
    contact: 'If you need to lift restrictions, please contact system administrator for help',
  },
}
