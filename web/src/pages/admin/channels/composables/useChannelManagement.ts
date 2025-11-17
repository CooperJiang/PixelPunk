import { useChannels } from './useChannels'
import { useChannelForm } from './useChannelForm'
import { useChannelOperations } from './useChannelOperations'
import { useFileOperations } from './useFileOperations'
import { useDragAndDrop } from './useDragAndDrop'
import type { StorageChannel } from '@/api/types/index'

export function useChannelManagement() {
  const channelData = useChannels()

  const channelForm = useChannelForm()

  const channelOperations = useChannelOperations()

  const fileOperations = useFileOperations()

  const dragAndDrop = useDragAndDrop(
    '.channels-content',
    (file: File) => {
      fileOperations.handleFileImportInternal(file, channelData.fetchChannels)
    },
    {
      isDisabled: () => channelForm.showAddChannelDialog.value === true,
    }
  )

  const handleChannelOperationWithRefresh = {
    toggleStatus: (channel: StorageChannel) => channelOperations.toggleChannelStatus(channel, channelData.fetchChannels),
    setDefault: (channel: StorageChannel) => channelOperations.handleSetDefault(channel, channelData.fetchChannels),
    deleteChannel: (channel: StorageChannel) => channelOperations.handleDeleteChannel(channel, channelData.fetchChannels),
    testConnection: channelOperations.testChannelConnection,
  }

  const handleFormOperationWithRefresh = {
    saveChannel: () => channelForm.saveChannel(channelData.fetchChannels),
    editChannel: channelForm.editChannel,
    resetForm: channelForm.resetForm,
  }

  const handleFileOperationWithRefresh = {
    importConfig: fileOperations.importChannelConfig,
    exportChannel: fileOperations.handleExportChannel,
    exportAll: fileOperations.handleExportAllChannels,
    handleFileImport: (event: Event) => fileOperations.handleFileImport(event, channelData.fetchChannels),
  }

  return {
    ...channelData,

    ...channelForm,

    fileInputRef: fileOperations.fileInputRef,

    isDragOver: dragAndDrop.isDragOver,

    channelOperations: handleChannelOperationWithRefresh,
    formOperations: handleFormOperationWithRefresh,
    fileOperations: handleFileOperationWithRefresh,
  }
}
