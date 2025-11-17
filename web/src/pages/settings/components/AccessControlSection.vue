<script setup lang="ts">
  import { computed, onMounted, reactive, ref } from 'vue'
  import { createOrUpdateAccessControl, getAccessControlConfig, resetAccessControlConfig } from '@/api/user'
  import type { AccessControlConfig } from '@/api/types/index'
  import { useToast } from '@/components/Toast/useToast'
  import {
    getAccessControlModeOptions,
    getBlockActionOptions,
    getDomainControlModeOptions,
    getRestrictionModeOptions,
  } from '@/constants'
  import { useTexts } from '@/composables/useTexts'
  import { showConfirm } from '@/utils/dialog'

  const toast = useToast()
  const loading = ref(false)
  const { $t } = useTexts()

  const defaultConfig: AccessControlConfig = {
    enabled: false,
    ip_mode: 'whitelist',
    domain_mode: 'whitelist',
    restriction_mode: 'strict',
    block_action: 'block',
    redirect_url: '',
    custom_error_message: '',
    ip_list: '',
    domain_list: '',
  }

  const config = reactive<AccessControlConfig>({ ...defaultConfig })

  const testIP = ref('')
  const testDomain = ref('')
  const testResults = reactive({
    ip: '',
    ipAllowed: false,
    domain: '',
    domainAllowed: false,
  })

  const errors = reactive({
    ip_list: '',
    domain_list: '',
    redirect_url: '',
  })

  const ipModeOptions = computed(() => getAccessControlModeOptions($t))

  const domainModeOptions = computed(() => getDomainControlModeOptions($t))

  const restrictionModeOptions = computed(() => getRestrictionModeOptions($t))

  const blockActionOptions = computed(() => getBlockActionOptions($t))

  const ipModeDescription = computed(() => $t(`settings.accessControl.descriptions.ip.${config.ip_mode}`))

  const domainModeDescription = computed(() => $t(`settings.accessControl.descriptions.domain.${config.domain_mode}`))

  const restrictionDescription = computed(() => $t(`settings.accessControl.descriptions.restriction.${config.restriction_mode}`))

  const blockActionWarning = computed(() => $t(`settings.accessControl.warnings.${config.block_action}`))

  const fetchAccessControlConfig = async () => {
    try {
      loading.value = true
      const result = await getAccessControlConfig()
      if (result.success) {
        Object.assign(config, result.data)
      }
    } catch (_error) {
    } finally {
      loading.value = false
    }
  }

  const saveConfig = async () => {
    if (!validateConfig()) {
      return
    }

    try {
      loading.value = true
      const result = await createOrUpdateAccessControl(config)
      if (result.success) {
        Object.assign(config, result.data)
        toast.success($t('settings.accessControl.messages.updateSuccess'))
      }
    } catch (_error) {
    } finally {
      loading.value = false
    }
  }

  const resetConfig = async () => {
    if (!showConfirm($t('settings.accessControl.confirm.reset'))) {
      return
    }

    try {
      loading.value = true
      const result = await resetAccessControlConfig()
      if (result.success) {
        Object.assign(config, defaultConfig)
        toast.success($t('settings.accessControl.messages.resetSuccess'))
      }
    } catch (_error) {
    } finally {
      loading.value = false
    }
  }

  const validateConfig = () => {
    let isValid = true

    errors.ip_list = ''
    errors.domain_list = ''
    errors.redirect_url = ''

    if (config.ip_list.trim() !== '') {
      const ips = config.ip_list.split(',')
      for (const ip of ips) {
        if (!isValidIP(ip.trim())) {
          errors.ip_list = $t('settings.accessControl.errors.invalidIp', { value: ip.trim() })
          isValid = false
          break
        }
      }
    }

    if (config.domain_list.trim() !== '') {
      const domains = config.domain_list.split(',')
      for (const domain of domains) {
        if (!isValidDomain(domain.trim())) {
          errors.domain_list = $t('settings.accessControl.errors.invalidDomain', {
            value: domain.trim(),
          })
          isValid = false
          break
        }
      }
    }

    if (config.block_action === 'redirect' && (!config.redirect_url || !isValidUrl(config.redirect_url))) {
      errors.redirect_url = $t('settings.accessControl.errors.invalidRedirect')
      isValid = false
    }

    return isValid
  }

  const validateAndSave = () => {
    if (validateConfig()) {
      saveConfig()
    }
  }

  const testAccess = (type: 'ip' | 'domain') => {
    if (type === 'ip') {
      if (!testIP.value) {
        testResults.ip = $t('settings.accessControl.test.ipRequired')
        return
      }

      if (!isValidIP(testIP.value)) {
        testResults.ip = $t('settings.accessControl.test.ipInvalid')
        return
      }

      const ipList = config.ip_list
        .split(',')
        .map((ip) => ip.trim())
        .filter((ip) => ip)
      const ipExists = ipList.includes(testIP.value)

      if (config.ip_mode === 'whitelist') {
        testResults.ipAllowed = ipExists
        testResults.ip = ipExists ? $t('settings.accessControl.test.ipAllowed') : $t('settings.accessControl.test.ipBlocked')
      } else {
        testResults.ipAllowed = !ipExists
        testResults.ip = ipExists ? $t('settings.accessControl.test.ipBlocked') : $t('settings.accessControl.test.ipAllowed')
      }
    } else {
      if (!testDomain.value) {
        testResults.domain = $t('settings.accessControl.test.domainRequired')
        return
      }

      if (!isValidDomain(testDomain.value)) {
        testResults.domain = $t('settings.accessControl.test.domainInvalid')
        return
      }

      const domainList = config.domain_list
        .split(',')
        .map((domain) => domain.trim())
        .filter((domain) => domain)
      const domainExists = domainList.some((domain) => {
        if (domain.startsWith('*.')) {
          const suffix = domain.substring(1)
          return testDomain.value.endsWith(suffix)
        }
        return domain === testDomain.value
      })

      if (config.domain_mode === 'whitelist') {
        testResults.domainAllowed = domainExists
        testResults.domain = domainExists
          ? $t('settings.accessControl.test.domainAllowed')
          : $t('settings.accessControl.test.domainBlocked')
      } else {
        testResults.domainAllowed = !domainExists
        testResults.domain = domainExists
          ? $t('settings.accessControl.test.domainBlocked')
          : $t('settings.accessControl.test.domainAllowed')
      }
    }
  }

  const isValidIP = (ip: string): boolean => {
    if (ip === 'localhost') {
      return true
    }

    const ipv4Regex = /^(\d{1,3})\.(\d{1,3})\.(\d{1,3})\.(\d{1,3})$/
    const cidrRegex = /^(\d{1,3})\.(\d{1,3})\.(\d{1,3})\.(\d{1,3})\/(\d{1,2})$/
    const ipv6Regex =
      /^([0-9a-fA-F]{1,4}:){7}[0-9a-fA-F]{1,4}$|^([0-9a-fA-F]{1,4}:){1,7}:|^([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}$|^([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}$|^([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}$|^([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}$|^([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}$|^[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})$|^:((:[0-9a-fA-F]{1,4}){1,7}|:)$/

    if (ipv4Regex.test(ip)) {
      const parts = ip.match(ipv4Regex)?.slice(1)
      if (!parts) {
        return false
      }

      return parts.every((part) => {
        const num = parseInt(part, 10)
        return num >= 0 && num <= 255
      })
    }

    if (cidrRegex.test(ip)) {
      const parts = ip.match(cidrRegex)?.slice(1)
      if (!parts) {
        return false
      }

      const prefix = parseInt(parts[4], 10)
      if (prefix < 0 || prefix > 32) {
        return false
      }

      return parts.slice(0, 4).every((part) => {
        const num = parseInt(part, 10)
        return num >= 0 && num <= 255
      })
    }

    return ipv6Regex.test(ip)
  }

  const isValidDomain = (domain: string): boolean => {
    if (domain === 'localhost') {
      return true
    }

    if (domain.startsWith('*.')) {
      domain = domain.substring(2)
    }

    const domainRegex = /^(?:[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?\.)+[a-z0-9][a-z0-9-]{0,61}[a-z0-9]$/i
    return domainRegex.test(domain)
  }

  const isValidUrl = (url: string): boolean => {
    try {
      new URL(url)
      return true
    } catch (_error) {
      return false
    }
  }

  onMounted(() => {
    fetchAccessControlConfig()
  })

  defineExpose({
    fetchAccessControlConfig,
  })
</script>

<template>
  <div class="access-control-section">
    <div class="section-header">
      <div class="header-content">
        <div class="title-wrapper">
          <div class="title-icon">
            <i class="fas fa-shield-alt" />
          </div>
          <div class="title-text">
            <h2 class="title">{{ $t('settings.accessControl.title') }}</h2>
            <p class="subtitle">{{ $t('settings.accessControl.description') }}</p>
          </div>
        </div>
        <div class="header-controls">
          <CyberSwitch v-model="config.enabled" @change="saveConfig" />
        </div>
      </div>
    </div>

    <div class="config-content" :class="{ 'config-disabled': !config.enabled }">
      <div class="control-grid">
        <div class="control-card">
          <div class="card-header">
            <div class="card-title">
              <i class="fas fa-network-wired card-icon" />
              <span>{{ $t('settings.accessControl.sections.ip.title') }}</span>
            </div>
          </div>

          <div class="card-content">
            <div class="form-group">
              <label class="form-label">
                {{ $t('settings.accessControl.sections.ip.modeLabel') }}
              </label>
              <CyberDropdown v-model="config.ip_mode" :clearable="false" :options="ipModeOptions" @change="saveConfig" />
              <div class="form-help">{{ ipModeDescription }}</div>
            </div>

            <div class="form-group">
              <label class="form-label">
                {{ $t('settings.accessControl.sections.ip.listLabel') }}
                <span class="form-hint"> ({{ $t('settings.accessControl.sections.ip.listHint') }}) </span>
              </label>
              <CyberInput
                v-model="config.ip_list"
                type="textarea"
                :rows="3"
                :placeholder="$t('settings.accessControl.sections.ip.placeholder')"
                @change="validateAndSave"
              />
              <div v-if="errors.ip_list" class="form-error">
                {{ errors.ip_list }}
              </div>
              <div class="form-help">
                {{ $t('settings.accessControl.sections.ip.helper') }}
              </div>
            </div>
          </div>
        </div>

        <div class="control-card">
          <div class="card-header">
            <div class="card-title">
              <i class="fas fa-globe card-icon" />
              <span>{{ $t('settings.accessControl.sections.domain.title') }}</span>
            </div>
          </div>

          <div class="card-content">
            <div class="form-group">
              <label class="form-label">
                {{ $t('settings.accessControl.sections.domain.modeLabel') }}
              </label>
              <CyberDropdown v-model="config.domain_mode" :clearable="false" :options="domainModeOptions" @change="saveConfig" />
              <div class="form-help">{{ domainModeDescription }}</div>
            </div>

            <div class="form-group">
              <label class="form-label">
                {{ $t('settings.accessControl.sections.domain.listLabel') }}
                <span class="form-hint"> ({{ $t('settings.accessControl.sections.domain.listHint') }}) </span>
              </label>
              <CyberInput
                v-model="config.domain_list"
                type="textarea"
                :rows="3"
                :placeholder="$t('settings.accessControl.sections.domain.placeholder')"
                @change="validateAndSave"
              />
              <div v-if="errors.domain_list" class="form-error">
                {{ errors.domain_list }}
              </div>
              <div class="form-help">
                {{ $t('settings.accessControl.sections.domain.helper') }}
              </div>
            </div>
          </div>
        </div>
      </div>

      <div class="advanced-card">
        <div class="card-header">
          <div class="card-title">
            <i class="fas fa-cogs card-icon" />
            <span>{{ $t('settings.accessControl.advanced.title') }}</span>
          </div>
        </div>

        <div class="card-content">
          <div class="advanced-grid">
            <div class="form-group">
              <label class="form-label">
                {{ $t('settings.accessControl.advanced.restrictionLabel') }}
              </label>
              <CyberDropdown
                v-model="config.restriction_mode"
                :clearable="false"
                :options="restrictionModeOptions"
                @change="saveConfig"
              />
              <div class="form-help">{{ restrictionDescription }}</div>
            </div>

            <div class="form-group">
              <label class="form-label">
                {{ $t('settings.accessControl.advanced.blockActionLabel') }}
              </label>
              <CyberDropdown
                v-model="config.block_action"
                :clearable="false"
                :options="blockActionOptions"
                @change="saveConfig"
              />
              <div class="form-warning">{{ blockActionWarning }}</div>
            </div>
          </div>

          <div v-if="config.block_action === 'redirect'" class="form-group">
            <label class="form-label">
              {{ $t('settings.accessControl.advanced.redirectLabel') }}
            </label>
            <CyberInput
              v-model="config.redirect_url"
              :placeholder="$t('settings.accessControl.advanced.redirectPlaceholder')"
              @change="validateAndSave"
            />
            <div v-if="errors.redirect_url" class="form-error">
              {{ errors.redirect_url }}
            </div>
            <div class="form-help">
              {{ $t('settings.accessControl.advanced.redirectHelper') }}
            </div>
          </div>

          <div v-if="config.block_action === 'block'" class="form-group">
            <label class="form-label">
              {{ $t('settings.accessControl.advanced.customErrorLabel') }}
            </label>
            <CyberInput
              v-model="config.custom_error_message"
              :placeholder="$t('settings.accessControl.advanced.customErrorPlaceholder')"
              @change="validateAndSave"
            />
            <div class="form-help">
              {{ $t('settings.accessControl.advanced.customErrorHelper') }}
            </div>
          </div>
        </div>
      </div>

      <div class="test-card">
        <div class="card-header">
          <div class="card-title">
            <i class="fas fa-flask card-icon" />
            <span>{{ $t('settings.accessControl.test.title') }}</span>
          </div>
        </div>

        <div class="card-content">
          <div class="test-grid">
            <div class="test-group">
              <label class="form-label">
                {{ $t('settings.accessControl.test.ipLabel') }}
              </label>
              <div class="test-input-group">
                <CyberInput v-model="testIP" class="flex-1" :placeholder="$t('settings.accessControl.test.ipPlaceholder')" />
                <CyberButton type="secondary" @click="testAccess('ip')">
                  {{ $t('settings.accessControl.test.button') }}
                </CyberButton>
              </div>
              <div v-if="testResults.ip" class="test-result" :class="testResults.ipAllowed ? 'test-success' : 'test-error'">
                <i :class="testResults.ipAllowed ? 'fas fa-check-circle' : 'fas fa-times-circle'" />
                {{ testResults.ip }}
              </div>
            </div>

            <div class="test-group">
              <label class="form-label">
                {{ $t('settings.accessControl.test.domainLabel') }}
              </label>
              <div class="test-input-group">
                <CyberInput
                  v-model="testDomain"
                  class="flex-1"
                  :placeholder="$t('settings.accessControl.test.domainPlaceholder')"
                />
                <CyberButton type="secondary" @click="testAccess('domain')">
                  {{ $t('settings.accessControl.test.button') }}
                </CyberButton>
              </div>
              <div
                v-if="testResults.domain"
                class="test-result"
                :class="testResults.domainAllowed ? 'test-success' : 'test-error'"
              >
                <i :class="testResults.domainAllowed ? 'fas fa-check-circle' : 'fas fa-times-circle'" />
                {{ testResults.domain }}
              </div>
            </div>
          </div>
        </div>
      </div>

      <div class="action-buttons">
        <CyberButton type="danger" @click="resetConfig">
          <i class="fas fa-undo" />
          {{ $t('settings.accessControl.actions.reset') }}
        </CyberButton>
        <CyberButton type="primary" :loading="loading" @click="saveConfig">
          <i class="fas fa-save" />
          {{ $t('settings.accessControl.actions.save') }}
        </CyberButton>
      </div>
    </div>
  </div>
</template>

<style scoped>
  .access-control-section {
    @apply space-y-6;
  }

  .section-header {
    @apply relative overflow-hidden;
    border-radius: var(--radius-sm);
    background: linear-gradient(135deg, rgba(var(--color-background-800-rgb), 0.9), rgba(var(--color-background-900-rgb), 0.95));
    border: 1px solid var(--color-border-subtle);
    box-shadow:
      0 8px 32px rgba(0, 0, 0, 0.12),
      inset 0 1px 0 rgba(255, 255, 255, 0.05);
  }

  .header-content {
    @apply relative flex items-start justify-between p-6;
    z-index: 1;
  }

  .title-wrapper {
    @apply flex items-start gap-4;
  }

  .title-icon {
    @apply flex h-12 w-12 items-center justify-center;
    border-radius: var(--radius-sm);
    background: linear-gradient(135deg, rgba(var(--color-brand-500-rgb), 0.18), rgba(var(--color-brand-500-rgb), 0.08));
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.24);
    color: var(--color-brand-400);
    font-size: 18px;
    box-shadow:
      0 4px 12px rgba(var(--color-brand-500-rgb), 0.18),
      inset 0 1px 0 rgba(255, 255, 255, 0.1);
  }

  .title-text {
    @apply flex flex-col gap-1;
  }

  .title {
    @apply text-xl font-semibold;
    color: var(--color-content-heading);
    margin: 0;
  }

  .subtitle {
    @apply text-sm leading-relaxed;
    color: var(--color-content-muted);
    margin: 0;
    max-width: 600px;
  }

  .header-controls {
    @apply flex items-center;
  }

  .config-content {
    @apply space-y-6 transition-all duration-300;
  }

  .config-disabled {
    @apply pointer-events-none opacity-50;
  }

  .control-grid {
    @apply grid grid-cols-1 gap-6 lg:grid-cols-2;
  }

  .advanced-grid {
    @apply grid grid-cols-1 gap-4 md:grid-cols-2;
  }

  .test-grid {
    @apply grid grid-cols-1 gap-6 lg:grid-cols-2;
  }

  .control-card,
  .advanced-card,
  .test-card {
    @apply relative overflow-hidden;
    border-radius: var(--radius-sm);
    background: linear-gradient(135deg, rgba(var(--color-background-700-rgb), 0.8), rgba(var(--color-background-800-rgb), 0.9));
    border: 1px solid var(--color-border-subtle);
    box-shadow:
      0 4px 16px rgba(0, 0, 0, 0.08),
      inset 0 1px 0 rgba(255, 255, 255, 0.02);
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  }

  .control-card:hover,
  .advanced-card:hover,
  .test-card:hover {
    border-color: var(--color-border-default);
    box-shadow:
      0 8px 32px rgba(0, 0, 0, 0.12),
      inset 0 1px 0 rgba(255, 255, 255, 0.03);
    transform: translateY(-1px);
  }

  .card-header {
    @apply border-b px-6 py-4;
    border-color: var(--color-border-subtle);
    background: linear-gradient(135deg, rgba(var(--color-background-600-rgb), 0.3), rgba(var(--color-background-700-rgb), 0.5));
  }

  .card-title {
    @apply flex items-center gap-3 text-base font-medium;
    color: var(--color-content-heading);
  }

  .card-icon {
    @apply text-sm;
    color: var(--color-brand-400);
  }

  .card-content {
    @apply space-y-4 p-6;
  }

  .form-group {
    @apply space-y-2;
  }

  .form-label {
    @apply block text-sm font-medium;
    color: var(--color-content);
  }

  .form-hint {
    @apply text-xs;
    color: var(--color-content-subtle);
  }

  .form-help {
    @apply text-xs leading-relaxed;
    color: var(--color-content-muted);
  }

  .form-error {
    @apply text-xs;
    color: var(--color-error-500);
  }

  .form-warning {
    @apply text-xs;
    color: var(--color-warning-500);
  }

  .test-group {
    @apply space-y-3;
  }

  .test-input-group {
    @apply flex gap-3;
  }

  .test-input-group :deep(.cyber-button) {
    white-space: nowrap;
  }

  .test-result {
    @apply flex items-center gap-2 text-xs font-medium;
  }

  .test-success {
    color: var(--color-success-500);
  }

  .test-error {
    color: var(--color-error-500);
  }

  .action-buttons {
    @apply flex items-center justify-end gap-3 p-6;
    border-radius: var(--radius-sm);
    background: linear-gradient(135deg, rgba(var(--color-background-700-rgb), 0.5), rgba(var(--color-background-800-rgb), 0.7));
    border: 1px solid var(--color-border-subtle);
  }

  @media (max-width: 768px) {
    .header-content {
      @apply flex-col gap-4;
    }

    .title-wrapper {
      @apply w-full;
    }

    .header-controls {
      @apply w-full justify-start;
    }

    .control-grid,
    .test-grid {
      @apply grid-cols-1;
    }

    .advanced-grid {
      @apply grid-cols-1;
    }
  }
</style>
