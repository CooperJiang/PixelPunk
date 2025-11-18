import { get, post } from '@/utils/network/http'
import type { ApiResult } from '@/utils/network/http-types'

export interface GithubOAuthConfig {
  enabled: boolean
  client_id: string
  client_secret: string
  redirect_uri: string
  scope: string
}

export interface GoogleOAuthConfig {
  enabled: boolean
  client_id: string
  client_secret: string
  redirect_uri: string
  scope: string
}

export interface LinuxdoOAuthConfig {
  enabled: boolean
  client_id: string
  client_secret: string
  redirect_uri: string
  scope: string
}

export interface OAuthConfigResponse {
  github: GithubOAuthConfig
  google: GoogleOAuthConfig
  linuxdo: LinuxdoOAuthConfig
}

export interface GithubOAuthLoginRequest {
  code: string
}

export interface OAuthLoginResponse {
  token: string
  userInfo: {
    id: number
    username: string
    email: string
    avatar: string
    role: number
    status: number
  }
  email: string
}

export function getOAuthConfig(): Promise<ApiResult<OAuthConfigResponse>> {
  return get<OAuthConfigResponse>('/common/settings/oauth')
}

export function githubOAuthLogin(code: string): Promise<ApiResult<OAuthLoginResponse>> {
  return post<OAuthLoginResponse>('/auth/oauth/github/login', { code }, { silent: true })
}

export function googleOAuthLogin(code: string): Promise<ApiResult<OAuthLoginResponse>> {
  return post<OAuthLoginResponse>('/auth/oauth/google/login', { code }, { silent: true })
}

export function linuxdoOAuthLogin(code: string): Promise<ApiResult<OAuthLoginResponse>> {
  return post<OAuthLoginResponse>('/auth/oauth/linuxdo/login', { code }, { silent: true })
}

export interface TestProxyRequest {
  proxy_dynamic?: boolean
  proxy_api_url?: string
  proxy_type: string
  proxy_host: string
  proxy_port: string
  proxy_username?: string
  proxy_password?: string
}

export interface TestProxyResponse {
  success: boolean
  message: string
  latency?: number
}

export function testProxy(params: TestProxyRequest): Promise<ApiResult<TestProxyResponse>> {
  return post<TestProxyResponse>('/settings/test-proxy', params)
}
