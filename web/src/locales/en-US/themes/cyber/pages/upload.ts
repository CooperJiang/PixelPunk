/**
 * Upload Page Text - Cyber Style
 */
export const upload = {
  title: 'Data Injection',

  settings: {
    title: 'Injection Settings',
    advanced: 'Advanced Config',
    folder: 'Data Folder',
    folderHint: 'Supports multi-level directories, searchable by folder name',
    guestFolderInfo: 'Guest injected data units will be stored in system folder, default {level} access',
    guestLimitInfo: 'Guest mode only supports single data unit injection (concurrency limit 1)',
    guestDisabled: 'Guest injection feature disabled, please contact administrator',
    accessLevelLabel: 'Access Permission Settings',
    accessLevelPlaceholder: 'Select access permission',
    accessLevel: {
      public: '🌍 Public',
      publicDesc: 'Can be recommended and displayed, anyone can access via link',
      private: '🔗 Private',
      privateDesc: 'Will not be recommended, but can be shared and externally referenced normally',
      protected: '🔒 Protected',
      protectedDesc: 'Only visible after you log in, external links cannot access',
    },
    storageDurationLabel: 'Storage Duration',
    storageDurationPlaceholder: 'Please select storage duration',
    storageDurationHint: 'Guest injection must select storage duration, data units will be automatically deleted after expiration',
    required: '(Required)',
    autoOptimize: 'Auto Optimize Data Units (Recommended)',
    autoOptimizeHint: 'Auto compress and optimize thumbnails, reduce data unit size without significantly lowering quality',
    autoRemove: 'Auto Remove After Success',
    autoRemoveHint: 'Successfully injected data units will be automatically removed from list, suitable for keeping list clean during batch injection',
    duration: {
      permanent: 'Permanent',
    },
  },

  uploadProgress: {
    status: {
      uploadSuccess: 'Injection Successful',
      instantSuccess: 'Instant Injection Successful',
      uploadFailed: 'Injection Failed',
      analyzing: 'Analyzing data unit...',
      paused: 'Found unfinished injection, can continue',
      retrying: 'Retrying injection...',
      pending: 'Waiting for Injection',
    },
    time: {
      seconds: '{count}s',
      minutes: '{count}min',
      hoursMinutes: '{hours}h{minutes}min',
    },
    remaining: 'Remaining',
    actions: {
      openInNewTab: 'Open in New Tab',
      copyLink: 'Copy Data Unit Link',
      resume: 'Continue Injection',
      retry: 'Retry Injection',
      remove: 'Remove Data Unit',
    },
    toast: {
      linkCopied: 'Link copied',
      copyFailed: 'Copy failed',
      uploadSuccessWithRemainingCount: 'Data injection successful! {count} injection opportunities remaining',
      watermarkApplyFailed: 'Data unit injected, but watermark application failed: {reason}',
      thumbnailGenerationFailed: 'Data unit injected, but thumbnail generation failed: {reason}',
      fileTypeNotSupported: 'Data type not supported: {name}',
      fileSizeExceeded: 'Data size exceeded: {name} ({size})',
      instantUploadSuccess: 'Instant upload successful: {name}',
      cacheCleared: 'Cache cleared',
    },
  },

  dropZone: {
    title: 'Drag data units here to inject',
    orText: 'or',
    clickToSelect: 'Click to select data',
    features: {
      batch: 'Batch Injection',
      autoOptimize: 'Auto Optimize',
      secureStorage: 'Secure Storage',
    },
    shortcut: {
      label: 'Shortcut',
      paste: 'Paste to Inject',
    },
  },

  actions: {
    startUpload: 'Start Injection',
    selectFiles: 'Select Data',
    copyLinks: 'Copy Links',
    clearQueue: 'Clear Queue',
    cancelUpload: 'Cancel Injection',
    clearCache: 'Clear resumable upload cache, free up local storage space',
  },

  queue: {
    title: 'Injection Queue',
    filesCount: '{count} data units',
    concurrent: '{current}/{max} concurrent',
    stats: {
      pending: 'Pending',
      uploading: 'Injecting',
      success: 'Success',
      failed: 'Failed',
      total: 'Total',
      totalSize: 'Total Size',
    },
    empty: {
      title: 'Queue Empty',
      desc: 'Drag data, click to select, or paste screenshots to start injection',
      action: 'Select Data',
    },
  },

  guide: {
    title: 'Injection Tips',
    multipleWays: {
      title: 'Multiple Injection Methods',
      desc: 'Supports drag and drop, click to select, or use {key}+V to paste clipboard data',
    },
    smartChunking: {
      title: 'Smart Chunking',
      desc: 'Large data units automatically chunked for transfer, supports resumable upload',
    },
    oneCopy: {
      title: 'One-Click Copy',
      desc: 'After batch injection, you can copy all data links with one click',
    },
  },

  advanced: {
    title: 'Advanced Config',
    storageDuration: 'Storage Duration',
    storageDurationDesc: 'Time data units are kept on server',
    watermark: 'Add Data Watermark',
    watermarkDesc: 'Add watermark identifier to injected data',
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
      desc: 'Access system to use permanent storage, custom watermarks, and other features',
      action: 'Access Now',
    },
  },

  toast: {
    folderSelected: 'Data directory selected: {path}',
    folderPreselected: 'Data directory preselected, can start data injection directly',
    sensitiveAutoDelete: 'Violating data will be automatically cleared, please ensure data units comply with protocol',
    sensitivePendingReview: 'Sensitive data will enter review queue, temporarily inaccessible during review',
    defaultHint: 'Please ensure injected data complies with platform protocol',
    guestDurationRequired: 'Guest injection must select storage duration',
    uploadCancelled: 'Data injection cancelled',
    noFilesInQueue: 'No data units to inject in queue',
    allLinksCopied: 'All links copied to clipboard',
    noLinksToCopy: 'No data links to copy yet',
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
    configure: 'Configure Data Watermark',
  },

  form: {
    storageDuration: {
      placeholder: 'Select storage duration',
    },
  },

  smartUpload: {
    error: {
      fileTypeNotSupported: 'Data type not supported: {name} ({type})',
      fileSizeExceeded: 'Data unit exceeded limit: {name} ({size})',
      chunkedUploadNotReady: 'Chunked transfer not initialized, please retry later',
      uploadFailed: 'Transfer Failed',
      analysisFailed: 'Data Analysis Failed',
      analysisTimeout: 'Data Analysis Timeout',
    },
    status: {
      preparing: 'Preparing Injection',
      uploading: 'Injecting...',
      uploadingWithProgress: 'Injecting {progress}% - {speed} - {remaining} remaining',
      completed: 'Injection Complete',
      failed: 'Injection Failed: {error}',
      paused: 'Paused',
      cancelled: 'Cancelled',
      retrying: 'Preparing Retry',
      userPaused: 'User Paused',
      userCancelled: 'User Cancelled',
      analyzing: 'Analyzing data...',
      analyzingProgress: 'Analyzing data {progress}%',
      checkingInstant: 'Checking instant upload...',
      checkingInstantProgress: 'Checking instant upload {progress}%',
      instantCompleted: 'Instant Upload Complete',
      foundResumable: 'Found unfinished transfer, can continue',
      analyzingCompleted: '✓ Analysis complete, ready to inject',
      analyzingFailed: 'Analysis failed: {error}',
      retryingUpload: 'Retrying injection...',
      resumingUpload: 'Resuming injection...',
      continuingUpload: 'Continuing injection...',
      restartingUpload: 'Restarting injection...',
      preparingSession: 'Preparing transfer session...',
      verifyingIntegrity: 'Verifying data integrity...',
      completingUpload: 'Completing injection...',
      chunkUploadFailed: 'Chunk {chunk} transfer failed: {error}',
      chunkRetrying: 'Chunk {chunk} retrying... ({attempt}/{maxRetry}) - retry in {delay}s',
    },
    time: {
      minutes: '{minutes}m {seconds}s',
      seconds: '{seconds}s',
    },
  },
}
