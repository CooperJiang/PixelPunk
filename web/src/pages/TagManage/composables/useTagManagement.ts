/**
 * 标签管理业务逻辑
 */
import { ref, reactive, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useToast } from '@/components/Toast/useToast'
import { useConfirmDialog } from '@/composables/useConfirmDialog'
import { useTexts } from '@/composables/useTexts'
import {
  getTagList,
  createTag,
  updateTag,
  deleteTag,
  batchDeleteTags,
  mergeTags,
  getTagStats,
  type TagStats,
  type UserTag,
} from '@/api/user/tag'

export function useTagManagement() {
  const toast = useToast()
  const router = useRouter()
  const { showConfirmDialog } = useConfirmDialog()
  const { $t } = useTexts()

  const isLoading = ref(false)
  const allTags = ref<UserTag[]>([])
  const totalTags = ref(0)
  const selectedTagIds = ref<number[]>([])
  const stats = ref<TagStats | null>(null)
  const searchKeyword = ref('')

  const showCreateModal = ref(false)
  const showEditModal = ref(false)
  const showMergeModal = ref(false)

  const editingTag = ref<UserTag | null>(null)

  const isBatchMode = ref(false)

  const tagForm = reactive({
    name: '',
  })

  const mergeForm = reactive({
    targetId: null as number | null,
  })

  const showFormModal = computed({
    get: () => showCreateModal.value || showEditModal.value,
    set: (value) => {
      if (!value) {
        closeFormModal()
      }
    },
  })

  const filteredTags = computed(() => {
    if (!searchKeyword.value.trim()) {
      return allTags.value
    }
    const keyword = searchKeyword.value.toLowerCase().trim()
    return allTags.value.filter((tag) => tag.name.toLowerCase().includes(keyword))
  })

  const selectedTags = computed(() => {
    return allTags.value.filter((tag) => selectedTagIds.value.includes(tag.id))
  })

  const mergeTargetOptions = computed(() => {
    return selectedTags.value.map((tag) => ({
      label: `${tag.name} (${$t('tagManage.tagCloud.usageCount', { count: tag.file_count })})`,
      value: tag.id,
    }))
  })

  const loadTags = async () => {
    try {
      isLoading.value = true
      const response = await getTagList({
        keyword: '',
        page: 1,
        size: 10000,
      })
      allTags.value = response.data.items || []
      totalTags.value = response.data.total

      const newIds = allTags.value.map((tag) => tag.id)
      selectedTagIds.value = selectedTagIds.value.filter((id) => newIds.includes(id))
    } catch (_error: any) {
    } finally {
      isLoading.value = false
    }
  }

  const loadStats = async () => {
    try {
      const response = await getTagStats()
      stats.value = response.data
    } catch (_error: any) {
      console.error('Failed to fetch tag stats:', _error)
    }
  }

  const openCreateModal = () => {
    tagForm.name = ''
    showCreateModal.value = true
  }

  const openEditModal = (tag: UserTag) => {
    editingTag.value = tag
    tagForm.name = tag.name
    showEditModal.value = true
  }

  const closeFormModal = () => {
    showCreateModal.value = false
    showEditModal.value = false
    editingTag.value = null
    tagForm.name = ''
  }

  const submitTagForm = async () => {
    if (!tagForm.name.trim()) {
      toast.warning($t('tagManage.toast.nameRequired'))
      return
    }

    try {
      if (showEditModal.value && editingTag.value) {
        await updateTag({
          id: editingTag.value.id,
          name: tagForm.name.trim(),
        })
        toast.success($t('tagManage.toast.updateSuccess'))
      } else {
        await createTag({
          name: tagForm.name.trim(),
        })
        toast.success($t('tagManage.toast.createSuccess'))
      }

      closeFormModal()
      loadTags()
      loadStats()
    } catch (_error: any) {}
  }

  const handleDelete = async (tag: UserTag) => {
    const confirmed = await showConfirmDialog({
      title: $t('tagManage.dialog.delete.title'),
      message: $t('tagManage.dialog.delete.message', { name: tag.name, count: tag.file_count }),
      type: 'danger',
      confirmText: $t('tagManage.dialog.delete.confirmText'),
    })

    if (!confirmed) return

    try {
      await deleteTag({ id: tag.id })
      toast.success($t('tagManage.toast.deleteSuccess'))
      loadTags()
      loadStats()

      const index = selectedTagIds.value.indexOf(tag.id)
      if (index !== -1) {
        selectedTagIds.value.splice(index, 1)
      }
    } catch (_error: any) {}
  }

  const toggleTagSelection = (tag: UserTag) => {
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

  const handleBatchDelete = async () => {
    if (selectedTagIds.value.length === 0) {
      toast.warning($t('tagManage.toast.selectRequired'))
      return
    }

    const confirmed = await showConfirmDialog({
      title: $t('tagManage.dialog.batchDelete.title'),
      message: $t('tagManage.dialog.batchDelete.message', { count: selectedTagIds.value.length }),
      type: 'danger',
      confirmText: $t('tagManage.dialog.batchDelete.confirmText'),
    })

    if (!confirmed) return

    try {
      const result = await batchDeleteTags({ ids: selectedTagIds.value })
      toast.success($t('tagManage.toast.batchDeleteSuccess', { count: result.data.deleted_count }))
      selectedTagIds.value = []
      loadTags()
      loadStats()
    } catch (_error: any) {}
  }

  const openMergeModal = () => {
    if (selectedTagIds.value.length < 2) {
      toast.warning($t('tagManage.toast.mergeMinRequired'))
      return
    }
    mergeForm.targetId = null
    showMergeModal.value = true
  }

  const submitMerge = async () => {
    if (!mergeForm.targetId) {
      toast.warning($t('tagManage.toast.mergeTargetRequired'))
      return
    }

    const targetId = Number(mergeForm.targetId)
    if (!selectedTagIds.value.includes(targetId)) {
      toast.warning($t('tagManage.toast.mergeTargetInvalid'))
      return
    }

    const sourceIds = selectedTagIds.value.filter((id) => id !== targetId)
    const targetTag = selectedTags.value.find((tag) => tag.id === targetId)

    const confirmed = await showConfirmDialog({
      title: $t('tagManage.dialog.mergeConfirm.title'),
      message: [
        $t('tagManage.dialog.mergeConfirm.message.0', { count: sourceIds.length, targetName: targetTag?.name }),
        $t('tagManage.dialog.mergeConfirm.message.1'),
        $t('tagManage.dialog.mergeConfirm.message.2'),
      ],
      type: 'warning',
      confirmText: $t('tagManage.dialog.mergeConfirm.confirmText'),
    })

    if (!confirmed) return

    try {
      await mergeTags({
        source_ids: sourceIds,
        target_id: targetId,
      })
      toast.success($t('tagManage.toast.mergeSuccess'))
      selectedTagIds.value = []
      loadTags()
      loadStats()
    } catch (_error: any) {}
  }

  const toggleBatchMode = () => {
    isBatchMode.value = !isBatchMode.value
    if (!isBatchMode.value) {
      selectedTagIds.value = []
    }
  }

  const viewSelectedTagsFiles = () => {
    if (selectedTagIds.value.length === 0) {
      toast.warning($t('tagManage.toast.batchSelectRequired'))
      return
    }

    const tagIds = selectedTagIds.value.join(',')
    const tagNames = selectedTags.value.map((tag) => tag.name).join(',')

    router.push({
      path: '/resource',
      query: {
        tags: tagIds, // 传递标签ID
        tag_names: tagNames, // 传递标签名称用于显示
      },
    })
  }

  return {
    isLoading,
    allTags,
    filteredTags,
    totalTags,
    selectedTagIds,
    stats,
    searchKeyword,
    showCreateModal,
    showEditModal,
    showMergeModal,
    showFormModal,
    editingTag,
    isBatchMode,
    tagForm,
    mergeForm,
    selectedTags,
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
  }
}
