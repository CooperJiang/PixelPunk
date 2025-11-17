/**
 * 公告管理模块文案 - 赛博风格
 */
export const announcements = {
  title: '公告管理',
  subtitle: '管理系统公告，支持Markdown和HTML格式',

  actions: {
    create: '创建公告',
    refresh: '刷新',
    filter: '筛选',
    edit: '编辑',
    delete: '删除',
    pin: '置顶',
    unpin: '取消置顶',
    cancel: '取消',
    update: '更新',
  },

  status: {
    draft: '草稿',
    published: '已发布',
    archived: '已归档',
  },

  table: {
    title: '公告标题',
    summary: '摘要',
    status: '状态',
    viewCount: '浏览次数',
    publishedAt: '发布时间',
    createdAt: '创建时间',
    actions: '操作',
    noSummary: '无摘要',
    notSet: '未设置',
    loadingText: '加载中...',
  },

  empty: {
    title: '暂无公告',
    description: '点击"创建公告"按钮开始创建第一条公告',
  },

  form: {
    title: {
      label: '标题',
      placeholder: '请输入公告标题',
    },
    summary: {
      label: '摘要',
      placeholder: '请输入公告摘要（可选）',
    },
    content: {
      label: '内容',
      placeholder: '请输入Markdown格式的内容...',
      hint: '支持Markdown格式，可实时预览效果。支持拖拽或粘贴图片上传',
    },
    status: {
      label: '状态',
    },
    options: {
      label: '选项',
      pinned: '置顶显示',
      pinnedHint: '首次访问自动弹窗',
    },
  },

  dialog: {
    createTitle: '创建公告',
    editTitle: '编辑公告',
  },

  messages: {
    fetchError: '加载公告列表失败: ',
    deleteConfirm: '确定要删除这条公告吗？此操作不可撤销。',
    deleteSuccess: '删除成功',
    deleteError: '删除失败: ',
    pinSuccess: '{action}成功',
    pinError: '{action}失败: ',
    createSuccess: '创建公告成功',
    createError: '创建失败: ',
    updateSuccess: '更新公告成功',
    updateError: '更新失败: ',
    operationError: '操作失败: ',
    uploadSuccess: '成功上传 {count} 张图片',
    uploadError: '图片上传失败: ',
    uploadFailed: '上传失败',
  },

  validation: {
    titleRequired: '请输入公告标题',
    contentRequired: '请输入公告内容',
  },

  filter: {
    status: {
      label: '状态',
      all: '全部',
    },
    pinned: {
      label: '置顶',
      all: '全部',
      onlyPinned: '仅置顶',
      notPinned: '非置顶',
    },
    search: {
      label: '搜索',
      placeholder: '搜索标题或摘要...',
    },
    actions: {
      query: '查询',
      reset: '重置',
    },
  },

  settings: {
    systemToggle: '系统开关',
    drawerPosition: {
      label: '抽屉位置',
      left: '左侧',
      right: '右侧',
    },
    displayLimit: '显示数量',
    autoShowDelay: '弹窗延迟',
    saveButton: '保存配置',
    messages: {
      fetchError: '获取配置失败:',
      limitRange: '显示数量必须在 1-20 之间',
      saveSuccess: '配置保存成功',
      saveError: '保存失败: ',
    },
    units: {
      second: '秒',
    },
  },
}
