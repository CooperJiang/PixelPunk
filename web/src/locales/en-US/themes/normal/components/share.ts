/**
 * Share component
 */
export const share = {
  createShare: 'Create Share',
  shareSelected: 'Share Selected',
  cancel: 'Cancel',
  passwordDialog: {
    title: 'This Share Requires Password Verification',
    message: 'This share content is password protected, please enter password to continue',
    placeholder: 'Please enter share password',
    verify: 'Verify Password',
  },
  selectionToolbar: {
    selectedPrefix: 'Selected',
    selectedSuffix: 'items',
    selectAll: 'Select All',
    batchDownload: 'Pack Download',
  },
  selectionInfo: {
    folders: 'folders',
    images: 'resources',
  },
  messages: {
    loginRequired: 'Please login first before sharing',
    noItemsSelected: 'Please select folders or files to share first',
  },
  image: {
    alt: 'File',
    downloadError: 'Download failed',
    loadTimeout: 'File load timeout',
  },
  errors: {
    invalidLink: 'Invalid Share Link',
    invalidLinkDesc: 'Please check if your link is complete, or this share may have been deleted',
    loadFailed: 'Load Failed',
    loadFailedDesc: 'Unable to load share content, please try again later',
    maxViewsReached: 'Max Views Reached',
    maxViewsReachedDesc: 'This share has reached the maximum view count set by creator',
    shareExpired: 'Share Expired',
    shareExpiredDesc: 'This share link has expired, please contact sharer to reshare',
    accessRestricted: 'Access Restricted',
    accessRestrictedDesc: 'You do not have permission to access this share',
    shareNotFound: 'Share Not Found',
    shareNotFoundDesc: 'This share has been deleted or link is invalid',
    shareInvalid: 'Share Invalid',
    shareInvalidDesc: 'This share has expired, please contact sharer to reshare',
    passwordVerifySuccess: 'Password verification successful, loading content...',
    passwordVerifyNoToken: 'Password verification successful, but failed to get access token',
    passwordVerifyFailed: 'Password verification failed',
    passwordVerifyFailedRetry: 'Password verification failed, please try again',
    visitorSubmitSuccess: 'Information submitted successfully, thank you for your cooperation',
    visitorSubmitFailed: 'Failed to submit information, please try again',
  },
  defaults: {
    shareName: 'Share',
    shareTitle: 'File Share',
    shareDescription: 'Share contains {count} files',
  },
}
