<script setup lang="ts">
  import { computed, nextTick, onMounted, onUnmounted, ref } from 'vue'
  import { useRoute } from 'vue-router'
  import UploadProgress from './components/UploadProgress.vue'
  import UploadSettings from './components/UploadSettings.vue'
  import FolderConfirmDialog from './components/FolderConfirmDialog.vue'
  import { useToast } from '@/components/Toast/useToast'
  import { useSettingsStore } from '@/store/settings'
  import { useAuthStore } from '@/store/auth'
  import { useGlobalUpload } from '@/composables/useGlobalUpload'
  import { useGlobalSettings } from '@/composables/useGlobalSettings'
  import { useFolderPath } from '@/hooks/useFolderPath'
  import { useUploadStore } from '@/store/upload'
  import { useLayoutStore } from '@/store/layout'
  import { useUploadConfig } from '@/composables/useUploadConfig'
  import { DEFAULT_WATERMARK_CONFIG, type WatermarkConfig as WatermarkConfigType } from '@/components/WatermarkConfig/types'
  import { logger } from '@/utils/system/logger'
  import { useTexts } from '@/composables/useTexts'
  import {
    extractFilesFromFileList,
    extractFilesFromItems,
    filterValidFiles,
    getFolderStats,
    hasFolder,
    supportsFolderUpload,
    type FolderStats,
  } from '@/utils/file/folderReader'

  const toast = useToast()
  const settingsStore = useSettingsStore()
  const authStore = useAuthStore()
  const route = useRoute()
  const { $t } = useTexts()

  const { settings: globalSettings } = useGlobalSettings()

  const { getFolderPathString } = useFolderPath()

  const uploadConfig = useUploadConfig()

  const isDismissed = ref(false)

  const fileInput = ref<HTMLInputElement | null>(null)
  const folderInput = ref<HTMLInputElement | null>(null)
  const isDragging = ref(false)

  const showFolderDialog = ref(false)
  const folderStats = ref<FolderStats>({ totalFiles: 0, validFiles: 0, invalidFiles: 0 })
  const pendingFolderFiles = ref<File[]>([])

  const isFolderUploadSupported = supportsFolderUpload()

  // 防止重复触发
  let isSelectingFolder = false

  const folderId = ref<string | null>(null)
  const presetFolderPath = ref<string>('') // 预设文件夹路径，避免重复查询
  const accessLevel = ref<'public' | 'private' | 'protected'>('public')
  const optimize = ref<boolean>(true)
  const autoRemove = ref<boolean>(false)
  const storageDuration = ref<string>(authStore.isLoggedIn ? 'permanent' : '7d') // 存储时长
  const showAdvancedSettings = ref<boolean>(false) // 高级设置折叠状态
  const advancedSettingsTrigger = ref<HTMLElement | null>(null) // 高级设置触发按钮引用
  const panelPosition = ref({}) // 面板位置

  const watermarkEnabled = ref<boolean>(false)
  const watermarkConfig = ref<WatermarkConfigType>({ ...DEFAULT_WATERMARK_CONFIG })
  const showWatermarkConfig = ref<boolean>(false)

  const uploadStore = useUploadStore()
  const layoutStore = useLayoutStore()

  const globalUpload = useGlobalUpload({
    folderId: computed(() => folderId.value),
    accessLevel: computed(() => accessLevel.value),
    optimize: computed(() => optimize.value),
  })

  const {
    uploadQueue,
    globalProgress,
    addFiles,
    startUpload,
    resumeUpload,
    cancelUpload,
    retryUpload,
    removeFile,
    clearQueue,
    clearAllSessions,
    formatFileSize,
    pendingCount,
    uploadingCount,
    successCount,
    errorCount,
    hasPendingFiles,
    hasUploadingFiles,
    hasSuccessFiles,
    totalFileSize,
    maxConcurrentUploads,
    copyAllUrls: copyAllUrlsFromStore,
  } = globalUpload

  const setupPreselectedFolder = async (folderIdParam: string) => {
    try {
      const pathString = await getFolderPathString(folderIdParam)
      presetFolderPath.value = pathString

      uploadStore.setGlobalOptions({ folderId: folderIdParam })

      if (pathString) {
        toast.success($t('upload.toast.folderSelected', { path: pathString }))
      } else {
        toast.info($t('upload.toast.folderPreselected'))
      }
    } catch (error) {
      logger.error('[Upload Page] Error setting up preselected folder:', error)
      toast.info($t('upload.toast.folderPreselected'))
    }
  }

  const shouldShowContentDetectionWarning = computed(() => {
    if (isDismissed.value) {
      return false
    }
    if (!globalSettings.value?.upload) {
      return false
    }

    const { upload } = globalSettings.value
    const contentDetectionEnabled = upload.content_detection_enabled
    const sensitiveHandling = upload.sensitive_content_handling

    return contentDetectionEnabled && sensitiveHandling && sensitiveHandling !== 'mark_only'
  })

  const simpleWarningIcon = computed(() => {
    const upload = globalSettings.value?.upload
    const sensitiveHandling = upload?.sensitive_content_handling

    if (sensitiveHandling === 'auto_delete') {
      return 'fas fa-exclamation-triangle text-red-400'
    } else if (sensitiveHandling === 'pending_review') {
      return 'fas fa-clock text-amber-400'
    }
    return 'fas fa-info-circle text-brand-500'
  })

  const simpleWarningMessage = computed(() => {
    const upload = globalSettings.value?.upload
    const sensitiveHandling = upload?.sensitive_content_handling

    if (sensitiveHandling === 'auto_delete') {
      return $t('upload.toast.sensitiveAutoDelete')
    } else if (sensitiveHandling === 'pending_review') {
      return $t('upload.toast.sensitivePendingReview')
    }
    return $t('upload.toast.defaultHint')
  })

  const acceptFormats = computed(() => {
    const formats = settingsStore.allowedImageFormats
    if (formats.length === 0) {
      return 'image/*'
    }
    return formats.map((format) => `image/${format.toLowerCase()}`).join(',')
  })

  const _pendingFilesSize = computed(() =>
    uploadQueue.value.filter((file) => file.status === 'pending').reduce((total, file) => total + file.size, 0)
  )

  const storageDurationOptionsForLoggedUser = computed(() => {
    if (!authStore.isLoggedIn) {
      return []
    }

    const allowedDurations = globalSettings.value?.upload?.user_allowed_storage_durations
    const options: { label: string; value: string }[] = []

    if (Array.isArray(allowedDurations) && allowedDurations.length > 0) {
      allowedDurations.forEach((duration) => {
        if (duration === 'permanent') {
          options.push({
            label: $t('upload.duration.permanent'),
            value: duration,
          })
        } else {
          options.push({
            label: duration,
            value: duration,
          })
        }
      })
    } else {
      options.push({ label: $t('upload.duration.permanent'), value: 'permanent' })
    }

    return options
  })

  const pasteShortcut = computed(() => {
    const isMac = /Mac|iPod|iPhone|iPad/.test(navigator.userAgent)
    return {
      key: isMac ? 'Cmd' : 'Ctrl',
      isMac,
    }
  })

  const _dismissWarning = () => {
    isDismissed.value = true
  }

  const triggerFileInput = () => {
    if (isSelectingFolder) {
      return
    }
    fileInput.value?.click()
  }

  const triggerFolderInput = () => {
    if (isSelectingFolder) {
      return
    }

    if (!isFolderUploadSupported) {
      toast.warning($t('upload.folder.browserNotSupported'))
      return
    }

    if (!folderInput.value) {
      toast.error($t('upload.folder.initFailed'))
      return
    }

    isSelectingFolder = true

    const input = folderInput.value

    if (input.hasAttribute('webkitdirectory')) {
      input.click()
    } else {
      input.setAttribute('webkitdirectory', '')
      input.setAttribute('mozdirectory', '')
      input.setAttribute('directory', '')
      input.click()
    }

    setTimeout(() => {
      isSelectingFolder = false
    }, 3000)
  }

  const handleFolderChange = async (event: Event) => {
    const target = event.target as HTMLInputElement
    isSelectingFolder = false

    if (target.files && target.files.length > 0) {
      await processFolderFiles(target.files)
      target.value = ''
    }
  }

  const processFolderFiles = async (fileList: FileList) => {
    try {
      if (!uploadConfig.isConfigLoaded.value) {
        await uploadConfig.loadConfig()
      }

      const allFiles = extractFilesFromFileList(fileList)

      if (allFiles.length === 0) {
        toast.warning($t('upload.folder.noFiles'))
        return
      }

      const validFiles = filterValidFiles(allFiles, uploadConfig.isAllowedFileType, uploadConfig.isAllowedFileSize)

      if (validFiles.length === 0) {
        toast.error($t('upload.folder.noValidFiles'))
        return
      }

      const stats = getFolderStats(allFiles, validFiles)
      folderStats.value = stats
      pendingFolderFiles.value = validFiles

      showFolderDialog.value = true
    } catch (error) {
      logger.error('[Upload] processFolderFiles error:', error)
      toast.error($t('upload.folder.processingError'))
    }
  }

  const handleConfirmFolderUpload = async () => {
    if (pendingFolderFiles.value.length > 0) {
      await addFilesWithWatermark(pendingFolderFiles.value as any)
      toast.success($t('upload.folder.addSuccess', { count: pendingFolderFiles.value.length }))
      pendingFolderFiles.value = []
    }
  }

  const handleCancelFolderUpload = () => {
    pendingFolderFiles.value = []
  }

  const onDragOver = () => {
    isDragging.value = true
  }

  const onDragLeave = () => {
    isDragging.value = false
  }

  const onDrop = async (event: DragEvent) => {
    isDragging.value = false
    const items = event.dataTransfer?.items

    if (!items) {
      const files = event.dataTransfer?.files
      if (files) {
        await addFilesWithWatermark(files)
      }
      return
    }

    const containsFolder = hasFolder(items)

    if (containsFolder) {
      await handleFolderDrop(items)
    } else {
      const files = event.dataTransfer?.files
      if (files) {
        await addFilesWithWatermark(files)
      }
    }
  }

  const handleFolderDrop = async (items: DataTransferItemList) => {
    try {
      toast.info($t('upload.folder.parsing'))

      if (!uploadConfig.isConfigLoaded.value) {
        await uploadConfig.loadConfig()
      }

      const allFiles = await extractFilesFromItems(items)

      if (allFiles.length === 0) {
        toast.warning($t('upload.folder.noFiles'))
        return
      }

      const validFiles = filterValidFiles(allFiles, uploadConfig.isAllowedFileType, uploadConfig.isAllowedFileSize)

      if (validFiles.length === 0) {
        toast.error($t('upload.folder.noValidFiles'))
        return
      }

      const stats = getFolderStats(allFiles, validFiles)
      folderStats.value = stats
      pendingFolderFiles.value = validFiles

      showFolderDialog.value = true
    } catch (error) {
      logger.error('[Upload] handleFolderDrop error:', error)
      toast.error($t('upload.folder.processingError'))
    }
  }

  const handleFileChange = async (event: Event) => {
    const target = event.target as HTMLInputElement

    if (target.files) {
      await addFilesWithWatermark(target.files)
      target.value = ''
    }
  }

  const handlePaste = async (e: ClipboardEvent) => {
    const items = e.clipboardData?.items
    if (!items) {
      return
    }

    const files: File[] = []
    for (let i = 0; i < items.length; i++) {
      const item = items[i]
      if (item.type.indexOf('image') !== -1) {
        const file = item.getAsFile()
        if (file) {
          files.push(file)
        }
      }
    }

    if (files.length > 0) {
      const dt = new DataTransfer()
      files.forEach((file) => dt.items.add(file))

      await addFilesWithWatermark(dt.files)
    }
  }

  const addFilesWithWatermark = async (files: FileList) => {
    const syncedWatermarkConfig = {
      ...watermarkConfig.value,
      enabled: watermarkEnabled.value,
    }

    watermarkConfig.value = syncedWatermarkConfig

    await addFiles(files)
  }

  const startUploadWrapper = async () => {
    const syncedWatermarkConfig = {
      ...watermarkConfig.value,
      enabled: watermarkEnabled.value,
    }

    if (!authStore.isLoggedIn) {
      if (!storageDuration.value || storageDuration.value === 'permanent') {
        logger.warn('[Upload] Guest mode storage duration validation failed')
        toast.error($t('upload.toast.guestDurationRequired'))
        return
      }
    }

    if (!hasPendingFiles.value) {
      showNoFilesMessage()
      return
    }

    watermarkConfig.value = syncedWatermarkConfig

    const uploadWatermarkConfig = {
      ...syncedWatermarkConfig,
      type: 'image' as const, // 上传时强制转为image类型
    }

    uploadStore.setGlobalOptions({
      watermarkEnabled: watermarkEnabled.value,
      watermarkConfig: uploadWatermarkConfig,
    })

    await startUpload()
  }

  const handleCancelUpload = () => {
    cancelUpload()
    toast.info($t('upload.toast.uploadCancelled'))
  }

  const retryUploadWrapper = (index: number) => {
    retryUpload(index)
  }

  const resumeUploadWrapper = (index: number) => {
    resumeUpload(index.toString())
  }

  const removeFileWrapper = (index: number) => {
    removeFile(index)
  }

  const showNoFilesMessage = () => {
    toast.warning($t('upload.toast.noFilesInQueue'))
  }

  const copyAllUrls = async () => {
    const success = await copyAllUrlsFromStore()
    if (success) {
      toast.success($t('upload.toast.allLinksCopied'))
    } else {
      toast.warning($t('upload.toast.noLinksToCopy'))
    }
  }

  const calculatePanelPosition = () => {
    if (!advancedSettingsTrigger.value) {
      return {}
    }

    const triggerRect = advancedSettingsTrigger.value.getBoundingClientRect()
    const panelWidth = 320 // 320px
    const panelHeight = authStore.isLoggedIn ? 200 : 280 // 根据内容调整高度
    const viewportWidth = window.innerWidth
    const viewportHeight = window.innerHeight

    let left = triggerRect.right - panelWidth
    let top = triggerRect.bottom + 8

    if (left < 20) {
      left = 20
    }

    if (left + panelWidth > viewportWidth - 20) {
      left = viewportWidth - panelWidth - 20
    }

    if (top + panelHeight > viewportHeight - 20) {
      top = triggerRect.top - panelHeight - 8
    }

    if (top < 20) {
      top = 20
    }

    return {
      left: `${left}px`,
      top: `${top}px`,
    }
  }

  const toggleAdvancedSettings = () => {
    if (!showAdvancedSettings.value) {
      panelPosition.value = calculatePanelPosition()
    }
    showAdvancedSettings.value = !showAdvancedSettings.value
  }
  const onSettingsChange = (settings: any) => {
    const uploadWatermarkConfig = settings.watermarkConfig
      ? {
          ...settings.watermarkConfig,
          type: 'image' as const,
        }
      : settings.watermarkConfig

    uploadStore.setGlobalOptions({
      folderId: settings.folderId,
      accessLevel: settings.accessLevel,
      optimize: settings.optimize,
      autoRemove: settings.autoRemove,
      storageDuration: settings.storageDuration,
      watermarkEnabled: settings.watermarkEnabled,
      watermarkConfig: uploadWatermarkConfig,
    })

    if (settings.storageDuration) {
      storageDuration.value = settings.storageDuration
    }
    if (settings.watermarkEnabled !== undefined) {
      watermarkEnabled.value = settings.watermarkEnabled
    }
    if (settings.watermarkConfig) {
      watermarkConfig.value = settings.watermarkConfig
    }
  }

  const getPositionText = (position: string) => {
    const positionMap: Record<string, string> = {
      'top-left': $t('upload.watermark.position.top-left'),
      'top-center': $t('upload.watermark.position.top-center'),
      'top-right': $t('upload.watermark.position.top-right'),
      'middle-left': $t('upload.watermark.position.middle-left'),
      'middle-center': $t('upload.watermark.position.middle-center'),
      'middle-right': $t('upload.watermark.position.middle-right'),
      'bottom-left': $t('upload.watermark.position.bottom-left'),
      'bottom-center': $t('upload.watermark.position.bottom-center'),
      'bottom-right': $t('upload.watermark.position.bottom-right'),
      custom: $t('upload.watermark.position.custom'),
    }
    return positionMap[position] || position
  }

  const handleWatermarkConfigConfirm = (newConfig: WatermarkConfigType) => {
    watermarkConfig.value = { ...newConfig }
    showWatermarkConfig.value = false
  }

  watch(storageDuration, (newValue) => {
    uploadStore.setGlobalOptions({
      storageDuration: newValue,
    })
  })

  watch(showWatermarkConfig, (newValue) => {
    if (newValue && showAdvancedSettings.value) {
      showAdvancedSettings.value = false
    }
  })

  watch(
    storageDurationOptionsForLoggedUser,
    (newOptions) => {
      if (newOptions && newOptions.length > 0 && globalSettings.value && authStore.isLoggedIn) {
        const defaultDuration = globalSettings.value.upload?.user_default_storage_duration || 'permanent'

        const currentValueValid = storageDuration.value && newOptions.some((opt) => opt.value === storageDuration.value)
        const shouldSetDefault = !currentValueValid || !storageDuration.value

        if (shouldSetDefault && defaultDuration && newOptions.some((opt) => opt.value === defaultDuration)) {
          storageDuration.value = defaultDuration
        } else if (shouldSetDefault && newOptions.length > 0) {
          storageDuration.value = newOptions[0].value
        }
      }
    },
    { immediate: true }
  )

  onMounted(async () => {
    await nextTick()

    if (folderInput.value && isFolderUploadSupported) {
      const input = folderInput.value as HTMLInputElement
      // 通过 JavaScript 设置属性，这是最可靠的方式
      input.setAttribute('webkitdirectory', '')
      input.setAttribute('mozdirectory', '')
      input.setAttribute('directory', '')

      // 同时设置 DOM 属性（某些浏览器需要）
      try {
        Object.defineProperty(input, 'webkitdirectory', {
          value: true,
          writable: true,
          configurable: true,
        })
      } catch (e) {
        // 忽略错误
      }
    }

    if (!authStore.isLoggedIn) {
      if (globalSettings.value?.guest?.guest_default_access_level) {
        accessLevel.value = globalSettings.value.guest.guest_default_access_level as 'public' | 'private' | 'protected'
      }

      if (globalSettings.value?.guest?.guest_default_storage_duration) {
        storageDuration.value = globalSettings.value.guest.guest_default_storage_duration as string
      } else {
        const allowedDurations = globalSettings.value?.guest?.guest_allowed_storage_durations
        if (Array.isArray(allowedDurations) && allowedDurations.length > 0) {
          storageDuration.value = allowedDurations[0]
        }
      }
    } else {
      if (globalSettings.value?.upload?.user_default_storage_duration) {
        storageDuration.value = globalSettings.value.upload.user_default_storage_duration as string
      }
    }

    document.addEventListener('paste', handlePaste)

    const queryFolderId = route.query.folderId as string
    if (queryFolderId) {
      folderId.value = queryFolderId
      setupPreselectedFolder(queryFolderId)
    }
  })

  onUnmounted(() => {
    document.removeEventListener('paste', handlePaste)
  })
</script>

<template>
  <div class="upload-page min-h-screen">
    <div
      class="upload-container relative z-10 flex w-full flex-col pt-4"
      :class="{
        'mx-auto': layoutStore.isTopLayout,
      }"
    >
      <div class="upload-top-panel mb-4 rounded-lg">
        <div class="grid grid-cols-1 gap-4 p-4 md:grid-cols-2">
          <div class="flex flex-col">
            <div class="mb-3">
              <h1 class="mb-0.5 text-xl font-bold text-content-heading">{{ $t('upload.title') }}</h1>
              <p class="mt-1 text-sm leading-relaxed text-content-muted">{{ settingsStore.uploadLimitText }}</p>
            </div>

            <div
              class="upload-drop-zone group relative min-h-48 flex-grow cursor-pointer rounded-lg border-2 border-dashed border-default transition-all duration-300 hover:border-brand-500 hover:bg-background-700 hover:shadow-lg"
              :class="{
                'border-brand-500 bg-background-700 shadow-lg': isDragging,
              }"
              @click="triggerFileInput"
              @dragover.prevent="onDragOver"
              @dragleave.prevent="onDragLeave"
              @drop.prevent="onDrop"
            >
              <input ref="fileInput" type="file" class="hidden" :accept="acceptFormats" multiple @change="handleFileChange" />
              <input ref="folderInput" type="file" class="hidden" multiple @change="handleFolderChange" />

              <div class="absolute inset-0 flex flex-col items-center justify-center p-4">
                <div class="upload-icon-container relative mb-5">
                  <div class="upload-icon relative p-6">
                    <i
                      class="fas fa-cloud-upload-alt to-accent-500 bg-gradient-to-r from-brand-500 bg-clip-text text-5xl text-transparent"
                    />
                  </div>
                </div>

                <div class="mb-3 text-center">
                  <p class="mb-1 text-lg font-semibold text-content">{{ $t('upload.dropZone.title') }}</p>
                  <p class="text-content-secondary text-base">
                    {{ $t('upload.dropZone.orText')
                    }}<span class="font-medium text-brand-500 transition-colors hover:text-content-heading">
                      {{ $t('upload.dropZone.clickToSelect') }}</span
                    >
                  </p>
                  <p v-if="isFolderUploadSupported" class="text-content-secondary mt-1 text-sm">
                    <i class="fas fa-folder-open text-brand-500" />
                    {{ $t('upload.dropZone.supportFolder') }}
                  </p>
                </div>

                <div class="text-content-secondary mb-3 flex items-center gap-4 text-xs">
                  <div class="flex items-center gap-1">
                    <i class="fas fa-images text-brand-500" />
                    <span>{{ $t('upload.dropZone.features.batch') }}</span>
                  </div>
                  <div class="flex items-center gap-1">
                    <i class="fas fa-compress text-success-500" />
                    <span>{{ $t('upload.dropZone.features.autoOptimize') }}</span>
                  </div>
                  <div class="flex items-center gap-1">
                    <i class="fas fa-shield-alt text-content-heading" />
                    <span>{{ $t('upload.dropZone.features.secureStorage') }}</span>
                  </div>
                </div>
                <div class="text-content-secondary flex items-center justify-center gap-2 text-xs">
                  <i class="fas fa-keyboard text-brand-500" />
                  <span>{{ $t('upload.dropZone.shortcut.label') }}</span>
                  <div class="flex items-center gap-1">
                    <kbd class="rounded border border-brand-500 bg-background-600 px-1.5 py-0.5 font-medium text-content">
                      {{ pasteShortcut.key }}
                    </kbd>
                    <span>+</span>
                    <kbd class="rounded border border-brand-500 bg-background-600 px-1.5 py-0.5 font-medium text-content">
                      V
                    </kbd>
                  </div>
                  <span>{{ $t('upload.dropZone.shortcut.paste') }}</span>
                </div>
              </div>
            </div>

            <div v-if="shouldShowContentDetectionWarning" class="mt-3">
              <div class="flex items-center justify-center rounded border border-default px-3 py-2">
                <i :class="simpleWarningIcon" class="mr-2 text-xs opacity-70" />
                <span class="text-content-secondary text-xs">{{ simpleWarningMessage }}</span>
              </div>
            </div>
          </div>

          <div class="flex flex-col">
            <div class="settings-panel settings-no-scrollbar h-full rounded-md">
              <div class="border-b border-default p-3">
                <h2 class="flex items-center text-base font-semibold text-content">
                  <i class="fas fa-cogs mr-1.5 text-sm text-brand-500" /> {{ $t('upload.settings.title') }}
                </h2>
              </div>

              <div class="p-3">
                <UploadSettings
                  v-model:folder-id="folderId"
                  v-model:access-level="accessLevel"
                  v-model:optimize="optimize"
                  v-model:auto-remove="autoRemove"
                  v-model:storage-duration="storageDuration"
                  v-model:watermark-enabled="watermarkEnabled"
                  v-model:watermark-config="watermarkConfig"
                  :folder-path="presetFolderPath"
                  @change="onSettingsChange"
                />

                <div class="mt-4 border-t border-default pt-3">
                  <button
                    ref="advancedSettingsTrigger"
                    class="flex w-full items-center text-base font-semibold text-content transition-colors hover:text-content-heading"
                    @click="toggleAdvancedSettings"
                  >
                    <i class="fas fa-sliders-h mr-1.5 text-sm text-brand-500" />
                    <span>{{ $t('upload.settings.advanced') }}</span>
                    <i class="fas fa-external-link-alt ml-auto text-xs" />
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>

        <div
          class="quick-actions mt-4 rounded-b-lg border-t border-default bg-gradient-to-r from-background-700 to-background-700 backdrop-blur-sm"
        >
          <div class="flex flex-wrap items-center justify-center gap-3 px-4 py-4">
            <div class="flex flex-wrap items-center gap-2">
              <button
                class="cyber-button cyber-button-gradient flex items-center gap-2 rounded-lg px-4 py-2 font-medium transition-all hover:scale-105 hover:shadow-lg hover:shadow-error-200"
                @click.stop="hasPendingFiles ? startUploadWrapper() : showNoFilesMessage()"
              >
                <i class="fas fa-rocket text-sm" />
                <span>{{ $t('upload.actions.startUpload') }}</span>
              </button>
              <button
                class="cyber-button cyber-button-outline flex items-center gap-2 rounded-lg px-4 py-2 font-medium transition-all hover:scale-105"
                @click.stop="triggerFileInput"
              >
                <i class="fas fa-plus text-sm" />
                <span>{{ $t('upload.actions.selectFiles') }}</span>
              </button>
              <button
                v-if="isFolderUploadSupported"
                class="cyber-button cyber-button-outline flex items-center gap-2 rounded-lg px-4 py-2 font-medium transition-all hover:scale-105"
                @click.stop="triggerFolderInput"
              >
                <i class="fas fa-folder-open text-sm" />
                <span>{{ $t('upload.actions.selectFolder') }}</span>
              </button>
            </div>

            <div class="flex items-center gap-2">
              <button
                class="cyber-button cyber-button-outline flex items-center gap-2 rounded-lg px-3 py-2 text-sm font-medium transition-all hover:scale-105"
                @click="clearQueue"
              >
                <i class="fas fa-trash" />
                <span>{{ $t('upload.actions.clearQueue') }}</span>
              </button>
              <button
                class="cyber-button cyber-button-outline flex items-center gap-2 rounded-lg px-3 py-2 text-sm font-medium transition-all hover:scale-105"
                :title="$t('upload.actions.clearCache')"
                @click="clearAllSessions"
              >
                <i class="fas fa-broom" />
                <span>{{ $t('upload.actions.clearCache') }}</span>
              </button>
            </div>
          </div>
        </div>
      </div>

      <div class="upload-queue-panel mb-4 flex-1">
        <div class="relative flex h-full flex-col">
          <div
            class="bg-background-800/95 relative z-10 flex h-full flex-col rounded-xl border border-default shadow-lg backdrop-blur-sm"
          >
            <div
              class="upload-header bg-background-800/95 sticky top-0 z-10 flex items-center justify-between gap-4 rounded-t-xl border-b border-default px-5 py-3.5 backdrop-blur"
            >
              <div class="flex flex-wrap items-center gap-3">
                <div class="flex items-center gap-2.5 text-base font-semibold text-content-heading">
                  <i class="fas fa-list text-brand-400" />
                  {{ $t('upload.queue.title') }}
                </div>
                <div class="flex items-center gap-2.5 text-xs">
                  <span v-if="uploadQueue.length > 0" class="queue-pill bg-brand-500/15 text-brand-400">
                    {{ $t('upload.queue.filesCount', { count: uploadQueue.length }) }}
                  </span>
                  <span v-if="uploadingCount > 0" class="queue-pill bg-error-500/15 text-error-400">
                    <i class="fas fa-rocket mr-1" />
                    {{ $t('upload.queue.concurrent', { current: uploadingCount, max: maxConcurrentUploads }) }}
                  </span>
                </div>
              </div>

              <div class="flex items-center gap-2">
                <CyberIconButton
                  v-if="hasSuccessFiles"
                  type="cyber"
                  size="small"
                  :tooltip="$t('upload.actions.copyLinks')"
                  tooltip-placement="top"
                  @click="copyAllUrls"
                >
                  <i class="fas fa-copy" />
                </CyberIconButton>
                <CyberIconButton
                  v-if="hasUploadingFiles"
                  type="danger"
                  size="small"
                  :tooltip="$t('upload.actions.cancelUpload')"
                  tooltip-placement="top"
                  @click="handleCancelUpload"
                >
                  <i class="fas fa-ban" />
                </CyberIconButton>
              </div>
            </div>

            <div v-if="uploadQueue.length > 0" class="custom-scrollbar flex-1 overflow-y-auto px-4 py-3">
              <div class="upload-stats bg-background-700/90 mb-3 rounded-lg border border-default p-3 shadow-sm">
                <div class="grid gap-3 md:grid-cols-5">
                  <div class="queue-stat">
                    <span class="queue-stat__label">{{ $t('upload.queue.stats.pending') }}</span>
                    <span class="queue-stat__value">
                      <i class="fas fa-clock text-warning-400" />
                      {{ pendingCount }}
                    </span>
                  </div>
                  <div class="queue-stat">
                    <span class="queue-stat__label">{{ $t('upload.queue.stats.uploading') }}</span>
                    <span class="queue-stat__value">
                      <i class="fas fa-arrow-circle-up text-brand-400" />
                      {{ uploadingCount }}
                    </span>
                  </div>
                  <div class="queue-stat">
                    <span class="queue-stat__label">{{ $t('upload.queue.stats.success') }}</span>
                    <span class="queue-stat__value">
                      <i class="fas fa-check-circle text-success-400" />
                      {{ successCount }}
                    </span>
                  </div>
                  <div class="queue-stat">
                    <span class="queue-stat__label">{{ $t('upload.queue.stats.failed') }}</span>
                    <span class="queue-stat__value">
                      <i class="fas fa-times-circle text-error-400" />
                      {{ errorCount }}
                    </span>
                  </div>
                  <div class="queue-stat">
                    <span class="queue-stat__label">{{ $t('upload.queue.stats.total') }}</span>
                    <span class="queue-stat__value">
                      <i class="fas fa-layer-group text-info-400" />
                      {{ uploadQueue.length }}
                    </span>
                  </div>
                </div>

                <div class="queue-meta">
                  <span class="queue-meta__label">{{ $t('upload.queue.stats.totalSize') }}</span>
                  <span class="queue-meta__value">
                    {{ formatFileSize(totalFileSize) }}
                    <span v-if="globalProgress > 0"> · {{ globalProgress }}%</span>
                  </span>
                </div>
              </div>

              <UploadProgress
                :files="uploadQueue"
                @remove="removeFileWrapper"
                @upload="startUploadWrapper"
                @cancel="handleCancelUpload"
                @retry="retryUploadWrapper"
                @resume="resumeUploadWrapper"
              />
            </div>

            <div v-else class="flex flex-1 flex-col items-center justify-center p-8">
              <div class="max-w-md p-8 text-center">
                <div class="mb-6 flex items-center justify-center">
                  <div class="relative">
                    <div class="bg-brand-500/20 absolute inset-0 animate-pulse rounded-full blur-xl" />
                    <div
                      class="from-brand-500/20 to-brand-600/20 relative flex h-20 w-20 items-center justify-center rounded-full bg-gradient-to-br shadow-inner"
                    >
                      <i class="fas fa-cloud-upload-alt text-4xl text-brand-400" />
                    </div>
                  </div>
                </div>
                <h3 class="mb-3 text-xl font-semibold text-content-heading">{{ $t('upload.queue.empty.title') }}</h3>
                <p class="mb-6 text-sm leading-relaxed text-content-muted">{{ $t('upload.queue.empty.desc') }}</p>
                <div class="flex justify-center">
                  <button class="queue-empty-primary" @click="triggerFileInput">
                    <i class="fas fa-images mr-2" /> {{ $t('upload.queue.empty.action') }}
                  </button>
                </div>
              </div>
            </div>

            <div class="upload-tips bg-background-800/90 rounded-b-2xl border-t border-subtle backdrop-blur">
              <div class="p-4">
                <div class="mb-4 flex items-center">
                  <div class="to-accent-500 mr-3 h-6 w-2 rounded-full bg-gradient-to-b from-brand-500" />
                  <h4 class="text-base font-semibold text-content">{{ $t('upload.guide.title') }}</h4>
                </div>

                <div class="grid grid-cols-1 gap-3 lg:grid-cols-2 xl:grid-cols-3">
                  <div
                    class="rounded-lg border border-default p-3 transition-all duration-200 hover:border-strong hover:bg-background-700"
                  >
                    <div class="mb-2 flex items-center">
                      <div class="mr-2 flex h-6 w-6 items-center justify-center rounded-md bg-brand-500">
                        <i class="fas fa-upload text-xs text-content-heading" />
                      </div>
                      <span class="text-sm font-medium text-content">{{ $t('upload.guide.multipleWays.title') }}</span>
                    </div>
                    <p class="text-content-secondary text-xs leading-relaxed">
                      {{ $t('upload.guide.multipleWays.desc', { key: pasteShortcut.key }) }}
                    </p>
                  </div>

                  <div
                    class="rounded-lg border border-default p-3 transition-all duration-200 hover:border-strong hover:bg-background-700"
                  >
                    <div class="mb-2 flex items-center">
                      <div class="mr-2 flex h-6 w-6 items-center justify-center rounded-md bg-error-500">
                        <i class="fas fa-puzzle-piece text-xs text-content-heading" />
                      </div>
                      <span class="text-sm font-medium text-content">{{ $t('upload.guide.smartChunking.title') }}</span>
                    </div>
                    <p class="text-content-secondary text-xs leading-relaxed">{{ $t('upload.guide.smartChunking.desc') }}</p>
                  </div>

                  <div
                    class="rounded-lg border border-default p-3 transition-all duration-200 hover:border-strong hover:bg-background-700"
                  >
                    <div class="mb-2 flex items-center">
                      <div class="mr-2 flex h-6 w-6 items-center justify-center rounded-md bg-brand-500">
                        <i class="fas fa-magic text-xs text-content-heading" />
                      </div>
                      <span class="text-sm font-medium text-content">{{ $t('upload.guide.oneCopy.title') }}</span>
                    </div>
                    <p class="text-content-secondary text-xs leading-relaxed">{{ $t('upload.guide.oneCopy.desc') }}</p>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <teleport to="body">
      <div
        v-if="showAdvancedSettings"
        class="advanced-settings-panel fixed z-[9999] w-80 rounded-lg border border-default bg-background-900 p-4 shadow-2xl backdrop-blur-md"
        :style="panelPosition"
        @click.stop
      >
        <div class="mb-4 flex items-center justify-between">
          <h3 class="flex items-center text-sm font-semibold text-content">
            <i class="fas fa-sliders-h mr-2 text-brand-500" />
            {{ $t('upload.advanced.title') }}
          </h3>
          <button
            class="text-content-secondary transition-colors hover:text-content-heading"
            @click="showAdvancedSettings = false"
          >
            <i class="fas fa-times text-sm" />
          </button>
        </div>

        <div v-if="authStore.isLoggedIn" class="space-y-3">
          <div>
            <label class="mb-1.5 block text-sm font-medium text-content">{{ $t('upload.advanced.storageDuration') }}</label>
            <cyberDropdown
              v-model="storageDuration"
              :options="storageDurationOptionsForLoggedUser"
              :placeholder="$t('upload.form.storageDuration.placeholder')"
              class="compact-dropdown text-sm"
            >
              <template #option-icon>
                <i class="fas fa-calendar-alt text-sm text-content" />
              </template>

              <template #selected-icon>
                <i class="fas fa-calendar-alt text-sm text-content" />
              </template>
            </cyberDropdown>
            <p class="text-content-secondary mt-1 text-xs">{{ $t('upload.advanced.storageDurationDesc') }}</p>
          </div>

          <div>
            <div class="flex items-center justify-between">
              <cyberCheckbox v-model="watermarkEnabled"> {{ $t('upload.advanced.watermark') }} </cyberCheckbox>
              <button
                v-if="watermarkEnabled"
                class="config-btn"
                :title="$t('upload.watermark.configure')"
                @click="showWatermarkConfig = true"
              >
                <i class="fas fa-cog" />
              </button>
            </div>
            <p class="text-content-secondary ml-6 mt-1 text-xs">{{ $t('upload.advanced.watermarkDesc') }}</p>

            <div v-if="watermarkEnabled && watermarkConfig.type" class="ml-6 mt-2">
              <div class="text-content-secondary space-y-1 text-xs">
                <div class="flex items-center">
                  <i class="fas fa-tag mr-1" />
                  <span
                    >{{ $t('upload.advanced.watermarkInfo.type') }}
                    {{
                      watermarkConfig.type === 'text'
                        ? $t('upload.advanced.watermarkType.text')
                        : $t('upload.advanced.watermarkType.image')
                    }}</span
                  >
                </div>
                <div v-if="watermarkConfig.text" class="flex items-center">
                  <i class="fas fa-font mr-1" />
                  <span>{{ $t('upload.advanced.watermarkInfo.content') }} {{ watermarkConfig.text }}</span>
                </div>
                <div class="flex items-center">
                  <i class="fas fa-map-marker-alt mr-1" />
                  <span>{{ $t('upload.advanced.watermarkInfo.position') }} {{ getPositionText(watermarkConfig.position) }}</span>
                </div>
              </div>
            </div>
          </div>
        </div>

        <div v-else class="py-6 text-center">
          <i class="fas fa-star mb-3 text-2xl text-brand-500" />
          <p class="mb-2 text-sm text-content">{{ $t('upload.advanced.loginPrompt.title') }}</p>
          <p class="text-content-secondary mb-4 text-xs">{{ $t('upload.advanced.loginPrompt.desc') }}</p>
          <button
            class="rounded bg-brand-500 px-4 py-2 text-sm text-content-heading transition-colors hover:bg-brand-600"
            @click="$router.push('/auth')"
          >
            {{ $t('upload.advanced.loginPrompt.action') }}
          </button>
        </div>
      </div>

      <div v-if="showAdvancedSettings" class="fixed inset-0 z-[9998]" @click="showAdvancedSettings = false" />
    </teleport>

    <cyberWatermarkConfig
      v-model:visible="showWatermarkConfig"
      :config="watermarkConfig"
      @confirm="handleWatermarkConfigConfirm"
    />

    <FolderConfirmDialog
      v-model:visible="showFolderDialog"
      :stats="folderStats"
      @confirm="handleConfirmFolderUpload"
      @cancel="handleCancelFolderUpload"
    />
  </div>
</template>

<style scoped>
  .upload-container {
    min-height: calc(100vh - 60px);
  }

  .upload-top-panel {
    background: rgba(var(--color-background-800-rgb), 0.95);
    backdrop-filter: blur(10px);
    border: 1.5px solid rgba(var(--color-brand-500-rgb), 0.2);
    box-shadow:
      0 4px 12px var(--color-overlay-medium),
      0 2px 6px var(--color-overlay-light);
  }

  .custom-scrollbar::-webkit-scrollbar {
    width: 4px;
    height: 4px;
  }

  .custom-scrollbar::-webkit-scrollbar-track {
    background: rgba(var(--color-background-900-rgb), 0.1);
    border-radius: var(--radius-sm);
  }

  .custom-scrollbar::-webkit-scrollbar-thumb {
    background: linear-gradient(to bottom, rgba(var(--color-brand-500-rgb), 0.5), rgba(var(--color-brand-600-rgb), 0.5));
    border-radius: var(--radius-sm);
  }

  .custom-scrollbar::-webkit-scrollbar-thumb:hover {
    background: linear-gradient(to bottom, rgba(var(--color-brand-500-rgb), 0.7), rgba(var(--color-brand-600-rgb), 0.7));
  }

  .settings-panel {
    background: rgba(var(--color-background-700-rgb), 0.6);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.15);
  }

  .settings-no-scrollbar {
    max-height: 100%;
    overflow: hidden;
  }

  .settings-no-scrollbar::-webkit-scrollbar {
    width: 0;
    height: 0;
    display: none;
  }

  .text-gradient-blue-pink {
    background: linear-gradient(to right, var(--color-brand-500), var(--color-fuchsia-600));
    -webkit-background-clip: text;
    background-clip: text;
    color: transparent;
  }

  .cyber-button {
    font-size: 0.75rem;
    font-weight: 500;
    padding: 0.35rem 0.75rem;
    border-radius: var(--radius-sm);
    cursor: pointer;
    transition: all 0.2s ease;
    display: inline-flex;
    align-items: center;
    justify-content: center;
  }

  .cyber-button:hover {
    transform: translateY(-1px);
    box-shadow: 0 2px 8px var(--color-overlay-medium);
  }

  .cyber-button:active {
    transform: translateY(0);
  }

  .cyber-button-xs {
    font-size: 0.65rem;
    padding: 0.2rem 0.5rem;
  }

  .cyber-button-small {
    font-size: 0.7rem;
    padding: 0.25rem 0.65rem;
  }

  .cyber-button-blue {
    background-color: rgba(var(--color-brand-500-rgb), 0.15);
    color: var(--color-brand-500);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.3);
  }

  .cyber-button-blue:hover {
    background-color: rgba(var(--color-brand-500-rgb), 0.25);
    border-color: rgba(var(--color-brand-500-rgb), 0.5);
  }

  .cyber-button-danger {
    background-color: rgba(var(--color-error-rgb), 0.15);
    color: var(--color-error-500);
    border: 1px solid rgba(var(--color-error-rgb), 0.3);
  }

  .cyber-button-danger:hover {
    background-color: rgba(var(--color-error-rgb), 0.25);
    border-color: rgba(var(--color-error-rgb), 0.5);
  }

  .cyber-button-outline {
    background-color: transparent;
    color: var(--color-brand-500);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.3);
  }

  .cyber-button-outline:hover {
    background-color: rgba(var(--color-brand-500-rgb), 0.05);
    border-color: rgba(var(--color-brand-500-rgb), 0.5);
  }

  .cyber-button-gradient {
    background: linear-gradient(to right, rgba(var(--color-brand-500-rgb), 0.7), rgba(var(--color-brand-600-rgb), 0.7));
    color: var(--color-text-on-brand);
    border: none;
  }

  .cyber-button-gradient:hover {
    background: linear-gradient(to right, rgba(var(--color-brand-500-rgb), 0.8), rgba(var(--color-brand-600-rgb), 0.8));
  }

  .cyber-button-gradient-outline {
    background: transparent;
    color: var(--color-brand-500);
    border: 1px solid;
    border-image: linear-gradient(to right, rgba(var(--color-brand-500-rgb), 0.5), rgba(var(--color-error-rgb), 0.5)) 1;
    position: relative;
    overflow: hidden;
    z-index: 1;
  }

  .cyber-dark-lighter {
    background-color: rgba(var(--color-background-700-rgb), 0.5);
  }

  .upload-drop-zone:hover .upload-icon,
  .upload-drop-zone.border-brand-500 .upload-icon {
    background: linear-gradient(to right, rgba(var(--color-brand-500-rgb), 0.1), rgba(var(--color-error-rgb), 0.1));
    border-radius: var(--radius-sm);
    transform: scale(1.05);
  }

  .upload-drop-zone:hover .upload-icon-container {
    transform: translateY(-2px);
  }

  .upload-drop-zone:hover {
    box-shadow: 0 8px 32px rgba(var(--color-brand-500-rgb), 0.1);
  }

  .upload-queue-panel {
    min-height: 220px;
  }

  .upload-queue-panel > .relative > div {
    background: rgba(var(--color-background-800-rgb), 0.95) !important;
    border: 1.5px solid rgba(var(--color-brand-500-rgb), 0.2);
    box-shadow:
      0 4px 12px var(--color-overlay-medium),
      0 2px 6px var(--color-overlay-light);
  }

  .text-2xs {
    font-size: 0.65rem;
    line-height: 1rem;
  }

  .queue-pill {
    @apply inline-flex items-center rounded-full border border-default px-2.5 py-1 font-medium shadow-sm;
  }

  .queue-stat {
    border-radius: var(--radius-sm);
    padding: 0.75rem 0.875rem;
    background-color: rgba(var(--color-background-800-rgb), 0.8);
    border: 1px solid var(--color-border-default);
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
    transition: all 0.2s ease;
  }

  .queue-stat:hover {
    background-color: rgba(var(--color-background-800-rgb), 0.95);
    border-color: rgba(var(--color-brand-500-rgb), 0.3);
  }

  .queue-stat__label {
    font-size: 0.75rem;
    font-weight: 500;
    color: var(--color-content-muted);
    letter-spacing: 0.02em;
  }

  .queue-stat__value {
    display: inline-flex;
    align-items: center;
    gap: 0.5rem;
    font-size: 0.875rem;
    font-weight: 600;
    color: var(--color-content);
  }

  .queue-stat__value i {
    font-size: 0.875rem;
  }

  .queue-meta {
    margin-top: 1rem;
    border-radius: var(--radius-sm);
    background-color: rgba(var(--color-background-800-rgb), 0.8);
    padding: 0.75rem 1rem;
    display: flex;
    justify-content: space-between;
    align-items: center;
    font-size: 0.8125rem;
    color: var(--color-content-muted);
    border: 1px solid var(--color-border-subtle);
  }

  .queue-meta__label {
    font-weight: 500;
  }

  .queue-meta__value {
    font-weight: 600;
    color: var(--color-content);
  }

  .queue-primary,
  .queue-secondary {
    @apply inline-flex items-center justify-center gap-1.5 rounded-lg px-3 py-1.5 text-xs font-medium transition-colors duration-200;
  }

  .queue-primary {
    background-color: rgba(var(--color-brand-500-rgb), 0.16);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.36);
    color: var(--color-brand-400);
  }

  .queue-primary:hover {
    background-color: rgba(var(--color-brand-500-rgb), 0.22);
  }

  .queue-secondary {
    border: 1px solid var(--color-border-subtle);
    color: var(--color-content-muted);
    background-color: rgba(var(--color-background-700-rgb), 0.4);
  }

  .queue-secondary:hover {
    color: var(--color-content-heading);
    border-color: var(--color-border-default);
  }

  .queue-empty-primary,
  .queue-empty-secondary {
    @apply inline-flex items-center justify-center rounded-lg px-5 py-2.5 text-sm font-semibold transition-all duration-200;
  }

  .queue-empty-primary {
    background: linear-gradient(135deg, rgba(var(--color-brand-500-rgb), 0.18), rgba(var(--color-brand-600-rgb), 0.22));
    border: 1.5px solid rgba(var(--color-brand-500-rgb), 0.4);
    color: var(--color-brand-400);
    box-shadow: 0 2px 8px rgba(var(--color-brand-500-rgb), 0.15);
  }

  .queue-empty-primary:hover {
    background: linear-gradient(135deg, rgba(var(--color-brand-500-rgb), 0.25), rgba(var(--color-brand-600-rgb), 0.3));
    border-color: rgba(var(--color-brand-500-rgb), 0.5);
    transform: translateY(-1px);
    box-shadow: 0 4px 12px rgba(var(--color-brand-500-rgb), 0.25);
  }

  .queue-empty-secondary {
    border: 1.5px solid var(--color-border-default);
    color: var(--color-content);
    background-color: rgba(var(--color-background-800-rgb), 0.6);
  }

  .queue-empty-secondary:hover {
    color: var(--color-content-heading);
    border-color: rgba(var(--color-brand-500-rgb), 0.3);
    background-color: rgba(var(--color-background-800-rgb), 0.8);
    transform: translateY(-1px);
  }

  :deep(.compact-dropdown) {
    --dropdown-height: 28px;
    --item-padding: 4px 8px;
    font-size: 0.75rem;
  }

  :deep(.compact-dropdown .dropdown-trigger) {
    height: var(--dropdown-height);
    min-height: var(--dropdown-height);
  }

  :deep(.compact-dropdown .dropdown-item) {
    padding: var(--item-padding);
  }

  :deep(.compact-dropdown .selected-item) {
    padding: var(--item-padding);
  }

  .advanced-settings-panel {
    animation: fadeInScale 0.2s ease-out;
    box-shadow:
      0 10px 25px var(--color-overlay-heavy),
      0 0 20px rgba(var(--color-brand-500-rgb), 0.1);
  }

  @keyframes fadeInScale {
    from {
      opacity: 0;
      transform: scale(0.95) translateY(-5px);
    }
    to {
      opacity: 1;
      transform: scale(1) translateY(0);
    }
  }

  .config-btn {
    @apply flex h-6 w-6 cursor-pointer items-center justify-center rounded border transition-all duration-200;
    border-color: var(--color-border-default);
    background: var(--color-background-800);
    color: var(--color-brand-500);
  }

  .config-btn:hover {
    border-color: var(--color-brand-500);
    background: var(--color-hover-bg);
    color: var(--color-brand-500);
    transform: scale(1.05);
  }
</style>
