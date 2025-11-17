<script setup lang="ts">
  import { onMounted, reactive, ref, watch } from 'vue'
  import { get } from '@/utils/network/http'
  import { defaultSettings, type Setting } from '@/api/admin/settings'
  import SettingItem from './SettingItem.vue'
  import { useTexts } from '@/composables/useTexts'

  const { $t } = useTexts()

  const props = defineProps<{
    settings: Setting[]
  }>()

  const emit = defineEmits<{
    (e: 'update', settings: Setting[]): void
  }>()

  /* 从defaultSettings获取上传设置的默认值 */
  const uploadDefaults = defaultSettings.upload.reduce(
    (acc, setting) => {
      acc[setting.key] = setting.value
      return acc
    },
    {} as Record<string, unknown>
  )

  /* 本地设置对象（扁平化）- 使用导入的默认值 */
  const localSettings = reactive({
    allowed_image_formats: [...(uploadDefaults.allowed_image_formats || [])],
    max_file_size: uploadDefaults.max_file_size || 20,
    max_batch_size: uploadDefaults.max_batch_size || 100,
    thumbnail_quality: uploadDefaults.thumbnail_quality || 50,
    thumbnail_max_width: uploadDefaults.thumbnail_max_width || 200,
    thumbnail_max_height: uploadDefaults.thumbnail_max_height || 200,
    preserve_exif: uploadDefaults.preserve_exif ?? true,
    daily_upload_limit: uploadDefaults.daily_upload_limit || 50,
    client_max_concurrent_uploads: uploadDefaults.client_max_concurrent_uploads || 3,
    chunked_upload_enabled: uploadDefaults.chunked_upload_enabled ?? true,
    chunked_threshold: uploadDefaults.chunked_threshold || 20,
    chunk_size: uploadDefaults.chunk_size || 2,
    max_concurrency: uploadDefaults.max_concurrency || 5,
    session_timeout: uploadDefaults.session_timeout || 24,
    retry_count: uploadDefaults.retry_count || 3,
    cleanup_interval: uploadDefaults.cleanup_interval || 60,
    content_detection_enabled: uploadDefaults.content_detection_enabled ?? false,
    sensitive_content_handling: uploadDefaults.sensitive_content_handling || 'mark_only',
    ai_analysis_enabled: uploadDefaults.ai_analysis_enabled ?? false,
    user_allowed_storage_durations: uploadDefaults.user_allowed_storage_durations || ['1h', '3d', '7d', '30d', 'permanent'],
    user_default_storage_duration: uploadDefaults.user_default_storage_duration || 'permanent',
  })

  /* 文件格式选项（动态） */
  const imageFormatOptions = ref<Array<{ label: string; value: string }>>([
    { label: 'JPG', value: 'jpg' },
    { label: 'JPEG', value: 'jpeg' },
    { label: 'PNG', value: 'png' },
    { label: 'GIF', value: 'gif' },
    { label: 'WebP', value: 'webp' },
    { label: 'BMP', value: 'bmp' },
    { label: 'SVG', value: 'svg' },
    { label: 'ICO', value: 'ico' },
    { label: 'APNG', value: 'apng' },
    { label: 'TIFF', value: 'tiff' },
    { label: 'TIF', value: 'tif' },
    { label: 'HEIC', value: 'heic' },
    { label: 'HEIF', value: 'heif' },
  ])

  /* 系统能力（只读展示） */
  const caps = ref<Record<string, unknown> | null>(null)

  /* 将扁平化的设置转换为Setting数组 - 使用defaultSettings中的类型和描述信息 */
  const getSettingsArray = (): Setting[] =>
    defaultSettings.upload.map((defaultSetting) => {
      const settingKey = defaultSetting.key as keyof typeof localSettings
      return {
        key: defaultSetting.key,
        value: localSettings[settingKey],
        type: defaultSetting.type,
        group: 'upload',
        description: defaultSetting.description,
      }
    })

  const applySettings = (settings: Setting[]) => {
    settings.forEach((setting) => {
      const key = setting.key as keyof typeof localSettings

      if (key in localSettings) {
        if (key === 'allowed_image_formats' && Array.isArray(setting.value)) {
          const validFormats = setting.value.filter((format) =>
            imageFormatOptions.value.some((option) => option.value === format)
          )

          localSettings.allowed_image_formats = validFormats
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

  const ensurePermanentSelected = () => {
    if (!localSettings.user_allowed_storage_durations.includes('permanent')) {
      localSettings.user_allowed_storage_durations = [...localSettings.user_allowed_storage_durations, 'permanent']
    }
  }

  watch(
    () => localSettings.user_allowed_storage_durations,
    (newDurations) => {
      ensurePermanentSelected()

      if (newDurations && localSettings.user_default_storage_duration !== 'permanent') {
        if (!newDurations.includes(localSettings.user_default_storage_duration)) {
          localSettings.user_default_storage_duration = newDurations[0] || 'permanent'
        }
      }
    },
    { deep: true }
  )

  const loadCapabilities = async () => {
    try {
      const resp = await get<{ supported_extensions?: string[] }>('/config/upload/capabilities')
      caps.value = resp?.data || null
      const exts: string[] = resp?.data?.supported_extensions || []
      if (Array.isArray(exts) && exts.length > 0) {
        imageFormatOptions.value = exts.map((e) => ({ label: e.toUpperCase(), value: e }))
        if (props.settings && props.settings.length > 0) {
          applySettings(props.settings)
        }
      }
    } catch {}
  }

  const handleAddUserOption = (_option: Record<string, unknown>) => {}

  const handleRemoveUserOption = (value: string) => {
    if (localSettings.user_default_storage_duration === value) {
      localSettings.user_default_storage_duration = 'permanent'
    }
  }

  const handleValidationError = () => {}

  onMounted(() => {
    loadCapabilities()
    if (props.settings && props.settings.length > 0) {
      applySettings(props.settings)
    }
    ensurePermanentSelected()
  })
</script>

<template>
  <div class="space-y-5">
    <div class="space-y-5">
      <h3 class="mb-4 border-b border-subtle pb-2 text-lg font-medium text-content-heading">
        {{ $t('admin.settings.upload.general.title') }}
      </h3>

      <div class="space-y-4">
        <SettingItem
          :label="$t('admin.settings.upload.formats.label')"
          icon="file-image"
          :description="$t('admin.settings.upload.formats.description')"
        >
          <CyberDropdown
            v-model="localSettings.allowed_image_formats"
            :options="imageFormatOptions"
            :placeholder="$t('admin.settings.upload.formats.placeholder')"
            multiple
            style="width: 500px"
          />
        </SettingItem>

        <SettingItem
          :label="$t('admin.settings.upload.sizeLimit.label')"
          icon="weight-hanging"
          :description="$t('admin.settings.upload.sizeLimit.description')"
        >
          <div class="flex items-center gap-3">
            <div class="flex items-center">
              <span class="text-content-content-muted mr-2 whitespace-nowrap text-sm"
                >{{ $t('admin.settings.upload.sizeLimit.singleFile') }}:</span
              >
              <CyberInput
                v-model.number="localSettings.max_file_size"
                type="number"
                :placeholder="$t('admin.settings.upload.sizeLimit.singlePlaceholder')"
                width="120px"
                min="0"
                max="1024"
              >
                <template #unit>MB</template>
              </CyberInput>
            </div>
            <div class="ml-4 flex items-center">
              <span class="text-content-content-muted mr-2 whitespace-nowrap text-sm"
                >{{ $t('admin.settings.upload.multiFile') }}:</span
              >
              <CyberInput
                v-model.number="localSettings.max_batch_size"
                type="number"
                :placeholder="$t('admin.settings.upload.batchSize.placeholder')"
                width="120px"
                min="0"
                max="2048"
              >
                <template #unit>MB</template>
              </CyberInput>
            </div>
          </div>
        </SettingItem>

        <SettingItem
          :label="$t('admin.settings.upload.thumbnailSize.label')"
          icon="arrows-alt"
          :description="$t('admin.settings.upload.thumbnailSize.description')"
        >
          <div class="flex items-center gap-3">
            <div class="flex items-center">
              <span class="text-content-content-muted mr-2 whitespace-nowrap text-sm"
                >{{ $t('admin.settings.upload.thumbnailSize.maxWidth') }}:</span
              >
              <CyberInput
                v-model.number="localSettings.thumbnail_max_width"
                type="number"
                :placeholder="$t('admin.settings.upload.thumbnailSize.widthPlaceholder')"
                width="120px"
                min="100"
                :max="10000"
              >
                <template #unit>px</template>
              </CyberInput>
            </div>
            <div class="ml-4 flex items-center">
              <span class="text-content-content-muted mr-2 whitespace-nowrap text-sm"
                >{{ $t('admin.settings.upload.thumbnailSize.maxHeight') }}:</span
              >
              <CyberInput
                v-model.number="localSettings.thumbnail_max_height"
                type="number"
                :placeholder="$t('admin.settings.upload.thumbnailSize.heightPlaceholder')"
                width="120px"
                min="100"
                :max="10000"
              >
                <template #unit>px</template>
              </CyberInput>
            </div>
          </div>
        </SettingItem>

        <SettingItem
          :label="$t('admin.settings.upload.thumbnailQuality.label')"
          icon="compress"
          :description="$t('admin.settings.upload.thumbnailQuality.description')"
        >
          <div class="flex items-center">
            <CyberSlider v-model="localSettings.thumbnail_quality" :min="0" :max="100" :step="1" width="460px" />
            <span class="text-content-content-muted ml-3 text-sm">{{ localSettings.thumbnail_quality }}%</span>
          </div>
        </SettingItem>

        <SettingItem
          :label="$t('admin.settings.upload.preserveExif.label')"
          icon="info-circle"
          :description="$t('admin.settings.upload.preserveExif.description')"
        >
          <div class="flex items-center">
            <CyberSwitch v-model="localSettings.preserve_exif" />
            <span class="text-content-content-muted ml-3 text-sm">{{
              localSettings.preserve_exif
                ? $t('admin.settings.upload.preserveExif.enabled')
                : $t('admin.settings.upload.preserveExif.disabled')
            }}</span>
          </div>
        </SettingItem>

        <SettingItem
          :label="$t('admin.settings.upload.dailyUploadLimit.label')"
          icon="upload"
          :description="$t('admin.settings.upload.dailyUploadLimit.description')"
        >
          <CyberInput
            v-model.number="localSettings.daily_upload_limit"
            type="number"
            :placeholder="$t('admin.settings.upload.dailyLimit.placeholder')"
            width="500px"
            min="0"
          />
        </SettingItem>

        <SettingItem
          :label="$t('admin.settings.upload.clientMaxConcurrent.label')"
          icon="rocket"
          :description="$t('admin.settings.upload.clientMaxConcurrent.description')"
        >
          <CyberInput
            v-model.number="localSettings.client_max_concurrent_uploads"
            type="number"
            :placeholder="$t('admin.settings.upload.maxConcurrent.placeholder')"
            width="500px"
            min="1"
            max="10"
          />
        </SettingItem>
      </div>
    </div>

    <!-- 分片上传设置 -->
    <div class="space-y-5">
      <h3 class="mb-4 border-b border-subtle pb-2 text-lg font-medium text-content-heading">
        {{ $t('admin.settings.upload.chunkedUpload.title') }}
      </h3>

      <div class="space-y-4">
        <SettingItem
          :label="$t('admin.settings.upload.chunkedUpload.enabled.label')"
          icon="toggle-on"
          :description="$t('admin.settings.upload.chunkedUpload.enabled.description')"
        >
          <div class="flex items-center">
            <CyberSwitch v-model="localSettings.chunked_upload_enabled" />
            <span class="text-content-content-muted ml-3 text-sm">{{
              localSettings.chunked_upload_enabled
                ? $t('admin.settings.upload.chunkedUpload.enabled.on')
                : $t('admin.settings.upload.chunkedUpload.enabled.off')
            }}</span>
          </div>
        </SettingItem>

        <SettingItem
          :label="$t('admin.settings.upload.chunkedUpload.threshold.label')"
          icon="ruler"
          :description="$t('admin.settings.upload.chunkedUpload.threshold.description')"
        >
          <CyberInput
            v-model.number="localSettings.chunked_threshold"
            type="number"
            :placeholder="$t('admin.settings.upload.chunkedThreshold.placeholder')"
            width="500px"
            min="1"
            max="100"
            :disabled="!localSettings.chunked_upload_enabled"
          >
            <template #unit>MB</template>
          </CyberInput>
        </SettingItem>

        <SettingItem
          :label="$t('admin.settings.upload.chunkedUpload.size.label')"
          icon="cut"
          :description="$t('admin.settings.upload.chunkedUpload.size.description')"
        >
          <CyberInput
            v-model.number="localSettings.chunk_size"
            type="number"
            :placeholder="$t('admin.settings.upload.chunkSize.placeholder')"
            width="500px"
            min="1"
            max="10"
            :disabled="!localSettings.chunked_upload_enabled"
          >
            <template #unit>MB</template>
          </CyberInput>
        </SettingItem>

        <SettingItem
          :label="$t('admin.settings.upload.chunkedUpload.maxConcurrency.label')"
          icon="layer-group"
          :description="$t('admin.settings.upload.chunkedUpload.maxConcurrency.description')"
        >
          <div class="flex items-center">
            <CyberSlider
              v-model="localSettings.max_concurrency"
              :min="1"
              :max="10"
              :step="1"
              width="460px"
              :disabled="!localSettings.chunked_upload_enabled"
            />
            <span class="text-content-content-muted ml-3 text-sm"
              >{{ localSettings.max_concurrency }} {{ $t('admin.settings.upload.units.count') }}</span
            >
          </div>
        </SettingItem>

        <SettingItem
          :label="$t('admin.settings.upload.chunkedUpload.retryCount.label')"
          icon="redo"
          :description="$t('admin.settings.upload.chunkedUpload.retryCount.description')"
        >
          <div class="flex items-center">
            <CyberSlider
              v-model="localSettings.retry_count"
              :min="1"
              :max="10"
              :step="1"
              width="460px"
              :disabled="!localSettings.chunked_upload_enabled"
            />
            <span class="text-content-content-muted ml-3 text-sm"
              >{{ localSettings.retry_count }} {{ $t('admin.settings.upload.units.times') }}</span
            >
          </div>
        </SettingItem>

        <SettingItem
          :label="$t('admin.settings.upload.chunkedUpload.sessionTimeout.label')"
          icon="clock"
          :description="$t('admin.settings.upload.chunkedUpload.sessionTimeout.description')"
        >
          <CyberInput
            v-model.number="localSettings.session_timeout"
            type="number"
            :placeholder="$t('admin.settings.upload.sessionTimeout.placeholder')"
            width="500px"
            min="1"
            max="168"
            :disabled="!localSettings.chunked_upload_enabled"
          >
            <template #unit>{{ $t('admin.settings.upload.units.hours') }}</template>
          </CyberInput>
        </SettingItem>

        <SettingItem
          :label="$t('admin.settings.upload.chunkedUpload.cleanupInterval.label')"
          icon="broom"
          :description="$t('admin.settings.upload.chunkedUpload.cleanupInterval.description')"
        >
          <CyberInput
            v-model.number="localSettings.cleanup_interval"
            type="number"
            :placeholder="$t('admin.settings.upload.cleanupInterval.placeholder')"
            width="500px"
            min="10"
            max="1440"
            :disabled="!localSettings.chunked_upload_enabled"
          >
            <template #unit>{{ $t('admin.settings.upload.units.minutes') }}</template>
          </CyberInput>
        </SettingItem>
      </div>
    </div>

    <!-- 存储时长设置 -->
    <div class="space-y-5">
      <h3 class="mb-4 border-b border-subtle pb-2 text-lg font-medium text-content-heading">
        {{ $t('admin.settings.upload.storageDuration.title') }}
      </h3>

      <div class="space-y-4">
        <div class="flex flex-col md:flex-row md:items-center">
          <label class="flex w-full items-center text-sm text-content md:w-48 md:py-2">
            <i class="fas fa-clock mr-2 text-brand-600" />{{ $t('admin.settings.upload.storageDuration.userAllowed') }}
          </label>
          <div class="flex-1">
            <CyberMultiSelector
              v-model="localSettings.user_allowed_storage_durations"
              input-id-prefix="user-duration"
              forced-value="permanent"
              size="sm"
              rounded="sm"
              default-icon="fas fa-clock"
              forced-icon="fas fa-infinity"
              :editable="true"
              :add-text="$t('admin.settings.upload.storageDuration.addNew')"
              :max-options="15"
              :is-guest="false"
              @add-option="handleAddUserOption"
              @remove-option="handleRemoveUserOption"
              @validation-error="handleValidationError"
            />
            <p class="text-content-content-disabled mt-1 text-xs">
              {{ $t('admin.settings.upload.storageDuration.userAllowedHint') }}
            </p>
          </div>
        </div>

        <div class="flex flex-col md:flex-row md:items-center">
          <label class="flex w-full items-center text-sm text-content md:w-48 md:py-2">
            <i class="fas fa-hourglass-half mr-2 text-brand-600" />{{
              $t('admin.settings.upload.storageDuration.defaultDuration')
            }}
          </label>
          <div class="flex-1">
            <div
              class="inline-flex h-[34px] w-[500px] items-center gap-2 rounded-lg border px-3 py-1.5"
              :style="{
                background: 'rgba(var(--color-background-800-rgb), 0.6)',
                borderColor: 'rgba(var(--color-brand-500-rgb), 0.2)',
              }"
            >
              <i class="fas fa-infinity text-sm text-brand-400" />
              <span class="text-content-default text-sm font-medium">{{
                $t('admin.settings.upload.storageDuration.permanent')
              }}</span>
              <i class="fas fa-lock ml-auto text-xs text-content-muted" />
            </div>
            <p class="text-content-content-disabled mt-1 text-xs">
              {{ $t('admin.settings.upload.storageDuration.defaultDurationHint') }}
            </p>
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
</style>
