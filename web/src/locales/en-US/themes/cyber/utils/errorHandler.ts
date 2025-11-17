/**
 * Error Handler Utils i18n - Cyber Theme
 */
export const errorHandler = {
  httpStatus: {
    400: 'Request parameter exception',
    401: 'Authentication credentials expired, please login again',
    403: 'Insufficient access permissions',
    404: 'Target resource not found',
    429: 'Request frequency too high, please try again later',
    500: 'System internal exception',
    502: 'Gateway connection exception',
    503: 'Service temporarily offline',
  },
  fallback: {
    unknown: 'Unknown exception',
    toastUnavailable: 'Toast system unavailable, error info:',
  },
}
