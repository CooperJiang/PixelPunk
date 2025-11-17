/**
 * Gallery Browsing Text - Cyber Style
 */
export const gallery = {
  pages: {
    gallery: 'Database',
    upload: 'Data Upload Center',
    login: 'Identity Verification',
    register: 'System Registration',
    settings: 'System Config',
    profile: 'User Profile',
    folders: 'Data Container',
    dashboard: 'Control Panel',
    admin: 'Admin Terminal',
  },
  viewer: {
    previous: 'Previous Frame',
    next: 'Next Frame',
    fitMode: 'Fit Viewport',
    fillMode: 'Fill Viewport',
    fullscreen: 'Immersive Mode',
    exitFullscreen: 'Exit Immersive',
    originalSize: 'Original Resolution',
    imageInfo: 'Data Properties',
    similarImages: 'Similar Data',
  },
  waterfallLayout: {
    loadFailed: 'Data Corrupted',
    collapse: 'Collapse',
  },
  file: {
    loadFailed: 'Data load failed',
    nsfwWarning: 'Restricted Content',
  },
  enhancedFilePreview: {
    defaultImageName: 'Data File',
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
    title: 'Data Details',
    untitled: 'Unnamed Data',

    basicInfo: {
      title: 'Basic Info',
      fileId: 'Data ID',
      originalName: 'Original Name',
      fileFormat: 'Data Format',
      fileSize: 'Data Size',
      imageDimensions: 'Data Dimensions',
      visibility: 'Visibility',
      uploader: 'Uploader',
      uploadTime: 'Upload Time',
      updateTime: 'Update Time',
      fileHash: 'Data Hash',
      viewCount: 'View Count',
      status: 'Status',
      times: 'times',
      count: 'items',
      unknownResolution: 'Unknown Resolution',
      none: 'None',
    },

    statusLabels: {
      duplicate: 'Duplicate Data',
      recommended: 'Recommended',
      nsfw: 'NSFW',
    },

    visibility: {
      public: 'Public Protocol',
      private: 'Private Archive',
      protected: 'Security Shield',
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
      high: 'High Risk',
      medium: 'Medium Risk',
      low: 'Low Risk',
    },

    visualAnalysis: {
      title: 'Visual Analysis',
      dominantColor: 'Dominant Color',
      imageDimensions: 'Data Dimensions',
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
      deleteImage: 'Clear Data',
      cancelRecommend: 'Cancel Recommendation',
      setRecommend: 'Set as Recommended',
      close: 'Disconnect',
    },
  },
  fileExpiryTag: {
    statuses: {
      expired: 'Data Expired',
      expiringSoon: 'Expiring Soon',
      active: 'Limited Time Storage',
    },
    labels: {
      storageDuration: 'Storage Duration',
      expiresAt: 'Expiration Time',
      expiredAtPrefix: 'Expired on',
      expiredAtSuffix: '',
    },
    units: {
      day: 'day',
      hour: 'hour',
      minute: 'minute',
    },
  },
  fileLoading: {
    loading: 'Transferring data',
  },
  fileActions: {
    copyLink: {
      protectedNotSupported: 'Protected data does not support external transfer, please view within terminal',
      publicLinkCopied: 'Public link copied to clipboard',
      privateLinkCopied: 'Private link copied to clipboard',
      copyFailed: 'Link copy failed',
    },
    delete: {
      success: 'Data deleted successfully',
      failed: 'Delete failed, please retry',
    },
    accessLevel: {
      public: 'Public Protocol',
      private: 'Private Archive',
      protected: 'Security Shield',
      switched: 'Data set to {level}',
      switchFailed: 'Switch failed, please retry',
    },
    batch: {
      noPublicFiles: 'No public resources to copy in selected data',
      filteredFiles: 'Filtered {count} private/protected data',
      linksCopied: 'Copied {count} data links',
      copyFailed: 'Batch copy failed',
      deleteSuccess: 'Successfully deleted {count} resources',
      deletePartial: 'Successfully deleted {successCount}/{totalCount} resources',
      deleteFailed: 'Delete failed, please retry',
    },
  },
}
