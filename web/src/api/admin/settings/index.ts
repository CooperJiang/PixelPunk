import {
  getSettings,
  getSetting,
  createSetting,
  updateSetting,
  batchUpsertSettings,
  deleteSetting,
  getGlobalSettings,
  getGlobalSettingsByGroup,
  getLegalDocuments,
} from './common'
import { testMailServer } from './mail'
import { testAIConfig, checkAIConfiguration } from './ai'
import { testVectorConfig, regenerateVectors } from './vector'
import { getOAuthConfig, githubOAuthLogin, googleOAuthLogin, linuxdoOAuthLogin } from './oauth'
import { defaultSettings } from './defaults'
import { initializeSettings } from './utils'

/* ==================== 类型定义 ==================== */
export type {
  SettingType,
  SettingGroup,
  Setting,
  SettingsListResponse,
  SettingsQuery,
  OAuthProviders,
  GlobalSettingsResponse,
  TestMailServerParams,
  VectorTestParams,
  VectorTestResult,
  RegenerateVectorsParams,
  RegenerateVectorsResult,
  AITestResult,
  AITestResponse,
  QdrantTestParams,
  QdrantTestResult,
  LegalDocumentsResponse,
} from './types'

/* ==================== 重新导出所有函数 ==================== */
export {
  getSettings,
  getSetting,
  createSetting,
  updateSetting,
  batchUpsertSettings,
  deleteSetting,
  getGlobalSettings,
  getGlobalSettingsByGroup,
  getLegalDocuments,
  testMailServer,
  testAIConfig,
  checkAIConfiguration,
  testVectorConfig,
  regenerateVectors,
  getOAuthConfig,
  githubOAuthLogin,
  googleOAuthLogin,
  linuxdoOAuthLogin,
  defaultSettings,
  initializeSettings,
}

/* ==================== 默认导出（向后兼容） ==================== */
export default {
  getSettings,
  getSetting,
  createSetting,
  updateSetting,
  batchUpsertSettings,
  deleteSetting,
  initializeSettings,
  testMailServer,
  testVectorConfig,
  getGlobalSettings,
  getGlobalSettingsByGroup,
  getLegalDocuments,
  getOAuthConfig,
  githubOAuthLogin,
  googleOAuthLogin,
  linuxdoOAuthLogin,
  defaultSettings,
}
