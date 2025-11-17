/**
 * 文件夹管理页面文案
 */
export const folders = {
  rootFolder: '根目录',
  emptyState: {
    title: '当前文件夹为空',
    description: '您可以上传资源或创建新的文件夹开始使用',
    createButton: '创建文件夹',
  },
  breadcrumb: {
    label: '文件夹导航',
  },
  share: {
    createSuccess: '创建分享成功',
    createSuccessWithCopy: '分享创建成功，链接已复制到剪贴板',
    manualCopy: '请手动复制分享链接',
    description: '分享描述',
  },
  shareSelectTip: {
    message: '请选择要分享的文件夹或图片',
  },
  folderGrid: {
    sectionTitle: '文件夹',
    dragTip: '拖拽排序',
    fileCount: '{count} 张',
    toast: {
      moveFailed: '移动文件夹失败',
      moveToRootSuccess: '已将文件夹移动到根目录',
      moveSuccess: '已将文件夹移动到 "{name}"',
    },
    contextMenu: {
      root: '根目录',
      moveTo: '移动到...',
      rename: '重命名/编辑',
      delete: '删除',
    },
    actions: {
      selectAll: '全选',
      deselectAll: '取消全选',
      selectAllTitle: '全选文件夹',
      deselectAllTitle: '取消全选文件夹',
      invert: '反选',
      invertTitle: '反选文件夹',
      deleteTitle: '删除文件夹',
    },
    visibility: {
      toggleTitle: '点击切换可见性',
      setPublic: '设为公开',
      setPrivate: '设为私密',
    },
  },
  createFolderDialog: {
    title: {
      create: '新建文件夹',
      edit: '编辑文件夹',
    },
    form: {
      name: {
        label: '文件夹名称',
        placeholder: '请输入文件夹名称',
        required: '请输入文件夹名称',
        maxLength: '文件夹名称不能超过30个字符',
      },
      description: {
        label: '描述（可选）',
        placeholder: '添加文件夹描述',
      },
      permission: {
        label: '权限设置',
        public: {
          name: '公开',
          desc: '他人进入我的主页可观看内容',
        },
        private: {
          name: '私密',
          desc: '对他人不展示此文件夹',
        },
      },
    },
    actions: {
      create: '创建',
      save: '保存',
    },
    toast: {
      createSuccess: '文件夹创建成功',
      updateSuccess: '文件夹更新成功',
      operationFailed: '文件夹操作失败',
    },
  },
  fileGrid: {
    sectionTitle: '文件',
    dragTip: '拖拽排序',
    selectedCount: '已选择 {count} 张',
    badge: {
      duplicate: '重复',
    },
    actions: {
      aiInfo: 'AI信息',
      preview: '预览',
      copyLink: '复制链接',
      download: '下载',
      delete: '删除',
      batchMode: '批量操作',
      batchDelete: '批量删除',
      cancelBatch: '取消批量操作',
      selectAll: '全选',
      deselectAll: '取消全选',
      selectAllTitle: '全选文件',
      deselectAllTitle: '取消全选文件',
      invert: '反选',
      invertTitle: '反选文件',
      gridView: '网格视图',
      listView: '列表视图',
    },
    meta: {
      size: '文件大小',
      dimensions: '文件尺寸',
      duplicateFile: '重复文件',
    },
    listView: {
      columns: {
        preview: '预览',
        filename: '文件名',
        size: '大小',
        dimensions: '尺寸',
        date: '上传时间',
        actions: '操作',
      },
      empty: '暂无文件',
    },
  },
  pageHeader: {
    title: '我的文件夹',
    breadcrumb: {
      root: '根目录',
      navigation: '文件夹导航',
    },
    upload: {
      default: '上传文件',
      toFolder: '上传到 "{name}"',
    },
    actions: {
      newFolder: '新建文件夹',
    },
  },
  createShareDialog: {
    title: '创建分享',
    sections: {
      basic: '基本信息',
      access: '访问控制',
      advanced: '高级设置',
      content: '分享内容',
    },
    form: {
      name: {
        label: '分享名称',
        placeholder: '请输入分享名称（可选）',
      },
      description: {
        placeholder: '请输入分享描述（可选）',
      },
      password: {
        label: '访问密码',
        placeholder: '请输入访问密码（可选）',
      },
      expiredDays: {
        label: '过期天数',
        placeholder: '0表示永不过期',
      },
      maxViews: {
        label: '最大访问次数',
        placeholder: '0表示不限制次数',
      },
      collectVisitor: {
        label: '收集访客信息',
        tip: '访客浏览时会弹出信息填写窗口，非必填项',
      },
      notification: {
        label: '访问通知',
        tip: '开启后将在访问量达到指定次数时收到邮件提醒',
      },
      notificationThreshold: {
        label: '通知阈值',
        placeholder: '设置访问次数达到多少时通知',
        tip: '当访问次数达到该值时，系统将发送邮件通知',
      },
    },
    emptyItems: '请选择要分享的文件夹或文件',
    cancel: '取消',
    submit: '创建分享',
    toast: {
      createSuccess: '创建分享成功',
      linkCopied: '分享创建成功，链接已复制到剪贴板',
      copyFailed: '复制到剪贴板失败',
      copyManually: '请手动复制分享链接',
    },
  },
  deleteConfirmDialogs: {
    title: {
      folder: '删除文件夹',
      file: '删除文件',
      batchFile: '批量删除文件',
    },
    message: {
      folder: '确定要删除文件夹 "{name}" 吗？',
      file: '确定要删除文件 "{name}" 吗？',
      batchFile: '确定要批量删除 {count} 个资源吗？',
    },
    warning: '此操作不可撤销，请谨慎操作。',
    confirmText: '确认删除',
  },
  fileContextMenu: {
    toast: {
      noFileSelected: '没有选中的文件',
      moveFailed: '移动文件失败',
      moveToRootSuccess: '已将文件移动到根目录',
      moveToFolderSuccess: '已将文件移动到 "{name}"',
      loadFoldersFailed: '加载文件夹列表失败',
      loadSubfoldersFailed: '加载子文件夹失败',
      noValidImageLink: '无有效原图链接',
      noThumbLink: '没有可复制的缩略图链接',
      thumbLinkCopied: '已复制缩略图链接',
      copyFailed: '复制失败',
      originalLinkCopied: '已复制原图链接',
      noOriginalLink: '没有可复制的原图链接',
      shortLinkCopied: '已复制短链',
      noShortLink: '没有可复制的短链',
      fileDeleted: '文件已删除',
      deleteFailed: '删除文件失败',
    },
    menu: {
      root: '根目录',
      preview: '预览',
      openNew: '新窗口打开',
      copyLink: '复制链接',
      copyThumb: '复制缩略图',
      copyOriginal: '复制原图',
      copyShort: '复制短链',
      download: '下载',
      moveTo: '移动到...',
      delete: '删除',
    },
    confirm: {
      deleteFile: '确定要删除文件 "{name}" 吗？此操作不可撤销。',
    },
  },
  folderManagement: {
    deleteSuccess: '文件夹 "{name}" 已删除',
    visibility: {
      public: '公开',
      private: '私密',
      switched: '文件夹 "{name}" 已设为{level}',
    },
  },
  dragSort: {
    folderSortUpdated: '文件夹排序已更新',
    fileSortUpdated: '文件排序已更新',
  },
  fileManagement: {
    copyLink: {
      protectedNotSupported: '受保护文件不支持外部分享，请在应用内查看',
      publicLinkCopied: '公开链接已复制到剪贴板',
      privateLinkCopied: '私有链接已复制到剪贴板',
    },
    download: {
      success: '文件 "{name}" 下载完成',
      cancelled: '下载被取消或失败',
    },
    delete: {
      success: '文件 "{name}" 已删除',
      selectFirst: '请先选择要删除的文件',
      batchSuccess: '成功删除 {successCount} 个文件',
      batchPartial: '成功删除 {successCount} 个文件，{failCount} 个删除失败',
      batchFailed: '删除失败，共 {failCount} 个文件删除失败',
    },
    accessLevel: {
      public: '公开',
      private: '私有',
      protected: '受保护',
      switched: '文件 "{name}" 已设为{level}',
    },
  },
}
