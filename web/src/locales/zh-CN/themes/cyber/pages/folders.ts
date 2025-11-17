/**
 * 文件夹管理页面文案 - 赛博风格
 */
export const folders = {
  rootFolder: '根节点',
  emptyState: {
    title: '当前数据目录为空',
    description: '您可以上传数据资源或创建新的数据目录开始使用',
    createButton: '创建数据目录',
  },
  breadcrumb: {
    label: '数据目录导航',
  },
  share: {
    createSuccess: '创建数据共享成功',
    createSuccessWithCopy: '数据共享创建成功，链接已复制到剪贴板',
    manualCopy: '请手动复制共享链接',
    description: '共享描述',
  },
  shareSelectTip: {
    message: '请选择要共享的数据目录或数据单元',
  },
  folderGrid: {
    sectionTitle: '数据目录',
    dragTip: '拖拽排序',
    fileCount: '{count} 个',
    toast: {
      moveFailed: '移动数据目录失败',
      moveToRootSuccess: '已将数据目录移动到根节点',
      moveSuccess: '已将数据目录移动到 "{name}"',
    },
    contextMenu: {
      root: '根节点',
      moveTo: '移动到...',
      rename: '重命名/编辑',
      delete: '删除',
    },
    actions: {
      selectAll: '全选',
      deselectAll: '取消全选',
      selectAllTitle: '全选数据目录',
      deselectAllTitle: '取消全选数据目录',
      invert: '反选',
      invertTitle: '反选数据目录',
      deleteTitle: '删除数据目录',
    },
    visibility: {
      toggleTitle: '点击切换可见性',
      setPublic: '设为公开',
      setPrivate: '设为私密',
    },
  },
  createFolderDialog: {
    title: {
      create: '新建数据目录',
      edit: '编辑数据目录',
    },
    form: {
      name: {
        label: '目录名称',
        placeholder: '请输入数据目录名称',
        required: '请输入数据目录名称',
        maxLength: '数据目录名称不能超过30个字符',
      },
      description: {
        label: '描述（可选）',
        placeholder: '添加数据目录描述',
      },
      permission: {
        label: '权限设置',
        public: {
          name: '公开',
          desc: '他人进入我的主页可观看内容',
        },
        private: {
          name: '私密',
          desc: '对他人不展示此数据目录',
        },
      },
    },
    actions: {
      create: '创建',
      save: '保存',
    },
    toast: {
      createSuccess: '数据目录创建成功',
      updateSuccess: '数据目录更新成功',
      operationFailed: '数据目录操作失败',
    },
  },
  fileGrid: {
    sectionTitle: '数据文件',
    dragTip: '拖拽排序',
    selectedCount: '已选择 {count} 个',
    badge: {
      duplicate: '重复',
    },
    actions: {
      aiInfo: 'AI分析',
      preview: '预览',
      copyLink: '复制链接',
      download: '下载',
      delete: '删除',
      batchMode: '批量协议操作',
      batchDelete: '批量删除',
      cancelBatch: '取消批量协议操作',
      selectAll: '全选',
      deselectAll: '取消全选',
      selectAllTitle: '全选数据单元',
      deselectAllTitle: '取消全选数据单元',
      invert: '反选',
      invertTitle: '反选数据单元',
      gridView: '网格视图',
      listView: '列表视图',
    },
    meta: {
      size: '数据单元大小',
      dimensions: '数据单元尺寸',
      duplicateFile: '重复数据单元',
    },
    listView: {
      columns: {
        preview: '预览',
        filename: '数据单元名称',
        size: '大小',
        dimensions: '尺寸',
        date: '注入时间',
        actions: '操作',
      },
      empty: '暂无数据单元',
    },
  },
  pageHeader: {
    title: '我的数据目录',
    breadcrumb: {
      root: '根节点',
      navigation: '数据目录导航',
    },
    upload: {
      default: '注入数据单元',
      toFolder: '注入到 "{name}"',
    },
    actions: {
      newFolder: '新建数据目录',
    },
  },
  createShareDialog: {
    title: '创建共享链接',
    sections: {
      basic: '基本信息节点',
      access: '访问控制节点',
      advanced: '高级设置节点',
      content: '共享内容节点',
    },
    form: {
      name: {
        label: '共享名称',
        placeholder: '请输入共享名称（可选）',
      },
      description: {
        placeholder: '请输入共享描述（可选）',
      },
      password: {
        label: '访问密码',
        placeholder: '请输入访问密码（可选）',
      },
      expiredDays: {
        label: '过期天数节点',
        placeholder: '0表示永不过期',
      },
      maxViews: {
        label: '最大访问次数节点',
        placeholder: '0表示不限制次数',
      },
      collectVisitor: {
        label: '收集访客信息节点',
        tip: '访客浏览时会弹出信息填写窗口节点，非必填项',
      },
      notification: {
        label: '访问通知节点',
        tip: '开启后将在访问量达到指定次数时收到邮件提醒节点',
      },
      notificationThreshold: {
        label: '通知阈值节点',
        placeholder: '设置访问次数达到多少时通知',
        tip: '当访问次数达到该值时，系统将发送邮件通知节点',
      },
    },
    emptyItems: '请选择要共享的数据目录或数据单元',
    cancel: '取消',
    submit: '创建共享链接',
    toast: {
      createSuccess: '创建共享链接成功',
      linkCopied: '共享链接创建成功，链接已复制到剪贴板',
      copyFailed: '复制到剪贴板失败',
      copyManually: '请手动复制共享链接',
    },
  },
  deleteConfirmDialogs: {
    title: {
      folder: '删除数据目录',
      file: '删除数据单元',
      batchFile: '批量删除数据单元',
    },
    message: {
      folder: '确定要删除数据目录 "{name}" 吗？',
      file: '确定要删除数据单元 "{name}" 吗？',
      batchFile: '确定要批量删除 {count} 个资源节点吗？',
    },
    warning: '此操作不可撤销节点，请谨慎操作。',
    confirmText: '确认删除',
  },
  fileContextMenu: {
    toast: {
      noFileSelected: '没有选中的数据单元',
      moveFailed: '移动数据单元失败',
      moveToRootSuccess: '已将数据单元移动到根节点',
      moveToFolderSuccess: '已将数据单元移动到 "{name}"',
      loadFoldersFailed: '加载数据目录列表失败',
      loadSubfoldersFailed: '加载子数据目录失败',
      noValidImageLink: '无有效原始数据链接',
      noThumbLink: '没有可复制的缩略图链接',
      thumbLinkCopied: '已复制缩略图链接',
      copyFailed: '复制失败',
      originalLinkCopied: '已复制原始数据链接',
      noOriginalLink: '没有可复制的原始数据链接',
      shortLinkCopied: '已复制短链',
      noShortLink: '没有可复制的短链',
      fileDeleted: '数据单元已删除',
      deleteFailed: '删除数据单元失败',
    },
    menu: {
      root: '根节点',
      preview: '预览',
      openNew: '新窗口打开',
      copyLink: '复制链接',
      copyThumb: '复制缩略图',
      copyOriginal: '复制原始数据',
      copyShort: '复制短链',
      download: '下载',
      moveTo: '移动到...',
      delete: '删除',
    },
    confirm: {
      deleteFile: '确定要删除数据单元 "{name}" 吗？此操作不可撤销。',
    },
  },
  folderManagement: {
    deleteSuccess: '数据目录 "{name}" 已删除',
    visibility: {
      public: '公开协议',
      private: '私密档案',
      switched: '数据目录 "{name}" 已设为{level}',
    },
  },
  dragSort: {
    folderSortUpdated: '数据目录排序已更新',
    fileSortUpdated: '数据单元排序已更新',
  },
  fileManagement: {
    copyLink: {
      protectedNotSupported: '受保护数据不支持外部传输，请在终端内查看',
      publicLinkCopied: '公开链接已复制到剪贴板',
      privateLinkCopied: '私有链接已复制到剪贴板',
    },
    download: {
      success: '数据单元 "{name}" 下载完成',
      cancelled: '下载被中止或失败',
    },
    delete: {
      success: '数据单元 "{name}" 已删除',
      selectFirst: '请先选择要删除的数据单元',
      batchSuccess: '成功删除 {successCount} 个数据单元',
      batchPartial: '成功删除 {successCount} 个数据单元，{failCount} 个删除失败',
      batchFailed: '删除失败，共 {failCount} 个数据单元删除失败',
    },
    accessLevel: {
      public: '公开协议',
      private: '私密档案',
      protected: '安全防护',
      switched: '数据单元 "{name}" 已设为{level}',
    },
  },
}
