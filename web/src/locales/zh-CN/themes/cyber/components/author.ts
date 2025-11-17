/**
 * 作者信息组件
 */
export const author = {
  loading: '数据加载中...',
  header: {
    joinedLabel: '接入时间 {days} 天',
    websiteLabel: '网络地址',
    stats: {
      views: '访问量',
      shares: '分享数据',
      images: '文件存储量',
    },
  },
  folders: {
    rootTitle: '数据容器',
    rootSubtitle: '浏览用户数据存储结构',
    childTitle: '子容器',
    childSubtitle: '当前容器下的子级存储单元',
    countLabel: '个数据单元',
    createdLabel: '{time}',
  },
  images: {
    sectionTitle: '文件数据库',
    sectionSubtitle: '浏览用户的文件数据集',
    countLabel: '个资源数据',
    loadingText: '文件数据加载中...',
    paginationInfo: '第 {current} 节，共 {total} 节',
    actions: {
      preview: {
        label: '数据预览',
        tooltip: '启动文件预览模式',
      },
      download: {
        label: '数据下载',
        tooltip: '获取原始文件数据',
      },
      detail: {
        label: '详细信息',
        tooltip: '复制数据链接',
      },
    },
  },
  shares: {
    sectionTitle: '分享网络',
    sectionSubtitle: '浏览用户的公开数据分享',
    statusBadge: '连接状态',
    meta: {
      viewsLabel: '访问次数',
      createdLabel: '{time}',
      keyLabel: '访问密钥',
    },
    empty: {
      title: '网络空间为空',
      description: '该用户尚未建立任何公开数据连接',
    },
  },
  empty: {
    title: '数据库空置',
    description: '该用户尚未上传任何数据内容',
    suggestion: '可稍后重新扫描',
    retry: '重新连接',
  },
  error: {
    title: '连接失败',
    retry: '重新尝试',
    invalidAuthor: '无效的用户ID',
    loadFailed: '用户档案加载失败',
    loadFolderFailed: '数据容器加载失败',
  },
  toolbar: {
    searchPlaceholder: '数据检索...',
    refresh: '重新同步',
    gridView: '矩阵视图',
    listView: '列表视图',
    fullscreenEnter: '全屏模式',
    fullscreenExit: '退出全屏',
  },
  breadcrumb: {
    root: '根目录',
  },
  toast: {
    copySuccess: '数据链接已复制到缓冲区',
    copyFailed: '数据复制操作失败',
  },
}
