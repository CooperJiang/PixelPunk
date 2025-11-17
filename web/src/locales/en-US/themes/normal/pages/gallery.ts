/**
 * Gallery browsing texts
 */
export const gallery = {
  pages: {
    gallery: 'Gallery',
    upload: 'Upload Files',
    login: 'Login',
    register: 'Register',
    settings: 'Settings',
    profile: 'Profile',
    folders: 'Folders',
    dashboard: 'Dashboard',
    admin: 'Admin Panel',
  },
  viewer: {
    previous: 'Previous',
    next: 'Next',
    fitMode: 'Fit Mode',
    fillMode: 'Fill Mode',
    fullscreen: 'Fullscreen',
    exitFullscreen: 'Exit Fullscreen',
    originalSize: 'Original Size',
    imageInfo: 'File Info',
    similarImages: 'Similar Files',
  },
  waterfallLayout: {
    loadFailed: 'Load Failed',
    collapse: 'Collapse',
  },
  file: {
    loadFailed: 'Load Failed',
    nsfwWarning: 'Inappropriate Content',
  },
  enhancedFilePreview: {
    defaultImageName: 'File',
    fillMode: 'Fill Mode',
    fitMode: 'Fit Mode',
    switchToFit: 'Switch to Fit',
    switchToFill: 'Switch to Fill',
    wheelZoom: 'Wheel Zoom',
    dragMove: 'Drag Move',
    spaceKey: 'Space',
    exitPreview: 'Exit Preview',
    exitFullscreen: 'Exit Fullscreen',
    exitPreviewEsc: 'Exit Preview (ESC)',
    exitFullscreenEsc: 'Exit Fullscreen (ESC)',
    enterFullscreen: 'Enter Browser Fullscreen',
    fullscreen: 'Fullscreen',
    preview: 'Preview',
    mode: 'Mode',
  },
  fileDetailModal: {
    title: 'File Details',
    untitled: 'Untitled',

    basicInfo: {
      title: 'Basic Information',
      fileId: 'File ID',
      originalName: 'Original Name',
      fileFormat: 'File Format',
      fileSize: 'File Size',
      imageDimensions: 'File Dimensions',
      visibility: 'Visibility',
      uploader: 'Uploader',
      uploadTime: 'Upload Time',
      updateTime: 'Update Time',
      fileHash: 'File Hash',
      viewCount: 'View Count',
      status: 'Status',
      times: 'times',
      count: 'items',
      unknownResolution: 'Unknown Resolution',
      none: 'None',
    },

    statusLabels: {
      duplicate: 'Duplicate',
      recommended: 'Recommended',
      nsfw: 'NSFW',
    },

    visibility: {
      public: 'Public',
      private: 'Private',
      protected: 'Protected',
      link: 'Link Visible',
    },

    qualityLevels: {
      high: 'High',
      medium: 'Medium',
      low: 'Low',
      excellent: 'Excellent',
      normal: 'Normal',
      poor: 'Poor',
    },

    contentSafety: {
      title: 'Content Safety',
      nsfwStatus: 'NSFW Status',
      inappropriate: 'Inappropriate',
      safe: 'Safe',
      nsfwScore: 'NSFW Score',
      high: 'High',
      medium: 'Medium',
      low: 'Low',
    },

    visualAnalysis: {
      title: 'Visual Analysis',
      dominantColor: 'Dominant Color',
      imageDimensions: 'File Dimensions',
      aspectRatio: 'Aspect Ratio',
      resolution: 'Resolution',
      colorPalette: 'Color Palette',
      composition: 'Composition',
      objectsCount: 'Object Count',
    },

    ai: {
      tags: 'AI Tags',
      description: 'AI Description',
    },

    actions: {
      deleteImage: 'Delete File',
      cancelRecommend: 'Cancel Recommendation',
      setRecommend: 'Set as Recommended',
      close: 'Close',
    },
  },
  fileExpiryTag: {
    statuses: {
      expired: 'Expired',
      expiringSoon: 'Expiring Soon',
      active: 'Time Limited',
    },
    labels: {
      storageDuration: 'Storage Duration',
      expiresAt: 'Expiration Time',
      expiredAtPrefix: 'Expired on',
      expiredAtSuffix: 'expired',
    },
    units: {
      day: 'days',
      hour: 'hours',
      minute: 'minutes',
    },
  },
  fileLoading: {
    loading: 'Loading',
  },
  fileActions: {
    copyLink: {
      protectedNotSupported: 'Protected files do not support external sharing, please view within the application',
      publicLinkCopied: 'Public link copied to clipboard',
      privateLinkCopied: 'Private link copied to clipboard',
      copyFailed: 'Failed to copy link',
    },
    delete: {
      success: 'File deleted successfully',
      failed: 'Delete failed, please try again',
    },
    accessLevel: {
      public: 'Public',
      private: 'Private',
      protected: 'Protected',
      switched: 'File set to {level}',
      switchFailed: 'Switch failed, please try again',
    },
    batch: {
      noPublicFiles: 'No public files available to copy in selected files',
      filteredFiles: 'Filtered out {count} private/protected files',
      linksCopied: 'Copied {count} file links',
      copyFailed: 'Batch copy failed',
      deleteSuccess: 'Successfully deleted {count} resources',
      deletePartial: 'Successfully deleted {successCount}/{totalCount} resources',
      deleteFailed: 'Delete failed, please try again',
    },
  },
}
