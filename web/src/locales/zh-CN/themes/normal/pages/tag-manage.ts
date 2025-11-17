/**
 * 标签管理页面文案
 */
export const tagManage = {
  title: '标签管理',
  subtitle: '管理您的标签，组织和整理您的资源',
  search: '搜索标签...',
  create: '创建标签',
  batchMode: '批量模式',
  exitBatchMode: '退出批量模式',
  loading: '标签加载中...',
  batch: {
    viewFiles: '查看选中标签文件',
    delete: '批量删除 ({count})',
    merge: '合并标签',
    cancel: '取消选中',
    modeHint: '批量操作模式',
    modeDescription: '点击标签可选择/取消选择，已选择 {count} 个标签',
  },
  multiSelect: {
    modeHint: '多选模式',
    modeDescription: '按住 {modifier} / Alt 键点击标签进行多选，已选择 {count} 个标签',
  },
  normalMode: {
    hint: '提示：单击标签快速搜索，按住 {modifier} / Alt 键进行多选',
  },
  actions: {
    edit: '编辑标签',
    delete: '删除标签',
  },
  dialog: {
    create: {
      title: '创建标签',
      submit: '创建',
    },
    edit: {
      title: '编辑标签',
      submit: '更新',
    },
    merge: {
      title: '合并标签',
      targetLabel: '目标标签',
      targetPlaceholder: '请选择目标标签',
    },
    delete: {
      title: '删除标签',
      message: '确定要删除标签"{name}"吗？该标签已被使用 {count} 次，删除后将移除所有文件的该标签。',
      confirmText: '确认删除',
    },
    batchDelete: {
      title: '批量删除标签',
      message: '确定要删除选中的 {count} 个标签吗？此操作不可恢复。',
      confirmText: '确认删除',
    },
    mergeConfirm: {
      title: '合并标签',
      message: [
        '确定要将 {count} 个标签合并到"{targetName}"吗？',
        '源标签将被删除，其文件关联将转移到目标标签。',
        '此操作不可恢复。',
      ],
      confirmText: '确认合并',
    },
  },
  form: {
    name: {
      label: '标签名称',
      placeholder: '请输入标签名称',
      hint: '最多50个字符，建议使用简短易记的名称',
    },
    cancel: '取消',
  },
  merge: {
    description: '将选中的标签合并到目标标签，源标签将被删除，其文件关联将转移到目标标签。',
    selectedInfo: '已选择 {count} 个标签，选择目标标签后，其他 {remaining} 个标签将被合并到目标标签中。',
    confirm: '确认合并',
    cancel: '取消',
  },
  empty: {
    title: '没有找到符合条件的标签',
    description: '尝试调整筛选条件或创建新标签',
  },
  tagCloud: {
    groupCount: '{count} 个标签',
    groupUsage: '使用 {usage} 次',
    letterNav: '{letter} ({count})',
    letterEmpty: '{letter} (无)',
    usageCount: '{count} 次使用',
  },
  toast: {
    nameRequired: '请输入标签名称',
    updateSuccess: '更新标签成功',
    createSuccess: '创建标签成功',
    deleteSuccess: '删除标签成功',
    selectRequired: '请先选择要删除的标签',
    mergeMinRequired: '至少需要选择2个标签才能合并',
    mergeTargetRequired: '请选择目标标签',
    mergeTargetInvalid: '目标标签必须在选中的标签中',
    mergeSuccess: '合并标签成功',
    batchSelectRequired: '请先选择标签',
  },
  source: {
    manual: '手动创建',
    ai: 'AI识别',
    system: '系统导入',
    unknown: '未知',
  },
}
