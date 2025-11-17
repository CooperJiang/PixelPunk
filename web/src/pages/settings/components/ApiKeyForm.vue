<script setup lang="ts">
  import { computed, defineEmits, defineProps, onMounted, reactive, ref, watch } from 'vue'
  import type { CreateApiKeyRequest, UpdateApiKeyRequest } from '@/api/types/index'
  import { getFileSizeOptions, getStorageSizeOptions } from '@/constants/storageOptions'
  import { useTexts } from '@/composables/useTexts'
  import { get } from '@/utils/network/http'

  interface ApiKeyFormData extends CreateApiKeyRequest, UpdateApiKeyRequest {
    name: string
    storage_limit: number
    single_file_limit: number
    upload_count_limit: number
    allowed_types: string[]
    folder_id: string
    expires_in_days?: number
  }

  const props = defineProps<{
    initialData?: Partial<ApiKeyFormData>
    isEditing: boolean
  }>()

  const emit = defineEmits(['submit', 'form-changed'])

  const { $t } = useTexts()

  /* 动态文件类型选项（从后端获取） */
  const fileTypeOptions = ref<Array<{ label: string; value: string }>>([
    { label: $t('settings.apiKeys.text.allFileTypes'), value: '*' },
  ])

  const storageSizeOptions = getStorageSizeOptions($t)
  const fileSizeOptions = getFileSizeOptions($t)

  const form = reactive<ApiKeyFormData>({
    name: '',
    storage_limit: 0, // 默认不限量，使用0
    single_file_limit: 10 * 1024 * 1024, // 默认10MB
    upload_count_limit: 0, // 默认不限次数，使用0
    allowed_types: [],
    folder_id: '',
    expires_in_days: undefined,
  })

  const selectedFolderPath = ref<string>('')

  const formatSize = (bytes: number) => {
    if (bytes === 0) {
      return $t('settings.apiKeys.text.unlimited')
    }

    const k = 1024
    const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
    const i = Math.floor(Math.log(bytes) / Math.log(k))

    return `${parseFloat((bytes / Math.pow(k, i)).toFixed(2))} ${sizes[i]}`
  }

  const loadFileTypes = async () => {
    try {
      const resp = await get<{ supported_extensions?: string[]; mime_map?: Record<string, string> }>(
        '/config/upload/capabilities'
      )
      const exts: string[] = resp?.data?.supported_extensions || []
      const mimeMap: Record<string, string> = resp?.data?.mime_map || {}

      if (Array.isArray(exts) && exts.length > 0) {
        const options = [{ label: $t('settings.apiKeys.text.allFileTypes'), value: '*' }]

        exts.forEach((ext) => {
          const mime = mimeMap[ext] || `image/${ext}`
          const label = `${ext.toUpperCase()}${$t('settings.apiKeys.text.file')}` // 简单格式：JPG文件、PNG文件等
          options.push({ label, value: mime })
        })

        fileTypeOptions.value = options
      }
    } catch (error) {
      console.warn('Failed to load file types, using default options:', error)
    }
  }

  const initForm = () => {
    form.storage_limit = 0 // 默认不限量
    form.single_file_limit = 10 * 1024 * 1024 // 默认10MB
    form.upload_count_limit = 0 // 默认不限次数

    if (props.initialData) {
      Object.keys(props.initialData).forEach((key) => {
        if (props.initialData[key] !== undefined) {
          form[key] = props.initialData[key]
        }
      })

      if (props.initialData.folder_id) {
        selectedFolderPath.value = props.initialData.folder_id
      }
    }
  }

  onMounted(async () => {
    await loadFileTypes() // 先加载文件类型选项
    initForm()
  })

  const handleFolderSelected = (folder: any) => {
    if (folder) {
      selectedFolderPath.value = folder.path || folder.fullPath || folder.name
    } else {
      selectedFolderPath.value = ''
    }
  }

  watch(
    form,
    (newForm) => {
      emit('form-changed', { ...newForm })
    },
    { deep: true }
  )

  const onSubmit = () => {
    const formData = { ...form }

    if (formData.allowed_types.includes('*')) {
      formData.allowed_types = ['*']
    }

    if (formData.allowed_types.length === 0) {
      formData.allowed_types = ['*']
    }

    emit('submit', formData)
  }

  defineExpose({
    getFormData: () => ({ ...form }),
    form,
  })
</script>

<template>
  <form @submit.prevent="onSubmit">
    <div class="mb-6 grid grid-cols-1 gap-6">
      <div>
        <label class="mb-2 block text-content">
          {{ $t('settings.apiKeys.form.nameLabel') }}
        </label>
        <CyberInput v-model="form.name" type="text" :placeholder="$t('settings.apiKeys.form.namePlaceholder')" required />
      </div>

      <div class="grid grid-cols-1 gap-4 md:grid-cols-3">
        <div>
          <label class="mb-2 block text-content">
            {{ $t('settings.apiKeys.form.storageLimitLabel') }}
          </label>
          <CyberDropdown
            v-model="form.storage_limit"
            :options="storageSizeOptions"
            searchable
            :placeholder="$t('settings.apiKeys.form.storageLimitPlaceholder')"
            allow-create
            :create-label="$t('settings.apiKeys.form.storageLimitCreateLabel')"
          />
          <p class="mt-1 text-xs text-brand-600">
            {{
              $t('settings.apiKeys.form.storageLimitHint', {
                value: formatSize(form.storage_limit),
              })
            }}
          </p>
        </div>

        <div>
          <label class="mb-2 block text-content">
            {{ $t('settings.apiKeys.form.singleFileLimitLabel') }}
          </label>
          <CyberDropdown
            v-model="form.single_file_limit"
            :options="fileSizeOptions"
            searchable
            :placeholder="$t('settings.apiKeys.form.singleFileLimitPlaceholder')"
            allow-create
            :create-label="$t('settings.apiKeys.form.storageLimitCreateLabel')"
          />
          <p class="mt-1 text-xs text-brand-600">
            {{
              $t('settings.apiKeys.form.singleFileLimitHint', {
                value: formatSize(form.single_file_limit),
              })
            }}
          </p>
        </div>

        <div>
          <label class="mb-2 block text-content">
            {{ $t('settings.apiKeys.form.uploadLimitLabel') }}
          </label>
          <CyberInput
            v-model.number="form.upload_count_limit"
            type="number"
            :placeholder="$t('settings.apiKeys.form.uploadLimitPlaceholder')"
            min="0"
          />
          <p class="mt-1 text-xs text-brand-600">
            {{
              $t('settings.apiKeys.form.uploadLimitHint', {
                suffix: form.upload_count_limit === 0 ? ` ${$t('settings.apiKeys.form.uploadLimitUnlimitedSuffix')}` : '',
              })
            }}
          </p>
        </div>
      </div>

      <div>
        <label class="mb-2 block text-content">
          {{ $t('settings.apiKeys.form.allowedTypesLabel') }}
        </label>
        <CyberDropdown
          v-model="form.allowed_types"
          :options="fileTypeOptions"
          multiple
          searchable
          :placeholder="$t('settings.apiKeys.form.allowedTypesPlaceholder')"
          class="w-full"
        />
        <p class="mt-1 text-xs text-brand-600">
          {{ $t('settings.apiKeys.form.allowedTypesHint') }}
        </p>
      </div>

      <div>
        <label class="mb-2 block text-content">
          {{ $t('settings.apiKeys.form.expiresLabel') }}
        </label>
        <CyberInput
          v-model.number="form.expires_in_days"
          type="number"
          :placeholder="$t('settings.apiKeys.form.expiresPlaceholder')"
          min="1"
        />
        <p class="mt-1 text-xs text-brand-600">
          {{ $t('settings.apiKeys.form.expiresHint') }}
        </p>
      </div>

      <div>
        <label class="mb-2 block text-content">
          {{ $t('settings.apiKeys.form.folderLabel') }}
        </label>
        <CyberFolderTree
          v-model="form.folder_id"
          class="compact-dropdown text-sm"
          :dropdown-z-index="6000"
          @folder-selected="handleFolderSelected"
        />
        <p class="mt-1 text-xs text-brand-600">
          {{ $t('settings.apiKeys.form.folderHint') }}
        </p>
      </div>
    </div>
  </form>
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
