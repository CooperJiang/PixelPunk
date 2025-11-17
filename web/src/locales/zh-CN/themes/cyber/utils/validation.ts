/**
 * 验证工具多语言配置 - 赛博主题
 */
export const validation = {
  errors: {
    required: '此字段为必填项',
    type: '字段类型不正确',
    min: '数值过小，最小值为 {min}',
    max: '数值过大，最大值为 {max}',
    length: '长度不符合规范',
    pattern: '格式校验失败',
    email: '请输入有效的电子邮箱',
    phone: '请输入有效的手机号码',
    url: '请输入有效的URL地址',
    username: '用户名只能包含字母、数字、下划线和连字符，长度3-20位',
    password: '密码必须包含大小写字母和数字，长度至少8位',
    custom: '验证失败',
    mustBeString: '必须是字符串类型',
    mustBeNumber: '必须是数字类型',
    mustBePositive: '必须是正数',
    minLength: '最少{n}个字符',
    maxLength: '最多{n}个字符',
  },
  password: {
    strength: {
      weak: '弱',
      medium: '中',
      strong: '强',
      veryStrong: '极强',
    },
    feedback: {
      minLength: '密码长度至少8位',
      requireLowercase: '需要包含小写字母',
      requireUppercase: '需要包含大写字母',
      requireNumber: '需要包含数字',
      suggestSpecial: '建议包含特殊字符',
      avoidRepeat: '避免连续重复字符',
    },
  },
  time: {
    days: '{n}天',
    hours: '{n}小时',
    minutes: '{n}分钟',
    seconds: '{n}秒',
    daysHours: '{d}天 {h}小时',
    hoursMinutes: '{h}小时 {m}分钟',
    minutesSeconds: '{m}分钟 {s}秒',
  },
}
