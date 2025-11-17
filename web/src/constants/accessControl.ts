/**
 * 访问控制常量定义
 */

import type { Composer } from '@/composables/useTexts'

export const getAccessControlModeOptions = ($t: Composer['$t']) => [
  { label: $t('constants.accessControl.modes.ipWhitelist'), value: 'whitelist' as const },
  { label: $t('constants.accessControl.modes.ipBlacklist'), value: 'blacklist' as const },
]

export const getDomainControlModeOptions = ($t: Composer['$t']) => [
  { label: $t('constants.accessControl.modes.domainWhitelist'), value: 'whitelist' as const },
  { label: $t('constants.accessControl.modes.domainBlacklist'), value: 'blacklist' as const },
]

export const getRestrictionModeOptions = ($t: Composer['$t']) => [
  { label: $t('constants.accessControl.restrictionModes.strict'), value: 'strict' as const },
  { label: $t('constants.accessControl.restrictionModes.moderate'), value: 'moderate' as const },
  { label: $t('constants.accessControl.restrictionModes.loose'), value: 'loose' as const },
]

export const getBlockActionOptions = ($t: Composer['$t']) => [
  { label: $t('constants.accessControl.blockActions.block'), value: 'block' as const },
  { label: $t('constants.accessControl.blockActions.redirect'), value: 'redirect' as const },
  { label: $t('constants.accessControl.blockActions.thumbnail'), value: 'thumbnail' as const },
]

export const getBlockActionWarnings = ($t: Composer['$t']) => ({
  block: $t('constants.accessControl.warnings.block'),
  redirect: $t('constants.accessControl.warnings.redirect'),
  thumbnail: $t('constants.accessControl.warnings.thumbnail'),
  watermark: $t('constants.accessControl.warnings.watermark'),
})

export const getModeDescriptions = ($t: Composer['$t']) => ({
  whitelist: {
    ip: $t('constants.accessControl.descriptions.ipWhitelist'),
    domain: $t('constants.accessControl.descriptions.domainWhitelist'),
  },
  blacklist: {
    ip: $t('constants.accessControl.descriptions.ipBlacklist'),
    domain: $t('constants.accessControl.descriptions.domainBlacklist'),
  },
})

export const getRestrictionModeDescriptions = ($t: Composer['$t']) => ({
  strict: $t('constants.accessControl.descriptions.strictMode'),
  moderate: $t('constants.accessControl.descriptions.moderateMode'),
  loose: $t('constants.accessControl.descriptions.looseMode'),
})
