/**
 * 智能Loading管理系统
 * 约定优于配置，自动检测和管理组件中的loading状态
 */
import { getCurrentInstance, type Ref } from 'vue'

/* 支持的loading变量命名模式 */
const LOADING_PATTERNS = [
  /^loading$/,
  /^isLoading$/,
  /^.*Loading$/,
  /^submitting$/,
  /^.*Submitting$/,
  /^fetching$/,
  /^.*Fetching$/,
  /^pending$/,
  /^.*Pending$/,
]

/* 全局loading状态管理 */
class SmartLoadingManager {
  private loadingStates = new Map<any, { count: number; startTime: number }>()

  detectLoadingRefs(): Ref<boolean>[] {
    const instance = getCurrentInstance()
    if (!instance) {
      return []
    }

    const foundRefs: Ref<boolean>[] = []

    if (instance.setupState) {
      Object.entries(instance.setupState).forEach(([key, value]) => {
        if (this.isLoadingPattern(key) && this.isRef(value)) {
          foundRefs.push(value as Ref<boolean>)
        }
      })
    }

    return foundRefs
  }

  private isLoadingPattern(name: string): boolean {
    return LOADING_PATTERNS.some((pattern) => pattern.test(name))
  }

  private isRef(value: any): boolean {
    return value && typeof value === 'object' && 'value' in value
  }

  startLoading(refs: Ref<boolean>[], _minTime: number = 300) {
    refs.forEach((ref) => {
      if (!this.loadingStates.has(ref)) {
        this.loadingStates.set(ref, { count: 0, startTime: Date.now() })
        ref.value = true
      }

      const state = this.loadingStates.get(ref)
      if (state) {
        state.count++
      }
    })
  }

  stopLoading(refs: Ref<boolean>[], minTime: number = 300) {
    refs.forEach((ref) => {
      const state = this.loadingStates.get(ref)
      if (!state) {
        return
      }

      state.count--

      if (state.count <= 0) {
        const elapsed = Date.now() - state.startTime

        if (elapsed < minTime) {
          setTimeout(() => {
            ref.value = false
            this.loadingStates.delete(ref)
          }, minTime - elapsed)
        } else {
          ref.value = false
          this.loadingStates.delete(ref)
        }
      }
    })
  }

  resetAll() {
    this.loadingStates.forEach((state, ref) => {
      ref.value = false
    })
    this.loadingStates.clear()
  }
}

export const smartLoadingManager = new SmartLoadingManager()
