<script setup lang="ts">
  import { computed, onMounted, ref, watch } from 'vue'
  import { useRouter } from 'vue-router'
  import setupApi from '@/api/setup'
  import type { DatabaseTestRequest, InstallRequest, InstallStatus } from '@/api/types/index'
  import { useToast } from '@/components/Toast/useToast'
  import { useAuthStore } from '@/store/auth'
  import { useSettingsStore } from '@/store/settings'
  import { useTexts } from '@/composables/useTexts'

  defineOptions({
    name: 'SetupPage',
  })

  const router = useRouter()
  const toast = useToast()
  const authStore = useAuthStore()
  const settingsStore = useSettingsStore()
  const { $t } = useTexts()

  const isLoading = ref(true)
  const isInstalling = ref(false)
  const isTestingConnection = ref(false)
  const installCompleted = ref(false)
  const showRedisConfig = ref(false)
  const showVectorConfig = ref(false)
  const connectionTestResult = ref<{ success: boolean; message: string } | null>(null)
  const vectorTestResult = ref<{ success: boolean; message: string } | null>(null)
  const isTestingVector = ref(false)
  const vectorMode = ref<'builtin' | 'external'>('builtin')
  const showAdvancedConfig = ref(false)
  const builtinPorts = ref({
    http: 6333,
    grpc: 6334,
  })
  const setupStatus = ref<InstallStatus | null>(null)

  const needsFullConfig = computed(() => {
    return setupStatus.value?.setup_level === 2
  })

  const isDockerMode = computed(() => {
    const deployMode = setupStatus.value?.deploy_mode
    return deployMode === 'docker' || deployMode === 'compose'
  })

  const form = ref<InstallRequest>({
    database: {
      type: 'mysql',
      host: '127.0.0.1',
      port: 3306,
      username: 'root',
      password: '',
      name: 'pixelpunk',
      path: './pixelpunk.db',
    },
    redis: {
      host: 'localhost',
      port: 6379,
      password: '',
      db: 0,
    },
    vector: {
      qdrant_url: 'http://localhost:6333',
      qdrant_timeout: 30,
      use_builtin: true, // 默认使用内置 Qdrant
      http_port: 6333,
      grpc_port: 6334,
    },
    admin_username: 'root',
    admin_password: '',
  })

  const checkInstallStatus = async () => {
    try {
      const result = await setupApi.getInstallStatus()
      if (result.success && result.data) {
        setupStatus.value = result.data
        if (result.data.installed) {
          installCompleted.value = true
        }
      }
    } catch (_error) {
    } finally {
      isLoading.value = false
    }
  }

  const testConnection = async () => {
    if (form.value.database.type === 'mysql') {
      if (!form.value.database.host || !form.value.database.username || !form.value.database.name) {
        connectionTestResult.value = { success: false, message: $t('setup.database.validation.mysqlRequired') }
        return
      }
    } else if (form.value.database.type === 'sqlite') {
      if (!form.value.database.path) {
        connectionTestResult.value = { success: false, message: $t('setup.database.validation.sqliteRequired') }
        return
      }
    }

    isTestingConnection.value = true
    connectionTestResult.value = null

    try {
      const result = await setupApi.testDatabaseConnection(form.value.database as DatabaseTestRequest)
      connectionTestResult.value = {
        success: result.success,
        message: result.success
          ? form.value.database.type === 'sqlite'
            ? $t('setup.database.test.success.sqlite')
            : $t('setup.database.test.success.mysql')
          : result.message || $t('setup.database.validation.connectionFailed'),
      }
    } catch (error: any) {
      connectionTestResult.value = {
        success: false,
        message: `${$t('setup.database.validation.connectionFailed')}: ${error?.response?.data?.message || error?.message || $t('setup.database.validation.unknownError')}`,
      }
    } finally {
      isTestingConnection.value = false
    }
  }

  const testVectorConnection = async () => {
    const testUrl = vectorMode.value === 'builtin' ? 'http://localhost:6333' : form.value.vector?.qdrant_url

    if (!testUrl) {
      vectorTestResult.value = { success: false, message: $t('setup.vector.validation.urlRequired') }
      return
    }

    isTestingVector.value = true
    vectorTestResult.value = null

    try {
      const url = new URL(testUrl)
      if (!url.protocol || !url.hostname) {
        throw new Error($t('setup.vector.validation.invalidUrl'))
      }

      vectorTestResult.value = {
        success: true,
        message:
          vectorMode.value === 'builtin' ? $t('setup.vector.test.success.builtin') : $t('setup.vector.test.success.external'),
      }
    } catch (error: any) {
      vectorTestResult.value = {
        success: false,
        message: `${$t('setup.vector.validation.urlInvalid')}: ${error?.message || $t('setup.vector.validation.formatError')}`,
      }
    } finally {
      isTestingVector.value = false
    }
  }

  const handleSubmit = async () => {
    if (needsFullConfig.value) {
      if (
        form.value.database.type === 'mysql' &&
        (!form.value.database.host || !form.value.database.username || !form.value.database.name)
      ) {
        toast.error($t('setup.admin.validation.mysqlInfoRequired'))
        return
      }
      if (form.value.database.type === 'sqlite' && !form.value.database.path) {
        toast.error($t('setup.admin.validation.sqlitePathRequired'))
        return
      }
    }
    if (!form.value.admin_username || form.value.admin_username.length < 3) {
      toast.error($t('setup.admin.validation.usernameRequired'))
      return
    }
    if (!form.value.admin_password) {
      toast.error($t('setup.admin.validation.passwordRequired'))
      return
    }
    if (form.value.admin_password.length < 6) {
      toast.error($t('setup.admin.validation.passwordTooShort'))
      return
    }

    isInstalling.value = true

    try {
      if (form.value.vector) {
        if (vectorMode.value === 'builtin') {
          form.value.vector.qdrant_url = `http://localhost:${builtinPorts.value.http}`
          form.value.vector.use_builtin = true
          form.value.vector.http_port = builtinPorts.value.http
          form.value.vector.grpc_port = builtinPorts.value.grpc // 使用用户配置的gRPC端口
        } else {
          form.value.vector.use_builtin = false
        }
      }

      const result = await setupApi.installSystem(form.value)
      if (result.success) {
        if (result.data?.token && result.data?.user) {
          authStore.setToken(result.data.token)
          authStore.setUserInfo(result.data.user)
        }
        try {
          await settingsStore.loadGlobalSettings()
        } catch (_e) {}
        installCompleted.value = true
      }
    } catch (_error) {
    } finally {
      isInstalling.value = false
    }
  }

  onMounted(() => checkInstallStatus())

  watch(
    () => form.value.database.type,
    () => {
      connectionTestResult.value = null
    }
  )
</script>

<template>
  <div class="setup-page">
    <div class="cyber-bg">
      <div class="cyber-grid" />
      <div class="scan-line" />
    </div>

    <div :class="['setup-container', { 'docker-mode': isDockerMode }]">
      <div v-if="isLoading" class="state-box loading-state">
        <div class="spinner" />
        <p class="state-text">{{ $t('setup.loading') }}</p>
      </div>

      <div v-else-if="installCompleted" class="success-page">
        <div class="success-header">
          <div class="success-icon-wrapper">
            <div class="success-icon-bg"></div>
            <i class="fas fa-check-circle success-icon" />
          </div>
          <h2 class="success-title">{{ $t('setup.success.title') }}</h2>
          <p class="success-subtitle">{{ $t('setup.success.subtitle') }}</p>
        </div>

        <div class="quick-entry-grid">
          <div class="entry-card" @click="router.push('/')">
            <div class="card-icon home-icon">
              <i class="fas fa-home" />
            </div>
            <h3 class="card-title">
              <i class="fas fa-home title-icon" />
              {{ $t('setup.success.entries.home.title') }}
            </h3>
            <p class="card-desc">{{ $t('setup.success.entries.home.desc') }}</p>
          </div>

          <div class="entry-card" @click="router.push('/upload')">
            <div class="card-icon upload-icon">
              <i class="fas fa-cloud-upload-alt" />
            </div>
            <h3 class="card-title">
              <i class="fas fa-upload title-icon" />
              {{ $t('setup.success.entries.upload.title') }}
            </h3>
            <p class="card-desc">{{ $t('setup.success.entries.upload.desc') }}</p>
          </div>

          <div class="entry-card" @click="router.push('/admin/dashboard')">
            <div class="card-icon dashboard-icon">
              <i class="fas fa-chart-line" />
            </div>
            <h3 class="card-title">
              <i class="fas fa-tachometer-alt title-icon" />
              {{ $t('setup.success.entries.dashboard.title') }}
            </h3>
            <p class="card-desc">{{ $t('setup.success.entries.dashboard.desc') }}</p>
          </div>

          <div class="entry-card" @click="router.push('/admin/settings')">
            <div class="card-icon settings-icon">
              <i class="fas fa-cog" />
            </div>
            <h3 class="card-title">
              <i class="fas fa-sliders-h title-icon" />
              {{ $t('setup.success.entries.settings.title') }}
            </h3>
            <p class="card-desc">{{ $t('setup.success.entries.settings.desc') }}</p>
          </div>
        </div>

        <div class="success-footer">
          <p class="footer-tip">
            <i class="fas fa-info-circle" />
            {{ $t('setup.success.tip') }}
          </p>
        </div>
      </div>

      <div v-else class="install-form">
        <div class="form-header">
          <div class="header-icon">
            <i class="fas fa-rocket" />
          </div>
          <h1 class="logo-title">{{ $t('setup.title') }}</h1>
          <p class="header-subtitle">{{ $t('setup.subtitle') }}</p>
          <div v-if="setupStatus?.deploy_mode && setupStatus.deploy_mode !== 'standalone'" class="deploy-mode-badge">
            <i class="fas fa-docker" />
            <span v-if="setupStatus.deploy_mode === 'docker'">Docker Mode</span>
            <span v-else-if="setupStatus.deploy_mode === 'compose'">Docker Compose Mode</span>
          </div>
        </div>

        <form @submit.prevent="handleSubmit" class="form-body">
          <div v-if="needsFullConfig" class="form-section db-section">
            <div class="section-header">
              <i class="fas fa-database" />
              <h3>{{ $t('setup.database.title') }}</h3>
            </div>

            <div class="db-type-tabs">
              <button
                type="button"
                :class="['db-tab', { active: form.database.type === 'mysql' }]"
                @click="form.database.type = 'mysql'"
              >
                <i class="fas fa-database" />
                {{ $t('setup.database.mysql') }}
              </button>
              <button
                type="button"
                :class="['db-tab', { active: form.database.type === 'sqlite' }]"
                @click="form.database.type = 'sqlite'"
              >
                <i class="fas fa-file-alt" />
                {{ $t('setup.database.sqlite') }}
              </button>
            </div>

            <div class="config-area">
              <div v-show="form.database.type === 'sqlite'" class="config-content">
                <div class="field">
                  <label>{{ $t('setup.database.fields.path') }}</label>
                  <CyberInput v-model="form.database.path" :placeholder="$t('setup.database.placeholders.path')" size="small" />
                </div>
              </div>

              <div v-show="form.database.type === 'mysql'" class="config-content">
                <div class="field-row">
                  <div class="field">
                    <label>{{ $t('setup.database.fields.host') }}</label>
                    <CyberInput v-model="form.database.host" :placeholder="$t('setup.database.placeholders.host')" size="small" />
                  </div>
                  <div class="field field-sm">
                    <label>{{ $t('setup.database.fields.port') }}</label>
                    <CyberInput
                      v-model.number="form.database.port"
                      type="number"
                      :placeholder="$t('setup.database.placeholders.port')"
                      size="small"
                    />
                  </div>
                </div>
                <div class="field-row">
                  <div class="field">
                    <label>{{ $t('setup.database.fields.user') }}</label>
                    <CyberInput
                      v-model="form.database.username"
                      :placeholder="$t('setup.database.placeholders.user')"
                      size="small"
                    />
                  </div>
                  <div class="field">
                    <label>{{ $t('setup.database.fields.password') }}</label>
                    <CyberInput
                      v-model="form.database.password"
                      type="password"
                      :placeholder="$t('setup.database.placeholders.password')"
                      size="small"
                    />
                  </div>
                </div>
                <div class="field">
                  <label>{{ $t('setup.database.fields.name') }}</label>
                  <CyberInput v-model="form.database.name" :placeholder="$t('setup.database.placeholders.name')" size="small" />
                </div>
              </div>
            </div>

            <div v-if="connectionTestResult" :class="['test-result', connectionTestResult.success ? 'success' : 'error']">
              <i :class="['fas', connectionTestResult.success ? 'fa-check-circle' : 'fa-times-circle']" />
              {{ connectionTestResult.message }}
            </div>
          </div>

          <div class="form-section admin-section">
            <div class="section-header">
              <i class="fas fa-user-shield" />
              <h3>{{ $t('setup.admin.title') }}</h3>
            </div>
            <div class="field-row">
              <div class="field">
                <label>{{ $t('setup.admin.fields.username') }}</label>
                <CyberInput v-model="form.admin_username" :placeholder="$t('setup.admin.placeholders.username')" size="small" />
              </div>
              <div class="field">
                <label>{{ $t('setup.admin.fields.password') }}</label>
                <CyberInput
                  v-model="form.admin_password"
                  type="password"
                  :placeholder="$t('setup.admin.placeholders.password')"
                  size="small"
                />
              </div>
            </div>
          </div>

          <div v-if="needsFullConfig" :class="['form-section', 'redis-section', { collapsed: !showRedisConfig }]">
            <div class="section-header clickable" @click="showRedisConfig = !showRedisConfig">
              <i class="fas fa-memory" />
              <h3>{{ $t('setup.redis.title') }}</h3>
              <span class="optional-badge">{{ $t('setup.redis.optional') }}</span>
              <i :class="['fas', 'chevron', showRedisConfig ? 'fa-chevron-up' : 'fa-chevron-down']" />
            </div>

            <transition name="redis-collapse">
              <div v-show="showRedisConfig" class="redis-config">
                <div class="field-row">
                  <div class="field">
                    <label>{{ $t('setup.redis.fields.host') }}</label>
                    <CyberInput v-model="form.redis.host" :placeholder="$t('setup.redis.placeholders.host')" size="small" />
                  </div>
                  <div class="field field-sm">
                    <label>{{ $t('setup.redis.fields.port') }}</label>
                    <CyberInput
                      v-model.number="form.redis.port"
                      type="number"
                      :placeholder="$t('setup.redis.placeholders.port')"
                      size="small"
                    />
                  </div>
                </div>
                <div class="field-row">
                  <div class="field">
                    <label>{{ $t('setup.redis.fields.password') }}</label>
                    <CyberInput
                      v-model="form.redis.password"
                      type="password"
                      :placeholder="$t('setup.redis.placeholders.password')"
                      size="small"
                    />
                  </div>
                  <div class="field field-sm">
                    <label>{{ $t('setup.redis.fields.db') }}</label>
                    <CyberInput
                      v-model.number="form.redis.db"
                      type="number"
                      :placeholder="$t('setup.redis.placeholders.db')"
                      size="small"
                    />
                  </div>
                </div>
              </div>
            </transition>
          </div>

          <div v-if="needsFullConfig" :class="['form-section', 'vector-section', { collapsed: !showVectorConfig }]">
            <div class="section-header clickable" @click="showVectorConfig = !showVectorConfig">
              <i class="fas fa-vector-square" />
              <h3>{{ $t('setup.vector.title') }}</h3>
              <span class="optional-badge">{{ $t('setup.vector.optional') }}</span>
              <i :class="['fas', 'chevron', showVectorConfig ? 'fa-chevron-up' : 'fa-chevron-down']" />
            </div>

            <transition name="vector-collapse">
              <div v-show="showVectorConfig" class="vector-config">
                <div class="config-hint">
                  <i class="fas fa-info-circle" />
                  <span>{{ $t('setup.vector.hint') }}</span>
                </div>

                <div class="vector-mode-selector">
                  <div :class="['mode-option', { active: vectorMode === 'builtin' }]" @click="vectorMode = 'builtin'">
                    <div class="mode-title">
                      <i class="fas fa-box" />
                      {{ $t('setup.vector.mode.builtin.title') }}
                    </div>
                    <div class="mode-desc">{{ $t('setup.vector.mode.builtin.desc') }}</div>
                    <div class="mode-badge">{{ $t('setup.vector.mode.builtin.badge') }}</div>
                  </div>

                  <div :class="['mode-option', { active: vectorMode === 'external' }]" @click="vectorMode = 'external'">
                    <div class="mode-title">
                      <i class="fas fa-network-wired" />
                      {{ $t('setup.vector.mode.external.title') }}
                    </div>
                    <div class="mode-desc">{{ $t('setup.vector.mode.external.desc') }}</div>
                  </div>
                </div>

                <template v-if="vectorMode === 'builtin'">
                  <div class="builtin-advanced">
                    <div
                      :class="['builtin-advanced-header', { expanded: showAdvancedConfig }]"
                      @click="showAdvancedConfig = !showAdvancedConfig"
                    >
                      <i class="fas fa-chevron-right" />
                      <span>{{ $t('setup.vector.advanced.title') }}</span>
                    </div>
                    <transition name="fade">
                      <div v-show="showAdvancedConfig" class="builtin-advanced-content">
                        <div class="port-hint">
                          <i class="fas fa-info-circle" />
                          <span>{{ $t('setup.vector.advanced.hint') }}</span>
                        </div>
                        <div class="field">
                          <label>{{ $t('setup.vector.advanced.httpPort.label') }}</label>
                          <CyberInput
                            v-model.number="builtinPorts.http"
                            type="number"
                            :placeholder="$t('setup.vector.placeholders.httpPort')"
                            size="small"
                          />
                          <p class="port-desc">{{ $t('setup.vector.advanced.httpPort.desc') }}</p>
                        </div>
                        <div class="field">
                          <label>{{ $t('setup.vector.advanced.grpcPort.label') }}</label>
                          <CyberInput
                            v-model.number="builtinPorts.grpc"
                            type="number"
                            :placeholder="$t('setup.vector.placeholders.grpcPort')"
                            size="small"
                          />
                          <p class="port-desc">{{ $t('setup.vector.advanced.grpcPort.desc') }}</p>
                        </div>
                      </div>
                    </transition>
                  </div>
                </template>

                <template v-if="vectorMode === 'external'">
                  <div class="field">
                    <label>{{ $t('setup.vector.fields.url') }}</label>
                    <CyberInput
                      v-model="form.vector!.qdrant_url"
                      :placeholder="$t('setup.vector.placeholders.url')"
                      size="small"
                    />
                  </div>
                </template>

                <div class="field field-sm">
                  <label>{{ $t('setup.vector.fields.timeout') }}</label>
                  <CyberInput
                    v-model.number="form.vector!.qdrant_timeout"
                    type="number"
                    :placeholder="$t('setup.vector.placeholders.timeout')"
                    size="small"
                  />
                </div>

                <div class="field">
                  <button type="button" :disabled="isTestingVector" @click="testVectorConnection" class="vector-test-btn">
                    <i :class="['fas', isTestingVector ? 'fa-spinner fa-spin' : 'fa-plug']" />
                    {{
                      isTestingVector
                        ? $t('setup.vector.test.testing')
                        : vectorMode === 'builtin'
                          ? $t('setup.vector.test.button.builtin')
                          : $t('setup.vector.test.button.external')
                    }}
                  </button>
                </div>

                <div v-if="vectorTestResult" :class="['test-result', vectorTestResult.success ? 'success' : 'error']">
                  <i :class="['fas', vectorTestResult.success ? 'fa-check-circle' : 'fa-times-circle']" />
                  {{ vectorTestResult.message }}
                </div>
              </div>
            </transition>
          </div>

          <div class="form-footer">
            <button
              v-if="needsFullConfig"
              type="button"
              :disabled="isTestingConnection"
              @click="testConnection"
              class="test-btn"
            >
              <i class="fas fa-plug" />
              {{ isTestingConnection ? $t('setup.database.test.testing') : $t('setup.database.test.button') }}
            </button>
            <button type="button" :disabled="isInstalling" @click="handleSubmit" class="install-btn">
              <i class="fas fa-rocket" />
              {{ isInstalling ? $t('setup.actions.installing') : $t('setup.actions.install') }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<style src="./styles.css" scoped></style>
