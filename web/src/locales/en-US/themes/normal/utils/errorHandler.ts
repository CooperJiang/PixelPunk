/**
 * Error handler utility i18n configuration - Normal theme
 */
export const errorHandler = {
  httpStatus: {
    400: 'Request Parameter Error',
    401: 'Login Expired, Please Login Again',
    403: 'No Access Permission',
    404: 'Requested Resource Not Found',
    429: 'Too Many Requests, Please Try Again Later',
    500: 'Internal Server Error',
    502: 'Bad Gateway',
    503: 'Service Temporarily Unavailable',
  },
  fallback: {
    unknown: 'Unknown Error',
    toastUnavailable: 'Toast unavailable, error message:',
  },
}
