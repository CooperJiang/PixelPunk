/**
 * Folder Management Page Text - Cyber Style
 */
export const folders = {
  rootFolder: 'Root Node',
  emptyState: {
    title: 'Current Data Directory Empty',
    description: 'You can upload data resources or create new data directories to get started',
    createButton: 'Create Data Directory',
  },
  breadcrumb: {
    label: 'Data Directory Navigation',
  },
  share: {
    createSuccess: 'Data share created successfully',
    createSuccessWithCopy: 'Data share created successfully, link copied to clipboard',
    manualCopy: 'Please manually copy share link',
    description: 'Share Description',
  },
  shareSelectTip: {
    message: 'Please select data directories or data units to share',
  },
  folderGrid: {
    sectionTitle: 'Data Directories',
    dragTip: 'Drag to Sort',
    fileCount: '{count} items',
    toast: {
      moveFailed: 'Failed to move data directory',
      moveToRootSuccess: 'Data directory moved to root node',
      moveSuccess: 'Data directory moved to "{name}"',
    },
    contextMenu: {
      root: 'Root Node',
      moveTo: 'Move to...',
      rename: 'Rename/Edit',
      delete: 'Delete',
    },
    actions: {
      selectAll: 'Select All',
      deselectAll: 'Deselect All',
      selectAllTitle: 'Select All Data Directories',
      deselectAllTitle: 'Deselect All Data Directories',
      invert: 'Invert Selection',
      invertTitle: 'Invert Data Directory Selection',
      deleteTitle: 'Delete Data Directory',
    },
    visibility: {
      toggleTitle: 'Click to toggle visibility',
      setPublic: 'Set as Public',
      setPrivate: 'Set as Private',
    },
  },
  createFolderDialog: {
    title: {
      create: 'New Data Directory',
      edit: 'Edit Data Directory',
    },
    form: {
      name: {
        label: 'Directory Name',
        placeholder: 'Please enter data directory name',
        required: 'Please enter data directory name',
        maxLength: 'Data directory name cannot exceed 30 characters',
      },
      description: {
        label: 'Description (Optional)',
        placeholder: 'Add data directory description',
      },
      permission: {
        label: 'Permission Settings',
        public: {
          name: 'Public',
          desc: 'Others can view content when visiting my homepage',
        },
        private: {
          name: 'Private',
          desc: 'This data directory will not be displayed to others',
        },
      },
    },
    actions: {
      create: 'Create',
      save: 'Save',
    },
    toast: {
      createSuccess: 'Data directory created successfully',
      updateSuccess: 'Data directory updated successfully',
      operationFailed: 'Data directory operation failed',
    },
  },
  fileGrid: {
    sectionTitle: 'Data Files',
    dragTip: 'Drag to Sort',
    selectedCount: '{count} selected',
    badge: {
      duplicate: 'Duplicate',
    },
    actions: {
      aiInfo: 'AI Analysis',
      preview: 'Preview',
      copyLink: 'Copy Link',
      download: 'Download',
      delete: 'Delete',
      batchMode: 'Batch Protocol Operation',
      batchDelete: 'Batch Delete',
      cancelBatch: 'Cancel Batch Protocol Operation',
      selectAll: 'Select All',
      deselectAll: 'Deselect All',
      selectAllTitle: 'Select All Data Units',
      deselectAllTitle: 'Deselect All Data Units',
      invert: 'Invert Selection',
      invertTitle: 'Invert Data Unit Selection',
      gridView: 'Grid View',
      listView: 'List View',
    },
    meta: {
      size: 'Data Unit Size',
      dimensions: 'Data Unit Dimensions',
      duplicateFile: 'Duplicate Data Unit',
    },
    listView: {
      columns: {
        preview: 'Preview',
        filename: 'Data Unit Name',
        size: 'Size',
        dimensions: 'Dimensions',
        date: 'Injection Time',
        actions: 'Actions',
      },
      empty: 'No data units yet',
    },
  },
  pageHeader: {
    title: 'My Data Directories',
    breadcrumb: {
      root: 'Root Node',
      navigation: 'Data Directory Navigation',
    },
    upload: {
      default: 'Inject Data Units',
      toFolder: 'Inject to "{name}"',
    },
    actions: {
      newFolder: 'New Data Directory',
    },
  },
  createShareDialog: {
    title: 'Create Share Link',
    sections: {
      basic: 'Basic Info Node',
      access: 'Access Control Node',
      advanced: 'Advanced Settings Node',
      content: 'Share Content Node',
    },
    form: {
      name: {
        label: 'Share Name',
        placeholder: 'Please enter share name (optional)',
      },
      description: {
        placeholder: 'Please enter share description (optional)',
      },
      password: {
        label: 'Access Password',
        placeholder: 'Please enter access password (optional)',
      },
      expiredDays: {
        label: 'Expiration Days Node',
        placeholder: '0 for never expire',
      },
      maxViews: {
        label: 'Max Access Count Node',
        placeholder: '0 for unlimited',
      },
      collectVisitor: {
        label: 'Collect Visitor Info Node',
        tip: 'A pop-up window node will appear when visitors browse, optional',
      },
      notification: {
        label: 'Access Notification Node',
        tip: 'Enable to receive email notification node when access count reaches specified number',
      },
      notificationThreshold: {
        label: 'Notification Threshold Node',
        placeholder: 'Set access count to trigger notification',
        tip: 'When access count reaches this value, system will send email notification node',
      },
    },
    emptyItems: 'Please select data directories or data units to share',
    cancel: 'Cancel',
    submit: 'Create Share Link',
    toast: {
      createSuccess: 'Share link created successfully',
      linkCopied: 'Share link created successfully, link copied to clipboard',
      copyFailed: 'Failed to copy to clipboard',
      copyManually: 'Please manually copy share link',
    },
  },
  deleteConfirmDialogs: {
    title: {
      folder: 'Delete Data Directory',
      file: 'Delete Data Unit',
      batchFile: 'Batch Delete Data Units',
    },
    message: {
      folder: 'Are you sure you want to delete data directory "{name}"?',
      file: 'Are you sure you want to delete data unit "{name}"?',
      batchFile: 'Are you sure you want to batch delete {count} resource nodes?',
    },
    warning: 'This operation is irreversible node, please proceed with caution.',
    confirmText: 'Confirm Delete',
  },
  fileContextMenu: {
    toast: {
      noFileSelected: 'No data unit selected',
      moveFailed: 'Failed to move data unit',
      moveToRootSuccess: 'Data unit moved to root node',
      moveToFolderSuccess: 'Data unit moved to "{name}"',
      loadFoldersFailed: 'Failed to load data directory list',
      loadSubfoldersFailed: 'Failed to load sub data directories',
      noValidImageLink: 'No valid original data link',
      noThumbLink: 'No thumbnail link to copy',
      thumbLinkCopied: 'Thumbnail link copied',
      copyFailed: 'Copy failed',
      originalLinkCopied: 'Original data link copied',
      noOriginalLink: 'No original data link to copy',
      shortLinkCopied: 'Short link copied',
      noShortLink: 'No short link to copy',
      fileDeleted: 'Data unit deleted',
      deleteFailed: 'Failed to delete data unit',
    },
    menu: {
      root: 'Root Node',
      preview: 'Preview',
      openNew: 'Open in New Window',
      copyLink: 'Copy Link',
      copyThumb: 'Copy Thumbnail',
      copyOriginal: 'Copy Original Data',
      copyShort: 'Copy Short Link',
      download: 'Download',
      moveTo: 'Move to...',
      delete: 'Delete',
    },
    confirm: {
      deleteFile: 'Are you sure you want to delete data unit "{name}"? This operation is irreversible.',
    },
  },
  folderManagement: {
    deleteSuccess: 'Data directory "{name}" deleted',
    visibility: {
      public: 'Public Protocol',
      private: 'Private Archive',
      switched: 'Data directory "{name}" set to {level}',
    },
  },
  dragSort: {
    folderSortUpdated: 'Data directory sort order updated',
    fileSortUpdated: 'Data unit sort order updated',
  },
  fileManagement: {
    copyLink: {
      protectedNotSupported: 'Protected data does not support external transfer, please view within terminal',
      publicLinkCopied: 'Public link copied to clipboard',
      privateLinkCopied: 'Private link copied to clipboard',
    },
    download: {
      success: 'Data unit "{name}" download complete',
      cancelled: 'Download cancelled or failed',
    },
    delete: {
      success: 'Data unit "{name}" deleted',
      selectFirst: 'Please select data units to delete first',
      batchSuccess: 'Successfully deleted {successCount} data units',
      batchPartial: 'Successfully deleted {successCount} data units, {failCount} failed',
      batchFailed: 'Delete failed, {failCount} data units failed to delete',
    },
    accessLevel: {
      public: 'Public Protocol',
      private: 'Private Archive',
      protected: 'Security Shield',
      switched: 'Data unit "{name}" set to {level}',
    },
  },
}
