/**
 * Tag Management Page Text - Cyber Style
 */
export const tagManage = {
  title: 'Tag Management',
  subtitle: 'Manage your tags, organize and sort your data units',
  search: 'Search tags...',
  create: 'Create Tag',
  batchMode: 'Batch Mode',
  exitBatchMode: 'Exit Batch Mode',
  loading: 'Loading tags...',
  batch: {
    viewFiles: 'View files with selected tags',
    delete: 'Batch Delete ({count})',
    merge: 'Merge Tags',
    cancel: 'Clear Selection',
    modeHint: 'Batch Operation Mode',
    modeDescription: 'Click tags to select/deselect, {count} tags selected',
  },
  multiSelect: {
    modeHint: 'Multi-select Mode',
    modeDescription: 'Hold {modifier} / Alt key and click tags to multi-select, {count} tags selected',
  },
  normalMode: {
    hint: 'Tip: Click tag to quickly search, hold {modifier} / Alt key to multi-select',
  },
  actions: {
    edit: 'Edit Tag',
    delete: 'Delete Tag',
  },
  dialog: {
    create: {
      title: 'Create Tag',
      submit: 'Create',
    },
    edit: {
      title: 'Edit Tag',
      submit: 'Update',
    },
    merge: {
      title: 'Merge Tags',
      targetLabel: 'Target Tag',
      targetPlaceholder: 'Please select target tag',
    },
    delete: {
      title: 'Delete Tag',
      message: 'Are you sure you want to delete tag "{name}"? This tag has been used {count} times, deleting it will remove this tag from all data units.',
      confirmText: 'Confirm Delete',
    },
    batchDelete: {
      title: 'Batch Delete Tags',
      message: 'Are you sure you want to delete the selected {count} tags? This operation is irreversible.',
      confirmText: 'Confirm Delete',
    },
    mergeConfirm: {
      title: 'Merge Tags',
      message: [
        'Are you sure you want to merge {count} tags into "{targetName}"?',
        'Source tags will be deleted, and their file associations will be transferred to the target tag.',
        'This operation is irreversible.',
      ],
      confirmText: 'Confirm Merge',
    },
  },
  form: {
    name: {
      label: 'Tag Name',
      placeholder: 'Please enter tag name',
      hint: 'Up to 50 characters, recommend using short and memorable names',
    },
    cancel: 'Cancel',
  },
  merge: {
    description: 'Merge selected tags into target tag. Source tags will be deleted, and their file associations will be transferred to the target tag.',
    selectedInfo: '{count} tags selected. After selecting target tag, the other {remaining} tags will be merged into target tag.',
    confirm: 'Confirm Merge',
    cancel: 'Cancel',
  },
  empty: {
    title: 'No matching tags found',
    description: 'Try adjusting filter conditions or creating new tags',
  },
  tagCloud: {
    groupCount: '{count} tags',
    groupUsage: 'Used {usage} times',
    letterNav: '{letter} ({count})',
    letterEmpty: '{letter} (None)',
    usageCount: '{count} uses',
  },
  toast: {
    nameRequired: 'Please enter tag name',
    updateSuccess: 'Tag updated successfully',
    createSuccess: 'Tag created successfully',
    deleteSuccess: 'Tag deleted successfully',
    selectRequired: 'Please select tags to delete first',
    mergeMinRequired: 'At least 2 tags must be selected to merge',
    mergeTargetRequired: 'Please select a target tag',
    mergeTargetInvalid: 'The target tag must be among the selected tags',
    mergeSuccess: 'Tags merged successfully',
    batchSelectRequired: 'Please select tags first',
  },
  source: {
    manual: 'Manual Creation',
    ai: 'AI Recognition',
    system: 'System Import',
    unknown: 'Unknown',
  },
}
