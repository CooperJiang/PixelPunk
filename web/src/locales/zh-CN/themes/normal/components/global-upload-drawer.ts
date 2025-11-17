/**
 * 全局上传抽屉组件
 */
export const globalUploadDrawer = {
  stats: {
    uploading: '上传中',
    pending: '等待中',
    completed: '已完成',
    failed: '失败',
    batchActions: {
      pauseAll: '暂停全部',
      resumeAll: '继续全部',
      clearCompleted: '清除已完成',
      clearAll: '清空全部',
      statusMessages: {
        uploading: '上传中',
        completed: '上传完成',
        failed: '上传失败',
        paused: '已暂停',
        pending: '等待上传',
        preparing: '准备中',
        actions: {
          pause: '暂停',
          resume: '继续',
          retry: '重试',
          copyLink: '复制链接',
          remove: '移除',
          messages: {
            pausedAll: '已暂停全部上传',
            resumedAll: '已继续全部上传',
            clearedCompleted: '已清除 {count} 个完成的任务',
            clearedAll: '已清空所有任务',
            linkCopied: '链接已复制',
            copyFailed: '复制失败',
            misc: {
              chunkedUpload: '分片上传',
              noTasks: '暂无上传任务',
              goToUpload: '去上传页面',
              viewFullPage: '查看完整页面',
            },
          },
        },
      },
    },
  },
}
