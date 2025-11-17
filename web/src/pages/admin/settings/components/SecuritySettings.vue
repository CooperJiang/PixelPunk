<script setup lang="ts">
  import { onMounted, reactive, ref, watch } from 'vue'
  import { defaultSettings, type Setting } from '@/api/admin/settings'
  import { useTexts } from '@/composables/useTexts'

  const { $t } = useTexts()

  const props = defineProps<{
    settings: Setting[]
  }>()

  const emit = defineEmits<{
    (e: 'update', settings: Setting[]): void
  }>()

  /* 从defaultSettings获取安全设置的默认值 */
  const securityDefaults = defaultSettings.security.reduce(
    (acc, setting) => {
      acc[setting.key] = setting.value
      return acc
    },
    {} as Record<string, unknown>
  )

  /* UI控制变量 */
  const showSecret = ref(false)
  const showAccessWarning = ref(true)
  const showStorageWarning = ref(false)

  /* 本地设置对象（扁平化）- 使用导入的默认值 */
  const localSettings = reactive({
    max_login_attempts: securityDefaults.max_login_attempts || 5,
    account_lockout_minutes: securityDefaults.account_lockout_minutes || 30,
    login_expire_hours: securityDefaults.login_expire_hours || 60,
    ip_whitelist: securityDefaults.ip_whitelist || '',
    ip_blacklist: securityDefaults.ip_blacklist || '',
    domain_whitelist: securityDefaults.domain_whitelist || '',
    domain_blacklist: securityDefaults.domain_blacklist || '',
    jwt_secret: securityDefaults.jwt_secret || '',
    hide_remote_url: securityDefaults.hide_remote_url !== undefined ? securityDefaults.hide_remote_url : true,
  })

  /* IP白名单和黑名单 */
  const whitelistIPs = ref<string[]>(localSettings.ip_whitelist ? localSettings.ip_whitelist.split(',').filter((ip) => ip) : [])
  const blacklistIPs = ref<string[]>(localSettings.ip_blacklist ? localSettings.ip_blacklist.split(',').filter((ip) => ip) : [])
  const newWhitelistIP = ref('')
  const newBlacklistIP = ref('')

  /* 域名白名单和黑名单 */
  const whitelistDomains = ref<string[]>(
    localSettings.domain_whitelist ? localSettings.domain_whitelist.split(',').filter((domain) => domain) : []
  )
  const blacklistDomains = ref<string[]>(
    localSettings.domain_blacklist ? localSettings.domain_blacklist.split(',').filter((domain) => domain) : []
  )
  const newWhitelistDomain = ref('')
  const newBlacklistDomain = ref('')

  /* IP地址验证函数 */
  const isValidIP = (ip: string): boolean => {
    if (!ip) {
      return false
    }

    const ipv4Regex =
      /^(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$/

    const cidrRegex =
      /^(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\/(3[0-2]|[1-2][0-9]|[0-9])$/

    const ipv6Regex =
      /^([0-9a-fA-F]{1,4}:){7}[0-9a-fA-F]{1,4}$|^([0-9a-fA-F]{1,4}:){1,7}:|^([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}$|^([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}$|^([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}$|^([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}$|^([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}$|^[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})$|^:((:[0-9a-fA-F]{1,4}){1,7}|:)$/

    return ipv4Regex.test(ip) || cidrRegex.test(ip) || ipv6Regex.test(ip)
  }

  const isValidDomain = (domain: string): boolean => {
    if (!domain) {
      return false
    }

    if (domain.startsWith('*.')) {
      domain = domain.substring(2)
    }

    const domainRegex = /^(?:[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?\.)+[a-z0-9][a-z0-9-]{0,61}[a-z0-9]$/i

    return domainRegex.test(domain)
  }

  const getSettingsArray = (): Setting[] => {
    updateIPLists()
    updateDomainLists()

    return defaultSettings.security.map((defaultSetting) => {
      const settingKey = defaultSetting.key as keyof typeof localSettings
      return {
        key: defaultSetting.key,
        value: localSettings[settingKey],
        type: defaultSetting.type,
        group: 'security',
        description: defaultSetting.description,
      }
    })
  }

  const applySettings = (settings: Setting[]) => {
    settings.forEach((setting) => {
      const key = setting.key as keyof typeof localSettings

      if (key in localSettings) {
        if (key === 'ip_whitelist' && typeof setting.value === 'string') {
          localSettings.ip_whitelist = setting.value
          whitelistIPs.value = setting.value ? setting.value.split(',').filter((ip) => ip) : []
        } else if (key === 'ip_blacklist' && typeof setting.value === 'string') {
          localSettings.ip_blacklist = setting.value
          blacklistIPs.value = setting.value ? setting.value.split(',').filter((ip) => ip) : []
        } else if (key === 'domain_whitelist' && typeof setting.value === 'string') {
          localSettings.domain_whitelist = setting.value
          whitelistDomains.value = setting.value ? setting.value.split(',').filter((domain) => domain) : []
        } else if (key === 'domain_blacklist' && typeof setting.value === 'string') {
          localSettings.domain_blacklist = setting.value
          blacklistDomains.value = setting.value ? setting.value.split(',').filter((domain) => domain) : []
        } else {
          const settingsRecord = localSettings as Record<string, unknown>
          settingsRecord[key] = setting.value
        }
      }
    })
  }

  const addWhitelistIP = () => {
    if (newWhitelistIP.value && isValidIP(newWhitelistIP.value)) {
      if (!whitelistIPs.value.includes(newWhitelistIP.value)) {
        whitelistIPs.value.push(newWhitelistIP.value)
        newWhitelistIP.value = ''
        updateIPLists()
      }
    }
  }

  const removeWhitelistIP = (index: number) => {
    whitelistIPs.value.splice(index, 1)
    updateIPLists()
  }

  const addBlacklistIP = () => {
    if (newBlacklistIP.value && isValidIP(newBlacklistIP.value)) {
      if (!blacklistIPs.value.includes(newBlacklistIP.value)) {
        blacklistIPs.value.push(newBlacklistIP.value)
        newBlacklistIP.value = ''
        updateIPLists()
      }
    }
  }

  const removeBlacklistIP = (index: number) => {
    blacklistIPs.value.splice(index, 1)
    updateIPLists()
  }

  const addWhitelistDomain = () => {
    if (newWhitelistDomain.value && isValidDomain(newWhitelistDomain.value)) {
      if (!whitelistDomains.value.includes(newWhitelistDomain.value)) {
        whitelistDomains.value.push(newWhitelistDomain.value)
        newWhitelistDomain.value = ''
        updateDomainLists()
      }
    }
  }

  const removeWhitelistDomain = (index: number) => {
    whitelistDomains.value.splice(index, 1)
    updateDomainLists()
  }

  const addBlacklistDomain = () => {
    if (newBlacklistDomain.value && isValidDomain(newBlacklistDomain.value)) {
      if (!blacklistDomains.value.includes(newBlacklistDomain.value)) {
        blacklistDomains.value.push(newBlacklistDomain.value)
        newBlacklistDomain.value = ''
        updateDomainLists()
      }
    }
  }

  const removeBlacklistDomain = (index: number) => {
    blacklistDomains.value.splice(index, 1)
    updateDomainLists()
  }

  const updateIPLists = () => {
    localSettings.ip_whitelist = whitelistIPs.value.filter((ip) => ip.trim()).join(',')
    localSettings.ip_blacklist = blacklistIPs.value.filter((ip) => ip.trim()).join(',')
  }

  const updateDomainLists = () => {
    localSettings.domain_whitelist = whitelistDomains.value.filter((domain) => domain.trim()).join(',')
    localSettings.domain_blacklist = blacklistDomains.value.filter((domain) => domain.trim()).join(',')
  }

  watch(
    [whitelistIPs, blacklistIPs],
    () => {
      updateIPLists()
    },
    { deep: true }
  )

  watch(
    [whitelistDomains, blacklistDomains],
    () => {
      updateDomainLists()
    },
    { deep: true }
  )

  watch(
    () => props.settings,
    (newSettings) => {
      if (newSettings && newSettings.length > 0) {
        applySettings(newSettings)
      }
    },
    { deep: true, immediate: true }
  )

  watch(
    localSettings,
    () => {
      emit('update', getSettingsArray())
    },
    { deep: true }
  )

  onMounted(() => {
    if (props.settings && props.settings.length > 0) {
      applySettings(props.settings)
    }
  })
</script>

<template>
  <div class="space-y-5">
    <div class="space-y-5">
      <h3 class="mb-3 border-b border-default pb-2 text-lg font-medium text-content-heading">
        {{ $t('admin.settings.security.loginSecurity.title') }}
      </h3>

      <div class="grid grid-cols-1 gap-4 md:grid-cols-2">
        <div class="flex flex-col space-y-1">
          <label class="text-sm text-content-muted">{{ $t('admin.settings.security.loginSecurity.maxAttempts') }}</label>
          <div class="flex items-center">
            <CyberInput v-model="localSettings.max_login_attempts" type="number" placeholder="5" class="w-24" />
          </div>
          <span class="text-xs text-content-muted">{{ $t('admin.settings.security.loginSecurity.maxAttemptsHint') }}</span>
        </div>

        <div class="flex flex-col space-y-1">
          <label class="text-sm text-content-muted">{{ $t('admin.settings.security.loginSecurity.lockoutTime') }}</label>
          <div class="flex items-center">
            <CyberInput v-model="localSettings.account_lockout_minutes" type="number" placeholder="30" class="w-24" />
          </div>
          <span class="text-xs text-content-muted">{{ $t('admin.settings.security.loginSecurity.lockoutTimeHint') }}</span>
        </div>

        <div class="flex flex-col space-y-1">
          <label class="text-sm text-content-muted">{{ $t('admin.settings.security.loginSecurity.expireTime') }}</label>
          <div class="flex items-center">
            <CyberInput v-model="localSettings.login_expire_hours" type="number" placeholder="60" class="w-24" />
          </div>
          <span class="text-xs text-content-muted">{{ $t('admin.settings.security.loginSecurity.expireTimeHint') }}</span>
        </div>

        <div class="flex flex-col space-y-1">
          <label class="text-sm text-content-muted">{{ $t('admin.settings.security.loginSecurity.jwtSecret') }}</label>
          <div class="flex items-center">
            <CyberInput
              v-model="localSettings.jwt_secret"
              :type="showSecret ? 'text' : 'password'"
              :placeholder="$t('admin.settings.security.loginSecurity.jwtSecretPlaceholder')"
              class="w-full"
            />
          </div>
          <span class="text-xs text-content-muted">{{ $t('admin.settings.security.loginSecurity.jwtSecretHint') }}</span>
        </div>
      </div>
    </div>

    <div class="rounded-xl border border-default bg-background-700 p-5">
      <div class="mb-3 flex items-center justify-between border-b border-default pb-2">
        <h3 class="text-lg font-medium text-content-heading">{{ $t('admin.settings.security.storage.title') }}</h3>

        <button
          class="flex items-center text-sm text-info-400 hover:text-info-300"
          @click="showStorageWarning = !showStorageWarning"
        >
          <i class="fas fa-info-circle mr-1" />
          <span>{{ $t('admin.settings.security.storage.viewInfo') }}</span>
        </button>
      </div>

      <div v-if="showStorageWarning" class="mb-4 rounded-md border border-default bg-background-700 p-2.5">
        <div class="flex items-start">
          <i class="fas fa-cloud mr-2 mt-0.5 text-info-400" />
          <div class="text-xs text-content">
            <p class="mb-1 font-medium text-info-300">{{ $t('admin.settings.security.storage.infoTitle') }}</p>
            <div class="space-y-1 text-content-muted">
              <p>
                <span class="text-success-400">✓ {{ $t('admin.settings.security.storage.secureLabel') }}：</span
                >{{ $t('admin.settings.security.storage.secureDesc') }}
              </p>
              <p>
                <span class="text-warning-400">! {{ $t('admin.settings.security.storage.openLabel') }}：</span
                >{{ $t('admin.settings.security.storage.openDesc') }}
              </p>
              <p>
                <span class="text-info-400">ℹ {{ $t('admin.settings.security.storage.suggestLabel') }}：</span
                >{{ $t('admin.settings.security.storage.suggestDesc') }}
              </p>
            </div>
          </div>
        </div>
      </div>

      <div class="flex items-center justify-between">
        <div>
          <div class="flex items-center text-sm text-content">
            <i class="fas fa-shield-alt mr-2 text-brand-600" />
            {{ $t('admin.settings.security.storage.hideRemoteUrl') }}
          </div>
          <div class="mt-1 text-xs text-content-muted">
            {{ $t('admin.settings.security.storage.hideRemoteUrlHint') }}
          </div>
        </div>
        <div class="flex items-center">
          <span class="mr-2 text-xs" :class="localSettings.hide_remote_url ? 'text-success-400' : 'text-warning-400'">
            {{
              localSettings.hide_remote_url
                ? $t('admin.settings.security.storage.globalSecure')
                : $t('admin.settings.security.storage.byChannel')
            }}
          </span>
          <CyberSwitch v-model="localSettings.hide_remote_url" />
        </div>
      </div>
    </div>

    <div class="rounded-xl border border-default bg-background-700 p-5">
      <div class="mb-3 flex items-center justify-between border-b border-default pb-2">
        <h3 class="text-lg font-medium text-content-heading">{{ $t('admin.settings.security.accessControl.title') }}</h3>

        <button
          class="flex items-center text-sm text-warning-400 hover:text-warning-300"
          @click="showAccessWarning = !showAccessWarning"
        >
          <i class="fas fa-exclamation-triangle mr-1" />
          <span>{{ $t('admin.settings.security.accessControl.viewWarning') }}</span>
        </button>
      </div>

      <div v-if="showAccessWarning" class="mb-4 rounded-md border border-default bg-background-700 p-2.5">
        <div class="flex items-start">
          <i class="fas fa-exclamation-triangle mr-2 mt-0.5 text-warning-400" />
          <div class="text-xs text-content">
            <p class="font-medium">
              {{ $t('admin.settings.security.accessControl.warningTitle')
              }}<span class="text-warning-400">{{ $t('admin.settings.security.accessControl.allUsers') }}</span
              >{{ $t('admin.settings.security.accessControl.warningEffect') }}
            </p>
            <p class="text-content-muted">{{ $t('admin.settings.security.accessControl.warningDesc') }}</p>
          </div>
        </div>
      </div>

      <div class="mb-4 grid grid-cols-1 gap-4 md:grid-cols-2">
        <div class="flex flex-col">
          <label class="mb-1.5 text-sm text-content-muted">{{ $t('admin.settings.security.accessControl.ipWhitelist') }}</label>
          <div class="rounded-lg border border-default bg-background-700 p-2.5">
            <div class="mb-2 flex min-h-[30px] flex-wrap gap-1.5">
              <div
                v-for="(ip, index) in whitelistIPs"
                :key="index"
                class="flex items-center rounded border border-brand-500 bg-background-600 px-1.5 py-0.5 text-xs text-content"
              >
                {{ ip }}
                <button class="ml-1.5 text-content-muted hover:text-error-400" @click="() => removeWhitelistIP(index)">
                  <i class="fas fa-times text-xs" />
                </button>
              </div>
              <div v-if="whitelistIPs.length === 0" class="text-xs italic text-content-muted">
                {{ $t('admin.settings.security.accessControl.ipWhitelistEmpty') }}
              </div>
            </div>
            <div class="flex items-center border-t border-default pt-2.5">
              <CyberInput
                v-model="newWhitelistIP"
                :placeholder="$t('admin.settings.security.accessControl.ipPlaceholder')"
                class="mr-2 flex-1 text-xs"
                @keyup.enter="addWhitelistIP"
              />
              <CyberButton
                type="secondary"
                :disabled="!isValidIP(newWhitelistIP) && newWhitelistIP !== ''"
                @click="addWhitelistIP"
              >
                <i class="fas fa-plus" />
              </CyberButton>
            </div>
            <div v-if="newWhitelistIP && !isValidIP(newWhitelistIP)" class="mt-1 text-xs text-error-400">
              {{ $t('admin.settings.security.accessControl.invalidIP') }}
            </div>
          </div>
        </div>

        <div class="flex flex-col">
          <label class="mb-1.5 text-sm text-content-muted">{{ $t('admin.settings.security.accessControl.ipBlacklist') }}</label>
          <div class="rounded-lg border border-default bg-background-700 p-2.5">
            <div class="mb-2 flex min-h-[30px] flex-wrap gap-1.5">
              <div
                v-for="(ip, index) in blacklistIPs"
                :key="index"
                class="flex items-center rounded border border-brand-500 bg-background-600 px-1.5 py-0.5 text-xs text-content"
              >
                {{ ip }}
                <button class="ml-1.5 text-content-muted hover:text-error-400" @click="() => removeBlacklistIP(index)">
                  <i class="fas fa-times text-xs" />
                </button>
              </div>
              <div v-if="blacklistIPs.length === 0" class="text-xs italic text-content-muted">
                {{ $t('admin.settings.security.accessControl.ipBlacklistEmpty') }}
              </div>
            </div>
            <div class="flex items-center border-t border-default pt-2.5">
              <CyberInput
                v-model="newBlacklistIP"
                :placeholder="$t('admin.settings.security.accessControl.ipPlaceholder')"
                class="mr-2 flex-1 text-xs"
                @keyup.enter="addBlacklistIP"
              />
              <CyberButton
                type="secondary"
                :disabled="!isValidIP(newBlacklistIP) && newBlacklistIP !== ''"
                @click="addBlacklistIP"
              >
                <i class="fas fa-plus" />
              </CyberButton>
            </div>
            <div v-if="newBlacklistIP && !isValidIP(newBlacklistIP)" class="mt-1 text-xs text-error-400">
              {{ $t('admin.settings.security.accessControl.invalidIP') }}
            </div>
          </div>
        </div>
      </div>

      <div class="grid grid-cols-1 gap-4 md:grid-cols-2">
        <div class="flex flex-col">
          <label class="mb-1.5 text-sm text-content-muted">{{
            $t('admin.settings.security.accessControl.domainWhitelist')
          }}</label>
          <div class="rounded-lg border border-default bg-background-700 p-2.5">
            <div class="mb-2 flex min-h-[30px] flex-wrap gap-1.5">
              <div
                v-for="(domain, index) in whitelistDomains"
                :key="index"
                class="flex items-center rounded border border-brand-500 bg-background-600 px-1.5 py-0.5 text-xs text-content"
              >
                {{ domain }}
                <button class="ml-1.5 text-content-muted hover:text-error-400" @click="() => removeWhitelistDomain(index)">
                  <i class="fas fa-times text-xs" />
                </button>
              </div>
              <div v-if="whitelistDomains.length === 0" class="text-xs italic text-content-muted">
                {{ $t('admin.settings.security.accessControl.domainWhitelistEmpty') }}
              </div>
            </div>
            <div class="flex items-center border-t border-default pt-2.5">
              <CyberInput
                v-model="newWhitelistDomain"
                :placeholder="$t('admin.settings.security.accessControl.domainPlaceholder')"
                class="mr-2 flex-1 text-xs"
                @keyup.enter="addWhitelistDomain"
              />
              <CyberButton
                type="secondary"
                :disabled="!isValidDomain(newWhitelistDomain) && newWhitelistDomain !== ''"
                @click="addWhitelistDomain"
              >
                <i class="fas fa-plus" />
              </CyberButton>
            </div>
            <div v-if="newWhitelistDomain && !isValidDomain(newWhitelistDomain)" class="mt-1 text-xs text-error-400">
              {{ $t('admin.settings.security.accessControl.invalidDomain') }}
            </div>
          </div>
        </div>

        <div class="flex flex-col">
          <label class="mb-1.5 text-sm text-content-muted">{{
            $t('admin.settings.security.accessControl.domainBlacklist')
          }}</label>
          <div class="rounded-lg border border-default bg-background-700 p-2.5">
            <div class="mb-2 flex min-h-[30px] flex-wrap gap-1.5">
              <div
                v-for="(domain, index) in blacklistDomains"
                :key="index"
                class="flex items-center rounded border border-brand-500 bg-background-600 px-1.5 py-0.5 text-xs text-content"
              >
                {{ domain }}
                <button class="ml-1.5 text-content-muted hover:text-error-400" @click="() => removeBlacklistDomain(index)">
                  <i class="fas fa-times text-xs" />
                </button>
              </div>
              <div v-if="blacklistDomains.length === 0" class="text-xs italic text-content-muted">
                {{ $t('admin.settings.security.accessControl.domainBlacklistEmpty') }}
              </div>
            </div>
            <div class="flex items-center border-t border-default pt-2.5">
              <CyberInput
                v-model="newBlacklistDomain"
                :placeholder="$t('admin.settings.security.accessControl.domainPlaceholder')"
                class="mr-2 flex-1 text-xs"
                @keyup.enter="addBlacklistDomain"
              />
              <CyberButton
                type="secondary"
                :disabled="!isValidDomain(newBlacklistDomain) && newBlacklistDomain !== ''"
                @click="addBlacklistDomain"
              >
                <i class="fas fa-plus" />
              </CyberButton>
            </div>
            <div v-if="newBlacklistDomain && !isValidDomain(newBlacklistDomain)" class="mt-1 text-xs text-error-400">
              {{ $t('admin.settings.security.accessControl.invalidDomain') }}
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
