/**
 * Open API page texts
 */
export const openApi = {
  page: {
    title: 'Open API Management',
    subtitle: 'Manage your random image API endpoints, provide image services for third-party applications',
    create: 'Create API',
    search: 'Search API name...',
  },
  loading: 'Loading...',
  footer: {
    total: 'Total {count} APIs',
  },
  empty: {
    title: 'No APIs',
    description: 'You have not created any API endpoints yet, click the button above to get started',
  },
  toast: {
    loadFailed: 'Load failed',
    createSuccess: 'Created successfully',
    createFailed: 'Create failed',
    updateSuccess: 'Updated successfully',
    updateFailed: 'Update failed',
    deleteSuccess: 'Deleted successfully',
    deleteFailed: 'Delete failed',
    copySuccess: 'API address copied',
    toggleSuccess: 'Status updated',
    toggleFailed: 'Status update failed',
  },
  dialog: {
    create: {
      title: 'Create Random Image API',
      cancel: 'Cancel',
      submit: 'Create',
      form: {
        name: {
          label: 'API Name',
          hint: 'Used to identify this API for easy management',
        },
        scope: {
          label: 'Image Scope',
        },
        folder: {
          label: 'Select Folder',
          hint: 'Select a folder, API will randomly return from public images in that folder',
        },
        returnType: {
          label: 'Return Method',
          redirect: {
            label: '302 Redirect',
            desc: 'Returns actual image URL, better performance, supports CDN acceleration',
          },
          direct: {
            label: 'Direct Return Image',
            desc: 'Each refresh is a new image, URL stays fixed, suitable for website background/wallpaper rotation',
          },
        },
      },
      hint: 'After creation, a unique API link will be generated, only returns images with public access',
    },
    edit: {
      title: 'Edit API Configuration',
      cancel: 'Cancel',
      submit: 'Save',
      form: {
        nameLabel: 'API Name:',
        scope: {
          label: 'Image Scope',
        },
        folder: {
          label: 'Select Folder',
          hint: 'Select a folder, API will randomly return from public images in that folder',
        },
        returnType: {
          label: 'Return Method',
          redirect: {
            label: '302 Redirect',
            desc: 'Returns actual image URL, better performance, supports CDN acceleration',
          },
          direct: {
            label: 'Direct Return Image',
            desc: 'Each refresh is a new image, URL stays fixed, suitable for website background/wallpaper rotation',
          },
        },
      },
    },
    delete: {
      title: 'Delete API',
      message: 'Are you sure you want to delete API "{name}"?',
      warning: 'Cannot be recovered after deletion, please proceed with caution',
      cancel: 'Cancel',
      confirm: 'Confirm Delete',
    },
  },
  form: {
    name: {
      label: 'API Name',
      placeholder: 'Please enter API name',
      required: 'Please enter API name',
    },
    description: {
      label: 'API Description',
      placeholder: 'Please enter API description (optional)',
    },
    categories: {
      label: 'Category Filter',
      placeholder: 'Select categories (multiple)',
      all: 'All Categories',
    },
    tags: {
      label: 'Tag Filter',
      placeholder: 'Select tags (multiple)',
      all: 'All Tags',
    },
    status: {
      label: 'Status',
      enabled: 'Enabled',
      disabled: 'Disabled',
    },
  },
  card: {
    status: {
      enabled: 'Enabled',
      disabled: 'Disabled',
    },
    actions: {
      open: 'Open API',
      copy: 'Copy Address',
      edit: 'Edit',
      delete: 'Delete',
      toggle: 'Toggle Status',
    },
    stats: {
      categories: 'Categories',
      tags: 'Tags',
      images: 'Images',
      calls: 'Calls',
    },
    filters: {
      allCategories: 'All Categories',
      allTags: 'All Tags',
    },
    lastCall: {
      never: 'Never Called',
      justNow: 'Just Now',
      minutesAgo: '{minutes} minutes ago',
      hoursAgo: '{hours} hours ago',
      daysAgo: '{days} days ago',
    },
    returnType: {
      direct: 'Direct Return',
      redirect: '302 Redirect',
    },
    scope: 'Scope',
    scopeAll: 'All My Public Images',
    scopeAllPublic: 'All Public Images',
    scopeFolder: 'Specified Folder',
    returnMethod: 'Return Method',
    times: 'times',
  },
}
