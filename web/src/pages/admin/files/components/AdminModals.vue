<script setup lang="ts">
  import { useTexts } from '@/composables/useTexts'
  import type { FileInfo } from '@/api/admin/files'

  defineProps<{
    showDetailModal: boolean
    showDeleteConfirm: boolean
    showBatchDeleteConfirm: boolean
    selectedFile: FileInfo | null
    selectedCount: number
    isDeleting: boolean
    isBatchProcessing: boolean
  }>()

  defineEmits<{
    (e: 'update:showDetailModal', value: boolean): void
    (e: 'update:showDeleteConfirm', value: boolean): void
    (e: 'update:showBatchDeleteConfirm', value: boolean): void
    (e: 'delete'): void
    (e: 'batchDelete'): void
    (e: 'recommend', id: string, recommended: boolean): void
    (e: 'confirmDelete', id?: string): void
  }>()

  const { $t } = useTexts()
</script>

<template>
  <div>
    <CyberFileDetailModal
      :model-value="showDetailModal"
      :file-info="selectedFile"
      is-admin
      @update:model-value="$emit('update:showDetailModal', $event)"
      @close="$emit('update:showDetailModal', false)"
      @delete="$emit('confirmDelete', selectedFile?.id)"
      @recommend="$emit('recommend', $event, !selectedFile?.is_recommended)"
    />

    <CyberDialog
      :model-value="showDeleteConfirm"
      :title="$t('admin.files.modals.deleteConfirm.title')"
      :show-footer="true"
      class="text-sm"
      @update:model-value="$emit('update:showDeleteConfirm', $event)"
    >
      <div class="p-3">
        <p class="text-sm text-content-heading">{{ $t('admin.files.modals.deleteConfirm.message') }}</p>
      </div>
      <template #footer>
        <div class="flex justify-end space-x-2 p-4">
          <CyberButton type="secondary" @click="$emit('update:showDeleteConfirm', false)">{{
            $t('admin.files.modals.deleteConfirm.cancel')
          }}</CyberButton>
          <CyberButton type="danger" :loading="isDeleting" @click="$emit('delete')">
            <i class="fas fa-trash-alt mr-1.5" />{{ $t('admin.files.modals.deleteConfirm.confirm') }}
          </CyberButton>
        </div>
      </template>
    </CyberDialog>

    <CyberDialog
      :model-value="showBatchDeleteConfirm"
      :title="$t('admin.files.modals.batchDeleteConfirm.title')"
      :show-footer="true"
      class="text-sm"
      @update:model-value="$emit('update:showBatchDeleteConfirm', $event)"
    >
      <div class="p-3">
        <div class="mb-3 flex items-center">
          <i class="fas fa-exclamation-triangle mr-3 text-xl text-warning" />
          <div>
            <p class="text-sm font-medium text-content-heading">
              {{ $t('admin.files.modals.batchDeleteConfirm.message', { count: selectedCount }) }}
            </p>
            <p class="mt-1 text-xs text-content-muted">{{ $t('admin.files.modals.batchDeleteConfirm.warning') }}</p>
          </div>
        </div>
        <div class="rounded border border-subtle bg-background-400 p-3">
          <p class="text-xs text-content">
            <i class="fas fa-info-circle mr-1" />
            {{ $t('admin.files.modals.batchDeleteConfirm.info') }}
          </p>
          <ul class="ml-4 mt-2 space-y-1 text-xs text-content-muted">
            <li v-for="(consequence, index) in $t('admin.files.modals.batchDeleteConfirm.consequences')" :key="index">
              • {{ consequence }}
            </li>
          </ul>
        </div>
      </div>
      <template #footer>
        <div class="flex justify-end space-x-2 p-4">
          <CyberButton type="secondary" @click="$emit('update:showBatchDeleteConfirm', false)">{{
            $t('admin.files.modals.batchDeleteConfirm.cancel')
          }}</CyberButton>
          <CyberButton type="danger" :loading="isBatchProcessing" @click="$emit('batchDelete')">
            <i class="fas fa-trash-alt mr-1.5" />{{
              $t('admin.files.modals.batchDeleteConfirm.confirm', { count: selectedCount })
            }}
          </CyberButton>
        </div>
      </template>
    </CyberDialog>
  </div>
</template>
