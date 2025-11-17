import { ref } from 'vue'
import { useToast } from '@/components/Toast/useToast'
import { useTexts } from '@/composables/useTexts'
import { exportAllChannels, exportChannel, importChannel } from '@/api/storage'
import type { StorageChannel } from '@/api/types/index'

/**
 * 文件操作 Composable
 * 负责渠道配置的导入、导出功能
 */
export function useFileOperations() {
  const toast = useToast()
  const { $t } = useTexts()

  const fileInputRef = ref<HTMLInputElement | null>(null)

  const importChannelConfig = () => {
    fileInputRef.value?.click()
  }

  const handleExportChannel = async (channel: StorageChannel) => {
    if (channel.type === 'local') {
      toast.warning($t('admin.channels.file.localNotSupport'))
      return
    }

    try {
      toast.info($t('admin.channels.file.preparingExport'))
      const result = await exportChannel(channel.id)
      if (result.success) {
        const { data } = result
        const json = JSON.stringify(data, null, 2)
        const blob = new Blob([json], { type: 'application/json' })
        const url = URL.createObjectURL(blob)

        const a = document.createElement('a')
        a.href = url
        const safeName = channel.name.replace(/[^a-zA-Z0-9\u4e00-\u9fa5_-]/g, '_')
        a.download = `${safeName}-export.json`
        document.body.appendChild(a)
        a.click()

        document.body.removeChild(a)
        URL.revokeObjectURL(url)

        toast.success($t('admin.channels.file.exportSuccess'))
      }
    } catch (error: unknown) {
      toast.error(error.message)
    }
  }

  const handleExportAllChannels = async () => {
    try {
      toast.info($t('admin.channels.file.preparingExportAll'))
      const result = await exportAllChannels()
      if (result.success) {
        const { data } = result
        const json = JSON.stringify(data, null, 2)
        const blob = new Blob([json], { type: 'application/json' })
        const url = URL.createObjectURL(blob)

        const a = document.createElement('a')
        a.href = url
        a.download = `all-channels-config-export.json`
        document.body.appendChild(a)
        a.click()

        document.body.removeChild(a)
        URL.revokeObjectURL(url)

        toast.success($t('admin.channels.file.exportAllSuccess', { count: data.total_count }))
      }
    } catch (error: unknown) {
      const errorMessage =
        (error as any)?.response?.data?.message || (error as any)?.message || $t('admin.channels.file.exportFailedDefault')
      toast.error(errorMessage)
    }
  }

  const handleFileImportInternal = async (file: File, onSuccess?: () => void) => {
    try {
      toast.info($t('admin.channels.file.importing'))
      const result = await importChannel(file)
      if (result.success) {
        toast.success($t('admin.channels.file.importSuccess'))
        onSuccess?.()
      }
    } catch (error: unknown) {
      const errorMessage =
        (error as any)?.response?.data?.message || (error as any)?.message || $t('admin.channels.file.importFailedDefault')
      toast.error(errorMessage)
    }
  }

  const handleFileImport = async (event: Event, onSuccess?: () => void) => {
    const target = event.target as HTMLInputElement
    const { files } = target

    if (!files || files.length === 0) {
      return
    }

    const file = files[0]
    await handleFileImportInternal(file, onSuccess)

    if (fileInputRef.value) {
      fileInputRef.value.value = ''
    }
  }

  return {
    fileInputRef,

    importChannelConfig,
    handleExportChannel,
    handleExportAllChannels,
    handleFileImportInternal,
    handleFileImport,
  }
}
