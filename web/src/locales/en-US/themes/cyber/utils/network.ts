/**
 * Network Utils Text - Cyber Style
 */
export const network = {
  http: {
    cancelDuplicate: 'Abort duplicate request: {method} {url}',
    requestTimeout: 'Request timeout cleanup',
    requestCancelled: 'Request aborted',
    pageSwitch: 'Page switch, abort all pending requests',
    cancelSpecificUrl: 'Abort request for specific URL',

    operationSuccess: 'Operation complete',
    operationFailed: 'Operation failed',
    requestFailed: 'Transfer failed',
    requestFailedWithStatus: 'Transfer failed ({status})',

    networkError: 'Neural link disconnected, check network',
    configError: 'Request config error',
    unknown: 'Unknown',

    ipNotInWhitelist: 'Your IP is not in access whitelist',
    ipInBlacklist: 'Your IP is in access blacklist',
    domainNotInWhitelist: 'Your domain is not in access whitelist',
    domainInBlacklist: 'Your domain is in access blacklist',
    userAccountDisabled: 'User account has been disabled',

    envVarReadError: 'Unable to read environment variable VITE_API_BASE_URL:',
    toastLoadError: 'Unable to load toast component',
  },
}
