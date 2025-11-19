/**
 * Upload page texts
 */
export const upload = {
  title: 'Upload Files',

  settings: {
    title: 'Upload Settings',
    advanced: 'Advanced Settings',
    folder: 'Folder',
    folderHint: 'Supports multi-level directories, can search folder names',
    guestFolderInfo: 'Guest uploaded files will be stored in system folder, default {level} access',
    guestLimitInfo: 'Guest mode only supports single file upload (concurrency limit is 1)',
    guestDisabled: 'Guest upload feature is disabled, please contact administrator',
    accessLevelLabel: 'Access Level Settings',
    accessLevelPlaceholder: 'Select access level',
    accessLevel: {
      public: '🌍 Public',
      publicDesc: 'Can be recommended and displayed, accessible by anyone via link',
      private: '🔗 Private',
      privateDesc: 'Will not be recommended, but can be shared and externally referenced normally',
      protected: '🔒 Protected',
      protectedDesc: 'Only visible after you login, external links cannot access',
    },
    storageDurationLabel: 'Storage Duration',
    storageDurationPlaceholder: 'Please select storage duration',
    storageDurationHint: 'Guest uploads must select storage duration, files will be automatically deleted after expiration',
    required: '(Required)',
    autoOptimize: 'Auto Optimize Files (Recommended)',
    autoOptimizeHint: 'Automatically compress and optimize thumbnails, reduce file size without noticeably reducing quality',
    autoRemove: 'Auto Remove After Upload Success',
    autoRemoveHint:
      'Successfully uploaded files will be automatically removed from list, suitable for keeping list clean during batch uploads',
    duration: {
      permanent: 'Permanent',
    },
  },

  uploadProgress: {
    status: {
      uploadSuccess: 'Upload Successful',
      instantSuccess: 'Instant Upload Successful',
      uploadFailed: 'Upload Failed',
      analyzing: 'Analyzing file...',
      paused: 'Found incomplete upload, can continue',
      retrying: 'Retrying...',
      pending: 'Waiting to Upload',
    },
    time: {
      seconds: '{count} seconds',
      minutes: '{count} minutes',
      hoursMinutes: '{hours} hours {minutes} minutes',
    },
    remaining: 'Remaining',
    actions: {
      openInNewTab: 'Open in New Tab',
      copyLink: 'Copy File Link',
      copyLinkTooltip: 'Original Link',
      copyMarkdownLink: 'Copy Markdown Link',
      copyMarkdownTooltip: 'MD Format',
      copyHtmlLink: 'Copy HTML Tag',
      copyHtmlTooltip: 'HTML Format',
      copyThumbnailLink: 'Copy Thumbnail Link',
      copyThumbnailTooltip: 'Thumbnail',
      resume: 'Resume Upload',
      retry: 'Retry Upload',
      remove: 'Remove File',
    },
    toast: {
      linkCopied: 'Link copied',
      markdownLinkCopied: 'Markdown link copied',
      htmlLinkCopied: 'HTML tag copied',
      thumbnailLinkCopied: 'Thumbnail link copied',
      copyFailed: 'Copy failed',
      uploadSuccessWithRemainingCount: 'Upload successful! {count} upload opportunities remaining',
      watermarkApplyFailed: 'File uploaded, but watermark application failed: {reason}',
      thumbnailGenerationFailed: 'File uploaded, but thumbnail generation failed: {reason}',
      fileTypeNotSupported: 'File type not supported: {name}',
      fileSizeExceeded: 'File size exceeded: {name} ({size})',
      instantUploadSuccess: 'Instant upload successful: {name}',
      cacheCleared: 'Cache cleared',
    },
  },

  dropZone: {
    title: 'Drag files here to upload',
    orText: 'or',
    clickToSelect: 'Click to select files',
    supportFolder: 'Supports dragging entire folder',
    features: {
      batch: 'Batch Upload',
      autoOptimize: 'Auto Optimize',
      secureStorage: 'Secure Storage',
    },
    shortcut: {
      label: 'Shortcut',
      paste: 'Paste Upload',
    },
  },

  actions: {
    startUpload: 'Start Upload',
    selectFiles: 'Select Files',
    selectFolder: 'Select Folder',
    copyLinks: 'Copy Links',
    copyMarkdownLinks: 'Copy All Markdown Links',
    copyHtmlLinks: 'Copy All HTML Tags',
    copyThumbnailLinks: 'Copy All Thumbnail Links',
    clearQueue: 'Clear Queue',
    cancelUpload: 'Cancel Upload',
    clearCache: 'Clear Cache',
  },

  queue: {
    title: 'Upload Queue',
    filesCount: '{count} files',
    concurrent: '{current}/{max} concurrent',
    stats: {
      pending: 'Pending',
      uploading: 'Uploading',
      success: 'Success',
      failed: 'Failed',
      total: 'Total',
      totalSize: 'Total Size',
    },
    empty: {
      title: 'Queue is Empty',
      desc: 'Drag files, click to select or paste screenshots to start uploading',
      action: 'Select Files',
    },
  },

  guide: {
    title: 'Upload Tips',
    multipleWays: {
      title: 'Multiple Upload Methods',
      desc: 'Supports drag-and-drop, click to select or use {key}+V to paste clipboard images',
    },
    smartChunking: {
      title: 'Smart Chunking',
      desc: 'Large files automatically chunked for upload, supports resumable upload',
    },
    oneCopy: {
      title: 'One-Click Copy',
      desc: 'After batch upload, copy all file links with one click',
    },
  },

  advanced: {
    title: 'Advanced Settings',
    storageDuration: 'Storage Duration',
    storageDurationDesc: 'File storage time on server',
    watermark: 'Add Watermark',
    watermarkDesc: 'Add watermark to uploaded images',
    watermarkInfo: {
      type: 'Type:',
      content: 'Content:',
      position: 'Position:',
    },
    watermarkType: {
      text: 'Text',
      image: 'Image',
    },
    loginPrompt: {
      title: 'More Advanced Features',
      desc: 'After login, you can use permanent storage, custom watermark and other features',
      action: 'Login Now',
    },
  },

  folder: {
    browserNotSupported: 'Your browser does not support folder upload',
    noValidFiles: 'No uploadable files found in folder',
    parsing: 'Parsing folder...',
    addSuccess: 'Added {count} files from folder',
    dialog: {
      title: 'Folder Upload Confirmation',
      totalFiles: 'Total Files',
      validFiles: 'Uploadable Files',
      invalidFiles: 'Unsupported Files',
      confirmMessage: 'Found <strong class="highlight-text">{count}</strong> uploadable files, add to upload queue?',
      warningMessage: 'Filtered {count} unsupported files',
      cancel: 'Cancel',
      confirm: 'Confirm Upload',
    },
  },

  toast: {
    folderSelected: 'Folder selected: {path}',
    folderPreselected: 'Folder preselected, you can start uploading directly',
    sensitiveAutoDelete: 'Inappropriate content will be automatically deleted, please ensure files comply with regulations',
    sensitivePendingReview: 'Sensitive content will enter review queue, temporarily inaccessible during review',
    defaultHint: 'Please ensure uploaded content complies with platform regulations',
    guestDurationRequired: 'Guest uploads must select storage duration',
    uploadCancelled: 'Upload cancelled',
    noFilesInQueue: 'No files in queue to upload',
    allLinksCopied: 'All links copied to clipboard',
    markdownLinksCopied: 'All Markdown links copied',
    htmlLinksCopied: 'All HTML tags copied',
    thumbnailLinksCopied: 'All thumbnail links copied',
    noLinksToCopy: 'No file links available to copy',
  },

  duration: {
    permanent: 'Permanent',
  },

  watermark: {
    position: {
      'top-left': 'Top Left',
      'top-center': 'Top Center',
      'top-right': 'Top Right',
      'middle-left': 'Middle Left',
      'middle-center': 'Center',
      'middle-right': 'Middle Right',
      'bottom-left': 'Bottom Left',
      'bottom-center': 'Bottom Center',
      'bottom-right': 'Bottom Right',
      custom: 'Custom Position',
    },
    configure: 'Configure Watermark',
  },

  form: {
    storageDuration: {
      placeholder: 'Select storage duration',
    },
  },

  smartUpload: {
    error: {
      fileTypeNotSupported: 'File type not supported: {name} ({type})',
      fileSizeExceeded: 'File size exceeded: {name} ({size})',
      chunkedUploadNotReady: 'Chunked upload not initialized, please try again later',
      uploadFailed: 'Upload failed',
      analysisFailed: 'File analysis failed',
      analysisTimeout: 'File analysis timeout',
    },
    status: {
      preparing: 'Preparing Upload',
      uploading: 'Uploading...',
      uploadingWithProgress: 'Uploading {progress}% - {speed} - {remaining} remaining',
      completed: 'Upload Completed',
      failed: 'Upload Failed: {error}',
      paused: 'Paused',
      cancelled: 'Cancelled',
      retrying: 'Preparing Retry',
      userPaused: 'User Paused',
      userCancelled: 'User Cancelled',
      analyzing: 'Analyzing file...',
      analyzingProgress: 'Analyzing {progress}%',
      checkingInstant: 'Checking instant upload...',
      checkingInstantProgress: 'Checking instant upload {progress}%',
      instantCompleted: 'Instant Upload Completed',
      foundResumable: 'Found incomplete upload, can continue',
      analyzingCompleted: '✓ Analysis complete, ready to upload',
      analyzingFailed: 'Analysis failed: {error}',
      retryingUpload: 'Retrying upload...',
      resumingUpload: 'Resuming upload...',
      continuingUpload: 'Continuing upload...',
      restartingUpload: 'Restarting upload...',
      preparingSession: 'Preparing upload session...',
      verifyingIntegrity: 'Verifying upload integrity...',
      completingUpload: 'Completing upload...',
      chunkUploadFailed: 'Chunk {chunk} upload failed: {error}',
      chunkRetrying: 'Chunk {chunk} retrying... ({attempt}/{maxRetry}) - retry in {delay} seconds',
    },
    time: {
      minutes: '{minutes}m {seconds}s',
      seconds: '{seconds}s',
    },
  },
}
