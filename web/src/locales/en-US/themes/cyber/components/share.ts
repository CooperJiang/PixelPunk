/**
 * Share Component
 */
export const share = {
  createShare: 'Establish Share Link',
  shareSelected: 'Share Selected Items',
  cancel: 'Abort Operation',
  passwordDialog: {
    title: 'Access Restricted: Password Required',
    message: 'This share has access encryption enabled, enter access key to continue',
    placeholder: 'Enter access key',
    verify: 'Verify Key',
  },
  selectionToolbar: {
    selectedPrefix: 'Selected',
    selectedSuffix: ' items',
    selectAll: 'Select All',
    batchDownload: 'Batch Download',
  },
  selectionInfo: {
    folders: ' data containers',
    images: ' data files',
  },
  messages: {
    loginRequired: 'Please establish neural link before sharing',
    noItemsSelected: 'Please select data containers or data files to share',
  },
  image: {
    alt: 'Data File',
    downloadError: 'Data extraction failed',
    loadTimeout: 'Data load timeout',
  },
  errors: {
    invalidLink: 'Invalid share link',
    invalidLinkDesc: 'Link format error or share has been cleared',
    loadFailed: 'Data load failed',
    loadFailedDesc: 'Cannot establish link, please retry later',
    maxViewsReached: 'Access count exhausted',
    maxViewsReachedDesc: 'This share has reached maximum access count limit',
    shareExpired: 'Share expired',
    shareExpiredDesc: 'This link has expired, contact source node to reshare',
    accessRestricted: 'Access Restricted',
    accessRestrictedDesc: 'Insufficient permissions, cannot access this share',
    shareNotFound: 'Share not found',
    shareNotFoundDesc: 'This share has been cleared or link is invalid',
    shareInvalid: 'Share invalid',
    shareInvalidDesc: 'This share has expired, contact source node to reshare',
    passwordVerifySuccess: 'Key verification successful, loading data...',
    passwordVerifyNoToken: 'Key verification successful, but access token not obtained',
    passwordVerifyFailed: 'Key verification failed',
    passwordVerifyFailedRetry: 'Key verification failed, please retry',
    visitorSubmitSuccess: 'Info submission successful, thanks for cooperation',
    visitorSubmitFailed: 'Info submission failed, please retry',
  },
  defaults: {
    shareName: 'Share',
    shareTitle: 'Data Share',
    shareDescription: 'Share contains {count} data files',
  },
}
