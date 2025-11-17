/**
 * 自动化任务页面文案 - 赛博风格
 */
export const automation = {
  page: {
    title: '自动化指令',
    subtitle: '实时监控 AI 标记与向量化任务执行状态',
    refresh: '刷新',
    tip: '队列将自动处理注入的数据单元，无需人工干预。系统按序列依次执行所有指令。',
  },
  tagging: {
    title: 'AI 标记任务',
    subtitle: '自动识别数据单元内容并生成标签',
    statusRunning: '运行中',
    statusPaused: '已暂停',
    progress: {
      label: '执行进度',
      completed: '{done} / {total} 已完成',
    },
    queue: {
      position: '您的数据单元在队列第 {position} 位',
      pending: '队列中有 {count} 个数据单元等待处理',
    },
    status: {
      none: '未处理',
      pending: '队列中',
      processing: '执行中',
      done: '已完成',
      failed: '失败',
      ignored: '已忽略',
    },
    empty: {
      title: '暂无 AI 标记任务',
      subtitle: '注入数据单元后将自动开始 AI 分析',
    },
  },
  vector: {
    title: '向量化任务',
    subtitle: '将数据单元转换为向量用于相似度检索',
    statusRunning: '运行中',
    statusPaused: '已暂停',
    progress: {
      label: '执行进度',
      completed: '{done} / {total} 已完成',
    },
    queue: {
      position: '您的向量任务在队列第 {position} 位',
      pending: '队列中有 {count} 个向量任务等待处理',
    },
    status: {
      pending: '待处理',
      reset: '重置',
      processing: '执行中',
      completed: '已完成',
      failed: '失败',
    },
    empty: {
      title: '暂无向量化任务',
      subtitle: '数据单元完成 AI 标记后将自动开始向量化',
    },
  },
  taggingHistory: {
    filter: {
      statusPlaceholder: '选择状态',
    },
  },
  history: {
    title: '标记历史记录',
    recordCount: '(共 {count} 条)',
    refresh: '刷新',
    empty: '暂无标记记录',
    columns: {
      preview: '预览',
      filename: '文件名',
      status: '状态',
      format: '格式',
      size: '大小',
      resolution: '分辨率',
      createdAt: '创建时间',
      updatedAt: '更新时间',
      error: '错误信息',
    },
  },
  toast: {
    loadHistoryError: '加载标记记录失败',
    loadDataError: '加载任务数据失败',
  },
}
