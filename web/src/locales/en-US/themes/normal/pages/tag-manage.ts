/**
 * Tag management page texts
 */
export const tagManage = {
  title: 'Tag Management',
  subtitle: 'Manage your tags, organize and sort your resources',
  search: 'Search tags...',
  create: 'Create Tag',
  batchMode: 'Batch Mode',
  exitBatchMode: 'Exit Batch Mode',
  loading: 'Loading tags...',
  batch: {
    viewFiles: 'View Selected Tag Files',
    delete: 'Batch Delete ({count})',
    merge: 'Merge Tags',
    cancel: 'Cancel Selection',
    modeHint: 'Batch Operation Mode',
    modeDescription: 'Click tags to select/deselect, {count} tags selected',
  },
  multiSelect: {
    modeHint: 'Multi-select Mode',
    modeDescription: 'Hold {modifier} / Alt key and click tags for multi-select, {count} tags selected',
  },
  normalMode: {
    hint: 'Tip: Click tag to quickly search, hold {modifier} / Alt key for multi-select',
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
      message:
        'Are you sure you want to delete tag "{name}"? This tag has been used {count} times, deleting will remove this tag from all files.',
      confirmText: 'Confirm Delete',
    },
    batchDelete: {
      title: 'Batch Delete Tags',
      message: 'Are you sure you want to delete {count} selected tags? This operation cannot be undone.',
      confirmText: 'Confirm Delete',
    },
    mergeConfirm: {
      title: 'Merge Tags',
      message: [
        'Are you sure you want to merge {count} tags into "{targetName}"?',
        'Source tags will be deleted, their file associations will be transferred to the target tag.',
        'This operation cannot be undone.',
      ],
      confirmText: 'Confirm Merge',
    },
  },
  form: {
    name: {
      label: 'Tag Name',
      placeholder: 'Please enter tag name',
      hint: 'Maximum 50 characters, recommend using short and memorable names',
    },
    cancel: 'Cancel',
  },
  merge: {
    description:
      'Merge selected tags into target tag, source tags will be deleted, their file associations will be transferred to target tag.',
    selectedInfo: '{count} tags selected, after selecting target tag, the other {remaining} tags will be merged into target tag.',
    confirm: 'Confirm Merge',
    cancel: 'Cancel',
  },
  empty: {
    title: 'No Tags Found Matching Criteria',
    description: 'Try adjusting filter conditions or create new tags',
  },
  tagCloud: {
    groupCount: '{count} tags',
    groupUsage: 'Used {usage} times',
    letterNav: '{letter} ({count})',
    letterEmpty: '{letter} (None)',
    usageCount: 'Used {count} times',
  },
  toast: {
    nameRequired: 'Please enter tag name',
    updateSuccess: 'Tag updated successfully',
    createSuccess: 'Tag created successfully',
    deleteSuccess: 'Tag deleted successfully',
    selectRequired: 'Please select tags to delete first',
    mergeMinRequired: 'At least 2 tags must be selected to merge',
    mergeTargetRequired: 'Please select target tag',
    mergeTargetInvalid: 'Target tag must be among selected tags',
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
