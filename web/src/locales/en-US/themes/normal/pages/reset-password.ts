/**
 * Reset password page texts
 */
export const resetPassword = {
  title: 'Reset Password',
  verifying: {
    title: 'Verifying Reset Token',
    loading: 'Verifying reset token...',
  },
  error: {
    missingToken: 'Reset token missing',
    invalidToken: 'Token invalid',
    invalidOrExpired: 'Token invalid or expired',
    hint: 'Please request password reset again',
  },
  form: {
    subtitle: 'Set new password for account {email}',
    newPassword: {
      label: 'New Password',
      placeholder: 'Please enter new password (6-50 characters)',
    },
    confirmPassword: {
      label: 'Confirm Password',
      placeholder: 'Please enter new password again',
    },
    submit: 'Reset Password',
    submitting: 'Resetting...',
  },
  validation: {
    passwordRequired: 'Please enter new password',
    passwordTooShort: 'Password must be at least 6 characters',
    passwordTooLong: 'Password cannot exceed 50 characters',
    confirmRequired: 'Please confirm new password',
    passwordMismatch: 'Passwords do not match',
  },
  actions: {
    backToLogin: 'Back to Login Page',
  },
  protocol: 'System Security Protocol Activated / SYSTEM SECURITY PROTOCOL',
}
