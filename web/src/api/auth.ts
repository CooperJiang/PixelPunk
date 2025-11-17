import { post } from '@/utils/network/http'

export interface ForgotPasswordRequest {
  email: string
}

export interface ForgotPasswordResponse {
  message: string
}

export interface VerifyResetTokenRequest {
  token: string
}

export interface VerifyResetTokenResponse {
  valid: boolean
  email?: string
  error?: string
}

export interface ResetPasswordRequest {
  token: string
  new_password: string
}

export interface ResetPasswordResponse {
  message: string
}

export function useAuthApi() {
  const forgotPassword = async (email: string): Promise<void> => {
    await post('/auth/forgot-password', { email }, { autoShowSuccess: true })
  }

  const verifyResetToken = async (token: string): Promise<{ valid: boolean; email: string }> => {
    const response = await post<{ valid: boolean; email: string }>('/auth/verify-reset-token', { token }, { silent: true })
    return response.data
  }

  const resetPassword = async (token: string, newPassword: string): Promise<void> => {
    await post(
      '/auth/reset-password-token',
      {
        token,
        new_password: newPassword,
      },
      { autoShowSuccess: true }
    )
  }

  return {
    forgotPassword,
    verifyResetToken,
    resetPassword,
  }
}
