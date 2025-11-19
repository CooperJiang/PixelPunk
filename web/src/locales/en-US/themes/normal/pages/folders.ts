/**
 * Folder management page texts
 */
export const folders = {
  rootFolder: 'Root',
  emptyState: {
    title: 'Current Folder is Empty',
    description: 'You can upload resources or create new folders to get started',
    createButton: 'Create Folder',
  },
  breadcrumb: {
    label: 'Folder Navigation',
  },
  share: {
    createSuccess: 'Share created successfully',
    createSuccessWithCopy: 'Share created successfully, link copied to clipboard',
    manualCopy: 'Please manually copy the share link',
    description: 'Share Description',
  },
  shareSelectTip: {
    message: 'Please select folders or images to share',
  },
  folderGrid: {
    sectionTitle: 'Folders',
    dragTip: 'Drag to sort',
    fileCount: '{count} files',
    toast: {
      moveFailed: 'Failed to move folder',
      moveToRootSuccess: 'Folder moved to root',
      moveSuccess: 'Folder moved to "{name}"',
    },
    contextMenu: {
      root: 'Root',
      moveTo: 'Move to...',
      rename: 'Rename/Edit',
      delete: 'Delete',
    },
    actions: {
      selectAll: 'Select All',
      deselectAll: 'Deselect All',
      selectAllTitle: 'Select All Folders',
      deselectAllTitle: 'Deselect All Folders',
      invert: 'Invert',
      invertTitle: 'Invert Folders',
      deleteTitle: 'Delete Folder',
    },
    visibility: {
      toggleTitle: 'Click to toggle visibility',
      setPublic: 'Set to Public',
      setPrivate: 'Set to Private',
    },
  },
  createFolderDialog: {
    title: {
      create: 'New Folder',
      edit: 'Edit Folder',
    },
    form: {
      name: {
        label: 'Folder Name',
        placeholder: 'Please enter folder name',
        required: 'Please enter folder name',
        maxLength: 'Folder name cannot exceed 30 characters',
      },
      description: {
        label: 'Description (Optional)',
        placeholder: 'Add folder description',
      },
      permission: {
        label: 'Permission Settings',
        public: {
          name: 'Public',
          desc: 'Others can view content when visiting my homepage',
        },
        private: {
          name: 'Private',
          desc: 'Do not show this folder to others',
        },
      },
    },
    actions: {
      create: 'Create',
      save: 'Save',
    },
    toast: {
      createSuccess: 'Folder created successfully',
      updateSuccess: 'Folder updated successfully',
      operationFailed: 'Folder operation failed',
    },
  },
  fileGrid: {
    sectionTitle: 'Files',
    dragTip: 'Drag to sort',
    selectedCount: '{count} files selected',
    badge: {
      duplicate: 'Duplicate',
    },
    actions: {
      aiInfo: 'AI Info',
      preview: 'Preview',
      copyLink: 'Copy Link',
      download: 'Download',
      delete: 'Delete',
      batchMode: 'Batch Operations',
      batchDelete: 'Batch Delete',
      cancelBatch: 'Cancel Batch Operations',
      selectAll: 'Select All',
      deselectAll: 'Deselect All',
      selectAllTitle: 'Select All Files',
      deselectAllTitle: 'Deselect All Files',
      invert: 'Invert',
      invertTitle: 'Invert Files',
      gridView: 'Grid View',
      listView: 'List View',
    },
    meta: {
      size: 'File Size',
      dimensions: 'File Dimensions',
      duplicateFile: 'Duplicate File',
    },
    listView: {
      columns: {
        preview: 'Preview',
        filename: 'Filename',
        size: 'Size',
        dimensions: 'Dimensions',
        date: 'Upload Time',
        actions: 'Actions',
      },
      empty: 'No files',
    },
  },
  pageHeader: {
    title: 'My Folders',
    breadcrumb: {
      root: 'Root',
      navigation: 'Folder Navigation',
    },
    upload: {
      default: 'Upload Files',
      toFolder: 'Upload to "{name}"',
    },
    actions: {
      newFolder: 'New Folder',
    },
  },
  createShareDialog: {
    title: 'Create Share',
    sections: {
      basic: 'Basic Information',
      access: 'Access Control',
      advanced: 'Advanced Settings',
      content: 'Share Content',
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
        label: 'Expiration Days',
        placeholder: '0 means never expires',
      },
      maxViews: {
        label: 'Max Views',
        placeholder: '0 means unlimited',
      },
      collectVisitor: {
        label: 'Collect Visitor Information',
        tip: 'Visitors will see an information form popup when browsing, not required',
      },
      notification: {
        label: 'Access Notification',
        tip: 'When enabled, you will receive email notifications when views reach the specified count',
      },
      notificationThreshold: {
        label: 'Notification Threshold',
        placeholder: 'Set when to notify based on view count',
        tip: 'When view count reaches this value, system will send email notification',
      },
    },
    emptyItems: 'Please select folders or files to share',
    cancel: 'Cancel',
    submit: 'Create Share',
    toast: {
      createSuccess: 'Share created successfully',
      linkCopied: 'Share created successfully, link copied to clipboard',
      copyFailed: 'Failed to copy to clipboard',
      copyManually: 'Please manually copy the share link',
    },
  },
  deleteConfirmDialogs: {
    title: {
      folder: 'Delete Folder',
      file: 'Delete File',
      batchFile: 'Batch Delete Files',
    },
    message: {
      folder: 'Are you sure you want to delete folder "{name}"?',
      file: 'Are you sure you want to delete file "{name}"?',
      batchFile: 'Are you sure you want to batch delete {count} resources?',
    },
    warning: 'This operation cannot be undone, please proceed with caution.',
    confirmText: 'Confirm Delete',
  },
  fileContextMenu: {
    toast: {
      noFileSelected: 'No file selected',
      moveFailed: 'Failed to move file',
      moveToRootSuccess: 'File moved to root',
      moveToFolderSuccess: 'File moved to "{name}"',
      loadFoldersFailed: 'Failed to load folder list',
      loadSubfoldersFailed: 'Failed to load subfolders',
      noValidImageLink: 'No valid original image link',
      noThumbLink: 'No thumbnail link to copy',
      thumbLinkCopied: 'Thumbnail link copied',
      copyFailed: 'Copy failed',
      originalLinkCopied: 'Original image link copied',
      noOriginalLink: 'No original image link to copy',
      shortLinkCopied: 'Short link copied',
      noShortLink: 'No short link to copy',
      markdownLinkCopied: 'Markdown link copied',
      htmlLinkCopied: 'HTML tag copied',
      fileDeleted: 'File deleted',
      deleteFailed: 'Failed to delete file',
    },
    menu: {
      root: 'Root',
      preview: 'Preview',
      openNew: 'Open in New Window',
      copyLink: 'Copy Link',
      copyThumb: 'Copy Thumbnail',
      copyOriginal: 'Copy Original',
      copyShort: 'Copy Short Link',
      copyMarkdown: 'Copy Markdown',
      copyHtml: 'Copy HTML',
      download: 'Download',
      moveTo: 'Move to...',
      delete: 'Delete',
    },
    confirm: {
      deleteFile: 'Are you sure you want to delete file "{name}"? This operation cannot be undone.',
    },
  },
  folderManagement: {
    deleteSuccess: 'Folder "{name}" deleted',
    visibility: {
      public: 'Public',
      private: 'Private',
      switched: 'Folder "{name}" set to {level}',
    },
  },
  dragSort: {
    folderSortUpdated: 'Folder sort updated',
    fileSortUpdated: 'File sort updated',
  },
  fileManagement: {
    copyLink: {
      protectedNotSupported: 'Protected files do not support external sharing, please view within the application',
      publicLinkCopied: 'Public link copied to clipboard',
      privateLinkCopied: 'Private link copied to clipboard',
    },
    download: {
      success: 'File "{name}" download completed',
      cancelled: 'Download cancelled or failed',
    },
    delete: {
      success: 'File "{name}" deleted',
      selectFirst: 'Please select files to delete first',
      batchSuccess: 'Successfully deleted {successCount} files',
      batchPartial: 'Successfully deleted {successCount} files, {failCount} failed',
      batchFailed: 'Delete failed, {failCount} files failed to delete',
    },
    accessLevel: {
      public: 'Public',
      private: 'Private',
      protected: 'Protected',
      switched: 'File "{name}" set to {level}',
    },
  },
}
