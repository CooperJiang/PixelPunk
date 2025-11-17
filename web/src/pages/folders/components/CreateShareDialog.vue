<script setup lang="ts">
  import { computed, ref } from 'vue'
  import { shareApi } from '@/api'
  import { useToast } from '@/components/Toast/useToast'
  import { useTexts } from '@/composables/useTexts'
  import { copyToClipboard } from '@/utils/file/clipboard'
  import type { FolderInfo } from '@/api/types/index'
  import type { CreateShareParams } from '@/api/share/types'

  const { $t } = useTexts()

  const props = defineProps({
    modelValue: {
      type: Boolean,
      default: false,
    },
    selectedFolders: {
      type: Array as () => FolderInfo[],
      default: () => [],
    },
    selectedImages: {
      type: Array as () => any[],
      default: () => [],
    },
  })

  const emit = defineEmits(['update:modelValue', 'created', 'update:selectedFolders', 'update:selectedImages'])

  const toast = useToast()

  /* 对话框显示状态 */
  const visible = computed({
    get: () => props.modelValue,
    set: (value) => emit('update:modelValue', value),
  })

  /* 表单数据 */
  const formData = ref<CreateShareParams>({
    name: '',
    description: '',
    password: '',
    expired_days: 7,
    max_views: 0,
    items: [],
    collect_visitor_info: false,
    notification_on_access: false,
    notification_threshold: 100,
  })

  const submitting = ref(false)

  /* 是否禁用提交按钮 */
  const isSubmitDisabled = computed(
    () => (props.selectedFolders.length === 0 && props.selectedImages.length === 0) || submitting.value
  )

  const handleClose = () => {
    visible.value = false
  }

  const removeFolder = (folder: FolderInfo) => {
    const index = props.selectedFolders.findIndex((f) => f.id === folder.id)
    if (index !== -1) {
      const updatedFolders = [...props.selectedFolders]
      updatedFolders.splice(index, 1)
      emit('update:selectedFolders', updatedFolders)
    }
  }

  const removeImage = (image: any) => {
    const index = props.selectedImages.findIndex((img) => img.id === image.id)
    if (index !== -1) {
      const updatedImages = [...props.selectedImages]
      updatedImages.splice(index, 1)
      emit('update:selectedImages', updatedImages)
    }
  }

  const handleSubmit = async () => {
    submitting.value = true

    const items = [
      ...props.selectedFolders.map((folder) => ({
        item_type: 'folder' as const,
        item_id: folder.id,
      })),
      ...props.selectedImages.map((image) => ({
        item_type: 'file' as const,
        item_id: image.id,
      })),
    ]

    formData.value.items = items

    try {
      const result = await shareApi.createShare(formData.value)

      if (result.success && result.data) {
        toast.success($t('folders.share.createSuccess'))

        try {
          await copyToClipboard(result.data.share_url)
          toast.success($t('folders.share.createSuccessWithCopy'))
        } catch (error) {
          console.error('Failed to copy to clipboard:', error)
          toast.info($t('folders.share.manualCopy'))
        }

        visible.value = false
        emit('created', result.data)
      }
    } catch (error) {
      toast.error(error.message)
    }

    submitting.value = false
  }
</script>

<template>
  <CyberDialog
    v-model="visible"
    :title="$t('folders.createShareDialog.title')"
    width="500px"
    :append-to-body="true"
    @close="handleClose"
  >
    <div class="create-share-dialog">
      <div class="share-form">
        <div class="form-section">
          <h3 class="section-title">{{ $t('folders.createShareDialog.sections.basic') }}</h3>
          <div class="form-row">
            <label>{{ $t('folders.createShareDialog.form.name.label') }}</label>
            <CyberInput v-model="formData.name" :placeholder="$t('folders.createShareDialog.form.name.placeholder')" />
          </div>
          <div class="form-row">
            <label>{{ $t('folders.share.description') }}</label>
            <CyberInput
              v-model="formData.description"
              :placeholder="$t('folders.createShareDialog.form.description.placeholder')"
              type="textarea"
              :maxlength="200"
              :show-word-limit="true"
              :rows="2"
            />
          </div>
        </div>

        <div class="form-section">
          <h3 class="section-title">{{ $t('folders.createShareDialog.sections.access') }}</h3>
          <div class="form-row">
            <label>{{ $t('folders.createShareDialog.form.password.label') }}</label>
            <CyberInput
              v-model="formData.password"
              :placeholder="$t('folders.createShareDialog.form.password.placeholder')"
              type="password"
            />
          </div>
          <div class="form-row">
            <label>{{ $t('folders.createShareDialog.form.expiredDays.label') }}</label>
            <CyberInput
              v-model.number="formData.expired_days"
              :placeholder="$t('folders.createShareDialog.form.expiredDays.placeholder')"
              type="number"
              min="0"
            />
          </div>
          <div class="form-row">
            <label>{{ $t('folders.createShareDialog.form.maxViews.label') }}</label>
            <CyberInput
              v-model.number="formData.max_views"
              :placeholder="$t('folders.createShareDialog.form.maxViews.placeholder')"
              type="number"
              min="0"
            />
          </div>
        </div>

        <div class="form-section">
          <h3 class="section-title">{{ $t('folders.createShareDialog.sections.advanced') }}</h3>
          <div class="form-row">
            <label>{{ $t('folders.createShareDialog.form.collectVisitor.label') }}</label>
            <div class="input-wrapper">
              <CyberSwitch v-model="formData.collect_visitor_info" />
              <div class="input-tip">{{ $t('folders.createShareDialog.form.collectVisitor.tip') }}</div>
            </div>
          </div>
          <div class="form-row">
            <label>{{ $t('folders.createShareDialog.form.notification.label') }}</label>
            <div class="input-wrapper">
              <CyberSwitch v-model="formData.notification_on_access" />
              <div class="input-tip">{{ $t('folders.createShareDialog.form.notification.tip') }}</div>
            </div>
          </div>
          <div v-if="formData.notification_on_access" class="form-row notification-threshold-row">
            <label>{{ $t('folders.createShareDialog.form.notificationThreshold.label') }}</label>
            <div class="input-wrapper">
              <CyberInput
                v-model.number="formData.notification_threshold"
                :placeholder="$t('folders.createShareDialog.form.notificationThreshold.placeholder')"
                type="number"
                min="1"
              />
              <div class="input-tip">{{ $t('folders.createShareDialog.form.notificationThreshold.tip') }}</div>
            </div>
          </div>
        </div>

        <div class="form-section">
          <h3 class="section-title">{{ $t('folders.createShareDialog.sections.content') }}</h3>
          <div class="share-items">
            <template v-if="selectedFolders.length === 0 && selectedImages.length === 0">
              <div class="empty-items">{{ $t('folders.createShareDialog.emptyItems') }}</div>
            </template>
            <div v-else class="item-list">
              <!-- 文件夹列表 -->
              <div v-for="folder in selectedFolders" :key="`folder-${folder.id}`" class="share-item">
                <div class="item-icon">
                  <i class="fas fa-folder" />
                </div>
                <div class="item-name">{{ folder.name }}</div>
                <button class="remove-btn" @click="removeFolder(folder)">
                  <i class="fas fa-times" />
                </button>
              </div>

              <!-- 文件列表 -->
              <div v-for="image in selectedImages" :key="`image-${image.id}`" class="share-item">
                <div class="item-icon">
                  <i class="fas fa-image" />
                </div>
                <div class="item-name">{{ image.display_name || image.file_name }}</div>
                <button class="remove-btn" @click="removeImage(image)">
                  <i class="fas fa-times" />
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <template #footer>
      <div class="dialog-footer p-4">
        <CyberButton type="outlined" @click="handleClose">{{ $t('folders.createShareDialog.cancel') }}</CyberButton>
        <CyberButton type="primary" :disabled="isSubmitDisabled" :loading="submitting" @click="handleSubmit">
          {{ $t('folders.createShareDialog.submit') }}
        </CyberButton>
      </div>
    </template>
  </CyberDialog>
</template>

<style scoped>
  .create-share-dialog {
    max-height: 60vh;
    overflow-y: auto;
    padding-right: 10px;
  }

  .create-share-dialog::-webkit-scrollbar {
    width: 4px;
  }

  .create-share-dialog::-webkit-scrollbar-track {
    background: rgba(var(--color-brand-500-rgb), 0.03);
    border-radius: var(--radius-sm);
  }

  .create-share-dialog::-webkit-scrollbar-thumb {
    background: rgba(var(--color-brand-500-rgb), 0.15);
    border-radius: var(--radius-sm);
  }

  .create-share-dialog::-webkit-scrollbar-thumb:hover {
    background: rgba(var(--color-brand-500-rgb), 0.3);
  }

  .share-form {
    display: flex;
    flex-direction: column;
    gap: 1.5rem;
  }

  .form-section {
    display: flex;
    flex-direction: column;
    gap: 0.75rem;
  }

  .section-title {
    margin: 0;
    font-size: 1rem;
    color: rgba(var(--color-brand-500-rgb), 0.9);
    border-bottom: 1px solid rgba(var(--color-brand-500-rgb), 0.1);
    padding-bottom: 0.5rem;
  }

  .form-row {
    display: flex;
    align-items: flex-start;
    gap: 1rem;
  }

  .form-row label {
    width: 100px;
    color: rgba(255, 255, 255, 0.85);
    font-size: 0.9rem;
    flex-shrink: 0;
    padding-top: 5px;
  }

  .input-wrapper {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }

  .input-tip {
    font-size: 0.8rem;
    color: rgba(var(--color-content-muted-rgb), 0.6);
    line-height: 1.2;
  }

  .share-items {
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.1);
    border-radius: var(--radius-sm);
    background: rgba(var(--color-background-900-rgb), 0.5);
    padding: 0.75rem;
    max-height: 200px;
    overflow-y: auto;
  }

  .share-items::-webkit-scrollbar {
    width: 4px;
  }

  .share-items::-webkit-scrollbar-track {
    background: rgba(var(--color-brand-500-rgb), 0.03);
    border-radius: var(--radius-sm);
  }

  .share-items::-webkit-scrollbar-thumb {
    background: rgba(var(--color-brand-500-rgb), 0.15);
    border-radius: var(--radius-sm);
  }

  .share-items::-webkit-scrollbar-thumb:hover {
    background: rgba(var(--color-brand-500-rgb), 0.3);
  }

  .empty-items {
    color: rgba(255, 255, 255, 0.5);
    font-size: 0.9rem;
    text-align: center;
    padding: 1rem;
  }

  .item-list {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }

  .share-item {
    display: flex;
    align-items: center;
    background: rgba(15, 25, 35, 0.8);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.1);
    border-radius: var(--radius-sm);
    padding: 0.5rem 0.75rem;
    gap: 0.5rem;
  }

  .item-icon {
    color: rgba(var(--color-brand-500-rgb), 0.9);
    width: 20px;
    flex-shrink: 0;
  }

  .item-name {
    flex: 1;
    color: rgba(var(--color-content-default-rgb), 0.9);
    font-size: 0.9rem;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .remove-btn {
    width: 20px;
    height: 20px;
    border-radius: var(--radius-full);
    border: none;
    background: rgba(var(--color-error-rgb), 0.2);
    color: rgba(var(--color-error-rgb), 0.9);
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    transition: all 0.2s;
    flex-shrink: 0;
  }

  .remove-btn:hover {
    background: rgba(var(--color-error-rgb), 0.3);
    transform: scale(1.1);
  }

  .dialog-footer {
    display: flex;
    justify-content: flex-end;
    gap: 0.75rem;
  }

  .dialog-footer :deep(.cyber-btn) {
    height: auto !important;
  }

  .notification-threshold-row {
    margin-top: 0.5rem;
    margin-left: 1.5rem;
    padding-left: 0.5rem;
    border-left: 2px solid rgba(var(--color-brand-500-rgb), 0.3);
  }
</style>
