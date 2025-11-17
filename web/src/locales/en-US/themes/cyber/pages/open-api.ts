/**
 * Open API Page Text - Cyber Style
 */
export const openApi = {
  page: {
    title: 'Open API Management',
    subtitle: 'Manage random image API interfaces, provide data services for third-party applications',
    create: 'Create API',
    search: 'Search API name...',
  },
  loading: 'Loading...',
  footer: {
    total: 'Total {count} APIs',
  },
  empty: {
    title: 'No APIs Yet',
    description: 'No API interfaces created yet, click button above to start creating',
  },
  toast: {
    loadFailed: 'Load failed',
    createSuccess: 'Created successfully',
    createFailed: 'Creation failed',
    updateSuccess: 'Updated successfully',
    updateFailed: 'Update failed',
    deleteSuccess: 'Deleted successfully',
    deleteFailed: 'Deletion failed',
    copySuccess: 'API address copied',
    toggleSuccess: 'Status updated',
    toggleFailed: 'Status update failed',
  },
  dialog: {
    create: {
      title: 'Create Random Image API Node',
      cancel: 'Cancel',
      submit: 'Create',
      form: {
        name: {
          label: 'API Node Name',
          hint: 'Used to identify this API node for easy management',
        },
        scope: {
          label: 'Image Scope Node',
        },
        folder: {
          label: 'Select Folder Node',
          hint: 'Select a folder node, API will randomly return public images from this folder',
        },
        returnType: {
          label: 'Return Method Node',
          redirect: {
            label: '302 Redirect Node',
            desc: 'Returns real image address node, better performance, supports CDN acceleration',
          },
          direct: {
            label: 'Direct Return Image Node',
            desc: 'Each refresh is a new image node, URL remains fixed, suitable for website backgrounds/wallpaper carousels',
          },
        },
      },
      hint: 'After creation, a unique API link node will be generated, only returning publicly accessible images',
    },
    edit: {
      title: 'Edit API Config Node',
      cancel: 'Cancel',
      submit: 'Save',
      form: {
        nameLabel: 'API Node Name:',
        scope: {
          label: 'Image Scope Node',
        },
        folder: {
          label: 'Select Folder Node',
          hint: 'Select a folder node, API will randomly return public images from this folder',
        },
        returnType: {
          label: 'Return Method Node',
          redirect: {
            label: '302 Redirect Node',
            desc: 'Returns real image address node, better performance, supports CDN acceleration',
          },
          direct: {
            label: 'Direct Return Image Node',
            desc: 'Each refresh is a new image node, URL remains fixed, suitable for website backgrounds/wallpaper carousels',
          },
        },
      },
    },
    delete: {
      title: 'Delete API Node',
      message: 'Are you sure you want to delete API node "{name}"?',
      warning: 'Cannot be recovered after deletion node, please proceed with caution',
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
      placeholder: 'Select categories (multiple selection allowed)',
      all: 'All Categories',
    },
    tags: {
      label: 'Tag Filter',
      placeholder: 'Select tags (multiple selection allowed)',
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
      images: 'Data Units',
      calls: 'Calls',
    },
    filters: {
      allCategories: 'All Categories',
      allTags: 'All Tags',
    },
    lastCall: {
      never: 'Never Called',
      justNow: 'Just now',
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
