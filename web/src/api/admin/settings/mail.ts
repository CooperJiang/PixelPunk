/**
 * 邮件设置相关 API

 */
import type { TranslationFunction } from '@/composables/useTexts'
import { post } from '@/utils/network/http'
import type { ApiResult } from '@/utils/network/http-types'
import type { TestMailServerParams, Setting } from './types'

/* 测试邮件服务器 */
export function testMailServer(params: TestMailServerParams): Promise<ApiResult<{ success: boolean }>> {
  return post<{ success: boolean }>('/settings/test-mail', params)
}

/* 获取邮件设置（带翻译） */
export function getMailDefaults($t: TranslationFunction): Setting[] {
  return [
    {
      key: 'smtp_host',
      value: 'smtp.example.com',
      type: 'string',
      group: 'mail',
      description: $t('api.settingsDefaults.mail.smtp_host'),
      is_system: true,
    },
    {
      key: 'smtp_port',
      value: '587',
      type: 'number',
      group: 'mail',
      description: $t('api.settingsDefaults.mail.smtp_port'),
      is_system: true,
    },
    {
      key: 'smtp_encryption',
      value: 'tls',
      type: 'string',
      group: 'mail',
      description: $t('api.settingsDefaults.mail.smtp_encryption'),
      is_system: true,
    },
    {
      key: 'smtp_username',
      value: 'noreply@example.com',
      type: 'string',
      group: 'mail',
      description: $t('api.settingsDefaults.mail.smtp_username'),
      is_system: true,
    },
    {
      key: 'smtp_password',
      value: '',
      type: 'string',
      group: 'mail',
      description: $t('api.settingsDefaults.mail.smtp_password'),
      is_system: true,
    },
    {
      key: 'smtp_from_address',
      value: 'noreply@example.com',
      type: 'string',
      group: 'mail',
      description: $t('api.settingsDefaults.mail.smtp_from_address'),
      is_system: true,
    },
    {
      key: 'smtp_from_name',
      value: 'PixelPunk',
      type: 'string',
      group: 'mail',
      description: $t('api.settingsDefaults.mail.smtp_from_name'),
      is_system: true,
    },
  ]
}

/* 向后兼容 - 默认中文 */
export const mailDefaults: Setting[] = [
  {
    key: 'smtp_host',
    value: 'smtp.example.com',
    type: 'string',
    group: 'mail',
    description: 'SMTP服务器',
    is_system: true,
  },
  {
    key: 'smtp_port',
    value: '587',
    type: 'number',
    group: 'mail',
    description: 'SMTP端口',
    is_system: true,
  },
  {
    key: 'smtp_encryption',
    value: 'tls',
    type: 'string',
    group: 'mail',
    description: '加密类型',
    is_system: true,
  },
  {
    key: 'smtp_username',
    value: 'noreply@example.com',
    type: 'string',
    group: 'mail',
    description: 'SMTP用户名',
    is_system: true,
  },
  {
    key: 'smtp_password',
    value: '',
    type: 'string',
    group: 'mail',
    description: 'SMTP密码',
    is_system: true,
  },
  {
    key: 'smtp_from_address',
    value: 'noreply@example.com',
    type: 'string',
    group: 'mail',
    description: '发件人地址',
    is_system: true,
  },
  {
    key: 'smtp_from_name',
    value: 'PixelPunk',
    type: 'string',
    group: 'mail',
    description: '发件人名称',
    is_system: true,
  },
]
