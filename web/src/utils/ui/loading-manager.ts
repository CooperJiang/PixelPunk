/**
 * Loading状态管理器
 * 解决多API调用共享loading状态的问题
 */
import type { Ref as _Ref } from 'vue'

interface LoadingRef {
  value: boolean
  _refCount?: number
  _minTime?: number
  _startTime?: number
  _timeoutId?: any
}

/* 全局Loading管理器 */
class LoadingManager {
  private loadingGroups = new Map<string, Set<LoadingRef>>()

  incrementLoading(
    loadingRef: LoadingRef,
    options: {
      minLoadingTime?: number
      groupId?: string
    } = {}
  ) {
    const { minLoadingTime = 300, groupId } = options

    if (!loadingRef._refCount) {
      loadingRef._refCount = 0
    }

    loadingRef._refCount++

    if (loadingRef._refCount === 1) {
      loadingRef.value = true
      loadingRef._startTime = Date.now()
      loadingRef._minTime = minLoadingTime
    }

    if (groupId) {
      if (!this.loadingGroups.has(groupId)) {
        this.loadingGroups.set(groupId, new Set())
      }
      this.loadingGroups.get(groupId)?.add(loadingRef)
    }
  }

  decrementLoading(loadingRef: LoadingRef, groupId?: string) {
    if (!loadingRef._refCount || loadingRef._refCount <= 0) {
      return
    }

    loadingRef._refCount--

    if (loadingRef._refCount === 0) {
      this.resetLoadingWithMinTime(loadingRef)
    }

    if (groupId) {
      this.loadingGroups.get(groupId)?.delete(loadingRef)
    }
  }

  private resetLoadingWithMinTime(loadingRef: LoadingRef) {
    const elapsed = Date.now() - (loadingRef._startTime || 0)
    const minTime = loadingRef._minTime || 300

    if (loadingRef._timeoutId) {
      clearTimeout(loadingRef._timeoutId)
    }

    if (elapsed < minTime) {
      loadingRef._timeoutId = setTimeout(() => {
        loadingRef.value = false
        loadingRef._refCount = 0
      }, minTime - elapsed)
    } else {
      loadingRef.value = false
      loadingRef._refCount = 0
    }
  }

  forceResetLoading(loadingRef: LoadingRef, groupId?: string) {
    if (loadingRef._timeoutId) {
      clearTimeout(loadingRef._timeoutId)
    }

    loadingRef.value = false
    loadingRef._refCount = 0

    if (groupId) {
      this.loadingGroups.get(groupId)?.delete(loadingRef)
    }
  }

  resetGroup(groupId: string) {
    const group = this.loadingGroups.get(groupId)
    if (group) {
      group.forEach((loadingRef) => {
        this.forceResetLoading(loadingRef)
      })
      this.loadingGroups.delete(groupId)
    }
  }
}

export const loadingManager = new LoadingManager()

export default loadingManager
