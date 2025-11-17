/**
 * Network utility texts
 */
export const network = {
  http: {
    cancelDuplicate: 'Cancel duplicate request: {method} {url}',
    requestTimeout: 'Request timeout cleanup',
    requestCancelled: 'Request cancelled',
    pageSwitch: 'Page switched, cancel all pending requests',
    cancelSpecificUrl: 'Cancel requests for specific URL',

    operationSuccess: 'Operation Successful',
    operationFailed: 'Operation Failed',
    requestFailed: 'Request Failed',
    requestFailedWithStatus: 'Request Failed ({status})',

    networkError: 'Network connection failed, please check network',
    configError: 'Request configuration error',
    unknown: 'Unknown',

    ipNotInWhitelist: 'Your IP is not in the access whitelist',
    ipInBlacklist: 'Your IP is in the access blacklist',
    domainNotInWhitelist: 'Your domain is not in the access whitelist',
    domainInBlacklist: 'Your domain is in the access blacklist',
    userAccountDisabled: 'User account has been disabled',

    envVarReadError: 'Cannot read environment variable VITE_API_BASE_URL:',
    toastLoadError: 'Cannot load toast component',
  },
}
