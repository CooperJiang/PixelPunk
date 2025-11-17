/**
 * 自动化任务页面文案
 */
export const automation = {
  page: {
    title: '自动任务',
    subtitle: '实时监控您的 AI 打标和向量化任务进度',
    refresh: '刷新',
    tip: '队列会自动处理您上传的图片，无需手动操作。系统会按顺序依次完成所有任务。',
  },
  tagging: {
    title: 'AI 打标任务',
    subtitle: '自动识别图片内容并生成标签',
    statusRunning: '运行中',
    statusPaused: '已暂停',
    progress: {
      label: '总体进度',
      completed: '{done} / {total} 已完成',
    },
    queue: {
      position: '您的图片在队列第 {position} 位',
      pending: '队列中有 {count} 张图片等待处理',
    },
    status: {
      none: '未处理',
      pending: '队列中',
      processing: '处理中',
      done: '已完成',
      failed: '失败',
      ignored: '已忽略',
    },
    empty: {
      title: '暂无 AI 打标任务',
      subtitle: '上传图片后会自动开始 AI 分析',
    },
  },
  vector: {
    title: '向量化任务',
    subtitle: '将图片转换为向量用于相似度搜索',
    statusRunning: '运行中',
    statusPaused: '已暂停',
    progress: {
      label: '总体进度',
      completed: '{done} / {total} 已完成',
    },
    queue: {
      position: '您的向量任务在队列第 {position} 位',
      pending: '队列中有 {count} 个向量任务等待处理',
    },
    status: {
      pending: '待处理',
      reset: '重置',
      processing: '处理中',
      completed: '已完成',
      failed: '失败',
    },
    empty: {
      title: '暂无向量化任务',
      subtitle: '图片完成 AI 打标后会自动开始向量化',
    },
  },
  taggingHistory: {
    filter: {
      statusPlaceholder: '选择状态',
    },
  },
  history: {
    title: '打标历史记录',
    recordCount: '(共 {count} 条)',
    refresh: '刷新',
    empty: '暂无打标记录',
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
    loadHistoryError: '加载打标记录失败',
    loadDataError: '加载任务数据失败',
  },
}
