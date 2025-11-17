/**
 * 全局上传抽屉组件
 */
export const globalUploadDrawer = {
  stats: {
    uploading: '传输中',
    pending: '队列中',
    completed: '已完成',
    failed: '失败',
    batchActions: {
      pauseAll: '暂停全部',
      resumeAll: '恢复全部',
      clearCompleted: '清除完成项',
      clearAll: '清空队列',
      statusMessages: {
        uploading: '数据传输中',
        completed: '传输完成',
        failed: '传输失败',
        paused: '传输暂停',
        pending: '等待传输',
        preparing: '准备传输',
        actions: {
          pause: '暂停',
          resume: '恢复',
          retry: '重试',
          copyLink: '复制链接',
          remove: '移除',
          messages: {
            pausedAll: '全部传输已暂停',
            resumedAll: '全部传输已恢复',
            clearedCompleted: '已清除 {count} 个完成任务',
            clearedAll: '传输队列已清空',
            linkCopied: '链接已复制至剪贴板',
            copyFailed: '复制操作失败',
            misc: {
              chunkedUpload: '分块传输',
              noTasks: '无传输任务',
              goToUpload: '前往传输控制台',
              viewFullPage: '查看完整控制台',
            },
          },
        },
      },
    },
  },
}
