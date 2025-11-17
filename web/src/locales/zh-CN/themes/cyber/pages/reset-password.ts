/**
 * 重置密码页面文案 - 赛博风格
 */
export const resetPassword = {
  title: '密码重置协议',
  verifying: {
    title: '验证重置令牌',
    loading: '正在验证重置令牌...',
  },
  error: {
    missingToken: '缺少重置令牌',
    invalidToken: '令牌无效',
    invalidOrExpired: '令牌无效或已过期',
    hint: '请重新申请密码重置',
  },
  form: {
    subtitle: '为账号 {email} 设置新密码',
    newPassword: {
      label: '新密码',
      placeholder: '请输入新密码 (6-50位)',
    },
    confirmPassword: {
      label: '确认密码',
      placeholder: '请再次输入新密码',
    },
    submit: '重置密码',
    submitting: '重置中...',
  },
  validation: {
    passwordRequired: '请输入新密码',
    passwordTooShort: '密码长度至少为6位',
    passwordTooLong: '密码长度不能超过50位',
    confirmRequired: '请确认新密码',
    passwordMismatch: '两次输入的密码不一致',
  },
  actions: {
    backToLogin: '返回登录页面',
  },
  protocol: '系统安全协议已启动 / SYSTEM SECURITY PROTOCOL',
}
