/**
 * 作者信息组件
 */
export const author = {
  loading: '加载中...',
  header: {
    joinedLabel: '加入 {days} 天',
    websiteLabel: '个人网站',
    stats: {
      views: '浏览量',
      shares: '分享数',
      images: '文件数',
    },
  },
  folders: {
    rootTitle: '文件夹',
    rootSubtitle: '浏览作者的文件夹',
    childTitle: '子文件夹',
    childSubtitle: '当前文件夹下的子文件夹',
    countLabel: '个文件',
    createdLabel: '{time}',
  },
  images: {
    sectionTitle: '文件',
    sectionSubtitle: '浏览作者的文件作品',
    countLabel: '个资源',
    loadingText: '加载文件中...',
    paginationInfo: '第 {current} 页，共 {total} 页',
    actions: {
      preview: {
        label: '预览',
        tooltip: '点击预览文件',
      },
      download: {
        label: '下载',
        tooltip: '下载原图',
      },
      detail: {
        label: '详情',
        tooltip: '复制文件链接',
      },
    },
  },
  shares: {
    sectionTitle: '分享',
    sectionSubtitle: '浏览作者的公开分享',
    statusBadge: '状态',
    meta: {
      viewsLabel: '浏览量',
      createdLabel: '{time}',
      keyLabel: '分享密钥',
    },
    empty: {
      title: '暂无分享',
      description: '该作者还没有创建任何公开分享',
    },
  },
  empty: {
    title: '暂无内容',
    description: '该作者还没有上传任何内容',
    suggestion: '可以稍后再来查看',
    retry: '重新加载',
  },
  error: {
    title: '加载失败',
    retry: '重试',
    invalidAuthor: '无效的作者ID',
    loadFailed: '加载作者信息失败',
    loadFolderFailed: '加载文件夹内容失败',
  },
  toolbar: {
    searchPlaceholder: '搜索内容...',
    refresh: '刷新',
    gridView: '网格视图',
    listView: '列表视图',
    fullscreenEnter: '全屏',
    fullscreenExit: '退出全屏',
  },
  breadcrumb: {
    root: '根目录',
  },
  toast: {
    copySuccess: '链接已复制到剪贴板',
    copyFailed: '复制失败',
  },
}
