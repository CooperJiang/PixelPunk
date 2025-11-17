import { useToast } from '@/components/Toast/useToast'
import { useTexts } from '@/composables/useTexts'
import { deleteChannel, disableChannel, enableChannel, setDefaultChannel, testChannel } from '@/api/storage'
import type { StorageChannel } from '@/api/types/index'

/**
 * 渠道操作 Composable
 * 负责渠道的启用、禁用、测试、设为默认、删除等操作
 */
export function useChannelOperations() {
  const toast = useToast()
  const { $t } = useTexts()

  const toggleChannelStatus = async (channel: StorageChannel, onSuccess?: () => void) => {
    try {
      if (channel.status === 1) {
        const result = await disableChannel(channel.id)
        if (result.success) {
          toast.success($t('admin.channels.operations.disabled'))
        }
      } else {
        const result = await enableChannel(channel.id)
        if (result.success) {
          toast.success($t('admin.channels.operations.enabled'))
        }
      }
      onSuccess?.()
    } catch (error: unknown) {
      // 全局 HTTP 拦截器已经显示了错误提示，这里不再重复显示
    }
  }

  const testChannelConnection = async (channel: StorageChannel) => {
    if (!channel) {
      toast.warning($t('admin.channels.operations.testWarning'))
      return
    }

    try {
      toast.info($t('admin.channels.operations.testing'))
      const result = await testChannel(channel.id)
      if (result.success) {
        toast.success($t('admin.channels.operations.testSuccess'))
      }
    } catch (error: unknown) {
      const errorMessage =
        (error as any)?.message || (error as any)?.response?.data?.message || $t('admin.channels.operations.testConnectionFailed')
      toast.error($t('admin.channels.operations.testFailed', { message: errorMessage }))
    }
  }

  const handleSetDefault = async (channel: StorageChannel, onSuccess?: () => void) => {
    if (channel.is_default) {
      toast.info($t('admin.channels.operations.alreadyDefault'))
      return
    }

    try {
      const result = await setDefaultChannel(channel.id)
      if (result.success) {
        toast.success($t('admin.channels.operations.setDefaultSuccess'))
        onSuccess?.()
      }
    } catch (error: unknown) {
      // 全局 HTTP 拦截器已经显示了错误提示，这里不再重复显示
    }
  }

  const handleDeleteChannel = async (channel: StorageChannel, onSuccess?: () => void) => {
    try {
      const result = await deleteChannel(channel.id)
      if (result.success) {
        toast.success($t('admin.channels.operations.deleteSuccess'))
        onSuccess?.()
      }
    } catch (error: unknown) {
      // 全局 HTTP 拦截器已经显示了错误提示，这里不再重复显示
    }
  }

  return {
    toggleChannelStatus,
    testChannelConnection,
    handleSetDefault,
    handleDeleteChannel,
  }
}
