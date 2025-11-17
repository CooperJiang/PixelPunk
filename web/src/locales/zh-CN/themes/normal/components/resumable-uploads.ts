/**
 * 可恢复上传组件
 */
export const resumableUploads = {
  title: '恢复上传',
  relativeTime: {
    justNow: '刚刚',
    minutesAgo: '{count}分钟前',
    hoursAgo: '{count}小时前',
    daysAgo: '{count}天前',
    empty: {
      title: '没有可恢复的上传',
      subtitle: '当前没有中断的上传任务',
      labels: {
        chunks: '分片',
        actions: {
          selectAll: '全选',
          clearSelection: '清空选择',
          deleteSelected: '删除选中',
          footer: {
            cancel: '取消',
            resume: '恢复上传',
          },
        },
      },
    },
  },
}
