import { computed, ref } from 'vue'
import { useToast } from '@/components/Toast/useToast'
import { getChannelList } from '@/api/storage'
import type { StorageChannel } from '@/api/types/index'

export function useChannels() {
  const toast = useToast()

  const channels = ref<StorageChannel[]>([])
  const loading = ref(false)

  const sortedChannels = computed(() =>
    [...(channels.value || [])].sort((a, b) => {
      if (a.type === 'local') {
        return -1
      }
      if (b.type === 'local') {
        return 1
      }
      return new Date(b.created_at).getTime() - new Date(a.created_at).getTime()
    })
  )

  const fetchChannels = async () => {
    try {
      loading.value = true
      const result = await getChannelList()
      if (result.success) {
        channels.value = result.data || []
      }
    } catch (error: unknown) {
      toast.error(error.message)
      channels.value = []
    } finally {
      loading.value = false
    }
  }

  return {
    channels,
    loading,
    sortedChannels,

    fetchChannels,
  }
}
