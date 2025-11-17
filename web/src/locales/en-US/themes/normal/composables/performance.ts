/**
 * Composables texts
 */
export const performance = {
  suggestions: {
    slowPageLoad: 'Page load time is too long, suggest optimizing file size and code splitting',
    slowRender: 'Component render time is too long, suggest using Vue.memo or optimizing computed properties',
    highMemory: 'Memory usage is too high, suggest checking for memory leaks or optimizing data structures',
  },
}

export const chunkedUpload = {
  errors: {
    sessionNotFound: 'Session ID does not exist, cannot verify upload status',
    sessionDataMissing: 'Session ID or file MD5 missing, cannot create session data',
    sessionIdMissing: 'Session ID does not exist, cannot upload chunk',
    unknownError: 'Unknown error',
    networkError: ' (Network error, please check network connection)',
    serverError: ' (Server error, please try again later)',
    rateLimited: ' (Request too frequent, please try again later)',
  },
}
