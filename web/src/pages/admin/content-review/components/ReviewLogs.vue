<script setup lang="ts">
  import { getCurrentLocale } from '@/utils/locale'
  import { computed, h, onMounted, reactive, ref, watch } from 'vue'
  import { useToast } from '@/components/Toast/useToast'
  import { useTexts } from '@/composables/useTexts'
  import {
    batchHardDeleteReviewedImages,
    batchRestoreReviewedImages,
    getReviewLogs,
    hardDeleteReviewedImage,
    restoreReviewedImage,
    type ReviewLog,
    type ReviewLogQuery,
  } from '@/api/admin/content-review'

  const { $t } = useTexts()

  defineOptions({
    name: 'ReviewLogs',
  })

  const props = withDefaults(
    defineProps<{
      isActive?: boolean
    }>(),
    {
      isActive: false,
    }
  )

  /* 工具 */
  const toast = useToast()

  /* 数据状态 */
  const isLoading = ref(false)
  const showSkeleton = ref(false)
  const logs = ref<ReviewLog[]>([])
  const pagination = ref<{ page: number; page_size: number; total: number; total_page: number }>({
    page: 1,
    page_size: 10,
    total: 0,
    total_page: 1,
  })

  /* 选择模式状态 */
  const selectMode = ref(false)
  const selectedRows = ref<(string | number)[]>([])

  /* 骨架屏表格列配置 - 使用 getter 函数避免打包时的循环依赖 */
  const getSkeletonColumns = () => [
    {
      type: 'image' as const,
      width: '60px',
      align: 'center' as const,
      title: $t('admin.contentReview.logs.columns.file'),
      imageSize: 'small' as const,
    },
    {
      type: 'text' as const,
      width: '180px',
      align: 'center' as const,
      title: $t('admin.contentReview.logs.columns.fileInfo'),
      textWidth: '70%',
    },
    {
      type: 'badge' as const,
      width: '120px',
      align: 'center' as const,
      title: $t('admin.contentReview.logs.columns.auditAction'),
      badgeStyle: 'normal' as const,
    },
    {
      type: 'badge' as const,
      width: '120px',
      align: 'center' as const,
      title: $t('admin.contentReview.logs.columns.deleteStatus'),
      badgeStyle: 'warning' as const,
    },
    {
      type: 'text' as const,
      width: '120px',
      align: 'center' as const,
      title: $t('admin.contentReview.logs.columns.auditor'),
      textWidth: '60%',
    },
    {
      type: 'text' as const,
      width: '140px',
      align: 'center' as const,
      title: $t('admin.contentReview.logs.columns.auditTime'),
      textWidth: '80%',
    },
    {
      type: 'text' as const,
      width: '120px',
      align: 'center' as const,
      title: $t('admin.contentReview.logs.columns.reason'),
      textWidth: '75%',
    },
    {
      type: 'number' as const,
      width: '100px',
      align: 'center' as const,
      title: $t('admin.contentReview.logs.columns.nsfwScore'),
    },
    {
      type: 'button' as const,
      width: '120px',
      align: 'center' as const,
      title: $t('admin.contentReview.logs.columns.operations'),
    },
  ]

  /* 表格列定义 - 使用 getter 函数避免打包时的循环依赖 */
  const getTableColumns = () => [
    {
      key: 'image',
      title: $t('admin.contentReview.logs.columns.file'),
      width: 60,
      align: 'center',
      render: (_value: unknown, record: ReviewLog) =>
        h('div', { class: 'image-cell' }, [
          h('div', { class: 'image-wrapper' }, [
            record.image && (record.image.full_thumb_url || record.image.thumb_url)
              ? h('img', {
                  src: record.image.full_thumb_url || record.image.thumb_url,
                  alt: record.image.original_name,
                  class: 'table-image-thumb',
                  onError: handleImageError,
                })
              : h('div', { class: 'image-placeholder' }, [h('i', { class: 'fas fa-image' })]),
          ]),
        ]),
    },
    {
      key: 'imageInfo',
      title: $t('admin.contentReview.logs.columns.fileInfo'),
      width: 180,
      ellipsis: true,
      align: 'center',
      render: (_value: unknown, record: ReviewLog) =>
        h('div', { class: 'image-info-cell' }, [
          h(
            'div',
            {
              class: 'image-name',
              title: record.image?.original_name || record.file_id,
            },
            record.image?.original_name || $t('admin.contentReview.logs.text.fileId', { id: record.file_id })
          ),
          record.image
            ? h(
                'div',
                { class: 'image-meta' },
                [
                  record.image.size_formatted
                    ? h('span', { class: 'meta-item' }, [
                        h('i', { class: 'fas fa-weight mr-1 text-xs text-brand-700' }),
                        record.image.size_formatted,
                      ])
                    : null,
                  record.image.width && record.image.height
                    ? h('span', { class: 'meta-item' }, [
                        h('i', { class: 'fas fa-ruler-combined mr-1 text-xs text-brand-700' }),
                        `${record.image.width}×${record.image.height}`,
                      ])
                    : null,
                ].filter(Boolean)
              )
            : h('div', { class: 'image-meta' }, [
                h('span', { class: 'meta-item no-info' }, [
                  h('i', { class: 'fas fa-exclamation-triangle mr-1 text-xs text-yellow-400/70' }),
                  $t('admin.contentReview.logs.text.fileInfoDeleted'),
                ]),
              ]),
        ]),
    },
    {
      key: 'action',
      dataIndex: 'action',
      title: $t('admin.contentReview.logs.columns.auditAction'),
      width: 120,
      align: 'center',
      sortable: true,
      render: (value: string, record: ReviewLog) =>
        h('div', { class: 'action-cell' }, [
          h(
            'span',
            {
              class: `inline-flex items-center px-2 py-1 rounded text-xs font-medium ${
                record.action === 'approve'
                  ? 'bg-green-500/20 text-green-400 border border-green-500/30'
                  : 'bg-red-500/20 text-red-400 border border-red-500/30'
              }`,
            },
            [
              h('i', {
                class: record.action === 'approve' ? 'fas fa-check-circle mr-1' : 'fas fa-times-circle mr-1',
              }),
              record.action === 'approve'
                ? $t('admin.contentReview.logs.actions.approve')
                : $t('admin.contentReview.logs.actions.reject'),
            ]
          ),
        ]),
    },
    {
      key: 'softDelete',
      title: $t('admin.contentReview.logs.columns.deleteStatus'),
      width: 120,
      align: 'center',
      sortable: true,
      render: (_value: unknown, record: ReviewLog) =>
        h('div', { class: 'soft-delete-cell' }, [
          record.action === 'reject' && record.delete_type
            ? h(
                'span',
                {
                  class: `inline-flex items-center px-2 py-1 rounded text-xs font-medium ${getDeleteStatusVariant(record.delete_type, record.image?.status)}`,
                },
                [
                  h('i', {
                    class: `${getDeleteStatusIcon(record.delete_type, record.image?.status)} mr-1`,
                  }),
                  getDeleteStatusText(record.delete_type, record.image?.status),
                ]
              )
            : h('div', { class: 'no-delete-status' }, [h('span', { class: 'status-placeholder' }, '-')]),
        ]),
    },
    {
      key: 'auditor',
      dataIndex: 'auditor',
      title: $t('admin.contentReview.logs.columns.auditor'),
      width: 120,
      align: 'center',
      sortable: true,
      render: (_value: unknown, record: ReviewLog) =>
        h('div', { class: 'auditor-cell' }, [
          h('div', { class: 'auditor-name' }, record.auditor?.username || $t('admin.contentReview.logs.text.unknownAuditor')),
        ]),
    },
    {
      key: 'auditTime',
      dataIndex: 'created_at',
      title: $t('admin.contentReview.logs.columns.auditTime'),
      width: 140,
      align: 'center',
      sortable: true,
      render: (value: string, record: ReviewLog) =>
        h('div', { class: 'audit-time-cell' }, [h('div', { class: 'audit-time' }, formatDateTime(record.created_at))]),
    },
    {
      key: 'reason',
      dataIndex: 'reason',
      title: $t('admin.contentReview.logs.columns.reason'),
      width: 120,
      ellipsis: true,
      align: 'center',
      render: (value: string, record: ReviewLog) =>
        h('div', { class: 'reason-cell' }, [
          record.reason
            ? h(
                'div',
                {
                  class: 'reason-text',
                  title: record.reason,
                },
                truncateText(record.reason, 40)
              )
            : h('div', { class: 'no-reason' }, $t('admin.contentReview.logs.text.noDescription')),
        ]),
    },
    {
      key: 'nsfw',
      dataIndex: 'nsfw_score',
      title: $t('admin.contentReview.logs.columns.nsfwScore'),
      width: 100,
      align: 'center',
      sortable: true,
      render: (value: number, record: ReviewLog) =>
        h('div', { class: 'nsfw-cell' }, [
          record.nsfw_score !== null
            ? h('div', { class: 'nsfw-score' }, [
                h(
                  'span',
                  {
                    class: `score-value ${getNsfwClass(record.nsfw_score)}`,
                  },
                  `${(record.nsfw_score * 100).toFixed(1)}%`
                ),
              ])
            : h('div', { class: 'no-nsfw' }, '-'),
        ]),
    },
    {
      key: 'operations',
      title: $t('admin.contentReview.logs.columns.operations'),
      width: 120,
      align: 'center',
      render: (_value: unknown, record: ReviewLog) =>
        h('div', { class: 'operations-cell' }, [
          !selectMode.value && record.action === 'reject' && record.delete_type === 'soft' && record.image?.status === 'deleted'
            ? h('div', { class: 'operation-buttons' }, [
                h(
                  'button',
                  {
                    class:
                      'inline-flex items-center justify-center w-7 h-7 text-xs bg-green-500/20 hover:bg-green-500/30 border border-green-500/30 text-green-400 rounded transition-all duration-200 mr-2',
                    title: $t('admin.contentReview.logs.tooltips.restore'),
                    disabled: isRestoring[record.file_id],
                    onClick: (e: Event) => {
                      e.stopPropagation()
                      confirmRestore(record)
                    },
                  },
                  [isRestoring[record.file_id] ? h('i', { class: 'fas fa-spinner fa-spin' }) : h('i', { class: 'fas fa-undo' })]
                ),

                h(
                  'button',
                  {
                    class:
                      'inline-flex items-center justify-center w-7 h-7 text-xs bg-red-500/20 hover:bg-red-500/30 border border-red-500/30 text-red-400 rounded transition-all duration-200',
                    title: $t('admin.contentReview.logs.tooltips.hardDelete'),
                    disabled: isHardDeleting[record.file_id],
                    onClick: (e: Event) => {
                      e.stopPropagation()
                      confirmHardDelete(record)
                    },
                  },
                  [
                    isHardDeleting[record.file_id]
                      ? h('i', { class: 'fas fa-spinner fa-spin' })
                      : h('i', { class: 'fas fa-trash-alt' }),
                  ]
                ),
              ])
            : h('div', { class: 'no-operations' }, [
                h('span', { class: 'text-xs text-content-disabled' }, $t('admin.contentReview.logs.text.noOperations')),
              ]),
        ]),
    },
  ]

  /* 使用 computed 包装 getter 函数以支持响应式更新 */
  const tableColumns = computed(() => getTableColumns())
  const _skeletonColumns = computed(() => getSkeletonColumns())

  /* 筛选状态 */
  const showFilter = ref(false)
  const searchKeyword = ref('')
  const actionFilter = ref('')
  const auditorFilter = ref(0)
  const dateRange = ref<[Date | null, Date | null]>([null, null])

  /* 筛选选项 */
  const actionOptions = computed(() => [
    { label: $t('admin.contentReview.logs.filter.allActions'), value: '' },
    { label: $t('admin.contentReview.logs.actions.approve'), value: 'approve' },
    { label: $t('admin.contentReview.logs.actions.reject'), value: 'reject' },
  ])

  const auditorOptions = computed(() => {
    const auditors = new Map()
    logs.value.forEach((log) => {
      if (log.auditor) {
        auditors.set(log.auditor.id, log.auditor)
      }
    })

    const options = [{ label: $t('admin.contentReview.logs.filter.allAuditors'), value: 0 }]
    auditors.forEach((auditor) => {
      options.push({
        label: auditor.username,
        value: auditor.id,
      })
    })

    return options
  })

  const hasValidSelection = computed(() => selectedValidLogs.value.length > 0)

  const selectedValidLogs = computed(() =>
    logs.value.filter((log) => selectedRows.value.includes(log.id) && isLogSelectable(log))
  )

  const isHardDeleting = reactive<Record<string, boolean>>({})
  const showHardDeleteDialog = ref(false)
  const selectedLog = ref<ReviewLog | null>(null)
  const isExecutingHardDelete = ref(false)

  const isRestoring = reactive<Record<string, boolean>>({})
  const showRestoreDialog = ref(false)
  const isExecutingRestore = ref(false)

  const showBatchRestoreDialog = ref(false)
  const isBatchRestoring = ref(false)
  const showBatchHardDeleteDialog = ref(false)
  const isBatchHardDeleting = ref(false)

  const handleDateRangeChange = (range: [Date | null, Date | null]) => {
    dateRange.value = range
    applyFilters()
  }

  const loadData = async (page = 1) => {
    try {
      isLoading.value = true
      showSkeleton.value = true

      const params: ReviewLogQuery = {
        page,
        size: pagination.value.page_size,
        keyword: searchKeyword.value || undefined,
        action: (actionFilter.value as string | undefined) || undefined,
        auditor_id: auditorFilter.value || undefined,
        date_from: dateRange.value[0] ? formatDateForAPI(dateRange.value[0]) : undefined,
        date_to: dateRange.value[1] ? formatDateForAPI(dateRange.value[1]) : undefined,
      }

      const response = await getReviewLogs(params)

      if (response.success && response.data) {
        logs.value = response.data.data || []
        const p = response.data.pagination || {}
        pagination.value = {
          page: p.page || 1,
          page_size: p.page_size || pagination.value.page_size || 10,
          total: p.total || 0,
          total_page: p.total_page || 1,
        }
      } else {
        logs.value = []
        pagination.value = {
          page: 1,
          page_size: pagination.value.page_size,
          total: 0,
          total_page: 1,
        }
      }
    } catch {
      toast.error($t('admin.contentReview.logs.messages.loadError'))
      logs.value = []
      pagination.value = { page: 1, page_size: pagination.value.page_size, total: 0, total_page: 1 }
    } finally {
      isLoading.value = false
      showSkeleton.value = false
    }
  }

  const formatDateForAPI = (date: Date): string => {
    const year = date.getFullYear()
    const month = (date.getMonth() + 1).toString().padStart(2, '0')
    const day = date.getDate().toString().padStart(2, '0')
    return `${year}-${month}-${day}`
  }

  const refreshData = () => {
    loadData(pagination.value.page)
  }

  const applyFilters = () => {
    pagination.value.page = 1
    loadData(1)
  }

  const resetFilters = () => {
    searchKeyword.value = ''
    actionFilter.value = ''
    auditorFilter.value = 0
    dateRange.value = [null, null]
    pagination.value.page = 1
    loadData(1)
  }

  const handlePageChange = (page: number) => {
    pagination.value.page = page
    loadData(page)
  }

  const startSelectMode = () => {
    selectMode.value = true
    selectedRows.value = []
  }

  const exitSelectMode = () => {
    selectMode.value = false
    selectedRows.value = []
  }

  const isLogSelectable = (log: ReviewLog): boolean =>
    log.action === 'reject' && log.delete_type === 'soft' && log.image?.status === 'deleted'

  const handleRowClick = (_row: ReviewLog, _index: number) => {}

  const _toggleRowSelection = (rowId: number) => {
    const index = selectedRows.value.indexOf(rowId)
    if (index > -1) {
      selectedRows.value.splice(index, 1)
    } else {
      selectedRows.value.push(rowId)
    }
  }

  const confirmHardDelete = (log: ReviewLog) => {
    selectedLog.value = log
    showHardDeleteDialog.value = true
  }

  const confirmRestore = (log: ReviewLog) => {
    selectedLog.value = log
    showRestoreDialog.value = true
  }

  const executeRestore = async () => {
    if (!selectedLog.value?.file_id) {
      return
    }

    const fileId = selectedLog.value.file_id // 保存fileId

    try {
      isExecutingRestore.value = true
      isRestoring[fileId] = true

      const response = await restoreReviewedImage(fileId)
      if (response.success) {
        toast.success($t('admin.contentReview.logs.messages.restoreSuccess'))
        showRestoreDialog.value = false
        selectedLog.value = null
        refreshData()
      } else {
        toast.error(response.message || $t('admin.contentReview.logs.messages.restoreError'))
      }
    } catch {
      toast.error($t('admin.contentReview.logs.messages.restoreError'))
    } finally {
      isExecutingRestore.value = false
      isRestoring[fileId] = false // 使用保存的fileId
    }
  }

  const confirmBatchRestore = () => {
    showBatchRestoreDialog.value = true
  }

  const executeBatchRestore = async () => {
    if (selectedValidLogs.value.length === 0) {
      return
    }

    try {
      isBatchRestoring.value = true

      const fileIds = selectedValidLogs.value.map((log) => log.file_id)
      const response = await batchRestoreReviewedImages(fileIds)

      if (response.success) {
        const successCount = response.data?.success_count || 0
        const failCount = response.data?.fail_count || 0

        if (successCount > 0) {
          const failText =
            failCount > 0
              ? $t('admin.contentReview.logs.messages.batchRestoreSuccess', {
                  success: successCount,
                  failText: $t('admin.contentReview.logs.messages.failCount', { count: failCount }),
                })
              : $t('admin.contentReview.logs.messages.batchRestoreSuccess', { success: successCount, failText: '' })
          toast.success(failText)
        } else {
          toast.error($t('admin.contentReview.logs.messages.batchRestoreError'))
        }
      } else {
        toast.error(response.message || $t('admin.contentReview.logs.messages.batchRestoreError'))
      }

      showBatchRestoreDialog.value = false
      exitSelectMode()
      refreshData()
    } catch {
      toast.error($t('admin.contentReview.logs.messages.batchRestoreError'))
    } finally {
      isBatchRestoring.value = false
    }
  }

  const confirmBatchHardDelete = () => {
    showBatchHardDeleteDialog.value = true
  }

  const executeHardDelete = async () => {
    if (!selectedLog.value?.file_id) {
      return
    }

    const fileId = selectedLog.value.file_id // 保存fileId

    try {
      isExecutingHardDelete.value = true
      isHardDeleting[fileId] = true

      const response = await hardDeleteReviewedImage(fileId)
      if (response.success) {
        toast.success($t('admin.contentReview.logs.messages.hardDeleteSuccess'))
        showHardDeleteDialog.value = false
        selectedLog.value = null
        refreshData()
      } else {
        toast.error(response.message || $t('admin.contentReview.logs.messages.hardDeleteError'))
      }
    } catch {
      toast.error($t('admin.contentReview.logs.messages.hardDeleteError'))
    } finally {
      isExecutingHardDelete.value = false
      isHardDeleting[fileId] = false // 使用保存的fileId
    }
  }

  const executeBatchHardDelete = async () => {
    if (selectedValidLogs.value.length === 0) {
      return
    }

    try {
      isBatchHardDeleting.value = true

      const fileIds = selectedValidLogs.value.map((log) => log.file_id)
      const response = await batchHardDeleteReviewedImages(fileIds)

      if (response.success) {
        const successCount = response.data?.success_count || 0
        const failCount = response.data?.fail_count || 0

        if (successCount > 0) {
          const failText =
            failCount > 0
              ? $t('admin.contentReview.logs.messages.batchHardDeleteSuccess', {
                  success: successCount,
                  failText: $t('admin.contentReview.logs.messages.failCount', { count: failCount }),
                })
              : $t('admin.contentReview.logs.messages.batchHardDeleteSuccess', { success: successCount, failText: '' })
          toast.success(failText)
        } else {
          toast.error($t('admin.contentReview.logs.messages.batchHardDeleteError'))
        }
      } else {
        toast.error(response.message || $t('admin.contentReview.logs.messages.batchHardDeleteError'))
      }

      showBatchHardDeleteDialog.value = false
      exitSelectMode()
      refreshData()
    } catch {
      toast.error($t('admin.contentReview.logs.messages.batchHardDeleteError'))
    } finally {
      isBatchHardDeleting.value = false
    }
  }

  const formatDateTime = (dateStr: string) => {
    const date = new Date(dateStr)
    return date.toLocaleDateString(getCurrentLocale(), {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    })
  }

  const _getActionClass = (action: string) => (action === 'approve' ? 'action-approve' : 'action-reject')

  const getDeleteStatusVariant = (deleteType: string, imageStatus: string) => {
    if (deleteType === 'hard') {
      return 'bg-red-500/20 text-red-400 border border-red-500/30'
    } else if (deleteType === 'soft') {
      return imageStatus === 'deleted'
        ? 'bg-yellow-500/20 text-yellow-400 border border-yellow-500/30'
        : 'bg-blue-500/20 text-blue-400 border border-blue-500/30'
    }
    return 'bg-blue-500/20 text-blue-400 border border-blue-500/30'
  }

  const getDeleteStatusIcon = (deleteType: string, imageStatus: string) => {
    if (deleteType === 'hard') {
      return 'fa-trash'
    } else if (deleteType === 'soft') {
      return imageStatus === 'deleted' ? 'fa-recycle' : 'fa-check-circle'
    }
    return 'fa-question'
  }

  const getDeleteStatusText = (deleteType: string, imageStatus: string) => {
    if (deleteType === 'hard') {
      return $t('admin.contentReview.logs.deleteStatus.hardDelete')
    } else if (deleteType === 'soft') {
      return imageStatus === 'deleted'
        ? $t('admin.contentReview.logs.deleteStatus.softDeleted')
        : $t('admin.contentReview.logs.deleteStatus.restored')
    }
    return $t('admin.contentReview.logs.deleteStatus.unknown')
  }

  const getNsfwClass = (score: number) => {
    if (score >= 0.8) {
      return 'nsfw-high'
    }
    if (score >= 0.6) {
      return 'nsfw-medium'
    }
    if (score >= 0.3) {
      return 'nsfw-low'
    }
    return 'nsfw-safe'
  }

  const truncateText = (text: string, maxLength: number) => {
    if (text.length <= maxLength) {
      return text
    }
    return `${text.substring(0, maxLength)}...`
  }

  const handleImageError = (event: Event) => {
    const img = event.target as HTMLImageElement
    img.src = '/placeholder-image.png'
  }

  watch(
    () => props.isActive,
    (newVal) => {
      if (newVal) {
        loadData(1)
      }
    },
    { immediate: true }
  )

  onMounted(() => {
    if (props.isActive) {
      loadData(1)
    }
  })
</script>

<template>
  <div class="review-logs-container">
    <div class="logs-header">
      <div class="header-info"></div>

      <div class="header-actions">
        <CyberButton type="secondary" :loading="isLoading" @click="refreshData">
          <i class="fas fa-sync-alt mr-1.5" />{{ $t('admin.contentReview.logs.buttons.refresh') }}
        </CyberButton>

        <CyberButton v-if="!selectMode" type="secondary" :disabled="logs.length === 0" @click="startSelectMode">
          <i class="fas fa-tasks mr-1.5" />{{ $t('admin.contentReview.logs.buttons.batchOperation') }}
        </CyberButton>

        <div v-else class="batch-controls">
          <div class="selection-info">
            <span class="selection-count"
              >{{ $t('admin.contentReview.logs.batch.selected') }}
              {{ $t('admin.contentReview.logs.batch.selectedCount', { current: selectedRows.length, total: logs.length }) }}</span
            >
          </div>

          <CyberButton
            type="success"
            size="small"
            :disabled="selectedRows.length === 0 || !hasValidSelection"
            :loading="isBatchRestoring"
            @click="confirmBatchRestore"
          >
            <i class="fas fa-undo mr-1" />{{ $t('admin.contentReview.logs.buttons.batchRestore') }} ({{ selectedRows.length }})
          </CyberButton>
          <CyberButton
            type="danger"
            size="small"
            :disabled="selectedRows.length === 0 || !hasValidSelection"
            :loading="isBatchHardDeleting"
            @click="confirmBatchHardDelete"
          >
            <i class="fas fa-trash mr-1" />{{ $t('admin.contentReview.logs.buttons.batchHardDelete') }} ({{
              selectedRows.length
            }})
          </CyberButton>
          <CyberButton type="secondary" size="small" icon="times" @click="exitSelectMode">{{
            $t('admin.contentReview.logs.buttons.cancel')
          }}</CyberButton>
        </div>

        <CyberButton
          type="secondary"
          :class="{ 'border-error-500 text-error-500': showFilter }"
          icon="filter"
          @click="showFilter = !showFilter"
        >
          {{ $t('admin.contentReview.logs.buttons.filter') }}
        </CyberButton>
      </div>
    </div>

    <div v-show="showFilter" class="filter-panel mb-4 rounded-xl border border-subtle bg-background-600 p-4">
      <div class="mb-3 flex items-center justify-between">
        <h4 class="text-base font-bold text-content-heading">
          <i class="fas fa-filter mr-1.5 text-content" />{{ $t('admin.contentReview.logs.filter.title') }}
        </h4>
        <div class="flex space-x-2">
          <CyberButton type="secondary" size="small" icon="undo-alt" @click="resetFilters">{{
            $t('admin.contentReview.logs.filter.reset')
          }}</CyberButton>
          <CyberButton type="primary" size="small" icon="search" @click="applyFilters">{{
            $t('admin.contentReview.logs.filter.apply')
          }}</CyberButton>
        </div>
      </div>

      <div class="grid grid-cols-1 gap-4 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5">
        <div class="filter-item">
          <label class="mb-1 block text-sm text-content">{{ $t('admin.contentReview.logs.filter.keyword') }}</label>
          <CyberInput
            v-model="searchKeyword"
            :placeholder="$t('admin.contentReview.logs.filter.keywordPlaceholder')"
            @keyup.enter="applyFilters"
          />
        </div>

        <div class="filter-item">
          <label class="mb-1 block text-sm text-content">{{ $t('admin.contentReview.logs.filter.auditAction') }}</label>
          <CyberDropdown
            v-model="actionFilter"
            :options="actionOptions"
            :placeholder="$t('admin.contentReview.logs.filter.auditActionPlaceholder')"
            @change="applyFilters"
          />
        </div>

        <div class="filter-item">
          <label class="mb-1 block text-sm text-content">{{ $t('admin.contentReview.logs.filter.auditor') }}</label>
          <CyberDropdown
            v-model="auditorFilter"
            :options="auditorOptions"
            :placeholder="$t('admin.contentReview.logs.filter.auditorPlaceholder')"
            @change="applyFilters"
          />
        </div>

        <div class="filter-item">
          <label class="mb-1 block text-sm text-content">{{ $t('admin.contentReview.logs.filter.dateRange') }}</label>
          <CyberDatePicker
            v-model="dateRange"
            type="daterange"
            :placeholder="$t('admin.contentReview.logs.filter.dateRangePlaceholder')"
            :start-placeholder="$t('admin.contentReview.logs.filter.startDatePlaceholder')"
            :end-placeholder="$t('admin.contentReview.logs.filter.endDatePlaceholder')"
            @change="handleDateRangeChange"
          />
        </div>
      </div>
    </div>

    <div class="content-area">
      <CyberSkeleton type="table" :count="pagination.page_size" :loading="showSkeleton" />

      <CyberDataTable
        v-if="!showSkeleton"
        :columns="tableColumns"
        :data="logs"
        :loading="false"
        :selectable="selectMode"
        :selected-row-keys="selectedRows"
        :hoverable="true"
        :striped="true"
        :bordered="true"
        max-height="calc(100vh - 350px)"
        :empty-text="$t('admin.contentReview.logs.empty.title')"
        :loading-text="$t('admin.contentReview.logs.loading.text')"
        @row-click="handleRowClick"
      />

      <div v-if="pagination.total_page > 1" class="pagination-wrapper">
        <CyberPagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.page_size"
          :total="pagination.total"
          @update:current-page="handlePageChange"
        />
      </div>
    </div>

    <CyberDialog
      v-model="showHardDeleteDialog"
      :title="$t('admin.contentReview.logs.dialogs.hardDelete.title')"
      :loading="isExecutingHardDelete"
      @confirm="executeHardDelete"
      @cancel="showHardDeleteDialog = false"
    >
      <div class="space-y-4">
        <div class="warning-box">
          <i class="fas fa-exclamation-triangle mr-2 text-red-400" />
          <span class="font-medium text-red-400">{{ $t('admin.contentReview.logs.dialogs.hardDelete.warning') }}</span>
        </div>
        <p class="text-content">
          {{ $t('admin.contentReview.logs.dialogs.hardDelete.confirmMessage', { name: selectedLog?.image?.original_name }) }}
        </p>
        <p class="text-content-content-muted text-sm">{{ $t('admin.contentReview.logs.dialogs.hardDelete.description') }}</p>
      </div>
    </CyberDialog>

    <CyberDialog
      v-model="showRestoreDialog"
      :title="$t('admin.contentReview.logs.dialogs.restore.title')"
      :loading="isExecutingRestore"
      @confirm="executeRestore"
      @cancel="showRestoreDialog = false"
    >
      <div class="space-y-4">
        <p class="text-content">
          {{ $t('admin.contentReview.logs.dialogs.restore.confirmMessage', { name: selectedLog?.image?.original_name }) }}
        </p>
        <p class="text-content-content-muted text-sm">{{ $t('admin.contentReview.logs.dialogs.restore.description') }}</p>
      </div>
    </CyberDialog>

    <CyberDialog
      v-model="showBatchRestoreDialog"
      :title="$t('admin.contentReview.logs.dialogs.batchRestore.title')"
      :loading="isBatchRestoring"
      @confirm="executeBatchRestore"
      @cancel="showBatchRestoreDialog = false"
    >
      <div class="space-y-4">
        <p class="text-content">
          {{ $t('admin.contentReview.logs.dialogs.batchRestore.confirmMessage', { count: selectedValidLogs.length }) }}
        </p>
        <p class="text-content-content-muted text-sm">{{ $t('admin.contentReview.logs.dialogs.batchRestore.description') }}</p>
        <div class="mt-3 max-h-32 overflow-y-auto rounded bg-background-400 p-2">
          <div v-for="log in selectedValidLogs" :key="log.id" class="mb-1 text-xs text-content-muted">
            • {{ log.image?.original_name || log.file_id }}
          </div>
        </div>
      </div>
    </CyberDialog>

    <CyberDialog
      v-model="showBatchHardDeleteDialog"
      :title="$t('admin.contentReview.logs.dialogs.batchHardDelete.title')"
      :loading="isBatchHardDeleting"
      @confirm="executeBatchHardDelete"
      @cancel="showBatchHardDeleteDialog = false"
    >
      <div class="space-y-4">
        <div class="warning-box">
          <i class="fas fa-exclamation-triangle mr-2 text-red-400" />
          <span class="font-medium text-red-400">{{ $t('admin.contentReview.logs.dialogs.batchHardDelete.warning') }}</span>
        </div>
        <p class="text-content">
          {{ $t('admin.contentReview.logs.dialogs.batchHardDelete.confirmMessage', { count: selectedValidLogs.length }) }}
        </p>
        <p class="text-content-content-muted text-sm">{{ $t('admin.contentReview.logs.dialogs.batchHardDelete.description') }}</p>
        <div class="mt-3 max-h-32 overflow-y-auto rounded bg-background-400 p-2">
          <div v-for="log in selectedValidLogs" :key="log.id" class="mb-1 text-xs text-content-muted">
            • {{ log.image?.original_name || log.file_id }}
          </div>
        </div>
      </div>
    </CyberDialog>
  </div>
</template>

<style scoped lang="scss">
  .review-logs-container {
    @apply flex flex-col;
    flex: 1;
    min-height: 0;
    overflow: hidden;
  }

  .logs-header {
    @apply mb-6 flex flex-shrink-0 items-start justify-between;
  }

  .header-info {
    @apply space-y-1;
  }

  .header-actions {
    @apply flex items-center gap-3;
  }

  .filter-item {
    @apply space-y-1;
  }

  .content-area {
    @apply flex min-h-0 flex-1 flex-col;
    overflow: hidden;
  }

  .pagination-wrapper {
    @apply flex-shrink-0 border-t border-subtle p-4;
  }

  .warning-box {
    @apply flex items-center rounded border border-red-500/30 bg-red-500/10 p-3;
  }

  .batch-controls {
    @apply flex items-center gap-2;
  }

  .selection-info {
    @apply mr-2 text-sm text-content-muted;
  }

  .selection-count {
    @apply rounded bg-brand-200 px-2 py-1 text-xs text-content;
  }

  .image-cell {
    @apply flex justify-center;
  }

  .image-wrapper {
    @apply relative h-8 w-8 overflow-hidden rounded-md;
    background: rgba(var(--color-background-900-rgb), 0.2);
  }

  .table-image-thumb {
    @apply h-full max-h-[80px] w-full object-cover;
  }

  .image-placeholder {
    @apply flex h-full w-full items-center justify-center;
    @apply bg-background-400 text-brand-500;

    i {
      @apply text-lg;
    }
  }

  .format-badge {
    @apply absolute right-0.5 top-0.5 rounded bg-brand-800 px-1 py-0.5 text-content-heading;
    font-size: 10px;
    line-height: 1;
  }

  .image-info-cell {
    @apply space-y-1;
  }

  .image-name {
    @apply truncate text-sm font-medium text-content-heading;
  }

  .image-meta {
    @apply flex flex-wrap gap-2;
  }

  .meta-item {
    @apply flex items-center text-xs text-content-muted;
    @apply rounded bg-brand-500 px-2 py-0.5;

    &.no-info {
      @apply bg-yellow-500/10 text-yellow-400/80;
    }
  }

  .action-cell {
    @apply flex justify-center;
  }

  .soft-delete-cell {
    @apply flex justify-center;
  }

  .no-delete-status {
    @apply flex items-center justify-center;
  }

  .status-placeholder {
    @apply text-sm font-medium text-content-disabled;
  }

  .auditor-cell {
    @apply space-y-1;
  }

  .auditor-name {
    @apply text-sm font-medium text-content-heading;
  }

  .audit-time-cell {
    @apply space-y-1;
  }

  .audit-time {
    @apply text-xs text-content-muted;
  }

  .reason-cell {
    @apply space-y-1;
  }

  .reason-text {
    @apply text-sm text-content;
  }

  .no-reason {
    @apply text-sm italic text-content-disabled;
  }

  .nsfw-cell {
    @apply space-y-1;
  }

  .nsfw-score {
    @apply space-y-0.5;
  }

  .score-value {
    @apply text-sm font-medium;

    &.nsfw-safe {
      @apply text-green-400;
    }

    &.nsfw-low {
      @apply text-yellow-400;
    }

    &.nsfw-medium {
      @apply text-orange-400;
    }

    &.nsfw-high {
      @apply text-red-400;
    }
  }

  .no-nsfw {
    @apply text-sm text-content-disabled;
  }

  .operations-cell {
    @apply flex items-center justify-center;
  }

  .operation-buttons {
    @apply flex items-center gap-2;
  }

  .no-operations {
    @apply text-center;
  }

  .table-empty-state {
    @apply flex flex-col items-center justify-center py-20;
  }

  .cyber-checkbox {
    @apply h-4 w-4 rounded border-2 border-neutral-600;
    @apply bg-transparent checked:border-brand-500 checked:bg-brand-500;
    @apply cursor-pointer transition-all duration-200;

    &:checked {
      background-image: url("data:image/svg+xml,%3csvg viewBox='0 0 16 16' fill='white' xmlns='http://www.w3.org/2000/svg'%3e%3cpath d='m13.854 3.646-7.5 7.5a.5.5 0 0 1-.708 0l-3.5-3.5a.5.5 0 1 1 .708-.708L6 10.293l7.146-7.147a.5.5 0 0 1 .708.708z'/%3e%3c/svg%3e");
    }

    &:disabled {
      @apply cursor-not-allowed opacity-50;
    }
  }
</style>
