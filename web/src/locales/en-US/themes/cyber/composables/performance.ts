/**
 * Composables - Cyberpunk Style
 */
export const performance = {
  suggestions: {
    slowPageLoad: 'Page load delay too high, recommend optimizing data package size and code splitting',
    slowRender: 'Component render delay too high, recommend using Vue.memo or optimizing computed properties',
    highMemory: 'Memory usage too high, recommend checking memory leaks or optimizing data structures',
  },
}

export const chunkedUpload = {
  errors: {
    sessionNotFound: 'Session ID not found, cannot verify upload status',
    sessionDataMissing: 'Session ID or data fingerprint missing, cannot create session data',
    sessionIdMissing: 'Session ID not found, cannot upload data chunks',
    unknownError: 'Unknown error',
    networkError: ' (Network connection error, please check network)',
    serverError: ' (Server response error, please retry later)',
    rateLimited: ' (Request frequency too high, please retry later)',
  },
}
