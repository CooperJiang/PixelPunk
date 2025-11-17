/**
 * 可恢复上传组件
 */
export const resumableUploads = {
  title: '恢复数据传输',
  relativeTime: {
    justNow: '刚刚',
    minutesAgo: '{count}分钟前',
    hoursAgo: '{count}小时前',
    daysAgo: '{count}天前',
    empty: {
      title: '无中断传输任务',
      subtitle: '当前没有需要恢复的数据传输',
      labels: {
        chunks: '数据块',
        actions: {
          selectAll: '全部选择',
          clearSelection: '清空选择',
          deleteSelected: '删除已选',
          footer: {
            cancel: '取消',
            resume: '恢复传输',
          },
        },
      },
    },
  },
}
