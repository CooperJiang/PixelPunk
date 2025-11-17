/**
 */
export const multiSelector = {
  addOption: '注入选项',
  deleteOption: '移除指令',
  confirmAdd: '确认注入',
  cancel: '中止',
  placeholders: {
    basic: '示例: 30m, 2h, 1d',
    withPermanent: '示例: 30m, 2h, 1d, permanent',
  },
  duration: {
    permanent: '永久存储',
    minute: '分钟',
    hour: '小时',
    day: '天',
  },
  hints: {
    format: '支持 30m(分钟) / 1h(小时) / 7d(天)，最低 1 分钟',
    permanent: ' · permanent(永久)',
  },
  errors: {
    empty: '请输入时长参数',
    guestNoPermanent: '游客模式禁用永久存储',
    invalidFormat: '格式错误，应为 30m / 1h / 7d',
    mustBePositive: '时长必须大于 0',
    minutesMin: '分钟数不能小于 1',
    minutesMax: '分钟数上限为 365 天',
    hoursMax: '小时数上限为 365 天',
    guestHoursMax: '游客上限为 720 小时(30天)',
    daysMax: '天数上限为 365 天',
    guestDaysMax: '游客上限为 30 天',
    exists: '该时长已存在',
  },
}
