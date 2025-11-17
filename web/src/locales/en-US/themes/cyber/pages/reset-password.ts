/**
 * Reset Password Page Text - Cyber Style
 */
export const resetPassword = {
  title: 'Password Reset Protocol',
  verifying: {
    title: 'Verifying Reset Token',
    loading: 'Verifying reset token...',
  },
  error: {
    missingToken: 'Missing reset token',
    invalidToken: 'Token invalid',
    invalidOrExpired: 'Token invalid or expired',
    hint: 'Please re-request password reset',
  },
  form: {
    subtitle: 'Set new password for account {email}',
    newPassword: {
      label: 'New Password',
      placeholder: 'Please enter new password (6-50 characters)',
    },
    confirmPassword: {
      label: 'Confirm Password',
      placeholder: 'Please re-enter new password',
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
  protocol: 'SYSTEM SECURITY PROTOCOL ACTIVATED / SYSTEM SECURITY PROTOCOL',
}
