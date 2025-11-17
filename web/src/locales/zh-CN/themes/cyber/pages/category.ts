/**
 * 分类管理页面文案 - 赛博风格
 */
export const category = {
  page: {
    title: '分类管理',
    subtitle: '管理数据分类，组织和整理资源节点',
    create: '创建分类',
    refresh: '刷新',
    filter: '筛选',
    hideFilter: '收起筛选',
  },
  table: {
    columns: {
      name: '分类名称',
      description: '描述',
      source: '来源',
      fileCount: '数据单元数量',
      sortOrder: '排序',
      status: '状态',
      createdAt: '创建时间',
      actions: '操作',
    },
    noDescription: '暂无描述',
    clickToView: '点击查看该分类下的所有数据单元',
    loading: '正在加载分类数据...',
  },
  source: {
    system: '系统',
    user: '用户',
    imported: '导入',
  },
  status: {
    active: '启用',
    archived: '归档',
  },
  actions: {
    edit: '编辑分类',
    archive: '归档分类',
    activate: '激活分类',
    delete: '删除分类',
  },
  empty: {
    title: '暂无分类数据',
    description: '尚未创建任何分类，点击「创建分类」按钮开始创建您的第一个分类',
    action: '创建分类',
  },
  dialog: {
    edit: {
      title: '编辑分类',
      cancel: '取消',
      save: '保存更改',
    },
    delete: {
      title: '删除分类',
      message: '确定要删除此分类吗？',
    },
    create: {
      title: '创建分类',
      cancel: '取消',
      create: '创建分类',
    },
  },
  form: {
    name: {
      label: '分类名称',
      placeholder: '请输入分类名称',
      required: '请输入分类名称',
    },
    description: {
      label: '分类描述',
      placeholder: '请输入分类描述（可选）',
    },
    sortOrder: {
      label: '排序值',
      placeholder: '请输入排序值',
      hint: '排序值越小，分类在列表中显示越靠前',
    },
    status: {
      label: '状态',
      active: '启用',
      archived: '归档',
    },
  },
  filter: {
    title: '筛选条件',
    sortBy: '排序方式',
    status: {
      label: '状态',
      all: '全部状态',
      active: '激活',
      archived: '归档',
    },
    source: {
      label: '来源',
      all: '全部',
      system: '系统',
      user: '用户',
      imported: '导入',
    },
    search: {
      placeholder: '检索分类名称或描述...',
    },
    apply: '应用',
    reset: '重置',
  },
  batch: {
    selected: '已选择 {count} 项',
    clearSelection: '取消选择',
    clear: '清空选择',
    delete: '批量删除',
    archive: '批量归档',
    activate: '批量激活',
    deleteConfirm: '确定要删除选中的 {count} 个分类吗？',
    archiveConfirm: '确定要归档选中的 {count} 个分类吗？',
    activateConfirm: '确定要激活选中的 {count} 个分类吗？',
  },
  toast: {
    createSuccess: '分类创建成功',
    createError: '创建分类失败',
    createWarning: '请输入分类名称',
    updateSuccess: '分类更新成功',
    updateError: '更新分类失败',
    deleteSuccess: '分类删除成功',
    deleteError: '删除分类失败',
    archiveSuccess: '归档成功',
    activateSuccess: '激活成功',
    archiveError: '归档失败',
    activateError: '激活失败',
    batchDeleteSuccess: '成功删除 {count} 个分类',
    batchDeleteError: '批量删除失败',
    batchArchiveSuccess: '批量归档成功',
    batchActivateSuccess: '批量激活成功',
    fetchError: '获取分类列表失败',
    selectWarning: '请先选择要删除的分类',
  },
  sourceInfo: {
    manual: '手动创建',
    aiSuggestion: 'AI建议',
    system: '系统模板',
    imported: '批量导入',
  },
  statusInfo: {
    active: '激活',
    archived: '归档',
  },
  placeholders: {
    search: '搜索分类名称或描述...',
    status: '状态',
    name: '请输入分类名称',
    description: '请输入分类描述（选填）',
    sortOrder: '数值越小排序越靠前',
  },
}
