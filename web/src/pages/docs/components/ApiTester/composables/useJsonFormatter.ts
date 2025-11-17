import { getCurrentLocale } from '@/utils/locale'
import { computed, type Ref } from 'vue'
import type { ApiResponse } from '../types'

export function useJsonFormatter(response: Ref<ApiResponse | null>) {
  const formattedResponse = computed(() => {
    if (!response.value?.data) return ''
    return JSON.stringify(response.value.data, null, 2)
  })

  const formatDate = (dateString: string) => {
    if (!dateString) return ''
    const date = new Date(dateString)
    return date.toLocaleDateString(getCurrentLocale(), {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit',
    })
  }

  return {
    formattedResponse,
    formatDate,
  }
}
