/**
 * 多级缓存Composable
 * 提供内存缓存和LocalStorage缓存
 */
import { readonly, ref } from 'vue'
import { StorageUtil } from '@/utils/storage/storage'

interface CacheItem {
  data: any
  timestamp: number
  ttl: number
}

class MultiLevelCache {
  private memoryCache = new Map<string, CacheItem>()
  private storagePrefix = 'cyber_cache_'

  async get<T = any>(key: string): Promise<T | null> {
    const memoryItem = this.memoryCache.get(key)
    if (memoryItem && Date.now() - memoryItem.timestamp < memoryItem.ttl) {
      return memoryItem.data
    }

    try {
      const item = StorageUtil.get<CacheItem>(this.storagePrefix + key)
      if (item) {
        if (Date.now() - item.timestamp < item.ttl) {
          this.memoryCache.set(key, item)
          return item.data
        }
        StorageUtil.remove(this.storagePrefix + key)
      }
    } catch (_error) {}

    return null
  }

  set(key: string, data: any, ttl = 5 * 60 * 1000) {
    const item: CacheItem = {
      data,
      timestamp: Date.now(),
      ttl,
    }

    this.memoryCache.set(key, item)

    try {
      StorageUtil.set(this.storagePrefix + key, item)
    } catch (_error) {}
  }

  clear(key?: string) {
    if (key) {
      this.memoryCache.delete(key)
      StorageUtil.remove(this.storagePrefix + key)
    } else {
      this.memoryCache.clear()
      StorageUtil.keys().forEach((k) => {
        if (k.startsWith(this.storagePrefix)) {
          StorageUtil.remove(k)
        }
      })
    }
  }

  cleanup() {
    const now = Date.now()

    for (const [key, item] of this.memoryCache.entries()) {
      if (now - item.timestamp >= item.ttl) {
        this.memoryCache.delete(key)
      }
    }

    StorageUtil.keys().forEach((key) => {
      if (key.startsWith(this.storagePrefix)) {
        try {
          const item = StorageUtil.get<CacheItem>(key)
          if (item && now - item.timestamp >= item.ttl) {
            StorageUtil.remove(key)
          }
        } catch {
          StorageUtil.remove(key)
        }
      }
    })
  }
}

const globalCache = new MultiLevelCache()

if (typeof window !== 'undefined') {
  setInterval(
    () => {
      globalCache.cleanup()
    },
    10 * 60 * 1000
  ) // 每10分钟清理一次
}

export function useCache() {
  const isLoading = ref(false)

  const getCached = async <T = any>(key: string): Promise<T | null> => globalCache.get<T>(key)

  const setCached = (key: string, data: any, ttl?: number) => {
    globalCache.set(key, data, ttl)
  }

  const clearCache = (key?: string) => {
    globalCache.clear(key)
  }

  const withCache = async <T = any>(
    key: string,
    asyncFn: () => Promise<T>,
    options: { ttl?: number; forceRefresh?: boolean } = {}
  ): Promise<T> => {
    const { ttl = 5 * 60 * 1000, forceRefresh = false } = options

    if (!forceRefresh) {
      const cached = await getCached<T>(key)
      if (cached !== null) {
        return cached
      }
    }

    isLoading.value = true
    try {
      const result = await asyncFn()
      setCached(key, result, ttl)
      return result
    } finally {
      isLoading.value = false
    }
  }

  return {
    isLoading: readonly(isLoading),
    getCached,
    setCached,
    clearCache,
    withCache,
  }
}
