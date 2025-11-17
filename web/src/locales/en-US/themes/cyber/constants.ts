/**
 * Constants Text - Cyber Style
 */
export const constants = {
  fileTypes: {
    image: 'Data file',
    video: 'Video stream',
    document: 'Document data',
    archive: 'Archive package',
    audio: 'Audio stream',
    other: 'Unknown type',
  },

  accessLevels: {
    private: 'Private mode',
    public: 'Public access',
    protected: 'Protected',
  },

  sortOptions: {
    created_at: 'Creation time',
    updated_at: 'Update time',
    size: 'Data size',
    name: 'File name',
    views: 'Access count',
    downloads: 'Download count',
  },

  categorySortBy: {
    sort_order: 'By priority',
    name: 'By name',
    usage_count: 'By usage frequency',
    created_at: 'By creation time',
  },

  sortOrder: {
    asc: 'Ascending order',
    desc: 'Descending order',
  },

  categoryType: {
    all: 'All categories',
    hot: 'Hot tags',
    normal: 'Normal tags',
  },

  batchOperations: {
    delete: 'Data purge',
    move: 'Data transfer',
    update_access: 'Permission change',
    update_folder: 'Move to folder',
    download: 'Data extract',
  },

  accessControl: {
    modes: {
      ipWhitelist: 'IP whitelist',
      ipBlacklist: 'IP blacklist',
      domainWhitelist: 'Domain whitelist',
      domainBlacklist: 'Domain blacklist',
    },
    restrictionModes: {
      strict: 'Strict mode',
      moderate: 'Moderate mode',
      loose: 'Loose mode',
    },
    blockActions: {
      block: 'Complete block',
      redirect: 'Redirect',
      thumbnail: 'Return thumbnail',
      watermark: 'Add watermark',
    },
    warnings: {
      block: 'Will return 403 error to visitors, completely denying access to data resources, may affect user experience.',
      redirect: 'Will redirect visitors to your specified URL, ensure target URL is accessible.',
      thumbnail: 'Only shows thumbnail version to unauthorized visitors, protecting original data from theft.',
      watermark: 'Automatically adds watermark to files, can mark your ownership but will change original appearance.',
    },
    descriptions: {
      ipWhitelist: 'Whitelist mode: Only allow IPs in list to access',
      ipBlacklist: 'Blacklist mode: Block IPs in list from access',
      domainWhitelist: 'Whitelist mode: Only allow domains in list to reference',
      domainBlacklist: 'Blacklist mode: Block domains in list from reference',
      strictMode: 'Strict mode: Must completely match rules',
      moderateMode: 'Moderate mode: Partial match accepted',
      looseMode: 'Loose mode: Minimal restrictions',
    },
  },

  storageOptions: {
    unlimited: 'Unlimited',
    storageTypes: {
      local: 'Local storage',
      s3: 'S3 cloud',
      aliyun: 'Aliyun OSS',
      tencent: 'Tencent COS',
    },
    channelTypes: {
      oss: 'Aliyun OSS',
      cos: 'Tencent COS',
      rainyun: 'RainYun',
      local: 'Local storage',
      unknown: 'Unknown node',
    },
    storageClasses: {
      standard: 'Standard storage',
      standardMultiAz: 'Standard storage (Multi-AZ)',
      infrequent: 'Infrequent storage',
      infrequentMultiAz: 'Infrequent storage (Multi-AZ)',
      infrequentAccess: 'Infrequent access',
      intelligentTiering: 'Intelligent tiering',
      intelligentTieringMultiAz: 'Intelligent tiering (Multi-AZ)',
      archive: 'Archive storage',
      deepArchive: 'Deep archive',
      coldArchive: 'Cold archive',
      glacier: 'Glacier storage',
    },
  },

  filterOptions: {
    sortOptions: {
      newest: 'Latest injection',
      oldest: 'Earliest injection',
      name: 'Name',
      size: 'Data size',
      width: 'Width',
      height: 'Height',
      quality: 'Data quality',
      nsfw_score: 'Sensitivity score',
    },
    resolutionOptions: {
      low: 'Low resolution',
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
      all: 'All data types',
      jpeg: 'JPEG data',
      png: 'PNG data',
      gif: 'GIF data',
      webp: 'WebP data',
      svg: 'SVG data',
      bmp: 'BMP data',
      ico: 'ICO icon',
    },
  },

  category: {
    sortByOptions: {
      sort_order: 'By order',
      name: 'By name',
      usage_count: 'By usage count',
      created_at: 'By creation time',
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
      all: 'All actions',
      approve: 'Pass verification',
      reject: 'Data rejection',
    },
    sort: {
      newest: 'Latest injection',
      oldest: 'Earliest injection',
      size: 'Data size',
      nsfw: 'Risk score',
    },
    nsfwLevels: {
      safe: 'Safe level',
      low: 'Low risk',
      medium: 'Medium risk',
      high: 'High risk',
    },
    columns: {
      file: 'Data file',
      fileInfo: 'Data info',
      action: 'Review action',
      softDelete: 'Mark status',
      auditor: 'Auditor',
      auditTime: 'Review time',
      reason: 'Review reason',
      nsfw: 'Risk score',
      operations: 'System operations',
    },
    defaultReasons: {
      approve: 'Data verification passed',
      reject: 'Does not meet system specifications',
      batchApprove: 'Batch verification passed',
      batchReject: 'Batch data rejection',
    },
    messages: {
      approveSuccess: 'Data verification passed',
      rejectSuccess: 'Data rejected',
      batchApproveSuccess: 'Batch verification successful, processed {count} data items',
      batchRejectSuccess: 'Batch rejection successful, processed {count} data items',
      restoreSuccess: 'Data restored',
      hardDeleteSuccess: 'Data completely purged',
      batchRestoreSuccess: 'Batch restore complete: {success} successful',
      batchRestoreWithFail: 'Batch restore complete: {success} successful, {fail} failed',
      batchHardDeleteSuccess: 'Batch purge complete: {success} successful',
      batchHardDeleteWithFail: 'Batch purge complete: {success} successful, {fail} failed',
      loadError: 'Data load failed',
      operationError: 'Operation execution failed',
      noSelection: 'Please select data items to operate on',
      queueEmpty: 'Queue cleared',
      queueEmptyDesc: 'All data has completed review',
      logsEmpty: 'No review logs',
      logsEmptyDesc: 'No review records in system yet',
    },
    warnings: {
      hardDelete: 'Warning: This operation is irreversible',
      hardDeleteDesc: 'Data will be completely purged from server, cannot be recovered.',
      batchOperation: 'This operation will affect multiple data items',
    },
    descriptions: {
      approveEffects: ['Data will display in public database', 'Data status changes to "Verified"', 'Users can normally share and distribute data'],
      restoreEffects: ['Data will be restored from marked deletion status to normal status'],
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
        label: 'Grid layout',
        description: 'Neat grid display',
      },
      waterfall: {
        label: 'Waterfall',
        description: 'Auto-arrange by data height',
      },
      masonry: {
        label: 'Masonry layout',
        description: 'Left-aligned waterfall',
      },
      large: {
        label: 'Large mode',
        description: 'View larger thumbnails',
      },
    },
  },

  navigation: {
    admin: {
      dashboard: 'Control panel',
      files: 'Data management',
      tags: 'Tag system',
      categories: 'Category system',
      shares: 'Share management',
      ai: 'Intelligence system',
      contentReview: 'Content review',
      tagging: 'Auto tagging',
      vectors: 'Vector system',
      users: 'User management',
      channels: 'Channel management',
      construction: 'Site construction',
      announcements: 'Announcement system',
      settings: 'System settings',
    },
    docs: {
      overview: 'API overview',
      authentication: 'Authentication mechanism',
      upload: 'Upload API',
      limits: 'Limits description',
      examples: 'Code examples',
      tester: 'API tester',
      faq: 'FAQ',
    },
    main: {
      home: 'Home',
      explore: 'Database',
      upload: 'Upload',
      random: 'Random',
      docs: 'API',
    },
    user: {
      dashboard: 'Control panel',
      myFiles: 'My data',
      folders: 'Folders',
      tagManage: 'Tag management',
      categoryManage: 'Category management',
      shares: 'My shares',
      settings: 'Settings',
    },
    settings: {
      profile: 'Personal info',
      security: 'Security settings',
      apikey: 'API keys',
      storage: 'Storage config',
    },
    adminSettings: {
      website: 'Website settings',
      upload: 'Upload settings',
      security: 'Security settings',
      registration: 'Registration settings',
      mail: 'Mail settings',
      ai: 'AI settings',
      vector: 'Vector settings',
    },
  },

  api: {
    errors: {
      badRequest: 'Request parameter exception',
      unauthorized: 'Unauthorized access',
      forbidden: 'Access denied',
      notFound: 'Resource not found',
      methodNotAllowed: 'Method not allowed',
      conflict: 'Data conflict',
      payloadTooLarge: 'Data package too large',
      unprocessableEntity: 'Data validation failed',
      tooManyRequests: 'Request frequency too high',
      unavailableForLegalReasons: 'IP address banned',
      internalServerError: 'System internal error',
      badGateway: 'Gateway exception',
      serviceUnavailable: 'Service unavailable',
      gatewayTimeout: 'Gateway timeout',
    },
  },

  systemOptions: {
    sensitiveContent: {
      markOnly: 'Mark only',
      pendingReview: 'System review',
      autoDelete: 'Auto purge',
    },
    sizeUnits: {
      bytes: 'Bytes (B)',
      kilobytes: 'Kilobytes (KB)',
      megabytes: 'Megabytes (MB)',
      gigabytes: 'Gigabytes (GB)',
      terabytes: 'Terabytes (TB)',
    },
    userRoles: {
      normal: 'Normal user',
      advanced: 'Advanced user',
      admin: 'System administrator',
      superAdmin: 'Super administrator',
    },
    userStatus: {
      normal: 'Normal',
      disabled: 'Disabled',
      pending: 'Pending review',
    },
    imageStatus: {
      normal: 'Normal',
      hidden: 'Hidden',
      pending: 'Pending review',
      deleted: 'Deleted',
    },
    shareStatus: {
      normal: 'Normal',
      disabled: 'Disabled',
      expired: 'Expired',
    },
  },
}
