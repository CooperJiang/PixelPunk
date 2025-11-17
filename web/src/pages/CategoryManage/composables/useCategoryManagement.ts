/**
 * 分类管理业务逻辑
 */
import { ref, reactive } from 'vue'
import { useToast } from '@/components/Toast/useToast'
import { useConfirmDialog } from '@/composables/useConfirmDialog'
import { useTexts } from '@/composables/useTexts'
import type { ImageCategory } from '@/api/types/category'
import { getCategoryList, updateCategory, deleteCategory, batchDelete, updateCategoryStatus } from '@/api/user/category'

export function useCategoryManagement() {
  const { $t } = useTexts()
  const toast = useToast()
  const { showConfirmDialog } = useConfirmDialog()

  const categories = ref<ImageCategory[]>([])
  const loading = ref(false)
  const saveLoading = ref(false)

  const totalCategories = ref(0)
  const currentPage = ref(1)
  const pageSize = ref(10)

  const showFilter = ref(false)
  const currentFilters = reactive({
    keyword: '',
    status: '',
    sort_by: 'sort_order',
    sort_order_dir: 'asc',
  })

  const showEditDialog = ref(false)
  const currentEditCategory = ref<ImageCategory | null>(null)

  const showCreateDialog = ref(false)

  const selectedCategoryIds = ref<number[]>([])

  const fetchCategoryList = async () => {
    loading.value = true

    const params = {
      page: currentPage.value,
      size: pageSize.value,
      keyword: currentFilters.keyword || undefined,
      status: currentFilters.status || undefined,
      sort_by: currentFilters.sort_by as any,
      sort_order_dir: currentFilters.sort_order_dir as 'asc' | 'desc',
    }

    try {
      const result = await getCategoryList(params)
      if (result.success) {
        const { data } = result
        categories.value = data.categories
        totalCategories.value = data.total

        selectedCategoryIds.value = selectedCategoryIds.value.filter((id) =>
          categories.value.some((category) => category.id === id)
        )
      }
    } catch {
      toast.error($t('category.toast.fetchError'))
    }

    loading.value = false
  }

  const handleFilter = (filters: typeof currentFilters) => {
    Object.assign(currentFilters, filters)
    currentPage.value = 1
    selectedCategoryIds.value = []
    fetchCategoryList()
    showFilter.value = false
  }

  const handlePageChange = (page: number) => {
    currentPage.value = page
    fetchCategoryList()
  }

  const handlePageSizeChange = (size: number) => {
    pageSize.value = size
    currentPage.value = 1
    fetchCategoryList()
  }

  const editCategory = (category: ImageCategory) => {
    currentEditCategory.value = JSON.parse(JSON.stringify(category))
    showEditDialog.value = true
  }

  const saveCategoryChanges = async () => {
    if (!currentEditCategory.value) {
      return
    }

    saveLoading.value = true

    try {
      const result = await updateCategory({
        id: currentEditCategory.value.id,
        name: currentEditCategory.value.name,
        description: currentEditCategory.value.description,
        sort_order: currentEditCategory.value.sort_order,
      })

      if (result.success) {
        toast.success($t('category.toast.updateSuccess'))
        showEditDialog.value = false
        fetchCategoryList()
      }
    } catch {
      toast.error($t('category.toast.updateError'))
    }
    saveLoading.value = false
  }

  const toggleCategoryStatus = async (category: ImageCategory) => {
    const newStatus = category.status === 'active' ? 'archived' : 'active'
    const isArchiving = newStatus === 'archived'

    const confirmed = await showConfirmDialog({
      title: $t(`category.actions.${isArchiving ? 'archive' : 'activate'}`),
      message: $t(`category.batch.${isArchiving ? 'archiveConfirm' : 'activateConfirm'}`, { count: 1 }),
      type: 'warning',
    })

    if (!confirmed) {
      return
    }

    try {
      const result = await updateCategoryStatus({
        id: category.id,
        status: newStatus,
      })
      if (result.success) {
        toast.success($t(`category.toast.${isArchiving ? 'archiveSuccess' : 'activateSuccess'}`))
        fetchCategoryList()
      }
    } catch {
      toast.error($t(`category.toast.${isArchiving ? 'archiveError' : 'activateError'}`))
    }
  }

  const deleteCategoryItem = async (id: number) => {
    try {
      const result = await deleteCategory({ id })
      if (result.success) {
        toast.success($t('category.toast.deleteSuccess'))
        fetchCategoryList()
      }
    } catch {
      toast.error($t('category.toast.deleteError'))
    }
  }

  const handleBatchDelete = async () => {
    if (selectedCategoryIds.value.length === 0) {
      toast.warning($t('category.toast.selectWarning'))
      return
    }

    const confirmed = await showConfirmDialog({
      title: $t('category.batch.delete'),
      message: $t('category.batch.deleteConfirm', { count: selectedCategoryIds.value.length }),
      type: 'danger',
      confirmText: $t('common.confirmDelete'),
    })

    if (!confirmed) {
      return
    }

    try {
      const result = await batchDelete({ ids: selectedCategoryIds.value })
      if (result.success) {
        toast.success($t('category.toast.batchDeleteSuccess', { count: result.data.deleted_count }))
        selectedCategoryIds.value = []
        fetchCategoryList()
      }
    } catch {
      toast.error($t('category.toast.batchDeleteError'))
    }
  }

  const handleRefresh = () => {
    selectedCategoryIds.value = []
    fetchCategoryList()
  }

  const handleCategoryCreated = () => {
    toast.success($t('category.toast.createSuccess'))
    fetchCategoryList()
  }

  const handleBatchCompleted = () => {
    selectedCategoryIds.value = []
    fetchCategoryList()
  }

  const handleSelectionChange = (selectedRows: (string | number)[]) => {
    selectedCategoryIds.value = selectedRows.map((id) => Number(id))
  }

  return {
    categories,
    loading,
    saveLoading,
    totalCategories,
    currentPage,
    pageSize,
    showFilter,
    currentFilters,
    showEditDialog,
    currentEditCategory,
    showCreateDialog,
    selectedCategoryIds,

    fetchCategoryList,
    handleFilter,
    handlePageChange,
    handlePageSizeChange,
    editCategory,
    saveCategoryChanges,
    toggleCategoryStatus,
    deleteCategoryItem,
    handleBatchDelete,
    handleRefresh,
    handleCategoryCreated,
    handleBatchCompleted,
    handleSelectionChange,
  }
}
