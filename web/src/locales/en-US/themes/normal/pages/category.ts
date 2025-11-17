/**
 * Category management page texts
 */
export const category = {
  page: {
    title: 'Category Management',
    subtitle: 'Manage your file categories, organize and sort your resources',
    create: 'New Category',
    refresh: 'Refresh',
    filter: 'Filter',
    hideFilter: 'Hide Filter',
  },
  table: {
    columns: {
      name: 'Category Name',
      description: 'Description',
      source: 'Source',
      fileCount: 'File Count',
      sortOrder: 'Sort Order',
      status: 'Status',
      createdAt: 'Created At',
      actions: 'Actions',
    },
    noDescription: 'No description',
    clickToView: 'Click to view all files in this category',
    loading: 'Loading category data...',
  },
  source: {
    system: 'System',
    user: 'User',
    imported: 'Imported',
  },
  status: {
    active: 'Active',
    archived: 'Archived',
  },
  actions: {
    edit: 'Edit Category',
    archive: 'Archive Category',
    activate: 'Activate Category',
    delete: 'Delete Category',
  },
  empty: {
    title: 'No Category Data',
    description: 'No categories created yet. Click "New Category" to create your first category',
    action: 'New Category',
  },
  dialog: {
    edit: {
      title: 'Edit Category',
      cancel: 'Cancel',
      save: 'Save Changes',
    },
    delete: {
      title: 'Delete Category',
      message: 'Are you sure you want to delete this category?',
    },
    create: {
      title: 'New Category',
      cancel: 'Cancel',
      create: 'Create Category',
    },
  },
  form: {
    name: {
      label: 'Category Name',
      placeholder: 'Please enter category name',
      required: 'Please enter category name',
    },
    description: {
      label: 'Category Description',
      placeholder: 'Please enter category description (optional)',
    },
    sortOrder: {
      label: 'Sort Order',
      placeholder: 'Please enter sort order',
      hint: 'Lower values appear first in the list',
    },
    status: {
      label: 'Status',
      active: 'Active',
      archived: 'Archived',
    },
  },
  filter: {
    title: 'Filter Conditions',
    sortBy: 'Sort By',
    status: {
      label: 'Status',
      all: 'All Status',
      active: 'Active',
      archived: 'Archived',
    },
    source: {
      label: 'Source',
      all: 'All',
      system: 'System',
      user: 'User',
      imported: 'Imported',
    },
    search: {
      placeholder: 'Search category name or description...',
    },
    apply: 'Apply',
    reset: 'Reset',
  },
  batch: {
    selected: '{count} items selected',
    clearSelection: 'Clear Selection',
    clear: 'Clear Selection',
    delete: 'Batch Delete',
    archive: 'Batch Archive',
    activate: 'Batch Activate',
    deleteConfirm: 'Are you sure you want to delete {count} selected categories?',
    archiveConfirm: 'Are you sure you want to archive {count} selected categories?',
    activateConfirm: 'Are you sure you want to activate {count} selected categories?',
  },
  toast: {
    createSuccess: 'Category created successfully',
    createError: 'Failed to create category',
    createWarning: 'Please enter category name',
    updateSuccess: 'Category updated successfully',
    updateError: 'Failed to update category',
    deleteSuccess: 'Category deleted successfully',
    deleteError: 'Failed to delete category',
    archiveSuccess: 'Archived successfully',
    activateSuccess: 'Activated successfully',
    archiveError: 'Failed to archive',
    activateError: 'Failed to activate',
    batchDeleteSuccess: 'Successfully deleted {count} categories',
    batchDeleteError: 'Batch delete failed',
    batchArchiveSuccess: 'Batch archive successful',
    batchActivateSuccess: 'Batch activate successful',
    fetchError: 'Failed to fetch category list',
    selectWarning: 'Please select categories to delete first',
  },
  sourceInfo: {
    manual: 'Manual Creation',
    aiSuggestion: 'AI Suggestion',
    system: 'System Template',
    imported: 'Batch Import',
  },
  statusInfo: {
    active: 'Active',
    archived: 'Archived',
  },
  placeholders: {
    search: 'Search category name or description...',
    status: 'Status',
    name: 'Please enter category name',
    description: 'Please enter category description (optional)',
    sortOrder: 'Lower values appear first',
  },
}
