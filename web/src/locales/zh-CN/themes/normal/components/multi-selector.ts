/**
 */
export const multiSelector = {
  addOption: '添加选项',
  deleteOption: '删除选项',
  confirmAdd: '确认添加',
  cancel: '取消',
  placeholders: {
    basic: '如: 30m, 2h, 1d',
    withPermanent: '如: 30m, 2h, 1d, permanent',
  },
  duration: {
    permanent: '永久',
    minute: '分钟',
    hour: '小时',
    day: '天',
  },
  hints: {
    format: '支持格式：30m(分钟), 1h(小时), 7d(天)，最低1分钟',
    permanent: ' · permanent(永久)',
  },
  errors: {
    empty: '请输入时长',
    guestNoPermanent: '游客不支持永久存储',
    invalidFormat: '格式错误，支持：30m(分钟)/1h(小时)/7d(天)',
    mustBePositive: '时长必须大于0',
    minutesMin: '分钟数不能小于1',
    minutesMax: '分钟数最多365天',
    hoursMax: '小时数最多365天',
    guestHoursMax: '游客最多720小时(30天)',
    daysMax: '天数最多365天',
    guestDaysMax: '游客最多30天',
    exists: '该时长已存在',
  },
}
