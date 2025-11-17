<script setup lang="ts">
  import { onMounted, onUnmounted, ref } from 'vue'
  import { useRouter } from 'vue-router'
  import { useConfirmDialog } from '@/composables/useConfirmDialog'
  import { useTagManagement } from './composables/useTagManagement'
  import type { UserTag } from '@/api/user/tag'
  import TagCloud from './components/TagCloud.vue'
  import { useTexts } from '@/composables/useTexts'

  defineOptions({
    name: 'TagManage',
  })

  const { $t } = useTexts()

  const router = useRouter()
  const { dialogState, handleConfirm, handleCancel } = useConfirmDialog()

  const {
    isLoading,
    allTags,
    filteredTags,
    totalTags,
    selectedTagIds,
    stats,
    showFormModal,
    showEditModal,
    showMergeModal,
    isBatchMode,
    searchKeyword,
    tagForm,
    mergeForm,
    mergeTargetOptions,

    loadTags,
    loadStats,
    openCreateModal,
    openEditModal,
    closeFormModal,
    submitTagForm,
    handleDelete,
    toggleTagSelection,
    clearSelection,
    handleBatchDelete,
    openMergeModal,
    submitMerge,
    toggleBatchMode,
    viewSelectedTagsFiles,
  } = useTagManagement()

  /* 多选模式状态 */
  const isMultiSelectMode = ref(false) // 是否按住修饰键
  const tempSelectedTags = ref<number[]>([]) // 临时选中的标签ID

  /* 检测操作系统 */
  const isMac = ref(false)
  const modifierKeyName = ref('Ctrl')

  const detectOS = () => {
    const userAgent = window.navigator.userAgent.toLowerCase()
    const platform = window.navigator.platform.toLowerCase()

    isMac.value = /mac|iphone|ipad|ipod/.test(userAgent) || /mac/.test(platform)
    modifierKeyName.value = isMac.value ? 'Command' : 'Ctrl'
  }

  const goToResourceByTag = (tagList: UserTag | UserTag[]) => {
    const tagArray = Array.isArray(tagList) ? tagList : [tagList]
    const tagIds = tagArray.map((t) => t.id).join(',')
    const tagNames = tagArray.map((t) => t.name).join(',')

    router.push({
      path: '/resource',
      query: {
        tags: tagIds,
        tag_names: tagNames,
      },
    })
  }

  const handleTagClick = (tag: UserTag, _event?: MouseEvent) => {
    if (isBatchMode.value) {
      toggleTagSelection(tag)
      return
    }

    if (isMultiSelectMode.value) {
      const index = tempSelectedTags.value.indexOf(tag.id)
      if (index === -1) {
        tempSelectedTags.value.push(tag.id)
      } else {
        tempSelectedTags.value.splice(index, 1)
      }
      return
    }

    goToResourceByTag(tag)
  }

  const checkModifierKey = (event: KeyboardEvent | MouseEvent) => {
    return event.altKey || event.ctrlKey || event.metaKey
  }

  const handleKeyDown = (event: KeyboardEvent) => {
    if (checkModifierKey(event) && !isBatchMode.value) {
      isMultiSelectMode.value = true
    }
  }

  const handleKeyUp = (event: KeyboardEvent) => {
    if (!event.altKey && !event.ctrlKey && !event.metaKey) {
      isMultiSelectMode.value = false

      if (tempSelectedTags.value.length > 0) {
        const selectedTagList = allTags.value.filter((t) => tempSelectedTags.value.includes(t.id))
        goToResourceByTag(selectedTagList)
        tempSelectedTags.value = []
      }
    }
  }

  const handleWindowBlur = () => {
    isMultiSelectMode.value = false
    tempSelectedTags.value = []
  }

  onMounted(() => {
    loadTags()
    loadStats()
    detectOS() // 检测操作系统
    window.addEventListener('keydown', handleKeyDown)
    window.addEventListener('keyup', handleKeyUp)
    window.addEventListener('blur', handleWindowBlur)
  })

  onUnmounted(() => {
    window.removeEventListener('keydown', handleKeyDown)
    window.removeEventListener('keyup', handleKeyUp)
    window.removeEventListener('blur', handleWindowBlur)
  })
</script>

<template>
  <div class="tag-manage-page">
    <div class="page-header">
      <div class="flex items-center justify-between gap-4">
        <div class="flex items-center gap-3">
          <div class="icon-container">
            <i class="fas fa-tags text-lg text-brand-400" />
          </div>
          <div>
            <div class="flex items-center gap-2">
              <h1 class="text-lg font-bold text-content">{{ $t('tagManage.title') }}</h1>
              <span v-if="stats" class="bg-brand-500/10 rounded-full px-2 py-0.5 text-xs font-medium text-brand-400">
                {{ stats.total_tags }}
              </span>
            </div>
            <p class="mt-0.5 text-xs text-content-muted">{{ $t('tagManage.subtitle') }}</p>
          </div>
        </div>

        <div class="flex flex-wrap items-center gap-2">
          <div class="search-input-wrapper">
            <CyberInput v-model="searchKeyword" :placeholder="$t('tagManage.search')" size="small" clearable>
              <template #prefix>
                <i class="fas fa-search text-content-muted"></i>
              </template>
            </CyberInput>
          </div>

          <CyberButton type="primary" size="small" @click="openCreateModal">
            <i class="fas fa-plus mr-1.5" />
            {{ $t('tagManage.create') }}
          </CyberButton>

          <CyberButton :type="isBatchMode ? 'success' : 'outlined'" size="small" @click="toggleBatchMode">
            <i class="fas fa-tasks mr-1.5" />
            {{ isBatchMode ? $t('tagManage.exitBatchMode') : $t('tagManage.batchMode') }}
          </CyberButton>

          <CyberButton v-if="isBatchMode && selectedTagIds.length > 0" type="primary" size="small" @click="viewSelectedTagsFiles">
            <i class="fas fa-images mr-1.5" />
            {{ $t('tagManage.batch.viewFiles') }}
          </CyberButton>

          <CyberButton v-if="isBatchMode && selectedTagIds.length > 0" type="warning" size="small" @click="handleBatchDelete">
            <i class="fas fa-trash mr-1.5" />
            {{ $t('tagManage.batch.delete', { count: selectedTagIds.length }) }}
          </CyberButton>

          <CyberButton v-if="isBatchMode && selectedTagIds.length >= 2" type="info" size="small" @click="openMergeModal">
            <i class="fas fa-code-branch mr-1.5" />
            {{ $t('tagManage.batch.merge') }}
          </CyberButton>

          <CyberButton v-if="selectedTagIds.length > 0" type="text" size="small" @click="clearSelection">
            <i class="fas fa-times mr-1" />
            {{ $t('tagManage.batch.cancel') }}
          </CyberButton>
        </div>
      </div>
    </div>

    <!-- 批量模式提示 -->
    <div v-if="isBatchMode" class="batch-mode-notice">
      <i class="fas fa-info-circle mr-2 text-brand-400" />
      <span class="font-medium text-brand-400">{{ $t('tagManage.batch.modeHint') }}</span>
      - {{ $t('tagManage.batch.modeDescription', { count: selectedTagIds.length }) }}
    </div>

    <!-- 多选模式提示 -->
    <div v-if="isMultiSelectMode" class="multi-select-notice">
      <i class="fas fa-hand-pointer mr-2 text-success-400" />
      <span class="font-medium text-success-400">{{ $t('tagManage.multiSelect.modeHint') }}</span>
      - {{ $t('tagManage.multiSelect.modeDescription', { modifier: modifierKeyName, count: tempSelectedTags.length }) }}
    </div>

    <!-- 普通模式提示 -->
    <div v-if="!isBatchMode && !isMultiSelectMode && allTags.length > 0" class="normal-mode-notice">
      <i class="fas fa-lightbulb mr-2 text-warning-400" />
      <span class="text-content-muted">{{ $t('tagManage.normalMode.hint', { modifier: modifierKeyName }) }}</span>
    </div>

    <!-- 标签云容器 -->
    <div class="tags-content">
      <TagCloud
        :tags="filteredTags"
        :selected-tag-ids="selectedTagIds"
        :temp-selected-tag-ids="tempSelectedTags"
        :is-loading="isLoading"
        :total-tags="totalTags"
        :is-batch-mode="isBatchMode"
        @select="handleTagClick"
        @edit="openEditModal"
        @delete="handleDelete"
      />
    </div>

    <!-- 创建/编辑标签模态框 -->
    <CyberDialog
      v-model="showFormModal"
      :title="showEditModal ? $t('tagManage.dialog.edit.title') : $t('tagManage.dialog.create.title')"
      @close="closeFormModal"
    >
      <div class="space-y-4 pt-2">
        <div>
          <label class="mb-2 block text-sm font-medium text-content">
            <i class="fas fa-tag mr-1" />
            {{ $t('tagManage.form.name.label') }} <span class="text-error-400">*</span>
          </label>
          <CyberInput
            v-model="tagForm.name"
            :placeholder="$t('tagManage.form.name.placeholder')"
            :maxlength="50"
            prefix-icon="tag"
            @keyup.enter="submitTagForm"
          />
          <p class="mt-2 text-xs text-content-muted">{{ $t('tagManage.form.name.hint') }}</p>
        </div>
      </div>

      <template #footer>
        <div class="flex justify-end gap-3 p-4">
          <CyberButton type="secondary" @click="closeFormModal">{{ $t('tagManage.form.cancel') }}</CyberButton>
          <CyberButton type="primary" @click="submitTagForm">
            <i class="fas fa-check mr-2" />
            {{ showEditModal ? $t('tagManage.dialog.edit.submit') : $t('tagManage.dialog.create.submit') }}
          </CyberButton>
        </div>
      </template>
    </CyberDialog>

    <!-- 合并标签模态框 -->
    <CyberDialog v-model="showMergeModal" :title="$t('tagManage.dialog.merge.title')" @close="showMergeModal = false">
      <div class="space-y-4 pt-2">
        <p class="text-sm text-content-muted">{{ $t('tagManage.merge.description') }}</p>

        <div>
          <label class="mb-2 block text-sm font-medium text-content">
            <i class="fas fa-bullseye mr-1" />
            {{ $t('tagManage.dialog.merge.targetLabel') }} <span class="text-error-400">*</span>
          </label>
          <CyberDropdown
            v-model="mergeForm.targetId"
            :options="mergeTargetOptions"
            :placeholder="$t('tagManage.dialog.merge.targetPlaceholder')"
            searchable
          />
        </div>

        <div class="mt-3 rounded-lg bg-background-700 p-3">
          <p class="text-xs text-content-muted">
            <i class="fas fa-info-circle mr-1 text-brand-400" />
            {{ $t('tagManage.merge.selectedInfo', { count: selectedTagIds.length, remaining: selectedTagIds.length - 1 }) }}
          </p>
        </div>
      </div>

      <template #footer>
        <div class="flex justify-end gap-3 p-4">
          <CyberButton type="secondary" @click="showMergeModal = false">{{ $t('tagManage.merge.cancel') }}</CyberButton>
          <CyberButton type="primary" @click="submitMerge">
            <i class="fas fa-code-branch mr-2" />
            {{ $t('tagManage.merge.confirm') }}
          </CyberButton>
        </div>
      </template>
    </CyberDialog>

    <!-- 全局确认对话框 -->
    <CyberConfirmDialog
      v-model="dialogState.isVisible"
      :title="dialogState.config.title"
      :message="dialogState.config.message"
      :type="dialogState.config.type"
      :confirm-text="dialogState.config.confirmText"
      :cancel-text="dialogState.config.cancelText"
      :width="dialogState.config.width"
      :require-input="dialogState.config.requireInput"
      :input-label="dialogState.config.inputLabel"
      :input-placeholder="dialogState.config.inputPlaceholder"
      :expected-input="dialogState.config.expectedInput"
      @confirm="handleConfirm"
      @cancel="handleCancel"
    />
  </div>
</template>

<style scoped lang="scss">
  .tag-manage-page {
    display: flex;
    flex-direction: column;
    min-height: 100vh;
  }

  .search-input-wrapper {
    width: 240px;

    :deep(.cyber-input) {
      font-size: var(--text-sm);
    }
  }

  .page-header {
    background: rgba(var(--color-background-800-rgb), 0.6);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.15);
    border-radius: var(--radius-sm);
    padding: 1rem 1.25rem;
    margin-bottom: 1rem;
    box-shadow: var(--shadow-sm);
    backdrop-filter: var(--backdrop-blur-md);

    .icon-container {
      width: 2.25rem;
      height: 2.25rem;
      display: flex;
      align-items: center;
      justify-content: center;
      background: linear-gradient(135deg, rgba(var(--color-brand-500-rgb), 0.12), rgba(var(--color-brand-500-rgb), 0.06));
      border: 1px solid rgba(var(--color-brand-500-rgb), 0.2);
      border-radius: var(--radius-sm);
    }
  }

  .batch-mode-notice,
  .multi-select-notice,
  .normal-mode-notice {
    display: flex;
    align-items: center;
    border-radius: var(--radius-sm);
    padding: 0.75rem 1rem;
    margin-bottom: 1rem;
    font-size: var(--text-sm);
    color: var(--color-content-muted);
    backdrop-filter: var(--backdrop-blur-md);
    animation: slideDown 0.3s var(--ease-out);

    kbd {
      display: inline-block;
      padding: var(--space-xs) var(--space-sm);
      font-size: var(--text-xs);
      font-family: var(--font-mono);
      line-height: 1;
      color: var(--color-content);
      background: var(--color-background-700);
      border: 1px solid var(--color-border-default);
      border-radius: var(--radius-sm);
      box-shadow: var(--shadow-xs);
    }
  }

  .batch-mode-notice {
    background: linear-gradient(
      135deg,
      rgba(var(--color-brand-500-rgb), 0.12) 0%,
      rgba(var(--color-brand-500-rgb), 0.08) 50%,
      rgba(var(--color-brand-500-rgb), 0.05) 100%
    );
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.2);
  }

  .multi-select-notice {
    background: linear-gradient(
      135deg,
      rgba(var(--color-success-primary-rgb), 0.12) 0%,
      rgba(var(--color-success-primary-rgb), 0.08) 50%,
      rgba(var(--color-success-primary-rgb), 0.05) 100%
    );
    border: 1px solid rgba(var(--color-success-primary-rgb), 0.3);
  }

  .normal-mode-notice {
    background: linear-gradient(
      135deg,
      rgba(var(--color-warning-primary-rgb), 0.08) 0%,
      rgba(var(--color-warning-primary-rgb), 0.05) 50%,
      rgba(var(--color-warning-primary-rgb), 0.03) 100%
    );
    border: 1px solid rgba(var(--color-warning-primary-rgb), 0.15);
  }

  @keyframes slideDown {
    from {
      opacity: 0;
      transform: translateY(-8px);
    }
    to {
      opacity: 1;
      transform: translateY(0);
    }
  }

  @keyframes pulse {
    0% {
      transform: scale(1);
    }
    50% {
      transform: scale(1.05);
    }
    100% {
      transform: scale(1);
    }
  }

  .tags-content {
    display: flex;
    flex-direction: column;
    gap: var(--space-md);
    flex: 1;
    min-height: 0;
  }

  @media (max-width: 768px) {
    .tag-manage-page {
    }

    .search-input-wrapper {
      width: 200px;
    }

    .page-header {
      padding: 0.875rem 1rem;

      > div:first-child {
        flex-direction: column;
        align-items: flex-start;
        gap: 0.75rem;
      }

      .icon-container {
        width: 2rem;
        height: 2rem;

        i {
          font-size: 1rem;
        }
      }
    }

    .tags-content {
      gap: var(--space-md);
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
  }
</style>
