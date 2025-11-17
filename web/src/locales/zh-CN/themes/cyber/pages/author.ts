/**
 * 作者页面文案 - 赛博风格
 */
export const author = {
  loading: '数据加载中...',

  header: {
    joinedLabel: '接入时间',
    websiteLabel: '网络节点',
    stats: {
      views: '访问量',
      shares: '共享数',
      images: '数据单元',
    },
  },

  folders: {
    rootTitle: '所有目录',
    rootSubtitle: '浏览作者的所有数据目录',
    childTitle: '子目录',
    childSubtitle: '当前目录下的子目录',
    countLabel: '单元数',
    createdLabel: '创建时间',
    unknownTime: '未知',
    rootFolder: '根目录',
  },

  shares: {
    sectionTitle: '公开共享',
    sectionSubtitle: '作者的公开数据共享',
    statusBadge: '状态',
    meta: {
      viewsLabel: '访问',
      createdLabel: '创建于',
      keyLabel: '共享密钥',
    },
    empty: {
      title: '暂无共享',
      description: '该作者还没有创建任何公开共享',
    },
  },

  toolbar: {
    searchPlaceholder: '搜索目录...',
    refresh: '刷新',
    fullscreenEnter: '全屏模式',
    fullscreenExit: '退出全屏',
  },

  empty: {
    title: '暂无数据',
    description: '该作者还没有上传任何数据',
    suggestion: '请稍后再来查看',
    retry: '重试',
    folderEmpty: {
      title: '目录为空',
      description: '该目录暂时没有任何内容',
      suggestion: '返回上级目录或浏览其他内容',
    },
  },

  error: {
    title: '加载失败',
    retry: '重试',
    invalidAuthor: '无效的作者ID',
    loadFailed: '加载作者信息失败',
    loadFolderFailed: '加载目录失败',
  },

  toast: {
    copySuccess: '复制成功',
    copyFailed: '复制失败',
  },
}
