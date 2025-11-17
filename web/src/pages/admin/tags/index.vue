<script setup lang="ts">
  import { computed, onMounted, ref } from 'vue'
  import { useToast } from '@/components/Toast/useToast'
  import { useTexts } from '@/composables/useTexts'
  import { tagApi } from '@/api'
  import type { BatchTagOperation, TagInfo } from '@/api/types/index'
  import { useRouter } from 'vue-router'

  import TagCloud from './components/TagCloud.vue'
  import TagFormModal from './components/TagFormModal.vue'
  import BatchOperationModal from './components/BatchOperationModal.vue'
  // import { COMMON_ACTIONS } from '@/constants'

  const { $t } = useTexts()
  const toast = useToast()
  const router = useRouter()

  const isLoading = ref(false)
  const isSubmitting = ref(false)
  const isDeleting = ref(false)
  const isBatchProcessing = ref(false)
  const isBatchMode = ref(false)
  const totalTags = ref(0)
  const allTags = ref<TagInfo[]>([])
  const searchKeyword = ref('')
  const selectedTagIds = ref<number[]>([])

  const showCreateModal = ref(false)
  const showEditModal = ref(false)
  const showDeleteConfirm = ref(false)
  const showBatchOperationModal = ref(false)

  const editingTag = ref<TagInfo | null>(null)
  const deletingTag = ref<TagInfo | null>(null)

  const filteredTags = computed(() => {
    if (!searchKeyword.value.trim()) {
      return allTags.value
    }
    const keyword = searchKeyword.value.toLowerCase().trim()
    return allTags.value.filter((tag) => tag.name.toLowerCase().includes(keyword))
  })
  const showModalVisible = computed({
    get: () => showCreateModal.value || showEditModal.value,
    set: (value) => {
      if (!value) {
        closeFormModal()
      }
    },
  })
  const loadTags = async () => {
    try {
      isLoading.value = true
      const response = await tagApi.getTagList({
        keyword: '',
        page: 1,
        size: 10000,
      })
      allTags.value = response.data.items || []
      totalTags.value = response.data.total

      const newIds = allTags.value.map((tag) => tag.id)
      selectedTagIds.value = selectedTagIds.value.filter((id) => newIds.includes(id))
    } catch (_error) {
    } finally {
      isLoading.value = false
    }
  }
  const toggleTagSelection = (tag: TagInfo) => {
    const index = selectedTagIds.value.indexOf(tag.id)
    if (index === -1) {
      selectedTagIds.value.push(tag.id)
    } else {
      selectedTagIds.value.splice(index, 1)
    }
  }

  const clearSelection = () => {
    selectedTagIds.value = []
  }

  const openEditModal = (tag: TagInfo) => {
    editingTag.value = tag
    showEditModal.value = true
  }

  const closeFormModal = () => {
    showCreateModal.value = false
    showEditModal.value = false
    editingTag.value = null
  }

  const submitTagForm = async (formData: Record<string, unknown>) => {
    try {
      isSubmitting.value = true

      if (showEditModal.value && editingTag.value) {
        await tagApi.updateTag({
          id: editingTag.value.id,
          name: formData.name,
        })
        toast.success($t('admin.tags.messages.tagUpdated'))
      } else {
        await tagApi.createTag({
          name: formData.name,
        })
        toast.success($t('admin.tags.messages.tagCreated'))
      }

      closeFormModal()
      loadTags()
    } catch (_error) {
    } finally {
      isSubmitting.value = false
    }
  }

  const confirmDeleteTag = (tag: TagInfo) => {
    deletingTag.value = tag
    showDeleteConfirm.value = true
  }

  const deleteTag = async () => {
    if (!deletingTag.value) {
      return
    }

    try {
      isDeleting.value = true
      await tagApi.deleteTag(deletingTag.value.id)
      toast.success($t('admin.tags.messages.tagDeleted'))
      showDeleteConfirm.value = false

      const index = selectedTagIds.value.indexOf(deletingTag.value.id)
      if (index !== -1) {
        selectedTagIds.value.splice(index, 1)
      }

      loadTags()
    } catch (_error) {
    } finally {
      isDeleting.value = false
    }
  }

  const getSelectedTags = () => (allTags.value || []).filter((tag) => selectedTagIds.value.includes(tag.id))

  const doBatchOperation = async (operation: BatchTagOperation) => {
    try {
      isBatchProcessing.value = true
      await tagApi.batchOperateTags(operation)

      if (operation.action === 'delete') {
        toast.success($t('admin.tags.messages.batchDeleted'))
      } else {
        toast.success($t('admin.tags.messages.batchMerged'))
      }
      showBatchOperationModal.value = false
      selectedTagIds.value = []
      loadTags()
    } catch (_error) {
    } finally {
      isBatchProcessing.value = false
    }
  }

  const toggleBatchMode = () => {
    isBatchMode.value = !isBatchMode.value
    if (!isBatchMode.value) {
      selectedTagIds.value = []
    }
  }

  const viewSelectedTagsImages = () => {
    if (selectedTagIds.value.length === 0) {
      return
    }

    const tagIds = selectedTagIds.value.join(',')

    router.push({
      name: 'adminFiles',
      query: { tag_id: tagIds },
    })
  }

  onMounted(() => {
    loadTags()
  })
</script>

<template>
  <div class="admin-tags-page admin-page-container">
    <CyberAdminWrapper :title="$t('admin.tags.title')" :subtitle="$t('admin.tags.subtitle')" icon="fas fa-tags">
      <template #actions>
        <div class="search-input-wrapper">
          <CyberInput v-model="searchKeyword" :placeholder="$t('admin.tags.searchPlaceholder')" size="small" clearable>
            <template #prefix>
              <i class="fas fa-search text-content-muted"></i>
            </template>
          </CyberInput>
        </div>
        <CyberButton type="primary" class="text-xs" @click="showCreateModal = true">
          <i class="fas fa-plus mr-1.5" />{{ $t('admin.tags.createTag') }}
        </CyberButton>
        <CyberButton :type="isBatchMode ? 'success' : 'outlined'" class="text-xs" @click="toggleBatchMode">
          <i class="fas fa-tasks mr-1.5" />{{ isBatchMode ? $t('admin.tags.exitBatchMode') : $t('admin.tags.batchMode') }}
        </CyberButton>
        <CyberButton
          v-if="isBatchMode"
          type="primary"
          class="text-xs"
          :disabled="selectedTagIds.length === 0"
          @click="viewSelectedTagsImages"
        >
          <i class="fas fa-images mr-1.5" />{{ $t('admin.tags.viewSelectedFiles') }}
        </CyberButton>
        <CyberButton
          v-if="isBatchMode"
          type="warning"
          class="text-xs"
          :disabled="selectedTagIds.length === 0"
          @click="showBatchOperationModal = true"
        >
          <i class="fas fa-cogs mr-1.5" />{{ $t('admin.tags.executeBatchOperation') }}
        </CyberButton>
        <CyberButton v-if="selectedTagIds.length > 0" type="text" @click="clearSelection">
          <i class="fas fa-times mr-1" />{{ $t('admin.tags.cancelSelection') }}
        </CyberButton>
      </template>

      <template #content>
        <div class="admin-content-wrapper">
          <div class="tags-content">
            <div v-if="isBatchMode" class="batch-mode-notice">
              <span class="font-medium text-brand-400">{{ $t('admin.tags.batchModeNotice') }}</span> -
              {{ $t('admin.tags.batchModeHint') }}
            </div>

            <TagCloud
              :tags="filteredTags"
              :selected-tag-ids="selectedTagIds"
              :is-loading="isLoading"
              :total-tags="totalTags"
              :is-batch-mode="isBatchMode"
              @select="toggleTagSelection"
              @edit="openEditModal"
              @delete="confirmDeleteTag"
            />
          </div>
        </div>
      </template>
    </CyberAdminWrapper>

    <TagFormModal
      v-model="showModalVisible"
      :is-editing="showEditModal"
      :editing-tag="editingTag"
      :is-submitting="isSubmitting"
      @close="closeFormModal"
      @submit="submitTagForm"
    />

    <CyberDialog v-model="showDeleteConfirm" :title="$t('admin.tags.deleteDialog.title')" @close="showDeleteConfirm = false">
      <div class="mb-4">
        <p class="text-content">
          {{ $t('admin.tags.deleteDialog.confirmText') }}
          <span class="font-medium text-brand-400">{{ deletingTag?.name }}</span
          >{{ $t('common.questionMark') }}
        </p>
        <div class="mt-3 rounded-lg bg-background-700 p-3 text-xs">
          <p class="text-content-muted">
            {{ $t('admin.tags.deleteDialog.usageInfo') }}
            <span class="font-medium text-content">{{ deletingTag?.count || 0 }}</span>
            {{ $t('admin.tags.deleteDialog.usageTimes') }}
          </p>
          <p class="mt-2 font-medium text-error-400">{{ $t('admin.tags.deleteDialog.warning') }}</p>
        </div>
      </div>

      <template #footer>
        <div class="flex justify-end gap-2 p-4">
          <CyberButton type="secondary" @click="showDeleteConfirm = false">{{
            $t('admin.tags.deleteDialog.cancel')
          }}</CyberButton>
          <CyberButton type="danger" :loading="isDeleting" :disabled="isDeleting" @click="deleteTag">{{
            $t('admin.tags.deleteDialog.delete')
          }}</CyberButton>
        </div>
      </template>
    </CyberDialog>

    <BatchOperationModal
      v-model="showBatchOperationModal"
      :selected-tags="getSelectedTags()"
      :available-tags="allTags"
      :is-submitting="isBatchProcessing"
      @close="showBatchOperationModal = false"
      @submit="doBatchOperation"
    />
  </div>
</template>

<style scoped lang="scss">
  .search-input-wrapper {
    width: 240px;

    :deep(.cyber-input) {
      font-size: var(--text-sm);
    }
  }

  .admin-content-wrapper {
    background: var(--color-background-900);
    padding: 0;
    display: flex;
    flex-direction: column;
    flex: 1;
    min-height: 0;
  }

  .tags-content {
    display: flex;
    flex-direction: column;
    gap: var(--space-md);
    flex: 1;
    min-height: 0;
  }

  .batch-mode-notice {
    background: linear-gradient(
      135deg,
      rgba(var(--color-brand-500-rgb), 0.12) 0%,
      rgba(var(--color-brand-500-rgb), 0.08) 50%,
      rgba(var(--color-brand-500-rgb), 0.05) 100%
    );
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.2);
    border-radius: var(--radius-md);
    padding: var(--space-xs) var(--space-sm);
    font-size: var(--text-xs);
    color: var(--color-content-muted);
    flex-shrink: 0;
    backdrop-filter: var(--backdrop-blur-md);
    box-shadow:
      var(--shadow-sm),
      inset 0 1px 0 rgba(255, 255, 255, 0.05);
  }

  @media (max-width: 768px) {
    .search-input-wrapper {
      width: 200px;
    }

    .tags-content {
      gap: var(--space-md);
    }

    .batch-mode-notice {
      padding: var(--space-xs) var(--space-sm);
      font-size: var(--text-sm);
    }
  }

  @media (max-width: 480px) {
    .search-input-wrapper {
      width: 100%;
      max-width: 160px;
    }

    .tags-content {
      gap: var(--space-sm);
    }

    .batch-mode-notice {
      padding: var(--space-xs);
      font-size: var(--text-xs);
    }
  }
</style>
