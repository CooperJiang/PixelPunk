<script setup lang="ts">
  import { defineExpose, defineProps, onMounted, reactive, ref } from 'vue'
  import { useToast } from '@/components/Toast/useToast'
  import api from '@/api'
  import type { ApiKeyInfo, CreateApiKeyRequest, UpdateApiKeyRequest } from '@/api/types/index'

  import SizeConverter from './SizeConverter.vue'
  import ApiKeyList from './ApiKeyList.vue'
  import ApiKeyForm from './ApiKeyForm.vue'
  import ApiKeyDetail from './ApiKeyDetail.vue'
  import { useTexts } from '@/composables/useTexts'
  import Validator, { getValidationRules } from '@/utils/validation/validator'

  const props = defineProps({
    lazyLoad: {
      type: Boolean,
      default: false,
    },
  })

  const toast = useToast()
  const { $t } = useTexts()
  const R = getValidationRules($t)

  const apiKeys = ref<ApiKeyInfo[]>([])
  const isLoading = ref(false)
  const pagination = reactive({
    page: 1,
    size: 10,
    total: 0,
    total_pages: 0,
  })

  const changePage = (page: number) => {
    if (page < 1 || page > pagination.total_pages) {
      return
    }
    pagination.page = page
    fetchApiKeys()
  }

  const showDialog = ref(false)
  const isEditing = ref(false)
  const isSubmitting = ref(false)
  const apiKeyForm = reactive<CreateApiKeyRequest & UpdateApiKeyRequest>({
    name: '',
    storage_limit: 0, // 不限量
    single_file_limit: 10485760, // 10MB
    upload_count_limit: 0,
    allowed_types: [],
    folder_id: '',
    expires_in_days: undefined,
  })
  const editingKeyId = ref('')

  const apiKeyFormRef = ref(null)

  const currentFormData = ref({})

  const handleFormChanged = (formData: Record<string, unknown>) => {
    currentFormData.value = formData
  }

  const showKeyDialog = ref(false)
  const selectedKey = ref<ApiKeyInfo | null>(null)
  const apiKeyValue = ref('')
  const copied = ref(false)

  const showDeleteDialog = ref(false)
  const isDeleting = ref(false)

  const fetchApiKeys = async () => {
    try {
      isLoading.value = true

      const result = await api.apikey.getApiKeyList({
        page: pagination.page,
        size: pagination.size,
      })

      if (result.success) {
        const { data } = result
        apiKeys.value = data.items
        pagination.total = data.pagination.total
        pagination.page = data.pagination.page
        pagination.size = data.pagination.size
        pagination.total_pages = data.pagination.total_pages
      }
    } catch (_error) {
    } finally {
      isLoading.value = false
    }
  }

  onMounted(() => {
    if (!props.lazyLoad) {
      fetchApiKeys()
    }
  })

  const openCreateDialog = () => {
    isEditing.value = false
    resetForm()
    showDialog.value = true
  }

  const editApiKey = (key: ApiKeyInfo) => {
    showKeyDialog.value = false
    isEditing.value = true
    editingKeyId.value = key.id

    apiKeyForm.name = key.name
    apiKeyForm.storage_limit = key.storage_limit
    apiKeyForm.single_file_limit = key.single_file_limit
    apiKeyForm.upload_count_limit = key.upload_count_limit
    apiKeyForm.allowed_types = [...(key.allowed_types || [])]
    apiKeyForm.folder_id = key.folder_id || ''
    apiKeyForm.expires_in_days = key.expires_in_days || undefined

    showDialog.value = true
  }

  const closeDialog = () => {
    showDialog.value = false
    resetForm()
  }

  const resetForm = () => {
    apiKeyForm.name = ''
    apiKeyForm.storage_limit = 0 // 不限量
    apiKeyForm.single_file_limit = 10485760 // 10MB
    apiKeyForm.upload_count_limit = 0
    apiKeyForm.allowed_types = []
    apiKeyForm.folder_id = ''
    apiKeyForm.expires_in_days = undefined
    editingKeyId.value = ''
  }

  const submitApiKey = async (formData: Record<string, unknown>) => {
    const nameCheck = Validator.validate(String(formData?.name || ''), [R.required], $t)
    if (!nameCheck.valid) return toast.error(nameCheck.message || $t('settings.apiKeys.toast.nameRequired'))

    try {
      isSubmitting.value = true

      let result

      if (isEditing.value) {
        result = await api.apikey.updateApiKey(editingKeyId.value, formData)
      } else {
        result = await api.apikey.createApiKey(formData)
      }

      if (result.success) {
        const { data } = result

        toast.success(isEditing.value ? $t('settings.apiKeys.toasts.updated') : $t('settings.apiKeys.toasts.created'))

        if (!isEditing.value && data.key) {
          selectedKey.value = data
          apiKeyValue.value = data.key
          showDialog.value = false
          showKeyDialog.value = true
        } else {
          closeDialog()
        }

        fetchApiKeys()
      }
    } catch (_error) {
    } finally {
      isSubmitting.value = false
    }
  }

  const submitFormFromButton = () => {
    if (currentFormData.value && Object.keys(currentFormData.value).length > 0) {
      submitApiKey(currentFormData.value)
    } else if (apiKeyFormRef.value && apiKeyFormRef.value.getFormData) {
      submitApiKey(apiKeyFormRef.value.getFormData())
    } else {
      submitApiKey(apiKeyForm)
    }
  }

  const showApiKey = async (key: ApiKeyInfo) => {
    selectedKey.value = key
    apiKeyValue.value = '' // 清空之前的密钥值
    showKeyDialog.value = true

    if (key.key) {
      apiKeyValue.value = key.key
    }
  }

  const closeKeyDialog = () => {
    showKeyDialog.value = false
    apiKeyValue.value = ''
    copied.value = false
  }

  const toggleApiKeyStatus = async (key: ApiKeyInfo) => {
    try {
      const result = await api.apikey.toggleApiKey(key.id, { is_active: !key.is_active })

      if (result.success) {
        toast.success(key.is_active ? $t('settings.apiKeys.toasts.disabled') : $t('settings.apiKeys.toasts.enabled'))
        fetchApiKeys() // 刷新列表
      }
    } catch (_error) {}
  }

  const confirmDeleteApiKey = (key: ApiKeyInfo) => {
    selectedKey.value = key
    showDeleteDialog.value = true
  }

  const closeDeleteDialog = () => {
    showDeleteDialog.value = false
    selectedKey.value = null
  }

  const deleteApiKey = async () => {
    if (!selectedKey.value) {
      return
    }

    try {
      isDeleting.value = true

      const result = await api.apikey.deleteApiKey(selectedKey.value.id)

      if (result.success) {
        toast.success($t('settings.apiKeys.toasts.deleted'))
        closeDeleteDialog()
        fetchApiKeys() // 刷新列表
      }
    } catch (_error) {
    } finally {
      isDeleting.value = false
    }
  }

  defineExpose({
    fetchApiKeys,
  })
</script>

<template>
  <div class="apikey-manager p-6">
    <SizeConverter />

    <div class="mb-6">
      <ApiKeyList
        :api-keys="apiKeys"
        :is-loading="isLoading"
        :pagination="pagination"
        @create="openCreateDialog"
        @edit="editApiKey"
        @detail="showApiKey"
        @delete="confirmDeleteApiKey"
        @toggle="toggleApiKeyStatus"
        @change-page="changePage"
      />
    </div>

    <CyberDialog
      id="create-edit-dialog"
      v-model="showDialog"
      :title="isEditing ? $t('settings.apiKeys.dialog.editTitle') : $t('settings.apiKeys.dialog.createTitle')"
      width="700px"
      :loading="isSubmitting"
      :show-default-footer="false"
      @cancel="closeDialog"
    >
      <ApiKeyForm
        ref="apiKeyFormRef"
        :initial-data="apiKeyForm"
        :is-editing="isEditing"
        @submit="submitApiKey"
        @form-changed="handleFormChanged"
      />

      <template #footer>
        <div class="flex justify-end gap-3 p-4">
          <CyberButton type="outlined" @click="closeDialog">
            {{ $t('settings.apiKeys.dialog.cancel') }}
          </CyberButton>
          <CyberButton type="primary" :loading="isSubmitting" loading-mode="inline" icon="save" @click="submitFormFromButton">
            {{
              isSubmitting
                ? $t('settings.apiKeys.dialog.submitting')
                : isEditing
                  ? $t('settings.apiKeys.dialog.saveAction')
                  : $t('settings.apiKeys.dialog.createAction')
            }}
          </CyberButton>
        </div>
      </template>
    </CyberDialog>

    <CyberDialog
      id="detail-dialog"
      v-model="showKeyDialog"
      :title="$t('settings.apiKeys.dialog.detailTitle')"
      width="750px"
      :show-default-footer="false"
      @cancel="closeKeyDialog"
    >
      <ApiKeyDetail :api-key="selectedKey" :api-key-value="apiKeyValue" @copy="copied = true" />

      <template #footer>
        <div class="flex justify-end gap-3 p-4">
          <CyberButton type="primary" @click="closeKeyDialog">
            {{ $t('settings.apiKeys.dialog.close') }}
          </CyberButton>
        </div>
      </template>
    </CyberDialog>

    <CyberDialog
      id="delete-dialog"
      v-model="showDeleteDialog"
      :title="$t('settings.apiKeys.dialog.deleteTitle')"
      width="450px"
      :show-default-footer="false"
      @cancel="closeDeleteDialog"
    >
      <p class="mb-6">
        {{
          $t('settings.apiKeys.dialog.deleteMessage', {
            name: selectedKey?.name || '',
          })
        }}
      </p>

      <template #footer>
        <div class="flex justify-end gap-3 p-4">
          <CyberButton type="ghost" @click="closeDeleteDialog">
            {{ $t('settings.apiKeys.dialog.cancel') }}
          </CyberButton>
          <CyberButton type="danger" icon="trash" :loading="isDeleting" @click="deleteApiKey">
            {{ isDeleting ? $t('status.deleting') : $t('settings.apiKeys.dialog.confirmDelete') }}
          </CyberButton>
        </div>
      </template>
    </CyberDialog>
  </div>
</template>

<style scoped>
  .apikey-manager {
    --settings-accent: var(--color-brand-500);
    --settings-accent-rgb: var(--color-brand-500-rgb);
  }

  .apikey-list {
    max-height: 60vh;
    overflow-y: auto;
  }

  .apikey-editor-panel p {
    margin: 0;
  }

  .cyber-btn {
    position: relative;
    overflow: hidden;
  }

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

  :deep(.cyber-input[type='number']) {
    -webkit-appearance: none;
    -moz-appearance: textfield;
  }

  :deep(.cyber-input[type='number']::-webkit-outer-spin-button),
  :deep(.cyber-input[type='number']::-webkit-inner-spin-button) {
    -webkit-appearance: none;
    margin: 0;
  }

  .apikey-details {
    @apply text-sm;
  }

  .apikey-details-content {
    @apply flex flex-col gap-4;
  }

  .info-item {
    @apply mb-2;
  }

  .info-label {
    @apply mb-1 text-xs text-brand-600;
  }

  .info-value {
    @apply text-content;
  }

  .status-badge {
    @apply inline-flex items-center rounded px-2 py-0.5 text-xs;
    border: 1px solid transparent;
  }

  .status-active {
    background-color: rgba(var(--color-success-rgb), 0.16);
    color: var(--color-success-400);
    border-color: rgba(var(--color-success-rgb), 0.3);
  }

  .status-inactive {
    background-color: rgba(var(--color-error-rgb), 0.16);
    color: var(--color-error-400);
    border-color: rgba(var(--color-error-rgb), 0.3);
  }

  .status-expired {
    background-color: rgba(var(--color-warning-rgb), 0.18);
    color: var(--color-warning-400);
    border-color: rgba(var(--color-warning-rgb), 0.32);
  }

  .folder-path {
    @apply flex items-center;
    color: var(--color-content-muted);
  }

  .usage-card {
    @apply flex flex-col border bg-background-800 p-3;
    border-radius: var(--radius-sm);
    border-color: var(--color-border-subtle);
  }

  .usage-header {
    @apply mb-2 flex items-center justify-between;
  }

  .usage-title {
    @apply flex items-center text-sm font-medium;
  }

  .usage-percentage {
    @apply font-mono text-sm;
  }

  .progress-bar {
    @apply mb-3 h-2 w-full overflow-hidden bg-background-600;
    border-radius: var(--radius-sm);
  }

  .progress-fill {
    @apply h-2;
    border-radius: var(--radius-sm);
  }

  .usage-details {
    @apply mt-2 flex flex-col space-y-1 text-xs;
  }

  .file-limit-display {
    @apply mb-2 flex h-12 items-center justify-center;
  }

  .file-limit-desc {
    @apply text-center text-xs text-brand-600;
  }

  .filetype-container {
    @apply flex flex-wrap gap-2;
  }

  .filetype-badge {
    @apply flex items-center rounded bg-brand-200 px-2 py-1 text-xs text-content;
  }
</style>
