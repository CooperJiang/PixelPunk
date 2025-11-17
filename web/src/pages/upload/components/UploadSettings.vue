<script setup lang="ts">
  import { computed, ref, watch } from 'vue'
  import { useAuthStore } from '@/store/auth'
  import { useGlobalSettings } from '@/composables/useGlobalSettings'
  import { DEFAULT_WATERMARK_CONFIG, type WatermarkConfig as WatermarkConfigType } from '@/components/WatermarkConfig/types'
  import { useTexts } from '@/composables/useTexts'

  const { $t } = useTexts()

  const props = defineProps<{
    folderId?: string
    folderPath?: string // 新增：预设文件夹路径，避免重复查询
    accessLevel?: 'public' | 'private' | 'protected'
    optimize?: boolean
    autoRemove?: boolean
    storageDuration?: string
    watermarkEnabled?: boolean
    watermarkConfig?: WatermarkConfigType
  }>()

  const emit = defineEmits<{
    (e: 'update:folderId', value: string | undefined): void
    (e: 'update:accessLevel', value: 'public' | 'private' | 'protected'): void
    (e: 'update:optimize', value: boolean): void
    (e: 'update:autoRemove', value: boolean): void
    (e: 'update:storageDuration', value: string): void
    (e: 'update:watermarkEnabled', value: boolean): void
    (e: 'update:watermarkConfig', value: WatermarkConfigType): void
    (
      e: 'change',
      settings: {
        folderId?: string
        accessLevel: 'public' | 'private' | 'protected'
        optimize: boolean
        autoRemove: boolean
        storageDuration?: string
        watermarkEnabled: boolean
        watermarkConfig: WatermarkConfigType
      }
    ): void
  }>()

  const authStore = useAuthStore()
  const { settings: globalSettings } = useGlobalSettings()
  const localFolder = ref<string>(props.folderId || '')
  const localAccessLevel = ref<'public' | 'private' | 'protected'>(props.accessLevel || 'private')
  const localOptimize = ref<boolean>(props.optimize !== undefined ? props.optimize : true)
  const localAutoRemove = ref<boolean>(props.autoRemove !== undefined ? props.autoRemove : false)
  const localStorageDuration = ref<string>(props.storageDuration || '')
  const localWatermarkEnabled = ref<boolean>(props.watermarkEnabled || false)
  const watermarkConfig = ref<WatermarkConfigType>(props.watermarkConfig || { ...DEFAULT_WATERMARK_CONFIG })
  const selectedFolderPath = ref<string>('')

  const showWatermarkConfig = ref<boolean>(false)

  watch(
    [globalSettings, () => authStore.isLoggedIn],
    ([newSettings, isLoggedIn]) => {
      if (newSettings && !props.storageDuration) {
        let defaultDuration = ''
        if (isLoggedIn) {
          defaultDuration = newSettings.upload?.user_default_storage_duration || 'permanent'
        } else {
          defaultDuration = newSettings.guest?.guest_default_storage_duration || ''
        }

        if (!localStorageDuration.value || localStorageDuration.value !== defaultDuration) {
          localStorageDuration.value = defaultDuration
          emitChangeEvent()
        }
      }
    },
    { immediate: true }
  )

  const accessLevelOptions = computed(() => [
    {
      label: $t('upload.settings.accessLevel.public'),
      value: 'public',
      description: $t('upload.settings.accessLevel.publicDesc'),
    },
    {
      label: $t('upload.settings.accessLevel.private'),
      value: 'private',
      description: $t('upload.settings.accessLevel.privateDesc'),
    },
    {
      label: $t('upload.settings.accessLevel.protected'),
      value: 'protected',
      description: $t('upload.settings.accessLevel.protectedDesc'),
    },
  ])

  const storageDurationOptions = computed(() => {
    if (authStore.isLoggedIn) {
      const allowedDurations = globalSettings.value?.upload?.user_allowed_storage_durations
      const options: { label: string; value: string }[] = []

      if (Array.isArray(allowedDurations) && allowedDurations.length > 0) {
        allowedDurations.forEach((duration) => {
          if (duration === 'permanent') {
            options.push({
              label: $t('upload.settings.duration.permanent'),
              value: duration,
            })
          } else {
            options.push({
              label: duration,
              value: duration,
            })
          }
        })
      } else {
        options.push({ label: $t('upload.settings.duration.permanent'), value: 'permanent' })
      }

      return options
    }
    const allowedDurations = globalSettings.value?.guest?.guest_allowed_storage_durations
    if (Array.isArray(allowedDurations) && allowedDurations.length > 0) {
      return allowedDurations.map((duration) => ({
        label: duration,
        value: duration,
      }))
    }
    return [
      { label: '3d', value: '3d' },
      { label: '7d', value: '7d' },
      { label: '30d', value: '30d' },
    ]
  })

  watch(
    storageDurationOptions,
    (newOptions) => {
      if (newOptions && newOptions.length > 0 && globalSettings.value) {
        const { isLoggedIn } = authStore
        const defaultDuration = isLoggedIn
          ? globalSettings.value.upload?.user_default_storage_duration || 'permanent'
          : globalSettings.value.guest?.guest_default_storage_duration || ''

        const currentValueValid = localStorageDuration.value && newOptions.some((opt) => opt.value === localStorageDuration.value)
        const shouldSetDefault = !currentValueValid || !localStorageDuration.value

        if (shouldSetDefault && defaultDuration && newOptions.some((opt) => opt.value === defaultDuration)) {
          localStorageDuration.value = defaultDuration
          emitChangeEvent()
        } else if (shouldSetDefault && newOptions.length > 0) {
          localStorageDuration.value = newOptions[0].value
          emitChangeEvent()
        }
      }
    },
    { immediate: true }
  )

  const guestAccessLevelText = computed(
    () => $t('constants.accessLevels.public')
  )

  const guestIsEnabled = computed(() => globalSettings.value?.guest?.enable_guest_upload !== false)

  watch(localFolder, (newValue) => {
    emit('update:folderId', newValue)
    emitChangeEvent()
  })

  watch(localAccessLevel, (newValue) => {
    emit('update:accessLevel', newValue)
    emitChangeEvent()
  })

  watch(localOptimize, (newValue) => {
    emit('update:optimize', newValue)
    emitChangeEvent()
  })

  watch(localAutoRemove, (newValue) => {
    emit('update:autoRemove', newValue)
    emitChangeEvent()
  })

  watch(localStorageDuration, (newValue) => {
    emit('update:storageDuration', newValue)
    emitChangeEvent()
  })

  watch(localWatermarkEnabled, (newValue) => {
    emit('update:watermarkEnabled', newValue)
    emitChangeEvent()
  })

  watch(
    watermarkConfig,
    (newValue) => {
      emit('update:watermarkConfig', newValue)
      emitChangeEvent()
    },
    { deep: true }
  )

  watch(
    () => props.folderId,
    (newValue) => {
      if (newValue !== localFolder.value) {
        localFolder.value = newValue || ''
      }
    }
  )

  watch(
    () => props.accessLevel,
    (newValue) => {
      if (newValue !== localAccessLevel.value) {
        localAccessLevel.value = newValue || 'private'
      }
    }
  )

  watch(
    () => props.optimize,
    (newValue) => {
      if (newValue !== localOptimize.value) {
        localOptimize.value = newValue !== undefined ? newValue : true
      }
    }
  )

  watch(
    () => props.autoRemove,
    (newValue) => {
      if (newValue !== localAutoRemove.value) {
        localAutoRemove.value = newValue !== undefined ? newValue : false
      }
    }
  )

  watch(
    () => props.storageDuration,
    (newValue) => {
      if (newValue !== localStorageDuration.value) {
        localStorageDuration.value = newValue || 'permanent'
      }
    }
  )

  watch(
    () => props.watermarkEnabled,
    (newValue) => {
      if (newValue !== localWatermarkEnabled.value) {
        localWatermarkEnabled.value = newValue || false
      }
    }
  )

  watch(
    () => props.watermarkConfig,
    (newValue) => {
      if (newValue && JSON.stringify(newValue) !== JSON.stringify(watermarkConfig.value)) {
        watermarkConfig.value = { ...newValue }
      }
    },
    { deep: true }
  )

  const handleFolderSelected = (folder: any) => {
    selectedFolderPath.value = folder.path
    emitChangeEvent()
  }

  const emitChangeEvent = () => {
    emit('change', {
      folderId: localFolder.value || undefined,
      accessLevel: localAccessLevel.value,
      optimize: localOptimize.value,
      autoRemove: localAutoRemove.value,
      storageDuration: localStorageDuration.value,
      watermarkEnabled: localWatermarkEnabled.value,
      watermarkConfig: watermarkConfig.value,
    })
  }

  const _getPositionText = (position: string) => {
    const positionMap: Record<string, string> = {
      'top-left': $t('upload.watermark.position.top-left'),
      'top-center': $t('upload.watermark.position.top-center'),
      'top-right': $t('upload.watermark.position.top-right'),
      'middle-left': $t('upload.watermark.position.middle-left'),
      'middle-center': $t('upload.watermark.position.middle-center'),
      'middle-right': $t('upload.watermark.position.middle-right'),
      'bottom-left': $t('upload.watermark.position.bottom-left'),
      'bottom-center': $t('upload.watermark.position.bottom-center'),
      'bottom-right': $t('upload.watermark.position.bottom-right'),
      custom: $t('upload.watermark.position.custom'),
    }
    return positionMap[position] || position
  }

  const handleWatermarkConfigConfirm = (newConfig: WatermarkConfigType) => {
    watermarkConfig.value = { ...newConfig }
    showWatermarkConfig.value = false
  }
</script>

<template>
  <div class="upload-settings">
    <div class="space-y-4">
      <div v-if="authStore.isLoggedIn">
        <label class="text-heading mb-1.5 block text-sm font-medium">{{ $t('upload.settings.folder') }}</label>
        <cyberFolderTree
          v-model="localFolder"
          :preset-path="props.folderPath"
          class="compact-dropdown text-sm"
          @folder-selected="handleFolderSelected"
        />
        <p class="text-muted mt-1 text-xs">{{ $t('upload.settings.folderHint') }}</p>
      </div>

      <div v-else class="rounded bg-background-700 p-4">
        <div class="space-y-2">
          <div class="flex items-center">
            <i class="fas fa-info-circle mr-2 text-sm text-brand-500" />
            <p class="text-sm text-content">{{ $t('upload.settings.guestFolderInfo', { level: guestAccessLevelText }) }}</p>
          </div>
          <div class="flex items-center">
            <i class="fas fa-upload mr-2 text-sm text-brand-500" />
            <p class="text-muted text-sm">{{ $t('upload.settings.guestLimitInfo') }}</p>
          </div>
          <div v-if="guestIsEnabled === false" class="flex items-center">
            <i class="fas fa-exclamation-triangle mr-2 text-sm text-error-400" />
            <p class="text-sm text-error-400">{{ $t('upload.settings.guestDisabled') }}</p>
          </div>
        </div>
      </div>

      <div v-if="authStore.isLoggedIn">
        <label class="text-heading mb-1.5 block text-sm font-medium">{{ $t('upload.settings.accessLevelLabel') }}</label>
        <cyberDropdown
          v-model="localAccessLevel"
          :options="accessLevelOptions"
          :placeholder="$t('upload.settings.accessLevelPlaceholder')"
          class="compact-dropdown text-sm"
        />
        <div class="text-muted mt-1 rounded bg-background-600 p-2 text-xs">
          <div v-if="localAccessLevel === 'public'" class="flex items-start">
            <i class="fas fa-info-circle mr-1.5 mt-0.5 text-xs text-brand-500" />
            <span>{{ $t('upload.settings.accessLevel.publicDesc') }}</span>
          </div>
          <div v-else-if="localAccessLevel === 'private'" class="flex items-start">
            <i class="fas fa-info-circle mr-1.5 mt-0.5 text-xs text-warning-400" />
            <span>{{ $t('upload.settings.accessLevel.privateDesc') }}</span>
          </div>
          <div v-else-if="localAccessLevel === 'protected'" class="flex items-start">
            <i class="fas fa-shield-alt mr-1.5 mt-0.5 text-xs text-error-400" />
            <span>{{ $t('upload.settings.accessLevel.protectedDesc') }}</span>
          </div>
        </div>
      </div>

      <div v-if="!authStore.isLoggedIn">
        <label class="text-heading mb-1.5 block text-sm font-medium">
          {{ $t('upload.settings.storageDurationLabel') }}
          <span class="text-error-400">{{ $t('upload.settings.required') }}</span>
        </label>
        <cyberDropdown
          v-model="localStorageDuration"
          :options="storageDurationOptions"
          :placeholder="$t('upload.settings.storageDurationPlaceholder')"
          class="compact-dropdown text-sm"
        >
          <template #option-icon>
            <i class="fas fa-calendar-alt text-sm text-brand-500" />
          </template>

          <template #selected-icon>
            <i class="fas fa-calendar-alt text-sm text-brand-500" />
          </template>
        </cyberDropdown>
        <p class="text-muted mt-1 text-xs">{{ $t('upload.settings.storageDurationHint') }}</p>
      </div>

      <div class="pt-2">
        <cyberCheckbox v-model="localOptimize">{{ $t('upload.settings.autoOptimize') }}</cyberCheckbox>
        <p class="text-muted ml-6 mt-1 text-xs">{{ $t('upload.settings.autoOptimizeHint') }}</p>
      </div>

      <div class="pt-2">
        <cyberCheckbox v-model="localAutoRemove">{{ $t('upload.settings.autoRemove') }}</cyberCheckbox>
        <p class="text-muted ml-6 mt-1 text-xs">{{ $t('upload.settings.autoRemoveHint') }}</p>
      </div>
    </div>

    <cyberWatermarkConfig
      v-model:visible="showWatermarkConfig"
      :config="watermarkConfig"
      @confirm="handleWatermarkConfigConfirm"
    />
  </div>
</template>

<style scoped>
  :deep(.compact-dropdown) {
    --dropdown-height: 32px;
    --item-padding: 6px 12px;
    font-size: 0.875rem;
  }

  :deep(.compact-dropdown .dropdown-trigger) {
    height: var(--dropdown-height);
    min-height: var(--dropdown-height);
  }

  :deep(.compact-dropdown .dropdown-item) {
    padding: var(--item-padding);
  }

  :deep(.compact-dropdown .selected-item) {
    padding: var(--item-padding);
  }
</style>
