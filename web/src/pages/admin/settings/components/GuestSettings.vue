<script setup lang="ts">
  import { computed, onMounted, reactive, watch } from 'vue'
  import { defaultSettings, type Setting } from '@/api/admin/settings'
  import { useTexts } from '@/composables/useTexts'

  const { $t } = useTexts()

  const props = defineProps<{
    settings: Setting[]
  }>()

  const emit = defineEmits<{
    (e: 'update', settings: Setting[]): void
  }>()

  /* 访问级别选项 */
  const accessLevelOptions = computed(() => [
    { label: $t('admin.settings.guest.accessPublic'), value: 'public' },
    { label: $t('admin.settings.guest.accessPrivate'), value: 'private' },
  ])

  /* 从defaultSettings获取访客设置的默认值 */
  const guestDefaults = defaultSettings.guest.reduce(
    (acc, setting) => {
      acc[setting.key] = setting.value
      return acc
    },
    {} as Record<string, unknown>
  )

  /* 本地设置对象（扁平化） */
  const localSettings = reactive({
    enable_guest_upload: guestDefaults.enable_guest_upload ?? true,
    guest_daily_limit: guestDefaults.guest_daily_limit ?? 10,
    guest_default_access_level: 'public', // 强制为公开，不允许修改
    guest_allowed_storage_durations: guestDefaults.guest_allowed_storage_durations ?? ['3d', '7d', '30d'],
    guest_default_storage_duration: guestDefaults.guest_default_storage_duration ?? '7d',
    guest_ip_daily_limit: guestDefaults.guest_ip_daily_limit ?? 50,
  })

  /* 游客默认存储时长选项 - 从已选择的存储时长中动态生成 */
  const guestDefaultOptions = computed(() => {
    if (!localSettings.guest_allowed_storage_durations?.length) {
      return []
    }

    return localSettings.guest_allowed_storage_durations.map((duration) => ({
      value: duration,
      label: formatDurationLabel(duration),
    }))
  })

  const formatDurationLabel = (duration: string): string => {
    if (duration === 'permanent') {
      return $t('admin.settings.guest.durationPermanent')
    }
    if (duration.endsWith('m')) {
      return duration.replace('m', $t('admin.settings.guest.durationMinutes'))
    }
    if (duration.endsWith('h')) {
      return duration.replace('h', $t('admin.settings.guest.durationHours'))
    }
    if (duration.endsWith('d')) {
      return duration.replace('d', $t('admin.settings.guest.durationDays'))
    }
    return duration
  }

  const getSettingsArray = (): Setting[] =>
    defaultSettings.guest.map((defaultSetting) => {
      const settingKey = defaultSetting.key as keyof typeof localSettings
      return {
        key: defaultSetting.key,
        value: localSettings[settingKey],
        type: defaultSetting.type,
        group: 'guest',
        description: defaultSetting.description,
      }
    })

  const applySettings = (settings: Setting[]) => {
    settings.forEach((setting) => {
      const key = setting.key as keyof typeof localSettings

      if (key in localSettings) {
        if (key === 'guest_allowed_storage_durations' && Array.isArray(setting.value)) {
          localSettings[key] = setting.value
        } else {
          const settingsRecord = localSettings as Record<string, unknown>
          settingsRecord[key] = setting.value
        }
      }
    })
  }

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

  watch(
    () => localSettings.guest_allowed_storage_durations,
    (newDurations) => {
      if (newDurations && newDurations.length > 0) {
        if (!newDurations.includes(localSettings.guest_default_storage_duration)) {
          localSettings.guest_default_storage_duration = newDurations[0]
        }
      } else {
        localSettings.guest_default_storage_duration = ''
      }
    },
    { deep: true }
  )

  const handleAddGuestOption = (_option: Record<string, unknown>) => {}

  const handleRemoveGuestOption = (value: string) => {
    if (localSettings.guest_default_storage_duration === value) {
      const remaining = localSettings.guest_allowed_storage_durations.filter((d) => d !== value)
      localSettings.guest_default_storage_duration = remaining[0] || ''
    }
  }

  const handleGuestValidationError = () => {}

  onMounted(() => {
    if (props.settings && props.settings.length > 0) {
      applySettings(props.settings)
    }
  })
</script>

<template>
  <div class="space-y-5">
    <div class="space-y-5">
      <div class="space-y-5">
        <div class="flex flex-col md:flex-row md:items-center">
          <label class="flex w-full items-center text-sm text-content md:w-48 md:py-2">
            <i class="fas fa-user-secret mr-2 text-brand-600" />{{ $t('admin.settings.guest.enableUpload') }}
          </label>
          <div class="flex flex-1 items-center">
            <CyberSwitch v-model="localSettings.enable_guest_upload" />
            <span class="ml-3 text-sm text-content-muted">{{
              localSettings.enable_guest_upload
                ? $t('admin.settings.guest.uploadAllowed')
                : $t('admin.settings.guest.uploadDisabled')
            }}</span>
          </div>
        </div>

        <div class="flex flex-col md:flex-row md:items-center">
          <label class="flex w-full items-center text-sm text-content md:w-48 md:py-2">
            <i class="fas fa-calendar-day mr-2 text-brand-600" />{{ $t('admin.settings.guest.dailyLimit') }}
          </label>
          <div class="flex-1">
            <div class="flex items-center">
              <CyberInput
                v-model="localSettings.guest_daily_limit"
                type="number"
                :placeholder="$t('admin.settings.guest.dailyLimitPlaceholder')"
                width="500px"
                :disabled="!localSettings.enable_guest_upload"
              >
                <template #unit>{{ $t('admin.settings.guest.times') }}</template>
              </CyberInput>
            </div>
            <p class="text-content-subtle mt-1 text-xs">{{ $t('admin.settings.guest.dailyLimitDesc') }}</p>
          </div>
        </div>

        <!-- 默认访问级别 - 固定为公开，暂时隐藏选项 -->
        <div v-if="false" class="flex flex-col md:flex-row md:items-center">
          <label class="flex w-full items-center text-sm text-content md:w-48 md:py-2">
            <i class="fas fa-lock mr-2 text-brand-600" />{{ $t('admin.settings.guest.defaultAccessLevel') }}
          </label>
          <div class="flex-1">
            <CyberDropdown
              v-model="localSettings.guest_default_access_level"
              :options="accessLevelOptions"
              :placeholder="$t('admin.settings.guest.selectAccessLevel')"
              class="w-48"
              :disabled="!localSettings.enable_guest_upload"
            >
              <!-- 选项中的图标插槽 -->
              <template #option-icon="{ option }">
                <i v-if="option.value === 'public'" class="fas fa-globe text-sm text-content" />
                <i v-else-if="option.value === 'private'" class="fas fa-eye-slash text-sm text-content" />
                <i v-else class="fas fa-shield-alt text-sm text-content" />
              </template>

              <!-- 选中项的图标插槽 -->
              <template #selected-icon="{ option }">
                <i v-if="option.value === 'public'" class="fas fa-globe text-sm text-content" />
                <i v-else-if="option.value === 'private'" class="fas fa-eye-slash text-sm text-content" />
                <i v-else class="fas fa-shield-alt text-sm text-content" />
              </template>
            </CyberDropdown>
            <p class="text-content-subtle mt-1 text-xs">{{ $t('admin.settings.guest.accessLevelDesc') }}</p>
          </div>
        </div>

        <!-- 访问级别说明 - 新增固定说明 -->
        <div class="flex flex-col md:flex-row md:items-center">
          <label class="flex w-full items-center text-sm text-content md:w-48 md:py-2">
            <i class="fas fa-globe mr-2 text-brand-600" />{{ $t('admin.settings.guest.accessPermission') }}
          </label>
          <div class="flex-1">
            <div class="access-level-badge">
              <i class="fas fa-globe badge-icon" />
              <span class="badge-text">{{ $t('admin.settings.guest.publicAccess') }}</span>
            </div>
            <p class="text-content-subtle mt-1 text-xs">{{ $t('admin.settings.guest.publicAccessDesc') }}</p>
          </div>
        </div>

        <div class="flex flex-col md:flex-row md:items-center">
          <label class="flex w-full items-center text-sm text-content md:w-48 md:py-2">
            <i class="fas fa-clock mr-2 text-brand-600" />{{ $t('admin.settings.guest.allowedDurations') }}
          </label>
          <div class="flex-1">
            <CyberMultiSelector
              v-model="localSettings.guest_allowed_storage_durations"
              :disabled="!localSettings.enable_guest_upload"
              input-id-prefix="guest-duration"
              size="sm"
              rounded="sm"
              default-icon="fas fa-clock"
              :editable="true"
              :add-text="$t('admin.settings.guest.addDuration')"
              :max-options="10"
              :is-guest="true"
              @add-option="handleAddGuestOption"
              @remove-option="handleRemoveGuestOption"
              @validation-error="handleGuestValidationError"
            />
            <p class="text-content-subtle mt-1 text-xs">{{ $t('admin.settings.guest.allowedDurationsDesc') }}</p>
          </div>
        </div>

        <div class="flex flex-col md:flex-row md:items-center">
          <label class="flex w-full items-center text-sm text-content md:w-48 md:py-2">
            <i class="fas fa-hourglass-half mr-2 text-brand-600" />{{ $t('admin.settings.guest.defaultDuration') }}
          </label>
          <div class="flex-1">
            <CyberDropdown
              v-model="localSettings.guest_default_storage_duration"
              :options="guestDefaultOptions"
              :placeholder="$t('admin.settings.guest.selectDefaultDuration')"
              style="width: 500px"
              class="h-[30px]"
              :disabled="!localSettings.enable_guest_upload || localSettings.guest_allowed_storage_durations?.length === 0"
            >
              <!-- 选项中的图标插槽 -->
              <template #option-icon>
                <i class="fas fa-calendar-alt text-sm text-content" />
              </template>

              <!-- 选中项的图标插槽 -->
              <template #selected-icon>
                <i class="fas fa-calendar-alt text-sm text-content" />
              </template>
            </CyberDropdown>
            <p class="text-content-subtle mt-1 text-xs">{{ $t('admin.settings.guest.defaultDurationDesc') }}</p>
          </div>
        </div>

        <div class="flex flex-col md:flex-row md:items-center">
          <label class="flex w-full items-center text-sm text-content md:w-48 md:py-2">
            <i class="fas fa-network-wired mr-2 text-brand-600" />{{ $t('admin.settings.guest.ipDailyLimit') }}
          </label>
          <div class="flex-1">
            <div class="flex items-center">
              <CyberInput
                v-model="localSettings.guest_ip_daily_limit"
                type="number"
                :placeholder="$t('admin.settings.guest.ipDailyLimitPlaceholder')"
                width="500px"
                :disabled="!localSettings.enable_guest_upload"
              >
                <template #unit>{{ $t('admin.settings.guest.times') }}</template>
              </CyberInput>
            </div>
            <p class="text-content-subtle mt-1 text-xs">{{ $t('admin.settings.guest.ipDailyLimitDesc') }}</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped lang="scss">
  .hover:shadow-glow:hover {
    box-shadow: 0 0 15px rgba(var(--color-brand-500-rgb), 0.15);
  }

  .access-level-badge {
    display: inline-flex;
    align-items: center;
    gap: var(--space-sm);
    height: 32px;
    width: 500px;
    padding: var(--space-xs) var(--space-md);
    background: linear-gradient(135deg, rgba(var(--color-brand-500-rgb), 0.08), rgba(var(--color-brand-500-rgb), 0.05));
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.2);
    border-radius: var(--radius-md);
    box-shadow: 0 2px 4px rgba(var(--color-brand-500-rgb), 0.1);
    transition: all var(--transition-fast);

    &:hover {
      background: linear-gradient(135deg, rgba(var(--color-brand-500-rgb), 0.12), rgba(var(--color-brand-500-rgb), 0.08));
      border-color: rgba(var(--color-brand-500-rgb), 0.3);
      box-shadow: 0 4px 8px rgba(var(--color-brand-500-rgb), 0.15);
    }

    .badge-icon {
      font-size: var(--text-sm);
      color: var(--color-brand-400);
      filter: drop-shadow(0 0 4px rgba(var(--color-brand-500-rgb), 0.4));
    }

    .badge-text {
      font-size: var(--text-sm);
      font-weight: var(--font-semibold);
      color: var(--color-content-heading);
      text-shadow: 0 0 6px rgba(var(--color-brand-500-rgb), 0.1);
    }
  }
</style>
