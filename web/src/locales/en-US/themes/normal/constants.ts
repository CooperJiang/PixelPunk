/**
 * Constants
 */
export const constants = {
  fileTypes: {
    image: 'File',
    video: 'Video',
    document: 'Document',
    archive: 'Archive',
    audio: 'Audio',
    other: 'Other',
  },

  accessLevels: {
    private: 'Private',
    public: 'Public',
    protected: 'Protected',
  },

  sortOptions: {
    created_at: 'Creation Time',
    updated_at: 'Update Time',
    size: 'File Size',
    name: 'File Name',
    views: 'Views',
    downloads: 'Downloads',
  },

  categorySortBy: {
    sort_order: 'By Sort Order',
    name: 'By Name',
    usage_count: 'By Usage Count',
    created_at: 'By Creation Time',
  },

  sortOrder: {
    asc: 'Ascending',
    desc: 'Descending',
  },

  categoryType: {
    all: 'All',
    hot: 'Hot',
    normal: 'Normal',
  },

  batchOperations: {
    delete: 'Delete',
    move: 'Move',
    update_access: 'Update Permission',
    update_folder: 'Move to Folder',
    download: 'Download',
  },

  accessControl: {
    modes: {
      ipWhitelist: 'IP Whitelist',
      ipBlacklist: 'IP Blacklist',
      domainWhitelist: 'Domain Whitelist',
      domainBlacklist: 'Domain Blacklist',
    },
    restrictionModes: {
      strict: 'Strict Mode',
      moderate: 'Moderate Mode',
      loose: 'Loose Mode',
    },
    blockActions: {
      block: 'Block Completely',
      redirect: 'Redirect',
      thumbnail: 'Return Thumbnail',
      watermark: 'Add Watermark',
    },
    warnings: {
      block:
        'Will return a 403 error to visitors, completely denying access to file resources, which may affect user experience.',
      redirect: 'Redirects visitors to the URL you specify. Please ensure the target URL is accessible.',
      thumbnail: 'Only shows thumbnail versions of files to unauthorized visitors, protecting original images from theft.',
      watermark:
        'Automatically adds watermarks to files, which can identify your ownership but will change the appearance of the original image.',
    },
    descriptions: {
      ipWhitelist: 'Whitelist mode: Only allow IPs in the list to access',
      ipBlacklist: 'Blacklist mode: Block IPs in the list from accessing',
      domainWhitelist: 'Whitelist mode: Only allow domains in the list to reference',
      domainBlacklist: 'Blacklist mode: Block domains in the list from referencing',
      strictMode: 'Strict mode: Must fully match rules',
      moderateMode: 'Moderate mode: Partial match is sufficient',
      looseMode: 'Loose mode: Minimal restrictions',
    },
  },

  storageOptions: {
    unlimited: 'Unlimited',
    storageTypes: {
      local: 'Local Storage',
      s3: 'S3 Storage',
      aliyun: 'Alibaba Cloud OSS',
      tencent: 'Tencent Cloud COS',
    },
    channelTypes: {
      oss: 'Alibaba Cloud OSS',
      cos: 'Tencent Cloud COS',
      rainyun: 'RainYun',
      local: 'Local Storage',
      unknown: 'Unknown',
    },
    storageClasses: {
      standard: 'Standard Storage',
      standardMultiAz: 'Standard Storage (Multi-AZ)',
      infrequent: 'Infrequent Storage',
      infrequentMultiAz: 'Infrequent Storage (Multi-AZ)',
      infrequentAccess: 'Infrequent Access',
      intelligentTiering: 'Intelligent Tiering Storage',
      intelligentTieringMultiAz: 'Intelligent Tiering Storage (Multi-AZ)',
      archive: 'Archive Storage',
      deepArchive: 'Deep Archive Storage',
      coldArchive: 'Cold Archive Storage',
      glacier: 'Glacier Storage',
    },
  },

  filterOptions: {
    sortOptions: {
      newest: 'Newest Upload',
      oldest: 'Oldest Upload',
      name: 'Name',
      size: 'File Size',
      width: 'Width',
      height: 'Height',
      quality: 'File Quality',
      nsfw_score: 'NSFW Score',
    },
    resolutionOptions: {
      low: 'Low Resolution',
      '720p': '720p',
      '1080p': '1080p',
      '2k': '2k',
      '4k': '4K',
      '8k': '8K',
    },
    imageFormatOptions: {
      jpg: 'JPG',
      jpeg: 'JPEG',
      png: 'PNG',
      gif: 'GIF',
      webp: 'WebP',
      bmp: 'BMP',
      svg: 'SVG',
      ico: 'ICO',
    },
    fileTypeOptions: {
      all: 'All File Types',
      jpeg: 'JPEG Files',
      png: 'PNG Files',
      gif: 'GIF Files',
      webp: 'WebP Files',
      svg: 'SVG Files',
      bmp: 'BMP Files',
      ico: 'ICO Icons',
    },
  },

  category: {
    sortByOptions: {
      sort_order: 'By Sort Order',
      name: 'By Name',
      usage_count: 'By Usage Count',
      created_at: 'By Creation Time',
    },
    sortOrderOptions: {
      asc: 'Ascending',
      desc: 'Descending',
    },
    typeOptions: {
      all: 'All',
      popular: 'Popular',
      normal: 'Normal',
    },
  },

  contentReview: {
    actions: {
      all: 'All Actions',
      approve: 'Approve',
      reject: 'Reject',
    },
    sort: {
      newest: 'Newest Upload',
      oldest: 'Oldest Upload',
      size: 'File Size',
      nsfw: 'NSFW Score',
    },
    nsfwLevels: {
      safe: 'Safe',
      low: 'Low',
      medium: 'Medium',
      high: 'High',
    },
    columns: {
      file: 'File',
      fileInfo: 'File Info',
      action: 'Review Action',
      softDelete: 'Delete Status',
      auditor: 'Auditor',
      auditTime: 'Audit Time',
      reason: 'Review Reason',
      nsfw: 'NSFW Score',
      operations: 'Operations',
    },
    defaultReasons: {
      approve: 'File review approved',
      reject: 'Does not meet platform standards',
      batchApprove: 'Batch Approve',
      batchReject: 'Batch Reject',
    },
    messages: {
      approveSuccess: 'File approved',
      rejectSuccess: 'File rejected',
      batchApproveSuccess: 'Batch approval successful, processed {count} resources',
      batchRejectSuccess: 'Batch rejection successful, processed {count} resources',
      restoreSuccess: 'File restored',
      hardDeleteSuccess: 'File permanently deleted',
      batchRestoreSuccess: 'Batch restore completed: {success} successful',
      batchRestoreWithFail: 'Batch restore completed: {success} successful, {fail} failed',
      batchHardDeleteSuccess: 'Batch delete completed: {success} successful',
      batchHardDeleteWithFail: 'Batch delete completed: {success} successful, {fail} failed',
      loadError: 'Failed to load data',
      operationError: 'Operation failed',
      noSelection: 'Please select items to operate',
      queueEmpty: 'Queue is empty',
      queueEmptyDesc: 'All files have been reviewed',
      logsEmpty: 'No review records',
      logsEmptyDesc: 'No review operation records yet',
    },
    warnings: {
      hardDelete: 'Warning: This operation is irreversible',
      hardDeleteDesc: 'Files will be permanently deleted from the server and cannot be recovered.',
      batchOperation: 'This operation will affect multiple items',
    },
    descriptions: {
      approveEffects: [
        'Files will be visible in public gallery',
        'File status changed to "Reviewed"',
        'Users can share and distribute normally',
      ],
      restoreEffects: ['Files will be restored from soft-deleted status to normal status'],
    },
  },

  share: {
    sortOptions: {
      date: 'Date',
      name: 'Name',
      size: 'Size',
    },
    layoutOptions: {
      grid: {
        label: 'Grid Layout',
        description: 'Neat grid display',
      },
      waterfall: {
        label: 'Waterfall',
        description: 'Automatically arranged by file height',
      },
      masonry: {
        label: 'Masonry Layout',
        description: 'Left-aligned waterfall',
      },
      large: {
        label: 'Large Image Mode',
        description: 'View larger thumbnails',
      },
    },
  },

  navigation: {
    admin: {
      dashboard: 'Dashboard',
      files: 'File Management',
      tags: 'Tag Management',
      categories: 'Category Management',
      shares: 'Share Management',
      ai: 'AI Management',
      contentReview: 'Content Review',
      tagging: 'Tagging Management',
      vectors: 'Vector Management',
      users: 'User Management',
      channels: 'Channel Management',
      construction: 'Website Construction',
      announcements: 'Announcement Management',
      settings: 'Global Settings',
    },
    docs: {
      overview: 'API Overview',
      authentication: 'Authentication',
      upload: 'Upload API',
      limits: 'Limits',
      examples: 'Code Examples',
      tester: 'API Tester',
      faq: 'FAQ',
    },
    main: {
      home: 'Home',
      explore: 'Gallery',
      upload: 'Upload',
      random: 'Random',
      docs: 'API',
    },
    user: {
      dashboard: 'Dashboard',
      myFiles: 'My Files',
      folders: 'Folders',
      tagManage: 'Tag Management',
      categoryManage: 'Category Management',
      shares: 'My Shares',
      settings: 'Settings',
    },
    settings: {
      profile: 'Profile',
      security: 'Security Settings',
      apikey: 'API Key',
      storage: 'Storage Configuration',
    },
    adminSettings: {
      website: 'Site Settings',
      upload: 'Upload Settings',
      security: 'Security Settings',
      registration: 'Registration Settings',
      mail: 'Mail Settings',
      ai: 'AI Settings',
      vector: 'Vector Settings',
    },
  },

  api: {
    errors: {
      badRequest: 'Bad Request',
      unauthorized: 'Unauthorized',
      forbidden: 'Forbidden',
      notFound: 'Resource Not Found',
      methodNotAllowed: 'Method Not Allowed',
      conflict: 'Resource Conflict',
      payloadTooLarge: 'Payload Too Large',
      unprocessableEntity: 'Validation Failed',
      tooManyRequests: 'Too Many Requests',
      unavailableForLegalReasons: 'IP Address Disabled',
      internalServerError: 'Internal Server Error',
      badGateway: 'Bad Gateway',
      serviceUnavailable: 'Service Unavailable',
      gatewayTimeout: 'Gateway Timeout',
    },
  },

  systemOptions: {
    sensitiveContent: {
      markOnly: 'Mark Only',
      pendingReview: 'Admin Review',
      autoDelete: 'Auto Delete',
    },
    sizeUnits: {
      bytes: 'Bytes (B)',
      kilobytes: 'Kilobytes (KB)',
      megabytes: 'Megabytes (MB)',
      gigabytes: 'Gigabytes (GB)',
      terabytes: 'Terabytes (TB)',
    },
    userRoles: {
      normal: 'Normal User',
      advanced: 'Advanced User',
      admin: 'Administrator',
      superAdmin: 'Super Administrator',
    },
    userStatus: {
      normal: 'Normal',
      disabled: 'Disabled',
      pending: 'Pending',
    },
    imageStatus: {
      normal: 'Normal',
      hidden: 'Hidden',
      pending: 'Pending',
      deleted: 'Deleted',
    },
    shareStatus: {
      normal: 'Normal',
      disabled: 'Disabled',
      expired: 'Expired',
    },
  },
}
